# go-edu

Go learning repository built as a collection of **self-contained projects** under `minis/` (Go fundamentals) and `geth/` (Ethereum/go-ethereum-flavored exercises).

## Repository Structure

The repository is organized into two main tracks:

- **`minis/`**: 50 projects covering Go fundamentals from basics to advanced topics
- **`geth/`**: 25 projects focused on Ethereum development using go-ethereum

Each project follows a standard Go project layout for consistency and clarity.

## Standard Project Layout

Every project follows the same structure:

```text
<project>/
  .vscode/
    launch.json          # VS Code debug configurations
    settings.json        # VS Code workspace settings
  cmd/
    app/
      main.go           # Application entry point (CLI/server/demo)
    dev/
      main.go           # Debug harness with fixed inputs
  internal/
    <pkg>/
      exercise.go                    # Student-facing file with TODO lists
      exercise_test.go              # Test suite for the exercise
      solution.reference.go         # Complete reference implementation
      solution_no_err.reference.go # Reference without error handling
      (optional: types.go, *_bench_test.go, etc.)
  README.md             # Project-specific documentation
```

### Directory and File Explanations

#### `.vscode/`
Contains VS Code configuration files for debugging and development:
- **`launch.json`**: Debug configurations for running tests, debugging cmd/app and cmd/dev
- **`settings.json`**: Workspace-specific settings

#### `cmd/app/`
**Purpose**: Application entry point that accepts command-line arguments.

**When to use**: When you want to run the program with custom inputs or test it as a real application.

**Example usage**:
```bash
go run ./cmd/app "hello world"
go run ./cmd/app https://api.example.com 3
```

#### `cmd/dev/`
**Purpose**: Debug harness with fixed, deterministic inputs designed for stepping through code with breakpoints.

**When to use**: 
- Learning and debugging (recommended for beginners)
- Setting breakpoints and stepping through code
- Understanding how the code works step-by-step

**Why use cmd/dev**:
- Fixed inputs: No need to remember command-line arguments
- Deterministic: Same inputs every time, making debugging predictable
- Focused: Contains only essential code to test your implementation
- Breakpoint-friendly: Includes "// BREAKPOINT:" comments at key locations

#### `internal/<pkg>/exercise.go`
**Purpose**: Student-facing file where you implement your solution.

**Content**:
- Concise problem statement with key constraints
- Function signatures with complete type information
- TODO comments indicating what needs to be implemented
- Brief hints and step-by-step guidance
- References to `solution.reference.go` for detailed explanations

**Build tag**: `//go:build !solution && !reference`

**How to use**:
1. Open `exercise.go` in your editor
2. Read the problem statement and TODO comments
3. Implement the functions according to the specifications
4. Run tests: `go test ./...`
5. Debug using VS Code launch configurations

**Example TODO structure**:
```go
// FunctionName - TODO: implement this function
func FunctionName(param string) (string, error) {
    // TODO: Step 1 - Validate input
    // TODO: Step 2 - Process data
    // TODO: Step 3 - Return result
    // Refer to solution.reference.go for the complete implementation
    return "", nil
}
```

#### `internal/<pkg>/solution.reference.go`
**Purpose**: Complete reference implementation with extensive educational commentary.

**Content**:
- Full working implementation
- Detailed step-by-step explanations
- Computer science principles and concepts
- Debugging tips with breakpoint suggestions
- Memory layout and performance considerations
- Algorithm descriptions and trade-offs

**Build tag**: `//go:build reference`

**When to use**:
- After attempting the exercise yourself
- To understand the complete solution approach
- To learn debugging techniques and best practices
- To see detailed explanations of Go concepts

**How to view**: The file is excluded from normal builds but can be viewed directly or tested with:
```bash
go test -tags=reference ./...
```

#### `internal/<pkg>/exercise_test.go`
**Purpose**: Test suite that validates your implementation.

**How to use**:
```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test function
go test -v -run TestFunctionName ./...

# Run with reference implementation (to verify tests)
go test -tags=reference -v ./...
```

#### `README.md` (project-specific)
Each project includes its own README.md with:
- Project-specific requirements and context
- Prerequisites and dependencies
- Usage examples and CLI arguments
- Learning objectives

## Getting Started

### Prerequisites

- Go 1.21 or later
- VS Code with Go extension (recommended)
- Git

### Quick Start

1. **Clone the repository**:
```bash
git clone <repo-url>
cd go-edu
```

2. **Pick a project**:
```bash
cd minis/01-hello-strings
# or: cd geth/01-stack
```

3. **Open in VS Code**:
```bash
code .
```

