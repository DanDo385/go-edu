# Project 20: Select, Fan-In, and Fan-Out Patterns

## What Is This Project About?

This project teaches you **Go's `select` statement** and advanced **channel coordination patterns** that are the foundation of concurrent Go programs. You'll learn:

1. **Select statement** (multiplexing multiple channel operations)
2. **Fan-in patterns** (merging multiple channels into one)
3. **Fan-out patterns** (distributing work across multiple workers)
4. **Non-blocking operations** (try without blocking)
5. **Timeout patterns** (operations with time limits)
6. **Default cases** (fallback when no channel is ready)

By the end, you'll understand how to orchestrate complex concurrent workflows using channels and select, powering everything from web servers to distributed systems.

---

## The Fundamental Problem: Coordinating Multiple Channels

### First Principles: What Is the Select Statement?

Imagine you have multiple channels and you need to:
- Read from whichever channel has data first
- Send to whichever channel is ready
- Implement timeouts for operations
- Avoid blocking when channels aren't ready

**Without select**, you'd have to:
```go
// BAD: Can only wait on one channel at a time
result1 := <-ch1  // Blocks forever if ch1 never sends
result2 := <-ch2  // Never reached if ch1 blocks
```

**With select**, you can:
```go
// GOOD: Wait on multiple channels simultaneously
select {
case result := <-ch1:
    fmt.Println("Received from ch1:", result)
case result := <-ch2:
    fmt.Println("Received from ch2:", result)
case <-time.After(1 * time.Second):
    fmt.Println("Timeout!")
}
```

