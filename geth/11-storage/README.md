# 11-storage: Storage Slot Access

## Overview

Learn to read contract storage directly using eth_getStorageAt. Understand storage layout, slot calculation, and how to access private variables.

## Learning Objectives

- Read contract storage slots directly
- Calculate storage slot positions for different data types
- Understand Solidity storage layout rules
- Access mappings and dynamic arrays in storage

## Project Structure

```
11-storage/
├── cmd/
│   ├── app/main.go
│   └── dev/main.go
├── internal/storage/
│   ├── exercise.go
│   ├── exercise_test.go
│   ├── solution.reference.go
│   └── solution_no_err.reference.go
└── README.md
```

## CLI Arguments

```bash
# Read storage slot
go run ./cmd/app/main.go <RPC_URL> <CONTRACT> <SLOT>

# Example: Read slot 0
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b...B48 0

# Example: Read mapping slot (calculated)
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b...B48 --mapping-key 0xd8dA...045
```

## What the Dev Harness Demonstrates

1. **Simple Storage** - Read basic storage slots
2. **Mapping Access** - Calculate and read mapping slots
3. **Array Access** - Read dynamic array elements
4. **Packed Storage** - Handle packed variables
5. **Storage Layout** - Understand Solidity packing rules

## Key Concepts

### Storage Slot Calculation

**Simple variables:** Sequential from slot 0

**Mappings:** `keccak256(key . slot)`

**Dynamic arrays:** 
- Length at slot N
- Element i at `keccak256(N) + i`

## Next Steps

Proceed to **geth/12-proofs** for Merkle proof verification.
