# geth-16-concurrency

**Goal:** fetch multiple resources concurrently with a worker pool while respecting RPC limits.

## Big Picture

RPCs can be slow; fan-out speeds things up, but you need to avoid rate limits and handle cancellation. Worker pools with contexts let you balance throughput and safety. You’ll reuse this pattern in indexers (module 17) and monitors (module 24).

## Learning Objectives
- Build a simple worker pool (jobs/results channels + WaitGroup).
- Use contexts to cancel/timeout concurrent work.
- Understand rate limiting/backoff considerations and how to bound goroutines.

## Prerequisites
- Modules 01–10; Go concurrency basics (goroutines, channels, context).

## Real-World Analogy
- Multiple clerks fetching ledger pages in parallel; close the office when time is up (context cancel).

## Steps
1. Parse endpoints + worker count.
2. Spin up worker goroutines that probe endpoints via a `Prober`.
3. Feed jobs, close channel, wait for completion.
4. Aggregate successes/failures with a mutex.

## Fun Facts & Comparisons
- Too much fan-out can trigger provider rate limits; add backoff/token bucket in production.
- ethers.js/JS often uses Promise.all; Go idiom is worker pools with channels.
- Go’s goroutines are ~2KB stacks; still, unbounded creation can exhaust resources.

## Related Solidity-edu Modules
- 17-indexer — fan out block fetches.
- 24-monitor — continuous health checks.
- 06-eip1559 — reuse contexts to time-bound RPCs.

## Files
- **Exercise:** `exercise/exercise.go` - Your starting point with TODO comments guiding implementation
- **Solution:** `exercise/solution.go` - Complete implementation with detailed educational comments explaining every concept
- **Tests:** `exercise/exercise_test.go` - Test suite demonstrating patterns and verifying correctness

## How to Run Tests

To run the tests for this module:

```bash
# From the project root (go-edu/)
cd geth/16-concurrency
go test ./exercise/

# Run with verbose output to see test details
go test -v ./exercise/

# Run solution tests (build with solution tag)
go test -tags solution -v ./exercise/

# Run with race detector to catch concurrency bugs
go test -race ./exercise/

# Run specific test
go test -v ./exercise/ -run TestRun
```

## Code Structure & Patterns

### The Exercise File (`exercise/exercise.go`)

The exercise file contains TODO comments guiding you through the implementation. Each TODO represents a fundamental concept:

1. **Input Validation** - Learn defensive programming patterns (same as modules 01-stack and 06-eip1559)
2. **Worker Pool Creation** - Understand bounded concurrency with channels
3. **Context with Timeout** - Learn how to propagate cancellation and timeouts
4. **Goroutine Management** - Use WaitGroup for synchronization
5. **Concurrent Map Access** - Protect shared state with mutex
6. **Channel Communication** - Master producer-consumer patterns
7. **Error Aggregation** - Collect partial results even on timeout

### The Solution File (`exercise/solution.go`)

The solution file contains detailed educational comments explaining:
- **Why** each step is necessary (the reasoning behind the code)
- **How** concepts repeat and build on each other (pattern recognition)
- **What** fundamental principles are being demonstrated (computer science concepts)

### Key Patterns You'll Learn

#### Pattern 1: Worker Pool with Channels
```go
// Create buffered jobs channel
jobs := make(chan string, cfg.Workers)

// Launch workers that pull from channel
for i := 0; i < cfg.Workers; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        for endpoint := range jobs {
            // Process endpoint
        }
    }()
}

// Feed jobs and close when done
go func() {
    defer close(jobs)
    for _, ep := range endpoints {
        jobs <- ep
    }
}()

// Wait for all workers
wg.Wait()
```

**Why:** Bounded concurrency prevents resource exhaustion. Unlike unbounded goroutines (one per endpoint), worker pools limit concurrency to a fixed number of workers.

**Building on:** Context patterns from modules 01-stack and 06-eip1559. Now we apply contexts to concurrent operations.

**Repeats in:** Module 17 (indexer), monitors, any tool fetching data from multiple sources.

#### Pattern 2: Context Propagation in Concurrent Operations
```go
// Parent context with overall timeout
ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
defer cancel()

// Child context for each request (half the overall timeout)
reqCtx, cancelReq := context.WithTimeout(ctx, cfg.Timeout/2)
err := p.Probe(reqCtx, endpoint)
cancelReq()
```

**Why:** Prevents one slow operation from consuming the entire timeout budget. Each probe gets a bounded time slice.

**Building on:** Module 01-stack introduced context basics. Module 06-eip1559 used contexts for RPC calls. Now we nest contexts for fine-grained control.

