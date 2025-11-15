//go:build !solution
// +build !solution

// Package exercise contains hands-on exercises for atomic operations.

package exercise

import "sync/atomic"

// ============================================================================
// EXERCISE 1: Atomic Counter
// ============================================================================

// NewAtomicCounter creates a new atomic counter initialized to 0.
//
// REQUIREMENTS:
// - Return a pointer to AtomicCounter with value = 0
//
// EXAMPLE:
//   counter := NewAtomicCounter()
//   counter.Increment()
//   fmt.Println(counter.Value())  // 1
func NewAtomicCounter() *AtomicCounter {
	// TODO: Implement this
	return nil
}

// Increment atomically increments the counter by 1.
//
// REQUIREMENTS:
// - Use atomic.AddInt64 to increment value
//
// HINT: atomic.AddInt64(&c.value, 1)
func (c *AtomicCounter) Increment() {
	// TODO: Implement this
}

// Decrement atomically decrements the counter by 1.
//
// REQUIREMENTS:
// - Use atomic.AddInt64 to decrement value
//
// HINT: Use a negative delta
func (c *AtomicCounter) Decrement() {
	// TODO: Implement this
}

// Add atomically adds delta to the counter.
//
// REQUIREMENTS:
// - Use atomic.AddInt64 to add delta
//
// HINT: atomic.AddInt64(&c.value, delta)
func (c *AtomicCounter) Add(delta int64) {
	// TODO: Implement this
}

// Value atomically reads the counter value.
//
// REQUIREMENTS:
// - Use atomic.LoadInt64 to read value
//
// HINT: return atomic.LoadInt64(&c.value)
func (c *AtomicCounter) Value() int64 {
	// TODO: Implement this
	return 0
}

// Reset atomically sets the counter to 0 and returns the old value.
//
// REQUIREMENTS:
// - Use atomic.SwapInt64 to set value to 0 and return old value
//
// HINT: return atomic.SwapInt64(&c.value, 0)
func (c *AtomicCounter) Reset() int64 {
	// TODO: Implement this
	return 0
}

// ============================================================================
// EXERCISE 2: Atomic Flag
// ============================================================================

// NewAtomicFlag creates a new atomic flag initialized to false.
//
// REQUIREMENTS:
// - Return a pointer to AtomicFlag with value = 0 (false)
//
// EXAMPLE:
//   flag := NewAtomicFlag()
//   flag.Set()
//   fmt.Println(flag.IsSet())  // true
func NewAtomicFlag() *AtomicFlag {
	// TODO: Implement this
	return nil
}

// Set atomically sets the flag to true.
//
// REQUIREMENTS:
// - Use atomic.StoreInt64 to set value to 1
//
// HINT: atomic.StoreInt64(&f.value, 1)
func (f *AtomicFlag) Set() {
	// TODO: Implement this
}

// Clear atomically sets the flag to false.
//
// REQUIREMENTS:
// - Use atomic.StoreInt64 to set value to 0
func (f *AtomicFlag) Clear() {
	// TODO: Implement this
}

// IsSet atomically reads the flag value.
//
// REQUIREMENTS:
// - Use atomic.LoadInt64 to read value
// - Return true if value == 1, false otherwise
func (f *AtomicFlag) IsSet() bool {
	// TODO: Implement this
	return false
}

// TestAndSet atomically sets the flag to true and returns the old value.
//
// REQUIREMENTS:
// - Use atomic.SwapInt64 to set value to 1 and return old value
// - Return true if old value was 1, false if it was 0
//
// HINT: old := atomic.SwapInt64(&f.value, 1); return old == 1
func (f *AtomicFlag) TestAndSet() bool {
	// TODO: Implement this
	return false
}

// ============================================================================
// EXERCISE 3: Rate Limiter
// ============================================================================

