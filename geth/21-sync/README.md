# 21-sync

**Sync Progress Monitoring**

Inspect sync progress to determine if your Ethereum node is fully synced.

## What You'll Learn

- eth_syncing RPC method
- Understanding sync modes (snap, full, archive)
- Progress calculation
- Detecting sync completion

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run(ctx, client, cfg)` | Query and analyze sync progress |

## Project Structure

```
21-sync/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/sync/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   ├── solution.reference.go # Reference solution
│   └── types.go         # Types and interfaces
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/21-sync

# Check sync status
go run ./cmd/app/main.go https://eth.llamarpc.com
```

### Run the Debug Harness

```bash
go run ./cmd/dev/main.go
```

### Run Tests

```bash
go test -v ./...
```

## CLI Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `RPC_URL` | Yes | Ethereum RPC endpoint URL |

## Quick Copy & Paste

```bash
# Check if node is synced
go run ./cmd/app/main.go https://eth.llamarpc.com

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **eth_syncing**: Returns false when synced, object when syncing
2. **Starting/Current/Highest Block**: Sync progress indicators
3. **Snap Sync**: Fast sync by downloading state snapshots

## Next Steps

After completing this exercise, proceed to `geth/22-peers`.