4. **Read the exercise**:
   - Open `internal/<pkg>/exercise.go`
   - Read the problem statement and TODO comments

5. **Implement your solution**:
   - Fill in the TODO sections
   - Write your implementation

6. **Test your solution**:
```bash
go test ./...
```

7. **Debug your code**:
   - Set breakpoints in `exercise.go`
   - Press F5 and select "Debug: cmd/dev" (recommended) or "Test: internal package"
   - Step through your code using F10 (Step Over) and F11 (Step Into)

## Working with Exercises

### Understanding TODO Comments

Each `exercise.go` file contains TODO comments that guide you through the implementation:

```go
// FunctionName - TODO: implement this function
func FunctionName(input string) string {
    // TODO: Step 1 - Validate the input string
    // Hint: Check if input is empty
    
    // TODO: Step 2 - Process the input
    // Hint: Use strings package functions
    
    // TODO: Step 3 - Return the result
    // Refer to solution.reference.go for complete implementation
    
    return ""
}
```

**Guidelines for TODO comments**:
- Keep them brief and focused
- Provide step-by-step hints, not full solutions
- Reference `solution.reference.go` for detailed explanations
- Include breakpoint suggestions where helpful

### Filling Out TODO Lists

1. **Read the problem statement** at the top of `exercise.go`
2. **Understand the constraints** and requirements
3. **Review the function signatures** to understand inputs and outputs
4. **Follow the TODO comments** step by step
5. **Test frequently** with `go test ./...`
6. **Debug when stuck** using VS Code debugger

### Resetting Exercises

You can reset exercises to their initial TODO state using the `make todo` command:

```bash
# From root directory: reset all exercises
make todo all

# From geth/ directory: reset all geth exercises
cd geth
make todo all

# From a specific project: reset only that project
cd minis/01-hello-strings
make todo all

# From any location: reset specific path
make todo geth/01-stack
make todo minis/02-arrays-maps-basics
```

**Note**: This will erase any code you've written in `exercise.go` files and restore them to the initial TODO state.

## Testing and Debugging

### Testing Your Code

#### Running Tests

```bash
# Run all tests in current project
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test function
go test -v -run TestFunctionName ./...

# Run with reference implementation (to verify tests work)
go test -tags=reference -v ./...
```

#### Debugging Tests

1. Set breakpoint in test function or exercise code
2. Press F5 and select "Test: internal package"
3. Step through test execution
4. Watch variables in the Variables panel

### Debugging Workflows

#### Using cmd/dev (Recommended for Learning)

The `cmd/dev/main.go` file is a debug harness designed for stepping through code:

1. **Open** `cmd/dev/main.go` in VS Code
2. **Set breakpoints** at "// BREAKPOINT:" comments or anywhere in your code
3. **Press F5** and select "Debug: cmd/dev (Debug Harness)"
4. **Step through** using:
   - **F10** (Step Over): Execute line, don't enter functions
   - **F11** (Step Into): Enter function calls to see internals
   - **Shift+F11** (Step Out): Return to caller
   - **F5** (Continue): Run until next breakpoint
5. **Watch variables** in the Variables panel
6. **Use Debug Console** to evaluate expressions

#### Using cmd/app (CLI Arguments)

The `cmd/app/main.go` file accepts command-line arguments:

**Project-specific CLI arguments**:

**minis/ projects**:
- `01-hello-strings`: `[input_string]`
  - Example: `go run ./cmd/app "hello world"`
- `06-worker-pool-wordcount`: `[url1] [url2] ... [urlN]`
  - Example: `go run ./cmd/app https://example.com https://example.org`
- `08-http-client-retries`: `[url] [max-retries]`
  - Example: `go run ./cmd/app https://api.example.com 3`

**geth/ projects**:
- `01-stack`: `<RPC_URL> [block_number]`
  - Example: `go run ./cmd/app https://eth.llamarpc.com`
  - Example: `go run ./cmd/app https://eth.llamarpc.com 12345`
- `05-tx-nonces`: `<RPC_URL> <address>`
  - Example: `go run ./cmd/app https://eth.llamarpc.com 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045`
- `07-eth-call`: `<RPC_URL> <contract_address>`
  - Example: `go run ./cmd/app https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48`

**Debugging cmd/app**:
1. Open `.vscode/launch.json`
2. Find "Debug: cmd/app" configuration
3. Edit the `args` array to include your CLI arguments:
```json
"args": ["https://eth.llamarpc.com", "12345"]
```
4. Press F5 and select "Debug: cmd/app"
5. Set breakpoints and step through

