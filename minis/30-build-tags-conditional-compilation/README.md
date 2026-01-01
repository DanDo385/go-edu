# 30: Build Tags and Conditional Compilation

## What Is This Project About?

This module teaches you how to use build tags for conditional compilation.

## Key Concepts

- **Build tags**: //go:build constraints
- **GOOS/GOARCH**: Platform-specific code
- **Custom tags**: Feature flags
- **File naming**: _linux.go, _test.go

## How to Run

```bash
go run ./cmd/dev/main.go
go build -tags=mytag ./...
```
