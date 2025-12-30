// Package exercise provides complete solutions for semaphore exercises.
//
// This file contains reference implementations. Students should work in
// exercise.go and refer to these solutions only after attempting the exercises.

package exercise

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// SOLUTION 1: Basic Counting Semaphore
// ============================================================================

// SemaphoreSolution is a counting semaphore using buffered channels.
//
// IMPLEMENTATION NOTES:
// - Buffered channel capacity = max permits
// - Send (sem <-) = acquire permit (blocks when full)
// - Receive (<-sem) = release permit (makes space)
// - Simple, idiomatic Go pattern
type SemaphoreSolution struct {
	sem chan struct{}
}

// NewSemaphoreSolution creates a new counting semaphore.
func NewSemaphoreSolution(maxPermits int) *SemaphoreSolution {
	return &SemaphoreSolution{
		sem: make(chan struct{}, maxPermits),
	}
}

// Acquire acquires a permit, blocking if none available.
func (s *SemaphoreSolution) Acquire() {
	s.sem <- struct{}{}
}

// Release releases a permit back to the semaphore.
func (s *SemaphoreSolution) Release() {
	<-s.sem
}

// TryAcquire attempts to acquire without blocking.
func (s *SemaphoreSolution) TryAcquire() bool {
	select {
	case s.sem <- struct{}{}:
		return true
	default:
		return false
	}
}

// AcquireWithContext acquires with timeout/cancellation support.
func (s *SemaphoreSolution) AcquireWithContext(ctx context.Context) error {
	select {
	case s.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// ============================================================================
// SOLUTION 2: Rate Limiter
// ============================================================================

// RateLimiterSolution implements token bucket rate limiting.
//
// IMPLEMENTATION NOTES:
// - Tokens channel holds available permits (buffered = maxBurst)
// - Refill goroutine adds tokens at specified rate
// - Wait() blocks until token available
// - TryAcquire() non-blocking attempt
type RateLimiterSolution struct {
	tokens chan struct{}
	rate   time.Duration
	done   chan struct{}
}

// NewRateLimiterSolution creates a new rate limiter.
func NewRateLimiterSolution(maxBurst int, rate time.Duration) *RateLimiterSolution {
	rl := &RateLimiterSolution{
		tokens: make(chan struct{}, maxBurst),
		rate:   rate,
		done:   make(chan struct{}),
	}

	// Fill initial tokens (allow burst)
	for i := 0; i < maxBurst; i++ {
		rl.tokens <- struct{}{}
	}

	// Start refill goroutine
	go rl.refill()

	return rl
}

// refill periodically adds tokens to the bucket.
func (rl *RateLimiterSolution) refill() {
	ticker := time.NewTicker(rl.rate)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Try to add token (non-blocking)
			select {
			case rl.tokens <- struct{}{}:
				// Token added
			default:
				// Bucket full, skip
			}
		case <-rl.done:
			return
		}
	}
}

// Wait blocks until a token is available.
func (rl *RateLimiterSolution) Wait() {
	<-rl.tokens
}

// TryAcquire attempts non-blocking token acquisition.
func (rl *RateLimiterSolution) TryAcquire() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

// Stop stops the rate limiter.
func (rl *RateLimiterSolution) Stop() {
	close(rl.done)
}

// ============================================================================
// SOLUTION 3: Weighted Semaphore
// ============================================================================

// WeightedSemaphoreSolution allows acquiring multiple permits at once.
//
// IMPLEMENTATION NOTES:
// - Uses buffered channel where capacity = max total weight
// - Acquire(n) sends n items to channel
// - Release(n) receives n items from channel
// - Context support requires careful cleanup on partial acquisition
type WeightedSemaphoreSolution struct {
	permits chan struct{}
}

// NewWeightedSemaphoreSolution creates a weighted semaphore.
func NewWeightedSemaphoreSolution(maxWeight int) *WeightedSemaphoreSolution {
	return &WeightedSemaphoreSolution{
		permits: make(chan struct{}, maxWeight),
	}
}

// Acquire acquires the specified weight of permits.
func (ws *WeightedSemaphoreSolution) Acquire(weight int) {
	for i := 0; i < weight; i++ {
		ws.permits <- struct{}{}
	}
}

// Release releases the specified weight of permits.
func (ws *WeightedSemaphoreSolution) Release(weight int) {
	for i := 0; i < weight; i++ {
		<-ws.permits
	}
}

// AcquireWithContext acquires with context support.
//
// CRITICAL DETAIL: If context cancels during acquisition, we must
// release the permits we already acquired to avoid leaks.
func (ws *WeightedSemaphoreSolution) AcquireWithContext(ctx context.Context, weight int) error {
	acquired := 0

	// Acquire permits one at a time
	for i := 0; i < weight; i++ {
		select {
		case ws.permits <- struct{}{}:
			acquired++
		case <-ctx.Done():
			// Context cancelled: release what we acquired
			for j := 0; j < acquired; j++ {
				<-ws.permits
			}
			return ctx.Err()
		}
	}

	return nil
}

// ============================================================================
// SOLUTION 4: Worker Pool
// ============================================================================

// WorkerPoolSolution processes jobs with bounded concurrency.
//
// IMPLEMENTATION NOTES:
// - Uses semaphore to limit concurrent workers
// - Job queue channel for submitted jobs
// - Results channel for processed results
// - Graceful shutdown waits for all workers to finish
type WorkerPoolSolution struct {
	jobs      chan Job
	results   chan Result
	sem       chan struct{}
	numWorkers int
	processor func(Job) Result
	wg        sync.WaitGroup
	started   bool
	mu        sync.Mutex
}

