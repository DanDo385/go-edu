# 20-node: Node Interaction

## Overview

Learn to interact with Ethereum node administration APIs. Understand node status, peer management, and network diagnostics.

## Learning Objectives

- Query node information
- Check node sync status
- Monitor peer connections
- Use admin RPC methods
- Diagnose node issues

## Project Structure

```
20-node/
├── cmd/
│   ├── app/main.go
│   └── dev/main.go
├── internal/node/
│   ├── exercise.go
│   ├── exercise_test.go
│   ├── solution.reference.go
│   └── solution_no_err.reference.go
└── README.md
```

## CLI Arguments

```bash
# Node info
go run ./cmd/app/main.go <RPC_URL> info

# Sync status
go run ./cmd/app/main.go <RPC_URL> sync

# Peer info
go run ./cmd/app/main.go <RPC_URL> peers
```

## What the Dev Harness Demonstrates

1. **Node Info** - Client version, network info
2. **Sync Status** - Check if node is synced
3. **Peer Count** - Active peer connections
4. **Protocol Version** - Supported protocols
5. **Admin APIs** - Advanced node management

## Next Steps

Proceed to **geth/21-sync** for sync monitoring.
