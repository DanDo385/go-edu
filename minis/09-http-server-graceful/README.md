# 09-http-server-graceful

**Graceful HTTP Server**

Build an HTTP server with graceful shutdown support.

## What You'll Learn

- HTTP server setup
- Signal handling (SIGINT, SIGTERM)
- Graceful shutdown patterns
- Connection draining

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Run(ctx, cfg) error` | Start server with graceful shutdown |

## Project Structure

```
09-http-server-graceful/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/httpservergraceful/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd minis/09-http-server-graceful

# Start server on port 8080
go run ./cmd/app/main.go --port 8080

# Press Ctrl+C to trigger graceful shutdown
```

### Run the Debug Harness

```bash
go run ./cmd/dev/main.go
```

### Run Tests

```bash
go test -v ./...
```

## CLI Arguments

| Argument | Description |
|----------|-------------|
| `--port` | Port to listen on (default: 8080) |
| `--shutdown-timeout` | Max shutdown time (default: 30s) |

## Quick Copy & Paste

```bash
# Start server
go run ./cmd/app/main.go --port 8080

# In another terminal, test it
curl http://localhost:8080/health

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **signal.Notify**: Catch OS signals
2. **Server.Shutdown**: Graceful shutdown
3. **Context with Timeout**: Limit shutdown time
4. **Connection Draining**: Wait for active requests

## Next Steps

After completing this exercise, proceed to `minis/10-grpc-telemetry-service`.
