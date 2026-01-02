# 40-merkle-tree-basics

**Merkle Tree Implementation**

Build a Merkle tree for data integrity verification.

## What You'll Learn

- Merkle tree structure
- Building tree from leaves
- Proof generation
- Proof verification

## Functions to Implement

| Function | Description |
|----------|-------------|
| Build Merkle tree | From list of data items |
| Generate proof | Path from leaf to root |
| Verify proof | Validate inclusion |

## Project Structure

```
40-merkle-tree-basics/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness
├── internal/merkletreebasics/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/40-merkle-tree-basics

# Build tree and show root
go run ./cmd/app/main.go build "a" "b" "c" "d"

# Generate proof
go run ./cmd/app/main.go proof "a" "b" "c" "d" --item "b"

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Build tree
go run ./cmd/app/main.go build "tx1" "tx2" "tx3" "tx4"

# Get proof for item
go run ./cmd/app/main.go proof "tx1" "tx2" "tx3" "tx4" --item "tx2"

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Merkle Root**: Single hash of all data
2. **Merkle Proof**: O(log n) size proof
3. **Leaf Hash**: H(data)
4. **Internal Hash**: H(left || right)

## Next Steps

After completing this exercise, proceed to `minis/41-signed-transactions-ed25519`.
