# 18: Goroutines 1M Demo

## What Is This Project About?

This module demonstrates Go's ability to spawn millions of goroutines, showing the lightweight nature of Go's concurrency primitives.

## Key Concepts

- **Goroutine overhead**: ~2KB stack size
- **Scheduler**: M:N scheduling model
- **Memory usage**: Tracking goroutine costs
- **Practical limits**: When too many is too many

## How to Run

```bash
go run ./cmd/dev/main.go
```
