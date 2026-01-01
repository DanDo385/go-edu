# 09: HTTP Server with Graceful Shutdown

## What Is This Project About?

This module teaches you how to build an HTTP server that shuts down gracefully, allowing in-flight requests to complete before stopping. This is essential for production deployments.

## Why Is This Important?

Graceful shutdown ensures:
- No dropped requests
- Clean resource cleanup
- Safe rolling deployments
- Better user experience

## Key Concepts You'll Learn

- **HTTP server**: Building with net/http
- **Signal handling**: Catching SIGINT/SIGTERM
- **Graceful shutdown**: Completing in-flight requests
- **Context cancellation**: Propagating shutdown

## Prerequisites

- Completion of `minis/08-http-client-retries`

## How to Run

```bash
go run ./cmd/dev/main.go
```

## Testing

```bash
go test -v ./...
```
