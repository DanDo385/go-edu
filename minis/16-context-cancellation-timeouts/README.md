# 16-context-cancellation-timeouts

**Context Cancellation**

Use context for cancellation, timeouts, and request-scoped values.

## What You'll Learn

- context.WithCancel
- context.WithTimeout
- context.WithDeadline
- Propagating cancellation

## Functions to Implement

| Function | Description |
|----------|-------------|
| Demonstrate context patterns | Cancellation and timeouts |

## Project Structure

```
16-context-cancellation-timeouts/
├── cmd/
│   ├── app/main.go      # CLI demonstration
│   └── dev/main.go      # Debug harness
├── internal/contextcancellationtimeouts/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/16-context-cancellation-timeouts

# Run with timeout
go run ./cmd/app/main.go --timeout 5s

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Run with timeout
go run ./cmd/app/main.go --timeout 5s

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **context.Background()**: Root context
2. **WithTimeout**: Auto-cancel after duration
3. **ctx.Done()**: Channel that closes on cancel
4. **ctx.Err()**: Returns Canceled or DeadlineExceeded

## Next Steps

After completing this exercise, proceed to `minis/17-file-streaming-bufio`.
