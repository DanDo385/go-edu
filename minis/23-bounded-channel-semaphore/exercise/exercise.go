//go:build !solution
// +build !solution

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

package exercise

import (
	"context"
	"fmt"
	"time"
)

// ============================================================================
// EXERCISE 1: Basic Counting Semaphore
// ============================================================================

// Semaphore is a counting semaphore using buffered channels.
//
// TODO: Implement a semaphore that:
// - Limits concurrent access to maxPermits
// - Blocks on Acquire when full
// - Releases permits on Release
// - Uses buffered channel as underlying mechanism
//
// HINT: The channel capacity should equal maxPermits.
// HINT: Send = acquire (blocks when full), receive = release.
type Semaphore struct {
	// TODO: Add fields
	// You'll need a buffered channel
}

// NewSemaphore creates a new counting semaphore.
//
// TODO: Initialize the semaphore with the given capacity.
func NewSemaphore(maxPermits int) *Semaphore {
	// TODO: Implement
	return nil
}

// Acquire acquires a permit, blocking if none available.
//
// TODO: Block until a permit is available.
// HINT: Send to the buffered channel.
func (s *Semaphore) Acquire() {
	// TODO: Implement
}

// Release releases a permit back to the semaphore.
//
// TODO: Release a permit.
// HINT: Receive from the buffered channel.
func (s *Semaphore) Release() {
	// TODO: Implement
}

// ============================================================================
// EXERCISE 2: Try-Acquire (Non-Blocking)
// ============================================================================

// TryAcquire attempts to acquire without blocking.
//
// TODO: Try to acquire a permit without blocking.
// Returns true if acquired, false if would block.
//
// HINT: Use select with default case.
func (s *Semaphore) TryAcquire() bool {
	// TODO: Implement
	return false
}

// ============================================================================
// EXERCISE 3: Context-Aware Acquisition
// ============================================================================

// AcquireWithContext acquires with timeout/cancellation support.
//
// TODO: Acquire a permit, respecting context timeout/cancellation.
// Returns error if context is cancelled before acquiring.
//
// HINT: Use select to wait on both channel and context.Done().
func (s *Semaphore) AcquireWithContext(ctx context.Context) error {
	// TODO: Implement
	return nil
}

// ============================================================================
// EXERCISE 4: Rate Limiter
// ============================================================================

// RateLimiter limits operations to a maximum rate.
//
// TODO: Implement a token bucket rate limiter using channels.
// - Allow burst of maxBurst requests
// - Refill tokens at specified rate
// - Block when no tokens available
//
// HINT: Use buffered channel for tokens.
// HINT: Use goroutine with ticker to refill tokens.
type RateLimiter struct {
	// TODO: Add fields
	// You'll need: token channel, ticker, done channel
}

// NewRateLimiter creates a new rate limiter.
//
// TODO: Initialize rate limiter and start refill goroutine.
func NewRateLimiter(maxBurst int, rate time.Duration) *RateLimiter {
	// TODO: Implement
	// 1. Create token channel with capacity maxBurst
	// 2. Fill initial tokens
	// 3. Start goroutine to refill tokens periodically
	return nil
}

// Wait blocks until a token is available.
//
// TODO: Wait for and consume a token.
func (rl *RateLimiter) Wait() {
	// TODO: Implement
}

// TryAcquire attempts non-blocking token acquisition.
//
// TODO: Try to acquire token without blocking.
func (rl *RateLimiter) TryAcquire() bool {
	// TODO: Implement
	return false
}

// Stop stops the rate limiter.
//
// TODO: Clean up resources (stop refill goroutine).
func (rl *RateLimiter) Stop() {
	// TODO: Implement
}

// ============================================================================
// EXERCISE 5: Weighted Semaphore
// ============================================================================

// WeightedSemaphore allows acquiring multiple permits at once.
//
// TODO: Implement a weighted semaphore that:
// - Tracks total capacity
// - Allows acquiring N permits at once
// - Blocks if insufficient permits available
// - Properly handles context cancellation during partial acquisition
type WeightedSemaphore struct {
	// TODO: Add fields
}

// NewWeightedSemaphore creates a weighted semaphore.
//
// TODO: Initialize with max total weight.
func NewWeightedSemaphore(maxWeight int) *WeightedSemaphore {
	// TODO: Implement
	return nil
}

// Acquire acquires the specified weight of permits.
//
// TODO: Acquire 'weight' permits, blocking if insufficient.
// HINT: You'll need to acquire permits one at a time in a loop.
func (ws *WeightedSemaphore) Acquire(weight int) {
	// TODO: Implement
}

// Release releases the specified weight of permits.
//
// TODO: Release 'weight' permits.
func (ws *WeightedSemaphore) Release(weight int) {
	// TODO: Implement
}

// AcquireWithContext acquires with context support.
//
// TODO: Acquire 'weight' permits, respecting context.
// CRITICAL: If context cancels during acquisition, release what you got!
//
// HINT: Track how many permits acquired so far.
// HINT: On context cancel, release acquired permits before returning error.
func (ws *WeightedSemaphore) AcquireWithContext(ctx context.Context, weight int) error {
	// TODO: Implement
	return nil
}

// ============================================================================
// EXERCISE 6: Worker Pool
// ============================================================================

// WorkerPool processes jobs with bounded concurrency.
//
// TODO: Implement a worker pool that:
// - Accepts jobs via Submit()
// - Processes jobs concurrently (max numWorkers at a time)
// - Returns results via Results() channel
// - Supports graceful shutdown via Stop()
//
// HINT: Use semaphore to limit concurrent workers.
// HINT: Use channels for job queue and results.
type WorkerPool struct {
	// TODO: Add fields
	// You'll need: job channel, result channel, semaphore, etc.
}

// NewWorkerPool creates a worker pool.
//
// TODO: Initialize worker pool.
// processor is the function that processes each job.
func NewWorkerPool(numWorkers int, processor func(Job) Result) *WorkerPool {
	// TODO: Implement
	return nil
}

// Submit submits a job to the pool.
//
// TODO: Send job to job queue.
func (wp *WorkerPool) Submit(job Job) {
	// TODO: Implement
}

// Start starts processing jobs.
//
// TODO: Start goroutine that:
// 1. Reads jobs from queue
// 2. Acquires semaphore permit
// 3. Launches goroutine to process job
// 4. Sends result to results channel
// 5. Releases permit when done
func (wp *WorkerPool) Start() {
	// TODO: Implement
}

// Results returns the results channel.
//
// TODO: Return channel where results are sent.
func (wp *WorkerPool) Results() <-chan Result {
	// TODO: Implement
	return nil
}

// Stop gracefully stops the pool.
//
// TODO: Close job queue and wait for all workers to finish.
func (wp *WorkerPool) Stop() {
	// TODO: Implement
}

// ============================================================================
// HELPER: Process Function for Testing
// ============================================================================

// DefaultProcessor is a simple job processor for testing.
func DefaultProcessor(job Job) Result {
	return Result{
		JobID:  job.ID,
		Output: fmt.Sprintf("Processed: %s", job.Data),
		Err:    nil,
	}
}
