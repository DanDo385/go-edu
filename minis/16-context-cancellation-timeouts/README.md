# minis/16-context-cancellation-timeouts

## Problem

Problem: Retry operations with timeout and exponential backoff

We need to:
1. Retry a function up to maxRetries times
2. Each attempt has its own timeout
3. Use exponential backoff between retries
4. Respect parent context cancellation

Architecture:
- Loop from 0 to maxRetries
- For each attempt:
  - Create child context with timeout
  - Call function with child context
  - If success, return
  - If failure and not last attempt, backoff with context check
  - If parent context cancelled, return immediately

Complexity:
- Time: O(maxRetries * timeout + sum of backoffs)
- Space: O(1) - only stores error

Three-Input Iteration Table:

Input 1: Success on 3rd attempt
  Attempt 0: fn() fails, wait 100ms
  Attempt 1: fn() fails, wait 200ms
  Attempt 2: fn() succeeds → return nil

Input 2: All attempts fail
  Attempt 0: fn() fails, wait 100ms
  Attempt 1: fn() fails, wait 200ms
  Attempt 2: fn() fails → return error

Input 3: Parent context cancelled during backoff
  Attempt 0: fn() fails, wait 100ms
  Parent cancelled → return context.Canceled

## Quickstart

```bash
cd minis/16-context-cancellation-timeouts
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
go run ./cmd/app -fn "NewCache"
go run ./cmd/app -fn "NewRateLimiter" -n 10
go run ./cmd/dev
```

## Debug harness (cmd/dev)

`cmd/dev` runs the same CLI with preselected example arguments, so you can just hit Run/Debug and see what the project can do.

```bash
go run ./cmd/dev
```

## Files

- `internal/contextcancellationtimeouts/exercise.go`: implement the TODOs here
- `internal/contextcancellationtimeouts/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
