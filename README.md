# Go 10x Minis: A Practical Learning Journey

A comprehensive Go learning repository featuring **10 progressively challenging mini-projects** from beginner to advanced, each with complete implementations, extensive documentation, and production-quality tests.

## 🎯 Overview

This repository is designed to take you from Go fundamentals to real-world backend patterns through hands-on practice. Each project includes:

- **Stub exercises** (`exercise/exercise.go`) for you to implement
- **Complete reference solutions** (`exercise/solution.go`) with line-by-line explanations
- **Comprehensive tests** that demonstrate idiomatic Go testing patterns
- **Detailed READMEs** explaining concepts, trade-offs, and Go's unique strengths

## 🚀 Quick Start

### Prerequisites

- **Go ≥1.22** ([install guide](https://go.dev/doc/install))
- Basic command-line familiarity
- A text editor or IDE (VS Code with Go extension recommended)

### Setup

```bash
# Clone this repository
git clone https://github.com/example/go-10x-minis.git
cd go-10x-minis

# Initialize dependencies and verify all projects build
make setup

# List all available projects
make list
```

### Running Projects

```bash
# Run a specific project
make run P=01-hello-strings

# Run all tests
make test

# Run benchmarks (projects 7 and 10)
make bench

# Clean build cache
make clean
```

## 📚 The Learning Path

### Beginner (Projects 1-2)

**01. hello-strings** - String manipulation fundamentals
Learn: `strings` package, UTF-8 handling, table-driven tests
Time: 30 minutes

**02. arrays-maps-basics** - Data structures and file I/O
Learn: Slices, maps, `bufio`, sorting
Time: 45 minutes

### Easy-Medium (Projects 3-5)

**03. csv-stats** - Structured data processing
Learn: `encoding/csv`, structs, streaming I/O, error handling
Time: 1 hour

**04. jsonl-log-filter** - JSON parsing and custom types
Learn: `encoding/json`, custom type marshaling, sorting algorithms
Time: 1 hour

**05. cli-todo-files** - Persistent CLI application
Learn: `flag` package, JSON persistence, interfaces, file handling
Time: 1.5 hours

### Medium-Hard (Projects 6-8)

**06. worker-pool-wordcount** - Concurrency patterns
Learn: Goroutines, channels, worker pools, `context` cancellation
Time: 2 hours

**07. generic-lru-cache** - Advanced data structures
Learn: Generics, `sync.Mutex`, `container/list`, thread safety
Time: 2.5 hours

**08. http-client-retries** - Resilient HTTP clients
Learn: `net/http`, exponential backoff, jitter, context timeouts
Time: 2 hours

### Advanced (Projects 9-10)

**09. http-server-graceful** - Production HTTP servers
Learn: `http.ServeMux`, middleware, graceful shutdown, OS signals
Time: 2.5 hours

**10. grpc-telemetry-service** - Modern RPC and streaming
Learn: gRPC, protobuf, bidirectional streaming, service architecture
Time: 3-4 hours

## 🎓 Core Concepts Covered

### Fundamentals
- **Type system**: Structs, interfaces, type parameters (generics)
- **Error handling**: Idiomatic error returns, error wrapping with `fmt.Errorf`
- **Pointers & values**: When to use pointers vs. value semantics
- **Zero values**: Go's predictable initialization

### Standard Library Mastery
- **I/O**: `io.Reader`/`io.Writer`, `bufio`, file operations
- **Encoding**: JSON, CSV, Protocol Buffers
- **Text processing**: `strings`, `unicode/utf8`, regular expressions
- **HTTP**: Client and server patterns, `httptest` for testing
- **Concurrency**: Goroutines, channels, `sync` primitives, `context`

### Testing & Quality
- **Table-driven tests**: The idiomatic Go testing pattern
- **Subtests**: `t.Run()` for organized test output
- **Benchmarking**: `testing.B` for performance measurement
- **Test fixtures**: Using `testdata/` and `t.TempDir()`
- **HTTP testing**: `httptest.Server` for deterministic tests

### Production Patterns
- **Worker pools**: Bounded concurrency for resource management
- **Graceful shutdown**: Clean termination with signal handling
- **Retry logic**: Exponential backoff with jitter
- **Caching**: LRU eviction with TTL support
- **Middleware**: Composable HTTP request processing

## 🌟 Why Go Excels

### 1. **Simplicity + Power**
Go has just 25 keywords, yet handles concurrency, networking, and systems programming elegantly. Compare to Rust's steep learning curve or Python's GIL limitations.

### 2. **Deployment Joy**
Compile once → single static binary → run anywhere. No Docker layers with Python interpreters, no Node.js version conflicts, no JVM memory overhead.

### 3. **Predictable Concurrency**
Goroutines + channels provide structured concurrency that's easier to reason about than:
- Python's asyncio (callback hell)
- JavaScript promises (error handling complexity)
- Java threads (low-level, error-prone)

### 4. **Standard Library Excellence**
`net/http` includes a production-grade HTTP server **out of the box**. `encoding/json`, `database/sql`, `crypto/*` — all first-class. Most projects need zero third-party dependencies.

### 5. **Tooling Consistency**
- `go fmt` → universal formatting (no debates!)
- `go test` → testing built into the language
- `go mod` → dependency management without `package.json` chaos
- `go vet` → static analysis for free

### 6. **Performance + Safety**
Faster than Python/Ruby/Node.js, safer than C/C++. Memory-safe without garbage collection pauses like Java. Compile times measured in seconds, not minutes (looking at you, Rust).

## 🛠️ Project Structure

```
go-10x-minis/
├── README.md              # This file
├── go.mod                 # Module definition
├── Makefile               # Convenience commands
├── .gitignore
└── minis/
    ├── 01-hello-strings/
    │   ├── README.md                    # Project overview
    │   ├── cmd/hello-strings/main.go    # CLI runner
    │   └── exercise/
    │       ├── exercise.go              # Stub for you to implement
    │       ├── solution.go              # Reference implementation
    │       └── exercise_test.go         # Comprehensive tests
    ├── 02-arrays-maps-basics/
    │   ├── testdata/input.txt           # Test fixtures
    │   └── ...
    └── ... (projects 3-10)
```

## 💡 How to Use This Repository

**IMPORTANT**: Each project has both `exercise.go` (stubs for you to implement) and `solution.go` (complete reference). Since both define the same functions, **you must choose one approach**:

### Option 1: Implement First (Recommended)
1. **Rename `solution.go`** to `solution.go.reference` (prevents compilation conflicts)
2. Read the project's `README.md` to understand requirements
3. Implement your solution in `exercise/exercise.go` (replace the `TODO` comments)
4. Run tests: `go test ./minis/<project>/...`
5. Compare your code with `solution.go.reference` to learn alternative approaches

### Option 2: Study the Solution
1. **Rename `exercise.go`** to `exercise.go.bak`
2. Read `solution.go` carefully, understanding each comment
3. Run tests to see it working: `go test ./minis/<project>/...`
4. Delete `solution.go`, restore `exercise.go.bak`, and reimplement from memory

### Option 3: Test-Driven Development
1. **Rename `solution.go`** to `solution.go.reference`
2. Read tests in `exercise_test.go` to understand expected behavior
3. Implement just enough in `exercise.go` to make one test pass
4. Repeat until all tests pass
5. Refactor and compare with `solution.go.reference`

**Quick Start Command** (to prepare project 01 for implementation):
```bash
cd minis/01-hello-strings/exercise
mv solution.go solution.go.reference
# Now implement in exercise.go and run: go test
```

## 🤝 Contributing

Found a bug? Have a better solution? Want to add an 11th project?

1. Open an issue describing your idea
2. Submit a PR following the existing code style
3. Ensure all tests pass: `make test`

## 📖 Additional Resources

- [Effective Go](https://go.dev/doc/effective_go) - Official style guide
- [Go by Example](https://gobyexample.com/) - Annotated code examples
- [Go Blog](https://go.dev/blog/) - Deep dives by the Go team
- [Concurrency Patterns](https://go.dev/talks/2012/concurrency.slide) - Rob Pike's classic talk

## 📜 License

MIT License - feel free to use this for learning, teaching, or as a foundation for your own projects.

---

**Happy Coding!** 🎉 Start with `make run P=01-hello-strings` and work your way up!
