//go:build reference
// +build reference

/*
Problem: Fetch multiple URLs concurrently and count word frequencies

We need to implement:
1. Worker pool pattern for bounded concurrency
2. Concurrent HTTP fetching with proper error handling
3. Word tokenization and frequency counting
4. Graceful cancellation with context

Constraints:
- Fixed number of workers (bounded concurrency)
- Must handle errors from any worker
- Must cancel all work if one worker fails
- Aggregate results from all workers

Time/Space Complexity:
- Time: O(n) where n = total words across all URLs (concurrent fetching)
- Space: O(u) where u = number of unique words across all URLs

Why Go is well-suited:
- Goroutines: Lightweight concurrent execution
- Channels: Safe communication between goroutines
- Context: Built-in cancellation and timeout support
- sync.WaitGroup: Simple goroutine coordination
- errgroup: Simplified error handling for goroutines

DEBUGGING THIS FILE:
==================
This solution is instrumented with extensive debugging comments to teach you
how to use Go's debugger (dlv) and VS Code's debugging features.

Key debugging concepts covered:
1. Setting breakpoints in concurrent goroutines
2. Watching channel operations and synchronization
3. Understanding worker pool patterns
4. Inspecting goroutine state and coordination
5. Using the Debug Console for concurrent debugging
6. Understanding context cancellation propagation
*/

package workerpoolwordcount

import (
	"context"  // Context for cancellation and timeouts
	"fmt"      // Error formatting
	"io"       // Reading HTTP response bodies
	"net/http" // HTTP client
	"strings"  // String manipulation (Fields, ToLower, Map)
	"sync"     // WaitGroup for goroutine coordination
	"unicode"  // Character classification (IsLetter)

	"golang.org/x/sync/errgroup" // Simplified error handling and goroutine coordination
)

