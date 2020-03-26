package handler

import (
	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/database"
	"github.com/georlav/recipeapi/internal/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

type Handler struct {
	db       *database.Database
	cfg      *config.Config
	log      *logger.Logger
	schema   *schema.Decoder
	validate *validator.Validate
}

func NewHandler(db *database.Database, c *config.Config, l *logger.Logger) *Handler {
	return &Handler{
		db:       db,
		cfg:      c,
		log:      l,
		schema:   schema.NewDecoder(),
		validate: validator.New(),
	}
}
