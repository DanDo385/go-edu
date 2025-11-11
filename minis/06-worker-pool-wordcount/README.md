# Project 06: Worker Pool Wordcount

## 1. What Is This About?

### Real-World Scenario

Imagine you're building a search engine crawler that needs to analyze thousands of web pages simultaneously. You could:

**❌ Bad approach:** Create one thread per webpage
- With 10,000 pages, you'd create 10,000 threads
- Each thread uses ~2MB of memory = 20GB just for thread stacks!
- Your system crashes or becomes unresponsive

**✅ Better approach:** Use a **worker pool** with bounded parallelism
- Create exactly 10 worker threads
- Feed them 10,000 URLs through a queue
- Workers process URLs concurrently but with controlled resource usage
- Memory usage stays constant regardless of URL count

This project teaches you how to build concurrent systems that are **fast** (parallel processing) yet **safe** (bounded resources).

### What You'll Learn

1. **Goroutines**: Lightweight threads in Go (only ~2KB per goroutine)
2. **Channels**: Type-safe queues for communication between goroutines
3. **Worker pools**: Bounded concurrency pattern (prevent resource exhaustion)
4. **Context**: Cancellation propagation across goroutines
5. **Error handling**: How to handle errors in concurrent code
6. **sync.WaitGroup**: Coordinating goroutine completion

### The Challenge