/*
WordCount fetches URLs concurrently using a worker pool pattern.

ANALOGY: Assembly Line Factory
------------------------------
Imagine a factory assembly line:
- Jobs channel = Conveyor belt carrying work items (URLs)
- Workers = Factory workers (goroutines) processing items
- Results channel = Finished products conveyor belt
- WaitGroup = Factory manager tracking how many workers are still working
- Context = Emergency stop button (cancels all work if something breaks)

The factory has exactly N workers (bounded concurrency) to prevent
overwhelming the system with too many simultaneous operations.

DEBUGGING WORKFLOW:
===================
1. Set breakpoints in worker goroutines and main coordination code
2. Use Goroutines panel to see all running goroutines
3. Switch between goroutines to see their current state
4. Watch channel operations (sends/receives)
5. Observe context cancellation propagation
*/
func WordCount(ctx context.Context, urls []string, workers int) (map[string]int, error) {
	// ============================================================================
	// PART 1: MANUAL APPROACH (Understanding What Happens Under the Hood)
	// ============================================================================
	// This shows the manual way using WaitGroup, channels, and context.
	// We'll see the errgroup version later which simplifies this significantly.

	// Step 1: Create communication channels (like conveyor belts)
	jobs := make(chan string, workers)            // Channel for URLs to process (buffered)
	results := make(chan map[string]int, workers) // Channel for word counts (buffered)
	errCh := make(chan error, 1)                  // Channel for first error (size 1)

	// Step 2: Create cancellable context (like an emergency stop button)
	ctx, cancel := context.WithCancel(ctx) // Create child context with cancel function
	defer cancel()                         // Always cleanup when function exits

	// Step 3: Create WaitGroup (like a counter tracking active workers)
	var wg sync.WaitGroup // WaitGroup is a struct with an internal counter

	// Step 4: Launch worker goroutines (like hiring factory workers)
	for i := 0; i < workers; i++ { // Loop: i = 0, 1, 2, ..., workers-1
		wg.Add(1) // Increment counter BEFORE starting goroutine (critical!)

		go func(workerID int) { // Launch goroutine (runs concurrently)
			// CRITICAL: Pass i as parameter to avoid closure bug!
			// Without (i), all goroutines would see the same variable value

			defer wg.Done() // Decrement counter when goroutine exits (even on panic)

			// Worker's main loop: keep processing jobs until told to stop
			for { // Infinite loop (exits via return statements)
				select { // Select waits for ONE of these cases to be ready
				case <-ctx.Done(): // Check if context cancelled (emergency stop pressed)
					return // Exit immediately if cancelled

				case url, ok := <-jobs: // Try to receive a job from channel
					// Two-value receive: url = value, ok = channel still open?
					if !ok { // ok=false means channel closed (no more work)
						return // Exit: no more jobs coming
					}

					// Process the job: fetch URL and count words
					counts, err := fetchAndCount(ctx, url) // Call helper function
					if err != nil {                        // Check for errors
						// Error occurred - send it (non-blocking) and stop everything
						select { // Non-blocking send pattern
						case errCh <- fmt.Errorf("fetching %s: %w", url, err): // Try to send error
							cancel() // Cancel context (stops all other workers)
						default: // If errCh full, ignore (first error already recorded)
						}
						return // Exit this worker
					}

					// Success! Send results (but check cancellation first)
					select { // Interruptible send
					case <-ctx.Done(): // Check if cancelled while processing
						return // Exit if cancelled
					case results <- counts: // Send word counts to results channel
						// Successfully sent, continue to next iteration
					}
				}
			}
		}(i) // Pass loop variable i as parameter (prevents closure bug)
	}

	// Step 5: Send jobs in separate goroutine (like loading items onto conveyor belt)
	go func() { // Launch goroutine to send jobs
		defer close(jobs) // Close channel when done (signals "no more jobs")

		for _, url := range urls { // Range over URLs slice
			select { // Check cancellation before sending each URL
			case <-ctx.Done(): // Stop if cancelled
				return
			case jobs <- url: // Send URL to jobs channel
				// Successfully sent
			}
		}
	}()

	// Step 6: Close results channel when all workers finish
	go func() { // Launch goroutine to coordinate shutdown
		wg.Wait()      // Wait for all workers to finish (blocks until counter=0)
		close(results) // Close channel (signals "no more results coming")
	}()

	// Step 7: Aggregate results (like collecting finished products)
	finalCounts := make(map[string]int) // Create map to store final word counts

	for counts := range results { // Range over channel (exits when channel closed)
		for word, count := range counts { // Range over map (word=key, count=value)
			finalCounts[word] += count // Accumulate: add count to existing total
		}
	}

	// Step 8: Check for errors (non-blocking receive)
	select { // Check if error occurred
	case err := <-errCh: // Try to receive error
		return nil, err // Return error if found
	default: // No error available
		// Continue to success return
	}

	return finalCounts, nil // Return aggregated word counts
}

