# 24-monitor

**Chain Monitoring**

Implement node health monitoring by checking block freshness and detecting lag.

## What You'll Learn

- Block timestamp freshness checks
- Detecting node lag
- Health check patterns
- Alerting thresholds

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run(ctx, client, cfg)` | Monitor node health and block freshness |

## Project Structure

```
24-monitor/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/monitor/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   ├── solution.reference.go # Reference solution
│   └── types.go         # Types and interfaces
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/24-monitor

# Monitor node health
go run ./cmd/app/main.go https://eth.llamarpc.com --threshold 60
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
| `RPC_URL` | Yes | Ethereum RPC endpoint URL |
| `--threshold` | No | Max block age in seconds (default: 60) |

## Quick Copy & Paste

```bash
# Monitor with 60 second threshold
go run ./cmd/app/main.go https://eth.llamarpc.com --threshold 60

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Block Freshness**: Time since last block
2. **Lag Detection**: Node falling behind chain tip
3. **Health Checks**: Automated monitoring patterns

## Next Steps

After completing this exercise, proceed to `geth/25-toolbox`.
