# Project 02: arrays-maps-basics

## What You're Building

A word frequency counter that reads text from a file (one word per line), builds a frequency map, and identifies the most common word. This project introduces Go's core data structures (slices and maps) and demonstrates idiomatic file I/O patterns.

## Concepts Covered

- Slices vs arrays (why slices are almost always preferred)
- Maps for key-value storage
- `io.Reader` interface for testable I/O
- `bufio.Scanner` for line-by-line reading
- Error handling with multiple return values
- Sorting maps by value (requires converting to slices)
- Case normalization with `strings.ToLower()`

## How to Run

```bash
# Run the program
make run P=02-arrays-maps-basics

# Or directly:
go run ./minis/02-arrays-maps-basics/cmd/arrays-maps-basics

# Run tests
go test ./minis/02-arrays-maps-basics/...

# Run tests with verbose output
go test -v ./minis/02-arrays-maps-basics/...
```

## Solution Explanation

### FreqFromReader Algorithm

1. **Read words line-by-line**: Use `bufio.Scanner` to read from any `io.Reader` (file, string, network). This is more efficient than reading the entire file into memory at once.

2. **Normalize to lowercase**: `strings.ToLower()` ensures "Hello" and "hello" count as the same word. Use `strings.TrimSpace()` to ignore blank lines.

3. **Build frequency map**: Go maps are hash tables with O(1) average lookup/insert. The pattern `map[key]++` automatically handles missing keys (zero value is 0).

4. **Find max**: Iterate the map once to find the word with the highest count. For ties, we return the first encountered (arbitrary but deterministic).

### Why `io.Reader` instead of `*os.File`?

By accepting `io.Reader`, our function works with:
- Files (`os.Open()`)
- Strings (`strings.NewReader()`)
- Network connections (`net.Conn`)
- Compressed data (`gzip.Reader`)

This is **dependency injection** via interfaces—a key Go pattern for testability.

## Where Go Shines

**Go vs Python:**
- Python: `Counter(words).most_common(1)` is one line, but requires importing `collections`
- Go: More verbose, but zero dependencies and faster execution
- Go's explicit error handling forces you to handle the "file not found" case

**Go vs JavaScript:**
- JS: No built-in file I/O in the browser; Node.js requires callbacks or promises
- Go: Synchronous I/O by default (simpler mental model), goroutines for concurrency when needed

**Go vs Rust:**
- Rust: `HashMap` and `BufReader` are similar, but requires handling UTF-8 errors explicitly
- Go: UTF-8 is assumed; simpler for text processing
- Both are fast, but Go's `map` syntax is less verbose

## Stretch Goals

1. **Add support for top N words**: Modify to return the top 5 most common words (requires sorting)
2. **Ignore common words**: Filter out "the", "a", "an" (stopwords) before counting
3. **Add word length statistics**: Track min/max/average word length
4. **Handle punctuation**: "hello," and "hello" should count as the same word (use `strings.TrimFunc`)
