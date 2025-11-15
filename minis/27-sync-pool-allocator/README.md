# Project 27: sync.Pool - Object Pooling and Allocation Reduction

## 1. What Is This About?

### Real-World Scenario

Imagine you're building a high-performance JSON API server that handles 100,000 requests per second. Each request needs temporary buffers for parsing and encoding JSON:

**❌ Without object pooling (naive approach):**
```go
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    // Allocate new buffer for every request
    buf := new(bytes.Buffer)  // Heap allocation!

    // Parse request
    var data RequestData
    json.NewDecoder(r.Body).Decode(&data)

    // Encode response
    json.NewEncoder(buf).Encode(data)
    w.Write(buf.Bytes())

    // buf is discarded → garbage collector must clean it up
}
```

**Problem**: With 100,000 requests/second:
- **100,000 allocations/second** (creates memory pressure)
- **Frequent GC pauses** (stop-the-world events)
- **Higher latency** (p99 latency spikes during GC)
- **Wasted CPU** (GC overhead instead of serving requests)

**Profiling results:**
```
GC overhead:     15% of CPU time
Allocations:     100,000/sec
GC pause time:   10-50ms (unacceptable for low-latency APIs)
Throughput:      Only 70,000 req/sec (CPU bottlenecked by GC)
```

**✅ With sync.Pool (object pooling):**
```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
    // Get buffer from pool (reuses existing buffer)
    buf := bufferPool.Get().(*bytes.Buffer)
    buf.Reset()  // Clear previous data
    defer bufferPool.Put(buf)  // Return to pool

    // Parse request
    var data RequestData
    json.NewDecoder(r.Body).Decode(&data)

    // Encode response
    json.NewEncoder(buf).Encode(data)
    w.Write(buf.Bytes())

    // buf is returned to pool → reused by next request
}
```

**Result**:
```
GC overhead:     2% of CPU time (7.5x reduction!)
Allocations:     ~1,000/sec (99% reduction!)
GC pause time:   1-2ms (5-25x improvement!)
Throughput:      100,000 req/sec (43% improvement!)
```

This project teaches you **sync.Pool** and **object pooling patterns**:
- **sync.Pool**: Go's built-in object pool for temporary objects
- **Allocation reduction**: Minimizing heap allocations and GC pressure
- **Memory profiling**: Understanding allocation hotspots with pprof
- **GC tuning**: Optimizing garbage collector behavior
- **Performance patterns**: When and how to use pooling effectively

### What You'll Learn

1. **sync.Pool API**: Get(), Put(), New function
2. **Object lifecycle**: Pool management, resetting state, GC interaction
3. **Allocation profiling**: Finding allocation hotspots with pprof
4. **GC fundamentals**: How Go's garbage collector works
5. **Performance optimization**: Measuring impact of pooling with benchmarks
6. **Common patterns**: Buffer pools, connection pools, worker pools
7. **Anti-patterns**: When NOT to use sync.Pool

### The Challenge

Build high-performance systems that demonstrate:
- Reducing allocations by 90%+ using sync.Pool
- Measuring GC overhead before and after pooling
- Creating custom pools for different object types
- Understanding pool sizing and GC reclamation
- Avoiding common pooling pitfalls

---

## 2. First Principles: Understanding Memory Allocation

### What is Heap Allocation?

When you create objects in Go, they can be stored in two places:

**1. Stack** (fast, automatic cleanup):
```go
func process() {
    x := 42  // Stack allocation (local variable)
    y := make([]int, 10)  // MIGHT be stack (escape analysis)
    // Automatically freed when function returns
}
```

**2. Heap** (slower, requires GC):
```go
func createBuffer() *bytes.Buffer {
    buf := new(bytes.Buffer)  // Heap allocation (escapes to heap)
    return buf  // Pointer returned → must live beyond function
}
```

**How does Go decide?**

**Escape analysis** determines whether a variable can stay on the stack:

```bash
go build -gcflags="-m" main.go
```

Output:
```
./main.go:10:6: moved to heap: buf
./main.go:15:2: x does not escape
```

**Rules of thumb:**
- **Stack**: Local variables that don't escape the function
- **Heap**:
  - Variables returned by pointer
  - Variables stored in closures
  - Variables too large for stack (>10MB)
  - Variables with unknown size at compile time

### The Cost of Heap Allocation

**Stack allocation** (fast):
```
Time:    ~1 nanosecond (just moving stack pointer)
Cleanup: Automatic (stack unwind)
```

**Heap allocation** (slower):
```
Time:    ~50-100 nanoseconds
Process:
  1. Find free memory in heap
  2. Update allocator metadata
  3. Zero the memory
  4. Return pointer
Cleanup: Garbage collector must trace and free
```

**Example: 1 million allocations**
```
Stack:   1ms total
Heap:    50-100ms total + GC overhead
```

### Why Heap Allocations Matter

