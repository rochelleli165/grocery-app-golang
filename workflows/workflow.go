package workflows

import (
	"time"
	"backend/main/models"
	"backend/main/config"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
)

func RetrieveAds(ctx workflow.Context) (stores []AdProcessInput) {
	logger := config.Logger
	err := workflow.ExecuteActivity(ctx, GetExpiredAdStores).Get(ctx, &stores)
	if err != nil {
		logger.Error("Failed to get ads", zap.Error(err))
	}
	return stores
}

func AdProcess(ctx workflow.Context, store_id int, zip_code string, store_flipp_name string) (error) {
	logger := config.Logger
	
	// RetryPolicy specifies how to automatically handle retries if an Activity fails.
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     2.0,
		MaximumInterval:        100 * time.Second,
		MaximumAttempts:        10, // 0 is unlimited retries
	}

	options := workflow.ActivityOptions {
		StartToCloseTimeout: 5 * time.Minute,
		RetryPolicy: retrypolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var rawData RequestData

	err := workflow.ExecuteActivity(ctx, GetFlippData, zip_code, store_flipp_name).Get(ctx, &rawData)
	if err != nil {
		logger.Error("Failed to get Flipp data", zap.Error(err))
		return err
	}
	
	var result RetrieveTranslationsResult
	err = workflow.ExecuteActivity(ctx, RetrieveTranslations, store_id, rawData).Get(ctx, &result)
	if err != nil {
		logger.Error("Failed to retrieve translations", zap.Error(err))
		return err
	}
	ad := result.Ad
	untranslatedIngredients := result.UntranslatedIngredients

	var ingredientMap map[string]uint
	err = workflow.ExecuteActivity(ctx, GetIngredientNamesAndIds).Get(ctx, &ingredientMap)
	if err != nil {
		logger.Error("Failed to get ingredient names and ids", zap.Error(err))
		return err
	}

	var translatedAdIngredients []models.AdIngredient
	if len(untranslatedIngredients) > 0 {
		err = workflow.ExecuteActivity(ctx, AddTranslations, untranslatedIngredients, ingredientMap).Get(ctx, &translatedAdIngredients)
		if err != nil {
			logger.Error("Failed to add translations", zap.Error(err))
			return err
		}
	}
	
	ad.Ingredient = append(ad.Ingredient, translatedAdIngredients...)
	
	err = workflow.ExecuteActivity(ctx, CreateAd, ad).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to create ad", zap.Error(err))
	}
	return nil
}

