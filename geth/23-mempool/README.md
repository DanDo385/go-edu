# 23-mempool: Mempool Monitoring

## Overview

Monitor the transaction mempool (txpool). Learn about pending transactions, gas price markets, and transaction lifecycle.

## Learning Objectives

- Query pending transactions
- Monitor mempool size
- Analyze gas price distribution
- Track transaction propagation
- Implement MEV strategies

## Project Structure

```
23-mempool/
├── cmd/
│   ├── app/main.go
│   └── dev/main.go
├── internal/mempool/
│   ├── exercise.go
│   ├── exercise_test.go
│   ├── solution.reference.go
│   └── solution_no_err.reference.go
└── README.md
```

## CLI Arguments

```bash
# Monitor mempool
go run ./cmd/app/main.go <RPC_URL>

# Show pending transactions
go run ./cmd/app/main.go <RPC_URL> --show-pending

# Gas price analysis
go run ./cmd/app/main.go <RPC_URL> --gas-analysis
```

## What the Dev Harness Demonstrates

1. **Mempool Stats** - Size and composition
2. **Pending Transactions** - List unconfirmed txs
3. **Gas Price Market** - Current gas prices
4. **Transaction Lifecycle** - From submission to inclusion
5. **MEV Opportunities** - Identify arbitrage

## Next Steps

Proceed to **geth/24-monitor** for comprehensive monitoring.
