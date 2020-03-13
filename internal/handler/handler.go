package handler

import (
	"log"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/database"
	"github.com/gorilla/schema"
)

type Handler struct {
	db      *database.Database
	decoder *schema.Decoder
	cfg     config.Config
	log     *log.Logger
}

func NewHandler(db *database.Database, c config.Config, l *log.Logger) *Handler {
	return &Handler{
		db:      db,
		decoder: schema.NewDecoder(),
		cfg:     c,
		log:     l,
	}
}
