# 06-worker-pool-wordcount

**Concurrent Word Counter**

Fetch multiple URLs concurrently and count word frequencies.

## What You'll Learn

- Worker pool pattern
- Bounded concurrency
- Error propagation from workers
- Context cancellation

## Functions to Implement

| Function | Description |
|----------|-------------|
| `WordCount(ctx, urls, workers) (map[string]int, error)` | Count words with manual workers |
| `WordCountWithErrGroup(ctx, urls, workers)` | Using errgroup package |
| `fetchAndCount(ctx, url)` | Fetch URL and count words |
| `tokenizeAndCount(text)` | Split text and count |

## Project Structure

```
06-worker-pool-wordcount/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/workerpoolwordcount/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd minis/06-worker-pool-wordcount

# Count words from URLs
go run ./cmd/app/main.go --workers 3 https://example.com https://golang.org
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
| `--workers` | Number of concurrent workers (default: 5) |
| `URLS...` | One or more URLs to fetch |

## Quick Copy & Paste

```bash
# Fetch and count from multiple URLs
go run ./cmd/app/main.go --workers 3 https://example.com https://golang.org

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Worker Pool**: Fixed goroutine count
2. **Channels**: Work distribution and result collection
3. **errgroup.Group**: Coordinated error handling
4. **Context Cancellation**: Stop all workers on error

## Next Steps

After completing this exercise, proceed to `minis/07-generic-lru-cache`.
