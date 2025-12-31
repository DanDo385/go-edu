# go-edu

**Learn Go by building real systems.**

This repository is a comprehensive educational resource for mastering Go through hands-on projects that mirror real-world production patterns and Ethereum/blockchain tooling.

## 🎯 Philosophy

**go-edu** teaches Go through two complementary tracks:

1. **`minis/`** — Fundamental Go patterns, concurrency primitives, HTTP servers, crypto basics
2. **`geth/`** — Ethereum-focused projects inspired by go-ethereum architecture and patterns

Every project is self-contained, production-ready in structure, and designed for deep understanding through debugging and experimentation.

---

## 📂 Repository Structure

```
go-edu/
├── go.mod                    # Single module at root
├── README.md                 # This file
├── minis/                    # 50 mini-projects covering core Go
│   ├── 01-hello-strings/
│   │   ├── README.md
│   │   ├── cmd/
│   │   │   ├── app/          # Realistic CLI application
│   │   │   │   └── main.go
│   │   │   └── dev/          # Debug harness (fixed inputs)
│   │   │       └── main.go
│   │   └── internal/
│   │       └── <pkg>/        # Self-contained package
│   │           ├── exercise.go
│   │           ├── exercise_test.go
│   │           ├── solution.reference.go         # Reference (excluded from builds)
│   │           └── solution_no_err.reference.go  # Alternative reference
│   └── ...
└── geth/                     # 25 Ethereum/Geth-focused projects
    ├── 01-stack/
    │   ├── README.md
    │   ├── cmd/
    │   │   ├── app/          # Realistic RPC client app
    │   │   │   └── main.go
    │   │   └── dev/          # Debug harness (fixed RPC endpoint)
    │   │       └── main.go
    │   └── internal/
    │       └── <pkg>/
    │           ├── exercise.go
    │           ├── exercise_test.go
    │           ├── types.go  # Geth-style type definitions
    │           └── solution.reference.go
    └── ...
```

---

## 🏗️ Project Structure Explained

Every project follows the **self-contained layout**:

### Directory Layout

| Directory/File | Purpose |
|----------------|---------|
| `cmd/app/main.go` | **Realistic application** — accepts CLI arguments, simulates production usage |
| `cmd/dev/main.go` | **Debug harness** — fixed, deterministic inputs for stepping through with a debugger |
| `internal/<pkg>/exercise.go` | **The only buildable implementation** — all production logic lives here |
| `internal/<pkg>/*_test.go` | Tests and benchmarks |
| `internal/<pkg>/*.reference.go` | Reference solutions (excluded from normal builds via `//go:build reference`) |
| `README.md` | Project-specific documentation and learning objectives |

### Build Tags and Implementation Model

#### Core Invariant: One Buildable Implementation

- `exercise.go` is the **only** non-test, non-reference file that participates in builds
- All production logic must be consolidated into `exercise.go`
- This enforces clarity: there is one source of truth for implementation

#### Reference Files Are Inert

Reference implementations exist **only** for learning:

```go
//go:build reference
// +build reference

package mypackage

// This file is NEVER compiled during normal builds
```

To view reference implementations:

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

## 📚 Learning Tracks

### `minis/` — Core Go Fundamentals (50 Projects)

| Range | Topics |
|-------|--------|
| **01-05** | Strings, arrays, maps, CSV, JSONL, file I/O |
| **06-10** | Worker pools, generics, HTTP client/server, gRPC |
| **11-17** | Slices internals, pointers, interfaces, methods, errors, context, bufio |
| **18-27** | Goroutines, channels, select, race detection, sync primitives (mutex, RWMutex, Once, Pool, atomic) |
| **28-30** | Profiling (pprof), escape analysis, build tags |
| **31-38** | HTTP servers, WebSockets, TCP, rate limiting, JWT, proxies, middleware, config |
| **39-45** | Cryptography (SHA256, Merkle trees, ed25519), blockchain basics (PoW, mempool, p2p gossip) |
| **46-50** | Generics, plugins, reflection, state machines, full-featured mini-service |

### `geth/` — Ethereum & Go-Ethereum Patterns (25 Projects)

| Range | Topics |
|-------|--------|
| **01-06** | RPC stack, chain ID, keys/addresses, accounts, transaction nonces, EIP-1559 |
| **07-12** | eth_call, abigen, events, filters, storage proofs, Merkle-Patricia tries |
| **13-17** | Tracing, block explorer, receipts, concurrency patterns, indexers |
| **18-25** | Reorgs, devnets, node setup, sync modes, peer management, mempool monitoring, toolbox utilities |

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

**Important:** Reference files are for learning only. They do **not** participate in tests or builds unless you explicitly use `-tags=reference`.

---

## 🧪 Testing and Verification

### Run All Tests

```bash
# From repository root
go test ./...
```

### Run Benchmarks

```bash
go test -bench=. ./...
```

### Run Tests for a Specific Project

```bash
cd minis/06-worker-pool-wordcount
go test ./...
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

## 🐛 Debugging

Every project includes:

- **`cmd/dev/main.go`** — Fixed inputs for reproducible debugging
- Inline comments marking good breakpoint locations
- VS Code launch configurations (`.vscode/launch.json` where applicable)

### Debugging Workflow

1. Open project in your editor (e.g., VS Code)
2. Set breakpoints in `internal/<pkg>/exercise.go`
3. Run `cmd/dev/main.go` with debugger (F5 in VS Code)
4. Step through code, inspect variables, explore execution flow

See individual project `README.md` files for project-specific debugging tips.

---

## 📖 Project-Specific READMEs

Every project has a `README.md` with:

- **Description** — What the project teaches
- **Concepts Covered** — Go features, CS principles, production patterns
- **How to Run** — Commands for tests, benchmarks, dev harness, app
- **Learning Objectives** — What you'll master by completing the project
- **Debugging Tips** — Suggested breakpoints and exploration strategies

Navigate to any project and read its `README.md`:

```bash
cd minis/07-generic-lru-cache
cat README.md
```

---

## 🎓 Educational Philosophy

### Why This Structure?

1. **Self-contained projects** — Each project is an isolated, complete system
2. **cmd/app vs cmd/dev** — Separates "production usage" from "debugging harness"
3. **One buildable implementation** — No confusion about which code is "active"
4. **Reference files are inert** — Learn by comparing, not by accidentally compiling multiple implementations
5. **Real-world patterns** — Every project mirrors production Go codebases (especially Geth and Ethereum tooling)

### Who Is This For?

- **New Go developers** transitioning from other languages
- **Intermediate Go developers** wanting to master concurrency, interfaces, and production patterns
- **Blockchain developers** learning go-ethereum internals and Ethereum RPC patterns
- **Anyone** seeking hands-on, project-based Go education

---

## 🔗 Additional Resources

- [Official Go Documentation](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)
- [go-ethereum (Geth) Documentation](https://geth.ethereum.org/docs)
- [Ethereum JSON-RPC Specification](https://ethereum.org/en/developers/docs/apis/json-rpc/)

---

## 🤝 Contributing

This is an educational repository. If you find issues or have suggestions:

1. Open an issue describing the problem or enhancement
2. For typos or small fixes, submit a pull request
3. For new projects or major changes, discuss in an issue first

---

## 📄 License

See `LICENSE` file for details.

---

## 🌟 Learning Path Recommendation

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

**Happy Learning! 🚀**

Build real things. Understand deeply. Master Go.
