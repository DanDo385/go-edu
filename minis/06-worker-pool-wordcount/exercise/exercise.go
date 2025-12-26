//go:build !solution
// +build !solution

package exercise

// TODO: Import required packages
// You'll need:
// - "context" for cancellation and timeouts
// - "fmt" for error formatting
// - "io" for reading response bodies
// - "net/http" for making HTTP requests
// - "strings" for text processing
// - "sync" for WaitGroup coordination
// - "unicode" for character classification
//
import (
    "context"
    "fmt"
    "io"
    "net/http"
    "strings"
    "sync"
    "unicode"
)

// ============================================================================
// WORKER POOL PATTERN: Bounded Concurrency for URL Processing
// ============================================================================
//
// A worker pool is a concurrency pattern that:
// 1. Limits the number of concurrent goroutines (workers)
// 2. Distributes work across a fixed number of workers
// 3. Prevents resource exhaustion with thousands of concurrent operations
//
// Why use a worker pool:
// - Controlled resource usage (limit concurrent HTTP connections)
// - Better performance than unbounded goroutines (less context switching)
// - Backpressure handling (slow down producers if workers are busy)
// - Clean error propagation and cancellation
//
// Architecture (Fan-Out/Fan-In Pattern):
// 1. Main goroutine: Sends work items (URLs) to jobs channel
// 2. Worker goroutines: Fetch URLs from jobs channel, process, send results
// 3. Aggregator: Collects results from workers and merges them
// 4. Error handling: First error cancels context, stops all workers
//
// Memory considerations:
// - Channels are bounded to prevent unbounded memory growth
// - Each worker has its own goroutine stack (2-8KB, grows as needed)
// - Maps are allocated per-result, then merged (O(vocabulary) space)
//
// ============================================================================

// ============================================================================
// Exercise 1: Implement WordCount with Worker Pool
// ============================================================================

// WordCount fetches URLs concurrently using a worker pool, tokenizes response bodies,
// and returns overall word frequencies.
//
// TODO: Implement a worker pool that:
// 1. Creates exactly 'workers' goroutines (bounded concurrency)
// 2. Distributes URL fetching across workers
// 3. Aggregates word counts from all workers
// 4. Cancels all work on first error
// 5. Respects context cancellation
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
// Memory allocation:
// - Channels: jobs and results channels (buffered to worker count)
// - Each worker allocates a map for its URL's word counts
// - Final aggregation merges all maps into one
//
// Worker Pool Pattern:

func WordCount(ctx context.Context, urls []string, workers int) (map[string]int, error) {
// 1. Create buffered channels:   
   jobs := make(chan string, workers)       		// URLs to fetch
   results := make(chan map[string]int, workers) 	// Word counts per URL
   errCh := make(chan error, 1)              		// First error only

// 2. Create cancellable context:
   ctx, cancel := context.WithCancel(ctx)
   defer cancel()  // Cleanup on return
// 3. Start worker goroutines:
   var wg sync.WaitGroup
   for i := 0; i < workers; i++ { // 
       wg.Add(1)
       go func() {
           defer wg.Done()
           // Process jobs until channel closed or context cancelled
           for {
               select {
               case <-ctx.Done():
                   return  // Context cancelled, stop worker
               case url, ok := <-jobs:
                   if !ok {
                       return  // Jobs channel closed, no more work
                   }
                   // Fetch and count words
                   counts, err := fetchAndCount(ctx, url)
                   if err != nil {
                       // Send error (non-blocking) and cancel context
                       select {
                       case errCh <- fmt.Errorf("fetching %s: %w", url, err):
                           cancel()  // Stop all other workers
                       default:
                           // Error channel already has error, ignore
                       }
                       return
                   }
                   // Send results (with context check)
                   select {
                   case <-ctx.Done():
                       return
                   case results <- counts:
                   }
               }
           }
       }()
   }
}
//
// 4. Send jobs in separate goroutine:
//    go func() {
//        for _, url := range urls {
//            select {
//            case <-ctx.Done():
//                return
//            case jobs <- url:
//            }
//        }
//        close(jobs)  // Signal no more work
//    }()
//
// 5. Close results channel when all workers finish:
//    go func() {
//        wg.Wait()
//        close(results)
//    }()
//
// 6. Aggregate results in main goroutine:
//    finalCounts := make(map[string]int)
//    for counts := range results {
//        for word, count := range counts {
//            finalCounts[word] += count
//        }
//    }
//
// 7. Check for errors:
//    select {
//    case err := <-errCh:
//        return nil, err
//    default:
//        return finalCounts, nil
//    }
//
// func WordCount(ctx context.Context, urls []string, workers int) (map[string]int, error) {
//     // Implementation here
// }

// ============================================================================
// Exercise 2: Implement fetchAndCount
// ============================================================================

