//go:build solution
// +build solution

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

/*
WordCount fetches a list of URLs concurrently and aggregates word frequencies.

This function is a *textbook worker-pool pattern*:

  - A bounded number of worker goroutines (N = workers)
  - A jobs channel distributes work
  - A results channel collects outputs
  - A shared context coordinates cancellation
  - A WaitGroup coordinates shutdown

IMPORTANT CONCURRENCY RULE:
---------------------------
At any instant:
  - At most `workers` HTTP requests are in flight
  - At most `workers` goroutines are actively tokenizing text
  - No goroutine ever mutates shared state unsafely
*/
func WordCount(ctx context.Context, urls []string, workers int) (map[string]int, error) {

	/*
	   CHANNEL DESIGN
	   --------------

	   jobs:
	     - carries *units of work*
	     - here: a single URL (string)
	     - produced by main goroutine
	     - consumed by workers

	   results:
	     - carries *units of output*
	     - here: a per-URL word frequency map
	     - produced by workers
	     - consumed by main goroutine

	   errCh:
	     - carries exactly ONE error (the first one)
	     - buffered to size 1 so the first sender never blocks
	*/

	jobs := make(chan string, workers)
	results := make(chan map[string]int, workers)
	errCh := make(chan error, 1)

	/*
	   CONTEXT
	   -------

	   We derive a cancellable context from the caller’s context.

	   Why?
	   - If *any* worker fails, we want to:
	       * stop issuing new jobs
	       * abort in-flight HTTP requests
	       * make other workers exit quickly

	   Context cancellation is Go’s standard "broadcast shutdown signal".
	*/
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // Safety net: ensures cleanup on return

	/*
	   WAITGROUP
	   ---------

	   The WaitGroup tracks *only worker goroutines*.

	   Why not track:
	     - job sender goroutine?
	     - result aggregator?

	   Because:
	     - workers are the only goroutines whose lifetime is data-dependent
	     - once workers are done, results will stop naturally
	*/
	var wg sync.WaitGroup

	/*
	   WORKER SPAWN LOOP
	   -----------------

	   This loop runs SEQUENTIALLY.
	   Each iteration *launches* a goroutine.
	   The goroutines themselves run CONCURRENTLY.
	*/
	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()

			/*
			   Each worker runs this loop until:
			     - jobs channel is closed, OR
			     - context is cancelled
			*/
			for {
				select {

				// Case 1: Someone cancelled the context
				case <-ctx.Done():
					return

				// Case 2: Try to receive a job
				case url, ok := <-jobs:
					if !ok {
						// Channel closed → no more work
						return
					}

					/*
					   This worker now *owns* this URL.
					   No other worker will see it.
					*/
					counts, err := fetchAndCount(ctx, url)
					if err != nil {
						/*
						   ERROR PROPAGATION STRATEGY
						   --------------------------

						   1. Try to send the error (non-blocking)
						   2. Cancel the context
						   3. Exit immediately

						   Why non-blocking?
						   - Multiple workers might error
						   - Only the *first* error matters
						*/
						select {
						case errCh <- fmt.Errorf("fetching %s: %w", url, err):
							cancel()
						default:
						}
						return
					}

					/*
					   Send successful results unless cancelled mid-flight.
					*/
					select {
					case <-ctx.Done():
						return
					case results <- counts:
					}
				}
			}
		}(i)
	}

	/*
	   JOB PRODUCER GOROUTINE
	   ---------------------

	   This goroutine feeds URLs into the jobs channel.

	   It exists so:
	     - main goroutine can immediately start aggregating results
	     - workers and sender overlap in time
	*/
	go func() {
		for _, url := range urls {
			select {
			case <-ctx.Done():
				return
			case jobs <- url:
			}
		}
		// Closing jobs tells workers: "no more URLs are coming"
		close(jobs)
	}()

	/*
	   RESULTS CLOSER GOROUTINE
	   -----------------------

	   This goroutine waits until ALL workers finish,
	   then closes the results channel.

	   This is crucial:
	     - closing results is how the aggregator knows when to stop
	     - workers must NEVER close results themselves
	*/
	go func() {
		wg.Wait()
		close(results)
	}()

	/*
	   AGGREGATION (MAIN GOROUTINE)
	   ----------------------------

	   This loop is SEQUENTIAL.
	   Only one goroutine mutates finalCounts → no locks needed.
	*/
	finalCounts := make(map[string]int)

	for counts := range results {
		for word, count := range counts {
			finalCounts[word] += count
		}
	}

	/*
	   ERROR CHECK
	   -----------

	   At this point:
	     - all workers are done
	     - all results have been merged

	   If an error occurred, it MUST be in errCh.
	*/
	select {
	case err := <-errCh:
		return nil, err
	default:
	}

	return finalCounts, nil
}

/*
fetchAndCount performs:
  1. Context-aware HTTP GET
  2. Full body read
  3. Tokenization

Key idea:
---------
Passing ctx into NewRequestWithContext ensures that:
  cancel() → HTTP request aborted at the transport level
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return tokenizeAndCount(string(body)), nil
}

/*
tokenizeAndCount is CPU-bound, deterministic, and side-effect-free.

This makes it:
  - safe to run concurrently
  - easy to reason about
*/
func tokenizeAndCount(text string) map[string]int {
	counts := make(map[string]int)

	/*
	   strings.Fields splits on Unicode whitespace.
	   Example:
	     "hi,\nthere\tfriend" → ["hi,", "there", "friend"]
	*/
	for _, word := range strings.Fields(text) {
		word = strings.ToLower(word)

		/*
		   strings.Map runs this function on EACH RUNE.

		   Returning:
		     - the rune → keeps it
		     - -1        → deletes it

		   This avoids regex and handles Unicode correctly.
		*/
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
