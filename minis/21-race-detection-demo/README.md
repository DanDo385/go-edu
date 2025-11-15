# Project 21: Race Detection Demo - Finding and Fixing Data Races

## What Is This Project About?

This project teaches **data races** and how to detect and fix them using Go's built-in **race detector**. You'll learn:

1. **What data races are** (concurrent unsynchronized access to shared memory)
2. **Why races are dangerous** (undefined behavior, corruption, crashes)
3. **The Go race detector** (how it works, how to use it)
4. **Common race patterns** (loop variables, lazy initialization, double-checked locking)
5. **How to fix races** (mutexes, channels, atomic operations, copying)
6. **Race-free concurrency patterns** (immutability, ownership transfer, CSP)

By the end, you'll understand why data races are the most insidious concurrency bugs, how to detect them systematically, and how to write race-free concurrent code.

---

## The Fundamental Problem: Unsynchronized Shared Memory Access

### First Principles: What Is a Data Race?

A **data race** occurs when:
1. **Two or more goroutines** access the same memory location
2. **At least one access is a write**
3. **The accesses are not synchronized** (no happens-before relationship)

**Example (classic race):**
```go
var counter int  // Shared variable

// Goroutine 1
counter++  // Read, increment, write

// Goroutine 2
counter++  // Read, increment, write
```

**The problem:**
- `counter++` is **not atomic** (3 operations: read, add, write)
- Both goroutines can read the same value simultaneously
- Both write back the same incremented value
- Result: Lost update (counter=1 instead of 2)

**Assembly breakdown:**
```asm
; counter++ translates to:
MOV  AX, [counter]   ; Read current value
INC  AX              ; Increment
MOV  [counter], AX   ; Write back

; If two goroutines interleave:
G1: MOV  AX, [counter]   ; AX=0
G2: MOV  BX, [counter]   ; BX=0
G1: INC  AX              ; AX=1
G2: INC  BX              ; BX=1
G1: MOV  [counter], AX   ; counter=1
G2: MOV  [counter], BX   ; counter=1 (should be 2!)
```

---

## Why Data Races Are Dangerous

### 1. Undefined Behavior

In languages like C/C++, data races invoke **undefined behavior**. The compiler can:
- Reorder instructions
- Remove code it thinks is redundant
- Assume no race exists and optimize accordingly

**Result:** Unpredictable crashes, corruption, or seemingly impossible bugs.

### 2. Non-Atomic Writes Can Corrupt Data

Even in Go (with defined memory model), races cause corruption:

```go
type Config struct {
    Host string
    Port int
}

var config Config

// Goroutine 1 (writer)
config = Config{Host: "localhost", Port: 8080}

// Goroutine 2 (reader)
fmt.Println(config)  // Could print partial state!
// Example: {Host: "localhost", Port: 0}
// Or: {Host: "", Port: 8080}
```

**Why?** Struct assignment is **not atomic**. The write happens field-by-field.

### 3. Compiler/CPU Reordering

Modern CPUs and compilers reorder instructions for performance:

```go
var ready bool
var data int

// Writer
data = 42
ready = true

// Reader
if ready {
    fmt.Println(data)  // Might print 0!
}
```

**Why?** CPU can reorder the writes (`ready=true` before `data=42`).

### 4. Races Are Non-Deterministic

Races often "work fine" in testing but fail in production:
- Different CPU architectures have different memory models
- Different load patterns cause different interleavings
- Heisenbugs: Adding `fmt.Println` or debugger changes timing and "fixes" the race

**This makes races the hardest bugs to debug.**

---

## The Go Race Detector: How It Works

### What Is the Race Detector?

The **race detector** is a runtime tool that:
1. **Instruments your code** (tracks all memory accesses)
2. **Records synchronization events** (mutexes, channels, atomics)
3. **Detects happens-before violations** (unsynchronized accesses)
4. **Reports races with stack traces** (where the race occurred)

**Key point:** It's a **dynamic analysis tool**. It only detects races that **actually happen during execution**.

### The Happens-Before Relationship

Two events A and B are **synchronized** if:
- A happens-before B (e.g., send before receive)
- B happens-before A
- Both are protected by the same lock

**Examples of happens-before:**
- `mu.Lock()` in goroutine A → `mu.Unlock()` → `mu.Lock()` in goroutine B
- Send on channel in A → Receive in B
- `wg.Add()` in A → `wg.Done()` in A → `wg.Wait()` returns in B
- Goroutine creation: `go f()` → start of `f()`

**No happens-before = potential race.**

### How to Use the Race Detector

