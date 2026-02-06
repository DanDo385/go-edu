# 01: Hello Strings

## Core Concepts

- The concrete problem in Hello Strings and the correctness invariants it depends on.
- How values, pointers, slices, maps, or channels behave in this module's runtime path.
- Why this lesson's implementation pattern is the right next step in the learning arc.

## CS Connection

- Data representation and state transitions: what is copied, what is shared, and what can race.
- API boundaries: where we validate, where we propagate errors, and where we normalize output.
- Algorithmic tradeoffs in this lesson (latency, throughput, memory, and complexity).

## End-State Understanding

- Diagnose and implement Hello Strings patterns without relying on hidden framework behavior.
- Explain memory and concurrency implications of the final implementation choices.
- Compare learner code and reference code by invariants, not only by syntax.

## Why This Lesson Exists Here

Problem statement:
This lesson turns the previous module's concepts into a reusable engineering pattern for hello strings.

At this point in the arc:
Lesson 01 introduces a sharper systems concern so later modules can assume this mental model is stable.

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


Welcome to your first lesson! This challenge introduces you to Go's powerful and nuanced approach to strings. You'll learn why a "character" isn't always what it seems and how to handle text correctly in a global, multilingual world.

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

- The difference between **bytes** and **runes** in Go.
- How to correctly count characters in a string.
- How to reverse a string that contains multi-byte Unicode characters.

## The Challenge: Bytes vs. Runes

In Go, a `string` is an **immutable sequence of bytes**. This is efficient, but it has a critical implication: `len("gö")` is `3`, not `2`.

Why? Because Go uses UTF-8, a variable-width encoding:
- English characters like `'g'` take **1 byte**.
- Other characters like `'ö'` can take **2, 3, or even 4 bytes**.

This means you cannot assume that the *i-th byte* of a string is the *i-th character*. To handle characters correctly, Go gives us a special type: `rune`. A `rune` is an alias for `int32` and represents a single Unicode code point (i.e., a character).

| What you have | What it means | How you get it |
|---|---|---|
| `string` | A sequence of **bytes** | `s := "hello"` |
| `[]rune` | A sequence of **characters** | `runes := []rune(s)` |

**Rule of Thumb:** Use `string` by default. Convert to `[]rune` only when you need to work with individual characters.

## Your Task

Your task is to implement three functions in the `internal/hellostrings/exercise.go` file.

1.  **`Reverse(s string) string`**: This function must reverse a string. To do this correctly for all languages, you will need to convert the string to a slice of runes, reverse the slice, and then convert it back to a string.
2.  **`RuneLen(s string) int`**: This function should return the number of *runes* (characters) in the string, not the number of bytes. The `unicode/utf8` package will be helpful here.
3.  **`TitleCase(s string) string`**: This function should capitalize the first letter of each word in a string. The `strings` and `unicode` packages are your friends.

Open `internal/hellostrings/exercise.go` and fill in the `// TODO` sections.

## How to Verify Your Work

Once you've implemented the functions, you can test your work by running the following command from this directory (`minis/01-hello-strings`):

```bash
go test -v ./...
```

If all tests pass, you have successfully completed the lesson! You are ready to move on to the next one.

## Related Lessons
- Previous: Start here
- Next: `minis/02-arrays-maps-basics`
