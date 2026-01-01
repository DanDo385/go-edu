# 07: Contract Calls (eth_call)

## What Is This Project About?

Make read-only contract calls using eth_call. This module teaches you to query contract state without sending transactions, building on the console concepts from geth/06-smart-contracts.

## Why Is This Important?

Most contract interactions are read-only queries (checking balances, getting prices, verifying permissions). eth_call is free, instant, and essential for building responsive UIs and analytics tools.

## Real-World Problems This Solves

- **Querying ERC20 token balances for wallet UIs**
- **Reading oracle prices for DeFi applications**
- **Checking user permissions before showing UI actions**

## Key Concepts You'll Learn

- **eth_call vs eth_sendTransaction (read-only vs state-changing)**: eth_call vs eth_sendTransaction (read-only vs state-changing)
- **ABI encoding of function calls**: ABI encoding of function calls
- **ABI decoding of return values**: ABI decoding of return values
- **Gas estimation for view functions**: Gas estimation for view functions

## Prerequisites

Completion of geth/01-stack through geth/06-smart-contracts (especially the console tutorial)

## Project Structure

```
geth/07-eth-call/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application
│   └── dev/
│       └── main.go          # Debug harness
├── internal/
│   └── ethcall/
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
go run ./cmd/app/main.go <RPC_URL> <contract_address>
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

See `internal/ethcall/exercise.go` for implementation details and exercises.

## Additional Resources

- [Go-Ethereum Documentation](https://geth.ethereum.org/docs)
- [Ethereum JSON-RPC API](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [Ethereum Development Documentation](https://ethereum.org/en/developers/)

## Next Steps

- **geth/08-abigen**
- **geth/09-events**
