package handler

import (
	"fmt"
	"net/http"
	"strings"

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

	// health endpoint
	r.HandleFunc("/health", nil).Methods("GET").Name("health")

	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}
		fmt.Println()
		return nil
	})

	return r
}
