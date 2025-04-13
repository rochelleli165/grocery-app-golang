package models

import (
	"context"
	"fmt"
	"strconv"
	"unicode"

	"backend/main/config"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type FoodType int

const (
	Fruit FoodType = iota
	Vegetable
	Meat
	Seafood
	Dairy
	Grain
	Condiments_Spices
	Bakery
	Baking
	Frozen
	Snacks
	Deli
	Canned_Goods
	Beverage
	Other
)

var foodTypeName = map[FoodType]string{
	Fruit:              "Fruit",
	Vegetable:          "Vegetable",
	Meat:               "Meat",
	Seafood:       		"Seafood",
	Dairy:              "Dairy",
	Grain: 				"GrainCereal",
	Condiments_Spices:  "Condiments/Spices",
	Bakery:             "Bakery",
	Baking:             "Baking",
	Frozen:             "Frozen",
	Snacks:             "Snack",
	Deli:               "Deli",
	Canned_Goods:       "Canned Goods",
	Beverage:           "Beverage",
	Other:              "Other",
}

func (ff FoodType) String() string {
	return foodTypeName[ff]
}

// IngredientModel defines a struct for ingredient service
type IngredientModel struct {
	PostgreSQL *pgx.Conn
	Logger    zap.Logger
}

// Ingredient defines a struct for ingredient data
type Ingredient struct {
	ID       *uint      `json:"id"`
	Name     string    `json:"name"`
	Season   *[]int    `json:"season"`
	Type     FoodType  `json:"type"`
	SourceOf *[]string `json:"sourceof"`
}

func NewIngredientModel(PostgreSQL *pgx.Conn, logger zap.Logger) *IngredientModel {
	return &IngredientModel{
		PostgreSQL: PostgreSQL,
		Logger: logger,
	}
}

// Constructor for Ingredient
func NewIngredient(id *uint, name string, season *[]int, foodtype FoodType, sourceOf *[]string) *Ingredient {
	return &Ingredient{
		ID:       id,
		Name:     name,
		Season:   season,
		Type:     foodtype,
		SourceOf: sourceOf,
	}
}

// ToFoodType converts a string to a FoodType
func (i *IngredientModel) ToFoodType(s string) FoodType {
	r := []rune(s)
    r[0] = unicode.ToUpper(r[0])
    s = string(r)
	if s == "Condiment" || s == "Spice" {
		return Condiments_Spices
	}
	for k, v := range foodTypeName {
		if v == s {
			return k
		}
	}
	return Other
}

func (i *IngredientModel) GetFoodTypes() []string {
	var foodTypes []string
	for _, v := range foodTypeName {
		foodTypes = append(foodTypes, v)
	}
	return foodTypes
}

// GetIngredientByID to find ingredient by ID from database
func (i *IngredientModel) GetIngredientByID(id uint) (Ingredient, error) {
	var ingredient Ingredient
	var typeStr string
	err := i.PostgreSQL.QueryRow(context.Background(), 
	"SELECT * FROM ingredient WHERE id = $1", id).Scan(&ingredient.ID, &ingredient.Name, &typeStr, &ingredient.Season)
	if err != nil {
		config.Logger.Error("Error processing query", zap.Error(err), zap.String("function", "GetIngredientByID"))
		return ingredient, err
	}
	ingredient.Type = i.ToFoodType(typeStr)
	return ingredient, nil
}

// GetAllIngredients retrieves all ingredients from database
func (i *IngredientModel) GetAllIngredients() ([]Ingredient, error) {
	var ingredients []Ingredient
	rows, err := i.PostgreSQL.Query(context.Background(),
	"SELECT * FROM ingredient")
	if err != nil {
		config.Logger.Error("Error getting Postgres client", zap.Error(err), zap.String("function", "GetAllIngredients"))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ing Ingredient
		var typeStr string
		err := rows.Scan(&ing.ID, &ing.Name, &typeStr, &ing.Season,)
		if err != nil {
			config.Logger.Error("Error scanning row", zap.Error(err), zap.String("function", "GetAllIngredients"))
			return nil, err
		}
		ing.Type = i.ToFoodType(typeStr)
		ingredients = append(ingredients, ing)
	}
	if err := rows.Err(); err != nil {
		config.Logger.Error("Error processing rows", zap.Error(err), zap.String("function", "GetAllIngredients"))
		return nil, err
	}
	return ingredients, nil
}

// GetAllIngredientsNameID retrieves all ingredient ids and names from database
func (i *IngredientModel) GetAllIngredientsNameID() (map[string]uint, error) {
	var mapIngredients = make(map[string]uint)
	rows, err := i.PostgreSQL.Query(context.Background(),
	"SELECT id, name FROM ingredient")
	if err != nil {
		config.Logger.Error("Error getting Firestore client", zap.Error(err), zap.String("function", "GetAllIngredientsNameID"))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var i struct {
			ID   *uint
			Name string
		}
		err := rows.Scan(&i.ID, &i.Name)
		if err != nil {
			config.Logger.Error("Error scanning row", zap.Error(err), zap.String("function", "GetAllIngredientsNameID"))
			return nil, err
		}
		mapIngredients[i.Name] = *i.ID
	}
	if err := rows.Err(); err != nil {
		config.Logger.Error("Error processing rows", zap.Error(err), zap.String("function", "GetAllIngredientsNameID"))
		return nil, err
	}
	return mapIngredients, nil
}

// CreateIngredient creates a new ingredient in the database and returns the ID
func (i *IngredientModel) CreateIngredient(ingredient Ingredient) (id uint, err error) {
	if ingredient.Season == nil {
		ingredient.Season = &[]int{}
	}
	err = i.PostgreSQL.QueryRow(context.Background(),
		"INSERT INTO ingredient (name, season, type) VALUES ($1, $2, $3) RETURNING id",
		ingredient.Name, ingredient.Season, ingredient.Type.String()).Scan(&id)
	if err != nil {
		config.Logger.Error("Error creating ingredient in database", zap.Error(err), zap.String("function", "CreateIngredient"))
		return 0, err
	}
	return id, nil
}

// CreateIngredients creates a new ingredients in the database
func (i *IngredientModel) CreateIngredients(ingredients []Ingredient) error {
	for _, ingredient := range ingredients {
		if ingredient.Season == nil {
			ingredient.Season = &[]int{}
		} 
	}
	rows := [][]interface{}{}
	for _, ingredient := range ingredients {
		rows = append(rows, []interface{}{ingredient.Name, ingredient.Season, ingredient.Type.String()})
	}

	copyCount, err := i.PostgreSQL.CopyFrom(
		context.Background(),
		pgx.Identifier{"ingredient"},
		[]string{"name", "season", "type"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		config.Logger.Error("Error creating ingredients in database", zap.Error(err), zap.String("function", "CreateIngredients"))
		return err
	}
	fmt.Println(copyCount)
	return nil
}

// UpdateIngredient updates an existing ingredient in the database
func (i *IngredientModel) UpdateIngredient(ingredient Ingredient) error {
	seasonStr := fmt.Sprintf("%v", *ingredient.Season)
	_, err := i.PostgreSQL.Query(context.Background(), 
	"UPDATE ingredients SET name =" + ingredient.Name + ", season =" + seasonStr + ", type =" + 
	ingredient.Type.String() + " WHERE id =" + strconv.FormatUint(uint64(*ingredient.ID), 10),)
	if err != nil {
		config.Logger.Error("Error updating ingredient to database", zap.Error(err), zap.String("function", "UpdateIngredient"))
		return err
	}
	return nil
}

// UpdateIngredients updates existing ingredients in the database
func (i *IngredientModel) UpdateIngredients(ingredients []Ingredient) error {
	query := ""
	for _, ingredient := range ingredients {
		seasonStr := fmt.Sprintf("%v", *ingredient.Season)
		query += "UPDATE ingredients SET name =" + ingredient.Name + ", season =" + seasonStr + 
		", type =" + ingredient.Type.String() + " WHERE id =" + strconv.FormatUint(uint64(*ingredient.ID), 10)
	}
	_, err := i.PostgreSQL.Query(context.Background(), query)
	if err != nil {
		config.Logger.Error("Error updating ingredient to database", zap.Error(err), zap.String("function", "UpdateIngredients"))
		return err
	}
	return nil
}

// DeleteIngredient deletes an existing ingredient from the database
func (i *IngredientModel) DeleteIngredient(id uint) error {
	client, err := config.FirebaseApp.Firestore(context.Background())
	if err != nil {
		config.Logger.Error("Error getting Firestore client", zap.Error(err), zap.String("function", "DeleteIngredient"))
		return err
	}
	idStr := strconv.FormatUint(uint64(id), 10)
	_, err = client.Collection("ingredients").Doc(idStr).Delete(context.Background())
	if err != nil {
		config.Logger.Error("Error deleting document", zap.Error(err), zap.String("function", "DeleteIngredient"))
		return err
	}
	return nil
}
