# Quick Start Guide

Welcome to Go Mini Projects! This guide will get you started in 5 minutes.

## Running Tests

Each project has comprehensive tests. Here's how to use them:

### Test Your Solution

```bash
# Test a specific project
./test-runner.sh 1          # Test project 01-hello-strings
./test-runner.sh 06         # Test project 06-worker-pool-wordcount

# Test all projects
./test-runner.sh all

# Verbose output (see all test cases)
./test-runner.sh 1 -v
```

### Test Reference Solution

Want to see if the tests work correctly? Run them against the reference solution:

```bash
./test-runner.sh 1 --solution    # Test project 01 with reference solution
./test-runner.sh all --solution  # Test all projects with reference solutions
```

### Run Benchmarks

Some projects (07, 10) include performance benchmarks:

```bash
./test-runner.sh 7 --bench            # Run benchmarks for project 07
./test-runner.sh 7 --solution --bench # Run benchmarks for reference solution
```

## Workflow

### 1. Pick a Project

Start with Project 01 (hello-strings) and work sequentially:

```bash
cd minis/01-hello-strings/
cat README.md  # Read the comprehensive guide
```

### 2. Write Your Solution

Edit `exercise/exercise.go`:

```go
// Before
func TitleCase(s string) string {
    // TODO: implement
    return ""
}

// After (your implementation)
func TitleCase(s string) string {
    words := strings.Fields(s)
    for i, word := range words {
        runes := []rune(word)
        if len(runes) > 0 {
            runes[0] = unicode.ToUpper(runes[0])
        }
        words[i] = string(runes)
    }
    return strings.Join(words, " ")
}
```

### 3. Run Tests

```bash
cd exercise/
go test          # Run tests against your code
go test -v       # See detailed output
```

Or from the repository root:

```bash
./test-runner.sh 1
```

### 4. Fix Failures

Tests show you:
- What input was tested
- What you returned
- What was expected

Example:
```
exercise_test.go:67: TitleCase("hello world") = "", want "Hello World"
                      ^^^^^^^^^ ^^^^^^^^^^^^    ^^        ^^^^^^^^^^^^
                      Function  Input          Your       Expected
                                                Output    Output
```

### 5. Verify Success

When all tests pass:

```
PASS
ok  	github.com/example/go-10x-minis/minis/01-hello-strings/exercise	0.007s
```

### 6. Compare with Reference

After passing tests, check the reference solution:

```bash
# View reference solution
cat exercise/solution.go

# Test reference solution
go test -tags=solution
```

## Alternative: Manual Testing

You can also use `go test` commands directly:

```bash
cd minis/01-hello-strings/exercise

# Basic test
go test

# Verbose output
go test -v

# Run specific test
go test -run TestReverse

# Run specific subtest
go test -run TestReverse/with_emoji

# Check coverage
go test -cover

# Test reference solution
go test -tags=solution

# Run benchmarks (projects 07, 10)
go test -bench=. -benchmem
```

## Project Structure

```
minis/01-hello-strings/
├── README.md              # Comprehensive learning guide (300+ lines)
├── exercise/
│   ├── exercise.go        # Your implementation (TODO stubs)
│   ├── exercise_test.go   # Test suite (runs against exercise.go by default)
│   └── solution.go        # Reference solution (only compiled with -tags=solution)
└── cmd/                   # Optional: runnable programs (some projects)
```

## Getting Help

1. **Read the README**: Each project has a detailed guide explaining everything
2. **Run tests verbosely**: `go test -v` shows all test cases
3. **Check test failures**: Error messages tell you exactly what's wrong
4. **Study reference solution**: `exercise/solution.go` has detailed comments
5. **See full testing guide**: Read `TESTING_GUIDE.md` for advanced techniques

## Quick Reference

| Command | Description |
|---------|-------------|
| `./test-runner.sh 1` | Test your solution for project 01 |
| `./test-runner.sh 1 -v` | Test with verbose output |
| `./test-runner.sh 1 --solution` | Test reference solution |
| `./test-runner.sh all` | Test all projects |
| `./test-runner.sh 7 --bench` | Run benchmarks |
| `cd minis/01-hello-strings/exercise && go test` | Direct test command |
| `go test -run TestReverse` | Run specific test |
| `go test -cover` | Check test coverage |

## What You'll Learn

- **Project 01-05**: Go basics (strings, slices, maps, files, JSON, CSV)
- **Project 06**: Concurrency (goroutines, channels, worker pools)
- **Project 07**: Generics and data structures (LRU cache)
- **Project 08**: Network resilience (retries, exponential backoff)
- **Project 09**: HTTP servers (middleware, graceful shutdown)
- **Project 10**: gRPC and Protocol Buffers (streaming RPC)

## Next Steps

1. Read `TESTING_GUIDE.md` for comprehensive testing documentation
2. Start with Project 01: `cd minis/01-hello-strings && cat README.md`
3. Work through projects sequentially
4. Check your understanding by completing stretch goals
5. Build real-world projects using these patterns!

Happy coding! 🚀
