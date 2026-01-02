# 46-generics-map-reduce

**Generic Map/Reduce**

Implement generic functional programming utilities.

## What You'll Learn

- Go generics for collections
- Map, Filter, Reduce patterns
- Type constraints
- Functional composition

## Functions to Implement

| Function | Description |
|----------|-------------|
| `Map[T, U]([]T, func(T) U) []U` | Transform elements |
| `Filter[T]([]T, func(T) bool) []T` | Select elements |
| `Reduce[T, U]([]T, U, func(U, T) U) U` | Aggregate elements |

## Project Structure

```
46-generics-map-reduce/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness
├── internal/genericsmapreduce/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/46-generics-map-reduce

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

1. **Type Parameters**: `[T any]`
2. **Constraints**: `comparable`, custom interfaces
3. **Map**: Transform each element
4. **Reduce**: Fold into single value

## Next Steps

After completing this exercise, proceed to `minis/47-plugin-system-hot-reload`.
