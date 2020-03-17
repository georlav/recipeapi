package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Routes initializes api routes and shared middleware
func Routes(h *Handler) *mux.Router {
	r := mux.NewRouter()
	r.Use(h.HeadersMiddleware)

	// recipe routes
	articleRoutes := r.PathPrefix("/api/recipes").Subrouter()
	articleRoutes.Use(h.AuthorizationMiddleware)
	articleRoutes.HandleFunc("/{id:[0-9]+}", h.Recipe).Methods(http.MethodGet).Name("recipe")
	articleRoutes.HandleFunc("", h.Recipes).Methods(http.MethodGet).Name("recipes")
	articleRoutes.HandleFunc("", h.Create).Methods(http.MethodPost).Name("create")

	// user routes
	r.HandleFunc("/api/user/signin", h.SignIn).Methods(http.MethodPost).Name("signin")
	r.HandleFunc("/api/user/signup", h.SignUp).Methods(http.MethodPost).Name("signup")
	user := r.Path("/api/user").Subrouter()
	user.Use(h.AuthorizationMiddleware)
	user.HandleFunc("", h.User).Methods(http.MethodGet).Name("profile")

	return r
}
