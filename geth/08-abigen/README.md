# 08: abigen - Typed Contract Bindings

## What Is This Project About?

This module teaches how to use `go-ethereum`'s `BoundContract` pattern for type-safe contract calls with automatic ABI encoding/decoding. Instead of manually encoding/decoding like in `geth/07-eth-call`, you'll use the `abi` package to handle encoding/decoding automatically, resulting in cleaner, safer contract interactions.

This builds on `geth/06-smart-contracts` (console concepts) and `geth/07-eth-call` (manual encoding), showing you how libraries abstract away the low-level details while maintaining the same underlying concepts.

## Why Is This Important?

Using typed contract bindings provides:

- **Type safety**: Compile-time checks instead of runtime errors
- **Cleaner code**: No manual encoding/decoding boilerplate
- **Better IDE support**: Autocomplete and type checking
- **Less error-prone**: Library handles edge cases automatically

## Real-World Problems This Solves

- **Production applications**: Typed bindings are standard in production code
- **Large codebases**: Easier to maintain than manual encoding
- **Team collaboration**: Type safety catches errors before runtime
- **Complex contracts**: Automatic handling of nested types and arrays

## Key Concepts You'll Learn

- **BoundContract**: Type-safe wrapper around contract calls
- **ABI Parsing**: Converting JSON ABI to Go types
- **Automatic Encoding/Decoding**: Library handles ABI details
- **Call vs Transaction**: Same distinction as console (geth/06-smart-contracts)
- **Adapter Pattern**: High-level interface wrapping low-level RPC

## Prerequisites

- Completion of `geth/06-smart-contracts` (understanding contract interaction concepts)
- Completion of `geth/07-eth-call` (understanding manual ABI encoding)
- Understanding of Go interfaces and type safety

## Project Structure

```
08-abigen/
├── cmd/
│   ├── app/          # Application entry point (CLI arguments)
│   └── dev/          # Debug harness (fixed inputs)
├── internal/
│   └── abigen/       # Exercise implementation
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
# Query ERC20 contract using typed bindings
go run ./cmd/app/main.go <RPC_URL> <contract_address>
```

### Using cmd/dev/main.go (Debug Harness)

```bash
# Run with fixed test inputs
go run ./cmd/dev/main.go

# Or use VS Code debugger (F5) with "Debug: cmd/dev" configuration
```

## How to Debug

1. Set breakpoints at `// BREAKPOINT:` comments
2. Use VS Code debugger (F5) and select:
   - **"Debug: cmd/app"** - Debug with CLI arguments
   - **"Debug: cmd/dev"** - Debug with fixed inputs (recommended)
   - **"Test: Run All Tests"** - Debug tests
3. Step through BoundContract usage
4. Compare with manual encoding from geth/07-eth-call

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

Implement the `Run` function in `internal/abigen/exercise.go`:

1. **Parse ABI**: Convert JSON ABI string to `abi.ABI` type
2. **Create BoundContract**: Wrap contract address with ABI
3. **Call name()**: Use BoundContract.Call for type-safe call
4. **Call symbol()**: Same pattern with automatic decoding
5. **Call decimals()**: Handle uint8 return type
6. **Call totalSupply()**: Handle uint256 return type

## Connection to Previous Modules

- **geth/06-smart-contracts**: Console tutorial on contract interaction concepts (Call vs Transaction)
- **geth/07-eth-call**: Manual ABI encoding/decoding (what this module abstracts)
- **geth/01-stack**: RPC connection patterns

## Where This Goes Next

- **geth/09-events**: Listening to contract events and logs
- **geth/10-filters**: Advanced event filtering

## Additional Resources

- [go-ethereum abi Package](https://pkg.go.dev/github.com/ethereum/go-ethereum/accounts/abi)
- [go-ethereum bind Package](https://pkg.go.dev/github.com/ethereum/go-ethereum/accounts/abi/bind)
- [ABI Specification](https://docs.soliditylang.org/en/latest/abi-spec.html)
