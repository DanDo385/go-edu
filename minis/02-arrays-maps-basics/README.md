# 02: Arrays and Maps Fundamentals

## What Is This Project About?

Master Go's array and map data structures. Learn when to use arrays vs slices, how to work with maps, and understand the performance characteristics of each.

## Why Is This Important?

Arrays and maps are fundamental data structures. Understanding their internals, performance characteristics, and idiomatic usage is essential for writing efficient Go code.

## Real-World Problems This Solves

- **Storing and looking up data efficiently**
- **Choosing the right data structure for your use case**
- **Understanding when to use arrays vs slices vs maps**

## Key Concepts You'll Learn

- **Array fixed-size semantics**: Array fixed-size semantics
- **Map (hash table) implementation**: Map (hash table) implementation
- **Slice header structure**: Slice header structure
- **Map iteration order (non-deterministic)**: Map iteration order (non-deterministic)

## Prerequisites

- Basic Go syntax knowledge
- Understanding of previous minis modules (if sequential)

## Project Structure

```
minis/02-arrays-maps-basics/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application
│   └── dev/
│       └── main.go          # Debug harness
├── internal/
│   └── arraysmapsbasics/
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
go run ./cmd/app/main.go
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

See `internal/arraysmapsbasics/exercise.go` for implementation details and exercises.

## Additional Resources

- [Go Documentation](https://go.dev/doc/)
- [Go Standard Library](https://pkg.go.dev/std)
- [Effective Go](https://go.dev/doc/effective_go)

## Next Steps

Continue with the next module to build on these concepts!
