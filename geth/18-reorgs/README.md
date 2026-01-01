# 18: Chain Reorganizations

## What Is This Project About?

This module teaches you how to handle chain reorganizations (reorgs) in Ethereum applications. Reorgs occur when the canonical chain changes, invalidating previously confirmed blocks.

## Why Is This Important?

Handling reorgs is critical for:
- Data integrity in indexers
- Transaction confirmation safety
- Finality-aware applications
- Exchange and payment systems

## Key Concepts You'll Learn

- **What causes reorgs**: Network latency, competing blocks
- **Detecting reorgs**: Comparing block hashes
- **Recovery strategies**: Rolling back and replaying
- **Confirmation depth**: Waiting for finality

## Prerequisites

- Completion of `geth/17-indexer`

## How to Run

```bash
go run ./cmd/app/main.go https://eth.llamarpc.com
```

## Testing

```bash
go test -v ./...
```
