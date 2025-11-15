//go:build !solution
// +build !solution

package exercise

import (
	"context"
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
	// TODO: Add fields
	// Hint: You'll need channels for jobs and results,
	// worker count, and synchronization primitives
}

// NewWorkerPool creates a new worker pool
//
// Parameters:
//   - queueSize: Maximum number of jobs that can be queued
//   - numWorkers: Number of concurrent workers
//
// Returns a configured WorkerPool ready to start
func NewWorkerPool(queueSize, numWorkers int) *WorkerPool {
	// TODO: implement
	return nil
}

// Start begins processing jobs
//
// Parameters:
//   - ctx: Context for cancellation
//   - process: Function to process each job
//
// Behavior:
//   - Starts numWorkers goroutines
//   - Each worker processes jobs from the queue
//   - Workers stop when context is cancelled or jobs channel is closed
//   - Results channel is closed when all workers finish
func (p *WorkerPool) Start(ctx context.Context, process func(Job) Result) {
	// TODO: implement
}

// Submit attempts to add a job to the queue (non-blocking)
//
// Parameters:
//   - job: The job to submit
//
// Returns:
//   - error: ErrQueueFull if queue is at capacity, nil otherwise
//
// Behavior:
//   - Uses select with default to avoid blocking
//   - Returns immediately if queue is full (backpressure)
func (p *WorkerPool) Submit(job Job) error {
	// TODO: implement
	return nil
}

// SubmitWithTimeout attempts to add a job with a timeout
//
// Parameters:
//   - ctx: Context for cancellation
//   - job: The job to submit
//   - timeout: Maximum time to wait
//
// Returns:
//   - error: ErrQueueFull, context error, or nil
//
// Behavior:
//   - Waits up to timeout duration for queue space
//   - Returns error if timeout expires or context is cancelled
func (p *WorkerPool) SubmitWithTimeout(ctx context.Context, job Job, timeout time.Duration) error {
	// TODO: implement
	return nil
}

// Results returns a read-only channel of results
func (p *WorkerPool) Results() <-chan Result {
	// TODO: implement
	return nil
}

// Close signals no more jobs will be submitted
// Workers will finish processing queued jobs and then stop
func (p *WorkerPool) Close() {
	// TODO: implement
}

// QueueDepth returns current number of jobs in queue
func (p *WorkerPool) QueueDepth() int {
	// TODO: implement
	return 0
}

// QueueUtilization returns queue fullness as a percentage (0.0 to 1.0)
func (p *WorkerPool) QueueUtilization() float64 {
	// TODO: implement
	return 0.0
}

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	// TODO: Add fields
	// Hint: You'll need a tokens channel, rate, and capacity
}

// NewRateLimiter creates a rate limiter
//
// Parameters:
//   - requestsPerSecond: Maximum requests allowed per second
//
// Returns a configured RateLimiter
//
// Behavior:
//   - Bucket starts full (allows initial burst)
//   - Tokens are added at constant rate
//   - Bucket never exceeds capacity
func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	// TODO: implement
	return nil
}

// Wait blocks until a token is available or context is cancelled
//
// Parameters:
//   - ctx: Context for cancellation
//
// Returns:
//   - error: Context error if cancelled, nil if token acquired
func (rl *RateLimiter) Wait(ctx context.Context) error {
	// TODO: implement
	return nil
}

// TryAcquire attempts to get a token without blocking
//
// Returns:
//   - bool: true if token acquired, false if none available
func (rl *RateLimiter) TryAcquire() bool {
	// TODO: implement
	return false
}

// Stop stops the rate limiter's token refill goroutine
func (rl *RateLimiter) Stop() {
	// TODO: implement
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
