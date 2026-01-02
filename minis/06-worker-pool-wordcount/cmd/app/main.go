package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/example/go-10x-minis/minis/06-worker-pool-wordcount/internal/workerpoolwordcount"
)

/*
Concurrent Word Counter CLI

Usage:
  go run ./cmd/app/main.go [--workers N] <url1> <url2> ...

Arguments:
  --workers N   Number of concurrent workers (default: 5)
  urls          One or more URLs to fetch and count words from

Examples:
  go run ./cmd/app/main.go https://example.com https://golang.org
  go run ./cmd/app/main.go --workers 3 https://example.com
*/
func main() {
	workers := 5
	var urls []string

	// Parse arguments
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		if args[i] == "--workers" && i+1 < len(args) {
			w, err := strconv.Atoi(args[i+1])
			if err == nil && w > 0 {
				workers = w
			}
			i++
		} else {
			urls = append(urls, args[i])
		}
	}

	if len(urls) == 0 {
		printUsage()
		os.Exit(1)
	}

	fmt.Printf("Counting words from %d URLs with %d workers...\n\n", len(urls), workers)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	freq, err := workerpoolwordcount.WordCount(ctx, urls, workers)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Show results
	fmt.Printf("=== Word Count Results ===\n")
	fmt.Printf("Total unique words: %d\n\n", len(freq))

	// Show top 20 words
	fmt.Println("Top words:")
	count := 0
	for word, cnt := range freq {
		if count >= 20 {
			fmt.Println("  ...")
			break
		}
		fmt.Printf("  %q: %d\n", word, cnt)
		count++
	}
}

func printUsage() {
	fmt.Println("Concurrent Word Counter CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run ./cmd/app/main.go [--workers N] <url1> <url2> ...")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  --workers N   Number of concurrent workers (default: 5)")
	fmt.Println("  urls          One or more URLs to fetch and count words from")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run ./cmd/app/main.go https://example.com https://golang.org")
	fmt.Println("  go run ./cmd/app/main.go --workers 3 https://example.com")
}
