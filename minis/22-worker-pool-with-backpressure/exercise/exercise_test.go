package exercise

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPool_BasicOperation(t *testing.T) {
	pool := NewWorkerPool(10, 3)
	if pool == nil {
		t.Fatal("NewWorkerPool returned nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Process function that doubles the job ID
	process := func(job Job) Result {
		return Result{
			JobID: job.ID,
			Data:  job.Payload + "_processed",
		}
	}

	pool.Start(ctx, process)

	// Submit jobs
	for i := 0; i < 5; i++ {
		err := pool.Submit(Job{ID: i, Payload: "job"})
		if err != nil {
			t.Errorf("Submit failed: %v", err)
		}
	}

	pool.Close()

	// Collect results
	results := make(map[int]Result)
	for result := range pool.Results() {
		results[result.JobID] = result
	}

	if len(results) != 5 {
		t.Errorf("Expected 5 results, got %d", len(results))
	}

	for i := 0; i < 5; i++ {
		result, ok := results[i]
		if !ok {
			t.Errorf("Missing result for job %d", i)
			continue
		}
		expected := "job_processed"
		if result.Data != expected {
			t.Errorf("Job %d: expected %q, got %q", i, expected, result.Data)
		}
	}
}

func TestWorkerPool_Backpressure(t *testing.T) {
	queueSize := 5
	pool := NewWorkerPool(queueSize, 1)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Slow processor
	process := func(job Job) Result {
		time.Sleep(100 * time.Millisecond)
		return Result{JobID: job.ID}
	}

	pool.Start(ctx, process)

	// Fill queue
	for i := 0; i < queueSize; i++ {
		err := pool.Submit(Job{ID: i})
		if err != nil {
			t.Errorf("Submit %d should succeed, got error: %v", i, err)
		}
	}

	// Next submit should fail (queue full)
	err := pool.Submit(Job{ID: queueSize})
	if !errors.Is(err, ErrQueueFull) {
		t.Errorf("Expected ErrQueueFull, got: %v", err)
	}

	pool.Close()

	// Drain results
	for range pool.Results() {
	}
}

func TestWorkerPool_SubmitWithTimeout(t *testing.T) {
	pool := NewWorkerPool(1, 1) // Queue size 1, 1 worker

	ctx := context.Background()

	// Block processing with a channel that never sends
	block := make(chan struct{})

	process := func(job Job) Result {
		<-block // Block forever
		return Result{JobID: job.ID}
	}

	pool.Start(ctx, process)

	time.Sleep(20 * time.Millisecond) // Let worker start

	// Submit first job - worker grabs it and blocks
	err := pool.Submit(Job{ID: 0})
	if err != nil {
		t.Fatalf("First submit failed: %v", err)
	}

	time.Sleep(20 * time.Millisecond)

	// Submit second job - goes into queue (1/1)
	err = pool.Submit(Job{ID: 1})
	if err != nil {
		t.Fatalf("Second submit failed: %v", err)
	}

	time.Sleep(20 * time.Millisecond)

	// Now worker is blocked with job 0, and queue has job 1
	// Next submit with timeout should fail
	err = pool.SubmitWithTimeout(ctx, Job{ID: 2}, 50*time.Millisecond)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}

	close(block) // Unblock
	pool.Close()
	for range pool.Results() {
	}
}

func TestWorkerPool_QueueMetrics(t *testing.T) {
	queueSize := 10
	pool := NewWorkerPool(queueSize, 1)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Block workers
	var mu sync.Mutex
	process := func(job Job) Result {
		mu.Lock()
		defer mu.Unlock()
		return Result{JobID: job.ID}
	}

	mu.Lock() // Block processing
	pool.Start(ctx, process)

	// Submit some jobs
	for i := 0; i < 5; i++ {
		_ = pool.Submit(Job{ID: i})
	}

	time.Sleep(50 * time.Millisecond) // Let jobs queue up

	depth := pool.QueueDepth()
	if depth < 4 || depth > 5 {
		t.Errorf("Expected queue depth 4-5, got %d", depth)
	}

	util := pool.QueueUtilization()
	if util < 0.4 || util > 0.6 {
		t.Errorf("Expected utilization 0.4-0.6, got %.2f", util)
	}

	mu.Unlock() // Unblock
	pool.Close()
	for range pool.Results() {
	}
}

func TestWorkerPool_ContextCancellation(t *testing.T) {
	pool := NewWorkerPool(20, 2) // Only 2 workers for 20 jobs

	ctx, cancel := context.WithCancel(context.Background())

	var processed atomic.Int32
	process := func(job Job) Result {
		time.Sleep(100 * time.Millisecond) // Slow processing
		processed.Add(1)
		return Result{JobID: job.ID}
	}

	pool.Start(ctx, process)

	// Submit many jobs
	for i := 0; i < 20; i++ {
		_ = pool.Submit(Job{ID: i})
	}

	time.Sleep(150 * time.Millisecond) // Let a few process
	cancel()                           // Cancel context (should stop workers)

	time.Sleep(50 * time.Millisecond) // Let cancellation propagate

	pool.Close()
	for range pool.Results() {
	}

	// Should have processed only a few (2 workers * 1-2 jobs each)
	count := processed.Load()
	if count >= 10 {
		t.Errorf("Expected fewer than 10 processed (due to cancellation), got %d", count)
	}
	if count < 2 {
		t.Errorf("Expected at least 2 processed, got %d", count)
	}
}

func TestRateLimiter_BasicOperation(t *testing.T) {
	rps := 10 // 10 requests per second
	limiter := NewRateLimiter(rps)
	if limiter == nil {
		t.Fatal("NewRateLimiter returned nil")
	}
	defer limiter.Stop()

	ctx := context.Background()

	// Should be able to acquire initial tokens quickly
	start := time.Now()
	for i := 0; i < rps; i++ {
		err := limiter.Wait(ctx)
		if err != nil {
			t.Errorf("Wait %d failed: %v", i, err)
		}
	}
	elapsed := time.Since(start)

	// Initial burst should be fast (< 100ms)
	if elapsed > 100*time.Millisecond {
		t.Errorf("Initial burst took too long: %v", elapsed)
	}

	// Next request should wait for refill
	start = time.Now()
	err := limiter.Wait(ctx)
	if err != nil {
		t.Errorf("Wait after burst failed: %v", err)
	}
	elapsed = time.Since(start)

	// Should wait ~100ms (1/10 second)
	if elapsed < 50*time.Millisecond || elapsed > 200*time.Millisecond {
		t.Errorf("Expected wait ~100ms, got %v", elapsed)
	}
}

func TestRateLimiter_TryAcquire(t *testing.T) {
	rps := 5
	limiter := NewRateLimiter(rps)
	if limiter == nil {
		t.Fatal("NewRateLimiter returned nil")
	}
	defer limiter.Stop()

	// Exhaust initial tokens
	acquired := 0
	for i := 0; i < rps+1; i++ {
		if limiter.TryAcquire() {
			acquired++
		}
	}

	if acquired != rps {
		t.Errorf("Expected to acquire %d tokens, got %d", rps, acquired)
	}

	// Should not be able to acquire more immediately
	if limiter.TryAcquire() {
		t.Error("Should not acquire token when bucket is empty")
	}

	// Wait for refill
	time.Sleep(250 * time.Millisecond)

	// Should be able to acquire again
	if !limiter.TryAcquire() {
		t.Error("Should acquire token after refill")
	}
}

func TestRateLimiter_ContextCancellation(t *testing.T) {
	limiter := NewRateLimiter(1)
	defer limiter.Stop()

	ctx, cancel := context.WithCancel(context.Background())

	// Exhaust token
	_ = limiter.Wait(ctx)

	// Cancel context
	cancel()

	// Should return error immediately
	err := limiter.Wait(ctx)
	if err == nil {
		t.Error("Expected context error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestRateLimiter_RateEnforcement(t *testing.T) {
	rps := 20
	limiter := NewRateLimiter(rps)
	defer limiter.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Try to do 40 requests (should take ~2 seconds at 20 rps)
	requests := 40
	start := time.Now()

	for i := 0; i < requests; i++ {
		err := limiter.Wait(ctx)
		if err != nil {
			t.Fatalf("Wait failed at request %d: %v", i, err)
		}
	}

	elapsed := time.Since(start)
	expectedMin := time.Duration(requests/rps-1) * time.Second
	expectedMax := time.Duration(requests/rps+1) * time.Second

	if elapsed < expectedMin || elapsed > expectedMax {
		t.Errorf("Expected duration %v-%v for %d requests at %d rps, got %v",
			expectedMin, expectedMax, requests, rps, elapsed)
	}
}

func BenchmarkWorkerPool_Submit(b *testing.B) {
	pool := NewWorkerPool(1000, 10)
	ctx := context.Background()

	process := func(job Job) Result {
		return Result{JobID: job.ID}
	}

	pool.Start(ctx, process)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pool.Submit(Job{ID: i})
	}

	pool.Close()
	for range pool.Results() {
	}
}

func BenchmarkRateLimiter_Wait(b *testing.B) {
	limiter := NewRateLimiter(1000) // High rate to minimize waiting
	defer limiter.Stop()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = limiter.Wait(ctx)
	}
}
