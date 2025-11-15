package exercise

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCache_GetSet(t *testing.T) {
	cache := NewCache(10, 1*time.Minute)

	// Test miss
	_, exists := cache.Get("key1")
	if exists {
		t.Error("Expected cache miss, got hit")
	}

	// Test set and hit
	entry := &CacheEntry{
		Body:       []byte("test data"),
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Expiry:     time.Now().Add(1 * time.Minute),
	}
	cache.Set("key1", entry)

	retrieved, exists := cache.Get("key1")
	if !exists {
		t.Fatal("Expected cache hit, got miss")
	}
	if string(retrieved.Body) != "test data" {
		t.Errorf("Expected 'test data', got %q", string(retrieved.Body))
	}
}

func TestCache_Expiry(t *testing.T) {
	cache := NewCache(10, 100*time.Millisecond)

	entry := &CacheEntry{
		Body:       []byte("test"),
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Expiry:     time.Now().Add(100 * time.Millisecond),
	}
	cache.Set("key1", entry)

	// Should exist initially
	_, exists := cache.Get("key1")
	if !exists {
		t.Error("Expected cache hit before expiry")
	}

	// Wait for expiry
	time.Sleep(150 * time.Millisecond)

	// Should be expired now
	_, exists = cache.Get("key1")
	if exists {
		t.Error("Expected cache miss after expiry")
	}
}

func TestCache_LRU_Eviction(t *testing.T) {
	cache := NewCache(3, 1*time.Minute) // Max 3 entries

	// Add 3 entries
	for i := 0; i < 3; i++ {
		entry := &CacheEntry{
			Body:       []byte(fmt.Sprintf("data%d", i)),
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Expiry:     time.Now().Add(1 * time.Minute),
		}
		cache.Set(fmt.Sprintf("key%d", i), entry)
	}

	// All should exist
	for i := 0; i < 3; i++ {
		if _, exists := cache.Get(fmt.Sprintf("key%d", i)); !exists {
			t.Errorf("Expected key%d to exist", i)
		}
	}

	// Add 4th entry, should evict key0 (LRU)
	entry := &CacheEntry{
		Body:       []byte("data3"),
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Expiry:     time.Now().Add(1 * time.Minute),
	}
	cache.Set("key3", entry)

	// key0 should be evicted
	if _, exists := cache.Get("key0"); exists {
		t.Error("Expected key0 to be evicted")
	}

	// Others should exist
	for i := 1; i <= 3; i++ {
		if _, exists := cache.Get(fmt.Sprintf("key%d", i)); !exists {
			t.Errorf("Expected key%d to exist", i)
		}
	}
}

func TestCache_LRU_AccessOrder(t *testing.T) {
	cache := NewCache(3, 1*time.Minute)

	// Add 3 entries
	for i := 0; i < 3; i++ {
		entry := &CacheEntry{
			Body:       []byte(fmt.Sprintf("data%d", i)),
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Expiry:     time.Now().Add(1 * time.Minute),
		}
		cache.Set(fmt.Sprintf("key%d", i), entry)
	}

	// Access key0 (makes it most recently used)
	cache.Get("key0")

	// Add 4th entry, should evict key1 (now LRU, not key0)
	entry := &CacheEntry{
		Body:       []byte("data3"),
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Expiry:     time.Now().Add(1 * time.Minute),
	}
	cache.Set("key3", entry)

	// key1 should be evicted (was LRU)
	if _, exists := cache.Get("key1"); exists {
		t.Error("Expected key1 to be evicted")
	}

	// key0 should still exist (was accessed)
	if _, exists := cache.Get("key0"); !exists {
		t.Error("Expected key0 to still exist")
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	cache := NewCache(100, 1*time.Minute)

	var wg sync.WaitGroup
	numGoroutines := 10
	numOps := 100

	// Concurrent sets
	wg.Add(numGoroutines)
	for g := 0; g < numGoroutines; g++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < numOps; i++ {
				key := fmt.Sprintf("key%d-%d", id, i)
				entry := &CacheEntry{
					Body:       []byte(fmt.Sprintf("data%d", i)),
					StatusCode: http.StatusOK,
					Header:     make(http.Header),
					Expiry:     time.Now().Add(1 * time.Minute),
				}
				cache.Set(key, entry)
			}
		}(g)
	}
	wg.Wait()

	// Concurrent gets
	wg.Add(numGoroutines)
	for g := 0; g < numGoroutines; g++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < numOps; i++ {
				key := fmt.Sprintf("key%d-%d", id, i)
				cache.Get(key)
			}
		}(g)
	}
	wg.Wait()
}

