//go:build reference

package goroutines1mdemo

/*
Reference Solution - Goroutines and Concurrent Programming
=========================================================

This file demonstrates Go's concurrency model with goroutines and channels.
Goroutines are lightweight threads managed by the Go runtime, enabling massive
concurrency (millions of goroutines) without the overhead of OS threads.

This connects to the broader Go ecosystem by showing:
- How goroutines enable scalable server architectures (net/http, database drivers)
- Why channels provide safer concurrency than shared memory
- How sync/atomic enables lock-free programming for performance
- Why context enables cooperative cancellation across goroutine trees

The exercise builds understanding of:
- Goroutine lifecycle: creation, scheduling, and termination
- Race conditions: why shared mutable state needs synchronization
- Work distribution: dividing problems across concurrent workers
- Synchronization primitives: WaitGroups, atomics, channels
- Performance characteristics: when concurrency helps vs overhead

Teaching notes:
- Memory/ownership: goroutines can share memory through pointers/channels, but
  this creates race conditions unless properly synchronized. The runtime doesn't
  prevent data races - programmers must use synchronization primitives.
- Invariants: concurrent programs must establish happens-before relationships
  using channels, mutexes, or atomics. Without these, execution is non-deterministic.
- Error surfaces: goroutine panics don't crash the program but can be caught
  with recover(). Network timeouts and context cancellation provide graceful
  failure modes for concurrent operations.
*/

/*
Project 18: Goroutines and Concurrency - Complete Solutions

This file demonstrates advanced concurrency patterns that power production Go systems.
The "1M demo" shows how Go can efficiently manage massive concurrency with minimal resources.

Key concurrency concepts demonstrated:
1. **Goroutine lifecycle**: Creation with `go` keyword, automatic scheduling by runtime
2. **Work distribution**: Dividing computational work across parallel workers
3. **Atomic operations**: Lock-free thread-safe counters for performance
4. **Channel patterns**: Fan-in/fan-out for dataflow orchestration
5. **Graceful shutdown**: Context-based cancellation for clean termination

DEBUGGING CONCURRENCY:
- Set breakpoints at goroutine creation to observe parallel execution
- Use runtime.NumGoroutine() to track active goroutines
- Watch atomic variables change from multiple goroutines simultaneously
- Step through channel operations to understand blocking/synchronization
- Observe race conditions with `go test -race` before fixes
*/

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

/*
ParallelSum - Concurrent Work Distribution and Atomic Accumulation

This function demonstrates the complete lifecycle of concurrent computation:
work division, parallel execution, and thread-safe result aggregation.

The algorithm: sum numbers from 1 to n using multiple worker goroutines.
Each worker sums a contiguous range, then atomically adds to a shared total.

Why this matters:
- Shows how CPU-bound work can be parallelized for performance
- Demonstrates goroutine overhead vs parallelism benefits
- Illustrates atomic operations for lock-free shared state
- Foundation for map-reduce style distributed computing

Parameters:
- n: upper bound of summation (1 + 2 + ... + n)
- numWorkers: number of concurrent goroutines to use

Returns: the mathematical sum as int64 (to handle large values)
*/
func ParallelSum(n int, numWorkers int) int64 {
	// SYNCHRONIZATION SETUP
	// Atomic counter for thread-safe accumulation across goroutines
	// Unlike mutex-protected variables, atomics have no lock contention
	// Perfect for high-frequency updates from many goroutines
	var total atomic.Int64

	// WaitGroup coordinates goroutine lifecycle
	// Ensures main goroutine waits for all workers to complete
	// Without this, function would return before workers finish
	var wg sync.WaitGroup

	// WORK DISTRIBUTION CALCULATION
	// Divide the range [1..n] into approximately equal chunks
	// Integer division ensures fair distribution with possible remainder
	rangeSize := n / numWorkers      // Base size per worker
	remainder := n % numWorkers      // Extra work for last worker

	// WORKER LAUNCH LOOP
	// Launch numWorkers goroutines, each handling one range
	// The `go` keyword creates a new goroutine (lightweight thread)
	// Goroutines are scheduled by Go runtime, not OS scheduler
	for i := 0; i < numWorkers; i++ {
		// Track this worker in WaitGroup before launching
		// wg.Add(1) must happen before goroutine starts
		wg.Add(1)

		// Calculate this worker's range boundaries
		// Ranges are contiguous and non-overlapping for correctness
		start := i*rangeSize + 1        // Start of this worker's range
		end := (i + 1) * rangeSize      // End of this worker's range

		// Last worker gets remainder to ensure all numbers are summed
		// Without this, numbers in remainder would be missed
		if i == numWorkers-1 {
			end += remainder
		}

		// LAUNCH GOROUTINE (per .cursorrules: closure capture)
		// go func(s, e int) — we PASS start/end as parameters, not capture i.
		// If we wrote go func() { ... use start, end } and used loop vars
		// start, end directly, all goroutines could share the SAME variables
		// (race!). By passing (start, end), each goroutine gets its own COPY.
		// BEFORE: loop has start, end. AFTER: goroutine runs with private s, e.
		go func(s, e int) {
			// Ensure WaitGroup is notified when worker completes
			// defer executes when function returns (normal or panic)
			defer wg.Done()

			// COMPUTE PARTIAL SUM
			// Each worker sums its assigned range independently
			// No shared state here - each worker has private variables
			var sum int64
			for j := s; j <= e; j++ {
				sum += int64(j)
			}

			// ATOMIC ACCUMULATION
			// Thread-safe addition to shared total
			// atomic.Add() is lock-free and wait-free
			// Multiple workers can call this simultaneously without issues
			total.Add(sum)
		}(start, end)
	}

	// SYNCHRONIZATION BARRIER
	// Block until all worker goroutines call wg.Done()
	// This ensures the sum is complete before returning
	// Without wg.Wait(), function would return 0 (initial atomic value)
	wg.Wait()

	// RETURN FINAL RESULT
	// Atomic load ensures we see all previous adds
	// Returns the complete sum of range [1..n]
	return total.Load()
}

