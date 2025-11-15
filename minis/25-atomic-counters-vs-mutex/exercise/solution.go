//go:build solution
// +build solution

// Package exercise contains hands-on exercises for atomic operations.

package exercise

import (
	"math"
	"sync/atomic"
	"time"
)

// ============================================================================
// SOLUTION 1: Atomic Counter
// ============================================================================

func NewAtomicCounter() *AtomicCounter {
	return &AtomicCounter{value: 0}
}

func (c *AtomicCounter) Increment() {
	atomic.AddInt64(&c.value, 1)
}

func (c *AtomicCounter) Decrement() {
	atomic.AddInt64(&c.value, -1)
}

func (c *AtomicCounter) Add(delta int64) {
	atomic.AddInt64(&c.value, delta)
}

func (c *AtomicCounter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

func (c *AtomicCounter) Reset() int64 {
	return atomic.SwapInt64(&c.value, 0)
}

// ============================================================================
// SOLUTION 2: Atomic Flag
// ============================================================================

func NewAtomicFlag() *AtomicFlag {
	return &AtomicFlag{value: 0}
}

func (f *AtomicFlag) Set() {
	atomic.StoreInt64(&f.value, 1)
}

func (f *AtomicFlag) Clear() {
	atomic.StoreInt64(&f.value, 0)
}

func (f *AtomicFlag) IsSet() bool {
	return atomic.LoadInt64(&f.value) == 1
}

func (f *AtomicFlag) TestAndSet() bool {
	old := atomic.SwapInt64(&f.value, 1)
	return old == 1
}

// ============================================================================
// SOLUTION 3: Rate Limiter
// ============================================================================

func NewRateLimiter(capacity, tokensPerSecond int64) *RateLimiter {
	return &RateLimiter{
		tokens:     capacity,
		maxTokens:  capacity,
		refillRate: tokensPerSecond,
		lastRefill: time.Now().Unix(),
	}
}

func (rl *RateLimiter) Allow() bool {
	now := time.Now().Unix()
	last := atomic.LoadInt64(&rl.lastRefill)

	// Try to refill tokens based on elapsed time
	if now > last {
		if atomic.CompareAndSwapInt64(&rl.lastRefill, last, now) {
			// We won the CAS, calculate new tokens
			elapsed := now - last
			newTokens := elapsed * rl.refillRate

			// Add tokens, but don't exceed max
			for {
				current := atomic.LoadInt64(&rl.tokens)
				updated := current + newTokens
				if updated > rl.maxTokens {
					updated = rl.maxTokens
				}
				if atomic.CompareAndSwapInt64(&rl.tokens, current, updated) {
					break
				}
			}
		}
	}

	// Try to consume a token
	for {
		tokens := atomic.LoadInt64(&rl.tokens)
		if tokens <= 0 {
			return false
		}
		if atomic.CompareAndSwapInt64(&rl.tokens, tokens, tokens-1) {
			return true
		}
	}
}

// ============================================================================
// SOLUTION 4: Atomic Max Tracker
// ============================================================================

func NewAtomicMax() *AtomicMax {
	return &AtomicMax{max: math.MinInt64}
}

func (am *AtomicMax) Update(value int64) {
	for {
		current := atomic.LoadInt64(&am.max)
		if value <= current {
			return // No update needed
		}
		if atomic.CompareAndSwapInt64(&am.max, current, value) {
			return // Successfully updated
		}
		// CAS failed, retry
	}
}

func (am *AtomicMax) Max() int64 {
	return atomic.LoadInt64(&am.max)
}

// ============================================================================
// SOLUTION 5: SpinLock
// ============================================================================

func NewSpinLock() *SpinLock {
	return &SpinLock{state: 0}
}

func (sl *SpinLock) Lock() {
	for atomic.SwapInt64(&sl.state, 1) != 0 {
		// Spin (busy-wait)
		// In production, you'd use runtime.Gosched() here
	}
}

func (sl *SpinLock) Unlock() {
	atomic.StoreInt64(&sl.state, 0)
}

// ============================================================================
// SOLUTION 6: Atomic State Machine
// ============================================================================

func NewAtomicState() *AtomicState {
	return &AtomicState{current: StateIdle}
}

func (as *AtomicState) CurrentState() int64 {
	return atomic.LoadInt64(&as.current)
}

func (as *AtomicState) Transition(expectedCurrent, newState int64) bool {
	return atomic.CompareAndSwapInt64(&as.current, expectedCurrent, newState)
}

// ============================================================================
// SOLUTION 7: Reference Counter
// ============================================================================

func NewReferenceCounter() *ReferenceCounter {
	return &ReferenceCounter{count: 0}
}

func (rc *ReferenceCounter) Acquire() {
	atomic.AddInt64(&rc.count, 1)
}

func (rc *ReferenceCounter) Release() bool {
	newCount := atomic.AddInt64(&rc.count, -1)
	return newCount == 0
}

func (rc *ReferenceCounter) Count() int64 {
	return atomic.LoadInt64(&rc.count)
}

// ============================================================================
// SOLUTION 8: Config Manager
// ============================================================================

func NewConfigManager() *ConfigManager {
	cm := &ConfigManager{}
	cm.config.Store(&Config{
		Timeout:    30,
		MaxRetries: 3,
		BatchSize:  100,
	})
	return cm
}

func (cm *ConfigManager) Update(newConfig *Config) {
	cm.config.Store(newConfig)
}

func (cm *ConfigManager) Get() *Config {
	return cm.config.Load().(*Config)
}

// ============================================================================
// SOLUTION 9: Load Balancer
// ============================================================================

func NewLoadBalancer(numWorkers int) *LoadBalancer {
	return &LoadBalancer{
		counter: 0,
		workers: int64(numWorkers),
	}
}

func (lb *LoadBalancer) NextWorker() int64 {
	val := atomic.AddInt64(&lb.counter, 1)
	return (val - 1) % lb.workers
}

// ============================================================================
// SOLUTION 10: Atomic Bitmap
// ============================================================================

func NewAtomicBitmap() *AtomicBitmap {
	return &AtomicBitmap{}
}

func (ab *AtomicBitmap) SetBit(bitIndex int) {
	if bitIndex < 0 || bitIndex >= 256 {
		return
	}

	word := bitIndex / 64
	bit := uint(bitIndex % 64)

	for {
		old := atomic.LoadInt64(&ab.bits[word])
		new := old | (1 << bit)
		if atomic.CompareAndSwapInt64(&ab.bits[word], old, new) {
			return
		}
	}
}

func (ab *AtomicBitmap) ClearBit(bitIndex int) {
	if bitIndex < 0 || bitIndex >= 256 {
		return
	}

	word := bitIndex / 64
	bit := uint(bitIndex % 64)

	for {
		old := atomic.LoadInt64(&ab.bits[word])
		new := old & ^(1 << bit)
		if atomic.CompareAndSwapInt64(&ab.bits[word], old, new) {
			return
		}
	}
}

func (ab *AtomicBitmap) TestBit(bitIndex int) bool {
	if bitIndex < 0 || bitIndex >= 256 {
		return false
	}

	word := bitIndex / 64
	bit := uint(bitIndex % 64)

	val := atomic.LoadInt64(&ab.bits[word])
	return (val & (1 << bit)) != 0
}

// ============================================================================
// HELPER FUNCTIONS FOR TESTING
// ============================================================================

func IncrementCounterNonAtomic(counter *int64) {
	*counter++
}

func IncrementCounterAtomic(counter *int64) {
	atomic.AddInt64(counter, 1)
}

// ============================================================================
// BONUS SOLUTION: Circular Buffer
// ============================================================================

func NewCircularBuffer(capacity int) *CircularBuffer {
	return &CircularBuffer{
		buffer:   make([]int64, capacity),
		head:     0,
		tail:     0,
		size:     0,
		capacity: int64(capacity),
	}
}

func (cb *CircularBuffer) Push(value int64) bool {
	for {
		size := atomic.LoadInt64(&cb.size)
		if size >= cb.capacity {
			return false // Buffer is full
		}

		// Try to increment size
		if !atomic.CompareAndSwapInt64(&cb.size, size, size+1) {
			continue // Retry
		}

		// Get tail position and increment
		tail := atomic.AddInt64(&cb.tail, 1) - 1
		index := tail % cb.capacity

		// Store value
		atomic.StoreInt64(&cb.buffer[index], value)
		return true
	}
}

func (cb *CircularBuffer) Pop() (int64, bool) {
	for {
		size := atomic.LoadInt64(&cb.size)
		if size <= 0 {
			return 0, false // Buffer is empty
		}

		// Try to decrement size
		if !atomic.CompareAndSwapInt64(&cb.size, size, size-1) {
			continue // Retry
		}

		// Get head position and increment
		head := atomic.AddInt64(&cb.head, 1) - 1
		index := head % cb.capacity

		// Load value
		value := atomic.LoadInt64(&cb.buffer[index])
		return value, true
	}
}
