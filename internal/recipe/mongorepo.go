package recipe

import (
	"context"
	"fmt"
	"time"

	"github.com/georlav/recipeapi/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRepo object
type MongoRepo struct {
	recipes *mongo.Collection
	cfg     config.Mongo
}

// NewMongoRepo create a new Mongo repository object
func NewMongoRepo(recipe *mongo.Collection, c config.Mongo) *MongoRepo {
	return &MongoRepo{
		recipes: recipe,
		cfg:     c,
	}
}

// GetOne returns a recipe by its unique identifier
func (m MongoRepo) GetOne(id string) (r Recipe, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return r, fmt.Errorf("invalid signal id (%s), %w", id, err)
	}

	result := m.recipes.FindOne(
		context.Background(),
		bson.D{{Key: "_id", Value: oid}},
	)
	if err := result.Decode(&r); err != nil {
		return r, err
	}

	return r, nil
}

// GetMany return signals filtered and paginated
func (m MongoRepo) GetMany(qp QueryParams) (r Recipes, totalCount int64, err error) {
	// find options
	fOpts := options.Find()
	fOpts.SetLimit(10)
	if qp.Page > 0 {
		qp.Page--
	}
	fOpts.SetSkip(qp.Page * 10)

	// filter options
	filters := paramsToFilters(qp)

	// Do query
	cursor, err := m.recipes.Find(
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
	for cursor.Next(context.TODO()) {
		recipe := Recipe{}
		err = cursor.Decode(&recipe)
		if err != nil {
			return nil, 0, err
		}

		recipe.ID = recipe.MID.Hex()
		recipe.CreatedAt = cursor.Current.Lookup("createdAt").Time().Format(time.RFC3339)
		recipe.UpdatedAt = cursor.Current.Lookup("updatedAt").Time().Format(time.RFC3339)

		r = append(r, recipe)
	}

	// Count total results
	count, err := m.recipes.CountDocuments(context.Background(), filters)
	if err != nil {
		return nil, 0, err
	}

	return r, count, nil
}

func paramsToFilters(qp QueryParams) bson.D {
	AndItems := bson.A{}

	// Filter by term
	if qp.Term != "" {
		orItems := bson.A{
			bson.D{{Key: "$text", Value: bson.D{
				{Key: "$search", Value: qp.Term}}},
			},
		}

		orGroup := bson.D{{Key: "$or", Value: orItems}}
		AndItems = append(AndItems, orGroup)
	}

	// Ingredients filter
	if len(qp.Ingredients) > 0 {
		orItems := bson.A{}
		for i := range qp.Ingredients {
			orItems = append(orItems, bson.D{{Key: "ingredients", Value: qp.Ingredients[i]}})
		}

		orGroup := bson.D{{Key: "$or", Value: orItems}}
		AndItems = append(AndItems, orGroup)
	}

	filters := bson.D{}
	if len(AndItems) > 0 {
		AndBlock := bson.E{Key: "$and", Value: AndItems}
		filters = append(filters, AndBlock)
	}

	return filters
}

// Insert updates or inserts a new record in signals collection
func (m MongoRepo) Insert(recipes ...Recipe) error {
	// Options to update record or insert if exists (upsert)
	updateOpts := options.Update()
	updateOpts.SetUpsert(true)

	for i := range recipes {
		if _, err := m.recipes.UpdateOne(
			context.Background(),
			bson.D{
				{Key: "title", Value: recipes[i].Title},
			},
			bson.D{
				{Key: "$set", Value: recipes[i]},
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
	}

	return nil
}
