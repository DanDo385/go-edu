# minis/02-arrays-maps-basics

## Problem

Problem: Count word frequencies from text input and find the most common word

Given an io.Reader containing one word per line, we need to:
1. Build a frequency map (word → count)
2. Find the word with the highest count
3. Handle errors gracefully (I/O failures, empty input)

Constraints:
- Normalize to lowercase ("Hello" == "hello")
- Ignore blank lines
- For ties, return any of the tied words (arbitrary but deterministic)

Time/Space Complexity:
- Time: O(n) where n = number of words (one pass to build map, one to find max)
- Space: O(u) where u = number of unique words (map storage)

Why Go is well-suited:
- Built-in maps with clean syntax: `map[string]int` and `count++` patterns
- `io.Reader` interface enables testing without real files
- `bufio.Scanner` handles line-by-line reading efficiently
- Zero value semantics: missing map keys return 0 (perfect for counting!)

DEBUGGING THIS FILE:
==================
This solution is instrumented with extensive debugging comments to teach you
how to use Go's debugger (dlv) and VS Code's debugging features.

Key debugging concepts covered:
1. Setting breakpoints at critical transformation points
2. Watching variable state changes in the Variables panel
3. Using F10 (Step Over) vs F11 (Step Into) effectively
4. Inspecting map data structure transformations
5. Using the Debug Console to evaluate expressions
6. Understanding the call stack at each level

## Quickstart

```bash
cd minis/02-arrays-maps-basics
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
printf "hello\nworld\n" | go run ./cmd/app -fn "FreqFromReader" -stdin
go run ./cmd/dev
```

## Debug harness (cmd/dev)

`cmd/dev` runs the same CLI with preselected example arguments, so you can just hit Run/Debug and see what the project can do.

```bash
go run ./cmd/dev
```

## Files

- `internal/arraysmapsbasics/exercise.go`: implement the TODOs here
- `internal/arraysmapsbasics/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
