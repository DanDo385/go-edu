# 15-receipts

**Transaction Receipts**

Fetch and analyze transaction receipts for execution results and logs.

## What You'll Learn

- Receipt structure (status, gas used, logs)
- Log decoding
- Contract deployment detection
- Gas analysis from receipts

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run(ctx, client, cfg)` | Fetch and analyze transaction receipts |

## Project Structure

```
15-receipts/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/receipts/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   ├── solution.reference.go # Reference solution
│   └── types.go         # Types and interfaces
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/15-receipts

# Get transaction receipt
go run ./cmd/app/main.go https://eth.llamarpc.com 0x<tx_hash>
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
| `TX_HASH` | Yes | Transaction hash |

## Quick Copy & Paste

```bash
# Get receipt for a transaction
go run ./cmd/app/main.go https://eth.llamarpc.com 0xTRANSACTION_HASH

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Receipt Status**: Success (1) or failure (0)
2. **Gas Used**: Actual gas consumed
3. **Logs**: Events emitted during execution
4. **Contract Address**: Set if transaction created contract

## Next Steps

After completing this exercise, proceed to `geth/16-concurrency`.
