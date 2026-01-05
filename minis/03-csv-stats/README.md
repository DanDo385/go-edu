# 03: Processing CSV Data

This project marks a significant step up in complexity. You'll move from simple collections to processing structured data from files—a core competency for any backend or data-focused engineer. You will build a mini ETL (Extract, Transform, Load) pipeline that reads real-world data, processes it, and produces valuable insights.

## Table of Contents

- [Learning Objectives](#learning-objectives)
- [The Big Picture: From Raw Data to Insight](#the-big-picture-from-raw-data-to-insight)
- [First Principles: Structured Data and Streaming](#first-principles-structured-data-and-streaming)
- [Project Structure](#project-structure)
- [Key Go Concepts in This Project](#key-go-concepts-in-this-project)
  - [Data Flow at a Glance](#data-flow-at-a-glance)
  - [`encoding/csv` Package](#encodingcsv-package)
  - [`struct` for Custom Data Types](#struct-for-custom-data-types)
  - [The Read-Modify-Write Pattern for Map Values](#the-read-modify-write-pattern-for-map-values)
  - [`strconv`: String Conversion](#strconv-string-conversion)
  - [Error Wrapping](#error-wrapping)
- [Progression: Building on Previous Concepts](#progression-building-on-previous-concepts)
- [How to Run and Test](#how-to-run-and-test)
- [Key Takeaways](#key-takeaways)
- [Further Reading](#further-reading)

## Learning Objectives

By the end of this project, you will be able to:

-   **Parse structured text files** using Go's `encoding/csv` package.
-   **Understand and apply streaming** to process files efficiently, regardless of their size.
-   **Define and use `struct` types** to model complex data.
-   **Master the "read-modify-write" pattern** for updating structs within a map.
-   **Perform robust string-to-number conversions** using `strconv`.
-   **Implement error wrapping** to provide meaningful context to errors in a data pipeline.

## The Big Picture: From Raw Data to Insight

In the real world, data rarely comes in a perfect, ready-to-use format. It's often stored in structured text files like CSV (Comma-Separated Values), JSON, or XML. The ability to read, parse, and extract meaning from these files is a fundamental skill in software engineering.

This project takes you one step further than the last. You're no longer just counting words; you're now parsing a structured dataset—a CSV file of financial transactions—and aggregating the data to produce meaningful statistics. This is a mini version of a real-world ETL (Extract, Transform, Load) pipeline.

## First Principles: Structured Data and Streaming

1.  **Structured Data**: CSV is a simple format for representing a table of data. Each line is a **record** (or row), and commas separate the **fields** (or columns) within that record. The first line is often a **header**, which gives a name to each field. This structure turns a simple text file into a rich dataset.

2.  **Streaming vs. In-Memory Processing**: When dealing with files, especially potentially large ones, you have two basic approaches:
    *   **In-Memory**: Read the entire file into memory at once. This is simple but can fail if the file is larger than the available RAM.
    *   **Streaming**: Read and process the file piece by piece (e.g., line by line). This uses a small, constant amount of memory, regardless of the file size. It is far more scalable and robust.

    Go's `encoding/csv` package provides a **streaming parser**, which is the professional way to handle this task.

## Project Structure

```
.
├── cmd/
│   └── dev/
│       └── main.go       # A simple program to test your functions manually.
├── internal/
│   └── stats.go      # Where you will implement the CSV processing logic.
└── testdata/
    └── transactions.csv # Sample data used for testing.
```
-   **`cmd/dev`**: An entry point for a development harness. Run `go run ./cmd/dev` to see your function in action.
-   **`internal/`**: Contains the core logic.
-   **`testdata/`**: Holds data files needed for your tests.

## Key Go Concepts in This Project

### Data Flow at a Glance

Here's how data moves through our program:

```
+----------------+      +----------------+      +-------------------+      +-------------------------+
|   CSV File     |----->|  csv.Reader    |----->|  For Loop         |----->|  map[string]Stat        |
| (io.Reader)    |      | (Streaming)    |      | (Read-Modify-Write) |      |  (Final Aggregation)    |
+----------------+      +----------------+      +-------------------+      +-------------------------+
       |                       |                      | (data)                 | (stats)
       | (text)                | (record)             |                        |
       v                       v                      v                        v
 "Groceries,25.50\n" --> `[]string{"Groceries", "25.50"}` --> s.Sum += 25.50 --> `stats["Groceries"] = Stat{...}`
 "Rent,800.00\n"   --> `[]string{"Rent", "800.00"}` --> s.Sum += 800.00 --> `stats["Rent"] = Stat{...}`
```

1.  **Input**: An `io.Reader` provides the raw, byte-stream of the CSV data.
2.  **Parsing**: `csv.Reader` reads from the `io.Reader` and decodes the stream one record (`[]string`) at a time.
3.  **Processing**: Our `for` loop iterates over the records. Inside the loop, we:
    *   Parse the amount string into a `float64`.
    *   Perform the crucial **read-modify-write** pattern on our `stats` map.
4.  **Output**: The final map, containing the aggregated `Stat` for each category, is returned.

### `encoding/csv` Package

Go's standard library provides a dedicated CSV parser. Key features:
*   **Streaming API**: The `csv.Reader` reads one record at a time with its `Read()` method. This keeps memory usage low and constant.
*   **Robustness**: It automatically handles complexities like quoted fields containing commas (`"Doe, John"`) and escaped quotes.

### `struct` for Custom Data Types

While a `map[string]int` was enough to count words, here we need to store more information for each category (a count, a sum, and eventually an average). A `struct` is the perfect tool for this. It lets you define a new, composite type with named fields.

```go
// Stat holds the aggregated statistics for one category.
// It's a custom data type we've defined.
type Stat struct {
    Count int
    Sum   float64
    Avg   float64 // Note: This is calculated at the end.
}
```

This makes your code more readable and maintainable. `stats["groceries"].Sum` is much clearer than a hypothetical `stats["groceries"][1]`.

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

You **cannot** modify the struct field directly in the map (e.g., `stats[category].Count++`). This is because map elements are not addressable in Go. The compiler prevents this to ensure memory safety, as the map might need to move its data during a resize, invalidating any pointers to its elements. By enforcing the read-modify-write pattern, Go guarantees safety.

### `strconv`: String Conversion

Parsing data almost always involves converting strings to other types. Go's `strconv` package is your tool for this. `strconv.ParseFloat()` and `strconv.Atoi()` are common functions that take a string and return a number and an `error`, forcing you to handle cases where the string is not a valid number.

### Error Wrapping

In this project, errors can happen at different stages (I/O error, CSV format error, data conversion error). Instead of just returning the last error you saw, it's better to add context. Go 1.13 introduced **error wrapping**:

```go
// If strconv.ParseFloat fails, we wrap its error with our own context.
if err != nil {
    // The original error `err` is "wrapped" in a new, more descriptive error.
    return nil, fmt.Errorf("row %d: invalid amount %q: %w", rowNum, amountStr, err)
}
```
The `%w` verb preserves the original error. This is incredibly useful for debugging. A caller can see the full story: "could not calculate stats: failed on row 5: invalid amount \"twenty-five\": `strconv.ParseFloat`: parsing \"twenty-five\": invalid syntax".

## Progression: Building on Previous Concepts

*   **Project 01 (Strings)**: Your string manipulation skills are used here for parsing and converting data.
*   **Project 02 (Maps & Structs)**: This project's core aggregation logic is a more advanced application of the `map` data structure, now using a `struct` as the value type.
*   **`io.Reader`**: We continue to use the `io.Reader` interface, reinforcing its importance for writing flexible, testable code that can read from a file, a network connection, or an in-memory string buffer.

Parsing structured data is a gateway skill. It's the first step in building web servers that accept JSON, services that interact with databases, and any tool that performs data analysis.

## How to Run and Test

```bash
# The dev harness runs your function against a sample in-memory CSV.
go run ./cmd/dev

# The tests cover multiple scenarios, including valid data,
# malformed rows, and edge cases.
gotest -v ./...
```

## Key Takeaways

-   **Process large files by streaming** (`csv.Reader`) to keep memory usage low.
-   **Use `structs` to group related data** into custom types.
-   Updating a struct in a map **requires the read-modify-write pattern**.
-   **`strconv` is the standard tool** for converting strings to numbers and other types.
-   **Wrap errors with `fmt.Errorf` and `%w`** to create a chain of context for easier debugging.

## Further Reading

*   [**Package `encoding/csv`**](https://pkg.go.dev/encoding/csv): Official documentation for the CSV reader and writer.
*   [**Package `strconv`**](https://pkg.go.dev/strconv): Documentation for string conversion functions.
*   [**A Tour of Go: Structs**](https://go.dev/tour/moretypes/2): An interactive introduction to structs.
*   [**Go by Example: Error Wrapping**](https://gobyexample.com/error-wrapping): A concise example of how and why to wrap errors.