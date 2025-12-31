//go:build !solution
// +build !solution

package workerpoolwithbackpressure

// ============================================================================
// WORKER POOL WITH BACKPRESSURE: Understanding Queue Management and Flow Control
// ============================================================================
//
// A worker pool is a concurrency pattern that:
// 1. Maintains a fixed number of worker goroutines
// 2. Processes jobs from a queue (channel)
// 3. Enforces backpressure when queue is full
// 4. Prevents unbounded memory growth
//
// Why backpressure matters:
// - Without it: Queue grows unbounded, memory exhaustion
// - With it: Producers slow down when consumers can't keep up
// - Real-world: HTTP servers return 503, message queues block
//
// Backpressure strategies:
// 1. REJECT: Fail fast (non-blocking send with select+default)
// 2. TIMEOUT: Wait with limit (select with time.After)
// 3. BLOCK: Wait indefinitely (regular channel send)
//
// Go advantages:
// - Buffered channels provide bounded queues naturally
// - Select statement enables non-blocking operations
// - Context propagation for cancellation
// - Goroutines are cheap (can have thousands)
//
// Memory considerations:
// - Buffered channel allocates queue array on heap
// - Size: (queueSize * sizeof(Job)) bytes
// - Each goroutine: ~2KB stack (grows as needed)
// - Total: queue memory + (numWorkers * 2KB)
//
// Real-world applications:
// - HTTP servers (prevent overload, return 503)
// - Message queue consumers (bounded processing)
// - Database connection pools (limit concurrent queries)
// - API clients (respect rate limits)
//
// ============================================================================

// TODO: Import required packages
// You'll need:
// - "context" for cancellation
// - "sync" for WaitGroup
// - "time" for timeout support
//
// import (
//     "context"
//     "sync"
//     "time"
// )

// ============================================================================
// Exercise 1: Worker Pool with Backpressure
// ============================================================================

// Job represents work to be processed
// TODO: This type is already defined, just understand the structure
//
// Memory layout:
// - ID: 8 bytes (int)
// - Payload: 16 bytes (string header: ptr + len)
// - Total: ~24 bytes per job (plus string data on heap)
type Job struct {
	ID      int
	Payload string
}

// Result represents the outcome of processing a job
// TODO: This type is already defined, understand error handling
//
// Memory layout:
// - JobID: 8 bytes (int)
// - Data: 16 bytes (string header)
// - Err: 16 bytes (interface: type + ptr)
// - Total: ~40 bytes per result (plus string/error data)
type Result struct {
	JobID int
	Data  string
	Err   error
}

// WorkerPool manages a bounded pool of workers with backpressure
// TODO: Add fields for the worker pool
//
// Required fields:
// - jobs: Buffered channel for pending jobs (bounded queue)
//   * Type: chan Job
//   * Size: queueSize (enforces backpressure)
//   * Memory: ~queueSize * 24 bytes
//
// - results: Buffered channel for completed results
//   * Type: chan Result
//   * Size: queueSize (prevents worker blocking on send)
//   * Memory: ~queueSize * 40 bytes
//
// - numWorkers: Number of concurrent workers
//   * Type: int
//   * Determines parallelism level
//
// - wg: WaitGroup to track worker completion
//   * Type: sync.WaitGroup
//   * Coordinates graceful shutdown
//
// Why pointer receiver?
// - Pool must be shared across goroutines
// - Cannot copy WaitGroup or channels (breaks synchronization)
// - Methods modify pool state
//
// type WorkerPool struct {
//     jobs       chan Job       // Bounded job queue (backpressure point)
//     results    chan Result    // Result channel (buffered to prevent blocking)
//     numWorkers int            // Worker count (parallelism level)
//     wg         sync.WaitGroup // Tracks worker goroutines
// }

type WorkerPool struct {
	jobs       chan Job
	results    chan Result
	numWorkers int
	wg         sync.WaitGroup
}

func NewWorkerPool(queueSize, numWorkers int) *WorkerPool {
	// TODO: Implement this function.
	// - Initialize and return a new `*WorkerPool`.
	// - `jobs`: This should be a buffered channel of `Job` with a capacity of `queueSize`. This is the core of your bounded queue.
	// - `results`: This should also be a buffered channel of `Result` with a capacity of `queueSize`. A buffer here is important so that workers don't block waiting for a consumer to read the result.
	// - `numWorkers`: Store this for the `Start` method.
	return nil
}

