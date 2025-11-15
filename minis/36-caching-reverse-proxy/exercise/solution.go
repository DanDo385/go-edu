//go:build solution
// +build solution

/*
Problem: Build a caching reverse proxy to reduce backend load and improve response times

Requirements:
1. Cache GET responses in memory
2. Respect Cache-Control headers
3. Implement LRU eviction
4. Support TTL expiration
5. Thread-safe for concurrent requests
6. Provide cache statistics

Algorithm:
- Check if request method is GET
- Generate cache key from URL
- Check cache for entry
- If hit and not expired: serve from cache
- If miss: forward to backend, cache response, serve to client
- Evict LRU entry if cache is full

Time Complexity:
- Get: O(1) average case
- Set: O(1) average case
- Eviction: O(1)

Space Complexity: O(n) where n is maxSize

Why Go is well-suited:
- net/http/httputil: Built-in reverse proxy
- sync.RWMutex: Efficient read-heavy locking
- container/list: Optimized doubly-linked list
- Goroutines: Natural concurrency model

Compared to other languages:
- Python: No built-in reverse proxy, GIL limits concurrency
- Node.js: Good for proxying, but weak typing for cache entries
- Rust: Excellent performance, but more complex ownership model
*/

package exercise

import (
	"bytes"
	"container/list"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// CacheEntry represents a cached HTTP response.
type CacheEntry struct {
	Body       []byte
	StatusCode int
	Header     http.Header
	Expiry     time.Time
	AccessTime time.Time
}

// lruEntry is used in the LRU linked list.
type lruEntry struct {
	key   string
	value *CacheEntry
}

// Cache is a thread-safe LRU cache with TTL support.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
	lru     *list.List
	lruMap  map[string]*list.Element
	maxSize int
	ttl     time.Duration
	hits    int64
	misses  int64
}

// NewCache creates a new cache.
func NewCache(maxSize int, ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]*CacheEntry),
		lru:     list.New(),
		lruMap:  make(map[string]*list.Element),
		maxSize: maxSize,
		ttl:     ttl,
	}
}

// Get retrieves an entry from the cache.
func (c *Cache) Get(key string) (*CacheEntry, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if exists
	elem, exists := c.lruMap[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	entry := elem.Value.(*lruEntry).value
	if time.Now().After(entry.Expiry) {
		// Expired, remove from cache
		c.lru.Remove(elem)
		delete(c.lruMap, key)
		delete(c.entries, key)
		return nil, false
	}

	// Mark as recently used
	c.lru.MoveToFront(elem)
	entry.AccessTime = time.Now()

	return entry, true
}

// Set adds or updates an entry in the cache.
func (c *Cache) Set(key string, entry *CacheEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if already exists
	if elem, exists := c.lruMap[key]; exists {
		// Update existing entry
		elem.Value.(*lruEntry).value = entry
		c.lru.MoveToFront(elem)
		return
	}

	// Check if cache is full
	if c.lru.Len() >= c.maxSize {
		// Evict LRU
		oldest := c.lru.Back()
		if oldest != nil {
			lruEntry := oldest.Value.(*lruEntry)
			c.lru.Remove(oldest)
			delete(c.lruMap, lruEntry.key)
			delete(c.entries, lruEntry.key)
		}
	}

	// Add new entry
	elem := c.lru.PushFront(&lruEntry{key: key, value: entry})
	c.lruMap[key] = elem
	c.entries[key] = entry
}

// Delete removes an entry from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, exists := c.lruMap[key]; exists {
		c.lru.Remove(elem)
		delete(c.lruMap, key)
		delete(c.entries, key)
	}
}

// Clear removes all entries from the cache.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]*CacheEntry)
	c.lru = list.New()
	c.lruMap = make(map[string]*list.Element)
	atomic.StoreInt64(&c.hits, 0)
	atomic.StoreInt64(&c.misses, 0)
}

// Stats returns cache statistics.
func (c *Cache) Stats() map[string]interface{} {
	c.mu.RLock()
	size := c.lru.Len()
	c.mu.RUnlock()

	hits := atomic.LoadInt64(&c.hits)
	misses := atomic.LoadInt64(&c.misses)
	total := hits + misses
	hitRate := 0.0
	if total > 0 {
		hitRate = float64(hits) / float64(total)
	}

	return map[string]interface{}{
		"hits":     hits,
		"misses":   misses,
		"total":    total,
		"size":     size,
		"max_size": c.maxSize,
		"hit_rate": hitRate,
	}
}

// ResponseRecorder captures an HTTP response for caching.
type ResponseRecorder struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
	header http.Header
}

// NewResponseRecorder creates a new response recorder.
func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{
		ResponseWriter: w,
		status:         http.StatusOK,
		body:           new(bytes.Buffer),
		header:         make(http.Header),
	}
}

// WriteHeader captures the status code.
func (rr *ResponseRecorder) WriteHeader(status int) {
	rr.status = status
}

// Write captures the response body.
func (rr *ResponseRecorder) Write(b []byte) (int, error) {
	return rr.body.Write(b)
}

