package handler_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/georlav/recipeapi/config"
	"github.com/georlav/recipeapi/handler"
	"github.com/georlav/recipeapi/mongoclient"
	"github.com/georlav/recipeapi/recipe"
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
		} else {
			ingredients = append(ingredients, "in5")
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

func TestHandler_Recipes(t *testing.T) {
	testData := []struct {
		params  url.Values
		results int
	}{
		{url.Values{"p": []string{}}, 10},
		{url.Values{"p": []string{"1"}}, 10},
		{url.Values{"p": []string{"2"}}, 5},
		{url.Values{"q": []string{`"test recipe 2"`}}, 1},
		{url.Values{"q": []string{"test recipe"}}, 10},
		{url.Values{"i": []string{"in1"}}, 10},
		{url.Values{"i": []string{"in2"}, "p": []string{"2"}}, 5},
		{url.Values{"i": []string{"in4"}}, 7},
		{url.Values{"i": []string{"in5"}}, 8},
		{url.Values{"i": []string{"in6"}}, 0},
		{url.Values{"q": []string{"some term with no results"}, "i": []string{"in1"}}, 0},
	}

	r, _, err := recipeRepo()
	if err != nil {
		t.Fatal(err)
	}
	h := handler.NewHandler(r, &config.Config{}, &log.Logger{})

	for i := range testData {
		tc := testData[i]

		t.Run(fmt.Sprintf(`Test Case %+v`, tc.params), func(t *testing.T) {
			req := httptest.NewRequest("GET", "/recipes?"+tc.params.Encode(), nil)

			// init response recorder
			rr := httptest.NewRecorder()
			rh := http.HandlerFunc(h.Recipes)
			rh.ServeHTTP(rr, req)

			if http.StatusOK != rr.Code {
				t.Fatalf("Wrong status code got %d expected %d", http.StatusOK, rr.Code)
			}

			if actualLen := strings.Count(rr.Body.String(), "id"); tc.results != actualLen {
				t.Fatalf("Expected %d results got %d", tc.results, actualLen)
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
		Database:                  "recipes-testdb",
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
