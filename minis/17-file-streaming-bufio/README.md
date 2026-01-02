# 17-file-streaming-bufio: File Streaming with bufio

## Overview

File Streaming with bufio exercise focusing on practical Go programming patterns.

## Learning Objectives

- Master file streaming with bufio concepts
- Implement production-ready code
- Write efficient and idiomatic Go
- Handle edge cases and errors

## Quick Start

```bash
# Implement the exercise
# Open internal/file/streaming/bufio/exercise.go and complete the TODOs

# Run tests
go test -v ./...

# Run application
go run ./cmd/app/main.go [args]

# Debug with fixed inputs
go run ./cmd/dev/main.go
```

## CLI Arguments

```bash
# Run the application (see cmd/app/main.go for specific arguments)
go run ./cmd/app/main.go
```

## What the Dev Harness Demonstrates

The `cmd/dev/main.go` provides a debug-friendly environment with:
- Fixed input values for consistent testing
- Clear step-by-step execution
- Breakpoint markers for debugging
- Expected output examples

## Key Concepts

See `internal/file/streaming/bufio/solution.reference.go` for:
- Complete implementation with explanations
- Best practices and patterns
- Performance considerations
- Common pitfalls and how to avoid them

## Testing

```bash
# Run all tests
go test -v ./...

# Run specific test
go test -v -run TestName ./...

# Run with coverage
go test -v -cover ./...

# Run benchmarks (if available)
go test -bench=. -benchmem ./...
```

## Next Steps

After completing this exercise, proceed to the next project in sequence.

## Resources

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Blog](https://blog.golang.org/)
