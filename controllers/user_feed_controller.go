package controllers

import (
	"backend/main/services"
	"backend/main/models"
	"backend/main/config"
)

func InitUserFeedController() *services.UserFeedService {
	userModel := models.NewUserModel(config.PostgreSQL, *config.Logger)
	storeModel := models.NewStoreModel(config.PostgreSQL, *config.Logger)
	adModel := models.NewAdModel(config.PostgreSQL, *config.Logger)
	ingredientModel := models.NewIngredientModel(config.PostgreSQL, *config.Logger)

	return services.NewUserFeedService(userModel, storeModel, adModel, ingredientModel)
}