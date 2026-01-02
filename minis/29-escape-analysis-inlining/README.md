# 29-escape-analysis-inlining

**Escape Analysis & Inlining**

Understand Go compiler optimizations.

## What You'll Learn

- Escape analysis (stack vs heap)
- Function inlining
- Compiler flags for inspection
- Optimization patterns

## Functions to Implement

| Function | Description |
|----------|-------------|
| Demonstrate compiler optimizations | Escape and inlining analysis |

## Project Structure

```
29-escape-analysis-inlining/
├── cmd/
│   ├── app/main.go      # CLI demonstration
│   └── dev/main.go      # Debug harness
├── internal/escapeanalysisinlining/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/29-escape-analysis-inlining

# Build with escape analysis output
go build -gcflags="-m" ./...

# Build with more detail
go build -gcflags="-m -m" ./...
```

## Quick Copy & Paste

```bash
# See escape analysis
go build -gcflags="-m" ./cmd/app/

# See inlining decisions
go build -gcflags="-m -m" ./cmd/app/

# Run demo
go run ./cmd/app/main.go

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Escape Analysis**: Decides stack vs heap
2. **-gcflags="-m"**: Show escape decisions
3. **Inlining**: Small functions are inlined
4. **//go:noinline**: Prevent inlining

## Next Steps

After completing this exercise, proceed to `minis/30-build-tags-conditional-compilation`.
