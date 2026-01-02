# minis/14-methods-value-vs-pointer-receivers

## Problem

Problem: Understanding method receivers in Go (value vs pointer)

Requirements:
1. Choose correct receiver type for mutation
2. Understand interface satisfaction rules
3. Handle nil receivers safely
4. Optimize for performance (large structs)
5. Maintain API consistency

Data Structure:
- Value receiver: Operates on copy (8-64 bytes typical)
- Pointer receiver: Operates on original (8 bytes pointer)
- Method set: T has methods with receiver T, *T has both T and *T

Algorithm: Receiver Selection
- Mutation needed: Use pointer receiver
- Large struct (>64 bytes): Use pointer receiver
- Small immutable value: Use value receiver
- Interface satisfaction: Consider both T and *T

Why receiver type matters:
- Value receiver: Safe copy, can't modify original
- Pointer receiver: Can modify, more efficient for large types
- Mixed receivers cause interface satisfaction issues

## Quickstart

```bash
cd minis/14-methods-value-vs-pointer-receivers
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
go run ./cmd/app -fn "NewSafeCounterMapSolution"
go run ./cmd/dev
```

## Debug harness (cmd/dev)

`cmd/dev` runs the same CLI with preselected example arguments, so you can just hit Run/Debug and see what the project can do.

```bash
go run ./cmd/dev
```

## Files

- `internal/methodsvaluevspointerreceivers/exercise.go`: implement the TODOs here
- `internal/methodsvaluevspointerreceivers/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
