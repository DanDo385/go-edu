# Project 19: Channels Basics - Communication Between Goroutines

## What Is This Project About?

This project teaches **channels**—Go's fundamental primitive for safe communication between goroutines. You'll learn:

1. **What channels are** (typed message queues with synchronization)
2. **Unbuffered vs buffered channels** (blocking semantics, capacity)
3. **Send and receive operations** (syntax, blocking behavior)
4. **Channel closing** (signaling completion, detecting closure)
5. **Range over channels** (consuming values until closed)
6. **Select statement** (multiplexing channel operations)
7. **Pipeline patterns** (composing stages with channels)

By the end, you'll understand why channels are Go's idiomatic way to share memory by communicating, rather than communicate by sharing memory.

---

## The Fundamental Problem: Sharing Data Between Goroutines

### First Principles: The Concurrency Challenge

When multiple goroutines need to share data, you face two fundamental approaches:

**Approach 1: Shared Memory (Traditional)**
```go
var counter int  // Shared variable
var mu sync.Mutex

// Goroutine 1
mu.Lock()
counter++
mu.Unlock()

// Goroutine 2
mu.Lock()
counter++
mu.Unlock()
```

**Problems:**
- **Easy to forget locks** → race conditions, data corruption
- **Deadlocks** when acquiring multiple locks in wrong order
- **Hard to reason about** complex locking patterns
- **Not composable** (difficult to build abstractions)

