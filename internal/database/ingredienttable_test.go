package database_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/georlav/recipeapi/internal/database"
)

func TestIngredientTable_Get(t *testing.T) {
	testCases := []struct {
		desc   string
		input  uint64
		output *database.Ingredient
		error  error
	}{
		{
			"Should get a recipe ingredient",
			1,
			&database.Ingredient{
				ID:       1,
				RecipeID: 1,
				Name:     "champagne",
			},
			nil,
		},
		{
			"Should fail to get an ingredient recipe",
			0,
			nil,
			sql.ErrNoRows,
		},
	}

	db, err := db()
	if err != nil {
		t.Fatal(err)
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.desc, func(t *testing.T) {
			ingredient, err := db.Ingredient.Get(tc.input)
			if err != nil && !errors.Is(err, tc.error) {
				t.Fatal(err)
			}
			if err == nil && uint64(ingredient.ID) != tc.input {
				t.Fatalf("Invalid id, expected %d got %d", tc.input, ingredient.ID)
			}
			if err == nil && ingredient.Name != tc.output.Name {
				t.Fatalf("Invalid title, expected %s got %s", tc.output.Name, ingredient.Name)
			}
		})
	}
}
