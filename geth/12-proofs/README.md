# 12-proofs: Merkle Proofs

## Overview

Learn about Merkle proofs for verifying account and storage data. Essential for light clients and trustless verification.

## Learning Objectives

- Understand Merkle Patricia Tries
- Generate and verify account proofs
- Generate and verify storage proofs
- Use eth_getProof RPC method

## Project Structure

```
12-proofs/
├── cmd/
│   ├── app/main.go
│   └── dev/main.go
├── internal/proofs/
│   ├── exercise.go
│   ├── exercise_test.go
│   ├── solution.reference.go
│   └── solution_no_err.reference.go
└── README.md
```

## CLI Arguments

```bash
# Get account proof
go run ./cmd/app/main.go <RPC_URL> <ADDRESS> --block <N>

# Get storage proof
go run ./cmd/app/main.go <RPC_URL> <ADDRESS> --storage-key <KEY> --block <N>

# Example
go run ./cmd/app/main.go https://eth.llamarpc.com 0xd8dA...045 --block 18000000
```

## What the Dev Harness Demonstrates

1. **Account Proofs** - Prove account existence and balance
2. **Storage Proofs** - Prove storage values
3. **Proof Verification** - Verify proofs against state root
4. **Light Client Pattern** - How light clients verify data

## Next Steps

Proceed to **geth/13-trace** for transaction tracing.
