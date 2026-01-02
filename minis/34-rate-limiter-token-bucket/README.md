# 34-rate-limiter-token-bucket

**Token Bucket Rate Limiter**

Implement rate limiting with token bucket algorithm.

## What You'll Learn

- Token bucket algorithm
- Rate limiting patterns
- golang.org/x/time/rate
- HTTP middleware

## Functions to Implement

| Function | Description |
|----------|-------------|
| Implement rate limiter | Token bucket algorithm |

## Project Structure

```
34-rate-limiter-token-bucket/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness
├── internal/ratelimitertokenbucket/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/34-rate-limiter-token-bucket

# Start server with rate limiting
go run ./cmd/app/main.go --rate 10 --burst 20

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Start rate-limited server
go run ./cmd/app/main.go --rate 10 --burst 20

# Test rate limiting
for i in {1..30}; do curl -s http://localhost:8080/ & done

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Token Bucket**: Steady rate with bursts
2. **Rate**: Tokens added per second
3. **Burst**: Maximum tokens (bucket size)
4. **Wait vs Allow**: Block or reject

## Next Steps

After completing this exercise, proceed to `minis/35-jwt-auth-middleware`.
