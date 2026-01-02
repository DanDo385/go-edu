# 47 Plugin System Hot Reload

## Problem

Dynamic plugin loading and execution

## Constraints

- Plugins must be built with the same Go version as the host

## Complexity

- plugin.Open(): O(n) where n = size of .so file (disk I/O + linking)

## Project Structure

```
47-plugin-system-hot-reload/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application with custom arguments
│   └── dev/
│       └── main.go          # Debug harness with auto-demo
├── internal/
│   └── pluginsystemhotreload/
│       ├── exercise.go      # YOUR CODE GOES HERE
│       ├── exercise_test.go # Test cases
│       └── solution.reference.go  # Reference implementation
└── README.md
```

## Getting Started

1. Navigate to this directory:
   ```bash
   cd minis/47-plugin-system-hot-reload
   ```

2. Open the exercise file:
   ```bash
   code internal/pluginsystemhotreload/exercise.go
   ```

3. Implement the TODO functions

4. Run tests:
   ```bash
   go test -v ./...
   ```

## CLI Usage

### Using cmd/app/main.go

This project's CLI application accepts the following arguments:

```bash
go run ./cmd/app/main.go [arguments]
```

Run without arguments to see usage information:

```bash
go run ./cmd/app/main.go
```

### Using cmd/dev/main.go

The `cmd/dev/main.go` file automatically demonstrates the project's capabilities
by running through different scenarios with pre-configured inputs.

**Run the demo:**

```bash
go run ./cmd/dev/main.go
```

**Debug with VS Code:**

1. Open `cmd/dev/main.go`
2. Set breakpoints at `// BREAKPOINT:` comments
3. Press F5 and select "Debug cmd/dev (Debug Harness)"

## Testing

Run all tests:

```bash
go test -v ./...
```

Run specific test:

```bash
go test -v -run TestFunctionName ./...
```

## Reference Solution

If you get stuck, check `internal/pluginsystemhotreload/solution.reference.go` for a complete
implementation with detailed explanations.