// NewRateLimiter creates a new rate limiter.
//
// REQUIREMENTS:
// - Set maxTokens to capacity
// - Set refillRate to tokensPerSecond
// - Set tokens to capacity (start full)
// - Set lastRefill to current Unix timestamp
// - Return pointer to RateLimiter
//
// EXAMPLE:
//   limiter := NewRateLimiter(10, 5)  // 10 max tokens, refill 5/sec
//   limiter.Allow()  // true (consumes 1 token)
func NewRateLimiter(capacity, tokensPerSecond int64) *RateLimiter {
	// TODO: Implement this
	return nil
}

// Allow attempts to consume one token.
//
// REQUIREMENTS:
// - Refill tokens based on elapsed time since lastRefill
// - Try to consume one token atomically
// - Return true if successful, false if no tokens available
//
// ALGORITHM:
// 1. Get current time (Unix seconds)
// 2. Load lastRefill atomically
// 3. If time has passed, try to update lastRefill with CAS
// 4. If successful, calculate new tokens and store
// 5. Try to consume a token using CAS loop
//
// HINT: Use atomic.LoadInt64, atomic.CompareAndSwapInt64
func (rl *RateLimiter) Allow() bool {
	// TODO: Implement this
	return false
}

// ============================================================================
// EXERCISE 4: Atomic Max Tracker
// ============================================================================

// NewAtomicMax creates a new atomic max tracker.
//
// REQUIREMENTS:
// - Initialize max to minimum possible int64 value
// - Return pointer to AtomicMax
//
// HINT: Use math.MinInt64 or -9223372036854775808
func NewAtomicMax() *AtomicMax {
	// TODO: Implement this
	return nil
}

// Update atomically updates the maximum if value is greater.
//
// REQUIREMENTS:
// - Use CAS loop to update max only if value > current max
// - Don't update if value <= current max
//
// ALGORITHM:
// 1. Load current max
// 2. If value <= max, return (no update needed)
// 3. Try to CAS max from old to value
// 4. If CAS fails, retry (another goroutine updated max)
//
// HINT: Use atomic.LoadInt64 and atomic.CompareAndSwapInt64 in a loop
func (am *AtomicMax) Update(value int64) {
	// TODO: Implement this
}

// Max atomically reads the current maximum.
//
// REQUIREMENTS:
// - Use atomic.LoadInt64 to read max
func (am *AtomicMax) Max() int64 {
	// TODO: Implement this
	return 0
}

// ============================================================================
// EXERCISE 5: SpinLock
// ============================================================================

// NewSpinLock creates a new spinlock.
//
// REQUIREMENTS:
// - Initialize state to 0 (unlocked)
// - Return pointer to SpinLock
func NewSpinLock() *SpinLock {
	// TODO: Implement this
	return nil
}

// Lock acquires the spinlock (busy-waits if locked).
//
// REQUIREMENTS:
// - Use atomic.SwapInt64 to try to set state to 1
// - If old value was 0, lock acquired
// - If old value was 1, spin (loop and retry)
//
// ALGORITHM:
// 1. Loop: Swap state with 1
// 2. If old value was 0, break (acquired)
// 3. If old value was 1, continue spinning
//
// HINT: for atomic.SwapInt64(&sl.state, 1) != 0 { }
func (sl *SpinLock) Lock() {
	// TODO: Implement this
}

// Unlock releases the spinlock.
//
// REQUIREMENTS:
// - Use atomic.StoreInt64 to set state to 0
func (sl *SpinLock) Unlock() {
	// TODO: Implement this
}

// ============================================================================
// EXERCISE 6: Atomic State Machine
// ============================================================================

// NewAtomicState creates a new atomic state machine.
//
// REQUIREMENTS:
// - Initialize current to StateIdle
// - Return pointer to AtomicState
func NewAtomicState() *AtomicState {
	// TODO: Implement this
	return nil
}

// CurrentState atomically reads the current state.
//
// REQUIREMENTS:
// - Use atomic.LoadInt64 to read current
func (as *AtomicState) CurrentState() int64 {
	// TODO: Implement this
	return 0
}

