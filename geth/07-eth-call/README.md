# 07-eth-call: Read-Only Contract Calls

## Overview

Master the eth_call RPC method for reading data from smart contracts without sending transactions. This is how you query contract state.

## Learning Objectives

- Make read-only contract calls using eth_call
- Encode function calls with ABI
- Decode return values
- Query contract state efficiently

## Project Structure

```
07-eth-call/
├── cmd/
│   ├── app/main.go
│   └── dev/main.go
├── internal/ethcall/
│   ├── exercise.go
│   ├── exercise_test.go
│   ├── solution.reference.go
│   └── solution_no_err.reference.go
└── README.md
```

## CLI Arguments

```bash
# Call a contract function
go run ./cmd/app/main.go <RPC_URL> <CONTRACT> <FUNCTION> [ARGS...]

# Example: Get ERC20 balance
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b...B48 balanceOf 0xd8dA...045

# Example: Get ERC20 total supply
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b...B48 totalSupply
```

## What the Dev Harness Demonstrates

1. **ERC20 Token Queries** - balanceOf, totalSupply, decimals
2. **Manual Call Data Encoding** - Build calldata by hand
3. **Result Decoding** - Parse return values
4. **Multiple Contract Calls** - Batch multiple reads
5. **Error Handling** - Handle reverts and errors

## Next Steps

Proceed to **geth/08-abigen** to generate type-safe Go bindings.
