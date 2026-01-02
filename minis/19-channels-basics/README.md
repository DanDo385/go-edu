# 19-channels-basics

**Channel Fundamentals**

Learn Go channels for goroutine communication.

## What You'll Learn

- Unbuffered vs buffered channels
- Send and receive operations
- Channel closing
- Range over channels

## Functions to Implement

| Function | Description |
|----------|-------------|
| Demonstrate channel patterns | Basic send/receive |

## Project Structure

```
19-channels-basics/
├── cmd/
│   ├── app/main.go      # CLI demonstration
│   └── dev/main.go      # Debug harness
├── internal/channelsbasics/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/19-channels-basics

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

1. **Unbuffered**: Synchronous, blocks until received
2. **Buffered**: `make(chan T, size)`, async until full
3. **close(ch)**: Signal no more values
4. **for v := range ch**: Receive until closed

## Next Steps

After completing this exercise, proceed to `minis/20-select-fanin-fanout`.
