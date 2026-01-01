# go-edu

**Learn Go by building real systems.**

This repository is a comprehensive educational resource for mastering Go through hands-on projects that mirror real-world production patterns and Ethereum/blockchain tooling.

## 🎯 Philosophy

**go-edu** teaches Go through two complementary tracks:

1. **`minis/`** — 50 projects covering fundamental Go patterns, concurrency primitives, HTTP servers, and crypto basics
2. **`geth/`** — 25 Ethereum-focused projects inspired by go-ethereum architecture and patterns

Every project is self-contained, production-ready in structure, and designed for deep understanding through debugging and experimentation.

---

## 📂 Repository Structure

```
go-edu/
├── go.mod                    # Single module at root
├── go.sum                    # Dependencies
├── .gitignore                # Build artifacts excluded
├── LICENSE                   # Repository license
├── README.md                 # This file (single source of truth)
├── minis/                    # 50 mini-projects covering core Go
│   ├── 01-hello-strings/
│   │   ├── .vscode/
│   │   │   └── launch.json   # VS Code debug configurations
│   │   ├── cmd/
│   │   │   ├── app/          # Realistic CLI application
│   │   │   │   └── main.go
│   │   │   └── dev/          # Debug harness (fixed inputs)
│   │   │       └── main.go
│   │   └── internal/
│   │       └── hellostrings/ # Self-contained package
│   │           ├── exercise.go            # YOUR implementation
│   │           ├── exercise_test.go       # Tests to verify
│   │           ├── solution.reference.go  # Reference (excluded)
│   │           └── solution_no_err.reference.go
│   └── ...
└── geth/                     # 25 Ethereum/Geth-focused projects
    ├── 01-stack/
    │   ├── .vscode/
    │   │   └── launch.json
    │   ├── cmd/
    │   │   ├── app/          # Realistic RPC client app
    │   │   │   └── main.go
    │   │   └── dev/          # Debug harness
    │   │       └── main.go
    │   └── internal/
    │       └── stack/
    │           ├── exercise.go
    │           ├── exercise_test.go
    │           ├── types.go              # Geth-style type definitions
    │           ├── solution.reference.go
    │           └── solution_no_err.reference.go
    └── ...
```

---

## 🏗️ Project Structure Explained

Every project follows the **self-contained layout**:

### Directory Layout

| Directory/File | Purpose |
|----------------|---------|
| `.vscode/launch.json` | VS Code debug configurations for the project |
| `cmd/app/main.go` | **Realistic application** — accepts CLI arguments, simulates production usage |
| `cmd/dev/main.go` | **Debug harness** — fixed, deterministic inputs for stepping through with a debugger |
| `internal/<pkg>/exercise.go` | **The only buildable implementation** — all production logic lives here |
| `internal/<pkg>/*_test.go` | Tests and benchmarks |
| `internal/<pkg>/*.reference.go` | Reference solutions (excluded from normal builds) |

### Build Tags and Implementation Model

#### Core Invariant: One Buildable Implementation

- `exercise.go` is the **only** non-test, non-reference file that participates in builds
- All production logic must be consolidated into `exercise.go`
- This enforces clarity: there is one source of truth for implementation

#### Build Tag Rules

**exercise.go files:**
```go
//go:build !solution && !reference

package mypackage
```

**Reference files (*.reference.go):**
```go
//go:build reference

package mypackage
```

> **Note:** Only use `//go:build` (the modern syntax). Do NOT include the legacy `// +build` line.

#### Reference Files Are Inert

Reference implementations exist **only** for learning:

```bash
# Normal build: only compiles exercise.go
go build ./...
go test ./...

# View reference (does not affect tests):
go build -tags=reference ./...
```

Reference files include:
- `solution.reference.go` — Complete, production-quality reference implementation
- `solution_no_err.reference.go` — Alternative reference (e.g., simplified error handling)

---

## 🚀 Getting Started

### Prerequisites

