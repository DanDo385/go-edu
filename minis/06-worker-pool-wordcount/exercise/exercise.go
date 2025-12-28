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
	//
	//   Example:
	//     g, ctx := errgroup.WithContext(ctx)
	g, ctx := errgroup.WithContext(ctx) // Create errgroup and cancellable context

	jobs := make(chan string, workers) // Create jobs channel	
	results := make(chan map[string]int, workers) // Create results channel with buffered channel

	// TODO: Step 3 - Launch workers using errgroup
	//   Loop from i := 0 to workers-1:
	//     g.Go(func() error {
	//         // Worker implementation
	//     })
	//
	//   Inside worker:
	//     - Infinite loop: for { ... }
	//     - Select on ctx.Done() and jobs channel
	//     - Call fetchAndCount(ctx, url)
	//     - If error: return err (errgroup handles cancellation automatically!)
	//     - If success: send to results channel
	//
	//   Key points:
	//   - g.Go() automatically calls WaitGroup.Add(1) and defer Done()
	//   - Just return error - errgroup cancels context automatically
	//   - No error channel needed!
	for i := 0; i < workers; i++ {
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case url := <-jobs:
					counts, err := fetchAndCount(ctx, url)
				}
			}
		})
	}
	// TODO: Step 4 - Send jobs in separate goroutine
	//   Launch goroutine that:
	//     - Ranges over urls: for _, url := range urls
	//     - Checks ctx.Done() before sending each URL
	//     - Sends to jobs channel: jobs <- url
	//     - Closes jobs channel when done: close(jobs)
	//
	//   Example:
	//     go func() {
	//         defer close(jobs)
	//         for _, url := range urls {
	//             select {
	//             case <-ctx.Done():
	//                 return
	//             case jobs <- url:
	//             }
	//         }
	//     }()
	go func() { // Launch goroutine to send jobs
		defer close(jobs)          // Close jobs channel when done
		for _, url := range urls { // Range over urls
			select { // Check cancellation before sending each URL
			case <-ctx.Done(): // Stop if cancelled
				return // Exit if cancelled
			case jobs <- url: // Send URL to jobs channel
				// Successfully sent
			}
		}
	}()

	// TODO: Step 5 - Close results when workers finish
	//   Launch goroutine that:
	//     - Calls g.Wait() to wait for all workers
	//     - Closes results channel: close(results)
	//
	//   Example:
	//     go func() {
	//         g.Wait()
	//         close(results)
	//     }()
	go func() { // Launch goroutine to close results channel
		g.Wait()       // Wait for all workers to finish
		close(results) // Close results channel
	}()
	// TODO: Step 6 - Aggregate results
	//   Create map and range over results channel:
	//     finalCounts := make(map[string]int)
	//     for counts := range results {
	//         for word, count := range counts {
	//             finalCounts[word] += count
	//         }
	//     }
	finalCounts := make(map[string]int)
	for counts := range results { // Range over results channel
		for word, count := range counts { // Range over map
			finalCounts[word] += count // Accumulate counts for each word
		}
	}
	// TODO: Step 7 - Check for errors
	//   Call g.Wait() and check error:
	//     if err := g.Wait(); err != nil {
	//         return nil, err
	//     }
	if err := g.Wait(); err != nil { // Wait for all workers to finish
		return nil, err // Return error if any occurred
	}
	// TODO: Step 8 - Return results
	//   return finalCounts, nil
	return finalCounts, nil // Return final counts
	// Placeholder to prevent compilation errors
	// return nil, fmt.Errorf("not implemented")
}

/*
fetchAndCount fetches a URL and counts words.

ANALOGY: Package Delivery
-------------------------
Like ordering a package:
1. Create order (HTTP request) with tracking (context)
2. Wait for delivery (HTTP response)
3. Open package (read body)
4. Process contents (tokenize and count)
*/
func fetchAndCount(ctx context.Context, url string) (map[string]int, error) { // Fetch and count words from URL
	// TODO: Implement fetchAndCount
	//   1. Create HTTP request: http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil { // Error creating request
		return nil, err // Return error
	}
	//   2. Execute request: http.DefaultClient.Do(req) with error handling
	resp, err := http.DefaultClient.Do(req) // Execute request
	if err != nil { // Error executing request
		return nil, err // Return error
	}
	defer resp.Body.Close() // Close response body
	//   3. Check error and status code (must be 200 OK)
	if resp.StatusCode != http.StatusOK { // Status code is not 200 OK
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode) // Return error
	}
	//   4. Read body: io.ReadAll(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//   5. Don't forget: defer resp.Body.Close()
	defer resp.Body.Close()
	//   6. Convert to string and call tokenizeAndCount
	text := string(body)
	counts := tokenizeAndCount(text)
	//   7. Return counts and error
	return counts, nil
}

/*
tokenizeAndCount splits text into words and counts frequencies.

ANALOGY: Word Processor
----------------------
Like a word processor's word count:
1. Split text into words (by spaces)
2. Normalize (lowercase, remove punctuation)
3. Count frequencies
*/
func tokenizeAndCount(text string) map[string]int {
	// TODO: Implement tokenizeAndCount
	//   1. Create map: counts := make(map[string]int)
	counts := make(map[string]int)

	//   2. Split text: strings.Fields(text)
	words := strings.Fields(text)
	//   3. For each word:
	//      a. Lowercase: strings.ToLower(word)
	//      b. Remove non-letters: strings.Map(func(r rune) rune { ... }, word)
	//         - Use unicode.IsLetter(r) to check if letter
	//         - Return r to keep, -1 to delete
	//      c. Skip empty: if word == "" { continue }
	//      d. Count: counts[word]++
	for _, word := range words {
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
	//   4. Return counts
	return counts
}
