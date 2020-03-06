package handler

import (
	"log"

	"github.com/georlav/recipeapi/internal/config"
)

type Handler struct {
	cfg *config.Config
	log *log.Logger
}

func NewHandler(c *config.Config, l *log.Logger) *Handler {
	return &Handler{
		cfg: c,
		log: l,
	}
}
