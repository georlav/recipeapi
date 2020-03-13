package database_test

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/georlav/recipeapi/internal/database"

	"github.com/georlav/recipeapi/internal/config"
)

func TestNewRecipeTable_Get(t *testing.T) {
	testCases := []struct {
		desc   string
		input  uint64
		output *database.Recipe
		error  error
	}{
		{
			"Should get a recipe",
			1,
			&database.Recipe{
				ID:    1,
				Title: "Ginger Champagne",
				Ingredients: database.Ingredients{
					{Name: "champagne"},
					{Name: "ginger"},
					{Name: "ice"},
					{Name: "vodka"},
				},
				URL:       "http://allrecipes.com/Recipe/Ginger-Champagne/Detail.aspx",
				Thumbnail: "http://img.recipepuppy.com/1.jpg",
			},
			nil,
		},
		{
			"Should fail to get a recipe",
			0,
			nil,
			sql.ErrNoRows,
		},
	}

	rt, _, err := recipeTable()
	if err != nil {
		t.Fatal(err)
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.desc, func(t *testing.T) {
			recipe, err := rt.Get(tc.input)
			if err != nil && !errors.Is(err, tc.error) {
				t.Fatal(err)
			}
			if err == nil && uint64(recipe.ID) != tc.input {
				t.Fatalf("Invalid id, expected %d got %d", tc.input, recipe.ID)
			}
			if err == nil && recipe.Title != tc.output.Title {
				t.Fatalf("Invalid title, expected %s got %s", tc.output.Title, recipe.Title)
			}
			if err == nil && len(tc.output.Ingredients) != len(recipe.Ingredients) {
				t.Fatalf(
					"Invalid ingredient length, expected %d got %d",
					len(tc.output.Ingredients),
					len(recipe.Ingredients),
				)
			}
		})
	}
}

func TestNewRecipeTable_Insert(t *testing.T) {
	testCases := []struct {
		desc  string
		input database.Recipe
		error error
	}{
		{
			"Should insert a recipe",
			database.Recipe{
				Title: "test",
				Ingredients: database.Ingredients{
					{Name: "champagne"},
					{Name: "ginger"},
					{Name: "ice"},
					{Name: "vodka"},
				},
				URL:       "http://allrecipes.com/Recipe/Ginger-Champagne/Detail.aspx",
				Thumbnail: "http://img.recipepuppy.com/1.jpg",
			},
			nil,
		},
		{
			"Should fail to insert recipe as it already exists",
			database.Recipe{
				Title:     "Ginger Champagne",
				URL:       "http://allrecipes.com/Recipe/Ginger-Champagne/Detail.aspx",
				Thumbnail: "http://img.recipepuppy.com/1.jpg",
			},
			errors.New("recipe error, Error 1062: Duplicate entry 'Ginger Champagne' for key 'recipe_title_uindex'"),
		},
	}

	table, mdb, err := recipeTable()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := mdb.Exec("delete from recipe where id = 23"); err != nil {
			t.Fatal(err)
		}
	}()

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.desc, func(t *testing.T) {
			id, err := table.Insert(tc.input)
			if err != nil && err.Error() != tc.error.Error() {
				t.Fatal(err)
			}
			if err == nil && id == 0 {
				t.Fatal("Recipe expected to have an id")
			}
		})
	}
}

func TestRecipeTable_Paginate(t *testing.T) {
	testCases := []struct {
		page             uint64
		filters          *database.RecipeFilters
		resultCount      int64
		resultTotalCount int64
	}{
		{0, nil, 10, 22},
		{1, nil, 10, 22},
		{2, nil, 10, 22},
		{3, nil, 2, 22},
		{1, &database.RecipeFilters{Term: "Ginger Champagne"}, 1, 1},
		{1, &database.RecipeFilters{Term: "potato"}, 4, 4},
		{1, &database.RecipeFilters{Term: "onion"}, 1, 1},
		{1, &database.RecipeFilters{Term: "onion", Ingredients: []string{"onions"}}, 1, 1},
		{1, &database.RecipeFilters{Ingredients: []string{"onions"}}, 8, 8},
		{1, &database.RecipeFilters{Ingredients: []string{"eggs"}}, 5, 5},
		{1, &database.RecipeFilters{Ingredients: []string{"eggs", "onions"}}, 10, 12},
		{1, &database.RecipeFilters{Term: "pork"}, 3, 3},
		{1, &database.RecipeFilters{Term: "pork", Ingredients: []string{"garlic"}}, 2, 2},
		{1, &database.RecipeFilters{Term: "pork", Ingredients: []string{"garlic", "brown sugar"}}, 2, 2},
		{1, &database.RecipeFilters{Term: "park", Ingredients: []string{"garlic", "brown sugar"}}, 0, 0},
		{1, &database.RecipeFilters{Term: "potato", Ingredients: []string{"eggs"}}, 1, 1},
		{1, &database.RecipeFilters{Term: "Spaghetti code"}, 0, 0},
	}

	rc, _, err := recipeTable()
	if err != nil {
		t.Fatal(err)
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(fmt.Sprintf("Request page %d with filters %+v", tc.page, tc.filters), func(t *testing.T) {
			t.Parallel()

			recipes, total, err := rc.Paginate(tc.page, tc.filters)
			if err != nil {
				t.Fatal(err)
			}
			if actualCount := int64(len(recipes)); tc.resultCount != actualCount {
				t.Fatalf("Should have found %d results got %d", tc.resultCount, actualCount)
			}
			if tc.resultTotalCount != total {
				t.Fatalf("Should have found %d total results got %d", tc.resultTotalCount, total)
			}
		})
	}
}

func recipeTable() (*database.RecipeTable, *sql.DB, error) {
	cfg, err := config.Load("testdata/config.json")
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewClient(cfg.MySQL)
	if err != nil {
		log.Fatal(err)
	}

	return database.NewRecipeTable(db), db, err
}
