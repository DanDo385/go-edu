# Project 06: Worker Pool Wordcount

## What Is This About?

### Real-World Scenario

Imagine you're building a web crawler that needs to analyze thousands of pages. You could:

**❌ Bad approach:** Create one goroutine per page
- With 10,000 pages, you'd create 10,000 goroutines
- Each uses ~2KB memory = 20MB just for stacks
- System becomes unresponsive

**✅ Better approach:** Use a **worker pool** with bounded parallelism
- Create exactly N worker goroutines
- Feed them URLs through a channel (like a conveyor belt)
- Workers process URLs concurrently but with controlled resource usage
- Memory usage stays constant regardless of URL count

### What You'll Learn

1. **Goroutines**: Lightweight threads (~2KB vs 2MB for OS threads)
2. **Channels**: Type-safe queues for goroutine communication
3. **Worker pools**: Bounded concurrency pattern
4. **Context**: Cancellation propagation
5. **errgroup**: Simplified error handling and coordination

---

## Key Concepts

### Goroutines

**Analogy: Factory Workers**

Think of goroutines like factory workers:
- **OS Threads** = Full-time employees (expensive, limited)
- **Goroutines** = Tasks that workers can do (cheap, unlimited)
- **Go Scheduler** = Factory manager assigning tasks

```go
go doWork()  // Launch goroutine (runs concurrently)
```

**Key facts:**
- ~2KB initial stack (grows as needed)
- Created in microseconds (vs milliseconds for OS threads)
- Can have millions running simultaneously

### Channels

**Analogy: Conveyor Belts**

Channels are like conveyor belts in a factory:
- One goroutine puts items on (sender)
- Another takes items off (receiver)
- Belt has a size limit (buffer)

```go
ch := make(chan string)      // Unbuffered (blocks until receiver ready)
ch := make(chan string, 10)  // Buffered (holds 10 items)

ch <- "hello"  // Put value on belt (send)
msg := <-ch    // Take value off belt (receive)
```

**Channel states:**

| Operation | Nil | Open | Closed |
|-----------|-----|------|--------|
| Send | Blocks forever | Sends (or blocks) | **Panic** |
| Receive | Blocks forever | Receives (or blocks) | Returns zero + ok=false |
| Close | **Panic** | Succeeds | **Panic** |

**Rule:** Only the sender should close a channel.

### Context

**Analogy: Emergency Stop Button**

Context is like an emergency stop button:
- Press it once → all machines stop immediately
- All workers listening to it exit
- No need to tell each worker individually

```go
ctx, cancel := context.WithCancel(ctx)  // Create stop button
defer cancel()                           // Always cleanup

// Workers check:
select {
case <-ctx.Done():  // Button pressed?
    return           // Stop working
case work := <-jobs:
    // Process work
}
```

**What is `ctx`?**
- A value type (~48 bytes) containing cancellation state
- `ctx.Done()` returns a channel that closes when cancelled
- Reading from closed channel (`<-ctx.Done()`) returns immediately

### WaitGroup

**Analogy: Worker Counter**

WaitGroup is like a counter tracking active workers:

```go
var wg sync.WaitGroup  // Counter starts at 0

wg.Add(1)   // Increment: 0 → 1
wg.Done()   // Decrement: 1 → 0
wg.Wait()   // Block until counter == 0
```

**Pattern:**
```go
wg.Add(1)  // BEFORE starting goroutine
go func() {
    defer wg.Done()  // When goroutine exits
    // work
}()
wg.Wait()  // Wait for all to finish
```

**Common mistake:**
```go
// ❌ WRONG - Race condition!
go func() {
    wg.Add(1)  // Too late! Wait() might run first
}()

// ✅ CORRECT
wg.Add(1)  // Before go
go func() {
    defer wg.Done()
}()
```

### errgroup

**Analogy: Smart Factory Manager**

errgroup is like a smart factory manager who:
- Automatically tracks workers (no manual WaitGroup)
- Automatically stops everything on first error (no error channel)
- Automatically cleans up (no manual defer cancel)

```go
g, ctx := errgroup.WithContext(ctx)  // Create manager

g.Go(func() error {  // Launch worker (Add/Done automatic!)
    return doWork()  // Just return error - manager handles it
})

if err := g.Wait(); err != nil {  // Wait and check error
    return err
}
```

**Benefits:**
- ~40% less code
- Impossible to forget WaitGroup operations
- Automatic error propagation
- Automatic context cleanup

