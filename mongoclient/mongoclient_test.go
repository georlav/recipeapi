package mongoclient_test

import (
	"testing"

	"github.com/georlav/recipeapi/mongoclient"

	"github.com/georlav/recipeapi/config"
)

func TestNewMongoClient(t *testing.T) {
	cfg := config.Mongo{
		Host:           "127.0.0.1",
		Port:           27017,
		Username:       "root",
		Password:       "toor",
		Database:       "recipes",
		PoolSize:       100,
		Timeout:        15,
		SetRetryWrites: false,
	}

	_, err := mongoclient.NewClient(cfg)
	if err != nil {
		t.Fatalf("Mongo client failure, %s", err)
	}
}
