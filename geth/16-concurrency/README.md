# 16-concurrency: Concurrent RPC Calls

## Overview

Learn to make concurrent RPC calls efficiently using goroutines, error groups, and proper synchronization. Essential for high-performance Ethereum applications.

## Learning Objectives

- Make concurrent RPC calls safely
- Use errgroup for parallel operations
- Handle rate limiting
- Implement connection pooling
- Batch RPC requests

## Project Structure

```
16-concurrency/
├── cmd/
│   ├── app/main.go
│   └── dev/main.go
├── internal/concurrency/
│   ├── exercise.go
│   ├── exercise_test.go
│   ├── solution.reference.go
│   └── solution_no_err.reference.go
└── README.md
```

## CLI Arguments

```bash
# Concurrent queries
go run ./cmd/app/main.go <RPC_URL> --workers <N> --addresses <FILE>

# Example: Query 100 balances with 10 workers
go run ./cmd/app/main.go https://eth.llamarpc.com --workers 10 --addresses addresses.txt
```

## What the Dev Harness Demonstrates

1. **Parallel Queries** - Multiple concurrent RPC calls
2. **Error Handling** - Proper error propagation
3. **Rate Limiting** - Avoid overwhelming RPC
4. **Batch Requests** - eth_batch_request
5. **Performance Comparison** - Serial vs parallel

## Next Steps

Proceed to **geth/17-indexer** for blockchain indexing.
