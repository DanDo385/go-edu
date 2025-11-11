# Testing Guide for Go Mini Projects

This guide explains how to test your solutions for each mini project.

## Quick Start

Each project comes with comprehensive test files located in the `exercise/` directory. You can run tests individually or use the convenient test runner script.

### Method 1: Using the Test Runner Script (Recommended)

```bash
# Test a specific project (using project number)
./test-runner.sh 1           # Test project 01-hello-strings
./test-runner.sh 06          # Test project 06-worker-pool-wordcount

# Test all projects at once
./test-runner.sh all

# Test with verbose output (shows each test case)
./test-runner.sh 1 -v

# Run benchmarks (for projects that have them)
./test-runner.sh 7 --bench   # Run benchmarks for LRU cache
./test-runner.sh all --bench # Run all tests + benchmarks
```

### Method 2: Using `go test` Directly

Navigate to any project's `exercise/` directory and run:

```bash
cd minis/01-hello-strings/exercise
go test                    # Run all tests
go test -v                 # Verbose output (see each test case)
go test -run TestReverse   # Run only tests matching "TestReverse"
```

## Understanding Test Output

### ✅ Passing Tests

```
$ ./test-runner.sh 1
Testing minis/01-hello-strings/exercise
✓ Tests passed
```

### ❌ Failing Tests

When tests fail, you'll see detailed error messages:

```
$ go test
--- FAIL: TestReverse (0.00s)
    --- FAIL: TestReverse/with_emoji (0.00s)
        exercise_test.go:131: Reverse("Hello 👋 World") = "dlroW  olleH", want "dlroW 👋 olleH"
FAIL
FAIL    example.com/minis/01-hello-strings/exercise    0.001s
```

This tells you:
- Which test failed: `TestReverse/with_emoji`
- What input was used: `"Hello 👋 World"`
- What you returned: `"dlroW  olleH"`
- What was expected: `"dlroW 👋 olleH"`

## Project-by-Project Test Overview

### Project 01: hello-strings

**What's tested:**
- `TitleCase()` - Capitalizing first letter of each word
- `Reverse()` - Reversing strings (UTF-8 aware)
- `RuneLen()` - Counting characters (not bytes)

**Key test cases:**
- ASCII strings (simple cases)
- Unicode emoji (🎉, 👋)
- Accented characters (café, résumé)
- Japanese characters (日本語)
- Edge cases (empty strings, single characters)

**Run tests:**
```bash
./test-runner.sh 1
# or
cd minis/01-hello-strings/exercise && go test
```

### Project 02: arrays-maps-basics

**What's tested:**
- `FreqFromReader()` - Word frequency counting

**Key test cases:**
- Simple word lists
- Case insensitivity (Go, GO, go → all counted as "go")
- Whitespace handling
- Empty input
- Unicode words

**Run tests:**
```bash
./test-runner.sh 2
```

### Project 03: csv-stats

**What's tested:**
- CSV parsing with custom delimiters
- Statistical calculations (count, sum, average, min, max)
- Column handling

**Key test cases:**
- Standard comma-separated CSV
- Tab-separated values
- Missing/empty values
- Numeric parsing

**Run tests:**
```bash
./test-runner.sh 3
```

### Project 04: jsonl-log-filter

**What's tested:**
- JSONL (newline-delimited JSON) parsing
- Log level filtering (ERROR, WARN, INFO, DEBUG)
- Timestamp formatting

**Key test cases:**
- Multiple log levels
- Malformed JSON handling
- Empty input
- Edge cases

**Run tests:**
```bash
./test-runner.sh 4
```

### Project 05: cli-todo-files

**What's tested:**
- Todo CRUD operations (Create, Read, Update, Delete)
- File persistence
- Status tracking (pending/completed)

**Key test cases:**
- Adding todos
- Listing todos
- Completing todos
- Deleting todos
- File I/O operations

**Run tests:**
```bash
./test-runner.sh 5
```

### Project 06: worker-pool-wordcount

**What's tested:**
- Concurrent HTTP requests
- Worker pool with bounded parallelism
- Word counting across multiple URLs
- Context cancellation
- Error handling

**Key test cases:**
- Basic word counting from multiple servers
- Empty URL list
- Context timeout/cancellation
- Server errors (500 status)
- Punctuation removal
- Case insensitivity
- Multiple workers (concurrency testing)

**Run tests:**
```bash
./test-runner.sh 6
```

### Project 07: generic-lru-cache

**What's tested:**
- Generic type parameters
- LRU eviction policy
- Get/Put operations
- TTL (time-to-live)
- Thread safety

**Key test cases:**
- Basic Get/Put operations
- Eviction when cache is full
- LRU ordering (least recently used evicted first)
- TTL expiration
- Generic types (string/int, int/string, etc.)

**Run tests:**
```bash
./test-runner.sh 7

# Also includes benchmarks!
./test-runner.sh 7 --bench
```

**Benchmark output example:**
```
BenchmarkCacheGet-8        5000000      250 ns/op      0 B/op   0 allocs/op
BenchmarkCachePut-8        3000000      400 ns/op     50 B/op   2 allocs/op
```

### Project 08: http-client-retries

**What's tested:**
- Exponential backoff algorithm
- Jitter (randomization)
- Retry logic
- Error classification (retryable vs permanent)

