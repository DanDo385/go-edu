package main

import (
	"context"           // Context for cancellation and timeouts
	"fmt"               // Formatted I/O
	"log"               // Logging
	"net/http"          // HTTP client and server
	"net/http/httptest" // Testing utilities for HTTP servers
	"os"                // Operating system interface
	"strconv"           // String to int conversion
	"time"              // Time operations

	// Import the exercise package from parent directory
	"github.com/example/go-10x-minis/minis/06-worker-pool-wordcount/internal/workerpoolwordcount"
)

/*
===================================================================
Worker Pool: Concurrent Word Count
===================================================================

This program demonstrates concurrency patterns in Go by implementing
a worker pool that fetches URLs concurrently and counts word frequencies.

USAGE:
    go run ./cmd/app/main.go
    go run ./cmd/app/main.go 5
    go run ./cmd/app/main.go 3 5s
    go run ./cmd/app/main.go 10 10s 20

Arguments:
    [workers] - Number of worker goroutines (default: 3)
    [timeout] - Overall timeout, e.g., "10s", "5m" (default: "10s")
    [top]     - Number of top words to display (default: 10)
    [servers] - Number of test HTTP servers (default: 5)

Examples:
    go run ./cmd/app/main.go
    go run ./cmd/app/main.go 5
    go run ./cmd/app/main.go 3 5s
    go run ./cmd/app/main.go 10 10s 20 3

DEBUGGING:
- Set breakpoints at "// BREAKPOINT:" comments
- Use F5 to start debugging (select "Debug: Run main.go")
- Step Into (F11) to enter workerpoolwordcount.WordCount function
- Watch Variables panel to see goroutine coordination
- See /RUN_DEBUG.md for comprehensive debugging guide

PACKAGE STRUCTURE:
- This file (cmd/app/main.go) is package main
- It imports the exercise package from internal/exercise
- workerpoolwordcount.go, solution.go are in package exercise
- Functions must be exported (capitalized) to be called from main
*/

