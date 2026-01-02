# 18-reorgs: Chain Reorganization Handling

## Overview

Master chain reorganization detection and handling. Critical for applications that need reliable transaction confirmation.

## Learning Objectives

- Detect chain reorganizations
- Implement rollback logic
- Track block confirmations
- Handle orphaned transactions
- Design reorg-resilient systems

## Project Structure

```
18-reorgs/
├── cmd/
│   ├── app/main.go
│   └── dev/main.go
├── internal/reorgs/
│   ├── exercise.go
│   ├── exercise_test.go
│   ├── solution.reference.go
│   └── solution_no_err.reference.go
└── README.md
```

## CLI Arguments

```bash
# Monitor for reorgs
go run ./cmd/app/main.go <RPC_URL> --depth <N>

# Example: Monitor with 12 block confirmation depth
go run ./cmd/app/main.go https://eth.llamarpc.com --depth 12

# Watch specific transaction
go run ./cmd/app/main.go https://eth.llamarpc.com --tx 0x123...abc
```

## What the Dev Harness Demonstrates

1. **Reorg Detection** - Identify when chain reorganizes
2. **Confirmation Tracking** - Monitor block depth
3. **Rollback Logic** - Handle invalidated data
4. **Uncle Blocks** - Understand uncle/ommer blocks
5. **Safe Finality** - Determine safe confirmation depth

## Next Steps

Proceed to **geth/19-devnets** for local development networks.
