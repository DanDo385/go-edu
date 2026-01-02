# 01-stack: Ethereum Stack Connectivity

## Overview

This project demonstrates basic Ethereum connectivity by querying chain ID, network ID, and the latest block header. This is the first thing any Ethereum application should do - verify it can talk to the network.

## Learning Objectives

- Connect to an Ethereum node via RPC
- Query basic chain information (Chain ID, Network ID)
- Retrieve and parse block headers
- Understand the difference between Chain ID and Network ID

## Project Structure

```
01-stack/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application with custom arguments
│   └── dev/
│       └── main.go          # Debug harness with fixed inputs
├── internal/
│   └── stack/
│       ├── exercise.go      # YOUR CODE GOES HERE
│       ├── exercise_test.go # Test cases
│       ├── solution.reference.go        # Complete solution
│       └── solution_no_err.reference.go # Error-free variant
└── README.md               # This file
```

## Quick Start

### 1. Implement the Exercise

Open `internal/stack/exercise.go` and implement the `Run` function. The function should:
- Query the chain ID using `client.ChainID(ctx)`
- Query the network ID using `client.NetworkID(ctx)`
- Retrieve the latest block header using `client.HeaderByNumber(ctx, nil)`

### 2. Run Tests

```bash
go test -v ./...
```

### 3. Test with CLI

```bash
# Query latest block on mainnet
go run ./cmd/app/main.go https://eth.llamarpc.com

# Query specific block number
go run ./cmd/app/main.go https://eth.llamarpc.com 12345678
```

### 4. Debug with Dev Harness

```bash
go run ./cmd/dev/main.go
```

Or use VS Code debugger (F5) with breakpoints.

## CLI Arguments (cmd/app/main.go)

### Syntax

```bash
go run ./cmd/app/main.go <RPC_URL> [BLOCK_NUMBER]
```

### Arguments

- `RPC_URL` - Ethereum RPC endpoint (required)
- `BLOCK_NUMBER` - Specific block number to query (optional, default: latest)

### Example Commands

```bash
# Connect to mainnet via public RPC
go run ./cmd/app/main.go https://eth.llamarpc.com

# Connect to Infura (requires API key)
go run ./cmd/app/main.go https://mainnet.infura.io/v3/YOUR_KEY

# Connect to Sepolia testnet
go run ./cmd/app/main.go https://rpc.sepolia.org

# Query specific block
go run ./cmd/app/main.go https://eth.llamarpc.com 18000000
```

## What the Dev Harness Demonstrates

The `cmd/dev/main.go` file automatically demonstrates:

1. **Connection Establishment** - Connects to a public RPC endpoint
2. **Chain Identification** - Retrieves Chain ID and Network ID
3. **Block Header Retrieval** - Gets the latest block header
4. **Data Display** - Shows block number, hash, timestamp, and gas usage

Run it with:

```bash
go run ./cmd/dev/main.go
```

## Key Concepts

### Chain ID vs Network ID

- **Chain ID**: Used for replay protection (EIP-155). Prevents transactions from one chain being replayed on another.
- **Network ID**: Legacy identifier used for P2P networking. Usually the same as Chain ID but not always.

### Block Headers

Block headers contain:
- Block number and hash
- Parent block hash (creates the blockchain)
- Timestamp
- Gas limit and gas used
- Miner address
- State root, transactions root, receipts root (Merkle roots)

## Common Issues

### "Failed to connect"

- Check your internet connection
- Verify the RPC URL is correct and accessible
- Try a different public RPC endpoint
- Some RPC providers require API keys

### "RPC connection timeout"

- The RPC endpoint might be rate limiting you
- Try increasing the timeout in the code
- Use a different RPC provider

## Next Steps

After completing this exercise, proceed to:
- **geth/02-rpc-basics** - Learn about different RPC methods and error handling

## Resources

- [Ethereum JSON-RPC Specification](https://ethereum.github.io/execution-apis/api-documentation/)
- [EIP-155: Simple replay attack protection](https://eips.ethereum.org/EIPS/eip-155)
- [go-ethereum Documentation](https://geth.ethereum.org/docs)
