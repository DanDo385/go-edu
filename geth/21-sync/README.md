# 21: Sync Status

## What Is This Project About?

This module teaches you how to monitor node synchronization status and understand different sync modes. Knowing whether your node is synced is crucial for data reliability.

## Why Is This Important?

Sync monitoring enables:
- Service health checks
- Data reliability verification
- Sync progress tracking
- Troubleshooting sync issues

## Key Concepts You'll Learn

- **Sync status**: Current vs highest block
- **Sync progress**: Percentage completion
- **Sync modes**: Snap, full, light
- **eth_syncing**: RPC method for status

## Prerequisites

- Completion of `geth/20-node`

## How to Run

```bash
go run ./cmd/app/main.go https://eth.llamarpc.com
```

## Testing

```bash
go test -v ./...
```
