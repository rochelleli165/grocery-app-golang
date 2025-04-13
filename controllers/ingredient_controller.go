package controllers

import (
	"encoding/json"
	"strconv"

	"backend/main/config"
	"backend/main/models"
	"net/http"

	"go.uber.org/zap"
)

// IngredientController defines a struct for ingredient controller
type IngredientController struct {
	IngredientModel *models.IngredientModel
}

// NewIngredientController is a constructor for IngredientController
func NewIngredientController(model *models.IngredientModel) *IngredientController {
	return &IngredientController{
		IngredientModel: model,
	}
}

// GetIngredientByID is a function to get ingredient by ID
func (i *IngredientController) GetIngredientByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	ingredientID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	ingredient, err := i.IngredientModel.GetIngredientByID(uint(ingredientID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ingredient)
}

// CreateIngredients is a function to create new ingredient(s)
func (i *IngredientController) CreateIngredients(w http.ResponseWriter, r *http.Request) {
	var ingredients []models.Ingredient
	err := json.NewDecoder(r.Body).Decode(&ingredients)
	if err != nil {
		config.Logger.Error("Error getting Firestore client", zap.Error(err), zap.String("function", "FindIngredientByID"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = i.IngredientModel.CreateIngredients(ingredients)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// UpdateIngredient is a function to update an existing ingredient
func (i *IngredientController) UpdateIngredient(w http.ResponseWriter, r *http.Request) {
	var ingredient models.Ingredient
	err := json.NewDecoder(r.Body).Decode(&ingredient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = i.IngredientModel.UpdateIngredient(ingredient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteIngredient is a function to delete an existing ingredient
func (i *IngredientController) DeleteIngredient(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	ingredientID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	err = i.IngredientModel.DeleteIngredient(uint(ingredientID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetAllIngredients is a function to get all ingredients
func (i *IngredientController) GetAllIngredients(w http.ResponseWriter, r *http.Request) {
	ingredients, err := i.IngredientModel.GetAllIngredients()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ingredients)
}

