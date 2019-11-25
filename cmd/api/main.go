package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"github.com/georlav/recipeapi/mongoclient"

	"github.com/georlav/recipeapi/recipe"

	"github.com/georlav/recipeapi/handler"

	"github.com/georlav/recipeapi/config"
)

func main() {
	// Load configuration from file
	cfg, err := config.Load("config.json")
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration, %s", err))
	}

	// Initialize logger
	logger := log.New(
		os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile,
	)

	// Disable logger from writing to stdout
	if !cfg.APP.Debug {
		logger.SetOutput(ioutil.Discard)
	}

	// Mongo client
	client, err := mongoclient.NewClient(cfg.Mongo)
	if err != nil {
		logger.Fatalf(`mongo client error, %s`, err)
	}

	// Initialize mongo
	db := client.Database(cfg.Mongo.Database)
	rCollection := db.Collection(cfg.Mongo.RecipeCollection)

	// Create searchable index
	iv := rCollection.Indexes()
	if _, err := iv.CreateOne(context.Background(), mongo.IndexModel{
		Keys: bsonx.Doc{{Key: "title", Value: bsonx.String("text")}},
	}); err != nil {
		logger.Fatal(err)
	}

	// initialize repository
	rr := recipe.NewMongoRepo(rCollection, cfg.Mongo)

	// Initialize handlers
	h := handler.NewHandler(rr, cfg, logger)

	// Initialize API routes
	r := handler.Routes(h)

	// Initialize server
	s := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
		Handler:      r,
	}

	// Start listening to incoming requests
	go func() {
		logger.Printf("Started web server at %s://%s%s", cfg.Server.Scheme, cfg.Server.Host, s.Addr)
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatalf("Server error, %s", err)
		}
	}()

	// Keep application open, close on termination signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	// Gracefully Shutdown server
	logger.Println("Application received a termination signal. Shutting down.")

	if err := s.Shutdown(context.Background()); err != nil {
		logger.Fatalf("Failed to gracefully shutdown http server, %s", err)
	}
}