**Key test cases:**
- Successful request (no retries needed)
- Transient failures (503, network errors)
- Permanent failures (404, 400)
- Max retries exceeded
- Context cancellation
- Backoff timing verification

**Run tests:**
```bash
./test-runner.sh 8
```

### Project 09: http-server-graceful

**What's tested:**
- HTTP request handling
- Middleware composition
- Logging middleware
- Graceful shutdown

**Key test cases:**
- Basic request handling
- Middleware execution order
- Request logging
- Clean shutdown without dropping requests

**Run tests:**
```bash
./test-runner.sh 9
```

### Project 10: grpc-telemetry-service

**What's tested:**
- gRPC server setup
- Client streaming
- Server streaming
- Bidirectional streaming
- Time-windowed data aggregation

**Key test cases:**
- Sending telemetry data points
- Retrieving statistics (count, sum, avg, min, max)
- Multiple metrics
- Time window behavior
- Concurrent client streams

**Run tests:**
```bash
./test-runner.sh 10

# Also includes benchmarks!
./test-runner.sh 10 --bench
```

## Development Workflow

### 1. Read the Project README

Each project has a comprehensive README explaining:
- What the project is about
- First principles behind the concepts
- Complete solution walkthrough
- Key concepts explained
- Real-world applications

### 2. Implement Your Solution

Edit the `exercise/exercise.go` file in each project:

```go
// Before (starter code)
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

### 3. Run the Tests

```bash
./test-runner.sh 1
```

### 4. Fix Failures

Read the error messages carefully:
- They show what input was tested
- What your function returned
- What was expected

### 5. Iterate Until All Tests Pass

Once you see:
```
✓ Tests passed
```

You've successfully completed the project! 🎉

## Advanced Testing Techniques

### Running Specific Tests

```bash
cd minis/01-hello-strings/exercise

# Run only the Reverse tests
go test -run TestReverse

# Run only the "with emoji" subtest
go test -run TestReverse/with_emoji

# Run tests with verbose output
go test -v -run TestReverse
```

### Running Tests in Parallel

```bash
# Run tests across multiple CPU cores (faster for large test suites)
go test -parallel 4
```

### Checking Test Coverage

```bash
# See which parts of your code are tested
go test -cover

# Generate detailed coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Running Benchmarks

For projects with benchmarks (07, 10):

```bash
cd minis/07-generic-lru-cache/exercise

# Run benchmarks
go test -bench=.

# Run benchmarks with memory allocation stats
go test -bench=. -benchmem

# Run specific benchmark
go test -bench=BenchmarkCacheGet
```

## Troubleshooting

### Tests Won't Run

**Problem:** `go: cannot find main module`

**Solution:** Make sure you're in the `exercise/` directory:
```bash
cd minis/01-hello-strings/exercise
go test
```

### Import Errors

**Problem:** `package X is not in GOROOT`

**Solution:** Run `go mod tidy` to download dependencies:
```bash
cd minis/10-grpc-telemetry-service/exercise
go mod tidy
go test
```

### Tests Timeout

**Problem:** Tests hang or timeout (common in projects 06, 08, 09, 10)

**Solution:** Check for:
- Infinite loops
- Missing channel closes
- Deadlocks in goroutines
- Context not being passed correctly

Add timeout flag:
```bash
go test -timeout 30s
```

## Understanding Table-Driven Tests

Most tests use the "table-driven" pattern, which is idiomatic in Go:

```go
func TestTitleCase(t *testing.T) {
    tests := []struct {
        name string  // Test case name
        in   string  // Input
        want string  // Expected output
    }{
        {
            name: "simple lowercase words",
            in:   "hello world",
            want: "Hello World",
        },
        {
            name: "empty string",
            in:   "",
            want: "",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := TitleCase(tt.in)
            if got != tt.want {
                t.Errorf("TitleCase(%q) = %q, want %q", tt.in, got, tt.want)
            }
        })
    }
}
```

**Benefits:**
- Easy to add new test cases (just add to the slice)
- Clear separation of test data and test logic
- Subtests can be run individually: `go test -run TestTitleCase/empty_string`
- Failures show exactly which test case failed

## Test Quality Standards

All test files follow these standards:

✅ **Comprehensive coverage:**
- Happy path (normal inputs)
- Edge cases (empty, nil, boundary values)
- Error cases (invalid input, network failures)
- Unicode/internationalization (emoji, non-ASCII)

✅ **Clear test names:**
- Use descriptive subtest names
- Name shows what is being tested

✅ **Good error messages:**
- Show input, output, and expected values
- Use `%q` for strings (shows quotes and whitespace)

✅ **Proper cleanup:**
- Close servers, files, connections
- Use `defer` for cleanup
- Cancel contexts

## Next Steps

1. **Start with Project 01:** It's the simplest and teaches core string handling
2. **Work sequentially:** Projects build on each other conceptually
3. **Read the README first:** Each project has detailed explanations
4. **Run tests frequently:** Get fast feedback on your implementation
5. **Study the reference solution:** After passing tests, compare with `solution.go`

## Getting Help

If you're stuck:

1. **Read the test failure message carefully** - it usually tells you exactly what's wrong
2. **Check the project README** - it has detailed explanations
3. **Look at the reference solution** - `exercise/solution.go` has a complete implementation with comments
4. **Run tests verbosely** - `go test -v` shows all test cases, even passing ones

Happy coding! 🚀