### VS Code Debug Configurations

Each project includes `.vscode/launch.json` with these configurations:

| Configuration | Description |
|---------------|-------------|
| **Debug: cmd/app** | Debug application with CLI arguments |
| **Debug: cmd/dev** | Debug harness with fixed inputs (recommended for learning) |
| **Test: internal package** | Run all tests with debugger |
| **Test: internal package (reference)** | Run tests with solution.reference.go |

**Using launch configurations**:
1. Open any `.go` file in the project
2. Set breakpoints where you want to pause
3. Press F5 to open the debug configuration menu
4. Select the appropriate configuration
5. Use debugging controls (F10, F11, etc.)

### Tips for Effective Debugging

1. **Start with cmd/dev/main.go** - It's designed for learning
2. **Use breakpoints liberally** - Set them at function entry points
3. **Watch the Variables panel** - See how data transforms
4. **Use Call Stack panel** - Understand function call hierarchy
5. **Step Into (F11)** - Enter function implementations
6. **Step Over (F10)** - Execute line by line
7. **Step Out (Shift+F11)** - Return to caller
8. **Use Debug Console** - Evaluate expressions and inspect values
9. **Add Watch Expressions** - Monitor specific variables or expressions

### Common Issues

#### "Cannot find package"
- Run `go mod tidy` in project directory
- Ensure you're in the correct directory

#### "Build constraints exclude all Go files"
- Check build tags in exercise.go (should be `//go:build !solution && !reference`)
- Ensure you're not building with conflicting tags

#### "RPC connection failed" (geth projects)
- Check RPC URL is correct and accessible
- Try a different public RPC endpoint
- Ensure network connectivity

#### "Geth console not working" (geth/06-smart-contracts)
- Ensure Geth is installed: `geth version`
- Check Geth is running: `geth attach` should connect
- Verify RPC endpoint is accessible
- Check firewall settings

## Build Tags

The repository uses Go build tags to manage which files are compiled:

### Exercise Implementation (`exercise.go`)
```go
//go:build !solution && !reference
```
This is the default - your implementation file.

### Reference Implementation (`solution.reference.go`)
```go
//go:build reference
```
To run tests using reference implementations:
```bash
go test -tags=reference ./...
```

**How build tags work**:
- **Default build** (`go build` or `go test`): Uses `exercise.go` files
- **Reference build** (`go build -tags=reference`): Uses `solution.reference.go` files
- This ensures students work with exercise files by default, while reference implementations remain available for review

## Makefile Commands

The repository includes a Makefile with helpful commands:

### Setup & Discovery
```bash
make setup           # Initialize dependencies and verify builds
make list            # Show all available projects (both tracks)
make list-minis      # Show only minis/ projects
make list-geth       # Show only geth/ projects
```

### Running Projects
```bash
make run P=<path>    # Run specific project (auto-detects track)
                     # Examples:
                     #   make run P=minis/01-hello-strings
                     #   make run P=geth/01-stack
                     #   make run P=01-hello-strings  (assumes minis/)

make run-minis P=XX  # Run minis project explicitly
make run-geth P=XX   # Run geth project explicitly
```

### Testing
```bash
make test            # Run all tests (both tracks)
make test P=<path>   # Test specific project
                     # Examples:
                     #   make test P=minis/03-csv-stats
                     #   make test P=geth/02-rpc-basics
```

### Benchmarking
```bash
make bench           # Run all benchmarks
make bench P=<path>  # Benchmark specific project
```

### Exercise Management
```bash
# Contextual todo command - works from any directory
make todo all                    # Reset all exercises (from root)
make todo all                    # Reset all in current directory (from geth/ or minis/)
make todo all                    # Reset current project (from project directory)
make todo <path>                 # Reset specific path (from any location)
                                 # Examples:
                                 #   make todo geth/01-stack
                                 #   make todo minis/02-arrays-maps-basics
```

### Cleanup
```bash
make clean           # Clean build cache
```

### Help
```bash
make help            # Show all available commands
```

## Project Index

### minis/ (Go Fundamentals)

