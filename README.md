# go-edu

Go learning repository with **self-contained projects** organized into two tracks:
- **`minis/`** — Go fundamentals (50 exercises)
- **`geth/`** — Ethereum/go-ethereum development (26 exercises)

## Repository Structure

```
go-edu/
├── geth/                    # Ethereum development track
│   ├── 01-stack/
│   ├── 02-rpc-basics/
│   └── ...
├── minis/                   # Go fundamentals track
│   ├── 01-hello-strings/
│   ├── 02-arrays-maps-basics/
│   └── ...
├── .vscode/
│   └── launch.json          # Root-level debug configurations
├── Makefile                 # Build automation and exercise management
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

## Project Layout

Each project follows a standard Go project structure:

```
<project>/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application entry point
│   └── dev/
│       └── main.go          # Debug harness with fixed inputs
├── internal/
│   └── <pkg>/
│       ├── exercise.go      # YOUR CODE GOES HERE (TODO items)
│       ├── exercise_test.go # Test cases
│       ├── solution.reference.go    # Complete implementation with explanations
│       └── solution_no_err.reference.go  # (optional) Error-free variant
└── .vscode/
    └── launch.json          # Project-specific debug configurations
```

### Key Files Explained

| File | Purpose |
|------|---------|
| **`exercise.go`** | Your working file. Contains function signatures with TODO comments. Implement the functions here. |
| **`exercise_test.go`** | Test cases that validate your implementation. Run with `go test ./...` |
| **`solution.reference.go`** | Complete working implementation with detailed explanations, debugging tips, and educational commentary. Reference this when stuck. |
| **`cmd/app/main.go`** | CLI entry point that accepts command-line arguments |
| **`cmd/dev/main.go`** | Debug harness with fixed inputs—ideal for stepping through code with breakpoints |

### Build Tags

The project uses Go build tags to manage which implementation is compiled:

```go
//go:build !solution && !reference   // exercise.go (default)
//go:build reference                  // solution.reference.go
```

- **Default build** (`go build` or `go test`): Uses your `exercise.go`
- **Reference build** (`go test -tags=reference`): Uses the solution for verification

## How to Complete Exercises

### Step 1: Navigate to a Project

```bash
cd minis/01-hello-strings
# or
cd geth/01-stack
```

### Step 2: Open the Exercise File

Open `internal/<pkg>/exercise.go`. You'll see function signatures with TODO comments:

```go
// TitleCase - TODO: implement this function
func TitleCase(s string) string {
    // TODO: Implement this function
    // Refer to solution.reference.go for the complete implementation
    return ""
}
```

### Step 3: Implement the Functions

Replace the TODO placeholders with your implementation. The function signature and return type are already defined—you just need to write the logic.

### Step 4: Run Tests

```bash
# Run tests to check your implementation
go test -v ./...

# Run a specific test
go test -v -run TestTitleCase ./...
```

### Step 5: Debug If Needed

Use the debug harness for step-through debugging:

```bash
go run ./cmd/dev/main.go
```

Or use VS Code's debugger (see Debugging section below).

### Step 6: Check the Solution

If stuck, reference `solution.reference.go` for the complete implementation with detailed explanations.

## Makefile Commands

### Exercise Management

```bash
# Reset exercise files to their original TODO state
make reset T=all              # Reset all exercises
make reset T=geth             # Reset all geth exercises
make reset T=minis            # Reset all minis exercises
make reset T=geth/01-stack    # Reset specific project
```

The `make reset` command uses `git checkout` to restore exercise.go files to their original TODO state. Any uncommitted changes will be lost.

### Running Projects

```bash
make run P=minis/01-hello-strings   # Run specific project
make run P=geth/01-stack            # Run geth project
make run P=01-hello-strings         # Assumes minis/ prefix
```

### Testing

```bash
make test                           # Run all tests
make test P=minis/01-hello-strings  # Test specific project
make bench P=minis/07-generic-lru-cache  # Run benchmarks
```

### Code Quality

```bash
make lint                  # Run golangci-lint
make check                 # Run go vet and staticcheck
```

### Other Commands

```bash
make setup                 # Initialize dependencies and verify builds
make list                  # Show all available projects
make list-minis            # Show only minis projects
make list-geth             # Show only geth projects
make clean                 # Clean build cache
make help                  # Show all commands
```

## Debugging

### Using VS Code

1. **Open the project folder** in VS Code
2. **Set breakpoints** at `// BREAKPOINT:` comments or anywhere in your code
3. **Press F5** and select a debug configuration:
   - **"Debug Current Package"** — Debug the package you're currently in
   - **"Debug Tests (Current Package)"** — Debug tests with breakpoints

