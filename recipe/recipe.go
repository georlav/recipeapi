package recipe

import "go.mongodb.org/mongo-driver/bson/primitive"

type Recipe struct {
	ID          string             `json:"id" bson:"-"`
	MID         primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	URL         string             `json:"url" bson:"url"`
	Ingredients []string           `json:"ingredients" bson:"ingredients"`
	Thumbnail   string             `json:"thumbnail" bson:"thumbnail"`
	UpdatedAt   string             `json:"updatedAt" bson:"-"`
	CreatedAt   string             `json:"createdAt" bson:"-"`
}

type Recipes []Recipe
