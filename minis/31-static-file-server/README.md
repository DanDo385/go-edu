# 31-static-file-server

**Static File Server**

Build an HTTP static file server.

## What You'll Learn

- http.FileServer
- Custom handlers
- Security considerations
- Embedding files

## Functions to Implement

| Function | Description |
|----------|-------------|
| Implement static server | Serve files from directory |

## Project Structure

```
31-static-file-server/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness
├── internal/staticfileserver/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/31-static-file-server

# Serve current directory
go run ./cmd/app/main.go --port 8080 --dir .

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Serve files on port 8080
go run ./cmd/app/main.go --port 8080 --dir ./public

# Test it
curl http://localhost:8080/

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **http.FileServer**: Serves static files
2. **http.StripPrefix**: URL path manipulation
3. **embed.FS**: Embed files in binary
4. **Directory Traversal**: Security prevention

## Next Steps

After completing this exercise, proceed to `minis/32-websocket-chatroom`.
