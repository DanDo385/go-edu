//go:build !solution
// +build !solution

package exercise

import "context"

// WordCount fetches URLs concurrently using a worker pool, tokenizes response bodies,
// and returns overall word frequencies.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - urls: List of URLs to fetch
//   - workers: Number of concurrent workers (goroutines)
//
// Returns:
//   - map[string]int: Word frequencies across all fetched pages
//   - error: Non-nil if any fetch fails (cancels all other fetches)
//
// Behavior:
//   - Words are normalized to lowercase
//   - Only alphabetic characters are kept (punctuation removed)
//   - Empty words are ignored
//   - If any fetch fails, all in-flight requests are cancelled
//
// Example:
//   counts, err := WordCount(ctx, []string{"http://example.com"}, 3)
func WordCount(ctx context.Context, urls []string, workers int) (map[string]int, error) {
	// TODO: implement
	return nil, nil
}
