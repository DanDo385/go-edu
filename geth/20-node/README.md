# 20-node

**Node Interaction**

Deep interaction with Ethereum node internals.

## What You'll Learn

- Admin RPC methods
- Node info queries
- Peer management
- Node configuration

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run()` | Interact with node admin APIs (placeholder) |

## Project Structure

```
20-node/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/node/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/20-node

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
# Query node info
go run ./cmd/app/main.go

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **admin_nodeInfo**: Node identity and protocols
2. **admin_peers**: Connected peer information
3. **web3_clientVersion**: Node client version

## Next Steps

After completing this exercise, proceed to `geth/21-sync`.
