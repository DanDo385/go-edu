# 23-bounded-channel-semaphore

**Channel-Based Semaphore**

Use channels as semaphores for bounded concurrency.

## What You'll Learn

- Semaphore pattern with channels
- Bounded concurrency control
- Acquire/release semantics
- golang.org/x/sync/semaphore

## Functions to Implement

| Function | Description |
|----------|-------------|
| Implement channel semaphore | Bounded parallel execution |

## Project Structure

```
23-bounded-channel-semaphore/
├── cmd/
│   ├── app/main.go      # CLI demonstration
│   └── dev/main.go      # Debug harness
├── internal/boundedchannelsemaphore/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/23-bounded-channel-semaphore

# Run with semaphore
go run ./cmd/app/main.go --max-concurrent 5

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Run semaphore demo
go run ./cmd/app/main.go --max-concurrent 5

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Channel Semaphore**: `make(chan struct{}, n)`
2. **Acquire**: `sem <- struct{}{}`
3. **Release**: `<-sem`
4. **Weighted Semaphore**: x/sync/semaphore

## Next Steps

After completing this exercise, proceed to `minis/24-sync-mutex-vs-rwmutex`.
