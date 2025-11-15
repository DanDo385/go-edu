//go:build solution
// +build solution

/*
Problem: Worker pool with backpressure and rate limiting

We need to:
1. Implement a bounded worker pool that prevents unbounded queue growth
2. Provide backpressure when queue is full (reject or timeout)
3. Support non-blocking submission (fail fast when full)
4. Support timeout-based submission (wait with limit)
5. Implement token bucket rate limiting
6. Gracefully handle context cancellation

Constraints:
- Queue size must be bounded (prevent memory exhaustion)
- Workers must respect context cancellation
- Non-blocking operations must use select with default
- Rate limiter must enforce throughput limits

Time/Space Complexity:
- Submit: O(1) - channel send or immediate return
- Worker processing: O(1) per job
- Rate limiter: O(1) per token acquisition
- Space: O(queueSize + numWorkers) for channels and goroutines

Why Go is well-suited:
- Buffered channels provide built-in bounded queues with blocking
- Select statement enables non-blocking operations
- Context propagation for cancellation
- Lightweight goroutines for workers

Real-world applications:
- HTTP servers (prevent overload with 503 responses)
- Message queue consumers (acknowledge only when processed)
- Database connection pools (bounded connections)
- API rate limiting (comply with third-party limits)
*/

package exercise

import (
	"context"
	"sync"
	"time"
)

// Job represents work to be processed
type Job struct {
	ID      int
	Payload string
}

// Result represents the outcome of processing a job
type Result struct {
	JobID int
	Data  string
	Err   error
}

// WorkerPool manages a bounded pool of workers with backpressure
type WorkerPool struct {
	jobs       chan Job
	results    chan Result
	numWorkers int
	wg         sync.WaitGroup
}

// NewWorkerPool creates a new worker pool
//
// Go Concepts Demonstrated:
// - Buffered channels with fixed capacity
// - Struct initialization
// - Resource allocation
//
// Parameters:
//   - queueSize: Maximum number of jobs that can be queued (backpressure threshold)
//   - numWorkers: Number of concurrent workers (parallelism level)
//
// Design decisions:
//   - jobs channel size = queueSize (enforces backpressure)
//   - results channel size = queueSize (prevents worker blocking on send)
//   - Store numWorkers to start workers later
func NewWorkerPool(queueSize, numWorkers int) *WorkerPool {
	return &WorkerPool{
		jobs:       make(chan Job, queueSize),
		results:    make(chan Result, queueSize),
		numWorkers: numWorkers,
	}
}

// Start begins processing jobs
//
// Go Concepts Demonstrated:
// - Goroutines for concurrent execution
// - WaitGroup for coordinating goroutine completion
// - Select for context cancellation
// - Defer for cleanup
//
// Architecture:
//   - Spawns numWorkers goroutines
//   - Each worker processes jobs from shared channel
//   - Workers stop on context cancellation or channel close
//   - Results channel closes when all workers finish
//
// Parameters:
//   - ctx: Context for cancellation (allows stopping all workers)
//   - process: Function to process each job (user-provided logic)
func (p *WorkerPool) Start(ctx context.Context, process func(Job) Result) {
	// Start worker goroutines
	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1)
		go func(workerID int) {
			defer p.wg.Done()

			// Worker loop: process jobs until stopped
			for {
				select {
				case <-ctx.Done():
					// Context cancelled, stop immediately
					return
				case job, ok := <-p.jobs:
					if !ok {
						// Jobs channel closed, no more work
						return
					}

					// Process job (could panic, so production code would use recover)
					result := process(job)

					// Send result (use select to respect cancellation)
					select {
					case <-ctx.Done():
						return
					case p.results <- result:
						// Result sent successfully
					}
				}
			}
		}(i)
	}

	// Close results channel when all workers finish
	// This allows consumers to range over results channel
	go func() {
		p.wg.Wait()
		close(p.results)
	}()
}

// Submit attempts to add a job to the queue (non-blocking)
//
// Go Concepts Demonstrated:
// - Select with default (non-blocking channel operation)
// - Error return for backpressure signaling
//
// Backpressure strategy: REJECT
//   - If queue is full, reject immediately
//   - Returns error so caller can decide (retry, drop, defer)
//   - Never blocks (guaranteed O(1) time)
//
// Parameters:
//   - job: The job to submit
//
// Returns:
//   - error: ErrQueueFull if queue is at capacity, nil otherwise
func (p *WorkerPool) Submit(job Job) error {
	select {
	case p.jobs <- job:
		return nil // Submitted successfully
	default:
		return ErrQueueFull // Queue full, backpressure applied
	}
}