**Approach 2: Message Passing (Go's Way)**
```go
ch := make(chan int)

// Goroutine 1
ch <- 42  // Send value

// Goroutine 2
value := <-ch  // Receive value
```

**Benefits:**
- **No shared state** → no locks needed
- **Type-safe** → compiler enforces types
- **Synchronization built-in** → sender blocks until receiver is ready
- **Composable** → easy to build pipelines and patterns

**Go's Philosophy:**
> "Don't communicate by sharing memory; share memory by communicating."

---

## What Is a Channel? (The Core Concept)

A **channel** is a **typed, thread-safe queue** that allows goroutines to send and receive values.

### Creating Channels

```go
// Unbuffered channel (capacity 0)
ch := make(chan int)

// Buffered channel (capacity N)
ch := make(chan int, 100)
```

### Basic Operations

```go
ch := make(chan int)

// Send (blocks until receiver is ready)
ch <- 42

// Receive (blocks until sender sends)
value := <-ch

// Receive and check if channel is closed
value, ok := <-ch
if !ok {
    // Channel is closed
}

// Close channel (signals no more values)
close(ch)
```

### Type Signature

Channels have a **direction** in their type:

```go
chan T      // Can send and receive T
chan<- T    // Send-only (can only send T)
<-chan T    // Receive-only (can only receive T)
```

**Why directional types?**
- **API clarity**: Function signature shows intent
- **Compiler enforcement**: Prevents misuse
- **Best practice**: Pass channels with restricted permissions

**Example:**
```go
// Producer only sends
func producer(out chan<- int) {
    out <- 42
    close(out)
}

// Consumer only receives
func consumer(in <-chan int) {
    value := <-in
    fmt.Println(value)
}

ch := make(chan int)
go producer(ch)  // chan int implicitly converts to chan<- int
consumer(ch)     // chan int implicitly converts to <-chan int
```

---

## Unbuffered vs Buffered Channels

### Unbuffered Channels (Capacity 0)

**Creation:**
```go
ch := make(chan int)  // No capacity argument
```

**Behavior:**
- **Send blocks** until a receiver is ready
- **Receive blocks** until a sender sends
- **Synchronization point**: Send and receive happen simultaneously

**Visual:**
```
Sender                  Receiver
  |                        |
  | ch <- 42               |
  | (BLOCKS)               |
  |                        | value := <-ch
  | ← UNBLOCKS ─────────→  | (receives 42)
  | (send completes)       |
```

**Use cases:**
- Signaling events (completion, errors)
- Synchronizing goroutines
- Handoff patterns (one sender, one receiver)

**Example:**
```go
done := make(chan struct{})  // Unbuffered signal channel

go func() {
    doWork()
    done <- struct{}{}  // Signal completion (blocks until main receives)
}()

<-done  // Wait for signal (blocks until worker sends)
fmt.Println("Work complete!")
```

### Buffered Channels (Capacity > 0)

**Creation:**
```go
ch := make(chan int, 3)  // Capacity of 3
```

**Behavior:**
- **Send blocks** only when buffer is full
- **Receive blocks** only when buffer is empty
- **Asynchronous**: Sender and receiver don't need to synchronize

**Visual:**
```
Channel buffer: [_, _, _] (capacity 3)

Send 1:   [1, _, _]  ← Send succeeds immediately
Send 2:   [1, 2, _]  ← Send succeeds immediately
Send 3:   [1, 2, 3]  ← Send succeeds immediately
Send 4:   BLOCKS (buffer full)

Receive:  [_, 2, 3]  ← Receive gets 1, send 4 can now proceed
```

**Use cases:**
- Decoupling producer/consumer speeds
- Work queues (bounded capacity prevents unbounded growth)
- Batching (accumulate values before processing)

**Example:**
```go
jobs := make(chan int, 100)  // Buffered queue

// Producer doesn't block until 100 jobs queued
go func() {
    for i := 0; i < 1000; i++ {
        jobs <- i  // Blocks only when buffer full
    }
    close(jobs)
}()

// Consumer processes at its own pace
for job := range jobs {
    process(job)
}
```

### Choosing Buffer Size

**Zero (unbuffered):**
- Strictest synchronization
- Use when you want guaranteed handoff

**Small (1-100):**
- Smooth out temporary speed differences
- Bounded memory usage

**Large (1000+):**
- Act as a queue/buffer
- Risk: Can hide backpressure issues

**Rule of thumb:** Start with unbuffered, add buffering only if profiling shows benefit.

---

## Send and Receive Semantics

### Send Operation: `ch <- value`

**What happens:**
1. **Unbuffered**: Blocks until a receiver is ready, then transfers value and unblocks
2. **Buffered (not full)**: Copies value to buffer, returns immediately
3. **Buffered (full)**: Blocks until space available, then copies value

**Panics if:**
- Sending on a **closed** channel → `panic: send on closed channel`
- Sending on a **nil** channel → blocks forever (deadlock)

**Example:**
```go
ch := make(chan int, 2)

ch <- 1  // Success (buffer: [1])
ch <- 2  // Success (buffer: [1, 2])
ch <- 3  // BLOCKS (buffer full)

close(ch)
ch <- 4  // PANIC: send on closed channel
```

### Receive Operation: `<-ch`

**What happens:**
1. **Unbuffered**: Blocks until a sender sends, then receives value and unblocks
2. **Buffered (not empty)**: Removes value from buffer, returns immediately
3. **Buffered (empty)**: Blocks until value available

**Behavior on closed channel:**
- Returns the **zero value** for the channel's type
- Never blocks (even if empty)

**Detecting closure:**
```go
value, ok := <-ch
if !ok {
    // Channel is closed and empty
}
```

**Panics if:**
- Receiving from a **nil** channel → blocks forever (deadlock)

**Example:**
```go
ch := make(chan int, 2)
ch <- 1
ch <- 2
close(ch)

v1 := <-ch     // 1 (from buffer)
v2 := <-ch     // 2 (from buffer)
v3 := <-ch     // 0 (closed, returns zero value)
v4 := <-ch     // 0 (still returns zero value)

v5, ok := <-ch // v5=0, ok=false (detects closure)
```

### Nil Channels

**Important edge case:**
```go
var ch chan int  // nil channel (zero value)

ch <- 42   // Blocks forever (deadlock)
<-ch       // Blocks forever (deadlock)
close(ch)  // PANIC: close of nil channel
```

**Use case for nil channels:**
Disabling a case in a `select` statement (covered later).

---

## Closing Channels

### Why Close Channels?

Closing a channel signals **"no more values will be sent"**. This allows receivers to:
- Detect when to stop processing
- Range over channels (loop exits when closed)
- Avoid goroutine leaks (don't wait forever)

### How to Close

```go
ch := make(chan int)
close(ch)
```

**Rules:**
1. **Only the sender should close** (receiver doesn't know if sender is done)
2. **Close is optional** (GC collects unclosed channels)
3. **Close exactly once** (closing a closed channel panics)
4. **Don't send after close** (panics)

### Detecting Closure (Two-Value Receive)

```go
value, ok := <-ch
if ok {
    // Channel is open, value is valid
} else {
    // Channel is closed
}
```

### Common Pattern: Signal Channel

```go
done := make(chan struct{})  // Empty struct uses 0 bytes

go func() {
    doWork()
    close(done)  // Signal by closing (don't need to send a value)
}()

<-done  // Wait for signal
```

**Why `struct{}`?**
- **0 bytes** (no memory overhead)
- **Clear intent**: We only care about closure, not the value

### Closing Multiple Receivers (Broadcast)

```go
done := make(chan struct{})

// Launch many goroutines
for i := 0; i < 1000; i++ {
    go func() {
        <-done  // All block on same channel
        // Cleanup...
    }()
}

// Signal all at once
close(done)  // All receivers unblock immediately!
```

**This is powerful:** One close operation wakes up all waiting receivers.

---

## Range Over Channels

### The Pattern

```go
ch := make(chan int)

// Producer
go func() {
    for i := 0; i < 10; i++ {
        ch <- i
    }
    close(ch)  // MUST close, or range will block forever
}()

// Consumer
for value := range ch {
    fmt.Println(value)  // Prints 0..9
}
// Loop exits when ch is closed
```

### How It Works

**Equivalent code:**
```go
for {
    value, ok := <-ch
    if !ok {
        break  // Channel closed
    }
    fmt.Println(value)
}
```

**Rules:**
1. Range **loops until channel is closed**
2. If channel is **never closed**, range **blocks forever** (goroutine leak)
3. Range **only works on receive** (`for range` on send-only channel is compile error)

### Common Mistake: Forgetting to Close

```go
// BAD: Goroutine leak
ch := make(chan int)

go func() {
    for i := 0; i < 10; i++ {
        ch <- i
    }
    // FORGOT TO CLOSE!
}()

for value := range ch {  // Blocks forever after receiving 10 values
    fmt.Println(value)
}
```

**Fix:**
```go
go func() {
    for i := 0; i < 10; i++ {
        ch <- i
    }
    close(ch)  // Signal completion
}()
```

---

## Select Statement: Multiplexing Channels

### The Problem

How do you wait on **multiple channels** simultaneously?

```go
ch1 := make(chan int)
ch2 := make(chan string)

// Want to receive from whichever is ready first
```

### The Solution: `select`

```go
select {
case v := <-ch1:
    fmt.Println("Received int:", v)
case v := <-ch2:
    fmt.Println("Received string:", v)
}
```

**How it works:**
1. **Waits** on all cases simultaneously
2. **Executes** the first case that's ready
3. **Random choice** if multiple cases ready (fairness)
4. **Blocks** if no cases ready (unless there's a `default`)

### Select with Default (Non-Blocking)

```go
select {
case v := <-ch:
    fmt.Println("Received:", v)
default:
    fmt.Println("No value ready")  // Executes immediately if ch empty
}
```

**Use case:** Polling without blocking.

### Select with Timeout

```go
select {
case v := <-ch:
    fmt.Println("Received:", v)
case <-time.After(1 * time.Second):
    fmt.Println("Timeout!")
}
```

**How it works:**
- `time.After(d)` returns a channel that receives after duration `d`
- If `ch` receives first, execute that case
- If timeout expires first, execute timeout case

### Select for Cancellation

```go
ctx, cancel := context.WithCancel(context.Background())

for {
    select {
    case <-ctx.Done():
        return  // Exit when cancelled
    case work := <-workCh:
        process(work)
    }
}
```

### Select with Send and Receive

```go
select {
case ch1 <- value:
    fmt.Println("Sent to ch1")
case v := <-ch2:
    fmt.Println("Received from ch2:", v)
case <-time.After(1 * time.Second):
    fmt.Println("Timeout")
}
```

### Disabling a Case with Nil Channel

```go
var ch1 chan int  // nil
ch2 := make(chan int)

select {
case <-ch1:  // This case is NEVER ready (nil channel blocks forever)
    // Never executes
case v := <-ch2:
    fmt.Println(v)
}
```

**Use case:** Dynamically enable/disable select cases.

**Example:**
```go
ch1 := make(chan int)
var ch2 chan int  // nil (disabled)

// ...later, enable ch2
ch2 = make(chan int)

select {
case v := <-ch1:
    // ...
case v := <-ch2:  // Now enabled
    // ...
}
```

---

## Pipeline Patterns

### What Is a Pipeline?

A **pipeline** is a series of stages connected by channels, where:
1. Each stage receives values from upstream
2. Performs a transformation
3. Sends results downstream
4. Runs in its own goroutine

**Benefits:**
- **Parallelism**: Stages run concurrently
- **Modularity**: Each stage is independent
- **Backpressure**: Slow stages naturally slow down fast stages

### Simple Pipeline

```go
// Stage 1: Generate numbers
func generate(n int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for i := 1; i <= n; i++ {
            out <- i
        }
    }()
    return out
}

// Stage 2: Square numbers
func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for v := range in {
            out <- v * v
        }
    }()
    return out
}

// Stage 3: Sum numbers
func sum(in <-chan int) int {
    total := 0
    for v := range in {
        total += v
    }
    return total
}

// Build pipeline
nums := generate(10)        // 1, 2, 3, ..., 10
squares := square(nums)     // 1, 4, 9, ..., 100
total := sum(squares)       // 385
```

**Visual:**
```
generate  ──→  square  ──→  sum
  (1..10)      (1..100)     (385)
   ↓            ↓             ↓
[goroutine]  [goroutine]  [main]
```

### Fan-Out, Fan-In Pattern

**Fan-Out:** Distribute work to multiple workers
**Fan-In:** Collect results from multiple workers

```go
// Fan-out: Multiple workers process same channel
func worker(id int, jobs <-chan int, results chan<- int) {
    for job := range jobs {
        results <- process(id, job)
    }
}

jobs := make(chan int, 100)
results := make(chan int, 100)

// Launch 10 workers (fan-out)
for i := 0; i < 10; i++ {
    go worker(i, jobs, results)
}

// Send jobs
for j := 0; j < 100; j++ {
    jobs <- j
}
close(jobs)

// Collect results (fan-in)
for i := 0; i < 100; i++ {
    <-results
}
```

### Cancellation in Pipelines

**Problem:** How to stop a pipeline early?

**Solution:** Use `context` or a `done` channel.

```go
func generate(done <-chan struct{}, n int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for i := 1; i <= n; i++ {
            select {
            case out <- i:
                // Sent successfully
            case <-done:
                return  // Exit early
            }
        }
    }()
    return out
}

done := make(chan struct{})
nums := generate(done, 1000000)

// Process a few values
for i := 0; i < 10; i++ {
    fmt.Println(<-nums)
}

close(done)  // Cancel pipeline
```

---

## Common Patterns and Idioms

### 1. Worker Pool

```go
func workerPool(numWorkers int, jobs <-chan int, results chan<- int) {
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                results <- process(job)
            }
        }()
    }
    wg.Wait()
    close(results)
}
```

### 2. Timeout Pattern

```go
select {
case result := <-ch:
    return result, nil
case <-time.After(5 * time.Second):
    return nil, errors.New("timeout")
}
```

### 3. Quit Channel

```go
quit := make(chan struct{})

go func() {
    for {
        select {
        case <-quit:
            return
        default:
            doWork()
        }
    }
}()

// Later...
close(quit)  // Signal quit
```

### 4. Future/Promise Pattern

```go
type Future chan int

func asyncCompute() Future {
    future := make(chan int, 1)
    go func() {
        future <- expensiveComputation()
    }()
    return future
}

future := asyncCompute()
// Do other work...
result := <-future  // Wait for result when needed
```

### 5. Semaphore (Limit Concurrency)

```go
sem := make(chan struct{}, 10)  // Max 10 concurrent

for i := 0; i < 100; i++ {
    sem <- struct{}{}  // Acquire
    go func() {
        defer func() { <-sem }()  // Release
        doWork()
    }()
}
```

---

## Common Pitfalls

### Pitfall 1: Send on Closed Channel

```go
ch := make(chan int)
close(ch)
ch <- 42  // PANIC: send on closed channel
```

**Rule:** Only send to open channels.

### Pitfall 2: Close from Receiver

```go
// BAD: Receiver closes
go func() {
    for v := range ch {
        process(v)
    }
    close(ch)  // WRONG! Sender might still be sending
}()
```

**Rule:** Only the sender should close.

### Pitfall 3: Forgetting to Close (Goroutine Leak)

```go
ch := make(chan int)

go func() {
    for i := 0; i < 10; i++ {
        ch <- i
    }
    // FORGOT close(ch)
}()

for v := range ch {  // Blocks forever after 10 values
    fmt.Println(v)
}
```

**Fix:** Always close when done sending.

### Pitfall 4: Buffered Channel as Unbuffered

```go
ch := make(chan int, 1)  // Buffered!

go func() {
    ch <- 42  // Doesn't block (buffer has space)
}()

// If you don't receive, goroutine exits but value is lost
// No synchronization guarantee
```

**Rule:** Use unbuffered for synchronization, buffered for queuing.

### Pitfall 5: Select with No Default (Accidental Block)

```go
select {
case v := <-ch1:
    // ...
case v := <-ch2:
    // ...
}
// Blocks if both channels are empty
```

**Fix:** Add `default` if you want non-blocking behavior.

---

## How to Run

```bash
# Run the demo program
cd minis/19-channels-basics
go run cmd/channels-demo/main.go

# Run exercises
cd exercise
go test -v

# Run a specific test
go test -v -run TestUnbufferedChannel
```

---

## Expected Output (Demo Program)

```
=== Channel Basics Demonstration ===

=== Unbuffered Channels ===
Sender: Sending value...
Receiver: Waiting for value...
Receiver: Received 42
Sender: Send complete

=== Buffered Channels ===
Sent 1 (buffer: [1])
Sent 2 (buffer: [1, 2])
Sent 3 (buffer: [1, 2, 3])
Sent 4 would block, launching receiver...
Received 1
Sent 4 (buffer: [2, 3, 4])

=== Select Statement ===
Received from ch1: hello
Received from ch2: 42
Timeout (no data)

=== Pipeline Pattern ===
Generated: 1, 2, 3, 4, 5
Squared:   1, 4, 9, 16, 25
Sum:       55
```

---

## Key Takeaways

1. **Channels are typed, thread-safe queues** for goroutine communication
2. **Unbuffered channels synchronize** (send blocks until receive, vice versa)
3. **Buffered channels decouple** (send blocks only when full, receive when empty)
4. **Close signals completion** (only sender closes, closing is optional)
5. **Range loops until closed** (must close or goroutine leaks)
6. **Select multiplexes channels** (waits on multiple, executes first ready)
7. **Pipelines compose stages** (each stage transforms and forwards)
8. **Direction matters** (`chan<-` send-only, `<-chan` receive-only)

---

## Connections to Other Projects

- **Project 18 (goroutines-1M-demo)**: Goroutines need channels for communication
- **Project 20 (select-fanin-fanout)**: Advanced select and pipeline patterns
- **Project 06 (worker-pool-wordcount)**: Worker pool using channels
- **Project 16 (context-cancellation-timeouts)**: Context uses channels internally
- **Project 24 (sync-mutex-vs-rwmutex)**: Alternative to channels (shared memory)

---

## Stretch Goals

1. **Rate limiter**: Implement a token bucket rate limiter using channels
2. **Ordered fan-in**: Merge multiple channels while preserving order
3. **Cancellable pipeline**: Build a 3-stage pipeline with context cancellation
4. **Benchmark**: Compare channel overhead vs mutex for different access patterns
5. **Visualize blocking**: Add timing logs to show when sends/receives block
