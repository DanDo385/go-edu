# 02-rpc-basics: JSON-RPC Fundamentals

## Overview

Learn the fundamentals of Ethereum's JSON-RPC interface by making various RPC calls to query blockchain data. This exercise covers the most common RPC methods and error handling patterns.

## Learning Objectives

- Make raw JSON-RPC calls to Ethereum nodes
- Understand different RPC methods and their use cases
- Handle RPC errors gracefully
- Query block data, account balances, and transaction counts

## Project Structure

```
02-rpc-basics/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application with custom arguments
│   └── dev/
│       └── main.go          # Debug harness with fixed inputs
├── internal/
│   └── rpcbasics/
│       ├── exercise.go      # YOUR CODE GOES HERE
│       ├── exercise_test.go # Test cases
│       ├── solution.reference.go        # Complete solution
│       └── solution_no_err.reference.go # Error-free variant
└── README.md               # This file
```

## Quick Start

### 1. Implement the Exercise

Open `internal/rpcbasics/exercise.go` and implement the required functions.

### 2. Run Tests

```bash
go test -v ./...
```

### 3. Test with CLI

```bash
# Query basic RPC information
go run ./cmd/app/main.go https://eth.llamarpc.com

# Query specific address balance
go run ./cmd/app/main.go https://eth.llamarpc.com 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045
```

### 4. Debug with Dev Harness

```bash
go run ./cmd/dev/main.go
```

## CLI Arguments (cmd/app/main.go)

### Syntax

```bash
go run ./cmd/app/main.go <RPC_URL> [ADDRESS]
```

### Arguments

- `RPC_URL` - Ethereum RPC endpoint (required)
- `ADDRESS` - Ethereum address to query (optional)

### Example Commands

```bash
# Basic RPC info
go run ./cmd/app/main.go https://eth.llamarpc.com

# Query Vitalik's address
go run ./cmd/app/main.go https://eth.llamarpc.com 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045

# Query a contract address
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
```

## What the Dev Harness Demonstrates

The `cmd/dev/main.go` automatically runs through:

1. **Client Version** - eth_clientVersion
2. **Network Info** - net_version, eth_chainId
3. **Block Number** - eth_blockNumber
4. **Account Balance** - eth_getBalance
5. **Transaction Count** - eth_getTransactionCount
6. **Error Handling** - How to handle various RPC errors

## Key Concepts

### Common RPC Methods

- `eth_chainId` - Get the chain ID
- `eth_blockNumber` - Get latest block number
- `eth_getBalance` - Get account balance
- `eth_getTransactionCount` - Get nonce (transaction count)
- `eth_getCode` - Check if address is a contract
- `eth_call` - Execute a read-only contract call

### Error Handling

RPC calls can fail for various reasons:
- Network connectivity issues
- Invalid parameters
- Rate limiting
- Node sync issues

Always check for errors and provide informative messages.

## Next Steps

After completing this exercise, proceed to:
- **geth/03-keys-addresses** - Learn about key generation and address derivation

## Resources

- [Ethereum JSON-RPC API](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [go-ethereum Client Documentation](https://pkg.go.dev/github.com/ethereum/go-ethereum/ethclient)
