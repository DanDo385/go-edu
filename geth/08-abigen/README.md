# 08-abigen

**Type-Safe Contract Bindings**

Use BoundContract for type-safe contract calls with automatic ABI encoding/decoding.

## What You'll Learn

- Using go-ethereum's ABI package
- Creating BoundContract instances
- Type-safe method calls
- Comparing manual vs generated bindings

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run(ctx, backend, cfg)` | Query contract using BoundContract |
| `callString(contract, opts, method, params...)` | Type-safe string getter |
| `callUint8(contract, opts, method, params...)` | Type-safe uint8 getter |
| `callUint256(contract, opts, method, params...)` | Type-safe uint256 getter |

## Project Structure

```
08-abigen/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/abigen/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   ├── solution.reference.go # Reference solution
│   └── types.go         # Types and interfaces
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/08-abigen

# Query token with type-safe bindings
go run ./cmd/app/main.go https://eth.llamarpc.com 0xdAC17F958D2ee523a2206206994597C13D831ec7
```

### Run the Debug Harness

```bash
go run ./cmd/dev/main.go
```

### Run Tests

```bash
go test -v ./...
```

## CLI Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `RPC_URL` | Yes | Ethereum RPC endpoint URL |
| `TOKEN_ADDRESS` | Yes | ERC20 token contract address |

## Quick Copy & Paste

```bash
# Query with bindings
go run ./cmd/app/main.go https://eth.llamarpc.com 0xdAC17F958D2ee523a2206206994597C13D831ec7

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **abigen Tool**: Generate Go bindings from ABI
2. **BoundContract**: Runtime binding without code generation
3. **Type Safety**: Compile-time type checking

## Next Steps

After completing this exercise, proceed to `geth/09-events`.
