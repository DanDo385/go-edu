# geth/12-proofs

## Problem

Problem: Fetch Merkle-Patricia trie proofs for accounts and storage slots.

Merkle-Patricia trie proofs are cryptographic receipts that prove "account X has
balance Y and storage slot Z has value W" without downloading the entire blockchain
state. This enables:
  - Light clients that verify state without full sync
  - Cross-chain bridges that prove state on one chain to another
  - Indexers that verify indexed data is correct
  - Trust-minimized verification of contract state

Computer science principles highlighted:
  - Merkle trees provide logarithmic proof size (log N instead of N)
  - Cryptographic commitment (root hash commits to all data)
  - Path-based verification (prove membership by providing path to root)

## Quickstart

```bash
cd geth/12-proofs
go test ./...
```

## CLI (cmd/app)

This project includes a small CLI wrapper around the exercise package.

### Flags

- **`-rpc`**: Ethereum JSON-RPC URL (defaults to a public RPC; can also be set via env vars like `RPC_URL`)
- **`-timeout`**: RPC timeout (default `30s`)
- **`-json`**: print result (best-effort) as a single Go-struct dump

Project config flags (if present in `Config`):

- **`-account`**: `Config.Account` (common.Address)
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

- `internal/proofs/exercise.go`: implement the TODOs here
- `internal/proofs/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