func TestCachingProxy_CacheHit(t *testing.T) {
	var backendCalls int32

	// Backend server
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&backendCalls, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"success"}`))
	}))
	defer backend.Close()

	// Create cache and proxy
	cache := NewCache(10, 1*time.Minute)
	proxy := httptest.NewServer(cache.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Forward to backend
		req, _ := http.NewRequest(r.Method, backend.URL+r.URL.Path, nil)
		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()

		// Copy response
		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		w.WriteHeader(resp.StatusCode)
		body := make([]byte, 1024)
		n, _ := resp.Body.Read(body)
		w.Write(body[:n])
	})))
	defer proxy.Close()

	// First request - cache miss
	resp1, _ := http.Get(proxy.URL + "/test")
	resp1.Body.Close()

	if atomic.LoadInt32(&backendCalls) != 1 {
		t.Errorf("Expected 1 backend call, got %d", backendCalls)
	}

	// Second request - should be cache hit
	resp2, _ := http.Get(proxy.URL + "/test")
	resp2.Body.Close()

	if atomic.LoadInt32(&backendCalls) != 1 {
		t.Errorf("Expected still 1 backend call (cache hit), got %d", backendCalls)
	}

	// Check cache header
	if resp2.Header.Get("X-Cache") != "HIT" {
		t.Errorf("Expected X-Cache: HIT, got %q", resp2.Header.Get("X-Cache"))
	}
}

func TestCachingProxy_OnlyGET(t *testing.T) {
	var backendCalls int32

	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&backendCalls, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	cache := NewCache(10, 1*time.Minute)
	proxy := httptest.NewServer(cache.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, _ := http.NewRequest(r.Method, backend.URL, nil)
		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()
		w.WriteHeader(resp.StatusCode)
	})))
	defer proxy.Close()

	// POST request should not be cached
	http.Post(proxy.URL, "application/json", nil)
	http.Post(proxy.URL, "application/json", nil)

	if atomic.LoadInt32(&backendCalls) != 2 {
		t.Errorf("POST should not be cached, expected 2 backend calls, got %d", backendCalls)
	}
}

func TestCache_Stats(t *testing.T) {
	cache := NewCache(10, 1*time.Minute)

	entry := &CacheEntry{
		Body:       []byte("test"),
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Expiry:     time.Now().Add(1 * time.Minute),
	}
	cache.Set("key1", entry)

	// Generate some hits and misses
	cache.Get("key1")      // hit
	cache.Get("key1")      // hit
	cache.Get("nonexist")  // miss
	atomic.AddInt64(&cache.hits, 2)
	atomic.AddInt64(&cache.misses, 1)

	stats := cache.Stats()

	if stats["hits"].(int64) != 2 {
		t.Errorf("Expected 2 hits, got %d", stats["hits"])
	}
	if stats["misses"].(int64) != 1 {
		t.Errorf("Expected 1 miss, got %d", stats["misses"])
	}
	if stats["size"].(int) != 1 {
		t.Errorf("Expected size 1, got %d", stats["size"])
	}
}

func TestCache_Clear(t *testing.T) {
	cache := NewCache(10, 1*time.Minute)

	// Add entries
	for i := 0; i < 5; i++ {
		entry := &CacheEntry{
			Body:       []byte(fmt.Sprintf("data%d", i)),
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Expiry:     time.Now().Add(1 * time.Minute),
		}
		cache.Set(fmt.Sprintf("key%d", i), entry)
	}

	// Clear cache
	cache.Clear()

	// All entries should be gone
	for i := 0; i < 5; i++ {
		if _, exists := cache.Get(fmt.Sprintf("key%d", i)); exists {
			t.Errorf("Expected key%d to be cleared", i)
		}
	}

	stats := cache.Stats()
	if stats["size"].(int) != 0 {
		t.Errorf("Expected size 0 after clear, got %d", stats["size"])
	}
}

func TestResponseRecorder(t *testing.T) {
	w := httptest.NewRecorder()
	recorder := NewResponseRecorder(w)

	recorder.Header().Set("X-Test", "value")
	recorder.WriteHeader(http.StatusCreated)
	recorder.Write([]byte("test body"))

	if recorder.status != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", recorder.status)
	}
	if recorder.body.String() != "test body" {
		t.Errorf("Expected body 'test body', got %q", recorder.body.String())
	}
	if recorder.header.Get("X-Test") != "value" {
		t.Errorf("Expected header X-Test: value, got %q", recorder.header.Get("X-Test"))
	}
}

func BenchmarkCache_Get(b *testing.B) {
	cache := NewCache(1000, 1*time.Minute)

	entry := &CacheEntry{
		Body:       []byte("test data"),
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Expiry:     time.Now().Add(1 * time.Minute),
	}
	cache.Set("key1", entry)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("key1")
	}
}

func BenchmarkCache_Set(b *testing.B) {
	cache := NewCache(1000, 1*time.Minute)

	entry := &CacheEntry{
		Body:       []byte("test data"),
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Expiry:     time.Now().Add(1 * time.Minute),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(fmt.Sprintf("key%d", i), entry)
	}
}

func BenchmarkCache_Concurrent(b *testing.B) {
	cache := NewCache(1000, 1*time.Minute)

	entry := &CacheEntry{
		Body:       []byte("test data"),
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Expiry:     time.Now().Add(1 * time.Minute),
	}

	// Pre-populate
	for i := 0; i < 100; i++ {
		cache.Set(fmt.Sprintf("key%d", i), entry)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key%d", i%100)
			if i%2 == 0 {
				cache.Get(key)
			} else {
				cache.Set(key, entry)
			}
			i++
		}
	})
}
