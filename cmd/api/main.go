package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/database"
	"github.com/georlav/recipeapi/internal/handler"
	"github.com/georlav/recipeapi/internal/logger"
)

// @title Recipe API
// @version 1.0
// @description Simple api that serves recipes for puppies. This project is a step by step guide on how to create a
// @description simple api using Go programming language. The purpose of the project is to demonstrate to new comers
// @description the language basic features and concepts.
// @termsOfService http://swagger.io/terms/

// @license.name MIT License
// @license.url https://raw.githubusercontent.com/georlav/recipeapi/master/LICENSE

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Load configuration from file
	cfg, err := config.Load("config.json")
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration, %s", err))
	}

	// Init logger
	log := logger.NewLogger(cfg.Logger)

	// Init database
	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize handlers
	h := handler.NewHandler(db, cfg, log)

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
		log.Printf("Started web server at %s://%s%s", cfg.Server.Scheme, cfg.Server.Host, s.Addr)
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error, %s", err)
		}
	}()

	// Keep application open, close on termination signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	// Gracefully Shutdown server
	log.Println("Application received a termination signal. Shutting down.")

	if err := s.Shutdown(context.Background()); err != nil {
		log.Fatalf("Failed to gracefully shutdown http server, %s", err)
	}
}
