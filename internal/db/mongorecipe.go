package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Recipe object
type RecipeCollection struct {
	collection *mongo.Collection
	pageSize   uint64
}

// RecipeCollection creates a recipe object
func NewRecipeCollection(recipe *mongo.Collection) *RecipeCollection {
	return &RecipeCollection{
		collection: recipe,
		pageSize:   10,
	}
}

// Get a single recipe by id
func (r *RecipeCollection) Get(id string) (*Recipe, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid signal id (%s), %w", id, err)
	}

	result := r.collection.FindOne(
		context.Background(),
		bson.D{{Key: "_id", Value: oid}},
	)

	var rcp Recipe
	if err := result.Decode(&rcp); err != nil {
		return nil, err
	}

	return &rcp, nil
}

// Paginate get paginated recipes
func (r *RecipeCollection) Paginate(page uint64, flt *Filters) (Recipes, int64, error) {
	// find options
	fOpts := options.Find()
	fOpts.SetLimit(int64(r.pageSize))
	if page > 0 {
		page--
	}
	fOpts.SetSkip(int64(page * r.pageSize))

	// if there are filters convert them to mongo filters
	filters := bson.D{}
	if flt != nil {
		filters = r.paramsToFilters(*flt)
	}

	// Do query
	cursor, err := r.collection.Find(
		context.Background(),
		filters,
		fOpts,
	)
	if err != nil {
		return nil, 0, err
	}
	// nolint[:errcheck, gosec]
	defer cursor.Close(context.TODO())

	// Iterate results
	var recipes Recipes
	for cursor.Next(context.TODO()) {
		recipe := Recipe{}
		err = cursor.Decode(&recipe)
		if err != nil {
			return nil, 0, err
		}

		recipe.ID = recipe.MID.Hex()
		recipe.CreatedAt = cursor.Current.Lookup("createdAt").Time().Format(time.RFC3339)
		recipe.UpdatedAt = cursor.Current.Lookup("updatedAt").Time().Format(time.RFC3339)

		recipes = append(recipes, recipe)
	}

	// Count total results
	count, err := r.collection.CountDocuments(context.Background(), filters)
	if err != nil {
		return nil, 0, err
	}

	return recipes, count, nil
}

// Insert updates or insert a new recipe
func (r *RecipeCollection) Insert(recipe Recipe) error {
	// Options to insert a record or update it if already exists
	updateOpts := options.Update()
	updateOpts.SetUpsert(true)

	if _, err := r.collection.UpdateOne(
		context.Background(),
		bson.D{
			{Key: "title", Value: recipe.Title},
		},
		bson.D{
			{Key: "$set", Value: recipe},
			{Key: "$setOnInsert", Value: bson.D{
				{Key: "createdAt", Value: time.Now().UTC()},
			}},
			{Key: "$currentDate", Value: bson.D{
				{Key: "updatedAt", Value: true},
			}},
		},
		updateOpts,
	); err != nil {
		return err
	}

	return nil
}

func (r *RecipeCollection) paramsToFilters(flt Filters) bson.D {
	AndItems := bson.A{}

	// Filter by term
	if flt.Term != "" {
		orItems := bson.A{
			bson.D{{Key: "$text", Value: bson.D{
				{Key: "$search", Value: flt.Term}}},
			},
		}

		orGroup := bson.D{{Key: "$or", Value: orItems}}
		AndItems = append(AndItems, orGroup)
	}

	// Ingredients filter
	if len(flt.Ingredients) > 0 {
		orItems := bson.A{}
		for i := range flt.Ingredients {
			orItems = append(orItems, bson.D{{Key: "ingredients", Value: flt.Ingredients[i]}})
		}

		andGroup := bson.D{{Key: "$and", Value: orItems}}
		AndItems = append(AndItems, andGroup)
	}

	filters := bson.D{}
	if len(AndItems) > 0 {
		AndBlock := bson.E{Key: "$and", Value: AndItems}
		filters = append(filters, AndBlock)
	}

	return filters
}
