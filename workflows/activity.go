package workflows

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"strings"
	
	"backend/main/config"
	"backend/main/models"
	
	"github.com/google/generative-ai-go/genai"
	"go.uber.org/zap"
)

func GetFlippData(ctx context.Context, zipCode string, merchant string) (RequestData, error) {
	httpClient := &http.Client{}
	resp, err := httpClient.Get("https://backflipp.wishabi.com/flipp/items/search?locale=en&postal_code=" + zipCode + "&q=" + merchant)
	if err != nil {
		return RequestData{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RequestData{}, err
	}
	var requestData RequestData
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		return RequestData{}, err
	}
	foodItemsRequestData := RequestData{}
	for _, value := range requestData.Items {
		if value.L2 == "Food Items" {
			foodItemsRequestData.Items = append(foodItemsRequestData.Items, value)
		}
	}
	return foodItemsRequestData, nil
}

func RetrieveTranslations(store_id uint, data RequestData) (result RetrieveTranslationsResult, err error) {
	translationModel := models.NewTranslationModel(config.PostgreSQL, *config.Logger)
	if translationModel == nil {
		config.Logger.Error("Failed to create translation model")
		return RetrieveTranslationsResult{}, err
	}
	result.Ad.StoreID = store_id
	result.Ad.SaleStart = data.Items[0].ValidFrom
	result.Ad.SaleEnd = data.Items[0].ValidTo
	var adIngredientMap = make(map[string]models.AdIngredient)
	for _, item := range data.Items {
		ingredient_id, err := translationModel.GetTranslationByName(item.Name)
		if err != nil {
			config.Logger.Error("Failed to get translation by name", zap.Error(err))
			return RetrieveTranslationsResult{}, err
		}
		if ingredient_id == -1 {
			result.UntranslatedIngredients = append(result.UntranslatedIngredients, item)
			continue
		}
		adIngredientMap[item.Name] = models.AdIngredient{
			IngredientID: uint(ingredient_id),
			Price:        item.CurrentPrice,
			Sale:		 item.PostPriceText,
			Name:		 item.Name,
		}
	}
	for _, value := range adIngredientMap {
		result.Ad.Ingredient = append(result.Ad.Ingredient, models.AdIngredient{
			IngredientID: value.IngredientID,
			Price:        value.Price,
			Sale:		 value.Sale,
			Name:		 value.Name,
		})
	}
	return result, nil
}

func GetIngredientNamesAndIds() (ingredientMap map[string]uint, err error) {
	ingredientModel := models.NewIngredientModel(config.PostgreSQL, *config.Logger)
	ingredientMap, err = ingredientModel.GetAllIngredientsNameID()
	if err != nil {
		return nil, err
	}
	return ingredientMap, nil
}

