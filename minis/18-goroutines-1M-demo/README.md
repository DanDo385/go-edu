# 18-goroutines-1M-demo

**Million Goroutines Demo**

Demonstrate Go's lightweight goroutines by spawning a million.

## What You'll Learn

- Goroutine memory footprint (~2KB stack)
- Scaling to millions of goroutines
- Runtime scheduling
- sync.WaitGroup for coordination

## Functions to Implement

| Function | Description |
|----------|-------------|
| Spawn and coordinate goroutines | Million goroutine demo |

## Project Structure

```
18-goroutines-1M-demo/
├── cmd/
│   ├── app/main.go      # CLI demonstration
│   └── dev/main.go      # Debug harness
├── internal/goroutines1mdemo/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/18-goroutines-1M-demo

# Spawn goroutines
go run ./cmd/app/main.go --count 1000000

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Spawn 1 million goroutines
go run ./cmd/app/main.go --count 1000000

# Smaller test
go run ./cmd/app/main.go --count 10000

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **Goroutine Stack**: Starts at 2KB, grows as needed
2. **M:N Scheduling**: Many goroutines on few OS threads
3. **sync.WaitGroup**: Wait for all to complete
4. **GOMAXPROCS**: Number of OS threads

## Next Steps

After completing this exercise, proceed to `minis/19-channels-basics`.
