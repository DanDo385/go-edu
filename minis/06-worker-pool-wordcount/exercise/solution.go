//go:build solution
// +build solution

/*
Problem: Concurrent URL fetching with word frequency counting using worker pools

We need to:
1. Fetch multiple URLs concurrently (but with bounded parallelism)
2. Extract and tokenize text from each response
3. Aggregate word counts across all pages
4. Handle errors gracefully (cancel all work if any fetch fails)
5. Respect context cancellation and timeouts

Constraints:
- Use exactly N worker goroutines (no unbounded concurrency)
- Normalize words to lowercase
- Remove punctuation (keep only letters)
- Cancel all work on first error
- Thread-safe aggregation

Time/Space Complexity:
- Time: O(n/w * f + n*t) where n=URLs, w=workers, f=fetch time, t=tokenize time
- Space: O(w + v) where w=workers (in-flight requests), v=unique vocabulary

Why Go is well-suited:
- Goroutines: Lightweight threads (2KB stack vs 2MB for OS threads)
- Channels: Type-safe, built-in communication primitive
- context.Context: Standardized cancellation across goroutines
- No callback hell: Sequential code within goroutines

Compared to other languages:
- Python: asyncio/threading, but GIL prevents true parallelism
- JavaScript: Promises are single-threaded (no true parallelism)
- Rust: tokio is powerful but complex (async/await, Send/Sync traits)
- Java: Thread pools are similar but more verbose (ExecutorService)
*/

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

// WordCount fetches URLs concurrently and returns word frequencies.
//
// Go Concepts Demonstrated:
// - Goroutines: Lightweight concurrent execution
// - Channels: Communication between goroutines
// - sync.WaitGroup: Waiting for goroutine completion
// - context.Context: Cancellation propagation
// - Error handling in concurrent code
//
// Architecture:
// 1. Main goroutine: Sends URLs to jobs channel, waits for completion
// 2. Worker goroutines: Fetch URLs, tokenize, send results
// 3. Aggregator goroutine: Merges word counts from all workers
// 4. Error channel: First error cancels context, stops all workers
//
// Three-Input Iteration Table:
//
// Input 1: 3 URLs, 2 workers (happy path)
//   Main: Send URL1 → jobs, URL2 → jobs, URL3 → jobs
//   Worker1: Fetch URL1 → tokenize → send counts to results
//   Worker2: Fetch URL2 → tokenize → send counts to results
//   Worker1: Fetch URL3 → tokenize → send counts to results
//   Aggregator: Merge all counts → final result
//   Result: Complete word frequency map
//
// Input 2: 0 URLs (edge case)
//   Main: Close jobs immediately
//   Workers: Exit (no work)
//   Aggregator: Return empty map
//   Result: {}, nil
//
// Input 3: One URL fails (error propagation)
//   Worker1: Fetch URL1 → error → send to errCh → cancel context
//   Worker2: Fetch URL2 → context cancelled → stop
//   Result: nil, error
func WordCount(ctx context.Context, urls []string, workers int) (map[string]int, error) {
	// Create channels
	// Buffer size = workers to avoid blocking when all workers are busy
	jobs := make(chan string, workers)       // URLs to fetch
	results := make(chan map[string]int, workers) // Word counts per URL
	errCh := make(chan error, 1)             // First error (buffered to avoid blocking)

	// Create a cancellable context
	// If any worker encounters an error, we'll cancel this context
	// This signals all other workers to stop immediately
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // Ensure cleanup

	// WaitGroup to track worker goroutines
	var wg sync.WaitGroup

	// Start worker goroutines
	// Each worker:
	// 1. Reads URLs from jobs channel
	// 2. Fetches the URL
	// 3. Tokenizes the response
	// 4. Sends word counts to results channel
	// 5. Exits when jobs channel is closed
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// Process jobs until channel is closed or context is cancelled
			for {
				select {
				case <-ctx.Done():
					// Context was cancelled (error occurred elsewhere)
					return
				case url, ok := <-jobs:
					if !ok {
						// Jobs channel closed, no more work
						return
					}

					// Fetch and process this URL
					counts, err := fetchAndCount(ctx, url)
					if err != nil {
						// Send error and cancel context
						// Use non-blocking send to avoid deadlock (errCh is buffered size 1)
						select {
						case errCh <- fmt.Errorf("fetching %s: %w", url, err):
							cancel() // Cancel context to stop all other workers
						default:
							// Error channel already has an error, ignore this one
						}
						return
					}

					// Send results
					select {
					case <-ctx.Done():
						return
					case results <- counts:
						// Successfully sent
					}
				}
			}
		}(i)
	}

	// Send all jobs
	// This runs concurrently with workers processing jobs
	go func() {
		for _, url := range urls {
			select {
			case <-ctx.Done():
				// Context cancelled, stop sending jobs
				return
			case jobs <- url:
				// Job sent successfully
			}
		}
		close(jobs) // Signal no more jobs
	}()

	// Close results channel when all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Aggregate results
	// This runs in the main goroutine
	finalCounts := make(map[string]int)
	for counts := range results {
		// Merge counts from this URL into final result
		for word, count := range counts {
			finalCounts[word] += count
		}
	}

	// Check for errors
	select {
	case err := <-errCh:
		return nil, err
	default:
		// No error
	}

	return finalCounts, nil
}

