# 27-sync-pool-allocator

**Object Pooling with sync.Pool**

Reduce allocations with object pooling.

## What You'll Learn

- sync.Pool for object reuse
- Reducing GC pressure
- Pool lifecycle (may be cleared)
- Buffer pooling patterns

## Functions to Implement

| Function | Description |
|----------|-------------|
| Implement object pool | Reuse expensive objects |

## Project Structure

```
27-sync-pool-allocator/
├── cmd/
│   ├── app/main.go      # CLI demonstration
│   └── dev/main.go      # Debug harness
├── internal/syncpoolallocator/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/27-sync-pool-allocator

# Run demonstration
go run ./cmd/app/main.go

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Run demo
go run ./cmd/app/main.go

# Run benchmarks
go test -bench=. -benchmem ./internal/syncpoolallocator/

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **sync.Pool**: Object cache, may be cleared
2. **Get/Put**: Acquire and return objects
3. **New Function**: Create when pool empty
4. **Reset Before Put**: Clear state before returning

## Next Steps

After completing this exercise, proceed to `minis/28-pprof-cpu-mem-benchmarks`.
