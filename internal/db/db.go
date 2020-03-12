package db

import (
	"fmt"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/db/mongodb"
	"github.com/georlav/recipeapi/internal/db/mysql"
)

func New(cfg config.Config) (Queryable, error) {
	switch cfg.APP.Database {
	case "MySQL":
		mdb, err := mysql.New(cfg.MySQL)
		if err != nil {
			return nil, err
		}

		return mysql.NewRecipe(mdb), nil
	case "MongoDB":
		mc, err := mongodb.New(cfg.Mongo)
		if err != nil {
			return nil, err
		}
		mdb := mc.Database(cfg.Mongo.Database)
		rCollection := mdb.Collection(cfg.Mongo.RecipeCollection)

		return mongodb.NewRecipe(rCollection), nil
	default:
		return nil, fmt.Errorf("invalid database option, supported options: MySQL, MongoDB")
	}
}
