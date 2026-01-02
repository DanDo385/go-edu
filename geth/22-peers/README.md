# 22-peers: Peer Discovery

## Overview

Learn about Ethereum's P2P networking layer. Discover and manage peer connections, understand the devp2p protocol.

## Learning Objectives

- Query connected peers
- Understand peer discovery mechanisms
- Monitor peer health
- Use admin_peers RPC method
- Analyze network topology

## Project Structure

```
22-peers/
├── cmd/
│   ├── app/main.go
│   └── dev/main.go
├── internal/peers/
│   ├── exercise.go
│   ├── exercise_test.go
│   ├── solution.reference.go
│   └── solution_no_err.reference.go
└── README.md
```

## CLI Arguments

```bash
# List peers
go run ./cmd/app/main.go <RPC_URL> list

# Peer details
go run ./cmd/app/main.go <RPC_URL> info <PEER_ID>

# Monitor connections
go run ./cmd/app/main.go <RPC_URL> monitor
```

## What the Dev Harness Demonstrates

1. **Peer Listing** - Get all connected peers
2. **Peer Info** - Detailed peer information
3. **Connection Quality** - Latency and reliability
4. **Network Distribution** - Geographic distribution
5. **Protocol Versions** - Supported capabilities

## Next Steps

Proceed to **geth/23-mempool** for mempool monitoring.
