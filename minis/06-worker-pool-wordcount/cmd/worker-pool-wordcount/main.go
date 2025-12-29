// Package main declares this as an executable program (not a library)
// When you run "go run" or "go build", Go looks for a main package with a main() function
package main

import (
	"context"           // Context for cancellation and timeouts (allows stopping operations)
	"fmt"               // Formatted I/O (Printf, Println, Sprintf, etc.)
	"log"               // Logging functions (log.Fatalf exits program on error)
	"net/http"          // HTTP client and server functionality
	"net/http/httptest" // Testing utilities for HTTP servers (creates mock servers)
	"time"              // Time operations (Duration, Now, Since, etc.)

	// Import our exercise package containing the WordCount function
	"github.com/example/go-10x-minis/minis/06-worker-pool-wordcount/exercise"
)

/*
main is the entry point of the program.
When you run this program, Go automatically calls main().

ANALOGY: Restaurant Kitchen
---------------------------
Think of this like a restaurant kitchen:
1. Set up test servers = Prepare ingredients (httptest servers)
2. Extract URLs = List what needs to be cooked
3. Create context = Set a timer (10 second timeout)
4. Run WordCount = Cook all dishes concurrently (worker pool)
5. Display results = Present the finished dishes (sorted word counts)
*/
func main() {
	// ============================================================================
	// STEP 1: Create Test HTTP Servers
	// ============================================================================
	// httptest.NewServer creates mock HTTP servers for testing
	// These simulate real websites without needing internet access
	// Each server returns a simple text response when requested
	servers := createTestServers() // Returns slice of *httptest.Server pointers

	// defer schedules this function to run when main() exits (even on panic)
	// This ensures servers are always cleaned up, preventing resource leaks
	// defer executes in reverse order (LIFO - Last In, First Out)
	defer closeServers(servers)

	// ============================================================================
	// STEP 2: Extract URLs from Test Servers
	// ============================================================================
	// make([]string, len(servers)) creates a slice with capacity = len(servers)
	// A slice is a dynamic array - it can grow/shrink, unlike fixed-size arrays
	// Syntax: make(type, length) - creates slice with specified length
	urls := make([]string, len(servers)) // Create slice of strings with length = server count

	// Range loop iterates over slice, providing index (i) and value (srv)
	// Syntax: for index, value := range collection
	// srv is a pointer to httptest.Server (*httptest.Server)
	for i, srv := range servers { // Loop: i = 0, 1, 2, ...; srv = servers[0], servers[1], ...
		urls[i] = srv.URL // srv.URL is a string containing the server's URL (e.g., "http://127.0.0.1:54321")
	}

	// Printf formats and prints a string
	// %d = decimal integer, \n = newline
	// This prints: "Fetching 5 URLs with 3 workers...\n\n"
	fmt.Printf("Fetching %d URLs with 3 workers...\n\n", len(urls))

	// ============================================================================
	// STEP 3: Create Context with Timeout
	// ============================================================================
	// context.Background() creates a root context (no parent, never cancelled)
	// context.WithTimeout creates a child context that cancels after specified duration
	// Returns: ctx (the cancellable context) and cancel (function to cancel manually)
	// Syntax: context.WithTimeout(parent, duration)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // 10 second timeout

	// defer cancel() ensures context is cancelled when main() exits
	// This releases resources and stops any ongoing operations
	defer cancel() // Cleanup: cancel context when function exits

	// ============================================================================
	// STEP 4: Run WordCount with Worker Pool
	// ============================================================================
	// time.Now() captures current time (used to measure execution duration)
	start := time.Now() // Record start time

	// Call WordCount function from exercise package
	// Parameters:
	//   - ctx: cancellable context (allows timeout/cancellation)
	//   - urls: slice of URLs to fetch
	//   - 3: number of worker goroutines (bounded concurrency)
	// Returns:
	//   - counts: map[string]int (word -> count mapping)
	//   - err: error (nil if successful)
	counts, err := exercise.WordCount(ctx, urls, 3) // Execute worker pool wordcount

	// Error handling: if WordCount returned an error, log it and exit
	// log.Fatalf prints formatted error and calls os.Exit(1)
	// %v = default format (works for any type)
	if err != nil { // Check if error occurred
		log.Fatalf("WordCount failed: %v", err) // Print error and exit program
	}

	// time.Since(start) calculates elapsed time since start
	// Returns time.Duration (e.g., "2.5s", "150ms")
	elapsed := time.Since(start) // Calculate execution time

	// ============================================================================
	// STEP 5: Display Results - Sort and Print Top 10 Words
	// ============================================================================
	fmt.Println("Top 10 words:") // Print header

	// Define a struct type (named type) for word-count pairs
	// struct is a collection of named fields (like a lightweight class)
	// Syntax: type Name struct { field1 type1; field2 type2 }
	type pair struct {
		word  string // Word text
		count int    // Word frequency count
	}

	// Declare a slice of pair structs (initially empty/nil)
	// Syntax: var name []type
	var pairs []pair // Slice to hold all word-count pairs

	// Range over map: for key, value := range map
	// Maps in Go are unordered (iteration order is random)
	// This converts map[string]int into []pair slice
	for word, count := range counts { // Iterate over word -> count mapping
		// append adds element(s) to slice, returns new slice
		// Syntax: append(slice, element1, element2, ...)
		// pair{word, count} creates a struct literal (initializes struct fields)
		pairs = append(pairs, pair{word, count}) // Add each word-count pair to slice
	}

	// ============================================================================
	// STEP 6: Sort Pairs by Count (Selection Sort - Top 10 Only)
	// ============================================================================
	// Selection sort: find max, swap to front, repeat
	// We only sort top 10 (partial sort) for efficiency
	// Loop condition: i < len(pairs) && i < 10
	//   - i < len(pairs): don't go beyond slice bounds
	//   - i < 10: only sort first 10 elements
	for i := 0; i < len(pairs) && i < 10; i++ { // Loop: i = 0, 1, 2, ..., 9 (max 10 iterations)
		maxIdx := i // Assume current position has maximum count

		// Find index of element with maximum count in remaining unsorted portion
		// Start from i+1 (next element) to end of slice
		for j := i + 1; j < len(pairs); j++ { // Loop: j = i+1, i+2, ..., len(pairs)-1
			// Compare counts: if pairs[j] has higher count than current max
			if pairs[j].count > pairs[maxIdx].count { // Check if this element has higher count
				maxIdx = j // Update max index to this position
			}
		}

		// Swap elements: move maximum to position i
		// Go supports multiple assignment: a, b = b, a (swaps values)
		// This is a common Go idiom for swapping
		pairs[i], pairs[maxIdx] = pairs[maxIdx], pairs[i] // Swap: move max to front

		// Print formatted output
		// Printf format specifiers:
		//   %2d = right-aligned integer, width 2 (e.g., " 1", "10")
		//   %-15s = left-aligned string, width 15 (e.g., "hello          ")
		//   %d = decimal integer
		//   \n = newline
		fmt.Printf("%2d. %-15s %d\n", i+1, pairs[i].word, pairs[i].count) // Print rank, word, count
	}

	// ============================================================================
	// STEP 7: Print Summary Statistics
	// ============================================================================
	// len(counts) returns number of key-value pairs in map (unique words)
	fmt.Printf("\nTotal unique words: %d\n", len(counts)) // Print total unique word count

	// %v = default format (for time.Duration, prints like "2.5s" or "150ms")
	fmt.Printf("Completed in: %v\n", elapsed) // Print execution time
}

