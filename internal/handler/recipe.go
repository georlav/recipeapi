package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/georlav/recipeapi/internal/database"
	"gopkg.in/go-playground/validator.v9"
)

// Recipe godoc
// @Summary Get a recipe
// @Description Get a recipe by ID
// @ID get-recipe-by-int
// @Accept  application/x-www-form-urlencoded
// @Produce  json
// @Param id path int true "Recipe ID"
// @Success 200 {object} handler.RecipeResponseItem
// @Failure 400 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Security ApiKeyAuth
// @Router /recipes/{id} [get]
func (h Handler) Recipe(w http.ResponseWriter, r *http.Request) {
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

// Recipes godoc
// @Summary Get recipes
// @Description Get a list of recipes
// @ID get-recipes
// @Accept  application/x-www-form-urlencoded
// @Produce  json
// @Success 200 {object} handler.RecipesResponse
// @Failure 400 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Security ApiKeyAuth
// @Router /recipes [get]
func (h Handler) Recipes(w http.ResponseWriter, r *http.Request) {
	// Map request to struct
	rr := RecipesRequest{}
	if err := h.decoder.Decode(&rr, r.URL.Query()); err != nil {
		h.respondError(w, APIError{Message: http.StatusText(http.StatusBadRequest), StatusCode: http.StatusBadRequest})
		return
	}

	// Pass request data to filters
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
func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	// Map request to struct
	rc := RecipeCreateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&rc); err != nil {
		h.respondError(w, APIError{Message: http.StatusText(http.StatusBadRequest), StatusCode: http.StatusBadRequest})
		return
	}

	// validate data in struct
	v := validator.New()
	if err := v.Struct(rc); err != nil {
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
