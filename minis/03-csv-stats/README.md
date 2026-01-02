# 03-csv-stats

**CSV Statistics**

Compute per-category statistics from a CSV of financial transactions.

## What You'll Learn

- CSV parsing with encoding/csv
- Aggregating data by category
- Float parsing and arithmetic
- Error handling for malformed data

## Functions to Implement

| Function | Description |
|----------|-------------|
| `SummarizeCSV(r io.Reader) (map[string]Stat, error)` | Parse CSV and compute category stats |

## Project Structure

```
03-csv-stats/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/csvstats/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
├── testdata/
│   └── transactions.csv # Sample CSV file
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd minis/03-csv-stats

# Summarize CSV file
go run ./cmd/app/main.go testdata/transactions.csv
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
| `FILE` | Path to CSV file |

## Quick Copy & Paste

```bash
# Summarize transactions
go run ./cmd/app/main.go testdata/transactions.csv

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **encoding/csv**: Standard library CSV parser
2. **Struct Aggregation**: Count, Sum, Average
3. **strconv.ParseFloat**: Convert strings to floats

## Next Steps

After completing this exercise, proceed to `minis/04-jsonl-log-filter`.