### Using cmd/dev/main.go

The `cmd/dev/main.go` file is a debug harness designed for learning:

- **Fixed inputs**: No need to remember CLI arguments
- **Deterministic**: Same inputs every time
- **Breakpoint-friendly**: Includes `// BREAKPOINT:` comments at key locations

```bash
# Run the debug harness
go run ./cmd/dev/main.go

# Or debug in VS Code
# 1. Open cmd/dev/main.go
# 2. Set breakpoints in exercise.go
# 3. Press F5
```

### Debugging Tips

1. **Use breakpoints liberally** — Set them at function entry points
2. **Watch variables** — Use the Variables panel to see state changes
3. **Step Into (F11)** — Enter function implementations
4. **Step Over (F10)** — Execute line by line
5. **Step Out (Shift+F11)** — Return to caller
6. **Debug Console** — Evaluate expressions like `len(s)` or `fmt.Sprintf("%#v", obj)`

### VS Code Launch Configurations

The root `.vscode/launch.json` provides these configurations:

| Configuration | Purpose |
|--------------|---------|
| **Debug Current Package** | Auto-detects if it's a program or test |
| **Debug Tests (Current Package)** | Runs tests with debugger attached |

## Learning Paths

### Go Fundamentals (`minis/`)

Progress sequentially. Each project builds on previous concepts.

| Level | Projects | Topics |
|-------|----------|--------|
| **Beginner** | 01-10 | Strings, arrays, maps, I/O, HTTP basics |
| **Intermediate** | 11-30 | Concurrency, performance, internals, sync primitives |
| **Advanced** | 31-50 | Networking, crypto, production patterns |

### Ethereum Development (`geth/`)

Start with `geth/01-stack` and progress sequentially.

| Level | Projects | Topics |
|-------|----------|--------|
| **Foundational** | 01-06 | Connectivity, accounts, transactions |
| **Contracts** | 07-10 | Manual calls, abigen, events, filters |
| **Advanced** | 11-25 | Storage, proofs, tracing, indexing, networking |

## Project Index

### minis/ (Go Fundamentals)

| # | Project | Topic |
|---|---------|-------|
| 01 | hello-strings | String operations, UTF-8 handling |
| 02 | arrays-maps-basics | Arrays, slices, maps |
| 03 | csv-stats | CSV parsing, file I/O |
| 04 | jsonl-log-filter | JSON streaming, log processing |
| 05 | cli-todo-files | CLI tools, file persistence |
| 06 | worker-pool-wordcount | Concurrency patterns |
| 07 | generic-lru-cache | Generics, caching |
| 08 | http-client-retries | HTTP client, retry logic |
| 09 | http-server-graceful | HTTP server, graceful shutdown |
| 10 | grpc-telemetry-service | gRPC, telemetry |
| 11 | slices-internals-capacity-growth | Slice internals |
| 12 | pointers-zero-values-nil-gotchas | Pointers, nil handling |
| 13 | interfaces-duck-typing | Interfaces |
| 14 | methods-value-vs-pointer-receivers | Methods, receivers |
| 15 | error-wrapping-sentinel-errors | Error handling |
| 16 | context-cancellation-timeouts | Context, cancellation |
| 17 | file-streaming-bufio | File streaming |
| 18 | goroutines-1M-demo | Goroutine scaling |
| 19 | channels-basics | Channels |
| 20 | select-fanin-fanout | Select, fan patterns |
| 21 | race-detection-demo | Race conditions |
| 22 | worker-pool-with-backpressure | Backpressure |
| 23 | bounded-channel-semaphore | Semaphores |
| 24 | sync-mutex-vs-rwmutex | Mutex types |
| 25 | atomic-counters-vs-mutex | Atomics |
| 26 | sync-once-singleton | Singleton pattern |
| 27 | sync-pool-allocator | Object pooling |
| 28 | pprof-cpu-mem-benchmarks | Profiling |
| 29 | escape-analysis-inlining | Compiler optimizations |
| 30 | build-tags-conditional-compilation | Build tags |
| 31 | static-file-server | Static file serving |
| 32 | websocket-chatroom | WebSockets |
| 33 | tcp-echo-server-client | TCP networking |
| 34 | rate-limiter-token-bucket | Rate limiting |
| 35 | jwt-auth-middleware | JWT authentication |
| 36 | caching-reverse-proxy | Reverse proxy |
| 37 | http-middleware-chain | Middleware patterns |
| 38 | config-loader-env-yaml | Configuration |
| 39 | sha256-hasher | Cryptographic hashing |
| 40 | merkle-tree-basics | Merkle trees |
| 41 | signed-transactions-ed25519 | Digital signatures |
| 42 | simple-block-struct-hashing | Block structures |
| 43 | proof-of-work-demo | Proof of work |
| 44 | mempool-in-memory | Mempool simulation |
| 45 | p2p-gossip-mock-network | P2P networking |
| 46 | generics-map-reduce | Generic algorithms |
| 47 | plugin-system-hot-reload | Plugin systems |
| 48 | reflection-introspection | Reflection |
| 49 | state-machine-pattern | State machines |
| 50 | mini-service-all-features | Full service |

