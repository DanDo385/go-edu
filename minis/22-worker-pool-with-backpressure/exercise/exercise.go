//go:build !solution
// +build !solution

package exercise

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
	// TODO: Add fields here
}

// NewWorkerPool creates a new worker pool
// TODO: Initialize the worker pool with bounded queues
//
// Parameters:
// - queueSize: Maximum number of jobs that can be queued (backpressure threshold)
//   * If producers exceed this, Submit() will return error
//   * Prevents unbounded memory growth
//
// - numWorkers: Number of concurrent workers (parallelism level)
//   * More workers = higher throughput (up to CPU limit)
//   * Too many = context switching overhead
//   * Rule of thumb: numWorkers ≈ number of CPU cores
//
// Memory allocation:
// - jobs channel: heap allocation for buffer
// - results channel: heap allocation for buffer
// - WorkerPool struct: escapes to heap (we return pointer)
//
// Why buffered channels?
// - jobs: Bounded queue enforces backpressure
// - results: Prevents workers from blocking on send
//   * Workers can keep processing while consumer drains results
//
// func NewWorkerPool(queueSize, numWorkers int) *WorkerPool {
//     return &WorkerPool{
//         jobs:       make(chan Job, queueSize),    // Bounded job queue
//         results:    make(chan Result, queueSize), // Buffered results
//         numWorkers: numWorkers,                   // Store for Start()
//     }
// }

func NewWorkerPool(queueSize, numWorkers int) *WorkerPool {
	// TODO: Implement
	return nil
}

// Start begins processing jobs
// TODO: Start worker goroutines that process jobs from the queue
//
// Parameters:
// - ctx: Context for cancellation (allows stopping all workers)
//   * When ctx.Done() closes, workers should exit
//   * Enables graceful shutdown
//
// - process: Function to process each job (user-provided logic)
//   * Takes: Job
//   * Returns: Result
//   * May be slow (I/O, computation, etc.)
//
// Architecture:
// - Spawns numWorkers goroutines
// - Each worker independently reads from jobs channel
// - Shared jobs channel provides load balancing (workers pull when ready)
// - Workers stop on context cancellation or channel close
//
// Worker lifecycle:
// 1. Worker starts in goroutine
// 2. Loop:
//    a. Wait for job from jobs channel OR context cancellation
//    b. Process job (call user-provided function)
//    c. Send result to results channel
// 3. Worker exits when:
//    - jobs channel is closed AND drained
//    - context is cancelled
// 4. wg.Done() signals completion
//
// Concurrency patterns:
// - Fan-in: Multiple workers send to single results channel
// - Context cancellation: Clean shutdown without deadlock
// - WaitGroup: Know when all workers have finished
//
// Why select for both receive and send?
// - jobs receive: Respect context cancellation while waiting for work
// - results send: Respect context cancellation while sending result
// - Without select: Workers could block forever on cancelled context
//
// func (p *WorkerPool) Start(ctx context.Context, process func(Job) Result) {
//     // Start worker goroutines
//     for i := 0; i < p.numWorkers; i++ {
//         p.wg.Add(1)
//         go func(workerID int) {
//             defer p.wg.Done()
//
//             // Worker loop: process jobs until stopped
//             for {
//                 select {
//                 case <-ctx.Done():
//                     // Context cancelled, stop immediately
//                     return
//                 case job, ok := <-p.jobs:
//                     if !ok {
//                         // Jobs channel closed, no more work
//                         return
//                     }
//
//                     // Process job (this is where the actual work happens)
//                     result := process(job)
//
//                     // Send result (use select to respect cancellation)
//                     select {
//                     case <-ctx.Done():
//                         // Context cancelled while sending
//                         return
//                     case p.results <- result:
//                         // Result sent successfully
//                     }
//                 }
//             }
//         }(i)
//     }
//
//     // Close results channel when all workers finish
//     // This allows consumers to range over results channel
//     go func() {
//         p.wg.Wait()        // Wait for all workers to exit
//         close(p.results)   // Signal no more results
//     }()
// }

