# geth/16-concurrency

## Problem

Problem: Probe multiple endpoints concurrently using a bounded worker pool.

When building Ethereum tooling, you often need to query multiple RPC endpoints,
check health of multiple nodes, or fetch data from multiple sources. Doing this
sequentially is slow. Doing it with unbounded goroutines risks exhausting resources
or hitting rate limits. A worker pool is the idiomatic Go solution.

Computer science principles highlighted:
  - Concurrency patterns: Worker pool with channels prevents unbounded goroutine creation
  - Resource management: Bounded workers respect system limits and RPC rate limits
  - Context propagation: Timeouts and cancellation cascade through concurrent operations
  - Safe concurrent access: Mutex-protected maps prevent data races when aggregating results

## Quickstart

```bash
cd geth/16-concurrency
go test ./...
```

## CLI (cmd/app)

This project includes a small CLI wrapper around the exercise package.

### Flags

- **`-rpc`**: Ethereum JSON-RPC URL (defaults to a public RPC; can also be set via env vars like `RPC_URL`)
- **`-timeout`**: RPC timeout (default `30s`)
- **`-json`**: print result (best-effort) as a single Go-struct dump

Project config flags (if present in `Config`):

- **`-workers`**: `Config.Workers` (int)
- **`-timeout`**: `Config.Timeout` (time.Duration)

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

- `internal/concurrency/exercise.go`: implement the TODOs here
- `internal/concurrency/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
