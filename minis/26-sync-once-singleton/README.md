# 26-sync-once-singleton

**Singleton with sync.Once**

Implement thread-safe singleton initialization.

## What You'll Learn

- sync.Once for lazy initialization
- Thread-safe singletons
- Double-checked locking (don't do it!)
- Initialization patterns

## Functions to Implement

| Function | Description |
|----------|-------------|
| Implement singleton pattern | Thread-safe lazy init |

## Project Structure

```
26-sync-once-singleton/
├── cmd/
│   ├── app/main.go      # CLI demonstration
│   └── dev/main.go      # Debug harness
├── internal/synconcesingleton/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/26-sync-once-singleton

# Run demonstration
go run ./cmd/app/main.go

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Run demo
go run ./cmd/app/main.go

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **sync.Once**: Exactly once execution
2. **Lazy Init**: Initialize on first use
3. **Thread Safety**: Handles concurrent calls
4. **No Double-Check**: Go's Once is the right way

## Next Steps

After completing this exercise, proceed to `minis/27-sync-pool-allocator`.
