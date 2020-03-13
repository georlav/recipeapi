package database

// Ingredient entity
type Ingredient struct {
	ID        int64
	RecipeID  int64
	Name      string
	CreatedAt string
	UpdatedAt string
}

// Ingredients slice or recipe ingredient entities
type Ingredients []Ingredient
