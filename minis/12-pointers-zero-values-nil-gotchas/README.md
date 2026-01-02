# minis/12-pointers-zero-values-nil-gotchas

## Problem

Problem: Understanding Go's pointer semantics and nil handling

Requirements:
1. Safe pointer dereferencing with nil checks
2. In-place value swapping using pointers
3. Nil map initialization and safe usage
4. Linked list operations with nil receivers
5. Pointer vs value receiver trade-offs

Data Structure:
- Pointer: 8 bytes (memory address on 64-bit systems)
- nil: Zero value for pointers, slices, maps, channels, interfaces
- Linked list: Recursive structure using pointers

Algorithm: Nil-Safe Operations
- Always check nil before dereferencing
- Methods can be called on nil receivers
- nil map: Can read but NOT write (panic!)
- nil slice: Can read, append (allocates)

Why pointers are essential:
- Modify original value (not a copy)
- Share large structs efficiently
- Represent optional values (nil = absent)
- Build recursive data structures

## Quickstart

```bash
cd minis/12-pointers-zero-values-nil-gotchas
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

- `internal/pointerszerovaluesnilgotchas/exercise.go`: implement the TODOs here
- `internal/pointerszerovaluesnilgotchas/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
