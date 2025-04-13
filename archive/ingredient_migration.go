package archive

import (
	"backend/main/config"
	"backend/main/models"
	"context"

	"fmt"
    "backend/main/controllers"
)

type FirebaseIngredient struct {
	ID       *uint      `json:"id"`
	Name     string    `json:"name"`
	Season   *[]int    `json:"season"`
	Type     string `json:"type"`
	SourceOf *[]string `json:"sourceof"`
}

func main1() {
fmt.Printf("Hello, world.\n")
    config.ConnectPostgreSQL()
    config.InitFirebase()
    config.InitLogger()
    ctx := context.Background()
    client, err := config.FirebaseApp.Database(ctx)
    if err != nil {
        fmt.Print(err)
        return
    }
    raw_ingredients := client.NewRef("ingredients")

    var result []FirebaseIngredient
    if err := raw_ingredients.Get(ctx, &result); err != nil {
        fmt.Print(err)
        return
    }
    ingredients := []models.Ingredient{}
    i := models.NewIngredientModel(config.PostgreSQL, *config.Logger)
    for _, raw_item := range result {
        ingredient := models.NewIngredient(raw_item.ID, raw_item.Name, raw_item.Season, i.ToFoodType(raw_item.Type), raw_item.SourceOf)
        ingredients = append(ingredients, *ingredient)
    }
    model := models.NewIngredientModel(config.PostgreSQL, *config.Logger)
    controller := controllers.NewIngredientController(model)

    controller.IngredientModel.CreateIngredients(ingredients)
}