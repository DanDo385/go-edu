# 21-race-detection-demo

**Race Detection**

Understand data races and use Go's race detector.

## What You'll Learn

- What data races are
- Go's -race flag
- Common race patterns
- Fixing races with synchronization

## Functions to Implement

| Function | Description |
|----------|-------------|
| Demonstrate race conditions | Show detection and fixes |

## Project Structure

```
21-race-detection-demo/
├── cmd/
│   ├── app/main.go      # CLI demonstration
│   └── dev/main.go      # Debug harness
├── internal/racedetectiondemo/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/21-race-detection-demo

# Run with race detector
go run -race ./cmd/app/main.go

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Run with race detector (will show races)
go run -race ./cmd/app/main.go

# Run tests with race detector
go test -race -v ./...

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Data Race**: Concurrent access, at least one write
2. **-race Flag**: Enable race detector
3. **sync.Mutex**: Mutual exclusion
4. **Channels**: Safe concurrent communication

## Next Steps

After completing this exercise, proceed to `minis/22-worker-pool-with-backpressure`.
