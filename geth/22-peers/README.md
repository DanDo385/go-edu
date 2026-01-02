# geth/22-peers

## Problem

Problem: Query the number of connected peers to assess node connectivity health.

In Ethereum's peer-to-peer network, nodes connect to other nodes (peers) to gossip
transactions and blocks. The number of connected peers is a basic health indicator:
too few peers means slow propagation of data, while zero peers means complete isolation.

The net_peerCount RPC method returns a hexadecimal string representing the count,
which the ethclient library automatically converts to uint64.

Computer science principles highlighted:
  - P2P network topology (decentralized mesh network)
  - Gossip protocols (how information spreads)
  - Health metrics and observability (monitoring system state)

## Quickstart

```bash
cd geth/22-peers
go test ./...
```

## CLI (cmd/app)

This project includes a small CLI wrapper around the exercise package.

### Flags

- **`-rpc`**: Ethereum JSON-RPC URL (defaults to a public RPC; can also be set via env vars like `RPC_URL`)
- **`-timeout`**: RPC timeout (default `30s`)
- **`-json`**: print result (best-effort) as a single Go-struct dump

Project config flags (if present in `Config`):


### Usage

```bash
go run ./cmd/app -h
```

### Copy/paste examples

```bash
go run ./cmd/app -rpc "https://eth.llamarpc.com"
go run ./cmd/dev
```

## Debug harness (cmd/dev)

`cmd/dev` runs the same CLI with preselected example arguments, so you can just hit Run/Debug and see what the project can do.

```bash
go run ./cmd/dev
```

## Files

- `internal/peers/exercise.go`: implement the TODOs here
- `internal/peers/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
