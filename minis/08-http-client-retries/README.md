# 08-http-client-retries

**Resilient HTTP Client**

Build a resilient HTTP client with retries and exponential backoff.

## What You'll Learn

- Retry patterns
- Exponential backoff with jitter
- Context timeout handling
- Generic JSON decoding

## Functions to Implement

| Function | Description |
|----------|-------------|
| `GetJSON[T](ctx, client, url) (T, error)` | GET with retries and JSON decode |
| `doRequest[T](ctx, client, url) (T, error)` | Single request attempt |
| `isRetryable(err) bool` | Determine if error is retryable |

## Project Structure

```
08-http-client-retries/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/httpclientretries/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd minis/08-http-client-retries

# Fetch JSON from URL
go run ./cmd/app/main.go https://api.github.com/users/golang
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
| `URL` | URL to fetch JSON from |
| `--retries` | Max retry attempts (default: 3) |

## Quick Copy & Paste

```bash
# Fetch JSON
go run ./cmd/app/main.go https://api.github.com/users/golang

# With custom retries
go run ./cmd/app/main.go --retries 5 https://api.github.com/users/golang

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Exponential Backoff**: delay = base * 2^attempt
2. **Jitter**: Random variation to prevent thundering herd
3. **Retryable Errors**: 5xx, timeouts, network errors
4. **Non-Retryable**: 4xx client errors

## Next Steps

After completing this exercise, proceed to `minis/09-http-server-graceful`.
