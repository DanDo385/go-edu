# Project 23: Bounded Channels as Semaphores - Resource Limiting Patterns

## What Is This Project About?

This project teaches **buffered channels as semaphores**—a fundamental pattern for controlling concurrent resource access in Go. You'll learn:

1. **What semaphores are** (counting vs binary, resource limiting)
2. **Buffered channels as semaphores** (idiomatic Go concurrency control)
3. **Concurrency limiting patterns** (max N concurrent operations)
4. **Resource pool patterns** (database connections, file handles, API rate limits)
5. **Weighted semaphores** (resources with different costs)
6. **Timeout and cancellation** (context-aware semaphore acquisition)
7. **Comparison with sync packages** (when to use channels vs sync.Semaphore)

By the end, you'll master one of Go's most elegant concurrency patterns: using simple buffered channels to elegantly solve complex resource limiting problems.

---

## The Fundamental Problem: Limiting Concurrent Access

### First Principles: What Is Resource Limiting?

Imagine you have:
- **100 goroutines** trying to access a resource
- **Only 10 concurrent connections allowed** (database, API, file handles)

**Without limiting:**
```go
// BAD: All 100 goroutines try to connect at once
for i := 0; i < 100; i++ {
    go func() {
        conn := db.Connect()  // Might exceed connection limit!
        defer conn.Close()
        query(conn)
    }()
}
// Result: Connection pool exhausted, errors, crashes
```

**With semaphore:**
```go
// GOOD: At most 10 concurrent connections
sem := make(chan struct{}, 10)  // Buffered channel = semaphore

for i := 0; i < 100; i++ {
    go func() {
        sem <- struct{}{}        // Acquire (blocks if 10 already acquired)
        defer func() { <-sem }() // Release

        conn := db.Connect()     // Now safe: max 10 concurrent
        defer conn.Close()
        query(conn)
    }()
}
```

**The pattern**: A buffered channel acts as a **counting semaphore** where:
- **Channel capacity** = maximum concurrent permits
- **Send to channel** = acquire permit (blocks when full)
- **Receive from channel** = release permit (makes space)

---

## What Is a Semaphore? (Computer Science Foundations)

### Semaphore Basics

A **semaphore** is a synchronization primitive with two operations:
1. **Acquire** (P, wait, down): Decrement counter, block if 0
2. **Release** (V, signal, up): Increment counter, wake waiting goroutines

**Types:**

1. **Binary Semaphore** (mutex)
   - Counter: 0 or 1
   - Use: Mutual exclusion (only 1 goroutine at a time)
   ```go
   sem := make(chan struct{}, 1)  // Binary semaphore

   sem <- struct{}{}  // Acquire (lock)
   // Critical section
   <-sem              // Release (unlock)
   ```

2. **Counting Semaphore**
   - Counter: 0 to N
   - Use: Limit concurrent access to N goroutines
   ```go
   sem := make(chan struct{}, 10)  // Counting semaphore (max 10)

   sem <- struct{}{}  // Acquire 1 permit
   // Do work with resource
   <-sem              // Release 1 permit
   ```

### Classical Semaphore Problems

**1. Bounded Buffer Problem (Producer-Consumer)**
```go
// Buffer with capacity 10
buffer := make(chan int, 10)
empty := 10   // Counting semaphore: empty slots
full := 0     // Counting semaphore: full slots

// Producer
sem <- struct{}{}  // Acquire empty slot
buffer <- item
<-fullSem         // Release (signal full slot)

// Consumer
fullSem <- struct{}{}  // Acquire full slot
item := <-buffer
<-sem                  // Release (signal empty slot)
```

**2. Readers-Writers Problem**
```go
// Max 5 concurrent readers, 1 writer
readers := make(chan struct{}, 5)
writers := make(chan struct{}, 1)

// Read
readers <- struct{}{}
defer func() { <-readers }()
read()

// Write
writers <- struct{}{}
defer func() { <-writers }()
write()
```

**3. Dining Philosophers Problem**
```go
// 5 forks (resources) for 5 philosophers
forks := make(chan struct{}, 5)

// Philosopher acquires 2 forks
forks <- struct{}{}  // Left fork
forks <- struct{}{}  // Right fork
eat()
<-forks  // Release left
<-forks  // Release right
```

---

## Problem 1: Buffered Channel as Semaphore (Core Pattern)

### Why Buffered Channels Make Perfect Semaphores

**Buffered channel properties map directly to semaphore semantics:**