**Memory pressure**:
```go
// Allocates 1GB/sec
for {
    data := make([]byte, 1024*1024)  // 1MB allocation
    process(data)
    // data is now garbage
}
```

**What happens:**
1. **Heap fills up** (default: GC triggered at 2x heap size)
2. **GC pauses program** (stop-the-world)
3. **GC scans memory** (traces all live objects)
4. **GC frees dead objects** (reclaims memory)
5. **Program resumes**

**GC pause times:**
```
Small heap (10MB):    0.1-1ms
Medium heap (1GB):    1-10ms
Large heap (10GB):    10-100ms
```

For low-latency systems, even 1ms is **too long**.

---

## 3. First Principles: Understanding sync.Pool

### What is sync.Pool?

**sync.Pool** is a **set of temporary objects** that may be individually saved and retrieved.

**Key characteristics:**
- **Thread-safe**: Safe for concurrent Get/Put operations
- **Temporary**: Objects may be automatically freed by GC at any time
- **No size limit**: Pool can grow unbounded
- **No ownership**: You don't control when objects are freed

**Conceptual model:**
```
sync.Pool = Thread-local caches + Global shared queue

Thread 1:  [obj1, obj2, obj3] → Local cache
Thread 2:  [obj4, obj5]       → Local cache
Thread 3:  []                 → Local cache

Global:    [obj6, obj7, obj8] → Shared queue

Get():  Try local cache → Try global queue → Call New()
Put():  Put to local cache (if space) → Put to global queue
```

### The sync.Pool API

**Creating a pool:**
```go
var pool = sync.Pool{
    New: func() interface{} {
        // Called when pool is empty
        return new(bytes.Buffer)
    },
}
```

**Getting an object:**
```go
obj := pool.Get()  // Returns interface{}
buf := obj.(*bytes.Buffer)  // Type assertion required
```

**Returning an object:**
```go
pool.Put(buf)  // Returns object to pool
```

**Complete pattern:**
```go
buf := pool.Get().(*bytes.Buffer)
buf.Reset()  // Clear any previous data
defer pool.Put(buf)  // Ensure object is returned

// Use buf...
```

### How sync.Pool Works Internally

**Per-P (processor) caches**:
```
Go's scheduler runs goroutines on P (processor) structures.
Each P has its own local pool cache (lock-free access).

P0: [private: obj1] [shared: obj2, obj3, obj4]
P1: [private: obj5] [shared: obj6, obj7]
P2: [private: nil]  [shared: obj8]
```

**Get() algorithm:**
```go
1. Check if private slot has object → return it
2. Check local shared slice → pop from tail
3. Try stealing from other P's shared slices
4. Call New() to create new object
```

**Put() algorithm:**
```go
1. If private slot is empty → store there
2. Otherwise, push to local shared slice
3. If local shared is full → drop object
```

**GC interaction:**
```go
// During GC cycle
poolCleanup() {
    for each Pool {
        for each P {
            // Move current objects to victim cache
            p.victim = p.local
            p.local = nil
        }
    }
}

// Next GC cycle
poolCleanup() {
    for each Pool {
        for each P {
            // Free victim cache (objects from 2 GCs ago)
            p.victim = nil
        }
    }
}
```

**Key insight**: Objects survive at least **two GC cycles** before being freed.

### Pool Lifecycle Example

```go
var pool = sync.Pool{New: func() interface{} { return new(int) }}

// Time T0: Pool is empty
obj1 := pool.Get()  // Calls New(), returns new(int)

// Time T1: Put object back
pool.Put(obj1)
// Pool: [obj1]

// Time T2: Get returns same object
obj2 := pool.Get()  // Returns obj1 (reused!)
// Pool: []

// Time T3: GC runs
runtime.GC()
// Pool: [] → Objects moved to victim cache

// Time T4: GC runs again
runtime.GC()
// Victim cache cleared → Objects freed

// Time T5: Pool is empty again
obj3 := pool.Get()  // Calls New() again
```

---

## 4. First Principles: Object Pooling Patterns

### Pattern 1: Buffer Pool

**Problem**: HTTP handlers allocate buffers for every request.

**Solution:**
```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func handler(w http.ResponseWriter, r *http.Request) {
    buf := bufferPool.Get().(*bytes.Buffer)
    buf.Reset()  // Important: Clear previous data
    defer bufferPool.Put(buf)

    // Write to buffer
    fmt.Fprintf(buf, "Response data")

    // Send response
    w.Write(buf.Bytes())
}
```

**Why it works:**
- Buffers are **temporary** (only needed during request handling)
- Buffers are **reusable** (same type, different data)
- High **allocation rate** (many requests/second)

### Pattern 2: Slice Pool

**Problem**: Frequent allocations of temporary slices.

