package handler

import (
	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/logger"
)

type Handler struct {
	cfg *config.Config
	log *logger.Logger
}

func NewHandler(c *config.Config, l *logger.Logger) *Handler {
	return &Handler{
		cfg: c,
		log: l,
	}
}
