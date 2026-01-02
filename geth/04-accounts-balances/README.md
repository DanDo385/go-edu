# 04-accounts-balances: Account Queries and Balance Management

## Overview

Learn how to query Ethereum account information including balances, transaction counts (nonces), and differentiate between externally owned accounts (EOAs) and contract accounts.

## Learning Objectives

- Query account balances in Wei and convert to Ether
- Retrieve account nonces (transaction counts)
- Distinguish between EOAs and smart contracts
- Handle balance queries at different block heights

## Project Structure

```
04-accounts-balances/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application with custom arguments
│   └── dev/
│       └── main.go          # Debug harness with fixed inputs
├── internal/
│   └── accountsbalances/
│       ├── exercise.go      # YOUR CODE GOES HERE
│       ├── exercise_test.go # Test cases
│       ├── solution.reference.go        # Complete solution
│       └── solution_no_err.reference.go # Error-free variant
└── README.md               # This file
```

## Quick Start

### 1. Implement the Exercise

Open `internal/accountsbalances/exercise.go` and implement the required functions.

### 2. Run Tests

```bash
go test -v ./...
```

### 3. Test with CLI

```bash
# Query account balance
go run ./cmd/app/main.go https://eth.llamarpc.com 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045

# Query at specific block
go run ./cmd/app/main.go https://eth.llamarpc.com 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045 18000000

# Query contract account
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
```

### 4. Debug with Dev Harness

```bash
go run ./cmd/dev/main.go
```

## CLI Arguments (cmd/app/main.go)

### Syntax

```bash
go run ./cmd/app/main.go <RPC_URL> <ADDRESS> [BLOCK_NUMBER]
```

### Arguments

- `RPC_URL` - Ethereum RPC endpoint (required)
- `ADDRESS` - Ethereum address to query (required)
- `BLOCK_NUMBER` - Specific block number (optional, default: latest)

### Example Commands

```bash
# Query Vitalik's address at latest block
go run ./cmd/app/main.go https://eth.llamarpc.com 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045

# Query USDC contract
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48

# Query at specific block height
go run ./cmd/app/main.go https://eth.llamarpc.com 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045 18000000

# Query zero address
go run ./cmd/app/main.go https://eth.llamarpc.com 0x0000000000000000000000000000000000000000
```

## What the Dev Harness Demonstrates

The `cmd/dev/main.go` automatically demonstrates:

1. **Balance Queries** - Gets ETH balance for multiple addresses
2. **Wei to Ether Conversion** - Shows balance in both Wei and Ether
3. **Nonce Queries** - Gets transaction count (nonce) for accounts
4. **Contract Detection** - Identifies if an address is a contract
5. **Historical Queries** - Queries balances at past block heights

## Key Concepts

### Wei and Ether

- 1 Ether = 10^18 Wei
- Balances are always stored in Wei (smallest unit)
- Use `big.Int` for precise arithmetic
- Convert to float64 for display (be aware of precision loss)

### Account Types

- **EOA (Externally Owned Account)**: Controlled by private key, no code
- **Contract Account**: Has associated bytecode, controlled by code logic

Check if account is a contract using `client.CodeAt()`

### Nonces

- Nonce tracks the number of transactions sent from an address
- Starts at 0 for new accounts
- Increments with each transaction
- Used to prevent replay attacks

### Historical Queries

You can query balances at any past block:
- `nil` = latest block
- Specific number = that block height
- Be aware some RPC providers limit historical queries

## Next Steps

After completing this exercise, proceed to:
- **geth/05-tx-nonces** - Deep dive into transaction nonce management

## Resources

- [Ethereum Accounts](https://ethereum.org/en/developers/docs/accounts/)
- [Units of Ether](https://ethereum.org/en/developers/docs/intro-to-ether/#denominations)
- [eth_getBalance RPC Method](https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_getbalance)
