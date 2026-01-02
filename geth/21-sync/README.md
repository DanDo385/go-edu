# 21-sync: Sync Monitoring

## Overview

Monitor Ethereum node synchronization progress. Learn about different sync modes and how to track sync status.

## Learning Objectives

- Monitor sync progress
- Understand sync modes (full, fast, snap, light)
- Track block synchronization
- Estimate sync completion time
- Handle unsynchronized nodes

## Project Structure

```
21-sync/
├── cmd/
│   ├── app/main.go
│   └── dev/main.go
├── internal/sync/
│   ├── exercise.go
│   ├── exercise_test.go
│   ├── solution.reference.go
│   └── solution_no_err.reference.go
└── README.md
```

## CLI Arguments

```bash
# Monitor sync
go run ./cmd/app/main.go <RPC_URL>

# With refresh interval
go run ./cmd/app/main.go <RPC_URL> --refresh 5s
```

## What the Dev Harness Demonstrates

1. **Sync Status** - Check synchronization state
2. **Progress Tracking** - Current vs highest block
3. **Sync Mode** - Identify sync strategy
4. **Estimated Time** - Calculate time to completion
5. **Waiting Logic** - Wait for node to sync

## Next Steps

Proceed to **geth/22-peers** for peer discovery.