// fetchAndCount fetches a URL and returns word frequencies.
//
// TODO: Implement URL fetching with context awareness:
// 1. Create HTTP request with context (allows cancellation)
// 2. Execute request using http.DefaultClient
// 3. Check status code (200 OK)
// 4. Read response body
// 5. Tokenize and count words
//
// Context-aware HTTP requests:
// - Use http.NewRequestWithContext to create request
// - Context cancellation propagates to HTTP client
// - Ongoing connections are closed when context is cancelled
//
// Pattern:
//   req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
//   if err != nil {
//       return nil, err
//   }
//
//   resp, err := http.DefaultClient.Do(req)
//   if err != nil {
//       return nil, err
//   }
//   defer resp.Body.Close()
//
//   if resp.StatusCode != http.StatusOK {
//       return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
//   }
//
//   body, err := io.ReadAll(resp.Body)
//   if err != nil {
//       return nil, err
//   }
//
//   return tokenizeAndCount(string(body)), nil
//
// Memory considerations:
// - io.ReadAll allocates buffer for entire response body
// - For large responses, consider streaming with bufio.Scanner
// - Response body must be closed (defer resp.Body.Close())
//
// func fetchAndCount(ctx context.Context, url string) (map[string]int, error) {
//     // Implementation here
// }

// ============================================================================
// Exercise 3: Implement tokenizeAndCount
// ============================================================================

// tokenizeAndCount extracts words and returns frequencies.
//
// TODO: Implement word tokenization and counting:
// 1. Split text into words (strings.Fields splits on whitespace)
// 2. Normalize each word to lowercase
// 3. Remove non-letter characters (punctuation)
// 4. Skip empty words
// 5. Count word frequencies in a map
//
// Text processing in Go:
// - strings.Fields: Split on any whitespace (space, tab, newline)
// - strings.ToLower: Convert to lowercase
// - strings.Map: Apply function to each rune (character)
// - unicode.IsLetter: Check if rune is a letter
//
// Pattern:
//   counts := make(map[string]int)
//
//   for _, word := range strings.Fields(text) {
//       // Normalize to lowercase
//       word = strings.ToLower(word)
//
//       // Remove non-letters
//       word = strings.Map(func(r rune) rune {
//           if unicode.IsLetter(r) {
//               return r
//           }
//           return -1  // Drop this character
//       }, word)
//
//       // Skip empty words
//       if word == "" {
//           continue
//       }
//
//       counts[word]++
//   }
//
//   return counts
//
// Map accumulation:
// - Maps in Go are hash tables (O(1) average insert/lookup)
// - Zero value for int is 0, so counts[word]++ works for new keys
// - Maps grow dynamically (rehashing when load factor is high)
//
// func tokenizeAndCount(text string) map[string]int {
//     // Implementation here
// }

