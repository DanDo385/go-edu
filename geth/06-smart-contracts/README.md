# 06-smart-contracts: Smart Contract Basics

## Overview

Introduction to interacting with smart contracts on Ethereum. Learn how to read contract bytecode, ABI basics, and make simple contract calls.

## Learning Objectives

- Understand smart contract bytecode and deployment
- Learn about Contract ABIs (Application Binary Interface)
- Detect and identify smart contracts on-chain
- Make basic contract calls

## Project Structure

```
06-smart-contracts/
├── cmd/
│   ├── app/main.go
│   └── dev/main.go
├── internal/smartcontracts/
│   ├── exercise.go
│   ├── exercise_test.go
│   ├── solution.reference.go
│   └── solution_no_err.reference.go
└── README.md
```

## CLI Arguments

```bash
# Inspect a contract
go run ./cmd/app/main.go <RPC_URL> <CONTRACT_ADDRESS>

# Example with USDC
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
```

## What the Dev Harness Demonstrates

1. Contract bytecode retrieval
2. Contract detection
3. Basic contract properties
4. ABI fundamentals

## Next Steps

Proceed to **geth/07-eth-call** for read-only contract calls.
