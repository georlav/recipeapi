package db

import "go.mongodb.org/mongo-driver/bson/primitive"

// Recipe entity
type Recipe struct {
	ID          string             `json:"-" bson:"-"`
	MID         primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	URL         string             `json:"url" bson:"url"`
	Thumbnail   string             `json:"thumbnail" bson:"thumbnail"`
	Ingredients Ingredients        `json:"ingredients" bson:"ingredients"`
	UpdatedAt   string             `json:"-" bson:"-"`
	CreatedAt   string             `json:"-" bson:"-"`
}

// Recipes slice or recipe entities
type Recipes []Recipe

// Ingredients
type Ingredients []string
