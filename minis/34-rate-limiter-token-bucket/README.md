# minis/34-rate-limiter-token-bucket

## Problem

Problem: Implement a production-grade rate limiter using the token bucket algorithm

Requirements:
1. Thread-safe token bucket with atomic operations
2. Per-client rate limiting with independent buckets
3. HTTP middleware integration
4. Automatic token refill based on elapsed time
5. Proper client IP extraction (handle proxies/load balancers)

Why Go is well-suited:
- sync/atomic: Lock-free atomic operations for high performance
- Goroutines: Each HTTP request runs concurrently, needs thread safety
- sync.RWMutex: Efficient read-heavy access patterns
- http.Handler: Clean middleware composition

Compared to other languages:
- Python: GIL limits true concurrency, harder to implement lock-free algorithms
- Node.js: Single-threaded, easier but can't utilize multiple cores effectively
- Rust: More control, but more complex with ownership/borrowing
- Java: Similar capabilities, but more verbose

Token Bucket Algorithm:
- Bucket has maximum capacity (allows bursts)
- Tokens refill at constant rate (sustained rate limit)
- Each request costs 1 token
- If no tokens available, request is denied (429)

Real-world usage:
- AWS API Gateway: Token bucket rate limiting
- Stripe API: ~25 req/s sustained, 100 req/s burst
- GitHub API: 5000 req/hour authenticated

## Quickstart

```bash
cd minis/34-rate-limiter-token-bucket
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
go run ./cmd/dev
```

## Debug harness (cmd/dev)

`cmd/dev` runs the same CLI with preselected example arguments, so you can just hit Run/Debug and see what the project can do.

```bash
go run ./cmd/dev
```

## Files

- `internal/ratelimitertokenbucket/exercise.go`: implement the TODOs here
- `internal/ratelimitertokenbucket/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
