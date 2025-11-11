# Project 03: csv-stats

## What You're Building

A CSV transaction analyzer that computes per-category statistics (count, sum, average) from financial data. This project demonstrates structured data parsing, aggregate calculations, and Go's approach to error handling in data pipelines.

## Concepts Covered

- `encoding/csv` package for structured text parsing
- Structs for data modeling
- `strconv` for string-to-number conversion
- Streaming I/O (process line-by-line, not all-at-once)
- Error handling: when to fail fast vs. accumulate errors
- Map-based aggregation
- Floating-point arithmetic caveats

## How to Run

```bash
# Run the program
make run P=03-csv-stats

# Or directly:
go run ./minis/03-csv-stats/cmd/csv-stats

# Run tests
go test ./minis/03-csv-stats/...

# Run tests with verbose output
go test -v ./minis/03-csv-stats/...
```

## Solution Explanation

### Algorithm Overview

1. **Parse CSV headers**: Use `csv.Reader.Read()` to get the first line and validate expected columns (id, category, amount)
2. **Stream records**: Read line-by-line to keep memory usage constant (don't load entire file)
3. **Aggregate by category**: Maintain a map of category → running totals (count, sum)
4. **Compute averages**: After processing all rows, divide sum by count for each category
5. **Error handling**: Return error immediately if any row is malformed (fail-fast approach)

### Why Streaming?

For a 10-row CSV, loading the entire file is fine. But for a 10GB transaction log with millions of rows, streaming is essential. Go's `csv.Reader` yields one record at a time, keeping memory usage constant regardless of file size.

### Floating-Point Precision

We use `float64` for amounts, which can accumulate rounding errors. For financial applications, consider:
- Using `int64` cents instead of `float64` dollars
- The `github.com/shopspring/decimal` package for arbitrary precision
- Rounding to 2 decimal places when displaying

## Where Go Shines

**Go vs Python:**
- Python: `pandas.read_csv()` is powerful but loads entire file into memory
- Go: Streaming CSV reader uses constant memory; perfect for large files
- Go's static typing catches schema errors at compile time

**Go vs JavaScript:**
- JS: Requires external libraries (`csv-parser`, `papaparse`) for CSV parsing
- Go: `encoding/csv` is built into the standard library
- Go's error handling is explicit; JS callbacks/promises can hide failures

**Go vs Rust:**
- Rust: `csv` crate is excellent and zero-copy where possible
- Go: Simpler API, but slightly slower due to GC and interface overhead
- Both handle errors explicitly (no exceptions!)

## Stretch Goals

1. **Add median calculation**: Track all amounts per category (requires storing values, not just sum/count)
2. **Support multiple CSV formats**: Accept column order via flags (e.g., `--schema amount,category,id`)
3. **Add date filtering**: Include a `date` column and filter by date range
4. **Output JSON**: Add a `--format json` flag to print results as JSON
5. **Parallel processing**: Split file into chunks and aggregate with goroutines (advanced!)
