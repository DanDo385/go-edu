# 42-simple-block-struct-hashing

**Block Structure Hashing**

Build simple blockchain block structures.

## What You'll Learn

- Block header structure
- Hash chaining
- Merkle root for transactions
- Block serialization

## Functions to Implement

| Function | Description |
|----------|-------------|
| Create block | With header and transactions |
| Hash block | Deterministic hashing |
| Verify chain | Check parent hashes |

## Project Structure

```
42-simple-block-struct-hashing/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness
├── internal/simpleblockstructhashing/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/42-simple-block-struct-hashing

# Create genesis block
go run ./cmd/app/main.go genesis

# Add block to chain
go run ./cmd/app/main.go add --parent <hash> --data "tx1,tx2"

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Create genesis
go run ./cmd/app/main.go genesis

# Add blocks
go run ./cmd/app/main.go add --data "transaction data"

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Block Header**: Parent hash, merkle root, timestamp
2. **Hash Chaining**: Each block references parent
3. **Genesis Block**: First block, no parent
4. **Immutability**: Changing data changes hash

## Next Steps

After completing this exercise, proceed to `minis/43-proof-of-work-demo`.
