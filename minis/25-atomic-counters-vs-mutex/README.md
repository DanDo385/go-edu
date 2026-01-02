# 25-atomic-counters-vs-mutex

**Atomic vs Mutex**

Compare atomic operations with mutex for counters.

## What You'll Learn

- sync/atomic package
- Lock-free programming
- When atomics are faster
- Memory ordering

## Functions to Implement

| Function | Description |
|----------|-------------|
| Compare atomic vs mutex | Benchmark counters |

## Project Structure

```
25-atomic-counters-vs-mutex/
├── cmd/
│   ├── app/main.go      # CLI demonstration
│   └── dev/main.go      # Debug harness
├── internal/atomiccountersvsmutex/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/25-atomic-counters-vs-mutex

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
go test -bench=. -benchmem ./internal/atomiccountersvsmutex/

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **atomic.AddInt64**: Lock-free increment
2. **atomic.LoadInt64**: Safe read
3. **atomic.StoreInt64**: Safe write
4. **Compare-and-Swap**: Conditional update

## Next Steps

After completing this exercise, proceed to `minis/26-sync-once-singleton`.
