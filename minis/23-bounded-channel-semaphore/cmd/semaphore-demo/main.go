// Package main demonstrates buffered channels as semaphores for resource limiting.
//
// MACRO-COMMENT: What This Program Teaches
// =========================================
// This program demonstrates:
// 1. Buffered channels as counting semaphores (idiomatic Go pattern)
// 2. Limiting concurrent access to resources (databases, APIs, files)
// 3. Rate limiting patterns (token bucket, leaky bucket)
// 4. Weighted semaphores (variable resource costs)
// 5. Context-aware acquisition (timeouts, cancellation)
// 6. Try-acquire patterns (non-blocking, graceful degradation)
// 7. Worker pool with semaphore-based concurrency control
//
// CORE INSIGHT:
// A buffered channel IS a counting semaphore:
//   - Channel capacity = max permits
//   - Send (ch <- x) = acquire permit
//   - Receive (<-ch) = release permit
//   - Full buffer = no permits available (blocks)
//
// This is one of Go's most elegant patterns: using simple primitives
// (channels) to solve complex problems (resource limiting).

package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// SECTION 1: Basic Semaphore Pattern
// ============================================================================

// demonstrateBasicSemaphore shows buffered channel as counting semaphore.
//
// MACRO-COMMENT: Counting Semaphore Fundamentals
// A counting semaphore allows N goroutines to hold a "permit" concurrently.
//
// OPERATIONS:
// - Acquire: Decrement available permits (blocks if 0)
// - Release: Increment available permits (wakes blocked goroutines)
//
// BUFFERED CHANNEL MAPPING:
// - Buffer capacity = max permits (e.g., 5)
// - Items in buffer = acquired permits
// - Empty slots = available permits
// - Send = acquire (blocks when full = no permits)
// - Receive = release (makes space = returns permit)
//
// VISUAL:
//   make(chan struct{}, 5) → [_____] (5 available)
//   Acquire x3             → [XXX__] (2 available)
//   Acquire x2             → [XXXXX] (0 available, full)
//   Acquire x1             → BLOCKS (must wait for release)
//   Release x1             → [XXXX_] (1 available, blocked goroutine wakes)
func demonstrateBasicSemaphore() {
	fmt.Println("=== Basic Semaphore Pattern ===")

	const (
		numTasks      = 10
		maxConcurrent = 3
	)

	// SEMAPHORE: Buffered channel with capacity = max concurrent
	// This limits concurrent execution to 3 goroutines at a time
	sem := make(chan struct{}, maxConcurrent)

	// MICRO-COMMENT: Track which tasks are running
	var running atomic.Int32

	fmt.Printf("  Launching %d tasks (max %d concurrent)...\n\n", numTasks, maxConcurrent)

	for i := 1; i <= numTasks; i++ {
		// ACQUIRE: Send to semaphore (blocks if 3 already running)
		// This is the "P" operation in classical semaphore terminology
		sem <- struct{}{}

		go func(id int) {
			// RELEASE: Always release on exit (even if panic)
			// This is the "V" operation in classical semaphore terminology
			defer func() { <-sem }()

			active := running.Add(1)
			fmt.Printf("    Task %2d: Started  (active: %d)\n", id, active)

			// MICRO-COMMENT: Simulate work
			time.Sleep(time.Duration(500+rand.Intn(500)) * time.Millisecond)

			active = running.Add(-1)
			fmt.Printf("    Task %2d: Complete (active: %d)\n", id, active)
		}(i)
	}

	// MACRO-COMMENT: Wait for All Tasks to Complete
	// To ensure all goroutines finish, we acquire all permits.
	// When we successfully acquire all N permits, we know all tasks are done
	// because each task holds exactly one permit while running.
	fmt.Println("\n  Waiting for all tasks to complete...")
	for i := 0; i < maxConcurrent; i++ {
		sem <- struct{}{}
	}

	fmt.Println("  All tasks complete!\n")
}