**Solution:**
```go
var slicePool = sync.Pool{
    New: func() interface{} {
        // Allocate with initial capacity
        s := make([]byte, 0, 4096)
        return &s
    },
}

func processData(data []byte) {
    bufPtr := slicePool.Get().(*[]byte)
    buf := *bufPtr
    buf = buf[:0]  // Reset length, keep capacity
    defer slicePool.Put(bufPtr)

    // Use buf for temporary storage
    buf = append(buf, data...)
    result := transform(buf)
    // ...
}
```

**Important**: Pool stores **pointers to slices**, not slices directly.

### Pattern 3: Struct Pool

**Problem**: Allocating temporary structs for processing.

**Solution:**
```go
type ProcessingContext struct {
    Buffer     []byte
    TempData   map[string]string
    Result     interface{}
}

var contextPool = sync.Pool{
    New: func() interface{} {
        return &ProcessingContext{
            Buffer:   make([]byte, 0, 1024),
            TempData: make(map[string]string),
        }
    },
}

func process(input []byte) interface{} {
    ctx := contextPool.Get().(*ProcessingContext)
    defer contextPool.Put(ctx)

    // Reset state
    ctx.Buffer = ctx.Buffer[:0]
    for k := range ctx.TempData {
        delete(ctx.TempData, k)
    }
    ctx.Result = nil

    // Use context...
    return ctx.Result
}
```

**Critical**: Always reset mutable state before reuse!

### Pattern 4: Writer Pool

**Problem**: Creating encoders/writers repeatedly.

**Solution:**
```go
var writerPool = sync.Pool{
    New: func() interface{} {
        return gzip.NewWriter(ioutil.Discard)
    },
}

func compressData(data []byte, w io.Writer) error {
    gzWriter := writerPool.Get().(*gzip.Writer)
    defer writerPool.Put(gzWriter)

    gzWriter.Reset(w)  // Reset to write to new destination

    _, err := gzWriter.Write(data)
    if err != nil {
        return err
    }

    return gzWriter.Close()
}
```

**Why it works**: gzip.NewWriter allocates compression buffers (~256KB).

### Pattern 5: Size-Classed Pools

**Problem**: Objects of different sizes need different pools.

**Solution:**
```go
var bufferPools = [...]sync.Pool{
    // 1KB pool
    {New: func() interface{} { b := make([]byte, 1024); return &b }},
    // 4KB pool
    {New: func() interface{} { b := make([]byte, 4096); return &b }},
    // 16KB pool
    {New: func() interface{} { b := make([]byte, 16384); return &b }},
    // 64KB pool
    {New: func() interface{} { b := make([]byte, 65536); return &b }},
}

func getBuffer(size int) *[]byte {
    // Select appropriate pool based on size
    switch {
    case size <= 1024:
        return bufferPools[0].Get().(*[]byte)
    case size <= 4096:
        return bufferPools[1].Get().(*[]byte)
    case size <= 16384:
        return bufferPools[2].Get().(*[]byte)
    default:
        return bufferPools[3].Get().(*[]byte)
    }
}

func putBuffer(buf *[]byte) {
    size := cap(*buf)
    switch {
    case size == 1024:
        bufferPools[0].Put(buf)
    case size == 4096:
        bufferPools[1].Put(buf)
    case size == 16384:
        bufferPools[2].Put(buf)
    case size == 65536:
        bufferPools[3].Put(buf)
    }
}
```

---

## 5. Understanding Go's Garbage Collector

### How Go's GC Works

**Go uses a concurrent mark-and-sweep garbage collector:**

**1. Mark phase** (finds live objects):
```
Start with roots (globals, stack variables, registers)
→ Trace all reachable objects
→ Mark them as "live"
```

**2. Sweep phase** (frees dead objects):
```
Scan heap memory
→ Free unmarked (dead) objects
→ Return memory to allocator
```

**3. Concurrent execution**:
```
Most GC work happens concurrently with application
Some phases require "stop-the-world" pauses:
  - Root scanning (very brief, ~10-100μs)
  - Write barrier coordination
```

### GC Triggers

**Automatic triggers:**
```go
// GC runs when heap size grows by GOGC%
// Default: GOGC=100 (GC at 2x live heap size)

Live heap: 100MB
Next GC:   200MB (100MB + 100% of 100MB)

After GC, live heap: 150MB
Next GC:   300MB (150MB + 100% of 150MB)
```

**Manual trigger:**
```go
runtime.GC()  // Force immediate GC (for testing)
```

**Tuning GOGC:**
```bash
# More frequent GC (lower memory, more CPU)
GOGC=50 ./myapp

# Less frequent GC (higher memory, less CPU)
GOGC=200 ./myapp

# Disable automatic GC (manual control only)
GOGC=off ./myapp
```

### GC Metrics

