//go:build !solution && !reference

package workerpoolwithbackpressure

import (
	"context"
	"sync"
	"time"
)

type Job struct {
	ID      int
	Payload string
}

type Result struct {
	JobID int
	Data  string
	Err   error
}

type WorkerPool struct {
	jobs       chan Job
	results    chan Result
	numWorkers int
	wg         sync.WaitGroup
}

type RateLimiter struct {
	tokens   chan struct{}
	rate     time.Duration
	capacity int
	stop     chan struct{}
	wg       sync.WaitGroup
}

type QueueFullError struct{}

// NewWorkerPool implements the exercise.
//
// TODO: Implement this function
func NewWorkerPool(queueSize int, numWorkers int) *WorkerPool {
	// TODO: Implement
	return nil
}

// Start implements the exercise.
//
// TODO: Implement this function
func (p *WorkerPool) Start(ctx context.Context, process func(Job) Result) {
	// TODO: Implement
}

// Submit implements the exercise.
//
// TODO: Implement this function
func (p *WorkerPool) Submit(job Job) error {
	// TODO: Implement
	return nil
}

// SubmitWithTimeout implements the exercise.
//
// TODO: Implement this function
func (p *WorkerPool) SubmitWithTimeout(ctx context.Context, job Job, timeout time.Duration) error {
	// TODO: Implement
	return nil
}

// Results implements the exercise.
//
// TODO: Implement this function
func (p *WorkerPool) Results() <-chan Result {
	// TODO: Implement
	return Result{}
}

// Close implements the exercise.
//
// TODO: Implement this function
func (p *WorkerPool) Close() {
	// TODO: Implement
}

// QueueDepth implements the exercise.
//
// TODO: Implement this function
func (p *WorkerPool) QueueDepth() int {
	// TODO: Implement
	return 0
}

// QueueUtilization implements the exercise.
//
// TODO: Implement this function
func (p *WorkerPool) QueueUtilization() float64 {
	// TODO: Implement
	return 0
}

// NewRateLimiter implements the exercise.
//
// TODO: Implement this function
func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	// TODO: Implement
	return nil
}

// Wait implements the exercise.
//
// TODO: Implement this function
func (rl *RateLimiter) Wait(ctx context.Context) error {
	// TODO: Implement
	return nil
}

// TryAcquire implements the exercise.
//
// TODO: Implement this function
func (rl *RateLimiter) TryAcquire() bool {
	// TODO: Implement
	return false
}

// Stop implements the exercise.
//
// TODO: Implement this function
func (rl *RateLimiter) Stop() {
	// TODO: Implement
}

// Error implements the exercise.
//
// TODO: Implement this function
func (e *QueueFullError) Error() string {
	// TODO: Implement
	return ""
}
