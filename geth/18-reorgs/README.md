# 18-reorgs

**Chain Reorganization Handling**

Detect and handle blockchain reorganizations.

## What You'll Learn

- What chain reorgs are
- Detecting reorgs via parent hash comparison
- Safe block confirmation depths
- Handling reorgs in indexers

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run()` | Detect and handle chain reorgs (placeholder) |

## Project Structure

```
18-reorgs/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/reorgs/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/18-reorgs

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
# Run reorg handler
go run ./cmd/app/main.go

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Chain Reorg**: When chain tip changes
2. **Canonical Chain**: The "main" chain
3. **Confirmation Depth**: Blocks to wait before considering final

## Next Steps

After completing this exercise, proceed to `geth/19-devnets`.
