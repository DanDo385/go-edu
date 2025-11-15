// Package exercise provides tests for semaphore exercises.
package exercise

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ============================================================================
// TEST 1: Basic Semaphore
// ============================================================================

func TestSemaphoreBasic(t *testing.T) {
	sem := NewSemaphore(3)
	if sem == nil {
		t.Fatal("NewSemaphore returned nil")
	}

	// Should be able to acquire 3 permits
	sem.Acquire()
	sem.Acquire()
	sem.Acquire()

	// Release and acquire again
	sem.Release()
	sem.Acquire()

	// Clean up
	sem.Release()
	sem.Release()
	sem.Release()
}

func TestSemaphoreConcurrency(t *testing.T) {
	const (
		maxConcurrent = 5
		numTasks      = 20
	)

	sem := NewSemaphore(maxConcurrent)
	if sem == nil {
		t.Skip("NewSemaphore not implemented")
	}

	var (
		active    atomic.Int32
		maxActive int32
		mu        sync.Mutex
	)

	var wg sync.WaitGroup

	for i := 0; i < numTasks; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			sem.Acquire()
			defer sem.Release()

			// Track concurrent goroutines
			current := active.Add(1)
			defer active.Add(-1)

			// Update max
			mu.Lock()
			if current > maxActive {
				maxActive = current
			}
			mu.Unlock()

			time.Sleep(10 * time.Millisecond)
		}()
	}

	wg.Wait()

	if maxActive > maxConcurrent {
		t.Errorf("Max concurrent (%d) exceeded limit (%d)", maxActive, maxConcurrent)
	}

	if maxActive < maxConcurrent {
		t.Logf("Warning: Peak usage (%d) below capacity (%d) - might be timing", maxActive, maxConcurrent)
	}
}

func TestSemaphoreBlocking(t *testing.T) {
	sem := NewSemaphore(1)
	if sem == nil {
		t.Skip("NewSemaphore not implemented")
	}

	// Acquire the only permit
	sem.Acquire()

	blocked := make(chan bool, 1)

	// Try to acquire in goroutine (should block)
	go func() {
		blocked <- true
		sem.Acquire()
		blocked <- false
	}()

	// Verify it's blocked
	select {
	case <-blocked:
		// Good, goroutine started
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Goroutine didn't start")
	}

	// Should still be blocked
	select {
	case <-blocked:
		t.Fatal("Acquire didn't block")
	case <-time.After(100 * time.Millisecond):
		// Good, still blocked
	}

	// Release and verify unblocked
	sem.Release()

	select {
	case <-blocked:
		// Good, unblocked
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Acquire didn't unblock after release")
	}

	// Clean up
	sem.Release()
}

// ============================================================================
// TEST 2: Try-Acquire
// ============================================================================

func TestTryAcquire(t *testing.T) {
	sem := NewSemaphore(2)
	if sem == nil {
		t.Skip("NewSemaphore not implemented")
	}

	// Should succeed when permits available
	if !sem.TryAcquire() {
		t.Error("TryAcquire failed when permits available")
	}

	if !sem.TryAcquire() {
		t.Error("TryAcquire failed when permits available")
	}

	// Should fail when full
	if sem.TryAcquire() {
		t.Error("TryAcquire succeeded when semaphore full")
	}

	// Release and try again
	sem.Release()

	if !sem.TryAcquire() {
		t.Error("TryAcquire failed after release")
	}

	// Clean up
	sem.Release()
	sem.Release()
}

func TestTryAcquireNonBlocking(t *testing.T) {
	sem := NewSemaphore(1)
	if sem == nil {
		t.Skip("NewSemaphore not implemented")
	}

	sem.Acquire()

	// TryAcquire should return immediately, not block
	done := make(chan bool)

	go func() {
		sem.TryAcquire()
		done <- true
	}()

	select {
	case <-done:
		// Good, returned immediately
	case <-time.After(100 * time.Millisecond):
		t.Fatal("TryAcquire blocked (should be non-blocking)")
	}

	sem.Release()
}

// ============================================================================
// TEST 3: Context-Aware Acquisition
// ============================================================================

func TestAcquireWithContext(t *testing.T) {
	sem := NewSemaphore(1)
	if sem == nil {
		t.Skip("NewSemaphore not implemented")
	}

	// Should succeed with valid context
	ctx := context.Background()
	if err := sem.AcquireWithContext(ctx); err != nil {
		t.Errorf("AcquireWithContext failed: %v", err)
	}
	sem.Release()
}

func TestAcquireWithContextTimeout(t *testing.T) {
	sem := NewSemaphore(1)
	if sem == nil {
		t.Skip("NewSemaphore not implemented")
	}

	// Acquire the permit
	sem.Acquire()

	// Try to acquire with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := sem.AcquireWithContext(ctx)
	if err == nil {
		t.Error("AcquireWithContext should timeout")
	}

	if err != context.DeadlineExceeded {
		t.Errorf("Expected DeadlineExceeded, got %v", err)
	}

	sem.Release()
}

