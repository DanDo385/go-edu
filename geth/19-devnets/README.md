# 19: Devnets

## What Is This Project About?

This project is part of the `geth/` track and teaches an Ethereum concept using Go + go-ethereum.

You’ll connect to an RPC node, query chain data, and learn the underlying primitives (blocks, transactions, calls, logs) that most higher-level tooling is built on.

## Why Is This Important?

Ethereum development is ultimately about understanding what the node exposes over JSON-RPC and how those APIs map to on-chain state. Learning the concepts in small, targeted projects makes later contract interaction and infra work dramatically easier.

## Real-World Problems This Solves

- Building reliable RPC-based tooling (indexers, monitors, analyzers)
- Debugging on-chain behavior by inspecting canonical node data structures
- Bridging console/RPC intuition into production Go services

## Key Concepts You’ll Learn

- Devnets: the core idea behind this module
- Debugging with deterministic inputs (`cmd/dev`) and tests (`go test`)
- Reading and reasoning about results + common failure modes

## Prerequisites

- Completion of earlier geth modules in sequence (recommended)

## Project Structure

```text
19-devnets/
  cmd/
    app/  # Application entry point (CLI arguments)
    dev/  # Debug harness (fixed inputs)
  internal/
    <package>/  # Exercise implementation
      exercise.go
      exercise_test.go
      solution.reference.go
      solution_no_err.reference.go
  .vscode/
    launch.json  # Debug configurations
```

## How to Run

### Using `cmd/app/main.go` (CLI arguments)

```bash
# from this project directory
go run ./cmd/app
```

### Using `cmd/dev/main.go` (debug harness)

```bash
# from this project directory
go run ./cmd/dev
```

### How to Debug

- Set breakpoints at `// BREAKPOINT:` comments
- Press **F5** and select:
  - **Debug: cmd/app (with RPC_URL argument support)** (geth) / the closest matching config (minis)
  - **Debug: cmd/dev (Debug Harness)**
  - **Test: Run All Tests** / **Test: Current Test Function**

## Testing

```bash
go test ./...
go test -v ./...
go test -v -run TestName ./...
go test -tags=reference -v ./...
```

## Exercises

- Implement the functions in `internal/<package>/exercise.go` until tests pass.

## Additional Resources

- Go testing: `go help test`
- go-ethereum docs (geth track): `github.com/ethereum/go-ethereum`
