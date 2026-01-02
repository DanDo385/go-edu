# 19-devnets

**Development Networks**

Work with local development networks like Anvil and Hardhat.

## What You'll Learn

- Setting up local devnets
- Using Anvil (Foundry) and Hardhat
- Testing with funded accounts
- Time manipulation for testing

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run()` | Interact with development networks (placeholder) |

## Project Structure

```
19-devnets/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/devnets/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/19-devnets

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
# Start Anvil first: anvil
go run ./cmd/app/main.go http://127.0.0.1:8545

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Anvil**: Foundry's fast local Ethereum node
2. **Hardhat Node**: JavaScript-based dev node
3. **Funded Accounts**: Pre-funded test accounts

## Next Steps

After completing this exercise, proceed to `geth/20-node`.
