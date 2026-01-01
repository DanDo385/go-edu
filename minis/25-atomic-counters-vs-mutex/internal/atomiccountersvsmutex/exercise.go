//go:build !solution && !reference

// Package exercise contains hands-on exercises for atomic operations.

package atomiccountersvsmutex

import "sync/atomic"

// ============================================================================
// EXERCISE 1: Atomic Counter
// ============================================================================

// NewAtomicCounter creates a new atomic counter initialized to 0.
func NewAtomicCounter() *AtomicCounter {
	// TODO: Implement NewAtomicCounter
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Increment atomically increments the counter by 1.
func (c *AtomicCounter) Increment() {
	// TODO: Implement this function.
	// - Use `atomic.AddInt64` to safely add 1 to the value.
	// - `atomic.AddInt64` takes a pointer to the integer and the delta.
	// - It's a single, uninterruptible hardware instruction, making it much faster than a mutex for simple arithmetic.
}

// Decrement atomically decrements the counter by 1.
func (c *AtomicCounter) Decrement() {
	// TODO: Implement this function.
	// - Use `atomic.AddInt64` with a negative delta.
}

// Add atomically adds delta to the counter.
func (c *AtomicCounter) Add(delta int64) {
	// TODO: Implement this function.
	// - Use `atomic.AddInt64` with the given delta.
}

// Value atomically reads the counter value.
func (c *AtomicCounter) Value() int64 {
	// TODO: Implement this function.
	// - A simple `return c.value` would be a "dirty read" and is not safe.
	// - Use `atomic.LoadInt64` to safely read the value without being affected by another goroutine's write operation.
	return 0
}

// Reset atomically sets the counter to 0 and returns the old value.
func (c *AtomicCounter) Reset() int64 {
	// TODO: Implement this function.
	// - Use `atomic.SwapInt64` to set the value to 0 and get the previous value in a single atomic operation.
	return 0
}

// ============================================================================
// EXERCISE 2: Atomic Flag
// ============================================================================

// NewAtomicFlag creates a new atomic flag initialized to false.
func NewAtomicFlag() *AtomicFlag {
	// TODO: Implement NewAtomicFlag
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Set atomically sets the flag to true.
func (f *AtomicFlag) Set() {
	// TODO: Implement this function.
	// - Use `atomic.StoreInt64` to unconditionally set the value to 1 (representing `true`).
}

// Clear atomically sets the flag to false.
func (f *AtomicFlag) Clear() {
	// TODO: Implement this function.
	// - Use `atomic.StoreInt64` to set the value to 0 (representing `false`).
}

// IsSet atomically reads the flag value.
func (f *AtomicFlag) IsSet() bool {
	// TODO: Implement this function.
	// - Use `atomic.LoadInt64` to safely read the value.
	// - Return `true` if the loaded value is 1.
	return false
}

// TestAndSet atomically sets the flag to true and returns the old value.
func (f *AtomicFlag) TestAndSet() bool {
	// TODO: Implement this function.
	// - This is a common atomic operation, often called "get and set".
	// - Use `atomic.SwapInt64` to set the value to 1 and get the old value in a single atomic step.
	// - Return `true` if the old value was 1.
	return false
}

// ============================================================================
// EXERCISE 3: Rate Limiter
// ============================================================================

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(capacity, tokensPerSecond int64) *RateLimiter {
	// TODO: Implement NewRateLimiter
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Allow attempts to consume one token.
func (rl *RateLimiter) Allow() bool {
	// TODO: Implement this function.

	// This is a lock-free implementation of the token bucket algorithm. It's more complex than a mutex-based one but can have higher performance under heavy contention.

	// Step 1: Refill tokens if necessary.
	// - `now := time.Now().Unix()`
	// - `last := atomic.LoadInt64(&rl.lastRefill)`
	// - If `now > last`, it's time to refill.
	// - To prevent multiple goroutines from refilling at the same time, use `atomic.CompareAndSwapInt64(&rl.lastRefill, last, now)`. Only the goroutine that succeeds in this CAS will perform the refill.
	// - The winning goroutine should calculate the number of new tokens (`elapsed * rl.refillRate`) and add them to `rl.tokens` using another CAS loop to handle concurrent updates. Be careful not to exceed `rl.maxTokens`.

	// Step 2: Try to consume a token.
	// - Use a `for` loop and CAS to decrement `rl.tokens`.
	// - `tokens := atomic.LoadInt64(&rl.tokens)`
	// - If `tokens <= 0`, return `false`.
	// - `if atomic.CompareAndSwapInt64(&rl.tokens, tokens, tokens-1)`, the operation succeeded, so return `true`.
	// - If the CAS fails, the loop will retry, loading the new value of `tokens`.
	return false
}

// ============================================================================
// EXERCISE 4: Atomic Max Tracker
// ============================================================================

// NewAtomicMax creates a new atomic max tracker.
func NewAtomicMax() *AtomicMax {
	// TODO: Implement NewAtomicMax
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Update atomically updates the maximum if value is greater.
func (am *AtomicMax) Update(value int64) {
	// TODO: Implement this function.

	// This function uses a "Compare-And-Swap" (CAS) loop, which is a common pattern in lock-free programming.

	// Step 1: Start an infinite `for` loop.
	// Step 2: Inside the loop, atomically load the current maximum value.
	// - `current := atomic.LoadInt64(&am.max)`
	// Step 3: Compare the new `value` with the `current` max.
	// - If `value <= current`, then there's nothing to update. You can `return` from the function.
	// Step 4: If the new `value` is greater, try to update the max.
	// - `if atomic.CompareAndSwapInt64(&am.max, current, value)`
	//   - This operation says: "If the value at `&am.max` is still `current`, then update it to `value`." It does this in a single, atomic step.
	//   - If it returns `true`, you successfully updated the max. You can `return` from the function.
	//   - If it returns `false`, it means another goroutine changed `am.max` between your `Load` and your `CompareAndSwap`. The loop will then repeat, loading the new current value and trying again.
}

// Max atomically reads the current maximum.
func (am *AtomicMax) Max() int64 {
	// TODO: Implement this function.
	// - Use `atomic.LoadInt64` for a safe read.
	return 0
}

// ============================================================================
// EXERCISE 5: SpinLock
// ============================================================================

// NewSpinLock creates a new spinlock.
func NewSpinLock() *SpinLock {
	// TODO: Implement NewSpinLock
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Lock acquires the spinlock (busy-waits if locked).
func (sl *SpinLock) Lock() {
	// TODO: Implement this function.

	// A spinlock repeatedly tries to acquire a lock in a tight loop without sleeping.
	// This is efficient if the lock is held for a very short time, but wastes CPU if the lock is held for a long time.

	// Step 1: Use a `for` loop that attempts to acquire the lock.
	// - `atomic.SwapInt64(&sl.state, 1)` attempts to set the state to 1 (locked) and returns the *old* state.
	// - If the old state was 0 (unlocked), then you have successfully acquired the lock. The loop condition `... != 0` will be false, and the loop will terminate.
	// - If the old state was 1 (locked), it means someone else holds the lock. The swap still sets the state to 1, and the loop condition `1 != 0` is true, so the loop continues to "spin".
	// - A common implementation is `for atomic.SwapInt64(&sl.state, 1) != 0 {}`.
	// - In real-world code, you might add `runtime.Gosched()` inside the loop to yield the processor, preventing the spinner from starving other goroutines.
}

// Unlock releases the spinlock.
func (sl *SpinLock) Unlock() {
	// TODO: Implement this function.
	// - To unlock, simply set the state back to 0.
	// - Use `atomic.StoreInt64(&sl.state, 0)`.
}

// ============================================================================
// EXERCISE 6: Atomic State Machine
// ============================================================================

// NewAtomicState creates a new atomic state machine.
func NewAtomicState() *AtomicState {
	// TODO: Implement NewAtomicState
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// CurrentState atomically reads the current state.
func (as *AtomicState) CurrentState() int64 {
	// TODO: Implement this function.
	// - Use `atomic.LoadInt64` to safely read the `current` state.
	return 0
}

// Transition atomically transitions from expectedCurrent to newState.
func (as *AtomicState) Transition(expectedCurrent, newState int64) bool {
	// TODO: Implement this function.
	// - This is the core of the state machine's safety.
	// - Use `atomic.CompareAndSwapInt64`.
	// - It will only set the state to `newState` if the `current` state is equal to `expectedCurrent`.
	// - It returns `true` if the swap was successful, and `false` otherwise. This allows the caller to know if their attempted state transition was valid.
	return false
}

// ============================================================================
// EXERCISE 7: Reference Counter
// ============================================================================

// NewReferenceCounter creates a new reference counter.
func NewReferenceCounter() *ReferenceCounter {
	// TODO: Implement NewReferenceCounter
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Acquire increments the reference count.
func (rc *ReferenceCounter) Acquire() {
	// TODO: Implement this function.
	// - Atomically increment the `count`.
	// - Use `atomic.AddInt64`.
}

// Release decrements the reference count and returns true if count reached 0.
func (rc *ReferenceCounter) Release() bool {
	// TODO: Implement this function.
	// - Atomically decrement the `count`.
	// - `atomic.AddInt64` returns the *new* value after the addition.
	// - You can check this new value to see if the count has reached zero.
	// - Return `true` if the new count is 0, indicating that the resource can now be freed.
	return false
}

// Count atomically reads the current count.
func (rc *ReferenceCounter) Count() int64 {
	// TODO: Implement this function.
	// - Use `atomic.LoadInt64` for a safe read.
	return 0
}

// ============================================================================
// EXERCISE 8: Config Manager
// ============================================================================

// NewConfigManager creates a new config manager.
func NewConfigManager() *ConfigManager {
	// TODO: Implement NewConfigManager
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// Update atomically updates the configuration.
func (cm *ConfigManager) Update(newConfig *Config) {
	// TODO: Implement this function.
	// - Use `cm.config.Store(newConfig)` to atomically swap the old config pointer with the new one.
	// - All subsequent calls to `Get()` will receive this new pointer.
}

// Get atomically reads the current configuration.
func (cm *ConfigManager) Get() *Config {
	// TODO: Implement this function.
	// - Use `cm.config.Load()` to get the currently stored value.
	// - The value returned by `Load()` is of type `interface{}`, so you must type-assert it back to a `*Config`.
	// - `return cm.config.Load().(*Config)`
	return nil
}

// ============================================================================
// EXERCISE 9: Load Balancer
// ============================================================================

// NewLoadBalancer creates a new load balancer.
func NewLoadBalancer(numWorkers int) *LoadBalancer {
	// TODO: Implement NewLoadBalancer
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// NextWorker returns the ID of the next worker using round-robin.
func (lb *LoadBalancer) NextWorker() int64 {
	// TODO: Implement this function.

	// This implements a thread-safe round-robin scheduler.

	// Step 1: Atomically increment the counter and get the new value.
	// - `val := atomic.AddInt64(&lb.counter, 1)`

	// Step 2: Use the modulo operator to map the counter value to a worker ID.
	// - `return (val - 1) % lb.workers`
	// - We use `val - 1` because the first value returned by `AddInt64` will be 1, and we want our worker IDs to be 0-indexed.
	return 0
}

// ============================================================================
// EXERCISE 10: Atomic Bitmap
// ============================================================================

// NewAtomicBitmap creates a new atomic bitmap.
func NewAtomicBitmap() *AtomicBitmap {
	// TODO: Implement NewAtomicBitmap
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// SetBit atomically sets a bit to 1.
func (ab *AtomicBitmap) SetBit(bitIndex int) {
	// TODO: Implement this function.

	// A bitmap uses individual bits within a block of integers to store boolean flags.
	// This is a memory-efficient way to store a large set of flags.

	// Step 1: Validate the index and calculate the word and bit position.
	// - `word := bitIndex / 64` (Which `int64` in our array holds the bit)
	// - `bit := uint(bitIndex % 64)` (Which bit within that `int64`)

	// Step 2: Use a CAS loop to set the bit.
	// - `for { ... }`
	// - `old := atomic.LoadInt64(&ab.bits[word])`
	// - `new := old | (1 << bit)` (Use bitwise OR to set the bit)
	// - `if atomic.CompareAndSwapInt64(&ab.bits[word], old, new) { return }`
}

// ClearBit atomically sets a bit to 0.
func (ab *AtomicBitmap) ClearBit(bitIndex int) {
	// TODO: Implement this function.

	// Step 1: Calculate `word` and `bit` as in `SetBit`.
	// Step 2: Use a CAS loop.
	// - `new := old & ^(1 << bit)` (Use bitwise AND NOT to clear the bit)
	// - The rest of the loop is the same as `SetBit`.
}

// TestBit atomically reads a bit value.
func (ab *AtomicBitmap) TestBit(bitIndex int) bool {
	// TODO: Implement this function.

	// This is a read-only operation.

	// Step 1: Calculate `word` and `bit` as in `SetBit`.
	// Step 2: Atomically load the word.
	// - `val := atomic.LoadInt64(&ab.bits[word])`
	// Step 3: Check if the specific bit is set.
	// - `return (val & (1 << bit)) != 0`
	return false
}

// ============================================================================
// HELPER FUNCTIONS FOR TESTING
// ============================================================================

// IncrementCounterNonAtomic is a BUGGY implementation for testing.
// This will fail under concurrent access (race condition).
func IncrementCounterNonAtomic(counter *int64) {
	// TODO: Implement IncrementCounterNonAtomic
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// IncrementCounterAtomic is the CORRECT implementation for testing.
func IncrementCounterAtomic(counter *int64) {
	// TODO: Implement IncrementCounterAtomic
	// See solution.reference.go for reference implementation
	panic("not implemented")
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
	// TODO: Implement NewCircularBuffer
	// See solution.reference.go for reference implementation
	panic("not implemented")
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
