package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("=== Mini Service Demo ===")
	fmt.Println()
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"Welcome to Mini Service"}`))
	})
	
	go func() {
		log.Printf("Server starting on :%s", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()
	
	// Wait for interrupt
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	
	fmt.Println("\nShutting down...")
}
