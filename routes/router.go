package routes

import (
	"backend/main/config"
	"backend/main/controllers"
	"backend/main/models"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	router := mux.NewRouter()

	// Initialize the ingredient controller
	ingredientModel := models.NewIngredientModel(config.PostgreSQL, *config.Logger)
	ingredientController := controllers.NewIngredientController(ingredientModel)

	router.HandleFunc("/api/GetAllIngredients", ingredientController.GetAllIngredients).Methods("GET")
	router.HandleFunc("/api/GetIngredient", ingredientController.GetIngredientByID).Methods("GET")
	router.HandleFunc("/api/CreateIngredient", ingredientController.CreateIngredients).Methods("POST")
	router.HandleFunc("/api/UpdateIngredient", ingredientController.UpdateIngredient).Methods("PUT")

	recipeModel := models.NewRecipeModel(config.PostgreSQL, *config.Logger)
	recipeController := controllers.NewRecipeController(recipeModel)
	router.HandleFunc("/api/GetAllRecipes", recipeController.GetAllRecipes).Methods("GET")


	pantryModel := models.NewPantryModel(config.PostgreSQL, *config.Logger)
	pantryController := controllers.NewPantryController(pantryModel)

	router.HandleFunc("/api/GetPantryItems", pantryController.GetPantryItems).Methods("GET")
	router.HandleFunc("/api/AddPantryIngredients", pantryController.AddPantryIngredients).Methods("POST")

	storeModel := models.NewStoreModel(config.PostgreSQL, *config.Logger)
	storeController := controllers.NewStoreController(storeModel)

	router.HandleFunc("/api/AddStore/{store_name,store_address}/", storeController.AddStore).Methods("POST")
	router.HandleFunc("/api/SubscribeStore/", storeController.SubscribeStore).Methods("POST")
	router.HandleFunc("/api/GetSubscriptions/{user_id}/", storeController.GetSubscriptions).Methods("GET")

	adModel := models.NewAdModel(config.PostgreSQL, *config.Logger)
	adController := controllers.NewAdController(adModel)

	router.HandleFunc("/api/GetAd/", adController.GetAd).Methods("GET")


/*
	router.HandleFunc("/api/GetTranslation/{ingredient_name}", )
	router.HandleFunc("/api/UpdateTranslations/{translations}", )
	router.HandleFunc("/api/SetRecipe/{recipe_title,recipe_link,ingredients}", )
	
	router.HandleFunc("/api/OptimizeAndGetRecipes/{ingredients}",)

*/
	return router
}