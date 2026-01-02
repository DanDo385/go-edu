# 07-eth-call

**Read-Only Contract Calls**

Query ERC20 token metadata using manual ABI encoding/decoding with eth_call.

## What You'll Learn

- Making eth_call requests
- ABI encoding for function calls
- Decoding various return types
- Understanding ERC20 standard interface

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run(ctx, client, cfg)` | Query ERC20 metadata via eth_call |
| `selector(sig)` | Calculate function selector |
| `decodeString(data)` | Decode ABI string |
| `decodeUint8(data)` | Decode ABI uint8 |
| `decodeUint256(data)` | Decode ABI uint256 |

## Project Structure

```
07-eth-call/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/ethcall/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   ├── solution.reference.go # Reference solution
│   └── types.go         # Types and interfaces
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/07-eth-call

# Query ERC20 token
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
# Query USDT
go run ./cmd/app/main.go https://eth.llamarpc.com 0xdAC17F958D2ee523a2206206994597C13D831ec7

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **eth_call**: Simulate contract execution without transaction
2. **View Functions**: Read-only contract methods
3. **ERC20 Interface**: name(), symbol(), decimals(), totalSupply()

## Next Steps

After completing this exercise, proceed to `geth/08-abigen`.
