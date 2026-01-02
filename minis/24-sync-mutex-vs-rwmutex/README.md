# 24-sync-mutex-vs-rwmutex

**Mutex vs RWMutex**

Compare sync.Mutex and sync.RWMutex for different workloads.

## What You'll Learn

- sync.Mutex for exclusive access
- sync.RWMutex for read-heavy workloads
- When to use each
- Benchmarking lock contention

## Functions to Implement

| Function | Description |
|----------|-------------|
| Compare mutex types | Benchmark different workloads |

## Project Structure

```
24-sync-mutex-vs-rwmutex/
├── cmd/
│   ├── app/main.go      # CLI demonstration
│   └── dev/main.go      # Debug harness
├── internal/syncmutexvsrwmutex/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/24-sync-mutex-vs-rwmutex

# Run comparison
go run ./cmd/app/main.go

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Run comparison demo
go run ./cmd/app/main.go

# Run benchmarks
go test -bench=. -benchmem ./internal/syncmutexvsrwmutex/

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **sync.Mutex**: One reader OR one writer
2. **sync.RWMutex**: Many readers OR one writer
3. **Read-Heavy**: RWMutex wins
4. **Write-Heavy**: Mutex may be better

## Next Steps

After completing this exercise, proceed to `minis/25-atomic-counters-vs-mutex`.