### geth/ (Ethereum Development)

| # | Project | Topic |
|---|---------|-------|
| 01 | stack | Ethereum connectivity basics |
| 02 | rpc-basics | JSON-RPC fundamentals |
| 03 | keys-addresses | Key management |
| 04 | accounts-balances | Account queries |
| 05 | tx-nonces | Transaction nonces |
| 06 | smart-contracts | Contract fundamentals |
| 06 | eip1559 | EIP-1559 transactions |
| 07 | eth-call | Read-only contract calls |
| 08 | abigen | Go bindings generation |
| 09 | events | Event parsing |
| 10 | filters | Log filtering |
| 11 | storage | Storage slot access |
| 12 | proofs | Merkle proofs |
| 13 | trace | Transaction tracing |
| 14 | explorer | Block explorer |
| 15 | receipts | Receipt handling |
| 16 | concurrency | Concurrent RPC |
| 17 | indexer | Blockchain indexing |
| 18 | reorgs | Reorg handling |
| 19 | devnets | Development networks |
| 20 | node | Node interaction |
| 21 | sync | Sync monitoring |
| 22 | peers | Peer discovery |
| 23 | mempool | Mempool monitoring |
| 24 | monitor | Chain monitoring |
| 25 | toolbox | Utility collection |

## Commentary Guidelines

### exercise.go

Keep commentary minimal:
- Brief TODO comments with step-by-step hints
- Function signatures with clear parameter names
- References to `solution.reference.go` for details
- Preserve `// BREAKPOINT:` comments for debugging

### solution.reference.go

Move detailed explanations here:
- Full working implementation
- Step-by-step explanations of the algorithm
- Computer science principles
- Debugging tips with breakpoint suggestions
- Memory layout explanations
- Alternative approaches and trade-offs

## Quick Start

```bash
# 1. Clone and setup
git clone <repo>
cd go-edu
make setup

# 2. Pick your first project
cd minis/01-hello-strings

# 3. Open the exercise file
code internal/hellostrings/exercise.go

# 4. Implement the TODOs

# 5. Run tests
go test -v ./...

# 6. Debug if needed
# Set breakpoints and press F5 in VS Code
```

## Common Issues

### "Cannot find package"
```bash
go mod tidy
```

### "Build constraints exclude all Go files"
Check build tags in exercise.go (should be `//go:build !solution && !reference`)

### "RPC connection failed" (geth projects)
- Check RPC URL is accessible
- Try a different public RPC endpoint
- Export: `export INFURA_RPC_URL=https://mainnet.infura.io/v3/YOUR_KEY`

## License

See [LICENSE](./LICENSE) for details.
