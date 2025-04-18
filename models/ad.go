package models

import (
	"context"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// RecipeModel defines a struct for Recipe service
type AdModel struct {
	PostgreSQL *pgx.Conn
	Logger     zap.Logger
}

// Recipe defines a struct for Recipe data
type Ad struct {
	ID     *uint   `json:"id"`
	StoreID uint `json:"store_id"`
	SaleStart  string  `json:"sale_start"`
	SaleEnd   string `json:"sale_end"`
	Ingredient []AdIngredient 
}

type AdIngredient struct {
	ID    *uint `json:"id"`
	IngredientID     uint `json:"ingredient_id"`
	Price		   *float32 `json:"amount"`
	Sale		   *string `json:"unit"`
	Name 		 string `json:"name"`
}

func NewAdModel(PostgreSQL *pgx.Conn, logger zap.Logger) *AdModel {
	return &AdModel{
		PostgreSQL: PostgreSQL,
		Logger:     logger,
	}
}

// Constructor for Recipe
func NewAd(id *uint, store_id uint, sale_start string, sale_end string, ingredients []AdIngredient) *Ad{
	return &Ad{
		ID:     id,
		StoreID: store_id,
		SaleStart:  sale_start,
		SaleEnd:   sale_end,
		Ingredient: ingredients,
	}
}

// Constructor for RecipeIngredient
func NewAdIngredient(id *uint, ingredient_id uint, price *float32, sale *string, name string) *AdIngredient {
	return &AdIngredient{
		ID: id,
		IngredientID: ingredient_id,
		Name: name,
		Price: price,
		Sale: sale,
	}
}

// CreateAd adds an ad to the database. This adds to the ad and ad_ingredient tables
func (i *AdModel) CreateAd(ad Ad) error {
	
	var adID int
	err := i.PostgreSQL.QueryRow(context.Background(),
		"INSERT INTO ad (store_id, sale_start, sale_end) VALUES ($1, $2, $3) RETURNING id;",
		ad.StoreID, ad.SaleStart, ad.SaleEnd).Scan(&adID)
	if err != nil {
		i.Logger.Error("Error adding recipe to database", zap.Error(err))
		return err
	}

	rows := make([][]interface{}, len(ad.Ingredient))
	ingredient_ids := make([]uint, len(ad.Ingredient))
	for i, ingredient := range ad.Ingredient {
		rows[i] = []interface{}{adID, ingredient.IngredientID, ingredient.Name, ingredient.Price, ingredient.Sale}
		ingredient_ids[i] = ingredient.IngredientID
	}

	_, err = i.PostgreSQL.CopyFrom(
		context.Background(),
		pgx.Identifier{"ad_ingredient"},
		[]string{"ad_id", "ingredient_id", "name", "price", "sale"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		i.Logger.Error("Error adding ad_ingredients to database", 
		zap.Error(err), zap.Int("ad", adID),
		zap.Any("ingredient_ids", ingredient_ids))
		return err
	}
	
	i.Logger.Info("Successfully added ad to database")
	return nil

}

// GetRecentAd
func (i * AdModel) GetRecentAd(storeID uint) (ad Ad, err error) {

	err = i.PostgreSQL.QueryRow(context.Background(), 
		"SELECT id, store_id, sale_start::text, sale_end::text FROM ad WHERE store_id = $1 order by sale_start desc limit 1", storeID).Scan(&ad.ID, &ad.StoreID, &ad.SaleStart, &ad.SaleEnd)
	if err != nil {
		i.Logger.Error("Error getting ad by ID", zap.Error(err))
		return Ad{}, err
	}
	rows, err := i.PostgreSQL.Query(context.Background(),
		"SELECT ingredient_id, name, price, sale FROM ad_ingredient WHERE ad_id = $1", ad.ID)
	if err != nil {
		i.Logger.Error("Error getting ad ingredients", zap.Error(err))
		return Ad{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var ai AdIngredient
		err := rows.Scan(&ai.IngredientID, &ai.Name, &ai.Price, &ai.Sale)
		if err != nil {
			i.Logger.Error("Error scanning row", zap.Error(err))
			return Ad{}, err
		}
		ad.Ingredient = append(ad.Ingredient, ai)
	}
	return ad, nil
}