| Semaphore Concept | Buffered Channel Equivalent |
|-------------------|----------------------------|
| Counter value     | Number of items in buffer  |
| Max count (N)     | Channel capacity           |
| Acquire           | Send to channel (`sem <- struct{}{}`) |
| Release           | Receive from channel (`<-sem`) |
| Block on acquire  | Send blocks when buffer full |
| Wake on release   | Receive makes space, unblocks sender |

**Key insight**: The channel **buffer** tracks available permits.
- **Full buffer** = no permits available (acquire blocks)
- **Empty buffer** = all permits available
- **Send** = take permit (decrements available)
- **Receive** = return permit (increments available)

### Basic Pattern: Limiting Concurrency

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    const (
        numTasks      = 20
        maxConcurrent = 5
    )

    // SEMAPHORE: Buffered channel with capacity = max concurrent
    sem := make(chan struct{}, maxConcurrent)

    for i := 1; i <= numTasks; i++ {
        // ACQUIRE: Send to channel (blocks if 5 already running)
        sem <- struct{}{}

        go func(id int) {
            // RELEASE: Always release on exit
            defer func() { <-sem }()

            fmt.Printf("Task %d: Started\n", id)
            time.Sleep(1 * time.Second)  // Simulate work
            fmt.Printf("Task %d: Complete\n", id)
        }(i)
    }

    // WAIT FOR ALL: Acquire all permits (ensures all goroutines done)
    for i := 0; i < maxConcurrent; i++ {
        sem <- struct{}{}
    }

    fmt.Println("All tasks complete!")
}
```

**Output pattern:**
```
Task 1: Started   ┐
Task 2: Started   │
Task 3: Started   ├─ First batch (5 concurrent)
Task 4: Started   │
Task 5: Started   ┘
Task 1: Complete  ← Release permits
Task 6: Started   ← Next task acquires released permit
Task 2: Complete
Task 7: Started
...
```

**Why `struct{}`?**
- **Zero bytes**: `struct{}{}` uses no memory
- **Signaling only**: We don't need data, just synchronization
- **Idiomatic**: Standard Go pattern for "signal only" channels

### Alternative: Empty Interface (More Flexible)

```go
sem := make(chan interface{}, 10)

// Acquire
sem <- nil

// Release
<-sem
```

Use when you might want to send data later (e.g., error results).

---

## Problem 2: Real-World Pattern - Database Connection Pool

### The Scenario

You have:
- **Database with max 10 connections**
- **100 concurrent requests**
- Need to **limit concurrent queries** to avoid exhausting connections

### Naive Approach (No Limiting)

```go
// BAD: Can exceed connection pool limit
func processRequests(requests []Request) {
    for _, req := range requests {
        go func(r Request) {
            conn := db.Connect()  // Might fail if pool exhausted!
            defer conn.Close()
            query(conn, r)
        }(req)
    }
}
```

**Problem**: If all 100 goroutines run simultaneously, you'll try to open 100 connections when only 10 are allowed.

### Semaphore Solution

```go
type DBPool struct {
    sem chan struct{}
}

func NewDBPool(maxConns int) *DBPool {
    return &DBPool{
        sem: make(chan struct{}, maxConns),
    }
}

func (p *DBPool) Acquire() {
    p.sem <- struct{}{}  // Blocks if maxConns already acquired
}

func (p *DBPool) Release() {
    <-p.sem
}

func (p *DBPool) Query(ctx context.Context, query string) (Result, error) {
    // Acquire connection permit
    p.Acquire()
    defer p.Release()

    // Now safe to open connection (guaranteed < maxConns)
    conn := db.Connect()
    defer conn.Close()

    return conn.Query(query)
}

// Usage
pool := NewDBPool(10)

for _, req := range requests {
    go func(r Request) {
        result, err := pool.Query(context.Background(), r.Query)
        if err != nil {
            log.Printf("Query failed: %v", err)
            return
        }
        process(result)
    }(req)
}
```

**How it works:**
1. **First 10 goroutines**: Acquire succeeds immediately (buffer has space)
2. **11th goroutine**: Acquire blocks (buffer full)
3. **First goroutine completes**: Calls Release, making space
4. **11th goroutine**: Unblocks, acquires permit, proceeds

### Enhanced Pattern: Context-Aware Acquisition

```go
func (p *DBPool) AcquireWithContext(ctx context.Context) error {
    select {
    case p.sem <- struct{}{}:
        return nil  // Acquired
    case <-ctx.Done():
        return ctx.Err()  // Context cancelled/timeout
    }
}