**The select statement is like a switch statement for channels**:
- It waits on multiple channel operations simultaneously
- It executes the case that becomes ready first
- If multiple cases are ready, it picks one **at random**
- If no cases are ready, it blocks (unless there's a `default` case)

### The Core Syntax

```go
select {
case v := <-ch1:
    // Receive from ch1
    fmt.Println("Received:", v)

case ch2 <- value:
    // Send to ch2
    fmt.Println("Sent to ch2")

case <-time.After(1 * time.Second):
    // Timeout after 1 second
    fmt.Println("Timeout")

default:
    // No channel operation ready (non-blocking)
    fmt.Println("Nothing ready")
}
```

**Key behaviors**:
1. **Blocks until one case is ready** (unless there's a `default`)
2. **Random selection** if multiple cases are ready (prevents starvation)
3. **Evaluates all channel expressions once** before selecting
4. **Executes only one case** per select statement

---

## Problem 1: The Select Statement Fundamentals

### What Problem Does Select Solve?

**Scenario**: You're building a web server that needs to handle:
- Incoming HTTP requests (from one channel)
- Shutdown signals (from another channel)
- Health checks (from a third channel)

Without select, you'd have to:
1. Read from requests channel (blocks if no requests)
2. Check shutdown channel (but you're already blocked!)
3. Never reach health check channel

**With select**, you handle all three simultaneously:

```go
for {
    select {
    case req := <-requests:
        handleRequest(req)
    case <-shutdown:
        return
    case <-healthCheck:
        respondHealthy()
    }
}
```

### Select Behavior Deep Dive

**1. Blocking Behavior**

```go
ch1 := make(chan int)
ch2 := make(chan int)

// This blocks until one channel receives a value
select {
case v := <-ch1:
    fmt.Println("ch1:", v)
case v := <-ch2:
    fmt.Println("ch2:", v)
}
// Program blocks here forever if no goroutines send to ch1 or ch2
```

**2. Random Selection (Fairness)**

```go
ch1 := make(chan int, 1)
ch2 := make(chan int, 1)

// Both channels are buffered and have values
ch1 <- 1
ch2 <- 2

// This select will randomly choose between ch1 and ch2
// Run multiple times, you'll see different outputs
select {
case v := <-ch1:
    fmt.Println("ch1:", v)  // Sometimes this
case v := <-ch2:
    fmt.Println("ch2:", v)  // Sometimes this
}
```

**Why random?** To prevent starvation. If select always chose the first ready case, you could have unfairness where one channel is always starved.

**3. Non-Blocking with Default**

```go
ch := make(chan int)

// This never blocks
select {
case v := <-ch:
    fmt.Println("Received:", v)
default:
    fmt.Println("No value available")
    // Executes immediately if ch has no value
}
```

**Use case**: Polling channels without blocking the goroutine.

**4. Send Operations**

```go
ch := make(chan int)

// Can select on send operations too
select {
case ch <- 42:
    fmt.Println("Sent 42")
default:
    fmt.Println("Channel not ready to receive")
}
```

**5. Nil Channel Behavior**

```go
var ch chan int  // nil channel

// This case is NEVER selected (nil channels block forever)
select {
case v := <-ch:
    fmt.Println("Never executes")
default:
    fmt.Println("Always executes (ch is nil)")
}
```

**Use case**: Dynamically enable/disable cases by setting channels to nil.

---

## Problem 2: Non-Blocking Operations

### The Pattern: Try Without Blocking

Sometimes you want to **try** an operation without blocking:
- Try to send without blocking (skip if channel is full)
- Try to receive without blocking (skip if channel is empty)

**Pattern: Non-blocking send**

```go
ch := make(chan int, 1)
ch <- 1  // Fill the buffer

value := 2
select {
case ch <- value:
    fmt.Println("Sent:", value)
default:
    fmt.Println("Channel full, skipping send")
}
// Output: "Channel full, skipping send"
```

**Pattern: Non-blocking receive**

```go
ch := make(chan int)

select {
case v := <-ch:
    fmt.Println("Received:", v)
default:
    fmt.Println("No value available")
}
// Output: "No value available"
```

**Pattern: Check if channel is closed (non-blocking)**

```go
ch := make(chan int)
close(ch)

select {
case v, ok := <-ch:
    if !ok {
        fmt.Println("Channel closed")
    } else {
        fmt.Println("Received:", v)
    }
default:
    fmt.Println("Nothing available")
}
// Output: "Channel closed" (closed channels are always ready)
```

### Real-World Example: Rate Limiting

```go
type RateLimiter struct {
    tokens chan struct{}
}

func (rl *RateLimiter) TryAcquire() bool {
    select {
    case <-rl.tokens:
        return true  // Got token
    default:
        return false // No tokens available
    }
}

// Usage:
if limiter.TryAcquire() {
    processRequest()
} else {
    rejectRequest("rate limit exceeded")
}
```

---

## Problem 3: Timeout Patterns

### The Pattern: Operations with Time Limits

**Pattern 1: Timeout for a single operation**

```go
ch := make(chan int)

select {
case v := <-ch:
    fmt.Println("Received:", v)
case <-time.After(1 * time.Second):
    fmt.Println("Timeout after 1 second")
}
```

**How it works**:
- `time.After(duration)` returns a channel
- That channel receives a value after the duration elapses
- Select picks whichever case happens first

**Pattern 2: Timeout with context**

```go
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()

select {
case v := <-ch:
    fmt.Println("Received:", v)
case <-ctx.Done():
    fmt.Println("Timeout:", ctx.Err())
}
```

**Pattern 3: Periodic timeout (ticker)**

```go
ticker := time.NewTicker(500 * time.Millisecond)
defer ticker.Stop()

for {
    select {
    case v := <-ch:
        fmt.Println("Received:", v)
    case <-ticker.C:
        fmt.Println("Heartbeat (no data for 500ms)")
    }
}
```

### Real-World Example: Database Query with Timeout

```go
func QueryWithTimeout(query string, timeout time.Duration) (Result, error) {
    resultCh := make(chan Result, 1)
    errCh := make(chan error, 1)

    go func() {
        result, err := database.Query(query)
        if err != nil {
            errCh <- err
        } else {
            resultCh <- result
        }
    }()

    select {
    case result := <-resultCh:
        return result, nil
    case err := <-errCh:
        return Result{}, err
    case <-time.After(timeout):
        return Result{}, fmt.Errorf("query timeout after %v", timeout)
    }
}
```

---

## Problem 4: Fan-In Pattern (Merging Multiple Channels)

### What Is Fan-In?

**Fan-in** is the pattern of merging multiple input channels into a single output channel.

**Real-world analogy**: Multiple checkout lines at a grocery store merging into a single exit.

**Use cases**:
- Collecting results from multiple workers
- Aggregating logs from multiple services
- Merging search results from multiple sources

### The Pattern

```go
func fanIn(ch1, ch2 <-chan int) <-chan int {
    out := make(chan int)

    go func() {
        defer close(out)
        for {
            select {
            case v, ok := <-ch1:
                if !ok {
                    ch1 = nil  // Disable this case
                    continue
                }
                out <- v
            case v, ok := <-ch2:
                if !ok {
                    ch2 = nil  // Disable this case
                    continue
                }
                out <- v
            }

            // Exit when both channels are closed
            if ch1 == nil && ch2 == nil {
                return
            }
        }
    }()

    return out
}
```

**How it works**:
1. Create output channel
2. Start goroutine that reads from all input channels
3. Use select to read from whichever channel has data
4. Forward data to output channel
5. When input channel closes, set it to nil (disables that case)
6. Exit when all input channels are closed

### Variation: Fan-In with Reflection (N channels)

```go
func fanInMany(channels ...<-chan int) <-chan int {
    out := make(chan int)

    var wg sync.WaitGroup

    // Start goroutine for each input channel
    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for v := range c {
                out <- v
            }
        }(ch)
    }

    // Close output when all inputs are exhausted
    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

### Real-World Example: Aggregating Logs

```go
func AggregateLogs(sources ...LogSource) <-chan LogEntry {
    out := make(chan LogEntry)

    var wg sync.WaitGroup
    for _, source := range sources {
        wg.Add(1)
        go func(s LogSource) {
            defer wg.Done()
            for entry := range s.Stream() {
                out <- entry
            }
        }(source)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}

// Usage:
logs := AggregateLogs(serverLogs, appLogs, dbLogs)
for entry := range logs {
    processLog(entry)
}
```

---

## Problem 5: Fan-Out Pattern (Distributing Work)

### What Is Fan-Out?

**Fan-out** is the pattern of distributing work from a single input channel across multiple workers.

**Real-world analogy**: A single task queue distributed across multiple workers.

**Use cases**:
- Parallel processing of jobs
- Load balancing across workers
- Distributing HTTP requests across handlers

### The Pattern

```go
func fanOut(input <-chan Task, numWorkers int) []<-chan Result {
    outputs := make([]<-chan Result, numWorkers)

    for i := 0; i < numWorkers; i++ {
        out := make(chan Result)
        outputs[i] = out

        go func() {
            defer close(out)
            for task := range input {
                result := process(task)
                out <- result
            }
        }()
    }

    return outputs
}
```

**How it works**:
1. Create multiple output channels (one per worker)
2. Start goroutine for each worker
3. Each worker reads from shared input channel
4. Each worker sends results to its own output channel
5. Go's runtime ensures fair distribution (no worker starves)

### Better Pattern: Fan-Out + Fan-In (Worker Pool)

```go
func WorkerPool(input <-chan Task, numWorkers int) <-chan Result {
    // Fan-out: distribute tasks to workers
    workerOutputs := make([]<-chan Result, numWorkers)

    for i := 0; i < numWorkers; i++ {
        out := make(chan Result)
        workerOutputs[i] = out

        go func(workerID int) {
            defer close(out)
            for task := range input {
                result := processTask(task, workerID)
                out <- result
            }
        }(i)
    }

    // Fan-in: merge worker outputs into single channel
    return fanInMany(workerOutputs...)
}
```

### Real-World Example: Image Processing Pipeline

```go
func ProcessImages(images <-chan Image, numWorkers int) <-chan ProcessedImage {
    results := make(chan ProcessedImage)

    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for img := range images {
                // Process image (expensive operation)
                processed := resize(img)
                processed = applyFilters(processed)
                processed = compress(processed)

                results <- processed
            }
        }(i)
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    return results
}

// Usage:
images := loadImages()
processed := ProcessImages(images, 8)  // 8 workers

for img := range processed {
    saveImage(img)
}
```

---

## Problem 6: Combining Patterns (Pipeline Architecture)

### The Pattern: Multi-Stage Pipeline

A **pipeline** chains together fan-out, processing, and fan-in stages:

```
Input → Stage 1 (fan-out) → Stage 2 (fan-out) → Stage 3 (fan-in) → Output
        [workers]            [workers]            [merge]
```

**Example: Data Processing Pipeline**

```go
// Stage 1: Generate data
func generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

// Stage 2: Square numbers (fan-out)
func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n
        }
    }()
    return out
}

