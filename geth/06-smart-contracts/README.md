# 06-smart-contracts

**Smart Contract Interaction**

Interact with smart contracts using manual ABI encoding/decoding.

## What You'll Learn

- Function selector calculation (first 4 bytes of keccak256)
- Manual ABI encoding for call data
- Decoding return values (strings, uint8, uint256)
- Low-level eth_call usage

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run(ctx, client, cfg)` | Query contract data using manual ABI encoding |
| `selector(sig)` | Calculate 4-byte function selector |
| `decodeString(data)` | Decode ABI-encoded string |
| `decodeUint8(data)` | Decode ABI-encoded uint8 |
| `decodeUint256(data)` | Decode ABI-encoded uint256 |

## Project Structure

```
06-smart-contracts/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/smartcontracts/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   ├── solution.reference.go # Reference solution
│   └── types.go         # Types and interfaces
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/06-smart-contracts

# Query ERC20 token info
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
| `CONTRACT` | Yes | Contract address to query |

## Quick Copy & Paste

```bash
# Query USDT token info
go run ./cmd/app/main.go https://eth.llamarpc.com 0xdAC17F958D2ee523a2206206994597C13D831ec7

# Query USDC token info
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Function Selector**: keccak256("name()")[0:4]
2. **ABI Encoding**: Standard encoding for contract calls
3. **eth_call**: Read-only contract execution

## Next Steps

After completing this exercise, proceed to `geth/07-eth-call`.