func (p *DBPool) QueryWithTimeout(query string, timeout time.Duration) (Result, error) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    // Acquire with timeout
    if err := p.AcquireWithContext(ctx); err != nil {
        return nil, fmt.Errorf("acquire timeout: %w", err)
    }
    defer p.Release()

    conn := db.Connect()
    defer conn.Close()

    return conn.QueryContext(ctx, query)
}
```

**Benefits:**
- ✅ **Timeout**: Don't wait forever for a permit
- ✅ **Cancellation**: Respect context cancellation
- ✅ **Graceful degradation**: Fail fast if resource unavailable

---

## Problem 3: Rate Limiting with Semaphores

### Token Bucket Rate Limiter

**Goal**: Allow at most N requests per time window.

```go
type RateLimiter struct {
    tokens chan struct{}
    rate   time.Duration
}

func NewRateLimiter(maxBurst int, rate time.Duration) *RateLimiter {
    rl := &RateLimiter{
        tokens: make(chan struct{}, maxBurst),
        rate:   rate,
    }

    // Fill initial tokens
    for i := 0; i < maxBurst; i++ {
        rl.tokens <- struct{}{}
    }

    // Refill tokens periodically
    go func() {
        ticker := time.NewTicker(rate)
        defer ticker.Stop()

        for range ticker.C {
            select {
            case rl.tokens <- struct{}{}:
                // Added token
            default:
                // Bucket full, skip
            }
        }
    }()

    return rl
}

func (rl *RateLimiter) Allow() bool {
    select {
    case <-rl.tokens:
        return true  // Token available
    default:
        return false  // No tokens (rate limited)
    }
}

func (rl *RateLimiter) Wait() {
    <-rl.tokens  // Block until token available
}

// Usage
limiter := NewRateLimiter(10, 100*time.Millisecond)

for i := 0; i < 100; i++ {
    go func(id int) {
        limiter.Wait()  // Wait for rate limit token
        makeAPIRequest(id)
    }(i)
}
```

**Behavior:**
- **Burst**: Allow up to 10 requests immediately
- **Sustained rate**: Refill 1 token every 100ms (10 req/sec)
- **Smoothing**: Spreads requests over time

### Leaky Bucket Pattern

```go
type LeakyBucket struct {
    sem chan struct{}
}

func NewLeakyBucket(rate time.Duration) *LeakyBucket {
    lb := &LeakyBucket{
        sem: make(chan struct{}, 1),
    }

    // Leak (allow) one request per rate period
    go func() {
        ticker := time.NewTicker(rate)
        defer ticker.Stop()

        for range ticker.C {
            select {
            case lb.sem <- struct{}{}:
            default:
            }
        }
    }()

    return lb
}

func (lb *LeakyBucket) Wait() {
    <-lb.sem  // Wait for next leak slot
}

// Usage: Exactly 1 request per 100ms
bucket := NewLeakyBucket(100 * time.Millisecond)

for i := 0; i < 50; i++ {
    bucket.Wait()
    makeRequest(i)
}
```

---

## Problem 4: Weighted Semaphores (Variable Resource Costs)

### The Problem

Sometimes resources have **different costs**:
- Small query: 1 connection
- Large batch query: 3 connections
- Report generation: 5 connections

**Goal**: Track total resource usage, not just count.

### Weighted Semaphore Pattern

```go
type WeightedSemaphore struct {
    permits chan struct{}
}

func NewWeightedSemaphore(maxWeight int) *WeightedSemaphore {
    return &WeightedSemaphore{
        permits: make(chan struct{}, maxWeight),
    }
}

func (ws *WeightedSemaphore) Acquire(weight int) {
    // Acquire 'weight' permits
    for i := 0; i < weight; i++ {
        ws.permits <- struct{}{}
    }
}

func (ws *WeightedSemaphore) Release(weight int) {
    // Release 'weight' permits
    for i := 0; i < weight; i++ {
        <-ws.permits
    }
}

// Usage
sem := NewWeightedSemaphore(10)  // Total capacity: 10

// Small task (cost: 1)
go func() {
    sem.Acquire(1)
    defer sem.Release(1)
    doSmallTask()
}()

// Large task (cost: 5)
go func() {
    sem.Acquire(5)
    defer sem.Release(5)
    doLargeTask()
}()
```

**Scenario:**
```
Capacity: 10 permits
[__________] (10 free)

Acquire(3):
[XXX_______] (7 free)

