# 12-pointers-zero-values-nil-gotchas

**Pointers and Nil Gotchas**

Master pointers, zero values, and common nil pitfalls.

## What You'll Learn

- Pointer vs value semantics
- Zero values for all types
- nil interface gotcha
- Safe nil handling

## Functions to Implement

| Function | Description |
|----------|-------------|
| Demonstrate pointer patterns | Show common nil pitfalls and fixes |

## Project Structure

```
12-pointers-zero-values-nil-gotchas/
├── cmd/
│   ├── app/main.go      # CLI demonstration
│   └── dev/main.go      # Debug harness
├── internal/pointerszerovaluesnilgotchas/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/12-pointers-zero-values-nil-gotchas

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

1. **Zero Values**: 0, "", false, nil
2. **nil Interface**: Interface with nil concrete value is NOT nil
3. **Pointer Receivers**: Can be called on nil
4. **Safe Dereference**: Check nil before accessing

## Next Steps

After completing this exercise, proceed to `minis/13-interfaces-duck-typing`.
