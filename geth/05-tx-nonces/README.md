# 05: Transaction Nonces

## What Is This Project About?

This module deep-dives into transaction nonces—the sequential counter that orders transactions from an account. You'll learn the difference between confirmed and pending nonces, common nonce-related issues, and strategies for managing nonces in concurrent applications.

## Why Is This Important?

Nonce management is one of the most common sources of bugs in Ethereum applications:
- Transactions get stuck when nonces are skipped
- Double-spending attempts are prevented by nonce checking
- Concurrent transaction sending requires careful nonce coordination

## Real-World Problems This Solves

- **Stuck transactions**: Diagnose and fix nonce gaps
- **Transaction replacement**: Use same nonce to replace pending txs
- **Concurrent sending**: Safely send multiple transactions in parallel
- **Transaction ordering**: Ensure transactions execute in desired order

## Key Concepts You'll Learn

- **Confirmed nonce**: Number of mined transactions
- **Pending nonce**: Next nonce including mempool transactions
- **Nonce gaps**: What happens when nonces are skipped
- **Transaction replacement**: Using nonce for tx replacement (EIP-1559)

## Prerequisites

- Completion of `geth/01-stack` through `geth/04-accounts-balances`

## How to Run

```bash
go run ./cmd/app/main.go https://eth.llamarpc.com 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045
```

## Testing

```bash
go test -v ./...
```
