# 36-caching-reverse-proxy

**Caching Reverse Proxy**

Build a reverse proxy with response caching.

## What You'll Learn

- httputil.ReverseProxy
- Response caching
- Cache invalidation
- Conditional requests

## Functions to Implement

| Function | Description |
|----------|-------------|
| Implement caching proxy | Proxy with response cache |

## Project Structure

```
36-caching-reverse-proxy/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness
├── internal/cachingreverseproxy/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/36-caching-reverse-proxy

# Start caching proxy
go run ./cmd/app/main.go --target https://api.example.com --cache-ttl 60s

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Start proxy to httpbin
go run ./cmd/app/main.go --target https://httpbin.org --cache-ttl 60s

# Test through proxy
curl http://localhost:8080/get

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **ReverseProxy**: Forward requests to backend
2. **ModifyResponse**: Intercept and cache
3. **Cache Key**: URL + headers
4. **TTL**: Time-based expiration

## Next Steps

After completing this exercise, proceed to `minis/37-http-middleware-chain`.
