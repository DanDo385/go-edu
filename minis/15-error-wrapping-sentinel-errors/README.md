# minis/15-error-wrapping-sentinel-errors

## Problem

Problem: Understanding Go's error handling patterns

Requirements:
1. Use sentinel errors for simple conditions
2. Wrap errors with context using %w
3. Check error identity with errors.Is
4. Extract error types with errors.As
5. Handle multiple errors

Data Structure:
- error: Interface with Error() string method
- Sentinel error: Pre-declared error value
- Wrapped error: Error chain with context
- Custom error: Struct implementing error interface

Algorithm: Error Chain Traversal
- errors.Is: Walk chain, check identity
- errors.As: Walk chain, extract type
- Unwrap(): Return next error in chain

Why error handling is critical:
- Explicit error returns (no exceptions)
- Error chains preserve context
- Type-safe error inspection
- Composable error handling

## Quickstart

```bash
cd minis/15-error-wrapping-sentinel-errors
go test ./...
```

## CLI (cmd/app)

This project includes a small CLI wrapper around the exercise package.

### Flags

- **`-list`**: list available exported functions
- **`-fn`**: function name to run
- **`-in`**: string input (for `func(string) ...`)
- **`-n`**: int input (for `func(int) ...`)
- **`-f`**: float64 input (for `func(float64) ...`)
- **`-b`**: bool input (for `func(bool) ...`)
- **`-file`** / **`-stdin`**: input sources for `func(io.Reader) ...`

### Usage

```bash
go run ./cmd/app -h
```

### Copy/paste examples

```bash
go run ./cmd/app -list
go run ./cmd/app -fn "CreateUser" -in "Hello, 世界 👋"
go run ./cmd/app -fn "FindUser" -n 10
go run ./cmd/app -fn "GetUserWithFallback" -n 10
go run ./cmd/dev
```

## Debug harness (cmd/dev)

`cmd/dev` runs the same CLI with preselected example arguments, so you can just hit Run/Debug and see what the project can do.

```bash
go run ./cmd/dev
```

## Files

- `internal/errorwrappingsentinelerrors/exercise.go`: implement the TODOs here
- `internal/errorwrappingsentinelerrors/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