---

## Worker Pool Pattern

### Architecture

**Analogy: Assembly Line Factory**

```
┌─────────────────────────────────────┐
│ Main Goroutine                      │
│ (Factory Manager)                   │
│                                     │
│  URLs → [Jobs Channel] → Workers   │
│         (Conveyor Belt)             │
└──────────────────┬──────────────────┘
                   │
        ┌──────────┼──────────┐
        │          │          │
        ▼          ▼          ▼
    ┌──────┐   ┌──────┐   ┌──────┐
    │Worker│   │Worker│   │Worker│
    │  1   │   │  2   │   │  3   │
    └──┬───┘   └──┬───┘   └──┬───┘
       │          │          │
       └──────────┼──────────┘
                  │
                  ▼
         [Results Channel]
         (Finished Products)
                  │
                  ▼
         Aggregate Results
```

### Manual Approach (Understanding Under the Hood)

**Step-by-step:**

1. **Create channels** (conveyor belts)
   ```go
   jobs := make(chan string, workers)        // Jobs conveyor
   results := make(chan map[string]int, workers)  // Results conveyor
   errCh := make(chan error, 1)               // Error box (size 1)
   ```

2. **Create cancellable context** (emergency stop)
   ```go
   ctx, cancel := context.WithCancel(ctx)
   defer cancel()  // Always cleanup
   ```

3. **Create WaitGroup** (worker counter)
   ```go
   var wg sync.WaitGroup
   ```

4. **Launch workers** (hire factory workers)
   ```go
   for i := 0; i < workers; i++ {
       wg.Add(1)  // Increment counter
       go func(workerID int) {  // Pass i to avoid closure bug!
           defer wg.Done()  // Decrement when done
           // Worker loop
       }(i)
   }
   ```

5. **Send jobs** (load conveyor belt)
   ```go
   go func() {
       defer close(jobs)  // Close when done
       for _, url := range urls {
           select {
           case <-ctx.Done():
               return
           case jobs <- url:
           }
       }
   }()
   ```

6. **Close results** (when all workers done)
   ```go
   go func() {
       wg.Wait()      // Wait for workers
       close(results) // Close channel
   }()
   ```

7. **Aggregate results** (collect finished products)
   ```go
   finalCounts := make(map[string]int)
   for counts := range results {
       for word, count := range counts {
           finalCounts[word] += count
       }
   }
   ```

8. **Check errors** (non-blocking)
   ```go
   select {
   case err := <-errCh:
       return nil, err
   default:
   }
   ```

### errgroup Approach (Simplified)

**Much simpler!**

```go
g, ctx := errgroup.WithContext(ctx)  // Create manager

// Launch workers
for i := 0; i < workers; i++ {
    g.Go(func() error {  // Add/Done automatic!
        // Worker loop
        if err != nil {
            return err  // Just return - manager handles cancellation!
        }
    })
}

// Send jobs (same as manual)
go func() {
    defer close(jobs)
    for _, url := range urls {
        jobs <- url
    }
}()

// Close results (same as manual)
go func() {
    g.Wait()
    close(results)
}()

// Aggregate (same as manual)
finalCounts := make(map[string]int)
for counts := range results {
    for word, count := range counts {
        finalCounts[word] += count
    }
}

// Check errors (simplified!)
if err := g.Wait(); err != nil {
    return nil, err
}
```

**Key simplifications:**
- No error channel needed
- No manual WaitGroup management
- No manual context cleanup
- Just return errors - errgroup handles everything

---

## Understanding Syntax

### `go func() { ... }()`

**What it means:**
- `go` = Launch goroutine (runs concurrently)
- `func() { ... }` = Anonymous function (function without name)
- `()` = Call the function immediately

**Example:**
```go
go func() {
    fmt.Println("Hello")
}()  // Launches goroutine that prints "Hello"
```

### `defer`

**What it means:**
- Schedules function call to run when enclosing function returns
- Runs even if function panics
- Multiple defers execute in reverse order (LIFO)

**Example:**
```go
defer cleanup()  // cleanup() runs when function exits
```

### `select`

**What it means:**
- Like `switch` but for channels
- Waits until ONE case can proceed
- If multiple ready, chooses randomly
- If none ready, blocks (unless `default` exists)

**Example:**
```go
select {
case msg := <-ch1:  // Receive from ch1
    handle(msg)
case ch2 <- value:  // Send to ch2
    // Sent
case <-ctx.Done():  // Context cancelled
    return
default:  // No case ready
    // Do something else
}
```

