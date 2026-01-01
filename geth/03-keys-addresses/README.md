# 03: Keys and Addresses

## What Is This Project About?

Learn to generate private keys, derive public keys, and create Ethereum addresses. This module teaches the cryptographic foundations of Ethereum accounts using ECDSA (secp256k1 curve).

## Why Is This Important?

Understanding key generation and address derivation is fundamental to Ethereum development. Whether you're building wallets, signing transactions, or verifying signatures, you need to understand how keys and addresses work.

## Real-World Problems This Solves

- **Generating secure private keys for new accounts**
- **Deriving addresses from public keys for wallet applications**
- **Verifying that an address corresponds to a given private key**

## Key Concepts You'll Learn

- **ECDSA cryptography on the secp256k1 curve**: ECDSA cryptography on the secp256k1 curve
- **Private key generation using crypto/rand**: Private key generation using crypto/rand
- **Public key derivation from private keys**: Public key derivation from private keys
- **Ethereum address format (Keccak256 hash of public key)**: Ethereum address format (Keccak256 hash of public key)

## Prerequisites

Completion of geth/01-stack and geth/02-rpc-basics

## Project Structure

```
geth/03-keys-addresses/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application
│   └── dev/
│       └── main.go          # Debug harness
├── internal/
│   └── keysaddresses/
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
go run ./cmd/app/main.go
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

See `internal/keysaddresses/exercise.go` for implementation details and exercises.

## Additional Resources

- [Go-Ethereum Documentation](https://geth.ethereum.org/docs)
- [Ethereum JSON-RPC API](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [Ethereum Development Documentation](https://ethereum.org/en/developers/)

## Next Steps

- **geth/04-accounts-balances**
- **geth/05-tx-nonces**
