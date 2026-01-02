# 28: pprof CPU and Memory Benchmarks

## What Is This Project About?

This module teaches you how to profile Go applications using pprof.

## Key Concepts

- **CPU profiling**: Finding hot spots
- **Memory profiling**: Allocation analysis
- **Benchmarks**: Performance testing
- **Visualization**: Flame graphs

## How to Run

```bash
go run ./cmd/dev/main.go
go test -bench=. -cpuprofile=cpu.prof
```