func main() {
	// ============================================================================
	// STEP 1: Parse Command-Line Arguments
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to inspect command-line arguments
	// DEBUG: os.Args[0] is the program name
	// DEBUG: os.Args[1] is the first argument (if provided)
	// DEBUG: In Variables panel, expand 'os.Args' to see all arguments

	// Default values
	workers := 3
	timeout := 10 * time.Second
	top := 10
	servers := 5

	// Override defaults if arguments provided
	if len(os.Args) > 1 {
		if w, err := strconv.Atoi(os.Args[1]); err == nil {
			workers = w
		}
	}
	if len(os.Args) > 2 {
		if d, err := time.ParseDuration(os.Args[2]); err == nil {
			timeout = d
		}
	}
	if len(os.Args) > 3 {
		if t, err := strconv.Atoi(os.Args[3]); err == nil {
			top = t
		}
	}
	if len(os.Args) > 4 {
		if s, err := strconv.Atoi(os.Args[4]); err == nil {
			servers = s
		}
	}

	// BREAKPOINT: Set breakpoint here to inspect parsed values in Variables panel
	// DEBUG: After parsing, variables now hold actual command-line values
	// DEBUG: Watch 'workers', 'timeout', 'top', 'servers'

	// ============================================================================
	// STEP 2: Create Test HTTP Servers
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before creating servers
	// DEBUG: Step Into (F11) on createTestServers() to see httptest.NewServer creation
	// DEBUG: Watch how mock HTTP servers are initialized

	testServers := createTestServers(servers)

	// BREAKPOINT: Set breakpoint here after creating servers
	// DEBUG: Inspect 'testServers' slice - expand to see all *httptest.Server pointers
	// DEBUG: Each server has URL field containing address like "http://127.0.0.1:12345"

	// defer ensures servers are closed when main() exits (cleanup)
	defer closeServers(testServers)

	// ============================================================================
	// STEP 3: Extract URLs from Test Servers
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before extracting URLs
	// DEBUG: We're creating a slice to hold server URLs
	// DEBUG: Watch 'urls' slice being built from server URL fields

	urls := make([]string, len(testServers))
	for i, srv := range testServers {
		// BREAKPOINT: Set breakpoint here inside URL extraction loop
		// DEBUG: Watch 'i' increment and 'srv.URL' change with each iteration
		// DEBUG: See how we're copying URL strings to urls slice

		urls[i] = srv.URL
	}

	// BREAKPOINT: Set breakpoint here after URL extraction
	// DEBUG: Inspect 'urls' slice - expand to see all URL strings
	// DEBUG: len(urls) should equal number of test servers

	fmt.Printf("=== Worker Pool Word Count Demo ===\n")
	fmt.Printf("Fetching %d URLs with %d workers...\n", len(urls), workers)
	fmt.Printf("Timeout: %v\n", timeout)
	fmt.Println()

	// ============================================================================
	// STEP 4: Create Context with Timeout
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before creating context
	// DEBUG: context.WithTimeout creates a cancellable context
	// DEBUG: If timeout expires, ctx.Done() channel will close

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	// BREAKPOINT: Set breakpoint here after creating context
	// DEBUG: Inspect 'ctx' - see Deadline field (when timeout expires)
	// DEBUG: 'cancel' is a function to manually cancel context

	defer cancel() // Cleanup: cancel context when main() exits

	// ============================================================================
	// STEP 5: Run WordCount with Worker Pool
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before starting word count
	// DEBUG: Step Into (F11) on workerpoolwordcount.WordCount() to see worker pool implementation
	// DEBUG: Watch how goroutines are created and coordinated with channels

	start := time.Now()

	counts, err := workerpoolwordcount.WordCount(ctx, urls, workers)

	// BREAKPOINT: Set breakpoint here after word count completes
	// DEBUG: Inspect 'counts' map - expand to see word → count mappings
	// DEBUG: Check 'err' - should be nil if successful
	// DEBUG: If err != nil, operation timed out or failed

	if err != nil {
		log.Fatalf("WordCount failed: %v", err)
	}

	elapsed := time.Since(start)

	// BREAKPOINT: Set breakpoint here after timing calculation
	// DEBUG: Inspect 'elapsed' - time.Duration showing execution time
	// DEBUG: Compare elapsed vs timeout to see if we completed in time

	// ============================================================================
	// STEP 6: Sort Results by Frequency
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before sorting
	// DEBUG: Maps are unordered - we need to convert to slice for sorting
	// DEBUG: Watch how we build pairs slice from counts map

	type pair struct {
		word  string
		count int
	}

	var pairs []pair

	// BREAKPOINT: Set breakpoint here before map iteration
	// DEBUG: Expand 'counts' map in Variables panel
	// DEBUG: We're about to convert map to slice

	for word, count := range counts {
		// BREAKPOINT: Set breakpoint here inside map iteration
		// DEBUG: Watch 'word' and 'count' change with each iteration
		// DEBUG: See how we append pair{} struct to pairs slice

		pairs = append(pairs, pair{word, count})
	}

	// BREAKPOINT: Set breakpoint here after building pairs slice
	// DEBUG: Inspect 'pairs' slice - should have same length as counts map
	// DEBUG: Expand pairs to see all word-count pairs

	// ============================================================================
	// STEP 7: Partial Selection Sort (Top N Only)
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before sorting
	// DEBUG: We're doing partial selection sort - only sort top N items
	// DEBUG: Watch how maxIdx is found and swapped in each iteration

	topN := top
	if topN > len(pairs) {
		topN = len(pairs)
	}

	// BREAKPOINT: Set breakpoint here before sort loop
	// DEBUG: Watch 'topN' - number of iterations
	// DEBUG: Selection sort: find max, swap to front, repeat

	for i := 0; i < topN; i++ {
		maxIdx := i

		// BREAKPOINT: Set breakpoint here inside outer sort loop
		// DEBUG: Watch 'i' increment with each iteration (0, 1, 2, ...)
		// DEBUG: 'maxIdx' starts as i, will be updated if higher count found

		for j := i + 1; j < len(pairs); j++ {
			// BREAKPOINT: Set breakpoint here inside inner sort loop
			// DEBUG: Watch 'j' scan through remaining unsorted portion
			// DEBUG: Watch maxIdx update when higher count is found

			if pairs[j].count > pairs[maxIdx].count {
				maxIdx = j
			}
		}

		// BREAKPOINT: Set breakpoint here before swap
		// DEBUG: Watch 'maxIdx' - index of element with highest count
		// DEBUG: Step Over (F10) to see swap occur

		pairs[i], pairs[maxIdx] = pairs[maxIdx], pairs[i]

		// BREAKPOINT: Set breakpoint here after swap
		// DEBUG: Inspect pairs[i] - should have highest count in remaining portion
		// DEBUG: Top of slice is now sorted (highest counts first)
	}

	// ============================================================================
	// STEP 8: Display Results
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before display
	// DEBUG: 'pairs' is now sorted by count (descending)
	// DEBUG: We'll print top N items

	fmt.Printf("Top %d words:\n", topN)

	for i := 0; i < topN; i++ {
		// BREAKPOINT: Set breakpoint here inside display loop
		// DEBUG: Watch 'i' increment and pairs[i] change
		// DEBUG: See word and count for each top result

		fmt.Printf("%2d. %-15s %d\n", i+1, pairs[i].word, pairs[i].count)
	}

	// ============================================================================
	// STEP 9: Print Summary Statistics
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before summary
	// DEBUG: Inspect len(counts) - total unique words
	// DEBUG: Inspect elapsed - total execution time

	fmt.Println()
	fmt.Printf("Total unique words: %d\n", len(counts))
	fmt.Printf("Completed in: %v\n", elapsed)
	fmt.Printf("Worker goroutines: %d\n", workers)

	// ============================================================================
	// STEP 10: Usage Note about Solution
	// ============================================================================
	if *useSolution {
		fmt.Println()
		fmt.Println("NOTE: To use solution.go, build/run with -tags=solution:")
		fmt.Println("  go run -tags=solution ./cmd/worker-pool-wordcount/main.go")
		fmt.Println("  go test -tags=solution")
	}

	// ============================================================================
	// DEBUGGING TIPS
	// ============================================================================
	//
	// KEY DEBUGGING TECHNIQUES FOR CONCURRENCY:
	//
	// 1. Set breakpoints at "// BREAKPOINT:" comments
	// 2. Press F5 to start debugging
	// 3. Use F10 (Step Over) to execute line by line
	// 4. Use F11 (Step Into) to enter workerpoolwordcount.WordCount function
	// 5. Watch Variables panel to see goroutine coordination
	// 6. Use Goroutines panel to see all running goroutines
	// 7. Add watch expressions: len(counts), elapsed, *workers
	//
	// UNDERSTANDING GOROUTINES:
	// - Step Into (F11) workerpoolwordcount.WordCount to see worker pool
	// - Watch Goroutines panel to see workers spawn
	// - Each worker is a goroutine running concurrently
	// - Channels coordinate communication between goroutines
	//
	// UNDERSTANDING CHANNELS:
	// - Channels are typed conduits for sending/receiving values
	// - Sending: ch <- value (blocks if channel full)
	// - Receiving: val := <-ch (blocks if channel empty)
	// - Closing: close(ch) (signals no more values)
	//
	// UNDERSTANDING CONTEXT:
	// - ctx.Done() returns a channel that closes on timeout/cancel
	// - Use select with ctx.Done() to respect cancellation
	// - Step Into (F11) to see context handling in worker pool
	//
	// TESTING DIFFERENT WORKER COUNTS:
	// - Try: go run main.go -workers=1 (sequential)
	// - Try: go run main.go -workers=10 (high concurrency)
	// - Compare execution times to see concurrency benefits
	//
	// See /RUN_DEBUG.md for comprehensive debugging guide
}

