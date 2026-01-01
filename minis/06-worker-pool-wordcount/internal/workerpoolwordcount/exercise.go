//go:build !solution && !reference

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


func WordCount(ctx context.Context, urls []string, workers int) (map[string]int, error) {
	// ============================================================================
	// PART 1: MANUAL APPROACH (Understanding What Happens Under the Hood)
	// TODO: Implement

	// ============================================================================
	// This shows the manual way using WaitGroup, channels, and context.
	// TODO: Implement

	panic("unimplemented")
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
	// TODO: Implement this function
	panic("unimplemented")
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
	// TODO: Implement this function
	panic("unimplemented")
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
	// TODO: Implement this function
	panic("unimplemented")
}
