//go:build !solution
// +build !solution

package exercise

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"unicode"
)

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
	return runWordCount(ctx, urls, workers)
}

// --- Implementation below (split for readability/testing helpers) ---

func runWordCount(ctx context.Context, urls []string, workers int) (map[string]int, error) {
	jobs := make(chan string, workers)
	results := make(chan map[string]int, workers)
	errCh := make(chan error, 1)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case url, ok := <-jobs:
					if !ok {
						return
					}
					counts, err := fetchAndCount(ctx, url)
					if err != nil {
						select {
						case errCh <- fmt.Errorf("fetching %s: %w", url, err):
							cancel()
						default:
						}
						return
					}
					select {
					case <-ctx.Done():
						return
					case results <- counts:
					}
				}
			}
		}()
	}

	go func() {
		for _, url := range urls {
			select {
			case <-ctx.Done():
				return
			case jobs <- url:
			}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	finalCounts := make(map[string]int)
	for counts := range results {
		for w, c := range counts {
			finalCounts[w] += c
		}
	}

	select {
	case err := <-errCh:
		return nil, err
	default:
	}

	return finalCounts, nil
}

func fetchAndCount(ctx context.Context, url string) (map[string]int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return tokenizeAndCount(string(body)), nil
}

func tokenizeAndCount(text string) map[string]int {
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
