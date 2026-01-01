# 01: Ethereum Stack - RPC Connectivity

## What Is This Project About?

This module teaches the foundational skill of connecting to an Ethereum RPC endpoint and retrieving basic network information. You'll learn how to establish a connection, retrieve chain ID and network ID (critical for replay protection), and fetch block headers—the lightweight cryptographic commitments that define the Ethereum execution stack.

This is the first module in the geth track and establishes the pattern you'll use throughout: connect to RPC, validate inputs, make calls, handle errors, and return results safely.

## Why Is This Important?

Understanding RPC connectivity is essential because:

- **Every Ethereum application starts here**: Before you can do anything, you need to connect to a node
- **Chain ID is critical**: Prevents replay attacks across different networks
- **Headers are efficient**: Lightweight way to verify blockchain state without downloading full blocks
- **Foundation for everything**: All subsequent modules build on these concepts

## Real-World Problems This Solves

- **Multi-chain applications**: Need to identify which chain you're connected to
- **Transaction signing**: Chain ID prevents replay attacks
- **Block verification**: Headers let you verify state without full block data
- **Network diagnostics**: Quickly check if your RPC connection is working

## Key Concepts You'll Learn

- **RPC Client Connection**: How to dial an Ethereum RPC endpoint
- **Chain ID vs Network ID**: Understanding replay protection and network identification
- **Block Headers**: Lightweight block metadata with cryptographic commitments
- **Defensive Programming**: Input validation and error handling patterns
- **Context Usage**: Timeout and cancellation handling in Go

## Prerequisites

- Basic Go knowledge (functions, structs, interfaces, error handling)
- Understanding of HTTP/RPC concepts
- Familiarity with command-line tools

## Project Structure

```
01-stack/
├── cmd/
│   ├── app/          # Application entry point (CLI arguments)
│   └── dev/          # Debug harness (fixed inputs)
├── internal/
│   └── stack/        # Exercise implementation
│       ├── exercise.go
│       ├── exercise_test.go
│       ├── solution.reference.go
│       ├── solution_no_err.reference.go
│       └── types.go
└── .vscode/
    └── launch.json   # Debug configurations
```

## How to Run

### Using cmd/app/main.go (CLI Arguments)

```bash
# Get latest block info
go run ./cmd/app/main.go https://eth.llamarpc.com

# Get specific block info
go run ./cmd/app/main.go https://eth.llamarpc.com 12345
```

### Using cmd/dev/main.go (Debug Harness)

```bash
# Run with fixed test inputs
go run ./cmd/dev/main.go

# Or use VS Code debugger (F5) with "Debug: cmd/dev" configuration
```

## How to Debug

1. Set breakpoints at `// BREAKPOINT:` comments in `cmd/app/main.go` or `cmd/dev/main.go`
2. Use VS Code debugger (F5) and select appropriate configuration:
   - **"Debug: cmd/app"** - Debug with CLI arguments
   - **"Debug: cmd/dev"** - Debug with fixed inputs (recommended for learning)
   - **"Test: Run All Tests"** - Debug tests
3. Step through code using F10 (Step Over) and F11 (Step Into)
4. Watch variables in the Variables panel

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -v -run TestFunctionName ./...

# Run with reference implementation
go test -tags=reference -v ./...
```

## Exercises

Implement the `Run` function in `internal/stack/exercise.go`:

1. **Input Validation**: Handle nil context and client
2. **Retrieve Chain ID**: Get chain ID from RPC client
3. **Retrieve Network ID**: Get network ID from RPC client
4. **Fetch Block Header**: Get header for specified block (or latest)
5. **Defensive Copying**: Return independent copies of data to prevent mutations

## Additional Resources

- [go-ethereum ethclient Documentation](https://pkg.go.dev/github.com/ethereum/go-ethereum/ethclient)
- [EIP-155: Simple replay attack protection](https://eips.ethereum.org/EIPS/eip-155)
- [Ethereum JSON-RPC API](https://ethereum.org/en/developers/docs/apis/json-rpc/)
