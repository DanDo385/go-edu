package exercise

import (
	"sync"
	"testing"
	"time"
)

// Test Exercise 1: Thread-Safe Counter
func TestCounter(t *testing.T) {
	counter := NewCounter()

	// Test single increment
	counter.Increment()
	if got := counter.Value(); got != 1 {
		t.Errorf("After 1 increment, got %d, want 1", got)
	}

	// Test multiple increments
	counter.Increment()
	counter.Increment()
	if got := counter.Value(); got != 3 {
		t.Errorf("After 3 increments, got %d, want 3", got)
	}

	// Test decrement
	counter.Decrement()
	if got := counter.Value(); got != 2 {
		t.Errorf("After 1 decrement, got %d, want 2", got)
	}

	// Test reset
	counter.Reset()
	if got := counter.Value(); got != 0 {
		t.Errorf("After reset, got %d, want 0", got)
	}
}

func TestCounterConcurrent(t *testing.T) {
	counter := NewCounter()
	var wg sync.WaitGroup

	// 100 goroutines, each incrementing 100 times
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				counter.Increment()
			}
		}()
	}

	wg.Wait()

	expected := 10000
	if got := counter.Value(); got != expected {
		t.Errorf("Concurrent increments: got %d, want %d", got, expected)
	}
}

func TestCounterConcurrentMixed(t *testing.T) {
	counter := NewCounter()
	var wg sync.WaitGroup

	// 50 goroutines incrementing
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				counter.Increment()
			}
		}()
	}

	// 25 goroutines decrementing
	for i := 0; i < 25; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				counter.Decrement()
			}
		}()
	}

	wg.Wait()

	expected := 2500 // 5000 increments - 2500 decrements
	if got := counter.Value(); got != expected {
		t.Errorf("Mixed concurrent ops: got %d, want %d", got, expected)
	}
}

// Test Exercise 2: Thread-Safe Cache
func TestCache(t *testing.T) {
	cache := NewCache[string, int]()

	// Test Set and Get
	cache.Set("age", 25)
	value, ok := cache.Get("age")
	if !ok || value != 25 {
		t.Errorf("Get after Set: got (%d, %v), want (25, true)", value, ok)
	}

	// Test Get non-existent key
	value, ok = cache.Get("nonexistent")
	if ok {
		t.Errorf("Get non-existent key: got (%d, %v), want (0, false)", value, ok)
	}

	// Test Len
	cache.Set("name", 100)
	cache.Set("city", 200)
	if got := cache.Len(); got != 3 {
		t.Errorf("Len: got %d, want 3", got)
	}

	// Test Delete
	cache.Delete("age")
	value, ok = cache.Get("age")
	if ok {
		t.Errorf("Get after Delete: got (%d, %v), want (0, false)", value, ok)
	}

	if got := cache.Len(); got != 2 {
		t.Errorf("Len after Delete: got %d, want 2", got)
	}

	// Test Clear
	cache.Clear()
	if got := cache.Len(); got != 0 {
		t.Errorf("Len after Clear: got %d, want 0", got)
	}
}

func TestCacheConcurrent(t *testing.T) {
	cache := NewCache[int, int]()
	var wg sync.WaitGroup

	// Pre-populate
	for i := 0; i < 100; i++ {
		cache.Set(i, i*10)
	}

	// 50 readers
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				// Just read, don't check value since writers may update
				cache.Get(j % 100)
			}
		}()
	}

	// 10 writers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				cache.Set(j, j*20)
			}
		}()
	}

	wg.Wait()

	// Verify cache still works and has correct final values
	if got := cache.Len(); got != 100 {
		t.Errorf("Final Len: got %d, want 100", got)
	}

	// Check that final values are from writers
	for i := 0; i < 100; i++ {
		value, ok := cache.Get(i)
		if !ok {
			t.Errorf("Key %d missing after concurrent operations", i)
		}
		if value != i*20 {
			t.Errorf("Final value for key %d: got %d, want %d", i, value, i*20)
		}
	}
}

// Test Exercise 3: Expiring Cache
func TestExpiringCache(t *testing.T) {
	cache := NewExpiringCache[string, string]()

	// Test Set and Get before expiration
	cache.Set("key1", "value1", 100*time.Millisecond)
	value, ok := cache.Get("key1")
	if !ok || value != "value1" {
		t.Errorf("Get before expiration: got (%s, %v), want (value1, true)", value, ok)
	}

	// Test Get after expiration
	time.Sleep(150 * time.Millisecond)
	value, ok = cache.Get("key1")
	if ok {
		t.Errorf("Get after expiration: got (%s, %v), want (\"\", false)", value, ok)
	}
}

func TestExpiringCacheCleanup(t *testing.T) {
	cache := NewExpiringCache[string, int]()
	cache.StartCleanup(50 * time.Millisecond)
	defer cache.StopCleanup()

	// Add entries with short TTL
	for i := 0; i < 10; i++ {
		cache.Set(string(rune('a'+i)), i, 100*time.Millisecond)
	}

	// Wait for cleanup to run
	time.Sleep(200 * time.Millisecond)

	// Verify entries were cleaned up
	for i := 0; i < 10; i++ {
		if _, ok := cache.Get(string(rune('a' + i))); ok {
			t.Errorf("Entry %d should have been cleaned up", i)
		}
	}
}

