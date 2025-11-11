# Project 09: http-server-graceful

## What You're Building

An HTTP server with routes, middleware, and graceful shutdown that waits for in-flight requests to complete before stopping.

## Concepts Covered

- `http.ServeMux` for routing
- Middleware pattern (request interceptors)
- Graceful shutdown with `os/signal`
- `context.Context` for shutdown coordination
- Interface-based storage abstraction

## How to Run

```bash
# Run the server
go run ./minis/09-http-server-graceful/cmd/http-server-graceful

# Test with curl
curl -X POST http://localhost:8080/kv -d '{"key":"name","val":"Go"}'
curl http://localhost:8080/kv?k=name

# Stop with Ctrl+C (graceful shutdown)
```

## Solution Explanation

Graceful shutdown prevents:
- Dropping in-flight requests
- Data corruption mid-write
- Client connection errors

Process:
1. Receive SIGINT/SIGTERM
2. Stop accepting new connections
3. Wait for active requests (timeout 5s)
4. Clean shutdown

## Stretch Goals

1. Add metrics endpoint (/metrics)
2. Implement rate limiting middleware
3. Add TLS support
