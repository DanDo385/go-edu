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
          │(goroutine)    │(goroutine)    │(goroutine)
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

**First Principles: What is `ctx` Actually?**

When you see `ctx context.Context` as a function parameter, you're passing a **value** that contains:
1. **Cancellation channel**: A channel that gets closed when cancellation is requested
2. **Deadline**: Optional time when the operation should timeout
3. **Values**: Optional key-value pairs (like request IDs)

**What does "passing context" actually mean?**

```go
func myFunction(ctx context.Context) {
    // ctx is a VALUE (struct), not a pointer
    // But internally, it contains pointers to shared cancellation state
    // All goroutines using the same context share the same cancellation signal
}
```

**How Context Works Internally:**

```
┌─────────────────────────────────────────┐
│ context.Context (value type)             │
│ ┌─────────────────────────────────────┐   │
│ │ done chan struct{} (pointer)       │   │  ← Shared across all goroutines
│ │ err error                          │   │
│ │ deadline time.Time                 │   │
│ │ values map[interface{}]interface{} │   │
│ └─────────────────────────────────────┘   │
└─────────────────────────────────────────┘
```

When you call `cancel()`:
1. The `done` channel is **closed**
2. Reading from `<-ctx.Done()` **unblocks immediately**
3. All goroutines waiting on `ctx.Done()` wake up and can exit

**Analogy**: Think of context like a "stop button" that you can press to cancel all related work:
- You start a task with a context
- If something goes wrong, you cancel the context
- All subtasks that are listening to that context stop immediately

**Understanding `<-ctx.Done()`:**

```go
// ctx.Done() returns a channel
doneChannel := ctx.Done()

// Reading from it: <-ctx.Done()
select {
case <-ctx.Done():
    // This case executes when:
    // 1. cancel() was called (channel is closed)
    // 2. Deadline expired (channel is closed)
    // Reading from a closed channel returns immediately with zero value
    return
}
```

**Memory Management:**
- Context is a small struct (~48 bytes)
- Passed by value (copied), but internally contains pointers
- The `done` channel is shared across all goroutines using the context
- Very memory-efficient cancellation mechanism

```go
ctx, cancel := context.WithCancel(context.Background())

// Later, when you want to stop everything:
cancel()  // Closes the done channel

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

**Detailed Workflow Chart - Complete System:**

```
┌─────────────────────────────────────────────────────────────────┐
│ MAIN GOROUTINE: WordCount()                                     │
│                                                                 │
│ 1. Create channels: jobs, results, errCh                        │
│ 2. Create cancellable context: ctx, cancel                     │
│ 3. Create WaitGroup: wg                                         │
│                                                                 │
│ ┌───────────────────────────────────────────────────────────┐  │
│ │ WORKER SPAWN LOOP (Sequential)                            │  │
│ │ for i := 0; i < workers; i++ {                            │  │
│ │     wg.Add(1)                                             │  │
│ │     go func(workerID int) {  ← (i) passed here           │  │
│ │         defer wg.Done()                                   │  │
│ │         // Worker loop (concurrent)                       │  │
│ │     }(i)  ← LINE 163: Why (i) is needed                  │  │
│ │ }                                                          │  │
│ └───────────────────────────────────────────────────────────┘  │
│                                                                 │
│ ┌───────────────────────────────────────────────────────────┐  │
│ │ JOB SENDER GOROUTINE (Concurrent)                         │  │
│ │ go func() {                                               │  │
│ │     for _, url := range urls {                            │  │
│ │         jobs <- url  ← Send URL to channel                │  │
│ │     }                                                      │  │
│ │     close(jobs)  ← Signal no more jobs                    │  │
│ │ }()                                                        │  │
│ └───────────────────────────────────────────────────────────┘  │
│                                                                 │
│ ┌───────────────────────────────────────────────────────────┐  │
│ │ RESULTS CLOSER GOROUTINE (Concurrent)                     │  │
│ │ go func() {                                               │  │
│ │     wg.Wait()  ← Wait for all workers                    │  │
│ │     close(results)  ← Signal no more results               │  │
│ │ }()                                                        │  │
│ └───────────────────────────────────────────────────────────┘  │
│                                                                 │
│ ┌───────────────────────────────────────────────────────────┐  │
│ │ AGGREGATOR (Main goroutine, Sequential)                   │  │
│ │ for counts := range results {                             │  │
│ │     for word, count := range counts {                      │  │
│ │         finalCounts[word] += count                         │  │
│ │     }                                                      │  │
│ │ }                                                          │  │
│ └───────────────────────────────────────────────────────────┘  │
│                                                                 │
│ 5. Check errCh for errors                                      │
│ 6. Return finalCounts or error                                 │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ WORKER GOROUTINE (N instances, all concurrent)                  │
│                                                                 │
│ for {  ← Infinite loop until exit condition                    │
│     select {                                                    │
│     case <-ctx.Done():  ← Check cancellation                    │
│         return  ← Exit if cancelled                            │
│                                                                 │
│     case url, ok := <-jobs:  ← Receive job                     │
│         if !ok {                                                │
│             return  ← Channel closed, exit                      │
│         }                                                       │
│                                                                 │
│         counts, err := fetchAndCount(ctx, url)                  │
│         ↓                                                       │
│         ┌──────────────────────────────────────┐               │
│         │ fetchAndCount()                      │               │
│         │ 1. HTTP request with ctx             │               │
│         │ 2. Read response body                │               │
│         │ 3. tokenizeAndCount()                │               │
│         │    ↓                                  │               │
│         │    ┌──────────────────────────────┐  │               │
│         │    │ tokenizeAndCount()           │  │               │
│         │    │ 1. strings.Fields()         │  │               │
│         │    │ 2. ToLower each word         │  │               │
│         │    │ 3. Remove punctuation       │  │               │
│         │    │ 4. Count in map              │  │               │
│         │    └──────────────────────────────┘  │               │
│         │    Returns: map[string]int           │               │
│         └──────────────────────────────────────┘               │
│         ↓                                                       │
│         if err {                                                │
│             errCh <- err  ← Send error (non-blocking)          │
│             cancel()  ← Cancel context                          │
│             return  ← Exit worker                               │
│         }                                                       │
│                                                                 │
│         select {                                                │
│         case <-ctx.Done():                                      │
│             return  ← Cancelled while processing               │
│         case results <- counts:  ← Send result                  │
│         }                                                       │
│     }                                                           │
│ }  ← Loop back to check for next job                           │
└─────────────────────────────────────────────────────────────────┘
```

**Variable Flow Through Functions:**

```
┌─────────────────────────────────────────────────────────────┐
│ Variable: urls ([]string)                                   │
│ Location: Main goroutine                                    │
│ Memory: O(num_urls * avg_url_length)                        │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│ Variable: url (string)                                      │
│ Location: jobs channel buffer                                │
│ Memory: ~16-32 bytes per URL (string header)                │
│ Transformation: Each URL copied from slice to channel       │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│ Variable: url (string)                                      │
│ Location: Worker goroutine                                  │
│ Memory: ~16-32 bytes (copy of string header)               │
│ Received from: jobs channel                                 │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│ Function: fetchAndCount(ctx, url)                           │
│ Variables created:                                          │
│   - req (*http.Request) ~200 bytes                          │
│   - resp (*http.Response) ~500 bytes                        │
│   - body ([]byte) O(response_size)                         │
│   - text (string) O(response_size)                         │
│   - counts (map[string]int) O(vocabulary_size)              │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│ Variable: counts (map[string]int)                           │
│ Location: results channel buffer                            │
│ Memory: ~8 bytes (map header pointer)                       │
│ Transformation: Map reference copied to channel              │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│ Variable: counts (map[string]int)                           │
│ Location: Main goroutine (aggregator)                       │
│ Memory: Reference to worker's map                           │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│ Variable: finalCounts (map[string]int)                      │
│ Location: Main goroutine                                    │
│ Memory: O(total_unique_words)                               │
│ Transformation: Accumulates counts from all workers          │
│   finalCounts[word] += count  ← Accumulation pattern         │
└─────────────────────────────────────────────────────────────┘
```

**Understanding Counters and Iteration:**

**Loop Counter `i` in Worker Spawn:**
```go
for i := 0; i < workers; i++ {
    // i starts at 0
    // Each iteration: i = 0, 1, 2, ..., workers-1
    // After loop: i = workers (but loop doesn't execute)
    
    wg.Add(1)  // Counter increment: 1, 2, 3, ..., workers
    
    go func(workerID int) {
        // workerID is a COPY of i's value at this iteration
        // Worker 0 gets workerID=0
        // Worker 1 gets workerID=1
        // etc.
    }(i)  // ← Passing i as parameter creates a copy
}
```

**Why `(i)` on line 163?**

This is a **closure variable capture** issue. Without passing `i` as a parameter:

```go
// ❌ WRONG - All workers see the same value!
for i := 0; i < workers; i++ {
    go func() {
        fmt.Println(i)  // BUG: All print "workers" (final value)
    }()
}
// Problem: All goroutines capture the SAME variable i by reference.
// By the time goroutines run, loop finished, i = workers.

