package recipe_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/georlav/recipeapi/config"
	"github.com/georlav/recipeapi/mongoclient"
	"github.com/georlav/recipeapi/recipe"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// Setup data for mongo related tests
func TestMain(m *testing.M) {
	mc, col, err := recipeRepo()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Import 15 recipes
	recipes := recipe.Recipes{}
	for i := 1; i <= 15; i++ {
		ingredients := []string{"in1", "in2", "in3"}
		if i%2 == 0 {
			ingredients = append(ingredients, "in4")
		}

		recipes = append(recipes, recipe.Recipe{
			Title:       fmt.Sprintf("test recipe %d", i),
			URL:         fmt.Sprintf("http://test%d.dev", i),
			Ingredients: ingredients,
			Thumbnail:   "http://img.recipepuppy.com/1.jpg",
		})
	}

	if err = mc.Insert(recipes...); err != nil {
		fmt.Printf(`Mongo Import error, %s`, err)
		os.Exit(1)
	}

	ec := m.Run()

	// remove imported data
	if _, err = col.DeleteMany(context.Background(), bson.D{}); err != nil {
		fmt.Printf(`Mongo mass delete error, %s`, err)
		os.Exit(1)
	}

	os.Exit(ec)
}

func TestMongoDBRepo_GetMany(t *testing.T) {
	testCases := []struct {
		params  recipe.QueryParams
		results int
		total   int64
	}{
		{recipe.QueryParams{}, 10, 15},
		{recipe.QueryParams{Page: 1}, 10, 15},
		{recipe.QueryParams{Page: 2}, 5, 15},
		{recipe.QueryParams{Term: `"test recipe 2"`}, 1, 1},
		{recipe.QueryParams{Term: "test recipe"}, 10, 15},
		{recipe.QueryParams{Term: "recipe"}, 10, 15},
		{recipe.QueryParams{Term: "recipe", Page: 2}, 5, 15},
		{recipe.QueryParams{Term: "recipe", Page: 2, Ingredients: []string{"in1"}}, 5, 15},
		{recipe.QueryParams{Term: "recipe", Page: 2, Ingredients: []string{"in1", "in2"}}, 5, 15},
		{recipe.QueryParams{Term: "recipe", Ingredients: []string{"in5"}}, 0, 0},
		{recipe.QueryParams{Term: "recipe", Ingredients: []string{"in1", "unknown"}}, 0, 0},
		{recipe.QueryParams{Term: "recipe", Ingredients: []string{"in4"}}, 7, 7},
		{recipe.QueryParams{Term: "Spaghetti code"}, 0, 0},
	}

	rr, _, err := recipeRepo()
	if err != nil {
		t.Fatal(err)
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(fmt.Sprintf("Quering with %+v", tc.params), func(t *testing.T) {
			t.Parallel()

			s, c, err := rr.GetMany(tc.params)
			if err != nil {
				t.Fatal(err)
			}

			if tc.results != len(s) {
				t.Fatalf("Should have found %d results got %d", tc.results, len(s))
			}

			if tc.total != c {
				t.Fatalf("Should have found %d total results got %d", tc.total, c)
			}
		})
	}
}

func recipeRepo() (*recipe.MongoRepo, *mongo.Collection, error) {
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
	client, err := mongoclient.NewClient(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf(`mongo client error, %s`, err)
	}

	// Select a database collection and inject it to repo
	db := client.Database(cfg.Database + "-testdb")
	rCollection := db.Collection(cfg.RecipeCollection)

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

	return recipe.NewMongoRepo(rCollection, cfg), rCollection, nil
}
