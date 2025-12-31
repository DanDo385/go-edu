# Minis — Go Fundamentals (50 Projects)

Learn Go from first principles through hands-on projects covering core language features, concurrency patterns, networking, and blockchain basics.

## 📚 Project Overview

Each project in this directory is **self-contained** and follows the same structure:

```
project-name/
├── README.md                    # (removed for cleaner structure)
├── cmd/
│   ├── app/main.go             # Realistic CLI application
│   └── dev/main.go             # Debug harness (fixed inputs)
└── internal/
    └── <pkg>/
        ├── exercise.go         # YOUR implementation (stubs with TODO comments)
        ├── exercise_test.go    # Tests to verify your work
        ├── solution.reference.go        # Reference implementation (excluded from builds)
        └── solution_no_err.reference.go # Alternative reference
```

## 🎯 How to Use

### 1. Navigate to a Project

```bash
cd minis/01-hello-strings
```

### 2. Implement in `exercise.go`

Open `internal/<pkg>/exercise.go` and implement the functions marked with `TODO` comments.

### 3. Run Tests

```bash
go test ./...
```

### 4. Debug with `cmd/dev`

```bash
go run ./cmd/dev
```

Fixed, deterministic inputs — perfect for stepping through with a debugger.

### 5. Run the Application

```bash
go run ./cmd/app <arguments>
```

Realistic CLI usage demonstrating how your code would be consumed.

### 6. Compare with Reference

```bash
cat internal/<pkg>/solution.reference.go
```

Reference implementations are **excluded** from normal builds (build tag: `reference`).

---

## 📋 Complete Project List

### Fundamentals (01-05)
- **01-hello-strings** — String manipulation, UTF-8, runes
- **02-arrays-maps-basics** — Arrays, slices, maps
- **03-csv-stats** — CSV parsing, file I/O
- **04-jsonl-log-filter** — JSONL parsing, filtering
- **05-cli-todo-files** — File operations, CLI

### HTTP & Networking (06-10)
- **06-worker-pool-wordcount** — Concurrency, worker pools
- **07-generic-lru-cache** — Generics, LRU caching
- **08-http-client-retries** — HTTP client, retries
- **09-http-server-graceful** — HTTP server, graceful shutdown
- **10-grpc-telemetry-service** — gRPC, Protocol Buffers

### Deep Dives (11-17)
- **11-slices-internals-capacity-growth** — Slice internals
- **12-pointers-zero-values-nil-gotchas** — Pointers, nil
- **13-interfaces-duck-typing** — Interfaces, duck typing
- **14-methods-value-vs-pointer-receivers** — Method receivers
- **15-error-wrapping-sentinel-errors** — Error handling
- **16-context-cancellation-timeouts** — Context, cancellation
- **17-file-streaming-bufio** — Streaming, bufio

### Concurrency Patterns (18-27)
- **18-goroutines-1M-demo** — Goroutines at scale
- **19-channels-basics** — Channels fundamentals
- **20-select-fanin-fanout** — Select, fan-in, fan-out
- **21-race-detection-demo** — Race detection
- **22-worker-pool-with-backpressure** — Backpressure
- **23-bounded-channel-semaphore** — Semaphores
- **24-sync-mutex-vs-rwmutex** — Mutex vs RWMutex
- **25-atomic-counters-vs-mutex** — Atomics vs Mutex
- **26-sync-once-singleton** — sync.Once, singletons
- **27-sync-pool-allocator** — sync.Pool, object pooling

### Performance & Profiling (28-30)
- **28-pprof-cpu-mem-benchmarks** — pprof, benchmarking
- **29-escape-analysis-inlining** — Escape analysis
- **30-build-tags-conditional-compilation** — Build tags

### Advanced HTTP (31-38)
- **31-static-file-server** — Static files
- **32-websocket-chatroom** — WebSockets
- **33-tcp-echo-server-client** — TCP networking
- **34-rate-limiter-token-bucket** — Rate limiting
- **35-jwt-auth-middleware** — JWT authentication
- **36-caching-reverse-proxy** — Reverse proxy, caching
- **37-http-middleware-chain** — Middleware patterns
- **38-config-loader-env-yaml** — Configuration

### Cryptography & Blockchain (39-45)
- **39-sha256-hasher** — SHA256 hashing
- **40-merkle-tree-basics** — Merkle trees
- **41-signed-transactions-ed25519** — Digital signatures
- **42-simple-block-struct-hashing** — Block structures
- **43-proof-of-work-demo** — Proof of Work
- **44-mempool-in-memory** — Mempool implementation
- **45-p2p-gossip-mock-network** — P2P networking

### Advanced Topics (46-50)
- **46-generics-map-reduce** — Generics, map/reduce
- **47-plugin-system-hot-reload** — Plugin systems
- **48-reflection-introspection** — Reflection
- **49-state-machine-pattern** — State machines
- **50-mini-service-all-features** — Complete microservice

---

## 🚀 Quick Start

```bash
# Start with project 01
cd minis/01-hello-strings

# Implement functions in internal/hellostrings/exercise.go
# Run tests
go test ./...

# Debug with fixed inputs
go run ./cmd/dev

# Run application
go run ./cmd/app "hello world"
```

## 🔑 Key Concepts

- **Self-contained projects** — Each project is fully independent
- **One buildable implementation** — Only `exercise.go` compiles
- **Reference files are inert** — Build tag `reference` excludes them
- **cmd/app vs cmd/dev** — Production vs debugging harnesses

## 📖 Learning Path

**Beginners:** Start with 01-05, then 06, then 11-17
**Intermediate:** Jump to 18-27 (concurrency), then 31-38 (HTTP)
**Blockchain:** Do 01-05 basics, then 39-45 (crypto/blockchain)

---

See the [root README](../README.md) for complete documentation.
