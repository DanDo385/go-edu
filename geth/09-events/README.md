# 09-events

**Event Parsing**

Query and decode ERC20 Transfer events from blockchain logs.

## What You'll Learn

- Event topics and signatures
- Log filtering by topic
- Decoding indexed vs non-indexed parameters
- ERC20 Transfer event structure

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run(ctx, client, cfg)` | Query and parse Transfer events |
| `addressTopic(addr)` | Convert address to 32-byte topic |
| `decodeTransferLog(log)` | Parse Transfer event from log |

## Project Structure

```
09-events/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/events/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   ├── solution.reference.go # Reference solution
│   └── types.go         # Types and interfaces
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/09-events

# Query Transfer events for a token
go run ./cmd/app/main.go https://eth.llamarpc.com 0xdAC17F958D2ee523a2206206994597C13D831ec7 18000000 18000100
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
| `TOKEN` | Yes | ERC20 token address |
| `FROM_BLOCK` | Yes | Start block number |
| `TO_BLOCK` | Yes | End block number |

## Quick Copy & Paste

```bash
# Query USDT transfers in block range
go run ./cmd/app/main.go https://eth.llamarpc.com 0xdAC17F958D2ee523a2206206994597C13D831ec7 18000000 18000100

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Event Signature**: keccak256("Transfer(address,address,uint256)")
2. **Topics**: Indexed parameters for efficient filtering
3. **Data Field**: Non-indexed parameters

## Next Steps

After completing this exercise, proceed to `geth/10-filters`.
