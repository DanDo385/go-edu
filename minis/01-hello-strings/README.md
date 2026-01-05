# 01: Hello, Strings!

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