**Reading GC stats:**
```go
var stats runtime.MemStats
runtime.ReadMemStats(&stats)

fmt.Printf("Alloc: %d MB\n", stats.Alloc/1024/1024)
fmt.Printf("TotalAlloc: %d MB\n", stats.TotalAlloc/1024/1024)
fmt.Printf("NumGC: %d\n", stats.NumGC)
fmt.Printf("PauseTotal: %v\n", time.Duration(stats.PauseTotalNs))
fmt.Printf("LastGC: %v ago\n", time.Since(time.Unix(0, int64(stats.LastGC))))
```

**Important metrics:**
- `Alloc`: Current heap usage
- `TotalAlloc`: Cumulative allocations (shows allocation rate)
- `Mallocs/Frees`: Number of allocations/frees
- `NumGC`: Number of GC cycles
- `PauseTotalNs`: Total time in stop-the-world pauses

### GC Impact on Performance

**Scenario: API server with high allocation rate**

**Without pooling:**
```
Allocation rate: 1GB/sec
Live heap:       100MB
GC frequency:    Every 100ms (1GB / (100MB * 1.0))
GC pause:        5-10ms per cycle
Total pause:     50-100ms per second (5-10% overhead!)
```

**With pooling:**
```
Allocation rate: 100MB/sec (90% reduction)
Live heap:       100MB
GC frequency:    Every 1 second
GC pause:        5-10ms per cycle
Total pause:     5-10ms per second (0.5-1% overhead)
```

**Result**: 10x reduction in GC overhead!

---

## 6. Allocation Profiling with pprof

### Enabling Allocation Profiling

**In production servers:**
```go
import _ "net/http/pprof"

func main() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()

    // Your application...
}
```

**In benchmarks:**
```go
func BenchmarkWithPooling(b *testing.B) {
    for i := 0; i < b.N; i++ {
        processWithPool()
    }
}

// Run: go test -bench=. -benchmem -memprofile=mem.out
```

### Capturing Allocation Profiles

**Heap profile (current allocations):**
```bash
# Capture profile
curl http://localhost:6060/debug/pprof/heap > heap.out

# Analyze
go tool pprof heap.out
```

**Allocation profile (all allocations):**
```bash
# Capture profile
curl http://localhost:6060/debug/pprof/allocs > allocs.out

# Analyze
go tool pprof allocs.out
```

### Analyzing Profiles

**Interactive pprof commands:**
```bash
go tool pprof mem.out

# Show top allocation sites
(pprof) top

# Show top with bytes
(pprof) top -cum

# List source code for function
(pprof) list functionName

# Show graph
(pprof) web

# Show percentages
(pprof) top -cum -sample_index=alloc_space
```

**Example output:**
```
      flat  flat%   sum%        cum   cum%
  512.50MB 45.67% 45.67%   512.50MB 45.67%  bytes.makeSlice
  256.25MB 22.84% 68.51%   256.25MB 22.84%  encoding/json.Unmarshal
  128.00MB 11.41% 79.92%   128.00MB 11.41%  net/http.(*conn).serve
```

**Interpretation:**
- `flat`: Allocations in this function only
- `cum`: Cumulative allocations (this function + callees)
- Focus on high `cum` values

### Comparing Before/After

```bash
# Capture before optimization
go test -bench=. -memprofile=before.out

# Make changes (add pooling)

# Capture after optimization
go test -bench=. -memprofile=after.out

# Compare
go tool pprof -base before.out after.out
```

---

## 7. When to Use (and NOT Use) sync.Pool

### When to Use sync.Pool ✅

**1. High allocation rate** (>10,000/sec):
```go
// HTTP server handling many requests
var bufferPool = sync.Pool{
    New: func() interface{} { return new(bytes.Buffer) },
}
```

**2. Temporary objects** (short-lived):
```go
// Objects used only during request/task
func handleRequest() {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer bufferPool.Put(buf)
    // Use buf...
}
```

**3. Expensive to create** (heavy initialization):
```go
// Compression writers allocate large buffers
var gzipPool = sync.Pool{
    New: func() interface{} {
        return gzip.NewWriter(ioutil.Discard)  // ~256KB allocation
    },
}
```

**4. Uniform size/type** (same shape):
```go
// All objects have similar structure
type Request struct {
    Buffer []byte
    Data   map[string]string
}
```

**5. Profiling shows impact** (measure first!):
```bash
# Allocation profile shows hotspot
go test -bench=. -memprofile=mem.out
go tool pprof mem.out
(pprof) top
```

### When NOT to Use sync.Pool ❌

**1. Long-lived objects** (persistent state):
```go
// BAD: Database connections should use a real pool
var dbPool = sync.Pool{  // Wrong!
    New: func() interface{} { return openDBConnection() },
}
// Use database/sql.DB instead
```

**2. Fixed pool size needed** (bounded resources):
```go
// BAD: sync.Pool has no size limit
var connPool = sync.Pool{  // Wrong!
    New: func() interface{} { return dialServer() },
}
// Use buffered channel instead
```

