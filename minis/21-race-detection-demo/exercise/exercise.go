//go:build !solution
// +build !solution

package exercise

// TODO: Import required packages
// You'll need:
// - "fmt" for error formatting
// - "sync" for Mutex, RWMutex, WaitGroup, Once
// - "sync/atomic" for atomic operations (Go 1.19+)
//
// import (
//     "fmt"
//     "sync"
//     "sync/atomic"  // For atomic.Int64
// )

// ============================================================================
// RACE DETECTION: Understanding Data Races and Memory Visibility
// ============================================================================
//
// A data race occurs when:
// 1. Two or more goroutines access the same memory location
// 2. At least one of them is writing
// 3. The accesses are NOT synchronized
//
// Why races are dangerous:
// - Unpredictable behavior (works sometimes, fails others)
// - Memory corruption (partially written values)
// - No compiler or runtime guarantees about what you'll read
// - Can cause crashes, data loss, security vulnerabilities
//
// How to detect races:
// - Run: go test -race
// - The race detector instruments your code to track all memory accesses
// - Reports races with full stack traces
// - Only catches races that actually execute during the test
//
// How to fix races:
// 1. Mutexes: sync.Mutex (exclusive access) or sync.RWMutex (reader/writer)
// 2. Atomic operations: atomic.Int64, atomic.Bool (lock-free)
// 3. Channels: Communicate by passing data, not sharing it
// 4. sync.Once: Guarantee one-time initialization
//
// ============================================================================

// ============================================================================
// Exercise 1: Fix the Counter Race
// ============================================================================

// SafeCounter is a thread-safe counter.
// TODO: Add fields needed for synchronization.
//
// Memory considerations:
// - The counter value must be protected from concurrent access
// - Option 1: Use sync.Mutex + int64 field
//   * Mutex is a struct containing internal state (lock word + wait queue)
//   * Size: ~16 bytes (platform dependent)
//   * Always use as a field, not embedded (makes copying accidental)
//
// - Option 2: Use atomic.Int64 (preferred for simple counters)
//   * Lock-free (uses CPU atomic instructions like LOCK XADD)
//   * Size: 8 bytes (just the int64)
//   * Much faster than mutex for uncontended access
//   * Must be aligned to 8-byte boundary (Go handles this automatically)
//
// type SafeCounter struct {
//     value atomic.Int64  // Preferred: lock-free atomic counter
//     // OR
//     value int64         // Alternative: regular int64
//     mu    sync.Mutex    // protected by this mutex
// }

// NewSafeCounter creates a new thread-safe counter.
// TODO: Initialize the counter with necessary fields.
//
// Memory allocation:
// - Returns a pointer (*SafeCounter) because:
//   1. Counters are meant to be shared and modified
//   2. Mutexes/atomics must not be copied (copying breaks synchronization)
//   3. Returning *SafeCounter allows multiple goroutines to share same instance
//
// - The struct is allocated on the heap (escapes because we return a pointer)
// - Caller receives a pointer to heap memory
// - GC will clean up when all references are gone
//
// func NewSafeCounter() *SafeCounter {
//     return &SafeCounter{}  // Zero value is valid for atomic.Int64
// }

// Increment safely increments the counter by 1.
// TODO: Implement thread-safe increment.
//
// This function will be called concurrently from multiple goroutines.
// Without synchronization, this would be a race:
//   goroutine 1: reads value (100)
//   goroutine 2: reads value (100)  <- both read same value!
//   goroutine 1: writes 101
//   goroutine 2: writes 101         <- lost update! should be 102
//
// Solution 1: Atomic operations (recommended for counters)
// - Use: c.value.Add(1)
// - This is a single atomic CPU instruction (LOCK XADD on x86)
// - No locks, no contention, very fast
// - Guarantees all goroutines see updates in consistent order
//
// Solution 2: Mutex-based
// - Use: c.mu.Lock(); c.value++; c.mu.Unlock()
// - Slower than atomic (needs to acquire/release lock)
// - But works for complex updates (multiple fields)
//
// func (c *SafeCounter) Increment() {
//     c.value.Add(1)  // Atomic: no lock needed
// }

