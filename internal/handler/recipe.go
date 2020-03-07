package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/georlav/recipeapi/internal/recipe"
)

func (h Handler) Recipes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Parse get params
	q := r.URL.Query()
	qp := recipe.QueryParams{Term: q.Get("q"), Ingredients: q["i"]}
	page, _ := strconv.Atoi(q.Get("p"))
	qp.Page = int64(page)

	// use recipe repository
	recipes, _, err := h.recipes.GetMany(qp)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), 500)
	}

	resp := RecipeResponse{Title: "Recipe Puppy Clone", Version: 1, Href: "", Results: recipes}
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), 500)
	}
}
