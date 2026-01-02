# minis/03-csv-stats

## Problem

Problem: Compute per-category statistics from a CSV of financial transactions

Given a CSV with columns (id, category, amount), we need to:
1. Parse the CSV line-by-line (streaming for memory efficiency)
2. Group transactions by category
3. Compute count, sum, and average for each category
4. Handle malformed data gracefully

Constraints:
- CSV has a header row that must be validated
- Amounts are decimal numbers (use float64)
- Missing or invalid amounts should cause an error (fail-fast)
- Empty categories should be treated as an error

Time/Space Complexity:
- Time: O(n) where n = number of rows (single pass)
- Space: O(c) where c = number of unique categories (map storage)

Why Go is well-suited:
- `encoding/csv` in stdlib: No external dependencies for CSV parsing
- Streaming I/O: Process line-by-line for constant memory usage
- Strong typing: Compile-time detection of struct field mismatches
- Explicit error handling: No silent data corruption

DEBUGGING THIS FILE:
==================
This solution is instrumented with extensive debugging comments to teach you
how to use Go's debugger (dlv) and VS Code's debugging features.

Key debugging concepts covered:
1. Setting breakpoints at critical CSV parsing points
2. Watching struct field transformations in the Variables panel
3. Using F10 (Step Over) vs F11 (Step Into) effectively
4. Inspecting map aggregation patterns
5. Using the Debug Console to evaluate expressions
6. Understanding streaming I/O and memory usage

## Quickstart

```bash
cd minis/03-csv-stats
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
printf "hello\nworld\n" | go run ./cmd/app -fn "SummarizeCSV" -stdin
go run ./cmd/dev
```

## Debug harness (cmd/dev)

`cmd/dev` runs the same CLI with preselected example arguments, so you can just hit Run/Debug and see what the project can do.

```bash
go run ./cmd/dev
```

## Files

- `internal/csvstats/exercise.go`: implement the TODOs here
- `internal/csvstats/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
