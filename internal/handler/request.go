package handler

// RecipesRequest object to map incoming request for Recipes handler
type RecipesRequest struct {
	Page        uint64   `schema:"page"`
	Term        string   `schema:"term"`
	Ingredients []string `schema:"ingredient"`
}

// CreateRecipeRequest object to map incoming request for Create handler
type RecipeCreateRequest struct {
	Title       string   `validate:"required,min=2" schema:"title"`
	URL         string   `validate:"required,min=10" schema:"url"`
	Thumbnail   string   `schema:"thumbnail"`
	Ingredients []string `validate:"required,max=30,min=1" schema:"ingredients"`
}
