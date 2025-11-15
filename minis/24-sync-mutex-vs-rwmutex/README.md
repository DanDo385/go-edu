# Project 24: Sync Mutex vs RWMutex

## 1. What Is This About?

### Real-World Scenario

Imagine you're building a web server that caches user sessions in memory. Your cache is accessed by thousands of concurrent requests:

**❌ Without synchronization (data race):**
```go
type Cache struct {
    data map[string]string
}

func (c *Cache) Get(key string) string {
    return c.data[key] // RACE: Multiple goroutines reading
}

func (c *Cache) Set(key, value string) {
    c.data[key] = value // RACE: Write conflicts with reads/writes
}
```

**Problem**: When multiple goroutines access the map concurrently, you get:
- **Data corruption**: Map internals become inconsistent
- **Crashes**: `fatal error: concurrent map read and map write`
- **Undefined behavior**: Lost updates, stale reads

**✅ With Mutex (exclusive locking):**
```go
type Cache struct {
    mu   sync.Mutex
    data map[string]string
}

func (c *Cache) Get(key string) string {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.data[key]
}

func (c *Cache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}
```

**Problem solved**: No more races! But there's a bottleneck:
- **Every read blocks every other read** (even though reads are safe!)
- With 1000 concurrent reads, they execute **sequentially** (very slow)

**✅✅ With RWMutex (shared reads, exclusive writes):**
```go
type Cache struct {
    mu   sync.RWMutex  // Read-Write Mutex
    data map[string]string
}

func (c *Cache) Get(key string) string {
    c.mu.RLock()  // Read lock (shared)
    defer c.mu.RUnlock()
    return c.data[key]
}

func (c *Cache) Set(key, value string) {
    c.mu.Lock()   // Write lock (exclusive)
    defer c.mu.Unlock()
    c.data[key] = value
}
```

**Result**:
- **Multiple reads execute concurrently** (no blocking each other)
- **Writes still block everything** (ensures consistency)
- Performance boost: 10-100x faster for read-heavy workloads!

This project teaches you **Go's synchronization primitives** for protecting shared memory:
- **sync.Mutex**: Exclusive lock (only one goroutine at a time)
- **sync.RWMutex**: Read-write lock (many readers OR one writer)
- **Lock contention**: What happens when goroutines compete for locks
- **Performance trade-offs**: When to use each lock type

### What You'll Learn

1. **Data races**: What they are and why they're dangerous
2. **sync.Mutex**: Exclusive locking for mutual exclusion
3. **sync.RWMutex**: Optimized for read-heavy workloads
4. **Lock contention**: Performance impact of competing goroutines
5. **Deadlocks**: How to avoid circular dependencies
6. **Performance benchmarking**: Measuring throughput under contention
7. **Best practices**: When to use Mutex vs RWMutex vs channels

### The Challenge

Build thread-safe data structures that demonstrate:
- Protecting shared state with Mutex
- Optimizing read-heavy workloads with RWMutex
- Measuring performance differences with benchmarks
- Avoiding deadlocks and race conditions
- Choosing the right synchronization primitive

---

## 2. First Principles: Understanding Race Conditions

### What is a Data Race?

A **data race** occurs when:
1. Two or more goroutines access the same memory
2. At least one access is a write
3. The accesses are not synchronized

**Example**:
```go
var counter int

func increment() {
    counter++  // RACE: Read-modify-write is not atomic
}

func main() {
    go increment()
    go increment()
    // Result: undefined! Could be 1 or 2
}
```

**Why is `counter++` not safe?**

It's actually three operations:
```go
// What Go does internally:
tmp := counter     // 1. Read
tmp = tmp + 1      // 2. Modify
counter = tmp      // 3. Write
```

**Concurrent execution:**
```
Goroutine 1          Goroutine 2          counter
-----------          -----------          -------
tmp1 = counter (0)                        0
                     tmp2 = counter (0)   0
tmp1 = tmp1 + 1                           0
                     tmp2 = tmp2 + 1      0
counter = tmp1 (1)                        1
                     counter = tmp2 (1)   1  ← Lost one increment!
```

**Result**: Counter is 1 instead of 2. One increment was **lost**.

### Why Go's Map is Not Thread-Safe

Maps are complex data structures:
- Hash collisions require linked lists
- Resizing requires copying all entries
- Internal pointers change during operations

**What happens with concurrent map access:**

```go
// Goroutine 1
map[key] = value1

// Goroutine 2
map[key] = value2

// Goroutine 3
v := map[key]  // May read corrupted pointer!
```

If a read happens during a write:
- Read may see **half-updated** internal state
- Pointers may be **invalid**
- Program **crashes** or returns **garbage**

**Go's solution**: Detect races and panic:
```
fatal error: concurrent map writes
```

This is a **failsafe**, not a fix. You must use synchronization.

### Detecting Races

Go has a built-in **race detector**:

```bash
# Run with race detector
go test -race ./...
go run -race main.go

# Build with race detector
go build -race
```

