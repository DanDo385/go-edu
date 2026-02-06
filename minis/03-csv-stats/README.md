# 03: Csv Stats

## Core Concepts

- The concrete problem in Csv Stats and the correctness invariants it depends on.
- How values, pointers, slices, maps, or channels behave in this module's runtime path.
- Why this lesson's implementation pattern is the right next step in the learning arc.

## CS Connection

- Data representation and state transitions: what is copied, what is shared, and what can race.
- API boundaries: where we validate, where we propagate errors, and where we normalize output.
- Algorithmic tradeoffs in this lesson (latency, throughput, memory, and complexity).

## End-State Understanding

- Diagnose and implement Csv Stats patterns without relying on hidden framework behavior.
- Explain memory and concurrency implications of the final implementation choices.
- Compare learner code and reference code by invariants, not only by syntax.

## Why This Lesson Exists Here

Problem statement:
This lesson turns the previous module's concepts into a reusable engineering pattern for csv stats.

At this point in the arc:
Lesson 03 introduces a sharper systems concern so later modules can assume this mental model is stable.

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


This lesson is a step up in complexity. You'll build a mini data-processing pipeline that reads a CSV file, parses it, and calculates statistics. This is a core skill for any backend or data-focused engineer.

## CS Connection

- Memory layout drives behavior: variables store values, and some values are addresses into other storage.
- Go is pass-by-value, including pointers; copying a pointer value copies an address, not the pointee.
- Correctness depends on understanding copying versus aliasing (`*T`, slices, maps, and channels) and enforcing synchronization when concurrent access exists.

## End-State Understanding

- Implement the exercise with explicit reasoning about correctness, edge cases, and error paths.

## What You'll Learn

- How to **parse structured text files** (CSV).
- How to use **structs** to define your own custom data types.
- The **"read-modify-write" pattern** for updating structs inside a map.
- How to convert strings to numbers using the `strconv` package.
- How to **wrap errors** to provide more context.

## Core Concepts

### Structs: Custom Data Types

While a `map[string]int` was enough for our word counter, what if we need to store more than just a single number? A `struct` lets you define a new type with named fields.

```go
// Stat holds the aggregated statistics for one category.
type Stat struct {
    Count int
    Sum   float64
    Avg   float64
}
```
Now we can have a `map[string]Stat` to store our aggregated data for each category.

### The Read-Modify-Write Pattern

A critical concept arises when your map values are structs. You **cannot** modify a field of a struct that is inside a map directly.

```go
// THIS WILL NOT COMPILE
stats[category].Count++
```

Instead, you must follow the "read-modify-write" pattern:
```go
// 1. Read: Get a *copy* of the Stat struct for the category.
s := stats[category]

// 2. Modify: Update the fields on the local copy `s`.
s.Count++
s.Sum += amount

// 3. Write: Put the modified copy *back* into the map.
stats[category] = s
```
This is a safety feature in Go related to how maps manage their memory.

### Streaming with `io.Reader`

Just like in the last lesson, our function will take an `io.Reader`. This allows us to "stream" the data, meaning we process it piece by piece instead of reading the entire file into memory at once. This is much more scalable. The `encoding/csv` package provides a streaming reader that is perfect for this task.

## Your Task

Your task is to implement the `SummarizeCSV(io.Reader) (map[string]Stat, error)` function in `internal/csvstats/exercise.go`.

The function should:
1.  Create a `csv.Reader` from the `io.Reader`.
2.  Read the CSV records one by one in a loop.
3.  Skip the header row.
4.  For each data row, parse the category and the transaction amount.
5.  Use `strconv.ParseFloat` to convert the amount string to a `float64`.
6.  Use the read-modify-write pattern to update a `map[string]Stat` with the statistics for each category.
7.  After the loop, calculate the average for each category.
8.  Return the final map and any errors that occurred. Remember to wrap your errors for better context!

Open `internal/csvstats/exercise.go` and fill in the `// TODO` sections.

## How to Verify Your Work

Run the following command from this directory (`minis/03-csv-stats`):

```bash
go test -v ./...
```
If the tests pass, you have successfully completed the lesson.

## Related Lessons
- Previous: `minis/02-arrays-maps-basics`
- Next: `minis/04-jsonl-log-filter`
