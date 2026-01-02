# 21: Race Detection Demo

## What Is This Project About?

This module teaches you about Go's race detector and how to identify and fix data races.

## Key Concepts

- **Data races**: Concurrent access without synchronization
- **Race detector**: go run -race
- **Fixing races**: Mutexes, channels, atomic operations
- **Best practices**: Avoiding race conditions

## How to Run

```bash
go run -race ./cmd/dev/main.go
```
