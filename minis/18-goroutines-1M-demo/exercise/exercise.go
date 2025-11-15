//go:build !solution
// +build !solution

// Package exercise contains hands-on exercises for goroutines and concurrency.

package exercise

import (
	"context"
)

// ParallelSum calculates the sum of numbers from 1 to n using multiple goroutines.
//
// REQUIREMENTS:
// - Split the work among numWorkers goroutines
// - Each worker sums a portion of the range
// - Return the total sum
//
// EXAMPLE:
//   ParallelSum(100, 4)  // 4 workers each sum 25 numbers
//   Returns: 5050 (sum of 1..100)
//
// HINT: Use a WaitGroup to wait for all workers to finish
// HINT: Use atomic operations or a mutex to safely accumulate the sum
func ParallelSum(n int, numWorkers int) int64 {
	// TODO: Implement this
	return 0
}

// FanOut sends values from a single input channel to multiple worker channels.
//
// REQUIREMENTS:
// - Create numWorkers output channels
// - Launch goroutines that read from input and distribute to outputs (round-robin or random)
// - Close output channels when input is closed
// - Return a slice of the output channels
//
// EXAMPLE:
//   input := make(chan int)
//   outputs := FanOut(input, 3)  // Creates 3 output channels
//   // Values sent to input are distributed to outputs[0], outputs[1], outputs[2]
//
// HINT: Use a goroutine to read from input and distribute values
// HINT: Don't forget to close all output channels when input is closed
func FanOut(input <-chan int, numWorkers int) []<-chan int {
	// TODO: Implement this
	return nil
}

// FanIn merges multiple input channels into a single output channel.
//
// REQUIREMENTS:
// - Create a single output channel
// - Launch a goroutine for each input channel that forwards values to output
// - Close output when ALL inputs are closed
// - Return the output channel
//
// EXAMPLE:
//   ch1 := make(chan int)
//   ch2 := make(chan int)
//   output := FanIn(ch1, ch2)
//   // Values from ch1 and ch2 are merged into output
//
// HINT: Use a WaitGroup to track when all input goroutines finish
// HINT: Close output in a separate goroutine after WaitGroup is done
func FanIn(inputs ...<-chan int) <-chan int {
	// TODO: Implement this
	return nil
}

// NewWorkerPool creates a new worker pool with the given number of workers.
//
// REQUIREMENTS:
// - Create a pool with numWorkers goroutines
// - Workers receive jobs from a channel and execute them
// - The returned pool should support Submit(job) and Stop()
//
// EXAMPLE:
//   pool := NewWorkerPool(5)
//   pool.Submit(func() { fmt.Println("job 1") })
//   pool.Submit(func() { fmt.Println("job 2") })
//   pool.Stop()  // Wait for all jobs to finish, then stop workers
func NewWorkerPool(numWorkers int) *WorkerPool {
	// TODO: Implement this
	return nil
}

// Submit adds a job to the worker pool.
//
// REQUIREMENTS:
// - Send the job to the job channel
// - If the pool is stopped, this should not block forever
//
// HINT: Check if the pool is stopped before sending
func (p *WorkerPool) Submit(job func()) {
	// TODO: Implement this
}

// Stop gracefully shuts down the worker pool.
//
// REQUIREMENTS:
// - Wait for all currently running jobs to finish
// - Prevent new jobs from being submitted
// - Clean up resources (close channels, etc.)
func (p *WorkerPool) Stop() {
	// TODO: Implement this
}

// NewRateLimiter creates a rate limiter that allows maxOps operations per second.
//
// REQUIREMENTS:
// - Allow at most maxOps operations per second
// - Return a RateLimiter that supports Wait()
// - Wait() should block until an operation slot is available
// - Use a ticker to release slots at a steady rate
//
// EXAMPLE:
//   limiter := NewRateLimiter(10)  // 10 ops/sec
//   for i := 0; i < 100; i++ {
//       limiter.Wait()  // Blocks until allowed
//       doExpensiveOperation()
//   }
func NewRateLimiter(maxOps int) *RateLimiter {
	// TODO: Implement this
	return nil
}

// Wait blocks until an operation is allowed.
func (r *RateLimiter) Wait() {
	// TODO: Implement this
}

// Increment atomically increments the counter by 1.
//
// REQUIREMENTS:
// - Use atomic operations to ensure thread-safety
// - Increment the counter's value by 1
//
// EXAMPLE:
//   counter := &ConcurrentCounter{}
//   counter.Increment()
//   counter.Increment()
//   fmt.Println(counter.Value())  // 2
func (c *ConcurrentCounter) Increment() {
	// TODO: Implement this
}

// Decrement atomically decrements the counter by 1.
func (c *ConcurrentCounter) Decrement() {
	// TODO: Implement this
}

// Value returns the current value of the counter.
func (c *ConcurrentCounter) Value() int64 {
	// TODO: Implement this
	return 0
}

// NewGracefulWorker creates a new graceful worker.
//
// REQUIREMENTS:
// - Create a worker that runs until context is cancelled
// - The worker should support Start() to begin work
// - Track the total work done (accessible via WorkDone())
//
// EXAMPLE:
//   ctx, cancel := context.WithCancel(context.Background())
//   worker := NewGracefulWorker(ctx)
//   worker.Start()
//   time.Sleep(100 * time.Millisecond)
//   cancel()  // Stop the worker
//   fmt.Println(worker.WorkDone())
func NewGracefulWorker(ctx context.Context) *GracefulWorker {
	return &GracefulWorker{ctx: ctx}
}

// Start begins the worker's execution loop.
//
// REQUIREMENTS:
// - Run in a goroutine
// - Increment workDone on each iteration
// - Check ctx.Done() and exit when cancelled
//
// HINT: Use select to check for cancellation
func (w *GracefulWorker) Start() {
	// TODO: Implement this
}

// WorkDone returns the total work completed.
func (w *GracefulWorker) WorkDone() int64 {
	return w.workDone.Load()
}

// Pipeline creates a processing pipeline with multiple stages.
//
// REQUIREMENTS:
// - Each stage is a function that transforms input to output
// - Values flow through stages: input -> stage1 -> stage2 -> ... -> output
// - Each stage runs in its own goroutine
// - Close channels appropriately when upstream closes
//
// EXAMPLE:
//   double := func(in <-chan int) <-chan int {
//       out := make(chan int)
//       go func() {
//           for v := range in {
//               out <- v * 2
//           }
//           close(out)
//       }()
//       return out
//   }
//   input := make(chan int)
//   output := Pipeline(input, double, double)  // Quadruples values
//
// HINT: Chain the stages by connecting each output to the next input
func Pipeline(input <-chan int, stages ...func(<-chan int) <-chan int) <-chan int {
	// TODO: Implement this
	return nil
}

// BoundedParallel executes functions concurrently with a maximum concurrency limit.
//
// REQUIREMENTS:
// - Execute all functions in fns
// - At most maxConcurrent functions run simultaneously
// - Wait for all functions to complete before returning
//
// EXAMPLE:
//   fns := []func(){
//       func() { time.Sleep(100 * time.Millisecond) },
//       func() { time.Sleep(100 * time.Millisecond) },
//       // ... 100 more functions
//   }
//   BoundedParallel(10, fns...)  // Only 10 run at a time
//
// HINT: Use a semaphore (buffered channel) to limit concurrency
func BoundedParallel(maxConcurrent int, fns ...func()) {
	// TODO: Implement this
}
