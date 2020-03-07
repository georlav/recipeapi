package mongoclient_test

import (
	"testing"
	"time"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/mongoclient"
)

func TestNewMongoClient(t *testing.T) {
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

	_, err := mongoclient.NewClient(cfg)
	if err != nil {
		t.Fatalf("Mongo client failure, %s", err)
	}
}
