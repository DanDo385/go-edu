# 05-tx-nonces

**Transaction Nonces**

Understand transaction nonces, the anti-replay mechanism for Ethereum transactions.

## What You'll Learn

- What nonces are and why they matter
- Querying pending vs confirmed nonces
- Nonce management strategies
- Building and signing legacy transactions

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run(ctx, client, cfg)` | Query nonces and demonstrate transaction building |

## Project Structure

```
05-tx-nonces/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/txnonces/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   ├── solution.reference.go # Reference solution
│   └── types.go         # Types and interfaces
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/05-tx-nonces

# Query nonce for address
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
| `ADDRESS` | Yes | Address to query nonce for |

## Quick Copy & Paste

```bash
# Query nonce
go run ./cmd/app/main.go https://eth.llamarpc.com 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Nonce**: Sequential counter per account
2. **Pending vs Confirmed**: Different nonce views
3. **Replay Protection**: How nonces prevent double-spending

## Next Steps

After completing this exercise, proceed to `geth/06-eip1559`.