// Transition atomically transitions from expectedCurrent to newState.
//
// REQUIREMENTS:
// - Use atomic.CompareAndSwapInt64 to transition
// - Return true if transition succeeded, false otherwise
//
// EXAMPLE:
//   state := NewAtomicState()  // StateIdle
//   ok := state.Transition(StateIdle, StateRunning)  // true
//   ok = state.Transition(StateIdle, StateRunning)   // false (not in StateIdle)
//
// HINT: return atomic.CompareAndSwapInt64(&as.current, expectedCurrent, newState)
func (as *AtomicState) Transition(expectedCurrent, newState int64) bool {
	// TODO: Implement this
	return false
}

// ============================================================================
// EXERCISE 7: Reference Counter
// ============================================================================

// NewReferenceCounter creates a new reference counter.
//
// REQUIREMENTS:
// - Initialize count to 0
// - Return pointer to ReferenceCounter
func NewReferenceCounter() *ReferenceCounter {
	// TODO: Implement this
	return nil
}

// Acquire increments the reference count.
//
// REQUIREMENTS:
// - Use atomic.AddInt64 to increment count by 1
func (rc *ReferenceCounter) Acquire() {
	// TODO: Implement this
}

// Release decrements the reference count and returns true if count reached 0.
//
// REQUIREMENTS:
// - Use atomic.AddInt64 to decrement count by 1
// - Return true if new count is 0, false otherwise
//
// HINT: newCount := atomic.AddInt64(&rc.count, -1); return newCount == 0
func (rc *ReferenceCounter) Release() bool {
	// TODO: Implement this
	return false
}

// Count atomically reads the current count.
//
// REQUIREMENTS:
// - Use atomic.LoadInt64 to read count
func (rc *ReferenceCounter) Count() int64 {
	// TODO: Implement this
	return 0
}

// ============================================================================
// EXERCISE 8: Config Manager
// ============================================================================

// NewConfigManager creates a new config manager.
//
// REQUIREMENTS:
// - Initialize with default config (Timeout=30, MaxRetries=3, BatchSize=100)
// - Use atomic.Value.Store to store the config
// - Return pointer to ConfigManager
//
// HINT: cm := &ConfigManager{}; cm.config.Store(&Config{...}); return cm
func NewConfigManager() *ConfigManager {
	// TODO: Implement this
	return nil
}

// Update atomically updates the configuration.
//
// REQUIREMENTS:
// - Use atomic.Value.Store to store the new config pointer
//
// HINT: cm.config.Store(newConfig)
func (cm *ConfigManager) Update(newConfig *Config) {
	// TODO: Implement this
}

// Get atomically reads the current configuration.
//
// REQUIREMENTS:
// - Use atomic.Value.Load to load the config
// - Type assert to *Config
// - Return the config pointer
//
// HINT: return cm.config.Load().(*Config)
func (cm *ConfigManager) Get() *Config {
	// TODO: Implement this
	return nil
}

// ============================================================================
// EXERCISE 9: Load Balancer
// ============================================================================

// NewLoadBalancer creates a new load balancer.
//
// REQUIREMENTS:
// - Initialize counter to 0
// - Set workers to numWorkers
// - Return pointer to LoadBalancer
func NewLoadBalancer(numWorkers int) *LoadBalancer {
	// TODO: Implement this
	return nil
}

// NextWorker returns the ID of the next worker using round-robin.
//
// REQUIREMENTS:
// - Use atomic.AddInt64 to increment counter
// - Return (counter - 1) % workers as the worker ID
//
// ALGORITHM:
// 1. Atomically increment counter and get new value
// 2. Return (value - 1) % numWorkers
//
// EXAMPLE:
//   lb := NewLoadBalancer(3)  // 3 workers
//   lb.NextWorker()  // 0
//   lb.NextWorker()  // 1
//   lb.NextWorker()  // 2
//   lb.NextWorker()  // 0 (wraps around)
//
// HINT: val := atomic.AddInt64(&lb.counter, 1); return (val - 1) % lb.workers
func (lb *LoadBalancer) NextWorker() int64 {
	// TODO: Implement this
	return 0
}

