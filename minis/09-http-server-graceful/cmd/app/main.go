package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/example/go-10x-minis/minis/09-http-server-graceful/internal/httpservergraceful"
)

func main() {
	// Create a new ServeMux to handle routes
	mux := http.NewServeMux()

	// A simple handler for the root path
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello! Try the /slow endpoint.")
	})

	// This handler takes 5 seconds to complete, simulating a long-running task.
	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request for /slow")
		time.Sleep(5 * time.Second)
		log.Println("Finished /slow request")
		fmt.Fprintln(w, "Slow request complete!")
	})

	// Create the http.Server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Server starting on :8080. Press Ctrl+C to shut down.")

	// Run the server with graceful shutdown logic
	err := httpservergraceful.RunGracefulServer(context.Background(), srv)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}

	log.Println("Server exited gracefully.")
}