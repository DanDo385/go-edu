# 06-eip1559: EIP-1559 Transaction Types

## Overview

Learn about EIP-1559 transactions including base fee, priority fee, and dynamic fee markets. Understand how to construct, estimate, and optimize gas parameters for modern Ethereum transactions.

## Learning Objectives

- Understand EIP-1559 fee structure (base fee + priority fee)
- Estimate gas fees using historical data
- Construct Type 2 (EIP-1559) transactions
- Compare legacy vs EIP-1559 transactions
- Implement fee estimation strategies

## Project Structure

```
06-eip1559/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application with custom arguments
│   └── dev/
│       └── main.go          # Debug harness with fixed inputs
├── internal/
│   └── eip1559/
│       ├── exercise.go      # YOUR CODE GOES HERE
│       ├── exercise_test.go # Test cases
│       ├── solution.reference.go        # Complete solution
│       └── solution_no_err.reference.go # Error-free variant
└── README.md               # This file
```

## Quick Start

### 1. Implement the Exercise

Open `internal/eip1559/exercise.go` and implement the EIP-1559 transaction functions.

### 2. Run Tests

```bash
go test -v ./...
```

### 3. Test with CLI

```bash
# Estimate current fees
go run ./cmd/app/main.go https://eth.llamarpc.com estimate

# Suggest fees with priority level
go run ./cmd/app/main.go https://eth.llamarpc.com suggest --priority fast

# Show fee history
go run ./cmd/app/main.go https://eth.llamarpc.com history --blocks 10
```

### 4. Debug with Dev Harness

```bash
go run ./cmd/dev/main.go
```

## CLI Arguments (cmd/app/main.go)

### Syntax

```bash
go run ./cmd/app/main.go <RPC_URL> <COMMAND> [OPTIONS]
```

### Commands

- `estimate` - Estimate current gas fees
- `suggest` - Suggest fees with priority level
- `history` - Show historical fee data

### Options

- `--priority <level>` - Priority level: low, medium, fast (default: medium)
- `--blocks <n>` - Number of blocks for history (default: 10)

### Example Commands

```bash
# Current fee estimation
go run ./cmd/app/main.go https://eth.llamarpc.com estimate

# Fast transaction suggestion
go run ./cmd/app/main.go https://eth.llamarpc.com suggest --priority fast

# Low priority (cheaper)
go run ./cmd/app/main.go https://eth.llamarpc.com suggest --priority low

# View 20 blocks of history
go run ./cmd/app/main.go https://eth.llamarpc.com history --blocks 20
```

## What the Dev Harness Demonstrates

The `cmd/dev/main.go` automatically demonstrates:

1. **Base Fee Query** - Gets current base fee from latest block
2. **Fee History** - Analyzes past blocks for fee patterns
3. **Fee Suggestion** - Calculates recommended fees for different priorities
4. **Transaction Construction** - Builds EIP-1559 transaction with proper fees
5. **Gas Estimation** - Estimates gas limit for transactions

## Key Concepts

### EIP-1559 Fee Structure

```
Max Fee = Base Fee + Priority Fee (tip)
```

- **Base Fee**: Algorithmically determined, burned (not given to miners)
- **Priority Fee (Tip)**: Goes to validators, incentivizes inclusion
- **Max Fee Per Gas**: Maximum you're willing to pay
- **Max Priority Fee Per Gas**: Maximum tip you're willing to pay

### Fee Market Dynamics

The base fee adjusts based on block utilization:
- If block is >50% full → base fee increases
- If block is <50% full → base fee decreases
- Adjusts by up to 12.5% per block

### Fee Estimation Strategy

**Conservative (Slow)**
```
maxPriorityFee = 1.5 gwei
maxFee = (2 * baseFee) + maxPriorityFee
```

**Standard (Medium)**
```
maxPriorityFee = 2 gwei
maxFee = (2 * baseFee) + maxPriorityFee
```

**Aggressive (Fast)**
```
maxPriorityFee = 3+ gwei
maxFee = (2.5 * baseFee) + maxPriorityFee
```

### Transaction Types

- **Type 0**: Legacy transactions (gasPrice)
- **Type 1**: EIP-2930 with access lists
- **Type 2**: EIP-1559 with dynamic fees (recommended)

## Common Issues

### "Max fee per gas less than block base fee"
- Your maxFeePerGas is too low
- Increase to at least current baseFee + tip

### "Transaction underpriced"
- Priority fee is too low
- Increase maxPriorityFeePerGas

### "Exceeds block gas limit"
- Gas limit too high for a single transaction
- Check your gas estimation logic

## Next Steps

After completing this exercise, proceed to:
- **geth/07-eth-call** - Make read-only contract calls

## Resources

- [EIP-1559 Specification](https://eips.ethereum.org/EIPS/eip-1559)
- [EIP-1559 Fee Market](https://ethereum.org/en/developers/docs/gas/#post-london)
- [eth_feeHistory RPC Method](https://docs.alchemy.com/reference/eth-feehistory)
- [Transaction Types](https://ethereum.org/en/developers/docs/transactions/#types-of-transactions)