Acquire(5):
[XXXYYYYY__] (2 free)

Acquire(3): BLOCKS (only 2 free, need 3)

Release(3):
[___YYYYY__] (5 free)

Acquire(3): PROCEEDS
[ZZZYYYYY__] (2 free)
```

### Context-Aware Weighted Acquisition

```go
func (ws *WeightedSemaphore) AcquireWithContext(ctx context.Context, weight int) error {
    acquired := 0

    // Acquire permits one at a time (so we can check context)
    for i := 0; i < weight; i++ {
        select {
        case ws.permits <- struct{}{}:
            acquired++
        case <-ctx.Done():
            // Timeout/cancelled: release what we acquired
            for j := 0; j < acquired; j++ {
                <-ws.permits
            }
            return ctx.Err()
        }
    }

    return nil
}

// Usage with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := sem.AcquireWithContext(ctx, 7); err != nil {
    return fmt.Errorf("failed to acquire: %w", err)
}
defer sem.Release(7)
```

**Critical detail**: If acquisition fails partway (context cancelled), we must **release what we acquired** to avoid permit leaks.

---

## Problem 5: Try-Acquire Pattern (Non-Blocking)

### The Pattern

Sometimes you want to **try** to acquire without blocking:
- If permit available: Acquire and proceed
- If not available: Skip or fallback

```go
type Semaphore struct {
    sem chan struct{}
}

func (s *Semaphore) TryAcquire() bool {
    select {
    case s.sem <- struct{}{}:
        return true  // Acquired
    default:
        return false  // Full (would block)
    }
}

func (s *Semaphore) Acquire() {
    s.sem <- struct{}{}
}

func (s *Semaphore) Release() {
    <-s.sem
}

// Usage: Graceful degradation
sem := NewSemaphore(5)

if sem.TryAcquire() {
    defer sem.Release()
    processRequest()
} else {
    returnError(http.StatusTooManyRequests, "Server busy, try again")
}
```

### With Timeout

```go
func (s *Semaphore) AcquireTimeout(timeout time.Duration) bool {
    select {
    case s.sem <- struct{}{}:
        return true  // Acquired
    case <-time.After(timeout):
        return false  // Timeout
    }
}

// Usage: Wait up to 1 second
if sem.AcquireTimeout(1 * time.Second) {
    defer sem.Release()
    process()
} else {
    log.Println("Timeout waiting for resource")
}
```

---

## Problem 6: Worker Pool with Semaphore

### Classic Worker Pool Pattern

```go
type WorkerPool struct {
    sem  chan struct{}
    work chan func()
    done chan struct{}
}

func NewWorkerPool(maxWorkers int) *WorkerPool {
    wp := &WorkerPool{
        sem:  make(chan struct{}, maxWorkers),
        work: make(chan func(), 100),  // Buffered work queue
        done: make(chan struct{}),
    }

    // Fill semaphore (all workers available initially)
    for i := 0; i < maxWorkers; i++ {
        wp.sem <- struct{}{}
    }

    return wp
}

func (wp *WorkerPool) Submit(job func()) {
    wp.work <- job
}

func (wp *WorkerPool) Start() {
    go func() {
        for job := range wp.work {
            <-wp.sem  // Acquire worker

            go func(j func()) {
                defer func() { wp.sem <- struct{}{} }()  // Release worker
                j()
            }(job)
        }
        close(wp.done)
    }()
}

func (wp *WorkerPool) Stop() {
    close(wp.work)
    <-wp.done
}

// Usage
pool := NewWorkerPool(10)
pool.Start()

for i := 0; i < 100; i++ {
    pool.Submit(func() {
        processTask()
    })
}

pool.Stop()
```

**Alternative: Pre-allocated Workers**

```go
func NewWorkerPool(numWorkers int, jobQueue <-chan Job) {
    sem := make(chan struct{}, numWorkers)

    for job := range jobQueue {
        sem <- struct{}{}  // Acquire worker slot

        go func(j Job) {
            defer func() { <-sem }()  // Release worker slot
            process(j)
        }(job)
    }

    // Wait for all workers to finish
    for i := 0; i < numWorkers; i++ {
        sem <- struct{}{}
    }
}
```

---

## Comparison: Channels vs sync.Semaphore

### Go's Extended Package

Go provides `golang.org/x/sync/semaphore` for more advanced use cases:

```go
import "golang.org/x/sync/semaphore"

sem := semaphore.NewWeighted(10)

