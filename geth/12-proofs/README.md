# 12: Merkle Proofs

## What Is This Project About?

This module teaches you how Merkle proofs work in Ethereum for verifying account and storage state. You'll learn how to request proofs from nodes and verify them locally, which is fundamental for light clients and bridges.

## Why Is This Important?

Merkle proofs enable:
- Light client verification
- Cross-chain bridges
- Trustless state verification
- Storage proof systems

## Key Concepts You'll Learn

- **Account proofs**: Proving account existence and state
- **Storage proofs**: Proving storage slot values
- **Patricia Merkle Tries**: Ethereum's state tree structure
- **eth_getProof**: RPC method for proof requests

## Prerequisites

- Completion of `geth/11-storage`

## How to Run

```bash
go run ./cmd/app/main.go https://eth.llamarpc.com
```

## Testing

```bash
go test -v ./...
```
