# Quick Start Guide: Project 26 - sync.Once and Singleton Pattern

## Running the Demo

View all demonstrations of sync.Once and singleton patterns:

```bash
cd /home/user/go-edu/minis/26-sync-once-singleton
go run ./cmd/once-demo/main.go
```

### Run with Race Detector (recommended!)

See how the race detector catches concurrency bugs:

```bash
go run -race ./cmd/once-demo/main.go
```

The race detector will show warnings for the naive implementation, demonstrating why sync.Once is needed.

## Working on Exercises

### 1. Run Tests (to see what needs to be implemented)

```bash
cd /home/user/go-edu/minis/26-sync-once-singleton/exercise
go test -v
```

You'll see failures showing which exercises need implementation.

### 2. Implement the Exercises

Open `exercise/exercise.go` and implement the TODOs:

- Exercise 1: Basic sync.Once with Counter
- Exercise 2: Configuration Singleton  
- Exercise 3: Database Singleton with Error Handling
- Exercise 4: Logger Singleton
- Exercise 5: Cache Singleton
- Exercise 6: Metrics Singleton
- Exercise 7: Lazy Field Initialization
- Exercise 8: Idempotent Initialization
- Exercise 9: Resettable Once (for testing)
- Exercise 10: Factory Singleton

### 3. Test Your Implementation

```bash
go test -v
```

### 4. Test with Race Detector

```bash
go test -race
```

This catches concurrency bugs in your implementation.

### 5. Run Benchmarks

```bash
go test -bench=. -benchmem
```

See the performance characteristics of your singleton implementations.

## Viewing the Solution

If you get stuck, you can run tests with the solution:

```bash
go test -tags solution -v
```

Or view the solution file directly:

```bash
cat exercise/solution.go
```

## Learning Path

1. **Read the README.md** - Understand sync.Once from first principles
2. **Run the demo** - See sync.Once in action
3. **Run with -race** - Understand why thread-safety matters
4. **Work through exercises** - Build muscle memory
5. **Read the solution** - Learn best practices

## Key Commands Reference

```bash
# Run demo
go run ./cmd/once-demo/main.go

# Run demo with race detector
go run -race ./cmd/once-demo/main.go

# Run tests (your implementation)
go test ./exercise/...

# Run tests with race detector
go test -race ./exercise/...

# Run tests with solution
go test -tags solution ./exercise/...

# Run benchmarks
go test -bench=. -benchmem ./exercise/...

# Verbose test output
go test -v ./exercise/...
```

## What You'll Learn

- ✅ How sync.Once guarantees exactly-once execution
- ✅ Implementing the singleton pattern correctly
- ✅ Lazy initialization for expensive resources
- ✅ Thread-safe initialization without races
- ✅ Error handling in singleton initialization
- ✅ Performance optimization (fast path)
- ✅ Memory ordering and visibility guarantees
- ✅ Testing strategies for global state

## Next Steps

After completing this project, you'll be ready for:

- Project 27: sync.Map (concurrent map implementations)
- Project 28: Worker pools with initialization
- Advanced concurrency patterns in production code

Happy learning! 🚀