**Output when race detected:**
```
==================
WARNING: DATA RACE
Write at 0x00c000014098 by goroutine 7:
  main.increment()
      /path/to/main.go:10 +0x3c

Previous write at 0x00c000014098 by goroutine 6:
  main.increment()
      /path/to/main.go:10 +0x3c
==================
```

**Key insight**: The race detector uses **happens-before** analysis to detect unsynchronized access.

---

## 3. First Principles: Understanding Mutex

### What is a Mutex?

**Mutex** = **Mut**ual **Ex**clusion

It's a lock that ensures **only one goroutine** can access protected code at a time.

**Analogy**: A mutex is like a **bathroom lock**:
- Only one person can use the bathroom at a time
- Others must wait in line
- When you're done, you unlock (next person can enter)

**Core operations:**

```go
var mu sync.Mutex

mu.Lock()      // Acquire lock (wait if locked)
// Critical section (only one goroutine here)
mu.Unlock()    // Release lock
```

### How Mutex Works Internally

**Simple conceptual model:**

```go
type Mutex struct {
    locked bool
    queue  []goroutine  // Waiting goroutines
}

func (m *Mutex) Lock() {
    if m.locked {
        // Add self to queue
        // Sleep until unlocked
    }
    m.locked = true
}

func (m *Mutex) Unlock() {
    m.locked = false
    // Wake up next goroutine in queue
}
```

**Actual implementation** is more complex:
- Uses **atomic compare-and-swap** (CAS) for fast-path
- **Semaphore** for blocking goroutines
- **Barging vs fairness** tradeoffs
- **Starvation prevention** (Go 1.9+)

### Mutex States

A Mutex can be in different states:

```
         Lock()
Unlocked ─────────► Locked
                      │
                      │ Unlock()
                      ▼
                   Unlocked
```

With contention:

```
         Lock() (first goroutine)
Unlocked ─────────────────────────► Locked
                                      │
         Lock() (second goroutine)    │
         ───────────────────────────► Locked (queued)
                                      │
                                      │ Unlock()
                                      ▼
                                    Locked (second goroutine runs)
                                      │
                                      │ Unlock()
                                      ▼
                                    Unlocked
```

### Mutex Example: Thread-Safe Counter

```go
type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    c.value++
    c.mu.Unlock()
}

func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}
```

**Why this is safe:**
- `Increment()` locks before modifying
- `Value()` locks before reading
- No two goroutines can execute critical section simultaneously

**Pattern: Always use `defer` for unlock:**

```go
func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()  // Ensures unlock even if panic
    return c.value
}
```

**Why defer?**
- If function panics, lock is still released
- If function has multiple return paths, you don't need `Unlock()` at each

---

## 4. First Principles: Understanding RWMutex

### What is RWMutex?

**RWMutex** = **R**ead-**W**rite Mutex

It's a lock with **two modes**:
1. **Read lock** (shared): Multiple goroutines can hold it
2. **Write lock** (exclusive): Only one goroutine can hold it

**Key rules:**
- Multiple read locks can be held simultaneously
- Write lock excludes all other locks (read or write)
- Readers block writers, writers block readers

**Analogy**: RWMutex is like a **library**:
- **Reading**: Many people can read books simultaneously
- **Writing**: Only one person can update the catalog (everyone else must wait)

**Core operations:**

```go
var mu sync.RWMutex

// Read operations
mu.RLock()      // Acquire read lock (shared)
// Read-only critical section
mu.RUnlock()    // Release read lock

// Write operations
mu.Lock()       // Acquire write lock (exclusive)
// Read-write critical section
mu.Unlock()     // Release write lock
```

### How RWMutex Works Internally

**Simplified model:**

```go
type RWMutex struct {
    readerCount  int       // Number of active readers
    writerWaiting bool     // Writer is waiting
    writerSem    semaphore // Writer semaphore
    readerSem    semaphore // Reader semaphore
}

func (m *RWMutex) RLock() {
    if m.writerWaiting {
        // Block until writer finishes
    }
    m.readerCount++
}

func (m *RWMutex) RUnlock() {
    m.readerCount--
    if m.readerCount == 0 && m.writerWaiting {
        // Wake up writer
    }
}

func (m *RWMutex) Lock() {
    m.writerWaiting = true
    // Wait for all readers to finish
    // Wait for previous writer to finish
}

func (m *RWMutex) Unlock() {
    m.writerWaiting = false
    // Wake up readers or next writer
}
```

**Actual implementation** uses:
- **Atomic operations** for fast-path
- **Two semaphores** (one for readers, one for writer)
- **Reader count tracking** with atomics
- **Writer preference** to prevent starvation

### RWMutex States

```
Initial state:
  Unlocked (no readers, no writer)

Scenario 1: Readers first
  Unlocked ─RLock()─► 1 Reader ─RLock()─► 2 Readers ─RLock()─► 3 Readers
                                                                    │
                      1 Reader ◄─RUnlock()─ 2 Readers ◄─RUnlock()─ ┘
                         │
                         │ RUnlock()
                         ▼
                      Unlocked

Scenario 2: Writer blocks readers
  Unlocked ─Lock()─► Writer Active
                         │
                         │ (RLock() attempts block)
                         │
                         │ Unlock()
                         ▼
                      Unlocked ─RLock()─► Queued readers run

Scenario 3: Readers block writer
  Unlocked ─RLock()─► 1 Reader ─RLock()─► 2 Readers
                                              │
                                              │ Lock() (blocks!)
                                              │
                                              │ RUnlock() (still 1 reader)
                                              │
                                              │ RUnlock() (no readers)
                                              ▼
                                           Writer Active
```

