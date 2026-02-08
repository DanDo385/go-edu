//go:build reference

package contextcancellationtimeouts

/*
Reference Solution
==================

This file is the canonical reference for this exercise. It keeps failure paths
explicit when an operation can fail, so callers can decide how to handle
errors at API boundaries.

Read this alongside exercise.go and the tests to understand the intended data
flow, ownership boundaries, and invariants that keep behavior deterministic.

Teaching notes:
- Memory/ownership: make copies when returning mutable data that should not
  alias internal state; share references only when aliasing is intentional.
- Invariants: establish assumptions close to construction, and rely on them in
  smaller helper functions to keep logic easy to audit.
- Error surfaces: prefer explicit returns over hidden panics so learners can
  reason about control flow in production-style code.
*/

/*
Project 16: Context Cancellation and Timeouts - Solutions

This file contains complete solutions with extensive debugging support.

Key Go Concepts Demonstrated:
1. Context cancellation and propagation
2. Timeout and deadline handling
3. Goroutine lifecycle management
4. Channel coordination with context
5. Preventing resource leaks

DEBUGGING GUIDE:
- BREAKPOINT comments mark ideal breakpoint locations
- DEBUG comments explain what to observe in debugger
- Use Step Over (F10) to execute line by line
- Use Step Into (F11) to enter function calls
- Watch panel shows variable values in real-time
*/

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
	// BREAKPOINT: Set breakpoint here to start debugging retry logic
	// DEBUG: Watch variables: attempt, lastErr, maxRetries, timeout
	var lastErr error

	// BREAKPOINT: Set breakpoint in loop to observe each retry attempt
	// DEBUG: Watch attempt counter incrementing from 0 to maxRetries-1
	for attempt := 0; attempt < maxRetries; attempt++ {
		// DEBUG: Observe attemptCtx being created with timeout
		// DEBUG: Each attempt gets fresh context with same timeout duration
		attemptCtx, cancel := context.WithTimeout(ctx, timeout)

		// BREAKPOINT: Set breakpoint before fn() call to step into operation
		// DEBUG: Watch attemptCtx.Done() channel and deadline
		err := fn(attemptCtx)

		// DEBUG: Always call cancel() to prevent context leak
		// DEBUG: Without cancel(), timer goroutine stays in memory
		cancel()

		// BREAKPOINT: Set breakpoint here to check if operation succeeded
		// DEBUG: Watch err variable - nil means success
		if err == nil {
			// DEBUG: Early return on success, no more retries needed
			return nil
		}

		// DEBUG: Store error in case all attempts fail
		// DEBUG: lastErr will be final error if loop completes
		lastErr = err

		// DEBUG: Check if this is the last attempt
		// DEBUG: If yes, skip backoff and exit loop
		if attempt == maxRetries-1 {
			break
		}

		// BREAKPOINT: Set breakpoint to observe exponential backoff calculation
		// DEBUG: Watch backoff value: 100ms, 200ms, 400ms, 800ms, etc.
		// DEBUG: Bit shift (1 << attempt) doubles the delay each time
		backoff := time.Duration(100) * time.Millisecond * (1 << uint(attempt))

		// BREAKPOINT: Set breakpoint on select to observe backoff vs cancellation
		// DEBUG: Watch both channels - which one fires first?
		// DEBUG: time.After creates a channel that fires after backoff duration
		select {
		case <-time.After(backoff):
			// DEBUG: Backoff complete, will retry
			// DEBUG: Loop continues to next attempt
		case <-ctx.Done():
			// DEBUG: Parent context cancelled during backoff
			// DEBUG: Return immediately, don't continue retrying
			return ctx.Err()
		}
	}

	// BREAKPOINT: Set breakpoint to see final error when all attempts failed
	// DEBUG: Watch lastErr - this is the error from the last attempt
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

