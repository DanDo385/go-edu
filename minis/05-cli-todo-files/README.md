# 05-cli-todo-files

**Persistent TODO List**

Build a persistent TODO list with JSON file storage.

## What You'll Learn

- File-based persistence
- JSON encoding/decoding
- Interface-based design
- CRUD operations

## Functions to Implement

| Function | Description |
|----------|-------------|
| `NewFileStore(path string) Store` | Create new file-backed store |
| `Load() error` | Load items from file |
| `Save() error` | Save items to file |
| `Add(text string) Item` | Add new item |
| `Toggle(id int) (Item, bool)` | Toggle item completion |
| `List(onlyPending bool) []Item` | List items |

## Project Structure

```
05-cli-todo-files/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness with fixed inputs
├── internal/clitodofiles/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

### Run the CLI Application

```bash
cd minis/05-cli-todo-files

# Add a todo
go run ./cmd/app/main.go add "Buy groceries"

# List all todos
go run ./cmd/app/main.go list

# List only pending
go run ./cmd/app/main.go list --pending

# Toggle completion
go run ./cmd/app/main.go toggle 1
```

### Run the Debug Harness

```bash
go run ./cmd/dev/main.go
```

### Run Tests

```bash
go test -v ./...
```

## CLI Arguments

| Command | Arguments | Description |
|---------|-----------|-------------|
| `add` | `TEXT` | Add new todo item |
| `list` | `--pending` | List todos (optionally only pending) |
| `toggle` | `ID` | Toggle item completion |

## Quick Copy & Paste

```bash
# Add items
go run ./cmd/app/main.go add "Learn Go"
go run ./cmd/app/main.go add "Build something"

# List all
go run ./cmd/app/main.go list

# Mark as done
go run ./cmd/app/main.go toggle 1

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **json.Marshal/Unmarshal**: Serialization
2. **os.ReadFile/WriteFile**: Simple file I/O
3. **Interface Design**: Abstraction over storage

## Next Steps

After completing this exercise, proceed to `minis/06-worker-pool-wordcount`.
