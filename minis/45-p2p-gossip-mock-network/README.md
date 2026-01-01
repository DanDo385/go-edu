# 45: P2p Gossip Mock Network

## What Is This Project About?

This project is a focused, self-contained mini that teaches a single Go concept by making you implement and test it in isolation.

You’ll work in the `internal/` exercise package (where the tests live), then use `cmd/dev` to step through deterministic examples with a debugger, and `cmd/app` to run the same idea with real CLI arguments.

## Why Is This Important?

Go rewards building strong intuition around its core primitives (types, interfaces, concurrency, I/O, and the standard library). These minis are designed to build that intuition quickly by combining tight exercises with fast feedback from tests and a debug harness.

## Real-World Problems This Solves

- Turning a language feature into reliable, testable code
- Designing small, composable APIs and data structures
- Debugging correctness and performance issues with repeatable inputs

## Key Concepts You’ll Learn

- P2p Gossip Mock Network: the core idea behind this module
- Debugging with deterministic inputs (`cmd/dev`) and tests (`go test`)
- Reading and reasoning about results + common failure modes

## Prerequisites

- Basic Go syntax (functions, structs, slices/maps)
- Comfort running `go test` and `go run`

## Project Structure

```text
45-p2p-gossip-mock-network/
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
