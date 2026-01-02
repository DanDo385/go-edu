# minis/11-slices-internals-capacity-growth

## Problem

Problem: Understanding Go slice internals and capacity growth patterns

Requirements:
1. Track capacity changes during append operations
2. Detect when slices share backing arrays
3. Safely truncate slices to allow garbage collection
4. Compare pre-allocation vs dynamic growth
5. Create capacity-limited sub-slices

Data Structure:
- Slice: Pointer to array + Length + Capacity (24 bytes on 64-bit)
- Backing array: Contiguous memory holding actual elements
- Multiple slices can reference same backing array

Algorithm: Slice Growth Strategy
- For cap < 256: new_cap = old_cap * 2 (doubling)
- For cap >= 256: new_cap ≈ old_cap * 1.25 + 192 (slower growth)
- Reallocation: Allocate new array, copy elements, return new slice

Why slices are tricky:
- Slicing creates views into same array (memory sharing)
- append() may or may not reallocate (depends on capacity)
- Holding small slice can prevent large array from being GC'd

## Quickstart

```bash
cd minis/11-slices-internals-capacity-growth
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
go run ./cmd/app -fn "PreallocateVsDynamic" -n 10
go run ./cmd/dev
```

## Debug harness (cmd/dev)

`cmd/dev` runs the same CLI with preselected example arguments, so you can just hit Run/Debug and see what the project can do.

```bash
go run ./cmd/dev
```

## Files

- `internal/slicesinternalscapacitygrowth/exercise.go`: implement the TODOs here
- `internal/slicesinternalscapacitygrowth/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