func AddTranslations(ctx context.Context, untranslatedIngredients []ItemData, ingredientMap map[string]uint) (adIngredients []models.AdIngredient,err error) {
	geminiModel := config.GeminiModel
	translationModel := models.NewTranslationModel(config.PostgreSQL, *config.Logger)
	ingredientModel := models.NewIngredientModel(config.PostgreSQL, *config.Logger)
	var IngredientsString string = ""
	for _, item := range untranslatedIngredients {
		IngredientsString += item.Name + ", "
	}
	fmt.Println(geminiModel)
	resp, err := geminiModel.GenerateContent(context.Background(), 
	genai.Text("Given the following names of grocery ingredients, return a simple ingredient name. The ingredient name has to be a food. For example: Fresh Green Bell Pepper -> bell peppers. Fresh Antibiotic Free Family Pack Thin Sliced Chicken Breast -> chicken breast. The output should be only a map with the key being the original input ingredient name and the value being the output simple ingredient name. For example: {\"Fresh Green Bell Pepper\": \"bell peppers\", \"Fresh Antibiotic Free Family Pack Thin Sliced Chicken Breast\": \"chicken breast\"}. Do not return anything else except a json. If the item name is two or more items (ex: Green Peppers and Cucumbers, Salmon and Ocean Perch, etc.) generalize the food (vegetables, fish, etc.) If the item is a fruit or vegetable, make sure the returned ingredient is plural. For example fresh avocados -> avocados"), genai.Text(IngredientsString))
	if err != nil {
		return nil, err
	}
	rawResponse := config.PrintResponse(resp)
	cleanedJSON := strings.TrimPrefix(rawResponse, "```json\n")
	cleanedJSON = strings.ReplaceAll(cleanedJSON, "`", "")

	var foodMap map[string]string
	var translationMap = make(map[string]uint)
	var adIngredientMap = make(map[string]models.AdIngredient)
	err = json.Unmarshal([]byte(cleanedJSON), &foodMap)
	if err != nil {
		fmt.Println("Error decoding translations gemini JSON:", err)
		fmt.Println(cleanedJSON)
		return
	}

	for _, item := range untranslatedIngredients {
		var ingredientID uint
		simpleIngredientName := foodMap[item.Name]
		// if item.Name is not in ingredientMap
		if _, ok := ingredientMap[simpleIngredientName]; !ok {
			// add ingredient and get new ingredient id 
			newIngredient, err := CreateNewIngredientByName(simpleIngredientName)
			if err != nil {
				return nil, err
			}
			ingredientID, err = ingredientModel.CreateIngredient(*newIngredient)
			if err != nil {
				return nil, err
			}
			ingredientMap[simpleIngredientName] = ingredientID
		} else {
			ingredientID = ingredientMap[simpleIngredientName]
		}
		adIngredientMap[item.Name] = models.AdIngredient{
			IngredientID: ingredientID,
			Price:        item.CurrentPrice,
			Sale:		 item.PostPriceText,
			Name:		 item.Name,
		}
		translationMap[item.Name] = ingredientID
	}
	var translations []models.Translation
	for key, value := range translationMap {
		translations = append(translations, *models.NewTranslation(key, value))
	}
	for _, value := range adIngredientMap {
	adIngredients = append(adIngredients, models.AdIngredient{
		IngredientID: value.IngredientID,
		Price:        value.Price,
		Sale:		 value.Sale,
		Name:		 value.Name,
	})
	}
	// add translations to database
	err = translationModel.CreateTranslations(translations)
	if err != nil {
		return nil, err
	}

	return adIngredients, nil
}

func CreateNewIngredientByName(name string) (ingredient *models.Ingredient, err error) {
	ingredientModel := models.NewIngredientModel(config.PostgreSQL, *config.Logger)
	geminiModel := config.GeminiModel
	var foodTypesString = ""
	for _, foodType := range ingredientModel.GetFoodTypes() {
		foodTypesString += foodType + ", "
	}
	queryString := "Given the following name of an ingredient, return only a json for the type: choose out of " + 
	foodTypesString + ", " + " and if the type is fruit or vegetable, provide a season as an array of ints." +
	"This means if let's say brussel sprouts -> Vegetable for type -> [9, 10 , 11] for season."
	resp, err := geminiModel.GenerateContent(context.Background(), 
	
	genai.Text(queryString), genai.Text(name))
	if err != nil {
		return nil, err
	}
	rawResponse := config.PrintResponse(resp)
	cleanedJSON := strings.TrimPrefix(rawResponse, "```json\n")
	cleanedJSON = strings.ReplaceAll(cleanedJSON, "`", "")

	var ingredientRaw RawGeneratedIngredientData
	err = json.Unmarshal([]byte(cleanedJSON), &ingredientRaw)
	if err != nil {
		fmt.Println("Error decoding ingredient gemini JSON:", err)
		fmt.Println(cleanedJSON)
		return nil, err
	}
	if ingredientRaw.Type == "" {
		ingredientRaw.Type = "Other"
	}
	ingredient = &models.Ingredient{
		Name: name,
		Type: ingredientModel.ToFoodType(ingredientRaw.Type),
		Season: ingredientRaw.Season,
	}
	return ingredient, nil
}

func CreateAd(ad models.Ad) (err error) {
	adModel := models.NewAdModel(config.PostgreSQL, *config.Logger)
	err = adModel.CreateAd(ad)
	if err != nil {
		adModel.Logger.Error("Failed to create ad", zap.Error(err))
		return err
	}
	return nil
}

func GetExpiredAdStores() (stores []AdProcessInput, err error) {
	storeModel := models.NewStoreModel(config.PostgreSQL, *config.Logger)
	rawStores, err := storeModel.GetExpiredAdStores()
	for _, v := range rawStores {
		var store AdProcessInput;
		store.StoreFlippName = v.FlippMerchantName
		store.ZipCode = v.Location[len(v.Location)-5:]
		store.StoreID = int(*v.ID)
	}
	return stores, nil
}