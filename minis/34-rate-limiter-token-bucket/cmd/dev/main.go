package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/example/go-10x-minis/minis/34-rate-limiter-token-bucket/exercise"
)

func main() {
	// Create a rate limiter: 10 requests per second, burst capacity of 20
	limiter := exercise.NewRateLimiter(20, 10.0)

	// Set up routes
	mux := http.NewServeMux()

	// API endpoint - rate limited
	mux.HandleFunc("/api/data", handleData)

	// Health check - not rate limited
	mux.HandleFunc("/health", handleHealth)

	// Stats endpoint - shows rate limiter stats
	mux.HandleFunc("/stats", handleStats(limiter))

	// Apply rate limiting middleware to all routes except /health
	handler := selectiveRateLimit(limiter, mux)

	// Add logging middleware
	handler = loggingMiddleware(handler)

	// Start server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Server starting on :8080")
	log.Println("Rate limit: 10 requests/second, burst capacity: 20")
	log.Println("")
	log.Println("Try these commands:")
	log.Println("  curl http://localhost:8080/api/data")
	log.Println("  curl http://localhost:8080/stats")
	log.Println("  curl http://localhost:8080/health")
	log.Println("")
	log.Println("Test rate limiting with:")
	log.Println("  for i in {1..25}; do curl -w '\\n' http://localhost:8080/api/data; done")
	log.Println("")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// handleData simulates an API endpoint that returns some data
func handleData(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"message":   "Success! This endpoint is rate limited.",
		"timestamp": time.Now().Unix(),
		"tip":       "Try making many requests quickly to see rate limiting in action",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// handleHealth is a simple health check endpoint (not rate limited)
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

// handleStats returns statistics about the rate limiter
func handleStats(limiter *exercise.RateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := limiter.Stats()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}

// selectiveRateLimit applies rate limiting to all routes except /health
func selectiveRateLimit(limiter *exercise.RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip rate limiting for health checks
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		// Apply rate limiting middleware
		limiter.Middleware(next).ServeHTTP(w, r)
	})
}

// loggingMiddleware logs each request
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response recorder to capture status code
		rec := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rec, r)

		duration := time.Since(start)

		// Log with color based on status code
		var statusColor string
		switch {
		case rec.statusCode >= 500:
			statusColor = "\033[31m" // Red
		case rec.statusCode >= 400:
			statusColor = "\033[33m" // Yellow
		case rec.statusCode >= 300:
			statusColor = "\033[36m" // Cyan
		default:
			statusColor = "\033[32m" // Green
		}
		resetColor := "\033[0m"

		log.Printf("%s%d%s %s %s (%v)",
			statusColor, rec.statusCode, resetColor,
			r.Method,
			r.URL.Path,
			duration,
		)
	})
}

// responseRecorder wraps http.ResponseWriter to capture status code
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *responseRecorder) WriteHeader(statusCode int) {
	rec.statusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *responseRecorder) Write(b []byte) (int, error) {
	if rec.statusCode == 0 {
		rec.statusCode = http.StatusOK
	}
	return rec.ResponseWriter.Write(b)
}

// Example of how to use rate limiter programmatically
func exampleUsage() {
	// Create rate limiter: 100 requests per minute (burst of 20)
	limiter := exercise.NewRateLimiter(20, 100.0/60.0) // 100/60 = requests per second

	clientIP := "192.168.1.1"

	for i := 0; i < 25; i++ {
		if limiter.Allow(clientIP) {
			fmt.Printf("Request %d: Allowed\n", i+1)
		} else {
			fmt.Printf("Request %d: Rate limited\n", i+1)
		}
	}

	// Wait for tokens to refill
	time.Sleep(2 * time.Second)

	fmt.Println("\nAfter waiting 2 seconds:")
	if limiter.Allow(clientIP) {
		fmt.Println("Request allowed (tokens refilled)")
	}
}
