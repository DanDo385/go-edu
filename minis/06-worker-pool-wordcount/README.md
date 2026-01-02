# 06: Worker Pool Wordcount

## What Is This Project About?

This module teaches you how to build a worker pool for concurrent processing. You'll implement a word counter that processes multiple URLs concurrently.

## Why Is This Important?

Worker pools are fundamental for:
- Parallel processing
- Resource management
- Throughput optimization
- Scalable services

## Key Concepts You'll Learn

- **Goroutines**: Lightweight concurrent execution
- **Channels**: Communication between workers
- **Worker pools**: Managing concurrent workers
- **Fan-out/Fan-in**: Distributing and collecting work

## Prerequisites

- Completion of `minis/05-cli-todo-files`
- Basic understanding of concurrency

## How to Run

```bash
go run ./cmd/dev/main.go
```

## Testing

```bash
go test -v ./...
```