```bash
# Run your program with race detection
go run -race main.go

# Run tests with race detection
go test -race ./...

# Build a binary with race detection (slower, for testing)
go build -race -o myapp
./myapp
```

**Performance overhead:**
- **Memory:** ~5-10x increase (tracks shadow state)
- **CPU:** ~2-20x slower (instruments every access)

**Best practices:**
- ✅ Always run tests with `-race` in CI/CD
- ✅ Run load tests with `-race` on staging
- ❌ Don't deploy `-race` builds to production (too slow)

---

## Common Race Patterns (The Top 7)

### 1. Incrementing a Counter (Classic Race)

**Buggy code:**
```go
var counter int

func increment() {
    for i := 0; i < 1000; i++ {
        counter++  // RACE!
    }
}

func main() {
    go increment()
    go increment()
    time.Sleep(time.Second)
    fmt.Println(counter)  // Expected: 2000, Actual: varies!
}
```

**Race detector output:**
```
WARNING: DATA RACE
Write at 0x00c000010000 by goroutine 7:
  main.increment()
      /path/to/main.go:5 +0x4e

Previous write at 0x00c000010000 by goroutine 6:
  main.increment()
      /path/to/main.go:5 +0x4e
```

**Fixes:**

**Option 1: Mutex**
```go
var (
    counter int
    mu      sync.Mutex
)

func increment() {
    for i := 0; i < 1000; i++ {
        mu.Lock()
        counter++
        mu.Unlock()
    }
}
```

**Option 2: Atomic operations**
```go
var counter atomic.Int64

func increment() {
    for i := 0; i < 1000; i++ {
        counter.Add(1)  // Atomic, no lock needed
    }
}
```

**Option 3: Channel (message passing)**
```go
func main() {
    ch := make(chan int)

    // Worker goroutines
    go func() {
        for i := 0; i < 1000; i++ {
            ch <- 1
        }
    }()

    go func() {
        for i := 0; i < 1000; i++ {
            ch <- 1
        }
    }()

    // Aggregator (owns counter, no race)
    counter := 0
    for i := 0; i < 2000; i++ {
        counter += <-ch
    }
    fmt.Println(counter)
}
```

---

### 2. Loop Variable Capture (Goroutine Closure Race)

**Buggy code:**
```go
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i)  // RACE! Shares &i across goroutines
    }()
}
time.Sleep(time.Second)
// Output: Often prints "5 5 5 5 5" (all see final value)
```

**Why it's a race:**
- The loop variable `i` is **shared** across all goroutines
- The loop increments `i` while goroutines read it
- Unsynchronized read + write = race

**Fixes:**

**Option 1: Pass as argument (copy)**
```go
for i := 0; i < 5; i++ {
    go func(id int) {
        fmt.Println(id)  // Each goroutine gets a copy
    }(i)
}
```

**Option 2: Shadow variable (Go 1.22+ auto-fixes this)**
```go
for i := 0; i < 5; i++ {
    i := i  // Create new variable per iteration
    go func() {
        fmt.Println(i)
    }()
}
```

---

### 3. Map Concurrent Read/Write

**Buggy code:**
```go
var cache = make(map[string]int)

// Goroutine 1
cache["key"] = 42  // Write

// Goroutine 2
value := cache["key"]  // Read

// RACE! Map operations are NOT thread-safe
// Can cause "concurrent map read and map write" panic
```

**Why maps panic:**
- Go's map implementation is **not concurrent-safe**
- Concurrent writes can corrupt internal hash table structure
- Go **panics immediately** to prevent silent corruption

**Fixes:**

**Option 1: Mutex**
```go
var (
    cache = make(map[string]int)
    mu    sync.Mutex
)

// Write
mu.Lock()
cache["key"] = 42
mu.Unlock()

// Read
mu.Lock()
value := cache["key"]
mu.Unlock()
```

**Option 2: RWMutex (better for read-heavy workloads)**
```go
var (
    cache = make(map[string]int)
    mu    sync.RWMutex
)

// Write (exclusive)
mu.Lock()
cache["key"] = 42
mu.Unlock()

// Read (shared, allows multiple readers)
mu.RLock()
value := cache["key"]
mu.RUnlock()
```

**Option 3: sync.Map (built-in concurrent map)**
```go
var cache sync.Map

// Write
cache.Store("key", 42)

// Read
value, ok := cache.Load("key")

// Use when: Many reads, few writes, disjoint key sets
```

---

### 4. Lazy Initialization (Double-Checked Locking Anti-Pattern)

