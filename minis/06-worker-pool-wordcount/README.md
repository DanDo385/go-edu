# Project 06: worker-pool-wordcount

## What You're Building

A concurrent URL fetcher that uses a worker pool pattern to fetch multiple web pages in parallel, tokenize their content, and aggregate word frequencies. This project demonstrates Go's concurrency primitives and patterns for bounded parallelism.

## Concepts Covered

- Goroutines for concurrent execution
- Channels for communication between goroutines
- Worker pool pattern (bounded concurrency)
- `context.Context` for cancellation and timeouts
- `sync.WaitGroup` for coordinating goroutine completion
- `httptest.Server` for deterministic HTTP testing
- Error propagation in concurrent code

## How to Run

```bash
# Run the program
make run P=06-worker-pool-wordcount

# Run tests
go test ./minis/06-worker-pool-wordcount/...

# Run with race detector (detects data races)
go test -race ./minis/06-worker-pool-wordcount/...
```

## Solution Explanation

### Worker Pool Pattern

Instead of spawning unlimited goroutines (one per URL), we use a fixed pool of workers. This:
- **Prevents resource exhaustion**: Limits concurrent HTTP connections
- **Improves predictability**: Bounded memory and goroutine count
- **Enables backpressure**: Queue fills up if workers are saturated

**Architecture**:
1. Main goroutine sends URLs to a `jobs` channel
2. N worker goroutines read from `jobs`, fetch URLs, tokenize, send results to `results` channel
3. Aggregator goroutine reads from `results`, merges word counts
4. Context cancellation propagates errors immediately

### Why Context?

`context.Context` provides:
- **Cancellation**: If one fetch fails, cancel all in-flight requests
- **Timeouts**: Prevent hanging on slow servers
- **Request-scoped values**: Pass auth tokens, trace IDs, etc.

## Where Go Shines

**Go vs Python:**
- Python: `asyncio` or `threading` are options
  Cons: GIL limits true parallelism; async/await has color problems
- Go: Goroutines are lightweight (2KB stack) and truly parallel

**Go vs JavaScript:**
- JS: Promises with `Promise.all()` or worker threads
  Cons: Single-threaded by default; worker threads are heavyweight
- Go: Concurrency is a first-class language feature

**Go vs Rust:**
- Rust: `tokio` async runtime is powerful
  Cons: Async Rust has a steep learning curve (Pin, Send, Sync traits)
- Go: Simpler mental model; runtime scheduler handles details

## Stretch Goals

1. **Add progress reporting**: Count completed/failed fetches
2. **Implement retry logic**: Retry failed fetches with exponential backoff
3. **Add rate limiting**: Use `time.Ticker` to limit requests per second
4. **Support streaming**: Process response bodies without buffering entirely
5. **Add result caching**: Skip re-fetching recently-seen URLs
