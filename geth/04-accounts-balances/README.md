# 04-accounts-balances

**Account Balance Queries**

Query Ethereum account balances at various block heights. Understand state queries and Wei/Ether conversions.

## What You'll Learn

- Querying account balances via RPC
- Understanding historical state queries
- Wei to Ether conversions
- Handling big integers in Go

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run(ctx, client, cfg)` | Query account balances and format results |

## Project Structure

```
04-accounts-balances/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/accountsbalances/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   ├── solution.reference.go # Reference solution
│   └── types.go         # Types and interfaces
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/04-accounts-balances

# Query balance
go run ./cmd/app/main.go https://eth.llamarpc.com 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045
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
| `ADDRESS` | Yes | Ethereum address to query |
| `BLOCK` | No | Block number (default: latest) |

## Quick Copy & Paste

```bash
# Query Vitalik's balance
go run ./cmd/app/main.go https://eth.llamarpc.com 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Wei**: Smallest unit of Ether (1 ETH = 10^18 Wei)
2. **big.Int**: Go's arbitrary-precision integers
3. **Historical State**: Querying balances at specific blocks

## Next Steps

After completing this exercise, proceed to `geth/05-tx-nonces`.