// Value safely returns the current counter value.
// TODO: Implement thread-safe read.
//
// Reading is also a race if not synchronized!
// Without protection:
//   goroutine 1: writes value = 100 (bytes: 00 00 00 00 00 00 00 64)
//   goroutine 2: reads value → might see partial write!
//     (e.g., 00 00 00 00 00 00 00 00 if write isn't complete)
//
// Solution: Must synchronize with writes
// - Atomic: c.value.Load()  (reads with memory fence)
// - Mutex: c.mu.Lock(); v := c.value; c.mu.Unlock()
//
// func (c *SafeCounter) Value() int64 {
//     return c.value.Load()  // Atomic read
// }

// ============================================================================
// Exercise 2: Fix the Map Race
// ============================================================================

// SafeMap is a thread-safe map wrapper.
// TODO: Add fields for map storage and synchronization.
//
// Memory considerations:
// - Maps in Go are NOT thread-safe (concurrent access = crash or corruption)
// - Must protect all map operations (read, write, iteration)
// - Use sync.RWMutex for maps (many readers, few writers)
//
// type SafeMap struct {
//     data map[string]int  // The actual map (reference type, points to hash table)
//     mu   sync.RWMutex    // Protects data
// }
//
// Why RWMutex instead of Mutex?
// - RWMutex allows multiple concurrent readers (RLock)
// - Only one writer at a time (Lock)
// - Readers and writer are mutually exclusive
// - Perfect for read-heavy workloads (like caches)
//
// Pointer vs value receiver:
// - Methods use pointer receiver (*SafeMap) because:
//   1. Maps are reference types, but the mutex must not be copied
//   2. All methods must work on the same mutex instance
//   3. Value receivers would copy the struct (and the mutex) = broken sync

// NewSafeMap creates a new thread-safe map.
// TODO: Initialize the map with necessary fields.
//
// CRITICAL: Must use make() to initialize the map!
// - var m map[string]int creates a nil map (can read but NOT write)
// - make(map[string]int) creates an empty map (can read and write)
//
// Memory allocation:
// - make() allocates a hash table on the heap
// - Initial capacity is small (grows as needed)
// - Can specify initial capacity: make(map[string]int, 100)
//
// func NewSafeMap() *SafeMap {
//     return &SafeMap{
//         data: make(map[string]int),  // Initialize map!
//     }
// }

// Set safely stores a key-value pair.
// TODO: Implement thread-safe write to map.
//
// Map writes are NOT thread-safe:
// - Concurrent writes can corrupt the hash table
// - Go will crash with "fatal error: concurrent map writes"
//
// Solution: Use Lock (exclusive access)
// - mu.Lock() acquires exclusive access (blocks other Lock/RLock)
// - Modify map
// - mu.Unlock() releases lock
//
// Pattern:
//   m.mu.Lock()
//   m.data[key] = value  // Only one goroutine can do this at a time
//   m.mu.Unlock()
//
// func (m *SafeMap) Set(key string, value int) {
//     m.mu.Lock()
//     m.data[key] = value
//     m.mu.Unlock()
// }

// Get safely retrieves a value by key.
// TODO: Implement thread-safe read from map.
//
// Returns the value and a boolean indicating if the key exists.
//
// Map reads while writing = crash:
// - Even reading a map during a write can crash
// - Must synchronize reads with writes
//
// Solution: Use RLock (shared read access)
// - mu.RLock() allows multiple readers (blocks writers)
// - Read from map
// - mu.RUnlock() releases read lock
//
// Why RLock instead of Lock?
// - Multiple goroutines can read simultaneously (better performance)
// - Writers are blocked while readers are active (correctness)
//
// func (m *SafeMap) Get(key string) (int, bool) {
//     m.mu.RLock()
//     value, ok := m.data[key]
//     m.mu.RUnlock()
//     return value, ok
// }

// Len safely returns the number of entries in the map.
// TODO: Implement thread-safe length check.
//
// Even len(map) needs synchronization!
// - Reading map metadata during a write can give wrong results
//
// func (m *SafeMap) Len() int {
//     m.mu.RLock()
//     length := len(m.data)
//     m.mu.RUnlock()
//     return length
// }

// ============================================================================
// Exercise 3: Fix the Lazy Initialization Race
// ============================================================================

