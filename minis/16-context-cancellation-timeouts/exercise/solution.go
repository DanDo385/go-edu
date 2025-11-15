//go:build solution
// +build solution

/*
Project 16: Context Cancellation and Timeouts - Solutions

This file contains complete solutions to all exercises with detailed explanations.

Key Go Concepts Demonstrated:
1. Context cancellation and propagation
2. Timeout and deadline handling
3. Goroutine lifecycle management
4. Channel coordination with context
5. Preventing resource leaks

Why Go is well-suited for this:
- context.Context is a first-class citizen in the standard library
- All standard library functions that perform I/O accept context
- Goroutines make concurrent timeout handling trivial
- Channels + context provide clean cancellation patterns
- No callback hell, no promise chains

Compared to other languages:
- Python: asyncio.timeout() added in 3.11, but less integrated
- JavaScript: AbortController is similar but less standardized
- Rust: tokio has timeout, but more complex type system
- Java: CompletableFuture has timeout, but more verbose
*/

package exercise

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// Type definitions (shared between exercise and solution)

type Job struct {
	ID int
}

type Result struct {
	JobID  int
	Output string
	Error  error
}

type cacheEntry struct {
	value      interface{}
	expiration time.Time
}

// ============================================================================
// Exercise 1: RetryWithTimeout
// ============================================================================

/*
Problem: Retry operations with timeout and exponential backoff

We need to:
1. Retry a function up to maxRetries times
2. Each attempt has its own timeout
3. Use exponential backoff between retries
4. Respect parent context cancellation

Architecture:
- Loop from 0 to maxRetries
- For each attempt:
  - Create child context with timeout
  - Call function with child context
  - If success, return
  - If failure and not last attempt, backoff with context check
  - If parent context cancelled, return immediately

Complexity:
- Time: O(maxRetries * timeout + sum of backoffs)
- Space: O(1) - only stores error

Three-Input Iteration Table:

Input 1: Success on 3rd attempt
  Attempt 0: fn() fails, wait 100ms
  Attempt 1: fn() fails, wait 200ms
  Attempt 2: fn() succeeds → return nil

Input 2: All attempts fail
  Attempt 0: fn() fails, wait 100ms
  Attempt 1: fn() fails, wait 200ms
  Attempt 2: fn() fails → return error

Input 3: Parent context cancelled during backoff
  Attempt 0: fn() fails, wait 100ms
  Parent cancelled → return context.Canceled
*/

func RetryWithTimeout(
	ctx context.Context,
	fn func(context.Context) error,
	maxRetries int,
	timeout time.Duration,
) error {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		// Create child context with timeout for this attempt
		// This ensures each attempt has its own timeout
		attemptCtx, cancel := context.WithTimeout(ctx, timeout)

		// Try the operation
		err := fn(attemptCtx)
		cancel() // Always cancel to release resources

		if err == nil {
			// Success!
			return nil
		}

		// Store error for potential return
		lastErr = err

		// If this was the last attempt, don't backoff
		if attempt == maxRetries-1 {
			break
		}

		// Calculate exponential backoff: 100ms * 2^attempt
		// Attempt 0: 100ms, Attempt 1: 200ms, Attempt 2: 400ms, etc.
		backoff := time.Duration(100) * time.Millisecond * (1 << uint(attempt))

		// Wait with context awareness
		// If parent context is cancelled during backoff, stop immediately
		select {
		case <-time.After(backoff):
			// Backoff complete, continue to next attempt
		case <-ctx.Done():
			// Parent context cancelled, stop retrying
			return ctx.Err()
		}
	}

	// All attempts failed
	return lastErr
}

// ============================================================================
// Exercise 2: FetchAll
// ============================================================================

/*
Problem: Fetch multiple URLs concurrently with timeout

We need to:
1. Fetch all URLs concurrently
2. Enforce total timeout for all fetches
3. Cancel all on first error
4. Return results in same order as input URLs

Architecture:
- Create context with timeout
- Create channels for results and errors
- Start goroutine for each URL
- Use index to preserve order
- Collect results into slice (maintaining order)
- Cancel all on first error

Complexity:
- Time: O(slowest fetch) due to concurrency
- Space: O(n) for results and goroutines

Three-Input Iteration Table:

Input 1: All succeed
  Goroutine 0: Fetch URL0 → send result[0]
  Goroutine 1: Fetch URL1 → send result[1]
  Goroutine 2: Fetch URL2 → send result[2]
  Collect all → return results in order

Input 2: URL1 fails
  Goroutine 0: Fetch URL0 → send result[0]
  Goroutine 1: Fetch URL1 → error → cancel context
  Goroutine 2: Context cancelled → stop
  Return error

Input 3: Timeout
  All goroutines: Fetching...
  Timeout expires → context cancelled
  All goroutines: Stop
  Return context.DeadlineExceeded
*/

type fetchResult struct {
	index int
	body  string
	err   error
}

