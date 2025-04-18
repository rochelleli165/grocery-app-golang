package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"backend/main/models"
	"backend/main/config"
	"go.uber.org/zap"
)

// StoreController defines a struct for Store controller
type StoreController struct {
	StoreModel *models.StoreModel
}

// NewStoreController is a constructor for StoreController
func NewStoreController(model *models.StoreModel) *StoreController {
	return &StoreController{
		StoreModel: model,
	}
}

// Add Store adds store given store_name and address
func (sc *StoreController) AddStore(w http.ResponseWriter, r *http.Request) {
	var store models.Store
	err := json.NewDecoder(r.Body).Decode(&store)
	if err != nil {
		config.Logger.Error("Error getting Firestore client", zap.Error(err), zap.String("function", "AddStore"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = sc.StoreModel.CreateStore(store)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Subscribe store subscribes user to a store given user_id and store_id
func (sc *StoreController) SubscribeStore(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	storeID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	if err != nil {
		config.Logger.Error("Error getting Firestore client", zap.Error(err), zap.String("function", "SubscribeStore"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = sc.StoreModel.SubscribeStore(uint(userID), uint(storeID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetSubscriptions returns all store_ids of a subscried stores given user_id
func (sc *StoreController) GetSubscribedStores(w http.ResponseWriter, r *http.Request) {
	
}


