# 03: CSV Processing and Statistics

## What Is This Project About?

Parse CSV files and compute statistics. This module teaches file I/O, CSV parsing with encoding/csv, and basic statistical calculations.

## Why Is This Important?

CSV is one of the most common data exchange formats. Being able to parse, process, and analyze CSV data is essential for data pipelines, reporting, and analysis tools.

## Real-World Problems This Solves

- **Importing data from external systems (most support CSV export)**
- **Processing large datasets efficiently**
- **Generating reports and analytics**

## Key Concepts You'll Learn

- **CSV parsing with encoding/csv**: CSV parsing with encoding/csv
- **File I/O with os and io packages**: File I/O with os and io packages
- **Statistical calculations (mean, median, sum)**: Statistical calculations (mean, median, sum)
- **Error handling for file operations**: Error handling for file operations

## Prerequisites

- Basic Go syntax knowledge
- Understanding of previous minis modules (if sequential)

## Project Structure

```
minis/03-csv-stats/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application
│   └── dev/
│       └── main.go          # Debug harness
├── internal/
│   └── csvstats/
│       ├── exercise.go      # Your implementation
│       ├── exercise_test.go # Tests
│       ├── solution.reference.go      # Reference solution
│       └── solution_no_err.reference.go # Simplified reference
└── .vscode/
    └── launch.json          # Debug configurations
```

## How to Run

### Using cmd/app/main.go (CLI Arguments)

```bash
go run ./cmd/app/main.go <csv_file>
```

### Using cmd/dev/main.go (Debug Harness)

```bash
# Run with fixed test inputs (recommended for learning)
go run ./cmd/dev/main.go

# Or use VS Code debugger (F5)
```

## How to Debug

1. Set breakpoints at `// BREAKPOINT:` comments in the code
2. Press F5 in VS Code and select "Debug: cmd/dev (Debug Harness)"
3. Step through code:
   - F10 (Step Over) - Execute current line
   - F11 (Step Into) - Enter function calls
4. Watch variables in the Variables panel
5. Inspect call stack to understand execution flow

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -v -run TestFunctionName ./...

# Run with reference implementation
go test -tags=reference -v ./...
```

## Exercises

See `internal/csvstats/exercise.go` for implementation details and exercises.

## Additional Resources

- [Go Documentation](https://go.dev/doc/)
- [Go Standard Library](https://pkg.go.dev/std)
- [Effective Go](https://go.dev/doc/effective_go)

## Next Steps

Continue with the next module to build on these concepts!
