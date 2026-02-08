//go:build reference

package syncmutexvsrwmutex

/*
Reference Solution - Mutex vs RWMutex, Caches, Sharding
======================================================

sync.Mutex: exclusive lock for read-write access (Counter, Set/Delete).
sync.RWMutex: RLock for concurrent reads, Lock for writes (Cache.Get vs Set).
Use RWMutex when reads dominate; Mutex when reads/writes are similar.

ExpiringCache: RLock for Get, then Lock only when deleting expired entry.
ShardedMap: hash key to shard, each shard has own RWMutex (reduces contention).
Metrics: RLock for GetCounter/GetGauge/Snapshot; Lock for Increment/SetGauge.
*/

import (
	"hash/fnv"
	"sync"
	"time"
)

// Exercise 1: Thread-Safe Counter
type Counter struct {
	mu    sync.Mutex
	value int
}

// NewCounter implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func NewCounter() *Counter {
	return &Counter{}
}

// Increment implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

// Decrement implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c *Counter) Decrement() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value--
}

// Value implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// Reset implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = 0
}

// Exercise 2: Thread-Safe Cache with RWMutex
type Cache[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

// NewCache implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		data: make(map[K]V),
	}
}

// Get implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.data[key]
	return value, ok
}

// Set implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

// Delete implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// Len implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c *Cache[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.data)
}

// Clear implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c *Cache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[K]V)
}

// Exercise 3: Cache with Expiration
type ExpiringCache[K comparable, V any] struct {
	mu      sync.RWMutex
	data    map[K]*cacheEntry[V]
	stopCh  chan struct{}
	stopped bool
}

type cacheEntry[V any] struct {
	value      V
	expiration time.Time
}

// NewExpiringCache implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func NewExpiringCache[K comparable, V any]() *ExpiringCache[K, V] {
	return &ExpiringCache[K, V]{
		data:   make(map[K]*cacheEntry[V]),
		stopCh: make(chan struct{}),
	}
}

// Set implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c *ExpiringCache[K, V]) Set(key K, value V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = &cacheEntry[V]{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

// Get implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c *ExpiringCache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	entry, ok := c.data[key]
	c.mu.RUnlock()

	var zero V
	if !ok {
		return zero, false
	}

	// Check expiration
	if time.Now().After(entry.expiration) {
		// Entry expired, remove it
		c.mu.Lock()
		delete(c.data, key)
		c.mu.Unlock()
		return zero, false
	}

	return entry.value, true
}

// StartCleanup implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c *ExpiringCache[K, V]) StartCleanup(interval time.Duration) {
	c.mu.Lock()
	if !c.stopped {
		c.mu.Unlock()
		return // Already running
	}
	c.stopped = false
	c.stopCh = make(chan struct{})
	c.mu.Unlock()

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.cleanup()
			case <-c.stopCh:
				return
			}
		}
	}()
}

// StopCleanup implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c *ExpiringCache[K, V]) StopCleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.stopped {
		c.stopped = true
		close(c.stopCh)
	}
}

// cleanup implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c *ExpiringCache[K, V]) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.data {
		if now.After(entry.expiration) {
			delete(c.data, key)
		}
	}
}

// Exercise 4: Sharded Map
type ShardedMap[K comparable, V any] struct {
	shards [numShards]*shard[K, V]
}

type shard[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

const numShards = 16

// NewShardedMap implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func NewShardedMap[K comparable, V any]() *ShardedMap[K, V] {
	sm := &ShardedMap[K, V]{}
	for i := 0; i < numShards; i++ {
		sm.shards[i] = &shard[K, V]{
			data: make(map[K]V),
		}
	}
	return sm
}

// getShard implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (sm *ShardedMap[K, V]) getShard(key K) *shard[K, V] {
	// Hash the key
	h := fnv.New32a()

	// Convert key to bytes for hashing
	// This is a simple approach; for production use a better hash function
	var keyBytes []byte
	switch k := any(key).(type) {
	case string:
		keyBytes = []byte(k)
	case int:
		keyBytes = []byte{byte(k), byte(k >> 8), byte(k >> 16), byte(k >> 24)}
	case int64:
		keyBytes = []byte{byte(k), byte(k >> 8), byte(k >> 16), byte(k >> 24),
			byte(k >> 32), byte(k >> 40), byte(k >> 48), byte(k >> 56)}
	default:
		// Fallback: use string representation
		keyBytes = []byte(any(key).(string))
	}

	h.Write(keyBytes)
	shardIndex := h.Sum32() % numShards
	return sm.shards[shardIndex]
}

// Get implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (sm *ShardedMap[K, V]) Get(key K) (V, bool) {
	shard := sm.getShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()
	value, ok := shard.data[key]
	return value, ok
}

// Set implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (sm *ShardedMap[K, V]) Set(key K, value V) {
	shard := sm.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	shard.data[key] = value
}

// Delete implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (sm *ShardedMap[K, V]) Delete(key K) {
	shard := sm.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	delete(shard.data, key)
}

// Exercise 5: Metrics Collector
type Metrics struct {
	mu      sync.RWMutex
	metrics map[string]int64
}

// NewMetrics implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func NewMetrics() *Metrics {
	return &Metrics{
		metrics: make(map[string]int64),
	}
}

// IncrementCounter implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (m *Metrics) IncrementCounter(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metrics[name]++
}

// SetGauge implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (m *Metrics) SetGauge(name string, value int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metrics[name] = value
}

// GetCounter implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (m *Metrics) GetCounter(name string) int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.metrics[name]
}

// GetGauge implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (m *Metrics) GetGauge(name string) int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.metrics[name]
}

// Snapshot implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (m *Metrics) Snapshot() map[string]int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	snapshot := make(map[string]int64, len(m.metrics))
	for name, value := range m.metrics {
		snapshot[name] = value
	}
	return snapshot
}

// Exercise 6: Rate Limiter
type RateLimiter struct {
	mu           sync.Mutex
	rate         float64
	burst        int
	tokens       float64
	lastRefill   time.Time
	stopCh       chan struct{}
	refillTicker *time.Ticker
}

// NewRateLimiter implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func NewRateLimiter(rate float64, burst int) *RateLimiter {
	rl := &RateLimiter{
		rate:       rate,
		burst:      burst,
		tokens:     float64(burst),
		lastRefill: time.Now(),
		stopCh:     make(chan struct{}),
	}

	// Start refill goroutine
	go rl.refill()

	return rl
}

// Allow implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}
	return false
}

// refill implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (rl *RateLimiter) refill() {
	ticker := time.NewTicker(time.Second / 10) // Refill 10 times per second
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()

			now := time.Now()
			elapsed := now.Sub(rl.lastRefill).Seconds()
			rl.lastRefill = now

			// Add tokens based on elapsed time
			rl.tokens += elapsed * rl.rate

			// Cap at burst size
			if rl.tokens > float64(rl.burst) {
				rl.tokens = float64(rl.burst)
			}

			rl.mu.Unlock()

		case <-rl.stopCh:
			return
		}
	}
}
