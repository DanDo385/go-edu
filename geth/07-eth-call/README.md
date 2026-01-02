# geth/07-eth-call

## Problem

Problem: Query ERC20 token metadata using manual ABI encoding/decoding.

This module teaches you how to interact with contracts without using typed bindings.
You'll manually encode function selectors and decode return values, giving you a deep
understanding of how contract calls work at the ABI level.

Computer science principles highlighted:
  - ABI encoding/decoding: Understanding how function calls are encoded as bytes
  - Function selectors: First 4 bytes of keccak256(functionSignature)
  - eth_call: Simulating contract execution without sending transactions
  - Manual memory management: Decoding dynamic types (strings) from raw bytes

## Quickstart

```bash
cd geth/07-eth-call
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

- `internal/ethcall/exercise.go`: implement the TODOs here
- `internal/ethcall/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
