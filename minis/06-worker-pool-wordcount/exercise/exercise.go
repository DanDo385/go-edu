//go:build !solution
// +build !solution

package exercise
// TODO: Import required packages
// You'll need:
// - "context" for cancellation
// - "fmt" for error formatting
// - "io" for reading HTTP bodies
// - "net/http" for HTTP requests
// - "strings" for text processing
// - "unicode" for character classification
// - "golang.org/x/sync/errgroup" for simplified goroutine coordination
import (
	"context" // Context for cancellation
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode"

	"golang.org/x/sync/errgroup"
)

/*
WordCount fetches URLs concurrently using errgroup.

ANALOGY: Smart Factory Manager
------------------------------
errgroup is like a smart factory manager who:
- Automatically tracks workers (no manual WaitGroup management)
- Automatically stops everything on first error (no error channel needed)
- Automatically cleans up resources (no manual defer cancel needed)

This makes the code ~40% shorter and much safer!
*/

func WordCount(ctx context.Context, urls []string, workers int) (map[string]int, error) {
	// TODO: Step 1 - Create errgroup
	//   Use errgroup.WithContext(ctx) to get:
	//   - g: *errgroup.Group (manages goroutines and errors)
	//   - ctx: context.Context (cancellable context)
	g, ctx := errgroup.WithContext(ctx) // Create errgroup and cancellable context
	// TODO: Step 2 - Create channels
	jobs := make(chan string, workers) // Create jobs channel with buffered channel
	results := make(chan map[string]int, workers) // Create results channel with buffered channel

	// TODO: Step 3 - Launch workers using errgroup
	//   Key points:
	//   - g.Go() automatically calls WaitGroup.Add(1) and defer Done()
	//   - Just return error - errgroup cancels context automatically
	for i := 0; i < workers; i++ { // Loop over worker count
		g.Go(func() error { // Launch worker goroutine
			for { // Infinite loop
				select { // Select on ctx.Done() and jobs channel
				case <-ctx.Done(): // Stop if cancelled
					return ctx.Err() // Return error if cancelled
				case url, ok := <-jobs: // Receive URL from jobs channel
					if !ok { // Channel closed, no more jobs
						return nil // Exit normally
					}
					counts, err := fetchAndCount(ctx, url) // Fetch and count words
					// fetchAndCount takes context and url, returns map[string]int and error
					if err != nil { // Error fetching and counting words
						return fmt.Errorf("fetching %s: %w", url, err) // Return error (errgroup handles cancellation)
					}
					// Send results with cancellation check
					select {
					case <-ctx.Done(): // Check if cancelled while processing
						return ctx.Err()
					case results <- counts: // Sends word counts to results channel
						// Successfully sent
					}
				}
			}
		}) // Add worker to errgroup
	}
	// TODO: Step 4 - Send jobs in separate goroutine
	go func() { // Launch goroutine to send jobs
		defer close(jobs)          // Close jobs channel when done
		for _, url := range urls { // Range over urls
			select { // Check cancellation before sending each URL
			case <-ctx.Done(): // Stop if cancelled
				return // Exit if cancelled
			case jobs <- url: // Send URL to jobs channel
				// Successfully sent URL to jobs channel
			}
		}
	}()

	// TODO: Step 5 - Close results when workers finish
	go func() { // Launch goroutine to close results channel
		g.Wait()       // Wait for all workers to finish
		close(results) // Close results channel
	}()
	// TODO: Step 6 - Aggregate results
	finalCounts := make(map[string]int) // Create map to store final counts
	for counts := range results { // Range over results channel
		for word, count := range counts { // Range over map
			finalCounts[word] += count // Accumulate counts for each word
		}
	}
	// TODO: Step 7 - Check for errors
	if err := g.Wait(); err != nil { // Wait for all workers and check for errors
		return nil, err // Return error if any occurred
	}
	return finalCounts, nil // Return final counts
}

func fetchAndCount(ctx context.Context, url string) (map[string]int, error) { // Fetch and count words from URL
	// TODO: Implement fetchAndCount
	//   1. Create HTTP request: http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil) // Create HTTP request
	if err != nil { // Error creating request
		return nil, err // Return error
	}
	//   2. Execute request: http.DefaultClient.Do(req) with error handling
	resp, err := http.DefaultClient.Do(req) // Execute request
	if err != nil { // Error executing request
		return nil, err // Return error
	}
	defer resp.Body.Close() // Close response body
	//   3. Check if status code is not 200 OK
	if resp.StatusCode != http.StatusOK { // Status code is not 200 OK
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode) // Return error
	}
	//   4. Read body: io.ReadAll(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//   5. Convert to string and call tokenizeAndCount
	text := string(body)
	counts := tokenizeAndCount(text)
	//   7. Return counts and error
	return counts, nil
}

func tokenizeAndCount(text string) map[string]int {
	// TODO: Implement tokenizeAndCount
	//   1. Create map: counts := make(map[string]int)
	counts := make(map[string]int)
	//   2. Split text: strings.Fields(text)
	words := strings.Fields(text)
	//   3. For each word:
	for _, word := range words { // Range over words
		word = strings.ToLower(word) // Convert word to lowercase
		word = strings.Map(func(r rune) rune { // Map over word
			if unicode.IsLetter(r) { // Check if letter
				return r // Return letter
			}
			return -1 // Return -1 to delete
		}, word)
		// Skip empty words (might be empty after removing all non-letters)
		if word == "" { // Check if empty
			continue // Skip to next word
		}
		counts[word]++ // Count word
	}
	//   4. Return counts
	return counts
}