### RWMutex Example: Thread-Safe Cache

```go
type Cache struct {
    mu    sync.RWMutex
    items map[string]string
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()  // Multiple Gets can run concurrently
    defer c.mu.RUnlock()
    value, ok := c.items[key]
    return value, ok
}

func (c *Cache) Set(key, value string) {
    c.mu.Lock()   // Exclusive access for writes
    defer c.mu.Unlock()
    c.items[key] = value
}

func (c *Cache) Delete(key string) {
    c.mu.Lock()   // Exclusive access for deletes
    defer c.mu.Unlock()
    delete(c.items, key)
}
```

**Why RWMutex here?**
- Caches are **read-heavy** (90%+ reads)
- Multiple concurrent reads don't interfere with each other
- Only writes need exclusive access

---

## 5. Mutex vs RWMutex: Performance Comparison

### Throughput Comparison

**Scenario**: 100 goroutines, 1 million operations total

**Read-only workload (100% reads):**
```
Mutex:    100K ops/sec  (all reads serialize)
RWMutex:  5M ops/sec    (50x faster!)
Channels: 500K ops/sec
```

**Read-heavy workload (90% reads, 10% writes):**
```
Mutex:    200K ops/sec
RWMutex:  2M ops/sec    (10x faster!)
Channels: 400K ops/sec
```

**Balanced workload (50% reads, 50% writes):**
```
Mutex:    300K ops/sec
RWMutex:  400K ops/sec  (1.3x faster)
Channels: 300K ops/sec
```

**Write-heavy workload (10% reads, 90% writes):**
```
Mutex:    350K ops/sec
RWMutex:  300K ops/sec  (RWMutex slower!)
Channels: 250K ops/sec
```

**Key insight**: RWMutex has **overhead**:
- Tracking reader count
- Coordinating readers and writers
- More complex lock/unlock logic

For write-heavy workloads, this overhead outweighs the benefits.

### When to Use Mutex

**✅ Use Mutex when:**

1. **Write-heavy workload** (>30% writes)
   - RWMutex overhead not worth it
   - Example: Counters, accumulators

2. **Critical section is very short**
   - Lock contention is low anyway
   - Example: Incrementing a counter

3. **Simplicity matters**
   - Code is easier to understand with plain Mutex
   - No need to distinguish read vs write

4. **Uncertain read/write ratio**
   - Start with Mutex (simpler)
   - Optimize to RWMutex later if needed

**Example use cases:**
- Request counters
- Simple flags/state
- Short critical sections
- Configuration updates

### When to Use RWMutex

**✅ Use RWMutex when:**

1. **Read-heavy workload** (>70% reads)
   - Multiple readers can run concurrently
   - Example: Caches, configuration, lookup tables

2. **Critical section is long**
   - Reads take significant time (parsing, searching, etc.)
   - Serializing reads would be a bottleneck

3. **High concurrency**
   - Many goroutines competing for lock
   - Parallel reads significantly improve throughput

4. **Benchmark shows benefit**
   - Measure before optimizing
   - RWMutex not always faster (profile first!)

**Example use cases:**
- In-memory caches
- Configuration objects
- Routing tables
- Session stores
- Lookup tables

### Performance Gotchas

**RWMutex is NOT always faster:**

```go
// BAD: RWMutex overhead not worth it
type Counter struct {
    mu sync.RWMutex
    n  int
}

func (c *Counter) Inc() {
    c.mu.Lock()    // Very short critical section
    c.n++
    c.mu.Unlock()
}

func (c *Counter) Get() int {
    c.mu.RLock()   // Overhead of RLock not worth it
    defer c.mu.RUnlock()
    return c.n
}
```

**Why bad?**
- Critical section is **too short** (1 operation)
- RWMutex overhead is **higher** than the work itself
- Plain Mutex would be faster

**GOOD: Use atomic for counters:**

```go
type Counter struct {
    n atomic.Int64
}

func (c *Counter) Inc() {
    c.n.Add(1)  // Lock-free!
}

func (c *Counter) Get() int64 {
    return c.n.Load()  // Lock-free!
}
```

**Hierarchy of synchronization** (fastest to slowest):

```
1. No synchronization (if safe)
2. Atomic operations (lock-free)
3. RWMutex (for read-heavy)
4. Mutex (general purpose)
5. Channels (for coordination)
```

---

## 6. Lock Contention and Scalability

### What is Lock Contention?

**Lock contention** occurs when multiple goroutines compete for the same lock.

**Visualization:**

