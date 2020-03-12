package mysql

import (
	"database/sql"

	"github.com/georlav/recipeapi/internal/db"
)

// Ingredient object
type Ingredient struct {
	db *sql.DB
}

// NewRecipe creates a recipe object
func NewIngredient(sqlDB *sql.DB) *Ingredient {
	return &Ingredient{
		db: sqlDB,
	}
}

// Upsert update or insert a new ingredient
func (i *Ingredient) Upsert(ins db.Ingredients) error {
	tx, err := i.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO ingredient (name) VALUES (?) ON DUPLICATE id=id`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := range ins {
		if _, err := stmt.Exec(ins[i]); err != nil {
			return err
		}
	}

	return tx.Commit()
}
