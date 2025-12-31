//go:build !solution
// +build !solution

package workerpoolwordcount

// TODO: Import required packages
// You'll need:
// - "context" for cancellation
// - "fmt" for error formatting
// - "io" for reading HTTP bodies
// - "net/http" for HTTP requests
// - "strings" for text processing
// - "unicode" for character classification
// - "golang.org/x/sync/errgroup" for goroutine coordination
//
// Uncomment these imports when you implement the functions:
// import (
// 	"context"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"strings"
// 	"unicode"

// 	"golang.org/x/sync/errgroup"
// )

// Keep imports available for implementation
// These will be used when you implement the functions below
var (
	_ fmt.Formatter
	_ io.Reader
	_ http.Client
	_ strings.Builder
	_ unicode.RangeTable
	_ errgroup.Group
)

// WordCount fetches URLs concurrently using a worker pool pattern and returns
// aggregated word frequency counts from all URLs.
//
// Algorithm:
// 1. Create an errgroup for managing concurrent workers
// 2. Create job and result channels for communication
// 3. Launch worker goroutines that process jobs from the jobs channel
// 4. Send URLs to workers via jobs channel
// 5. Collect results from the results channel
// 6. Aggregate word counts into final map
// 7. Return aggregated counts and any errors
//
// CS Concepts:
// - Concurrency: Multiple workers process jobs in parallel
// - Channels: Thread-safe communication between goroutines
// - Worker Pool: Fixed number of workers handling variable workload
// - Aggregation: Combining results from multiple sources
// - Graceful Shutdown: Proper cleanup when done or on error
// TODO: Implement this function
func WordCount(ctx context.Context, urls []string, workers int) (map[string]int, error) {
	// TODO: Implement WordCount
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// fetchAndCount fetches a URL and returns word frequency counts
// TODO: Implement this function
//
// Algorithm:
// 1. Create HTTP request with context (allows cancellation)
// 2. Execute request and handle errors
// 3. Check response status code
// 4. Read response body
// 5. Tokenize and count words
// 6. Return word counts
//
// CS Concepts:
// - HTTP: Fetching remote resources
// - Error Handling: Checking status codes and I/O errors
// - Cancellation: Respecting context.Done() for graceful shutdown
func fetchAndCount(ctx context.Context, url string) (map[string]int, error) {
	// TODO: Step 1 - Create HTTP request with context

	// TODO: Step 2 - Execute request

	// TODO: Step 3 - Check status code

	// TODO: Step 4 - Read response body

	// TODO: Step 5 - Tokenize and count words

	// TODO: Step 6 - Return counts

	return nil, nil
}

// tokenizeAndCount splits text into words, normalizes them, and counts occurrences
// TODO: Implement this function
//
// Algorithm:
// 1. Create empty map for word counts
// 2. Split text into words using whitespace
// 3. For each word:
//    a. Convert to lowercase
//    b. Remove non-letter characters
//    c. Increment count if word is non-empty
// 4. Return word count map
//
// CS Concepts:
// - Tokenization: Breaking text into individual words
// - Normalization: Converting text to consistent form (lowercase, remove punctuation)
// - Frequency Counting: Using map to track occurrences
func tokenizeAndCount(text string) map[string]int {
	// TODO: Step 1 - Create word count map

	// TODO: Step 2 - Split text into words

	// TODO: Step 3 - Loop through words

	// TODO: Step 4 - Normalize word (lowercase, remove non-letters)

	// TODO: Step 5 - Skip empty words

	// TODO: Step 6 - Increment word count

	// TODO: Step 7 - Return map

	return nil
}

// After implementing:
// - Run: go test ./...
// - Check: go test -v for verbose output
// - Compare with solution.go to see detailed explanations