// FetchAll implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func FetchAll(ctx context.Context, urls []string, timeout time.Duration) ([]string, error) {
	// BREAKPOINT: Set breakpoint to observe context creation with timeout
	// DEBUG: Watch timeout duration and resulting deadline
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel() // DEBUG: Ensure cleanup even if function panics

	// BREAKPOINT: Set breakpoint to observe buffered channel creation
	// DEBUG: Watch resultCh capacity - matches len(urls) to prevent blocking
	// DEBUG: Buffered channels allow goroutines to send without waiting
	resultCh := make(chan fetchResult, len(urls))

	// BREAKPOINT: Set breakpoint in loop to observe goroutine launches
	// DEBUG: Watch i and url values being captured for each goroutine
	for i, url := range urls {
		// DEBUG: Each goroutine gets copy of index and url as parameters
		// DEBUG: This prevents race conditions from loop variable capture
		go func(index int, u string) {
			// BREAKPOINT: Set breakpoint inside goroutine to observe fetch
			// DEBUG: Watch ctx.Done() channel - may close during fetch
			body, err := fetchURL(ctx, u)

			// BREAKPOINT: Set breakpoint to observe result being sent
			// DEBUG: Watch fetchResult struct fields being populated
			// DEBUG: index preserves original URL order in results
			resultCh <- fetchResult{
				index: index,
				body:  body,
				err:   err,
			}

			// DEBUG: If error, cancel context to stop other fetches early
			// DEBUG: Fast-fail pattern: one failure cancels all
			if err != nil {
				cancel()
			}
		}(i, url)
	}

	// BREAKPOINT: Set breakpoint before result collection loop
	// DEBUG: Watch results slice being populated
	// DEBUG: Must receive from all goroutines to prevent leaks
	results := make([]string, len(urls))
	for i := 0; i < len(urls); i++ {
		// BREAKPOINT: Set breakpoint on receive to watch results arrive
		// DEBUG: Results may arrive out of order (concurrent execution)
		// DEBUG: Watch result.index to see which URL completed
		result := <-resultCh

		// BREAKPOINT: Set breakpoint to observe error handling
		// DEBUG: First error cancels context and returns immediately
		if result.err != nil {
			cancel()
			return nil, result.err
		}

		// DEBUG: Store result at original index to maintain URL order
		// DEBUG: Watch results[result.index] being assigned
		results[result.index] = result.body
	}

	// BREAKPOINT: Set breakpoint to observe successful completion
	// DEBUG: Watch complete results slice with all fetched bodies
	return results, nil
}

// fetchURL fetches a URL with context support
func fetchURL(ctx context.Context, url string) (string, error) {
	// BREAKPOINT: Set breakpoint to observe HTTP request creation
	// DEBUG: Watch req being created with context attached
	// DEBUG: Context cancellation will abort HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	// BREAKPOINT: Set breakpoint before HTTP request execution
	// DEBUG: Watch ctx.Done() channel - request stops if context cancelled
	// DEBUG: Step into Do() to see request execution
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// BREAKPOINT: Set breakpoint to observe HTTP status code
	// DEBUG: Watch resp.StatusCode - should be 200 for success
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// BREAKPOINT: Set breakpoint before reading response body
	// DEBUG: Watch body bytes being read from network
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
	// BREAKPOINT: Set breakpoint to observe worker pool initialization
	// DEBUG: Watch numWorkers and results channel capacity
	// DEBUG: Buffered channel prevents workers from blocking on send
	results := make(chan Result, numWorkers)

	// DEBUG: WaitGroup tracks how many workers are still running
	// DEBUG: Watch wg.counter to see active worker count
	var wg sync.WaitGroup

	// BREAKPOINT: Set breakpoint in loop to observe worker creation
	// DEBUG: Watch i incrementing as each worker starts
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		// DEBUG: Each worker gets its own goroutine with unique ID
		go func(workerID int) {
			defer wg.Done()

			// BREAKPOINT: Set breakpoint in worker loop to observe job processing
			// DEBUG: Watch workerID to identify which worker is executing
			for {
				// DEBUG: Select monitors both context cancellation and job arrival
				// DEBUG: Whichever channel is ready first will be executed
				select {
				case <-ctx.Done():
					// DEBUG: Context cancelled - immediate shutdown
					// DEBUG: Worker exits without processing remaining jobs
					return

				case job, ok := <-jobs:
					// BREAKPOINT: Set breakpoint to observe job receipt
					// DEBUG: Watch ok flag - false means channel closed
					if !ok {
						// DEBUG: Jobs channel closed - graceful shutdown
						// DEBUG: Worker exits after finishing current job
						return
					}

					// BREAKPOINT: Set breakpoint before job processing
					// DEBUG: Watch job.ID and step into processJob()
					result := processJob(ctx, job)

					// BREAKPOINT: Set breakpoint before sending result
					// DEBUG: Inner select checks context during send operation
					// DEBUG: This prevents blocking if context cancelled
					select {
					case results <- result:
						// DEBUG: Result sent successfully
					case <-ctx.Done():
						// DEBUG: Context cancelled while sending result
						// DEBUG: Exit immediately without sending
						return
					}
				}
			}
		}(i)
	}

	// BREAKPOINT: Set breakpoint in cleanup goroutine
	// DEBUG: This goroutine waits for all workers to finish
	go func() {
		// DEBUG: wg.Wait() blocks until all wg.Done() called
		// DEBUG: Watch wg.counter decreasing to zero
		wg.Wait()
		// DEBUG: Close results channel to signal no more results coming
		close(results)
	}()

	return results
}

