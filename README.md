# Go Educational Projects: From Fundamentals to Ethereum

A comprehensive Go learning repository featuring **two progressive tracks**: 50 general Go mini-projects and 25 Ethereum/go-ethereum projects. Each project includes complete implementations, extensive documentation, and production-quality tests.

## Overview

This repository teaches Go from first principles through hands-on practice, progressing from basic string manipulation to building production-grade Ethereum tooling. Each project includes:

- **Stub exercises** (`exercise/exercise.go`) for you to implement
- **Complete reference solutions** (`exercise/solution.go` or `solution/`) with line-by-line explanations
- **Comprehensive tests** that demonstrate idiomatic Go testing patterns
- **Detailed READMEs** explaining concepts, trade-offs, and Go's unique strengths

## Two Learning Tracks

### Track 1: Go Fundamentals (minis/)
**50 progressive projects** covering Go core concepts, concurrency, networking, and backend patterns.

**Topics covered:**
- Fundamentals: strings, arrays, maps, structs, interfaces
- File I/O: CSV, JSON, JSONL parsing and processing
- Concurrency: goroutines, channels, worker pools, race detection
- HTTP: clients, servers, middleware, graceful shutdown
- Advanced: generics, reflection, plugins, state machines
- Cryptography: SHA256, Merkle trees, Ed25519, proof-of-work
- Blockchain basics: transactions, blocks, mempool, P2P gossip

### Track 2: Ethereum Development (geth/)
**25 progressive projects** for building production-grade Ethereum tooling with go-ethereum.

**Topics covered:**
- Foundation: Ethereum stack, RPC, keys, addresses, accounts
- Transactions: nonces, signing, EIP-1559 dynamic fees
- Smart Contracts: ABI encoding, abigen, typed bindings
- Events & Logs: event decoding, filters, WebSocket subscriptions
- Storage & Proofs: storage slots, Merkle Patricia tries, cryptographic verification
- Advanced: transaction tracing, indexing, chain reorgs
- Operations: node management, sync status, mempool monitoring

## Quick Start

### Prerequisites

