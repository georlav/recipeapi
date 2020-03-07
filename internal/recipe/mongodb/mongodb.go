package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/georlav/recipeapi/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// New returns a new ready to use mongo client
func New(cfg config.Mongo) (*mongo.Client, error) {
	mcOpts, err := clientOptions(cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid client options %w", err)
	}

	// Mongo Client initialization
	mClient, err := mongo.Connect(context.Background(), mcOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize mongo client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// verify that the client can connect
	if err = mClient.Ping(ctx, readpref.PrimaryPreferred()); err != nil {
		return nil, fmt.Errorf("failed to ping mongo DB, client options: %+v, error: %w", mcOpts, err)
	}

	return mClient, nil
}

func clientOptions(cfg config.Mongo) (*options.ClientOptions, error) {
	cOpts := options.Client()
	cOpts.ApplyURI(
		fmt.Sprintf(`mongodb://%s:%s@%s:%d`,
			cfg.Username, cfg.Password, cfg.Host, cfg.Port,
		),
	)

	cOpts.SetServerSelectionTimeout(cfg.SetServerSelectionTimeout * time.Second)
	cOpts.SetMaxConnIdleTime(cfg.SetMaxConnIdleTime * time.Second)
	cOpts.SetRetryWrites(cfg.SetRetryWrites)

	if err := cOpts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid client options: %w", err)
	}

	return cOpts, nil
}