func (p *WorkerPool) Start(ctx context.Context, process func(Job) Result) {
	// TODO: Implement
}

// Submit attempts to add a job to the queue (non-blocking)
// TODO: Implement non-blocking job submission with backpressure
//
// Backpressure strategy: REJECT (fail fast)
// - If queue has space: Add job, return nil
// - If queue is full: Reject immediately, return ErrQueueFull
// - NEVER blocks caller
//
// Why non-blocking?
// - Caller can decide what to do: retry, drop, defer, log, etc.
// - HTTP servers can return 503 Service Unavailable
// - Message producers can apply exponential backoff
// - Prevents cascading slowdowns
//
// Pattern: select with default
//   select {
//   case p.jobs <- job:
//       return nil  // Job queued successfully
//   default:
//       return ErrQueueFull  // Queue full, backpressure applied
//   }
//
// How select with default works:
// 1. Try to send job to channel
// 2. If channel has space: Send succeeds, return nil
// 3. If channel is full: Don't block, execute default case
// 4. Return error to signal backpressure
//
// This is O(1) time complexity (no blocking, no loops)
//
// func (p *WorkerPool) Submit(job Job) error {
//     select {
//     case p.jobs <- job:
//         return nil // Job submitted successfully
//     default:
//         return ErrQueueFull // Queue full, backpressure signal
//     }
// }

func (p *WorkerPool) Submit(job Job) error {
	// TODO: Implement
	return nil
}

// SubmitWithTimeout attempts to add a job with a timeout
// TODO: Implement timeout-based job submission
//
// Backpressure strategy: TIMEOUT (wait with limit)
// - Wait up to timeout duration for queue space
// - If space becomes available: Add job, return nil
// - If timeout expires: Return ErrQueueFull
// - Respects context cancellation
//
// Why timeout-based?
// - More patient than reject (gives system time to drain)
// - Less aggressive than blocking forever
// - Caller controls wait duration (can use exponential backoff)
//
// Pattern: select with multiple cases
//   select {
//   case p.jobs <- job:
//       return nil  // Submitted immediately
//   case <-time.After(timeout):
//       return ErrQueueFull  // Timeout expired
//   case <-ctx.Done():
//       return ctx.Err()  // Context cancelled
//   }
//
// Select evaluation order:
// 1. If jobs channel has space: Send immediately
// 2. If timeout timer fires: Return timeout error
// 3. If context is cancelled: Return context error
// 4. If multiple ready: Go picks one pseudo-randomly
//
// Memory note:
// - time.After() creates a timer (small heap allocation)
// - Timer is garbage collected after select completes
// - For high-frequency calls, consider time.NewTimer (reusable)
//
// func (p *WorkerPool) SubmitWithTimeout(ctx context.Context, job Job, timeout time.Duration) error {
//     select {
//     case p.jobs <- job:
//         return nil // Submitted immediately
//     case <-time.After(timeout):
//         // Timeout expired, queue still full
//         return ErrQueueFull
//     case <-ctx.Done():
//         // Context cancelled (e.g., request timeout)
//         return ctx.Err()
//     }
// }

func (p *WorkerPool) SubmitWithTimeout(ctx context.Context, job Job, timeout time.Duration) error {
	// TODO: Implement
	return nil
}

// Results returns a read-only channel of results
// TODO: Return the results channel for consumers
//
// Returns: <-chan Result (receive-only channel)
//
// Why read-only channel type?
// - Encapsulation: Caller can't send to this channel
// - Type safety: Compiler enforces correct usage
// - Only worker pool can send results
//
// Usage pattern:
//   for result := range pool.Results() {
//       if result.Err != nil {
//           log.Printf("Job %d failed: %v", result.JobID, result.Err)
//       } else {
//           log.Printf("Job %d: %s", result.JobID, result.Data)
//       }
//   }
//
// This range loop exits when results channel is closed
// (which happens when all workers finish and wg.Wait() completes)
//
// func (p *WorkerPool) Results() <-chan Result {
//     return p.results
// }

