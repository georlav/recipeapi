package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/georlav/recipeapi/internal/db"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

func (h Handler) Recipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	id, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, "recipe id is required."), http.StatusBadRequest)
		return
	}

	recipe, err := h.recipes.Get(id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, "unknown recipe."), http.StatusNotFound)
		return
	}

	// Respond
	resp := NewRecipeResponse(recipe)
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusInternalServerError)
	}
}

func (h Handler) Recipes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Map request to struct
	rr := RecipesRequest{}
	if err := h.decoder.Decode(&rr, r.URL.Query()); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusBadRequest)
		return
	}

	// Pass request data to filters
	filters := db.Filters{
		Term:        rr.Term,
		Ingredients: rr.Ingredients,
	}

	// retrieve data from database
	recipes, total, err := h.recipes.Paginate(rr.Page, &filters)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusInternalServerError)
		return
	}

	// Respond
	resp := NewRecipesResponse("Recipe Puppy Clone", h.cfg.APP.Version, recipes, total)
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusInternalServerError)
	}
}

// Create a new recipe
func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Map request to struct
	rc := RecipeCreateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&rc); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusBadRequest)
		return
	}

	// validate data in struct
	v := validator.New()
	if err := v.Struct(rc); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusBadRequest)
		return
	}

	// Insert new recipe
	err := h.recipes.Insert(db.Recipe{
		Title:       rc.Title,
		URL:         rc.URL,
		Thumbnail:   rc.Thumbnail,
		Ingredients: rc.Ingredients,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, "failed to create recipe"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