/*
WordCountWithErrGroup demonstrates the errgroup approach - much simpler!

ANALOGY: Smart Factory Manager
------------------------------
errgroup is like a smart factory manager who:
- Automatically tracks workers (no manual WaitGroup.Add/Done)
- Automatically stops everything on first error (no manual error channel)
- Automatically cleans up (no manual defer cancel)

The code is ~40% shorter and less error-prone!
*/
func WordCountWithErrGroup(ctx context.Context, urls []string, workers int) (map[string]int, error) {
	// Step 1: Create errgroup (combines WaitGroup + Context + Error handling)
	g, ctx := errgroup.WithContext(ctx) // Returns group and cancellable context

	// Step 2: Create channels (still needed for worker pool pattern)
	jobs := make(chan string, workers)            // Jobs channel
	results := make(chan map[string]int, workers) // Results channel
	// No error channel needed! errgroup handles errors automatically

	// Step 3: Launch workers using errgroup
	for i := 0; i < workers; i++ { // Loop over worker count
		g.Go(func() error { // g.Go() automatically handles WaitGroup.Add/Done!
			// Worker loop
			for { // Infinite loop
				select { // Select on channels
				case <-ctx.Done(): // Check cancellation
					return ctx.Err() // Return error if cancelled

				case url, ok := <-jobs: // Receive job
					if !ok { // Channel closed
						return nil // Exit normally (no error)
					}

					// Process job
					counts, err := fetchAndCount(ctx, url) // Fetch and count
					if err != nil {                        // Check error
						return fmt.Errorf("fetching %s: %w", url, err) // Just return error!
						// errgroup automatically cancels context and stores first error
					}

					// Send results
					select { // Interruptible send
					case <-ctx.Done(): // Check cancellation
						return ctx.Err()
					case results <- counts: // Send result
						// Success
					}
				}
			}
		}) // No need to pass i - errgroup handles everything!
	}

	// Step 4: Send jobs
	go func() { // Launch goroutine
		defer close(jobs)          // Close when done
		for _, url := range urls { // Range over URLs
			select { // Check cancellation
			case <-ctx.Done(): // Stop if cancelled
				return
			case jobs <- url: // Send URL
				// Sent
			}
		}
	}()

	// Step 5: Close results when workers done
	go func() { // Launch goroutine
		_ = g.Wait()   // Wait for workers (errgroup handles WaitGroup internally)
		close(results) // Close channel
	}()

	// Step 6: Aggregate results
	finalCounts := make(map[string]int) // Create map
	for counts := range results {       // Range over channel
		for word, count := range counts { // Range over map
			finalCounts[word] += count // Accumulate
		}
	}

	// Step 7: Check for errors (errgroup makes this simple!)
	if err := g.Wait(); err != nil { // Wait() returns first error (or nil)
		return nil, err // Return error if any occurred
	}

	return finalCounts, nil // Return results
}

/*
fetchAndCount fetches a URL and counts words in the response.

ANALOGY: Package Delivery
-------------------------
Think of this like ordering a package:
1. Create order (HTTP request) with tracking (context)
2. Wait for delivery (HTTP response)
3. Open package (read body)
4. Process contents (tokenize and count words)

If you cancel the order (context cancelled), delivery stops immediately.
*/
func fetchAndCount(ctx context.Context, url string) (map[string]int, error) {
	// Create HTTP request with context (allows cancellation)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil) // GET request
	if err != nil {                                                       // Check error
		return nil, err // Return error if request creation failed
	}

	// Execute HTTP request (may take seconds)
	resp, err := http.DefaultClient.Do(req) // Send request, get response
	if err != nil {                         // Check error
		return nil, err // Return error if request failed
	}
	defer resp.Body.Close() // Always close body (releases resources)

	// Check status code (must be 200 OK)
	if resp.StatusCode != http.StatusOK { // Check if successful
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode) // Return error
	}

	// Read entire response body into memory
	body, err := io.ReadAll(resp.Body) // Read all bytes from body
	if err != nil {                    // Check error
		return nil, err // Return error if read failed
	}

	// Convert bytes to string and tokenize
	return tokenizeAndCount(string(body)), nil // Call tokenize function
}

/*
tokenizeAndCount splits text into words and counts frequencies.

ANALOGY: Word Processor
----------------------
Like a word processor's word count feature:
1. Split text into words (by spaces)
2. Normalize each word (lowercase, remove punctuation)
3. Count how many times each word appears

Example: "Hello, world! Hello Go" → {"hello": 2, "world": 1, "go": 1}
*/
func tokenizeAndCount(text string) map[string]int {
	counts := make(map[string]int) // Create map to store word counts

	// Split text into words (by whitespace: spaces, tabs, newlines)
	for _, word := range strings.Fields(text) { // Range over words slice
		word = strings.ToLower(word) // Convert to lowercase ("Hello" → "hello")

		// Remove non-letter characters (punctuation, numbers, etc.)
		word = strings.Map(func(r rune) rune { // Apply function to each character
			if unicode.IsLetter(r) { // Check if character is a letter
				return r // Keep letter
			}
			return -1 // Delete non-letter (returns -1 to remove character)
		}, word)

		// Skip empty words (might be empty after removing all non-letters)
		if word == "" { // Check if empty
			continue // Skip to next word
		}

		counts[word]++ // Increment count for this word (creates entry if doesn't exist)
	}

	return counts // Return word frequency map
}
