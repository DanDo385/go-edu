# Project 05: cli-todo-files - Building a Persistent CLI Application

## What Is This Project About?

Imagine you want to track your daily tasks from the command line. You need to:
- Add new tasks: `todo -add "Learn Go interfaces"`
- Mark tasks as done: `todo -toggle 1`
- List all tasks: `todo -list`
- Have everything persist between runs (survive computer restarts!)

This project teaches you how to build a **stateful CLI application** that saves data to disk. You'll learn about JSON persistence, command-line flags, and interface-based design.

## First Principles: What Makes an Application "Stateful"?

### Stateless vs Stateful

**Stateless program**: Processes input, produces output, forgets everything
```bash
echo "Hello" | wc -l    # Counts lines, then exits
# No memory of previous runs
```

**Stateful program**: Remembers data between runs
```bash
todo -add "Buy milk"     # Saves to disk
# Computer restarts...
todo -list               # Still shows "Buy milk"!
```

**The key**: Data must survive process termination. We need **persistence**.

### How Do We Persist Data?

Three common approaches:

1. **Database** (SQLite, PostgreSQL)
   - Pros: Structured queries, transactions, concurrent access
   - Cons: Requires setup, overkill for small apps

2. **Plain text file** (CSV, custom format)
   - Pros: Human-readable, easy to edit
   - Cons: Hard to parse complex structures

3. **JSON file** (our choice!)
   - Pros: Structured, easy to parse, human-readable
   - Cons: Must read/write entire file for updates

### What Is JSON Persistence?

Store data as JSON in a file:

```json
[
  {"id": 1, "text": "Learn Go", "done": false},
  {"id": 2, "text": "Build project", "done": true}
]
```

**Read**: Load JSON from file → Unmarshal to Go structs
**Modify**: Change structs in memory
**Write**: Marshal structs → Save JSON to file

## The Problem We're Solving

Build a TODO app with these operations:

| Command | Example | Effect |
|---------|---------|--------|
| Add | `todo -add "Task"` | Create new task |
| Toggle | `todo -toggle 1` | Mark task 1 as done/undone |
| List | `todo -list` | Show pending tasks |
| List all | `todo -list -all` | Show all tasks (including done) |

**Persistent storage**: All data saved to `todos.json` file

## Breaking Down the Solution (Step by Step)

### Step 1: Understanding Command-Line Flags

Flags are like function parameters for CLI programs:

```bash
todo -add "Buy milk" -file custom.json
     ↑     ↑          ↑      ↑
   flag  value      flag   value
```

In Go, we use the `flag` package:

```go
addText := flag.String("add", "", "Add a new todo")
filePath := flag.String("file", "todos.json", "Path to todo file")
flag.Parse()  // Actually parse the command-line arguments

if *addText != "" {
    // User provided -add flag
    fmt.Println("Adding:", *addText)
}
```

**Key points**:
- `flag.String()` returns a **pointer** to a string
- Must dereference with `*` to get the value
- `flag.Parse()` must be called to actually process `os.Args`

### Step 2: Designing the Data Model

Each TODO item has:
- **ID**: Unique identifier (integer)
- **Text**: The task description (string)
- **Done**: Completion status (boolean)

```go
type Item struct {
    ID   int    `json:"id"`
    Text string `json:"text"`
    Done bool   `json:"done"`
}
```

**JSON struct tags** (`json:"id"`) tell Go:
- When marshaling to JSON, use lowercase field names
- When unmarshaling, map JSON `"id"` to struct field `ID`

### Step 3: Designing the Interface

Instead of tying our code to a specific file format, we define an **interface**:

```go
type Store interface {
    Load() error
    Save() error
    Add(text string) Item
    Toggle(id int) (Item, bool)
    List(onlyPending bool) []Item
}
```

**Why an interface?**
- **Testability**: In tests, use in-memory storage instead of real files
- **Flexibility**: Could swap JSON for SQLite later without changing CLI code
- **Clear contract**: Documents exactly what a Store must do

### Step 4: Implementing FileStore

The `fileStore` struct implements our `Store` interface:

