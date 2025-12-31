//go:build !solution
// +build !solution

package goroutines1mdemo

// TODO: Import required packages
// You'll need:
// - "context" for graceful shutdown and cancellation
// - "sync" for sync.WaitGroup, sync.Mutex
// - "sync/atomic" for atomic.Int64, atomic.Bool
// - "time" for time.Duration, time.Ticker, time.Sleep
// - "fmt" for error formatting (optional)
//
// import (
//     "context"
//     "sync"
//     "sync/atomic"
//     "time"
// )

// ============================================================================
// GOROUTINES AND CONCURRENCY: Understanding Go's Lightweight Threads
// ============================================================================
//
// What are goroutines?
// - Lightweight threads managed by the Go runtime (not OS threads)
// - Initial stack: 2KB (grows/shrinks dynamically)
// - OS thread: 1-2MB fixed stack (much heavier)
// - Can create millions of goroutines on a single machine
//
// Memory model:
// - Each goroutine has its own stack (on heap, managed by runtime)
// - Goroutines multiplex onto OS threads (M:N threading)
// - GOMAXPROCS controls number of OS threads (default = CPU cores)
// - Scheduler uses work-stealing for load balancing
//
// Why Go can handle millions of goroutines:
// - Small stack size: 2KB vs 1MB for OS threads
// - 1 million goroutines = ~2GB RAM (vs 1TB for OS threads!)
// - Fast context switching (no kernel involvement)
// - Cooperative scheduling (yield at function calls, channel ops, syscalls)
//
// Goroutine lifecycle:
// - Created: go func() { ... }() allocates stack, adds to scheduler
// - Running: Executing on an OS thread
// - Waiting: Blocked on channel, I/O, or sleep
// - Exiting: Function returns, stack deallocated
//
// Common patterns:
// 1. Fire-and-forget: go doWork()
// 2. Wait for completion: use sync.WaitGroup
// 3. Graceful shutdown: use context.Context
// 4. Worker pool: fixed number of goroutines, shared work queue
// 5. Pipeline: chain of goroutines connected by channels
//
// Best practices:
// - Always ensure goroutines can exit (prevent leaks)
// - Use context for cancellation
// - Don't create unlimited goroutines (use worker pools)
// - Synchronize with channels, not shared memory (when possible)
//
// ============================================================================

// ============================================================================
// Exercise 1: ParallelSum
// ============================================================================

// ParallelSum calculates the sum of numbers from 1 to n using multiple goroutines.
// TODO: Implement parallel summation with atomic operations.
//
// Parameters:
//   - n: Upper bound of range (sum 1 to n)
//   - numWorkers: Number of goroutines to use
//
// Returns:
//   - int64: Sum of numbers from 1 to n
//
// Behavior:
//   - Split range [1, n] into numWorkers chunks
//   - Each worker sums its chunk
//   - Workers add their partial sums to shared total (atomically)
//   - Wait for all workers to finish
//   - Return total sum
//
// Example:
//   ParallelSum(100, 4)
//   - Worker 1: sum(1..25) = 325
//   - Worker 2: sum(26..50) = 950
//   - Worker 3: sum(51..75) = 1575
//   - Worker 4: sum(76..100) = 2200
//   - Total: 5050
//
// Architecture:
//   - Create atomic.Int64 for total sum
//   - Create sync.WaitGroup to track worker completion
//   - Calculate chunk size: rangeSize = n / numWorkers
//   - For each worker:
//     * Calculate start and end of its range
//     * Start goroutine that sums its range
//     * Atomically add partial sum to total
//     * Call wg.Done() when finished
//   - Wait for all workers
//   - Return total
//
// Memory:
//   - Each goroutine: ~2KB initial stack
//   - Total for 4 workers: ~8KB
//   - Atomic variable: 8 bytes
//   - WaitGroup: ~16 bytes
//
// Why atomic operations?
//   - Without atomics: race condition when adding to total
//   - Alternative: use mutex (slower for simple operations)
//   - Atomic operations use CPU instructions (LOCK prefix on x86)
//
// func ParallelSum(n int, numWorkers int) int64 {
//     var total atomic.Int64
//     var wg sync.WaitGroup
//
//     rangeSize := n / numWorkers
//     remainder := n % numWorkers
//
//     for i := 0; i < numWorkers; i++ {
//         wg.Add(1)
//
//         start := i*rangeSize + 1
//         end := (i + 1) * rangeSize
//
//         // Give remainder to last worker
//         if i == numWorkers-1 {
//             end += remainder
//         }
//
//         go func(s, e int) {
//             defer wg.Done()
//
//             // Sum this worker's range
//             var sum int64
//             for j := s; j <= e; j++ {
//                 sum += int64(j)
//             }
//
//             // Add to total atomically
//             total.Add(sum)
//         }(start, end)
//     }
//
//     wg.Wait()
//     return total.Load()
// }

