# 37-http-middleware-chain

**HTTP Middleware Chain**

Build composable HTTP middleware.

## What You'll Learn

- Middleware pattern
- Function composition
- Request/response modification
- Handler wrapping

## Functions to Implement

| Function | Description |
|----------|-------------|
| Build middleware chain | Composable handlers |

## Project Structure

```
37-http-middleware-chain/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness
├── internal/httpmiddlewarechain/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/37-http-middleware-chain

# Start server with middleware
go run ./cmd/app/main.go --port 8080

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Start server
go run ./cmd/app/main.go --port 8080

# Test with curl (see middleware effects)
curl -v http://localhost:8080/

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **func(http.Handler) http.Handler**: Middleware signature
2. **Chaining**: logging(auth(handler))
3. **Before/After**: Pre and post processing
4. **Context Values**: Pass data through chain

## Next Steps

After completing this exercise, proceed to `minis/38-config-loader-env-yaml`.
