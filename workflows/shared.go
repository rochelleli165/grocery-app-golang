package workflows

import (
	"net/http"
	"backend/main/models"

)

const AdProcessingTaskQueueName = "AD_PROCESSING_TASK_QUEUE"
const GetExpiredAdStoresTaskQueueName = "GET_EXPIRED_AD_STORES_TASK_QUEUE"

type HTTPGetter interface {
	Get(url string) (*http.Response, error)
}

type RequestData struct {
	Items []ItemData `json:"items"`
}

type ItemData struct {
	L2            string   `json:"_L2"`
	Name          string   `json:"name"`
	CurrentPrice  *float32 `json:"current_price"`
	PostPriceText *string  `json:"post_price_text"`
	ValidFrom	 string   `json:"valid_from"`
	ValidTo		 string   `json:"valid_to"`
}

type RawGeneratedIngredientData struct {
	Type string `json:"type"`
	Season *[]int `json:"season"`
}

type RetrieveTranslationsResult struct {
	Ad models.Ad
	UntranslatedIngredients []ItemData
}

type AdProcessInput struct {
	StoreID int
	ZipCode string
	StoreFlippName string
}