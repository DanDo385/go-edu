# minis/22-worker-pool-with-backpressure

## Problem

Problem: Worker pool with backpressure and rate limiting

We need to:
1. Implement a bounded worker pool that prevents unbounded queue growth
2. Provide backpressure when queue is full (reject or timeout)
3. Support non-blocking submission (fail fast when full)
4. Support timeout-based submission (wait with limit)
5. Implement token bucket rate limiting
6. Gracefully handle context cancellation

Constraints:
- Queue size must be bounded (prevent memory exhaustion)
- Workers must respect context cancellation
- Non-blocking operations must use select with default
- Rate limiter must enforce throughput limits

Time/Space Complexity:
- Submit: O(1) - channel send or immediate return
- Worker processing: O(1) per job
- Rate limiter: O(1) per token acquisition
- Space: O(queueSize + numWorkers) for channels and goroutines

Why Go is well-suited:
- Buffered channels provide built-in bounded queues with blocking
- Select statement enables non-blocking operations
- Context propagation for cancellation
- Lightweight goroutines for workers

Real-world applications:
- HTTP servers (prevent overload with 503 responses)
- Message queue consumers (acknowledge only when processed)
- Database connection pools (bounded connections)
- API rate limiting (comply with third-party limits)

## Quickstart

```bash
cd minis/22-worker-pool-with-backpressure
go test ./...
```

## CLI (cmd/app)

This project includes a small CLI wrapper around the exercise package.

### Flags

- **`-list`**: list available exported functions
- **`-fn`**: function name to run
- **`-in`**: string input (for `func(string) ...`)
- **`-n`**: int input (for `func(int) ...`)
- **`-f`**: float64 input (for `func(float64) ...`)
- **`-b`**: bool input (for `func(bool) ...`)
- **`-file`** / **`-stdin`**: input sources for `func(io.Reader) ...`

### Usage

```bash
go run ./cmd/app -h
```

### Copy/paste examples

```bash
go run ./cmd/app -list
go run ./cmd/app -fn "NewRateLimiter" -n 10
go run ./cmd/dev
```

## Debug harness (cmd/dev)

`cmd/dev` runs the same CLI with preselected example arguments, so you can just hit Run/Debug and see what the project can do.

```bash
go run ./cmd/dev
```

## Files

- `internal/workerpoolwithbackpressure/exercise.go`: implement the TODOs here
- `internal/workerpoolwithbackpressure/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