Fetch multiple URLs concurrently, count word frequencies across all pages, and aggregate the results - all while:
- Limiting concurrent HTTP requests (don't overwhelm servers)
- Cancelling all work if any request fails
- Safely merging results from multiple goroutines
- Handling errors gracefully

---

## 2. First Principles: Understanding Concurrency in Go

### What is a Goroutine?

A **goroutine** is a lightweight thread managed by the Go runtime.

**Analogy**: Imagine a restaurant kitchen:
- **OS Threads** = Full-time chefs (expensive, limited number, each has their own workspace)
- **Goroutines** = Tasks (chopping, stirring, plating) that chefs can work on
- **Go Scheduler** = Kitchen manager who assigns tasks to available chefs

```go
// Creating a goroutine is trivial
go doSomething()  // Runs concurrently with the rest of your code
```

**Key differences from OS threads:**

| Feature | OS Thread | Goroutine |
|---------|-----------|-----------|
| Stack size | 2MB (fixed) | 2KB (grows dynamically) |
| Creation cost | High (~1-2ms) | Low (~few microseconds) |
| Manageable count | Hundreds | Millions |
| Scheduling | OS kernel | Go runtime |

### What is a Channel?

A **channel** is a typed queue that allows goroutines to communicate safely.

**Analogy**: Channels are like conveyor belts in a factory:
- One goroutine puts items on the belt (sender)
- Another goroutine takes items off the belt (receiver)
- The belt has a size limit (buffer)

```go
ch := make(chan string)     // Unbuffered channel (blocks sender until receiver reads)
ch := make(chan string, 10) // Buffered channel (sender can add 10 items before blocking)

// Sending
ch <- "hello"  // Put value into channel

// Receiving
msg := <-ch   // Get value from channel
```

**Critical insight**: Channels provide **synchronization** for free:
- Unbuffered channel: Sender waits for receiver (rendezvous)
- Buffered channel: Sender only blocks when buffer is full

### What is the Worker Pool Pattern?

The **worker pool** pattern limits concurrency by using a fixed number of goroutines ("workers") that process jobs from a shared queue.

**Visual representation**:

```
                         ┌─────────────┐
    URLs to fetch ──────►│  Jobs Queue │
                         │  (channel)  │
                         └──────┬──────┘
                                │
                ┌───────────────┼───────────────┐
                │               │               │
                ▼               ▼               ▼
          ┌─────────┐     ┌─────────┐     ┌─────────┐
          │Worker 1 │     │Worker 2 │     │Worker 3 │
          │(goroutine)     │(goroutine)     │(goroutine)
          └────┬────┘     └────┬────┘     └────┬────┘
               │               │               │
               │   Fetch URL, count words      │
               │               │               │
               └───────────────┼───────────────┘
                               │
                               ▼
                         ┌─────────────┐
                         │   Results   │
                         │  (channel)  │
                         └─────────────┘
                               │
                               ▼
                         Aggregate counts
```

**Why this pattern?**

1. **Bounded resources**: Exactly N workers, not unlimited
2. **Backpressure**: If jobs arrive faster than workers can process, they queue up
3. **Clean shutdown**: Close the jobs channel when no more work
4. **Error handling**: Workers can signal errors to a central handler

### What is Context?

`context.Context` is Go's standard way to handle cancellation, timeouts, and request-scoped values.

**Analogy**: Think of context like a "stop button" that you can press to cancel all related work:
- You start a task with a context
- If something goes wrong, you cancel the context
- All subtasks that are listening to that context stop immediately

```go
ctx, cancel := context.WithCancel(context.Background())

// Later, when you want to stop everything:
cancel()

// Workers can check if context is cancelled:
select {
case <-ctx.Done():
    // Context cancelled, stop working
    return
}
```

---

## 3. Breaking Down the Solution

### Step 1: Understand the Problem

**Inputs**:
- `[]string` (list of URLs to fetch)
- `int` (number of workers)

**Outputs**:
- `map[string]int` (word frequencies across all URLs)
- `error` (if any fetch failed)

**Requirements**:
- Fetch URLs concurrently (for speed)
- But limit concurrency to N workers (for safety)
- Tokenize response bodies (extract words)
- Merge word counts from all pages
- Cancel all work if any request fails

### Step 2: Design the Architecture

We need **four types of goroutines**:

1. **Main goroutine**: Sends URLs to workers, waits for completion
2. **Worker goroutines** (N of them): Fetch URLs, count words, send results
3. **Aggregator goroutine**: Merges word counts from all workers
4. **Error handler**: Cancels context on first error

**Communication channels**:
- `jobs chan string`: Main → Workers (URLs to fetch)
- `results chan map[string]int`: Workers → Aggregator (word counts)
- `errCh chan error`: Workers → Main (error notifications)

### Step 3: Worker Lifecycle

Each worker follows this pattern:

```
1. Wait for a job from the jobs channel
2. If channel is closed → exit (no more work)
3. If context is cancelled → exit (error occurred elsewhere)
4. Fetch the URL
5. If error → send to errCh, cancel context, exit
6. Tokenize response → count words
7. Send word counts to results channel
8. Go back to step 1
```

### Step 4: Synchronization Strategy

**Problem**: How do we know when all workers are done?

**Solution**: Use `sync.WaitGroup`:
```go
var wg sync.WaitGroup
wg.Add(numWorkers)  // Increment counter

// Each worker:
defer wg.Done()  // Decrement counter when done

// Main goroutine:
wg.Wait()  // Block until counter reaches 0
```

**Problem**: When do we close the results channel?

**Solution**: Close it after all workers are done:
```go
go func() {
    wg.Wait()         // Wait for all workers
    close(results)    // Then close results channel
}()
```

This allows the aggregator to use `range` over the results channel:
```go
for counts := range results {
    // Merge counts
}
// Loop exits when results channel is closed
```

### Step 5: Error Propagation

**Challenge**: If one URL fails, we want to cancel all in-flight requests immediately.

**Solution**:
1. Create a cancellable context: `ctx, cancel := context.WithCancel(ctx)`
2. Pass this context to all HTTP requests
3. When a worker encounters an error:
   - Send error to `errCh`
   - Call `cancel()` to cancel the context
   - All other workers check `ctx.Done()` and exit

---

## 4. Complete Solution Walkthrough

Let's walk through the code step by step.

### Function Signature

```go
func WordCount(ctx context.Context, urls []string, workers int) (map[string]int, error)
```

**Why `context.Context` as first parameter?**
Go convention: Context is always the first parameter. This allows callers to set timeouts:
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
counts, err := WordCount(ctx, urls, 5)  // Entire operation must finish in 30s
```

### Creating Channels

```go
jobs := make(chan string, workers)
results := make(chan map[string]int, workers)
errCh := make(chan error, 1)
```

**Why buffer `jobs` and `results` with size `workers`?**
- Prevents deadlock: If all workers are busy, sender can still add jobs without blocking
- Performance: Reduces synchronization overhead

**Why buffer `errCh` with size 1?**
- Non-blocking send: First error is guaranteed to be sent even if no one is reading yet
- We only care about the first error (subsequent errors are ignored)

### Creating Cancellable Context

```go
ctx, cancel := context.WithCancel(ctx)
defer cancel()  // Clean up resources
```

**Why `defer cancel()`?**
Even if no error occurs, we must call `cancel()` to release resources. The `defer` ensures this happens no matter how the function exits.

### Starting Workers

```go
var wg sync.WaitGroup
for i := 0; i < workers; i++ {
    wg.Add(1)
    go func(workerID int) {
        defer wg.Done()
        // Worker logic
    }(i)
}
```

**Key pattern**:
- `wg.Add(1)` before `go`: Increment counter before starting goroutine
- `defer wg.Done()`: Decrement counter when goroutine exits (even if it panics)
- Pass `i` as parameter: Avoids closure variable capture issue

### Worker Logic

```go
for {
    select {
    case <-ctx.Done():
        return  // Context cancelled
    case url, ok := <-jobs:
        if !ok {
            return  // Jobs channel closed
        }

        counts, err := fetchAndCount(ctx, url)
        if err != nil {
            select {
            case errCh <- fmt.Errorf("fetching %s: %w", url, err):
                cancel()
            default:
                // Error channel already full, ignore
            }
            return
        }

        select {
        case <-ctx.Done():
            return
        case results <- counts:
            // Sent successfully
        }
    }
}
```

**Why two `select` statements?**

First `select`: Check context **or** receive job
- `case <-ctx.Done()`: Exit if context cancelled
- `case url, ok := <-jobs`: Receive next job (if available)

Second `select`: Send error without blocking
- `case errCh <- err`: Send error
- `default`: If channel is full (already has an error), ignore this error

Third `select`: Send results (interruptible)
- `case results <- counts`: Send results
- `case <-ctx.Done()`: Exit if context cancelled while sending

**What is `url, ok := <-jobs`?**
- `ok` is `false` if channel is closed and empty
- This is how we detect "no more jobs"

### Sending Jobs

```go
go func() {
    for _, url := range urls {
        select {
        case <-ctx.Done():
            return
        case jobs <- url:
            // Sent successfully
        }
    }
    close(jobs)  // Signal no more jobs
}()
```

**Why in a goroutine?**
If we did this in the main goroutine, we'd block until all jobs are sent. By using a goroutine, workers can start processing jobs while the rest are still being sent.

**Why `close(jobs)`?**
Signals to workers that no more jobs will arrive. Workers exit when they receive from a closed channel.

### Closing Results Channel

```go
go func() {
    wg.Wait()
    close(results)
}()
```

**Critical timing**:
1. All workers finish (`wg.Wait()` returns)
2. Then close results channel
3. This unblocks the aggregator's `for range results` loop

### Aggregating Results

```go
finalCounts := make(map[string]int)
for counts := range results {
    for word, count := range counts {
        finalCounts[word] += count
    }
}
```

**Why is this safe without locks?**
Only one goroutine (the main one) is writing to `finalCounts`. Workers send results through a channel, which provides synchronization.

**What happens when `results` is closed?**
The `for range` loop exits, and we proceed to error checking.

### Error Checking

```go
select {
case err := <-errCh:
    return nil, err
default:
    // No error
}
```

**Why non-blocking receive?**
If there's an error, we return it. If not, we proceed (don't wait).

---

## 5. Key Concepts Explained

### Concept 1: Goroutines vs OS Threads

**Why are goroutines so lightweight?**

| Feature | OS Thread | Goroutine |
|---------|-----------|-----------|
| **Stack** | 2MB fixed | 2KB initial (grows to 1GB max) |
| **Scheduling** | Kernel (expensive context switch) | Runtime (cheap, cooperative) |
| **Creation** | syscall (~1ms) | Function call (~few μs) |

**Example**:
```go
// This is fine in Go:
for i := 0; i < 1000000; i++ {
    go func() {
        // Do work
    }()
}
// 1 million goroutines ≈ 2GB RAM

// This would crash in most languages:
// 1 million OS threads ≈ 2TB RAM (impossible)
```

**How does the Go scheduler work?**

The Go runtime uses an **M:N scheduler**:
- **M** goroutines are multiplexed onto **N** OS threads
- Typically N = number of CPU cores
- Go scheduler decides which goroutine runs on which thread
- Goroutines are **cooperative**: They yield when blocking (I/O, channel ops, etc.)

### Concept 2: Channel Semantics

**Unbuffered vs Buffered Channels**

```go
// Unbuffered: Synchronous communication
ch := make(chan int)
ch <- 42  // BLOCKS until someone receives

// Buffered: Asynchronous up to buffer size
ch := make(chan int, 3)
ch <- 1   // Doesn't block (buffer has space)
ch <- 2   // Doesn't block
ch <- 3   // Doesn't block
ch <- 4   // BLOCKS (buffer full)
```

**Channel states**:

| Operation | Nil Channel | Open Channel | Closed Channel |
|-----------|-------------|--------------|----------------|
| Send | Block forever | Send (or block) | **Panic** |
| Receive | Block forever | Receive (or block) | Receive zero value + `ok=false` |
| Close | **Panic** | Succeed | **Panic** |

**Key rule**: Only the sender should close a channel.

**Why?**
If a receiver closes the channel, senders will panic when they try to send.

### Concept 3: Context Cancellation

**Context is like a cascading shutdown signal:**

```go
parent, cancel1 := context.WithCancel(context.Background())
child, cancel2 := context.WithCancel(parent)

cancel1()  // Cancels both parent and child
```

**Using context with HTTP requests:**

```go
req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
resp, err := http.DefaultClient.Do(req)
```

If `ctx` is cancelled during the request, the HTTP client aborts the request immediately (doesn't wait for timeout).

### Concept 4: sync.WaitGroup

**WaitGroup is a counter for goroutines:**

```go
var wg sync.WaitGroup

// Increment
wg.Add(1)  // Counter: 0 → 1

// Decrement
wg.Done()  // Counter: 1 → 0

// Wait for counter to reach 0
wg.Wait()  // Blocks until counter == 0
```

**Common mistake**:
```go
// ❌ WRONG
for i := 0; i < 10; i++ {
    go func() {
        wg.Add(1)  // Race condition!
        defer wg.Done()
        // work
    }()
}
wg.Wait()  // Might return before all goroutines start

// ✅ CORRECT
for i := 0; i < 10; i++ {
    wg.Add(1)  // In main goroutine, before starting worker
    go func() {
        defer wg.Done()
        // work
    }()
}
wg.Wait()
```

### Concept 5: Select Statement

`select` is like `switch` for channels:

```go
select {
case msg := <-ch1:
    fmt.Println("Received from ch1:", msg)
case msg := <-ch2:
    fmt.Println("Received from ch2:", msg)
case ch3 <- "hello":
    fmt.Println("Sent to ch3")
default:
    fmt.Println("No channel ready")
}
```

**Key behaviors**:
- If multiple cases are ready, one is chosen **randomly**
- If no case is ready and there's no `default`, `select` **blocks**
- If there's a `default`, it executes immediately when no case is ready

**Common patterns**:

```go
// Pattern 1: Timeout
select {
case result := <-ch:
    return result
case <-time.After(5 * time.Second):
    return errors.New("timeout")
}

// Pattern 2: Non-blocking send
select {
case ch <- value:
    // Sent successfully
default:
    // Channel full, drop value or handle otherwise
}

// Pattern 3: Cancellation
for {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case work := <-workQueue:
        process(work)
    }
}
```

---

## 6. Common Patterns You Can Reuse

### Pattern 1: Worker Pool (This Project)

**When to use**: Bounded parallelism for I/O-bound tasks

```go
func WorkerPool(jobs <-chan Job, workers int) <-chan Result {
    results := make(chan Result, workers)
    var wg sync.WaitGroup

    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                results <- process(job)
            }
        }()
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    return results
}
```

### Pattern 2: Fan-Out, Fan-In

**When to use**: Distribute work, then collect results

```go
func FanOut(input <-chan int, workers int) []<-chan int {
    outputs := make([]<-chan int, workers)
    for i := 0; i < workers; i++ {
        outputs[i] = worker(input)
    }
    return outputs
}

