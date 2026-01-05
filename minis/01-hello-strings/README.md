# 01: Hello, Strings!

This project is your introduction to Go's powerful and nuanced approach to strings. It's designed to build a solid foundation for handling text, which is a critical skill in virtually all areas of software development.

## Table of Contents

- [Learning Objectives](#learning-objectives)
- [The Big Picture: Why Start with Strings?](#the-big-picture-why-start-with-strings)
- [Computer Science First Principles: What Is a String?](#computer-science-first-principles-what-is-a-string)
- [Go's String Implementation: Immutability and the String Header](#gos-string-implementation-immutability-and-the-string-header)
- [Project Structure](#project-structure)
- [Key Concepts in This Project](#key-concepts-in-this-project)
- [Common Pitfalls: The Gotchas of Go Strings](#common-pitfalls-the-gotchas-of-go-strings)
- [Visualizing the Difference: "Hi 👋"](#visualizing-the-difference-hi-)
- [Progression: Where We Go from Here](#progression-where-we-go-from-here)
- [How to Run and Test](#how-to-run-and-test)
- [Key Takeaways](#key-takeaways)
- [Further Reading](#further-reading)

## Learning Objectives

By the end of this project, you will be able to:

-   **Explain the difference between bytes and runes** in Go.
-   **Understand and articulate why UTF-8 is a variable-width encoding**.
-   **Use `string` and `[]rune` appropriately** for different tasks.
-   **Avoid common pitfalls** like incorrect character counting and iteration.
-   **Implement basic string manipulation functions** that correctly handle Unicode characters.
-   **Write and run tests** for your Go code.

## The Big Picture: Why Start with Strings?

Welcome to your first Go project! We're starting with strings because they are the bedrock of modern computing. Nearly every application, from a simple command-line tool to a complex web service, communicates with the world through text. Understanding how to manipulate strings efficiently and correctly is not just a basic skill—it is a critical one.

This project introduces you to Go's powerful and nuanced approach to strings, which is designed for a global, multilingual internet.

## Computer Science First Principles: What *Is* a String?

At its core, a **string is a sequence of characters**. But this simple definition hides a lot of complexity:

1.  **Representation in Memory**: Computers only understand numbers. To represent characters, we use a mapping called a character encoding.
    *   **ASCII**: The classic American standard, it uses 7 bits for 128 characters (English letters, numbers, basic punctuation). Simple and small.
    *   **Unicode**: The modern global standard. It defines a unique number (a "code point" or "rune") for every character in nearly every language, plus emojis and symbols. Unicode can represent over 149,000 characters.

2.  **Encoding vs. Code Point**:
    *   **Code Point (Rune)**: The abstract idea of a character, represented by a number (e.g., `U+0041` for 'A', `U+20AC` for '€', `U+1F600` for '😀').
    *   **Encoding**: The way we store those numbers in bytes. **UTF-8** is the dominant encoding on the web and in Go. It's a *variable-width* encoding:
        *   ASCII characters ('A', 'z', '7') take just **1 byte**.
        *   Common European characters ('é', 'ü', 'ñ') take **2 bytes**.
        *   Most Asian characters (e.g., '世', '界') take **3 bytes**.
        *   Emojis and rare symbols ('👋', '😀') take **4 bytes**.

This variable-width nature is why `len("go")` is 2, but `len("gö")` is 3. The string "gö" contains two characters but three bytes. This is the most common "gotcha" for developers new to Go, and this project tackles it head-on.

## Go's String Implementation: Immutability and the String Header

In Go, a `string` is an **immutable sequence of bytes**. This has profound implications:

*   **Immutability**: Once a string is created, its contents cannot be changed. When you "modify" a string (e.g., concatenate, slice, or change its case), Go creates a *new* string in memory.
    *   **Pro**: This makes strings safe to share across different parts of your program (including concurrently) without fear of modification. It eliminates a whole class of bugs.
    *   **Con**: Frequent modifications can lead to many small memory allocations, which can impact performance. Go's runtime is highly optimized for this, but it's a trade-off to be aware of.

*   **The String Header**: A `string` variable is a small, 2-word struct. It contains:
    1.  A **pointer** to the underlying byte array (which is stored elsewhere in memory).
    2.  The **length** of the string in bytes.

```
// A string is just a small header
stringHeader = {
    Data *byte // Pointer to the bytes
    Len  int   // Number of bytes
}
```

When you pass a string to a function, you are only copying this small header, not the entire byte array. This makes string passing very efficient.

## Project Structure

This project follows a standard Go project layout:

```
.
├── cmd/
│   └── dev/
│       └── main.go   # A simple program to test your functions manually.
└── internal/
    └── strings.go    # Where you will implement your string functions.
```

-   **`cmd/dev`**: This is the entry point for a small development harness. You can run `go run ./cmd/dev` to see your functions in action with sample inputs.
-   **`internal/`**: This directory contains the core logic of your project. The `internal` name is a Go convention: code in this directory can only be imported by code within the same project (`01-hello-strings`), not by other projects. This is where you'll write your code.

## Key Concepts in This Project

This project requires you to implement three common string utilities, forcing you to confront the byte vs. character distinction.

| Function | Concept Demonstrated | Why it's Important |
| :--- | :--- | :--- |
| `TitleCase(s)` | Word boundaries and case conversion. | Shows how to decompose a string into logical units (words) and apply Unicode-aware transformations. |
| `Reverse(s)` | Character-by-character manipulation. | This is impossible to do correctly at the byte level for UTF-8. It forces you to convert the string to a `[]rune` slice, which is Go's way of representing a sequence of *characters* (code points). |
| `RuneLen(s)` | Correctly counting characters. | Highlights the difference between `len(s)` (byte count) and `utf8.RuneCountInString(s)` (character/rune count). |

### `string` vs. `[]rune`: A Quick Comparison

| Characteristic | `string` | `[]rune` |
| :--- | :--- | :--- |
| **Nature** | Immutable sequence of bytes | Mutable sequence of code points (int32) |
| **Representation** | UTF-8 encoded bytes | One `int32` per character |
| **`len()` behavior** | Returns byte count | Returns character (rune) count |
| **Modification** | Not possible; must create a new string | Possible; you can change elements in place |
| **Memory** | Compact for ASCII-heavy text | Uses 4 bytes for *every* character, regardless of complexity |

**When to use which?**
*   Use `string` for storing and passing text. This is the default and most efficient choice.
*   Convert to `[]rune` only when you need to inspect or modify individual characters within a string.

### Common Pitfalls: The Gotchas of Go Strings

When working with strings in Go, developers often encounter a few common traps, especially when coming from other languages. Awareness is the first step to avoiding them.

1.  **Incorrect Character Counting with `len()`**:
    *   **The Trap**: Using `len(s)` to find the number of characters in a string.
    *   **The Reality**: `len(s)` returns the number of **bytes**, not characters (runes). For pure ASCII, this is the same, but for any other character, it's not.
    *   **The Fix**: Use `utf8.RuneCountInString(s)` for an efficient character count.
    
    ```go
    s := "gö"
    fmt.Println(len(s)) // "3" (bytes)
    fmt.Println(utf8.RuneCountInString(s)) // "2" (characters/runes)
    ```

2.  **Incorrect Iteration with `for i := 0; i < len(s); i++`**:
    *   **The Trap**: Accessing `s[i]` in a loop to get the i-th character.
    *   **The Reality**: This accesses the i-th **byte**. For a multi-byte character, this will lead to corrupt or partial character data.
    *   **The Fix**: Use a `for range` loop (`for i, r := range s`). Go automatically decodes one rune per iteration, giving you the correct character and its starting byte index.

    ```go
    s := "gö"
    for i := 0; i < len(s); i++ {
        fmt.Printf("Byte at index %d: %c\n", i, s[i]) 
        // Output will be broken for 'ö'
    }
    // Correct way:
    for i, r := range s {
        fmt.Printf("Rune at byte index %d: %c\n", i, r)
    }
    ```

3.  **Slicing in the Middle of a Rune**:
    *   **The Trap**: Slicing a string with byte indices, like `s[0:3]`, without knowing the character boundaries.
    *   **The Reality**: If a slice boundary falls in the middle of a multi-byte rune, you will get invalid UTF-8, often represented as the replacement character ().
    *   **The Fix**: If you must slice at the byte level, be careful. If you need to slice based on character count, convert to a `[]rune` slice first: `string([]rune(s)[0:2])`.

    ```go
    s := "Hello, 世界" // "世界" are 3 bytes each
    fmt.Println(s[0:7]) // "Hello, " - Safe
    fmt.Println(s[0:8]) // "Hello, " - Corrupted!
    
    // Correct way to get first 8 characters:
    runes := []rune(s)
    fmt.Println(string(runes[0:8])) // "Hello, 世"
    ```

### Visualizing the Difference: "Hi 👋"

Let's visualize the string `"Hi 👋"` to make the distinction between bytes and runes crystal clear.

*   **Characters (Runes)**: 4
*   **Bytes (in UTF-8)**: 8

| Character | H | i | (space) | 👋 |
| :--- | :---: | :---: | :---: | :---: |
| **Unicode Code Point (Rune)** | `U+0048` | `U+0069` | `U+0020` | `U+1F44B` |
| **Rune (int32 value)** | `72` | `105` | `32` | `128075` |
| **UTF-8 Bytes (Hex)** | `48` | `69` | `20` | `F0 9F 91 8B` |
| **Byte Length** | 1 | 1 | 1 | 4 |

When you have the string `s := "Hi 👋"`:
*   `len(s)` returns **8** (the total number of bytes).
*   `utf8.RuneCountInString(s)` returns **4** (the number of runes).
*   `[]byte(s)` is `[72, 105, 32, 240, 159, 145, 139]`
*   `[]rune(s)` is `[72, 105, 32, 128075]`

## Progression: Where We Go from Here

Mastering strings is the first step. The concepts you learn here are foundational for many of the projects that follow:

*   **03-csv-stats** & **04-jsonl-log-filter**: You'll use string splitting and parsing to process structured text data.
*   **05-cli-todo-files**: You'll handle user input from the command line, which is always text.
*   **08-http-client-retries** & **09-http-server-graceful**: You'll build and parse HTTP requests and responses, which are heavily text-based.
*   **32-websocket-chatroom**: Real-time communication is all about sending and receiving text messages.

By understanding how Go handles text from first principles, you are building a solid foundation for creating robust, correct, and global-ready software.

## How to Run and Test

```bash
# Run the development harness to see your functions in action
go run ./cmd/dev

# Run the provided tests to check your implementation for correctness
go test -v ./...
```

## Key Takeaways

-   **Strings are immutable byte slices.**
-   **Go source code is UTF-8 by default.**
-   **`len()` gives byte count, not character count.**
-   **Use `for range` to iterate over runes in a string.**
-   **Convert to `[]rune` for character-level manipulation.**

## Further Reading

To dive deeper into the topics covered in this project, check out these official Go resources:
*   [**Blog Post: Strings, bytes, runes and characters in Go**](https://go.dev/blog/strings): The canonical and most important read on this topic.
*   [**Effective Go: Strings**](https://go.dev/doc/effective_go#strings): A great summary of idiomatic string usage in Go.
*   [**Package `unicode/utf8`**](https://pkg.go.dev/unicode/utf8): The official documentation for the `utf8` package, with functions for working with UTF-8 encoded text.
*   [**Package `strings`**](https://pkg.go.dev/strings): The standard library's main package for string manipulation.