// ============================================================================
// EXERCISE 10: Atomic Bitmap
// ============================================================================

// NewAtomicBitmap creates a new atomic bitmap.
//
// REQUIREMENTS:
// - Initialize all bits to 0
// - Return pointer to AtomicBitmap
func NewAtomicBitmap() *AtomicBitmap {
	// TODO: Implement this
	return nil
}

// SetBit atomically sets a bit to 1.
//
// REQUIREMENTS:
// - Calculate which int64 contains the bit (bitIndex / 64)
// - Calculate bit position within that int64 (bitIndex % 64)
// - Use CAS loop to set the bit
//
// ALGORITHM:
// 1. word := bitIndex / 64
// 2. bit := bitIndex % 64
// 3. Loop:
//    a. Load current value of bits[word]
//    b. newVal := oldVal | (1 << bit)
//    c. Try CAS, break if successful
//
// HINT: Use atomic.LoadInt64 and atomic.CompareAndSwapInt64
func (ab *AtomicBitmap) SetBit(bitIndex int) {
	// TODO: Implement this
}

// ClearBit atomically sets a bit to 0.
//
// REQUIREMENTS:
// - Similar to SetBit, but clear the bit instead
//
// HINT: newVal := oldVal & ^(1 << bit)
func (ab *AtomicBitmap) ClearBit(bitIndex int) {
	// TODO: Implement this
}

// TestBit atomically reads a bit value.
//
// REQUIREMENTS:
// - Calculate which int64 contains the bit
// - Calculate bit position within that int64
// - Use atomic.LoadInt64 to read the word
// - Return true if bit is set, false otherwise
//
// HINT: word := bitIndex / 64; bit := bitIndex % 64
//       val := atomic.LoadInt64(&ab.bits[word])
//       return (val & (1 << bit)) != 0
func (ab *AtomicBitmap) TestBit(bitIndex int) bool {
	// TODO: Implement this
	return false
}

// ============================================================================
// HELPER FUNCTIONS FOR TESTING
// ============================================================================

// IncrementCounterNonAtomic is a BUGGY implementation for testing.
// This will fail under concurrent access (race condition).
func IncrementCounterNonAtomic(counter *int64) {
	*counter++
}

// IncrementCounterAtomic is the CORRECT implementation for testing.
func IncrementCounterAtomic(counter *int64) {
	atomic.AddInt64(counter, 1)
}

// ============================================================================
// BONUS EXERCISE: Circular Buffer
// ============================================================================

// NewCircularBuffer creates a new lock-free circular buffer.
//
// REQUIREMENTS:
// - Allocate buffer of size capacity
// - Initialize head, tail, size to 0
// - Set capacity
// - Return pointer to CircularBuffer
//
// NOTE: This is a challenging exercise! Lock-free circular buffers
// require careful handling of head/tail indices and size.
func NewCircularBuffer(capacity int) *CircularBuffer {
	// TODO: Implement this
	return nil
}

// Push adds a value to the buffer (returns false if full).
//
// REQUIREMENTS:
// - Check if buffer is full (size == capacity)
// - If not full, add value at tail position
// - Increment tail (wrap around using modulo)
// - Increment size
// - Return true if successful, false if full
//
// NOTE: This needs careful synchronization! Consider using CAS loops.
func (cb *CircularBuffer) Push(value int64) bool {
	// TODO: Implement this (challenging!)
	return false
}

// Pop removes and returns a value from the buffer (returns 0, false if empty).
//
// REQUIREMENTS:
// - Check if buffer is empty (size == 0)
// - If not empty, read value at head position
// - Increment head (wrap around using modulo)
// - Decrement size
// - Return value and true if successful, 0 and false if empty
//
// NOTE: This needs careful synchronization! Consider using CAS loops.
func (cb *CircularBuffer) Pop() (int64, bool) {
	// TODO: Implement this (challenging!)
	return 0, false
}
