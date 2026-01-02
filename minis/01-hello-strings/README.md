# 01: Hello Strings

## What Is This Project About?

This is your first Go project in the series. You'll learn fundamental string operations including case conversion, string manipulation, and building string utilities. Strings are one of the most commonly used types in programming, and mastering string operations is essential for any Go developer.

## Why Is This Important?

String manipulation is foundational to:
- Text processing and parsing
- User input handling
- Data formatting
- Building CLI tools

## Real-World Problems This Solves

- **Text normalization**: Convert user input to consistent format
- **Title case generation**: Format titles and headers
- **String validation**: Check and transform text data

## Key Concepts You'll Learn

- **String basics**: Immutable byte sequences in Go
- **Unicode handling**: Working with runes vs bytes
- **strings package**: Standard library string utilities
- **Test-driven development**: Writing tests first

## Prerequisites

- Go 1.21+ installed
- Basic programming knowledge

## Project Structure

```
minis/01-hello-strings/
├── cmd/
│   ├── app/
│   │   └── main.go
│   └── dev/
│       └── main.go
├── internal/
│   └── hellostrings/
│       ├── exercise.go
│       ├── exercise_test.go
│       ├── solution.reference.go
│       └── solution_no_err.reference.go
└── .vscode/
    └── launch.json
```

## How to Run

```bash
# Using debug harness
go run ./cmd/dev/main.go

# Using CLI (if implemented)
go run ./cmd/app/main.go "hello world"
```

## Testing

```bash
go test -v ./...
```
