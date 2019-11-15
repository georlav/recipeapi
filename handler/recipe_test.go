package handler_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/georlav/recipeapi/config"
	"github.com/georlav/recipeapi/mongoclient"
	"github.com/georlav/recipeapi/recipe"

	"github.com/georlav/recipeapi/handler"
)

func TestHandler_Recipes(t *testing.T) {
	testData := []struct {
		params  url.Values
		results int
	}{
		{url.Values{"p": []string{}}, 10},
		{url.Values{"p": []string{"1"}}, 10},
		{url.Values{"p": []string{"2"}}, 5},
		{url.Values{"q": []string{"sometext"}}, 0},
		{url.Values{"q": []string{"13"}}, 1},
		{url.Values{"i": []string{"in1"}}, 10},
		{url.Values{"i": []string{"in2"}, "p": []string{"2"}}, 5},
	}

	for i := range testData {
		tc := testData[i]

		t.Run(fmt.Sprintf(`Test Case %+v`, tc.params), func(t *testing.T) {
			req := httptest.NewRequest("GET", "/recipes?"+tc.params.Encode(), nil)

			// init handlers
			h, err := handlers()
			if err != nil {
				t.Fatalf("failed to init handlers, %s", err)
			}

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

func handlers() (*handler.Handler, error) {
	r, err := repository()
	if err != nil {
		return nil, err
	}

	// initialize logger
	logger := log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile)

	// Import 15 recipes
	recipes := recipe.Recipes{}
	for i := 1; i <= 15; i++ {
		recipes = append(recipes, recipe.Recipe{
			Title:       fmt.Sprintf("test recipe %d", i),
			URL:         fmt.Sprintf("http://test%d.dev", i),
			Ingredients: []string{"in1", "in2", "in3"},
			Thumbnail:   "http://img.recipepuppy.com/1.jpg",
		})
	}
	err = r.Insert(recipes...)
	if err != nil {
		fmt.Printf(`Mongo Import error, %s`, err)
		os.Exit(1)
	}

	return handler.NewHandler(r, &config.Config{}, logger), nil
}

func repository() (*recipe.MongoRepo, error) {
	cfg := config.Mongo{
		Host:             "127.0.0.1",
		Port:             27017,
		Username:         "root",
		Password:         "toor",
		Database:         "recipes-testdb",
		RecipeCollection: "recipe",
		PoolSize:         100,
		Timeout:          15,
		SetRetryWrites:   false,
	}

	// Mongo client
	client, err := mongoclient.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	// mongo client/collection
	db := client.Database(cfg.Database)
	rCollection := db.Collection(cfg.RecipeCollection)

	// recipe repository
	return recipe.NewMongoRepo(rCollection, cfg), nil
}
