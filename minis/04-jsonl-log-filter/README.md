# 04-jsonl-log-filter

**JSON Lines Log Parser**

Parse and filter JSONL (JSON Lines) log entries by severity level.

## What You'll Learn

- JSONL format (one JSON per line)
- Custom JSON unmarshaling
- Log level filtering
- Sorting by timestamp

## Functions to Implement

| Function | Description |
|----------|-------------|
| `FilterLogs(r io.Reader, minLevel Level) ([]Entry, error)` | Parse and filter logs |
| `(l *Level) UnmarshalJSON(data []byte) error` | Custom JSON decoder for Level |

## Project Structure

```
04-jsonl-log-filter/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/jsonllogfilter/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
├── testdata/
│   └── logs.jsonl       # Sample log file
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd minis/04-jsonl-log-filter

# Filter logs by level
go run ./cmd/app/main.go testdata/logs.jsonl warn
go run ./cmd/app/main.go testdata/logs.jsonl error
```

### Run the Debug Harness

```bash
go run ./cmd/dev/main.go
```

### Run Tests

```bash
go test -v ./...
```

## CLI Arguments

| Argument | Description |
|----------|-------------|
| `FILE` | Path to JSONL file |
| `LEVEL` | Minimum level: debug, info, warn, error |

## Quick Copy & Paste

```bash
# Show only errors
go run ./cmd/app/main.go testdata/logs.jsonl error

# Show warnings and above
go run ./cmd/app/main.go testdata/logs.jsonl warn

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **JSONL**: One JSON object per line
2. **json.Unmarshal**: Parse JSON into structs
3. **Custom Unmarshaler**: Convert string → enum
4. **sort.Slice**: Sort by timestamp

## Next Steps

After completing this exercise, proceed to `minis/05-cli-todo-files`.
