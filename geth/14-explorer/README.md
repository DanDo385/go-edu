# 14-explorer

**Block Explorer**

Build a simple block explorer that retrieves and displays block information.

## What You'll Learn

- Fetching full blocks with transactions
- Transaction inspection
- Block navigation (parent/child)
- Formatting blockchain data for display

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run(ctx, client, cfg)` | Fetch and format block data for exploration |

## Project Structure

```
14-explorer/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/explorer/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   ├── solution.reference.go # Reference solution
│   └── types.go         # Types and interfaces
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/14-explorer

# Explore a block
go run ./cmd/app/main.go https://eth.llamarpc.com 18000000
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
| `BLOCK` | No | Block number or "latest" (default: latest) |

## Quick Copy & Paste

```bash
# Explore block 18000000
go run ./cmd/app/main.go https://eth.llamarpc.com 18000000

# Explore latest block
go run ./cmd/app/main.go https://eth.llamarpc.com latest

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Block Structure**: Header + transactions
2. **Transaction Types**: Legacy, EIP-2930, EIP-1559
3. **Block Navigation**: Using parent hash

## Next Steps

After completing this exercise, proceed to `geth/15-receipts`.
