# geth/21-sync

## Problem

Problem: Inspect sync progress to determine if your Ethereum node is fully synced.

When running an Ethereum node, the first critical check is whether it's finished
syncing the blockchain. A non-synced node returns stale data and shouldn't be used
for production queries. The SyncProgress RPC call returns nil when fully synced,
or a progress object with current/highest block numbers when syncing.

Computer science principles highlighted:
  - Nil as a sentinel value (nil = fully synced, non-nil = syncing)
  - Progress tracking via counters (current vs highest block)
  - State inspection without mutation (read-only health check)

## Quickstart

```bash
cd geth/21-sync
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

- `internal/sync/exercise.go`: implement the TODOs here
- `internal/sync/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