**Buggy code:**
```go
var (
    config *Config
    mu     sync.Mutex
)

func getConfig() *Config {
    if config == nil {  // Read without lock - RACE!
        mu.Lock()
        if config == nil {
            config = loadConfig()  // Write
        }
        mu.Unlock()
    }
    return config
}
```

**Why it's a race:**
- First `if config == nil` reads **without lock**
- Another goroutine could be writing `config` inside the lock
- Race between unsynchronized read and write

**Why it seems to "work":**
- Pointer reads/writes are atomic on most architectures
- But the race detector catches it (correctly!)
- Could still fail with subtle bugs (e.g., seeing partial struct state)

**Fixes:**

**Option 1: Always lock (simple, slightly slower)**
```go
func getConfig() *Config {
    mu.Lock()
    defer mu.Unlock()

    if config == nil {
        config = loadConfig()
    }
    return config
}
```

**Option 2: sync.Once (idiomatic, fast)**
```go
var (
    config *Config
    once   sync.Once
)

func getConfig() *Config {
    once.Do(func() {
        config = loadConfig()
    })
    return config
}
```

**Option 3: atomic.Value (for frequent reads)**
```go
var config atomic.Value  // Stores *Config

func getConfig() *Config {
    cfg := config.Load()
    if cfg == nil {
        newCfg := loadConfig()
        config.Store(newCfg)
        return newCfg
    }
    return cfg.(*Config)
}
```

---

### 5. Slice Append Race

**Buggy code:**
```go
var results []int

func worker(id int) {
    results = append(results, id)  // RACE!
}

func main() {
    for i := 0; i < 10; i++ {
        go worker(i)
    }
    time.Sleep(time.Second)
    fmt.Println(results)  // Missing results or panic
}
```

**Why it's a race:**
- `append` **reads** the slice header (len, cap, ptr)
- `append` **writes** a new slice header if it grows
- Multiple goroutines race on the slice header

**Fixes:**

**Option 1: Mutex**
```go
var (
    results []int
    mu      sync.Mutex
)

func worker(id int) {
    mu.Lock()
    results = append(results, id)
    mu.Unlock()
}
```

**Option 2: Channel (collect results)**
```go
func worker(id int, ch chan<- int) {
    ch <- id
}

func main() {
    ch := make(chan int, 10)

    for i := 0; i < 10; i++ {
        go worker(i, ch)
    }

    results := make([]int, 0, 10)
    for i := 0; i < 10; i++ {
        results = append(results, <-ch)
    }
}
```

**Option 3: Preallocate and use atomic index**
```go
var (
    results = make([]int, 10)
    idx     atomic.Int64
)

func worker(id int) {
    i := idx.Add(1) - 1
    results[i] = id
}
```

---

### 6. Struct Field Race

**Buggy code:**
```go
type Stats struct {
    Requests int
    Errors   int
}

var stats Stats

// Goroutine 1
stats.Requests++  // RACE!

// Goroutine 2
stats.Errors++    // RACE!
```

**Why it's a race (even on different fields!):**
- On some architectures, fields might share the same cache line
- False sharing: CPU cache coherency protocol causes contention
- Race detector flags it because struct is treated as a unit

**Fixes:**

**Option 1: Mutex (protects all fields)**
```go
var (
    stats Stats
    mu    sync.Mutex
)

mu.Lock()
stats.Requests++
mu.Unlock()
```

**Option 2: Atomic fields**
```go
type Stats struct {
    Requests atomic.Int64
    Errors   atomic.Int64
}

stats.Requests.Add(1)  // No lock needed
```

**Option 3: Separate variables (avoid false sharing)**
```go
var (
    requests atomic.Int64
    errors   atomic.Int64
)

requests.Add(1)
```

---

### 7. Channel Iteration Race

**Buggy code:**
```go
ch := make(chan int)

// Producer
go func() {
    for i := 0; i < 10; i++ {
        ch <- i
    }
    close(ch)
}()

// Consumer 1
go func() {
    for v := range ch {
        fmt.Println("C1:", v)
    }
}()

// Consumer 2
go func() {
    for v := range ch {
        fmt.Println("C2:", v)  // RACE on channel state!
    }
}()
```

**Why it's a race:**
- Ranging over a **closed** channel from multiple goroutines is undefined
- If both consumers read at the same time, they race on channel state
- Can lead to missed values or panics

**Actually, this is NOT a race!** Channels are thread-safe. Both consumers can range safely. But there's a **logic bug**—you don't know which consumer gets which value.

**If you want deterministic distribution:**

