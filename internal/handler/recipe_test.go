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

	"github.com/gorilla/mux"

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

func TestHandler_Recipe(t *testing.T) {
	testData := []struct {
		input        uint64
		resultTitle  string
		expectedCode int
	}{
		{1, "Ginger Champagne", http.StatusOK},
		{2, "Potato and Cheese Frittata", http.StatusOK},
		{9999, "souvlaki", http.StatusNotFound},
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

		t.Run(fmt.Sprintf(`Get recipe with id %d`, tc.input), func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", fmt.Sprintf("/recipes/%d", tc.input), nil)

			muxReq := mux.SetURLVars(req, map[string]string{
				"id": fmt.Sprintf("%d", tc.input),
			})

			// initialize response recorder to monitor handler response data
			rr := httptest.NewRecorder()
			rh := http.HandlerFunc(h.Recipe)
			rh.ServeHTTP(rr, muxReq)

			if rr.Code != tc.expectedCode {
				t.Fatalf("Wrong status code got %d expected %d, %s", http.StatusOK, rr.Code, rr.Body.String())
			}
			if http.StatusOK == tc.expectedCode && !strings.Contains(rr.Body.String(), tc.resultTitle) {
				t.Fatalf("Expected to have result with title %s", tc.resultTitle)
			}
		})
	}
}

func TestHandler_Recipes(t *testing.T) {
	testData := []struct {
		params  url.Values
		results int
	}{
		{url.Values{"page": []string{"0"}}, 10},
		{url.Values{"page": []string{"1"}}, 10},
		{url.Values{"page": []string{"2"}}, 10},
		{url.Values{"page": []string{"3"}}, 2},
		{url.Values{"page": []string{"1"}, "term": []string{"Ginger Champagne"}}, 1},
		{url.Values{"page": []string{"1"}, "term": []string{"potato"}}, 4},
		{url.Values{"page": []string{"1"}, "term": []string{"onion"}}, 1},
		{url.Values{"page": []string{"1"}, "term": []string{"onion"}, "ingredient": []string{"onions"}}, 1},
		{url.Values{"page": []string{"1"}, "ingredient": []string{"onions"}}, 8},
		{url.Values{"page": []string{"1"}, "ingredient": []string{"eggs"}}, 5},
		{url.Values{"page": []string{"1"}, "ingredient": []string{"onions", "eggs"}}, 10},
		{url.Values{"page": []string{"2"}, "ingredient": []string{"onions", "eggs"}}, 2},
		{url.Values{"page": []string{"1"}, "term": []string{"pork"}}, 3},
		{url.Values{"page": []string{"1"}, "term": []string{"pork"}, "ingredient": []string{"garlic"}}, 2},
		{url.Values{"page": []string{"1"}, "term": []string{"pork"}, "ingredient": []string{"garlic", "brown sugar"}}, 2},
		{url.Values{"page": []string{"1"}, "term": []string{"park"}, "ingredient": []string{"garlic", "brown sugar"}}, 0},
		{url.Values{"page": []string{"1"}, "term": []string{"potato"}, "ingredient": []string{"eggs"}}, 1},
		{url.Values{"page": []string{"1"}, "ingredient": []string{"Spaghetti code"}}, 0},
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
				t.Fatalf("Wrong status code got %d expected %d, %s", http.StatusOK, rr.Code, rr.Body.String())
			}
			if actualLen := strings.Count(rr.Body.String(), "createdAt"); tc.results != actualLen {
				t.Fatalf("Expected %d results got %d", tc.results, actualLen)
			}
		})
	}
}

func TestHandler_Create(t *testing.T) {
	testData := []struct {
		payload       string
		expectedCode  int
		expectedError string
	}{
		{
			`{"Title":"Ginger Champagne2","URL":"http://allrecipes.com/Recipe/Ginger-Champagne/Detail.aspx",
"Ingredients":["champagne","ginger","ice","vodka"],"Thumbnail":"http://img.recipepuppy.com/1.jpg"}`,
			http.StatusCreated,
			"",
		},
		{
			`{"Title":"Ginger Champagne2","URL":"http://allrecipes.com/Recipe/Ginger-Champagne/Detail.aspx",
"Ingredients":["champagne","ginger","ice","vodka"],"Thumbnail":"http://img.recipepuppy.com/1.jpg"}`,
			http.StatusInternalServerError,
			"failed to create recipe",
		},
		{
			`{"URL":"http://allrecipes.com/Recipe/Ginger-Champagne/Detail.aspx",
"Ingredients":["champagne","ginger","ice","vodka"],"Thumbnail":"http://img.recipepuppy.com/1.jpg"}`,
			http.StatusBadRequest,
			`Field validation for 'Title' failed on the 'required' tag`,
		},
		{
			`{"Title":"Ginger Champagne2", "URL":"http://allrecipes.com/Recipe/Ginger-Champagne/Detail.aspx",
"Ingredients":[],"Thumbnail":"http://img.recipepuppy.com/1.jpg"}`,
			http.StatusBadRequest,
			`Field validation for 'Ingredients' failed on the 'min' tag`,
		},
		{
			`{"Title":"t", "URL":"http://allrecipes.com/Recipe/Ginger-Champagne/Detail.aspx",
"Ingredients":[],"Thumbnail":"http://img.recipepuppy.com/1.jpg"}`,
			http.StatusBadRequest,
			`Field validation for 'Title' failed on the 'min' tag`,
		},
		{
			`{"Title":"Ginger Champagne2", "URL":"http://allrecipes.com/Recipe/Ginger-Champagne/Detail.aspx",
"Thumbnail":"http://img.recipepuppy.com/1.jpg"}`,
			http.StatusBadRequest,
			`Field validation for 'Ingredients' failed on the 'required' tag"`,
		},
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

		t.Run(`Sending payload`, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/recipes", strings.NewReader(tc.payload))

			// initialize response recorder to monitor handler response data
			rr := httptest.NewRecorder()
			rh := http.HandlerFunc(h.Create)
			rh.ServeHTTP(rr, req)

			if rr.Code != tc.expectedCode {
				t.Fatalf("Wrong status code got %d expected %d, %s", rr.Code, tc.expectedCode, rr.Body.String())
			}
		})
	}
}