```go
type fileStore struct {
    path  string   // Where to save JSON
    items []Item   // In-memory storage
}
```

**Critical insight**: We keep data in memory (`items []Item`), only touching the file on `Load()` and `Save()`.

**Why?**
- **Performance**: Reading/parsing JSON is slow; do it once
- **Simplicity**: Work with Go structs, not JSON strings
- **Atomic updates**: Modify multiple items, save once

### Step 5: Loading from JSON

```go
func (fs *fileStore) Load() error {
    // Read entire file
    data, err := os.ReadFile(fs.path)
    if err != nil {
        return err  // File doesn't exist or can't be read
    }

    // Unmarshal JSON array to slice of Items
    if err := json.Unmarshal(data, &fs.items); err != nil {
        return fmt.Errorf("parsing JSON: %w", err)
    }

    return nil
}
```

**What happens**:
1. `os.ReadFile()` reads entire file as `[]byte`
2. `json.Unmarshal()` parses JSON, populates `fs.items` slice
3. If file doesn't exist, return error (caller handles)

### Step 6: Saving to JSON

```go
func (fs *fileStore) Save() error {
    // Marshal items to JSON (with indentation for readability)
    data, err := json.MarshalIndent(fs.items, "", "  ")
    if err != nil {
        return fmt.Errorf("marshaling JSON: %w", err)
    }

    // Write to file (0644 = owner read/write, others read)
    if err := os.WriteFile(fs.path, data, 0644); err != nil {
        return fmt.Errorf("writing file: %w", err)
    }

    return nil
}
```

**File permissions** (`0644`):
- Owner: read + write (6 = 4+2)
- Group: read only (4)
- Others: read only (4)

### Step 7: Adding Items (ID Generation)

How do we generate unique IDs?

**Naive approach**: Use current timestamp
Problem: Two items added in same second get same ID!

**Our approach**: Find max existing ID + 1

```go
func (fs *fileStore) Add(text string) Item {
    // Find maximum ID
    maxID := 0
    for _, item := range fs.items {
        if item.ID > maxID {
            maxID = item.ID
        }
    }

    // Create new item
    newItem := Item{
        ID:   maxID + 1,
        Text: text,
        Done: false,
    }

    // Append to slice
    fs.items = append(fs.items, newItem)

    return newItem
}
```

**Production alternatives**:
- UUID: Globally unique, works in distributed systems
- Database auto-increment: Let DB generate IDs
- Snowflake IDs: Timestamp + machine ID + sequence number

### Step 8: Toggling Items

```go
func (fs *fileStore) Toggle(id int) (Item, bool) {
    // Search for item by ID
    for i := range fs.items {
        if fs.items[i].ID == id {
            // Toggle done status
            fs.items[i].Done = !fs.items[i].Done
            return fs.items[i], true
        }
    }

    // Not found
    return Item{}, false
}
```

