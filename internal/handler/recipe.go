package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/georlav/recipeapi/internal/database"
	"github.com/go-chi/chi"
)

func (h *Handler) Recipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	id := chi.URLParam(r, "id")
	nID, err := strconv.Atoi(id)
	if err != nil || id == "" {
		h.respondError(w, APIError{Message: "recipe id is required.", StatusCode: http.StatusBadRequest})
		return
	}

	recipe, err := h.db.Recipe.Get(uint64(nID))
	if err != nil {
		h.respondError(w, APIError{Message: "unknown recipe", StatusCode: http.StatusNotFound})
		return
	}

	// Respond
	h.respond(w, NewRecipeResponse(recipe), http.StatusOK)
}

func (h *Handler) Recipes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Map request to struct
	rr := RecipesRequest{Page: 1}
	if err := h.schema.Decode(&rr, r.URL.Query()); err != nil {
		h.respondError(w, APIError{Message: http.StatusText(http.StatusBadRequest), StatusCode: http.StatusBadRequest})
		return
	}

	// validate data in struct
	if err := h.validate.Struct(rr); err != nil {
		h.respondError(w, APIError{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	// Create db filters from validated request data
	filters := database.RecipeFilters{
		Term:        rr.Term,
		Ingredients: rr.Ingredients,
	}

	// retrieve data from database
	recipes, total, err := h.db.Recipe.Paginate(rr.Page, &filters)
	if err != nil {
		h.respondError(w, err)
		return
	}

	// Respond
	resp := NewRecipesResponse("Recipe Puppy Clone", h.cfg.APP.Version, recipes, total)
	h.respond(w, resp, http.StatusOK)
}

// Create a new recipe
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Map request to struct
	rc := RecipeCreateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&rc); err != nil {
		h.respondError(w, APIError{Message: http.StatusText(http.StatusBadRequest), StatusCode: http.StatusBadRequest})
		return
	}

	// validate data in struct
	if err := h.validate.Struct(rc); err != nil {
		h.respondError(w, APIError{Message: http.StatusText(http.StatusBadRequest), StatusCode: http.StatusBadRequest})
		return
	}

	// Create a slice of ingredients
	ingredients := func() (ing database.Ingredients) {
		for i := range rc.Ingredients {
			ing = append(ing, database.Ingredient{Name: rc.Ingredients[i]})
		}

		return ing
	}()

	// Insert new recipe
	if _, err := h.db.Recipe.Insert(database.Recipe{
		Title:       rc.Title,
		URL:         rc.URL,
		Thumbnail:   rc.Thumbnail,
		Ingredients: ingredients,
	}); err != nil {
		h.respondError(w, APIError{Message: "failed to create recipe", StatusCode: http.StatusInternalServerError})
		return
	}

	w.WriteHeader(http.StatusCreated)
}
