package controllers

import (
	"encoding/json"
	"strconv"
	
	
	"net/http"
	"backend/main/models"

)

// AdController defines a struct for Ad controller
type AdController struct {
	AdModel *models.AdModel
}

// NewAdController is a constructor for AdController
func NewAdController(model *models.AdModel) *AdController {
	return &AdController{
		AdModel: model,
	}
}

func (a *AdController) GetAd(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	adID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	ad, err := a.AdModel.GetRecentAd(uint(adID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ad)
}