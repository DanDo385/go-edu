# 11-slices-internals-capacity-growth

**Slice Internals**

Understand slice memory layout, capacity growth, and aliasing.

## What You'll Learn

- Slice header structure (ptr, len, cap)
- Capacity growth algorithm
- Slice aliasing pitfalls
- Memory efficiency

## Functions to Implement

| Function | Description |
|----------|-------------|
| Demonstrate slice internals | Show capacity growth and aliasing |

## Project Structure

```
11-slices-internals-capacity-growth/
├── cmd/
│   ├── app/main.go      # CLI demonstration
│   └── dev/main.go      # Debug harness
├── internal/slicesinternalscapacitygrowth/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/11-slices-internals-capacity-growth

# Run demonstration
go run ./cmd/app/main.go

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Run demo
go run ./cmd/app/main.go

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Slice Header**: struct { ptr *T; len, cap int }
2. **Growth**: ~2x until 1024, then ~1.25x
3. **Aliasing**: Multiple slices sharing backing array
4. **make([]T, len, cap)**: Pre-allocate capacity

## Next Steps

After completing this exercise, proceed to `minis/12-pointers-zero-values-nil-gotchas`.
