package controllers

import (
	"encoding/json"
	
	"net/http"
	"backend/main/models"

)

// RecipeController defines a struct for Recipe controller
type RecipeController struct {
	RecipeModel *models.RecipeModel
}

// NewRecipeController is a constructor for RecipeController
func NewRecipeController(model *models.RecipeModel) *RecipeController {
	return &RecipeController{
		RecipeModel: model,
	}
}

// SetRecipe creates a recipe given recipe title, link, and ingredients
func (rc *RecipeController) SetRecipe(w http.ResponseWriter, r *http.Request) {

	
}

// GetAllRecipes retrieves all recipes
func (rc *RecipeController) GetAllRecipes(w http.ResponseWriter, r *http.Request) {
	recipes, err := rc.RecipeModel.GetAllRecipes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(recipes)
}