package handler

import (
	"log"

	"github.com/gorilla/schema"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/db"
)

type Handler struct {
	recipes db.Queryable
	decoder *schema.Decoder
	cfg     *config.Config
	log     *log.Logger
}

func NewHandler(r db.Queryable, c *config.Config, l *log.Logger) *Handler {
	return &Handler{
		recipes: r,
		decoder: schema.NewDecoder(),
		cfg:     c,
		log:     l,
	}
}
