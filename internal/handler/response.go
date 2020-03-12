package handler

import "github.com/georlav/recipeapi/internal/db"

// RecipeResponse recipe response object
type RecipesResponse struct {
	Title    string               `json:"title"`
	Version  int                  `json:"version"`
	Href     string               `json:"href"`
	Data     *RecipeResponseItems `json:"results"`
	Metadata Metadata             `json:"metadata"`
}

type Metadata struct {
	Total int64
}

// NewRecipesResponse
func NewRecipesResponse(title string, version int, r db.Recipes, total int64) RecipesResponse {
	rr := RecipesResponse{
		Title:   title,
		Version: version,
		Metadata: Metadata{
			Total: total,
		},
	}

	var items RecipeResponseItems
	for i := range r {
		items = append(items, RecipeResponseItem{
			ID:          r[i].ID,
			Title:       r[i].Title,
			Ingredients: r[i].Ingredients,
			Thumbnail:   r[i].Thumbnail,
			CreatedAt:   r[i].CreatedAt,
			UpdatedAt:   r[i].UpdatedAt,
		})
	}
	rr.Data = &items

	return rr
}

// RecipeResponseItem object to map recipe items
type RecipeResponseItems []RecipeResponseItem

// NewRecipesResponse
func NewRecipeResponse(r *db.Recipe) RecipeResponseItem {
	return RecipeResponseItem{
		ID:          r.ID,
		Title:       r.Title,
		Ingredients: nil,
		Thumbnail:   r.Thumbnail,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// RecipeResponseItem object to map a recipe item
type RecipeResponseItem struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Href        string   `json:"href"`
	Ingredients []string `json:"ingredients"`
	Thumbnail   string   `json:"thumbnail"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}
