# 44-mempool-in-memory

**In-Memory Mempool**

Build an in-memory transaction mempool.

## What You'll Learn

- Priority queue for transactions
- Fee-based ordering
- Mempool limits
- Transaction expiration

## Functions to Implement

| Function | Description |
|----------|-------------|
| Add transaction | Insert with priority |
| Get pending | Retrieve by fee |
| Remove confirmed | Clean up included txs |

## Project Structure

```
44-mempool-in-memory/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness
├── internal/mempoolinmemory/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/44-mempool-in-memory

# Run mempool demo
go run ./cmd/app/main.go --max-size 1000

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Run demo
go run ./cmd/app/main.go --max-size 1000

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Priority Queue**: Order by gas price
2. **Size Limits**: Evict low-fee txs
3. **Nonce Ordering**: Per-account sequence
4. **TTL**: Expire old transactions

## Next Steps

After completing this exercise, proceed to `minis/45-p2p-gossip-mock-network`.