// ✅ CORRECT - Each worker gets its own copy
for i := 0; i < workers; i++ {
    go func(workerID int) {
        fmt.Println(workerID)  // Each prints 0, 1, 2, ...
    }(i)  // ← Pass i as argument, creates copy
}
// Solution: Each goroutine receives its own COPY of i's value.
```

**Where does `workerID` come from?**

`workerID` is the **parameter name** of the anonymous function:
```go
go func(workerID int) {  // ← workerID is the parameter name
    // Use workerID here
}(i)  // ← i is the argument passed to workerID
```

It's just a name - you could call it anything:
```go
go func(id int) { ... }(i)
go func(w int) { ... }(i)
go func(x int) { ... }(i)  // All equivalent
```

**Understanding `errCh`:**

`errCh` is a **buffered error channel** with size 1:

```go
errCh := make(chan error, 1)  // Buffer size = 1
```

**Why buffer size 1?**
- Ensures the **first** error sender never blocks
- If multiple workers error simultaneously, only first error is recorded
- Non-blocking send pattern prevents deadlock

**How it works:**
```go
// Worker encounters error:
select {
case errCh <- err:  // Try to send error
    cancel()  // Cancel context
default:  // Channel full (already has error), ignore
}
```

**Memory:** Only one error is stored (~16 bytes for error interface).

### Step 4: Synchronization Strategy

**Problem**: How do we know when all workers are done?

**Solution**: Use `sync.WaitGroup`:

**Understanding WaitGroup - First Principles:**

WaitGroup is a **counter** that tracks active goroutines:

```go
var wg sync.WaitGroup  // Counter starts at 0

wg.Add(1)   // Counter: 0 → 1
wg.Add(1)   // Counter: 1 → 2
wg.Done()   // Counter: 2 → 1 (decrements by 1)
wg.Done()   // Counter: 1 → 0
wg.Wait()   // Blocks until counter == 0
```

**How WaitGroup Works Internally:**

```
┌─────────────────────────────────────┐
│ sync.WaitGroup (struct)              │
│ ┌─────────────────────────────────┐  │
│ │ state [3]uint32                 │  │  ← Atomic counter + semaphore
│ │   [0] = counter (number waiting) │  │
│ │   [1] = waiter count            │  │
│ │   [2] = semaphore               │  │
│ └─────────────────────────────────┘  │
└─────────────────────────────────────┘

Operations:
- Add(delta): Atomically adds delta to counter
- Done(): Atomically decrements counter by 1
- Wait(): Blocks (using semaphore) until counter == 0
```

**Memory:** WaitGroup is ~12 bytes, uses atomic operations (lock-free, very fast).

**Pattern in our code:**
```go
var wg sync.WaitGroup

for i := 0; i < workers; i++ {
    wg.Add(1)  // Increment BEFORE starting goroutine
    
    go func() {
        defer wg.Done()  // Decrement when goroutine exits
        // work
    }()
}

