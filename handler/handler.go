package handler

import (
	"log"

	"github.com/georlav/recipeapi/recipe"

	"github.com/georlav/recipeapi/config"
)

type Handler struct {
	recipes recipe.Repository
	cfg     *config.Config
	log     *log.Logger
}

func NewHandler(r recipe.Repository, c *config.Config, l *log.Logger) *Handler {
	return &Handler{
		recipes: r,
		cfg:     c,
		log:     l,
	}
}