// LazyInit demonstrates thread-safe lazy initialization.
// TODO: Add fields for sync.Once and the value.
//
// Lazy initialization pattern:
// - Create expensive resource only when first needed
// - Common for singletons, database connections, config loading
//
// Problem without synchronization:
//   goroutine 1: if value == nil { value = init() }
//   goroutine 2: if value == nil { value = init() }  <- both see nil!
//   Result: init() called twice, waste of resources, possible race
//
// Solution: sync.Once
// - Guarantees the init function runs exactly once
// - Other goroutines block until first call completes
// - Subsequent calls return immediately (no lock contention)
//
// type LazyInit struct {
//     once  sync.Once      // Coordinates one-time initialization
//     value interface{}    // The lazily initialized value
// }
//
// Memory:
// - sync.Once contains: done uint32 + mutex
// - After first call, done=1, subsequent calls just read done (fast path)

// NewLazyInit creates a new lazy initializer.
// TODO: Initialize with necessary fields.
//
// func NewLazyInit() *LazyInit {
//     return &LazyInit{}  // Zero values are valid
// }

// GetOrInit returns the initialized value, initializing it if needed.
// TODO: Implement thread-safe lazy initialization.
//
// The init function should only be called ONCE, even when called concurrently.
//
// Pattern:
//   l.once.Do(func() {
//       l.value = init()
//   })
//   return l.value
//
// How sync.Once works:
// 1. First caller: Execute the function, set done=1
// 2. Concurrent callers: Block until first caller finishes
// 3. Future callers: See done=1, return immediately
//
// Memory visibility:
// - once.Do() includes memory barriers
// - Guarantees l.value write is visible to all goroutines after Do() returns
//
// func (l *LazyInit) GetOrInit(init func() interface{}) interface{} {
//     l.once.Do(func() {
//         l.value = init()
//     })
//     return l.value
// }

// ============================================================================
// Exercise 4: Fix the Slice Append Race
// ============================================================================

// SafeSlice is a thread-safe slice wrapper.
// TODO: Add fields for slice storage and synchronization.
//
// Slices are NOT thread-safe:
// - append() can reallocate the underlying array
// - Concurrent appends can corrupt the slice or lose data
//
// Why append is racy:
// - Slice is: {ptr *array, len int, cap int}
// - append() reads len, writes to ptr[len], updates len
// - Two goroutines appending:
//   G1: reads len=5, writes ptr[5], sets len=6
//   G2: reads len=5, writes ptr[5], sets len=6  <- overwrites G1's data!
//
// Solution: Protect with mutex
// - Use sync.Mutex (or RWMutex if you have many reads)
//
// type SafeSlice struct {
//     data []int        // The slice (header on stack, array on heap)
//     mu   sync.RWMutex // Protects data
// }

// NewSafeSlice creates a new thread-safe slice.
// TODO: Initialize with necessary fields.
//
// Pre-allocate capacity if you know approximate size:
// - make([]int, 0, 100) allocates array for 100 items
// - Reduces allocations during append
//
// func NewSafeSlice() *SafeSlice {
//     return &SafeSlice{
//         data: make([]int, 0),  // Empty slice
//     }
// }

// Append safely appends a value to the slice.
// TODO: Implement thread-safe append.
//
// Must protect: read length, append, update slice header
//
// func (s *SafeSlice) Append(value int) {
//     s.mu.Lock()
//     s.data = append(s.data, value)  // May reallocate array
//     s.mu.Unlock()
// }

// Get safely retrieves a value by index.
// TODO: Implement thread-safe indexed read.
//
// Returns the value and a boolean indicating if the index is valid.
//
// Race: Reading while another goroutine appends
// - Might read past end of slice (panic)
// - Might read from old array after reallocation (wrong data)
//
// func (s *SafeSlice) Get(index int) (int, bool) {
//     s.mu.RLock()
//     if index < 0 || index >= len(s.data) {
//         s.mu.RUnlock()
//         return 0, false
//     }
//     value := s.data[index]
//     s.mu.RUnlock()
//     return value, true
// }

// Len safely returns the length of the slice.
// TODO: Implement thread-safe length check.
//
// func (s *SafeSlice) Len() int {
//     s.mu.RLock()
//     length := len(s.data)
//     s.mu.RUnlock()
//     return length
// }

// ============================================================================
// Exercise 5: Fix the Loop Variable Capture Race
// ============================================================================

