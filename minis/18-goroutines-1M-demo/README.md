# Project 18: Goroutines - Running 1 Million Concurrent Tasks

## What Is This Project About?

This project demonstrates **goroutines**—Go's lightweight concurrency primitive that makes it trivial to run **millions of concurrent tasks** on commodity hardware. You'll learn:

1. **What goroutines are** (user-space threads vs OS threads)
2. **How they differ from OS threads** (stack size, context switching, scheduling)
3. **Stack growth mechanics** (how 2KB stacks can grow to megabytes)
4. **Runtime scheduling** (the Go scheduler, GOMAXPROCS, work stealing)
5. **Why 1 million goroutines is possible** (memory footprint calculations)
6. **Communication patterns** (channels, sync primitives, graceful shutdown)

By the end, you'll understand why Go can run 1,000,000 goroutines with less than 2GB of RAM, while trying the same with OS threads would require 1TB+ of memory and crash your system.

---

## The Fundamental Problem: OS Threads Are Too Heavy

### First Principles: What Is an OS Thread?

An **operating system thread** is a unit of execution managed by the OS kernel. Each thread:

1. **Has a fixed stack size** (typically 1-8 MB, allocated upfront)
2. **Requires kernel involvement** for creation, scheduling, and context switching
3. **Has OS-level overhead** (thread control block, kernel structures)

**Example (creating 10,000 threads):**
```c
// In C/pthreads:
pthread_t threads[10000];
for (int i = 0; i < 10000; i++) {
    pthread_create(&threads[i], NULL, worker, NULL);
}
// Memory required: 10,000 × 2 MB = 20 GB (just for stacks!)
// Context switch time: ~1-2 microseconds per switch
```

**The problems:**
- **Memory explosion**: Each thread reserves megabytes of stack space upfront
- **Slow creation**: Kernel must allocate resources, set up page tables
- **Expensive context switches**: CPU must save/restore registers, flush TLB
- **Limited scalability**: Most systems limit threads to thousands, not millions

This is why traditional server models (one thread per connection) struggle at high concurrency.

---

## What Are Goroutines? (The Core Concept)

A **goroutine** is a **user-space thread** managed entirely by the Go runtime, not the OS kernel.

### Key Characteristics

| Feature              | OS Thread           | Goroutine           |
|----------------------|---------------------|---------------------|
| **Stack size**       | 1-8 MB (fixed)      | 2 KB (grows dynamically) |
| **Creation cost**    | ~1-2 μs (kernel)    | ~1-2 ns (runtime)   |
| **Context switch**   | ~1-2 μs (kernel)    | ~200 ns (runtime)   |
| **Memory overhead**  | ~8 MB per thread    | ~2-4 KB per goroutine |
| **Maximum count**    | ~10,000s            | ~millions           |
| **Scheduler**        | OS kernel           | Go runtime (userspace) |

**The magic:** Goroutines are **multiplexed** onto a small number of OS threads (typically equal to the number of CPU cores).

### Creating a Goroutine

```go
// Synchronous function call (blocks)
doWork()

// Asynchronous goroutine (runs concurrently)
go doWork()
```

**That's it.** No thread objects, no thread pools, no manual management. Just prefix any function call with `go`.

---

## Goroutines vs OS Threads: The Technical Details

### 1. Stack Management

**OS Threads:**
- Allocate a **fixed stack** (e.g., 2 MB) at creation time
- Stack size is determined by OS defaults (cannot grow without re-creating thread)
- Wastes memory for threads that only use a few KB

**Goroutines:**
- Start with a **tiny stack** (2 KB in Go 1.4+)
- **Grow dynamically** as needed (explained in next section)
- Can shrink back down when stack usage decreases

**Memory comparison for 1 million tasks:**
```
OS threads:  1,000,000 × 2 MB = 2 TB (impossible on most systems)
Goroutines:  1,000,000 × 2 KB = 2 GB (easily fits in RAM)
```

### 2. Context Switching

**OS Threads:**
- Switching requires a **kernel mode transition** (expensive)
- CPU must:
  1. Save all registers (including floating-point, SSE, etc.)
  2. Flush TLB (translation lookaside buffer)
  3. Load new thread's registers and page tables
