package handler

import (
	"github.com/gorilla/mux"
)

// Routes initializes api routes and shared middleware
func Routes(h *Handler) *mux.Router {
	r := mux.NewRouter()
	r.Use(
		h.headersMiddleware,
	)

	// Product API handlers
	api := r.PathPrefix("/api/recipes").Subrouter()
	api.HandleFunc("/{id:[0-9]+}", h.Recipe).Methods("GET").Name("recipe")
	api.HandleFunc("", h.Recipes).Methods("GET").Name("recipes")
	api.HandleFunc("", h.Recipes).Methods("POST").Name("create")

	return r
}
