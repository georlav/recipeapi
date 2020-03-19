package handler_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/handler"
	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {
	// Init handlers
	h := handler.NewHandler(nil, config.Config{}, nil)

	// Init routes
	r := handler.Routes(h)

	expectedRoutes := map[string]struct{}{
		"/api/*/recipes/*/":            struct{}{},
		"/api/*/recipes/*/{id:[0-9]+}": struct{}{},
		"/api/*/user/*/":               struct{}{},
		"/api/*/user/*/signin":         struct{}{},
		"/api/*/user/*/signup":         struct{}{},
		"/swagger/*":                   struct{}{},
	}

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if _, ok := expectedRoutes[route]; !ok {
			return fmt.Errorf("route %s is not registered", route)
		}

		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		t.Fatalf("route error: %s", err)
	}
}
