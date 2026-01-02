# geth/06-eip1559

## Problem

Problem: Build and sign an EIP-1559 dynamic fee transaction with proper fee estimation.

EIP-1559 (London upgrade, August 2021) introduced a two-part fee structure:
  - Base Fee: Algorithmically determined, burned (removed from ETH supply)
  - Priority Fee (Tip): Paid to validators, incentivizes inclusion

This is more predictable than legacy transactions where users bid against each other.

Computer science principles highlighted:
  - Algorithm design: Base fee adjusts automatically based on block fullness (control theory)
  - Economic incentives: Fee burning aligns validator and user interests
  - Defensive copying: Protect mutable big.Int values from external mutation
  - Error handling: Validate all inputs and RPC responses

## Quickstart

```bash
cd geth/06-eip1559
go test ./...
```

## CLI (cmd/app)

This project includes a small CLI wrapper around the exercise package.

### Flags

- **`-rpc`**: Ethereum JSON-RPC URL (defaults to a public RPC; can also be set via env vars like `RPC_URL`)
- **`-timeout`**: RPC timeout (default `30s`)
- **`-json`**: print result (best-effort) as a single Go-struct dump

Project config flags (if present in `Config`):

- **`-to`**: `Config.To` (common.Address)
- **`-amount-wei`**: `Config.AmountWei` (*big.Int)
- **`-gas-limit`**: `Config.GasLimit` (uint64)
- **`-max-priority-fee`**: `Config.MaxPriorityFee` (*big.Int)
- **`-max-fee`**: `Config.MaxFee` (*big.Int)
- **`-no-send`**: `Config.NoSend` (bool)
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

- `internal/eip1559/exercise.go`: implement the TODOs here
- `internal/eip1559/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