// processJob simulates job processing
func processJob(ctx context.Context, job Job) Result {
	// BREAKPOINT: Set breakpoint to check context before work starts
	// DEBUG: Watch ctx.Done() channel - may already be closed
	select {
	case <-ctx.Done():
		// DEBUG: Context cancelled before work started
		// DEBUG: Return error result immediately
		return Result{
			JobID: job.ID,
			Error: ctx.Err(),
		}
	default:
		// DEBUG: Context still active, proceed with work
	}

	// BREAKPOINT: Set breakpoint on select to watch simulated work
	// DEBUG: Select races between work completion and context cancellation
	// DEBUG: Watch which case executes - time.After or ctx.Done()
	select {
	case <-time.After(10 * time.Millisecond):
		// DEBUG: Work completed successfully within timeout
		// DEBUG: Watch result being created with job output
		return Result{
			JobID:  job.ID,
			Output: fmt.Sprintf("Processed job %d", job.ID),
		}
	case <-ctx.Done():
		// DEBUG: Context cancelled during work processing
		// DEBUG: Return error result, abandon work
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

// NewCache implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func NewCache() *Cache {
	return &Cache{
		entries: make(map[string]cacheEntry),
	}
}

// Set implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	// BREAKPOINT: Set breakpoint to observe write lock acquisition
	// DEBUG: Lock() blocks if another goroutine holds read or write lock
	// DEBUG: Watch c.mu.state to see lock contention
	c.mu.Lock()
	defer c.mu.Unlock()

	// BREAKPOINT: Set breakpoint to watch cache entry creation
	// DEBUG: Watch expiration being calculated: now + ttl
	// DEBUG: Watch entry being stored in map
	c.entries[key] = cacheEntry{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

// Get implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c *Cache) Get(key string) (interface{}, bool) {
	// BREAKPOINT: Set breakpoint to observe read lock acquisition
	// DEBUG: RLock() allows multiple concurrent readers
	// DEBUG: Blocks only if writer holds lock
	c.mu.RLock()
	defer c.mu.RUnlock()

	// BREAKPOINT: Set breakpoint to watch map lookup
	// DEBUG: Watch exists flag - false if key not in map
	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}

	// BREAKPOINT: Set breakpoint to check expiration
	// DEBUG: Watch now and entry.expiration comparison
	// DEBUG: Entry expired if now > expiration
	if time.Now().After(entry.expiration) {
		// DEBUG: Entry exists but expired - return false
		return nil, false
	}

	// DEBUG: Entry exists and not expired - return value
	return entry.value, true
}