// Stage 3: Merge results (fan-in)
func merge(channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup

    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for v := range c {
                out <- v
            }
        }(ch)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}

// Build pipeline
func main() {
    // Generate numbers
    nums := generate(1, 2, 3, 4, 5)

    // Fan-out to 3 workers
    sq1 := square(nums)
    sq2 := square(nums)
    sq3 := square(nums)

    // Fan-in results
    results := merge(sq1, sq2, sq3)

    // Consume results
    for v := range results {
        fmt.Println(v)
    }
}
```

---

## Advanced Patterns

### Pattern 1: Quit Channel

Signal goroutines to stop:

```go
func worker(quit <-chan struct{}) {
    for {
        select {
        case <-quit:
            fmt.Println("Worker stopping")
            return
        default:
            doWork()
        }
    }
}

// Usage:
quit := make(chan struct{})
go worker(quit)

time.Sleep(5 * time.Second)
close(quit)  // Signal all workers to stop
```

### Pattern 2: Done Channel

Wait for goroutine completion:

```go
func worker(done chan<- struct{}) {
    defer func() { done <- struct{}{} }()
    doWork()
}

// Usage:
done := make(chan struct{})
go worker(done)
<-done  // Wait for completion
```

### Pattern 3: Select with Multiple Timeouts

Different timeouts for different operations:

```go
select {
case v := <-fastCh:
    // Process fast result
case v := <-slowCh:
    // Process slow result
case <-time.After(100 * time.Millisecond):
    // Fast operation timeout
case <-time.After(5 * time.Second):
    // Slow operation timeout
}
```

**Problem**: Both timeout channels are created, wasting resources.

**Better approach**:

```go
fastTimeout := time.After(100 * time.Millisecond)
slowTimeout := time.After(5 * time.Second)

