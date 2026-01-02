# minis/08-http-client-retries

## Problem

Problem: Build a resilient HTTP client with retries and exponential backoff

Requirements:
1. Retry failed requests automatically
2. Exponential backoff: delay increases exponentially
3. Jitter: add randomness to prevent thundering herd
4. Context-aware: respect timeouts and cancellation
5. Generic JSON decoding

Algorithm: Exponential Backoff with Jitter
- Attempt request
- If fails and retryable: wait backoff duration, retry
- If fails and non-retryable: return error immediately
- Backoff formula: BaseDelay * 2^attempt ± 20% jitter
- Repeat up to MaxRetries times

Why Exponential Backoff:
- Linear backoff (1s, 2s, 3s) is too aggressive
- Exponential (100ms, 200ms, 400ms, 800ms) backs off quickly
- Gives servers time to recover from overload
- Industry standard for distributed systems

Why Jitter (Randomness):
- Prevents thundering herd (synchronized retries)
- Spreads retries over time
- Reduces load spikes on recovering servers

Time Complexity: O(retries * request_time)
Space Complexity: O(1)

## Quickstart

```bash
cd minis/08-http-client-retries
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

- `internal/httpclientretries/exercise.go`: implement the TODOs here
- `internal/httpclientretries/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
