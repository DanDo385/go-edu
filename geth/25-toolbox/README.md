# geth/25-toolbox

## Problem

Problem: Build a Swiss Army knife CLI that combines multiple node operations.

This capstone module brings together patterns from all previous modules into a single
unified tool. Instead of separate programs for each operation, you'll have one tool
with subcommands (like git, docker, kubectl).

This demonstrates:
  - Command routing and dispatch
  - Code reuse across modules
  - Building production-ready tools
  - Composing simple operations into complex workflows

Computer science principles highlighted:
  - Command pattern (encapsulating operations)
  - Composition (building complex from simple)
  - Interface segregation (ToolboxClient combines many interfaces)

## Quickstart

```bash
cd geth/25-toolbox
go test ./...
```

## CLI (cmd/app)

This project includes a small CLI wrapper around the exercise package.

### Flags

- **`-rpc`**: Ethereum JSON-RPC URL (defaults to a public RPC; can also be set via env vars like `RPC_URL`)
- **`-timeout`**: RPC timeout (default `30s`)
- **`-json`**: print result (best-effort) as a single Go-struct dump

Project config flags (if present in `Config`):

- **`-command`**: `Config.Command` (string)

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

- `internal/toolbox/exercise.go`: implement the TODOs here
- `internal/toolbox/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
