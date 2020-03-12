package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/georlav/recipeapi/internal/db"
)

func (h Handler) Recipes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Parse query params
	q := r.URL.Query()
	qp := db.Filters{Term: q.Get("q"), Ingredients: q["i"]}
	page, _ := strconv.Atoi(q.Get("p"))

	// retrieve data from database
	recipes, _, err := h.recipes.Paginate(uint64(page), &qp)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), 500)
	}

	// Respond
	resp := NewRecipesResponse("Recipe Puppy Clone", h.cfg.APP.Version, recipes)
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), 500)
	}
}
