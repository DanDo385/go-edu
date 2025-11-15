// Package exercise provides semaphore implementation exercises.
package exercise

import (
	"context"
	"time"
)

// Job represents a unit of work to be processed.
type Job struct {
	ID   int
	Data string
}

// Result represents the outcome of processing a job.
type Result struct {
	JobID  int
	Output string
	Err    error
}

// Config holds rate limiter configuration.
type Config struct {
	MaxBurst int
	Rate     time.Duration
}

// Stats holds semaphore usage statistics.
type Stats struct {
	Acquired      int
	Capacity      int
	PeakUsage     int
	TotalAcquires int
	TotalReleases int
}

// SemaphoreInterface defines the semaphore operations.
type SemaphoreInterface interface {
	Acquire()
	Release()
	TryAcquire() bool
	AcquireWithContext(ctx context.Context) error
}

// RateLimiterInterface defines rate limiter operations.
type RateLimiterInterface interface {
	Wait()
	TryAcquire() bool
	Stop()
}

// WorkerPoolInterface defines worker pool operations.
type WorkerPoolInterface interface {
	Submit(job Job)
	Start()
	Stop()
	Results() <-chan Result
}