```
No contention (goroutines don't overlap):
G1: ████ (lock acquired immediately)
G2:      ████ (lock acquired immediately)
G3:           ████ (lock acquired immediately)
Time: ────────────────────►

High contention (goroutines overlap):
G1: ████████████ (holds lock long)
G2:     ⏸⏸⏸⏸⏸████ (waits, then acquires)
G3:     ⏸⏸⏸⏸⏸⏸⏸⏸⏸████ (waits even longer)
Time: ────────────────────►
```

**Effects of contention:**
- **Reduced throughput**: Goroutines spend time waiting
- **Increased latency**: Operations take longer
- **CPU waste**: Context switching between goroutines
- **Scalability problems**: Adding more goroutines doesn't help

### Measuring Contention

**Mutex contention metrics:**

```go
var mu sync.Mutex
var counter int64

func worker() {
    start := time.Now()
    mu.Lock()
    lockWaitTime := time.Since(start)  // How long we waited

    // Critical section
    counter++
    time.Sleep(1 * time.Millisecond)  // Simulate work

    mu.Unlock()

    if lockWaitTime > 5*time.Millisecond {
        log.Printf("High contention: waited %v", lockWaitTime)
    }
}
```

**Go's mutex profiling:**

```bash
# Build with mutex profiling
go test -mutexprofile=mutex.out

# View profile
go tool pprof mutex.out
```

### Reducing Contention

**Strategy 1: Shorten critical sections**

```go
// BAD: Lock held during expensive operation
func (c *Cache) Get(key string) string {
    c.mu.RLock()
    defer c.mu.RUnlock()

    value := c.data[key]

    // Expensive operation while holding lock!
    processed := expensiveProcessing(value)
    return processed
}

// GOOD: Release lock before expensive operation
func (c *Cache) Get(key string) string {
    c.mu.RLock()
    value := c.data[key]
    c.mu.RUnlock()  // Release early

    // Do expensive work without lock
    processed := expensiveProcessing(value)
    return processed
}
```

**Strategy 2: Shard the lock (partition data)**

```go
// BAD: Single lock for all keys
type Cache struct {
    mu   sync.RWMutex
    data map[string]string
}

// GOOD: Multiple locks (one per shard)
type ShardedCache struct {
    shards [256]*CacheShard  // 256 independent caches
}

type CacheShard struct {
    mu   sync.RWMutex
    data map[string]string
}

func (sc *ShardedCache) Get(key string) string {
    shard := sc.getShard(key)
    shard.mu.RLock()
    defer shard.mu.RUnlock()
    return shard.data[key]
}

func (sc *ShardedCache) getShard(key string) *CacheShard {
    hash := fnv.New32a()
    hash.Write([]byte(key))
    return sc.shards[hash.Sum32()%256]
}
```

**Benefits of sharding:**
- Reduces contention by **256x** (goroutines rarely compete for same shard)
- Scales with number of keys
- Trade-off: More memory, more complex code

**Strategy 3: Use RWMutex for read-heavy loads**

```go
// BAD: Mutex serializes all reads
type Config struct {
    mu     sync.Mutex
    values map[string]string
}

// GOOD: RWMutex allows concurrent reads
type Config struct {
    mu     sync.RWMutex
    values map[string]string
}

func (c *Config) Get(key string) string {
    c.mu.RLock()  // Many goroutines can hold this
    defer c.mu.RUnlock()
    return c.values[key]
}
```

**Strategy 4: Use atomic operations**

```go
// BAD: Mutex for simple counter
type Metrics struct {
    mu       sync.Mutex
    requests int64
}

func (m *Metrics) IncrementRequests() {
    m.mu.Lock()
    m.requests++
    m.mu.Unlock()
}

// GOOD: Atomic for simple counter
type Metrics struct {
    requests atomic.Int64
}

func (m *Metrics) IncrementRequests() {
    m.requests.Add(1)  // Lock-free!
}
```

**Strategy 5: Use sync.Map for concurrent maps**

```go
// BAD: Map with RWMutex (high contention)
type Registry struct {
    mu    sync.RWMutex
    items map[string]Item
}

// GOOD: sync.Map (optimized for concurrent access)
type Registry struct {
    items sync.Map  // Built-in concurrent map
}

func (r *Registry) Get(key string) (Item, bool) {
    value, ok := r.items.Load(key)
    if !ok {
        return Item{}, false
    }
    return value.(Item), true
}

func (r *Registry) Set(key string, item Item) {
    r.items.Store(key, item)
}
```

**When to use sync.Map:**
- Keys are only written once but read many times
- Multiple goroutines read, write, and overwrite different keys
- Benchmark shows benefit (not always faster!)

---

## 7. Deadlocks and How to Avoid Them

### What is a Deadlock?

**Deadlock** occurs when goroutines are waiting for each other in a cycle.

**Classic example: Two locks acquired in different orders**

```go
var mu1, mu2 sync.Mutex

// Goroutine 1
func transfer1to2() {
    mu1.Lock()
    // ... do work ...
    mu2.Lock()  // Acquire second lock
    // ... transfer ...
    mu2.Unlock()
    mu1.Unlock()
}

// Goroutine 2
func transfer2to1() {
    mu2.Lock()
    // ... do work ...
    mu1.Lock()  // Acquire first lock
    // ... transfer ...
    mu1.Unlock()
    mu2.Unlock()
}
```