// Later, in another goroutine:
wg.Wait()  // Block until all workers call Done()
```

**Why `defer wg.Done()`?**

`defer` ensures `wg.Done()` is called **even if the goroutine panics**:
```go
go func() {
    defer wg.Done()  // Always called, even on panic
    
    // If this panics, defer still executes
    riskyOperation()
}()
```

**Problem**: When do we close the results channel?

**Solution**: Close it after all workers are done:

**Why we need a separate goroutine to close results:**

```
Timeline:
---------
T0: Workers start processing URLs
T1: Worker 1 finishes → sends result → calls wg.Done()
T2: Worker 2 finishes → sends result → calls wg.Done()
T3: Worker 3 finishes → sends result → calls wg.Done()
...
TN: All workers done → wg.Wait() returns
TN+1: close(results) → Aggregator's range loop exits
```

If we closed results in main goroutine:
- We'd block on `wg.Wait()` (can't aggregate results)
- Workers would finish but results channel never closes
- Aggregator's `range` loop would block forever

**Solution:**
```go
go func() {
    wg.Wait()         // Wait for all workers (blocks here)
    close(results)    // Then close results channel
}()
// Main goroutine continues immediately to aggregate results
```

This allows the aggregator to use `range` over the results channel:
```go
for counts := range results {
    // Merge counts
    // Loop continues until results channel is closed
}
// Loop exits when results channel is closed AND empty
```

**Understanding `defer cancel()`:**

```go
ctx, cancel := context.WithCancel(ctx)
defer cancel()  // Always called when function returns
```

**Why defer?**
- Ensures cleanup happens even if function returns early
- Releases context resources (closes done channel)
- Prevents resource leaks

**What happens:**
1. Function starts: `defer cancel()` is scheduled
2. Function returns (normal or error): `cancel()` executes
3. Context's done channel is closed
4. All goroutines waiting on `<-ctx.Done()` wake up

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

### Worker Logic - Detailed Breakdown

```go
for {  // Infinite loop - exits via return statements
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

**Understanding the `select` Statement:**

`select` is like a `switch` for channels - it waits until ONE case can proceed:

```go
select {
case value := <-ch1:
    // ch1 has a value → receive it
case ch2 <- value:
    // ch2 has space → send value
case <-ctx.Done():
    // ctx.Done() channel is closed → execute this
default:
    // No case ready → execute immediately (non-blocking)
}
```

**Behavior:**
- If **multiple** cases are ready: ONE is chosen **randomly**
- If **no** case is ready and no `default`: **blocks** until one becomes ready
- If **no** case is ready but `default` exists: executes `default` immediately

**First `select`: Check cancellation OR receive job**

```go
select {
case <-ctx.Done():
    // Context cancelled → exit worker immediately
    return
    
case url, ok := <-jobs:
    // Try to receive from jobs channel
    // If channel has value: receive it, ok=true
    // If channel is empty but open: block until value arrives
    // If channel is closed: receive zero value, ok=false
}
```

**Understanding `url, ok := <-jobs`:**

This is the **two-value receive** syntax:
- `url`: The value received from channel (or zero value if closed)
- `ok`: Boolean indicating channel state
  - `true`: Successfully received value, channel is open
  - `false`: Channel is closed and empty

**Why check `ok`?**
```go
case url, ok := <-jobs:
    if !ok {
        // Channel closed → no more jobs coming
        // This is how workers know to exit
        return
    }
    // Channel open, got URL → process it
```

**Second `select`: Non-blocking error send**

```go
if err != nil {
    select {
    case errCh <- fmt.Errorf("fetching %s: %w", url, err):
        // Successfully sent error → cancel context
        cancel()
    default:
        // errCh is full (already has an error)
        // Ignore this error, first error already recorded
    }
    return
}
```

**Why non-blocking?**
- Multiple workers might error simultaneously
- We only care about the **first** error
- If errCh already has an error, we don't need to wait
- Prevents deadlock if error handler isn't reading

**Third `select`: Interruptible result send**

```go
select {
case <-ctx.Done():
    // Context cancelled while trying to send → exit
    return
case results <- counts:
    // Successfully sent result → continue loop
}
```

**Why check context here?**
- Cancellation might happen **between** receiving job and sending result
- Worker might have finished processing but context was cancelled
- We want to exit immediately, not send stale results

**Memory Flow in Worker:**

```
┌─────────────────────────────────────────┐
│ Worker receives: url (string)           │
│ Memory: ~16-32 bytes (string header)    │
└──────────────────┬──────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│ fetchAndCount(ctx, url)                  │
│ Allocates:                               │
│   - req: ~200 bytes                      │
│   - resp: ~500 bytes                     │
│   - body: O(response_size)               │
│   - counts: O(vocabulary_size)           │
└──────────────────┬──────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│ counts (map[string]int)                  │
│ Memory: Map header (~8 bytes) + data    │
│ Sent to: results channel                 │
└─────────────────────────────────────────┘
```

**Why three `select` statements?**

1. **First select**: Multiplex between cancellation check and job reception
   - Can't use separate if/else because both are blocking operations
   - Need to check BOTH simultaneously

2. **Second select**: Non-blocking error send
   - Don't want to block if errCh is full
   - `default` case makes it non-blocking

3. **Third select**: Interruptible result send
   - Check cancellation before/during send
   - Don't send results if context cancelled

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

### Aggregating Results - How Variables Morph

```go
finalCounts := make(map[string]int)  // Empty map
for counts := range results {
    for word, count := range counts {
        finalCounts[word] += count
    }
}
```

**Understanding `for counts := range results`:**

This is **range over channel** syntax:
- Receives values from `results` channel
- Continues until channel is **closed AND empty**
- Each iteration gets one `map[string]int` from a worker

**Variable Morphing Through Iterations:**

```
Iteration 1:
  counts = {"hello": 2, "world": 1}  (from worker 1)
  finalCounts = {}  (empty)
  After merge: finalCounts = {"hello": 2, "world": 1}

Iteration 2:
  counts = {"hello": 1, "go": 3}  (from worker 2)
  finalCounts = {"hello": 2, "world": 1}  (from iteration 1)
  After merge:
    finalCounts["hello"] += 1  → finalCounts["hello"] = 3
    finalCounts["go"] += 3     → finalCounts["go"] = 3
  Result: finalCounts = {"hello": 3, "world": 1, "go": 3}

Iteration 3:
  counts = {"world": 2, "go": 1}  (from worker 3)
  finalCounts = {"hello": 3, "world": 1, "go": 3}  (from iteration 2)
  After merge:
    finalCounts["world"] += 2  → finalCounts["world"] = 3
    finalCounts["go"] += 1     → finalCounts["go"] = 4
  Result: finalCounts = {"hello": 3, "world": 3, "go": 4}

... continues until results channel is closed ...
```

**Understanding `finalCounts[word] += count`:**

This is **map accumulation**:
```go
finalCounts[word] += count
```

**How it works:**
1. Look up `word` in `finalCounts`
2. If doesn't exist: Go creates entry with zero value (0 for int)
3. Add `count` to the value
4. Store result back in map

**Example:**
```go
finalCounts["hello"] = 2      // Set initial value
finalCounts["hello"] += 1     // Add 1
// finalCounts["hello"] = 3   // Result

finalCounts["new"] += 5       // "new" doesn't exist
// Go creates finalCounts["new"] = 0
// Then adds 5: finalCounts["new"] = 5
```

**Why is this safe without locks?**

Only **ONE goroutine** (main) writes to `finalCounts`:
- Workers send results through **channel** (synchronization built-in)
- Channels provide **thread-safe** communication
- No shared mutable state → no race conditions
- Main goroutine receives results sequentially (one at a time)

**Memory Management:**

```
┌─────────────────────────────────────────┐
│ finalCounts (map[string]int)             │
│ Memory: O(total_unique_words)           │
│ Grows as words are added                │
│ Rehashes when load factor > 6.5         │
└─────────────────────────────────────────┘

Each iteration:
  - counts: Reference to worker's map (not copied)
  - Iteration over counts: O(vocabulary_size) time
  - Accumulation: O(vocabulary_size) time
  - Total per iteration: O(vocabulary_size)
```

**What happens when `results` is closed?**

1. Channel is closed by results closer goroutine
2. `for range` receives remaining values (if any)
3. When channel is empty AND closed, `range` loop exits
4. Proceeds to error checking

**Timing:**
```
T0: All workers finish → wg.Wait() returns
T1: close(results) called
T2: Aggregator receives last result
T3: Aggregator's range loop exits (channel empty and closed)
T4: Error check executes
```

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

## 5. Understanding Tokenization

**What is Tokenization?**

Tokenization is the process of breaking text into individual words (tokens).

**Example:**
```
Input:  "Hello, world! Go is great."
Output: ["hello", "world", "go", "is", "great"]
```

**How Tokenization Works in Our Code:**

```
┌─────────────────────────────────────────┐
│ Input: text (string)                    │
│ "Hello, world! Go is great."            │
└──────────────────┬──────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│ Step 1: strings.Fields(text)           │
│ Splits on whitespace (space, tab, \n)   │
│ Result: ["Hello,", "world!", "Go", ...] │
└──────────────────┬──────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│ Step 2: strings.ToLower(word)           │
│ Converts to lowercase                   │
│ Result: ["hello,", "world!", "go", ...] │
└──────────────────┬──────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│ Step 3: strings.Map(...)                │
│ Remove non-letter characters            │
│ Result: ["hello", "world", "go", ...]   │
└──────────────────┬──────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│ Step 4: Count in map                    │
│ counts["hello"] = 1                     │
│ counts["world"] = 1                     │
│ counts["go"] = 1                        │
│ ...                                      │
└─────────────────────────────────────────┘
```

**Variable Transformation Through Tokenization:**

```
Iteration 1:
  word = "Hello,"        (from strings.Fields)
  word = "hello,"        (after ToLower)
  word = "hello"         (after Map removes comma)
  counts["hello"] = 1

Iteration 2:
  word = "world!"
  word = "world!"
  word = "world"
  counts["world"] = 1

Iteration 3:
  word = "Go"
  word = "go"
  word = "go"
  counts["go"] = 1

... continues for each word ...
```

**Memory Considerations:**

- `strings.Fields`: Allocates slice of strings (O(text_length))
- `strings.ToLower`: May allocate new string if changes needed
- `strings.Map`: Allocates new string for each word (O(word_length))
- Map: Grows dynamically as words are added (O(vocabulary_size))

**Total memory:** O(text_length + vocabulary_size)

---

## 6. Using errgroup: A Safer, More Concise Alternative

### What is errgroup?

`errgroup` (from `golang.org/x/sync/errgroup`) is a package that provides **synchronization, error propagation, and context cancellation** for groups of goroutines.

**Key benefits:**
- Automatic context cancellation on first error
- Built-in WaitGroup functionality
- Cleaner, more concise code
- Less boilerplate than manual WaitGroup + context management

### How errgroup Works Under the Hood

**Understanding the errgroup Package Internals:**

```go
// Simplified version of errgroup.Group (from golang.org/x/sync/errgroup)
type Group struct {
    cancel func()              // Context cancellation function
    wg     sync.WaitGroup     // WaitGroup for goroutine tracking
    errOnce sync.Once         // Ensures only first error is stored
    err    error              // First error encountered
    errMu  sync.Mutex         // Protects err field
}

func WithContext(ctx context.Context) (*Group, context.Context) {
    ctx, cancel := context.WithCancel(ctx)
    return &Group{cancel: cancel}, ctx
}

func (g *Group) Go(f func() error) {
    g.wg.Add(1)
    go func() {
        defer g.wg.Done()
        if err := f(); err != nil {
            g.errOnce.Do(func() {
                g.errMu.Lock()
                g.err = err
                g.errMu.Unlock()
                g.cancel()  // Cancel context on first error
            })
        }
    }()
}

func (g *Group) Wait() error {
    g.wg.Wait()
    g.cancel()  // Always cancel context when done
    return g.err
}
```

**How It Works Step-by-Step:**

```
┌─────────────────────────────────────────────────────────┐
│ errgroup.WithContext(ctx)                               │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ 1. Creates cancellable context:                     │ │
│ │    ctx, cancel := context.WithCancel(ctx)           │ │
│ │                                                      │ │
│ │ 2. Creates Group struct:                            │ │
│ │    g := &Group{                                     │ │
│ │        cancel: cancel,                               │ │
│ │        wg: sync.WaitGroup{},                        │ │
│ │        errOnce: sync.Once{},                        │ │
│ │        err: nil,                                    │ │
│ │        errMu: sync.Mutex{},                         │ │
│ │    }                                                 │ │
│ │                                                      │ │
│ │ 3. Returns group and new context                    │ │
│ └─────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────┐
│ g.Go(func() error { ... })                               │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ 1. g.wg.Add(1)  ← Increment WaitGroup counter       │ │
│ │                                                      │ │
│ │ 2. Launch goroutine:                                │ │
│ │    go func() {                                       │ │
│ │        defer g.wg.Done()  ← Decrement on exit      │ │
│ │                                                      │ │
│ │        if err := f(); err != nil {                  │ │
│ │            g.errOnce.Do(func() {                   │ │
│ │                // Only first error is stored        │ │
│ │                g.errMu.Lock()                       │ │
│ │                g.err = err                           │ │
│ │                g.errMu.Unlock()                    │ │
│ │                g.cancel()  ← Cancel context         │ │
│ │            })                                        │ │
│ │        }                                             │ │
│ │    }()                                               │ │
│ └─────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────┐
│ g.Wait()                                                 │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ 1. g.wg.Wait()  ← Block until all goroutines done   │ │
│ │                                                      │ │
│ │ 2. g.cancel()  ← Always cancel context              │ │
│ │                                                      │ │
│ │ 3. return g.err  ← Return first error (if any)     │ │
│ └─────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

**Key Mechanisms:**

1. **WaitGroup Integration:**
   - `g.Go()` automatically calls `wg.Add(1)` before starting goroutine
   - `defer wg.Done()` is handled internally
   - `g.Wait()` calls `wg.Wait()` to wait for all goroutines

2. **Error Handling:**
   - Uses `sync.Once` to ensure only the **first** error is stored
   - Uses `sync.Mutex` to protect the `err` field (thread-safe)
   - When first error occurs, `cancel()` is called automatically

3. **Context Cancellation:**
   - Context is created with `WithContext()`
   - Automatically cancelled when first error occurs
   - Always cancelled when `Wait()` returns (cleanup)

4. **Memory Management:**
   - Group struct: ~64 bytes (WaitGroup + mutexes + function pointer)
   - No additional channels needed (unlike manual approach)
   - More memory-efficient than manual error channel pattern

### Step-by-Step Execution: How errgroup Works in Practice

Let's trace through exactly what happens when we use errgroup with our WordCount function:

**Initialization:**

```go
g, ctx := errgroup.WithContext(ctx)
```

**What happens internally:**

```
1. errgroup.WithContext(ctx) is called
   ↓
2. Creates cancellable context:
   ctx, cancel := context.WithCancel(ctx)
   ↓
3. Creates Group struct:
   g := &Group{
       cancel: cancel,           // Function to cancel context
       wg: sync.WaitGroup{},     // Counter starts at 0
       errOnce: sync.Once{},      // Ensures only first error stored
       err: nil,                  // No error yet
       errMu: sync.Mutex{},       // Protects err field
   }
   ↓
4. Returns: (g, ctx)
```

**Launching Workers:**

```go
g.Go(func() error {
    counts, err := fetchAndCount(ctx, url)
    if err != nil {
        return err  // ← What happens here?
    }
    return nil
})
```

**What happens internally when g.Go() is called:**

```
1. g.Go(func() error { ... }) is called
   ↓
2. g.wg.Add(1)  ← Increment WaitGroup counter
   Counter: 0 → 1 → 2 → ... → workers
   ↓
3. Launch goroutine:
   go func() {
       defer g.wg.Done()  ← Schedule decrement on exit
       ↓
       if err := f(); err != nil {  ← Execute function
           ↓
           g.errOnce.Do(func() {  ← Only first error!
               g.errMu.Lock()     ← Lock mutex
               g.err = err        ← Store error
               g.errMu.Unlock()   ← Unlock mutex
               g.cancel()         ← Cancel context!
           })
       }
   }()
   ↓
4. Main goroutine continues immediately (doesn't wait)
```

**What happens when error occurs:**

```
Worker 1 encounters error:
  ↓
return fmt.Errorf("fetching %s: %w", url, err)
  ↓
g.errOnce.Do(func() {  ← sync.Once ensures this runs only once
    g.errMu.Lock()     ← Lock mutex (thread-safe)
    g.err = err        ← Store first error
    g.errMu.Unlock()   ← Unlock mutex
    g.cancel()         ← Cancel context!
})
  ↓
Context cancelled:
  - ctx.Done() channel is CLOSED
  - All other workers see <-ctx.Done() unblock
  - All other workers exit immediately
  ↓
Worker 1 calls defer g.wg.Done()
  Counter: workers → workers-1
  ↓
Other workers exit:
  - See ctx.Done() closed
  - Call defer g.wg.Done()
  Counter: workers-1 → workers-2 → ... → 0
```

**What happens when g.Wait() is called:**

```go
if err := g.Wait(); err != nil {
    return nil, err
}
```

**Internal execution:**

```
1. g.Wait() is called
   ↓
2. g.wg.Wait()  ← Block until counter reaches 0
   - If counter > 0: Block (wait for goroutines)
   - If counter == 0: Continue immediately
   ↓
3. g.cancel()  ← Always cancel context (cleanup)
   - Releases context resources
   - Closes done channel
   ↓
4. return g.err  ← Return first error (or nil)
   - If error occurred: return that error
   - If no error: return nil
```

### Detailed Execution Timeline

**Scenario: 3 workers, 5 URLs, 1 error on URL 3**

```
T0: g, ctx := errgroup.WithContext(ctx)
    - Group created, counter = 0
    - Context created (not cancelled)

T1: g.Go(worker1), g.Go(worker2), g.Go(worker3)
    - Counter: 0 → 1 → 2 → 3
    - All workers start processing

T2: Worker1 receives URL1, starts fetching
    Worker2 receives URL2, starts fetching
    Worker3 receives URL3, starts fetching

T3: Worker1 finishes URL1 successfully
    - Sends result to results channel
    - Receives URL4, starts fetching

T4: Worker3 encounters error on URL3
    - return fmt.Errorf("fetching URL3: %w", err)
    ↓
    g.errOnce.Do(func() {
        g.errMu.Lock()
        g.err = <error>  ← First error stored
        g.errMu.Unlock()
        g.cancel()  ← Context cancelled!
    })
    ↓
    ctx.Done() channel CLOSED
    ↓
    Worker1: <-ctx.Done() unblocks → return ctx.Err()
    Worker2: <-ctx.Done() unblocks → return ctx.Err()
    Worker3: defer g.wg.Done() → counter: 3 → 2

T5: Worker1 exits (saw ctx.Done())
    - defer g.wg.Done() → counter: 2 → 1

T6: Worker2 exits (saw ctx.Done())
    - defer g.wg.Done() → counter: 1 → 0

T7: g.Wait() called
    - g.wg.Wait() returns (counter == 0)
    - g.cancel() called (cleanup, already cancelled but safe)
    - return g.err ← Returns the error from URL3

T8: if err := g.Wait(); err != nil {
        return nil, err  ← Error returned to caller
    }
```

**Key Observations:**

1. **First Error Wins:**
   - Worker3's error is stored
   - sync.Once ensures only first error stored
   - Even if Worker1 or Worker2 error later, their errors are ignored

2. **Automatic Cancellation:**
   - `g.cancel()` called automatically
   - All workers see `ctx.Done()` and exit
   - No manual cancellation needed

3. **Automatic Cleanup:**
   - `g.Wait()` always calls `g.cancel()`
   - Even if no error occurred
   - Context resources always released

4. **Thread Safety:**
   - `sync.Mutex` protects `err` field
   - `sync.Once` ensures only first error stored
   - Multiple workers can error simultaneously safely

### Detailed Comparison: Manual vs errgroup Applied to WordCount

Let's compare the **actual WordCount function** implementations side-by-side to see exactly how errgroup simplifies the code.

#### **Manual Approach (Current Implementation):**

```go
func WordCount(ctx context.Context, urls []string, workers int) (map[string]int, error) {
    // STEP 1: Create channels
    jobs := make(chan string, workers)
    results := make(chan map[string]int, workers)
    errCh := make(chan error, 1)  // ← Manual error channel
    
    // STEP 2: Create cancellable context
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()  // ← Must remember to defer!
    
    // STEP 3: Create WaitGroup
    var wg sync.WaitGroup  // ← Manual WaitGroup
    
    // STEP 4: Start workers
    for i := 0; i < workers; i++ {
        wg.Add(1)  // ← Must call before go!
        go func(workerID int) {
            defer wg.Done()  // ← Must remember defer!
            for {
                select {
                case <-ctx.Done():
                    return
                case url, ok := <-jobs:
                    if !ok {
                        return
                    }
                    counts, err := fetchAndCount(ctx, url)
                    if err != nil {
                        // ← Complex error handling
                        select {
                        case errCh <- fmt.Errorf("fetching %s: %w", url, err):
                            cancel()  // ← Manual cancellation
                        default:
                            // ← Non-blocking send pattern
                        }
                        return
                    }
                    select {
                    case <-ctx.Done():
                        return
                    case results <- counts:
                    }
                }
            }
        }(i)
    }
    
    // STEP 5: Send jobs
    go func() {
        for _, url := range urls {
            select {
            case <-ctx.Done():
                return
            case jobs <- url:
            }
        }
        close(jobs)
    }()
    
    // STEP 6: Close results when done
    go func() {
        wg.Wait()  // ← Manual WaitGroup wait
        close(results)
    }()
    
    // STEP 7: Aggregate results
    finalCounts := make(map[string]int)
    for counts := range results {
        for word, count := range counts {
            finalCounts[word] += count
        }
    }
    
    // STEP 8: Check for errors
    select {  // ← Complex error checking
    case err := <-errCh:
        return nil, err
    default:
    }
    
    return finalCounts, nil
}
```

**Problems with Manual Approach:**

1. **Error Channel Complexity:**
   ```go
   errCh := make(chan error, 1)  // Must create channel
   // Later...
   select {
   case errCh <- err:  // Non-blocking send
       cancel()
   default:
   }
   // Even later...
   select {
   case err := <-errCh:  // Non-blocking receive
       return nil, err
   default:
   }
   ```
   - Must remember buffer size 1
   - Must use non-blocking send pattern
   - Must use non-blocking receive pattern
   - Easy to forget and cause deadlock

2. **WaitGroup Management:**
   ```go
   wg.Add(1)  // Must call BEFORE go
   go func() {
       defer wg.Done()  // Must remember defer
       // ...
   }()
   // Later...
   wg.Wait()  // Must call to wait
   ```
   - Easy to forget `wg.Add(1)` → race condition
   - Easy to forget `defer wg.Done()` → deadlock
   - Must coordinate Add/Done/Wait carefully

3. **Context Cleanup:**
   ```go
   ctx, cancel := context.WithCancel(ctx)
   defer cancel()  // Must remember!
   ```
   - Easy to forget `defer cancel()` → resource leak
   - Must remember even on success path

4. **Error Propagation:**
   ```go
   if err != nil {
       select {
       case errCh <- err:
           cancel()  // Must manually cancel
       default:
       }
       return
   }
   ```
   - Must manually send to error channel
   - Must manually cancel context
   - Complex non-blocking pattern

**Total Manual Boilerplate:**
- Error channel creation: 1 line
- Error channel send: 5 lines (select + default)
- Error channel receive: 4 lines (select + default)
- WaitGroup management: 3 lines (Add, Done, Wait)
- Context cleanup: 2 lines (WithCancel, defer cancel)
- **Total: ~15 lines of error-prone boilerplate**

#### **errgroup Approach (Simplified Implementation):**

```go
import "golang.org/x/sync/errgroup"

func WordCountWithErrGroup(ctx context.Context, urls []string, workers int) (map[string]int, error) {
    // STEP 1: Create errgroup with context
    g, ctx := errgroup.WithContext(ctx)  // ← Automatic context + cleanup!
    
    // STEP 2: Create channels (still needed for worker pool)
    jobs := make(chan string, workers)
    results := make(chan map[string]int, workers)
    // ← No error channel needed!
    
    // STEP 3: Start workers using errgroup
    for i := 0; i < workers; i++ {
        g.Go(func() error {  // ← Automatic WaitGroup management!
            for {
                select {
                case <-ctx.Done():
                    return ctx.Err()
                case url, ok := <-jobs:
                    if !ok {
                        return nil  // Success
                    }
                    counts, err := fetchAndCount(ctx, url)
                    if err != nil {
                        return fmt.Errorf("fetching %s: %w", url, err)  // ← Just return!
                    }
                    select {
                    case <-ctx.Done():
                        return ctx.Err()
                    case results <- counts:
                    }
                }
            }
        })
    }
    
    // STEP 4: Send jobs (same as manual)
    go func() {
        defer close(jobs)
        for _, url := range urls {
            select {
            case <-ctx.Done():
                return
            case jobs <- url:
            }
        }
    }()
    
    // STEP 5: Close results when done
    go func() {
        _ = g.Wait()  // ← Automatic WaitGroup wait + cleanup
        close(results)
    }()
    
    // STEP 6: Aggregate results (same as manual)
    finalCounts := make(map[string]int)
    for counts := range results {
        for word, count := range counts {
            finalCounts[word] += count
        }
    }
    
    // STEP 7: Check for errors (simplified!)
    if err := g.Wait(); err != nil {  // ← Just check error, no select!
        return nil, err
    }
    
    return finalCounts, nil
}
```

**Benefits of errgroup:**

1. **No Error Channel:**
   ```go
   // Before: errCh := make(chan error, 1)
   // After: Nothing! Just return error
   ```
   - No channel creation
   - No non-blocking send pattern
   - No non-blocking receive pattern
   - Just return error from goroutine

2. **Automatic WaitGroup Management:**
   ```go
   // Before:
   wg.Add(1)
   go func() {
       defer wg.Done()
       // ...
   }()
   
   // After:
   g.Go(func() error {
       // ... Add/Done handled automatically!
   })
   ```
   - Can't forget `wg.Add(1)` → impossible
   - Can't forget `defer wg.Done()` → automatic
   - `g.Wait()` handles everything

3. **Automatic Context Cleanup:**
   ```go
   // Before:
   ctx, cancel := context.WithCancel(ctx)
   defer cancel()
   
   // After:
   g, ctx := errgroup.WithContext(ctx)  // Cleanup automatic!
   ```
   - No `defer cancel()` needed
   - `g.Wait()` always cancels context
   - Even on success path

4. **Simplified Error Propagation:**
   ```go
   // Before:
   select {
   case errCh <- err:
       cancel()
   default:
   }
   
   // After:
   return err  // That's it!
   ```
   - Just return error
   - errgroup handles cancellation
   - Thread-safe (sync.Once + Mutex)

**Total errgroup Boilerplate:**
- Error handling: 1 line (just return)
- WaitGroup management: 0 lines (automatic)
- Context cleanup: 0 lines (automatic)
- Error checking: 2 lines (simple if statement)
- **Total: ~3 lines vs ~15 lines manual**

### Side-by-Side Code Comparison

**Error Handling:**

| Manual | errgroup |
|--------|----------|
| `errCh := make(chan error, 1)` | (not needed) |
| `select { case errCh <- err: cancel() default: }` | `return err` |
| `select { case err := <-errCh: return err default: }` | `if err := g.Wait(); err != nil { return err }` |

**WaitGroup Management:**

| Manual | errgroup |
|--------|----------|
| `wg.Add(1)` before `go` | (automatic in `g.Go()`) |
| `defer wg.Done()` | (automatic in `g.Go()`) |
| `wg.Wait()` | `g.Wait()` (also cleans up context) |

**Context Cleanup:**

| Manual | errgroup |
|--------|----------|
| `ctx, cancel := context.WithCancel(ctx)` | `g, ctx := errgroup.WithContext(ctx)` |
| `defer cancel()` | (automatic in `g.Wait()`) |

### Why errgroup is Safer

**1. Impossible to Forget WaitGroup Operations:**
```go
// Manual: Easy to forget
wg.Add(1)  // ← What if you forget this?
go func() {
    defer wg.Done()  // ← What if you forget this?
    // ...
}()

// errgroup: Impossible to forget
g.Go(func() error {  // ← Add/Done automatic, can't forget
    // ...
})
```

**2. Impossible to Forget Context Cleanup:**
```go
// Manual: Easy to forget
ctx, cancel := context.WithCancel(ctx)
// defer cancel()  // ← What if you forget this? Resource leak!

// errgroup: Automatic cleanup
g, ctx := errgroup.WithContext(ctx)
// g.Wait() always cancels, even if you forget to call it explicitly
```

**3. Thread-Safe Error Handling:**
```go
// Manual: Must be careful
errCh := make(chan error, 1)  // Buffer size 1
select {
case errCh <- err:  // What if multiple errors? Must use default
    cancel()
default:
}

// errgroup: Thread-safe by design
return err  // sync.Once ensures only first error stored
            // sync.Mutex protects err field
```

**4. Simpler Code = Fewer Bugs:**
- Less code to write = fewer places for bugs
- Less code to read = easier to understand
- Less code to test = easier to verify correctness

### Performance Comparison

**Memory:**

| Component | Manual | errgroup |
|-----------|--------|----------|
| Error channel | 16 bytes | 0 bytes |
| WaitGroup | 12 bytes | (included) |
| Context management | 48 bytes | (included) |
| Group struct | 0 bytes | 64 bytes |
| **Total** | **76 bytes** | **64 bytes** |

**Code Complexity:**

| Metric | Manual | errgroup |
|--------|--------|----------|
| Lines of coordination code | ~100 | ~70 |
| Error handling complexity | High | Low |
| Potential bugs | Many | Few |
| Cognitive load | High | Low |

### When to Use Each Approach

**Use Manual Approach When:**
- ✅ Need fine-grained control over error handling
- ✅ Complex coordination patterns
- ✅ Want to understand every detail
- ✅ Learning concurrency fundamentals

**Use errgroup When:**
- ✅ Want simpler, safer code
- ✅ Fan-out patterns (one goroutine per task)
- ✅ Need automatic error propagation
- ✅ Production code (less error-prone)

**For Worker Pools:**
- Both approaches work
- errgroup simplifies error handling
- Still need channels for job distribution
- errgroup reduces boilerplate significantly

### Why errgroup Makes Code More Efficient

**1. Reduced Cognitive Load:**

**Manual Approach Mental Model:**
```
Developer must track:
- WaitGroup counter (when to Add/Done/Wait)
- Error channel (buffer size, non-blocking patterns)
- Context cancellation (when to cancel, defer cleanup)
- Error propagation (how errors flow through channels)
- Coordination (multiple goroutines, channels, cleanup)

Total: 5+ concepts to coordinate
```

**errgroup Mental Model:**
```
Developer must track:
- g.Go() launches goroutine (handles WaitGroup)
- return error propagates automatically
- g.Wait() waits and returns error

Total: 3 concepts, simpler
```

**2. Compiler-Enforced Safety:**

**Manual Approach - Runtime Errors Possible:**
```go
// ❌ Can forget wg.Add(1) → race condition
go func() {
    defer wg.Done()  // Counter never incremented!
    // ...
}()

// ❌ Can forget defer wg.Done() → deadlock
wg.Add(1)
go func() {
    // No defer wg.Done() → counter never decrements
    // wg.Wait() blocks forever!
}()

// ❌ Can forget defer cancel() → resource leak
ctx, cancel := context.WithCancel(ctx)
// No defer cancel() → context never cleaned up
```

**errgroup - Compile-Time Safety:**
```go
// ✅ Impossible to forget wg.Add(1)
g.Go(func() error {  // Add() called automatically
    // ...
})

// ✅ Impossible to forget wg.Done()
g.Go(func() error {  // Done() called automatically via defer
    // ...
})

// ✅ Impossible to forget context cleanup
g, ctx := errgroup.WithContext(ctx)
// g.Wait() always cancels context
```

**3. Fewer Lines of Code = Fewer Bugs:**

**Code Complexity Metrics:**

| Metric | Manual | errgroup | Improvement |
|--------|--------|----------|-------------|
| Lines of coordination code | 100 | 70 | 30% reduction |
| Error handling lines | 9 | 1 | 89% reduction |
| WaitGroup management lines | 3 | 0 | 100% reduction |
| Context cleanup lines | 2 | 0 | 100% reduction |
| **Total boilerplate** | **~15 lines** | **~3 lines** | **80% reduction** |

**Bug Probability:**

Research shows bug density is roughly constant per line of code (~1-5 bugs per 1000 lines). Fewer lines = fewer bugs.

- Manual: ~100 lines → ~0.1-0.5 potential bugs
- errgroup: ~70 lines → ~0.07-0.35 potential bugs
- **~30% reduction in potential bugs**

**4. Better Error Handling:**

**Manual Approach - Complex Pattern:**
```go
errCh := make(chan error, 1)  // Must remember buffer size 1

// In goroutine:
select {
case errCh <- err:  // What if channel full?
    cancel()
default:  // Must remember default!
    // Error lost? Or ignored?
}

// Later:
select {
case err := <-errCh:  // What if no error?
    return err
default:  // Must remember default!
    // Continue?
}
```

**Problems:**
- Easy to forget `default` → deadlock
- Easy to use wrong buffer size → deadlock
- Complex mental model (non-blocking patterns)

**errgroup - Simple Pattern:**
```go
// In goroutine:
return err  // That's it!

// Later:
if err := g.Wait(); err != nil {
    return err
}
```

**Benefits:**
- No channels needed
- No non-blocking patterns
- Simple: just return error
- Thread-safe by design

**5. Automatic Resource Cleanup:**

**Manual Approach - Easy to Leak:**
```go
ctx, cancel := context.WithCancel(ctx)
// What if function returns early?
// What if panic occurs?
// Must remember defer cancel() everywhere!

defer cancel()  // ← Easy to forget
```

**errgroup - Always Clean:**
```go
g, ctx := errgroup.WithContext(ctx)
// g.Wait() ALWAYS cancels context
// Even if you forget to call it explicitly
// Even if function panics (if Wait() was called)
```

**6. Better Testing:**

**Manual Approach - Many Edge Cases:**
```go
// Must test:
// - What if wg.Add(1) forgotten?
// - What if defer wg.Done() forgotten?
// - What if defer cancel() forgotten?
// - What if error channel full?
// - What if multiple errors?
// - What if context cancelled before error?
```

**errgroup - Fewer Edge Cases:**
```go
// Must test:
// - What if error returned? (automatic cancellation)
// - What if no error? (automatic cleanup)
// - What if multiple errors? (first wins, automatic)
```

### Real-World Impact

**Example: Production Bug Prevention**

**Scenario:** Worker pool processing 10,000 URLs

**Manual Approach Bug:**
```go
// Developer forgot defer cancel()
ctx, cancel := context.WithCancel(ctx)
// No defer! → Resource leak

// After processing 10,000 URLs:
// - 10,000 contexts created
// - 0 contexts cleaned up
// - Memory leak grows over time
// - Eventually: Out of memory crash
```

**errgroup Prevention:**
```go
g, ctx := errgroup.WithContext(ctx)
// g.Wait() always cancels
// Impossible to forget cleanup
// No memory leak possible
```

**Example: Error Handling Bug**

**Manual Approach Bug:**
```go
// Developer used unbuffered error channel
errCh := make(chan error)  // ← Wrong! Should be size 1

// Worker tries to send error:
errCh <- err  // ← BLOCKS! (no receiver yet)
// Deadlock! All workers blocked waiting to send error
```

**errgroup Prevention:**
```go
// No error channel needed!
return err  // ← Never blocks
// errgroup handles it internally (non-blocking)
```

### Summary: Why errgroup is Better

**Safety:**
- ✅ Impossible to forget WaitGroup operations
- ✅ Impossible to forget context cleanup
- ✅ Thread-safe error handling built-in
- ✅ Compiler-enforced patterns

**Efficiency:**
- ✅ 30% less code
- ✅ 80% less boilerplate
- ✅ Simpler mental model
- ✅ Fewer potential bugs

**Maintainability:**
- ✅ Easier to read
- ✅ Easier to understand
- ✅ Easier to test
- ✅ Easier to modify

**When to Use:**
- ✅ Production code (safer)
- ✅ Fan-out patterns (ideal)
- ✅ Worker pools (simplifies error handling)
- ✅ Any concurrent code with error handling

**When Manual is OK:**
- ✅ Learning concurrency fundamentals
- ✅ Need fine-grained control
- ✅ Complex coordination patterns
- ✅ Understanding every detail

**However:** For worker pools, we still need channels for job distribution. errgroup is best for **fan-out patterns** where each goroutine handles one task.

### Better Use Case: Fan-Out Pattern with errgroup

**When errgroup Really Shines:**

```go
import (
    "context"
    "golang.org/x/sync/errgroup"
)

func WordCountFanOut(ctx context.Context, urls []string, maxConcurrency int) (map[string]int, error) {
    // Create errgroup with context and concurrency limit
    g, ctx := errgroup.WithContext(ctx)
    g.SetLimit(maxConcurrency)  // Limit concurrent goroutines
    
    // Channel for results
    results := make(chan map[string]int, len(urls))
    
    // Launch one goroutine per URL (fan-out)
    for _, url := range urls {
        url := url  // Capture loop variable
        g.Go(func() error {
            // Fetch and count words
            counts, err := fetchAndCount(ctx, url)
            if err != nil {
                return fmt.Errorf("fetching %s: %w", url, err)
            }
            
            // Send results
            select {
            case <-ctx.Done():
                return ctx.Err()
            case results <- counts:
                return nil
            }
        })
    }
    
    // Close results when all goroutines done
    go func() {
        g.Wait()
        close(results)
    }()
    
    // Aggregate results
    finalCounts := make(map[string]int)
    for counts := range results {
        for word, count := range counts {
            finalCounts[word] += count
        }
    }
    
    // Wait and return error if any
    if err := g.Wait(); err != nil {
        return nil, err
    }
    
    return finalCounts, nil
}
```

**Why This is Better:**

1. **Concurrency Limiting:**
   ```go
   g.SetLimit(maxConcurrency)
   ```
   - Automatically limits concurrent goroutines
   - No need for worker pool pattern
   - Simpler code

2. **Error Handling:**
   ```go
   g.Go(func() error {
       // Just return error - errgroup handles cancellation
       if err != nil {
           return err  // Automatically cancels context
       }
       return nil
   })
   ```
   - No manual error channel
   - No manual context cancellation
   - First error automatically cancels all others

3. **Cleanup:**
   ```go
   if err := g.Wait(); err != nil {
       return nil, err
   }
   ```
   - Context automatically cancelled
   - All goroutines guaranteed to finish
   - Error returned if any occurred

### Memory Comparison

**Manual Approach:**
```
- errCh channel: 1 * sizeof(error) ≈ 16 bytes
- Manual WaitGroup: ~12 bytes
- Context management: ~48 bytes
- Error handling logic: ~200 bytes (code)
Total overhead: ~276 bytes + coordination complexity
```

**errgroup Approach:**
```
- errgroup.Group: ~64 bytes
- Built-in WaitGroup: included
- Built-in context: included
- Error handling: included
Total overhead: ~64 bytes, simpler code
```

### When to Use errgroup vs Manual Approach

**Use errgroup when:**
- ✅ Fan-out pattern (one goroutine per task)
- ✅ Need automatic error propagation
- ✅ Want simpler code
- ✅ Don't need complex worker pool coordination

**Use Manual Approach when:**
- ✅ Worker pool pattern (reuse goroutines for multiple tasks)
- ✅ Need fine-grained control over error handling
- ✅ Complex coordination between goroutines
- ✅ Need custom synchronization patterns

### Complete errgroup Example

```go
package main

import (
    "context"
    "fmt"
    "golang.org/x/sync/errgroup"
    "net/http"
    "time"
)

func fetchURL(ctx context.Context, url string) (string, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return "", err
    }
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
    }
    
    // Read body...
    return "body content", nil
}

