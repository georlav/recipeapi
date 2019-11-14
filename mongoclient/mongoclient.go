package mongoclient

import (
	"context"
	"fmt"
	"time"

	"github.com/georlav/recipeapi/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewClient returns a new ready to use mongo client
func NewClient(cfg config.Mongo) (*mongo.Client, error) {
	// mongo client options, for aws managed mongoDB limits visit url
	// https://docs.aws.amazon.com/documentdb/latest/developerguide/limits.html
	mcOpts, err := clientOptions(cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid client options %w", err)
	}

	// Mongo Client initialization
	mClient, err := mongo.Connect(context.Background(), mcOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize mongo client: %w", err)
	}

	// verify that the client can connect
	err = mClient.Ping(context.Background(), nil)
	if err != nil {
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

	cOpts.SetServerSelectionTimeout(10 * time.Second)
	cOpts.SetConnectTimeout(10 * time.Second)
	cOpts.SetSocketTimeout(15 * time.Second)
	cOpts.SetMaxConnIdleTime(30 * time.Second)
	cOpts.SetRetryWrites(true)

	if err := cOpts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid client options: %w", err)
	}

	return cOpts, nil
}