// fetchAndCount fetches a URL and returns word frequencies.
//
// Go Concepts Demonstrated:
// - http.NewRequestWithContext: Context-aware HTTP requests
// - io.ReadAll: Read entire response body
// - String processing: tokenization without regex
func fetchAndCount(ctx context.Context, url string) (map[string]int, error) {
	// Create HTTP request with context
	// This allows cancellation to propagate to the HTTP client
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Execute request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Tokenize and count words
	return tokenizeAndCount(string(body)), nil
}

// tokenizeAndCount extracts words and returns frequencies.
//
// Go Concepts Demonstrated:
// - strings.Fields: Split on whitespace
// - unicode.IsLetter: Proper character classification
// - Map accumulation pattern
func tokenizeAndCount(text string) map[string]int {
	counts := make(map[string]int)

	// Split on whitespace
	words := strings.Fields(text)

	for _, word := range words {
		// Normalize to lowercase
		word = strings.ToLower(word)

		// Remove non-letters
		// strings.Map applies a function to each rune
		word = strings.Map(func(r rune) rune {
			if unicode.IsLetter(r) {
				return r
			}
			return -1 // Drop this rune
		}, word)

		// Skip empty words
		if word == "" {
			continue
		}

		counts[word]++
	}

	return counts
}

/*
Alternatives & Trade-offs:

1. Unbounded goroutines (one per URL):
   for _, url := range urls {
       go func(u string) { ... }(url)
   }
   Pros: Simpler code (no worker pool)
   Cons: Resource exhaustion with 10k URLs (too many concurrent connections)

2. Single goroutine (sequential):
   for _, url := range urls {
       counts, _ := fetchAndCount(ctx, url)
       merge(counts)
   }
   Pros: Simplest code; no race conditions
   Cons: Slow (no parallelism); one failure blocks all

3. Fan-out/fan-in without worker pool:
   Similar to our approach but spawn goroutine per URL with a semaphore:
   sem := make(chan struct{}, workers)
   Pros: Conceptually simpler (no job queue)
   Cons: All goroutines created upfront (higher memory if many URLs)

4. Use sync.Map instead of channel aggregation:
   var mu sync.Mutex
   globalCounts := make(map[string]int)
   // Each worker locks and merges directly
   Pros: No aggregator goroutine needed
   Cons: Lock contention; harder to reason about

5. Use errgroup from golang.org/x/sync:
   g, ctx := errgroup.WithContext(ctx)
   for _, url := range urls {
       g.Go(func() error { return fetchAndCount(ctx, url) })
   }
   if err := g.Wait(); err != nil { return nil, err }
   Pros: Cleaner error handling
   Cons: External dependency; still need aggregation logic

Go vs X:

Go vs Python (asyncio):
  async def word_count(urls, workers):
      sem = asyncio.Semaphore(workers)
      async def fetch(url):
          async with sem:
              async with aiohttp.ClientSession() as session:
                  async with session.get(url) as resp:
                      text = await resp.text()
                      return tokenize(text)
      tasks = [fetch(url) for url in urls]
      results = await asyncio.gather(*tasks)
      return merge(results)
  Pros: Similar structure
  Cons: GIL prevents true parallelism (CPU-bound tokenization is slow)
        async/await is viral (all callers must be async)
  Go: True parallelism; simpler code (no async coloring)

Go vs JavaScript (Node.js):
  async function wordCount(urls, workers) {
      const limit = pLimit(workers);
      const promises = urls.map(url => limit(() => fetchAndCount(url)));
      const results = await Promise.all(promises);
      return merge(results);
  }
  Pros: Similar brevity with promise libraries
  Cons: Single-threaded (no true parallelism)
        Requires external library (p-limit)
  Go: Built-in worker pools; true parallelism

Go vs Rust (tokio):
  use tokio::task::JoinSet;
  async fn word_count(urls: Vec<String>, workers: usize) -> Result<HashMap<String, usize>> {
      let sem = Arc::new(Semaphore::new(workers));
      let mut set = JoinSet::new();
      for url in urls {
          let permit = sem.clone().acquire_owned().await?;
          set.spawn(async move {
              let _permit = permit;
              fetch_and_count(&url).await
          });
      }
      let mut counts = HashMap::new();
      while let Some(res) = set.join_next().await {
          merge(&mut counts, res??);
      }
      Ok(counts)
  }
  Pros: Zero-cost abstractions; compile-time safety
  Cons: Much more complex (Arc, Semaphore, async traits, lifetime management)
        Steeper learning curve
  Go: Simpler code; faster development

Go vs Java (ExecutorService):
  ExecutorService executor = Executors.newFixedThreadPool(workers);
  List<Future<Map<String, Integer>>> futures = urls.stream()
      .map(url -> executor.submit(() -> fetchAndCount(url)))
      .collect(Collectors.toList());
  Map<String, Integer> counts = new HashMap<>();
  for (Future<Map<String, Integer>> future : futures) {
      merge(counts, future.get());
  }
  executor.shutdown();
  Pros: Similar worker pool concept
  Cons: Much more verbose (generics, Future.get() blocking)
        Heavyweight threads (not goroutines)
  Go: Lighter weight; cleaner syntax
*/
