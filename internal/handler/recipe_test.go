package handler_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/handler"
)

func TestHandler_Recipes(t *testing.T) {
	// Create request
	req := httptest.NewRequest("GET", "/recipes?p=1", nil)

	// init handlers
	h := handler.NewHandler(nil, &config.Config{}, &log.Logger{})

	// Create a response recorder
	rec := httptest.NewRecorder()
	rh := http.HandlerFunc(h.Recipes)
	rh.ServeHTTP(rec, req)

	// Status should be 200
	if http.StatusOK != rec.Code {
		t.Fatalf("Wrong status code got %d expected %d", http.StatusOK, rec.Code)
	}

	// Should find 10 times in response the word ingredients
	if actualLen := strings.Count(rec.Body.String(), "ingredients"); actualLen != 10 {
		t.Fatalf("Expected %d results got %d", 10, actualLen)
	}
}
