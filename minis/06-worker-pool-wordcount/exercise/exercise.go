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
	"sync"
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
	// TODO: Step 1 - Create buffered channels
	//   - jobs: chan string, buffer size = workers (URLs to fetch)
	//   - results: chan map[string]int, buffer size = workers (word counts per URL)
	//   - errCh: chan error, buffer size = 1 (first error only)
	//
	//   Example:
	//     jobs := make(chan string, workers)
	//     results := make(chan map[string]int, workers)
	//     errCh := make(chan error, 1)

	// TODO: Step 2 - Create cancellable context
	//   - Use context.WithCancel(ctx) to create child context
	//   - Use defer cancel() to ensure cleanup
	//
	//   Example:
	//     ctx, cancel := context.WithCancel(ctx)
	//     defer cancel()

	// TODO: Step 3 - Create WaitGroup
	//   - var wg sync.WaitGroup
	//   - This will track worker goroutine completion
	var wg sync.WaitGroup

	/*
		WORKER SPAWN LOOP - UNDERSTANDING ITERATION
		--------------------------------------------

		This loop runs SEQUENTIALLY in the main goroutine.
		Each iteration:
		  1. Increments WaitGroup counter: wg.Add(1)
		  2. Launches a goroutine (runs concurrently)
		  3. Increments loop counter: i++

		HOW THE COUNTER `i` WORKS:
		--------------------------
		- `i` starts at 0, increments each iteration: 0, 1, 2, ..., workers-1
		- Each iteration creates ONE goroutine
		- After `workers` iterations, we have `workers` goroutines running

		CRITICAL: Why we pass `i` as parameter `(i)`
		---------------------------------------------
		This avoids a closure variable capture bug!

		❌ WRONG (without parameter):
		for i := 0; i < workers; i++ {
		    go func() {
		        fmt.Println(i)  // BUG: All goroutines print the SAME value!
		    }()
		}
		Problem: All goroutines capture the SAME variable `i` by reference.
		        By the time goroutines run, loop has finished, i = workers.
		        All goroutines see i = workers (not 0, 1, 2, ...).

		✅ CORRECT (with parameter):
		for i := 0; i < workers; i++ {
		    go func(workerID int) {
		        fmt.Println(workerID)  // Each goroutine gets its own copy!
		    }(i)  // Pass i as argument
		}
		Solution: Each goroutine receives its own COPY of i's value.
		        Goroutine 0 gets workerID=0, goroutine 1 gets workerID=1, etc.

		VARIABLE MORPHING THROUGH THE LOOP:
		-----------------------------------
		Iteration 0: i=0 → spawn goroutine with workerID=0
		Iteration 1: i=1 → spawn goroutine with workerID=1
		Iteration 2: i=2 → spawn goroutine with workerID=2
		...
		Iteration N-1: i=workers-1 → spawn goroutine with workerID=workers-1

		Each goroutine has its own independent copy of workerID.
	*/
	// TODO: Step 4 - Start worker goroutines
	//   Loop from i := 0 to workers-1:
	//     1. Call wg.Add(1) BEFORE starting goroutine
	//     2. Launch goroutine: go func(workerID int) { ... }(i)
	//        - CRITICAL: Pass i as parameter to avoid closure bug!
	//        - Use defer wg.Done() inside goroutine
	//     3. Inside goroutine, create infinite loop:
	//        for {
	//            select {
	//            case <-ctx.Done():
	//                return  // Context cancelled
	//            case url, ok := <-jobs:
	//                if !ok {
	//                    return  // Channel closed
	//                }
	//                // Process job: fetchAndCount(ctx, url)
	//                // Handle errors: send to errCh, cancel context
	//                // Send results: results <- counts
	//            }
	//        }
	//
	//   Key points:
	//   - wg.Add(1) must be called BEFORE go statement
	//   - defer wg.Done() ensures counter decremented on exit
	//   - Pass i as parameter: go func(workerID int) { ... }(i)
	//   - Use select to check ctx.Done() and receive from jobs
	//   - On error: non-blocking send to errCh, then cancel()
	//   - On success: send counts to results channel

	// TODO: Implement worker loop here
	// for i := 0; i < workers; i++ {
	//     wg.Add(1)
	//     go func(workerID int) {
	//         defer wg.Done()
	//         // Worker implementation
	//     }(i)
	// }
	// TODO: Step 5 - Send jobs in separate goroutine
	//   Launch goroutine that:
	//     1. Iterates over urls: for _, url := range urls
	//     2. For each URL, use select to check ctx.Done() and send to jobs
	//     3. After all URLs sent, close(jobs) to signal no more work
	//
	//   Pattern:
	//     go func() {
	//         defer close(jobs)  // Or close at end
	//         for _, url := range urls {
	//             select {
	//             case <-ctx.Done():
	//                 return  // Stop if cancelled
	//             case jobs <- url:
	//                 // Sent successfully
	//             }
	//         }
	//         close(jobs)
	//     }()
	//
	//   Why separate goroutine?
	//   - Allows workers to start processing while URLs are still being sent
	//   - Prevents blocking main goroutine

	// TODO: Step 6 - Close results channel when all workers finish
	//   Launch goroutine that:
	//     1. Calls wg.Wait() to wait for all workers
	//     2. After workers finish, close(results)
	//
	//   Pattern:
	//     go func() {
	//         wg.Wait()
	//         close(results)
	//     }()
	//
	//   Why needed?
	//   - Aggregator uses `for counts := range results`
	//   - Range loop exits when channel is closed
	//   - Must wait for all workers before closing

	// TODO: Step 7 - Aggregate results in main goroutine
	//   1. Create finalCounts map: finalCounts := make(map[string]int)
	//   2. Range over results channel: for counts := range results
	//   3. For each counts map, iterate: for word, count := range counts
	//   4. Accumulate: finalCounts[word] += count
	//
	//   Pattern:
	//     finalCounts := make(map[string]int)
	//     for counts := range results {
	//         for word, count := range counts {
	//             finalCounts[word] += count
	//         }
	//     }
	//
	//   Why no locks needed?
	//   - Only main goroutine writes to finalCounts
	//   - Workers send through channel (thread-safe)
	//   - No shared mutable state

	// TODO: Step 8 - Check for errors
	//   Use non-blocking receive from errCh:
	//     select {
	//     case err := <-errCh:
	//         return nil, err
	//     default:
	//         // No error
	//     }
	//
	//   Why non-blocking?
	//   - errCh might be empty (no error occurred)
	//   - Blocking receive would wait forever

	// TODO: Step 9 - Return results
	//   return finalCounts, nil

	// Placeholder to prevent compilation errors
	return nil, fmt.Errorf("not implemented")
}

