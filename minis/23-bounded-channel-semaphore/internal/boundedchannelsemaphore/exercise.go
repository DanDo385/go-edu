//go:build !solution && !reference


// Package exercise provides semaphore implementation exercises.
//
// EXERCISES:
// 1. Implement basic counting semaphore
// 2. Add context-aware acquisition with timeout
// 3. Implement try-acquire (non-blocking)
// 4. Build rate limiter using semaphore
// 5. Create weighted semaphore for variable costs
// 6. Implement worker pool with semaphore-based concurrency control
//
// LEARNING GOALS:
// - Master buffered channels as semaphores
// - Understand acquire/release patterns
// - Handle timeouts and cancellation
// - Implement common concurrency patterns

package boundedchannelsemaphore

// import (
// 	"context"
// 	"fmt"
// 	"time"
// )

// ============================================================================
// EXERCISE 1: Basic Counting Semaphore
// ============================================================================

// Semaphore is a counting semaphore using buffered channels.
type Semaphore struct {
	sem chan struct{}
}

// NewSemaphore creates a new counting semaphore.
func NewSemaphore(maxPermits int) *Semaphore {
	// TODO: Implement NewSemaphore
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Acquire acquires a permit, blocking if none available.
func (s *Semaphore) Acquire() {
	// TODO: Implement this function.
	// - To "acquire" a permit, send a value to the channel.
	// - `s.sem <- struct{}{} `
	// - If the channel's buffer is full (meaning all permits are in use), this operation will block until another goroutine calls `Release`.
}

// Release releases a permit back to the semaphore.
func (s *Semaphore) Release() {
	// TODO: Implement this function.
	// - To "release" a permit, receive a value from the channel.
	// - `<-s.sem`
	// - This makes space in the buffer, allowing a blocked `Acquire` call to proceed.
}

// TryAcquire attempts to acquire without blocking.
func (s *Semaphore) TryAcquire() bool {
	// TODO: Implement this function.
	// - Use a `select` statement with a `default` case for a non-blocking send.
	// - `case s.sem <- struct{}{}:` -> Acquired successfully, return `true`.
	// - `default:` -> Would have blocked, return `false`.
	return false
}

// AcquireWithContext acquires with timeout/cancellation support.
func (s *Semaphore) AcquireWithContext(ctx context.Context) error {
	// TODO: Implement this function.
	// - Use a `select` statement to wait on two different channels.
	// - `case s.sem <- struct{}{}:` -> Acquired successfully, return `nil`.
	// - `case <-ctx.Done():` -> The context was cancelled or timed out. Return `ctx.Err()`.
	return nil
}

// ============================================================================
// EXERCISE 4: Rate Limiter
// ============================================================================

// RateLimiter limits operations to a maximum rate.
type RateLimiter struct {
	tokens chan struct{}
	rate   time.Duration
	done   chan struct{}
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(maxBurst int, rate time.Duration) *RateLimiter {
	// TODO: Implement NewRateLimiter
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Wait blocks until a token is available.
func (rl *RateLimiter) Wait() {
	// TODO: Implement this function.
	// - To wait for a token, simply receive from the `tokens` channel.
	// - This will block if the bucket is empty.
}

// TryAcquire attempts non-blocking token acquisition.
func (rl *RateLimiter) TryAcquire() bool {
	// TODO: Implement this function.
	// - Use a `select` with a `default` case to perform a non-blocking receive from the `tokens` channel.
	// - Return `true` if a token was received, `false` otherwise.
	return false
}

// Stop stops the rate limiter.
func (rl *RateLimiter) Stop() {
	// TODO: Implement this function.
	// - Close the `done` channel to signal the refill goroutine to stop.
}

// ============================================================================
// EXERCISE 5: Weighted Semaphore
// ============================================================================

// WeightedSemaphore allows acquiring multiple permits at once.
type WeightedSemaphore struct {
	permits chan struct{}
}

// NewWeightedSemaphore creates a weighted semaphore.
func NewWeightedSemaphore(maxWeight int) *WeightedSemaphore {
	// TODO: Implement NewWeightedSemaphore
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Acquire acquires the specified weight of permits.
func (ws *WeightedSemaphore) Acquire(weight int) {
	// TODO: Implement this function.
	// - To acquire a "weight", you need to acquire that many permits.
	// - Use a `for` loop to send `weight` times to the `permits` channel.
}

// Release releases the specified weight of permits.
func (ws *WeightedSemaphore) Release(weight int) {
	// TODO: Implement this function.
	// - To release a "weight", you need to release that many permits.
	// - Use a `for` loop to receive `weight` times from the `permits` channel.
}

// AcquireWithContext acquires with context support.
func (ws *WeightedSemaphore) AcquireWithContext(ctx context.Context, weight int) error {
	// TODO: Implement this function.

	// This is the most complex part of the weighted semaphore.

	// Step 1: Keep track of how many permits you have acquired.
	// - `var acquired int`

	// Step 2: Loop `weight` times to acquire each permit.
	// - Inside the loop, use a `select` statement.
	// - `case ws.permits <- struct{}{}:`
	//   - You successfully acquired one permit. Increment your `acquired` counter.
	// - `case <-ctx.Done():`
	//   - The context was cancelled *during* your acquisition attempt.
	//   - **CRITICAL:** You must release the permits you already acquired to avoid leaking them.
	//   - Loop `acquired` times and receive from the `permits` channel to give them back.
	//   - Return `ctx.Err()`.

	// Step 3: If the loop completes without the context being cancelled, return `nil`.
	return nil
}

// ============================================================================
// EXERCISE 6: Worker Pool
// ============================================================================

// WorkerPool processes jobs with bounded concurrency.
type WorkerPool struct {
	jobs      chan Job
	results   chan Result
	sem       *Semaphore
	processor func(Job) Result
	wg        sync.WaitGroup
}

// NewWorkerPool creates a worker pool.
func NewWorkerPool(numWorkers int, processor func(Job) Result) *WorkerPool {
	// TODO: Implement NewWorkerPool
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Submit submits a job to the pool.
func (wp *WorkerPool) Submit(job Job) {
	// TODO: Implement this function.
	// - Send the job to the `jobs` channel.
}

// Start starts processing jobs.
func (wp *WorkerPool) Start() {
	// TODO: Implement this function.

	// This method contains the main loop that orchestrates the workers.

	// Step 1: Launch a single goroutine to manage job distribution.
	// - `wp.wg.Add(1)`
	// - `go func() { ... }()`
	// - The goroutine should `defer wp.wg.Done()`.

	// Step 2: Implement the distribution logic.
	// - Use a `for job := range wp.jobs` loop to read submitted jobs.
	// - Inside the loop, acquire a semaphore permit: `wp.sem.Acquire()`. This will block if `numWorkers` jobs are already running.
	// - Once a permit is acquired, launch a new goroutine for the worker.
	//   - `wp.wg.Add(1)`
	//   - The worker goroutine should:
	//     - `defer wp.wg.Done()`
	//     - `defer wp.sem.Release()` to release the permit when it's finished.
	//     - Call the `wp.processor(job)` function.
	//     - Send the result to the `wp.results` channel.
}

// Results returns the results channel.
func (wp *WorkerPool) Results() <-chan Result {
	// TODO: Implement this function.
	// - Return the `results` channel.
	return nil
}

// Stop gracefully stops the pool.
func (wp *WorkerPool) Stop() {
	// TODO: Implement this function.
	// - Close the `jobs` channel. This will cause the `range` loop in the `Start` method to terminate once all jobs are read.
	// - Wait for all goroutines (the distributor and all workers) to finish: `wp.wg.Wait()`.
	// - After waiting, close the `results` channel to signal to consumers that no more results will be sent.
}

// ============================================================================
// HELPER: Process Function for Testing
// ============================================================================

// DefaultProcessor is a simple job processor for testing.
func DefaultProcessor(job Job) Result {
	// TODO: Implement DefaultProcessor
	// See solution.reference.go for reference implementation
	panic("not implemented")
}

