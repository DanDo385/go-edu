package main

import (
	"context"
	"fmt"
	"time"

	"github.com/example/go-10x-minis/minis/06-worker-pool-wordcount/internal/workerpoolwordcount"
)

/*
Debug Harness for Worker Pool Wordcount Module

This file demonstrates the concurrent word counter.
Set breakpoints in exercise.go and press F5 in VS Code.
*/
func main() {
	fmt.Println("=== Worker Pool Wordcount Debug Harness ===")
	fmt.Println()

	// Test URLs (using example.com which is reliable)
	urls := []string{
		"https://example.com",
		"https://www.iana.org/domains/reserved",
	}

	fmt.Println("--- Demo: Concurrent Word Count ---")
	fmt.Printf("URLs: %v\n", urls)
	fmt.Printf("Workers: 3\n\n")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start := time.Now()
	freq, err := workerpoolwordcount.WordCount(ctx, urls, 3)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println()
		fmt.Println("Note: This demo requires internet access.")
		fmt.Println("If offline, the functions still demonstrate the pattern.")
		return
	}

	fmt.Printf("Completed in %v\n", elapsed)
	fmt.Printf("Total unique words: %d\n\n", len(freq))

	// Show top 10 words
	fmt.Println("Top 10 words:")
	count := 0
	for word, cnt := range freq {
		if count >= 10 {
			break
		}
		fmt.Printf("  %q: %d\n", word, cnt)
		count++
	}
	fmt.Println()

	fmt.Println("=== What You Learned ===")
	fmt.Println("1. Worker pool pattern for bounded concurrency")
	fmt.Println("2. Channels for work distribution")
	fmt.Println("3. errgroup.Group for coordinated error handling")
	fmt.Println("4. Context cancellation propagation")
	fmt.Println()
	fmt.Println("Next: Proceed to minis/07-generic-lru-cache")
}
