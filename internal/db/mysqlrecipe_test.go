package db_test

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/georlav/recipeapi/internal/db"

	"github.com/georlav/recipeapi/internal/config"
)

func TestNewRecipeTable_Get(t *testing.T) {
	testCases := []struct {
		desc   string
		input  string
		output *db.Recipe
		error  error
	}{
		{
			"Should get a recipe",
			"1",
			&db.Recipe{
				ID:          "1",
				Title:       "Ginger Champagne",
				Ingredients: []string{"champagne", "ginger", "ice", "vodka"},
				URL:         "http://allrecipes.com/Recipe/Ginger-Champagne/Detail.aspx",
				Thumbnail:   "http://img.recipepuppy.com/1.jpg",
			},
			nil,
		},
		{
			"Should fail to get a recipe",
			"0",
			nil,
			sql.ErrNoRows,
		},
	}

	table, _, err := recipeTbl()
	if err != nil {
		t.Fatal(err)
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.desc, func(t *testing.T) {
			res, err := table.Get(tc.input)
			if err != nil && !errors.Is(err, tc.error) {
				t.Fatal(err)
			}
			if err == nil && res.ID != tc.input {
				t.Fatalf("Invalid id, expected %s got %s", tc.input, res.ID)
			}
			if err == nil && res.Title != tc.output.Title {
				t.Fatalf("Invalid title, expected %s got %s", tc.output.Title, res.Title)
			}
			if err == nil && !reflect.DeepEqual(tc.output.Ingredients, res.Ingredients) {
				t.Fatalf("Invalid ingredients, expected %+v got %+v", tc.output.Ingredients, res.Ingredients)
			}
		})
	}
}

func TestNewRecipeTable_Insert(t *testing.T) {
	testCases := []struct {
		desc  string
		input db.Recipe
		error error
	}{
		{
			"Should insert a recipe",
			db.Recipe{
				Title:       "test",
				Ingredients: []string{"champagne", "ginger", "ice", "vodka"},
				URL:         "http://allrecipes.com/Recipe/Ginger-Champagne/Detail.aspx",
				Thumbnail:   "http://img.recipepuppy.com/1.jpg",
			},
			nil,
		},
		{
			"Should fail to insert recipe as it already exists",
			db.Recipe{
				Title:     "Ginger Champagne",
				URL:       "http://allrecipes.com/Recipe/Ginger-Champagne/Detail.aspx",
				Thumbnail: "http://img.recipepuppy.com/1.jpg",
			},
			errors.New("recipe error, Error 1062: Duplicate entry 'Ginger Champagne' for key 'recipe_title_uindex'"),
		},
	}

	table, mdb, err := recipeTbl()
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
			err := table.Insert(tc.input)
			if err != nil && err.Error() != tc.error.Error() {
				t.Fatal(err)
			}
		})
	}
}

func TestRecipeTable_Paginate(t *testing.T) {
	testCases := []struct {
		page             uint64
		filters          *db.Filters
		resultCount      int64
		resultTotalCount int64
	}{
		{0, nil, 10, 22},
		{1, nil, 10, 22},
		{2, nil, 10, 22},
		{3, nil, 2, 22},
		{1, &db.Filters{Term: "Ginger Champagne"}, 1, 1},
		{1, &db.Filters{Term: "potato"}, 4, 4},
		{1, &db.Filters{Term: "onion"}, 1, 1},
		{1, &db.Filters{Term: "onion", Ingredients: []string{"onions"}}, 1, 1},
		{1, &db.Filters{Ingredients: []string{"onions"}}, 8, 8},
		{1, &db.Filters{Ingredients: []string{"eggs"}}, 5, 5},
		{1, &db.Filters{Ingredients: []string{"eggs", "onions"}}, 10, 12},
		{1, &db.Filters{Term: "pork"}, 3, 3},
		{1, &db.Filters{Term: "pork", Ingredients: []string{"garlic"}}, 2, 2},
		{1, &db.Filters{Term: "pork", Ingredients: []string{"garlic", "brown sugar"}}, 2, 2},
		{1, &db.Filters{Term: "park", Ingredients: []string{"garlic", "brown sugar"}}, 0, 0},
		{1, &db.Filters{Term: "potato", Ingredients: []string{"eggs"}}, 1, 1},
		{1, &db.Filters{Term: "Spaghetti code"}, 0, 0},
	}

	rc, _, err := recipeTbl()
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

func recipeTbl() (*db.RecipeTable, *sql.DB, error) {
	cfg, err := config.Load("testdata/config.json")
	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := db.NewMySQL(cfg.MySQL)
	if err != nil {
		return nil, nil, err
	}

	return db.NewRecipeTable(sqlDB), sqlDB, nil
}
