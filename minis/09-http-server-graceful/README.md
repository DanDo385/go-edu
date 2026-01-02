# minis/09-http-server-graceful

## Problem

Problem: Build an HTTP server with routes, middleware, and graceful shutdown

Requirements:
1. REST endpoints for key-value storage
2. Request counting middleware
3. Graceful shutdown on SIGINT/SIGTERM
4. JSON request/response handling

Algorithm: HTTP Request Handling
- Route requests to appropriate handlers
- Apply middleware for cross-cutting concerns
- Handle JSON encoding/decoding
- Coordinate graceful shutdown

Graceful Shutdown Algorithm:
- Receive shutdown signal (SIGINT/SIGTERM)
- Stop accepting new connections
- Wait for in-flight requests to complete
- Close server cleanly

Middleware Pattern:
- Wrap handler functions
- Execute before/after main handler
- Common uses: logging, metrics, authentication

## Quickstart

```bash
cd minis/09-http-server-graceful
go test ./...
```

## CLI (cmd/app)

This project includes a small CLI wrapper around the exercise package.

### Flags

- **`-list`**: list available exported functions
- **`-fn`**: function name to run
- **`-in`**: string input (for `func(string) ...`)
- **`-n`**: int input (for `func(int) ...`)
- **`-f`**: float64 input (for `func(float64) ...`)
- **`-b`**: bool input (for `func(bool) ...`)
- **`-file`** / **`-stdin`**: input sources for `func(io.Reader) ...`

### Usage

```bash
go run ./cmd/app -h
```

### Copy/paste examples

```bash
go run ./cmd/app -list
go run ./cmd/app -fn "NewMemStore"
go run ./cmd/dev
```

## Debug harness (cmd/dev)

`cmd/dev` runs the same CLI with preselected example arguments, so you can just hit Run/Debug and see what the project can do.

```bash
go run ./cmd/dev
```

## Files

- `internal/httpservergraceful/exercise.go`: implement the TODOs here
- `internal/httpservergraceful/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