func (p *WorkerPool) Start(ctx context.Context, process func(Job) Result) {
	// TODO: Implement this function.

	// This method launches the worker goroutines.

	// Step 1: Loop `p.numWorkers` times to start each worker.
	// - In each iteration, call `p.wg.Add(1)` and launch a goroutine.

	// Step 2: Implement the worker's logic inside the goroutine.
	// - The goroutine should `defer p.wg.Done()` to signal completion.
	// - Use an infinite `for { ... }` loop and a `select` statement.
	// - `case <-ctx.Done():` -> The context was cancelled, so `return` to stop the worker.
	// - `case job, ok := <-p.jobs:`
	//   - If `!ok`, the channel is closed and empty. `return` to stop the worker.
	//   - Otherwise, call the `process(job)` function.
	//   - Send the result to the `p.results` channel. **Important:** This send must also be in a `select` with a `case <-ctx.Done()` to prevent the worker from blocking forever if the context is cancelled while it's trying to send a result.

	// Step 3: Start a final goroutine to close the `results` channel.
	// - This goroutine waits for all workers to finish and then closes the `results` channel.
	// - `go func() { p.wg.Wait(); close(p.results) }()`
	// - This is what allows a consumer to safely `range` over the `Results()` channel.
}

func (p *WorkerPool) Submit(job Job) error {
	// TODO: Implement this function.

	// This method implements the "REJECT" backpressure strategy.

	// Step 1: Use a `select` statement with a `default` case.
	// - `case p.jobs <- job:` -> The job was successfully sent. Return `nil`.
	// - `default:` -> The `jobs` channel is full, so the send would block. Return `ErrQueueFull`.
	return nil
}

func (p *WorkerPool) SubmitWithTimeout(ctx context.Context, job Job, timeout time.Duration) error {
	// TODO: Implement this function.

	// This method implements the "TIMEOUT" backpressure strategy.

	// Step 1: Use a `select` statement with three cases.
	// - `case p.jobs <- job:` -> The job was sent successfully. Return `nil`.
	// - `case <-time.After(timeout):` -> The timeout was reached. Return `ErrQueueFull`.
	// - `case <-ctx.Done():` -> The context was cancelled. Return `ctx.Err()`.
	return nil
}

func (p *WorkerPool) Results() <-chan Result {
	// TODO: Implement this function.
	// - Simply return the `p.results` channel.
	// - The return type is `<-chan Result`, a receive-only channel, which prevents the caller from closing or sending to the channel, enforcing encapsulation.
	return nil
}

func (p *WorkerPool) Close() {
	// TODO: Implement this function.
	// - Simply close the `p.jobs` channel.
	// - This is the signal that no more jobs will be submitted. The workers will drain the channel and then exit.
	// - **Rule:** Only the sender should close a channel. In this design, the `WorkerPool` is the sender to the `jobs` channel.
}

func (p *WorkerPool) QueueDepth() int {
	// TODO: Implement this function.
	// - Return the number of items currently in the `jobs` channel's buffer.
	// - `len(p.jobs)` does this. It's a non-blocking, O(1) operation.
	return 0
}

func (p *WorkerPool) QueueUtilization() float64 {
	// TODO: Implement this function.
	// - Return the queue depth divided by the queue capacity.
	// - `float64(len(p.jobs)) / float64(cap(p.jobs))`
	// - `cap(p.jobs)` returns the buffer capacity of the channel.
	return 0.0
}

// ============================================================================
// Exercise 2: Token Bucket Rate Limiter
// ============================================================================

// RateLimiter implements token bucket rate limiting
// TODO: Add fields for rate limiting
//
// Token bucket algorithm:
// - Bucket holds tokens (each token = permission for one operation)
// - Tokens are consumed when operation is performed
// - Tokens are refilled at constant rate
// - Bucket has maximum capacity (allows burst)
//
// Required fields:
// - tokens: Buffered channel holding available tokens
//   * Type: chan struct{}
//   * Capacity: requestsPerSecond (max burst)
//   * struct{} is zero-size (memory efficient)
//
// - rate: Time between token refills
//   * Type: time.Duration
//   * = 1 second / requestsPerSecond
//
// - capacity: Maximum tokens in bucket
//   * Type: int
//   * Same as channel capacity
//
// - stop: Channel to signal refill goroutine to stop
//   * Type: chan struct{}
//   * Closed to signal shutdown
//
// - wg: WaitGroup to wait for refill goroutine
//   * Type: sync.WaitGroup
//   * Ensures clean shutdown
//
// Why buffered channel for tokens?
// - Buffer size = max tokens (burst capacity)
// - Sending = consuming token (blocks if empty)
// - Receiving in refill = adding token (blocks if full)
//
// Example: 10 requests/second
// - Bucket capacity: 10 tokens
// - Refill rate: 1 token per 100ms
// - Can burst 10 requests immediately
// - Sustained rate: 10 req/sec
//
// type RateLimiter struct {
//     tokens   chan struct{}  // Token bucket (buffered = capacity)
//     rate     time.Duration  // Time between refills
//     capacity int            // Max tokens in bucket
//     stop     chan struct{}  // Signal to stop refill goroutine
//     wg       sync.WaitGroup // Wait for refill goroutine to exit
// }

