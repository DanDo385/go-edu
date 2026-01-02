# 28-pprof-cpu-mem-benchmarks

**Profiling with pprof**

Profile Go programs for CPU and memory optimization.

## What You'll Learn

- CPU profiling
- Memory profiling
- Benchmark profiling
- Flame graphs

## Functions to Implement

| Function | Description |
|----------|-------------|
| Profile code performance | CPU and memory analysis |

## Project Structure

```
28-pprof-cpu-mem-benchmarks/
├── cmd/
│   ├── app/main.go      # CLI demonstration
│   └── dev/main.go      # Debug harness
├── internal/pprofcpumembenchmarks/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests (with benchmarks)
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/28-pprof-cpu-mem-benchmarks

# Run with CPU profile
go run ./cmd/app/main.go --cpuprofile cpu.prof

# Run with memory profile
go run ./cmd/app/main.go --memprofile mem.prof

# Analyze profile
go tool pprof cpu.prof
```

## Quick Copy & Paste

```bash
# Generate CPU profile
go run ./cmd/app/main.go --cpuprofile cpu.prof

# Generate memory profile
go run ./cmd/app/main.go --memprofile mem.prof

# Analyze with pprof
go tool pprof -http=:8080 cpu.prof

# Benchmark with profiling
go test -bench=. -cpuprofile cpu.prof -memprofile mem.prof ./...

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **CPU Profile**: Where time is spent
2. **Memory Profile**: Where allocations happen
3. **go tool pprof**: Analysis tool
4. **Flame Graphs**: Visual call stack

## Next Steps

After completing this exercise, proceed to `minis/29-escape-analysis-inlining`.
