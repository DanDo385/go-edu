# 08: Type-Safe Bindings (abigen)

## What Is This Project About?

Generate type-safe Go bindings from contract ABIs using abigen. This module teaches you to work with smart contracts using generated Go code instead of manual ABI encoding/decoding.

## Why Is This Important?

Manual ABI encoding/decoding is error-prone. abigen generates type-safe Go code that catches errors at compile-time and provides a natural Go API for contract interactions.

## Real-World Problems This Solves

- **Reducing bugs from manual ABI encoding errors**
- **Providing IDE autocomplete for contract methods**
- **Building maintainable codebases that handle ABI changes**

## Key Concepts You'll Learn

- **abigen tool usage and workflow**: abigen tool usage and workflow
- **Generated bindings**:  NewContract, contract.Call, contract.Transact
- **Type safety for contract parameters and return values**: Type safety for contract parameters and return values
- **Handling contract events through generated structs**: Handling contract events through generated structs

## Prerequisites

Completion of geth/06-smart-contracts and geth/07-eth-call

## Project Structure

```
geth/08-abigen/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application
│   └── dev/
│       └── main.go          # Debug harness
├── internal/
│   └── abigen/
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

See `internal/abigen/exercise.go` for implementation details and exercises.

## Additional Resources

- [Go-Ethereum Documentation](https://geth.ethereum.org/docs)
- [Ethereum JSON-RPC API](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [Ethereum Development Documentation](https://ethereum.org/en/developers/)

## Next Steps

- **geth/09-events**
- **geth/10-filters**