// Acquire
if err := sem.Acquire(ctx, 3); err != nil {
    // Context cancelled or error
}
defer sem.Release(3)
```

### When to Use Each

| Use Case | Buffered Channel | sync.Semaphore |
|----------|-----------------|----------------|
| **Simple counting semaphore** | ✅ Preferred (idiomatic) | Overkill |
| **Binary semaphore (mutex)** | ✅ Simple pattern | Use sync.Mutex instead |
| **Weighted semaphore** | Possible but clunky | ✅ Designed for this |
| **Context integration** | Manual select | ✅ Built-in |
| **Non-blocking try** | ✅ select with default | TryAcquire available |
| **Zero dependencies** | ✅ Built-in | Requires x/sync |

**Rule of thumb:**
- **Buffered channel**: Simple counting semaphore, idiomatic Go
- **sync.Semaphore**: Weighted semaphores, complex resource management
- **sync.Mutex**: Binary semaphore (mutual exclusion)

---

## Common Patterns and Idioms

### Pattern 1: Defer Release (Critical!)

```go
// ALWAYS use defer for release
sem <- struct{}{}
defer func() { <-sem }()

// Do work (might panic or return early)
// Release guaranteed to execute
```

### Pattern 2: Wait for All

```go
// Submit tasks
for i := 0; i < 100; i++ {
    sem <- struct{}{}
    go func() {
        defer func() { <-sem }()
        work()
    }()
}

// Wait for all to complete: acquire all permits
for i := 0; i < cap(sem); i++ {
    sem <- struct{}{}
}
```

### Pattern 3: Dynamic Capacity Adjustment

```go
type DynamicSemaphore struct {
    sem chan struct{}
    mu  sync.Mutex
}

func (ds *DynamicSemaphore) IncreaseCapacity(n int) {
    ds.mu.Lock()
    defer ds.mu.Unlock()

    // Add permits
    for i := 0; i < n; i++ {
        select {
        case ds.sem <- struct{}{}:
        default:
        }
    }
}

func (ds *DynamicSemaphore) DecreaseCapacity(n int) {
    ds.mu.Lock()
    defer ds.mu.Unlock()

    // Remove permits
    for i := 0; i < n; i++ {
        select {
        case <-ds.sem:
        default:
        }
    }
}
```

### Pattern 4: Semaphore Metrics

```go
type MonitoredSemaphore struct {
    sem     chan struct{}
    current atomic.Int64
    max     int64
}

func (ms *MonitoredSemaphore) Acquire() {
    ms.sem <- struct{}{}
    current := ms.current.Add(1)

    // Track peak usage
    if current > ms.max {
        ms.max = current
    }
}

func (ms *MonitoredSemaphore) Release() {
    <-ms.sem
    ms.current.Add(-1)
}

func (ms *MonitoredSemaphore) Usage() (current, capacity int) {
    return int(ms.current.Load()), cap(ms.sem)
}
```

---

## Common Mistakes and Pitfalls

### Mistake 1: Forgetting to Release

```go
// BAD: Permit leak if error occurs
sem <- struct{}{}
if err := doWork(); err != nil {
    return err  // Forgot to release!
}
<-sem

// GOOD: Defer ensures release
sem <- struct{}{}
defer func() { <-sem }()

if err := doWork(); err != nil {
    return err  // Release happens via defer
}
```

### Mistake 2: Deadlock from Acquiring Twice

```go
sem := make(chan struct{}, 1)

// BAD: Deadlock
sem <- struct{}{}
sem <- struct{}{}  // BLOCKS FOREVER (capacity 1)
<-sem
<-sem

// GOOD: Release before re-acquiring
sem <- struct{}{}
// Work
<-sem

sem <- struct{}{}
// More work
<-sem
```

### Mistake 3: Wrong Acquire/Release Order

```go
// BAD: Release before acquire
<-sem  // Receives garbage or blocks if empty
sem <- struct{}{}

// GOOD: Always acquire first
sem <- struct{}{}
defer func() { <-sem }()
```

### Mistake 4: Not Checking Context

```go
// BAD: Might wait forever
sem <- struct{}{}
defer func() { <-sem }()

// GOOD: Check context
select {
case sem <- struct{}{}:
    defer func() { <-sem }()
case <-ctx.Done():
    return ctx.Err()
}
```

### Mistake 5: Capacity Mismatch

```go
// BAD: Capacity doesn't match use case
sem := make(chan struct{}, 100)  // Too high: wastes resources
sem := make(chan struct{}, 1)    // Too low: excessive blocking