// FanOut distributes values to multiple channels.
func FanOut(input <-chan int, numWorkers int) []<-chan int {
	outputs := make([]chan int, numWorkers)
	readOnlyOutputs := make([]<-chan int, numWorkers)

	for i := 0; i < numWorkers; i++ {
		outputs[i] = make(chan int)
		readOnlyOutputs[i] = outputs[i]
	}

	// Distributor goroutine
	go func() {
		defer func() {
			for _, ch := range outputs {
				close(ch)
			}
		}()

		i := 0
		for v := range input {
			// Round-robin distribution
			outputs[i%numWorkers] <- v
			i++
		}
	}()

	return readOnlyOutputs
}

// FanIn merges multiple channels into one.
func FanIn(inputs ...<-chan int) <-chan int {
	output := make(chan int)
	var wg sync.WaitGroup

	// Launch a goroutine for each input
	for _, in := range inputs {
		wg.Add(1)
		go func(ch <-chan int) {
			defer wg.Done()
			for v := range ch {
				output <- v
			}
		}(in)
	}

	// Close output when all inputs are done
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

// NewWorkerPool creates a worker pool.
func NewWorkerPool(numWorkers int) *WorkerPool {
	pool := &WorkerPool{
		jobs: make(chan func(), 100), // Buffered channel
	}

	// Launch workers
	for i := 0; i < numWorkers; i++ {
		pool.wg.Add(1)
		go func() {
			defer pool.wg.Done()
			for job := range pool.jobs {
				job()
			}
		}()
	}

	return pool
}

// Submit adds a job to the pool.
func (p *WorkerPool) Submit(job func()) {
	if p.stopped.Load() {
		return // Pool is stopped, don't accept new jobs
	}
	p.jobs <- job
}

// Stop shuts down the pool.
func (p *WorkerPool) Stop() {
	p.stopped.Store(true)
	close(p.jobs) // Signal workers to stop
	p.wg.Wait()   // Wait for all workers to finish
}

// NewRateLimiter creates a rate limiter.
func NewRateLimiter(maxOps int) *RateLimiter {
	limiter := &RateLimiter{
		ticker: time.NewTicker(time.Second / time.Duration(maxOps)),
		tokens: make(chan struct{}, maxOps),
	}

	// Pre-fill tokens
	for i := 0; i < maxOps; i++ {
		limiter.tokens <- struct{}{}
	}

	// Refill tokens at rate
	go func() {
		for range limiter.ticker.C {
			select {
			case limiter.tokens <- struct{}{}:
			default:
				// Token bucket full, skip
			}
		}
	}()

	return limiter
}

// Wait blocks until a token is available.
func (r *RateLimiter) Wait() {
	<-r.tokens
}

// ConcurrentCounter implementation.

// Increment atomically increments.
func (c *ConcurrentCounter) Increment() {
	c.value.Add(1)
}

// Decrement atomically decrements.
func (c *ConcurrentCounter) Decrement() {
	c.value.Add(-1)
}

// Value returns current value.
func (c *ConcurrentCounter) Value() int64 {
	return c.value.Load()
}

// NewGracefulWorker creates a new graceful worker.
func NewGracefulWorker(ctx context.Context) *GracefulWorker {
	return &GracefulWorker{
		ctx:      ctx,
		workDone: atomic.Int64{},
	}
}

// Start begins execution.
func (w *GracefulWorker) Start() {
	go func() {
		for {
			select {
			case <-w.ctx.Done():
				return // Exit when cancelled
			default:
				// Do work
				w.workDone.Add(1)
				time.Sleep(1 * time.Microsecond)
			}
		}
	}()
}

// WorkDone returns the total work completed.
func (w *GracefulWorker) WorkDone() int64 {
	return w.workDone.Load()
}

// Pipeline chains stages together.
func Pipeline(input <-chan int, stages ...func(<-chan int) <-chan int) <-chan int {
	output := input
	for _, stage := range stages {
		output = stage(output)
	}
	return output
}

// BoundedParallel executes with limited concurrency.
func BoundedParallel(maxConcurrent int, fns ...func()) {
	sem := make(chan struct{}, maxConcurrent) // Semaphore
	var wg sync.WaitGroup

	for _, fn := range fns {
		wg.Add(1)

		go func(f func()) {
			defer wg.Done()

			// Acquire semaphore
			sem <- struct{}{}
			defer func() { <-sem }() // Release semaphore

			// Execute function
			f()
		}(fn)
	}

	wg.Wait()
}
