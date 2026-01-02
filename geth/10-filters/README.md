# 10-filters: Log Filtering and Queries

## Overview

Advanced log filtering techniques including block range optimization, topic filters, and efficient event queries. Learn to build production-ready event indexing systems.

## Learning Objectives

- Optimize log queries with block range strategies
- Use advanced topic filtering
- Handle pagination for large result sets
- Implement efficient event polling

## Project Structure

```
10-filters/
├── cmd/
│   ├── app/main.go
│   └── dev/main.go
├── internal/filters/
│   ├── exercise.go
│   ├── exercise_test.go
│   ├── solution.reference.go
│   └── solution_no_err.reference.go
└── README.md
```

## CLI Arguments

```bash
# Query with filters
go run ./cmd/app/main.go <RPC_URL> <CONTRACT> [OPTIONS]

# Example: Filter by sender
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b...B48 --from 0xd8dA...045 --from-block 18000000

# Example: Filter by receiver
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b...B48 --to 0xd8dA...045 --from-block 18000000

# Example: Both sender and receiver
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b...B48 --from 0xAAA --to 0xBBB
```

## What the Dev Harness Demonstrates

1. **Topic Filtering** - Filter by indexed parameters
2. **Block Range Optimization** - Efficient range queries
3. **Pagination** - Handle RPC result limits
4. **Multiple Contracts** - Query multiple addresses
5. **Real-time Monitoring** - Poll for new events

## Next Steps

Proceed to **geth/11-storage** to access contract storage directly.