**Deadlock scenario:**
```
Time  │ Goroutine 1        │ Goroutine 2
──────┼────────────────────┼────────────────────
t1    │ mu1.Lock() ✓       │
t2    │                    │ mu2.Lock() ✓
t3    │ mu2.Lock() ⏸      │ (G1 waits for mu2)
t4    │                    │ mu1.Lock() ⏸ (G2 waits for mu1)
t5+   │ ⏸ DEADLOCK ⏸     │ ⏸ DEADLOCK ⏸
```

Both goroutines wait forever.

### Preventing Deadlocks

**Rule 1: Lock ordering**

Always acquire locks in the **same order**:

```go
// GOOD: Consistent lock ordering
func transfer(from, to *Account) {
    // Always lock lower ID first
    first, second := from, to
    if from.ID > to.ID {
        first, second = to, from
    }

    first.mu.Lock()
    defer first.mu.Unlock()

    second.mu.Lock()
    defer second.mu.Unlock()

    // Transfer money
}
```

**Rule 2: Lock timeout**

Use `TryLock()` (Go 1.18+) with timeout:

```go
func transfer(from, to *Account) error {
    if !from.mu.TryLock() {
        return fmt.Errorf("could not lock source account")
    }
    defer from.mu.Unlock()

    if !to.mu.TryLock() {
        return fmt.Errorf("could not lock destination account")
    }
    defer to.mu.Unlock()

    // Transfer money
    return nil
}
```

**Rule 3: Avoid nested locks**

If possible, don't acquire multiple locks:

```go
// BAD: Nested locks
func (s *Store) UpdateBoth(id1, id2 string) {
    s.mu1.Lock()
    s.mu2.Lock()  // Risk of deadlock
    // ...
    s.mu2.Unlock()
    s.mu1.Unlock()
}

// GOOD: Single lock
func (s *Store) UpdateBoth(id1, id2 string) {
    s.mu.Lock()  // One lock protects both
    // ...
    s.mu.Unlock()
}
```

**Rule 4: Hold locks for shortest time**

```go
// BAD: Lock held during expensive operation
func (c *Cache) Process(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()

    value := c.data[key]
    result := expensiveOperation(value)  // Lock held too long!
    c.data[key] = result
}

// GOOD: Release lock before expensive operation
func (c *Cache) Process(key string) {
    c.mu.RLock()
    value := c.data[key]
    c.mu.RUnlock()

    result := expensiveOperation(value)  // No lock held

    c.mu.Lock()
    c.data[key] = result
    c.mu.Unlock()
}
```

### Detecting Deadlocks

**Go's deadlock detector**:

If all goroutines are blocked, Go panics:

```
fatal error: all goroutines are asleep - deadlock!
```

**However**, this only detects **global deadlocks** (all goroutines blocked).

**Partial deadlocks** (some goroutines blocked, others running) are **not detected**.

**Use timeouts to detect partial deadlocks:**

```go
func worker(mu *sync.Mutex) {
    done := make(chan struct{})

    go func() {
        mu.Lock()
        defer mu.Unlock()
        // ... work ...
        close(done)
    }()

    select {
    case <-done:
        // Success
    case <-time.After(5 * time.Second):
        log.Println("Possible deadlock: operation took too long")
    }
}
```

---

## 8. Common Patterns You Can Reuse

### Pattern 1: Thread-Safe Cache with RWMutex

```go
type Cache[K comparable, V any] struct {
    mu    sync.RWMutex
    items map[K]V
}

func NewCache[K comparable, V any]() *Cache[K, V] {
    return &Cache[K, V]{
        items: make(map[K]V),
    }
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    value, ok := c.items[key]
    return value, ok
}

func (c *Cache[K, V]) Set(key K, value V) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = value
}

func (c *Cache[K, V]) Delete(key K) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.items, key)
}

func (c *Cache[K, V]) Len() int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return len(c.items)
}
```

### Pattern 2: Sharded Map for High Concurrency

```go
type ShardedMap[K comparable, V any] struct {
    shards    []*Shard[K, V]
    shardMask uint32
}

type Shard[K comparable, V any] struct {
    mu    sync.RWMutex
    items map[K]V
}

func NewShardedMap[K comparable, V any](shardCount int) *ShardedMap[K, V] {
    // Ensure power of 2 for fast modulo
    shards := make([]*Shard[K, V], shardCount)
    for i := range shards {
        shards[i] = &Shard[K, V]{
            items: make(map[K]V),
        }
    }
    return &ShardedMap[K, V]{
        shards:    shards,
        shardMask: uint32(shardCount - 1),
    }
}

func (sm *ShardedMap[K, V]) getShard(key K) *Shard[K, V] {
    hash := sm.hash(key)
    return sm.shards[hash&sm.shardMask]
}

func (sm *ShardedMap[K, V]) Get(key K) (V, bool) {
    shard := sm.getShard(key)
    shard.mu.RLock()
    defer shard.mu.RUnlock()
    value, ok := shard.items[key]
    return value, ok
}

func (sm *ShardedMap[K, V]) Set(key K, value V) {
    shard := sm.getShard(key)
    shard.mu.Lock()
    defer shard.mu.Unlock()
    shard.items[key] = value
}
```

