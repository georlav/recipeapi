package handler

import (
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Routes initializes api routes and shared middleware
func Routes(h *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(
		middleware.RequestID,
		middleware.Timeout(60*time.Second),
	)

	// Recipe routes
	r.Route("/api/recipes", func(r chi.Router) {
		r.Get("/{id:[0-9]+}", h.Recipe)
		r.Get("/", h.Recipes)
		r.Post("/", nil)
	})

	return r
}