// ============================================================================
// SECTION 2: Real-World Pattern - Database Connection Pool
// ============================================================================

// DBPool simulates a database connection pool using a semaphore.
//
// MACRO-COMMENT: Connection Pool Pattern
// Problem: Database allows max N connections, but we have M >> N goroutines.
// Solution: Use semaphore to limit concurrent connections to N.
//
// BENEFITS:
// - Prevents connection exhaustion
// - Automatic queuing when pool busy
// - Graceful degradation under load
type DBPool struct {
	sem        chan struct{}
	maxConns   int
	activeConn atomic.Int32
}

// NewDBPool creates a connection pool with max connections.
func NewDBPool(maxConns int) *DBPool {
	return &DBPool{
		sem:      make(chan struct{}, maxConns),
		maxConns: maxConns,
	}
}

// Acquire gets a connection permit (blocks if pool full).
func (p *DBPool) Acquire() {
	p.sem <- struct{}{}
	p.activeConn.Add(1)
}

// Release returns a connection permit to the pool.
func (p *DBPool) Release() {
	<-p.sem
	p.activeConn.Add(-1)
}

// AcquireWithContext acquires with timeout/cancellation support.
//
// MICRO-COMMENT: Context-Aware Acquisition
// This prevents goroutines from waiting forever for a permit.
// Critical for production systems with request timeouts.
func (p *DBPool) AcquireWithContext(ctx context.Context) error {
	select {
	case p.sem <- struct{}{}:
		p.activeConn.Add(1)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Query simulates a database query with connection pooling.
func (p *DBPool) Query(id int, query string) error {
	// ACQUIRE: Get connection from pool
	p.Acquire()
	defer p.Release()

	active := p.activeConn.Load()
	fmt.Printf("    Query %2d: Executing (pool: %d/%d)\n", id, active, p.maxConns)

	// MICRO-COMMENT: Simulate query execution
	time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)

	return nil
}

// demonstrateConnectionPool shows connection pooling with semaphores.
func demonstrateConnectionPool() {
	fmt.Println("=== Connection Pool Pattern ===")

	const (
		numQueries = 15
		maxConns   = 5
	)

	pool := NewDBPool(maxConns)

	fmt.Printf("  Database: Max %d connections\n", maxConns)
	fmt.Printf("  Executing %d queries concurrently...\n\n", numQueries)

	var wg sync.WaitGroup

	for i := 1; i <= numQueries; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			if err := pool.Query(id, "SELECT * FROM users"); err != nil {
				fmt.Printf("    Query %2d: Failed - %v\n", id, err)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("\n  All queries complete!\n")
}

// ============================================================================
// SECTION 3: Rate Limiting Patterns
// ============================================================================

// TokenBucketLimiter implements token bucket rate limiting.
//
// MACRO-COMMENT: Token Bucket Algorithm
// Allows bursts up to maxBurst, then limits to sustained rate.
//
// HOW IT WORKS:
// - Bucket holds up to maxBurst tokens
// - Tokens refilled at fixed rate
// - Request consumes 1 token (blocks if empty)
//
// BEHAVIOR:
// - Burst: First maxBurst requests succeed immediately
// - Sustained: After burst, rate limited to 1/rate requests
//
// EXAMPLE:
//   maxBurst=10, rate=100ms → 10 immediate, then 10/sec sustained
type TokenBucketLimiter struct {
	tokens chan struct{}
	rate   time.Duration
	done   chan struct{}
}

// NewTokenBucketLimiter creates a rate limiter with burst and sustained rate.
func NewTokenBucketLimiter(maxBurst int, rate time.Duration) *TokenBucketLimiter {
	rl := &TokenBucketLimiter{
		tokens: make(chan struct{}, maxBurst),
		rate:   rate,
		done:   make(chan struct{}),
	}

	// MICRO-COMMENT: Fill initial tokens (allow burst)
	for i := 0; i < maxBurst; i++ {
		rl.tokens <- struct{}{}
	}

	// MACRO-COMMENT: Refill Goroutine
	// Periodically adds tokens back to bucket (rate limiting)
	go func() {
		ticker := time.NewTicker(rate)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// MICRO-COMMENT: Try to add token (non-blocking)
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
	}()

	return rl
}

// Wait blocks until a token is available.
func (rl *TokenBucketLimiter) Wait() {
	<-rl.tokens
}

// TryAcquire attempts non-blocking token acquisition.
func (rl *TokenBucketLimiter) TryAcquire() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

// Stop stops the rate limiter.
func (rl *TokenBucketLimiter) Stop() {
	close(rl.done)
}

// demonstrateRateLimiting shows token bucket rate limiting.
func demonstrateRateLimiting() {
	fmt.Println("=== Rate Limiting (Token Bucket) ===")

	const (
		numRequests = 20
		maxBurst    = 5
		rate        = 200 * time.Millisecond
	)

	limiter := NewTokenBucketLimiter(maxBurst, rate)
	defer limiter.Stop()

	fmt.Printf("  Rate limiter: burst=%d, rate=%v\n", maxBurst, rate)
	fmt.Printf("  Making %d requests...\n\n", numRequests)

	start := time.Now()

	for i := 1; i <= numRequests; i++ {
		limiter.Wait()

		elapsed := time.Since(start)
		fmt.Printf("    Request %2d: Allowed at %v\n", i, elapsed.Round(time.Millisecond))
	}

	fmt.Println()
}

// ============================================================================
// SECTION 4: Weighted Semaphore
// ============================================================================

// WeightedSemaphore allows acquiring multiple permits at once.
//
// MACRO-COMMENT: Weighted Semaphore Pattern
// Different operations have different resource costs:
// - Small task: weight 1
// - Medium task: weight 3
// - Large task: weight 5
//
// Total capacity enforced across all weights.
//
// EXAMPLE:
//   Capacity: 10
//   Acquire(3) → 7 remaining
//   Acquire(5) → 2 remaining
//   Acquire(3) → BLOCKS (only 2 available)
type WeightedSemaphore struct {
	permits chan struct{}
	mu      sync.Mutex
}

// NewWeightedSemaphore creates a weighted semaphore with max capacity.
func NewWeightedSemaphore(maxWeight int) *WeightedSemaphore {
	return &WeightedSemaphore{
		permits: make(chan struct{}, maxWeight),
	}
}

// Acquire acquires 'weight' permits (blocks if insufficient).
func (ws *WeightedSemaphore) Acquire(weight int) {
	// MICRO-COMMENT: Acquire permits one at a time
	// We can't send 'weight' items in one operation, so we loop
	for i := 0; i < weight; i++ {
		ws.permits <- struct{}{}
	}
}

// Release releases 'weight' permits.
func (ws *WeightedSemaphore) Release(weight int) {
	for i := 0; i < weight; i++ {
		<-ws.permits
	}
}

// AcquireWithContext acquires with context support.
//
// MACRO-COMMENT: Partial Acquisition Cleanup
// If context cancels while acquiring, we must release what we got
// to avoid permit leaks. This is tricky and error-prone!
func (ws *WeightedSemaphore) AcquireWithContext(ctx context.Context, weight int) error {
	acquired := 0

	for i := 0; i < weight; i++ {
		select {
		case ws.permits <- struct{}{}:
			acquired++
		case <-ctx.Done():
			// CRITICAL: Release what we acquired before returning error
			for j := 0; j < acquired; j++ {
				<-ws.permits
			}
			return ctx.Err()
		}
	}

	return nil
}

// demonstrateWeightedSemaphore shows variable resource costs.
func demonstrateWeightedSemaphore() {
	fmt.Println("=== Weighted Semaphore Pattern ===")

	const totalCapacity = 10

	sem := NewWeightedSemaphore(totalCapacity)

	fmt.Printf("  Total capacity: %d permits\n", totalCapacity)
	fmt.Println("  Running tasks with different weights...\n")

	var wg sync.WaitGroup

	// MICRO-COMMENT: Launch tasks with different weights
	tasks := []struct {
		id     int
		weight int
	}{
		{1, 2}, {2, 3}, {3, 1}, {4, 5},
		{5, 2}, {6, 4}, {7, 1}, {8, 3},
	}

	for _, task := range tasks {
		wg.Add(1)

		go func(id, weight int) {
			defer wg.Done()

			fmt.Printf("    Task %d: Waiting for %d permits...\n", id, weight)
			sem.Acquire(weight)
			defer sem.Release(weight)

			fmt.Printf("    Task %d: Acquired %d permits, working...\n", id, weight)
			time.Sleep(time.Duration(300+rand.Intn(200)) * time.Millisecond)
			fmt.Printf("    Task %d: Releasing %d permits\n", id, weight)
		}(task.id, task.weight)

		// MICRO-COMMENT: Stagger task starts for visibility
		time.Sleep(100 * time.Millisecond)
	}

	wg.Wait()
	fmt.Println("\n  All weighted tasks complete!\n")
}

// ============================================================================
// SECTION 5: Try-Acquire Pattern (Non-Blocking)
// ============================================================================

// Semaphore wraps a buffered channel with common semaphore operations.
type Semaphore struct {
	sem chan struct{}
}

// NewSemaphore creates a semaphore with max permits.
func NewSemaphore(maxPermits int) *Semaphore {
	return &Semaphore{
		sem: make(chan struct{}, maxPermits),
	}
}

// Acquire blocks until permit available.
func (s *Semaphore) Acquire() {
	s.sem <- struct{}{}
}

// Release returns a permit.
func (s *Semaphore) Release() {
	<-s.sem
}

// TryAcquire attempts non-blocking acquisition.
//
// MICRO-COMMENT: Non-Blocking Try Pattern
// Use select with default to test without blocking.
// Returns immediately with success/failure status.
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.sem <- struct{}{}:
		return true
	default:
		return false
	}
}

