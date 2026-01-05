# 04: JSON Lines Log Filter

## The Big Picture: The Power of Structured Logging

Imagine trying to find a specific error in a massive text file filled with messy, inconsistent log messages. It's a nightmare. This is the problem that **structured logging** solves.

Instead of writing plain text like `[ERROR] User login failed for 'admin'`, you write a machine-readable object, typically in JSON format:

```json
{"level": "error", "ts": "2024-05-10T14:32:01Z", "msg": "User login failed", "username": "admin"}
```

This is revolutionary because your logs are now a rich, queryable dataset. You can easily filter, aggregate, and visualize them using powerful tools. This project introduces you to **JSON Lines** (`.jsonl`), a common format for structured logs where each line is a self-contained JSON object. You will build a tool to parse, filter, and sort these logs—a core task in modern software operations.

## First Principles: Parsing and Data Mapping

1.  **JSON (JavaScript Object Notation)**: A lightweight, text-based data-interchange format. It represents data as a collection of key-value pairs (like Go maps) and ordered lists (like Go slices). Its simplicity and readability have made it the de facto standard for APIs and configuration files.

2.  **Schema on Read**: Unlike a rigid database schema, a log file often has an implicit structure. The process of parsing a JSON object and mapping it into a typed Go `struct` is a form of applying a "schema on read." You are enforcing the structure you expect, and deciding how to handle data that doesn't conform.

3.  **Enumerations (Enums)**: Log levels (`DEBUG`, `INFO`, `WARN`, `ERROR`) represent a fixed set of distinct values. The best way to represent this in code is with an **enumeration**, which maps human-readable names to efficient integer constants. This prevents typos and allows for fast, numeric comparisons (e.g., `level >= WARN`).

## Key Go Concepts in This Project

This project dives deep into Go's powerful `encoding/json` package.

### `encoding/json` and Struct Tags

Go's standard library makes JSON handling almost trivial. The `json.Unmarshal` function can automatically parse a JSON object into a Go `struct`. It does this using **reflection** and **struct tags**.

```go
type Entry struct {
    // This tag tells Unmarshal to look for a JSON key named "ts"
    // and put its value into the `TS` field of the struct.
    TS    time.Time `json:"ts"`
    Level Level     `json:"level"`
    Msg   string    `json:"msg"`
}
```

This declarative mapping is extremely powerful and readable. The `json` package handles all the details of converting JSON types (string, number, boolean) into Go types.

### Custom Unmarshaling with `json.Unmarshaler`

What happens when the default behavior isn't enough? In our logs, the `level` is a string (`"warn"`), but we want to store it as an efficient integer `Level` enum.

To do this, we can implement the `json.Unmarshaler` interface.

```go
// By implementing this method on a *pointer* to Level,
// we teach the json package how to parse a Level.
func (l *Level) UnmarshalJSON(data []byte) error {
    var s string
    json.Unmarshal(data, &s) // First, parse the JSON string (e.g., "warn")
    switch s {
    case "warn":
        *l = Warn // Assign the correct integer value
    // ... other cases
    }
    return nil
}
```
When `json.Unmarshal` sees a field of type `Level`, it checks if `*Level` has an `UnmarshalJSON` method. If so, it calls our custom logic instead of the default. This is a prime example of Go's philosophy of extending types with behavior via interfaces.

### `iota` for Creating Enums

Go doesn't have a dedicated `enum` keyword, but it provides a clean and idiomatic way to create them using `const` and the `iota` keyword. `iota` is a special constant that acts as an incrementing counter in a `const` block.

```go
type Level int

const (
    Debug Level = iota // 0
    Info               // 1
    Warn               // 2
    Error              // 3
)
```
This creates a new `Level` type (backed by `int`) and a set of named constants. This is far better than using raw strings, as it's type-safe, memory-efficient, and allows for trivial comparisons (`if level >= Warn`).

### `sort.Slice` with Closures

Once you have a slice of `Entry` structs, how do you sort them by timestamp? Go's `sort.Slice` is the perfect tool. It takes a slice and a "less" function. This function, often provided as a **closure** (an anonymous function), tells the sort algorithm how to compare any two elements.

```go
sort.Slice(entries, func(i, j int) bool {
    // The "less" function. Return true if element `i` should come before `j`.
    return entries[i].TS.Before(entries[j].TS)
})
```
This is incredibly expressive. You can sort by any field or combination of fields simply by changing the logic inside the closure.

## Progression: From Parsing to APIs

This project is a direct evolution of the previous one:
*   **Project 03 (CSV)**: You parsed line-delimited, structured text.
*   **Project 04 (JSONL)**: You are again parsing line-delimited, structured text, but the structure is now nested and more complex, requiring a more powerful parsing strategy (`json.Unmarshal` vs. `csv.Reader`).

The skills you learn here—especially `json.Unmarshal` and struct tags—are the **exact same skills** you will use to build web APIs, which almost universally use JSON as their data format. This project is your bridge from simple data processing to building web services.

## How to Run and Test

```bash
# The dev harness runs your function against a sample in-memory JSONL string.
go run ./cmd/dev

# The tests cover filtering, sorting, custom unmarshaling, and error handling.
go test -v ./...
```