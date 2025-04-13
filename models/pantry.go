package models

import (
	"context"


	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// PantryModel defines a struct for pantry service
type PantryModel struct {
	PostgreSQL *pgx.Conn
	Logger     zap.Logger
}

// Pantry defines a struct for pantry data
type Pantry struct {
	ID     *uint   `json:"id"`
	UserID  uint  `json:"user_id"`
	Ingredient []PantryIngredient 
}

type PantryIngredient struct {
	IngredientID     *uint `json:"id"`
	Quantity		 *string `json:"quantity"`
	Unit		   *string `json:"unit"`
}

func NewPantryModel(PostgreSQL *pgx.Conn, logger zap.Logger) *PantryModel {
	return &PantryModel{
		PostgreSQL: PostgreSQL,
		Logger:     logger,
	}
}

// Constructor for Pantry
func NewPantry(id *uint, userID uint, ingredients []PantryIngredient) *Pantry {
	return &Pantry{
		ID:    id,
		UserID: userID,
		Ingredient: ingredients,
	}
}

// Constructor for PantryIngredient
func NewPantryIngredient(id *uint, quantity *string, unit *string) *PantryIngredient {
	return &PantryIngredient{
		IngredientID: id,
		Quantity: quantity,
		Unit: unit,
	}
}

// GetPantry gets all pantry ingredients given user id
func (i *PantryModel) GetPantryItems(userID uint) ([]PantryIngredient, error) {
	var pantryIngredients []PantryIngredient
	rows, err := i.PostgreSQL.Query(context.Background(), 
		"SELECT pi.id, pi.ingredient_id, pi.amount, pi.unit, pi.name FROM Pantry_ingredient pi INNER JOIN Pantry p ON pi.Pantry_id = p.id WHERE p.user_id = $1", userID)
	if err != nil {
		i.Logger.Error("Error getting pantry ingredients", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pantryIngredient PantryIngredient
		err := rows.Scan(&pantryIngredient.IngredientID, &pantryIngredient.Quantity, &pantryIngredient.Unit)
		if err != nil {
			i.Logger.Error("Error scanning pantry ingredient", zap.Error(err))
			return nil, err
		}
		pantryIngredients = append(pantryIngredients, pantryIngredient)
	}

	return pantryIngredients, nil
}

// CreatePantrys adds Pantrys to the database. This adds to the Pantry and Pantry_ingredient tables
func (i *PantryModel) AddPantryIngredients(pantry Pantry) error {
	
	for _, ingredient := range pantry.Ingredient {
		var PantryID int
		err := i.PostgreSQL.QueryRow(context.Background(),
			"INSERT INTO Pantry (user_id, ingredient_id, quantity, unit) VALUES ($1, $2, $3) RETURNING id;",
			pantry.UserID, ingredient.IngredientID, ingredient.Quantity, ingredient.Unit).Scan(&PantryID)
		if err != nil {
			i.Logger.Error("Error adding pantry ingredients to database", zap.Error(err))
			return err
		}

	}
	i.Logger.Info("Successfully added Pantrys to database")
	return nil

}