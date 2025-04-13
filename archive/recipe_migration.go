package archive

import (
	"backend/main/config"
	"backend/main/models"
	"context"

	//"encoding/json"
	"fmt"
)

type FirebaseRecipeIngredient struct {
    Name     string    `json:"name"`
    Amount string `json:"amount"`
    Unit string `json:"unit"`
}

type FirebaseRecipe struct {
    Ingredients map[int]FirebaseRecipeIngredient `json:"ingredients"`
    RecipeTitle string `json:"recipe-title"`
    RecipeLink string `json:"recipe-link"`
}

func main3() {
    config.ConnectPostgreSQL()
    config.InitFirebase()
    config.InitLogger()
    ctx := context.Background()
    client, err := config.FirebaseApp.Database(ctx)
    if err != nil {
        fmt.Print(err)
        return
    }

    raw_recipes := client.NewRef("recipes")

    var result map[string]FirebaseRecipe
    if err := raw_recipes.Get(ctx, &result); err != nil {
        fmt.Print(err)
        return
    }
    recipes := []models.Recipe{}

    model := models.NewRecipeModel(config.PostgreSQL, *config.Logger)

    for _, raw_item := range result {
        ingredients := []models.RecipeIngredient{}
        for index, raw_ingredient := range raw_item.Ingredients {
            uintIndex := uint(index) + 1
            ingredient := models.NewRecipeIngredient(&uintIndex, raw_ingredient.Name, &raw_ingredient.Amount, &raw_ingredient.Unit)
            ingredients = append(ingredients, *ingredient)
        }
        recipe := models.NewRecipe(nil, raw_item.RecipeTitle, raw_item.RecipeLink, "Just One Cookbook", ingredients)
        recipes = append(recipes, *recipe)
    }
    model.CreateRecipes(recipes)

}