func main() {
    urls := []string{
        "https://example.com/1",
        "https://example.com/2",
        "https://example.com/3",
    }
    
    // Create errgroup with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    g, ctx := errgroup.WithContext(ctx)
    g.SetLimit(2)  // Max 2 concurrent requests
    
    results := make(chan string, len(urls))
    
    // Launch goroutines
    for _, url := range urls {
        url := url  // Capture loop variable
        g.Go(func() error {
            body, err := fetchURL(ctx, url)
            if err != nil {
                return fmt.Errorf("fetch %s: %w", url, err)
            }
            results <- body
            return nil
        })
    }
    
    // Close results when done
    go func() {
        g.Wait()
        close(results)
    }()
    
    // Collect results
    for result := range results {
        fmt.Println(result)
    }
    
    // Check for errors
    if err := g.Wait(); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

**Key Takeaways:**
1. errgroup combines WaitGroup + Context + Error handling
2. Automatically cancels context on first error
3. Simplifies error propagation (just return error)
4. Reduces boilerplate significantly
5. Best for fan-out patterns, less ideal for worker pools

---

## 7. Key Concepts Explained

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

### Concept 2: Channel Semantics - Deep Dive

**What is a Channel?**

A channel is a **typed, thread-safe queue** for communication between goroutines.

**Memory Layout:**

```
┌─────────────────────────────────────┐
│ chan T (channel of type T)          │
│ ┌─────────────────────────────────┐  │
│ │ buffer []T (circular buffer)    │  │  ← Buffered channels only
│ │ sendq waitq (waiting senders)   │  │
│ │ recvq waitq (waiting receivers) │  │
│ │ lock mutex (synchronization)    │  │
│ │ closed uint32 (closed flag)     │  │
│ └─────────────────────────────────┘  │
└─────────────────────────────────────┘
```

**Unbuffered vs Buffered Channels**

**Unbuffered (synchronous):**
```go
ch := make(chan int)  // Buffer size = 0

// Sender blocks until receiver is ready
ch <- 42  // BLOCKS until someone receives

// Receiver blocks until sender is ready
value := <-ch  // BLOCKS until someone sends
```

**How it works:**
- Sender and receiver must **meet** (rendezvous)
- No data is stored in channel (buffer size = 0)
- Direct handoff from sender to receiver

**Buffered (asynchronous up to buffer size):**
```go
ch := make(chan int, 3)  // Buffer size = 3

ch <- 1   // Doesn't block (buffer has space)
ch <- 2   // Doesn't block
ch <- 3   // Doesn't block
ch <- 4   // BLOCKS (buffer full, no receiver yet)
```

**How it works:**
- Sender can add up to `buffer_size` items without blocking
- Once buffer is full, sender blocks until receiver takes an item
- Receiver can take items immediately if buffer has data
- If buffer is empty, receiver blocks until sender adds item

**Memory:**
- Unbuffered: ~48 bytes (just synchronization structures)
- Buffered: ~48 bytes + buffer_size * sizeof(T)

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

### Concept 4: sync.WaitGroup - Counter Mechanics

**WaitGroup is a counter for goroutines:**

**How Counters Work:**

```go
var wg sync.WaitGroup  // Counter starts at 0

wg.Add(1)   // Counter: 0 → 1
wg.Add(1)   // Counter: 1 → 2
wg.Add(1)   // Counter: 2 → 3

wg.Done()   // Counter: 3 → 2 (decrements by 1)
wg.Done()   // Counter: 2 → 1
wg.Done()   // Counter: 1 → 0

wg.Wait()   // Blocks until counter == 0, then returns
```

**Internal Structure:**

```
┌─────────────────────────────────────┐
│ sync.WaitGroup                      │
│ ┌─────────────────────────────────┐ │
│ │ state [3]uint32                 │ │
│ │   [0] = counter (active count)  │ │  ← Atomic operations
│ │   [1] = waiter count            │ │
│ │   [2] = semaphore               │ │
│ └─────────────────────────────────┘ │
└─────────────────────────────────────┘

Operations use atomic operations (lock-free, very fast):
- Add(delta): Atomically adds delta to counter
- Done(): Atomically decrements counter by 1
- Wait(): Uses semaphore to block until counter == 0
```

**Memory:** WaitGroup is ~12 bytes. Uses atomic operations (no locks, very efficient).

**Pattern in Worker Pool:**

```go
var wg sync.WaitGroup

// Before starting workers:
for i := 0; i < workers; i++ {
    wg.Add(1)  // Increment counter BEFORE starting goroutine
               // Counter: 0 → 1 → 2 → ... → workers
    
    go func() {
        defer wg.Done()  // Decrement when goroutine exits
                         // Counter: workers → workers-1 → ... → 0
        // work
    }()
}

// In another goroutine:
wg.Wait()  // Blocks until counter == 0
           // Returns when all workers have called Done()
```

**Why `defer wg.Done()`?**

`defer` ensures `wg.Done()` is called **even if goroutine panics**:

```go
go func() {
    defer wg.Done()  // Always executes, even on panic
    
    riskyOperation()  // Might panic
    // If panic occurs, defer still runs
}()
```

**Common Mistake:**

```go
// ❌ WRONG - Race condition!
for i := 0; i < 10; i++ {
    go func() {
        wg.Add(1)  // Race: main might call Wait() before this runs
        defer wg.Done()
        work()
    }()
}
wg.Wait()  // Might return too early!

// ✅ CORRECT - Increment before starting
for i := 0; i < 10; i++ {
    wg.Add(1)  // BEFORE starting goroutine
    go func() {
        defer wg.Done()
        work()
    }()
}
wg.Wait()  // Safe: all goroutines guaranteed to be started
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

### Concept 5: Select Statement - Multiplexing Channels

`select` is like `switch` for channels - it waits until ONE case can proceed:

**Syntax:**

```go
select {
case msg := <-ch1:
    // ch1 has a value → receive it
    fmt.Println("Received from ch1:", msg)
case msg := <-ch2:
    // ch2 has a value → receive it
    fmt.Println("Received from ch2:", msg)
case ch3 <- "hello":
    // ch3 has space → send "hello"
    fmt.Println("Sent to ch3")
default:
    // No case ready → execute immediately (non-blocking)
    fmt.Println("No channel ready")
}
```

**How Select Works:**

```
┌─────────────────────────────────────┐
│ select statement                    │
│ ┌─────────────────────────────────┐ │
│ │ Check all cases:                │ │
│ │   - ch1 ready?                  │ │
│ │   - ch2 ready?                  │ │
│ │   - ch3 ready?                  │ │
│ └─────────────────────────────────┘ │
│                                     │
│ If multiple ready:                  │
│   → Choose ONE randomly             │
│                                     │
│ If none ready:                      │
│   → If default exists: execute it   │
│   → If no default: BLOCK            │
└─────────────────────────────────────┘
```

**Key behaviors**:
- If **multiple** cases are ready: ONE is chosen **randomly** (pseudo-random)
- If **no** case is ready and there's no `default`: `select` **blocks** until one becomes ready
- If **no** case is ready but `default` exists: executes `default` immediately (non-blocking)

**Memory:** Select uses a small amount of stack space (~100 bytes) to track case states.

**Common Patterns:**

**Pattern 1: Timeout**
```go
select {
case result := <-ch:
    return result
case <-time.After(5 * time.Second):
    return errors.New("timeout")
}
```

**Pattern 2: Non-blocking send**
```go
select {
case ch <- value:
    // Sent successfully
default:
    // Channel full, drop value or handle otherwise
}
```

**Pattern 3: Cancellation check**
```go
for {
    select {
    case <-ctx.Done():
        return ctx.Err()  // Context cancelled
    case work := <-workQueue:
        process(work)
    }
}
```

**Pattern 4: Multiplex multiple channels**
```go
select {
case msg1 := <-ch1:
    handle(msg1)
case msg2 := <-ch2:
    handle(msg2)
case msg3 := <-ch3:
    handle(msg3)
}
```

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

### What You Learned

**Core Concepts:**

1. **Goroutines:**
   - Lightweight threads (~2KB initial stack vs 2MB for OS threads)
   - Managed by Go runtime scheduler (M:N scheduling)
   - Can create millions of goroutines efficiently
   - Use `go func() { ... }()` to launch

2. **Channels:**
   - Type-safe queues for goroutine communication
   - Unbuffered: Synchronous (sender waits for receiver)
   - Buffered: Asynchronous up to buffer size
   - `ch <- value` sends, `value := <-ch` receives
   - `close(ch)` signals no more values

3. **Context:**
   - Standard way to handle cancellation and timeouts
   - `ctx.Done()` returns channel that closes on cancellation
   - `<-ctx.Done()` blocks until cancelled
   - Propagates cancellation through function calls
   - Always use `defer cancel()` for cleanup

4. **WaitGroup:**
   - Counter for tracking goroutine completion
   - `wg.Add(1)` increments, `wg.Done()` decrements
   - `wg.Wait()` blocks until counter reaches 0
   - Always call `Add()` BEFORE starting goroutine

5. **Select Statement:**
   - Multiplexes channel operations
   - Chooses randomly if multiple cases ready
   - Blocks if no case ready (unless `default` exists)
   - Enables non-blocking operations with `default`

6. **Worker Pool Pattern:**
   - Bounded concurrency (N workers, not unlimited)
   - Jobs channel distributes work
   - Results channel collects outputs
   - Prevents resource exhaustion

**Advanced Concepts:**

7. **Closure Variable Capture:**
   - Loop variables captured by reference in closures
   - Fix: Pass loop variable as function parameter
   - Example: `go func(i int) { ... }(i)` not `go func() { ... }()`

8. **Error Propagation:**
   - Use buffered error channel (size 1)
   - Non-blocking send with `select` + `default`
   - First error cancels context, stops all workers

9. **Memory Management:**
   - Channels: O(buffer_size * element_size)
   - Goroutines: ~2KB initial stack per goroutine
   - Maps: O(vocabulary_size) for word counts
   - HTTP responses: O(response_size) per concurrent request

10. **errgroup Package:**
    - Combines WaitGroup + Context + Error handling
    - Automatic context cancellation on first error
    - Simpler code for fan-out patterns
    - Less boilerplate than manual approach

### Key Patterns You Can Reuse

**Pattern 1: Worker Pool**
```go
jobs := make(chan Job, workers)
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
```

**Pattern 2: Context Cancellation**
```go
ctx, cancel := context.WithCancel(ctx)
defer cancel()

go func() {
    select {
    case <-ctx.Done():
        return
    case work := <-workQueue:
        process(work)
    }
}()
```

**Pattern 3: Error Propagation**
```go
errCh := make(chan error, 1)
// In goroutine:
select {
case errCh <- err:
    cancel()
default:
}
```

**Pattern 4: errgroup Fan-Out**
```go
g, ctx := errgroup.WithContext(ctx)
g.SetLimit(maxConcurrency)

for _, task := range tasks {
    g.Go(func() error {
        return process(task)
    })
}

if err := g.Wait(); err != nil {
    return err
}
```

### Common Mistakes to Avoid

1. ❌ **Unbounded goroutines** → ✅ Use worker pool
2. ❌ **Forgetting `defer cancel()`** → ✅ Always defer cleanup
3. ❌ **Closing channels multiple times** → ✅ Only sender closes
4. ❌ **WaitGroup race conditions** → ✅ Call `Add()` before `go`
5. ❌ **Closure variable capture bugs** → ✅ Pass loop variable as parameter
6. ❌ **Ignoring context in blocking ops** → ✅ Always check `ctx.Done()`

### Why This Matters

Go's concurrency model is one of its biggest strengths. The patterns you learned here (worker pools, fan-out/fan-in, pipelines) are used in production systems processing billions of requests per day.

**Real-World Applications:**
- Web crawlers (Google, Bing)
- Image processing pipelines (Instagram, Pinterest)
- Log aggregation (Splunk, Datadog)
- API rate limiting
- Batch processing (banks, e-commerce)

### Next Steps

- **Project 07:** Learn about generics and advanced data structures (LRU cache)
- **Project 08:** Apply concurrency to HTTP clients with retries
- **Project 09:** Build HTTP servers with graceful shutdown
- **Explore:** `golang.org/x/sync` package for more concurrency utilities

### Quick Reference

| Concept | Syntax | Purpose |
|---------|--------|---------|
| Goroutine | `go func() { ... }()` | Launch concurrent function |
| Channel | `ch := make(chan T, size)` | Communicate between goroutines |
| Context | `ctx, cancel := context.WithCancel(ctx)` | Cancellation propagation |
| WaitGroup | `wg.Add(1); defer wg.Done(); wg.Wait()` | Coordinate goroutine completion |
| Select | `select { case <-ch: ... }` | Multiplex channel operations |
| errgroup | `g, ctx := errgroup.WithContext(ctx)` | Simplified error handling |

Go forth and parallelize! 🚀
