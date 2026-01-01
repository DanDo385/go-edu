# 10: gRPC Service Implementation

## What Is This Project About?

Build a gRPC telemetry service with Protocol Buffers. Learn modern RPC patterns, code generation, and efficient binary protocols.

## Why Is This Important?

gRPC is the de-facto standard for microservice communication. It's used by Google, Netflix, Square, and thousands of other companies for high-performance RPC.

## Real-World Problems This Solves

- **Building efficient microservice APIs**
- **Implementing type-safe RPC calls**
- **Achieving better performance than JSON/REST**

## Key Concepts You'll Learn

- **Protocol Buffers schema definition**: Protocol Buffers schema definition
- **gRPC code generation with protoc**: gRPC code generation with protoc
- **Streaming RPC (unary, server-streaming, client-streaming, bidirectional)**: Streaming RPC (unary, server-streaming, client-streaming, bidirectional)
- **gRPC error handling and status codes**: gRPC error handling and status codes

## Prerequisites

- Basic Go syntax knowledge
- Understanding of previous minis modules (if sequential)

## Project Structure

```
minis/10-grpc-telemetry-service/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application
│   └── dev/
│       └── main.go          # Debug harness
├── internal/
│   └── grpctelemetryservice/
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

See `internal/grpctelemetryservice/exercise.go` for implementation details and exercises.

## Additional Resources

- [Go Documentation](https://go.dev/doc/)
- [Go Standard Library](https://pkg.go.dev/std)
- [Effective Go](https://go.dev/doc/effective_go)

## Next Steps

Continue with the next module to build on these concepts!
