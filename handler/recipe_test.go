package handler_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/georlav/recipeapi/config"
	"github.com/georlav/recipeapi/handler"
)

func TestHandler_Recipes(t *testing.T) {
	// Create request
	req := httptest.NewRequest("GET", "/recipes?p=1", nil)

	// init handlers
	h := handler.NewHandler(&config.Config{}, &log.Logger{})

	// Create a response recorder
	rr := httptest.NewRecorder()
	rh := http.HandlerFunc(h.Recipes)
	rh.ServeHTTP(rr, req)

	if http.StatusOK != rr.Code {
		t.Fatalf("Wrong status code got %d expected %d", http.StatusOK, rr.Code)
	}

	//if actualLen := strings.Count(rr.Body.String(), "id"); tc.results != actualLen {
	//	t.Fatalf("Expected %d results got %d", tc.results, actualLen)
	//}
}
