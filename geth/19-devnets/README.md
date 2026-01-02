# 19-devnets: Development Networks

## Overview

Learn to work with local development networks (ganache, hardhat, anvil) and testnets. Essential for development and testing.

## Learning Objectives

- Set up local Ethereum nodes
- Use development network features
- Fast-forward time and blocks
- Impersonate accounts
- Test with deterministic environments

## Project Structure

```
19-devnets/
├── cmd/
│   ├── app/main.go
│   └── dev/main.go
├── internal/devnets/
│   ├── exercise.go
│   ├── exercise_test.go
│   ├── solution.reference.go
│   └── solution_no_err.reference.go
└── README.md
```

## CLI Arguments

```bash
# Connect to devnet
go run ./cmd/app/main.go <DEVNET_URL>

# Examples
go run ./cmd/app/main.go http://localhost:8545  # Hardhat/Anvil
go run ./cmd/app/main.go http://localhost:7545  # Ganache
```

## What the Dev Harness Demonstrates

1. **Devnet Connection** - Connect to local networks
2. **Account Funding** - Get test ETH
3. **Time Manipulation** - Fast-forward blocks
4. **Snapshot/Revert** - Save and restore state
5. **Account Impersonation** - Test as any address

## Next Steps

Proceed to **geth/20-node** for node interaction.
