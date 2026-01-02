# 10-filters

**Log Filtering & Subscriptions**

Subscribe to new block headers and implement polling vs subscription patterns.

## What You'll Learn

- WebSocket subscriptions vs HTTP polling
- Block header notifications
- Implementing fallback patterns
- Resource management for subscriptions

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run(ctx, client, cfg)` | Subscribe to or poll for new headers |
| `subscribeHeads(ctx, client, cfg)` | WebSocket subscription approach |
| `pollHeads(ctx, client, cfg)` | HTTP polling fallback |

## Project Structure

```
10-filters/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/filters/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   ├── solution.reference.go # Reference solution
│   └── types.go         # Types and interfaces
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/10-filters

# Subscribe to new blocks (WebSocket required)
go run ./cmd/app/main.go wss://eth.llamarpc.com

# Poll mode with HTTP
go run ./cmd/app/main.go https://eth.llamarpc.com --poll
```

### Run the Debug Harness

```bash
go run ./cmd/dev/main.go
```

### Run Tests

```bash
go test -v ./...
```

## CLI Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `RPC_URL` | Yes | Ethereum RPC endpoint (ws/wss for subscriptions) |
| `--poll` | No | Force polling mode |
| `--max` | No | Maximum headers to collect (default: 5) |

## Quick Copy & Paste

```bash
# WebSocket subscription
go run ./cmd/app/main.go wss://eth.llamarpc.com

# HTTP polling fallback
go run ./cmd/app/main.go https://eth.llamarpc.com --poll

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **eth_subscribe**: WebSocket-based real-time notifications
2. **Polling**: HTTP-based fallback for providers without WS
3. **Graceful Degradation**: Handle subscription failures

## Next Steps

After completing this exercise, proceed to `geth/11-storage`.
