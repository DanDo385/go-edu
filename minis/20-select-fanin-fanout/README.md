# 20-select-fanin-fanout

**Select, Fan-In, Fan-Out**

Master select statement and fan-in/fan-out patterns.

## What You'll Learn

- select for multiplexing channels
- Fan-out: one source, many workers
- Fan-in: many sources, one destination
- Non-blocking operations

## Functions to Implement

| Function | Description |
|----------|-------------|
| Implement fan patterns | Select-based multiplexing |

## Project Structure

```
20-select-fanin-fanout/
├── cmd/
│   ├── app/main.go      # CLI demonstration
│   └── dev/main.go      # Debug harness
├── internal/selectfaninfanout/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/20-select-fanin-fanout

# Run demonstration
go run ./cmd/app/main.go

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Run demo
go run ./cmd/app/main.go

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **select**: Wait on multiple channels
2. **Fan-Out**: Distribute work to workers
3. **Fan-In**: Merge results from workers
4. **default case**: Non-blocking select

## Next Steps

After completing this exercise, proceed to `minis/21-race-detection-demo`.
