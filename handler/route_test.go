package handler_test

import (
	"log"
	"testing"

	"github.com/georlav/recipeapi/config"
	"github.com/georlav/recipeapi/handler"
)

func TestRoutes(t *testing.T) {
	// Init handlers
	h := handler.NewHandler(&config.Config{}, &log.Logger{})

	// Init routes
	r := handler.Routes(h)

	// Check if recipes route exist
	rr := r.GetRoute("recipes")
	if rr == nil {
		t.Fatal("recipes route is missing")
	}
}