func FetchAll(ctx context.Context, urls []string, timeout time.Duration) ([]string, error) {
	// Create context with timeout for entire operation
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel() // Ensure cleanup

	// Channel for results (buffered to prevent goroutines blocking)
	resultCh := make(chan fetchResult, len(urls))

	// Start goroutine for each URL
	for i, url := range urls {
		go func(index int, u string) {
			// Fetch URL
			body, err := fetchURL(ctx, u)

			// Send result with index to preserve order
			resultCh <- fetchResult{
				index: index,
				body:  body,
				err:   err,
			}

			// If error occurred, cancel context to stop other fetches
			if err != nil {
				cancel()
			}
		}(i, url)
	}

	// Collect results
	results := make([]string, len(urls))
	for i := 0; i < len(urls); i++ {
		result := <-resultCh

		if result.err != nil {
			// Cancel all other fetches
			cancel()
			return nil, result.err
		}

		// Store result at correct index
		results[result.index] = result.body
	}

	return results, nil
}

// fetchURL fetches a URL with context support
func fetchURL(ctx context.Context, url string) (string, error) {
	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	// Execute request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// ============================================================================
// Exercise 3: WorkerPool
// ============================================================================

/*
Problem: Worker pool with graceful shutdown

We need to:
1. Start numWorkers goroutines
2. Each worker processes jobs from channel
3. Stop when jobs channel is closed (graceful)
4. Stop when context is cancelled (immediate)
5. Close results channel when all workers exit

Architecture:
- Create results channel
- Start numWorkers goroutines
- Each worker:
  - Reads from jobs channel
  - Checks context cancellation
  - Processes job
  - Sends result
- Use WaitGroup to track workers
- Close results when all workers done

Complexity:
- Time: O(numJobs / numWorkers * processingTime)
- Space: O(numWorkers) for goroutines

Three-Input Iteration Table:

Input 1: All jobs processed
  Jobs: [J0, J1, J2, J3, J4]
  Worker0: Process J0 → result
  Worker1: Process J1 → result
  Worker2: Process J2 → result
  Worker0: Process J3 → result
  Worker1: Process J4 → result
  Jobs closed → workers exit → results closed

Input 2: Context cancelled mid-processing
  Jobs: [J0, J1, J2, J3, J4]
  Worker0: Processing J0
  Context cancelled
  Worker0: Check ctx.Done() → exit
  Worker1: Check ctx.Done() → exit
  Worker2: Check ctx.Done() → exit
  Results closed

Input 3: No jobs
  Jobs: closed immediately
  All workers: Read from closed channel → exit
  Results closed
*/

func WorkerPool(ctx context.Context, numWorkers int, jobs <-chan Job) <-chan Result {
	// Create results channel (buffered to prevent workers blocking)
	results := make(chan Result, numWorkers)

	// WaitGroup to track worker completion
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// Process jobs until channel is closed or context is cancelled
			for {
				select {
				case <-ctx.Done():
					// Context cancelled, stop immediately
					return

				case job, ok := <-jobs:
					if !ok {
						// Jobs channel closed, no more work
						return
					}

					// Process job
					result := processJob(ctx, job)

					// Send result (check context in case it was cancelled)
					select {
					case results <- result:
						// Sent successfully
					case <-ctx.Done():
						// Context cancelled while sending
						return
					}
				}
			}
		}(i)
	}

	// Close results channel when all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

// processJob simulates job processing
func processJob(ctx context.Context, job Job) Result {
	// Check if context is cancelled before processing
	select {
	case <-ctx.Done():
		return Result{
			JobID: job.ID,
			Error: ctx.Err(),
		}
	default:
	}

	// Simulate work (with context check)
	select {
	case <-time.After(10 * time.Millisecond):
		// Work completed
		return Result{
			JobID:  job.ID,
			Output: fmt.Sprintf("Processed job %d", job.ID),
		}
	case <-ctx.Done():
		// Cancelled during processing
		return Result{
			JobID: job.ID,
			Error: ctx.Err(),
		}
	}
}

// ============================================================================
// Exercise 4: CacheWithTTL
// ============================================================================

/*
Problem: Cache with automatic expiration

We need to:
1. Store key-value pairs with TTL
2. Automatically expire entries
3. Clean up expired entries periodically
4. Be thread-safe

Architecture:
- Map of key → cacheEntry (value + expiration time)
- Mutex for concurrent access
- Cleanup goroutine that periodically removes expired entries
- Context to stop cleanup goroutine

Complexity:
- Set: O(1)
- Get: O(1)
- Cleanup: O(n) where n = number of entries

Three-Input Iteration Table:

Input 1: Normal usage
  Set("key1", "val1", 1s)
  Get("key1") → "val1", true
  Sleep 1.5s
  Get("key1") → "", false (expired)

Input 2: Cleanup removes expired
  Set("key1", "val1", 100ms)
  Set("key2", "val2", 100ms)
  Cleanup runs every 50ms
  After 150ms: both entries removed

Input 3: Context cancels cleanup
  Cleanup running
  Context cancelled
  Cleanup exits
*/

