# 09-events: Event Parsing and Logs

## Overview

Learn to query, filter, and parse Ethereum events (logs). Events are the primary way smart contracts communicate state changes.

## Learning Objectives

- Query contract events using eth_getLogs
- Parse event data and topics
- Understand indexed vs non-indexed parameters
- Decode events with ABI

## Project Structure

```
09-events/
├── cmd/
│   ├── app/main.go
│   └── dev/main.go
├── internal/events/
│   ├── exercise.go
│   ├── exercise_test.go
│   ├── solution.reference.go
│   └── solution_no_err.reference.go
└── README.md
```

## CLI Arguments

```bash
# Query events from contract
go run ./cmd/app/main.go <RPC_URL> <CONTRACT> --from-block <N> --to-block <M>

# Example: Get Transfer events
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b...B48 --from-block 18000000 --to-block 18000100

# Filter by address
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b...B48 --from-block 18000000 --address 0xd8dA...045
```

## What the Dev Harness Demonstrates

1. **ERC20 Transfer Events** - Parse Transfer(from, to, value)
2. **Event Filtering** - Filter by address and block range
3. **Topic Parsing** - Understand indexed parameters
4. **Batch Processing** - Handle large event sets
5. **Event Signatures** - Calculate event signature hashes

## Key Concepts

### Event Topics

- **Topic 0**: Event signature hash (e.g., keccak256("Transfer(address,address,uint256)"))
- **Topic 1-3**: Indexed parameters (up to 3)
- **Data**: Non-indexed parameters (ABI-encoded)

### Event Signature

```solidity
event Transfer(address indexed from, address indexed to, uint256 value);
```

Results in:
- Topic[0] = keccak256("Transfer(address,address,uint256)")
- Topic[1] = from address
- Topic[2] = to address
- Data = ABI-encoded value

## Next Steps

Proceed to **geth/10-filters** for advanced log filtering.
