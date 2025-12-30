package exercise

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ============================================================================
// Tests for Exercise 1: RetryWithTimeout
// ============================================================================

func TestRetryWithTimeout_Success(t *testing.T) {
	ctx := context.Background()
	attempts := 0

	fn := func(ctx context.Context) error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary error")
		}
		return nil // Success on 3rd attempt
	}

	err := RetryWithTimeout(ctx, fn, 5, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestRetryWithTimeout_AllFail(t *testing.T) {
	ctx := context.Background()
	attempts := 0

	fn := func(ctx context.Context) error {
		attempts++
		return errors.New("permanent error")
	}

	err := RetryWithTimeout(ctx, fn, 3, 100*time.Millisecond)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestRetryWithTimeout_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	attempts := 0

	fn := func(ctx context.Context) error {
		attempts++
		return errors.New("error")
	}

	// Cancel after first attempt
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := RetryWithTimeout(ctx, fn, 10, 100*time.Millisecond)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled, got %v", err)
	}

	// Should stop early due to cancellation
	if attempts >= 10 {
		t.Errorf("Expected fewer than 10 attempts due to cancellation, got %d", attempts)
	}
}

func TestRetryWithTimeout_IndividualTimeout(t *testing.T) {
	ctx := context.Background()
	attempts := 0

	fn := func(ctx context.Context) error {
		attempts++
		// Simulate slow operation (200ms)
		select {
		case <-time.After(200 * time.Millisecond):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	// Each attempt has 50ms timeout (operation needs 200ms)
	err := RetryWithTimeout(ctx, fn, 3, 50*time.Millisecond)
	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

// ============================================================================
// Tests for Exercise 2: FetchAll
// ============================================================================

func TestFetchAll_Success(t *testing.T) {
	// Create test servers
	servers := createTestServers([]string{"response1", "response2", "response3"})
	defer closeServers(servers)

	urls := serversToURLs(servers)
	ctx := context.Background()

	bodies, err := FetchAll(ctx, urls, 5*time.Second)
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	if len(bodies) != 3 {
		t.Fatalf("Expected 3 bodies, got %d", len(bodies))
	}

	// Check responses are in correct order
	expected := []string{"response1\n", "response2\n", "response3\n"}
	for i, body := range bodies {
		if body != expected[i] {
			t.Errorf("Body %d: expected %q, got %q", i, expected[i], body)
		}
	}
}

func TestFetchAll_Timeout(t *testing.T) {
	// Create slow server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		fmt.Fprintln(w, "slow response")
	}))
	defer server.Close()

	ctx := context.Background()
	urls := []string{server.URL}

	_, err := FetchAll(ctx, urls, 100*time.Millisecond)
	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Expected DeadlineExceeded, got %v", err)
	}
}

func TestFetchAll_OneFailure(t *testing.T) {
	// Create servers: one normal, one error
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "success")
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server2.Close()

	ctx := context.Background()
	urls := []string{server1.URL, server2.URL}

	_, err := FetchAll(ctx, urls, 5*time.Second)
	if err == nil {
		t.Fatal("Expected error from failed server, got nil")
	}
}

func TestFetchAll_ContextCancelled(t *testing.T) {
	// Create slow server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		fmt.Fprintln(w, "response")
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())

	// Cancel after 50ms
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	urls := []string{server.URL}
	_, err := FetchAll(ctx, urls, 5*time.Second)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

// ============================================================================
// Tests for Exercise 3: WorkerPool
// ============================================================================

func TestWorkerPool_Basic(t *testing.T) {
	ctx := context.Background()
	jobs := make(chan Job, 10)
	results := WorkerPool(ctx, 3, jobs)

	// Send 10 jobs
	for i := 0; i < 10; i++ {
		jobs <- Job{ID: i}
	}
	close(jobs)

	// Collect results
	resultCount := 0
	for range results {
		resultCount++
	}

	if resultCount != 10 {
		t.Errorf("Expected 10 results, got %d", resultCount)
	}
}

func TestWorkerPool_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	jobs := make(chan Job, 100)
	results := WorkerPool(ctx, 3, jobs)

	// Send many jobs
	go func() {
		for i := 0; i < 100; i++ {
			jobs <- Job{ID: i}
		}
		close(jobs)
	}()

	// Cancel context after short delay
	time.Sleep(50 * time.Millisecond)
	cancel()

	// Collect results (should stop early)
	resultCount := 0
	for range results {
		resultCount++
	}

	// Should process fewer than 100 jobs due to cancellation
	if resultCount >= 100 {
		t.Errorf("Expected fewer than 100 results due to cancellation, got %d", resultCount)
	}
}

