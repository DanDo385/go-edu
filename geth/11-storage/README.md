# 11: Contract Storage

## What Is This Project About?

This module teaches you how to read raw storage slots directly from smart contracts. Understanding Solidity's storage layout allows you to inspect contract state without relying on view functions, which is invaluable for debugging and analysis.

## Why Is This Important?

Direct storage access enables:
- Reading private state variables
- Debugging storage-related issues
- Analyzing contract state without ABI
- Building storage verification tools

## Key Concepts You'll Learn

- **Storage slots**: 32-byte slots addressed by position
- **Solidity storage layout**: How variables map to slots
- **Mappings and arrays**: Dynamic storage slot calculation
- **eth_getStorageAt**: RPC method for storage reads

## Prerequisites

- Completion of `geth/01-stack` through `geth/10-filters`

## How to Run

```bash
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
```

## Testing

```bash
go test -v ./...
```
