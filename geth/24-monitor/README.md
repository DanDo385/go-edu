# geth/24-monitor

## Problem

Problem: Implement node health monitoring by checking block freshness and detecting lag.

Monitoring nodes is critical for production systems. A stale node (not receiving new blocks)
will return outdated data, causing issues for applications. By comparing the latest block's
timestamp to the current time, we can detect if a node is lagging behind the network.

Computer science principles highlighted:
  - Time-based health checks (staleness detection)
  - Threshold-based alerting (classify OK vs STALE)
  - Observability patterns (monitoring system health)

## Quickstart

```bash
cd geth/24-monitor
go test ./...
```

## CLI (cmd/app)

This project includes a small CLI wrapper around the exercise package.

### Flags

- **`-rpc`**: Ethereum JSON-RPC URL (defaults to a public RPC; can also be set via env vars like `RPC_URL`)
- **`-timeout`**: RPC timeout (default `30s`)
- **`-json`**: print result (best-effort) as a single Go-struct dump

Project config flags (if present in `Config`):

- **`-block-number`**: `Config.BlockNumber` (*big.Int)

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

- `internal/monitor/exercise.go`: implement the TODOs here
- `internal/monitor/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
