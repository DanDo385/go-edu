# geth/09-events

## Problem

Problem: Query and decode ERC20 Transfer events from blockchain logs.

This module teaches you how to work with Ethereum events/logs. Events are append-only
records emitted during contract execution. They're searchable via bloom filters and
provide an efficient way to track state changes without querying contract storage.

Computer science principles highlighted:
  - Event-driven architecture: Logs as append-only audit trail
  - Bloom filters: Probabilistic data structures for efficient searching
  - Indexed vs non-indexed parameters: Trade-off between searchability and cost
  - Log structure: Topics (indexed) vs Data (non-indexed)

## Quickstart

```bash
cd geth/09-events
go test ./...
```

## CLI (cmd/app)

This project includes a small CLI wrapper around the exercise package.

### Flags

- **`-rpc`**: Ethereum JSON-RPC URL (defaults to a public RPC; can also be set via env vars like `RPC_URL`)
- **`-timeout`**: RPC timeout (default `30s`)
- **`-json`**: print result (best-effort) as a single Go-struct dump

Project config flags (if present in `Config`):

- **`-token`**: `Config.Token` (common.Address)
- **`-from-block`**: `Config.FromBlock` (*big.Int)
- **`-to-block`**: `Config.ToBlock` (*big.Int)

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

- `internal/events/exercise.go`: implement the TODOs here
- `internal/events/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
