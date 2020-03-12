package db

import "go.mongodb.org/mongo-driver/bson/primitive"

// Recipe entity
type Recipe struct {
	ID          string             `bson:"-"`
	MID         primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	URL         string             `bson:"url"`
	Thumbnail   string             `bson:"thumbnail"`
	Ingredients Ingredients        `bson:"ingredients"`
	UpdatedAt   string             `bson:"-"`
	CreatedAt   string             `bson:"-"`
}

// Recipes slice or recipe entities
type Recipes []Recipe

// Ingredients
type Ingredients []string