// NewWorkerPoolSolution creates a worker pool.
func NewWorkerPoolSolution(numWorkers int, processor func(Job) Result) *WorkerPoolSolution {
	return &WorkerPoolSolution{
		jobs:       make(chan Job, numWorkers*2), // Buffered job queue
		results:    make(chan Result, numWorkers*2),
		sem:        make(chan struct{}, numWorkers),
		numWorkers: numWorkers,
		processor:  processor,
	}
}

// Submit submits a job to the pool.
func (wp *WorkerPoolSolution) Submit(job Job) {
	wp.jobs <- job
}

// Start starts processing jobs.
func (wp *WorkerPoolSolution) Start() {
	wp.mu.Lock()
	if wp.started {
		wp.mu.Unlock()
		return
	}
	wp.started = true
	wp.mu.Unlock()

	wp.wg.Add(1)

	go func() {
		defer wp.wg.Done()

		for job := range wp.jobs {
			wp.sem <- struct{}{} // Acquire worker slot

			go func(j Job) {
				defer func() { <-wp.sem }() // Release worker slot

				result := wp.processor(j)
				wp.results <- result
			}(job)
		}

		// Wait for all workers to finish
		for i := 0; i < wp.numWorkers; i++ {
			wp.sem <- struct{}{}
		}

		close(wp.results)
	}()
}

// Results returns the results channel.
func (wp *WorkerPoolSolution) Results() <-chan Result {
	return wp.results
}

// Stop gracefully stops the pool.
func (wp *WorkerPoolSolution) Stop() {
	close(wp.jobs)
	wp.wg.Wait()
}

// ============================================================================
// ADDITIONAL SOLUTIONS: Advanced Patterns
// ============================================================================

// MonitoredSemaphoreSolution tracks usage statistics.
type MonitoredSemaphoreSolution struct {
	sem           chan struct{}
	acquired      int
	capacity      int
	peakUsage     int
	totalAcquires int
	totalReleases int
	mu            sync.Mutex
}

// NewMonitoredSemaphoreSolution creates a semaphore with metrics.
func NewMonitoredSemaphoreSolution(capacity int) *MonitoredSemaphoreSolution {
	return &MonitoredSemaphoreSolution{
		sem:      make(chan struct{}, capacity),
		capacity: capacity,
	}
}

// Acquire acquires a permit and updates metrics.
func (ms *MonitoredSemaphoreSolution) Acquire() {
	ms.sem <- struct{}{}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.acquired++
	ms.totalAcquires++

	if ms.acquired > ms.peakUsage {
		ms.peakUsage = ms.acquired
	}
}

// Release releases a permit and updates metrics.
func (ms *MonitoredSemaphoreSolution) Release() {
	<-ms.sem

	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.acquired--
	ms.totalReleases++
}

// GetStats returns current statistics.
func (ms *MonitoredSemaphoreSolution) GetStats() Stats {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	return Stats{
		Acquired:      ms.acquired,
		Capacity:      ms.capacity,
		PeakUsage:     ms.peakUsage,
		TotalAcquires: ms.totalAcquires,
		TotalReleases: ms.totalReleases,
	}
}

// ============================================================================
// HELPER: Connection Pool Example
// ============================================================================

// ConnectionPoolSolution demonstrates real-world semaphore usage.
type ConnectionPoolSolution struct {
	sem      chan struct{}
	maxConns int
}

// NewConnectionPoolSolution creates a connection pool.
func NewConnectionPoolSolution(maxConns int) *ConnectionPoolSolution {
	return &ConnectionPoolSolution{
		sem:      make(chan struct{}, maxConns),
		maxConns: maxConns,
	}
}

// Acquire acquires a connection permit.
func (cp *ConnectionPoolSolution) Acquire() {
	cp.sem <- struct{}{}
}

// Release releases a connection permit.
func (cp *ConnectionPoolSolution) Release() {
	<-cp.sem
}

// AcquireWithTimeout acquires with timeout.
func (cp *ConnectionPoolSolution) AcquireWithTimeout(timeout time.Duration) error {
	select {
	case cp.sem <- struct{}{}:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("timeout acquiring connection")
	}
}

// ExecuteQuery simulates a database query with connection pooling.
func (cp *ConnectionPoolSolution) ExecuteQuery(ctx context.Context, query string) error {
	// Acquire connection
	select {
	case cp.sem <- struct{}{}:
		defer func() { <-cp.sem }()
	case <-ctx.Done():
		return ctx.Err()
	}

	// Simulate query execution
	time.Sleep(10 * time.Millisecond)
	return nil
}

// ============================================================================
// BENCHMARKING HELPERS
// ============================================================================

// BenchmarkHelper provides utilities for performance testing.
type BenchmarkHelper struct {
	sem *SemaphoreSolution
}

// NewBenchmarkHelper creates a benchmark helper.
func NewBenchmarkHelper(capacity int) *BenchmarkHelper {
	return &BenchmarkHelper{
		sem: NewSemaphoreSolution(capacity),
	}
}

// AcquireRelease performs acquire/release cycle.
func (bh *BenchmarkHelper) AcquireRelease() {
	bh.sem.Acquire()
	bh.sem.Release()
}

// ConcurrentAcquireRelease performs concurrent acquire/release.
func (bh *BenchmarkHelper) ConcurrentAcquireRelease(n int) {
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			bh.sem.Acquire()
			bh.sem.Release()
		}()
	}

	wg.Wait()
}
