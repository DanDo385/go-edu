# 02-rpc-basics

**JSON-RPC Fundamentals**

Learn the fundamentals of Ethereum JSON-RPC by fetching blocks, transactions, and understanding the RPC layer mechanics.

## What You'll Learn

- JSON-RPC protocol fundamentals
- Fetching blocks and transactions
- Understanding retry patterns
- Error handling in RPC calls

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run(ctx, client, cfg)` | Execute fundamental RPC queries with retry logic |

## Project Structure

```
02-rpc-basics/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/rpcbasics/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   ├── solution.reference.go # Reference solution
│   └── types.go         # Types and interfaces
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/02-rpc-basics

# Query RPC basics
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
# Mainnet
go run ./cmd/app/main.go https://eth.llamarpc.com

# With retry demonstration
go run ./cmd/app/main.go https://mainnet.infura.io/v3/YOUR_KEY
```

## Key Concepts

1. **JSON-RPC 2.0**: The protocol Ethereum nodes use
2. **Retry Patterns**: Handling transient failures
3. **Request/Response**: Understanding the RPC lifecycle

## Next Steps

After completing this exercise, proceed to `geth/03-keys-addresses`.
