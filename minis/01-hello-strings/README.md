# 01: String Manipulation Basics

## What Is This Project About?

Learn fundamental string operations in Go including case conversion, substring operations, and string comparison. This module introduces you to Go's string handling and the strings package.

## Why Is This Important?

String manipulation is fundamental to almost every program. Whether you're parsing user input, formatting output, or processing text data, you need to master string operations.

## Real-World Problems This Solves

- **Processing user input in CLI applications**
- **Formatting text for display in UIs**
- **Parsing and validating string data**

## Key Concepts You'll Learn

- **String immutability in Go**: String immutability in Go
- **strings package functions (ToUpper, ToLower, Contains, Split, etc.)**: strings package functions (ToUpper, ToLower, Contains, Split, etc.)
- **Rune vs byte handling for Unicode**: Rune vs byte handling for Unicode
- **String concatenation performance considerations**: String concatenation performance considerations

## Prerequisites

- Basic Go syntax knowledge
- Understanding of previous minis modules (if sequential)

## Project Structure

```
minis/01-hello-strings/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application
│   └── dev/
│       └── main.go          # Debug harness
├── internal/
│   └── hellostrings/
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
go run ./cmd/app/main.go "hello world" titlecase
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

See `internal/hellostrings/exercise.go` for implementation details and exercises.

## Additional Resources

- [Go Documentation](https://go.dev/doc/)
- [Go Standard Library](https://pkg.go.dev/std)
- [Effective Go](https://go.dev/doc/effective_go)

## Next Steps

Continue with the next module to build on these concepts!
