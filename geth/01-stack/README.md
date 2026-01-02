# geth/01-stack

## Problem

Problem: Prove RPC connectivity by reading the network identifiers and latest header.

The very first thing an Ethereum Go tool should do is dial an RPC endpoint,
retrieve the chain/network IDs (replay protection + legacy identifier), and
inspect a block header. Headers are lightweight (~500 bytes) yet contain the
state root, parent hash, and other cryptographic commitments that define the
execution stack you are about to interact with. This function mirrors the CLI
demo from module 01 but exposes it as a reusable library API.

Computer science principles highlighted:
  - Separation of configuration from code (cfg.BlockNumber allows deterministic tests)
  - Fault tolerance via context propagation—callers control cancellation/timeouts
  - Immutability via defensive copies (we never hand pointers owned by go-ethereum back to callers)

## Quickstart

```bash
cd geth/01-stack
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

- `internal/stack/exercise.go`: implement the TODOs here
- `internal/stack/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
