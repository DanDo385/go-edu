# 05-tx-nonces: Transaction Nonce Management

## Overview

Master transaction nonce management including pending transaction tracking, nonce gaps, and strategies for handling concurrent transactions. Understanding nonces is critical for reliable transaction submission.

## Learning Objectives

- Query pending and confirmed nonces
- Handle nonce gaps and stuck transactions
- Implement strategies for concurrent transaction submission
- Understand the difference between pending and latest nonces

## Project Structure

```
05-tx-nonces/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application with custom arguments
│   └── dev/
│       └── main.go          # Debug harness with fixed inputs
├── internal/
│   └── txnonces/
│       ├── exercise.go      # YOUR CODE GOES HERE
│       ├── exercise_test.go # Test cases
│       ├── solution.reference.go        # Complete solution
│       └── solution_no_err.reference.go # Error-free variant
└── README.md               # This file
```

## Quick Start

### 1. Implement the Exercise

Open `internal/txnonces/exercise.go` and implement the nonce management functions.

### 2. Run Tests

```bash
go test -v ./...
```

### 3. Test with CLI

```bash
# Query nonce for address
go run ./cmd/app/main.go https://eth.llamarpc.com 0xYourAddress

# Compare pending vs latest nonce
go run ./cmd/app/main.go https://eth.llamarpc.com 0xYourAddress --show-pending
```

### 4. Debug with Dev Harness

```bash
go run ./cmd/dev/main.go
```

## CLI Arguments (cmd/app/main.go)

### Syntax

```bash
go run ./cmd/app/main.go <RPC_URL> <ADDRESS> [OPTIONS]
```

### Arguments

- `RPC_URL` - Ethereum RPC endpoint (required)
- `ADDRESS` - Ethereum address to query (required)

### Options

- `--show-pending` - Show both pending and latest nonces
- `--block <number>` - Query nonce at specific block

### Example Commands

```bash
# Get current nonce
go run ./cmd/app/main.go https://eth.llamarpc.com 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045

# Show pending transactions
go run ./cmd/app/main.go https://eth.llamarpc.com 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045 --show-pending

# Historical nonce
go run ./cmd/app/main.go https://eth.llamarpc.com 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045 --block 18000000
```

## What the Dev Harness Demonstrates

The `cmd/dev/main.go` automatically demonstrates:

1. **Nonce Queries** - Gets latest confirmed nonce
2. **Pending Nonce** - Shows nonce including pending transactions
3. **Nonce Gap Detection** - Identifies stuck transactions
4. **Concurrent Nonce Management** - Strategies for multiple transactions
5. **Nonce Recovery** - How to handle failed transactions

## Key Concepts

### Latest vs Pending Nonce

- **Latest (Confirmed)**: `eth_getTransactionCount` with "latest" block
  - Only counts confirmed transactions
  - Safe for most use cases
  
- **Pending**: `eth_getTransactionCount` with "pending" block
  - Includes transactions in mempool
  - Use for rapid transaction submission

### Nonce Gaps

If you send transactions with nonces 5, 6, 8:
- Transaction 8 will be stuck until 7 is submitted
- All transactions must be sequential
- No gaps allowed in the sequence

### Concurrent Transaction Strategies

**Strategy 1: Local Nonce Tracking**
```
nonce := getInitialNonce()
for each transaction {
    send with nonce
    nonce++
}
```

**Strategy 2: Pending Nonce**
```
for each transaction {
    nonce := getPendingNonce()
    send with nonce
}
```

**Strategy 3: Retry with Increment**
```
nonce := getLatestNonce()
for {
    try send with nonce
    if nonce too low: nonce++
    if success: break
}
```

## Common Issues

### "Nonce too low"
- You're using a nonce that's already been used
- Get fresh nonce and try again

### "Nonce too high"
- You've created a nonce gap
- Submit transactions for missing nonces first

### "Replacement transaction underpriced"
- Trying to replace a pending transaction without higher gas price
- Increase gas price by at least 10%

## Next Steps

After completing this exercise, proceed to:
- **geth/06-eip1559** - Learn about EIP-1559 transaction types

## Resources

- [Transaction Nonces](https://ethereum.org/en/developers/docs/transactions/#nonce)
- [eth_getTransactionCount](https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_gettransactioncount)
- [Handling Nonce Issues](https://metamask.zendesk.com/hc/en-us/articles/360015489251-How-to-Speed-Up-or-Cancel-a-Pending-Transaction)
