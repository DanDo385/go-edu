# minis/13-interfaces-duck-typing

## Problem

Problem: Understanding Go's interface system and duck typing

Requirements:
1. Implement interfaces implicitly (no "implements" keyword)
2. Use type assertions to extract concrete types
3. Handle nil interface gotchas (type vs value)
4. Compose interfaces through embedding
5. Implement polymorphism with interface dispatch

Data Structure:
- Interface value: Type pointer + Data pointer (16 bytes on 64-bit)
- Type assertion: Runtime check of concrete type
- Type switch: Multi-way type-based branching

Algorithm: Dynamic Dispatch
- Interface stores concrete type metadata
- Method calls routed through virtual table
- Type assertions inspect type metadata

Why interfaces enable polymorphism:
- Decouple behavior from implementation
- Write code once, works with many types
- No inheritance needed
- Runtime flexibility with compile-time safety

## Quickstart

```bash
cd minis/13-interfaces-duck-typing
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
go run ./cmd/dev
```

## Debug harness (cmd/dev)

`cmd/dev` runs the same CLI with preselected example arguments, so you can just hit Run/Debug and see what the project can do.

```bash
go run ./cmd/dev
```

## Files

- `internal/interfacesducktyping/exercise.go`: implement the TODOs here
- `internal/interfacesducktyping/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