type Cache struct {
	mu      sync.RWMutex
	entries map[string]cacheEntry
}

func NewCache() *Cache {
	return &Cache{
		entries: make(map[string]cacheEntry),
	}
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = cacheEntry{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(entry.expiration) {
		return nil, false
	}

	return entry.value, true
}

func (c *Cache) Cleanup(ctx context.Context) {
	// Cleanup every 100ms
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Remove expired entries
			c.removeExpired()

		case <-ctx.Done():
			// Context cancelled, stop cleanup
			return
		}
	}
}

func (c *Cache) removeExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.entries {
		if now.After(entry.expiration) {
			delete(c.entries, key)
		}
	}
}

// ============================================================================
// Exercise 5: RateLimiter
// ============================================================================

/*
Problem: Context-aware rate limiter

We need to:
1. Allow 'rate' operations per second
2. Block when rate limit exceeded
3. Respect context cancellation

Architecture:
- Buffered channel as token bucket
- Fill tokens at rate/second
- Wait() consumes a token (blocks if none available)
- Context allows cancellation while waiting

Complexity:
- Wait: O(1) - just channel receive
- Space: O(rate) - token buffer

Three-Input Iteration Table:

Input 1: Normal rate limiting (10 ops/sec)
  Op 0: Token available → proceed
  Op 1: Token available → proceed
  ...
  Op 10: Token available → proceed
  Op 11: Wait for token (blocked)
  After 100ms: Token refilled → proceed

Input 2: Context cancelled while waiting
  Op 0: Token available → proceed
  Op 1: Wait for token (blocked)
  Context cancelled
  Op 1: Return context.Canceled

Input 3: High concurrency
  100 goroutines call Wait()
  Rate = 10/sec
  First 10: Proceed immediately
  Rest 90: Wait for tokens
  Tokens refilled over ~9 seconds
*/

type RateLimiter struct {
	tokens chan struct{}
}

func NewRateLimiter(rate int) *RateLimiter {
	rl := &RateLimiter{
		tokens: make(chan struct{}, rate),
	}

	// Fill bucket initially
	for i := 0; i < rate; i++ {
		rl.tokens <- struct{}{}
	}

	// Refill tokens at specified rate
	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(rate))
		defer ticker.Stop()

		for range ticker.C {
			// Try to add token (non-blocking)
			select {
			case rl.tokens <- struct{}{}:
				// Token added
			default:
				// Bucket full, skip
			}
		}
	}()

	return rl
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
	// Wait for token or context cancellation
	select {
	case <-rl.tokens:
		// Got token, can proceed
		return nil
	case <-ctx.Done():
		// Context cancelled
		return ctx.Err()
	}
}

/*
Alternatives & Trade-offs:

1. time.Ticker-based limiter:
   ticker := time.NewTicker(time.Second / rate)
   for range ticker.C { ... }
   Pros: Simpler implementation
   Cons: Doesn't allow bursts; strict 1 op every interval

2. Token bucket with time.Sleep:
   if lastOp + interval > now { time.Sleep(interval) }
   Pros: No goroutine for refilling
   Cons: Can't use context to interrupt sleep

3. Semaphore-based (golang.org/x/sync/semaphore):
   sem := semaphore.NewWeighted(rate)
   sem.Acquire(ctx, 1)
   Pros: Battle-tested library
   Cons: External dependency

4. Sliding window counter:
   Track ops in last second, reject if exceeds rate
   Pros: More accurate over time
   Cons: Higher memory usage

Go vs X:

Go vs Python (asyncio.Semaphore):
  async with semaphore:
      await do_work()
  Pros: Similar API
  Cons: No built-in rate limiting (need external lib)
  Go: Channels make rate limiting natural

Go vs JavaScript (bottleneck package):
  const limiter = new Bottleneck({ maxConcurrent: 10, minTime: 100 });
  await limiter.schedule(() => doWork());
  Pros: Good rate limiting library
  Cons: Requires external package
  Go: Easy to implement from scratch with channels

Go vs Rust (governor crate):
  let limiter = RateLimiter::direct(Quota::per_second(nonzero!(10u32)));
  limiter.until_ready().await;
  Pros: Type-safe, efficient
  Cons: More complex type system
  Go: Simpler implementation

Common Mistakes to Avoid:

1. Not buffering token channel:
   tokens := make(chan struct{}) // WRONG: unbuffered
   Causes: Refill goroutine blocks if no one is waiting

2. Not handling context cancellation:
   <-rl.tokens // WRONG: can't be interrupted
   Causes: Hangs if context is cancelled

3. Creating goroutine per Wait():
   go func() { <-rl.tokens }() // WRONG: goroutine leak
   Causes: Many goroutines if high concurrency

4. Not deferring cancel():
   ctx, cancel := context.WithTimeout(...)
   // Forgot defer cancel()
   Causes: Timer goroutine leak

5. Storing context in struct:
   type Limiter struct { ctx context.Context } // WRONG
   Causes: All users share same context
*/