// createTestServers creates N mock HTTP servers for testing
func createTestServers(n int) []*httptest.Server {
	// BREAKPOINT: Set breakpoint here to see server creation
	// DEBUG: Watch 'n' - number of servers to create
	// DEBUG: We're creating mock HTTP servers that return text

	content := []string{
		"Go is a great programming language for building scalable systems",
		"Concurrency in Go is simple and powerful with goroutines and channels",
		"Go programming makes it easy to build reliable and efficient software",
		"The Go standard library is comprehensive and well-designed",
		"Building web services with Go is straightforward and productive",
	}

	servers := make([]*httptest.Server, n)

	for i := 0; i < n; i++ {
		// BREAKPOINT: Set breakpoint here inside server creation loop
		// DEBUG: Watch 'i' increment with each server
		// DEBUG: Each server will return content[i % len(content)]

		text := content[i%len(content)]
		servers[i] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, text)
		}))
	}

	// BREAKPOINT: Set breakpoint here after creating servers
	// DEBUG: Inspect 'servers' slice - see all httptest.Server pointers
	// DEBUG: Each server is listening on a random port

	return servers
}

// closeServers closes all test servers to free resources
func closeServers(servers []*httptest.Server) {
	// BREAKPOINT: Set breakpoint here to see cleanup
	// DEBUG: This runs on defer - even if main() panics
	// DEBUG: Watch each server being closed

	for _, srv := range servers {
		// BREAKPOINT: Set breakpoint here inside close loop
		// DEBUG: Watch 'srv' change with each iteration
		// DEBUG: srv.Close() releases network resources

		srv.Close()
	}
}
