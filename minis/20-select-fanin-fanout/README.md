# minis/20-select-fanin-fanout

## Problem

Problem: Receive from whichever channel has data first

Architecture:
- Use select to wait on multiple channels simultaneously
- Add timeout case to prevent indefinite blocking
- Return value and success indicator

Complexity:
- Time: O(1) - select operation is constant time
- Space: O(1) - no additional data structures

Three-Input Iteration Table:

Input 1: ch1 has value immediately
  ch1: "fast"
  ch2: (empty)
  select: ch1 case executes
  return: "fast", true

Input 2: Timeout occurs
  ch1: (empty)
  ch2: (empty)
  time passes...
  select: timeout case executes
  return: "", false

Input 3: ch2 sends after delay
  ch1: (empty)
  ch2: sends "slow" after 500ms
  timeout: 1 second
  select: ch2 case executes at 500ms
  return: "slow", true

## Quickstart

```bash
cd minis/20-select-fanin-fanout
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
go run ./cmd/app -fn "NewRateLimiter" -n 10
go run ./cmd/dev
```

## Debug harness (cmd/dev)

`cmd/dev` runs the same CLI with preselected example arguments, so you can just hit Run/Debug and see what the project can do.

```bash
go run ./cmd/dev
```

## Files

- `internal/selectfaninfanout/exercise.go`: implement the TODOs here
- `internal/selectfaninfanout/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
