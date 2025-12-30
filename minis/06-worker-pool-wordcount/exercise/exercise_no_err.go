//go:build !solution
// +build !solution

package exercise

import (
	"context"
	"io"
	"net/http"
	"strings"
	"unicode"

	"golang.org/x/sync/errgroup"
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

// CoreWordCount demonstrates concurrent word counting without error handling
//
// Algorithm steps:
// 1. Create errgroup for worker management
// 2. Create channels for job distribution and result collection
// 3. Launch worker goroutines
// 4. Send URLs as jobs to workers
// 5. Collect results from all workers
// 6. Aggregate word counts
// 7. Return final aggregated counts
func CoreWordCount(ctx context.Context, urls []string, workers int) map[string]int {
	// TODO: Step 1 - Create errgroup

	// TODO: Step 2 - Create jobs and results channels

	// TODO: Step 3 - Launch worker goroutines

	// TODO: Step 4 - Send jobs in separate goroutine

	// TODO: Step 5 - Close results when workers finish

	// TODO: Step 6 - Aggregate results

	// TODO: Step 7 - Return aggregated counts

	return nil
}

// CoreFetchAndCount fetches URL and returns word counts (no error handling)
//
// Algorithm steps:
// 1. Create HTTP request
// 2. Execute request
// 3. Read response body
// 4. Tokenize and count words
// 5. Return counts
func CoreFetchAndCount(ctx context.Context, url string) map[string]int {
	// TODO: Step 1 - Create and execute HTTP request

	// TODO: Step 2 - Read response body

	// TODO: Step 3 - Tokenize and count words

	// TODO: Step 4 - Return counts

	return nil
}

// CoreTokenizeAndCount splits text into words and counts them (no error handling)
//
// Algorithm steps:
// 1. Create word count map
// 2. Split text into words
// 3. Normalize each word (lowercase, remove non-letters)
// 4. Increment count for each word
// 5. Return counts
func CoreTokenizeAndCount(text string) map[string]int {
	// TODO: Step 1 - Create map

	// TODO: Step 2 - Split text into words

	// TODO: Step 3 - Loop and normalize words

	// TODO: Step 4 - Count words

	// TODO: Step 5 - Return counts

	return nil
}
