# minis/36-caching-reverse-proxy

## Problem

Problem: Build a caching reverse proxy to reduce backend load and improve response times

Requirements:
1. Cache GET responses in memory
2. Respect Cache-Control headers
3. Implement LRU eviction
4. Support TTL expiration
5. Thread-safe for concurrent requests
6. Provide cache statistics

Algorithm:
- Check if request method is GET
- Generate cache key from URL
- Check cache for entry
- If hit and not expired: serve from cache
- If miss: forward to backend, cache response, serve to client
- Evict LRU entry if cache is full

Time Complexity:
- Get: O(1) average case
- Set: O(1) average case
- Eviction: O(1)

Space Complexity: O(n) where n is maxSize

Why Go is well-suited:
- net/http/httputil: Built-in reverse proxy
- sync.RWMutex: Efficient read-heavy locking
- container/list: Optimized doubly-linked list
- Goroutines: Natural concurrency model

Compared to other languages:
- Python: No built-in reverse proxy, GIL limits concurrency
- Node.js: Good for proxying, but weak typing for cache entries
- Rust: Excellent performance, but more complex ownership model

## Quickstart

```bash
cd minis/36-caching-reverse-proxy
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

- `internal/cachingreverseproxy/exercise.go`: implement the TODOs here
- `internal/cachingreverseproxy/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