// AcquireTimeout waits up to timeout for permit.
func (s *Semaphore) AcquireTimeout(timeout time.Duration) bool {
	select {
	case s.sem <- struct{}{}:
		return true
	case <-time.After(timeout):
		return false
	}
}

// demonstrateTryAcquire shows non-blocking semaphore patterns.
func demonstrateTryAcquire() {
	fmt.Println("=== Try-Acquire Pattern (Non-Blocking) ===")

	sem := NewSemaphore(3)

	fmt.Println("  Filling semaphore to capacity (3 permits)...")

	// MICRO-COMMENT: Acquire all permits
	for i := 0; i < 3; i++ {
		sem.Acquire()
		fmt.Printf("    Acquired permit %d\n", i+1)
	}

	fmt.Println("\n  Attempting non-blocking acquire on full semaphore...")

	// MICRO-COMMENT: Try to acquire (should fail)
	if sem.TryAcquire() {
		fmt.Println("    ✓ Acquired (unexpected!)")
		sem.Release()
	} else {
		fmt.Println("    ✗ Failed (expected - semaphore full)")
	}

	fmt.Println("\n  Releasing one permit...")
	sem.Release()

	fmt.Println("  Attempting non-blocking acquire again...")

	// MICRO-COMMENT: Try to acquire (should succeed)
	if sem.TryAcquire() {
		fmt.Println("    ✓ Acquired (expected - space available)")
		sem.Release()
	} else {
		fmt.Println("    ✗ Failed (unexpected!)")
	}

	// MICRO-COMMENT: Try with timeout
	fmt.Println("\n  Attempting acquire with 500ms timeout...")
	sem.Release() // Make space
	sem.Release()

	if sem.AcquireTimeout(500 * time.Millisecond) {
		fmt.Println("    ✓ Acquired within timeout")
		sem.Release()
	} else {
		fmt.Println("    ✗ Timeout")
	}

	fmt.Println()
}