func TestWorkerPool_NoJobs(t *testing.T) {
	ctx := context.Background()
	jobs := make(chan Job)
	results := WorkerPool(ctx, 3, jobs)

	// Close immediately (no jobs)
	close(jobs)

	// Should get no results
	resultCount := 0
	for range results {
		resultCount++
	}

	if resultCount != 0 {
		t.Errorf("Expected 0 results, got %d", resultCount)
	}
}

// ============================================================================
// Tests for Exercise 4: CacheWithTTL
// ============================================================================

func TestCache_SetGet(t *testing.T) {
	cache := NewCache()

	cache.Set("key1", "value1", 1*time.Second)
	val, ok := cache.Get("key1")
	if !ok {
		t.Fatal("Expected key1 to exist")
	}
	if val != "value1" {
		t.Errorf("Expected 'value1', got %v", val)
	}
}

func TestCache_Expiration(t *testing.T) {
	cache := NewCache()

	cache.Set("key1", "value1", 100*time.Millisecond)

	// Should exist immediately
	_, ok := cache.Get("key1")
	if !ok {
		t.Fatal("Expected key1 to exist")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should be expired
	_, ok = cache.Get("key1")
	if ok {
		t.Error("Expected key1 to be expired")
	}
}

func TestCache_Cleanup(t *testing.T) {
	cache := NewCache()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start cleanup goroutine
	go cache.Cleanup(ctx)

	// Add entries with short TTL
	for i := 0; i < 10; i++ {
		cache.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i), 50*time.Millisecond)
	}

	// Wait for cleanup
	time.Sleep(100 * time.Millisecond)

	// All entries should be cleaned up
	for i := 0; i < 10; i++ {
		if _, ok := cache.Get(fmt.Sprintf("key%d", i)); ok {
			t.Errorf("Expected key%d to be cleaned up", i)
		}
	}
}

func TestCache_CleanupStops(t *testing.T) {
	cache := NewCache()
	ctx, cancel := context.WithCancel(context.Background())

	// Start cleanup goroutine
	done := make(chan bool)
	go func() {
		cache.Cleanup(ctx)
		done <- true
	}()

	// Cancel context
	cancel()

	// Cleanup should stop
	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Error("Cleanup goroutine did not stop after context cancellation")
	}
}

// ============================================================================
// Tests for Exercise 5: RateLimiter
// ============================================================================

func TestRateLimiter_Basic(t *testing.T) {
	limiter := NewRateLimiter(10) // 10 ops/sec
	ctx := context.Background()

	// First 10 operations should be instant (bucket is pre-filled)
	start := time.Now()
	for i := 0; i < 10; i++ {
		if err := limiter.Wait(ctx); err != nil {
			t.Fatalf("Wait failed: %v", err)
		}
	}
	elapsed := time.Since(start)

	// First 10 operations should complete very quickly (tokens pre-filled)
	if elapsed > 100*time.Millisecond {
		t.Errorf("Expected fast completion, got %v", elapsed)
	}

	// Next 10 operations should take ~1 second (need to wait for refill)
	start = time.Now()
	for i := 0; i < 10; i++ {
		if err := limiter.Wait(ctx); err != nil {
			t.Fatalf("Wait failed: %v", err)
		}
	}
	elapsed = time.Since(start)

	// 10 operations at 10 ops/sec should take ~1 second
	if elapsed < 900*time.Millisecond || elapsed > 1100*time.Millisecond {
		t.Errorf("Expected ~1 second, got %v", elapsed)
	}
}

func TestRateLimiter_ContextCancellation(t *testing.T) {
	limiter := NewRateLimiter(1) // 1 op/sec
	ctx, cancel := context.WithCancel(context.Background())

	// First operation should succeed
	if err := limiter.Wait(ctx); err != nil {
		t.Fatalf("First Wait failed: %v", err)
	}

	// Cancel context
	cancel()

	// Second operation should fail with context error
	err := limiter.Wait(ctx)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestRateLimiter_Concurrent(t *testing.T) {
	limiter := NewRateLimiter(100) // 100 ops/sec
	ctx := context.Background()

	var count atomic.Int32

	// Run 100 concurrent operations
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := limiter.Wait(ctx); err != nil {
				t.Errorf("Wait failed: %v", err)
				return
			}
			count.Add(1)
		}()
	}

	wg.Wait()

	if count.Load() != 100 {
		t.Errorf("Expected 100 operations, got %d", count.Load())
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

func createTestServers(responses []string) []*httptest.Server {
	servers := make([]*httptest.Server, len(responses))
	for i, resp := range responses {
		response := resp // Capture loop variable
		servers[i] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, response)
		}))
	}
	return servers
}

func closeServers(servers []*httptest.Server) {
	for _, srv := range servers {
		srv.Close()
	}
}

func serversToURLs(servers []*httptest.Server) []string {
	urls := make([]string, len(servers))
	for i, srv := range servers {
		urls[i] = srv.URL
	}
	return urls
}
