# 15-receipts: Receipt Handling

## Overview

Master transaction receipts including status, gas used, logs, and revert reasons. Essential for confirming transaction outcomes.

## Learning Objectives

- Retrieve transaction receipts
- Check transaction success/failure
- Extract revert reasons
- Calculate effective gas price
- Parse receipt logs

## Project Structure

```
15-receipts/
├── cmd/
│   ├── app/main.go
│   └── dev/main.go
├── internal/receipts/
│   ├── exercise.go
│   ├── exercise_test.go
│   ├── solution.reference.go
│   └── solution_no_err.reference.go
└── README.md
```

## CLI Arguments

```bash
# Get receipt
go run ./cmd/app/main.go <RPC_URL> <TX_HASH>

# Example: Successful transaction
go run ./cmd/app/main.go https://eth.llamarpc.com 0x123...abc

# Example: Failed transaction
go run ./cmd/app/main.go https://eth.llamarpc.com 0x456...def --show-revert
```

## What the Dev Harness Demonstrates

1. **Receipt Retrieval** - Get transaction receipts
2. **Status Checking** - Success vs failure
3. **Revert Reasons** - Extract error messages
4. **Gas Accounting** - Effective gas price calculation
5. **Log Extraction** - Parse events from receipts

## Next Steps

Proceed to **geth/16-concurrency** for concurrent RPC patterns.
