package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/go-10x-minis/minis/09-http-server-graceful/exercise"
)

func main() {
	// Create in-memory store
	store := exercise.NewMemStore()

	// Set up routes
	mux := http.NewServeMux()
	exercise.RegisterRoutes(mux, store)

	// Create server
	srv := exercise.NewServer(":8080", mux)

	// Start server in background
	go func() {
		log.Println("Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutdown signal received...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := exercise.GracefulShutdown(ctx, srv); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}

	log.Println("Server stopped gracefully")
}
