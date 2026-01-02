# minis/07-generic-lru-cache

## Problem

Problem: Implement a thread-safe LRU cache with generics and TTL

Requirements:
1. O(1) Get and Set operations
2. Thread-safe (concurrent access from multiple goroutines)
3. LRU eviction when capacity is reached
4. Optional per-item TTL expiration
5. Generic over key and value types

Data Structure:
- Map: key → list element (O(1) lookup)
- Doubly-linked list: maintains recency order (front = most recent)
- Mutex: protects concurrent access

Time/Space Complexity:
- Get: O(1) average (map lookup + list move)
- Set: O(1) average (map insert + list append/evict)
- Space: O(capacity) for map + list

Algorithm: LRU Eviction Policy
- Track access order in doubly-linked list
- Most recent items at front
- Least recent items at back
- Evict from back when capacity exceeded

Why doubly-linked list:
- O(1) move to front (mark as recently used)
- O(1) remove from back (evict LRU item)
- O(1) remove arbitrary element (for updates)

## Quickstart

```bash
cd minis/07-generic-lru-cache
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

- `internal/genericlrucache/exercise.go`: implement the TODOs here
- `internal/genericlrucache/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
