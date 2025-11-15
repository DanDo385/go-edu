# Project 25: Atomic Counters vs Mutex - Lock-Free Synchronization

## What Is This Project About?

This project teaches **atomic operations**—lock-free primitives for managing shared state without mutexes. You'll learn:

1. **What atomic operations are** (CPU-level instructions, lock-free synchronization)
2. **sync/atomic package** (Add, CompareAndSwap, Load, Store, Swap operations)
3. **Memory ordering guarantees** (happens-before relationships, visibility)
4. **Atomic vs mutex performance** (when to use each, benchmarking)
5. **Common patterns** (counters, flags, state machines, lock-free data structures)
6. **Pitfalls and limitations** (ABA problem, complex state, composability)

By the end, you'll understand when atomic operations provide better performance than mutexes and how to use them correctly without introducing subtle bugs.

---

## The Fundamental Problem: Shared State Synchronization

### First Principles: The Race Condition

When multiple goroutines access shared memory without synchronization, you get **race conditions**:

```go
// DANGER: Race condition
var counter int

for i := 0; i < 1000; i++ {
    go func() {
        counter++  // Read-modify-write (NOT ATOMIC!)
    }()
}

// Expected: counter = 1000
// Actual:   counter = 500-900 (non-deterministic)
```

**Why it fails:**
```
counter++ is three operations:
1. LOAD:  temp = counter
2. ADD:   temp = temp + 1
3. STORE: counter = temp

Goroutine A              Goroutine B
LOAD (counter=0)
                         LOAD (counter=0)
ADD (temp=1)
                         ADD (temp=1)
STORE (counter=1)
                         STORE (counter=1)  ← Overwrites A's write!

Result: counter=1 (should be 2)
```

**Traditional solution: Mutex**
```go
var counter int
var mu sync.Mutex

for i := 0; i < 1000; i++ {
    go func() {
        mu.Lock()
        counter++
        mu.Unlock()
    }()
}

// Always: counter = 1000 (correct)
```

**The cost:**
- Mutex involves kernel syscalls (slow)
- Goroutines block waiting for lock (no parallelism)
- Context switching overhead
- ~50-100 ns per lock/unlock operation

**Better solution: Atomic operations**
```go
var counter int64

for i := 0; i < 1000; i++ {
    go func() {
        atomic.AddInt64(&counter, 1)  // Single CPU instruction
    }()
}

// Always: counter = 1000 (correct)
// Cost: ~5-10 ns (10x faster than mutex!)
```

---

## What Are Atomic Operations? (The Core Concept)

An **atomic operation** is an indivisible operation that appears to execute instantaneously from other threads' perspective.

### Hardware Support

Modern CPUs provide atomic instructions:

**x86/x64:**
- `LOCK XADD`: Atomic add
- `LOCK CMPXCHG`: Compare-and-swap
- `MOV`: Load/store with ordering guarantees
- `XCHG`: Atomic exchange

**ARM:**
- `LDXR/STXR`: Load-exclusive/store-exclusive
- `LDREX/STREX`: Older ARM versions
- `DMB`: Data memory barrier

**Key properties:**
1. **Atomicity**: Operation completes without interruption
2. **Visibility**: Changes are immediately visible to other cores
3. **Ordering**: Memory operations around atomics have defined order

### Go's sync/atomic Package

Go wraps CPU atomic instructions in a portable API:

```go
package atomic

// Add operations
func AddInt32(addr *int32, delta int32) int32
func AddInt64(addr *int64, delta int64) int64
func AddUint32(addr *uint32, delta uint32) uint32
func AddUint64(addr *uint64, delta uint64) uint64
func AddUintptr(addr *uintptr, delta uintptr) uintptr

// Load operations (atomic read)
func LoadInt32(addr *int32) int32
func LoadInt64(addr *int64) int64
func LoadPointer(addr *unsafe.Pointer) unsafe.Pointer

// Store operations (atomic write)
func StoreInt32(addr *int32, val int32)
func StoreInt64(addr *int64, val int64)
func StorePointer(addr *unsafe.Pointer, val unsafe.Pointer)

// Swap operations (exchange)
func SwapInt32(addr *int32, new int32) (old int32)
func SwapInt64(addr *int64, new int64) (old int64)

// Compare-and-swap (CAS)
func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)
func CompareAndSwapInt64(addr *int64, old, new int64) (swapped bool)
func CompareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool)
```