### `range` over channel

**What it means:**
- Iterates over channel values
- Continues until channel closed AND empty
- Exits when channel closed

**Example:**
```go
for value := range ch {  // Receive values from channel
    process(value)
}  // Exits when ch is closed
```

### `range` over map

**What it means:**
- Iterates over key-value pairs
- Order is random (by design)
- `word` = key, `count` = value

**Example:**
```go
for word, count := range counts {  // word=key, count=value
    fmt.Printf("%s: %d\n", word, count)
}
```

### Closure Variable Capture

**Problem:**
```go
// ❌ WRONG
for i := 0; i < 3; i++ {
    go func() {
        fmt.Println(i)  // All print "3"!
    }()
}
```

**Why:** All goroutines capture the SAME variable `i` by reference. By the time they run, loop finished, `i = 3`.

**Solution:**
```go
// ✅ CORRECT
for i := 0; i < 3; i++ {
    go func(id int) {  // Parameter creates copy
        fmt.Println(id)  // Each prints 0, 1, 2
    }(i)  // Pass i as argument
}
```

**Why:** Each goroutine gets its own COPY of `i`'s value.

---

## Memory Management

### Channel Memory

- **Unbuffered**: ~48 bytes (just synchronization structures)
- **Buffered**: ~48 bytes + buffer_size * element_size
- Example: `make(chan string, 10)` ≈ 48 + 10*16 = 208 bytes

### Goroutine Memory

- **Initial stack**: ~2KB per goroutine
- **Grows as needed**: Up to 1GB default
- Example: 100 goroutines ≈ 200KB initial

### Map Memory

- **Empty map**: ~48 bytes (hash table header)
- **Per entry**: ~8 bytes (key pointer) + ~8 bytes (value) + overhead
- **Grows dynamically**: Rehashes when load factor > 6.5

### HTTP Response Memory

- **Request**: ~200 bytes
- **Response**: ~500 bytes
- **Body**: O(response_size) - entire body read into memory
- Example: 1MB response = 1MB memory per concurrent request

---

## Common Mistakes

### 1. Unbounded Goroutines

```go
// ❌ WRONG
for _, url := range urls {
    go fetch(url)  // Creates 10,000 goroutines!
}

// ✅ CORRECT
jobs := make(chan string)
for i := 0; i < workers; i++ {
    go worker(jobs)  // Only N workers
}
```

### 2. Forgetting defer cancel()

```go
// ❌ WRONG
ctx, cancel := context.WithCancel(ctx)
// Forgot defer cancel() → resource leak!

// ✅ CORRECT
ctx, cancel := context.WithCancel(ctx)
defer cancel()  // Always cleanup
```

### 3. WaitGroup Race Condition

```go
// ❌ WRONG
go func() {
    wg.Add(1)  // Race: Wait() might run first
}()

// ✅ CORRECT
wg.Add(1)  // Before go
go func() {
    defer wg.Done()
}()
```

### 4. Closing Channels Multiple Times

```go
// ❌ WRONG
close(ch)
close(ch)  // PANIC!

// ✅ CORRECT
var once sync.Once
once.Do(func() { close(ch) })
```

### 5. Closure Variable Capture

```go
// ❌ WRONG
for i := 0; i < 3; i++ {
    go func() {
        fmt.Println(i)  // All print "3"
    }()
}

// ✅ CORRECT
for i := 0; i < 3; i++ {
    go func(id int) {
        fmt.Println(id)  // Prints 0, 1, 2
    }(i)
}
```

---

## Testing

Run tests:
```bash
go test -v ./minis/06-worker-pool-wordcount/...
```

Run with race detector:
```bash
go test -race ./minis/06-worker-pool-wordcount/...
```

---

## Summary

**What you learned:**
- ✅ Goroutines are lightweight threads (~2KB vs 2MB)
- ✅ Channels provide type-safe communication
- ✅ Worker pools limit concurrency for safety
- ✅ Context enables cancellation propagation
- ✅ errgroup simplifies error handling significantly

**Key patterns:**
- Worker pool: Bounded concurrency with channels
- Fan-out/Fan-in: Distribute work, collect results
- Error propagation: First error cancels all work

**Next steps:**
- Project 07: Generics and LRU cache
- Project 08: HTTP clients with retries
- Project 09: HTTP servers with graceful shutdown

Go forth and parallelize! 🚀