// Cleanup implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c *Cache) Cleanup(ctx context.Context) {
	// BREAKPOINT: Set breakpoint to observe ticker creation
	// DEBUG: Ticker fires every 100ms for periodic cleanup
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// BREAKPOINT: Set breakpoint in cleanup loop
	// DEBUG: Loop runs periodically until context cancelled
	for {
		select {
		case <-ticker.C:
			// DEBUG: Ticker fired - time to cleanup expired entries
			// DEBUG: removeExpired() acquires write lock and scans map
			c.removeExpired()

		case <-ctx.Done():
			// DEBUG: Context cancelled - stop cleanup goroutine
			// DEBUG: ticker.Stop() prevents goroutine leak
			return
		}
	}
}

// removeExpired implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c *Cache) removeExpired() {
	// BREAKPOINT: Set breakpoint to watch cleanup acquire write lock
	// DEBUG: Lock() ensures exclusive access during cleanup
	c.mu.Lock()
	defer c.mu.Unlock()

	// BREAKPOINT: Set breakpoint to observe expiration scan
	// DEBUG: Watch now timestamp
	// DEBUG: Watch entries being deleted as we iterate
	now := time.Now()
	for key, entry := range c.entries {
		// DEBUG: Compare entry expiration with current time
		// DEBUG: Watch delete() removing expired entries
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

// NewRateLimiter implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func NewRateLimiter(rate int) *RateLimiter {
	// BREAKPOINT: Set breakpoint to observe token bucket creation
	// DEBUG: Watch rate value and channel capacity (same value)
	// DEBUG: Buffered channel acts as token bucket
	rl := &RateLimiter{
		tokens: make(chan struct{}, rate),
	}

	// BREAKPOINT: Set breakpoint in fill loop
	// DEBUG: Watch tokens being added to channel (filling bucket)
	// DEBUG: Loop runs 'rate' times to fill bucket initially
	for i := 0; i < rate; i++ {
		rl.tokens <- struct{}{}
	}

	// BREAKPOINT: Set breakpoint before refill goroutine launch
	// DEBUG: Refill goroutine adds tokens at constant rate
	go func() {
		// DEBUG: Ticker interval = 1 second / rate
		// DEBUG: For rate=10, ticker fires every 100ms
		ticker := time.NewTicker(time.Second / time.Duration(rate))
		defer ticker.Stop()

		// BREAKPOINT: Set breakpoint in refill loop
		// DEBUG: Loop adds one token per ticker fire
		for range ticker.C {
			// DEBUG: Non-blocking send prevents overfilling bucket
			// DEBUG: Watch select choosing case vs default
			select {
			case rl.tokens <- struct{}{}:
				// DEBUG: Token added to bucket
				// DEBUG: Watch channel length increasing
			default:
				// DEBUG: Bucket full (at capacity), skip token
				// DEBUG: This prevents overflow beyond rate limit
			}
		}
	}()

	return rl
}

// Wait implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (rl *RateLimiter) Wait(ctx context.Context) error {
	// BREAKPOINT: Set breakpoint on select to watch token consumption
	// DEBUG: Select waits for either token available or context cancelled
	// DEBUG: Watch rl.tokens channel - blocks if empty
	select {
	case <-rl.tokens:
		// DEBUG: Token consumed from bucket
		// DEBUG: Operation allowed to proceed
		return nil
	case <-ctx.Done():
		// DEBUG: Context cancelled while waiting for token
		// DEBUG: Return cancellation error, don't consume token
		return ctx.Err()
	}
}

/*
Common Implementation Patterns:

1. Token Bucket (used here):
   - Buffered channel as token container
   - Goroutine refills at constant rate
   - Allows bursts up to bucket capacity
   - Context-aware with select statement

2. Ticker-Based Limiting:
   - Use time.Ticker for strict rate
   - No burst allowance
   - Simpler but less flexible

3. Semaphore Pattern:
   - Can use golang.org/x/sync/semaphore
   - Good for limiting concurrency
   - External dependency required

Critical Implementation Details:

- Always buffer token channel (capacity = rate)
- Use select with default for non-blocking refill
- Handle context cancellation in Wait()
- Call defer cancel() to prevent leaks
- Never store context in struct fields

Debugging Tips:

- Watch token channel length to see available capacity
- Observe ticker firing rate (should match rate parameter)
- Track context.Done() channel for cancellation events
- Monitor goroutine count to detect leaks
- Use -race flag to detect concurrent access issues
*/