select {
case v := <-fastCh:
    // ...
case v := <-slowCh:
    // ...
case <-fastTimeout:
    // ...
case <-slowTimeout:
    // ...
}
```

### Pattern 4: Priority Select (Preference for Certain Channels)

```go
// Give priority to quit channel
for {
    select {
    case <-quit:
        return
    default:
    }

    select {
    case v := <-ch:
        process(v)
    case <-quit:
        return
    }
}
```

The first select checks quit with default (non-blocking). If quit is ready, we return immediately. Otherwise, we proceed to the second select.

---

## Common Mistakes to Avoid

### Mistake 1: time.After in Tight Loop

```go
// BAD: Creates new timer every iteration (memory leak!)
for {
    select {
    case v := <-ch:
        process(v)
    case <-time.After(1 * time.Second):
        timeout()
    }
}
```

**Problem**: Each `time.After` creates a timer that lives until it fires. In a tight loop, you create thousands of timers.

**Fix**: Use `time.NewTimer` and reuse it:

```go
// GOOD: Reuse timer
timer := time.NewTimer(1 * time.Second)
defer timer.Stop()

for {
    timer.Reset(1 * time.Second)
    select {
    case v := <-ch:
        process(v)
    case <-timer.C:
        timeout()
    }
}
```

### Mistake 2: Forgetting to Close Channels

```go
// BAD: Never closes channels
func generate() <-chan int {
    out := make(chan int)
    go func() {
        for i := 0; i < 10; i++ {
            out <- i
        }
        // Forgot to close!
    }()
    return out
}

// Receiver hangs forever after 10 items
for v := range generate() {
    fmt.Println(v)  // Hangs after 10
}
```

**Fix**: Always close channels when done:

```go
// GOOD: Close when done
func generate() <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)  // ← Critical
        for i := 0; i < 10; i++ {
            out <- i
        }
    }()
    return out
}
```

### Mistake 3: Select on Nil Channel

```go
var ch chan int  // nil