**3. Objects need cleanup** (resources):
```go
// BAD: Files need explicit Close()
var filePool = sync.Pool{  // Wrong!
    New: func() interface{} { return os.Create("tmp") },
}
// Files might leak if not closed
```

**4. Low allocation rate** (<100/sec):
```go
// BAD: Premature optimization
var configPool = sync.Pool{  // Not worth it
    New: func() interface{} { return new(Config) },
}
// Allocation overhead is negligible
```

**5. Large objects** (>1MB):
```go
// BAD: Large buffers can cause memory bloat
var hugePool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 10*1024*1024)  // 10MB
    },
}
// GC might not free these for multiple cycles
```

### sync.Pool vs Other Pools

**sync.Pool (temporary objects):**
- ✅ Automatic cleanup (GC managed)
- ✅ Unbounded size
- ✅ Thread-safe
- ❌ No control over object lifetime
- ❌ No size limit

**Buffered channel pool (bounded resources):**
```go
type ConnPool struct {
    conns chan *Connection
}

func NewConnPool(size int) *ConnPool {
    p := &ConnPool{
        conns: make(chan *Connection, size),
    }
    for i := 0; i < size; i++ {
        p.conns <- newConnection()
    }
    return p
}

func (p *ConnPool) Get() *Connection {
    return <-p.conns
}

func (p *ConnPool) Put(c *Connection) {
    p.conns <- c
}
```

- ✅ Fixed size limit
- ✅ Guaranteed availability (if pre-warmed)
- ✅ Controlled cleanup
- ❌ Blocking on empty/full
- ❌ Manual management

**Custom pool (complex logic):**
- ✅ Full control
- ✅ Custom eviction policies
- ✅ Metrics and monitoring
- ❌ More complex to implement
- ❌ More code to maintain

---

## 8. Common Patterns You Can Reuse

### Pattern 1: Generic Buffer Pool

```go
type BufferPool struct {
    pool sync.Pool
}

func NewBufferPool() *BufferPool {
    return &BufferPool{
        pool: sync.Pool{
            New: func() interface{} {
                return new(bytes.Buffer)
            },
        },
    }
}

func (bp *BufferPool) Get() *bytes.Buffer {
    return bp.pool.Get().(*bytes.Buffer)
}

func (bp *BufferPool) Put(buf *bytes.Buffer) {
    buf.Reset()
    bp.pool.Put(buf)
}
```

### Pattern 2: Typed Pool (Generics)

```go
type Pool[T any] struct {
    pool sync.Pool
    reset func(*T)
}

func NewPool[T any](newFunc func() *T, resetFunc func(*T)) *Pool[T] {
    return &Pool[T]{
        pool: sync.Pool{
            New: func() interface{} {
                return newFunc()
            },
        },
        reset: resetFunc,
    }
}

func (p *Pool[T]) Get() *T {
    return p.pool.Get().(*T)
}

func (p *Pool[T]) Put(obj *T) {
    if p.reset != nil {
        p.reset(obj)
    }
    p.pool.Put(obj)
}

// Usage:
var stringSlicePool = NewPool(
    func() *[]string {
        s := make([]string, 0, 10)
        return &s
    },
    func(s *[]string) {
        *s = (*s)[:0]
    },
)
```

### Pattern 3: Pool with Metrics

```go
type MetricsPool struct {
    pool     sync.Pool
    gets     atomic.Int64
    puts     atomic.Int64
    news     atomic.Int64
}

func NewMetricsPool(newFunc func() interface{}) *MetricsPool {
    mp := &MetricsPool{}
    mp.pool.New = func() interface{} {
        mp.news.Add(1)
        return newFunc()
    }
    return mp
}

func (mp *MetricsPool) Get() interface{} {
    mp.gets.Add(1)
    return mp.pool.Get()
}

func (mp *MetricsPool) Put(obj interface{}) {
    mp.puts.Add(1)
    mp.pool.Put(obj)
}

func (mp *MetricsPool) Stats() (gets, puts, news int64) {
    return mp.gets.Load(), mp.puts.Load(), mp.news.Load()
}

func (mp *MetricsPool) Efficiency() float64 {
    g := mp.gets.Load()
    n := mp.news.Load()
    if g == 0 {
        return 0
    }
    return float64(g-n) / float64(g) * 100
}
```

### Pattern 4: Size-Limited Pool

```go
type BoundedPool struct {
    pool     sync.Pool
    semaphore chan struct{}
}

func NewBoundedPool(maxSize int, newFunc func() interface{}) *BoundedPool {
    return &BoundedPool{
        pool: sync.Pool{New: newFunc},
        semaphore: make(chan struct{}, maxSize),
    }
}

func (bp *BoundedPool) Get() interface{} {
    bp.semaphore <- struct{}{}
    return bp.pool.Get()
}

func (bp *BoundedPool) Put(obj interface{}) {
    bp.pool.Put(obj)
    <-bp.semaphore
}
```

