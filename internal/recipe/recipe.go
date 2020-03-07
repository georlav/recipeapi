package recipe

import "go.mongodb.org/mongo-driver/bson/primitive"

type Recipe struct {
	ID          string             `bson:"-"`
	MID         primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	URL         string             `bson:"url"`
	Ingredients []string           `bson:"ingredients"`
	Thumbnail   string             `bson:"thumbnail"`
	UpdatedAt   string             `bson:"-"`
	CreatedAt   string             `bson:"-"`
}

type Recipes []Recipe
