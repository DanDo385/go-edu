# 22: Worker Pool with Backpressure

## What Is This Project About?

Implement a worker pool that handles backpressure—preventing producers from overwhelming consumers. This is critical for building stable, production-ready systems.

## Why Is This Important?

Without backpressure, fast producers can overwhelm slow consumers, leading to memory exhaustion and crashes. This pattern is essential for production systems.

## Real-World Problems This Solves

- **Preventing memory exhaustion in high-throughput systems**
- **Building stable data pipelines**
- **Implementing rate limiting and flow control**

## Key Concepts You'll Learn

- **Backpressure mechanisms (bounded channels)**: Backpressure mechanisms (bounded channels)
- **Producer rate limiting**: Producer rate limiting
- **Consumer processing guarantees**: Consumer processing guarantees
- **System stability under load**: System stability under load

## Prerequisites

- Basic Go syntax knowledge
- Understanding of previous minis modules (if sequential)

## Project Structure

```
minis/22-worker-pool-with-backpressure/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application
│   └── dev/
│       └── main.go          # Debug harness
├── internal/
│   └── workerpoolwithbackpressure/
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

See `internal/workerpoolwithbackpressure/exercise.go` for implementation details and exercises.

## Additional Resources

- [Go Documentation](https://go.dev/doc/)
- [Go Standard Library](https://pkg.go.dev/std)
- [Effective Go](https://go.dev/doc/effective_go)

## Next Steps

Continue with the next module to build on these concepts!
