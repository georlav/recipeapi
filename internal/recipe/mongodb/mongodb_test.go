package mongodb_test

import (
	"testing"
	"time"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/recipe/mongodb"
)

func TestNew(t *testing.T) {
	cfg := config.Mongo{
		Host:                      "127.0.0.1",
		Port:                      27017,
		Username:                  "root",
		Password:                  "toor",
		Database:                  "recipes",
		RecipeCollection:          "recipe",
		PoolSize:                  10,
		Timeout:                   15 * time.Second,
		SetServerSelectionTimeout: 15 * time.Second,
		SetMaxConnIdleTime:        15 * time.Second,
		SetRetryWrites:            false,
	}

	_, err := mongodb.New(cfg)
	if err != nil {
		t.Fatalf("Mongo client failure, %s", err)
	}
}