/*
createTestServers creates mock HTTP servers for testing.

ANALOGY: Virtual Restaurants
-----------------------------
Think of these as virtual restaurants that serve text instead of food.
Each server responds to HTTP requests with a simple text message.

Returns: []*httptest.Server (slice of pointers to test servers)
*/
func createTestServers() []*httptest.Server {
	// Return a slice literal containing 5 httptest.Server instances
	// Syntax: []type{value1, value2, ...} creates a slice with initial values
	return []*httptest.Server{
		// httptest.NewServer creates a test HTTP server that listens on a random port
		// http.HandlerFunc converts a function to an http.Handler interface
		// The function signature must match: func(http.ResponseWriter, *http.Request)
		httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Handler function: called when server receives HTTP request
			// w = ResponseWriter (write response to client)
			// r = Request (contains request details - not used here)
			// fmt.Fprintln writes formatted string to writer, adds newline
			fmt.Fprintln(w, "Go is a great programming language for building scalable systems")
		})),

		// Second test server with different content
		httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Concurrency in Go is simple and powerful with goroutines and channels")
		})),

		// Third test server
		httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Go programming makes it easy to build reliable and efficient software")
		})),

		// Fourth test server
		httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "The Go standard library is comprehensive and well-designed")
		})),

		// Fifth test server
		httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Building web services with Go is straightforward and productive")
		})),
	}
}

/*
closeServers closes all test servers to free resources.

ANALOGY: Closing Restaurants
-----------------------------
Like closing all restaurants at the end of the day - stops accepting
new customers and releases resources (ports, sockets, etc.).

Parameters:
  - servers: slice of pointers to httptest.Server instances

Why this matters:
  - Servers hold network resources (ports, sockets)
  - Not closing them causes resource leaks
  - defer in main() ensures this always runs
*/
func closeServers(servers []*httptest.Server) {
	// Range over slice: for index, value := range slice
	// Using _ discards the index (we don't need it)
	// srv is each *httptest.Server pointer in the slice
	for _, srv := range servers { // Iterate over all servers
		// Close() stops the server and releases its resources
		// Safe to call multiple times (idempotent)
		srv.Close() // Shutdown server and free port/socket
	}
}