**Repeats in:** All concurrent operations that need per-item timeouts.

#### Pattern 3: Mutex-Protected Map Access
```go
var mu sync.Mutex
results := make(map[string]time.Duration)

// In worker goroutine
mu.Lock()
results[endpoint] = latency
mu.Unlock()
```

**Why:** Maps in Go are not safe for concurrent access. Multiple goroutines writing simultaneously causes data races and panics.

**Building on:** New concurrency concept. Previous modules didn't have concurrent writes.

**Repeats in:** All concurrent aggregation scenarios (multiple goroutines updating shared state).

#### Pattern 4: WaitGroup for Synchronization
```go
var wg sync.WaitGroup

// Before launching goroutine
wg.Add(1)
go func() {
    defer wg.Done()
    // Work here
}()

// After launching all goroutines
wg.Wait()
```

**Why:** Main goroutine must wait for workers to complete before returning results. Without WaitGroup, workers would be killed mid-processing.

**Building on:** New concurrency primitive. Essential for coordinating multiple goroutines.

**Repeats in:** All patterns using worker pools or coordinated goroutines.

## Deep Dive: Worker Pool Design

### Why Worker Pools?

**Problem:** You have 1000 endpoints to probe. What are your options?

1. **Sequential:** Probe one at a time. Takes 1000 × 2s = 2000s (33 minutes). Too slow!
2. **Unbounded concurrency:** Launch 1000 goroutines. Fast, but:
   - Exhausts file descriptors (each connection = 1 FD)
   - Hits rate limits (provider throttles you)
   - Wastes memory (1000 goroutines = 2MB+ stacks)
3. **Worker pool:** Launch 4-8 workers that process jobs from a queue. Fast enough, resource-safe.

**Worker pool wins:** Balances speed and resource consumption.

### Channel Buffering

**Unbuffered channel (size 0):**
```go
jobs := make(chan string) // Producer blocks until worker receives
```

**Buffered channel (size > 0):**
```go
jobs := make(chan string, 4) // Producer can enqueue 4 jobs without blocking
```

**Why buffer size = worker count?** Allows producer to stay ahead of workers. If workers are busy, producer enqueues jobs without blocking. Keeps pipeline full.

### Per-Request Timeout Strategy

**Why half the overall timeout?**

Math example:
- Overall timeout: 5s
- Workers: 4
- Per-request timeout: 2.5s

**Worst case:** All 4 workers hit slow endpoints simultaneously:
- Each takes 2.5s
- Total: 2.5s (parallel execution)
- Well under 5s overall timeout

**If we used full 5s per-request:** One slow probe could consume entire budget, leaving no time for retries or other endpoints.

## Error Handling

### Common Concurrency Errors

**1. "fatal error: concurrent map writes"**
```
Cause: Multiple goroutines writing to map without mutex
Solution: Protect all map access with mutex (Lock before write, Unlock after)
Prevention: Always use mutex or channels for shared mutable state
```

**2. "all goroutines are asleep - deadlock!"**
```
Cause: Common causes:
  - Forgot to close channel (workers wait forever)
  - Forgot wg.Done() (main goroutine waits forever)
  - Circular waiting (A waits for B, B waits for A)
Solution: Ensure channels are closed and WaitGroup is balanced
Prevention: Use defer close() and defer wg.Done()
```

**3. Context timeout with partial results**
```
Cause: Operation exceeded timeout before completing all probes
Behavior: Returns partial results + error
Solution: This is expected! Caller should handle partial results
Prevention: Increase timeout or reduce per-request timeout for faster probes
```

### Error Wrapping Strategy

```go
// Layer 1: Probe error
err := p.Probe(ctx, endpoint)
// Error: "connection refused"

// Layer 2: Add context
res.Failures[endpoint] = fmt.Errorf("probe failed: %w", err)
// Error: "probe failed: connection refused"

// Layer 3: Overall timeout
if ctx.Err() != nil {
    return res, fmt.Errorf("run aborted: %w", ctx.Err())
}
// Error: "run aborted: context deadline exceeded"
```

This creates a traceable error chain showing exactly what happened and why.

## Testing Strategy

The test file (`exercise_test.go`) demonstrates several important patterns:

1. **Mock implementations:** Mock Prober for controllable testing
2. **Race detector:** Use `go test -race` to catch data races
3. **Timeout testing:** Verify partial results on timeout
4. **Concurrency testing:** Ensure worker pool handles concurrent loads
5. **Error case testing:** Verify error handling works correctly

**Key insight:** Because we use interfaces (Prober), we can test logic without real network calls. This makes tests fast, reliable, and deterministic.

