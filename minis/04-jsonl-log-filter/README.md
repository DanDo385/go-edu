# Project 04: jsonl-log-filter

## What You're Building

A JSONL (JSON Lines) log parser that filters entries by severity level and sorts them chronologically. This project demonstrates custom type marshaling, streaming JSON parsing, and time handling in Go.

## Concepts Covered

- JSONL format (newline-delimited JSON)
- `encoding/json` for JSON parsing
- Custom `UnmarshalJSON` for enum-like types
- `time.Time` and RFC3339 timestamps
- `sort.Slice` with custom comparators
- Error accumulation vs. fail-fast strategies
- Streaming I/O for large log files

## How to Run

```bash
# Run the program
make run P=04-jsonl-log-filter

# Or directly:
go run ./minis/04-jsonl-log-filter/cmd/jsonl-log-filter

# Run tests
go test ./minis/04-jsonl-log-filter/...

# Run tests with verbose output
go test -v ./minis/04-jsonl-log-filter/...
```

## Solution Explanation

### Algorithm Overview

1. **Read line-by-line**: Use `bufio.Scanner` to process JSONL format (one JSON object per line)
2. **Parse each line**: `json.Unmarshal` each line into an `Entry` struct
3. **Custom unmarshaling**: Implement `UnmarshalJSON` for `Level` to convert strings ("info", "error") to enum values
4. **Filter by level**: Only keep entries >= minimum severity (e.g., "warn" and above excludes "debug" and "info")
5. **Handle parse errors**: Skip malformed lines but count them for reporting
6. **Sort by timestamp**: Use `sort.Slice` with a comparator function
7. **Return results**: Entries + error if any lines were skipped

### JSONL vs JSON Array

**JSONL** (one JSON object per line):
```
{"ts":"2024-01-01T12:00:00Z","level":"info","msg":"Server started"}
{"ts":"2024-01-01T12:00:05Z","level":"error","msg":"Database failed"}
```

**JSON Array**:
```json
[
  {"ts":"2024-01-01T12:00:00Z","level":"info","msg":"Server started"},
  {"ts":"2024-01-01T12:00:05Z","level":"error","msg":"Database failed"}
]
```

JSONL advantages:
- Streaming: Process one line at a time (constant memory)
- Append-friendly: Add new entries without modifying existing file
- Error isolation: One malformed line doesn't corrupt the entire file

## Where Go Shines

**Go vs Python:**
- Python: `json.loads(line)` per line is simple, but error handling is verbose (try/except)
- Go: Explicit error returns force you to handle parse failures
- Go's `time.Time` is more robust than Python's `datetime` (no timezone gotchas)

**Go vs JavaScript:**
- JS: `JSON.parse(line)` but requires catching exceptions for malformed JSON
- Go: Errors are values, not exceptions (more predictable control flow)
- Node.js streams are powerful but have a steeper learning curve

**Go vs Rust:**
- Rust: `serde_json` is excellent and zero-copy where possible
- Go: Simpler syntax; reflection-based JSON is easier to use
- Both handle errors explicitly (no exceptions!)

## Stretch Goals

1. **Add structured logging**: Use `log/slog` to output filtered results in JSONL format
2. **Support custom time ranges**: Filter by `--after` and `--before` flags
3. **Add full-text search**: Filter by regex on the `msg` field
4. **Colorize output**: Use ANSI codes to colorize levels (red for error, yellow for warn)
5. **Streaming output**: Write results as they're processed instead of buffering all
