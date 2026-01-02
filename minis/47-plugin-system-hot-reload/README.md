# 47-plugin-system-hot-reload

**Plugin System with Hot Reload**

Build a plugin system with dynamic loading.

## What You'll Learn

- Go plugin package
- Dynamic loading
- Interface-based plugins
- File watching for reload

## Functions to Implement

| Function | Description |
|----------|-------------|
| Load plugin | From .so file |
| Execute plugin | Call interface methods |
| Watch for changes | Reload on file change |

## Project Structure

```
47-plugin-system-hot-reload/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness
├── plugins/
│   └── example/         # Example plugin
├── internal/pluginsystemhotreload/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/47-plugin-system-hot-reload

# Build plugin
go build -buildmode=plugin -o plugins/example.so plugins/example/

# Run with plugin
go run ./cmd/app/main.go --plugin plugins/example.so

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Build a plugin (Linux/macOS only)
go build -buildmode=plugin -o plugins/example.so plugins/example/

# Load and run plugin
go run ./cmd/app/main.go --plugin plugins/example.so

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **plugin.Open**: Load shared object
2. **plugin.Lookup**: Get exported symbol
3. **Interface Cast**: Type assert to interface
4. **fsnotify**: Watch for changes

## Next Steps

After completing this exercise, proceed to `minis/48-reflection-introspection`.
