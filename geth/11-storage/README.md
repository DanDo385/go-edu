# 11-storage

**Contract Storage Access**

Read raw storage slots from Ethereum contracts, including mapping slots.

## What You'll Learn

- Storage slot layout in Solidity
- Calculating mapping slot positions
- Reading arbitrary storage slots
- Keccak256 for slot computation

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run(ctx, client, cfg)` | Read storage slots from contract |
| `slotToHash(slot)` | Convert big.Int slot to 32-byte hash |
| `mappingSlotHash(key, slot)` | Calculate slot for mapping entry |

## Project Structure

```
11-storage/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/storage/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   ├── solution.reference.go # Reference solution
│   └── types.go         # Types and interfaces
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd geth/11-storage

# Read storage slot 0 from contract
go run ./cmd/app/main.go https://eth.llamarpc.com 0xdAC17F958D2ee523a2206206994597C13D831ec7 0
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
| `CONTRACT` | Yes | Contract address |
| `SLOT` | Yes | Storage slot number |

## Quick Copy & Paste

```bash
# Read slot 0
go run ./cmd/app/main.go https://eth.llamarpc.com 0xdAC17F958D2ee523a2206206994597C13D831ec7 0

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Storage Layout**: Fixed slots for state variables
2. **Mapping Slots**: keccak256(key || slot)
3. **eth_getStorageAt**: Raw storage access

## Next Steps

After completing this exercise, proceed to `geth/12-proofs`.