- Cost: ~1-2 microseconds

**Goroutines:**
- Switching happens **entirely in user space** (no kernel involvement)
- Go runtime only saves/restores 3 registers (PC, SP, BP)
- No TLB flush (all goroutines share the same address space)
- Cost: ~200 nanoseconds (10x faster!)

### 3. Scheduling

**OS Threads:**
- Scheduled by the **kernel scheduler** (preemptive, time-sliced)
- Kernel has no knowledge of your application's semantics
- Context switches can happen at any time (even mid-critical-section)

**Goroutines:**
- Scheduled by the **Go runtime scheduler** (cooperative + preemptive)
- Scheduler understands Go semantics (channels, mutexes, etc.)
- Switches at known-safe points (function calls, channel operations, GC)

---

## Stack Growth Mechanics: How 2KB Becomes Megabytes

### The Segmented Stack Problem (Go 1.3 and earlier)

Early Go versions used **segmented stacks**:
1. Start with a small stack (e.g., 4 KB)
2. If a function needs more space, allocate a **new segment** and link it
3. If the function returns, unlink and free the segment

**The problem (hot split):**
```go
func recursive(n int) {
    if n == 0 {
        return
    }
    recursive(n - 1)  // If this crosses a segment boundary...
}

// Tight loop:
for {
    recursive(10)  // Allocate segment
    // Return, free segment
    recursive(10)  // Allocate segment again!
}
// Thrashes allocator, terrible performance
```

### The Solution: Stack Copying (Go 1.4+)

Modern Go uses **stack copying**:
1. Goroutine starts with a **contiguous 2 KB stack**
2. If a function call would overflow, the runtime:
   - Allocates a **new, larger stack** (typically 2x the old size)
   - **Copies** all data from the old stack to the new stack
   - **Updates all pointers** to point to the new stack
   - Frees the old stack
3. The goroutine continues with the larger stack

**Example:**
```go
func deep(n int) {
    if n == 0 {
        return
    }
    var buf [1024]byte  // Allocate stack space
    deep(n - 1)
}

go deep(1000)
// Initial stack: 2 KB
// After ~100 calls, stack grows to 4 KB
// After ~200 calls, stack grows to 8 KB
// ...continues growing as needed
```

**How pointer updates work:**
- Every goroutine's stack is scanned by the GC (already tracked)
- When copying, the runtime adjusts all pointers to the new addresses
- This is **safe** because Go has no pointer arithmetic (unlike C)

---

## Runtime Scheduling: The Go Scheduler

### The M:N Threading Model

Go uses an **M:N scheduler**, where:
- **M** = OS threads (machines)
- **N** = Goroutines
- **P** = Logical processors (contexts)

**The relationship:**
```
Goroutines (G) ──> Logical Processors (P) ──> OS Threads (M)
   (millions)         (GOMAXPROCS)              (few)
```

### The GMP Model

**G (Goroutine):**
- The task you want to run
- Contains: stack, instruction pointer, other metadata

**M (Machine):**
- An OS thread
- Executes goroutines by running their code

**P (Processor):**
- A scheduling context (think: "CPU core context")
- Owns a **local run queue** of goroutines
- Number of Ps = `GOMAXPROCS` (defaults to number of CPU cores)

**Visual:**
```
┌─────────────────────────────────────┐
│   Global Run Queue (overflow)       │
│   [G] [G] [G]                       │
└─────────────────────────────────────┘
           │
           ↓ (work stealing)
┌──────────────────┐  ┌──────────────────┐
│  P1              │  │  P2              │
│  Local Queue:    │  │  Local Queue:    │
│  [G] [G] [G] [G] │  │  [G] [G] [G]     │
│        ↓         │  │        ↓         │
│       M1         │  │       M2         │
│    (OS Thread)   │  │    (OS Thread)   │
└──────────────────┘  └──────────────────┘
```

### Scheduling Decisions

**When does a goroutine yield?**
1. **Explicit yield**: `runtime.Gosched()` (rare)
2. **Channel operation**: Send/receive on a channel
3. **System call**: File I/O, network I/O
4. **Garbage collection**: GC runs
5. **Function call**: (Preemption point, can be interrupted)