### Pattern 5: Pool with Validation

```go
type ValidatedPool[T any] struct {
    pool     sync.Pool
    validate func(*T) bool
    reset    func(*T)
}

func NewValidatedPool[T any](
    newFunc func() *T,
    validateFunc func(*T) bool,
    resetFunc func(*T),
) *ValidatedPool[T] {
    return &ValidatedPool[T]{
        pool: sync.Pool{
            New: func() interface{} {
                return newFunc()
            },
        },
        validate: validateFunc,
        reset:    resetFunc,
    }
}

func (vp *ValidatedPool[T]) Get() *T {
    for {
        obj := vp.pool.Get().(*T)
        if vp.validate == nil || vp.validate(obj) {
            return obj
        }
        // Object failed validation, try again
    }
}

func (vp *ValidatedPool[T]) Put(obj *T) {
    if vp.validate != nil && !vp.validate(obj) {
        return  // Don't return invalid objects
    }
    if vp.reset != nil {
        vp.reset(obj)
    }
    vp.pool.Put(obj)
}
```

---

## 9. Real-World Applications

### HTTP Server Buffer Pooling

```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "sync"
)

var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

type Response struct {
    Status  string `json:"status"`
    Data    interface{} `json:"data"`
    Message string `json:"message"`
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
    // Get buffer from pool
    buf := bufferPool.Get().(*bytes.Buffer)
    buf.Reset()
    defer bufferPool.Put(buf)

    // Create response
    resp := Response{
        Status:  "success",
        Data:    map[string]string{"key": "value"},
        Message: "Request processed",
    }

    // Encode to buffer
    if err := json.NewEncoder(buf).Encode(resp); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Write response
    w.Header().Set("Content-Type", "application/json")
    w.Write(buf.Bytes())
}

func main() {
    http.HandleFunc("/api", apiHandler)
    http.ListenAndServe(":8080", nil)
}
```

**Impact:**
- Before: 50,000 allocs/sec, 5ms GC pause
- After: 500 allocs/sec, 0.5ms GC pause

### Log Formatting

```go
package logger

import (
    "bytes"
    "fmt"
    "sync"
    "time"
)

var logBufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

type Logger struct {
    level string
}

func (l *Logger) Log(level, message string, args ...interface{}) {
    buf := logBufferPool.Get().(*bytes.Buffer)
    buf.Reset()
    defer logBufferPool.Put(buf)

    // Format log message
    fmt.Fprintf(buf, "[%s] %s: ", time.Now().Format(time.RFC3339), level)
    fmt.Fprintf(buf, message, args...)
    buf.WriteString("\n")

    // Write to output (stdout, file, etc.)
    fmt.Print(buf.String())
}
```

### JSON Encoding Pool

```go
package jsonpool

import (
    "bytes"
    "encoding/json"
    "sync"
)

var encoderPool = sync.Pool{
    New: func() interface{} {
        return json.NewEncoder(new(bytes.Buffer))
    },
}

func Marshal(v interface{}) ([]byte, error) {
    encoder := encoderPool.Get().(*json.Encoder)
    buf := encoder.(*json.Encoder).(*bytes.Buffer)
    buf.Reset()
    defer encoderPool.Put(encoder)

    if err := encoder.Encode(v); err != nil {
        return nil, err
    }

    // Copy bytes (buffer will be reused)
    result := make([]byte, buf.Len())
    copy(result, buf.Bytes())

    return result, nil
}
```

### CSV Processing

```go
package csvprocessor

import (
    "encoding/csv"
    "io"
    "sync"
)

var recordPool = sync.Pool{
    New: func() interface{} {
        // Pre-allocate slice with capacity
        rec := make([]string, 0, 10)
        return &rec
    },
}

func ProcessCSV(r io.Reader, handler func([]string) error) error {
    reader := csv.NewReader(r)

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }

        // Get pooled slice
        recPtr := recordPool.Get().(*[]string)
        rec := *recPtr
        rec = rec[:0]

        // Copy record data
        rec = append(rec, record...)

        // Process
        if err := handler(rec); err != nil {
            recordPool.Put(recPtr)
            return err
        }

        // Return to pool
        recordPool.Put(recPtr)
    }

    return nil
}
```

### gRPC Connection Pool

```go
package grpcpool

import (
    "compress/gzip"
    "io"
    "sync"
)

var gzipWriterPool = sync.Pool{
    New: func() interface{} {
        return gzip.NewWriter(io.Discard)
    },
}

var gzipReaderPool = sync.Pool{
    New: func() interface{} {
        // Reader needs to be reset before use
        return new(gzip.Reader)
    },
}

func CompressGRPC(w io.Writer, data []byte) error {
    gzw := gzipWriterPool.Get().(*gzip.Writer)
    defer gzipWriterPool.Put(gzw)

    gzw.Reset(w)

    if _, err := gzw.Write(data); err != nil {
        return err
    }

    return gzw.Close()
}

func DecompressGRPC(r io.Reader) ([]byte, error) {
    gzr := gzipReaderPool.Get().(*gzip.Reader)
    defer gzipReaderPool.Put(gzr)

    if err := gzr.Reset(r); err != nil {
        return nil, err
    }

    return io.ReadAll(gzr)
}
```

