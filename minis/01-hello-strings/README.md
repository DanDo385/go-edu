# Project 01: hello-strings

## What You're Building

A set of string manipulation utilities that handle UTF-8 text properly. You'll implement title-case conversion, string reversal (character-aware, not byte-aware), and rune counting. These exercises demonstrate Go's first-class support for Unicode and the distinction between bytes and runes.

## Concepts Covered

- `strings` package for text manipulation
- `unicode/utf8` for proper Unicode handling
- Difference between bytes and runes (critical in Go!)
- Rune slices for character-level operations
- Functions as first-class citizens
- Table-driven testing patterns

## How to Run

```bash
# Run the program
make run P=01-hello-strings

# Or directly:
go run ./minis/01-hello-strings/cmd/hello-strings

# Run tests
go test ./minis/01-hello-strings/...

# Run tests with verbose output
go test -v ./minis/01-hello-strings/...
```

## Solution Explanation

### TitleCase
Convert the first letter of each word to uppercase using `strings.Fields()` to split on whitespace, then `unicode.ToUpper()` for the first rune of each word. This is more robust than simple byte manipulation because it handles multi-byte UTF-8 characters correctly.

### Reverse
Reversing strings in Go requires understanding that strings are immutable sequences of bytes, but characters (runes) can be 1-4 bytes each. The solution converts the string to a `[]rune` slice, reverses the slice in-place with a two-pointer swap, then converts back to a string. This ensures that emoji and other multi-byte characters remain intact.

### RuneLen
Go's `len()` function returns the **byte count**, not the character count. For UTF-8 strings containing emoji or non-ASCII characters, these differ. We use `utf8.RuneCountInString()` from the standard library to count actual characters (code points).

## Where Go Shines

**Go vs Python/JavaScript:**
- Python's `str[::-1]` and JavaScript's `split('').reverse().join('')` both work because they handle Unicode by default, but Go's explicit byte/rune distinction forces you to think about encoding—preventing subtle bugs
- Go's standard library includes UTF-8 utilities (`unicode/utf8`) without external dependencies
- Go strings are immutable and efficiently shared (cheap to pass by value)

**Go vs Rust:**
- Rust's `String` and `&str` distinction is steeper to learn
- Go's `[]rune` conversion is more intuitive than Rust's `chars().rev().collect()`
- Both are memory-safe, but Go's garbage collector means no lifetime annotations

## Stretch Goals

1. **Add a `Palindrome(s string) bool` function** that checks if a string reads the same forwards and backwards (use your `Reverse` function!)
2. **Implement `TruncateWords(s string, maxWords int) string`** that limits text to N words, adding "..." if truncated
3. **Create `CountVowels(s string) int`** using `unicode.IsLetter()` and a vowel set
4. **Benchmark** different reverse implementations (byte-based vs rune-based) for ASCII-only strings
