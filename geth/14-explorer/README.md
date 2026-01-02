# 14-explorer: Block Explorer

## Overview

Build a simple block explorer that can query and display blocks, transactions, and addresses. Learn to correlate on-chain data.

## Learning Objectives

- Retrieve and display block data
- Parse transaction details
- Build address activity views
- Create a coherent data model

## Project Structure

```
14-explorer/
├── cmd/
│   ├── app/main.go
│   └── dev/main.go
├── internal/explorer/
│   ├── exercise.go
│   ├── exercise_test.go
│   ├── solution.reference.go
│   └── solution_no_err.reference.go
└── README.md
```

## CLI Arguments

```bash
# Explore block
go run ./cmd/app/main.go <RPC_URL> block <NUMBER>

# Explore transaction
go run ./cmd/app/main.go <RPC_URL> tx <HASH>

# Explore address
go run ./cmd/app/main.go <RPC_URL> address <ADDRESS>

# Examples
go run ./cmd/app/main.go https://eth.llamarpc.com block 18000000
go run ./cmd/app/main.go https://eth.llamarpc.com tx 0x123...abc
go run ./cmd/app/main.go https://eth.llamarpc.com address 0xd8dA...045
```

## What the Dev Harness Demonstrates

1. **Block Viewer** - Display block details
2. **Transaction List** - Show transactions in block
3. **Address History** - Transaction history for address
4. **Balance Tracking** - Historical balance changes
5. **Event Correlation** - Link events to transactions

## Next Steps

Proceed to **geth/15-receipts** for receipt handling.
