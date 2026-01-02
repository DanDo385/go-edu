# 02-arrays-maps-basics

**Word Frequency Counter**

Count word frequencies from text input and find the most common word.

## What You'll Learn

- Maps for counting
- Reading from io.Reader
- String normalization (lowercase)
- Finding max in a map

## Functions to Implement

| Function | Description |
|----------|-------------|
| `FreqFromReader(r io.Reader) (map[string]int, string, error)` | Count words and return most common |

## Project Structure

```
02-arrays-maps-basics/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/arraysmapsbasics/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
├── testdata/
│   └── input.txt        # Sample input file
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd minis/02-arrays-maps-basics

# Count words from file
go run ./cmd/app/main.go testdata/input.txt

# Count words from stdin
echo "hello world hello" | go run ./cmd/app/main.go -
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
| `FILE` | Path to text file, or `-` for stdin |

## Quick Copy & Paste

```bash
# From file
go run ./cmd/app/main.go testdata/input.txt

# From stdin
echo "the quick brown fox jumps over the lazy dog" | go run ./cmd/app/main.go -

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Maps**: Go's hash table implementation
2. **bufio.Scanner**: Efficient line/word reading
3. **strings.Fields**: Split on whitespace
4. **strings.ToLower**: Normalize case

## Next Steps

After completing this exercise, proceed to `minis/03-csv-stats`.