// ============================================================================
// Exercise 2: Implement fetchAndCount
// ============================================================================

// fetchAndCount fetches a URL and returns word frequencies.
//
// FUNCTION CALL STACK AND VARIABLE FLOW:
// --------------------------------------
// Input: ctx (context.Context), url (string)
//
//	↓
//
// http.NewRequestWithContext(ctx, ...)
//
//	↓
//
// req (*http.Request) - contains URL, method, context
//
//	↓
//
// http.DefaultClient.Do(req)
//
//	↓
//
// resp (*http.Response) - contains status, headers, body
//
//	↓
//
// io.ReadAll(resp.Body)
//
//	↓
//
// body ([]byte) - entire response body in memory
//
//	↓
//
// string(body) - convert bytes to string
//
//	↓
//
// tokenizeAndCount(string)
//
//	↓
//
// counts (map[string]int) - word frequencies
//
//	↓
//
// Return: counts, error
//
// MEMORY ALLOCATIONS:
// ------------------
// 1. req: ~200 bytes (HTTP request struct)
// 2. resp: ~500 bytes (HTTP response struct)
// 3. body: O(response_size) - entire body read into memory
// 4. string(body): O(response_size) - copy of body as string
// 5. counts map: O(vocabulary_size) - hash table for word counts
//
// Total per URL: ~700 bytes + 2 * response_size + vocabulary_size
//
// CONTEXT PROPAGATION:
// -------------------
// Passing ctx into NewRequestWithContext ensures that:
//   - If cancel() is called, the HTTP request is aborted immediately
//   - The underlying TCP connection is closed
//   - No time wasted waiting for slow/failed requests
//   - Resources are freed promptly
//
// Key idea:
// ---------
// Context cancellation propagates through the HTTP client stack.
// When ctx.Done() is closed, the HTTP transport layer detects it and
// aborts the connection, even if the request is in-flight.
//
// TODO: Implement URL fetching with context awareness:
// 1. Create HTTP request with context (allows cancellation)
// 2. Execute request using http.DefaultClient
// 3. Check status code (200 OK)
// 4. Read response body
// 5. Tokenize and count words
//
// Pattern:
//
//	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
//	if err != nil {
//	    return nil, err
//	}
//
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil {
//	    return nil, err
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//	    return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
//	}
//
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//	    return nil, err
//	}
//
//	return tokenizeAndCount(string(body)), nil
func fetchAndCount(ctx context.Context, url string) (map[string]int, error) {
	// TODO: Implement fetchAndCount
	//   1. Create HTTP request: http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	//   2. Execute request: http.DefaultClient.Do(req)
	//   3. Check error and status code (must be 200 OK)
	//   4. Read response body: io.ReadAll(resp.Body)
	//   5. Convert to string and call tokenizeAndCount
	//   6. Return counts and error
	//
	//   Don't forget:
	//   - defer resp.Body.Close() to free resources
	//   - Check resp.StatusCode == http.StatusOK
	//   - Handle all errors appropriately

	return nil, fmt.Errorf("not implemented")
}

// ============================================================================
// Exercise 3: Implement tokenizeAndCount
// ============================================================================

// tokenizeAndCount extracts words and returns frequencies.
//
// TOKENIZATION - WHAT IT MEANS:
// -----------------------------
// Tokenization is the process of breaking text into individual words (tokens).
// Example: "Hello, world!" → ["hello", "world"]
//
// This function:
//  1. Splits text into words (by whitespace)
//  2. Normalizes words (lowercase, remove punctuation)
//  3. Counts frequency of each word
//
// VARIABLE TRANSFORMATION THROUGH ITERATION:
// ------------------------------------------
// Input: text = "Hello, world! Go is great."
//
// Iteration 1:
//
//	word = "Hello," (from strings.Fields)
//	word = "hello," (after ToLower)
//	word = "hello" (after Map removes comma)
//	counts["hello"] = 1
//
// Iteration 2:
//
//	word = "world!"
//	word = "world!"
//	word = "world"
//	counts["world"] = 1
//
// Iteration 3:
//
//	word = "Go"
//	word = "go"
//	word = "go"
//	counts["go"] = 1
//
// ... continues for each word ...
//
// MEMORY CONSIDERATIONS:
// ---------------------
// - strings.Fields: Allocates slice of strings (O(text_length))
// - strings.ToLower: May allocate new string if changes needed
// - strings.Map: Allocates new string for each word (O(word_length))
// - Map: Grows dynamically as words are added (O(vocabulary_size))
//
// Total memory: O(text_length + vocabulary_size)
//
// TODO: Implement word tokenization and counting:
// 1. Split text into words (strings.Fields splits on whitespace)
// 2. Normalize each word to lowercase
// 3. Remove non-letter characters (punctuation)
// 4. Skip empty words
// 5. Count word frequencies in a map
//
// Pattern:
//
//	counts := make(map[string]int)
//
//	for _, word := range strings.Fields(text) {
//	    // Normalize to lowercase
//	    word = strings.ToLower(word)
//
//	    // Remove non-letters
//	    word = strings.Map(func(r rune) rune {
//	        if unicode.IsLetter(r) {
//	            return r
//	        }
//	        return -1  // Drop this character
//	    }, word)
//
//	    // Skip empty words
//	    if word == "" {
//	        continue
//	    }
//
//	    counts[word]++
//	}
//
//	return counts
//
// Map accumulation:
// - Maps in Go are hash tables (O(1) average insert/lookup)
// - Zero value for int is 0, so counts[word]++ works for new keys
// - Maps grow dynamically (rehashing when load factor is high)
func tokenizeAndCount(text string) map[string]int {
	// TODO: Implement tokenizeAndCount
	//   1. Create counts map: counts := make(map[string]int)
	//   2. Split text into words: strings.Fields(text)
	//   3. For each word:
	//      a. Convert to lowercase: strings.ToLower(word)
	//      b. Remove non-letters: strings.Map(func(r rune) rune { ... }, word)
	//         - Use unicode.IsLetter(r) to check if character is a letter
	//         - Return r to keep, -1 to delete
	//      c. Skip empty words: if word == "" { continue }
	//      d. Count: counts[word]++
	//   4. Return counts map
	//
	//   Pattern:
	//     counts := make(map[string]int)
	//     for _, word := range strings.Fields(text) {
	//         word = strings.ToLower(word)
	//         word = strings.Map(func(r rune) rune {
	//             if unicode.IsLetter(r) {
	//                 return r
	//             }
	//             return -1
	//         }, word)
	//         if word == "" {
	//             continue
	//         }
	//         counts[word]++
	//     }
	//     return counts

	return make(map[string]int) // Placeholder
}

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
