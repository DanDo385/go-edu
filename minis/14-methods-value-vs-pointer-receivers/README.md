# 14-methods-value-vs-pointer-receivers

**Value vs Pointer Receivers**

Understand when to use value vs pointer method receivers.

## What You'll Learn

- Value receiver semantics
- Pointer receiver semantics
- Method set rules
- Consistency guidelines

## Functions to Implement

| Function | Description |
|----------|-------------|
| Demonstrate receiver patterns | Value vs pointer behavior |

## Project Structure

```
14-methods-value-vs-pointer-receivers/
├── cmd/
│   ├── app/main.go      # CLI demonstration
│   └── dev/main.go      # Debug harness
├── internal/methodsvaluevspointerreceivers/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/14-methods-value-vs-pointer-receivers

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

1. **Value Receiver**: Gets copy, can't modify original
2. **Pointer Receiver**: Can modify, avoid large copies
3. **Method Set**: Value type doesn't include pointer methods
4. **Consistency**: Use same receiver type for all methods

## Next Steps

After completing this exercise, proceed to `minis/15-error-wrapping-sentinel-errors`.
