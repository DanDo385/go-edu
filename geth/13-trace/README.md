# geth/13-trace

## Problem

Problem: Trace transaction execution to see opcode-level details and gas usage.

Transaction tracing replays a transaction in the EVM and returns structured data
describing every operation (call, gas usage, storage changes, etc.). This is
essential for:
  - Debugging contract behavior (why did this revert?)
  - Analyzing gas usage (which operations are expensive?)
  - Understanding internal calls (what contracts were called?)
  - Building block explorers and analytics tools

Computer science principles highlighted:
  - Deterministic replay (same inputs → same execution trace)
  - Execution instrumentation (observing without changing behavior)
  - JSON as a universal interchange format for complex data

## Quickstart

```bash
cd geth/13-trace
go test ./...
```

## CLI (cmd/app)

This project includes a small CLI wrapper around the exercise package.

### Flags

- **`-rpc`**: Ethereum JSON-RPC URL (defaults to a public RPC; can also be set via env vars like `RPC_URL`)
- **`-timeout`**: RPC timeout (default `30s`)
- **`-json`**: print result (best-effort) as a single Go-struct dump

Project config flags (if present in `Config`):

- **`-tx-hash`**: `Config.TxHash` (common.Hash)

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

- `internal/trace/exercise.go`: implement the TODOs here
- `internal/trace/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
