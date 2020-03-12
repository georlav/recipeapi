package db_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/db"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestMongoDBRepo_Paginate(t *testing.T) {
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
		{1, &db.Filters{Ingredients: []string{"eggs", "onions"}}, 1, 1},
		{1, &db.Filters{Term: "pork"}, 3, 3},
		{1, &db.Filters{Term: "pork", Ingredients: []string{"garlic"}}, 2, 2},
		{1, &db.Filters{Term: "pork", Ingredients: []string{"garlic", "brown sugar"}}, 1, 1},
		{1, &db.Filters{Term: "park", Ingredients: []string{"garlic", "brown sugar"}}, 0, 0},
		{1, &db.Filters{Term: "potato", Ingredients: []string{"eggs"}}, 1, 1},
		{1, &db.Filters{Term: "Spaghetti code"}, 0, 0},
	}

	rc, _, err := recipeCollection()
	if err != nil {
		t.Fatal(err)
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(fmt.Sprintf("Quering with %+v", tc.filters), func(t *testing.T) {
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

func recipeCollection() (db.Queryable, *mongo.Collection, error) {
	cfg, err := config.Load("testdata/config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Mongo client
	client, err := db.NewMongoDB(cfg.Mongo)
	if err != nil {
		return nil, nil, fmt.Errorf(`mongo client error, %s`, err)
	}
	mdb := client.Database(cfg.Mongo.Database)
	col := mdb.Collection(cfg.Mongo.RecipeCollection)

	// Init Mongo queries
	cfg.APP.Database = "MongoDB"
	mongoQueries, err := db.New(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	return mongoQueries, col, nil
}
