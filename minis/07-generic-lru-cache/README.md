# 07-generic-lru-cache

**Generic LRU Cache**

Implement a thread-safe LRU cache with generics and TTL.

## What You'll Learn

- Go generics (type parameters)
- LRU eviction algorithm
- Thread-safe design with mutex
- TTL expiration

## Functions to Implement

| Function | Description |
|----------|-------------|
| `New[K, V](capacity, ttl) *Cache[K, V]` | Create new cache |
| `Get(key K) (V, bool)` | Get value (updates access time) |
| `Set(key K, val V)` | Set value (evicts if needed) |

## Project Structure

```
07-generic-lru-cache/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/genericlrucache/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   ├── cache_bench_test.go # Benchmarks
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd minis/07-generic-lru-cache

# Demo the cache
go run ./cmd/app/main.go --capacity 100 --ttl 5s
```

### Run the Debug Harness

```bash
go run ./cmd/dev/main.go
```

### Run Tests and Benchmarks

```bash
go test -v ./...
go test -bench=. -benchmem ./...
```

## CLI Arguments

| Argument | Description |
|----------|-------------|
| `--capacity` | Max items (default: 100) |
| `--ttl` | Time to live (default: 5m) |

## Quick Copy & Paste

```bash
# Run cache demo
go run ./cmd/app/main.go --capacity 10 --ttl 5s

# Run benchmarks
go test -bench=. -benchmem ./internal/genericlrucache/

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Generics**: `[K comparable, V any]`
2. **LRU Algorithm**: Doubly-linked list + map
3. **sync.Mutex**: Thread safety
4. **container/list**: Standard library linked list

## Next Steps

After completing this exercise, proceed to `minis/08-http-client-retries`.
