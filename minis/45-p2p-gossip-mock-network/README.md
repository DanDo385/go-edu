# 45-p2p-gossip-mock-network

**P2P Gossip Network**

Simulate a gossip protocol for message propagation.

## What You'll Learn

- Gossip protocol basics
- Message propagation
- Network simulation
- Deduplication

## Functions to Implement

| Function | Description |
|----------|-------------|
| Broadcast message | Gossip to peers |
| Handle receive | Process and forward |
| Deduplicate | Track seen messages |

## Project Structure

```
45-p2p-gossip-mock-network/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness
├── internal/p2pgossipmocknetwork/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/45-p2p-gossip-mock-network

# Simulate network with 10 nodes
go run ./cmd/app/main.go --nodes 10 --fanout 3

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Simulate gossip
go run ./cmd/app/main.go --nodes 10 --fanout 3

# Larger network
go run ./cmd/app/main.go --nodes 100 --fanout 5

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Gossip**: Random peer selection
2. **Fanout**: Peers to forward to
3. **Seen Set**: Prevent re-broadcasting
4. **Exponential Spread**: O(log n) rounds

## Next Steps

After completing this exercise, proceed to `minis/46-generics-map-reduce`.
