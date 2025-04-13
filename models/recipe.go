package models

import (
	"context"


	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// RecipeModel defines a struct for Recipe service
type RecipeModel struct {
	PostgreSQL *pgx.Conn
	Logger     zap.Logger
}

// Recipe defines a struct for Recipe data
type Recipe struct {
	ID     *uint   `json:"id"`
	Title  string  `json:"title"`
	Link   *string `json:"link"`
	Author *string `json:"author"`
	Ingredient []RecipeIngredient 
}

type RecipeIngredient struct {
	IngredientID     *uint `json:"id"`
	Amount		   *string `json:"amount"`
	Unit		   *string `json:"unit"`
	Name 		 string `json:"name"`
}

func NewRecipeModel(PostgreSQL *pgx.Conn, logger zap.Logger) *RecipeModel {
	return &RecipeModel{
		PostgreSQL: PostgreSQL,
		Logger:     logger,
	}
}

// Constructor for Recipe
func NewRecipe(id *uint, title string, link string, author string, ingredients []RecipeIngredient) *Recipe {
	return &Recipe{
		ID:     id,
		Title:  title,
		Link:   &link,
		Author: &author,
		Ingredient: ingredients,
	}
}

// Constructor for RecipeIngredient
func NewRecipeIngredient(id *uint, name string, amount *string, unit *string) *RecipeIngredient {
	return &RecipeIngredient{
		IngredientID: id,
		Name: name,
		Amount: amount,
		Unit: unit,
	}
}

// CreateRecipes adds recipes to the database. This adds to the recipe and recipe_ingredient tables
func (i *RecipeModel) CreateRecipes(recipes []Recipe) error {
	
	for _, recipe := range recipes {
		var recipeID int
		err := i.PostgreSQL.QueryRow(context.Background(),
			"INSERT INTO recipe (title, link, author) VALUES ($1, $2, $3) RETURNING id;",
			recipe.Title, recipe.Link, recipe.Author).Scan(&recipeID)
		if err != nil {
			i.Logger.Error("Error adding recipe to database", zap.Error(err))
			return err
		}

		rows := make([][]interface{}, len(recipe.Ingredient))
		ingredient_ids := make([]int, len(recipe.Ingredient))
		for i, ingredient := range recipe.Ingredient {
			rows[i] = []interface{}{recipeID, ingredient.IngredientID, ingredient.Amount, ingredient.Unit, ingredient.Name}
			ingredient_ids[i] = int(*ingredient.IngredientID)
		}

		_, err = i.PostgreSQL.CopyFrom(
			context.Background(),
			pgx.Identifier{"recipe_ingredient"},
			[]string{"recipe_id", "ingredient_id", "amount", "unit", "name"},
			pgx.CopyFromRows(rows),
		)
		if err != nil {
			i.Logger.Error("Error adding recipe_ingredients to database", 
			zap.Error(err), zap.String("recipe", recipe.Title),
			zap.Any("ingredient_ids", ingredient_ids))
			return err
		}
	}
	i.Logger.Info("Successfully added recipes to database")
	return nil

}

// GetAllRecipes returns all recipes from the database
func (i *RecipeModel) GetAllRecipes() ([]Recipe, error) {
	var recipes []Recipe
	rows, err := i.PostgreSQL.Query(context.Background(),
		"SELECT id, title, link, author FROM recipe")
	if err != nil {
		i.Logger.Error("Error getting all recipes", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var recipe Recipe
		err := rows.Scan(&recipe.ID, &recipe.Title, &recipe.Link, &recipe.Author)
		if err != nil {
			i.Logger.Error("Error scanning row", zap.Error(err))
			return nil, err
		}
		recipes = append(recipes, recipe)
	}
	if err := rows.Err(); err != nil {
		i.Logger.Error("Error processing rows", zap.Error(err))
		return nil, err
	}
	return recipes, nil
}