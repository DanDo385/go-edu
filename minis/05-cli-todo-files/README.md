# 05: Cli Todo Files

## Core Concepts

- The concrete problem in Cli Todo Files and the correctness invariants it depends on.
- How values, pointers, slices, maps, or channels behave in this module's runtime path.
- Why this lesson's implementation pattern is the right next step in the learning arc.

## CS Connection

- Data representation and state transitions: what is copied, what is shared, and what can race.
- API boundaries: where we validate, where we propagate errors, and where we normalize output.
- Algorithmic tradeoffs in this lesson (latency, throughput, memory, and complexity).

## End-State Understanding

- Diagnose and implement Cli Todo Files patterns without relying on hidden framework behavior.
- Explain memory and concurrency implications of the final implementation choices.
- Compare learner code and reference code by invariants, not only by syntax.

## Why This Lesson Exists Here

Problem statement:
This lesson turns the previous module's concepts into a reusable engineering pattern for cli todo files.

At this point in the arc:
Lesson 05 introduces a sharper systems concern so later modules can assume this mental model is stable.

## Step-by-Step Build Path

### Step 1: Problem This Step Solves
Define the smallest valid behavior and reject invalid input or impossible state early.

### Step 2: Why This Approach
Pick a direct design that keeps control flow and data flow visible for debugging and testing.

### Step 3: Memory / Pointer Impact
Call out where data is copied versus aliased, and where mutable shared state needs synchronization.

### Step 4: What Changed
Produce a stable result shape and explicit error behavior that downstream code can rely on.

## Pointer and Indirection

- Explain * and & in this module when they appear in code or docs.
- Show memory-before and memory-after when data ownership changes.
- Clarify common misconceptions: Go stays pass-by-value even when pointer values are copied.
- Primer link: docs/MEMORY_POINTERS_PRIMER.md

## Verify


a) learner path


go test -v ./...


b) reference path


go test -tags=reference -v ./...


This is a major milestone. You will build your first **complete, interactive application**, assembling all the skills you've learned so far. The central challenge is **persistence**: making your application remember data between runs by saving it to a JSON file.

## Core Concepts

- Value semantics in Go: what gets copied at function calls and what can still alias shared state.
- Ownership boundaries for mutation, especially when multiple code paths touch the same logical data.

## CS Connection

- Memory layout drives behavior: variables store values, and some values are addresses into other storage.
- Go is pass-by-value, including pointers; copying a pointer value copies an address, not the pointee.
- Correctness depends on understanding copying versus aliasing (`*T`, slices, maps, and channels) and enforcing synchronization when concurrent access exists.

## End-State Understanding

- Implement the exercise with explicit reasoning about correctness, edge cases, and error paths.

## What You'll Learn

- How to build an interactive **Command-Line Interface (CLI)** application.
- How to parse command-line arguments using the `flag` package.
- How to implement **persistence** by saving and loading application state to a file.
- How to use **interfaces** to design a clean, testable application architecture.
- Why methods that modify a struct are defined on a **pointer receiver**.

## Application Architecture: Separation of Concerns

Good software design involves separating your code into layers. This makes it easier to understand, test, and maintain.

```
+---------------------+      +---------------------+      +----------------------+
|     main.go         |----->|   Store Interface   |----->|   fileStore struct   |
| (The CLI Layer)     |      |    (The Contract)   |      | (The Data Layer)     |
+---------------------+      +---------------------+      +----------------------+
```

1.  **CLI Layer (`cmd/app/main.go`):** Knows what the user wants (e.g., "add an item"). It knows *nothing* about files or JSON. It only talks to the `Store` interface.
2.  **The Contract (`Store` interface):** Defines *what* can be done (e.g., `Add`, `List`), but not *how*.
3.  **Data Layer (`fileStore` struct):** Knows *how* to handle the data (e.g., read/write a JSON file). It knows *nothing* about the CLI.

This is powerful because we can easily swap out the `fileStore` for a different implementation (like one that saves to a database, or a fake one for testing) without changing the CLI layer at all.

### 🚨 Deep Dive: Pointer Receivers on Methods

In `exercise.go`, you'll see that all the methods are defined on `*fileStore`, not `fileStore`.
```go
func (fs *fileStore) Load() error { ... }
func (fs *fileStore) Add(text string) Item { ... }
```
This is the same concept as the pointer receiver in the last lesson, but it's even more important here. The `fileStore` struct holds the slice of todo items (`items []Item`). All of the methods (`Load`, `Save`, `Add`, `Toggle`) need to **modify this slice**.

- If the methods were defined on `(fs fileStore)`, `fs` would be a **copy** of the original `fileStore`.
- When `Add` appended an item to `fs.items`, it would be modifying the slice *within the copy*.
- The original `fileStore`'s `items` slice would remain unchanged.

By using a pointer receiver `(fs *fileStore)`, the methods get the memory address of the original `fileStore` and can modify its `items` slice directly.

**Rule of Thumb:** If a method needs to modify a field in its receiver struct, it **must** have a pointer receiver.

## Your Task

Your task is to implement the `NewFileStore` function and all the methods for the `fileStore` type in `internal/clitodofiles/exercise.go`.

1.  **`NewFileStore(path string) Store`**: This function should create and return a new `*fileStore`.
2.  **`Load() error`**: This method should read the JSON file from the `path`, use `json.Unmarshal` to decode it into `fs.items`.
3.  **`Save() error`**: This method should use `json.MarshalIndent` to encode `fs.items` into JSON and write it to the file.
4.  **`Add(text string) Item`**: This method should create a new `Item`, append it to `fs.items`, and return the new item. Remember to give it a unique ID!
5.  **`Toggle(id int) (Item, bool)`**: This method should find the item with the given `id` and flip its `Done` status.
6.  **`List(onlyPending bool) []Item`**: This method should return a new slice of `Item`s, filtered based on the `onlyPending` flag.

The `main` function in `cmd/app/main.go` is already written for you. It handles the flag parsing and calls the `Store` methods you are implementing.

## How to Verify Your Work

First, build the application from the project root (`go-edu`):
```bash
go build -o todo ./minis/05-cli-todo-files/cmd/app
```

Now, use it from your terminal!
```bash
# Add a new item
./todo -add "Finish the CLI project"

# Add another item
./todo -add "Celebrate"

# List all items
./todo -list

# Mark item 1 as complete
./todo -toggle 1

# List only pending items
./todo -list -pending
```

You can also run the automated tests from the lesson directory (`minis/05-cli-todo-files`):
```bash
go test -v ./...
```
If the tests pass and your CLI application works as described, you're ready for the next lesson!

## Related Lessons
- Previous: `minis/04-jsonl-log-filter`
- Next: `minis/06-worker-pool-wordcount`