func ParallelSum(n int, numWorkers int) int64 {
	// TODO: Implement ParallelSum
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// ============================================================================
// Exercise 2: FanOut
// ============================================================================

// FanOut distributes values from one input channel to multiple output channels.
// TODO: Implement fan-out pattern with round-robin distribution.
//
// Parameters:
//   - input: Input channel to read from
//   - numWorkers: Number of output channels to create
//
// Returns:
//   - []<-chan int: Slice of output channels (read-only)
//
// Behavior:
//   - Create numWorkers output channels
//   - Start distributor goroutine that:
//     * Reads from input channel
//     * Distributes values round-robin to outputs
//     * Closes all outputs when input closes
//   - Return slice of output channels
//
// Round-robin distribution:
//   - Value 0 → output[0]
//   - Value 1 → output[1]
//   - Value 2 → output[2]
//   - Value 3 → output[0] (wraps around)
//
// func FanOut(input <-chan int, numWorkers int) []<-chan int {
//     outputs := make([]chan int, numWorkers)
//     readOnlyOutputs := make([]<-chan int, numWorkers)
//
//     for i := 0; i < numWorkers; i++ {
//         outputs[i] = make(chan int)
//         readOnlyOutputs[i] = outputs[i]
//     }
//
//     // Distributor goroutine
//     go func() {
//         defer func() {
//             for _, ch := range outputs {
//                 close(ch)
//             }
//         }()
//
//         i := 0
//         for v := range input {
//             outputs[i%numWorkers] <- v
//             i++
//         }
//     }()
//
//     return readOnlyOutputs
// }

func FanOut(input <-chan int, numWorkers int) []<-chan int {
	// TODO: Implement FanOut
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// ============================================================================
// Exercise 3: FanIn
// ============================================================================

// FanIn merges multiple input channels into a single output channel.
// TODO: Implement fan-in pattern that combines multiple channels.
//
// Parameters:
//   - inputs: Variable number of input channels
//
// Returns:
//   - <-chan int: Output channel with merged values
//
// Behavior:
//   - Create output channel
//   - For each input channel, start goroutine that forwards to output
//   - When all inputs close, close output
//   - Values arrive in output in whatever order goroutines send them
//
// Architecture:
//   - Create output channel
//   - Create WaitGroup to track input goroutines
//   - For each input:
//     * Start goroutine that reads from input, sends to output
//     * Call wg.Done() when input closes
//   - Start goroutine that waits for all inputs, then closes output
//
// func FanIn(inputs ...<-chan int) <-chan int {
//     output := make(chan int)
//     var wg sync.WaitGroup
//
//     for _, in := range inputs {
//         wg.Add(1)
//         go func(ch <-chan int) {
//             defer wg.Done()
//             for v := range ch {
//                 output <- v
//             }
//         }(in)
//     }
//
//     go func() {
//         wg.Wait()
//         close(output)
//     }()
//
//     return output
// }

func FanIn(inputs ...<-chan int) <-chan int {
	// TODO: Implement FanIn
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// ============================================================================
// Exercise 4: WorkerPool
// ============================================================================

// WorkerPool is defined in types.go:
// type WorkerPool struct {
//     jobs    chan func()
//     wg      sync.WaitGroup
//     stopped atomic.Bool
// }

// NewWorkerPool creates a worker pool with a fixed number of workers.
// TODO: Implement worker pool that processes jobs concurrently.
//
// Parameters:
//   - numWorkers: Number of worker goroutines
//
// Returns:
//   - *WorkerPool: Pool instance with running workers
//
// Behavior:
//   - Create buffered jobs channel (capacity: 100)
//   - Start numWorkers goroutines
//   - Each worker:
//     * Reads jobs from channel
//     * Executes job function
//     * Continues until channel closes
//   - Workers wait on sync.WaitGroup
//
// Memory per worker:
//   - Goroutine stack: ~2KB
//   - 100 workers = ~200KB total
//
// func NewWorkerPool(numWorkers int) *WorkerPool {
//     pool := &WorkerPool{
//         jobs: make(chan func(), 100),
//     }
//
//     for i := 0; i < numWorkers; i++ {
//         pool.wg.Add(1)
//         go func() {
//             defer pool.wg.Done()
//             for job := range pool.jobs {
//                 job()
//             }
//         }()
//     }
//
//     return pool
// }

func NewWorkerPool(numWorkers int) *WorkerPool {
	// TODO: Implement NewWorkerPool
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Submit adds a job to the worker pool.
// TODO: Implement job submission with stopped check.
//
// Behavior:
//   - Check if pool is stopped (atomic load)
//   - If stopped, return immediately (don't block)
//   - Otherwise, send job to jobs channel
//
// func (p *WorkerPool) Submit(job func()) {
//     if p.stopped.Load() {
//         return
//     }
//     p.jobs <- job
// }

func (p *WorkerPool) Submit(job func()) {
	// TODO: Implement job submission
}

// Stop gracefully shuts down the worker pool.
// TODO: Implement graceful shutdown.
//
// Behavior:
//   - Set stopped flag (atomic store)
//   - Close jobs channel (signals workers to exit)
//   - Wait for all workers to finish (wg.Wait)
//
// func (p *WorkerPool) Stop() {
//     p.stopped.Store(true)
//     close(p.jobs)
//     p.wg.Wait()
// }

func (p *WorkerPool) Stop() {
	// TODO: Implement graceful shutdown
}

// ============================================================================
// Exercise 5: RateLimiter
// ============================================================================

// RateLimiter is defined in types.go:
// type RateLimiter struct {
//     ticker *time.Ticker
//     tokens chan struct{}
// }

// NewRateLimiter creates a rate limiter using token bucket algorithm.
// TODO: Implement rate limiter with token bucket.
//
// Parameters:
//   - maxOps: Maximum operations per second
//
// Returns:
//   - *RateLimiter: Rate limiter instance
//
// Token bucket algorithm:
//   - Bucket capacity = maxOps
//   - Initial tokens = maxOps (full bucket)
//   - Refill rate = 1 token per (1 second / maxOps)
//   - Each operation consumes 1 token
//   - If no tokens, operation blocks until refill
//
// func NewRateLimiter(maxOps int) *RateLimiter {
//     limiter := &RateLimiter{
//         ticker: time.NewTicker(time.Second / time.Duration(maxOps)),
//         tokens: make(chan struct{}, maxOps),
//     }
//
//     // Pre-fill tokens
//     for i := 0; i < maxOps; i++ {
//         limiter.tokens <- struct{}{}
//     }
//
//     // Refill tokens at rate
//     go func() {
//         for range limiter.ticker.C {
//             select {
//             case limiter.tokens <- struct{}{}:
//             default:
//                 // Bucket full, skip
//             }
//         }
//     }()
//
//     return limiter
// }

func NewRateLimiter(maxOps int) *RateLimiter {
	// TODO: Implement NewRateLimiter
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Wait blocks until an operation is allowed.
// TODO: Implement blocking wait for token.
//
// Behavior:
//   - Block on tokens channel
//   - When token available, consume it and return
//
// func (r *RateLimiter) Wait() {
//     <-r.tokens
// }

func (r *RateLimiter) Wait() {
	// TODO: Implement wait
}

// ============================================================================
// Exercise 6: ConcurrentCounter
// ============================================================================

// ConcurrentCounter is defined in types.go:
// type ConcurrentCounter struct {
//     value atomic.Int64
// }

// Increment atomically increments the counter by 1.
// TODO: Implement atomic increment.
//
// Why atomic?
//   - Without atomic: counter++ is three operations (read, add, write)
//   - Race: Two goroutines read same value, both write incremented value
//   - Result: Lost update (counter should be 2, but is 1)
//   - Atomic: Single CPU instruction, guaranteed correct
//
// func (c *ConcurrentCounter) Increment() {
//     c.value.Add(1)
// }

func (c *ConcurrentCounter) Increment() {
	// TODO: Implement atomic increment
}

// Decrement atomically decrements the counter by 1.
// TODO: Implement atomic decrement.
//
// func (c *ConcurrentCounter) Decrement() {
//     c.value.Add(-1)
// }

func (c *ConcurrentCounter) Decrement() {
	// TODO: Implement atomic decrement
}

// Value returns the current value of the counter.
// TODO: Implement atomic read.
//
// Why Load() instead of direct access?
//   - Direct access: c.value might see partial write
//   - Load(): Guarantees full value is read atomically
//
// func (c *ConcurrentCounter) Value() int64 {
//     return c.value.Load()
// }

func (c *ConcurrentCounter) Value() int64 {
	// TODO: Implement atomic read
	return 0
}

// ============================================================================
// Exercise 7: GracefulWorker
// ============================================================================

// GracefulWorker is defined in types.go:
// type GracefulWorker struct {
//     ctx      context.Context
//     workDone atomic.Int64
// }

// NewGracefulWorker creates a worker that can be gracefully cancelled.
// TODO: Initialize worker with context.
//
// func NewGracefulWorker(ctx context.Context) *GracefulWorker {
//     return &GracefulWorker{
//         ctx:      ctx,
//         workDone: atomic.Int64{},
//     }
// }

func NewGracefulWorker(ctx context.Context) *GracefulWorker {
	// TODO: Implement NewGracefulWorker
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Start begins the worker's execution loop.
// TODO: Implement worker loop with context cancellation.
//
// Behavior:
//   - Run in a goroutine
//   - Loop continuously, incrementing workDone
//   - Check ctx.Done() on each iteration
//   - Exit when context is cancelled
//
// Pattern:
//   go func() {
//       for {
//           select {
//           case <-w.ctx.Done():
//               return // Exit on cancellation
//           default:
//               // Do work
//               w.workDone.Add(1)
//               time.Sleep(1 * time.Microsecond)
//           }
//       }
//   }()
//
// func (w *GracefulWorker) Start() {
//     go func() {
//         for {
//             select {
//             case <-w.ctx.Done():
//                 return
//             default:
//                 w.workDone.Add(1)
//                 time.Sleep(1 * time.Microsecond)
//             }
//         }
//     }()
// }

func (w *GracefulWorker) Start() {
	// TODO: Implement worker loop with context
}

// WorkDone returns the total work completed.
func (w *GracefulWorker) WorkDone() int64 {
	return w.workDone.Load()
}

// ============================================================================
// Exercise 8: Pipeline
// ============================================================================

// Pipeline chains multiple processing stages together.
// TODO: Implement pipeline pattern.
//
// Parameters:
//   - input: Input channel
//   - stages: Functions that transform channels
//
// Returns:
//   - <-chan int: Output of final stage
//
// Behavior:
//   - Start with input channel
//   - For each stage, apply it to current channel
//   - Each stage creates new channel, starts goroutine
//   - Return final output channel
//
// Example:
//   double := func(in <-chan int) <-chan int {
//       out := make(chan int)
//       go func() {
//           defer close(out)
//           for v := range in {
//               out <- v * 2
//           }
//       }()
//       return out
//   }
//
//   input := make(chan int)
//   output := Pipeline(input, double, double) // Quadruples values
//
// func Pipeline(input <-chan int, stages ...func(<-chan int) <-chan int) <-chan int {
//     output := input
//     for _, stage := range stages {
//         output = stage(output)
//     }
//     return output
// }

func Pipeline(input <-chan int, stages ...func(<-chan int) <-chan int) <-chan int {
	// TODO: Implement Pipeline
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// ============================================================================
// Exercise 9: BoundedParallel
// ============================================================================

// BoundedParallel executes functions concurrently with a maximum concurrency limit.
// TODO: Implement bounded parallelism using semaphore pattern.
//
// Parameters:
//   - maxConcurrent: Maximum number of functions running simultaneously
//   - fns: Functions to execute
//
// Behavior:
//   - At most maxConcurrent functions run at the same time
//   - Use buffered channel as semaphore (token bucket)
//   - Each function acquires token before running
//   - Releases token when done
//   - Wait for all functions to complete
//
// Semaphore pattern:
//   sem := make(chan struct{}, maxConcurrent)
//   sem <- struct{}{} // Acquire
//   <-sem             // Release
//
// Why bounded?
//   - Prevent resource exhaustion (too many goroutines)
//   - Control load on external systems
//   - Example: Don't want 10,000 concurrent HTTP requests
//
// func BoundedParallel(maxConcurrent int, fns ...func()) {
//     sem := make(chan struct{}, maxConcurrent)
//     var wg sync.WaitGroup
//
//     for _, fn := range fns {
//         wg.Add(1)
//
//         go func(f func()) {
//             defer wg.Done()
//
//             // Acquire semaphore
//             sem <- struct{}{}
//             defer func() { <-sem }() // Release
//
//             // Execute function
//             f()
//         }(fn)
//     }
//
//     wg.Wait()
// }

func BoundedParallel(maxConcurrent int, fns ...func()) {
	// TODO: Implement BoundedParallel
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// After implementing all functions:
// - Run: go test -v ./...
// - Test with race detector: go test -race
// - Experiment with different numbers of workers
// - Monitor goroutine count with runtime.NumGoroutine()
// - Compare performance: 1 worker vs 4 workers vs 100 workers
//
// Key takeaways:
// 1. Goroutines are cheap, but not free
// 2. Use worker pools for bounded concurrency
// 3. Always ensure goroutines can exit (no leaks)
// 4. Atomic operations for simple shared state
// 5. Channels for complex coordination
// 6. Context for cancellation
// 7. WaitGroup to wait for completion