// ============================================================================
// SECTION 6: Worker Pool with Semaphore
// ============================================================================

// WorkerPool limits concurrent job execution.
//
// MACRO-COMMENT: Worker Pool Pattern
// Classic concurrency pattern:
// - Jobs submitted to queue
// - Workers (goroutines) process jobs
// - Semaphore limits concurrent workers
//
// BENEFITS:
// - Bounded concurrency (prevents goroutine explosion)
// - Job queuing (buffered channel)
// - Graceful shutdown
type WorkerPool struct {
	jobs    chan func()
	sem     chan struct{}
	workers int
	wg      sync.WaitGroup
}

// NewWorkerPool creates a worker pool with max concurrent workers.
func NewWorkerPool(numWorkers int, queueSize int) *WorkerPool {
	return &WorkerPool{
		jobs:    make(chan func(), queueSize),
		sem:     make(chan struct{}, numWorkers),
		workers: numWorkers,
	}
}

// Submit adds a job to the queue.
func (wp *WorkerPool) Submit(job func()) {
	wp.jobs <- job
}

// Start begins processing jobs.
func (wp *WorkerPool) Start() {
	wp.wg.Add(1)

	go func() {
		defer wp.wg.Done()

		for job := range wp.jobs {
			wp.sem <- struct{}{} // Acquire worker slot

			go func(j func()) {
				defer func() { <-wp.sem }() // Release worker slot
				j()
			}(job)
		}

		// MICRO-COMMENT: Wait for all workers to finish
		for i := 0; i < wp.workers; i++ {
			wp.sem <- struct{}{}
		}
	}()
}

