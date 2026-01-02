# 03 Csv Stats

## Problem

Compute per-category statistics from a CSV of financial transactions

## Constraints

- CSV has a header row that must be validated

## Complexity

- Time: O(n) where n = number of rows (single pass)

## Project Structure

```
03-csv-stats/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application with custom arguments
│   └── dev/
│       └── main.go          # Debug harness with auto-demo
├── internal/
│   └── csvstats/
│       ├── exercise.go      # YOUR CODE GOES HERE
│       ├── exercise_test.go # Test cases
│       └── solution.reference.go  # Reference implementation
└── README.md
```

## Getting Started

1. Navigate to this directory:
   ```bash
   cd minis/03-csv-stats
   ```

2. Open the exercise file:
   ```bash
   code internal/csvstats/exercise.go
   ```

3. Implement the TODO functions

4. Run tests:
   ```bash
   go test -v ./...
   ```

## CLI Usage

### Using cmd/app/main.go

This project's CLI application accepts the following arguments:

```bash
go run ./cmd/app/main.go [arguments]
```

Run without arguments to see usage information:

```bash
go run ./cmd/app/main.go
```

### Using cmd/dev/main.go

The `cmd/dev/main.go` file automatically demonstrates the project's capabilities
by running through different scenarios with pre-configured inputs.

**Run the demo:**

```bash
go run ./cmd/dev/main.go
```

**Debug with VS Code:**

1. Open `cmd/dev/main.go`
2. Set breakpoints at `// BREAKPOINT:` comments
3. Press F5 and select "Debug cmd/dev (Debug Harness)"

## Testing

Run all tests:

```bash
go test -v ./...
```

Run specific test:

```bash
go test -v -run TestFunctionName ./...
```

## Reference Solution

If you get stuck, check `internal/csvstats/solution.reference.go` for a complete
implementation with detailed explanations.

