package handler

import (
	"github.com/gorilla/mux"
)

// Routes initializes api routes and shared middleware
func Routes(h *Handler) *mux.Router {
	r := mux.NewRouter()
	r.Use(
	// Add here your global middleware
	)

	// Product API handlers
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/recipes", h.Recipes).Methods("GET").Name("recipes")
	api.HandleFunc("/recipes/{id}", h.Recipe).Methods("GET").Name("recipe")
	api.HandleFunc("/recipes", h.Create).Methods("POST").Name("create")

	return r
}
