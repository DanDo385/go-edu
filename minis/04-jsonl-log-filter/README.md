# minis/04-jsonl-log-filter

## Problem

Problem: Parse and filter JSONL (JSON Lines) log entries by severity level

Given a file with one JSON object per line, we need to:
1. Parse each line as a structured log entry
2. Filter by minimum severity level (debug < info < warn < error)
3. Sort results by timestamp
4. Handle malformed lines gracefully (skip but report count)

Constraints:
- JSONL format: one JSON object per line (not a JSON array!)
- Timestamps are RFC3339 format (e.g., "2024-01-01T12:00:00Z")
- Level is a string ("debug", "info", "warn", "error") that must map to an enum
- Malformed lines should be skipped, not cause total failure

Time/Space Complexity:
- Time: O(n log n) where n = number of valid entries (O(n) parse + O(n log n) sort)
- Space: O(n) to store filtered entries

Why Go is well-suited:
- `encoding/json`: Robust JSON parser with struct mapping via reflection
- Custom unmarshalers: `UnmarshalJSON` interface for enum-like types
- `time.Time`: First-class time support with zone awareness
- `sort.Slice`: Inline sorting with custom comparators
- Error accumulation: Handle partial failures gracefully without exceptions

DEBUGGING THIS FILE:
==================
This solution is instrumented with extensive debugging comments to teach you
how to use Go's debugger (dlv) and VS Code's debugging features.

Key debugging concepts covered:
1. Setting breakpoints at critical JSON parsing points
2. Watching custom unmarshaling transformations
3. Using F10 (Step Over) vs F11 (Step Into) effectively
4. Inspecting sort algorithm execution
5. Using the Debug Console to evaluate expressions
6. Understanding enum-like types and their memory representations

## Quickstart

```bash
cd minis/04-jsonl-log-filter
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

- `internal/jsonllogfilter/exercise.go`: implement the TODOs here
- `internal/jsonllogfilter/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