**Beginner** (01-10): Basics, I/O, HTTP
- [01-hello-strings](./minis/01-hello-strings/) - UTF-8 string manipulation
- [02-arrays-maps-basics](./minis/02-arrays-maps-basics/) - Arrays and maps fundamentals
- [03-csv-stats](./minis/03-csv-stats/) - CSV parsing and statistics
- [04-jsonl-log-filter](./minis/04-jsonl-log-filter/) - JSONL log filtering
- [05-cli-todo-files](./minis/05-cli-todo-files/) - CLI file operations
- [06-worker-pool-wordcount](./minis/06-worker-pool-wordcount/) - Worker pools and concurrency
- [07-generic-lru-cache](./minis/07-generic-lru-cache/) - Generic LRU cache implementation
- [08-http-client-retries](./minis/08-http-client-retries/) - HTTP client with retry logic
- [09-http-server-graceful](./minis/09-http-server-graceful/) - Graceful HTTP server shutdown
- [10-grpc-telemetry-service](./minis/10-grpc-telemetry-service/) - gRPC service implementation

**Intermediate** (11-30): Concurrency, performance, internals
- [11-slices-internals-capacity-growth](./minis/11-slices-internals-capacity-growth/) - Slice internals
- [12-pointers-zero-values-nil-gotchas](./minis/12-pointers-zero-values-nil-gotchas/) - Pointers and nil
- [13-interfaces-duck-typing](./minis/13-interfaces-duck-typing/) - Interfaces and duck typing
- [14-methods-value-vs-pointer-receivers](./minis/14-methods-value-vs-pointer-receivers/) - Method receivers
- [15-error-wrapping-sentinel-errors](./minis/15-error-wrapping-sentinel-errors/) - Error handling
- [16-context-cancellation-timeouts](./minis/16-context-cancellation-timeouts/) - Context package
- [17-file-streaming-bufio](./minis/17-file-streaming-bufio/) - File streaming with bufio
- [18-goroutines-1M-demo](./minis/18-goroutines-1M-demo/) - Goroutine scalability
- [19-channels-basics](./minis/19-channels-basics/) - Channel fundamentals
- [20-select-fanin-fanout](./minis/20-select-fanin-fanout/) - Select statement patterns
- [21-race-detection-demo](./minis/21-race-detection-demo/) - Race condition detection
- [22-worker-pool-with-backpressure](./minis/22-worker-pool-with-backpressure/) - Worker pools with backpressure
- [23-bounded-channel-semaphore](./minis/23-bounded-channel-semaphore/) - Bounded channels and semaphores
- [24-sync-mutex-vs-rwmutex](./minis/24-sync-mutex-vs-rwmutex/) - Mutex types
- [25-atomic-counters-vs-mutex](./minis/25-atomic-counters-vs-mutex/) - Atomic operations
- [26-sync-once-singleton](./minis/26-sync-once-singleton/) - Sync.Once pattern
- [27-sync-pool-allocator](./minis/27-sync-pool-allocator/) - Sync.Pool usage
- [28-pprof-cpu-mem-benchmarks](./minis/28-pprof-cpu-mem-benchmarks/) - Performance profiling
- [29-escape-analysis-inlining](./minis/29-escape-analysis-inlining/) - Compiler optimizations
- [30-build-tags-conditional-compilation](./minis/30-build-tags-conditional-compilation/) - Build tags

**Advanced** (31-50): Networking, crypto, production patterns
- [31-static-file-server](./minis/31-static-file-server/) - HTTP file server
- [32-websocket-chatroom](./minis/32-websocket-chatroom/) - WebSocket implementation
- [33-tcp-echo-server-client](./minis/33-tcp-echo-server-client/) - TCP networking
- [34-rate-limiter-token-bucket](./minis/34-rate-limiter-token-bucket/) - Rate limiting
- [35-jwt-auth-middleware](./minis/35-jwt-auth-middleware/) - JWT authentication
- [36-caching-reverse-proxy](./minis/36-caching-reverse-proxy/) - Reverse proxy with caching
- [37-http-middleware-chain](./minis/37-http-middleware-chain/) - Middleware patterns
- [38-config-loader-env-yaml](./minis/38-config-loader-env-yaml/) - Configuration management
- [39-sha256-hasher](./minis/39-sha256-hasher/) - Cryptographic hashing
- [40-merkle-tree-basics](./minis/40-merkle-tree-basics/) - Merkle tree implementation
- [41-signed-transactions-ed25519](./minis/41-signed-transactions-ed25519/) - Digital signatures
- [42-simple-block-struct-hashing](./minis/42-simple-block-struct-hashing/) - Block structures
- [43-proof-of-work-demo](./minis/43-proof-of-work-demo/) - Proof of work algorithm
- [44-mempool-in-memory](./minis/44-mempool-in-memory/) - In-memory transaction pool
- [45-p2p-gossip-mock-network](./minis/45-p2p-gossip-mock-network/) - P2P networking
- [46-generics-map-reduce](./minis/46-generics-map-reduce/) - Generic map-reduce
- [47-plugin-system-hot-reload](./minis/47-plugin-system-hot-reload/) - Plugin architecture
- [48-reflection-introspection](./minis/48-reflection-introspection/) - Reflection API
- [49-state-machine-pattern](./minis/49-state-machine-pattern/) - State machine implementation
- [50-mini-service-all-features](./minis/50-mini-service-all-features/) - Complete service example

