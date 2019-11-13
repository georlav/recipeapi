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

	"github.com/georlav/recipeapi/handler"

	"github.com/georlav/recipeapi/config"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration, %s", err))
	}

	// initialize logger
	logger := log.New(
		os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile,
	)

	// Disable logger from writing to stdout
	if !cfg.APP.Debug {
		logger.SetOutput(ioutil.Discard)
	}

	// Initialize handlers
	h := handler.NewHandler(cfg, logger)

	// Initialize api routes
	r := handler.Routes(h)

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
