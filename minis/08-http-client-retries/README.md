# 08: HTTP Client with Retries

## What Is This Project About?

This module teaches you how to build a robust HTTP client with automatic retry logic. You'll learn exponential backoff, error handling, and resilient network programming.

## Why Is This Important?

Robust HTTP clients are essential for:
- Microservice communication
- API integrations
- Handling transient failures
- Production reliability

## Key Concepts You'll Learn

- **HTTP client**: Making requests with net/http
- **Retry logic**: Automatic retry on failure
- **Exponential backoff**: Increasing delays between retries
- **Error classification**: Retryable vs permanent errors

## Prerequisites

- Completion of previous minis projects

## How to Run

```bash
go run ./cmd/app/main.go https://api.example.com 3
```

## Testing

```bash
go test -v ./...
```
