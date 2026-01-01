# 23-mempool

## What Is This Project About?
[2-3 paragraphs explaining the core concept in accessible language]

## Why Is This Important?
[Explain why this concept matters in real-world software development]

## Real-World Problems This Solves
[Provide 2-3 concrete examples of real-world problems this concept solves]

## Key Concepts You'll Learn
- Concept 1: [explanation]
- Concept 2: [explanation]
- Concept 3: [explanation]

## Prerequisites
[List any prerequisites]

## Project Structure
```
23-mempool/
├── cmd/
│   ├── app/           # Application entry point (CLI arguments)
│   └── dev/           # Debug harness (fixed inputs)
├── internal/
│   └── mempool/      # Exercise implementation
│       ├── exercise.go
│       ├── exercise_test.go
│       ├── solution.reference.go
│       └── solution_no_err.reference.go
└── .vscode/
    └── launch.json    # Debug configurations
```

## How to Run

### Using cmd/app/main.go (CLI Arguments)
```bash
# Basic usage
go run ./cmd/app/main.go [arguments]

# Examples:
go run ./cmd/app/main.go "hello world"
go run ./cmd/app/main.go https://eth.llamarpc.com
```

### Using cmd/dev/main.go (Debug Harness)
```bash
# Run with fixed test inputs
go run ./cmd/dev/main.go

# Or use VS Code debugger (F5) with "Debug: cmd/dev" configuration
```

## How to Debug
1. Set breakpoints at "// BREAKPOINT:" comments
2. Use VS Code debugger (F5) and select appropriate configuration:
   - "Debug: cmd/app" - Debug with CLI arguments
   - "Debug: cmd/dev" - Debug with fixed inputs (recommended for learning)
   - "Test: Run All Tests" - Debug tests
3. Step through code using F10 (Step Over) and F11 (Step Into)
4. Watch variables in the Variables panel

## Testing
```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -v -run TestFunctionName ./...

# Run with reference implementation
go test -tags=reference -v ./...
```

## Exercises
[If applicable, list the specific exercises/functions to implement]

## Additional Resources
[Links to relevant documentation, articles, or related projects]
