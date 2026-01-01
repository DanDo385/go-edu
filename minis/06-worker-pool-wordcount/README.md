# 06: Worker Pool Pattern for URL Processing

## What Is This Project About?

Implement a worker pool to concurrently fetch and process URLs. This module teaches you the worker pool concurrency pattern—one of the most important patterns in Go.

## Why Is This Important?

Worker pools are the foundation of building high-performance concurrent applications. This pattern appears everywhere: web servers, data processors, task queues, and distributed systems.

## Real-World Problems This Solves

- **Processing many tasks concurrently without overwhelming resources**
- **Implementing bounded concurrency (limit max goroutines)**
- **Building responsive applications that can process work efficiently**

## Key Concepts You'll Learn

- **Worker pool pattern (fixed number of workers)**: Worker pool pattern (fixed number of workers)
- **Channel-based work distribution**: Channel-based work distribution
- **Graceful shutdown of worker pools**: Graceful shutdown of worker pools
- **Bounded concurrency for resource management**: Bounded concurrency for resource management

## Prerequisites

- Basic Go syntax knowledge
- Understanding of previous minis modules (if sequential)

## Project Structure

```
minis/06-worker-pool-wordcount/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application
│   └── dev/
│       └── main.go          # Debug harness
├── internal/
│   └── workerpoolwordcount/
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
go run ./cmd/app/main.go https://example.com https://example.org
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

See `internal/workerpoolwordcount/exercise.go` for implementation details and exercises.

## Additional Resources

- [Go Documentation](https://go.dev/doc/)
- [Go Standard Library](https://pkg.go.dev/std)
- [Effective Go](https://go.dev/doc/effective_go)

## Next Steps

Continue with the next module to build on these concepts!
