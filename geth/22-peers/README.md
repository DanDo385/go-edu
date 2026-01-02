# 22: Peer Management

## What Is This Project About?

This module teaches you about Ethereum's P2P networking layer and peer management. Understanding peer connections is important for node operators and network researchers.

## Why Is This Important?

Peer management enables:
- Network topology analysis
- Connection troubleshooting
- Privacy considerations
- Network research

## Key Concepts You'll Learn

- **Peer discovery**: How nodes find each other
- **Connection management**: Max peers, trusted peers
- **Network protocols**: devp2p, discv5
- **admin_peers**: RPC method for peer info

## Prerequisites

- Completion of `geth/21-sync`

## How to Run

```bash
go run ./cmd/app/main.go https://eth.llamarpc.com
```

## Testing

```bash
go test -v ./...
```
