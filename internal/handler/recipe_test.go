package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/handler"
	"github.com/georlav/recipeapi/internal/logger"
)

func TestHandler_Recipes(t *testing.T) {
	// Create request
	req := httptest.NewRequest("GET", "/recipes?p=1", nil)

	// init handlers
	h := handler.NewHandler(&config.Config{}, &logger.Logger{})

	// Create a response recorder
	rr := httptest.NewRecorder()
	rh := http.HandlerFunc(h.Recipes)
	rh.ServeHTTP(rr, req)

	// Status should be 200
	if http.StatusOK != rr.Code {
		t.Fatalf("Wrong status code got %d expected %d", http.StatusOK, rr.Code)
	}

	// Should find 10 times in response the word ingredients
	if actualLen := strings.Count(rr.Body.String(), "ingredients"); actualLen != 10 {
		t.Fatalf("Expected %d results got %d", 10, actualLen)
	}
}
