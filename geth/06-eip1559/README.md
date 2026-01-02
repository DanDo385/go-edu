# 06-eip1559

**EIP-1559 Transactions**

Build and sign EIP-1559 dynamic fee transactions with proper fee estimation.

## What You'll Learn

- EIP-1559 fee mechanism (base fee + priority fee)
- Building dynamic fee transactions
- Fee estimation strategies
- Transaction signing with go-ethereum

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run(ctx, client, cfg)` | Build and sign EIP-1559 transactions with fee estimation |

## Project Structure

```
06-eip1559/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/eip1559/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   ├── solution.reference.go # Reference solution
│   └── types.go         # Types and interfaces
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/06-eip1559

# Estimate fees
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
# Estimate current fees
go run ./cmd/app/main.go https://eth.llamarpc.com

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Base Fee**: Protocol-determined minimum fee
2. **Priority Fee (Tip)**: Incentive for validators
3. **Max Fee**: Upper bound on total fee
4. **Fee Burning**: Base fee is burned (EIP-1559)

## Next Steps

After completing this exercise, proceed to `geth/06-smart-contracts`.
