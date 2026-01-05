# 03: Processing CSV Data

## The Big Picture: From Raw Data to Insight

In the real world, data rarely comes in a perfect, ready-to-use format. It's often stored in structured text files like CSV (Comma-Separated Values), JSON, or XML. The ability to read, parse, and extract meaning from these files is a fundamental skill in software engineering.

This project takes you one step further than the last. You're no longer just counting words; you're now parsing a structured dataset—a CSV file of financial transactions—and aggregating the data to produce meaningful statistics. This is a mini version of a real-world ETL (Extract, Transform, Load) pipeline.

## First Principles: Structured Data and Streaming

1.  **Structured Data**: CSV is a simple format for representing a table of data. Each line is a **record** (or row), and commas separate the **fields** (or columns) within that record. The first line is often a **header**, which gives a name to each field. This structure turns a simple text file into a rich dataset.

2.  **Streaming vs. In-Memory Processing**: When dealing with files, especially potentially large ones, you have two basic approaches:
    *   **In-Memory**: Read the entire file into memory at once. This is simple but can fail if the file is larger than the available RAM.
    *   **Streaming**: Read and process the file piece by piece (e.g., line by line). This uses a small, constant amount of memory, regardless of the file size. It is far more scalable and robust.

    Go's `encoding/csv` package provides a **streaming parser**, which is the professional way to handle this task.

## Key Go Concepts in This Project

This project introduces several powerful Go features for data processing.

### `encoding/csv` Package

Go's standard library provides a dedicated CSV parser. Key features:
*   **Streaming API**: The `csv.Reader` reads one record at a time with its `Read()` method. This keeps memory usage low and constant.
*   **Robustness**: It automatically handles complexities like quoted fields containing commas (`"Doe, John"`) and escaped quotes.

### `struct` for Custom Data Types

While a `map[string]int` was enough to count words, here we need to store more information for each category (a count, a sum, and an average). A `struct` is the perfect tool for this. It lets you define a new, composite type with named fields.

```go
// Stat holds the aggregated statistics for one category.
// It's a custom data type we've defined.
type Stat struct {
    Count int
    Sum   float64
    Avg   float64
}
```

This makes your code more readable and maintainable. `stats["groceries"].Sum` is much clearer than `stats["groceries"][1]`.

### The Read-Modify-Write Pattern for Map Values

A critical concept arises when your map values are structs. Consider this update logic:

```go
// 1. Read: Get a *copy* of the Stat struct for the category.
// If the key doesn't exist, Go returns a zero-value Stat.
s := stats[category]

// 2. Modify: Update the fields on the local copy `s`.
s.Count++
s.Sum += amount

// 3. Write: Put the modified copy *back* into the map.
// This step is essential!
stats[category] = s
```

You **cannot** modify the struct field directly in the map (e.g., `stats[category].Count++`). This is because map elements are not addressable in Go. The compiler prevents this to ensure memory safety, as the map might need to move its data during a resize, invalidating any pointers to its elements.

### `strconv`: String Conversion

Parsing data almost always involves converting strings to other types. Go's `strconv` package is your tool for this. `strconv.ParseFloat()` and `strconv.Atoi()` are common functions that take a string and return a number and an `error`, forcing you to handle cases where the string is not a valid number.

### Error Wrapping

In this project, errors can happen at different stages (I/O error, parsing error, invalid data). Instead of just returning the last error you saw, it's better to add context. Go 1.13 introduced **error wrapping**:

```go
// If strconv.ParseFloat fails, we wrap its error with our own context.
return nil, fmt.Errorf("row %d: invalid amount %q: %w", rowNum, amountStr, err)
```
The `%w` verb preserves the original error, allowing higher-level functions to inspect the entire chain of errors.

## Progression: Building on Previous Concepts

*   **Project 01 (Strings)**: Your string manipulation skills are used here for parsing and converting data.
*   **Project 02 (Maps)**: This project's core aggregation logic is a more advanced application of the `map` data structure you just learned, using a `struct` as the value type.
*   **`io.Reader`**: We continue to use the `io.Reader` interface, reinforcing its importance for writing flexible, testable code that can read from a file, a network connection, or an in-memory string buffer.

Parsing structured data is a gateway skill. It's the first step in building web servers that accept JSON, services that interact with databases, and any tool that performs data analysis.

## How to Run and Test

```bash
# The dev harness runs your function against a sample in-memory CSV.
go run ./cmd/dev

# The tests cover multiple scenarios, including valid data,
# malformed rows, and edge cases.
go test -v ./...
```