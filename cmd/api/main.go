package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	s := http.Server{
		Addr:         "127.0.0.1:8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
		Handler:      nil,
	}

	// Start listening to incoming requests
	go func() {
		fmt.Println("Starting web server at", "http://127.0.0.1:8080")
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error, %s", err)
		}
	}()

	// Keep application open, close on termination signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	// Gracefully Shutdown server
	fmt.Println("Application received a termination signal. Shutting down.")

	if err := s.Shutdown(context.Background()); err != nil {
		log.Fatalf("Failed to gracefully shutdown http server, %s", err)
	}
}
