package database

// Recipe entity
type Recipe struct {
	ID          int64
	Title       string
	URL         string
	Thumbnail   string
	Ingredients Ingredients
	CreatedAt   string
	UpdatedAt   string
}

// Recipes slice or recipe entities
type Recipes []Recipe