**Why return (Item, bool)?**
- Alternative: return `error` if not found
- But "not found" isn't exceptional—it's expected user behavior
- Pattern: `value, ok := map[key]` (Go's map lookup pattern)

### Step 9: Listing Items

```go
func (fs *fileStore) List(onlyPending bool) []Item {
    if !onlyPending {
        return fs.items  // All items
    }

    // Filter: only items where Done == false
    var result []Item
    for _, item := range fs.items {
        if !item.Done {
            result = append(result, item)
        }
    }

    return result
}
```

**Memory note**: We return `fs.items` directly (not a copy). This shares the underlying array—efficient but be careful not to modify the returned slice!

## The Complete Flow (User Perspective)

```bash
# First run: File doesn't exist yet
$ todo -add "Learn Go"
Added: [1] Learn Go

# Behind the scenes:
# 1. Load() fails (file doesn't exist) - that's OK
# 2. Add("Learn Go") creates Item{ID:1, Text:"Learn Go", Done:false}
# 3. Save() writes todos.json

# Second run: File exists
$ todo -list
[1] [ ] Learn Go

# Behind the scenes:
# 1. Load() succeeds, reads JSON, populates items
# 2. List(onlyPending=true) filters items
# 3. Display results

# Toggle an item
$ todo -toggle 1
Toggled: [1] Learn Go (done)

# Behind the scenes:
# 1. Load() reads existing data
# 2. Toggle(1) finds item 1, sets Done=true
# 3. Save() writes updated JSON

# List again
$ todo -list
(empty - no pending items)

$ todo -list -all
[1] [✓] Learn Go
```

## Key Concepts Explained

### Why Pointer Receivers?

Methods that modify struct fields need pointer receivers:

```go
func (fs *fileStore) Add(text string) Item {
    fs.items = append(fs.items, newItem)  // Modifies fs.items
    // ...
}
```

If we used value receiver `(fs fileStore)`, we'd modify a **copy** of the struct!

### Why Return Interface from Constructor?

```go
func NewFileStore(path string) Store {  // Returns interface
    return &fileStore{path: path, items: []Item{}}
}
```

**Benefits**:
- Caller doesn't know implementation details
- Easy to swap implementations (fileStore → dbStore)
- Encourages thinking about behavior, not structure

### Testing with Temporary Files

Tests use `t.TempDir()` to create temporary directories:

```go
func TestFileStore(t *testing.T) {
    tmpDir := t.TempDir()  // Auto-cleaned after test
    storePath := filepath.Join(tmpDir, "test.json")

    store := NewFileStore(storePath)
    store.Add("Test task")
    store.Save()

    // Verify file exists and contains correct JSON
}
```

## Common Patterns You're Learning

### Pattern 1: Flag Parsing
```go
add := flag.String("add", "", "Add task")
flag.Parse()

if *add != "" {
    // Use *add
}
```

### Pattern 2: JSON Persistence
```go
// Load
data, _ := os.ReadFile(path)
json.Unmarshal(data, &items)

// Save
data, _ := json.MarshalIndent(items, "", "  ")
os.WriteFile(path, data, 0644)
```

### Pattern 3: Interface-Based Design
```go
type Store interface {
    Load() error
    Save() error
    // ...
}

func NewFileStore(path string) Store {
    return &fileStore{path: path}
}
```

## Real-World Applications

1. **CLI Tools**: Git, npm, docker all use persistent config files
2. **Personal Productivity**: Task managers, note-taking apps
3. **System Configuration**: Storing application settings
4. **Data Migration Tools**: Import/export data as JSON
5. **Logging**: Structured log storage and rotation

## How to Run

```bash
# Prepare the project
cd minis/05-cli-todo-files/exercise
mv solution.go solution.go.reference

# Run tests
go test -v

# Build and use the CLI
cd ../cmd/cli-todo-files

# Add tasks
go run . -add "Learn interfaces"
go run . -add "Build CLI tool"

# List tasks
go run . -list

# Toggle a task
go run . -toggle 1

# List all (including done)
go run . -list -all

# Check the JSON file
cat todos.json
```

## Common Mistakes to Avoid

1. **Forgetting to call `flag.Parse()`**: Flags won't work without it!
2. **Not checking file existence**: Handle `os.IsNotExist(err)` gracefully
3. **Modifying return values**: Don't modify slices returned from `List()`
4. **Not saving after changes**: Add → Toggle → Save (don't forget the Save!)
5. **Concurrent access**: Our implementation isn't thread-safe (use locks for concurrent writes)

## Stretch Goals

1. **Add delete command**: `todo -delete 1`
   - Remove item from slice: `items = append(items[:i], items[i+1:]...)`

2. **Add edit command**: `todo -edit 1 "New text"`
   - Find item, update Text field

3. **Add priorities**: High/Medium/Low priority field
   - Sort by priority in List()

4. **Add due dates**: Parse dates with `time.Parse()`
   - Filter overdue items

5. **Atomic file writes**: Write to temp file, then rename
   - Prevents corruption if program crashes mid-write
   ```go
   tmpFile := path + ".tmp"
   os.WriteFile(tmpFile, data, 0644)
   os.Rename(tmpFile, path)  // Atomic on most filesystems
   ```