---

## 10. Common Mistakes to Avoid

### Mistake 1: Not Resetting Object State

**❌ Wrong**:
```go
var bufPool = sync.Pool{
    New: func() interface{} { return new(bytes.Buffer) },
}

func process() {
    buf := bufPool.Get().(*bytes.Buffer)
    defer bufPool.Put(buf)

    // BUG: Buffer might contain data from previous use!
    buf.WriteString("new data")
    return buf.String()  // Might include old data
}
```

**✅ Correct**:
```go
func process() {
    buf := bufPool.Get().(*bytes.Buffer)
    buf.Reset()  // Clear previous contents
    defer bufPool.Put(buf)

    buf.WriteString("new data")
    return buf.String()
}
```

### Mistake 2: Storing Pooled Objects

**❌ Wrong**:
```go
type Handler struct {
    buffer *bytes.Buffer  // BAD: Stored reference
}

func (h *Handler) Init() {
    h.buffer = bufPool.Get().(*bytes.Buffer)  // Object might be GC'd
}

func (h *Handler) Process() {
    h.buffer.WriteString("data")  // Might panic or corrupt data
}
```

**✅ Correct**:
```go
func (h *Handler) Process() {
    buf := bufPool.Get().(*bytes.Buffer)
    buf.Reset()
    defer bufPool.Put(buf)

    buf.WriteString("data")
    // Use buf locally only
}
```

### Mistake 3: Pooling Large Objects

**❌ Wrong**:
```go
var hugePool = sync.Pool{
    New: func() interface{} {
        // 100MB allocation!
        return make([]byte, 100*1024*1024)
    },
}

// If 10 goroutines use this, 1GB of memory!
// Objects won't be freed until 2 GC cycles later
```

**✅ Correct**:
```go
var reasonablePool = sync.Pool{
    New: func() interface{} {
        // Reasonable size (4KB)
        b := make([]byte, 4096)
        return &b
    },
}

// If you need larger buffers, grow them as needed
func getBuffer(size int) *bytes.Buffer {
    buf := bufPool.Get().(*bytes.Buffer)
    buf.Reset()
    buf.Grow(size)
    return buf
}
```

### Mistake 4: Putting Wrong Type

**❌ Wrong**:
```go
var pool = sync.Pool{
    New: func() interface{} { return new(bytes.Buffer) },
}

func process() {
    buf := pool.Get().(*bytes.Buffer)
    defer pool.Put("wrong type")  // BUG: Putting string instead of *bytes.Buffer
}
```

**✅ Correct**:
```go
func process() {
    buf := pool.Get().(*bytes.Buffer)
    defer pool.Put(buf)  // Correct type
}
```

### Mistake 5: Relying on Pool for Resource Management

**❌ Wrong**:
```go
var filePool = sync.Pool{
    New: func() interface{} {
        f, _ := os.Create("/tmp/file")
        return f  // Files might leak!
    },
}

func process() {
    f := filePool.Get().(*os.File)
    defer filePool.Put(f)  // File might not be closed!

    f.Write([]byte("data"))
}
```

**✅ Correct**:
```go
// Don't pool files - use proper resource management
func process() {
    f, err := os.Create("/tmp/file")
    if err != nil {
        return err
    }
    defer f.Close()  // Explicit cleanup

    f.Write([]byte("data"))
}
```

### Mistake 6: Type Assertion Without Check

**❌ Wrong**:
```go
func process() {
    buf := pool.Get().(*bytes.Buffer)  // Panics if wrong type
    // ...
}
```

**✅ Correct**:
```go
func process() {
    obj := pool.Get()
    buf, ok := obj.(*bytes.Buffer)
    if !ok {
        // Handle error
        return fmt.Errorf("unexpected type: %T", obj)
    }
    defer pool.Put(buf)
    // ...
}
```

### Mistake 7: Premature Optimization

**❌ Wrong**:
```go
// Adding pooling without profiling first
var configPool = sync.Pool{
    New: func() interface{} { return new(Config) },
}

// Config is created once at startup - pooling adds complexity for no benefit
```

**✅ Correct**:
```go
// Profile first, optimize hot paths
go test -bench=. -memprofile=mem.out
go tool pprof mem.out
(pprof) top

// Only add pooling if allocations are significant
```

---

## 11. Benchmarking Pooling Impact

### Writing Benchmarks

