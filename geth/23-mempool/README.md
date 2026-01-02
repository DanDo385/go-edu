# 23-mempool

**Mempool Monitoring**

Inspect the mempool (transaction pool) to understand pending transactions.

## What You'll Learn

- txpool_content and txpool_status RPC methods
- Understanding pending vs queued transactions
- Nonce gap detection
- MEV and transaction ordering

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run(ctx, client, cfg)` | Query mempool status and contents |

## Project Structure

```
23-mempool/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/mempool/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   ├── solution.reference.go # Reference solution
│   └── types.go         # Types and interfaces
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/23-mempool

# Check mempool status
go run ./cmd/app/main.go https://eth.llamarpc.com
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

## Quick Copy & Paste

```bash
# Get mempool status
go run ./cmd/app/main.go https://eth.llamarpc.com

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Pending Transactions**: Ready for inclusion
2. **Queued Transactions**: Waiting due to nonce gaps
3. **Transaction Ordering**: Priority fee determines position

## Next Steps

After completing this exercise, proceed to `geth/24-monitor`.
