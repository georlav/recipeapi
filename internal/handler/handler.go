package handler

import (
	"log"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/db"
)

type Handler struct {
	recipes db.Queryable
	cfg     *config.Config
	log     *log.Logger
}

func NewHandler(r db.Queryable, c *config.Config, l *log.Logger) *Handler {
	return &Handler{
		recipes: r,
		cfg:     c,
		log:     l,
	}
}