**Work stealing:**
- If a P's local queue is empty, it **steals** goroutines from other Ps
- This balances load across CPU cores automatically

### GOMAXPROCS

Controls the number of Ps (logical processors):

```go
import "runtime"

// Default: number of CPU cores
runtime.GOMAXPROCS(0)  // Returns current value

// Set to 4 (only 4 goroutines run in parallel)
runtime.GOMAXPROCS(4)
```

**When to change:**
- **CPU-bound work**: Use all cores (default is good)
- **I/O-bound work**: May benefit from more than cores (but rarely needed)
- **Testing**: Set to 1 to force serial execution (useful for debugging)

---

## How 1 Million Goroutines Is Possible

### Memory Footprint Calculation

**Minimal goroutine (idle, just blocking on a channel):**
- **Stack**: 2 KB (minimum, may grow if used)
- **G struct**: ~320 bytes (metadata: stack pointers, status, etc.)
- **Overhead**: ~100 bytes (scheduler data structures)

**Total per goroutine: ~2.5 KB**

**1 million goroutines:**
```
1,000,000 × 2.5 KB = 2.5 GB
```

This easily fits in modern RAM (most servers have 16-128 GB).

### Real-World Example

```go
func main() {
    for i := 0; i < 1_000_000; i++ {
        go func() {
            time.Sleep(10 * time.Minute)  // Just block
        }()
    }
    time.Sleep(10 * time.Minute)
}
```

**Actual memory usage:**
- Initial: ~50 MB (runtime overhead)
- After 1M goroutines: ~2.5 GB (stack + metadata)
- **CPU usage: near 0%** (all goroutines are sleeping!)

### What About Active Goroutines?

If each goroutine is **actively doing work**, memory grows:
- **Stack growth**: If recursion or large local variables, stacks grow
- **Heap allocations**: Goroutines allocate objects that live on the heap

**Example (each goroutine allocates 1 KB):**
```
1,000,000 × (2 KB stack + 1 KB heap) = 3 GB
```

Still very manageable.

---

## Communication and Synchronization

### Channels (Recommended)

Goroutines communicate via **channels** (type-safe message queues):

```go
ch := make(chan int)

// Sender goroutine
go func() {
    ch <- 42  // Send value
}()

// Receiver goroutine
result := <-ch  // Receive value
```

**Benefits:**
- **Type-safe**: Compiler enforces types
- **Synchronized**: Sender blocks until receiver is ready (or use buffered channels)
- **Composable**: Use `select` to handle multiple channels

### WaitGroups (Synchronization)

For waiting on multiple goroutines to finish:

```go
var wg sync.WaitGroup

for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        // Do work...
    }(i)
}

wg.Wait()  // Block until all goroutines call Done()
```

### Context (Cancellation)

For canceling goroutines:

```go
ctx, cancel := context.WithCancel(context.Background())

go func() {
    for {
        select {
        case <-ctx.Done():
            return  // Exit when canceled
        default:
            // Do work...
        }
    }
}()

// Later:
cancel()  // Signal all goroutines to stop
```

---

## Graceful Shutdown

**Problem:** How do you cleanly stop 1 million goroutines?

**Solution 1: Broadcast via context**
```go
ctx, cancel := context.WithCancel(context.Background())

for i := 0; i < 1_000_000; i++ {
    go worker(ctx, i)
}

// On shutdown:
cancel()  // All workers check ctx.Done() and exit
```

**Solution 2: Close a channel**
```go
done := make(chan struct{})

for i := 0; i < 1_000_000; i++ {
    go func() {
        <-done  // Block until closed
        // Cleanup...
    }()
}

// On shutdown:
close(done)  // Unblocks ALL receivers
```

**Solution 3: WaitGroup for completion tracking**
```go
var wg sync.WaitGroup

for i := 0; i < 1_000_000; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        // Work...
    }()
}

wg.Wait()  // Block until all finish
```

---

## Common Pitfalls

### Pitfall 1: Goroutine Leaks

```go
// BAD: Goroutine blocks forever
func leak() {
    ch := make(chan int)
    go func() {
        ch <- 42  // No receiver! Blocks forever
    }()
    // Goroutine never exits (memory leak)
}
```

