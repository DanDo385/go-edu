# 13-trace

**Transaction Tracing**

Trace transaction execution to see opcode-level details and gas usage.

## What You'll Learn

- Debug tracing via debug_traceTransaction
- Understanding EVM opcodes
- Gas profiling at instruction level
- Call graph analysis

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run(ctx, client, cfg)` | Trace transaction and analyze execution |

## Project Structure

```
13-trace/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/trace/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   ├── solution.reference.go # Reference solution
│   └── types.go         # Types and interfaces
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/13-trace

# Trace a transaction (requires archive node with debug API)
go run ./cmd/app/main.go https://your-archive-node 0x<tx_hash>
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
| `RPC_URL` | Yes | Ethereum RPC endpoint (archive node with debug API) |
| `TX_HASH` | Yes | Transaction hash to trace |

## Quick Copy & Paste

```bash
# Trace transaction (requires archive node)
go run ./cmd/app/main.go https://your-archive-node 0xTRANSACTION_HASH

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **debug_traceTransaction**: Full execution trace
2. **Struct Logs**: Opcode-by-opcode execution
3. **Gas Analysis**: Understanding gas consumption

## Next Steps

After completing this exercise, proceed to `geth/14-explorer`.
