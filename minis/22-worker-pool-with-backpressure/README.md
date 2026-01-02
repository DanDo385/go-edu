# 22-worker-pool-with-backpressure

**Worker Pool with Backpressure**

Implement bounded work queues with backpressure.

## What You'll Learn

- Bounded channels for backpressure
- Work queue patterns
- Graceful degradation
- Queue depth monitoring

## Functions to Implement

| Function | Description |
|----------|-------------|
| Implement backpressure pool | Bounded work processing |

## Project Structure

```
22-worker-pool-with-backpressure/
├── cmd/
│   ├── app/main.go      # CLI demonstration
│   └── dev/main.go      # Debug harness
├── internal/workerpoolwithbackpressure/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/22-worker-pool-with-backpressure

# Run with backpressure
go run ./cmd/app/main.go --workers 5 --queue-size 100

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Run with backpressure demo
go run ./cmd/app/main.go --workers 5 --queue-size 100

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Bounded Channel**: `make(chan T, size)`
2. **Backpressure**: Block producers when queue full
3. **Graceful Degradation**: Reject or wait
4. **Queue Monitoring**: len(ch) for depth

## Next Steps

After completing this exercise, proceed to `minis/23-bounded-channel-semaphore`.