// SubmitWithTimeout attempts to add a job with a timeout
//
// Go Concepts Demonstrated:
// - Select with multiple cases (channel send, timeout, cancellation)
// - time.After for timeout
// - Context for cancellation
//
// Backpressure strategy: TIMEOUT
//   - Waits up to timeout duration for space
//   - If space becomes available, submits
//   - If timeout expires, returns error
//   - Respects context cancellation
//
// Parameters:
//   - ctx: Context for cancellation
//   - job: The job to submit
//   - timeout: Maximum time to wait
//
// Returns:
//   - error: ErrQueueFull (timeout), context error (cancelled), or nil (success)
func (p *WorkerPool) SubmitWithTimeout(ctx context.Context, job Job, timeout time.Duration) error {
	select {
	case p.jobs <- job:
		return nil // Submitted immediately
	case <-time.After(timeout):
		// Timeout expired, queue still full
		return ErrQueueFull
	case <-ctx.Done():
		// Context cancelled
		return ctx.Err()
	}
}

// Results returns a read-only channel of results
//
// Go Concepts Demonstrated:
// - Read-only channel type (<-chan)
// - Encapsulation (caller can't send to this channel)
//
// Returns:
//   - <-chan Result: Read-only channel of results
//
// Usage:
//   for result := range pool.Results() {
//       handleResult(result)
//   }
func (p *WorkerPool) Results() <-chan Result {
	return p.results
}

// Close signals no more jobs will be submitted
//
// Go Concepts Demonstrated:
// - Closing channels to signal completion
// - Only sender should close (worker pool owns jobs channel)
//
// Behavior:
//   - Closes jobs channel
//   - Workers finish processing queued jobs
//   - Workers exit when jobs channel is drained
//   - Results channel closes when all workers exit (via Start's goroutine)
func (p *WorkerPool) Close() {
	close(p.jobs)
}

// QueueDepth returns current number of jobs in queue
//
// Go Concepts Demonstrated:
// - len() on channel (non-blocking query)
//
// Returns:
//   - int: Number of jobs currently waiting in queue
//
// Use case:
//   - Monitoring queue depth for metrics/alerting
//   - Adaptive scaling decisions
func (p *WorkerPool) QueueDepth() int {
	return len(p.jobs)
}

// QueueUtilization returns queue fullness as a percentage (0.0 to 1.0)
//
// Go Concepts Demonstrated:
// - len() and cap() on channels
// - Float division
//
// Returns:
//   - float64: Utilization ratio (0.0 = empty, 1.0 = full)
//
// Use case:
//   - Alert when utilization > 0.8 (approaching capacity)
//   - Scale workers when consistently high
func (p *WorkerPool) QueueUtilization() float64 {
	return float64(len(p.jobs)) / float64(cap(p.jobs))
}

// RateLimiter implements token bucket rate limiting
//
// Algorithm: Token Bucket
//   - Bucket holds tokens (permission to make request)
//   - Tokens are consumed when request is made
//   - Tokens are refilled at constant rate
//   - Bucket has maximum capacity (allows burst)
//
// Example: 10 requests/second
//   - Bucket capacity: 10 tokens
//   - Refill rate: 1 token per 100ms
//   - Can burst 10 requests immediately
//   - Sustained rate: 10 req/sec
type RateLimiter struct {
	tokens   chan struct{}
	rate     time.Duration
	capacity int
	stop     chan struct{}
	wg       sync.WaitGroup
}

// NewRateLimiter creates a rate limiter
//
// Go Concepts Demonstrated:
// - Buffered channel as semaphore
// - Ticker for periodic operations
// - Goroutine for background refilling
//
// Parameters:
//   - requestsPerSecond: Maximum requests allowed per second
//
// Returns:
//   - *RateLimiter: Configured rate limiter (starts immediately)
//
// Design:
//   - tokens channel size = requestsPerSecond (max burst)
//   - Refill rate = 1 second / requestsPerSecond
//   - Start with full bucket (allow immediate burst)
func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	rl := &RateLimiter{
		tokens:   make(chan struct{}, requestsPerSecond),
		rate:     time.Second / time.Duration(requestsPerSecond),
		capacity: requestsPerSecond,
		stop:     make(chan struct{}),
	}

	// Fill bucket initially (allow immediate burst)
	for i := 0; i < requestsPerSecond; i++ {
		rl.tokens <- struct{}{}
	}

	// Start background token refiller
	rl.wg.Add(1)
	go func() {
		defer rl.wg.Done()

		ticker := time.NewTicker(rl.rate)
		defer ticker.Stop()

		for {
			select {
			case <-rl.stop:
				return
			case <-ticker.C:
				// Try to add token (non-blocking)
				select {
				case rl.tokens <- struct{}{}:
					// Token added
				default:
					// Bucket full, drop token
				}
			}
		}
	}()

	return rl
}

