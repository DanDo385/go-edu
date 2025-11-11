package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/example/go-10x-minis/minis/06-worker-pool-wordcount/exercise"
)

func main() {
	// Create test HTTP servers (simulating real websites)
	servers := createTestServers()
	defer closeServers(servers)

	// Extract URLs
	urls := make([]string, len(servers))
	for i, srv := range servers {
		urls[i] = srv.URL
	}

	fmt.Printf("Fetching %d URLs with 3 workers...\n\n", len(urls))

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Run word count
	start := time.Now()
	counts, err := exercise.WordCount(ctx, urls, 3)
	if err != nil {
		log.Fatalf("WordCount failed: %v", err)
	}
	elapsed := time.Since(start)

	// Display results
	fmt.Println("Top 10 words:")
	type pair struct {
		word  string
		count int
	}
	var pairs []pair
	for word, count := range counts {
		pairs = append(pairs, pair{word, count})
	}
	// Simple sort (top 10)
	for i := 0; i < len(pairs) && i < 10; i++ {
		maxIdx := i
		for j := i + 1; j < len(pairs); j++ {
			if pairs[j].count > pairs[maxIdx].count {
				maxIdx = j
			}
		}
		pairs[i], pairs[maxIdx] = pairs[maxIdx], pairs[i]
		fmt.Printf("%2d. %-15s %d\n", i+1, pairs[i].word, pairs[i].count)
	}

	fmt.Printf("\nTotal unique words: %d\n", len(counts))
	fmt.Printf("Completed in: %v\n", elapsed)
}

func createTestServers() []*httptest.Server {
	return []*httptest.Server{
		httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Go is a great programming language for building scalable systems")
		})),
		httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Concurrency in Go is simple and powerful with goroutines and channels")
		})),
		httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Go programming makes it easy to build reliable and efficient software")
		})),
		httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "The Go standard library is comprehensive and well-designed")
		})),
		httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Building web services with Go is straightforward and productive")
		})),
	}
}

func closeServers(servers []*httptest.Server) {
	for _, srv := range servers {
		srv.Close()
	}
}
