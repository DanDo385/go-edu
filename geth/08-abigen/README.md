# 08-abigen: Go Bindings Generation

## Overview

Learn to use abigen to generate type-safe Go bindings for smart contracts. This eliminates manual ABI encoding and provides a clean, type-safe interface.

## Learning Objectives

- Generate Go bindings from contract ABI
- Use generated bindings for type-safe calls
- Understand the benefits of code generation
- Work with common contract standards (ERC20, ERC721)

## Project Structure

```
08-abigen/
├── cmd/
│   ├── app/main.go
│   └── dev/main.go
├── internal/abigen/
│   ├── exercise.go
│   ├── exercise_test.go
│   ├── solution.reference.go
│   └── solution_no_err.reference.go
└── README.md
```

## CLI Arguments

```bash
# Use generated bindings
go run ./cmd/app/main.go <RPC_URL> <CONTRACT_ADDRESS>

# Example with ERC20
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48

# Example with holder address
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b...B48 --holder 0xd8dA...045
```

## What the Dev Harness Demonstrates

1. **Binding Generation** - How to use abigen
2. **Type-Safe Calls** - Using generated Go code
3. **ERC20 Interface** - Common token operations
4. **Multiple Contracts** - Working with different ABIs
5. **Comparison** - Manual vs generated approach

## Generating Bindings

```bash
# Install abigen
go install github.com/ethereum/go-ethereum/cmd/abigen@latest

# Generate from ABI
abigen --abi contract.abi --pkg mypackage --out contract.go

# Generate from Solidity
abigen --sol contract.sol --pkg mypackage --out contract.go
```

## Next Steps

Proceed to **geth/09-events** to learn about event parsing.