- **Go 1.24+** (or version specified in `go.mod`)
- For `geth/` projects: access to an Ethereum RPC endpoint (e.g., [Infura](https://infura.io), [Alchemy](https://alchemy.com), or public endpoints like `https://eth.llamarpc.com`)

### Quick Start

```bash
# Clone the repository
git clone <repository-url> go-edu
cd go-edu

# Run tests for a specific project
cd minis/01-hello-strings
go test ./...

# Run the debug harness (fixed inputs, great for setting breakpoints)
go run ./cmd/dev

# Run the application (realistic usage)
go run ./cmd/app "your input here"

# Run benchmarks
go test -bench=. ./...

# Build the application binary
go build -o app ./cmd/app
./app "test input"
```

---

## 📚 Project Index

### `minis/` — Core Go Fundamentals (50 Projects)

#### Fundamentals (01-05)
| Project | Description |
|---------|-------------|
| **01-hello-strings** | String manipulation, UTF-8, runes |
| **02-arrays-maps-basics** | Arrays, slices, maps |
| **03-csv-stats** | CSV parsing, file I/O |
| **04-jsonl-log-filter** | JSONL parsing, filtering |
| **05-cli-todo-files** | File operations, CLI |

#### HTTP & Networking (06-10)
| Project | Description |
|---------|-------------|
| **06-worker-pool-wordcount** | Concurrency, worker pools |
| **07-generic-lru-cache** | Generics, LRU caching |
| **08-http-client-retries** | HTTP client, retries |
| **09-http-server-graceful** | HTTP server, graceful shutdown |
| **10-grpc-telemetry-service** | gRPC, Protocol Buffers |

#### Deep Dives (11-17)
| Project | Description |
|---------|-------------|
| **11-slices-internals-capacity-growth** | Slice internals |
| **12-pointers-zero-values-nil-gotchas** | Pointers, nil |
| **13-interfaces-duck-typing** | Interfaces, duck typing |
| **14-methods-value-vs-pointer-receivers** | Method receivers |
| **15-error-wrapping-sentinel-errors** | Error handling |
| **16-context-cancellation-timeouts** | Context, cancellation |
| **17-file-streaming-bufio** | Streaming, bufio |

#### Concurrency Patterns (18-27)
| Project | Description |
|---------|-------------|
| **18-goroutines-1M-demo** | Goroutines at scale |
| **19-channels-basics** | Channels fundamentals |
| **20-select-fanin-fanout** | Select, fan-in, fan-out |
| **21-race-detection-demo** | Race detection |
| **22-worker-pool-with-backpressure** | Backpressure |
| **23-bounded-channel-semaphore** | Semaphores |
| **24-sync-mutex-vs-rwmutex** | Mutex vs RWMutex |
| **25-atomic-counters-vs-mutex** | Atomics vs Mutex |
| **26-sync-once-singleton** | sync.Once, singletons |
| **27-sync-pool-allocator** | sync.Pool, object pooling |

#### Performance & Profiling (28-30)
| Project | Description |
|---------|-------------|
| **28-pprof-cpu-mem-benchmarks** | pprof, benchmarking |
| **29-escape-analysis-inlining** | Escape analysis |
| **30-build-tags-conditional-compilation** | Build tags |

#### Advanced HTTP (31-38)
| Project | Description |
|---------|-------------|
| **31-static-file-server** | Static files |
| **32-websocket-chatroom** | WebSockets |
| **33-tcp-echo-server-client** | TCP networking |
| **34-rate-limiter-token-bucket** | Rate limiting |
| **35-jwt-auth-middleware** | JWT authentication |
| **36-caching-reverse-proxy** | Reverse proxy, caching |
| **37-http-middleware-chain** | Middleware patterns |
| **38-config-loader-env-yaml** | Configuration |

#### Cryptography & Blockchain (39-45)
| Project | Description |
|---------|-------------|
| **39-sha256-hasher** | SHA256 hashing |
| **40-merkle-tree-basics** | Merkle trees |
| **41-signed-transactions-ed25519** | Digital signatures |
| **42-simple-block-struct-hashing** | Block structures |
| **43-proof-of-work-demo** | Proof of Work |
| **44-mempool-in-memory** | Mempool implementation |
| **45-p2p-gossip-mock-network** | P2P networking |

#### Advanced Topics (46-50)
| Project | Description |
|---------|-------------|
| **46-generics-map-reduce** | Generics, map/reduce |
| **47-plugin-system-hot-reload** | Plugin systems |
| **48-reflection-introspection** | Reflection |
| **49-state-machine-pattern** | State machines |
| **50-mini-service-all-features** | Complete microservice |

---

### `geth/` — Ethereum & Go-Ethereum Patterns (25 Projects)

#### Foundation (01-06)
| Project | Description |
|---------|-------------|
| **01-stack** | Ethereum stack overview, RPC connectivity |
| **02-rpc-basics** | Chain ID, network ID, headers |
| **03-keys-addresses** | secp256k1, private keys, addresses |
| **04-accounts-balances** | Account queries, balance checks |
| **05-tx-nonces** | Transaction nonces, replay protection |
| **06-eip1559** | EIP-1559 dynamic fees, priority fees |

#### Smart Contracts (07-09)
| Project | Description |
|---------|-------------|
| **07-eth-call** | eth_call, read-only contract calls |
| **08-abigen** | ABI encoding/decoding, typed bindings |
| **09-events** | Event logs, decoding, filtering |

#### State & Storage (10-12)
| Project | Description |
|---------|-------------|
| **10-filters** | Log filters, WebSocket subscriptions |
| **11-storage** | Storage slots, state access |
| **12-proofs** | Merkle-Patricia tries, cryptographic proofs |

#### Advanced Queries (13-17)
| Project | Description |
|---------|-------------|
| **13-trace** | Transaction tracing, internal calls |
| **14-explorer** | Block explorer queries |
| **15-receipts** | Transaction receipts, gas used, logs |
| **16-concurrency** | Concurrent RPC calls, rate limiting |
| **17-indexer** | Event indexing, historical data |

#### Operations (18-25)
| Project | Description |
|---------|-------------|
| **18-reorgs** | Chain reorganizations, handling reorgs |
| **19-devnets** | Local devnets, Ganache, Hardhat |
| **20-node** | Node management, peer discovery |
| **21-sync** | Sync status, sync modes (full, snap, light) |
| **22-peers** | Peer management, network topology |
| **23-mempool** | Mempool monitoring, pending transactions |
| **24-monitor** | Node monitoring, health checks |
| **25-toolbox** | Utility functions, helper libraries |

---

## 🛠️ How to Use This Repository

### 1. **Implement in `exercise.go`**

All your code goes in `internal/<pkg>/exercise.go`. This is the **only** file that compiles during normal builds.

### 2. **Run Tests**

```bash
go test ./...
```

Tests verify your implementation against expected behavior.

### 3. **Use the Debug Harness (`cmd/dev`)**

```bash
go run ./cmd/dev
```

- Fixed, deterministic inputs
- Perfect for setting breakpoints and stepping through logic
- No command-line argument parsing needed

### 4. **Use the Application (`cmd/app`)**

```bash
go run ./cmd/app <arguments>
```

- Mimics real-world usage
- Accepts dynamic inputs via CLI
- Demonstrates how your library would be consumed

### 5. **Compare with Reference Implementations**

Reference files (`.reference.go`) are **excluded** from normal builds. View them for guidance:

```bash
# View reference solution
cat internal/<pkg>/solution.reference.go

# Optionally build with reference code (for exploration only)
go build -tags=reference ./...
```

---

## 🐛 Debugging with VS Code

Every project includes a `.vscode/launch.json` with these configurations:

| Configuration | Description |
|---------------|-------------|
| **Debug: cmd/app** | Debug the application entry point |
| **Debug: cmd/dev (Debug Harness)** | Debug with fixed inputs — perfect for stepping through code |
| **Test: Run All Tests** | Run all tests with verbose output |
| **Test: Current Test Function** | Debug specific test — select test name in editor first |
| **Test: View Reference Implementation** | Run tests with reference implementation |
| **Debug: Current File** | Debug currently open file |

### Debugging Workflow

1. Open project in VS Code
2. Set breakpoints in `internal/<pkg>/exercise.go`
3. Press F5 and select "Debug: cmd/dev (Debug Harness)"
4. Step through code, inspect variables, explore execution flow

---

## 🧪 Testing and Verification

### Run All Tests

```bash
# From repository root
go test ./...
```

### Run Tests for a Specific Project

```bash
cd minis/06-worker-pool-wordcount
go test ./...
```

### Run Benchmarks

```bash
go test -bench=. ./...
```

### Verify Reference Files Are Excluded

```bash
# This should compile only exercise.go (not solution.reference.go)
go build ./...

# Verify:
go list -f '{{.GoFiles}}' ./minis/01-hello-strings/internal/hellostrings
# Output should NOT include solution.reference.go
```

---

## 🎓 Learning Path Recommendation

### Absolute Beginners

1. Start with `minis/01-hello-strings` through `minis/05-cli-todo-files`
2. Master basic syntax, slices, maps, file I/O
3. Move to `minis/06-worker-pool-wordcount` to learn concurrency
4. Progress through `minis/11-17` for deep language understanding

### Intermediate Go Developers

1. Jump to `minis/18-27` for concurrency patterns
2. Explore `minis/28-30` for performance and profiling
3. Try `geth/01-06` for Ethereum RPC basics
4. Build `minis/31-38` for HTTP/network services

### Blockchain/Ethereum Developers

1. Start with `geth/01-stack` to understand RPC connectivity
2. Progress through `geth/02-12` for core Ethereum concepts
3. Dive into `geth/13-25` for advanced indexing, tracing, and node operations
4. Complement with `minis/39-45` for crypto/blockchain fundamentals

---

## 🔗 Additional Resources

- [Official Go Documentation](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)
- [go-ethereum (Geth) Documentation](https://geth.ethereum.org/docs)
- [Ethereum JSON-RPC Specification](https://ethereum.org/en/developers/docs/apis/json-rpc/)

---

## 📄 License

See `LICENSE` file for details.

---

**Happy Learning! 🚀**

Build real things. Understand deeply. Master Go.
