# 16-concurrency

**Concurrent RPC Calls**

Probe multiple endpoints concurrently using a bounded worker pool.

## What You'll Learn

- Bounded worker pools with semaphores
- Concurrent RPC request patterns
- Error aggregation from multiple workers
- Context cancellation propagation

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run(ctx, prober, cfg)` | Probe endpoints concurrently with bounded workers |

## Project Structure

```
16-concurrency/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/concurrency/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   ├── solution.reference.go # Reference solution
│   └── types.go         # Types and interfaces
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/16-concurrency

# Probe multiple endpoints
go run ./cmd/app/main.go https://eth.llamarpc.com https://rpc.ankr.com/eth --workers 3
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
| `RPC_URLS` | Yes | One or more RPC endpoints |
| `--workers` | No | Number of concurrent workers (default: 5) |

## Quick Copy & Paste

```bash
# Probe multiple endpoints
go run ./cmd/app/main.go https://eth.llamarpc.com https://rpc.ankr.com/eth

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Worker Pool**: Fixed number of concurrent goroutines
2. **Semaphore**: Bounded channel for concurrency control
3. **Error Aggregation**: Collecting errors from workers

## Next Steps

After completing this exercise, proceed to `geth/17-indexer`.
