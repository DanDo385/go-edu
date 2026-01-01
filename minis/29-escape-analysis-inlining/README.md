# 29: Escape Analysis and Inlining

## What Is This Project About?

This module teaches you about Go compiler optimizations.

## Key Concepts

- **Escape analysis**: Stack vs heap allocation
- **Inlining**: Function call elimination
- **Compiler flags**: -gcflags for analysis
- **Performance impact**: Understanding allocations

## How to Run

```bash
go build -gcflags="-m" ./...
go run ./cmd/dev/main.go
```
