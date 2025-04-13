package models

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"github.com/jackc/pgx/v5"

)

// TranslationModel defines a struct for Translation service
type TranslationModel struct {
	PostgreSQL *pgx.Conn
	Logger    zap.Logger
}

// Translation defines a struct for translation data
type Translation struct {
	Name     string    `json:"name"`
	IngredientID uint `json:"ingredient_id"`
}

func NewTranslationModel(PostgreSQL *pgx.Conn, logger zap.Logger) *TranslationModel {
	return &TranslationModel{
		PostgreSQL: PostgreSQL,
		Logger: logger,
	}
}

// Constructor for Translation
func NewTranslation(name string, ingredient_id uint) *Translation {
	return &Translation{
		Name:     name,
		IngredientID: ingredient_id,
	}
}

// GetTranslationByName to find ingredient_id by name from database
func (i *TranslationModel) GetTranslationByName(name string) (int, error) {
	var ingredientID int
	err := i.PostgreSQL.QueryRow(context.Background(), 
		"SELECT ingredient_id FROM translation WHERE name = $1", name).Scan(&ingredientID)
	if err == pgx.ErrNoRows {
		return -1, nil
	}
	if err != nil {
		i.Logger.Error("Error getting translation by name", zap.Error(err))
		return -1, err
	}
	return ingredientID, nil
}

// CreateTranslations to add Translations to database
func (i *TranslationModel) CreateTranslations(translations []Translation) error {
	fmt.Println(translations)
	rows := [][]interface{}{}
	for _, translation := range translations {
		rows = append(rows, []interface{}{translation.Name, translation.IngredientID})
	}

	copyCount, err := i.PostgreSQL.CopyFrom(
		context.Background(),
		pgx.Identifier{"translation"},
		[]string{"name", "ingredient_id"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		i.Logger.Error("Error creating translations in database", zap.Error(err), zap.String("function", "CreateTranslations"))
		return err
	}

	i.Logger.Info("Translation added to database", zap.Int64("rows_added", copyCount))
	return nil
}