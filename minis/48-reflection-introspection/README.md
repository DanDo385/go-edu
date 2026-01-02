# 48-reflection-introspection

**Reflection and Introspection**

Use Go's reflect package for runtime introspection.

## What You'll Learn

- reflect.TypeOf and ValueOf
- Struct field introspection
- Dynamic method calls
- Tag parsing

## Functions to Implement

| Function | Description |
|----------|-------------|
| Introspect struct | List fields and types |
| Get/set fields | Dynamic field access |
| Call methods | Runtime method invocation |

## Project Structure

```
48-reflection-introspection/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness
├── internal/reflectionintrospection/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/48-reflection-introspection

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

1. **reflect.Type**: Type metadata
2. **reflect.Value**: Runtime value wrapper
3. **Field Tags**: `reflect.StructTag.Get()`
4. **SetValue**: Must be addressable

## Next Steps

After completing this exercise, proceed to `minis/49-state-machine-pattern`.