func TestAcquireWithContextCancellation(t *testing.T) {
	sem := NewSemaphore(1)
	if sem == nil {
		t.Skip("NewSemaphore not implemented")
	}

	sem.Acquire()

	ctx, cancel := context.WithCancel(context.Background())

	// Cancel immediately
	cancel()

	err := sem.AcquireWithContext(ctx)
	if err == nil {
		t.Error("AcquireWithContext should fail when context cancelled")
	}

	if err != context.Canceled {
		t.Errorf("Expected Canceled, got %v", err)
	}

	sem.Release()
}

// ============================================================================
// TEST 4: Rate Limiter
// ============================================================================

func TestRateLimiterBurst(t *testing.T) {
	limiter := NewRateLimiter(5, 100*time.Millisecond)
	if limiter == nil {
		t.Skip("NewRateLimiter not implemented")
	}
	defer limiter.Stop()

	// Should allow burst of 5 immediately
	start := time.Now()

	for i := 0; i < 5; i++ {
		limiter.Wait()
	}

	elapsed := time.Since(start)

	// Burst should complete quickly (< 50ms)
	if elapsed > 50*time.Millisecond {
		t.Errorf("Burst took too long: %v", elapsed)
	}
}

func TestRateLimiterSustained(t *testing.T) {
	limiter := NewRateLimiter(2, 100*time.Millisecond)
	if limiter == nil {
		t.Skip("NewRateLimiter not implemented")
	}
	defer limiter.Stop()

	// Consume initial burst
	limiter.Wait()
	limiter.Wait()

	// Next wait should take ~100ms (refill time)
	start := time.Now()
	limiter.Wait()
	elapsed := time.Since(start)

	if elapsed < 80*time.Millisecond {
		t.Errorf("Rate limiting not working: too fast (%v)", elapsed)
	}

	if elapsed > 150*time.Millisecond {
		t.Errorf("Rate limiting too slow: %v", elapsed)
	}
}

func TestRateLimiterTryAcquire(t *testing.T) {
	limiter := NewRateLimiter(2, 100*time.Millisecond)
	if limiter == nil {
		t.Skip("NewRateLimiter not implemented")
	}
	defer limiter.Stop()

	// Should succeed for burst
	if !limiter.TryAcquire() {
		t.Error("TryAcquire failed during burst")
	}

	if !limiter.TryAcquire() {
		t.Error("TryAcquire failed during burst")
	}

	// Should fail when exhausted
	if limiter.TryAcquire() {
		t.Error("TryAcquire succeeded when rate limit exhausted")
	}

	// Wait for refill and try again
	time.Sleep(150 * time.Millisecond)

	if !limiter.TryAcquire() {
		t.Error("TryAcquire failed after refill")
	}
}

// ============================================================================
// TEST 5: Weighted Semaphore
// ============================================================================

func TestWeightedSemaphoreBasic(t *testing.T) {
	sem := NewWeightedSemaphore(10)
	if sem == nil {
		t.Skip("NewWeightedSemaphore not implemented")
	}

	// Acquire different weights
	sem.Acquire(3)
	sem.Acquire(5)

	// Should have 2 remaining (10 - 3 - 5)
	// Try to acquire 3 (should block)
	acquired := make(chan bool, 1)

	go func() {
		sem.Acquire(3)
		acquired <- true
	}()

	// Should block
	select {
	case <-acquired:
		t.Error("Acquire(3) should block when only 2 available")
	case <-time.After(100 * time.Millisecond):
		// Good, blocked
	}

	// Release some
	sem.Release(3)

	// Now should succeed
	select {
	case <-acquired:
		// Good, unblocked
	case <-time.After(100 * time.Millisecond):
		t.Error("Acquire didn't unblock after release")
	}

	// Clean up
	sem.Release(5)
	sem.Release(3)
}

func TestWeightedSemaphoreContextCancel(t *testing.T) {
	sem := NewWeightedSemaphore(10)
	if sem == nil {
		t.Skip("NewWeightedSemaphore not implemented")
	}

	// Fill semaphore
	sem.Acquire(10)

	// Try to acquire with cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := sem.AcquireWithContext(ctx, 5)
	if err == nil {
		t.Error("AcquireWithContext should fail when cancelled")
	}

	if err != context.Canceled {
		t.Errorf("Expected Canceled, got %v", err)
	}

	// Semaphore should still be intact (no leaked permits)
	sem.Release(10)

	// Should be able to acquire full capacity
	if err := sem.AcquireWithContext(context.Background(), 10); err != nil {
		t.Errorf("Failed to acquire after cancelled acquisition: %v", err)
	}

	sem.Release(10)
}

