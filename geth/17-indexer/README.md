# 17: Block Indexer

## What Is This Project About?

This module teaches you how to build a block indexer—a service that processes blocks and stores relevant data for efficient querying. Indexers are fundamental infrastructure for blockchain applications.

## Why Is This Important?

Indexers enable:
- Fast data queries without scanning chain
- Custom data transformations
- Historical data analysis
- Backend for DApps

## Key Concepts You'll Learn

- **Block processing**: Iterating through block history
- **Data extraction**: Pulling relevant data from blocks
- **State management**: Tracking indexer progress
- **Recovery**: Handling restarts and gaps

## Prerequisites

- Completion of `geth/16-concurrency`

## How to Run

```bash
go run ./cmd/app/main.go https://eth.llamarpc.com
```

## Testing

```bash
go test -v ./...
```
