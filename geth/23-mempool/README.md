# 23: Mempool

## What Is This Project About?

This module teaches you about the transaction mempool (pending transaction pool). Understanding the mempool is essential for transaction monitoring, MEV, and gas optimization.

## Why Is This Important?

Mempool knowledge enables:
- Pending transaction monitoring
- Gas price analysis
- MEV research
- Transaction replacement

## Key Concepts You'll Learn

- **Transaction pool**: Pending and queued transactions
- **Transaction ordering**: Gas price priority
- **Pool limits**: Size constraints
- **txpool_content**: RPC method for mempool

## Prerequisites

- Completion of `geth/22-peers`

## How to Run

```bash
go run ./cmd/app/main.go https://eth.llamarpc.com
```

## Testing

```bash
go test -v ./...
```
