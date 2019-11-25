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

	// Api handlers
	apiV1Prefix := "v1"
	apiV1 := r.PathPrefix("/api/" + apiV1Prefix).Subrouter()
	apiV1.HandleFunc("/recipes/{id}", nil).Methods("GET").Name("recipe")
	apiV1.HandleFunc("/recipes", h.Recipes).Methods("GET").Name("recipes")

	// Health endpoint
	r.HandleFunc("/health", nil).Methods("GET").Name("health")

	return r
}