**Type restrictions:**
- Only works on specific types: int32, int64, uint32, uint64, uintptr, unsafe.Pointer
- NOT int, bool, float64, or custom types (use int64 and convert)
- Must pass **pointer** to variable (not the value)

---

## Atomic Operations Deep Dive

### 1. Add: Atomic Increment/Decrement

**Function:**
```go
func AddInt64(addr *int64, delta int64) int64
```

**What it does:**
1. Atomically adds `delta` to `*addr`
2. Returns the **new value** (after addition)
3. Equivalent to `*addr += delta; return *addr` (but atomic)

**Example:**
```go
var counter int64 = 0

// Increment by 1
newVal := atomic.AddInt64(&counter, 1)
// counter = 1, newVal = 1

// Decrement by 3
newVal = atomic.AddInt64(&counter, -3)
// counter = -2, newVal = -2

// Add 10
newVal = atomic.AddInt64(&counter, 10)
// counter = 8, newVal = 8
```

**Common patterns:**
```go
// Increment counter
atomic.AddInt64(&counter, 1)

// Decrement counter
atomic.AddInt64(&counter, -1)

// Get old value before increment
old := atomic.AddInt64(&counter, 1) - 1

// Check if counter hit limit
if atomic.AddInt64(&counter, 1) > limit {
    // Over limit!
}
```

### 2. Load: Atomic Read

**Function:**
```go
func LoadInt64(addr *int64) int64
```

**Why needed?**

Regular reads are NOT atomic for types > word size on some architectures:

```go
// On 32-bit systems:
var x int64 = 0x0000000100000002

// Goroutine A
x = 0x0000000300000004

// Goroutine B (WITHOUT atomic)
val := x  // Might read 0x0000000300000002 (torn read!)
```

**With atomic:**
```go
val := atomic.LoadInt64(&x)  // Always reads full 64 bits atomically
```

**When to use:**
- Reading shared int64/uint64 on 32-bit systems
- Ensuring visibility of writes from other goroutines
- Reading values that are atomically written

**Example:**
```go
var config int64

// Writer goroutine
atomic.StoreInt64(&config, newConfig)

// Reader goroutines
currentConfig := atomic.LoadInt64(&config)  // Always sees complete write
```

### 3. Store: Atomic Write

**Function:**
```go
func StoreInt64(addr *int64, val int64)
```

**What it does:**
1. Atomically writes `val` to `*addr`
2. Ensures visibility to other goroutines
3. No return value

**Example:**
```go
var flag int64

// Atomically set flag
atomic.StoreInt64(&flag, 1)

// Other goroutines will see this write
if atomic.LoadInt64(&flag) == 1 {
    // Flag is set
}
```

**Use cases:**
- Publishing configuration changes
- Setting flags/status values
- Updating shared state that's only read (rarely written)

### 4. Swap: Atomic Exchange

**Function:**
```go
func SwapInt64(addr *int64, new int64) (old int64)
```

**What it does:**
1. Atomically reads old value from `*addr`
2. Writes `new` to `*addr`
3. Returns the old value
4. Equivalent to: `old := *addr; *addr = new; return old` (but atomic)

**Example:**
```go
var state int64 = 1

old := atomic.SwapInt64(&state, 2)
// state = 2, old = 1

old = atomic.SwapInt64(&state, 0)
// state = 0, old = 2
```

**Use cases:**
- Acquiring a lock-free spinlock
- Resetting counters/flags while reading old value
- State transitions

**Example: Spinlock**
```go
var lock int64

// Acquire lock
for atomic.SwapInt64(&lock, 1) != 0 {
    // Lock was held (old value = 1), spin
    runtime.Gosched()  // Yield to other goroutines
}
// Lock acquired (old value was 0)

// Critical section...

// Release lock
atomic.StoreInt64(&lock, 0)
```

### 5. CompareAndSwap (CAS): The Building Block

**Function:**
```go
func CompareAndSwapInt64(addr *int64, old, new int64) (swapped bool)
```

