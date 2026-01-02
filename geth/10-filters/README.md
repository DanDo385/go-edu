# geth/10-filters

## Quickstart

```bash
cd geth/10-filters
go test ./...
```

## CLI (cmd/app)

This project includes a small CLI wrapper around the exercise package.

### Flags

- **`-rpc`**: Ethereum JSON-RPC URL (defaults to a public RPC; can also be set via env vars like `RPC_URL`)
- **`-timeout`**: RPC timeout (default `30s`)
- **`-json`**: print result (best-effort) as a single Go-struct dump

Project config flags (if present in `Config`):

- **`-max-heads`**: `Config.MaxHeads` (int)
- **`-poll-interval`**: `Config.PollInterval` (time.Duration)
- **`-poll-mode`**: `Config.PollMode` (bool)

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

- `internal/filters/exercise.go`: implement the TODOs here
- `internal/filters/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