func TestWeightedSemaphorePartialAcquisition(t *testing.T) {
	sem := NewWeightedSemaphore(10)
	if sem == nil {
		t.Skip("NewWeightedSemaphore not implemented")
	}

	// Acquire 8, leaving 2
	sem.Acquire(8)

	// Try to acquire 5 with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// This should fail (only 2 available, need 5)
	err := sem.AcquireWithContext(ctx, 5)
	if err == nil {
		t.Error("Should timeout waiting for permits")
	}

	// Critical test: semaphore should be intact
	// Release the 8, should be able to acquire 10
	sem.Release(8)

	if err := sem.AcquireWithContext(context.Background(), 10); err != nil {
		t.Errorf("Semaphore corrupted after partial acquisition: %v", err)
	}

	sem.Release(10)
}

// ============================================================================
// TEST 6: Worker Pool
// ============================================================================

func TestWorkerPoolBasic(t *testing.T) {
	pool := NewWorkerPool(3, DefaultProcessor)
	if pool == nil {
		t.Skip("NewWorkerPool not implemented")
	}

	pool.Start()

	// Submit jobs
	for i := 1; i <= 5; i++ {
		pool.Submit(Job{ID: i, Data: "test"})
	}

	// Collect results
	results := make(map[int]bool)
	timeout := time.After(1 * time.Second)

	for i := 0; i < 5; i++ {
		select {
		case result := <-pool.Results():
			if result.Err != nil {
				t.Errorf("Job %d failed: %v", result.JobID, result.Err)
			}
			results[result.JobID] = true
		case <-timeout:
			t.Fatal("Timeout waiting for results")
		}
	}

	pool.Stop()

	// Verify all jobs processed
	for i := 1; i <= 5; i++ {
		if !results[i] {
			t.Errorf("Job %d not processed", i)
		}
	}
}

func TestWorkerPoolConcurrency(t *testing.T) {
	var active atomic.Int32
	var maxActive int32
	var mu sync.Mutex

	processor := func(job Job) Result {
		current := active.Add(1)
		defer active.Add(-1)

		mu.Lock()
		if current > maxActive {
			maxActive = current
		}
		mu.Unlock()

		time.Sleep(50 * time.Millisecond)

		return Result{JobID: job.ID, Output: "done"}
	}

	pool := NewWorkerPool(3, processor)
	if pool == nil {
		t.Skip("NewWorkerPool not implemented")
	}

	pool.Start()

	// Submit many jobs
	for i := 1; i <= 10; i++ {
		pool.Submit(Job{ID: i, Data: "test"})
	}

	// Collect results
	for i := 0; i < 10; i++ {
		<-pool.Results()
	}

	pool.Stop()

	if maxActive > 3 {
		t.Errorf("Max concurrent workers (%d) exceeded limit (3)", maxActive)
	}

	if maxActive < 3 {
		t.Logf("Warning: Peak workers (%d) below limit (3)", maxActive)
	}
}

func TestWorkerPoolGracefulShutdown(t *testing.T) {
	pool := NewWorkerPool(2, DefaultProcessor)
	if pool == nil {
		t.Skip("NewWorkerPool not implemented")
	}

	pool.Start()

	// Submit jobs
	for i := 1; i <= 5; i++ {
		pool.Submit(Job{ID: i, Data: "test"})
	}

	// Stop should wait for all jobs to complete
	done := make(chan bool)

	go func() {
		pool.Stop()
		done <- true
	}()

	// Drain results
	count := 0
	timeout := time.After(2 * time.Second)

	for {
		select {
		case <-pool.Results():
			count++
			if count == 5 {
				// All results received, Stop should complete soon
				select {
				case <-done:
					return // Success
				case <-time.After(100 * time.Millisecond):
					t.Error("Stop() didn't complete after all jobs finished")
					return
				}
			}
		case <-timeout:
			t.Fatalf("Timeout: only %d/5 jobs completed", count)
		}
	}
}

// ============================================================================
// BENCHMARKS
// ============================================================================

func BenchmarkSemaphore(b *testing.B) {
	sem := NewSemaphore(100)
	if sem == nil {
		b.Skip("NewSemaphore not implemented")
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sem.Acquire()
			sem.Release()
		}
	})
}

func BenchmarkTryAcquire(b *testing.B) {
	sem := NewSemaphore(100)
	if sem == nil {
		b.Skip("NewSemaphore not implemented")
	}

	// Fill semaphore
	for i := 0; i < 100; i++ {
		sem.Acquire()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sem.TryAcquire()
	}
}

func BenchmarkWeightedSemaphore(b *testing.B) {
	sem := NewWeightedSemaphore(1000)
	if sem == nil {
		b.Skip("NewWeightedSemaphore not implemented")
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sem.Acquire(5)
			sem.Release(5)
		}
	})
}
