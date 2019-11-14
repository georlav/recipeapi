package recipe_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"github.com/georlav/recipeapi/mongoclient"

	"github.com/georlav/recipeapi/recipe"

	"github.com/georlav/recipeapi/config"
)

// Setup data for handlers related tests
func TestMain(m *testing.M) {
	mc, err := recipeRepo()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = mc.Insert(
		recipe.Recipe{
			Title:       "test recipe 1",
			URL:         "http://test1.dev",
			Ingredients: []string{"in1", "in2"},
			Thumbnail:   "http://img.recipepuppy.com/1.jpg",
		},
		recipe.Recipe{
			Title:       "test recipe 2",
			URL:         "http://test2.dev",
			Ingredients: []string{"in3", "in4"},
			Thumbnail:   "http://img.recipepuppy.com/2.jpg",
		},
	)
	if err != nil {
		fmt.Printf(`Mongo Import error, %s`, err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestMongoDBRepo_GetMany(t *testing.T) {
	testCases := []struct {
		params  recipe.QueryParams
		results int
	}{
		{recipe.QueryParams{}, 2},
		{recipe.QueryParams{Page: 1}, 2},
		{recipe.QueryParams{Term: `"test recipe 2"`}, 1},
		{recipe.QueryParams{Term: "test recipe"}, 2},
		{recipe.QueryParams{Term: "recipe"}, 2},
		{recipe.QueryParams{Term: "Spaghetti code"}, 0},
		{recipe.QueryParams{Page: 2}, 0},
	}

	rr, err := recipeRepo()
	if err != nil {
		t.Fatal(err)
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(fmt.Sprintf("Quering with %+v", tc.params), func(t *testing.T) {
			s, _, err := rr.GetMany(tc.params)
			if err != nil {
				t.Fatal(err)
			}

			if tc.results != len(s) {
				t.Fatalf("Should have found %d results got %d", tc.results, len(s))
			}

		})
	}
}

func recipeRepo() (*recipe.MongoRepo, error) {
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
		SetConnectTimeout:         15 * time.Second,
		SetSocketTimeout:          15 * time.Second,
		SetMaxConnIdleTime:        15 * time.Second,
		SetRetryWrites:            false,
	}

	// Init logger, discard output
	l := log.New(os.Stdout, "", 0)
	l.SetOutput(ioutil.Discard)

	// Mongo client
	client, err := mongoclient.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf(`mongo client error, %s`, err)
	}

	// Select a database collection and inject it to repo
	db := client.Database(cfg.Database + "-testdb")
	rCollection := db.Collection(cfg.RecipeCollection)

	iv := rCollection.Indexes()
	_, err = iv.CreateOne(context.Background(), mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "title", Value: bsonx.String("text")},
		},
	})
	if err != nil {
		return nil, err
	}

	return recipe.NewMongoRepo(rCollection, cfg), nil
}