**Option 1: Worker pool with WaitGroup**
```go
var wg sync.WaitGroup

for i := 0; i < 2; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        for v := range ch {
            fmt.Printf("C%d: %d\n", id, v)
        }
    }(i)
}

wg.Wait()
```

**Option 2: Fan-out pattern with separate channels**
```go
ch1, ch2 := make(chan int), make(chan int)

go func() {
    for i := 0; i < 10; i++ {
        if i%2 == 0 {
            ch1 <- i
        } else {
            ch2 <- i
        }
    }
    close(ch1)
    close(ch2)
}()
```

---

## How to Fix Data Races: The Strategies

### Strategy 1: Mutual Exclusion (Locks)

**Use when:** Multiple goroutines need to read/write shared state.

```go
var (
    balance int
    mu      sync.Mutex
)

func deposit(amount int) {
    mu.Lock()
    balance += amount
    mu.Unlock()
}

func withdraw(amount int) bool {
    mu.Lock()
    defer mu.Unlock()

    if balance >= amount {
        balance -= amount
        return true
    }
    return false
}
```

**Pros:**
- Simple, straightforward
- Works for complex state

**Cons:**
- Can cause contention (goroutines block)
- Risk of deadlocks with multiple locks

---

### Strategy 2: Channel Communication

**Use when:** Goroutines can send messages instead of sharing state.

```go
type Command struct {
    Op     string  // "deposit" or "withdraw"
    Amount int
    Result chan bool
}

func accountManager(ch <-chan Command) {
    balance := 0

    for cmd := range ch {
        switch cmd.Op {
        case "deposit":
            balance += cmd.Amount
            cmd.Result <- true
        case "withdraw":
            if balance >= cmd.Amount {
                balance -= cmd.Amount
                cmd.Result <- true
            } else {
                cmd.Result <- false
            }
        }
    }
}
```

**Pros:**
- No shared state
- Clear ownership (account manager owns balance)
- Composable

**Cons:**
- More complex for simple cases
- Overhead of channel operations

---

### Strategy 3: Atomic Operations

**Use when:** Simple numeric operations on shared variables.

```go
var counter atomic.Int64

func increment() {
    counter.Add(1)
}

func read() int64 {
    return counter.Load()
}
```

**Pros:**
- No locks, no contention
- Very fast (single CPU instruction)

**Cons:**
- Limited to simple types (int, uint, pointer)
- Can't protect complex operations

---

### Strategy 4: Immutability

**Use when:** Data doesn't need to change after creation.

```go
type Config struct {
    Host string
    Port int
}

var config *Config  // Pointer is swapped atomically

func updateConfig(newConfig Config) {
    atomic.StorePointer(
        (*unsafe.Pointer)(unsafe.Pointer(&config)),
        unsafe.Pointer(&newConfig),
    )
}

// Better: Use atomic.Value
var configValue atomic.Value  // Stores *Config

func updateConfig(newConfig *Config) {
    configValue.Store(newConfig)
}

func getConfig() *Config {
    return configValue.Load().(*Config)
}
```

**Pros:**
- Readers never block
- No need to copy data

**Cons:**
- Requires copying on updates
- Pointer swap is tricky (use `atomic.Value`)

---

### Strategy 5: Confinement (No Sharing)

**Use when:** Each goroutine can have its own data.

```go
func main() {
    // Each worker has its own counter
    var wg sync.WaitGroup

    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()

            counter := 0  // Not shared!
            for j := 0; j < 1000; j++ {
                counter++
            }
            fmt.Println(counter)  // Always 1000
        }()
    }

    wg.Wait()
}
```

**Pros:**
- Zero synchronization cost
- Simplest approach

**Cons:**
- Can't share results easily
- Need aggregation step if you need global state

---

## Race-Free Concurrency Patterns

### Pattern 1: Single-Owner Pattern

**Idea:** One goroutine owns the data, others send requests.

```go
type Request struct {
    Op    string
    Value int
    Reply chan int
}

func dataOwner(requests <-chan Request) {
    data := make(map[string]int)

    for req := range requests {
        switch req.Op {
        case "get":
            req.Reply <- data[req.Value]
        case "set":
            data[req.Value] = req.Value
            req.Reply <- 0
        }
    }
}
```

---

### Pattern 2: Pipeline Pattern

**Idea:** Pass data through stages, each owned by a goroutine.

