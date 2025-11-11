# Project 05: cli-todo-files

## What You're Building

A CLI TODO list application with JSON file persistence. Users can add tasks, mark them as done, and list all items. This project demonstrates interface-based design, JSON persistence, and command-line flag parsing.

## Concepts Covered

- `flag` package for CLI argument parsing
- Interfaces for abstraction and testability
- JSON marshaling/unmarshaling with `encoding/json`
- File I/O with `os.ReadFile` and `os.WriteFile`
- Method receivers (pointer vs value)
- Stateful struct design
- Testing with temporary files (`t.TempDir()`)

## How to Run

```bash
# Run with different commands
go run ./minis/05-cli-todo-files/cmd/cli-todo-files -add "Learn Go interfaces"
go run ./minis/05-cli-todo-files/cmd/cli-todo-files -list
go run ./minis/05-cli-todo-files/cmd/cli-todo-files -toggle 1
go run ./minis/05-cli-todo-files/cmd/cli-todo-files -list -all

# Run tests
go test ./minis/05-cli-todo-files/...
```

## Solution Explanation

### Architecture

**Store Interface**: Defines CRUD operations (Load, Save, Add, Toggle, List). This abstraction allows:
- Testing without real files (mock implementations)
- Future backends (SQLite, remote API) without changing CLI code
- Clear separation of concerns

**FileStore Implementation**: Concrete implementation backed by a JSON file. Uses:
- In-memory slice for fast operations
- Lazy loading (Load() called once at startup)
- Atomic writes (write to temp file, then rename)

### ID Generation

Items get unique IDs by finding the max existing ID and incrementing. For production, consider:
- UUIDs for distributed systems
- Database auto-increment
- Snowflake IDs for time-sortable IDs

## Where Go Shines

**Go vs Python:**
- Python: `argparse` is powerful but verbose; `json.dump()` is simpler
- Go: `flag` package is built-in and type-safe
- Go's interfaces enable testing without mocking frameworks

**Go vs Node.js:**
- JS: Requires libraries (`commander`, `fs-extra`) for CLI + file I/O
- Go: Everything in stdlib; single binary deployment

**Go vs Rust:**
- Rust: `clap` is excellent but adds compile time; `serde_json` is similar
- Go: Faster iteration (quick compile times); simpler for CRUD apps

## Stretch Goals

1. **Add delete command**: `-delete 1` removes an item
2. **Add edit command**: `-edit 1 "New text"` updates an item
3. **Add priorities**: Items have priority levels (low/medium/high)
4. **Add due dates**: `-add "Task" -due "2024-12-31"`
5. **Export formats**: `-export json|csv|markdown`
