//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
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
// TODO: implement NewWorkerPool.
func NewWorkerPool(queueSize, numWorkers int) *WorkerPool { panic("TODO: implement") }
// TODO: implement Start.
func (p *WorkerPool) Start(ctx context.Context, process func(Job) Result) { panic("TODO: implement") }
// TODO: implement Submit.
func (p *WorkerPool) Submit(job Job) error { panic("TODO: implement") }
// TODO: implement SubmitWithTimeout.
func (p *WorkerPool) SubmitWithTimeout(ctx context.Context, job Job, timeout time.Duration) error {
	panic("TODO: implement")
}
// TODO: implement Results.
func (p *WorkerPool) Results() <-chan Result { panic("TODO: implement") }
// TODO: implement Close.
func (p *WorkerPool) Close() { panic("TODO: implement") }
// TODO: implement QueueDepth.
func (p *WorkerPool) QueueDepth() int { panic("TODO: implement") }
// TODO: implement QueueUtilization.
func (p *WorkerPool) QueueUtilization() float64 { panic("TODO: implement") }

type RateLimiter struct {
	tokens   chan struct{}
	rate     time.Duration
	capacity int
	stop     chan struct{}
	wg       sync.WaitGroup
}
// TODO: implement NewRateLimiter.
func NewRateLimiter(requestsPerSecond int) *RateLimiter { panic("TODO: implement") }
// TODO: implement Wait.
func (rl *RateLimiter) Wait(ctx context.Context) error { panic("TODO: implement") }
// TODO: implement TryAcquire.
func (rl *RateLimiter) TryAcquire() bool { panic("TODO: implement") }
// TODO: implement Stop.
func (rl *RateLimiter) Stop() { panic("TODO: implement") }

var (
	ErrQueueFull = &QueueFullError{}
)

type QueueFullError struct{}
// TODO: implement Error.
func (e *QueueFullError) Error() string { panic("TODO: implement") }
