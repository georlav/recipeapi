package handler_test

import (
	"testing"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/handler"
	"github.com/georlav/recipeapi/internal/logger"
)

func TestRoutes(t *testing.T) {
	// Init handlers
	h := handler.NewHandler(&config.Config{}, &logger.Logger{})

	// Init routes
	r := handler.Routes(h)

	// Check if recipes route exist
	rr := r.GetRoute("recipes")
	if rr == nil {
		t.Fatal("recipes route is missing")
	}
}
