# 05: Transaction Nonces

## What Is This Project About?

Understand transaction nonces, sequencing, and how to handle pending transactions. Nonces prevent replay attacks and ensure transactions are processed in order.

## Why Is This Important?

Nonce management is critical for building reliable transaction-sending applications. Incorrect nonce handling causes stuck transactions, wasted gas, or security vulnerabilities.

## Real-World Problems This Solves

- **Sending multiple transactions from the same account without conflicts**
- **Handling stuck transactions by replacing them with higher gas**
- **Building transaction queues that respect nonce sequencing**

## Key Concepts You'll Learn

- **Nonce as transaction sequence number (prevents replays)**: Nonce as transaction sequence number (prevents replays)
- **Pending nonce vs confirmed nonce**: Pending nonce vs confirmed nonce
- **Transaction replacement (same nonce, higher gas)**: Transaction replacement (same nonce, higher gas)
- **Gap nonces and their effects**: Gap nonces and their effects

## Prerequisites

Completion of geth/01-stack through geth/04-accounts-balances

## Project Structure

```
geth/05-tx-nonces/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application
│   └── dev/
│       └── main.go          # Debug harness
├── internal/
│   └── txnonces/
│       ├── exercise.go      # Your implementation
│       ├── exercise_test.go # Tests
│       ├── solution.reference.go      # Reference solution
│       ├── solution_no_err.reference.go # Simplified reference
│       └── types.go         # Type definitions
└── .vscode/
    └── launch.json          # Debug configurations
```

## How to Run

### Using cmd/app/main.go (CLI Arguments)

```bash
go run ./cmd/app/main.go <RPC_URL> <address>
```

### Using cmd/dev/main.go (Debug Harness)

```bash
# Run with fixed test inputs
go run ./cmd/dev/main.go

# Or use VS Code debugger (F5)
```

## How to Debug

1. Set breakpoints at `// BREAKPOINT:` comments
2. Press F5 in VS Code, select "Debug: cmd/dev (Debug Harness)"
3. Step through with F10 (Step Over) and F11 (Step Into)
4. Watch variables in the Variables panel

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with reference implementation
go test -tags=reference -v ./...
```

## Exercises

See `internal/txnonces/exercise.go` for implementation details and exercises.

## Additional Resources

- [Go-Ethereum Documentation](https://geth.ethereum.org/docs)
- [Ethereum JSON-RPC API](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [Ethereum Development Documentation](https://ethereum.org/en/developers/)

## Next Steps

- **geth/06-smart-contracts**
- **geth/26-eip1559**
