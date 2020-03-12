package db

import (
	"context"
	"fmt"

	"github.com/georlav/recipeapi/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func New(cfg config.Config) (Queryable, error) {
	switch cfg.APP.Database {
	case "MySQL":
		mdb, err := NewMySQL(cfg.MySQL)
		if err != nil {
			return nil, err
		}

		return NewRecipeTable(mdb), nil
	case "MongoDB":
		mc, err := NewMongoDB(cfg.Mongo)
		if err != nil {
			return nil, err
		}
		mdb := mc.Database(cfg.Mongo.Database)
		rCollection := mdb.Collection(cfg.Mongo.RecipeCollection)

		iv := rCollection.Indexes()
		_, err = iv.CreateOne(context.Background(), mongo.IndexModel{
			Keys: bsonx.Doc{
				{Key: "title", Value: bsonx.String("text")},
			},
		})
		if err != nil {
			return nil, err
		}

		return NewRecipeCollection(rCollection), nil
	default:
		return nil, fmt.Errorf("invalid database option, supported options: MySQL, MongoDB")
	}
}