type RateLimiter struct {
	tokens   chan struct{}
	rate     time.Duration
	capacity int
	stop     chan struct{}
	wg       sync.WaitGroup
}

func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	// TODO: Implement this function.

	// This function creates a token bucket rate limiter.

	// Step 1: Initialize the RateLimiter struct.
	// - `tokens`: A buffered channel of `struct{}` with capacity `requestsPerSecond`. `struct{}` is used because it has zero size.
	// - `rate`: The duration between token refills, calculated as `time.Second / time.Duration(requestsPerSecond)`.
	// - `stop`: A channel for signaling the refill goroutine to stop.
	// - Also store `capacity` and a `sync.WaitGroup`.

	// Step 2: Fill the token bucket initially.
	// - Loop `requestsPerSecond` times and send an empty struct to the `tokens` channel. This allows for an initial burst of requests.

	// Step 3: Start the background refill goroutine.
	// - `rl.wg.Add(1)`
	// - `go func() { ... }()`
	// - The goroutine should `defer rl.wg.Done()`.
	// - Create a `time.Ticker` with the calculated `rate`. `defer ticker.Stop()`.
	// - Loop forever with a `select` statement.
	//   - `case <-rl.stop:` -> `return` to stop the goroutine.
	//   - `case <-ticker.C:` -> Time to add a new token. Use a non-blocking send (`select` with `default`) to add a token to the `tokens` channel, but only if it's not already full.

	return nil
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
	// TODO: Implement this function.

	// This is the blocking method to acquire a token.

	// Step 1: Use a `select` to wait for a token or for the context to be cancelled.
	// - `case <-rl.tokens:` -> A token was acquired. Return `nil`.
	// - `case <-ctx.Done():` -> The context was cancelled. Return `ctx.Err()`.
	return nil
}

func (rl *RateLimiter) TryAcquire() bool {
	// TODO: Implement this function.

	// This is the non-blocking method to acquire a token.

	// Step 1: Use a `select` with a `default` case.
	// - `case <-rl.tokens:` -> A token was acquired. Return `true`.
	// - `default:` -> No token was available. Return `false`.
	return false
}

func (rl *RateLimiter) Stop() {
	// TODO: Implement this function.

	// This gracefully stops the refill goroutine.

	// Step 1: Close the `stop` channel.
	// - `close(rl.stop)`
	// - This will be received by the `case <-rl.stop:` in the refill goroutine, causing it to exit.

	// Step 2: Wait for the goroutine to finish.
	// - `rl.wg.Wait()`
	// - This ensures that the goroutine has fully stopped before this function returns.
}

// Common errors
var (
	ErrQueueFull = &QueueFullError{}
)

// QueueFullError indicates the queue is at capacity
type QueueFullError struct{}

func (e *QueueFullError) Error() string {
	return "queue is full"
}

// ============================================================================
// After implementing all functions:
// - Run: go test -v ./... to verify correctness
// - Run: go test -race ./... to check for data races
// - Run: go test -bench=. to measure performance
// - Compare with solution.go to see optimizations
//
// Key learnings:
// 1. Bounded channels enforce backpressure naturally
// 2. Select with default enables non-blocking operations
// 3. Context propagation allows graceful cancellation
// 4. WaitGroup coordinates goroutine completion
// 5. Closing channels signals "no more data"
// 6. Token bucket rate limiting uses channels elegantly
//
// Common mistakes:
// - Forgetting to close channels (goroutine leaks)
// - Closing channels from receiver (panic)
// - Not respecting context cancellation (hangs)
// - Unbuffered results channel (workers block on send)
// - Not using select with timeout (operations hang forever)
// ============================================================================
