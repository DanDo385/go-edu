# 30-build-tags-conditional-compilation

**Build Tags**

Use build tags for conditional compilation.

## What You'll Learn

- Build constraint syntax
- Platform-specific code
- Feature flags
- Build tags vs file naming

## Functions to Implement

| Function | Description |
|----------|-------------|
| Platform-specific implementations | Conditional compilation |

## Project Structure

```
30-build-tags-conditional-compilation/
├── cmd/
│   ├── app/main.go      # CLI demonstration
│   └── dev/main.go      # Debug harness
├── internal/buildtagsconditionalcompilation/
│   ├── exercise.go         # Default implementation
│   ├── exercise_linux.go   # Linux-specific
│   ├── exercise_darwin.go  # macOS-specific
│   ├── exercise_test.go    # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/30-build-tags-conditional-compilation

# Build with tag
go build -tags=feature_x ./...

# Build for specific OS
GOOS=linux go build ./...
```

## Quick Copy & Paste

```bash
# Default build
go run ./cmd/app/main.go

# Build with custom tag
go build -tags=debug ./cmd/app/

# Cross-compile
GOOS=linux GOARCH=amd64 go build ./cmd/app/

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **//go:build tag**: New syntax (Go 1.17+)
2. **// +build tag**: Old syntax (still works)
3. **_linux.go**: File name suffix convention
4. **GOOS/GOARCH**: Cross-compilation

## Next Steps

After completing this exercise, proceed to `minis/31-static-file-server`.
