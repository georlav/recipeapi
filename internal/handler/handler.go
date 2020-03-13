package handler

import (
	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/database"
	"github.com/georlav/recipeapi/internal/logger"
)

type Handler struct {
	db  *database.Database
	cfg *config.Config
	log *logger.Logger
}

func NewHandler(db *database.Database, c *config.Config, l *logger.Logger) *Handler {
	return &Handler{
		db:  db,
		cfg: c,
		log: l,
	}
}
