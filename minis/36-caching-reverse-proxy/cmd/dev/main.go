package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/example/go-10x-minis/minis/36-caching-reverse-proxy/exercise"
)

func main() {
	// Backend server URL (what we're proxying to)
	// You can change this to any backend server
	// For testing, use the example_backend.go in this project:
	//   go run example_backend.go
	backendURL := "http://localhost:9000"

	// Or use httpbin.org for public testing:
	// backendURL := "http://httpbin.org"

	// Parse backend URL
	target, err := url.Parse(backendURL)
	if err != nil {
		log.Fatal(err)
	}

	// Create cache
	// maxSize: 100 entries, TTL: 5 minutes
	cache := exercise.NewCache(100, 5*time.Minute)

	// Create caching proxy
	proxy := cache.NewCachingProxy(target)

	// Routes
	http.Handle("/", proxy)
	http.HandleFunc("/stats", cache.StatsHandler())
	http.HandleFunc("/cache/clear", cache.ClearHandler())

	// Start server
	addr := ":8080"
	fmt.Printf("🚀 Caching reverse proxy running on http://localhost%s\n", addr)
	fmt.Printf("📊 Stats endpoint: http://localhost%s/stats\n", addr)
	fmt.Printf("🗑️  Clear cache: http://localhost%s/cache/clear\n", addr)
	fmt.Printf("🎯 Proxying to: %s\n\n", backendURL)
	fmt.Println("Try these requests:")
	fmt.Println("  # First request - cache miss (slow)")
	fmt.Printf("  curl http://localhost%s/api/data\n", addr)
	fmt.Println("\n  # Second request - cache hit (fast!)")
	fmt.Printf("  curl http://localhost%s/api/data\n", addr)
	fmt.Println("\n  # View cache statistics")
	fmt.Printf("  curl http://localhost%s/stats\n", addr)
	fmt.Println("\n  # Clear cache")
	fmt.Printf("  curl -X POST http://localhost%s/cache/clear\n\n", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}
