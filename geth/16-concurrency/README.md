# 16: Concurrent RPC Calls

## What Is This Project About?

This module teaches you how to make concurrent RPC calls to improve performance when querying multiple pieces of data from Ethereum. You'll learn Go concurrency patterns applied to blockchain data fetching.

## Why Is This Important?

Concurrent calls enable:
- Faster data fetching
- Efficient batch operations
- Responsive UIs
- High-throughput indexing

## Key Concepts You'll Learn

- **Goroutines**: Concurrent execution of RPC calls
- **Channels**: Collecting results from concurrent calls
- **Error handling**: Managing errors in concurrent code
- **Rate limiting**: Respecting RPC rate limits

## Prerequisites

- Completion of `geth/15-receipts`
- Understanding of Go concurrency (goroutines, channels)

## How to Run

```bash
go run ./cmd/app/main.go https://eth.llamarpc.com
```

## Testing

```bash
go test -v ./...
```
