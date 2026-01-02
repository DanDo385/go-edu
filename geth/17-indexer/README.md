# 17-indexer

**Blockchain Indexer**

Index blockchain data for efficient querying.

## What You'll Learn

- Block range processing strategies
- Data persistence patterns
- Incremental indexing
- Checkpoint management

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run()` | Index blockchain data (placeholder) |

## Project Structure

```
17-indexer/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/indexer/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/17-indexer

go run ./cmd/app/main.go
```

### Run the Debug Harness

```bash
go run ./cmd/dev/main.go
```

### Run Tests

```bash
go test -v ./...
```

## Quick Copy & Paste

```bash
# Run indexer
go run ./cmd/app/main.go

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Block Processing**: Sequential vs parallel
2. **Checkpointing**: Resume from last position
3. **Data Model**: Efficient storage for queries

## Next Steps

After completing this exercise, proceed to `geth/18-reorgs`.