**What it does:**
1. Reads current value from `*addr`
2. If current value == `old`, writes `new` to `*addr` and returns `true`
3. Otherwise, does nothing and returns `false`
4. The comparison and swap happen **atomically** (no race possible)

**Pseudo-code:**
```go
func CompareAndSwap(addr *int64, old, new int64) bool {
    // ALL OF THIS HAPPENS ATOMICALLY:
    current := *addr
    if current == old {
        *addr = new
        return true
    }
    return false
}
```

**Example:**
```go
var counter int64 = 0

// Try to increment from 0 to 1
if atomic.CompareAndSwapInt64(&counter, 0, 1) {
    fmt.Println("Success: counter was 0, now 1")
} else {
    fmt.Println("Failed: counter was not 0")
}

// Try again (will fail because counter is now 1)
if atomic.CompareAndSwapInt64(&counter, 0, 1) {
    fmt.Println("Success")
} else {
    fmt.Println("Failed: counter is", atomic.LoadInt64(&counter))
    // Output: Failed: counter is 1
}
```

**Use cases:**
- Implementing lock-free algorithms
- Conditional updates (only update if value hasn't changed)
- Retry loops for complex operations

**Pattern: CAS loop for complex operations**
```go
// Atomically double the counter (can't use AddInt64 for this)
for {
    old := atomic.LoadInt64(&counter)
    new := old * 2
    if atomic.CompareAndSwapInt64(&counter, old, new) {
        break  // Success!
    }
    // Failed (another goroutine modified counter), retry
}
```

---

## Memory Ordering and Happens-Before

### The Memory Ordering Problem

Modern CPUs and compilers **reorder** instructions for performance:

```go
// Original code
x = 1
y = 2

// CPU might execute as:
y = 2
x = 1
```

**This breaks synchronization:**
```go
var data int
var ready int64

// Writer goroutine
data = 42        // Write 1
atomic.StoreInt64(&ready, 1)  // Write 2

// Reader goroutine
if atomic.LoadInt64(&ready) == 1 {  // Read 2
    fmt.Println(data)  // Read 1 - might see 0!
}
```

**Without ordering guarantees**, the reader might see `ready=1` but `data=0` because:
1. CPU reordered writes (Write 2 before Write 1)
2. CPU cached `data` (didn't see the write yet)

### Go's Memory Model Guarantees

Go's atomic operations provide **happens-before** guarantees:

**Rule 1: Atomic writes are visible to atomic reads**
```go
var x int64

// Goroutine A
atomic.StoreInt64(&x, 1)

// Goroutine B
if atomic.LoadInt64(&x) == 1 {
    // Definitely saw the store
}
```

**Rule 2: Atomic operations synchronize non-atomic operations**
```go
var data int
var ready int64

// Writer
data = 42                         // (1) Regular write
atomic.StoreInt64(&ready, 1)      // (2) Atomic write

// Reader
if atomic.LoadInt64(&ready) == 1 {  // (3) Atomic read
    fmt.Println(data)                // (4) Regular read - ALWAYS sees 42!
}
```

**Why it works:**
- (1) happens-before (2) [program order in same goroutine]
- (2) happens-before (3) [atomic store synchronizes with atomic load]
- (3) happens-before (4) [program order in same goroutine]
- **Transitive:** (1) happens-before (4)

**This is called a "release-acquire" pattern:**
- Atomic store = **release** (publishes all prior writes)
- Atomic load = **acquire** (observes all prior writes)

### Comparison: Atomic vs Mutex Ordering

**Both provide the same happens-before guarantees:**

```go
// ATOMIC VERSION
var data int
var ready int64

data = 42
atomic.StoreInt64(&ready, 1)  // Release

if atomic.LoadInt64(&ready) == 1 {  // Acquire
    // Guaranteed to see data = 42
}

// MUTEX VERSION
var data int
var mu sync.Mutex

mu.Lock()
data = 42
mu.Unlock()  // Release

mu.Lock()    // Acquire
fmt.Println(data)  // Guaranteed to see data = 42
mu.Unlock()
```

**Difference:** Atomic operations are faster (no lock contention, no syscalls).

---

## Atomic vs Mutex: When to Use Each

### Performance Comparison

**Benchmark results (typical):**
```
BenchmarkMutexCounter       50,000,000    25 ns/op
BenchmarkAtomicCounter     200,000,000     7 ns/op

BenchmarkMutexCounter-8     20,000,000    80 ns/op  (high contention)
BenchmarkAtomicCounter-8   100,000,000    15 ns/op  (scales better)
```

**Key insights:**
1. **Atomic is 3-5x faster** than mutex for simple operations
2. **Atomic scales better** under contention (multiple goroutines)
3. **Mutex becomes slower** as contention increases (lock waiting)

### When to Use Atomic Operations

**Use atomic when:**

1. **Simple counters/flags**
   ```go
   var requestCount int64
   atomic.AddInt64(&requestCount, 1)
   ```

2. **Single-value state**
   ```go
   var connected int64  // 0 = disconnected, 1 = connected
   atomic.StoreInt64(&connected, 1)
   ```

3. **Lock-free data structures**
   ```go
   // Lock-free stack using CAS
   for {
       old := atomic.LoadPointer(&head)
       newNode.next = old
       if atomic.CompareAndSwapPointer(&head, old, newNode) {
           break
       }
   }
   ```

4. **High-contention scenarios**
   - Many goroutines incrementing a counter
   - Frequent reads, rare writes
   - Avoiding lock contention

5. **Simple state machines**
   ```go
   const (
       StateInit = iota
       StateRunning
       StateStopped
   )
   var state int64
   atomic.StoreInt64(&state, StateRunning)
   ```

### When to Use Mutexes

**Use mutex when:**

1. **Multiple related fields** (can't update atomically)
   ```go
   type Stats struct {
       mu      sync.Mutex
       count   int
       total   int64
       average float64
   }

   // Must update all three together
   s.mu.Lock()
   s.count++
   s.total += value
   s.average = float64(s.total) / float64(s.count)
   s.mu.Unlock()
   ```

2. **Complex operations**
   ```go
   // Can't express this with atomics alone
   mu.Lock()
   if balance >= amount {
       balance -= amount
       transactions = append(transactions, tx)
   }
   mu.Unlock()
   ```

3. **Composition** (multiple operations must be atomic)
   ```go
   mu.Lock()
   delete(map1, key)
   map2[key] = value
   counter++
   mu.Unlock()
   ```

4. **Read-write patterns** (use sync.RWMutex)
   ```go
   // Many readers, few writers
   rwmu.RLock()
   value := config.Get(key)
   rwmu.RUnlock()
   ```

5. **Non-integer types**
   - Slices, maps, structs, strings
   - Floats (no atomic operations for float64 in Go)
   - Booleans (can use int32, but mutex is clearer)

### Decision Matrix

| Scenario | Use Atomic | Use Mutex |
|----------|-----------|-----------|
| **Single counter** | ✓ | - |
| **Single flag/state** | ✓ | - |
| **Multiple related fields** | - | ✓ |
| **Complex conditions** | - | ✓ |
| **Map/slice updates** | - | ✓ |
| **High contention** | ✓ | - |
| **Lock-free algorithms** | ✓ | - |
| **Code clarity** | - | ✓ (often) |
| **Float64 operations** | - | ✓ |
| **Pointer updates** | ✓ (with CAS) | ✓ |

**Rule of thumb:**
- If you can express the operation with a single atomic function → use atomic
- If you need multiple steps to be atomic → use mutex

---

## Common Patterns

### Pattern 1: Reference Counting

```go
type Resource struct {
    refCount int64
    // ...
}

func (r *Resource) Acquire() {
    atomic.AddInt64(&r.refCount, 1)
}

func (r *Resource) Release() {
    if atomic.AddInt64(&r.refCount, -1) == 0 {
        // Last reference, clean up
        r.cleanup()
    }
}
```

### Pattern 2: Configuration Updates

```go
type Config struct {
    value atomic.Value  // Can store any type
}

func (c *Config) Update(newConfig map[string]string) {
    c.value.Store(newConfig)
}

func (c *Config) Get() map[string]string {
    return c.value.Load().(map[string]string)
}
```

### Pattern 3: Once Initialization (Manual)

```go
var initialized int64
var data []byte

func getData() []byte {
    if atomic.LoadInt64(&initialized) == 0 {
        if atomic.CompareAndSwapInt64(&initialized, 0, 1) {
            data = loadExpensiveData()
        }
    }
    return data
}
```

### Pattern 4: Lock-Free Stack

```go
type Node struct {
    value int
    next  unsafe.Pointer  // *Node
}

var head unsafe.Pointer  // *Node

func Push(value int) {
    node := &Node{value: value}
    for {
        node.next = atomic.LoadPointer(&head)
        if atomic.CompareAndSwapPointer(&head, node.next, unsafe.Pointer(node)) {
            return
        }
    }
}

func Pop() (int, bool) {
    for {
        headPtr := atomic.LoadPointer(&head)
        if headPtr == nil {
            return 0, false
        }
        head := (*Node)(headPtr)
        next := atomic.LoadPointer(&head.next)
        if atomic.CompareAndSwapPointer(&head, headPtr, next) {
            return head.value, true
        }
    }
}
```

### Pattern 5: Rate Limiting with Atomic

```go
type RateLimiter struct {
    tokens    int64
    maxTokens int64
    lastRefill int64  // Unix timestamp
}

func (rl *RateLimiter) Allow() bool {
    now := time.Now().Unix()
    last := atomic.LoadInt64(&rl.lastRefill)

    // Refill tokens based on elapsed time
    if now > last {
        if atomic.CompareAndSwapInt64(&rl.lastRefill, last, now) {
            elapsed := now - last
            atomic.StoreInt64(&rl.tokens, min(rl.maxTokens, atomic.LoadInt64(&rl.tokens)+elapsed))
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
```

---

## Common Pitfalls

### Pitfall 1: Using Regular int Instead of int64

```go
// WRONG
var counter int  // Platform-dependent size
atomic.AddInt64(&counter, 1)  // COMPILE ERROR or RUNTIME PANIC

// CORRECT
var counter int64
atomic.AddInt64(&counter, 1)
```

### Pitfall 2: Mixing Atomic and Non-Atomic Access

```go
// WRONG: Race condition
var counter int64

counter++  // Non-atomic
atomic.AddInt64(&counter, 1)  // Atomic

// CORRECT: All accesses must be atomic
atomic.AddInt64(&counter, 1)
atomic.AddInt64(&counter, 1)
```

**Run with `-race` to detect this!**

### Pitfall 3: The ABA Problem

```go
// Lock-free stack with ABA problem
func Pop() {
    for {
        head := atomic.LoadPointer(&head)  // Read A
        // ... other goroutine: Pop A, Push B, Pop B, Push A (back to A)
        next := (*Node)(head).next
        if atomic.CompareAndSwapPointer(&head, head, next) {
            // SUCCESS, but we might have missed B!
        }
    }
}
```

**Solution:** Use versioned pointers or hazard pointers (complex).

### Pitfall 4: False Sharing

```go
// WRONG: counter1 and counter2 likely share a cache line
var counter1 int64
var counter2 int64

// Goroutines incrementing counter1 and counter2 fight over the same cache line!
```

**Solution:** Pad structures to avoid cache line contention:
```go
type Counter struct {
    value int64
    _     [56]byte  // Padding to fill 64-byte cache line
}
```

### Pitfall 5: Assuming Atomics Provide Full Locking

```go
// WRONG: This is NOT atomic (two separate operations)
if atomic.LoadInt64(&balance) >= amount {
    atomic.AddInt64(&balance, -amount)  // RACE: balance might have changed!
}

// CORRECT: Use CAS loop
for {
    old := atomic.LoadInt64(&balance)
    if old < amount {
        return errors.New("insufficient funds")
    }
    if atomic.CompareAndSwapInt64(&balance, old, old-amount) {
        return nil
    }
}
```

### Pitfall 6: Not Checking CAS Return Value

```go
// WRONG: Ignoring failure
atomic.CompareAndSwapInt64(&state, StateOld, StateNew)
// Continue assuming it succeeded (might not have!)

// CORRECT: Check and retry
for !atomic.CompareAndSwapInt64(&state, StateOld, StateNew) {
    // Retry or handle failure
}
```

---

## atomic.Value: Type-Safe Any

**Problem:** Atomic operations only work on integers and pointers.

**Solution:** `atomic.Value` for arbitrary types:

```go
type Value struct {
    // contains filtered or unexported fields
}

func (v *Value) Load() (val any)
func (v *Value) Store(val any)
func (v *Value) Swap(new any) (old any)
func (v *Value) CompareAndSwap(old, new any) (swapped bool)
```

**Example:**
```go
var config atomic.Value

// Store a map
config.Store(map[string]int{"timeout": 30})

// Load (type assert required)
cfg := config.Load().(map[string]int)
timeout := cfg["timeout"]

// Update
config.Store(map[string]int{"timeout": 60})
```

**Rules:**
1. **Must store consistent types** (can't store `int` then `string`)
2. **Type assertion required** on Load()
3. **Stores the pointer** (not a copy), so use immutable values
4. **Thread-safe** but not magic (still need to think about visibility)

**Best practice:**
```go
type Config struct {
    Timeout int
    MaxRetries int
}

var config atomic.Value

// Always store pointer to immutable struct
config.Store(&Config{Timeout: 30, MaxRetries: 3})

// Load
cfg := config.Load().(*Config)
```

---

## How to Run

```bash
# Run the demo program
cd minis/25-atomic-counters-vs-mutex
go run cmd/atomic-demo/main.go

# Run exercises
cd exercise
go test -v

# Run with race detector
go test -race -v

# Benchmark atomic vs mutex
go test -bench=. -benchmem

# Run a specific test
go test -v -run TestAtomicCounter
```

---

## Expected Output (Demo Program)

```
=== Atomic Operations Demonstration ===

=== 1. Atomic Add ===
Starting with counter = 0
Launching 1000 goroutines to increment...
Final counter: 1000 (correct!)

=== 2. Atomic CompareAndSwap ===
Initial state: 0 (Idle)
Worker 1: Attempting to start...
Worker 1: Success! Changed state 0 → 1
Worker 2: Attempting to start...
Worker 2: Failed (state already 1)

=== 3. Load/Store ===
Writer: Storing config version 1
Reader: Loaded config version 1
Writer: Storing config version 2
Reader: Loaded config version 2

=== 4. Atomic vs Mutex Benchmark ===
Testing 10,000,000 increments with 8 goroutines...
Atomic:  145 ms
Mutex:   523 ms
Speedup: 3.6x

=== 5. Lock-Free Stack ===
Pushing 1, 2, 3...
Popping: 3, 2, 1 (LIFO order)
```

---

## Key Takeaways

1. **Atomic operations are lock-free** CPU instructions (faster than mutexes)
2. **Use int64/uint64/uintptr** (not int, bool, float)
3. **Five core operations:** Add, Load, Store, Swap, CompareAndSwap
4. **Memory ordering guaranteed:** Atomic stores synchronize with atomic loads
5. **3-5x faster than mutex** for simple operations, scales better under contention
6. **Use atomic for:** single values, counters, flags, simple state machines
7. **Use mutex for:** multiple related fields, complex logic, non-integer types
8. **Pitfalls:** ABA problem, false sharing, mixing atomic/non-atomic access
9. **atomic.Value** for type-safe storage of arbitrary types
10. **Always use `-race` detector** to catch concurrency bugs

---

## Connections to Other Projects

- **Project 18 (goroutines-1M-demo)**: Goroutines share state via atomics
- **Project 21 (race-detection-demo)**: Race detector catches non-atomic access
- **Project 24 (sync-mutex-vs-rwmutex)**: Mutexes for complex synchronization
- **Project 22 (worker-pool-with-backpressure)**: Atomics for task counters
- **Project 26 (sync-once-singleton)**: Uses atomic operations internally
- **Project 28 (pprof-cpu-mem-benchmarks)**: Benchmark atomic vs mutex

---

## Stretch Goals

1. **Implement a lock-free queue** using CAS operations
2. **Measure cache line effects** with padded vs unpadded counters
3. **Build a spinlock** using atomic.SwapInt64 and compare to sync.Mutex
4. **Create a lock-free LRU cache** with atomic reference counting
5. **Visualize ABA problem** with a reproducible test case
