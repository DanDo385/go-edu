//go:build !solution && !reference

package workerpoolwordcount

/*
Problem: Fetch multiple URLs concurrently and count word frequencies
Constraints:
- Fixed number of workers (bounded concurrency)
- Must handle errors from any worker
- Must cancel all work if one worker fails
- Aggregate results from all workers
Time/Space Complexity:
- Time: O(n) where n = total words across all URLs (concurrent fetching)
- Space: O(u) where u = number of unique words across all URLs
*/

import (
	"context"  // Context for cancellation and timeouts
	"fmt"      // Error formatting
	"io"       // Reading HTTP response bodies
	"net/http" // HTTP client
	"strings"  // String manipulation (Fields, ToLower, Map)

// WordCount - TODO: implement this function
func WordCount(ctx context.Context, urls []string, workers int) (map[string]int, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil, nil
}

// WordCountWithErrGroup - TODO: implement this function
func WordCountWithErrGroup(ctx context.Context, urls []string, workers int) (map[string]int, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil, nil
}

// fetchAndCount - TODO: implement this function
func fetchAndCount(ctx context.Context, url string) (map[string]int, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil, nil
}

// tokenizeAndCount - TODO: implement this function
func tokenizeAndCount(text string) map[string]int {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil
}