// Wait blocks until a token is available or context is cancelled
//
// Go Concepts Demonstrated:
// - Blocking receive from channel
// - Select for context cancellation
//
// Behavior:
//   - Consumes one token from bucket
//   - Blocks if no tokens available (rate limit enforcement)
//   - Returns immediately if token available
//   - Returns error if context is cancelled
//
// Parameters:
//   - ctx: Context for cancellation
//
// Returns:
//   - error: Context error if cancelled, nil if token acquired
func (rl *RateLimiter) Wait(ctx context.Context) error {
	select {
	case <-rl.tokens:
		return nil // Got token
	case <-ctx.Done():
		return ctx.Err() // Context cancelled
	}
}

// TryAcquire attempts to get a token without blocking
//
// Go Concepts Demonstrated:
// - Select with default (non-blocking receive)
//
// Behavior:
//   - Returns immediately
//   - Consumes token if available
//   - Returns false if no tokens (rate limited)
//
// Returns:
//   - bool: true if token acquired, false if none available
//
// Use case:
//   - Best-effort operations (drop if rate limited)
//   - Metrics collection (sample when not overloaded)
func (rl *RateLimiter) TryAcquire() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

// Stop stops the rate limiter's token refill goroutine
//
// Go Concepts Demonstrated:
// - Closing channel to signal goroutine
// - WaitGroup for waiting on goroutine completion
//
// Behavior:
//   - Stops token refill
//   - Waits for refiller goroutine to exit
//   - Should be called when rate limiter is no longer needed
func (rl *RateLimiter) Stop() {
	close(rl.stop)
	rl.wg.Wait()
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

/*
Alternatives & Trade-offs:

1. Unbounded queue (no backpressure):
   jobs := make(chan Job)  // No buffer limit
   Pros: Never rejects work, simple
   Cons: Memory exhaustion under load, no flow control
   Go: Buffered channels enforce bounded queues naturally

2. Semaphore-based rate limiting (no refill):
   sem := make(chan struct{}, N)
   Pros: Simpler (no ticker goroutine)
   Cons: No automatic refill, manual release required
   Use case: Connection pools, not rate limiting

3. Global mutex for rate limiting:
   var mu sync.Mutex
   var lastRequest time.Time
   Pros: Lower memory (no channel)
   Cons: Lock contention, harder to reason about
   Go: Channels provide cleaner API

4. Adaptive queue resizing:
   Dynamically grow/shrink queue
   Pros: Handles variable load
   Cons: Complex, still needs max size, breaks backpressure contract
   Go: Fixed-size channels are simpler and more predictable

5. Priority queue with multiple lanes:
   highPriority := make(chan Job, N)
   lowPriority := make(chan Job, N)
   Pros: Differentiate critical vs best-effort work
   Cons: More complex, can starve low priority
   Use case: Multi-tenant systems

Go vs X:

Go vs Java (Executors):
- Java uses ExecutorService with BlockingQueue
- Similar bounded queue concept
- More verbose, heavyweight threads, manual queue management
- Go: Channels integrate queue and synchronization

Go vs Python (asyncio.Queue):
- Python uses asyncio.Queue with maxsize parameter
- Similar non-blocking put (raises QueueFull exception)
- Single-threaded (no true parallelism), exception-based flow
- Go: Multi-threaded, error-based flow (no exceptions)

Go vs Rust (tokio + channel):
- Rust uses tokio mpsc channel with try_send
- Zero-cost abstractions, compile-time safety
- More complex (async traits, Send bounds, lifetimes)
- Go: Simpler, faster development

Go vs Node.js (p-queue):
- Node.js uses p-queue library for concurrency limiting
- Similar queue pattern
- Single-threaded, no true parallelism, requires library
- Go: Built-in, multi-threaded

Real-world examples:

1. NGINX (HTTP server):
   - Bounded connection queue
   - Returns 503 when queue full
   - Our pattern: Submit returns ErrQueueFull

2. Kafka (message broker):
   - Bounded log segments
   - Backpressure via ack delays
   - Our pattern: Workers drain at limited rate

3. AWS Lambda (serverless):
   - Concurrent execution limit
   - Throttles when limit reached
   - Our pattern: Rate limiter enforces limit

4. Database connection pools:
   - Fixed pool size
   - Wait or reject when exhausted
   - Our pattern: SubmitWithTimeout
*/
