# minis/31-static-file-server

## Problem

Problem: Build a production-ready static file server with modern HTTP features

Requirements:
1. Serve files from a root directory with security (path traversal prevention)
2. ETag generation and validation (cache validation)
3. Range request support (resumable downloads, video streaming)
4. Proper caching headers (Cache-Control, Last-Modified, max-age)
5. MIME type detection (extension-based + content sniffing)
6. Directory listing (optional, with HTML generation)
7. HTTP conditional requests (If-None-Match, If-Modified-Since)

Time/Space Complexity:
- File serving: O(1) time for small files (sendfile syscall), O(n) for large files
- Directory listing: O(n log n) for sorting n entries
- Space: O(1) for file serving (streaming), O(n) for directory listing

Why Go is well-suited:
- http.ServeContent: Built-in range request and conditional request handling
- Zero-copy sendfile() syscall for efficient file serving (Linux/Unix)
- mime package: MIME type detection by extension
- filepath package: Cross-platform secure path handling
- No external dependencies needed for production-grade file server

Compared to other languages:
- Node.js (express.static): Requires library, slower (no sendfile), async complexity
- Python (http.server): Built-in but single-threaded, no range requests by default
- Rust (actix-files): Faster (zero-cost abstractions), but more complex ownership
- Nginx: Production standard, but less flexible (requires config, no programmatic control)

Go Concepts Demonstrated:
- HTTP handlers and the http.Handler interface
- Pointer vs value receivers for methods
- Interface types (ResponseWriter, FileInfo)
- Security: Path traversal prevention
- I/O: File operations, streaming responses
- HTTP caching: ETags, Last-Modified, Cache-Control
- Error handling: Logging, status codes

## Quickstart

```bash
cd minis/31-static-file-server
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
go run ./cmd/dev
```

## Debug harness (cmd/dev)

`cmd/dev` runs the same CLI with preselected example arguments, so you can just hit Run/Debug and see what the project can do.

```bash
go run ./cmd/dev
```

## Files

- `internal/staticfileserver/exercise.go`: implement the TODOs here
- `internal/staticfileserver/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