### Pattern 3: Thread-Safe Slice

```go
type SafeSlice[T any] struct {
    mu    sync.RWMutex
    items []T
}

func NewSafeSlice[T any]() *SafeSlice[T] {
    return &SafeSlice[T]{
        items: make([]T, 0),
    }
}

func (s *SafeSlice[T]) Append(item T) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.items = append(s.items, item)
}

func (s *SafeSlice[T]) Get(index int) (T, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var zero T
    if index < 0 || index >= len(s.items) {
        return zero, false
    }
    return s.items[index], true
}

func (s *SafeSlice[T]) Len() int {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return len(s.items)
}

func (s *SafeSlice[T]) Snapshot() []T {
    s.mu.RLock()
    defer s.mu.RUnlock()

    // Return a copy to prevent races
    snapshot := make([]T, len(s.items))
    copy(snapshot, s.items)
    return snapshot
}
```

### Pattern 4: Copy-on-Write for Read-Heavy Data

```go
type Config struct {
    mu     sync.RWMutex
    data   *ConfigData  // Pointer to immutable data
}

type ConfigData struct {
    Settings map[string]string
    Version  int
}

func (c *Config) Get(key string) string {
    c.mu.RLock()
    data := c.data  // Read pointer (fast)
    c.mu.RUnlock()

    // No lock needed - data is immutable
    return data.Settings[key]
}

func (c *Config) Update(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()

    // Create a copy
    newSettings := make(map[string]string)
    for k, v := range c.data.Settings {
        newSettings[k] = v
    }
    newSettings[key] = value

    // Atomically swap pointer
    c.data = &ConfigData{
        Settings: newSettings,
        Version:  c.data.Version + 1,
    }
}
```

### Pattern 5: Periodic Sync with Mutex

```go
type PeriodicWriter struct {
    mu     sync.Mutex
    buffer []string
    dirty  bool
}

func (pw *PeriodicWriter) Add(item string) {
    pw.mu.Lock()
    defer pw.mu.Unlock()
    pw.buffer = append(pw.buffer, item)
    pw.dirty = true
}

func (pw *PeriodicWriter) StartPeriodicFlush(ctx context.Context, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            pw.flush()
        case <-ctx.Done():
            pw.flush()  // Final flush
            return
        }
    }
}

func (pw *PeriodicWriter) flush() {
    pw.mu.Lock()
    if !pw.dirty {
        pw.mu.Unlock()
        return
    }

    // Copy buffer
    items := make([]string, len(pw.buffer))
    copy(items, pw.buffer)
    pw.buffer = pw.buffer[:0]
    pw.dirty = false
    pw.mu.Unlock()

    // Write without holding lock
    writeToFile(items)
}
```

---

## 9. Real-World Applications

### In-Memory Caches

**Use case**: Cache HTTP responses, database queries, computed results

```go
type HTTPCache struct {
    mu      sync.RWMutex
    entries map[string]*CacheEntry
}

type CacheEntry struct {
    Response  []byte
    ExpiresAt time.Time
}

func (c *HTTPCache) Get(url string) ([]byte, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    entry, ok := c.entries[url]
    if !ok || time.Now().After(entry.ExpiresAt) {
        return nil, false
    }
    return entry.Response, true
}

func (c *HTTPCache) Set(url string, response []byte, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.entries[url] = &CacheEntry{
        Response:  response,
        ExpiresAt: time.Now().Add(ttl),
    }
}
```

**Companies**: Redis, Memcached, CDNs (Cloudflare, Fastly)

### Configuration Management

**Use case**: Thread-safe application configuration

```go
type AppConfig struct {
    mu     sync.RWMutex
    config map[string]interface{}
}

func (ac *AppConfig) Get(key string) interface{} {
    ac.mu.RLock()
    defer ac.mu.RUnlock()
    return ac.config[key]
}

func (ac *AppConfig) Reload(filepath string) error {
    newConfig, err := loadConfigFile(filepath)
    if err != nil {
        return err
    }

    ac.mu.Lock()
    defer ac.mu.Unlock()
    ac.config = newConfig
    return nil
}
```

**Companies**: Kubernetes, Consul, etcd

### Metrics Collection

**Use case**: Thread-safe counters, gauges, histograms

```go
type Metrics struct {
    mu      sync.RWMutex
    counters map[string]*atomic.Int64
}

func (m *Metrics) Inc(name string) {
    m.mu.RLock()
    counter, ok := m.counters[name]
    m.mu.RUnlock()

    if !ok {
        // Create new counter
        m.mu.Lock()
        counter = &atomic.Int64{}
        m.counters[name] = counter
        m.mu.Unlock()
    }

    counter.Add(1)
}

func (m *Metrics) Snapshot() map[string]int64 {
    m.mu.RLock()
    defer m.mu.RUnlock()

    snapshot := make(map[string]int64)
    for name, counter := range m.counters {
        snapshot[name] = counter.Load()
    }
    return snapshot
}
```