// Header returns the captured headers.
func (rr *ResponseRecorder) Header() http.Header {
	return rr.header
}

// NewCachingProxy creates a caching reverse proxy.
func (c *Cache) NewCachingProxy(target *url.URL) http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(target)
	return c.Handler(proxy)
}

// Handler wraps an http.Handler with caching.
func (c *Cache) Handler(backend http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only cache GET requests
		if r.Method != http.MethodGet {
			backend.ServeHTTP(w, r)
			return
		}

		// Generate cache key
		key := r.URL.String()

		// Try cache
		if entry, exists := c.Get(key); exists {
			c.serveFromCache(w, entry)
			atomic.AddInt64(&c.hits, 1)
			return
		}

		// Cache miss
		atomic.AddInt64(&c.misses, 1)

		// Capture response
		recorder := NewResponseRecorder(w)

		// Forward to backend
		backend.ServeHTTP(recorder, r)

		// Cache if cacheable
		if c.isCacheable(recorder) {
			entry := &CacheEntry{
				Body:       recorder.body.Bytes(),
				StatusCode: recorder.status,
				Header:     recorder.header,
				Expiry:     c.calculateExpiry(recorder.header),
				AccessTime: time.Now(),
			}
			c.Set(key, entry)
		}

		// Send to client
		c.copyResponseToWriter(w, recorder)
	})
}

// serveFromCache serves a cached response.
func (c *Cache) serveFromCache(w http.ResponseWriter, entry *CacheEntry) {
	// Copy headers
	for key, values := range entry.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Add cache status header
	w.Header().Set("X-Cache", "HIT")

	// Write status and body
	w.WriteHeader(entry.StatusCode)
	w.Write(entry.Body)
}

// copyResponseToWriter copies recorded response to actual writer.
func (c *Cache) copyResponseToWriter(w http.ResponseWriter, recorder *ResponseRecorder) {
	// Copy headers
	for key, values := range recorder.header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Add cache status header
	w.Header().Set("X-Cache", "MISS")

	// Write status and body
	w.WriteHeader(recorder.status)
	w.Write(recorder.body.Bytes())
}

// isCacheable determines if a response should be cached.
func (c *Cache) isCacheable(recorder *ResponseRecorder) bool {
	// Only cache successful responses
	if recorder.status != http.StatusOK {
		return false
	}

	// Check Cache-Control
	cacheControl := recorder.header.Get("Cache-Control")

	// no-store means never cache
	if strings.Contains(cacheControl, "no-store") {
		return false
	}

	// private means only browser can cache, not proxy
	if strings.Contains(cacheControl, "private") {
		return false
	}

	return true
}

// calculateExpiry calculates when a cache entry should expire.
func (c *Cache) calculateExpiry(header http.Header) time.Time {
	// For simplicity, use default TTL
	// In production, parse Cache-Control max-age
	return time.Now().Add(c.ttl)
}

// StatsHandler returns an HTTP handler for cache statistics.
func (c *Cache) StatsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		stats := c.Stats()

		fmt.Fprintf(w, "{\n")
		fmt.Fprintf(w, "  \"hits\": %d,\n", stats["hits"])
		fmt.Fprintf(w, "  \"misses\": %d,\n", stats["misses"])
		fmt.Fprintf(w, "  \"total\": %d,\n", stats["total"])
		fmt.Fprintf(w, "  \"hit_rate\": %.2f,\n", stats["hit_rate"])
		fmt.Fprintf(w, "  \"size\": %d,\n", stats["size"])
		fmt.Fprintf(w, "  \"max_size\": %d\n", stats["max_size"])
		fmt.Fprintf(w, "}\n")
	}
}

// ClearHandler returns an HTTP handler to clear the cache.
func (c *Cache) ClearHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.Clear()
		w.WriteHeader(http.StatusNoContent)
	}
}

/*
Alternatives and Trade-offs:

1. Write-through vs Write-aside caching:
   - We use write-aside (lazy loading)
   - Write-through would update cache on every write
   - Trade-off: Simplicity vs consistency

2. Distributed caching (Redis):
   - Pros: Shared cache across instances, persistence
   - Cons: Network latency, complexity
   - When to use: Multi-instance deployments

3. Cache stampede protection:
   - Problem: Multiple requests for expired key hit backend
   - Solution: Request coalescing (singleflight pattern)
   - Trade-off: Complexity vs thundering herd prevention

4. Compression:
   - Could compress cached entries to save memory
   - Trade-off: Memory vs CPU

Go vs Other Languages:

Python:
- No built-in reverse proxy
- GIL limits concurrent caching
- Dict-based cache similar but slower

Node.js:
- Good proxy libraries (http-proxy)
- Single-threaded limits CPU-bound caching
- Weak typing for cache entries

Rust:
- Excellent performance
- Ownership model adds complexity
- No built-in proxy like httputil

Java:
- Good caching libraries (Caffeine)
- More verbose, heavier runtime
- Similar performance to Go
*/
