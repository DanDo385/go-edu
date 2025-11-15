# Project 28: pprof CPU/Memory Profiling & Benchmarks - Performance Analysis

## What Is This Project About?

This project teaches **pprof profiling**—Go's built-in performance analysis tool for finding CPU and memory bottlenecks. You'll learn:

1. **What pprof is** (profiler, sampling, overhead, visualization)
2. **CPU profiling** (finding hot paths, optimizing algorithms, goroutine profiling)
3. **Memory profiling** (heap allocations, memory leaks, escape analysis)
4. **Benchmarking** (how to write benchmarks, interpreting results, optimizations)
5. **Flamegraphs** (visualizing call stacks, identifying bottlenecks)
6. **pprof HTTP endpoints** (runtime profiling of running services)
7. **Optimization strategies** (algorithmic, allocation reduction, concurrency)

By the end, you'll know how to profile Go applications, identify performance bottlenecks, and apply systematic optimizations backed by measurements.

---

## The Fundamental Problem: Performance Bottlenecks

### First Principles: Why Profile?

**Premature optimization is the root of all evil** —Donald Knuth

Before optimizing:
1. **Measure first** (don't guess where the bottleneck is)
2. **Profile to find hotspots** (focus on the 20% that matters)
3. **Optimize systematically** (change one thing, measure again)
4. **Know when to stop** (diminishing returns)

**Example of wrong optimization:**
```go
// Programmer thinks: "String concatenation is slow, use bytes.Buffer!"
func buildMessage(name string, age int) string {
    var buf bytes.Buffer
    buf.WriteString("User: ")
    buf.WriteString(name)
    buf.WriteString(", Age: ")
    buf.WriteString(strconv.Itoa(age))
    return buf.String()
}

// Actual bottleneck: Database query takes 50ms
// String building: 100ns (500,000x faster!)
// Optimization wasted time, added complexity, no impact.
```

**The lesson:** Profile before optimizing. The bottleneck is rarely where you think it is.

---

## What Is pprof? (The Core Concept)

### Overview

**pprof** is Go's built-in profiling tool that:
1. **Samples program execution** (records stack traces at intervals)
2. **Collects performance data** (CPU time, memory allocations, goroutines, blocks)
3. **Visualizes hotspots** (text reports, graphs, flamegraphs, web UI)
4. **Low overhead** (~5% CPU for CPU profiling, minimal for memory)

**Key insight:** pprof is a **sampling profiler**, not a tracing profiler. It doesn't record every event—it samples periodically.

### Types of Profiles

| Profile Type | What It Measures | When to Use |
|--------------|------------------|-------------|
| **CPU** | Where the program spends execution time | Finding slow algorithms, hot loops |
| **Heap** | Memory allocations on the heap | Reducing allocations, finding leaks |
| **Allocs** | All memory allocations (including stack) | Fine-grained allocation analysis |
| **Goroutine** | Number of running goroutines | Finding goroutine leaks |
| **Block** | Where goroutines block (channels, locks) | Finding contention |
| **Mutex** | Lock contention hotspots | Optimizing synchronization |
| **Threadcreate** | OS thread creation | Diagnosing thread issues |

---

## CPU Profiling: Finding Hot Paths

### How CPU Profiling Works

**Mechanism:**
1. **Signal-based sampling:** OS sends SIGPROF signal 100 times/second
2. **Record stack trace:** When signal arrives, Go records the current stack
3. **Aggregate samples:** After profiling, pprof counts samples per function
4. **Identify hotspots:** Functions with most samples are the slowest

**Visual:**
```
Time →
|-----|-----|-----|-----|-----|-----|-----|-----|
  Sample at 10ms intervals

If function F appears in 40/100 samples:
→ F consumes ~40% of CPU time
```

### Enabling CPU Profiling

**Method 1: In Code (One-Time Profiling)**
```go
import (
    "os"
    "runtime/pprof"
)

func main() {
    // Create CPU profile file
    f, err := os.Create("cpu.prof")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    // Start CPU profiling
    if err := pprof.StartCPUProfile(f); err != nil {
        panic(err)
    }
    defer pprof.StopCPUProfile()

    // Your program logic here
    doWork()
}
```

**Method 2: HTTP Endpoint (Running Server)**
```go
import (
    _ "net/http/pprof"  // Registers /debug/pprof/* handlers
    "net/http"
)

func main() {
    // Start HTTP server (pprof handlers auto-registered)
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()

    // Your server logic
    runServer()
}
```

**Collecting profile from HTTP:**
```bash
# Collect 30-second CPU profile
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Or download to file
curl http://localhost:6060/debug/pprof/profile?seconds=30 > cpu.prof
go tool pprof cpu.prof
```

**Method 3: Testing (Benchmarks)**
```bash
# Run benchmarks and generate CPU profile
go test -bench=. -cpuprofile=cpu.prof

# Analyze
go tool pprof cpu.prof
```

### Analyzing CPU Profiles

**Interactive pprof commands:**
```bash
$ go tool pprof cpu.prof
(pprof) top
# Shows top 10 functions by CPU time

(pprof) top -cum
# Shows top 10 by cumulative time (includes callees)

(pprof) list functionName
# Shows annotated source code with sample counts

(pprof) web
# Opens graph visualization in browser (requires graphviz)

(pprof) pdf > profile.pdf
# Generates PDF call graph

(pprof) text
# Text-based report
```

**Sample output:**
```
(pprof) top
Showing nodes accounting for 3.50s, 87.50% of 4.00s total
Showing top 10 nodes out of 45
      flat  flat%   sum%        cum   cum%
     1.50s 37.50% 37.50%      1.50s 37.50%  runtime.mallocgc
     0.80s 20.00% 57.50%      0.80s 20.00%  main.processData
     0.50s 12.50% 70.00%      2.30s 57.50%  main.findPrimes
     0.40s 10.00% 80.00%      0.40s 10.00%  runtime.scanobject
     0.30s  7.50% 87.50%      0.30s  7.50%  fmt.Sprintf

flat  = CPU time spent IN this function (exclusive)
cum   = CPU time spent IN this function + callees (inclusive)
```

**Interpretation:**
- `runtime.mallocgc` at 37.50%: Lots of allocations (optimize to reduce GC pressure)
- `main.findPrimes` at 12.50% flat, 57.50% cum: Algorithm itself is 12.50%, but calls expensive functions
- Focus optimization on `processData` and `findPrimes`

### CPU Profiling Example

**Before optimization:**
```go
func findPrimes(n int) []int {
    var primes []int
    for i := 2; i <= n; i++ {
        isPrime := true
        for j := 2; j < i; j++ {
            if i%j == 0 {
                isPrime = false
                break
            }
        }
        if isPrime {
            primes = append(primes, i)
        }
    }
    return primes
}
```

**Profile shows:**
```
flat  70% in nested loop (i%j check)
```

**After optimization (Sieve of Eratosthenes):**
```go
func findPrimes(n int) []int {
    isPrime := make([]bool, n+1)
    for i := 2; i <= n; i++ {
        isPrime[i] = true
    }

    for i := 2; i*i <= n; i++ {
        if isPrime[i] {
            for j := i * i; j <= n; j += i {
                isPrime[j] = false
            }
        }
    }

    var primes []int
    for i := 2; i <= n; i++ {
        if isPrime[i] {
            primes = append(primes, i)
        }
    }
    return primes
}
```

**Result:** 100x faster for n=100,000 (O(n²) → O(n log log n))

---

## Memory Profiling: Finding Allocation Hotspots

### How Memory Profiling Works

**Mechanism:**
1. **Sample allocations:** Go samples ~1 allocation per 512KB allocated (configurable)
2. **Record stack traces:** Captures where allocations happen
3. **Track in-use vs allocated:** Distinguish current heap usage from total allocated

**Key metrics:**
- **alloc_space:** Total bytes allocated (including freed memory)
- **alloc_objects:** Total number of allocations
- **inuse_space:** Current heap usage (live objects)
- **inuse_objects:** Number of live objects

### Enabling Memory Profiling

**Method 1: In Code**
```go
import (
    "os"
    "runtime/pprof"
)

func main() {
    doWork()

    // Write heap profile
    f, err := os.Create("mem.prof")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    pprof.WriteHeapProfile(f)
}
```

**Method 2: HTTP Endpoint**
```bash
# Snapshot of current heap
go tool pprof http://localhost:6060/debug/pprof/heap

# Download
curl http://localhost:6060/debug/pprof/heap > mem.prof
go tool pprof mem.prof
```

**Method 3: Testing**
```bash
go test -bench=. -memprofile=mem.prof
go tool pprof mem.prof
```

### Analyzing Memory Profiles

**Interactive analysis:**
```bash
$ go tool pprof -alloc_space mem.prof
(pprof) top
# Shows functions allocating most memory

(pprof) list functionName
# Shows line-by-line allocation counts

(pprof) web
# Visual call graph
```

**Sample output:**
```
(pprof) top -alloc_space
      flat  flat%   sum%        cum   cum%
  512.00MB 64.00% 64.00%   512.00MB 64.00%  main.processStrings
  128.00MB 16.00% 80.00%   128.00MB 16.00%  fmt.Sprintf
   64.00MB  8.00% 88.00%    64.00MB  8.00%  encoding/json.Marshal

(pprof) top -inuse_space
      flat  flat%   sum%        cum   cum%
   64.00MB 80.00% 80.00%    64.00MB 80.00%  main.cache
   16.00MB 20.00% 100.0%    16.00MB 20.00%  main.buffer
```

**Interpretation:**
- `alloc_space`: `processStrings` allocated 512MB total (may have been freed)
- `inuse_space`: `cache` holds 64MB currently (potential memory leak if unexpected)

### Memory Profiling Example

**Before optimization:**
```go
func processLogs(logs []string) []string {
    var results []string
    for _, log := range logs {
        // Creates new string allocation on each concat
        result := "Processed: " + log + " at " + time.Now().String()
        results = append(results, result)
    }
    return results
}
```

**Profile shows:**
```
1.2GB allocated in string concatenation
800MB allocated in time.Now().String()
```

**After optimization:**
```go
func processLogs(logs []string) []string {
    results := make([]string, 0, len(logs))  // Preallocate capacity
    now := time.Now().String()  // Call once, reuse
    var buf strings.Builder  // Reusable buffer

    for _, log := range logs {
        buf.Reset()
        buf.WriteString("Processed: ")
        buf.WriteString(log)
        buf.WriteString(" at ")
        buf.WriteString(now)
        results = append(results, buf.String())
    }
    return results
}
```

**Result:** 80% reduction in allocations (1.2GB → 240MB)

---

## Benchmarking: Measuring Performance

### Writing Benchmarks

**Basic benchmark:**
```go
func BenchmarkMyFunction(b *testing.B) {
    for i := 0; i < b.N; i++ {
        MyFunction()
    }
}
```

**How it works:**
1. Go runs the loop multiple times
2. Adjusts `b.N` until execution time is stable (~1 second)
3. Reports operations/second and ns/operation

**Running benchmarks:**
```bash
# Run all benchmarks
go test -bench=.

# Run specific benchmark
go test -bench=BenchmarkMyFunction

# With memory stats
go test -bench=. -benchmem

# Multiple runs for stability
go test -bench=. -benchtime=10s -count=5
```

**Sample output:**
```
BenchmarkFindPrimes-8        1000    1234567 ns/op    524288 B/op    1024 allocs/op

Explanation:
- BenchmarkFindPrimes-8: Benchmark name, 8 = GOMAXPROCS
- 1000: Number of iterations (b.N)
- 1234567 ns/op: Nanoseconds per operation
- 524288 B/op: Bytes allocated per operation
- 1024 allocs/op: Number of allocations per operation
```

### Benchmark Best Practices

**1. Reset timer for setup code:**
```go
func BenchmarkWithSetup(b *testing.B) {
    data := generateTestData(10000)  // Setup (not measured)
    b.ResetTimer()  // Start timing here

    for i := 0; i < b.N; i++ {
        processData(data)
    }
}
```

**2. Prevent compiler optimizations:**
```go
var result int  // Package-level variable

func BenchmarkFunction(b *testing.B) {
    var r int
    for i := 0; i < b.N; i++ {
        r = expensiveCalculation()  // Store to prevent dead code elimination
    }
    result = r  // Assign to global to prevent optimization
}
```

**3. Run sub-benchmarks for comparisons:**
```go
func BenchmarkStringBuilding(b *testing.B) {
    b.Run("Concatenation", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = "Hello" + "World"
        }
    })

    b.Run("fmt.Sprintf", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = fmt.Sprintf("%s%s", "Hello", "World")
        }
    })

    b.Run("strings.Builder", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            var buf strings.Builder
            buf.WriteString("Hello")
            buf.WriteString("World")
            _ = buf.String()
        }
    })
}
```

**4. Use benchstat for statistical comparison:**
```bash
# Run benchmarks multiple times, save to files
go test -bench=. -count=10 > old.txt
# (make changes)
go test -bench=. -count=10 > new.txt

# Install benchstat
go install golang.org/x/perf/cmd/benchstat@latest

# Compare
benchstat old.txt new.txt
```

**Output:**
```
name         old time/op    new time/op    delta
FindPrimes   1.23ms ± 2%    0.12ms ± 1%   -90.24%  (p=0.000 n=10+10)

name         old alloc/op   new alloc/op   delta
FindPrimes   512kB ± 0%     64kB ± 0%      -87.50%  (p=0.000 n=10+10)
```

---

## Flamegraphs: Visualizing Call Stacks

### What Are Flamegraphs?

**Flamegraphs** visualize profiling data as interactive SVG diagrams:
- **X-axis:** Alphabetical order (NOT time!)
- **Y-axis:** Call stack depth
- **Width:** Proportion of samples (wider = more CPU time)
- **Color:** Random (for visual distinction, no meaning)

**Structure:**
```
[main (100%)]
├── [processData (60%)]
│   ├── [parseJSON (40%)]
│   └── [validateData (20%)]
└── [writeResults (40%)]
```

**Benefits:**
- Instantly see bottlenecks (wide boxes)
- Understand call hierarchy
- Interactive (click to zoom)

### Generating Flamegraphs

**Method 1: pprof built-in (Go 1.20+):**
```bash
go tool pprof -http=:8080 cpu.prof
# Opens web UI with flamegraph tab
```

**Method 2: Uber's go-torch (legacy):**
```bash
# Install go-torch and FlameGraph
go install github.com/uber/go-torch@latest
git clone https://github.com/brendangregg/FlameGraph

# Generate flamegraph
go-torch cpu.prof
# Creates torch.svg
```

**Method 3: From live server:**
```bash
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/profile?seconds=30
```

### Reading Flamegraphs

**Example interpretation:**
```
[main: 100%]
├── [runtime.gcDrain: 40%]  ← 40% of time in GC!
└── [myApp.processData: 60%]
    ├── [myApp.parseJSON: 35%]  ← Biggest hotspot
    │   └── [encoding/json.Unmarshal: 35%]
    └── [myApp.computeHash: 25%]
        └── [crypto/sha256.Sum256: 25%]
```

**Optimization strategy:**
1. **40% in GC:** Reduce allocations (preallocate, reuse buffers, sync.Pool)
2. **35% in JSON parsing:** Consider faster parsers (jsoniter, sonic) or binary formats (protobuf)
3. **25% in SHA256:** Profile-specific optimization (e.g., cache hashes, use faster hash if security allows)

---

## pprof HTTP Endpoints

### Standard Endpoints

When you import `_ "net/http/pprof"`, these endpoints become available:

| Endpoint | Description |
|----------|-------------|
| `/debug/pprof/` | Index page with available profiles |
| `/debug/pprof/profile?seconds=30` | 30-second CPU profile |
| `/debug/pprof/heap` | Heap memory profile |
| `/debug/pprof/goroutine` | Stack traces of all goroutines |
| `/debug/pprof/block` | Blocking profile |
| `/debug/pprof/mutex` | Mutex contention profile |
| `/debug/pprof/allocs` | All memory allocations |
| `/debug/pprof/threadcreate` | Thread creation profile |
| `/debug/pprof/trace?seconds=5` | 5-second execution trace |

### Using HTTP Endpoints

**Browser access:**
```bash
# View index
http://localhost:6060/debug/pprof/

# View goroutines as text
http://localhost:6060/debug/pprof/goroutine?debug=1

# View heap summary
http://localhost:6060/debug/pprof/heap?debug=1
```

**CLI analysis:**
```bash
# Interactive CPU profile
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Save and analyze later
curl http://localhost:6060/debug/pprof/heap > heap.prof
go tool pprof heap.prof

# Compare before/after
curl http://localhost:6060/debug/pprof/heap > before.prof
# (make changes, wait for traffic)
curl http://localhost:6060/debug/pprof/heap > after.prof
go tool pprof -base=before.prof after.prof
```

**Web UI:**
```bash
# Opens interactive web interface
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/profile?seconds=30
```

---

## Common Optimization Patterns

### Pattern 1: Reduce Allocations

**Problem:** Frequent small allocations cause GC pressure.

**Solution 1: Preallocate slices**
```go
// Bad: Grows slice multiple times
func collectResults(n int) []int {
    var results []int
    for i := 0; i < n; i++ {
        results = append(results, i)
    }
    return results
}

// Good: Preallocate exact size
func collectResults(n int) []int {
    results := make([]int, 0, n)  // Capacity = n
    for i := 0; i < n; i++ {
        results = append(results, i)
    }
    return results
}
```

**Solution 2: Reuse buffers (sync.Pool)**
```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func processData(data []byte) string {
    buf := bufferPool.Get().(*bytes.Buffer)
    buf.Reset()
    defer bufferPool.Put(buf)

    buf.Write(data)
    // Process...
    return buf.String()
}
```

### Pattern 2: Avoid String Allocations

**Problem:** Strings are immutable; concatenation creates copies.

**Solution: Use strings.Builder**
```go
// Bad: Creates n intermediate strings
func buildMessage(parts []string) string {
    msg := ""
    for _, part := range parts {
        msg += part  // Allocates new string each time
    }
    return msg
}

// Good: Single allocation
func buildMessage(parts []string) string {
    var buf strings.Builder
    buf.Grow(len(parts) * 10)  // Preallocate (estimate size)
    for _, part := range parts {
        buf.WriteString(part)
    }
    return buf.String()
}
```

### Pattern 3: Optimize Hot Loops

**Problem:** Small inefficiencies in tight loops compound.

**Example:**
```go
// Bad: Expensive operations in loop
func sumEvenSquares(nums []int) int {
    sum := 0
    for i := 0; i < len(nums); i++ {  // len() called every iteration
        if nums[i]%2 == 0 {
            sum += int(math.Pow(float64(nums[i]), 2))  // Float conversion!
        }
    }
    return sum
}

// Good: Hoist invariants, use cheaper operations
func sumEvenSquares(nums []int) int {
    sum := 0
    n := len(nums)  // Hoist len() call
    for i := 0; i < n; i++ {
        if nums[i]%2 == 0 {
            sum += nums[i] * nums[i]  // Integer multiplication
        }
    }
    return sum
}
```

### Pattern 4: Cache Expensive Computations

**Solution: Memoization**
```go
type Cache struct {
    mu    sync.RWMutex
    cache map[string]Result
}

func (c *Cache) GetOrCompute(key string, compute func() Result) Result {
    // Fast path: read lock
    c.mu.RLock()
    if result, ok := c.cache[key]; ok {
        c.mu.RUnlock()
        return result
    }
    c.mu.RUnlock()

    // Slow path: write lock
    c.mu.Lock()
    defer c.mu.Unlock()

    // Double-check (another goroutine may have computed it)
    if result, ok := c.cache[key]; ok {
        return result
    }

    result := compute()
    c.cache[key] = result
    return result
}
```

### Pattern 5: Parallelize CPU-Bound Work

**Solution: Worker pool**
```go
func processInParallel(items []Item) []Result {
    numWorkers := runtime.NumCPU()
    jobs := make(chan Item, len(items))
    results := make(chan Result, len(items))

    // Start workers
    var wg sync.WaitGroup
    for w := 0; w < numWorkers; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for item := range jobs {
                results <- processItem(item)
            }
        }()
    }

    // Send jobs
    for _, item := range items {
        jobs <- item
    }
    close(jobs)

    // Wait and collect
    go func() {
        wg.Wait()
        close(results)
    }()

    var output []Result
    for result := range results {
        output = append(output, result)
    }
    return output
}
```

---

## Escape Analysis: Stack vs Heap

### What Is Escape Analysis?

**Escape analysis** determines whether a variable can live on the **stack** (fast, no GC) or must live on the **heap** (slower, requires GC).

**Rules:**
- Variables that **don't escape** the function → **stack**
- Variables that **escape** (returned, stored in globals, captured by closures) → **heap**

### Viewing Escape Analysis

```bash
go build -gcflags="-m" main.go
# Shows escape analysis decisions

# More verbose
go build -gcflags="-m -m" main.go
```

**Example:**
```go
func createUser(name string) *User {
    u := User{Name: name}  // Escapes to heap (returned as pointer)
    return &u
}

// Output:
// ./main.go:10:6: moved to heap: u
```

### Reducing Heap Allocations

**Pattern 1: Return values instead of pointers**
```go
// Bad: Allocates on heap
func createPoint() *Point {
    p := Point{X: 1, Y: 2}
    return &p  // Escapes
}

// Good: Stack allocation
func createPoint() Point {
    return Point{X: 1, Y: 2}  // Stack
}
```

**Pattern 2: Avoid capturing large variables**
```go
// Bad: Closure captures largeData
func processAsync(largeData []byte) {
    go func() {
        process(largeData)  // largeData escapes to heap
    }()
}

// Good: Pass copy to goroutine
func processAsync(largeData []byte) {
    data := largeData  // Make explicit copy
    go func() {
        process(data)
    }()
}
```

---

## How to Run

```bash
# Run the demo server
cd minis/28-pprof-cpu-mem-benchmarks
go run cmd/pprof-demo/main.go

# In another terminal, collect profiles
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=10
go tool pprof http://localhost:6060/debug/pprof/heap

# Run benchmarks
cd exercise
go test -bench=. -benchmem

# With profiling
go test -bench=. -cpuprofile=cpu.prof -memprofile=mem.prof
go tool pprof cpu.prof

# Run exercises
go test -v
```

---

## Expected Output (Demo Server)

```
=== pprof Profiling Demonstration ===

Starting HTTP server on :6060
Profile endpoints available at:
  http://localhost:6060/debug/pprof/

Available profiles:
  - CPU:       /debug/pprof/profile?seconds=30
  - Heap:      /debug/pprof/heap
  - Goroutine: /debug/pprof/goroutine
  - Block:     /debug/pprof/block
  - Mutex:     /debug/pprof/mutex

Running workload generators...
  [CPU] Computing prime numbers...
  [Memory] Generating large data structures...
  [Concurrent] Running 100 goroutines...

Press Ctrl+C to stop.
```

---

## Key Takeaways

1. **Profile before optimizing** (measure, don't guess)
2. **CPU profiling** finds slow algorithms (use `go tool pprof`, look at `flat%` and `cum%`)
3. **Memory profiling** finds allocation hotspots (reduce GC pressure)
4. **Benchmarks** measure performance (use `-benchmem`, compare with `benchstat`)
5. **Flamegraphs** visualize call stacks (wide = hotspot)
6. **HTTP endpoints** enable runtime profiling (`import _ "net/http/pprof"`)
7. **Common optimizations:** reduce allocations, preallocate, use sync.Pool, avoid string concat
8. **Escape analysis** shows stack vs heap (`go build -gcflags="-m"`)
9. **Low-hanging fruit:** Often 80% improvement from 20% effort (focus on top hotspots)
10. **Measure after every change** (optimizations can make things worse!)

---

## Connections to Other Projects

- **Project 21 (race-detection-demo)**: pprof's mutex/block profiles detect contention
- **Project 24 (sync-mutex-vs-rwmutex)**: Benchmark to compare lock performance
- **Project 25 (atomic-counters-vs-mutex)**: Benchmark atomic vs mutex
- **Project 22 (worker-pool-with-backpressure)**: Profile to tune pool size
- **Project 07 (generic-lru-cache)**: Profile cache hit/miss patterns
- **Project 06 (worker-pool-wordcount)**: Benchmark parallel vs sequential

---

## Stretch Goals

1. **Build a CPU-bound benchmark suite** comparing sorting algorithms (profile to find O(n²) bottlenecks)
2. **Implement a memory leak detector** using heap profiling (detect growing allocations)
3. **Create a live profiling dashboard** (pprof + Prometheus + Grafana)
4. **Write a custom pprof visualization** (parse profiles, generate custom reports)
5. **Optimize a real-world application** (find and fix 10x performance improvement)
6. **Compare Go profiling to other languages** (pprof vs perf, valgrind, VisualVM)
7. **Implement a flamegraph renderer** from scratch (parse pprof protobuf format)

---

## Further Reading

- [Go pprof documentation](https://pkg.go.dev/runtime/pprof)
- [Profiling Go Programs (official blog)](https://go.dev/blog/pprof)
- [Brendan Gregg's Flamegraph Guide](https://www.brendangregg.com/flamegraphs.html)
- [Go memory profiling (Julia Evans)](https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/)
- [High Performance Go Workshop](https://dave.cheney.net/high-performance-go-workshop/dotgo-paris.html)