// Stop gracefully shuts down the pool.
func (wp *WorkerPool) Stop() {
	close(wp.jobs)
	wp.wg.Wait()
}

// demonstrateWorkerPool shows worker pool with semaphore.
func demonstrateWorkerPool() {
	fmt.Println("=== Worker Pool with Semaphore ===")

	const (
		numWorkers = 3
		numJobs    = 10
	)

	pool := NewWorkerPool(numWorkers, 20)
	pool.Start()

	fmt.Printf("  Worker pool: %d workers, processing %d jobs...\n\n", numWorkers, numJobs)

	var activeJobs atomic.Int32

	for i := 1; i <= numJobs; i++ {
		jobID := i

		pool.Submit(func() {
			active := activeJobs.Add(1)
			fmt.Printf("    Job %2d: Started  (active: %d/%d)\n", jobID, active, numWorkers)

			// MICRO-COMMENT: Simulate work
			time.Sleep(time.Duration(300+rand.Intn(300)) * time.Millisecond)

			active = activeJobs.Add(-1)
			fmt.Printf("    Job %2d: Complete (active: %d/%d)\n", jobID, active, numWorkers)
		})
	}

	fmt.Println("\n  Waiting for all jobs to complete...")
	pool.Stop()
	fmt.Println("  All jobs complete!\n")
}

// ============================================================================
// SECTION 7: Binary Semaphore (Mutex Alternative)
// ============================================================================

// demonstrateBinarySemaphore shows using channel as mutex.
//
// MICRO-COMMENT: Binary Semaphore as Mutex
// A semaphore with capacity 1 is equivalent to a mutex.
// However, sync.Mutex is preferred for mutual exclusion.
//
// USE CASES WHERE CHANNEL IS BETTER:
// - Need select with other operations
// - Want to pass "lock" between goroutines
// - Implementing higher-level synchronization primitives
func demonstrateBinarySemaphore() {
	fmt.Println("=== Binary Semaphore (Mutex Pattern) ===")

	// BINARY SEMAPHORE: Capacity 1 = mutex
	mutex := make(chan struct{}, 1)

	counter := 0
	const numIncrements = 5

	fmt.Println("  Using binary semaphore to protect counter...\n")

	var wg sync.WaitGroup

	for i := 1; i <= numIncrements; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			// LOCK: Acquire binary semaphore
			mutex <- struct{}{}
			defer func() { <-mutex }() // UNLOCK: Release

			fmt.Printf("    Goroutine %d: Acquired lock (counter=%d)\n", id, counter)

			// CRITICAL SECTION: Only one goroutine at a time
			old := counter
			time.Sleep(100 * time.Millisecond)
			counter = old + 1

			fmt.Printf("    Goroutine %d: Releasing lock (counter=%d)\n", id, counter)
		}(i)
	}

	wg.Wait()
	fmt.Printf("\n  Final counter: %d (expected: %d)\n\n", counter, numIncrements)
}

