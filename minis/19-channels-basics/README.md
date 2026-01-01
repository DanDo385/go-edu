# 19: Channel Fundamentals

## What Is This Project About?

Master Go channels—the core primitive for goroutine communication. Learn buffered vs unbuffered channels, channel closing semantics, and the select statement.

## Why Is This Important?

Channels are Go's superpower. They enable safe communication between goroutines without explicit locks, making concurrent programming more manageable.

## Real-World Problems This Solves

- **Coordinating work between goroutines safely**
- **Implementing producer-consumer patterns**
- **Building concurrent pipelines**

## Key Concepts You'll Learn

- **Unbuffered vs buffered channels**: Unbuffered vs buffered channels
- **Channel send/receive semantics (blocking behavior)**: Channel send/receive semantics (blocking behavior)
- **Channel closing and range loops**: Channel closing and range loops
- **The select statement for multiplexing**: The select statement for multiplexing

## Prerequisites

- Basic Go syntax knowledge
- Understanding of previous minis modules (if sequential)

## Project Structure

```
minis/19-channels-basics/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application
│   └── dev/
│       └── main.go          # Debug harness
├── internal/
│   └── channelsbasics/
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

See `internal/channelsbasics/exercise.go` for implementation details and exercises.

## Additional Resources

- [Go Documentation](https://go.dev/doc/)
- [Go Standard Library](https://pkg.go.dev/std)
- [Effective Go](https://go.dev/doc/effective_go)

## Next Steps

Continue with the next module to build on these concepts!
