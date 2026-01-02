# 22-peers

**Peer Discovery**

Query the number of connected peers to assess node connectivity health.

## What You'll Learn

- net_peerCount RPC method
- Understanding peer connections
- Network health assessment
- P2P protocol basics

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run(ctx, client, cfg)` | Query peer count and connectivity info |

## Project Structure

```
22-peers/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/peers/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   ├── solution.reference.go # Reference solution
│   └── types.go         # Types and interfaces
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/22-peers

# Check peer count
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
# Get peer count
go run ./cmd/app/main.go https://eth.llamarpc.com

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **net_peerCount**: Number of connected peers
2. **Peer Limits**: Default max 50 peers for Geth
3. **Bootstrap Nodes**: Initial peer discovery

## Next Steps

After completing this exercise, proceed to `geth/23-mempool`.