**Companies**: Prometheus, Datadog, New Relic

### Session Management

**Use case**: Web session store

```go
type SessionStore struct {
    mu       sync.RWMutex
    sessions map[string]*Session
}

type Session struct {
    UserID    string
    Data      map[string]interface{}
    ExpiresAt time.Time
}

func (ss *SessionStore) Get(sessionID string) (*Session, bool) {
    ss.mu.RLock()
    defer ss.mu.RUnlock()

    session, ok := ss.sessions[sessionID]
    if !ok || time.Now().After(session.ExpiresAt) {
        return nil, false
    }
    return session, true
}

func (ss *SessionStore) Create(sessionID, userID string) {
    ss.mu.Lock()
    defer ss.mu.Unlock()

    ss.sessions[sessionID] = &Session{
        UserID:    userID,
        Data:      make(map[string]interface{}),
        ExpiresAt: time.Now().Add(30 * time.Minute),
    }
}
```

**Companies**: Express.js, Django, Rails

### Connection Pools

**Use case**: Database connection pooling

```go
type ConnectionPool struct {
    mu        sync.Mutex
    conns     []*Connection
    maxConns  int
    available chan *Connection
}

func (cp *ConnectionPool) Get() (*Connection, error) {
    select {
    case conn := <-cp.available:
        return conn, nil
    default:
        return cp.createConnection()
    }
}

func (cp *ConnectionPool) Put(conn *Connection) {
    select {
    case cp.available <- conn:
        // Returned to pool
    default:
        // Pool full, close connection
        conn.Close()
    }
}

func (cp *ConnectionPool) createConnection() (*Connection, error) {
    cp.mu.Lock()
    defer cp.mu.Unlock()

    if len(cp.conns) >= cp.maxConns {
        return nil, errors.New("connection pool exhausted")
    }

    conn := newConnection()
    cp.conns = append(cp.conns, conn)
    return conn, nil
}
```

**Companies**: PostgreSQL, MySQL drivers, HTTP client pools

---

## 10. Common Mistakes to Avoid

### Mistake 1: Copying Mutex

**❌ Wrong**:
```go
type Counter struct {
    mu    sync.Mutex
    value int
}

func (c Counter) Increment() {  // Value receiver copies mutex!
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}
```

**Problem**: `sync.Mutex` contains internal state. Copying it creates independent locks.

**✅ Correct**:
```go
func (c *Counter) Increment() {  // Pointer receiver
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}
```

**How to detect**: `go vet` will warn about copying locks.

### Mistake 2: Unlocking Unlocked Mutex

**❌ Wrong**:
```go
var mu sync.Mutex
mu.Unlock()  // PANIC: unlock of unlocked mutex
```

**Problem**: Unlocking an unlocked mutex panics.

**✅ Correct**: Always pair `Lock()` with `Unlock()`.

### Mistake 3: Forgetting defer Unlock

**❌ Wrong**:
```go
func (c *Cache) Get(key string) string {
    c.mu.RLock()

    if value, ok := c.data[key]; ok {
        return value  // BUG: Lock not released!
    }

    c.mu.RUnlock()
    return ""
}
```

**Problem**: Early return leaves lock held → deadlock.

**✅ Correct**:
```go
func (c *Cache) Get(key string) string {
    c.mu.RLock()
    defer c.mu.RUnlock()  // Always use defer

    if value, ok := c.data[key]; ok {
        return value
    }
    return ""
}
```

### Mistake 4: Lock Upgrade (RLock → Lock)

**❌ Wrong**:
```go
func (c *Cache) GetOrCreate(key string) string {
    c.mu.RLock()
    if value, ok := c.data[key]; ok {
        c.mu.RUnlock()
        return value
    }

    // BUG: Can't upgrade from RLock to Lock!
    c.mu.Lock()  // DEADLOCK: Already hold RLock
    c.data[key] = "new value"
    c.mu.Unlock()
    return "new value"
}
```

**Problem**: Can't upgrade read lock to write lock → deadlock.

**✅ Correct**:
```go
func (c *Cache) GetOrCreate(key string) string {
    // Try read lock first
    c.mu.RLock()
    value, ok := c.data[key]
    c.mu.RUnlock()

    if ok {
        return value
    }

    // Release read lock before acquiring write lock
    c.mu.Lock()
    defer c.mu.Unlock()

    // Check again (another goroutine may have created it)
    if value, ok := c.data[key]; ok {
        return value
    }

    c.data[key] = "new value"
    return "new value"
}
```

### Mistake 5: Holding Lock During Expensive Operation

**❌ Wrong**:
```go
func (c *Cache) Process(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()

    value := c.data[key]

    // BUG: Expensive operation while holding lock!
    result := expensiveComputation(value)  // Takes 1 second

    c.data[key] = result
}
```

**Problem**: Lock held during expensive operation → high contention.

**✅ Correct**:
```go
func (c *Cache) Process(key string) {
    c.mu.RLock()
    value := c.data[key]
    c.mu.RUnlock()  // Release lock early

    // Expensive operation without lock
    result := expensiveComputation(value)

    c.mu.Lock()
    c.data[key] = result
    c.mu.Unlock()
}
```

