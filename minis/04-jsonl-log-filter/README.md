# 04: Jsonl Log Filter

## Core Concepts

- The concrete problem in Jsonl Log Filter and the correctness invariants it depends on.
- How values, pointers, slices, maps, or channels behave in this module's runtime path.
- Why this lesson's implementation pattern is the right next step in the learning arc.

## CS Connection

- Data representation and state transitions: what is copied, what is shared, and what can race.
- API boundaries: where we validate, where we propagate errors, and where we normalize output.
- Algorithmic tradeoffs in this lesson (latency, throughput, memory, and complexity).

## End-State Understanding

- Diagnose and implement Jsonl Log Filter patterns without relying on hidden framework behavior.
- Explain memory and concurrency implications of the final implementation choices.
- Compare learner code and reference code by invariants, not only by syntax.

## Why This Lesson Exists Here

Problem statement:
This lesson turns the previous module's concepts into a reusable engineering pattern for jsonl log filter.

At this point in the arc:
Lesson 04 introduces a sharper systems concern so later modules can assume this mental model is stable.

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


This project introduces you to structured logging, a cornerstone of modern application observability. You will build a tool to parse, filter, and sort JSON Lines (`.jsonl`) data, a common format for logs.

## CS Connection

- Memory layout drives behavior: variables store values, and some values are addresses into other storage.
- Go is pass-by-value, including pointers; copying a pointer value copies an address, not the pointee.
- Correctness depends on understanding copying versus aliasing (`*T`, slices, maps, and channels) and enforcing synchronization when concurrent access exists.

## End-State Understanding

- Implement the exercise with explicit reasoning about correctness, edge cases, and error paths.

## What You'll Learn

- How to parse JSON into Go structs using `json.Unmarshal` and **struct tags**.
- How to implement the `json.Unmarshaler` interface for **custom parsing logic**.
- How to create type-safe **enums** using `iota`.
- How to write flexible sorting logic using `sort.Slice` and **closures**.
- Why a **pointer receiver** is used on methods that modify their receiver.

## Core Concepts

### Structured Logging and JSON Lines

Instead of messy text, structured logging writes logs as machine-readable JSON objects.
```json
{"level": "error", "ts": "2024-05-10T14:32:01Z", "msg": "User login failed"}
```
A **JSON Lines** (`.jsonl`) file is simply a file where each line is a separate JSON object, making it ideal for streaming.

### `json.Unmarshal` and Struct Tags

Go can automatically parse JSON into a `struct`. You use `struct tags` to tell the `json` package which JSON key corresponds to which struct field.
```go
type Entry struct {
    TS    time.Time `json:"ts"` // The `json:"ts"` is the struct tag.
    Level Level     `json:"level"`
    Msg   string    `json:"msg"`
}
```

### Enums with `iota`

Go doesn't have a dedicated `enum` keyword, but `iota` provides a way to create them. `iota` is a special constant that increments in a `const` block. This is perfect for log levels.
```go
type Level int

const (
    Debug Level = iota // 0
    Info               // 1
    Warn               // 2
    Error              // 3
)
```

### Custom Parsing with `json.Unmarshaler` and Pointer Receivers

What if you want to parse the string `"warn"` into your integer `Warn` enum? You can teach Go's JSON parser how to do this by implementing the `json.Unmarshaler` interface.

```go
func (l *Level) UnmarshalJSON(data []byte) error {
    // ... your parsing logic ...
}
```
This is where we must have a **deep explanation of the pointer receiver `(l *Level)`**.

---

### 🚨 Deep Dive: Pointer Receivers (`*`)

Notice the `*` in `(l *Level)`. This is a **pointer receiver**. It means the `UnmarshalJSON` method gets a *pointer to* the `Level` variable it's being called on, not a copy of it.

Let's break it down:
1.  **The Goal:** When the `json.Unmarshal` function finds a `level` field in the JSON, it needs to change the value of our `Level` variable in the `Entry` struct.
2.  **Passing by Value (The Problem):** If the method was `func (l Level) ...` (no `*`), `l` would be a *copy* of the `Level` variable. If we changed `l` inside the method, we would only be changing the copy. The original `Level` variable in the `Entry` struct would be unaffected.
3.  **Passing by Reference (The Solution):** The pointer receiver `(l *Level)` solves this.
    *   `l` is not the `Level` value itself, but a variable that holds the **memory address** of the original `Level` variable.
    *   Inside the method, when we use the `*` operator again, as in `*l = Warn`, it means: "take the memory address stored in `l`, go to that address, and put the value `Warn` there."

**In summary:** We use a pointer receiver because the method needs to **modify the original value**. Without the pointer, the method would operate on a copy, and its changes would be lost.

---

## Your Task

Your task is to implement the `FilterAndSort(io.Reader, Level, SortField) ([]Entry, error)` function in `internal/jsonllogfilter/exercise.go`.

This function should:
1.  Read the `io.Reader` line by line.
2.  For each line, use `json.Unmarshal` to parse the JSON object into an `Entry` struct. Your custom `UnmarshalJSON` method for the `Level` type will be called automatically.
3.  Filter the entries, keeping only those at or above the `minLevel`.
4.  Sort the resulting slice based on the `sortBy` field. Use `sort.Slice` with a closure to handle the different sort fields.
5.  Return the final, filtered, and sorted slice of `Entry` structs.

Open `internal/jsonllogfilter/exercise.go` and fill in the `// TODO` sections.

## How to Verify Your Work

Run the following command from this directory (`minis/04-jsonl-log-filter`):

```bash
go test -v ./...
```
If the tests pass, you have successfully completed the lesson.

## Related Lessons
- Previous: `minis/03-csv-stats`
- Next: `minis/05-cli-todo-files`