select {
case <-ch:  // Never selected (blocks forever)
    fmt.Println("Never prints")
}
```

**Use case**: This is actually useful for disabling cases:

```go
var ch1, ch2 <-chan int

if disableCh1 {
    ch1 = nil  // This case will be ignored
}

select {
case v := <-ch1:  // Won't be selected if ch1 is nil
    process1(v)
case v := <-ch2:
    process2(v)
}
```

### Mistake 4: Race Condition with Shared State

```go
// BAD: Race condition
count := 0

for i := 0; i < 10; i++ {
    go func() {
        count++  // RACE!
    }()
}
```

**Fix**: Use channels or mutexes:

```go
// GOOD: Use channel
countCh := make(chan int)

for i := 0; i < 10; i++ {
    go func() {
        countCh <- 1
    }()
}

count := 0
for i := 0; i < 10; i++ {
    count += <-countCh
}
```

---

## Performance Characteristics

### Select Statement

- **Time complexity**: O(n) where n = number of cases
  - Runtime must evaluate all channel expressions
  - Random selection among ready cases is O(n)

- **Space complexity**: O(n)
  - Must track all cases during evaluation

- **Fairness**: Random selection prevents starvation
  - If multiple cases are ready, random choice ensures fairness over time

### Fan-In/Fan-Out

- **Throughput**: Increases linearly with workers (up to CPU cores)
  - 1 worker: N tasks/sec
  - 4 workers: ~4N tasks/sec (if CPU-bound)

- **Latency**: Decreases with fan-out
  - More workers → shorter queue per worker → lower wait time

- **Overhead**: Each goroutine costs ~2-4 KB
  - 100 workers = ~400 KB (negligible)
  - 100,000 workers = ~400 MB (significant)

---

## Real-World Applications

### Web Server Request Handling

```go
func HandleRequests(requests <-chan Request, shutdown <-chan struct{}) {
    for {
        select {
        case req := <-requests:
            go handleRequest(req)
        case <-shutdown:
            return
        }
    }
}
```

### Load Balancer

```go
func LoadBalance(requests <-chan Request, backends []Backend) {
    for req := range requests {
        // Fan-out: send to least loaded backend
        backend := selectLeastLoaded(backends)
        backend.Handle(req)
    }
}
```

### Data Pipeline

```go
func Pipeline(input <-chan Data) <-chan Result {
    // Stage 1: Parse (fan-out)
    parsed := fanOut(input, parseWorker, 4)

    // Stage 2: Validate (fan-out)
    validated := fanOut(parsed, validateWorker, 4)

    // Stage 3: Transform (fan-out)
    transformed := fanOut(validated, transformWorker, 4)

    // Stage 4: Aggregate (fan-in)
    return fanIn(transformed)
}
```

---

## How to Run

```bash
# Run the demo
go run ./minis/20-select-fanin-fanout/cmd/select-demo/main.go

# Run tests
go test ./minis/20-select-fanin-fanout/...

# Run with verbose output
go test -v ./minis/20-select-fanin-fanout/...

# Run with race detector
go test -race ./minis/20-select-fanin-fanout/...
```

---

## Summary

**What you learned**:
- ✅ Select statement multiplexes multiple channel operations
- ✅ Random selection ensures fairness (prevents starvation)
- ✅ Default case enables non-blocking operations
- ✅ Timeout patterns prevent operations from hanging
- ✅ Fan-in merges multiple channels into one
- ✅ Fan-out distributes work across workers
- ✅ Pipelines chain stages for complex workflows
- ✅ Nil channels disable select cases (dynamic control)

**Key rules**:
1. Select picks one ready case at random (fairness)
2. Default makes select non-blocking
3. Nil channels are never ready (use to disable cases)
4. Close channels to signal completion
5. Always reuse timers in loops (avoid `time.After`)
6. Fan-in with WaitGroup ensures proper cleanup

**Next steps**:
- Project 21: Sync primitives (Mutex, RWMutex, Cond)
- Project 22: Atomic operations
- Project 24: Advanced concurrency patterns

Master select, master Go concurrency! 🚀
