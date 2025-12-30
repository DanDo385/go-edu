package exercise

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// ConcurrentCounter is a thread-safe counter.
type ConcurrentCounter struct {
	value atomic.Int64
}

// GracefulWorker is a worker that can be gracefully cancelled using context.
type GracefulWorker struct {
	ctx      context.Context
	workDone atomic.Int64
}

// WorkerPool implements a fixed-size pool of workers that process jobs.
type WorkerPool struct {
	jobs    chan func()
	wg      sync.WaitGroup
	stopped atomic.Bool
}

// RateLimiter limits the rate at which operations can be performed.
type RateLimiter struct {
	ticker *time.Ticker
	tokens chan struct{}
}
