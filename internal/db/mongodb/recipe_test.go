package mongodb_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/db"
	"github.com/georlav/recipeapi/internal/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// Setup data for handlers related tests
func TestMain(m *testing.M) {
	file, err := ioutil.ReadFile("../testdata/recipes.json")
	if err != nil {
		log.Fatal("failed to load test data", err)
	}

	var data struct{ Recipes db.Recipes }
	if err := json.Unmarshal(file, &data); err != nil {
		log.Fatal("failed to marshal testdata", err)
	}

	rc, col, err := recipeCollection()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Import data
	for i := range data.Recipes {
		if err := rc.Insert(data.Recipes[i]); err != nil {
			log.Fatalf("failed to insert test data %s", err)
		}
	}

	ec := m.Run()

	// remove imported data
	if _, err = col.DeleteMany(context.Background(), bson.D{}); err != nil {
		fmt.Printf(`Mongo mass delete error, %s`, err)
		os.Exit(1)
	}

	os.Exit(ec)
}

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
	cfg := config.Mongo{
		Host:                      "127.0.0.1",
		Port:                      27017,
		Username:                  "root",
		Password:                  "toor",
		Database:                  "recipes",
		RecipeCollection:          "recipe",
		PoolSize:                  100,
		Timeout:                   15 * time.Second,
		SetServerSelectionTimeout: 15 * time.Second,
		SetMaxConnIdleTime:        15 * time.Second,
		SetRetryWrites:            false,
	}

	// Mongo client
	client, err := mongodb.New(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf(`mongo client error, %s`, err)
	}

	// inject recipe collection to recipe object
	mdb := client.Database(cfg.Database + "-testdb")
	rCollection := mdb.Collection(cfg.RecipeCollection)

	// Create searchable index
	iv := rCollection.Indexes()
	_, err = iv.CreateOne(context.Background(), mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "title", Value: bsonx.String("text")},
		},
	})
	if err != nil {
		return nil, nil, err
	}

	return mongodb.NewRecipe(rCollection), rCollection, nil
}
