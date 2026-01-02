# 01-stack

**Ethereum Stack Connectivity**

Prove RPC connectivity by reading the network identifiers and latest block header. This is the foundational first step for any Ethereum Go application.

## What You'll Learn

- Connecting to an Ethereum node via JSON-RPC
- Understanding Chain ID vs Network ID
- Retrieving and inspecting block headers
- Defensive programming with nil checks and error handling
- Context propagation for timeouts and cancellation

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run(ctx, client, cfg)` | Connect to node, retrieve chain/network IDs and latest header |

## Project Structure

```
01-stack/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/stack/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   ├── solution.reference.go # Reference solution
│   └── types.go         # Types and interfaces
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
# From project directory
cd geth/01-stack

# Query latest block on mainnet
go run ./cmd/app/main.go https://eth.llamarpc.com

# Query specific block number
go run ./cmd/app/main.go https://eth.llamarpc.com 12345678
```

### Run the Debug Harness

```bash
# Auto-runs with fixed inputs for debugging
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
| `BLOCK_NUMBER` | No | Specific block number (default: latest) |

## Quick Copy & Paste

```bash
# Mainnet via public RPC
go run ./cmd/app/main.go https://eth.llamarpc.com

# Sepolia testnet
go run ./cmd/app/main.go https://rpc.sepolia.org

# With Infura (requires API key)
go run ./cmd/app/main.go https://mainnet.infura.io/v3/YOUR_KEY
```

## Key Concepts

1. **Chain ID (EIP-155)**: Unique identifier for replay protection
2. **Network ID**: Legacy P2P network identifier  
3. **Block Headers**: Lightweight (~500 bytes) cryptographic commitments
4. **Defensive Copying**: Prevents mutation of shared data

## Next Steps

After completing this exercise, proceed to `geth/02-rpc-basics`.
