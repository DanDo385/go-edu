# Project 08: http-client-retries

## What You're Building

A generic HTTP client with automatic retries, exponential backoff, jitter, and context-aware timeouts. This project demonstrates resilient network programming patterns.

## Concepts Covered

- `net/http` client configuration
- Generic functions with type parameters
- Exponential backoff with jitter
- `context.Context` for timeouts and cancellation
- Error classification (retryable vs fatal)
- `httptest` for deterministic testing

## How to Run

```bash
# Run tests
go test ./minis/08-http-client-retries/...

# Run with verbose output
go test -v ./minis/08-http-client-retries/...
```

## Solution Explanation

### Retry Strategy

**Exponential Backoff**: delay = baseDelay * (2^attempt)
- Attempt 1: 100ms
- Attempt 2: 200ms
- Attempt 3: 400ms

**Jitter**: Add randomness to prevent thundering herd
- ±20% variation: prevents all clients retrying simultaneously

## Where Go Shines

**Go vs other languages:**
- Built-in HTTP client with sensible defaults
- Context propagation is standardized
- Type-safe generics for JSON decoding

## Stretch Goals

1. Add circuit breaker pattern
2. Implement request hedging
3. Add metrics/tracing integration