### Mistake 6: Not Using RWMutex for Read-Heavy Workloads

**❌ Suboptimal**:
```go
type Config struct {
    mu     sync.Mutex  // All reads serialize!
    values map[string]string
}

func (c *Config) Get(key string) string {
    c.mu.Lock()  // Blocks other reads
    defer c.mu.Unlock()
    return c.values[key]
}
```

**✅ Better**:
```go
type Config struct {
    mu     sync.RWMutex  // Concurrent reads
    values map[string]string
}

func (c *Config) Get(key string) string {
    c.mu.RLock()  // Doesn't block other reads
    defer c.mu.RUnlock()
    return c.values[key]
}
```

### Mistake 7: Returning Pointers to Locked Data

**❌ Dangerous**:
```go
func (c *Cache) GetEntry(key string) *Entry {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.data[key]  // BUG: Returns pointer to internal data!
}

// Caller can modify without lock!
entry := cache.GetEntry("key")
entry.Value = "modified"  // RACE CONDITION
```

**✅ Correct**:
```go
func (c *Cache) GetEntry(key string) Entry {
    c.mu.RLock()
    defer c.mu.RUnlock()

    // Return a copy
    if entry, ok := c.data[key]; ok {
        return *entry
    }
    return Entry{}
}
```

---

## 11. Benchmarking and Profiling

### Writing Benchmarks

```go
func BenchmarkMutexRead(b *testing.B) {
    c := &MutexCache{data: make(map[string]string)}
    c.Set("key", "value")

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            c.Get("key")
        }
    })
}

func BenchmarkRWMutexRead(b *testing.B) {
    c := &RWMutexCache{data: make(map[string]string)}
    c.Set("key", "value")

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            c.Get("key")
        }
    })
}
```

**Run benchmarks:**
```bash
go test -bench=. -benchmem
```

**Output:**
```
BenchmarkMutexRead-8        10000000    150 ns/op    0 B/op    0 allocs/op
BenchmarkRWMutexRead-8      50000000     30 ns/op    0 B/op    0 allocs/op
```

**Interpretation**:
- RWMutex is **5x faster** for reads (30ns vs 150ns)
- No allocations for either (good!)

### Profiling Lock Contention

```bash
# Run with mutex profiling
go test -mutexprofile=mutex.out -bench=.

# View profile
go tool pprof mutex.out

# Commands in pprof:
# top      - Show functions with most contention
# list     - Show source code with contention
# web      - Visualize as graph
```

---

## 12. Stretch Goals

### Goal 1: Implement LRU Cache with RWMutex ⭐⭐

Add eviction when cache is full.

### Goal 2: Build Sharded Counter for High Throughput ⭐⭐

Partition counter across multiple shards to reduce contention.

### Goal 3: Create Copy-on-Write Config Store ⭐⭐⭐

Immutable config updates without blocking reads.

### Goal 4: Benchmark Mutex vs RWMutex vs sync.Map ⭐⭐

Compare performance under different workloads.

### Goal 5: Detect and Report Lock Contention ⭐⭐⭐

Instrument locks to measure wait times.

---

## How to Run

```bash
# Run the demo
go run ./minis/24-sync-mutex-vs-rwmutex/cmd/mutex-demo/main.go

# Run tests
go test ./minis/24-sync-mutex-vs-rwmutex/...

# Run benchmarks
go test -bench=. -benchmem ./minis/24-sync-mutex-vs-rwmutex/...

# Run with race detector
go test -race ./minis/24-sync-mutex-vs-rwmutex/...

# Profile lock contention
go test -mutexprofile=mutex.out -bench=. ./minis/24-sync-mutex-vs-rwmutex/...
go tool pprof mutex.out
```

---

## Summary

**What you learned**:
- ✅ Data races occur when concurrent access lacks synchronization
- ✅ sync.Mutex provides exclusive locking (only one goroutine)
- ✅ sync.RWMutex allows concurrent reads, exclusive writes
- ✅ RWMutex is faster for read-heavy workloads (70%+ reads)
- ✅ Lock contention reduces throughput and scalability
- ✅ Sharding, atomic operations, and short critical sections reduce contention
- ✅ Always use pointer receivers and defer unlock

**Why this matters**:
Synchronization is fundamental to concurrent programming. Every production Go application uses mutexes:
- Web servers (session management, caching)
- Databases (transaction isolation)
- Message queues (buffer management)
- Microservices (configuration, metrics)

**Key rules**:
1. Use pointer receivers for mutex methods
2. Always `defer unlock()` to prevent leaks
3. RWMutex for read-heavy (70%+ reads), Mutex otherwise
4. Keep critical sections short
5. Use atomic operations for counters
6. Benchmark before optimizing

**Next steps**:
- Project 25: Channels and select (alternative to mutexes)
- Project 26: sync.WaitGroup and sync.Once (coordination primitives)
- Project 27: Context and cancellation (combining with locks)

Master mutexes, master concurrent Go! 🚀
