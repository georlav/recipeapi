package handler

import (
	"log"
	"net/http"

	"github.com/georlav/recipeapi/config"
)

type Handler struct {
	cfg config.Config
	log *log.Logger
}

func (t Handler) Recipe(w http.ResponseWriter, r *http.Request)  {}
func (t Handler) Recipes(w http.ResponseWriter, r *http.Request) {}
