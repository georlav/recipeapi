package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/database"
	"github.com/georlav/recipeapi/internal/handler"
	"github.com/georlav/recipeapi/internal/logger"
	"github.com/go-chi/chi"
)

func TestHandler_Recipe(t *testing.T) {
	testData := []struct {
		input        uint64
		resultTitle  string
		expectedCode int
	}{
		{1, "Ginger Champagne", http.StatusOK},
		{2, "Potato and Cheese Frittata", http.StatusOK},
		{9999, "souvlaki", http.StatusNotFound},
		{0, "", http.StatusBadRequest},
	}

	cfg, err := config.New("config", "testdata")
	if err != nil {
		t.Fatal(err)
	}

	db, err := database.New(cfg.Database)
	if err != nil {
		t.Fatal(err)
	}

	h := handler.NewHandler(db, cfg, logger.NewLogger(cfg.Logger))

	for i := range testData {
		tc := testData[i]

		t.Run(fmt.Sprintf(`Get recipe with id %d`, tc.input), func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("GET", fmt.Sprintf("/recipes/%d", tc.input), nil)

			// Inject uri param
			ctx := chi.NewRouteContext()
			if tc.input != 0 {
				ctx.URLParams.Add("id", fmt.Sprintf(`%d`, tc.input))
			}

			// initialize response recorder to monitor handler response data
			rr := httptest.NewRecorder()
			rh := http.HandlerFunc(h.Recipe)
			rh.ServeHTTP(rr, req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx)))

			if rr.Code != tc.expectedCode {
				t.Fatalf(
					"Wrong status code expected %d got %d, Body: (%s)",
					tc.expectedCode,
					rr.Code,
					rr.Body.String(),
				)
			}

			if http.StatusOK == rr.Code {
				if !strings.Contains(rr.Body.String(), tc.resultTitle) {
					t.Fatalf("Expected to have result with title %s", tc.resultTitle)
				}

				respData := handler.RecipeResponseItem{}
				if err := json.Unmarshal(rr.Body.Bytes(), &respData); err != nil {
					t.Fatal(err)
				}

				if len(respData.Ingredients) == 0 {
					t.Fatal("Expected to have at least one ingredient")
				}
			}
		})
	}
}

func TestHandler_Recipes(t *testing.T) {
	testData := []struct {
		params     url.Values
		results    int
		statusCode int
	}{
		{url.Values{"page": []string{"0"}}, 10, http.StatusOK},
		{url.Values{"page": []string{"1"}}, 10, http.StatusOK},
		{url.Values{"page": []string{"2"}}, 10, http.StatusOK},
		{url.Values{"page": []string{"3"}}, 2, http.StatusOK},
		{url.Values{"page": []string{"1"}, "term": []string{"Ginger Champagne"}}, 1, http.StatusOK},
		{url.Values{"page": []string{"1"}, "term": []string{"potato"}}, 4, http.StatusOK},
		{url.Values{"page": []string{"1"}, "term": []string{"onion"}}, 1, http.StatusOK},
		{url.Values{"page": []string{"1"}, "term": []string{"onion"}, "ingredient": []string{"onions"}}, 1, http.StatusOK},
		{url.Values{"page": []string{"1"}, "ingredient": []string{"onions"}}, 8, http.StatusOK},
		{url.Values{"page": []string{"1"}, "ingredient": []string{"eggs"}}, 5, http.StatusOK},
		{url.Values{"page": []string{"1"}, "ingredient": []string{"onions", "eggs"}}, 10, http.StatusOK},
		{url.Values{"page": []string{"2"}, "ingredient": []string{"onions", "eggs"}}, 2, http.StatusOK},
		{url.Values{"page": []string{"1"}, "term": []string{"pork"}}, 3, http.StatusOK},
		{url.Values{"page": []string{"1"}, "term": []string{"pork"}, "ingredient": []string{"garlic"}}, 2, http.StatusOK},
		{url.Values{"page": []string{"1"}, "term": []string{"pork"}, "ingredient": []string{"garlic", "brown sugar"}}, 2, http.StatusOK},
		{url.Values{"page": []string{"1"}, "term": []string{"park"}, "ingredient": []string{"garlic", "brown sugar"}}, 0, http.StatusOK},
		{url.Values{"page": []string{"1"}, "term": []string{"potato"}, "ingredient": []string{"eggs"}}, 1, http.StatusOK},
		{url.Values{"page": []string{"1"}, "ingredient": []string{"Spaghetti code"}}, 0, http.StatusOK},
		{url.Values{"page": []string{"1"}, "ingredient": []string{"1", "2", "3", "4", "5", "6"}}, 0, http.StatusBadRequest},
		{url.Values{"term": []string{"ab"}}, 0, http.StatusBadRequest},
		{url.Values{"page": []string{"-5"}}, 0, http.StatusBadRequest},
	}

	cfg, err := config.New("config", "testdata")
	if err != nil {
		t.Fatal(err)
	}

	db, err := database.New(cfg.Database)
	if err != nil {
		t.Fatal(err)
	}

	h := handler.NewHandler(db, cfg, logger.NewLogger(cfg.Logger))

	for i := range testData {
		tc := testData[i]

		t.Run(fmt.Sprintf(`Test Case %+v`, tc.params), func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", "/recipes?"+tc.params.Encode(), nil)

			// initialize response recorder to monitor handler response data
			rr := httptest.NewRecorder()
			rh := http.HandlerFunc(h.Recipes)
			rh.ServeHTTP(rr, req)

			if rr.Code != tc.statusCode {
				t.Fatalf("Wrong status code got %d expected %d, %s", rr.Code, tc.statusCode, rr.Body.String())
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
			`{"title":"Ginger Champagne2","url":"http://allrecipes.com/Recipe/Ginger-Champagne/Detail.aspx",
"ingredients":["champagne","ginger","ice","vodka"],"thumbnail":"http://img.recipepuppy.com/1.jpg"}`,
			http.StatusCreated,
			"",
		},
		{
			`{"title":"Ginger Champagne2","url":"http://allrecipes.com/Recipe/Ginger-Champagne/Detail.aspx",
"ingredients":["champagne","ginger","ice","vodka"],"thumbnail":"http://img.recipepuppy.com/1.jpg"}`,
			http.StatusInternalServerError,
			"failed to create recipe",
		},
		{
			`{"url":"http://allrecipes.com/Recipe/Ginger-Champagne/Detail.aspx",
"ingredients":["champagne","ginger","ice","vodka"],"thumbnail":"http://img.recipepuppy.com/1.jpg"}`,
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
			`{"title":"t", "URL":"http://allrecipes.com/Recipe/Ginger-Champagne/Detail.aspx",
"ingredients":[],"thumbnail":"http://img.recipepuppy.com/1.jpg"}`,
			http.StatusBadRequest,
			`Field validation for 'Title' failed on the 'min' tag`,
		},
		{
			`{"title":"Ginger Champagne2", "url":"http://allrecipes.com/Recipe/Ginger-Champagne/Detail.aspx",
"thumbnail":"http://img.recipepuppy.com/1.jpg"}`,
			http.StatusBadRequest,
			`Field validation for 'Ingredients' failed on the 'required' tag"`,
		},
	}

	cfg, err := config.New("config", "testdata")
	if err != nil {
		t.Fatal(err)
	}

	db, err := database.New(cfg.Database)
	if err != nil {
		t.Fatal(err)
	}

	h := handler.NewHandler(db, cfg, logger.NewLogger(cfg.Logger))

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
