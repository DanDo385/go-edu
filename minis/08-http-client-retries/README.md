# 08: HTTP Client with Retry Logic

## What Is This Project About?

Build a resilient HTTP client with exponential backoff retry logic. Learn to handle transient failures gracefully and build fault-tolerant applications.

## Why Is This Important?

Network requests fail. Building reliable applications means handling failures gracefully with retries, timeouts, and proper error handling.

## Real-World Problems This Solves

- **Handling transient network failures**
- **Dealing with rate limits from APIs**
- **Building resilient microservices that can survive partial failures**

## Key Concepts You'll Learn

- **Exponential backoff retry strategy**: Exponential backoff retry strategy
- **Context-based timeouts and cancellation**: Context-based timeouts and cancellation
- **Idempotency considerations for retries**: Idempotency considerations for retries
- **Circuit breaker pattern (advanced)**: Circuit breaker pattern (advanced)

## Prerequisites

- Basic Go syntax knowledge
- Understanding of previous minis modules (if sequential)

## Project Structure

```
minis/08-http-client-retries/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application
│   └── dev/
│       └── main.go          # Debug harness
├── internal/
│   └── httpclientretries/
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
go run ./cmd/app/main.go https://api.example.com 3
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

See `internal/httpclientretries/exercise.go` for implementation details and exercises.

## Additional Resources

- [Go Documentation](https://go.dev/doc/)
- [Go Standard Library](https://pkg.go.dev/std)
- [Effective Go](https://go.dev/doc/effective_go)

## Next Steps

Continue with the next module to build on these concepts!