### geth/ (Ethereum Development)

**Foundational** (01-06): Connectivity, accounts, transactions, console basics
- [01-stack](./geth/01-stack/) - Ethereum client stack
- [02-rpc-basics](./geth/02-rpc-basics/) - JSON-RPC fundamentals
- [03-keys-addresses](./geth/03-keys-addresses/) - Key and address management
- [04-accounts-balances](./geth/04-accounts-balances/) - Account and balance queries
- [05-tx-nonces](./geth/05-tx-nonces/) - Transaction nonces
- [06-smart-contracts](./geth/06-smart-contracts/) - Smart contract console tutorial
- [06-eip1559](./geth/06-eip1559/) - EIP-1559 transaction types

**Contracts** (07-10): Manual calls, abigen, events, filters
- [07-eth-call](./geth/07-eth-call/) - Contract call implementation
- [08-abigen](./geth/08-abigen/) - Contract binding generation
- [09-events](./geth/09-events/) - Event logging and parsing
- [10-filters](./geth/10-filters/) - Event filtering

**Advanced** (11-25): Storage, proofs, tracing, indexing, networking
- [11-storage](./geth/11-storage/) - Contract storage inspection
- [12-proofs](./geth/12-proofs/) - Merkle proofs
- [13-trace](./geth/13-trace/) - Transaction tracing
- [14-explorer](./geth/14-explorer/) - Block explorer functionality
- [15-receipts](./geth/15-receipts/) - Transaction receipts
- [16-concurrency](./geth/16-concurrency/) - Concurrent operations
- [17-indexer](./geth/17-indexer/) - Blockchain indexing
- [18-reorgs](./geth/18-reorgs/) - Chain reorganizations
- [19-devnets](./geth/19-devnets/) - Development networks
- [20-node](./geth/20-node/) - Node management
- [21-sync](./geth/21-sync/) - Blockchain synchronization
- [22-peers](./geth/22-peers/) - Peer management
- [23-mempool](./geth/23-mempool/) - Transaction mempool
- [24-monitor](./geth/24-monitor/) - Network monitoring
- [25-toolbox](./geth/25-toolbox/) - Utility tools

## Learning Paths

### Go Fundamentals (minis/)

Start with `minis/01-hello-strings` and progress sequentially. Each project builds on previous concepts.

**Recommended progression**:
1. **Beginner** (01-10): Master basics, I/O, and HTTP
2. **Intermediate** (11-30): Deep dive into concurrency, performance, and internals
3. **Advanced** (31-50): Build production-ready networking, crypto, and service patterns

### Ethereum Development (geth/)

Start with `geth/01-stack` and progress sequentially. Prerequisites are listed in each project's README.

**Recommended progression**:
1. **Foundational** (01-06): Learn connectivity, accounts, transactions, and console basics
2. **Contracts** (07-10): Master manual calls, abigen, events, and filters
3. **Advanced** (11-25): Explore storage, proofs, tracing, indexing, and networking

**Important Note**: For smart contract interaction, complete `geth/06-smart-contracts` (console tutorial) before `geth/07-eth-call` (Go implementation). The console experience provides essential conceptual foundation.

## Commentary Guidelines

### exercise.go Files

Keep commentary minimal and focused:
- Brief TODO comments with step-by-step hints
- Concise problem statements
- Key constraints and requirements
- References to `solution.reference.go` for details
- Breakpoint suggestions where helpful

**Example**:
```go
// FunctionName - TODO: implement this function
func FunctionName(input string) string {
    // TODO: Step 1 - Validate input
    // Hint: Check for empty string
    
    // TODO: Step 2 - Process input
    // Refer to solution.reference.go for detailed algorithm
    
    return ""
}
```

### solution.reference.go Files

Include comprehensive explanations:
- Full implementation with detailed comments
- Step-by-step algorithm explanations
- Computer science principles
- Debugging tips and breakpoint guidance
- Memory layout and performance considerations
- Trade-offs and alternatives

## Contributing

This is an educational repository. If you find issues or have suggestions:

1. Check existing issues
2. Create a new issue with details
3. Follow the project's coding standards
4. Ensure all tests pass

## License

See [LICENSE](./LICENSE) file for details.
