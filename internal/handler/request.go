package handler

// RecipesRequest object to map incoming request for Recipes handler
type RecipesRequest struct {
	Page        uint64   `schema:"page" validate:"omitempty,min=1"`
	Term        string   `schema:"term" validate:"omitempty,min=3"`
	Ingredients []string `schema:"ingredient" validate:"omitempty,max=5"`
}

// CreateRecipeRequest object to map incoming request for Create handler
type RecipeCreateRequest struct {
	Title       string   `json:"title" validate:"required,min=2"`
	URL         string   `json:"url" validate:"required,min=10"`
	Thumbnail   string   `json:"thumbnail"`
	Ingredients []string `json:"ingredients" validate:"required,max=30,min=1"`
}
