# 13-interfaces-duck-typing

**Interfaces and Duck Typing**

Understand Go's implicit interface implementation.

## What You'll Learn

- Implicit interface satisfaction
- Interface composition
- The empty interface (any)
- Type assertions and switches

## Functions to Implement

| Function | Description |
|----------|-------------|
| Demonstrate interface patterns | Duck typing in action |

## Project Structure

```
13-interfaces-duck-typing/
├── cmd/
│   ├── app/main.go      # CLI demonstration
│   └── dev/main.go      # Debug harness
├── internal/interfacesducktyping/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/13-interfaces-duck-typing

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

1. **Duck Typing**: "If it walks like a duck..."
2. **No `implements`**: Implicit satisfaction
3. **Interface Composition**: Embedding interfaces
4. **Type Switch**: `switch v := x.(type)`

## Next Steps

After completing this exercise, proceed to `minis/14-methods-value-vs-pointer-receivers`.