// ============================================================================
// SECTION 8: Semaphore Metrics and Monitoring
// ============================================================================

// MonitoredSemaphore tracks usage statistics.
type MonitoredSemaphore struct {
	sem        chan struct{}
	acquired   atomic.Int64
	maxAcq     atomic.Int64
	totalAcq   atomic.Int64
	totalRel   atomic.Int64
}

// NewMonitoredSemaphore creates a semaphore with metrics.
func NewMonitoredSemaphore(capacity int) *MonitoredSemaphore {
	return &MonitoredSemaphore{
		sem: make(chan struct{}, capacity),
	}
}

// Acquire acquires a permit and updates metrics.
func (ms *MonitoredSemaphore) Acquire() {
	ms.sem <- struct{}{}

	current := ms.acquired.Add(1)
	ms.totalAcq.Add(1)

	// MICRO-COMMENT: Track peak usage
	for {
		max := ms.maxAcq.Load()
		if current <= max || ms.maxAcq.CompareAndSwap(max, current) {
			break
		}
	}
}

// Release releases a permit and updates metrics.
func (ms *MonitoredSemaphore) Release() {
	<-ms.sem
	ms.acquired.Add(-1)
	ms.totalRel.Add(1)
}

// Stats returns current statistics.
func (ms *MonitoredSemaphore) Stats() (acquired, capacity int64, peak, totalAcq, totalRel int64) {
	return ms.acquired.Load(),
		int64(cap(ms.sem)),
		ms.maxAcq.Load(),
		ms.totalAcq.Load(),
		ms.totalRel.Load()
}

// demonstrateMonitoring shows semaphore with metrics.
func demonstrateMonitoring() {
	fmt.Println("=== Semaphore Monitoring ===")

	sem := NewMonitoredSemaphore(5)

	fmt.Println("  Running tasks with metrics tracking...\n")

	var wg sync.WaitGroup

	for i := 1; i <= 15; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			sem.Acquire()
			defer sem.Release()

			// MICRO-COMMENT: Simulate varying work duration
			time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)
		}(i)

		time.Sleep(50 * time.Millisecond)
	}

	wg.Wait()

	// MICRO-COMMENT: Print final statistics
	acquired, capacity, peak, totalAcq, totalRel := sem.Stats()
	fmt.Printf("  Semaphore Statistics:\n")
	fmt.Printf("    Capacity:         %d\n", capacity)
	fmt.Printf("    Peak usage:       %d\n", peak)
	fmt.Printf("    Total acquires:   %d\n", totalAcq)
	fmt.Printf("    Total releases:   %d\n", totalRel)
	fmt.Printf("    Currently held:   %d\n", acquired)
	fmt.Println()
}

// ============================================================================
// MAIN FUNCTION
// ============================================================================

func main() {
	fmt.Println("=== Bounded Channel Semaphore Demonstration ===\n")

	// SECTION 1: Basics
	demonstrateBasicSemaphore()

	// SECTION 2: Real-world patterns
	demonstrateConnectionPool()
	demonstrateRateLimiting()

	// SECTION 3: Advanced patterns
	demonstrateWeightedSemaphore()
	demonstrateTryAcquire()

	// SECTION 4: Worker pool
	demonstrateWorkerPool()

	// SECTION 5: Special cases
	demonstrateBinarySemaphore()
	demonstrateMonitoring()

	fmt.Println("=== All Demonstrations Complete ===")
}
