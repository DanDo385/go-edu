# 14: Block Explorer

## What Is This Project About?

This module teaches you how to build block explorer functionality—querying blocks, transactions, and displaying chain data in a human-readable format. You'll combine skills from previous modules to create a comprehensive chain data viewer.

## Why Is This Important?

Block explorers are essential tools for:
- Debugging transactions
- Verifying contract deployments
- Analyzing chain activity
- User-facing blockchain UIs

## Key Concepts You'll Learn

- **Block queries**: Getting full blocks with transactions
- **Transaction details**: Decoding tx data and receipts
- **Data formatting**: Human-readable display of chain data
- **Pagination**: Handling large result sets

## Prerequisites

- Completion of `geth/01-stack` through `geth/13-trace`

## How to Run

```bash
go run ./cmd/app/main.go https://eth.llamarpc.com
```

## Testing

```bash
go test -v ./...
```