// ============================================================================
// Comparison with Other Approaches
// ============================================================================
//
// Alternative 1: Unbounded goroutines (one per URL)
//   for _, url := range urls {
//       go func(u string) {
//           counts, _ := fetchAndCount(ctx, u)
//           // Send to results channel
//       }(url)
//   }
//   Pros: Simpler code (no worker pool logic)
//   Cons: Resource exhaustion with 10,000 URLs (too many concurrent connections)
//         OS limits on file descriptors (each connection is a file descriptor)
//         Thundering herd problem (all requests start simultaneously)
//
// Alternative 2: Sequential processing (no concurrency)
//   for _, url := range urls {
//       counts, err := fetchAndCount(ctx, url)
//       if err != nil {
//           return nil, err
//       }
//       merge(counts)
//   }
//   Pros: Simplest code; no race conditions; predictable resource usage
//   Cons: Very slow (no parallelism)
//         One slow URL blocks all others
//         Wastes CPU while waiting for network I/O
//
// Alternative 3: Semaphore-based limiting
//   sem := make(chan struct{}, workers)
//   for _, url := range urls {
//       sem <- struct{}{}  // Acquire semaphore
//       go func(u string) {
//           defer func() { <-sem }()  // Release semaphore
//           counts, _ := fetchAndCount(ctx, u)
//       }(url)
//   }
//   Pros: Similar limiting without job queue
//   Cons: All goroutines created upfront (higher memory if many URLs)
//         More complex cleanup (need to drain semaphore)
//
// Alternative 4: sync.Map for aggregation
//   var m sync.Map
//   // Each worker updates m directly with locks
//   Pros: No aggregator goroutine
//   Cons: Lock contention (all workers compete for map lock)
//         Harder to reason about (shared mutable state)
//
// Alternative 5: errgroup package (golang.org/x/sync)
//   g, ctx := errgroup.WithContext(ctx)
//   g.SetLimit(workers)  // Go 1.20+
//   for _, url := range urls {
//       url := url
//       g.Go(func() error {
//           counts, err := fetchAndCount(ctx, url)
//           // Still need aggregation logic
//           return err
//       })
//   }
//   if err := g.Wait(); err != nil {
//       return nil, err
//   }
//   Pros: Cleaner error handling and concurrency limiting
//   Cons: External dependency (not in stdlib)
//         Still need separate aggregation logic
//
// ============================================================================
// Go vs Other Languages
// ============================================================================
//
// Go vs Python (asyncio):
//   async def word_count(urls, workers):
//       sem = asyncio.Semaphore(workers)
//       async def fetch(url):
//           async with sem:
//               async with aiohttp.ClientSession() as session:
//                   async with session.get(url) as resp:
//                       text = await resp.text()
//                       return tokenize(text)
//       tasks = [fetch(url) for url in urls]
//       results = await asyncio.gather(*tasks, return_exceptions=True)
//       return merge(results)
//
//   Pros: Similar structure with semaphore limiting
//   Cons: GIL prevents true parallelism (CPU-bound tokenization is slow)
//         async/await is viral (all callers must be async)
//         More verbose setup (aiohttp session management)
//   Go: True parallelism with goroutines
//       Simpler code (no async coloring)
//       Worker pool pattern is more explicit
//
// Go vs JavaScript (Node.js):
//   async function wordCount(urls, workers) {
//       const limit = pLimit(workers);
//       const promises = urls.map(url =>
//           limit(() => fetch(url).then(r => r.text()).then(tokenize))
//       );
//       const results = await Promise.all(promises);
//       return merge(results);
//   }
//
//   Pros: Similar brevity with promise libraries
//   Cons: Single-threaded (no true parallelism for tokenization)
//         Requires external library (p-limit)
//         Less explicit about worker lifecycle
//   Go: Built-in worker pools with channels
//       True parallelism across CPU cores
//       Explicit control over goroutine lifecycle
//
// Go vs Rust (tokio):
//   use tokio::task::JoinSet;
//   async fn word_count(urls: Vec<String>, workers: usize) -> Result<HashMap<String, usize>> {
//       let sem = Arc::new(Semaphore::new(workers));
//       let mut set = JoinSet::new();
//       for url in urls {
//           let permit = sem.clone().acquire_owned().await?;
//           set.spawn(async move {
//               let _permit = permit;
//               fetch_and_count(&url).await
//           });
//       }
//       let mut counts = HashMap::new();
//       while let Some(res) = set.join_next().await {
//           merge(&mut counts, res??);
//       }
//       Ok(counts)
//   }
//
//   Pros: Zero-cost abstractions; compile-time safety
//         Excellent async runtime (tokio)
//   Cons: Much more complex (Arc, Semaphore, async traits, lifetimes)
//         Steeper learning curve
//         More verbose error handling
//   Go: Simpler code; faster development
//       Less boilerplate for common patterns
//       Easier to reason about concurrency
//
// Go vs Java (ExecutorService):
//   ExecutorService executor = Executors.newFixedThreadPool(workers);
//   List<Future<Map<String, Integer>>> futures = urls.stream()
//       .map(url -> executor.submit(() -> fetchAndCount(url)))
//       .collect(Collectors.toList());
//   Map<String, Integer> counts = new HashMap<>();
//   for (Future<Map<String, Integer>> future : futures) {
//       counts.putAll(future.get());  // Blocking wait
//   }
//   executor.shutdown();
//   executor.awaitTermination(1, TimeUnit.MINUTES);
//
//   Pros: Similar worker pool concept
//   Cons: Much more verbose (generics, Future.get() blocking)
//         Heavyweight threads (not goroutines)
//         Manual executor lifecycle management
//   Go: Lighter weight goroutines
//       Cleaner syntax with channels
//       Automatic cleanup with defer
//
// ============================================================================
// After Implementation
// ============================================================================
//
// Testing your implementation:
// - Run: go test -v ./...
// - Run with race detector: go test -race ./...
// - Test with different worker counts (1, 10, 100)
// - Test error cases (invalid URLs, timeouts)
// - Compare performance: 1 worker vs N workers
//
// Performance tips:
// - Optimal worker count ≈ number of CPU cores for CPU-bound work
// - Optimal worker count ≈ 10-100 for I/O-bound work (HTTP fetching)
// - Too many workers: Context switching overhead, memory usage
// - Too few workers: Underutilized CPU/network
//
// Common mistakes:
// - Forgetting to close channels (causes deadlock)
// - Not using context for HTTP requests (can't cancel in-flight requests)
// - Unbuffered channels in hot path (causes unnecessary blocking)
// - Not using defer for cleanup (resource leaks)
// - Capturing loop variables incorrectly (Go 1.21 and earlier)
