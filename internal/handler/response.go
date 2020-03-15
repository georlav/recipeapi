package handler

import "github.com/georlav/recipeapi/internal/database"

// RecipeResponse recipe response object
type RecipesResponse struct {
	Title    string               `json:"title"`
	Version  int                  `json:"version"`
	Href     string               `json:"href"`
	Data     *RecipeResponseItems `json:"data"`
	Metadata Metadata             `json:"metadata"`
}

// Metadata
type Metadata struct {
	Total int64
}

// NewRecipesResponse
func NewRecipesResponse(title string, version int, r database.Recipes, total int64) RecipesResponse {
	rr := RecipesResponse{
		Title:   title,
		Version: version,
		Metadata: Metadata{
			Total: total,
		},
	}

	var items RecipeResponseItems
	for i := range r {
		item := RecipeResponseItem{
			ID:        r[i].ID,
			Title:     r[i].Title,
			Thumbnail: r[i].Thumbnail,
			CreatedAt: r[i].CreatedAt,
			UpdatedAt: r[i].UpdatedAt,
		}
		for j := range r[i].Ingredients {
			item.Ingredients = append(item.Ingredients, IngredientResponseItem{
				ID:   r[i].Ingredients[j].ID,
				Name: r[i].Ingredients[j].Name,
			})
		}
		items = append(items, item)
	}
	rr.Data = &items

	return rr
}

// RecipeResponseItem object to map recipe items
type RecipeResponseItems []RecipeResponseItem

// NewRecipesResponse
func NewRecipeResponse(r *database.Recipe) RecipeResponseItem {
	ingredients := IngredientResponse{}
	for i := range r.Ingredients {
		ingredients = append(ingredients, IngredientResponseItem{
			ID:   r.Ingredients[i].ID,
			Name: r.Ingredients[i].Name,
		})
	}

	return RecipeResponseItem{
		ID:          r.ID,
		Title:       r.Title,
		Ingredients: ingredients,
		Thumbnail:   r.Thumbnail,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// RecipeResponseItem object to map a recipe item
type RecipeResponseItem struct {
	ID          int64              `json:"id"`
	Title       string             `json:"title"`
	Href        string             `json:"href"`
	Ingredients IngredientResponse `json:"ingredients"`
	Thumbnail   string             `json:"thumbnail"`
	CreatedAt   string             `json:"createdAt"`
	UpdatedAt   string             `json:"updatedAt"`
}

type IngredientResponseItem struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type IngredientResponse []IngredientResponseItem
