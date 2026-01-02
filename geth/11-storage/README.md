# geth/11-storage

## Problem

Problem: Read raw storage slots from Ethereum contracts, including mapping slots.

Storage is the cryptographic database where contracts store their persistent state.
Every contract has 2^256 possible 32-byte slots organized as a Merkle-Patricia trie.
Understanding storage layout is essential for:
  - Debugging contract state
  - Building indexers that track specific contract data
  - Verifying storage proofs (module 12)
  - Optimizing gas costs (packed storage)

Computer science principles highlighted:
  - Cryptographic commitment via Merkle trees (storage root commits to all slots)
  - Deterministic slot calculation (mapping slots via keccak256 hash)
  - Key-value store abstraction (2^256 address space maps to 32-byte values)

## Quickstart

```bash
cd geth/11-storage
go test ./...
```

## CLI (cmd/app)

This project includes a small CLI wrapper around the exercise package.

### Flags

- **`-rpc`**: Ethereum JSON-RPC URL (defaults to a public RPC; can also be set via env vars like `RPC_URL`)
- **`-timeout`**: RPC timeout (default `30s`)
- **`-json`**: print result (best-effort) as a single Go-struct dump

Project config flags (if present in `Config`):

- **`-contract`**: `Config.Contract` (common.Address)
- **`-slot`**: `Config.Slot` (*big.Int)
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

- `internal/storage/exercise.go`: implement the TODOs here
- `internal/storage/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
