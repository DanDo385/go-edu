# 04: JSON Lines Log Filter

This project introduces you to the world of structured logging, a cornerstone of modern application observability. You will build a practical command-line tool to parse, filter, and sort JSON Lines (`.jsonl`) data, a common format for logs. This is a crucial skill for debugging and monitoring production systems.

## Table of Contents

- [Learning Objectives](#learning-objectives)
- [The Big Picture: The Power of Structured Logging](#the-big-picture-the-power-of-structured-logging)
- [First Principles: Parsing and Data Mapping](#first-principles-parsing-and-data-mapping)
- [Project Structure](#project-structure)
- [Key Go Concepts in This Project](#key-go-concepts-in-this-project)
  - [`encoding/json` and Struct Tags](#encodingjson-and-struct-tags)
  - [Custom Unmarshaling with `json.Unmarshaler`](#custom-unmarshaling-with-jsonunmarshaler)
  - [`iota` for Creating Enums](#iota-for-creating-enums)
  - [`sort.Slice` with Closures](#sortslice-with-closures)
- [Progression: From Parsing to APIs](#progression-from-parsing-to-apis)
- [How to Run and Test](#how-to-run-and-test)
- [Key Takeaways](#key-takeaways)
- [Further Reading](#further-reading)

## Learning Objectives

By the end of this project, you will be able to:

-   **Parse JSON into Go structs** using `json.Unmarshal` and struct tags.
-   **Implement custom JSON parsing logic** by satisfying the `json.Unmarshaler` interface.
-   **Create type-safe enumerations (enums)** using `iota` for concepts like log levels.
-   **Write flexible sorting logic** using `sort.Slice` and closures.
-   **Understand the benefits of structured logging** and the JSON Lines format.
-   **Build a complete, practical tool** for processing log files.

## The Big Picture: The Power of Structured Logging

Imagine trying to find a specific error in a massive text file filled with messy, inconsistent log messages. It's a nightmare. This is the problem that **structured logging** solves.

Instead of writing plain text like `[ERROR] User login failed for 'admin'`, you write a machine-readable object, typically in JSON format:

```json
{"level": "error", "ts": "2024-05-10T14:32:01Z", "msg": "User login failed", "username": "admin"}
```

This is revolutionary because your logs are now a rich, queryable dataset. You can easily filter, aggregate, and visualize them using powerful tools. This project introduces you to **JSON Lines** (`.jsonl`), a common format for structured logs where each line is a self-contained JSON object. You will build a tool to parse, filter, and sort these logs—a core task in modern software operations.

### JSON vs. JSON Lines (`.jsonl`)

-   **A standard JSON file** contains a *single* JSON object or array. The entire file must be read and parsed at once to be valid.
-   **A JSON Lines file** contains *multiple*, independent JSON objects, with exactly one object per line. Each line is a valid JSON object on its own, making it ideal for streaming.

## First Principles: Parsing and Data Mapping

1.  **JSON (JavaScript Object Notation)**: A lightweight, text-based data-interchange format. It represents data as key-value pairs (like Go maps) and ordered lists (like Go slices).
2.  **Schema on Read**: The process of parsing a JSON object and mapping it into a typed Go `struct` is a form of applying a "schema on read." You are enforcing the structure you expect and deciding how to handle data that doesn't conform.
3.  **Enumerations (Enums)**: Log levels (`DEBUG`, `INFO`, `WARN`, `ERROR`) represent a fixed set of distinct values. The best way to represent this in code is with an **enumeration**, which maps human-readable names to efficient integer constants.

## Project Structure

```
.
├── cmd/
│   └── dev/
│       └── main.go       # A simple program to test your functions manually.
├── internal/
│   └── filter.go     # Where you will implement the JSONL processing logic.
└── testdata/
    └── logs.jsonl    # Sample data used for testing.
```

-   **`cmd/dev`**: An entry point for a development harness. Run `go run ./cmd/dev` to see your function in action.
-   **`internal/`**: Contains the core logic.
-   **`testdata/`**: Holds data files needed for your tests.

## Key Go Concepts in This Project

### `encoding/json` and Struct Tags

Go's standard library makes JSON handling almost trivial. The `json.Unmarshal` function can automatically parse a JSON object into a Go `struct` using **reflection** and **struct tags**. A struct tag is a string literal that provides metadata to a package.

```go
type Entry struct {
    // This tag tells Unmarshal to look for a JSON key named "ts"
    // and put its value into the `TS` field of the struct.
    TS    time.Time `json:"ts"`
    Level Level     `json:"level"`
    Msg   string    `json:"msg"`
}
```

### Custom Unmarshaling with `json.Unmarshaler`

What happens when the default behavior isn't enough? In our logs, the `level` is a string (`"warn"`), but we want to store it as an efficient integer `Level` enum. To do this, we can implement the `json.Unmarshaler` interface.

```go
// By implementing this method on a *pointer* to Level,
// we teach the json package how to parse a Level from a string.
func (l *Level) UnmarshalJSON(data []byte) error {
    var s string
    // Step 1: Unmarshal the raw JSON (e.g., `"warn"`) into a plain string `s`.
    if err := json.Unmarshal(data, &s); err != nil {
        return fmt.Errorf("level must be a string: %w", err)
    }

    // Step 2: Convert the string to our internal integer representation.
    switch strings.ToLower(s) {
    case "debug":
        *l = Debug
    case "info":
        *l = Info
    case "warn":
        *l = Warn
    case "error":
        *l = Error
    default:
        return fmt.Errorf("unknown level %q", s)
    }
    return nil
}
```
When `json.Unmarshal` sees a field of type `Level`, it checks if `*Level` has an `UnmarshalJSON` method. If so, it calls our custom logic. This is a prime example of Go's philosophy of extending types with behavior via interfaces.

### `iota` for Creating Enums

Go doesn't have a dedicated `enum` keyword, but `iota` provides an idiomatic way to create them. `iota` is a special constant that acts as an incrementing counter in a `const` block.

```go
type Level int

// iota starts at 0 and increments for each line in the const block.
const (
    Debug Level = iota // 0
    Info               // 1
    Warn               // 2
    Error              // 3
)
```
This is far better than using raw strings (`"warn"`), as it's type-safe, memory-efficient, and allows for trivial comparisons (`if level >= Warn`).

### `sort.Slice` with Closures

Once you have a slice of `Entry` structs, how do you sort them? Go's `sort.Slice` is the perfect tool. It takes a slice and a "less" function. This function, often provided as a **closure** (an anonymous function), tells the sort algorithm how to compare any two elements.

```go
// Primary sort by timestamp (ascending)
sort.Slice(entries, func(i, j int) bool {
    return entries[i].TS.Before(entries[j].TS)
})

// Want to sort by level (descending) instead? Just change the closure!
sort.Slice(entries, func(i, j int) bool {
    // Higher level value means more severe, so it should come first.
    return entries[i].Level > entries[j].Level
})
```
This is incredibly expressive and flexible. You can sort by any field or combination of fields simply by changing the logic inside the closure.

## Progression: From Parsing to APIs

The skills you learn here—especially `json.Unmarshal` and struct tags—are the **exact same skills** you will use to build web APIs, which almost universally use JSON as their data format. This project is your bridge from simple data processing to building web services.

-   **Project 03 (CSV)**: You parsed line-delimited, structured text.
-   **Project 04 (JSONL)**: You are again parsing line-delimited, structured text, but the structure is now nested and more complex, requiring a more powerful parsing strategy (`json.Unmarshal` vs. `csv.Reader`).

## How to Run and Test

```bash
# The dev harness runs your function against a sample in-memory JSONL string.
go run ./cmd/dev

# The tests cover filtering, sorting, custom unmarshaling, and error handling.
go test -v ./...
```

## Key Takeaways

-   **Structured logging is essential** for modern software observability.
-   **`json.Unmarshal` with struct tags** is the idiomatic way to parse JSON in Go.
-   Implement the **`json.Unmarshaler` interface** for full control over parsing.
-   Use **`iota` to create efficient, type-safe enums**.
-   **`sort.Slice` with a closure** provides a flexible and powerful way to sort data.
-   The skills for parsing JSON logs are directly transferable to building JSON APIs.

## Further Reading

*   [**JSON Lines Website**](https://jsonlines.org/): The official spec and rationale for the format.
*   [**Blog Post: JSON and Go**](https://go.dev/blog/json): The Go team's guide to using the `encoding/json` package.
*   [**Package `encoding/json`**](https://pkg.go.dev/encoding/json): Official documentation, including details on `Unmarshal` and struct tags.
*   [**Go by Example: Sorting by Functions**](https://gobyexample.com/sorting-by-functions): A great example of using closures to implement custom sorts.