// ProcessIDs processes a list of IDs concurrently.
// TODO: Fix the race condition where all goroutines see the same ID.
//
// Common gotcha: Loop variable capture
// Problem:
//   for i, id := range ids {
//       go func() {
//           process(id)  // BUG: id is shared by all goroutines!
//       }()
//   }
//
// What happens:
// - Loop creates goroutines very fast
// - By the time goroutines run, loop has finished
// - All goroutines see the LAST value of id
//
// Solution 1: Pass as parameter (recommended)
//   for i, id := range ids {
//       go func(index int, value int) {
//           results[index] = process(value)
//       }(i, id)  // Pass as arguments
//   }
//
// Solution 2: Shadow the variable
//   for i, id := range ids {
//       i := i   // Create new variable
//       id := id // Create new variable
//       go func() {
//           results[i] = process(id)  // Now safe!
//       }()
//   }
//
// Note: Go 1.22+ does this shadowing automatically in for loops!
//
// Each goroutine should process a unique ID.
// Returns a slice of results in the same order as input.
//
// func ProcessIDs(ids []int, process func(int) int) []int {
//     var wg sync.WaitGroup
//     results := make([]int, len(ids))
//
//     for i, id := range ids {
//         wg.Add(1)
//         go func(index int, value int) {  // Pass as parameters!
//             defer wg.Done()
//             results[index] = process(value)
//         }(i, id)  // Arguments create copies
//     }
//
//     wg.Wait()
//     return results
// }

// ============================================================================
// Exercise 6: Concurrent URL Cache
// ============================================================================

// URLCache is a concurrent URL fetcher with caching.
// (Type definition is in types.go)

// Fetch fetches a URL's content, using cache if available.
// TODO: Implement thread-safe caching.
//
// Multiple goroutines may call this simultaneously.
// Each URL should only be fetched once (even if multiple goroutines request it simultaneously).
//
// Pattern (double-checked locking):
// 1. RLock, check cache, RUnlock (fast path for cache hit)
// 2. Lock, check cache again (another goroutine might have fetched it)
// 3. Fetch and store, Unlock
//
// Why double-check?
// - Many goroutines request same URL simultaneously
// - First goroutine: cache miss, acquires Lock, fetches
// - Other goroutines: wait for Lock, then see it's now in cache
// - Without double-check: all waiting goroutines would re-fetch!
//
// func (c *URLCache) Fetch(url string) (string, error) {
//     // Fast path: check cache with read lock
//     c.mu.RLock()
//     if content, ok := c.cache[url]; ok {
//         c.mu.RUnlock()
//         return content, nil
//     }
//     c.mu.RUnlock()
//
//     // Slow path: fetch with write lock
//     c.mu.Lock()
//     defer c.mu.Unlock()
//
//     // Double-check: another goroutine might have fetched it
//     if content, ok := c.cache[url]; ok {
//         return content, nil
//     }
//
//     // Fetch and cache
//     content, err := c.fetcher(url)
//     if err != nil {
//         return "", err
//     }
//     c.cache[url] = content
//     return content, nil
// }

// ============================================================================
// Exercise 7: Concurrent Metrics Tracking
// ============================================================================

// Metrics tracks application metrics concurrently.
// (Type definition is in types.go)

// IncrementRequests increments the request counter.
// TODO: Implement thread-safe increment.
//
// Use atomic operations for maximum performance:
// - m.requests.Add(1)
//
// func (m *Metrics) IncrementRequests() {
//     m.requests.Add(1)
// }

// IncrementErrors increments the error counter.
// TODO: Implement thread-safe increment.
//
// func (m *Metrics) IncrementErrors() {
//     m.errors.Add(1)
// }

// GetStats returns the current request and error counts.
// TODO: Implement thread-safe read of both counters.
//
// Note: Reading two atomic values separately is not atomic!
// - Another goroutine might update between our two reads
// - Result: Inconsistent snapshot (requests=100, errors=99, then errors increments)
// - For this use case, it's acceptable (eventual consistency)
// - For strict consistency, use mutex to protect both fields
//
// func (m *Metrics) GetStats() (requests int64, errors int64) {
//     return m.requests.Load(), m.errors.Load()
// }

// ============================================================================
// Exercise 8: Bank Account (Deposits and Withdrawals)
// ============================================================================

// BankAccount simulates a bank account with concurrent deposits/withdrawals.
// (Type definition is in types.go)

// Deposit adds money to the account.
// TODO: Implement thread-safe deposit.
//
// func (b *BankAccount) Deposit(amount int64) {
//     b.mu.Lock()
//     b.balance += amount
//     b.mu.Unlock()
// }