```go
package pooling_test

import (
    "bytes"
    "sync"
    "testing"
)

// Without pooling
func BenchmarkWithoutPool(b *testing.B) {
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        buf := new(bytes.Buffer)
        buf.WriteString("hello world")
        _ = buf.String()
    }
}

// With pooling
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func BenchmarkWithPool(b *testing.B) {
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        buf := bufferPool.Get().(*bytes.Buffer)
        buf.Reset()
        buf.WriteString("hello world")
        _ = buf.String()
        bufferPool.Put(buf)
    }
}

// Parallel without pooling
func BenchmarkParallelWithoutPool(b *testing.B) {
    b.ReportAllocs()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            buf := new(bytes.Buffer)
            buf.WriteString("hello world")
            _ = buf.String()
        }
    })
}

// Parallel with pooling
func BenchmarkParallelWithPool(b *testing.B) {
    b.ReportAllocs()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            buf := bufferPool.Get().(*bytes.Buffer)
            buf.Reset()
            buf.WriteString("hello world")
            _ = buf.String()
            bufferPool.Put(buf)
        }
    })
}
```

### Running Benchmarks

```bash
# Run benchmarks
go test -bench=. -benchmem

# With memory profiling
go test -bench=. -benchmem -memprofile=mem.out

# With CPU profiling
go test -bench=. -cpuprofile=cpu.out

# Compare before/after
benchstat before.txt after.txt
```

### Interpreting Results

```
BenchmarkWithoutPool-8           5000000    250 ns/op    64 B/op    1 allocs/op
BenchmarkWithPool-8             10000000    120 ns/op     0 B/op    0 allocs/op

BenchmarkParallelWithoutPool-8  20000000     80 ns/op    64 B/op    1 allocs/op
BenchmarkParallelWithPool-8     50000000     25 ns/op     0 B/op    0 allocs/op
```

**Analysis:**
- **2x faster** with pooling (250ns → 120ns)
- **Zero allocations** with pooling (1 → 0 allocs/op)
- **3.2x faster** in parallel workload (80ns → 25ns)
- Higher parallel speedup due to reduced GC contention

---

## 12. Stretch Goals

### Goal 1: Implement Generic Pool with Metrics ⭐⭐

Track hit rate, miss rate, and object creation count.

### Goal 2: Build Size-Classed Buffer Pool ⭐⭐⭐

Multiple pools for different buffer sizes (1KB, 4KB, 16KB, 64KB).

### Goal 3: Create Hybrid Pool with TTL ⭐⭐⭐

Combine sync.Pool with explicit cleanup for stale objects.

### Goal 4: Benchmark Pool vs No Pool Impact ⭐⭐

Measure allocation reduction and GC impact with pprof.

### Goal 5: Implement Custom Ring Buffer Pool ⭐⭐⭐⭐

Build a lock-free pool using ring buffer for ultra-low latency.

---

## How to Run

```bash
# Run the demo
go run ./minis/27-sync-pool-allocator/cmd/pool-demo/main.go

# Run tests
go test ./minis/27-sync-pool-allocator/...

# Run benchmarks
go test -bench=. -benchmem ./minis/27-sync-pool-allocator/...

# Memory profiling
go test -bench=. -memprofile=mem.out ./minis/27-sync-pool-allocator/...
go tool pprof mem.out

# Allocation profiling
go test -bench=. -benchmem -memprofile=mem.out ./minis/27-sync-pool-allocator/...
go tool pprof -alloc_space mem.out

# Compare benchmarks
go test -bench=. -benchmem > before.txt
# Make changes
go test -bench=. -benchmem > after.txt
benchstat before.txt after.txt
```

---

## Summary

**What you learned**:
- ✅ sync.Pool provides temporary object pooling with automatic cleanup
- ✅ Object pooling reduces allocations and GC pressure
- ✅ Pool objects must be reset before reuse to avoid data corruption
- ✅ Pooling is most effective for high-frequency temporary allocations
- ✅ Use pprof to identify allocation hotspots before optimizing
- ✅ Pool objects survive at least 2 GC cycles before being freed
- ✅ Don't use sync.Pool for long-lived or resource-managed objects

**Why this matters**:
Object pooling is critical for high-performance Go applications:
- **API servers**: Reduce allocation rate from 1M/sec to 10K/sec
- **Data processing**: Reuse buffers across millions of records
- **Real-time systems**: Minimize GC pauses for consistent latency
- **High-throughput systems**: Maximize CPU efficiency by reducing GC overhead

**Key rules**:
1. Profile first - don't optimize without data
2. Always reset object state before reuse
3. Use for temporary objects only (short-lived)
4. Don't store pooled objects - use and return immediately
5. Prefer sync.Pool for unbounded temporary objects
6. Use channel pools for bounded resources (connections, workers)

**Next steps**:
- Explore context.Context for request-scoped data
- Learn about custom allocators for specialized use cases
- Study memory management patterns in high-performance Go libraries
- Investigate escape analysis to reduce allocations at the source

Master sync.Pool, master Go performance optimization! 🚀
