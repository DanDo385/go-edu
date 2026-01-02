# minis/06-worker-pool-wordcount

## Problem

Problem: Fetch multiple URLs concurrently and count word frequencies

We need to implement:
1. Worker pool pattern for bounded concurrency
2. Concurrent HTTP fetching with proper error handling
3. Word tokenization and frequency counting
4. Graceful cancellation with context

Constraints:
- Fixed number of workers (bounded concurrency)
- Must handle errors from any worker
- Must cancel all work if one worker fails
- Aggregate results from all workers

Time/Space Complexity:
- Time: O(n) where n = total words across all URLs (concurrent fetching)
- Space: O(u) where u = number of unique words across all URLs

Why Go is well-suited:
- Goroutines: Lightweight concurrent execution
- Channels: Safe communication between goroutines
- Context: Built-in cancellation and timeout support
- sync.WaitGroup: Simple goroutine coordination
- errgroup: Simplified error handling for goroutines

DEBUGGING THIS FILE:
==================
This solution is instrumented with extensive debugging comments to teach you
how to use Go's debugger (dlv) and VS Code's debugging features.

Key debugging concepts covered:
1. Setting breakpoints in concurrent goroutines
2. Watching channel operations and synchronization
3. Understanding worker pool patterns
4. Inspecting goroutine state and coordination
5. Using the Debug Console for concurrent debugging
6. Understanding context cancellation propagation

## Quickstart

```bash
cd minis/06-worker-pool-wordcount
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
go run ./cmd/dev
```

## Debug harness (cmd/dev)

`cmd/dev` runs the same CLI with preselected example arguments, so you can just hit Run/Debug and see what the project can do.

```bash
go run ./cmd/dev
```

## Files

- `internal/workerpoolwordcount/exercise.go`: implement the TODOs here
- `internal/workerpoolwordcount/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