// Test Exercise 4: Sharded Map
func TestShardedMap(t *testing.T) {
	sm := NewShardedMap[string, int]()

	// Test Set and Get
	sm.Set("key1", 100)
	value, ok := sm.Get("key1")
	if !ok || value != 100 {
		t.Errorf("Get after Set: got (%d, %v), want (100, true)", value, ok)
	}

	// Test multiple keys
	for i := 0; i < 100; i++ {
		key := string(rune('a' + (i % 26)))
		sm.Set(key, i)
	}

	// Verify some keys
	value, ok = sm.Get("a")
	if !ok {
		t.Errorf("Get key 'a': got ok=%v, want true", ok)
	}

	// Test Delete
	sm.Delete("a")
	_, ok = sm.Get("a")
	if ok {
		t.Errorf("Get after Delete: got ok=%v, want false", ok)
	}
}

func TestShardedMapConcurrent(t *testing.T) {
	sm := NewShardedMap[int, int]()
	var wg sync.WaitGroup

	// 100 goroutines writing
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				key := id*100 + j
				sm.Set(key, key*10)
			}
		}(i)
	}

	// 50 goroutines reading
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				key := j % 10000
				sm.Get(key)
			}
		}(i)
	}

	wg.Wait()

	// Verify some values
	value, ok := sm.Get(0)
	if !ok || value != 0 {
		t.Errorf("Get(0): got (%d, %v), want (0, true)", value, ok)
	}

	value, ok = sm.Get(9999)
	if !ok || value != 99990 {
		t.Errorf("Get(9999): got (%d, %v), want (99990, true)", value, ok)
	}
}

// Test Exercise 5: Metrics Collector
func TestMetrics(t *testing.T) {
	metrics := NewMetrics()

	// Test counter
	metrics.IncrementCounter("requests")
	metrics.IncrementCounter("requests")
	if got := metrics.GetCounter("requests"); got != 2 {
		t.Errorf("GetCounter: got %d, want 2", got)
	}

	// Test gauge
	metrics.SetGauge("connections", 42)
	if got := metrics.GetGauge("connections"); got != 42 {
		t.Errorf("GetGauge: got %d, want 42", got)
	}

	// Test snapshot
	metrics.IncrementCounter("errors")
	metrics.SetGauge("memory", 1024)

	snapshot := metrics.Snapshot()
	if len(snapshot) != 4 {
		t.Errorf("Snapshot length: got %d, want 4", len(snapshot))
	}

	if snapshot["requests"] != 2 {
		t.Errorf("Snapshot requests: got %d, want 2", snapshot["requests"])
	}
}

func TestMetricsConcurrent(t *testing.T) {
	metrics := NewMetrics()
	var wg sync.WaitGroup

	// 100 goroutines incrementing counters
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				metrics.IncrementCounter("requests")
			}
		}()
	}

	wg.Wait()

	if got := metrics.GetCounter("requests"); got != 10000 {
		t.Errorf("Concurrent increments: got %d, want 10000", got)
	}
}

// Test Exercise 6: Rate Limiter
func TestRateLimiter(t *testing.T) {
	limiter := NewRateLimiter(10, 10)

	// Should allow up to burst
	allowed := 0
	for i := 0; i < 15; i++ {
		if limiter.Allow() {
			allowed++
		}
	}

	if allowed != 10 {
		t.Errorf("Initial burst: allowed %d, want 10", allowed)
	}

	// Wait for tokens to refill
	time.Sleep(1 * time.Second)

	// Should allow more after refill
	if !limiter.Allow() {
		t.Error("Should allow after refill")
	}
}

func TestRateLimiterConcurrent(t *testing.T) {
	limiter := NewRateLimiter(100, 50)
	var wg sync.WaitGroup
	var allowed int64

	// 10 goroutines trying to get tokens
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				if limiter.Allow() {
					allowed++
				}
				time.Sleep(10 * time.Millisecond)
			}
		}()
	}

	wg.Wait()

	// Should allow some operations (exact count depends on timing)
	if allowed == 0 {
		t.Error("Should allow some operations")
	}
}

// Benchmark Cache operations
func BenchmarkCacheGet(b *testing.B) {
	cache := NewCache[int, int]()
	for i := 0; i < 1000; i++ {
		cache.Set(i, i*10)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Get(i % 1000)
			i++
		}
	})
}

func BenchmarkCacheSet(b *testing.B) {
	cache := NewCache[int, int]()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Set(i%1000, i)
			i++
		}
	})
}

func BenchmarkShardedMapGet(b *testing.B) {
	sm := NewShardedMap[int, int]()
	for i := 0; i < 1000; i++ {
		sm.Set(i, i*10)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			sm.Get(i % 1000)
			i++
		}
	})
}

func BenchmarkShardedMapSet(b *testing.B) {
	sm := NewShardedMap[int, int]()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			sm.Set(i%1000, i)
			i++
		}
	})
}

func BenchmarkCounter(b *testing.B) {
	counter := NewCounter()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Increment()
		}
	})
}

func BenchmarkMetricsCounter(b *testing.B) {
	metrics := NewMetrics()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			metrics.IncrementCounter("requests")
		}
	})
}
