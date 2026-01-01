# go-edu

A comprehensive Go learning repository built as a collection of **self-contained projects** under `minis/` (Go fundamentals) and `geth/` (Ethereum/go-ethereum exercises).

## Table of Contents

- [Overview](#overview)
- [Repository Structure](#repository-structure)
- [Project Organization](#project-organization)
- [File Structure Explained](#file-structure-explained)
- [Getting Started](#getting-started)
- [How to Complete Exercises](#how-to-complete-exercises)
- [Testing and Debugging](#testing-and-debugging)
- [Makefile Commands](#makefile-commands)
- [Project Index](#project-index)
- [Learning Paths](#learning-paths)

---

## Overview

This repository contains **75 self-contained Go projects** (50 minis + 25 geth) designed to teach Go programming through hands-on exercises. Each project follows a consistent structure with exercise files, reference solutions, and comprehensive debugging support.

### Two Learning Tracks

1. **`minis/`** - Go fundamentals (50 projects)
   - Basic syntax and data structures
   - Concurrency patterns
   - HTTP servers and clients
   - Performance optimization
   - Cryptography and blockchain concepts

2. **`geth/`** - Ethereum development with go-ethereum (25 projects)
   - RPC connectivity and blockchain queries
   - Transaction handling
   - Smart contract interaction
   - Event monitoring and indexing
   - Node operations and networking

---

## Repository Structure

```
go-edu/
├── minis/                  # Go fundamentals track (50 projects)
│   ├── 01-hello-strings/
│   ├── 02-arrays-maps-basics/
│   ├── ...
│   └── 50-mini-service-all-features/
├── geth/                   # Ethereum development track (25 projects)
│   ├── 01-stack/
│   ├── 02-rpc-basics/
│   ├── ...
│   └── 25-toolbox/
├── scripts/                # Utility scripts
│   └── reset-exercises.sh  # Reset exercise files to TODO state
├── Makefile                # Project automation commands
├── README.md               # This file
└── LICENSE
```

---

## Project Organization

Every project follows the **standard Go project layout** with consistent structure:

```
<project>/
├── .vscode/
│   ├── launch.json        # VS Code debug configurations
│   └── settings.json      # Project-specific settings
├── cmd/
│   ├── app/
│   │   └── main.go        # CLI application entry point
│   └── dev/
│       └── main.go        # Debug harness with fixed inputs
├── internal/
│   └── <package>/
│       ├── exercise.go            # Your implementation goes here (TODO stubs)
│       ├── exercise_test.go       # Test cases
│       ├── solution.reference.go  # Complete reference solution with extensive commentary
│       └── (optional: types.go, helpers, benchmarks, etc.)
```

### Directory Purpose

- **`.vscode/`** - Per-project debug/run/test configurations for VS Code
- **`cmd/app/`** - Production-style entry point with CLI argument parsing
- **`cmd/dev/`** - Deterministic debug harness with fixed inputs (recommended for learning)
- **`internal/<pkg>/`** - The exercise package containing your implementation and tests

---

## File Structure Explained

### 1. `exercise.go` - Your Work Goes Here

**Purpose**: Student-facing file with minimal TODO lists and function signatures

**What's in it**:
- Concise problem statement
- Complete function signatures with proper types
- Brief TODO comments indicating what to implement
- Zero-value return statements (placeholders)
- Build tag: `//go:build !solution && !reference`

**Size**: Typically 10-60 lines (vs 150-700 lines in solution files)

**Example**:
```go
//go:build !solution && !reference

package hellostrings

/*
Problem: Implement UTF-8-aware string utilities in Go
Constraints:
- Must handle multi-byte UTF-8 characters (emoji, accented letters, CJK)
- Preserve all characters without corruption
- Use only the Go standard library
*/

// TitleCase - TODO: implement this function
func TitleCase(s string) string {
    // TODO: Implement this function
    // Refer to solution.reference.go for the complete implementation with detailed explanations
    return ""
}
```

**Your workflow**:
1. Read the problem statement
2. Implement the function body
3. Run tests to verify your solution
4. Compare with `solution.reference.go` if stuck

---

### 2. `solution.reference.go` - Complete Reference Implementation

**Purpose**: Fully working implementation with extensive educational commentary

**What's in it**:
- Complete working implementation
- Step-by-step explanations of the algorithm
- Computer science principles and patterns
- Debugging instructions with "BREAKPOINT:" comments
- Memory layout explanations
- Performance considerations
- Go idioms and best practices
- Build tag: `//go:build reference`

**Size**: 150-700 lines with comprehensive commentary

**When to use it**:
- When you're stuck on an implementation
- To learn debugging techniques
- To understand the "why" behind design decisions
- To see professional Go code organization

**Running tests with reference solution**:
```bash
go test -tags=reference ./...
```

---

### 3. `exercise_test.go` - Test Cases

**Purpose**: Comprehensive test suite to verify your implementation

**What's in it**:
- Table-driven tests (idiomatic Go testing pattern)
- Edge cases and boundary conditions
- Error condition testing
- Unicode/UTF-8 test cases
- Performance benchmarks (optional)

**Example test structure**:
```go
func TestTitleCase(t *testing.T) {
    tests := []struct {
        name string
        in   string
        want string
    }{
        {name: "simple lowercase", in: "hello world", want: "Hello World"},
        {name: "with emoji", in: "hello 👋 world", want: "Hello 👋 World"},
        // ... more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := TitleCase(tt.in)
            if got != tt.want {
                t.Errorf("got %q, want %q", got, tt.want)
            }
        })
    }
}
```

---

### 4. `cmd/app/main.go` - CLI Application

**Purpose**: Production-style entry point with command-line argument parsing

**What's in it**:
- Argument validation and error handling
- Usage instructions
- Integration with your exercise package
- Real-world input/output handling

**Example usage**:
```bash
# minis examples
go run ./cmd/app "hello world"
go run ./cmd/app https://example.com

# geth examples  
go run ./cmd/app https://eth.llamarpc.com
go run ./cmd/app https://eth.llamarpc.com 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb
```

**When to use it**:
- Testing with real-world inputs
- Building CLI tools
- Production-style error handling
- Demonstrating the completed project

---

### 5. `cmd/dev/main.go` - Debug Harness (Recommended for Learning)

**Purpose**: Deterministic debug environment with fixed inputs and breakpoint guidance

**What's in it**:
- Fixed test inputs (no CLI arguments to remember)
- Detailed console output explaining each step
- "BREAKPOINT:" comments at key locations
- Step-by-step debugging instructions
- Educational explanations inline

**Why use this instead of cmd/app**:
- ✅ **Deterministic** - Same inputs every run
- ✅ **No arguments** - Just press F5 to debug
- ✅ **Learning-focused** - Explains what's happening at each step
- ✅ **Breakpoint-friendly** - Pre-marked debugging locations
- ✅ **Quick iteration** - Faster than typing CLI arguments

**How to use**:
1. Open `cmd/dev/main.go` in VS Code
2. Set breakpoints at lines marked with `// BREAKPOINT:`
3. Press F5 and select "Debug: cmd/dev"
4. Step through using F10 (Step Over) and F11 (Step Into)
5. Watch variables transform in the Variables panel

**Example debug harness**:
```go
func main() {
    fmt.Println("=== Debug Harness ===")
    
    // Fixed test inputs
    // BREAKPOINT: Inspect these values
    input := "hello world"
    
    // BREAKPOINT: Step into this function
    result := hellostrings.TitleCase(input)
    
    // BREAKPOINT: Inspect the result
    fmt.Printf("Result: %s\n", result)
}
```

---

### 6. `.vscode/launch.json` - Debug Configurations

**Purpose**: Pre-configured debug settings for VS Code

**Available configurations**:

| Configuration | Purpose | When to Use |
|--------------|---------|-------------|
| **Debug: cmd/dev** | Debug harness with fixed inputs | **Recommended for learning** - No setup needed |
| **Debug: cmd/app** | Debug CLI with arguments | Testing with real inputs |
| **Test: internal package** | Run and debug tests | Debugging test failures |
| **Test: internal package (reference)** | Run tests with solution | Verify reference implementation works |

**How to use**:
1. Press F5 in VS Code
2. Select configuration from dropdown
3. Set breakpoints in your code
4. Step through with F10/F11

**Customizing cmd/app debugging**:
Edit `launch.json` to add CLI arguments:
```json
{
    "name": "Debug: cmd/app",
    "type": "go",
    "request": "launch",
    "mode": "debug",
    "program": "${workspaceFolder}/cmd/app",
    "args": ["https://eth.llamarpc.com", "12345"]  // Add your arguments here
}
```

---

## Getting Started

### Prerequisites

- **Go 1.21+** installed
- **VS Code** with Go extension (recommended)
- **Git** for cloning
- **Internet connection** (for geth projects requiring RPC access)

### Installation

```bash
# Clone the repository
git clone <repo-url>
cd go-edu

# Initialize dependencies
make setup

# List all available projects
make list
```

### Your First Project

```bash
# Navigate to first project
cd minis/01-hello-strings

# View the exercise file
cat internal/hellostrings/exercise.go

# Run tests (they will fail initially)
go test ./...

# Option 1: Implement in exercise.go and rerun tests
# Option 2: View reference solution
cat internal/hellostrings/solution.reference.go

# Run with debug harness
go run ./cmd/dev

# Or debug with VS Code (F5 → "Debug: cmd/dev")
```

---

## How to Complete Exercises

### Step 1: Read the Problem

Open `internal/<pkg>/exercise.go` and read the problem statement at the top:

```go
/*
Problem: Implement UTF-8-aware string utilities in Go
Constraints:
- Must handle multi-byte UTF-8 characters
- Use only Go standard library
Time/Space Complexity:
- O(n) time, O(n) space
*/
```

### Step 2: Locate the TODO Comments

Find functions with TODO comments:

```go
// TitleCase - TODO: implement this function
func TitleCase(s string) string {
    // TODO: Implement this function
    // Refer to solution.reference.go for the complete implementation with detailed explanations
    return ""  // Replace this with your implementation
}
```

### Step 3: Implement the Function

Replace the TODO and return statement with your implementation:

```go
func TitleCase(s string) string {
    words := strings.Fields(s)
    for i, word := range words {
        if len(word) > 0 {
            runes := []rune(word)
            runes[0] = unicode.ToUpper(runes[0])
            words[i] = string(runes)
        }
    }
    return strings.Join(words, " ")
}
```

### Step 4: Run Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -v -run TestTitleCase ./...
```

### Step 5: Debug if Tests Fail

**Option A: Use cmd/dev debug harness** (Recommended)
1. Press F5 in VS Code
2. Select "Debug: cmd/dev"
3. Set breakpoints in your implementation
4. Step through with F10/F11

**Option B: Debug tests directly**
1. Set breakpoint in test or implementation
2. Press F5 → "Test: internal package"
3. Step through test execution

**Option C: Check reference solution**
```bash
# View the complete solution
cat internal/<pkg>/solution.reference.go

# Run tests with reference implementation
go test -tags=reference -v ./...
```

### Step 6: Verify and Move On

Once all tests pass:
```bash
✓ All tests passed
```

Move to the next project or explore optimizations.

---

## Testing and Debugging

### Running Tests

```bash
# From project root - test everything
go test ./...

# From specific project directory
cd minis/01-hello-strings
go test ./...

# With verbose output
go test -v ./...

# Run specific test function
go test -v -run TestTitleCase ./...

# Run with reference implementation
go test -tags=reference -v ./...

# Run benchmarks
go test -bench=. -benchmem ./...
```

### Using the Debug Harness (cmd/dev/main.go)

**Best for**: Learning, stepping through code, understanding data flow

**Workflow**:
1. **Open** `cmd/dev/main.go` in VS Code
2. **Read** the setup - it shows you the test inputs
3. **Set breakpoints** at lines marked `// BREAKPOINT:`
4. **Press F5** and select "Debug: cmd/dev"
5. **Step through**:
   - F10 - Step Over (execute current line)
   - F11 - Step Into (enter function)
   - Shift+F11 - Step Out (return to caller)
6. **Watch Variables** panel to see data transform
7. **Use Debug Console** to evaluate expressions

**Example debug session**:
```go
// BREAKPOINT 1: Set here - inspect input
input := "hello world"

// BREAKPOINT 2: Set here - step into TitleCase
result := hellostrings.TitleCase(input)

// BREAKPOINT 3: Set here - inspect result
fmt.Printf("Result: %s\n", result)
```

### Debugging Tests

**Workflow**:
1. Set breakpoint in test function or implementation
2. Press F5 → "Test: internal package"
3. Step through test execution
4. Inspect test table values in Variables panel

### VS Code Debugging Tips

1. **Variables Panel** - Shows all local variables and their values
2. **Watch Panel** - Add expressions to monitor (e.g., `len(slice)`, `cap(slice)`)
3. **Call Stack** - See function call hierarchy
4. **Debug Console** - Evaluate Go expressions interactively
5. **Breakpoint Types**:
   - Regular breakpoint - Stop every time
   - Conditional breakpoint - Stop when condition is true
   - Log point - Print message without stopping

### Common Debugging Patterns

**Pattern 1: Inspect function entry**
```go
func TitleCase(s string) string {
    // BREAKPOINT: Set here, inspect 's' value and length
    words := strings.Fields(s)
    // ...
}
```

**Pattern 2: Before/after transformation**
```go
// BREAKPOINT: Set here, note 'words' before transformation
for i, word := range words {
    runes := []rune(word)
    runes[0] = unicode.ToUpper(runes[0])
    words[i] = string(runes)
    // BREAKPOINT: Set here, see word transformation
}
```

**Pattern 3: Inspect return values**
```go
result := doSomething()
// BREAKPOINT: Set here, verify result is correct
return result
```

---

## Makefile Commands

### Available Commands

```bash
# Setup and discovery
make setup              # Initialize dependencies and verify builds
make list               # Show all available projects
make list-minis         # Show only minis/ projects
make list-geth          # Show only geth/ projects

# Running projects
make run P=<project>    # Run specific project (auto-detects track)
                        # Examples:
                        #   make run P=minis/01-hello-strings
                        #   make run P=geth/01-stack
                        #   make run P=01-hello-strings  (assumes minis/)

# Testing
make test               # Run all tests
make test P=<project>   # Test specific project
                        # Example: make test P=minis/03-csv-stats

# Benchmarking
make bench              # Run all benchmarks
make bench P=<project>  # Benchmark specific project

# Cleanup
make clean              # Clean build cache

# Exercise management (reset to TODO state)
make todo               # Context-aware reset (see below)
make todo P=<path>      # Reset specific project

# Help
make help               # Show all commands
```

### Context-Aware `make todo` Command

The `make todo` command resets `exercise.go` files back to their initial TODO state, erasing any code you've written. This is useful when:
- You want to start over on an exercise
- You're preparing exercises for students
- You want to practice an exercise again

**Behavior depends on current directory**:

```bash
# From repository root
make todo               # Resets ALL exercise.go files (minis + geth)
make todo P=minis       # Resets all minis/ exercises
make todo P=geth        # Resets all geth/ exercises

# From geth/ directory
cd geth
make todo               # Resets all geth/ exercises only

# From specific project directory
cd minis/01-hello-strings
make todo               # Resets only this project's exercise.go

# From any location with explicit path
make todo P=geth/01-stack           # Resets geth/01-stack only
make todo P=minis/01-hello-strings  # Resets minis/01-hello-strings only
```

**What gets reset**:
- `internal/<pkg>/exercise.go` files are regenerated with TODO stubs
- Function signatures are preserved
- All your implementation code is erased
- `solution.reference.go` files are NEVER modified

**Safety**:
⚠️ **Warning**: This command will erase your work! Commit your changes to git before resetting.

**Under the hood**:
The reset process:
1. Reads `solution.reference.go` to extract function signatures
2. Generates minimal TODO-based `exercise.go` file
3. Preserves imports, types, and function signatures
4. Adds appropriate zero-value return statements

---

## Project Index

### minis/ (50 projects)

Go fundamentals, from basics to advanced patterns.

#### Basics (1-10)
- **01-hello-strings** - String manipulation, UTF-8 handling
- **02-arrays-maps-basics** - Collections and data structures
- **03-csv-stats** - File I/O and data processing
- **04-jsonl-log-filter** - JSON parsing and filtering
- **05-cli-todo-files** - Building CLI applications
- **06-worker-pool-wordcount** - Concurrency patterns
- **07-generic-lru-cache** - Generics and caching
- **08-http-client-retries** - HTTP clients with retry logic
- **09-http-server-graceful** - HTTP servers with graceful shutdown
- **10-grpc-telemetry-service** - gRPC services

#### Internals & Language Features (11-20)
- **11-slices-internals-capacity-growth** - Slice memory model
- **12-pointers-zero-values-nil-gotchas** - Pointer semantics
- **13-interfaces-duck-typing** - Interface design
- **14-methods-value-vs-pointer-receivers** - Method receivers
- **15-error-wrapping-sentinel-errors** - Error handling patterns
- **16-context-cancellation-timeouts** - Context package
- **17-file-streaming-bufio** - Streaming I/O
- **18-goroutines-1M-demo** - Goroutine scaling
- **19-channels-basics** - Channel fundamentals
- **20-select-fanin-fanout** - Channel patterns

#### Concurrency & Performance (21-30)
- **21-race-detection-demo** - Data races
- **22-worker-pool-with-backpressure** - Advanced worker pools
- **23-bounded-channel-semaphore** - Concurrency primitives
- **24-sync-mutex-vs-rwmutex** - Locking strategies
- **25-atomic-counters-vs-mutex** - Atomic operations
- **26-sync-once-singleton** - Singleton pattern
- **27-sync-pool-allocator** - Object pooling
- **28-pprof-cpu-mem-benchmarks** - Profiling and optimization
- **29-escape-analysis-inlining** - Compiler optimizations
- **30-build-tags-conditional-compilation** - Build tags

#### Networking & Web (31-38)
- **31-static-file-server** - File servers
- **32-websocket-chatroom** - WebSocket communication
- **33-tcp-echo-server-client** - TCP networking
- **34-rate-limiter-token-bucket** - Rate limiting
- **35-jwt-auth-middleware** - Authentication
- **36-caching-reverse-proxy** - Reverse proxies
- **37-http-middleware-chain** - Middleware patterns
- **38-config-loader-env-yaml** - Configuration management

#### Cryptography & Blockchain (39-50)
- **39-sha256-hasher** - Hashing
- **40-merkle-tree-basics** - Merkle trees
- **41-signed-transactions-ed25519** - Digital signatures
- **42-simple-block-struct-hashing** - Block structure
- **43-proof-of-work-demo** - Consensus algorithms
- **44-mempool-in-memory** - Transaction pools
- **45-p2p-gossip-mock-network** - P2P networking
- **46-generics-map-reduce** - Generic algorithms
- **47-plugin-system-hot-reload** - Plugin architecture
- **48-reflection-introspection** - Reflection
- **49-state-machine-pattern** - State machines
- **50-mini-service-all-features** - Complete service example

---

### geth/ (25 projects)

Ethereum development with go-ethereum library.

#### Foundation (01-06)
- **01-stack** - RPC connectivity, chain ID, network ID
- **02-rpc-basics** - Making RPC calls
- **03-keys-addresses** - Key generation and addresses
- **04-accounts-balances** - Querying account balances
- **05-tx-nonces** - Transaction nonce management
- **06-smart-contracts** - Smart contract fundamentals (console tutorial)

#### Smart Contracts (07-10)
- **06-eip1559** - EIP-1559 transaction types
- **07-eth-call** - Contract view calls
- **08-abigen** - Go bindings generation
- **09-events** - Event logs and filtering
- **10-filters** - Log filtering

#### Advanced Queries (11-13)
- **11-storage** - Storage slot reading
- **12-proofs** - Merkle proofs
- **13-trace** - Transaction tracing

#### Indexing & Monitoring (14-19)
- **14-explorer** - Block explorer
- **15-receipts** - Transaction receipts
- **16-concurrency** - Concurrent RPC calls
- **17-indexer** - Event indexer
- **18-reorgs** - Chain reorganizations
- **19-devnets** - Local development networks

#### Node Operations (20-25)
- **20-node** - Running a node
- **21-sync** - Sync status monitoring
- **22-peers** - P2P peer management
- **23-mempool** - Mempool monitoring
- **24-monitor** - Network monitoring
- **25-toolbox** - Utility toolkit

---

## Learning Paths

### Path 1: Go Fundamentals (minis/)

**Beginner** (1-10): Basics, I/O, HTTP
- Start with `01-hello-strings` for basic syntax
- Progress through data structures, file I/O, and HTTP
- Learn concurrency basics with worker pools

**Intermediate** (11-30): Concurrency, performance, internals
- Deep dive into Go's memory model
- Master goroutines, channels, and synchronization
- Learn profiling and optimization

**Advanced** (31-50): Networking, crypto, production patterns
- Build production-grade services
- Implement blockchain concepts
- Design plugin architectures

### Path 2: Ethereum Development (geth/)

**Prerequisites**: Complete minis 01-10 for Go fundamentals

**Foundational** (01-06): Connectivity, accounts, transactions
- Establish RPC connections
- Generate keys and query balances
- Understand transaction lifecycle

**Contracts** (07-10): Smart contract interaction
- Make view calls to contracts
- Generate Go bindings with abigen
- Monitor events and logs

**Advanced** (11-25): Storage, indexing, node operations
- Read storage slots directly
- Build event indexers
- Monitor chain reorgs
- Operate and monitor nodes

### Recommended Order

1. **Start with minis/01-hello-strings** - Get comfortable with the project structure
2. **Complete minis 01-10** - Build Go fundamentals
3. **Try geth/01-stack** - Your first Ethereum project
4. **Alternate between tracks** - Reinforce learning
5. **Revisit projects** - Use `make todo` to practice

---

## Build Tags

This repository uses Go build tags to manage which implementation is compiled:

### Default Build (Exercise Mode)
```bash
go build ./...
go test ./...
```
Uses `exercise.go` files (your implementations).

### Reference Build (Solution Mode)
```bash
go build -tags=reference ./...
go test -tags=reference ./...
```
Uses `solution.reference.go` files (complete implementations).

### Build Tag Syntax

**Exercise files** (`exercise.go`):
```go
//go:build !solution && !reference
```

**Solution files** (`solution.reference.go`):
```go
//go:build reference
```

This ensures the correct file is compiled based on build flags.

---

## Common Issues

### "Cannot find package"
```bash
cd <project-directory>
go mod tidy
```

### "Build constraints exclude all Go files"
Check build tags in `exercise.go`:
```go
//go:build !solution && !reference
```

### "RPC connection failed" (geth projects)
- Verify RPC URL is accessible
- Try alternative public endpoints:
  - `https://eth.llamarpc.com`
  - `https://cloudflare-eth.com`
  - `https://rpc.ankr.com/eth`

### Tests fail immediately
- Ensure you've implemented the function (not just the TODO stub)
- Check return types match function signature
- Run with `-v` flag to see detailed error messages

### VS Code debugger not working
- Install Go extension for VS Code
- Ensure `dlv` (Delve debugger) is installed: `go install github.com/go-delve/delve/cmd/dlv@latest`
- Check `.vscode/launch.json` exists in project directory

---

## Tips for Success

### Learning Strategy
1. **Read problem statements carefully** - Understand requirements before coding
2. **Start with cmd/dev** - Use debug harness for learning
3. **Write tests first** - Think about edge cases
4. **Use breakpoints liberally** - See data transform step-by-step
5. **Compare with reference** - Learn idiomatic patterns

### Debugging Workflow
1. **Set breakpoint at function entry** - Inspect inputs
2. **Step through transformations** - Watch data change
3. **Check return values** - Verify correctness
4. **Use Watch panel** - Monitor specific expressions
5. **Read reference solutions** - Learn debugging techniques

### Best Practices
- Commit your work before using `make todo`
- Complete projects sequentially for progressive learning
- Revisit earlier projects to reinforce concepts
- Read `solution.reference.go` for production patterns
- Use `make test` frequently to catch errors early

---

## Contributing

This repository is designed for education. If you find issues or have suggestions:

1. Ensure tests pass with reference implementation: `go test -tags=reference ./...`
2. Check your implementation against `solution.reference.go`
3. Use `make todo` to reset if needed
4. Review build tags and project structure

---

## License

See [LICENSE](./LICENSE) file for details.

---

## Quick Reference Card

```bash
# Setup
make setup                      # Initialize project
make list                       # See all projects

# Navigate
cd minis/01-hello-strings      # Go to project
cd geth/01-stack               # Go to geth project

# Implement
# Edit internal/<pkg>/exercise.go
# Fill in TODO comments

# Test
go test ./...                   # Run tests
go test -v ./...               # Verbose output
go test -tags=reference ./...  # Test reference solution

# Debug
go run ./cmd/dev               # Run debug harness
# Or: Press F5 in VS Code → "Debug: cmd/dev"

# Reset
make todo                       # Reset exercises (context-aware)

# Reference
cat internal/<pkg>/solution.reference.go  # View solution
```

**Happy learning! Start with `minis/01-hello-strings` and work your way up.** 🚀
