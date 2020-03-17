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
	user := r.PathPrefix("/api/user").Subrouter()
	user.HandleFunc("/signin", h.SignIn).Methods(http.MethodPost).Name("signin")
	user.HandleFunc("/signup", h.SignUp).Methods(http.MethodPost).Name("signup")

	userProtected := user.Path("/").Subrouter()
	userProtected.Use(h.AuthorizationMiddleware)
	userProtected.HandleFunc("", h.User).Methods(http.MethodGet).Name("profile")

	// health endpoint
	r.HandleFunc("/health", nil).Methods("GET").Name("health")

	return r
}
