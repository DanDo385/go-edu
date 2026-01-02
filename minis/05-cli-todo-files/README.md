# minis/05-cli-todo-files

## Problem

Problem: Build a persistent TODO list with JSON file storage

We need to implement:
1. CRUD operations (Create, Read, Update, Delete-ish with Toggle)
2. JSON serialization/deserialization for persistence
3. CLI interface with flag parsing
4. Atomic file writes (no partial corruption)

Constraints:
- Items have unique IDs (auto-incrementing)
- JSON file stores all items as an array
- Toggle operation is idempotent
- List can filter by completion status

Time/Space Complexity:
- Load/Save: O(n) where n = number of items (JSON marshal/unmarshal)
- Add: O(n) to find max ID, O(1) to append
- Toggle: O(n) to find item by ID
- List: O(n) to filter items

Why Go is well-suited:
- `flag` package for CLI parsing (built-in, type-safe)
- JSON marshal/unmarshal with struct tags (no external dependencies)
- Interfaces enable testing without real files
- Pointer receivers for mutable state (clear semantics)

DEBUGGING THIS FILE:
==================
This solution is instrumented with extensive debugging comments to teach you
how to use Go's debugger (dlv) and VS Code's debugging features.

Key debugging concepts covered:
1. Setting breakpoints at critical CRUD operations
2. Watching interface implementations and method calls
3. Using F10 (Step Over) vs F11 (Step Into) effectively
4. Inspecting JSON marshaling/unmarshaling
5. Using the Debug Console to evaluate expressions
6. Understanding pointer receivers and mutability

## Quickstart

```bash
cd minis/05-cli-todo-files
go test ./...
```

## CLI (cmd/app)

This project includes a small CLI wrapper around the exercise package.

### Flags

- **`-list`**: list available exported functions
- **`-fn`**: function name to run
- **`-in`**: string input (for `func(string) ...`)
- **`-n`**: int input (for `func(int) ...`)
- **`-f`**: float64 input (for `func(float64) ...`)
- **`-b`**: bool input (for `func(bool) ...`)
- **`-file`** / **`-stdin`**: input sources for `func(io.Reader) ...`

### Usage

```bash
go run ./cmd/app -h
```

### Copy/paste examples

```bash
go run ./cmd/app -list
go run ./cmd/app -fn "NewFileStore" -in "Hello, 世界 👋"
go run ./cmd/dev
```

## Debug harness (cmd/dev)

`cmd/dev` runs the same CLI with preselected example arguments, so you can just hit Run/Debug and see what the project can do.

```bash
go run ./cmd/dev
```

## Files

- `internal/clitodofiles/exercise.go`: implement the TODOs here
- `internal/clitodofiles/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
