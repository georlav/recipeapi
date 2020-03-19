package handler

import (
	"time"

	_ "github.com/georlav/recipeapi/api/swagger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Routes initializes api routes and shared middleware
func Routes(h *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(
		middleware.RequestID,
		middleware.Timeout(60*time.Second),
		h.CorsMiddleware,
		h.ContentTypeMiddleware,
	)

	// Recipe routes
	r.Route("/recipes", func(r chi.Router) {
		r.Use(h.AuthorizationMiddleware)
		r.Get("/{id:[0-9]+}", h.Recipe)
		r.Get("/", h.Recipes)
		r.Post("/", h.Create)
	})

	// User routes
	r.Route("/user", func(r chi.Router) {
		// Public
		r.Post("/signin", h.SignIn)
		r.Post("/signup", h.SignUp)

		// Need authentication
		r.With(h.AuthorizationMiddleware).Get("/", h.User)
	})

	// Swagger Docs
	rs := chi.NewRouter()
	rs.Handle("/swagger/*", httpSwagger.Handler())

	// Mount recipes and users under /api and merge with swagger routes
	rs.Mount("/api", r)

	return rs
}
