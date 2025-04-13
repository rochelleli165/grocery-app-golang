package controllers

import (
	"encoding/json"
	
	"net/http"
	"strconv"
	"backend/main/models"

)

// PantryController defines a struct for Pantry controller
type PantryController struct {
	PantryModel *models.PantryModel
}

// NewPantryController is a constructor for PantryController
func NewPantryController(model *models.PantryModel) *PantryController {
	return &PantryController{
		PantryModel: model,
	}
}

// GetPantry retrieves all pantry items
func (pc *PantryController) GetPantryItems(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	pantry, err := pc.PantryModel.GetPantryItems(uint(userID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pantry)
}

// AddToPantry takes in user_id and cart items and adds cart items to user's pantry
func (pc *PantryController) AddPantryIngredients(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var pantryItems []models.PantryIngredient
	err = json.NewDecoder(r.Body).Decode(&pantryItems)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pantry := models.NewPantry(nil, uint(userID), pantryItems)
	err = pc.PantryModel.AddPantryIngredients(*pantry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
	