**Example test case:**
```go
{
    name: "concurrent probes",
    cfg: Config{
        Endpoints: []string{"a", "b", "c", "d"},
        Workers:   2,
        Timeout:   5 * time.Second,
    },
    prober: &mockProber{
        delay: 100 * time.Millisecond, // Each probe takes 100ms
    },
    validate: func(t *testing.T, res *Result) {
        // All should succeed
        if len(res.Successes) != 4 {
            t.Errorf("expected 4 successes, got %d", len(res.Successes))
        }
        // Total time should be ~200ms (2 workers × 2 batches × 100ms)
        // Not 400ms (sequential)
    },
}
```

## Common Pitfalls & How to Avoid Them

### Pitfall 1: Concurrent Map Writes Without Mutex
```go
// BAD: Data race! Multiple goroutines write to map
results := make(map[string]time.Duration)
go func() { results["a"] = 1 * time.Second }()
go func() { results["b"] = 2 * time.Second }()

// GOOD: Mutex protects map
var mu sync.Mutex
results := make(map[string]time.Duration)
go func() {
    mu.Lock()
    results["a"] = 1 * time.Second
    mu.Unlock()
}()
```

**Why it's a problem:** Concurrent map access causes "fatal error: concurrent map writes" panic.

**Fix:** Always protect map access with mutex or use channels for aggregation.

### Pitfall 2: Forgetting to Close Channel
```go
// BAD: Workers block forever waiting for jobs
jobs := make(chan string)
for i := 0; i < workers; i++ {
    go func() {
        for ep := range jobs { /* process */ }
    }()
}
// Forgot close(jobs)! Workers never exit range loop.

// GOOD: Close signals workers to exit
go func() {
    defer close(jobs)
    for _, ep := range endpoints {
        jobs <- ep
    }
}()
```

**Why it's a problem:** Workers wait forever for more jobs, causing deadlock.

**Fix:** Always close channel when done sending. Use `defer close()` to ensure it happens.

### Pitfall 3: Forgetting wg.Done()
```go
// BAD: Main goroutine waits forever
var wg sync.WaitGroup
wg.Add(1)
go func() {
    // Forgot defer wg.Done()!
    doWork()
}()
wg.Wait() // Blocks forever

// GOOD: defer ensures Done is called
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    doWork()
}()
wg.Wait()
```

**Why it's a problem:** WaitGroup counter never reaches zero, causing deadlock.

**Fix:** Always `defer wg.Done()` immediately after `wg.Add()`.

### Pitfall 4: Context Leak
```go
// BAD: Context resources never freed
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
// Forgot defer cancel()!
doWork(ctx)

// GOOD: defer ensures cleanup
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel()
doWork(ctx)
```

**Why it's a problem:** Context implementation allocates resources (timers, goroutines) that leak if not canceled.

**Fix:** Always `defer cancel()` after creating context.

### Pitfall 5: Unbounded Goroutine Creation
```go
// BAD: Creates 1000 goroutines for 1000 endpoints
for _, ep := range endpoints {
    go probe(ep)
}

// GOOD: Worker pool limits to 4 concurrent operations
jobs := make(chan string, 4)
for i := 0; i < 4; i++ {
    go worker(jobs)
}
for _, ep := range endpoints {
    jobs <- ep
}
```

**Why it's a problem:** Exhausts file descriptors, memory, and may hit rate limits.

**Fix:** Use worker pool to bound concurrency.

## How Concepts Build on Each Other

This module builds on patterns from previous modules while introducing new concurrency concepts:

1. **From Module 01-stack:**
   - Context validation → Same pattern here
   - RPC call pattern → Now applied to multiple concurrent RPC calls
   - Error wrapping → Consistent usage for aggregated errors

2. **From Module 06-eip1559:**
   - Context timeouts → Now with nested timeouts (overall + per-request)
   - Defensive programming → Applied to concurrent operations
   - Config-based defaults → Workers and timeout configuration

3. **New in this module:**
   - Worker pools (channels + WaitGroup)
   - Concurrent map access (mutex protection)
   - Context propagation across goroutines
   - Per-request timeout strategy
   - Partial result aggregation

4. **Patterns that repeat throughout the course:**
   - Input validation → Every function
   - Context propagation → All operations
   - Error wrapping → All error returns
   - Worker pools → All concurrent operations (indexers, monitors)

**The progression:**
- Module 01: Single RPC call with context
- Module 06: Single transaction with context and fee calculation
- Module 16: Multiple concurrent RPC calls with worker pools
- Module 17: Concurrent block indexing (builds on this pattern)

Each module layers new concepts on top of existing patterns, building your understanding incrementally.
