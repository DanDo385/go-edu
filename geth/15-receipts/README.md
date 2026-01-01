# 15: Transaction Receipts

## What Is This Project About?

This module teaches you how to work with transaction receipts—the post-execution record of what happened when a transaction was mined. Receipts contain gas usage, status, logs, and other execution metadata.

## Why Is This Important?

Transaction receipts are essential for:
- Confirming transaction success
- Extracting emitted events
- Calculating actual gas costs
- Building confirmation workflows

## Key Concepts You'll Learn

- **Receipt structure**: Status, gas used, logs, etc.
- **Transaction status**: Success vs revert detection
- **Cumulative gas**: Understanding gas accounting
- **Log extraction**: Getting events from receipts

## Prerequisites

- Completion of `geth/01-stack` through `geth/14-explorer`

## How to Run

```bash
go run ./cmd/app/main.go https://eth.llamarpc.com
```

## Testing

```bash
go test -v ./...
```