- **Go ≥1.22** ([install guide](https://go.dev/doc/install))
- Basic command-line familiarity
- A text editor or IDE (VS Code with Go extension recommended)
- **For geth/ projects:** An Ethereum RPC endpoint (Infura, Alchemy, or local node)

### Setup

```bash
# Clone this repository
git clone <your-repo-url>
cd go-edu

# Initialize dependencies and verify all projects build
make setup

# List all available projects
make list           # List all projects
make list-minis     # List only minis projects
make list-geth      # List only geth projects
```

### Running Projects

```bash
# Run a specific project (works for both minis/ and geth/)
make run P=minis/01-hello-strings
make run P=geth/01-stack

# Shorter aliases
make run P=01-hello-strings    # Assumes minis/ if no prefix
make run-geth P=01-stack       # Explicit geth/

# Run tests
make test                      # Test all projects
make test P=minis/03-csv-stats # Test specific project
make test P=geth/02-rpc-basics

# Run benchmarks
make bench                     # Run all benchmarks
make bench P=minis/07-generic-lru-cache

# Clean build cache
make clean
```

### For Ethereum Projects (geth/)

Set up your RPC endpoint:

```bash
# Option 1: Environment variable
export INFURA_RPC_URL=https://mainnet.infura.io/v3/YOUR_PROJECT_ID

# Option 2: Create .env file
cp .env.example .env
# Edit .env with your RPC URL

# Run a geth project
INFURA_RPC_URL=... make run P=geth/01-stack
```

## Learning Paths

### Path 1: Complete Beginner
Start with minis track, then move to geth:

1. **minis/01-02**: Strings, arrays, maps (1-2 hours)
2. **minis/03-05**: CSV, JSON, file I/O (3-4 hours)
3. **minis/06**: Goroutines and channels (2 hours)
4. **minis/11-25**: Deep dive into Go internals and concurrency (8-10 hours)
5. **geth/01-10**: Ethereum fundamentals (6-8 hours)
6. Continue with advanced projects in both tracks

### Path 2: Ethereum-Focused
Already know Go? Jump to geth track:

1. **minis/01-05**: Quick Go refresher (2-3 hours)
2. **minis/06**: Concurrency patterns (2 hours)
3. **geth/01-10**: Ethereum basics (6-8 hours)
4. **geth/11-25**: Advanced Ethereum development (10-15 hours)

### Path 3: Blockchain-Focused
Learn cryptography and blockchain concepts:

1. **minis/01-05**: Go basics (2-3 hours)
2. **minis/39-45**: Cryptography and blockchain (SHA256, Merkle trees, PoW, P2P) (6-8 hours)
3. **geth/01-25**: Production Ethereum development (15-20 hours)

## Project Structure

```
go-edu/
├── README.md                  # This file - unified guide
├── Makefile                   # Common commands for both tracks
├── go.mod                     # Module definition
├── .gitignore
│
├── minis/                     # Track 1: Go Fundamentals
│   ├── 01-hello-strings/
│   │   ├── README.md                    # Project overview
│   │   ├── cmd/hello-strings/main.go    # CLI runner
│   │   └── exercise/
│   │       ├── exercise.go              # Your implementation
│   │       ├── solution.go              # Reference solution
│   │       └── exercise_test.go         # Tests
│   ├── 02-arrays-maps-basics/
│   │   ├── testdata/input.txt           # Test fixtures
│   │   └── ...
│   └── ... (projects 03-50)
│
└── geth/                      # Track 2: Ethereum Development
    ├── README.md              # Geth track overview
    ├── 01-stack/
    │   ├── README.md          # CS-first-principles explanations
    │   └── exercise/
    │       ├── exercise.go
    │       ├── solution.go
    │       └── exercise_test.go
    ├── 02-rpc-basics/
    └── ... (projects 03-25)
```

## How to Use This Repository

Each project has both `exercise.go` (stubs for you to implement) and `solution.go` (complete reference). **Choose your approach:**

### Option 1: Implement First (Recommended)

```bash
cd minis/01-hello-strings/exercise   # or geth/01-stack/exercise
mv solution.go solution.go.reference # Prevent compilation conflicts

# Read README, implement in exercise.go, then test
go test

# Compare with reference
cat solution.go.reference
```

### Option 2: Study First

```bash
cd minis/01-hello-strings/exercise
mv exercise.go exercise.go.bak

# Study solution.go, run tests
go test -tags=solution

# Reimplement from memory
rm solution.go
mv exercise.go.bak exercise.go
# Implement and test
```

### Option 3: Test-Driven Development

```bash
cd minis/01-hello-strings/exercise
mv solution.go solution.go.reference

# Read tests to understand requirements
cat exercise_test.go

# Implement just enough to pass one test at a time
go test
```

## Makefile Commands Reference

| Command | Description |
|---------|-------------|
| `make setup` | Initialize dependencies, verify builds |
| `make list` | List all projects (both tracks) |
| `make list-minis` | List only minis projects |
| `make list-geth` | List only geth projects |
| `make run P=<path>` | Run specific project |
| `make test` | Run all tests |
| `make test P=<path>` | Test specific project |
| `make bench` | Run all benchmarks |
| `make bench P=<path>` | Benchmark specific project |
| `make clean` | Clean build cache |
| `make help` | Show all available commands |

**Examples:**
```bash
make run P=minis/01-hello-strings
make run P=geth/02-rpc-basics
make test P=minis/06-worker-pool-wordcount
make bench P=minis/07-generic-lru-cache
```

## Key Go Concepts Covered

### Fundamentals (minis/01-25)
- Type system: structs, interfaces, generics
- Error handling: idiomatic returns, wrapping
- Pointers vs values, zero values
- Slices internals, capacity growth
- Maps, nil gotchas

### Standard Library Mastery
- I/O: `io.Reader`/`io.Writer`, `bufio`
- Encoding: JSON, CSV, Protocol Buffers
- HTTP: clients, servers, `httptest`
- Concurrency: goroutines, channels, `sync`, `context`

### Concurrency Patterns (minis/18-27)
- Worker pools with backpressure
- Fan-in, fan-out with `select`
- Race detection and fixes
- Mutex vs RWMutex vs atomics
- `sync.Once`, `sync.Pool`

### Performance & Profiling (minis/28-29)
- pprof: CPU and memory profiling
- Benchmarking with `testing.B`
- Escape analysis and inlining

### Ethereum Specifics (geth/01-25)
- JSON-RPC patterns and error handling
- secp256k1 cryptography, key management
- ABI encoding/decoding
- Event logs and Bloom filters
- Merkle Patricia tries and proofs
- Transaction tracing and debugging
- Chain reorganization handling
- Production monitoring patterns

## Why Go Excels

### 1. Simplicity + Power
25 keywords, yet handles concurrency, networking, and systems programming elegantly.

### 2. Deployment Joy
Single static binary → run anywhere. No runtime dependencies, no version conflicts.

### 3. Predictable Concurrency
Goroutines + channels provide structured concurrency that's easier than:
- Python's asyncio (callback complexity)
- JavaScript promises (error handling)
- Java threads (low-level, error-prone)

### 4. Standard Library Excellence
Production-grade HTTP server, JSON, database drivers — all built-in. Most projects need zero third-party dependencies.

### 5. Tooling Consistency
- `go fmt` → universal formatting
- `go test` → testing built into language
- `go mod` → dependency management
- `go vet` → static analysis

### 6. Performance + Safety
Faster than Python/Ruby/Node.js, safer than C/C++. Memory-safe without GC pauses like Java.

## Contributing

Found a bug? Have a better solution? Want to add a project?

1. Open an issue describing your idea
2. Submit a PR following existing code style
3. Ensure all tests pass: `make test`

## Additional Resources

- [Effective Go](https://go.dev/doc/effective_go) - Official style guide
- [Go by Example](https://gobyexample.com/) - Annotated examples
- [Go Blog](https://go.dev/blog/) - Deep dives by the Go team
- [go-ethereum Documentation](https://geth.ethereum.org/docs/)
- [Ethereum JSON-RPC Spec](https://ethereum.org/en/developers/docs/apis/json-rpc/)

## License

MIT License - use freely for learning, teaching, or building your own projects.

---

**Happy Coding!** Start with `make list` to see all projects, then dive in with `make run P=minis/01-hello-strings`
