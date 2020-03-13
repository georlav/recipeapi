package handler

import (
	"database/sql"
	"log"

	"github.com/georlav/recipeapi/internal/config"
)

type Handler struct {
	db  *sql.DB
	cfg *config.Config
	log *log.Logger
}

func NewHandler(db *sql.DB, c *config.Config, l *log.Logger) *Handler {
	return &Handler{
		db:  db,
		cfg: c,
		log: l,
	}
}
