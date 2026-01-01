# 01: Ethereum Stack Connectivity

## What Is This Project About?

This module teaches you how to establish your first connection to an Ethereum node and verify the connection by querying fundamental network identifiers. You'll learn how to retrieve the Chain ID (used for replay protection), Network ID (legacy P2P identifier), and the latest block header—proving that your Go application can communicate with the Ethereum network.

This is the foundation of all Ethereum development. Before you can send transactions, query balances, or interact with smart contracts, you must first connect to a node. This module ensures that fundamental skill is solid.

## Why Is This Important?

Every Ethereum application starts with a connection to a node. Understanding how to:
- Dial an RPC endpoint
- Handle connection errors gracefully
- Verify connectivity by querying chain metadata
- Use context for timeout/cancellation control

...is essential for building reliable Ethereum tooling. This module establishes patterns you'll use throughout the entire course.

## Real-World Problems This Solves

- **Service health checks**: Verify your RPC connection before processing requests
- **Multi-chain support**: Detect which network you're connected to via Chain ID
- **Block monitoring**: Get the latest block header to know the current chain state
- **Connection pooling**: Understand how to properly open and close client connections

## Key Concepts You'll Learn

- **Chain ID (EIP-155)**: Unique identifier for replay protection in transaction signatures
- **Network ID**: Legacy identifier used in P2P networking
- **Block Headers**: Lightweight representation containing cryptographic commitments to block data
- **Context handling**: Go's idiomatic way to handle timeouts and cancellation
- **Defensive copying**: Protecting against unintended mutations of shared data

## Prerequisites

- Basic Go programming knowledge
- Go 1.21+ installed
- Internet connection for RPC access

## Project Structure

```
geth/01-stack/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application with argument parsing
│   └── dev/
│       └── main.go          # Debug harness with fixed inputs
├── internal/
│   └── stack/
│       ├── exercise.go      # Your implementation
│       ├── exercise_test.go # Test cases
│       ├── solution.reference.go
│       ├── solution_no_err.reference.go
│       └── types.go         # Type definitions
└── .vscode/
    └── launch.json          # Debug configurations
```

## How to Run

### Using cmd/app/main.go (CLI Arguments)

```bash
# Query latest block on mainnet
go run ./cmd/app/main.go https://eth.llamarpc.com

# Query specific block number
go run ./cmd/app/main.go https://eth.llamarpc.com 12345678
```

### Using cmd/dev/main.go (Debug Harness)

```bash
# Run with fixed test inputs
go run ./cmd/dev/main.go

# Or use VS Code debugger (F5) with "Debug: cmd/dev" configuration
```

## How to Debug

1. Set breakpoints at `// BREAKPOINT:` comments in exercise.go
2. Use VS Code debugger (F5) and select appropriate configuration:
   - "Debug: cmd/app" - Debug with CLI arguments
   - "Debug: cmd/dev" - Debug with fixed inputs (recommended for learning)
   - "Test: Run All Tests" - Debug tests
3. Step through code using F10 (Step Over) and F11 (Step Into)
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

1. Implement the `Run` function in `internal/stack/exercise.go`
2. Handle nil context and client gracefully
3. Query Chain ID, Network ID, and latest block header
4. Return results with defensive copies to prevent mutation

## Additional Resources

- [go-ethereum ethclient Package](https://pkg.go.dev/github.com/ethereum/go-ethereum/ethclient)
- [EIP-155: Simple Replay Attack Protection](https://eips.ethereum.org/EIPS/eip-155)
- [Ethereum Block Headers](https://ethereum.org/en/developers/docs/blocks/)
