package database

import (
	"database/sql"
	"fmt"
)

const ingredientColumns = "i.id, i.recipe_id, i.name, i.created_at, i.updated_at"

// RecipeTable object
type IngredientTable struct {
	db   *sql.DB
	name string
}

// NewIngredientTable create a IngredientTable object
func NewIngredientTable(db *sql.DB) *IngredientTable {
	return &IngredientTable{
		db:   db,
		name: "ingredient i",
	}
}

// Get a recipe by id
func (it *IngredientTable) Get(id uint64) (*Ingredient, error) {
	// nolint:gosec
	query := fmt.Sprintf(`SELECT %s FROM ingredient i WHERE id = ?`, ingredientColumns)

	var i Ingredient
	if err := it.db.QueryRow(query, id).Scan(
		&i.ID, &i.RecipeID, &i.Name, &i.CreatedAt, &i.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &i, nil
}
