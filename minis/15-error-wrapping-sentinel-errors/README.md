# 15-error-wrapping-sentinel-errors

**Error Handling Patterns**

Master Go error handling with wrapping and sentinel errors.

## What You'll Learn

- Error wrapping with %w
- errors.Is and errors.As
- Sentinel errors
- Custom error types

## Functions to Implement

| Function | Description |
|----------|-------------|
| Demonstrate error patterns | Wrapping and inspection |

## Project Structure

```
15-error-wrapping-sentinel-errors/
├── cmd/
│   ├── app/main.go      # CLI demonstration
│   └── dev/main.go      # Debug harness
├── internal/errorwrappingsentinelerrors/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/15-error-wrapping-sentinel-errors

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

1. **Error Wrapping**: `fmt.Errorf("context: %w", err)`
2. **errors.Is**: Check for specific error in chain
3. **errors.As**: Extract typed error from chain
4. **Sentinel Errors**: `var ErrNotFound = errors.New(...)`

## Next Steps

After completing this exercise, proceed to `minis/16-context-cancellation-timeouts`.
