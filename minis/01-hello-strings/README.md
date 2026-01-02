# minis/01-hello-strings

## Problem

Problem: Implement UTF-8-aware string utilities in Go

We need three functions that correctly handle Unicode text:
1. TitleCase - Capitalize the first letter of each word
2. Reverse - Reverse a string character-by-character (not byte-by-byte!)
3. RuneLen - Count characters (runes), not bytes

Constraints:
- Must handle multi-byte UTF-8 characters (emoji, accented letters, CJK)
- Preserve all characters without corruption
- Use only the Go standard library

Time/Space Complexity:
- TitleCase: O(n) time, O(n) space (allocates new string)
- Reverse: O(n) time, O(n) space (allocates rune slice + result string)
- RuneLen: O(n) time, O(1) space (just counting)

Why Go is well-suited:
- Built-in UTF-8 support: strings are UTF-8 byte sequences by default
- Clear byte/rune distinction: prevents subtle encoding bugs
- Excellent stdlib: `unicode/utf8` and `strings` cover most needs
- Fast: no string copying overhead (immutable strings are shared internally)

DEBUGGING THIS FILE:
==================
This solution is instrumented with extensive debugging comments to teach you
how to use Go's debugger (dlv) and VS Code's debugging features.

Key debugging concepts covered:
1. Setting breakpoints at critical transformation points
2. Watching variable state changes in the Variables panel
3. Using F10 (Step Over) vs F11 (Step Into) effectively
4. Inspecting memory allocations and data structure transformations
5. Using the Debug Console to evaluate expressions
6. Understanding the call stack at each level

## Quickstart

```bash
cd minis/01-hello-strings
go test ./...
```

## CLI (cmd/app)

This project includes a small CLI wrapper around the exercise package.

### Flags

- **`-list`**: list available exported functions
- **`-fn`**: function name to run
- **`-in`**: string input (for `func(string) ...`)
- **`-n`**: int input (for `func(int) ...`)
- **`-f`**: float64 input (for `func(float64) ...`)
- **`-b`**: bool input (for `func(bool) ...`)
- **`-file`** / **`-stdin`**: input sources for `func(io.Reader) ...`

### Usage

```bash
go run ./cmd/app -h
```

### Copy/paste examples

```bash
go run ./cmd/app -list
go run ./cmd/app -fn "Reverse" -in "Hello, 世界 👋"
go run ./cmd/app -fn "RuneLen" -in "Hello, 世界 👋"
go run ./cmd/app -fn "TitleCase" -in "Hello, 世界 👋"
go run ./cmd/dev
```

## Debug harness (cmd/dev)

`cmd/dev` runs the same CLI with preselected example arguments, so you can just hit Run/Debug and see what the project can do.

```bash
go run ./cmd/dev
```

## Files

- `internal/hellostrings/exercise.go`: implement the TODOs here
- `internal/hellostrings/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