func (p *WorkerPool) Results() <-chan Result {
	// TODO: Implement
	return nil
}

// Close signals no more jobs will be submitted
// TODO: Close the jobs channel to signal completion
//
// Behavior:
// - Closes jobs channel
// - Workers finish processing queued jobs
// - Workers exit when jobs channel is drained
// - Results channel closes when all workers exit (via Start's goroutine)
//
// Why close jobs channel?
// - Signals to workers: "No more work coming"
// - Workers can drain remaining jobs and exit
// - Range over jobs channel will exit when closed
//
// Graceful shutdown sequence:
// 1. Caller calls Close()
// 2. jobs channel is closed
// 3. Workers process remaining jobs in queue
// 4. Workers see closed channel (ok=false) and exit
// 5. wg.Wait() completes when all workers exit
// 6. results channel is closed (by goroutine in Start)
// 7. Consumer's range loop exits
//
// Important: Only the sender should close a channel!
// - Worker pool owns jobs channel, so it closes it
// - Never close a channel from receiver side (causes panic)
//
// func (p *WorkerPool) Close() {
//     close(p.jobs)
// }

func (p *WorkerPool) Close() {
	// TODO: Implement
}

// QueueDepth returns current number of jobs in queue
// TODO: Implement queue depth monitoring
//
// Returns: Number of jobs waiting to be processed
//
// How len() works on channels:
// - Returns number of elements currently in channel buffer
// - O(1) operation (channel tracks this internally)
// - Safe to call concurrently
// - Snapshot in time (may change immediately)
//
// Use cases:
// - Monitoring/metrics: Track queue utilization over time
// - Adaptive scaling: Add workers if queue consistently full
// - Alerting: Alert when queue > 80% full for extended period
// - Debugging: Understand system behavior
//
// Example:
//   depth := pool.QueueDepth()
//   if depth > queueSize * 0.8 {
//       log.Warn("Queue approaching capacity: %d/%d", depth, queueSize)
//   }
//
// func (p *WorkerPool) QueueDepth() int {
//     return len(p.jobs)
// }

func (p *WorkerPool) QueueDepth() int {
	// TODO: Implement
	return 0
}

// QueueUtilization returns queue fullness as a percentage (0.0 to 1.0)
// TODO: Implement queue utilization calculation
//
// Returns: Ratio of (current depth / capacity)
//
// Useful for:
// - Dashboards: Visual representation of queue health
// - Auto-scaling: Scale workers when utilization > threshold
// - Capacity planning: Historical data shows if queue is sized right
//
// Example thresholds:
// - < 0.5: System healthy, plenty of capacity
// - 0.5-0.8: Moderate load, monitor closely
// - > 0.8: High load, consider scaling
// - = 1.0: Queue full, backpressure in effect
//
// func (p *WorkerPool) QueueUtilization() float64 {
//     return float64(len(p.jobs)) / float64(cap(p.jobs))
// }

func (p *WorkerPool) QueueUtilization() float64 {
	// TODO: Implement
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
	// TODO: Add fields here
}

