# geth/08-abigen

## Problem

Problem: Use BoundContract for type-safe contract calls with automatic ABI encoding/decoding.

This module teaches you how to use go-ethereum's BoundContract pattern for cleaner,
safer contract interactions. Instead of manually encoding/decoding like module 07,
you'll use the abi package to handle encoding/decoding automatically.

Computer science principles highlighted:
  - Adapter pattern: BoundContract wraps low-level RPC with high-level interface
  - Type safety: ABI definitions provide compile-time checks
  - Code reuse: Helper functions eliminate boilerplate
  - Separation of concerns: ABI encoding is separate from business logic

## Quickstart

```bash
cd geth/08-abigen
go test ./...
```

## CLI (cmd/app)

This project includes a small CLI wrapper around the exercise package.

### Flags

- **`-rpc`**: Ethereum JSON-RPC URL (defaults to a public RPC; can also be set via env vars like `RPC_URL`)
- **`-timeout`**: RPC timeout (default `30s`)
- **`-json`**: print result (best-effort) as a single Go-struct dump

Project config flags (if present in `Config`):

- **`-a-b-i`**: `Config.ABI` (string)
- **`-contract`**: `Config.Contract` (common.Address)
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

- `internal/abigen/exercise.go`: implement the TODOs here
- `internal/abigen/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