func FanIn(inputs ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup

    for _, in := range inputs {
        wg.Add(1)
        go func(ch <-chan int) {
            defer wg.Done()
            for val := range ch {
                out <- val
            }
        }(in)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

### Pattern 3: Pipeline

**When to use**: Multi-stage processing

```go
func Pipeline(input <-chan int) <-chan int {
    stage1 := make(chan int)
    stage2 := make(chan int)

    // Stage 1: Double
    go func() {
        defer close(stage1)
        for val := range input {
            stage1 <- val * 2
        }
    }()

    // Stage 2: Square
    go func() {
        defer close(stage2)
        for val := range stage1 {
            stage2 <- val * val
        }
    }()

    return stage2
}
```

### Pattern 4: Timeout Per Operation

**When to use**: Each operation has its own deadline

```go
func FetchWithTimeout(url string, timeout time.Duration) ([]byte, error) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}
```

### Pattern 5: Semaphore (Bounded Concurrency)

**When to use**: Limit number of concurrent operations without worker pool

```go
type Semaphore chan struct{}

func NewSemaphore(max int) Semaphore {
    return make(Semaphore, max)
}

func (s Semaphore) Acquire() {
    s <- struct{}{}
}

func (s Semaphore) Release() {
    <-s
}

// Usage:
sem := NewSemaphore(10)
for _, url := range urls {
    sem.Acquire()
    go func(u string) {
        defer sem.Release()
        fetch(u)
    }(url)
}
```

---

## 7. Real-World Applications

### Web Crawlers

**Use case**: Crawl millions of pages without exhausting resources

```go
crawler := NewCrawler(maxWorkers: 100)
for seed := range seedURLs {
    crawler.Enqueue(seed)
}
results := crawler.Run()
```

Companies using this: Google, Bing, Archive.org

### Image Processing Pipelines

**Use case**: Process uploaded images (resize, watermark, upload to CDN)

```
Upload → Validation → Resize → Watermark → CDN Upload → Database
  |          |          |          |           |           |
Worker    Worker    Worker    Worker      Worker      Worker
 Pool      Pool      Pool      Pool        Pool        Pool
```

Companies using this: Instagram, Pinterest, Shutterstock

### Log Aggregation

**Use case**: Collect logs from thousands of servers concurrently

```go
type LogCollector struct {
    workers int
    servers []Server
}

func (lc *LogCollector) Collect() []LogEntry {
    // Worker pool pattern
}
```

Companies using this: Splunk, Datadog, Elasticsearch

### API Rate Limiter

**Use case**: Make many API calls while respecting rate limits

```go
type RateLimitedClient struct {
    workers   int
    rateLimit int  // requests per second
}

func (c *RateLimitedClient) Fetch(urls []string) []Response {
    // Worker pool + rate limiting
}
```

Companies using this: Any company integrating with third-party APIs

### Batch Processing

**Use case**: Process millions of database records concurrently

```go
func ProcessRecords(db *sql.DB, workers int) error {
    rows, _ := db.Query("SELECT * FROM records")
    defer rows.Close()

    jobs := make(chan Record, workers)
    // Worker pool processes records
}
```

Companies using this: Banks (transaction processing), E-commerce (order fulfillment)

---

## 8. Common Mistakes to Avoid

### Mistake 1: Unbounded Goroutines

**❌ Wrong**:
```go
for _, url := range urls {
    go fetch(url)  // If urls has 1 million items, you create 1M goroutines!
}
```

**✅ Correct**:
```go
jobs := make(chan string)
for i := 0; i < workers; i++ {
    go worker(jobs)
}
for _, url := range urls {
    jobs <- url
}
```

**Why**: Even though goroutines are cheap, 1 million goroutines is wasteful and can exhaust resources (file descriptors, memory, etc.).

### Mistake 2: Closing Channels Multiple Times

**❌ Wrong**:
```go
close(ch)
close(ch)  // PANIC: close of closed channel
```

**✅ Correct**:
```go
var once sync.Once
once.Do(func() { close(ch) })  // Safe to call multiple times
```

Or better: Design so only one goroutine is responsible for closing.

### Mistake 3: Sending on Closed Channel

**❌ Wrong**:
```go
close(ch)
ch <- value  // PANIC: send on closed channel
```

**✅ Correct**: Only the sender should close the channel, and it should stop sending before closing.

### Mistake 4: WaitGroup Add/Done Imbalance

**❌ Wrong**:
```go
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    go func() {
        wg.Add(1)  // Race: main might call Wait() before this runs
        defer wg.Done()
        work()
    }()
}
wg.Wait()  // Might return too early
```

**✅ Correct**:
```go
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)  // BEFORE starting goroutine
    go func() {
        defer wg.Done()
        work()
    }()
}
wg.Wait()
```

### Mistake 5: Forgetting Context Cancellation

**❌ Wrong**:
```go
ctx, cancel := context.WithCancel(context.Background())
// Forget to call cancel()
// Context resources leak
```

**✅ Correct**:
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()  // Always defer cancel()
```

### Mistake 6: Shared Memory Without Synchronization

**❌ Wrong**:
```go
var counter int
for i := 0; i < 100; i++ {
    go func() {
        counter++  // DATA RACE!
    }()
}
```

**✅ Correct**:
```go
var mu sync.Mutex
var counter int
for i := 0; i < 100; i++ {
    go func() {
        mu.Lock()
        counter++
        mu.Unlock()
    }()
}
```

Or even better: Use channels or `sync/atomic`.

### Mistake 7: Ignoring Context in Blocking Operations

**❌ Wrong**:
```go
func worker(ctx context.Context, jobs <-chan Job) {
    for job := range jobs {  // Doesn't respect context cancellation
        process(job)
    }
}
```

**✅ Correct**:
```go
func worker(ctx context.Context, jobs <-chan Job) {
    for {
        select {
        case <-ctx.Done():
            return
        case job := <-jobs:
            process(job)
        }
    }
}
```

---

## 9. Stretch Goals

### Goal 1: Add Progress Reporting ⭐

Track and report how many URLs have been fetched, failed, and are pending.

**Hint**: Add a stats struct and a channel for status updates:
```go
type Stats struct {
    Completed int
    Failed    int
    Pending   int
}

func WordCountWithProgress(ctx context.Context, urls []string, workers int, progress chan<- Stats) (map[string]int, error) {
    // Send stats updates periodically
}
```

### Goal 2: Implement Retry Logic ⭐⭐

If a URL fetch fails with a transient error (e.g., 503 Service Unavailable), retry it up to 3 times with exponential backoff.

**Hint**: Use a retry helper:
```go
func fetchWithRetry(ctx context.Context, url string, maxRetries int) (map[string]int, error) {
    for attempt := 0; attempt < maxRetries; attempt++ {
        result, err := fetchAndCount(ctx, url)
        if err == nil {
            return result, nil
        }
        if !isRetryable(err) {
            return nil, err
        }
        time.Sleep(time.Duration(1<<attempt) * time.Second)  // Exponential backoff
    }
    return nil, fmt.Errorf("max retries exceeded")
}
```

### Goal 3: Add Rate Limiting ⭐⭐

Limit the worker pool to N requests per second to avoid overwhelming servers.

**Hint**: Use `time.Ticker`:
```go
type RateLimitedPool struct {
    workers int
    rps     int  // requests per second
}

func (p *RateLimitedPool) WordCount(ctx context.Context, urls []string) (map[string]int, error) {
    ticker := time.NewTicker(time.Second / time.Duration(p.rps))
    defer ticker.Stop()

    for _, url := range urls {
        <-ticker.C  // Wait for rate limiter
        jobs <- url
    }
}
```

### Goal 4: Support Streaming Responses ⭐⭐⭐

Instead of buffering the entire response body, tokenize it in a streaming fashion.

**Hint**: Use `bufio.Scanner`:
```go
func fetchAndCountStreaming(ctx context.Context, url string) (map[string]int, error) {
    resp, err := fetch(ctx, url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    counts := make(map[string]int)
    scanner := bufio.NewScanner(resp.Body)
    scanner.Split(bufio.ScanWords)

    for scanner.Scan() {
        word := normalize(scanner.Text())
        counts[word]++
    }

    return counts, scanner.Err()
}
```

### Goal 5: Add Result Caching ⭐⭐⭐

If the same URL appears multiple times, fetch it only once and reuse the results.

**Hint**: Use a cache with sync.Map or regular map with mutex:
```go
type CachedPool struct {
    cache sync.Map  // url → word counts
    pool  *WorkerPool
}

func (p *CachedPool) WordCount(ctx context.Context, urls []string) (map[string]int, error) {
    uniqueURLs := deduplicate(urls)

    for _, url := range uniqueURLs {
        if cached, ok := p.cache.Load(url); ok {
            // Use cached result
            continue
        }

        // Fetch and cache
        counts, _ := p.pool.Fetch(ctx, url)
        p.cache.Store(url, counts)
    }
}
```

---

## How to Run

```bash
# Run the program
make run P=06-worker-pool-wordcount

# Run tests
go test ./minis/06-worker-pool-wordcount/...

# Run with race detector (detects data races)
go test -race ./minis/06-worker-pool-wordcount/...

# Benchmark
go test -bench=. ./minis/06-worker-pool-wordcount/...
```

---

## Summary

**What you learned**:
- ✅ Goroutines are lightweight threads (2KB vs 2MB)
- ✅ Channels provide type-safe communication between goroutines
- ✅ Worker pools limit concurrency for bounded resource usage
- ✅ Context enables cancellation propagation across goroutines
- ✅ sync.WaitGroup coordinates goroutine completion
- ✅ `select` statement multiplexes channel operations

**Why this matters**:
Go's concurrency model is one of its biggest strengths. The patterns you learned here (worker pools, fan-out/fan-in, pipelines) are used in production systems processing billions of requests per day.

**Next steps**:
- Project 07: Learn about generics and advanced data structures (LRU cache)
- Project 08: Apply concurrency to HTTP clients with retries
- Project 09: Build HTTP servers with graceful shutdown

Go forth and parallelize! 🚀