// Withdraw removes money from the account.
// TODO: Implement thread-safe withdrawal.
//
// Returns true if successful, false if insufficient funds.
//
// Critical section must be atomic:
// - Check balance >= amount
// - Subtract amount
// - Both operations under same lock!
//
// Bad (race):
//   if b.Balance() >= amount {  // Read
//       b.mu.Lock()
//       b.balance -= amount       // Write
//       b.mu.Unlock()
//   }
//   Problem: Another goroutine might withdraw between check and update
//
// Good:
//   b.mu.Lock()
//   if b.balance >= amount {
//       b.balance -= amount
//       b.mu.Unlock()
//       return true
//   }
//   b.mu.Unlock()
//   return false
//
// func (b *BankAccount) Withdraw(amount int64) bool {
//     b.mu.Lock()
//     defer b.mu.Unlock()
//
//     if b.balance >= amount {
//         b.balance -= amount
//         return true
//     }
//     return false
// }

// Balance returns the current balance.
// TODO: Implement thread-safe read.
//
// func (b *BankAccount) Balance() int64 {
//     b.mu.Lock()
//     balance := b.balance
//     b.mu.Unlock()
//     return balance
// }

// ============================================================================
// Exercise 9: Pipeline Pattern (Race-Free)
// ============================================================================

// Pipeline implements a concurrent pipeline: generate -> square -> sum
// TODO: Implement a race-free pipeline using channels.
//
// Should process all numbers concurrently and return the sum of squares.
//
// Channels are Go's primary synchronization mechanism:
// - Sending blocks if channel is full
// - Receiving blocks if channel is empty
// - Closing a channel signals "no more data"
// - Range over channel receives until closed
//
// Pattern:
// 1. Stage 1 (generator): Send numbers to channel, close when done
// 2. Stage 2 (transformer): Receive, transform, send to next channel, close
// 3. Stage 3 (collector): Receive and accumulate results
//
// func Pipeline(numbers []int) int {
//     // Stage 1: Generate numbers
//     gen := make(chan int)
//     go func() {
//         defer close(gen)  // Close when done sending
//         for _, n := range numbers {
//             gen <- n
//         }
//     }()
//
//     // Stage 2: Square numbers
//     sq := make(chan int)
//     go func() {
//         defer close(sq)
//         for n := range gen {  // Receive until gen is closed
//             sq <- n * n
//         }
//     }()
//
//     // Stage 3: Sum all squared numbers
//     sum := 0
//     for n := range sq {  // Receive until sq is closed
//         sum += n
//     }
//
//     return sum
// }

// ============================================================================
// Exercise 10: Worker Pool (Race-Free)
// ============================================================================

// WorkerPool processes jobs using a fixed number of workers.
// TODO: Implement a race-free worker pool.
//
// Parameters:
//   - numWorkers: number of concurrent workers
//   - jobs: slice of jobs to process
//   - process: function to process each job
//
// Returns: slice of results (order doesn't matter)
//
// Pattern:
// 1. Create jobs channel (buffered for all jobs)
// 2. Create results channel (buffered for all results)
// 3. Start numWorkers goroutines reading from jobs channel
// 4. Send all jobs to jobs channel, close it
// 5. Wait for workers to finish, close results channel
// 6. Collect all results from results channel
//
// func WorkerPool(numWorkers int, jobs []int, process func(int) int) []int {
//     jobsCh := make(chan int, len(jobs))
//     resultsCh := make(chan int, len(jobs))
//
//     // Start workers
//     var wg sync.WaitGroup
//     for i := 0; i < numWorkers; i++ {
//         wg.Add(1)
//         go func() {
//             defer wg.Done()
//             for job := range jobsCh {  // Process until channel closed
//                 resultsCh <- process(job)
//             }
//         }()
//     }
//
//     // Send all jobs
//     for _, job := range jobs {
//         jobsCh <- job
//     }
//     close(jobsCh)  // Signal no more jobs
//
//     // Close results channel when all workers finish
//     go func() {
//         wg.Wait()
//         close(resultsCh)
//     }()
//
//     // Collect results
//     results := make([]int, 0, len(jobs))
//     for result := range resultsCh {
//         results = append(results, result)
//     }
//
//     return results
// }

// After implementing all functions:
// - Run: go test -race ./...  (CRITICAL: always test with -race!)
// - The race detector instruments your code to find data races
// - Only catches races that execute during tests
// - Check: go test -v for verbose output
// - Compare with solution.go to see detailed explanations