**Fix:** Use buffered channels or ensure receivers exist.

### Pitfall 2: Capturing Loop Variables

```go
// BAD: All goroutines see i=10
for i := 0; i < 10; i++ {
    go func() {
        fmt.Println(i)  // Captures &i, not the value!
    }()
}

// GOOD: Pass i as argument
for i := 0; i < 10; i++ {
    go func(id int) {
        fmt.Println(id)
    }(i)
}
```

### Pitfall 3: Unbounded Goroutine Creation

```go
// BAD: Could create millions of goroutines
for _, item := range hugeList {
    go process(item)  // No limit!
}
```

**Fix:** Use a worker pool with bounded concurrency.

### Pitfall 4: Ignoring Panics

```go
// BAD: Panic in goroutine crashes the program
go func() {
    panic("oops")  // Kills entire program!
}()
```

**Fix:** Use `defer recover()` in goroutines.

---

## Performance Characteristics

### Goroutine Creation

**Benchmark:**
```
Creating 1,000,000 goroutines: ~100 ms
Creating 10,000 OS threads: ~10 seconds (100x slower)
```

### Context Switching

**Benchmark:**
```
Goroutine switch: ~200 ns
OS thread switch: ~1-2 μs (5-10x slower)
```

### Memory Efficiency

**Benchmark:**
```
1M goroutines (idle): ~2.5 GB
1M OS threads: ~2 TB (impossible)
```

---

## How to Run

```bash
# Run the 1 million goroutine demo
cd minis/18-goroutines-1M-demo
go run cmd/goroutine-demo/main.go

# Run exercises
cd exercise
go test -v

# Monitor memory usage during demo
go run cmd/goroutine-demo/main.go &
PID=$!
while kill -0 $PID 2>/dev/null; do
    ps -p $PID -o rss,vsz,comm
    sleep 1
done
```

---

## Expected Output (Demo Program)

```
=== Goroutine Demonstration ===
Starting 1,000,000 goroutines...

Memory before: 5 MB
Launching goroutines...
Memory after: 2,500 MB
Time to launch: 150 ms

All goroutines launched!
Active goroutines: 1,000,001

Testing communication...
Received 100,000 messages in 50 ms

Initiating graceful shutdown...
All goroutines stopped.
Final memory: 10 MB
```

---

## Key Takeaways

1. **Goroutines are user-space threads**, not OS threads (tiny, fast, millions possible)
2. **Start with 2 KB stacks**, grow dynamically via stack copying (not segmentation)
3. **Scheduled by Go runtime**, not the kernel (faster context switches, cooperative)
4. **GOMAXPROCS controls parallelism** (defaults to CPU cores)
5. **1 million goroutines ≈ 2.5 GB**, 1 million threads ≈ 2 TB (1000x difference!)
6. **Communicate via channels**, synchronize via WaitGroups, cancel via context
7. **Graceful shutdown** requires broadcasting (context, closed channels)
8. **Watch for leaks** (blocked goroutines, captured loop variables)

---

## Connections to Other Projects

- **Project 06 (worker-pool-wordcount)**: Used goroutines in a worker pool pattern
- **Project 16 (context-cancellation-timeouts)**: Context for graceful goroutine shutdown
- **Project 19 (channels-basics)**: Deep dive into channel communication
- **Project 20 (select-fanin-fanout)**: Advanced goroutine patterns (fan-in, fan-out)
- **Project 24 (sync-mutex-vs-rwmutex)**: Synchronizing shared state between goroutines
- **Project 28 (pprof-cpu-mem-benchmarks)**: Profiling goroutine usage and memory

---

## Stretch Goals

1. **Measure stack growth** using `runtime.Stack()` and calculate actual sizes
2. **Benchmark goroutine creation** vs OS thread creation (pthread or Java threads)
3. **Visualize the GMP model** using `runtime.NumGoroutine()`, `GOMAXPROCS`, etc.
4. **Implement a goroutine pool** with bounded concurrency and backpressure
5. **Detect goroutine leaks** using `runtime.NumGoroutine()` before/after operations
