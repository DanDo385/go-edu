# 04: Accounts and Balances

## What Is This Project About?

Query account balances, nonces, and code. This module teaches you to read account state from the Ethereum blockchain, distinguishing between EOAs (Externally Owned Accounts) and contract accounts.

## Why Is This Important?

Reading account state is essential for building any Ethereum application—from displaying wallet balances to verifying contract deployments to checking nonces before sending transactions.

## Real-World Problems This Solves

- **Displaying wallet balances in user interfaces**
- **Verifying contract deployments by checking bytecode**
- **Determining if an address is an EOA or contract account**

## Key Concepts You'll Learn

- **EOA vs Contract accounts (code length = 0 vs > 0)**: EOA vs Contract accounts (code length = 0 vs > 0)
- **Balance queries using eth_getBalance**: Balance queries using eth_getBalance
- **Nonce queries using eth_getTransactionCount**: Nonce queries using eth_getTransactionCount
- **Code queries using eth_getCode**: Code queries using eth_getCode

## Prerequisites

Completion of geth/01-stack through geth/03-keys-addresses

## Project Structure

```
geth/04-accounts-balances/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application
│   └── dev/
│       └── main.go          # Debug harness
├── internal/
│   └── accountsbalances/
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

See `internal/accountsbalances/exercise.go` for implementation details and exercises.

## Additional Resources

- [Go-Ethereum Documentation](https://geth.ethereum.org/docs)
- [Ethereum JSON-RPC API](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [Ethereum Development Documentation](https://ethereum.org/en/developers/)

## Next Steps

- **geth/05-tx-nonces**
- **geth/06-smart-contracts**
