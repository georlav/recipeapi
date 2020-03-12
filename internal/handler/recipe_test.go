package handler_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/db"
	"github.com/georlav/recipeapi/internal/handler"
)

func TestMain(m *testing.M) {
	// load config
	cfg, err := config.Load("testdata/config.json")
	if err != nil {
		log.Fatal(err)
	}

	// load data to import
	jd, err := ioutil.ReadFile("testdata/recipes.json")
	if err != nil {
		log.Fatalf("failed to load test data, %s", err)
	}

	// Create recipes from data
	var data struct{ Recipes db.Recipes }
	if err := json.Unmarshal(jd, &data); err != nil {
		log.Fatalf("failed to marshal testdata, %s", err)
	}

	// Get a recipe handle
	recipeTbl, err := db.New(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Import data
	for i := range data.Recipes {
		if err := recipeTbl.Insert(data.Recipes[i]); err != nil {
			log.Fatalf("failed to insert test data, %s", err)
		}
	}

	code := m.Run()

	sqlDB, err := db.NewMySQL(cfg.MySQL)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := sqlDB.Exec(`SET FOREIGN_KEY_CHECKS = 0`); err != nil {
		log.Fatal(err)
	}
	if _, err := sqlDB.Exec(`TRUNCATE TABLE recipe`); err != nil {
		log.Fatal(err)
	}
	if _, err := sqlDB.Exec(`TRUNCATE TABLE ingredient`); err != nil {
		log.Fatal(err)
	}
	if _, err := sqlDB.Exec(`SET FOREIGN_KEY_CHECKS = 1`); err != nil {
		log.Fatal(err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

func TestHandler_Recipes(t *testing.T) {
	testData := []struct {
		params  url.Values
		results int
	}{
		{url.Values{"p": []string{"0"}}, 10},
		{url.Values{"p": []string{"1"}}, 10},
		{url.Values{"p": []string{"2"}}, 10},
		{url.Values{"p": []string{"3"}}, 2},
		{url.Values{"p": []string{"1"}, "q": []string{"Ginger Champagne"}}, 1},
		{url.Values{"p": []string{"1"}, "q": []string{"potato"}}, 4},
		{url.Values{"p": []string{"1"}, "q": []string{"onion"}}, 1},
		{url.Values{"p": []string{"1"}, "q": []string{"onion"}, "i": []string{"onions"}}, 1},
		{url.Values{"p": []string{"1"}, "i": []string{"onions"}}, 8},
		{url.Values{"p": []string{"1"}, "i": []string{"eggs"}}, 5},
		{url.Values{"p": []string{"1"}, "i": []string{"onions", "eggs"}}, 10},
		{url.Values{"p": []string{"2"}, "i": []string{"onions", "eggs"}}, 2},
		{url.Values{"p": []string{"1"}, "q": []string{"pork"}}, 3},
		{url.Values{"p": []string{"1"}, "q": []string{"pork"}, "i": []string{"garlic"}}, 2},
		{url.Values{"p": []string{"1"}, "q": []string{"pork"}, "i": []string{"garlic", "brown sugar"}}, 2},
		{url.Values{"p": []string{"1"}, "q": []string{"park"}, "i": []string{"garlic", "brown sugar"}}, 0},
		{url.Values{"p": []string{"1"}, "q": []string{"potato"}, "i": []string{"eggs"}}, 1},
		{url.Values{"p": []string{"1"}, "i": []string{"Spaghetti code"}}, 0},
	}

	cfg, err := config.Load("testdata/config.json")
	if err != nil {
		log.Fatal(err)
	}

	recipeTbl, err := db.New(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	h := handler.NewHandler(recipeTbl, &config.Config{}, &log.Logger{})

	for i := range testData {
		tc := testData[i]

		t.Run(fmt.Sprintf(`Test Case %+v`, tc.params), func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", "/recipes?"+tc.params.Encode(), nil)

			// initialize response recorder to monitor handler response data
			rr := httptest.NewRecorder()
			rh := http.HandlerFunc(h.Recipes)
			rh.ServeHTTP(rr, req)

			if http.StatusOK != rr.Code {
				t.Fatalf("Wrong status code got %d expected %d", http.StatusOK, rr.Code)
			}
			if actualLen := strings.Count(rr.Body.String(), "createdAt"); tc.results != actualLen {
				t.Fatalf("Expected %d results got %d", tc.results, actualLen)
			}
		})
	}
}
