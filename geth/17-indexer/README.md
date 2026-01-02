# 17-indexer: Blockchain Indexing

## Overview

Build a blockchain indexer that processes blocks and events, stores them in a database, and handles chain reorganizations.

## Learning Objectives

- Process blocks sequentially
- Index events and transactions
- Handle chain reorganizations
- Implement checkpointing and restart logic
- Design efficient data models

## Project Structure

```
17-indexer/
├── cmd/
│   ├── app/main.go
│   └── dev/main.go
├── internal/indexer/
│   ├── exercise.go
│   ├── exercise_test.go
│   ├── solution.reference.go
│   └── solution_no_err.reference.go
└── README.md
```

## CLI Arguments

```bash
# Start indexer
go run ./cmd/app/main.go <RPC_URL> --from-block <N> --to-block <M>

# Example: Index recent blocks
go run ./cmd/app/main.go https://eth.llamarpc.com --from-block 18000000 --to-block 18001000

# With specific contract
go run ./cmd/app/main.go https://eth.llamarpc.com --contract 0xA0b...B48 --from-block 18000000
```

## What the Dev Harness Demonstrates

1. **Block Processing** - Sequential block handling
2. **Event Indexing** - Extract and store events
3. **Progress Tracking** - Checkpointing
4. **Reorg Handling** - Rollback on reorganization
5. **Data Storage** - Efficient persistence

## Next Steps

Proceed to **geth/18-reorgs** for advanced reorg handling.
