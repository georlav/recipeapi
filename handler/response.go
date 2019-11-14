package handler

import "github.com/georlav/recipeapi/recipe"

type RecipeResponse struct {
	Title   string         `json:"title"`
	Version float64        `json:"version"`
	Href    string         `json:"href"`
	Results recipe.Recipes `json:"results"`
}
