# geth/23-mempool

## Problem

Problem: Inspect the mempool (transaction pool) to understand pending transactions.

The mempool contains transactions that have been broadcast to the network but not
yet included in a block. Monitoring the mempool helps you understand network congestion,
estimate gas prices, and track your own transactions.

However, mempool visibility is limited for privacy and security reasons. Many public
RPC endpoints don't expose pending transaction details.

Computer science principles highlighted:
  - Queue management (FIFO with priority)
  - Privacy/security trade-offs (transparency vs exploitation)
  - Resource management (mempool size limits)

## Quickstart

```bash
cd geth/23-mempool
go test ./...
```

## CLI (cmd/app)

This project includes a small CLI wrapper around the exercise package.

### Flags

- **`-rpc`**: Ethereum JSON-RPC URL (defaults to a public RPC; can also be set via env vars like `RPC_URL`)
- **`-timeout`**: RPC timeout (default `30s`)
- **`-json`**: print result (best-effort) as a single Go-struct dump

Project config flags (if present in `Config`):

- **`-limit`**: `Config.Limit` (int)

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

- `internal/mempool/exercise.go`: implement the TODOs here
- `internal/mempool/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
