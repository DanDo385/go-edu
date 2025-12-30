//go:build !solution
// +build !solution

package exercise

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"unicode"
)

/*
Core Logic Without Error Handling
==================================

This file demonstrates the core worker pool algorithm without
error handling complexity. It's designed to help you understand the
fundamental logic flow.

KEY CONCEPTS:
- Concurrent worker pools
- Channel-based communication
- Result aggregation
- Word tokenization and counting
*/

// coreWordCount demonstrates concurrent word counting without error handling
//
// Algorithm steps:
// 1. Create WaitGroup for worker management
// 2. Create channels for job distribution and result collection
// 3. Launch worker goroutines
// 4. Send URLs as jobs to workers
// 5. Collect results from all workers
// 6. Aggregate word counts
// 7. Return final aggregated counts
func coreWordCount(ctx context.Context, urls []string, workers int) map[string]int {
	// Step 1 - Create WaitGroup
	var wg sync.WaitGroup

	// Step 2 - Create jobs and results channels
	jobs := make(chan string, workers)
	results := make(chan map[string]int, workers)

	// Step 3 - Launch worker goroutines
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range jobs {
				results <- coreFetchAndCount(ctx, url)
			}
		}()
	}

	// Step 4 - Send jobs in separate goroutine
	go func() {
		defer close(jobs)
		for _, url := range urls {
			jobs <- url
		}
	}()

	// Step 5 - Close results when workers finish
	go func() { // Launch goroutine to coordinate shutdown
		wg.Wait()      // Wait for all workers to finish (blocks until counter=0)
		close(results) // Close results channel (signals "no more results coming")
	}() // Launch goroutine to coordinate shutdown

	// Step 6 - Aggregate results
	finalCounts := make(map[string]int)
	for counts := range results {
		for word, count := range counts {
			finalCounts[word] += count
		}
	}

	// Step 7 - Return aggregated counts
	return finalCounts
}

// coreFetchAndCount fetches URL and returns word counts (no error handling)
//
// Algorithm steps:
// 1. Create HTTP request
// 2. Execute request
// 3. Read response body
// 4. Tokenize and count words
// 5. Return counts
func coreFetchAndCount(ctx context.Context, url string) map[string]int {
	// Step 1 - Create and execute HTTP request (assumes success)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	// Step 2 - Read response body (assumes success)
	body, _ := io.ReadAll(resp.Body)

	// Step 3 - Tokenize and count words
	return coreTokenizeAndCount(string(body))
}

// coreTokenizeAndCount splits text into words and counts them (no error handling)
//
// Algorithm steps:
// 1. Create word count map
// 2. Split text into words
// 3. Normalize each word (lowercase, remove non-letters)
// 4. Increment count for each word
// 5. Return counts
func coreTokenizeAndCount(text string) map[string]int {
	counts := make(map[string]int)

	for _, word := range strings.Fields(text) {
		word = strings.ToLower(word)
		word = strings.Map(func(r rune) rune {
			if unicode.IsLetter(r) {
				return r
			}
			return -1
		}, word)

		if word == "" {
			continue
		}

		counts[word]++
	}

	return counts
}
