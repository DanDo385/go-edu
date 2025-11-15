// +build ignore

package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

var requestCount int64

func main() {
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt64(&requestCount, 1)

		// Simulate slow backend (200ms database query)
		time.Sleep(200 * time.Millisecond)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "max-age=60") // Cache for 60 seconds

		fmt.Fprintf(w, `{"message":"Hello from backend","request_number":%d,"timestamp":"%s"}`,
			count, time.Now().Format(time.RFC3339))
	})

	http.HandleFunc("/api/no-cache", func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt64(&requestCount, 1)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store") // Never cache

		fmt.Fprintf(w, `{"message":"This should not be cached","request_number":%d}`, count)
	})

	http.HandleFunc("/api/slow", func(w http.ResponseWriter, r *http.Request) {
		// Very slow endpoint (1 second)
		time.Sleep(1 * time.Second)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "max-age=300") // Cache for 5 minutes

		fmt.Fprintf(w, `{"message":"Slow response","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Backend server is running!\n")
		fmt.Fprintf(w, "Total requests handled: %d\n", atomic.LoadInt64(&requestCount))
	})

	addr := ":9000"
	fmt.Printf("Backend server running on http://localhost%s\n", addr)
	fmt.Println("Endpoints:")
	fmt.Println("  /api/data - Cacheable (200ms delay)")
	fmt.Println("  /api/no-cache - Not cacheable")
	fmt.Println("  /api/slow - Very slow (1s delay)")
	log.Fatal(http.ListenAndServe(addr, nil))
}