// GOOD: Match capacity to real limit
maxDBConns := 10
sem := make(chan struct{}, maxDBConns)
```

---

## Performance Characteristics

### Buffered Channel Semaphore

- **Acquire (send)**: O(1) when space available, blocks when full
- **Release (receive)**: O(1) always
- **Memory**: O(capacity) - each slot in buffer
- **Contention**: Optimized in Go runtime (fast path for uncontended)

### Comparison to Mutex

| Operation | Mutex | Buffered Channel (cap=1) |
|-----------|-------|-------------------------|
| Lock/Acquire | ~20 ns (uncontended) | ~25 ns (uncontended) |
| Unlock/Release | ~20 ns | ~25 ns |
| Contention | Slower (parking goroutines) | Slower (parking goroutines) |
| Use case | Mutual exclusion | Counting semaphore |

**Key insight**: For simple mutual exclusion, `sync.Mutex` is slightly faster. For counting semaphores, buffered channels are idiomatic and nearly as fast.

---

## Real-World Applications

### API Rate Limiter

```go
type APIClient struct {
    limiter *RateLimiter
    client  *http.Client
}

func (c *APIClient) Request(ctx context.Context, url string) (*http.Response, error) {
    if err := c.limiter.AcquireWithContext(ctx); err != nil {
        return nil, fmt.Errorf("rate limit: %w", err)
    }
    defer c.limiter.Release()

    return c.client.Get(url)
}
```

### Database Connection Pool

```go
type DB struct {
    connPool chan *sql.Conn
}

func NewDB(maxConns int) *DB {
    db := &DB{
        connPool: make(chan *sql.Conn, maxConns),
    }

    for i := 0; i < maxConns; i++ {
        db.connPool <- createConnection()
    }

    return db
}

func (db *DB) Query(query string) (Result, error) {
    conn := <-db.connPool  // Acquire connection
    defer func() { db.connPool <- conn }()  // Release

    return conn.Query(query)
}
```

### File Descriptor Limit

```go
type FilePool struct {
    sem chan struct{}
}

func NewFilePool(maxFiles int) *FilePool {
    return &FilePool{
        sem: make(chan struct{}, maxFiles),
    }
}

func (fp *FilePool) OpenFile(path string) (*os.File, error) {
    fp.sem <- struct{}{}  // Acquire file descriptor slot

    f, err := os.Open(path)
    if err != nil {
        <-fp.sem  // Release on error
        return nil, err
    }

    return &wrappedFile{File: f, pool: fp}, nil
}

type wrappedFile struct {
    *os.File
    pool *FilePool
}

func (wf *wrappedFile) Close() error {
    err := wf.File.Close()
    <-wf.pool.sem  // Release file descriptor slot
    return err
}
```

---

## How to Run

```bash
# Run the demonstration program
go run ./minis/23-bounded-channel-semaphore/cmd/semaphore-demo/main.go

# Run exercises
go test -v ./minis/23-bounded-channel-semaphore/exercise/...

# Run with race detector
go test -race ./minis/23-bounded-channel-semaphore/exercise/...

# Benchmark semaphore performance
go test -bench=. ./minis/23-bounded-channel-semaphore/exercise/...
```

---

## Summary

**What you learned**:
- ✅ Buffered channels are idiomatic counting semaphores in Go
- ✅ Channel capacity = max concurrent permits
- ✅ Send = acquire (blocks when full), receive = release
- ✅ Always use defer for release (prevents permit leaks)
- ✅ Context-aware acquisition prevents indefinite blocking
- ✅ Try-acquire enables graceful degradation
- ✅ Weighted semaphores for variable resource costs
- ✅ Common patterns: connection pools, rate limiters, worker pools

**Key patterns**:
1. **Basic**: `sem <- struct{}{}` (acquire), `defer func() { <-sem }()` (release)
2. **Non-blocking**: `select { case sem <- struct{}{}: ... default: ... }`
3. **Context-aware**: `select { case sem <- struct{}{}: ... case <-ctx.Done(): ... }`
4. **Wait for all**: Acquire all permits to ensure completion

**When to use**:
- Limiting concurrent access to resources (DB, files, APIs)
- Rate limiting requests
- Worker pool sizing
- Preventing resource exhaustion

**Next steps**:
- Project 24: sync.Mutex vs RWMutex (alternative to channels)
- Project 25: Atomic operations (lock-free synchronization)
- Project 27: sync.Pool (object pooling pattern)

Master semaphores, master resource control in concurrent Go!
