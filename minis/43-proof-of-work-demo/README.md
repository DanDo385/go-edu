# 43-proof-of-work-demo

**Proof of Work**

Implement Hashcash-style proof of work.

## What You'll Learn

- Proof of work concept
- Difficulty adjustment
- Nonce searching
- Mining simulation

## Functions to Implement

| Function | Description |
|----------|-------------|
| Mine block | Find nonce meeting difficulty |
| Verify PoW | Check hash meets target |

## Project Structure

```
43-proof-of-work-demo/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness
├── internal/proofofworkdemo/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/43-proof-of-work-demo

# Mine with difficulty 4 (4 leading zeros)
go run ./cmd/app/main.go mine --difficulty 4 --data "block data"

# Verify a solution
go run ./cmd/app/main.go verify --hash <hash> --difficulty 4

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Mine with low difficulty (fast)
go run ./cmd/app/main.go mine --difficulty 2 --data "test"

# Mine with higher difficulty (slower)
go run ./cmd/app/main.go mine --difficulty 5 --data "test"

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Difficulty**: Required leading zero bits
2. **Nonce**: Counter to vary hash
3. **Hash Target**: hash < target
4. **Exponential Work**: Each bit doubles effort

## Next Steps

After completing this exercise, proceed to `minis/44-mempool-in-memory`.