```go
func main() {
    // Stage 1: Generate numbers
    gen := func(nums ...int) <-chan int {
        out := make(chan int)
        go func() {
            for _, n := range nums {
                out <- n
            }
            close(out)
        }()
        return out
    }

    // Stage 2: Square numbers
    sq := func(in <-chan int) <-chan int {
        out := make(chan int)
        go func() {
            for n := range in {
                out <- n * n
            }
            close(out)
        }()
        return out
    }

    // Compose pipeline
    nums := gen(2, 3, 4)
    squares := sq(nums)

    for v := range squares {
        fmt.Println(v)  // 4, 9, 16
    }
}
```

---

### Pattern 3: Worker Pool Pattern

**Idea:** Fixed number of workers, bounded concurrency.

```go
func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)

    // Start workers
    var wg sync.WaitGroup
    for w := 0; w < 5; w++ {
        wg.Add(1)
        go worker(w, jobs, results, &wg)
    }

    // Send jobs
    for i := 0; i < 100; i++ {
        jobs <- i
    }
    close(jobs)

    // Wait and close results
    go func() {
        wg.Wait()
        close(results)
    }()

    // Collect results
    for r := range results {
        fmt.Println(r)
    }
}

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    for j := range jobs {
        results <- j * j
    }
}
```

---

## Race Detector Limitations

### 1. Only Detects Executed Races

The race detector is **dynamic**—it only finds races that actually happen.

```go
// If this branch never executes in your test, race won't be detected
if veryRareCondition {
    counter++  // RACE (but not executed)
}
```

**Solution:** Achieve high code coverage in tests.

### 2. False Negatives (Missed Races)

If a race doesn't happen during your test run, it won't be detected.

**Solution:** Use techniques like:
- Randomized timing (`time.Sleep(rand.Duration())`)
- Stress tests (run many iterations)
- Multiple GOMAXPROCS values

### 3. No False Positives

The race detector is **sound**—if it reports a race, it's a real race (not a false alarm).

### 4. Performance Overhead

Running with `-race` is slow (2-20x). Don't use in production.

**Solution:** Run race detection in CI/CD, staging, and load tests.

---

## How to Run

```bash
# Run the race demo
cd minis/21-race-detection-demo
go run -race cmd/race-demo/main.go

# Run exercises
cd exercise
go test -v

# Run with race detector (catch bugs!)
go test -race -v

# Run specific exercise
go test -run TestCounterRace -race -v
```

---

## Expected Output (Demo Program)

```
=== Race Detection Demonstration ===

[1] Counter Increment Race
Running buggy version...
WARNING: DATA RACE
Write at 0x00c000018090 by goroutine 7:
  main.buggyCounter()
      /path/to/main.go:10 +0x3c
...
Running fixed version (mutex)...
Final counter: 2000 (correct!)

[2] Map Concurrent Access
Running buggy version...
fatal error: concurrent map writes
...

[3] Loop Variable Capture
Running buggy version...
Output: 5 5 5 5 5 (all print final value - WRONG)
Running fixed version...
Output: 0 1 2 3 4 (correct!)

All demonstrations complete!
```

---

## Key Takeaways

1. **Data race = unsynchronized access** (at least one write) to shared memory
2. **Races cause undefined behavior**, corruption, and crashes (not just "wrong values")
3. **Use `-race` in tests** (catches races at runtime with stack traces)
4. **Common patterns:** counter increment, map access, loop capture, lazy init, slice append
5. **Fixes:** mutexes, channels, atomics, immutability, confinement
6. **Race-free patterns:** single owner, pipelines, worker pools
7. **Race detector is dynamic** (only finds executed races—achieve high coverage!)
8. **Channels > locks for Go** (prefer message passing to shared memory)

---

## Connections to Other Projects

- **Project 18 (goroutines-1M-demo)**: Where concurrency begins (goroutines need synchronization!)
- **Project 19 (channels-basics)**: Race-free communication primitive
- **Project 20 (select-fanin-fanout)**: Advanced race-free patterns
- **Project 24 (sync-mutex-vs-rwmutex)**: Deep dive into locking primitives
- **Project 25 (atomic-operations)**: Lock-free synchronization
- **Project 28 (pprof-cpu-mem-benchmarks)**: Profiling shows contention hotspots

---

## Stretch Goals

1. **Write a race detector from scratch** (simplified version using shadow state)
2. **Benchmark lock vs atomic vs channel** for different workloads
3. **Implement double-checked locking correctly** (using `sync.Once` and `atomic.Value`)
4. **Create a race-free LRU cache** (compare mutex vs channels vs sync.Map)
5. **Analyze false sharing** using CPU cache line alignment (`cacheline` package)
6. **Port a racy C++ program to Go** and use the race detector to find all races
