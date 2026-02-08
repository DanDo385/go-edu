//go:build reference

package workerpoolwordcount

/*
Reference Solution - Worker Pool for Parallel URL Word Count
==========================================================

This file demonstrates the classic worker-pool concurrency pattern: a fixed
number of goroutines process items from a job channel and send results to an
output channel. Applied to word counting across multiple URLs fetched in parallel.

This connects to the Go concurrency model:
- Goroutines: lightweight threads, spawned with `go fn()`
- Channels: typed conduits for communication between goroutines
- Context: cancellation and timeout propagation across goroutines
- sync.WaitGroup: coordination for "wait until all workers finish"

The exercise teaches:
- Worker pool architecture: bounded parallelism to avoid overwhelming resources
- Fan-out: distribute URLs to workers via channel
- Fan-in: merge results from workers into single map
- Context-aware cancellation: workers exit when context is done
- Channel lifecycle: who closes what, and when

Teaching notes:
- Memory/ownership: each worker owns its local counts map; we merge into a single
  map in the main goroutine. No shared mutable state between workers.
- Invariants: Channel ownership (per .cursorrules). Only the SENDER closes a channel.
  We have one producer (sends URLs to jobs) — it closes jobs. Workers are receivers.
  Results: workers send; a dedicated goroutine closes results after wg.Wait(). The
  main goroutine ranges over results until close. Never close a channel you didn't
  create or that has multiple senders — that causes panic.
- Error surfaces: first worker error cancels context; all workers exit; we return
  that error. Alternative: collect all errors, return partial results.
*/

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"unicode"
)

/*
WordCount - Parallel Word Count Across URLs

Fetches each URL, tokenizes content into words, counts occurrences, and merges
counts across all URLs. Uses a worker pool for bounded parallelism.

Parameters:
  - ctx: cancellation/timeout context; when done, all workers exit
  - urls: list of HTTP URLs to fetch and process
  - workers: number of concurrent workers (defaults to 1 if <= 0)

Returns merged word counts, or error on first fetch/parse failure.
*/
func WordCount(ctx context.Context, urls []string, workers int) (map[string]int, error) {
	// Guard against invalid worker count
	if workers <= 0 {
		workers = 1
	}

	// result carries either counts + nil, or nil + error
	type result struct {
		counts map[string]int
		err    error
	}

	// Derive cancellable context so we can cancel on first error
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Unbuffered channels: send blocks until receive
	// jobs: URLs to process; results: per-URL counts or errors
	jobs := make(chan string)
	results := make(chan result)

	var wg sync.WaitGroup

	// Worker goroutine: pulls URLs from jobs, fetches, counts, sends result
	workerFn := func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case url, ok := <-jobs:
				if !ok {
					// Channel closed - no more jobs, exit
					return
				}
				counts, err := fetchAndCount(ctx, url)
				// Send result, but respect cancellation (don't block forever if ctx done)
				select {
				case results <- result{counts: counts, err: err}:
				case <-ctx.Done():
					return
				}
			}
		}
	}

	// Start worker pool
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go workerFn()
	}

	// Producer: send URLs to jobs channel, close when done
	go func() {
		defer close(jobs)
		for _, url := range urls {
			select {
			case jobs <- url:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Closer: wait for all workers to finish, then close results
	// This allows the range over results to terminate
	go func() {
		wg.Wait()
		close(results)
	}()

	// Merge results from all workers
	merged := make(map[string]int)
	for res := range results {
		if res.err != nil {
			cancel()
			return nil, res.err
		}
		for word, count := range res.counts {
			merged[word] += count
		}
	}

	// If context was canceled (timeout, etc.) vs our explicit cancel, surface it
	if err := ctx.Err(); err != nil && err != context.Canceled {
		return nil, err
	}
	return merged, nil
}

/*
fetchAndCount - Fetch URL and Count Words

Creates an HTTP GET request with context (for cancellation/timeout), executes it,
validates status, reads body, and tokenizes/counts words. Returns (counts, nil)
or (nil, err) on failure.
*/
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

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return tokenizeAndCount(string(body)), nil
}

/*
tokenizeAndCount - Split Text into Words and Count

Uses bufio.Scanner with ScanWords to tokenize. ScanWords splits on whitespace.
Each token is normalized (lowercase, alphanumeric only) and counted.
*/
func tokenizeAndCount(text string) map[string]int {
	counts := make(map[string]int)
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		// ScanWords yields each whitespace-separated token
		token := normalizeToken(scanner.Text())
		if token == "" {
			continue
		}
		counts[token]++
	}
	return counts
}

/*
normalizeToken - Normalize Word for Consistent Counting

Keeps only letters and digits, converts to lowercase.
"Hello!" -> "hello", "CAN'T" -> "cant" (apostrophe stripped).
Empty string returned if no valid chars (e.g. "---").
*/
func normalizeToken(word string) string {
	var b strings.Builder
	for _, r := range word {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}