// NewRateLimiter creates a rate limiter
// TODO: Initialize rate limiter and start refill goroutine
//
// Parameters:
// - requestsPerSecond: Maximum requests allowed per second
//   * Also the burst capacity (can do N requests immediately)
//   * Refill rate = 1 token per (1s / N)
//
// Initialization:
// 1. Create token channel with capacity = requestsPerSecond
// 2. Fill bucket initially (allow immediate burst)
// 3. Calculate refill rate
// 4. Start background refill goroutine
//
// Why start with full bucket?
// - Allows immediate burst of requests
// - User doesn't wait for first token
// - More user-friendly behavior
//
// Refill goroutine:
// - Runs in background
// - Adds one token every (1s / requestsPerSecond)
// - Stops when stop channel is closed
//
// func NewRateLimiter(requestsPerSecond int) *RateLimiter {
//     rl := &RateLimiter{
//         tokens:   make(chan struct{}, requestsPerSecond),
//         rate:     time.Second / time.Duration(requestsPerSecond),
//         capacity: requestsPerSecond,
//         stop:     make(chan struct{}),
//     }
//
//     // Fill bucket initially (allow burst)
//     for i := 0; i < requestsPerSecond; i++ {
//         rl.tokens <- struct{}{}
//     }
//
//     // Start background refill goroutine
//     rl.wg.Add(1)
//     go func() {
//         defer rl.wg.Done()
//
//         ticker := time.NewTicker(rl.rate)
//         defer ticker.Stop()
//
//         for {
//             select {
//             case <-rl.stop:
//                 return // Shutdown signal
//             case <-ticker.C:
//                 // Try to add token (non-blocking)
//                 select {
//                 case rl.tokens <- struct{}{}:
//                     // Token added
//                 default:
//                     // Bucket full, drop token
//                 }
//             }
//         }
//     }()
//
//     return rl
// }

func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	// TODO: Implement
	return nil
}

// Wait blocks until a token is available or context is cancelled
// TODO: Implement blocking token acquisition
//
// Behavior:
// - Consumes one token from bucket
// - Blocks if no tokens available (rate limit enforcement)
// - Returns immediately if token available
// - Returns error if context is cancelled
//
// Pattern: select with channel receive
//   select {
//   case <-rl.tokens:
//       return nil  // Got token, operation allowed
//   case <-ctx.Done():
//       return ctx.Err()  // Context cancelled
//   }
//
// Use case:
// - Critical operations that must eventually succeed
// - Willing to wait for rate limit to allow operation
// - HTTP requests, database queries, API calls
//
// Example:
//   if err := limiter.Wait(ctx); err != nil {
//       return err  // Context cancelled or deadline exceeded
//   }
//   // Perform rate-limited operation
//   makeAPICall()
//
// func (rl *RateLimiter) Wait(ctx context.Context) error {
//     select {
//     case <-rl.tokens:
//         return nil // Token acquired
//     case <-ctx.Done():
//         return ctx.Err() // Context cancelled
//     }
// }

func (rl *RateLimiter) Wait(ctx context.Context) error {
	// TODO: Implement
	return nil
}

// TryAcquire attempts to get a token without blocking
// TODO: Implement non-blocking token acquisition
//
// Behavior:
// - Returns immediately
// - Consumes token if available (returns true)
// - Returns false if no tokens (rate limited)
//
// Pattern: select with default
//   select {
//   case <-rl.tokens:
//       return true  // Token acquired
//   default:
//       return false  // No tokens available
//   }
//
// Use case:
// - Best-effort operations (okay to skip if rate limited)
// - Metrics collection (sample when not overloaded)
// - Background tasks (can defer to later)
//
// Example:
//   if limiter.TryAcquire() {
//       sendMetrics()  // Send if allowed
//   }
//   // Otherwise skip this round
//
// func (rl *RateLimiter) TryAcquire() bool {
//     select {
//     case <-rl.tokens:
//         return true  // Token acquired
//     default:
//         return false  // Rate limited
//     }
// }

func (rl *RateLimiter) TryAcquire() bool {
	// TODO: Implement
	return false
}

// Stop stops the rate limiter's token refill goroutine
// TODO: Implement graceful shutdown
//
// Behavior:
// - Signals refill goroutine to stop
// - Waits for goroutine to exit
// - Should be called when rate limiter is no longer needed
//
// Shutdown sequence:
// 1. Close stop channel (broadcast to refill goroutine)
// 2. Refill goroutine sees closed channel in select
// 3. Refill goroutine exits
// 4. wg.Wait() returns when goroutine completes
//
// Why close instead of send?
// - Closing broadcasts to all listeners
// - Safe to close multiple times (panic on second close)
// - Idiomatic Go shutdown pattern
//
// func (rl *RateLimiter) Stop() {
//     close(rl.stop)
//     rl.wg.Wait()
// }

func (rl *RateLimiter) Stop() {
	// TODO: Implement
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
