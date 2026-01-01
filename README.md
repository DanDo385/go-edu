## go-edu

Go learning repo built as a collection of **self-contained projects** under `minis/` (Go fundamentals) and `geth/` (Ethereum/go-ethereum-flavored exercises).

### Project layout (every project)

Each project follows the same **standard Go project layout**:

```text
<project>/
  cmd/
    app/
      main.go
    dev/
      main.go
  internal/
    <pkg>/
      exercise.go
      exercise_test.go
      solution.reference.go
      solution_no_err.reference.go
      (optional: extra *_test.go, benchmarks, types.go, etc.)
```

- **`internal/<pkg>/exercise.go`**: **student-facing** file. This is where you write your solution by filling in the `// TODO:` blocks (keep signatures the same).
- **`internal/<pkg>/exercise_test.go`**: tests that describe the required behavior and guide your implementation.
- **`internal/<pkg>/solution.reference.go`** / `solution_no_err.reference.go`: **reference implementations** with deeper explanations (and often debugging notes).
- **`cmd/app/main.go`**: "real" entry point (CLI/server/demo). Often takes arguments.
- **`cmd/dev/main.go`**: deterministic **debug harness** with fixed inputs. Designed to be breakpoint-friendly.

Repo-wide VS Code debugging is configured in **`.vscode/launch.json`** at the repository root (some projects may also ship their own `.vscode/` folder).

### Getting started

- **Clone**:

```bash
git clone <repo>
cd <repo>
```

- **Pick a project**:

```bash
cd minis/01-hello-strings
# or: cd geth/01-stack
```

- **Run tests**:

```bash
go test ./...
```

- **Run the program**:

```bash
go run ./cmd/app
```

- **Debug in VS Code**:
  - Open the repository root in VS Code.
  - Use the repo-wide launch configs in `.vscode/launch.json`:
    - **"Debug: cmd/dev (project)"** (recommended)
    - **"Debug: cmd/app (project)"**
    - **"Test: go test ./... (project)"**

### Solving an exercise (how to use the `// TODO:` blocks)

Open the exercise package (usually `internal/<pkg>/exercise.go`) and implement the code *inside* the `// TODO:` blocks.

- Keep **function signatures** unchanged (tests and `cmd/` programs depend on them).
- Use the `// TODO:` comments as your step-by-step plan.
- Run tests frequently:

```bash
go test ./...
```

If you get stuck, read `solution.reference.go` / `solution_no_err.reference.go` for a full implementation and deeper explanations.

### Commentary guidelines (where documentation should live)

- **`exercise.go`**: keep commentary minimal and actionable (short `// TODO:` steps + brief debugging reminders).
- **`solution.reference.go`** and this `README.md`: put longer explanations, walkthroughs, and “why this works” material here.
- **Breakpoints/debugging**: keep and use `// BREAKPOINT:` markers in `cmd/dev` (and in exercises where they help).

### Resetting exercises (`make todo`)

This repo intentionally supports “start over” resets that **erase your in-progress code** and regenerate `exercise.go` files back to their starter TODO state from the reference solutions.

- From repo root: `make todo all` resets **both** `minis/` and `geth/`
- From `geth/`: `make todo all` resets **only** `geth/`
- From `minis/`: `make todo all` resets **only** `minis/`
- From a project folder (e.g. `geth/01-stack/`): `make todo all` resets **only that project**
- From anywhere: `make todo <path>` resets the project(s) under `<path>` (root-relative paths work)

---

## Testing and Debugging Guide

### Overview

Each project provides two ways to test and debug your implementation:

1. **cmd/app/main.go** - Application entry point with CLI arguments
2. **cmd/dev/main.go** - Debug harness with fixed inputs

### Using cmd/dev/main.go (Recommended for Learning)

The `cmd/dev/main.go` file is a debug harness designed for stepping through code with breakpoints.

#### Why Use cmd/dev/main.go?

- **Fixed inputs**: No need to remember command-line arguments
- **Deterministic**: Same inputs every time, making debugging predictable
- **Focused**: Contains only the essential code to test your implementation
- **Breakpoint-friendly**: Includes "// BREAKPOINT:" comments at key locations

#### How to Use:

1. **Open** `cmd/dev/main.go` in VS Code
2. **Set breakpoints** at "// BREAKPOINT:" comments or anywhere in your code
3. **Press F5** and select "Debug: cmd/dev (project)" (then enter the project path when prompted)
4. **Step through** using F10 (Step Over) and F11 (Step Into)
5. **Watch variables** in the Variables panel

#### Example Workflow:

```bash
# 1. Navigate to project
cd minis/01-hello-strings

# 2. Open cmd/dev/main.go in VS Code
# 3. Set breakpoint in internal/hellostrings/exercise.go
# 4. Press F5, select "Debug: cmd/dev (project)" and enter: minis/01-hello-strings
# 5. Debugger stops at breakpoint - step through implementation
```

### Using cmd/app/main.go (CLI Arguments)

The `cmd/app/main.go` file is the application entry point that accepts command-line arguments.

#### Project-Specific CLI Arguments

Each project has different CLI arguments based on its purpose:

##### minis/ Projects:

- **01-hello-strings**: `[input_string]`
  - Example: `go run ./cmd/app/main.go "hello world"`
- **06-worker-pool-wordcount**: `[url1] [url2] ... [urlN]`
  - Example: `go run ./cmd/app/main.go https://example.com https://example.org`
- **08-http-client-retries**: `[url] [max-retries]`
  - Example: `go run ./cmd/app/main.go https://api.example.com 3`

##### geth/ Projects:

- **01-stack**: `<RPC_URL> [block_number]`
  - Example: `go run ./cmd/app/main.go https://eth.llamarpc.com`
  - Example: `go run ./cmd/app/main.go https://eth.llamarpc.com 12345`
- **05-tx-nonces**: `<RPC_URL> <address>`
  - Example: `go run ./cmd/app/main.go https://eth.llamarpc.com 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045`
- **06-smart-contracts**: `[RPC_URL] [CONTRACT_ADDRESS]` (Optional - primarily console-based)
  - Note: This module is primarily console-based. See README.md for Geth console tutorial.
- **07-eth-call**: `<RPC_URL> <contract_address>`
  - Example: `go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48`
- **09-events**: `<RPC_URL> <contract_address>`
  - Example: `go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48`

#### Debugging cmd/app/main.go:

1. Press F5 and select **"Debug: cmd/app (project)"** (then enter the project path when prompted)
2. If your `cmd/app` expects CLI args, either:
   - Run it from a terminal (recommended for ad-hoc args): `go run ./cmd/app -- <args...>`
   - Or temporarily edit `.vscode/launch.json` to add an `"args": [...]` array for that configuration
3. Set breakpoints and step through

### Testing with go test

#### Running Tests:

```bash
# Run all tests in current project
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test function
go test -v -run TestFunctionName ./...

# Run with reference implementation
go test -tags=reference -v ./...
```

#### Debugging Tests:

1. Set breakpoint in test function or exercise code
2. Press F5 and select "Test: go test ./... (project)" (then enter the project path when prompted)
3. Step through test execution

### VS Code Debug Configurations

The repository includes `.vscode/launch.json` with these configurations:

| Configuration | Description |
|---------------|-------------|
| **Debug: cmd/app (project)** | Debug application with CLI arguments |
| **Debug: cmd/dev (project)** | Debug harness with fixed inputs (recommended) |
| **Test: go test ./... (project)** | Run all tests in the project with debugger |
| **Test: go test -tags=reference ./... (project)** | Run tests against reference implementations |

### Tips for Effective Debugging

1. **Start with cmd/dev/main.go** - It's designed for learning
2. **Use breakpoints liberally** - Set them at function entry points
3. **Watch the Variables panel** - See how data transforms
4. **Use Call Stack panel** - Understand function call hierarchy
5. **Step Into (F11)** - Enter function implementations
6. **Step Over (F10)** - Execute line by line
7. **Step Out (Shift+F11)** - Return to caller

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

---

### Build tags used in this repo

- **Exercise implementation** (`internal/<pkg>/exercise.go`):

```go
//go:build !solution && !reference
```

- **Reference implementations** (`*.reference.go`):

```go
//go:build reference
```

To run tests using reference implementations:

```bash
go test -tags=reference ./...
```

---

### Project index

#### minis/

- [minis/01-hello-strings](./minis/01-hello-strings/)
- [minis/02-arrays-maps-basics](./minis/02-arrays-maps-basics/)
- [minis/03-csv-stats](./minis/03-csv-stats/)
- [minis/04-jsonl-log-filter](./minis/04-jsonl-log-filter/)
- [minis/05-cli-todo-files](./minis/05-cli-todo-files/)
- [minis/06-worker-pool-wordcount](./minis/06-worker-pool-wordcount/)
- [minis/07-generic-lru-cache](./minis/07-generic-lru-cache/)
- [minis/08-http-client-retries](./minis/08-http-client-retries/)
- [minis/09-http-server-graceful](./minis/09-http-server-graceful/)
- [minis/10-grpc-telemetry-service](./minis/10-grpc-telemetry-service/)
- [minis/11-slices-internals-capacity-growth](./minis/11-slices-internals-capacity-growth/)
- [minis/12-pointers-zero-values-nil-gotchas](./minis/12-pointers-zero-values-nil-gotchas/)
- [minis/13-interfaces-duck-typing](./minis/13-interfaces-duck-typing/)
- [minis/14-methods-value-vs-pointer-receivers](./minis/14-methods-value-vs-pointer-receivers/)
- [minis/15-error-wrapping-sentinel-errors](./minis/15-error-wrapping-sentinel-errors/)
- [minis/16-context-cancellation-timeouts](./minis/16-context-cancellation-timeouts/)
- [minis/17-file-streaming-bufio](./minis/17-file-streaming-bufio/)
- [minis/18-goroutines-1M-demo](./minis/18-goroutines-1M-demo/)
- [minis/19-channels-basics](./minis/19-channels-basics/)
- [minis/20-select-fanin-fanout](./minis/20-select-fanin-fanout/)
- [minis/21-race-detection-demo](./minis/21-race-detection-demo/)
- [minis/22-worker-pool-with-backpressure](./minis/22-worker-pool-with-backpressure/)
- [minis/23-bounded-channel-semaphore](./minis/23-bounded-channel-semaphore/)
- [minis/24-sync-mutex-vs-rwmutex](./minis/24-sync-mutex-vs-rwmutex/)
- [minis/25-atomic-counters-vs-mutex](./minis/25-atomic-counters-vs-mutex/)
- [minis/26-sync-once-singleton](./minis/26-sync-once-singleton/)
- [minis/27-sync-pool-allocator](./minis/27-sync-pool-allocator/)
- [minis/28-pprof-cpu-mem-benchmarks](./minis/28-pprof-cpu-mem-benchmarks/)
- [minis/29-escape-analysis-inlining](./minis/29-escape-analysis-inlining/)
- [minis/30-build-tags-conditional-compilation](./minis/30-build-tags-conditional-compilation/)
- [minis/31-static-file-server](./minis/31-static-file-server/)
- [minis/32-websocket-chatroom](./minis/32-websocket-chatroom/)
- [minis/33-tcp-echo-server-client](./minis/33-tcp-echo-server-client/)
- [minis/34-rate-limiter-token-bucket](./minis/34-rate-limiter-token-bucket/)
- [minis/35-jwt-auth-middleware](./minis/35-jwt-auth-middleware/)
- [minis/36-caching-reverse-proxy](./minis/36-caching-reverse-proxy/)
- [minis/37-http-middleware-chain](./minis/37-http-middleware-chain/)
- [minis/38-config-loader-env-yaml](./minis/38-config-loader-env-yaml/)
- [minis/39-sha256-hasher](./minis/39-sha256-hasher/)
- [minis/40-merkle-tree-basics](./minis/40-merkle-tree-basics/)
- [minis/41-signed-transactions-ed25519](./minis/41-signed-transactions-ed25519/)
- [minis/42-simple-block-struct-hashing](./minis/42-simple-block-struct-hashing/)
- [minis/43-proof-of-work-demo](./minis/43-proof-of-work-demo/)
- [minis/44-mempool-in-memory](./minis/44-mempool-in-memory/)
- [minis/45-p2p-gossip-mock-network](./minis/45-p2p-gossip-mock-network/)
- [minis/46-generics-map-reduce](./minis/46-generics-map-reduce/)
- [minis/47-plugin-system-hot-reload](./minis/47-plugin-system-hot-reload/)
- [minis/48-reflection-introspection](./minis/48-reflection-introspection/)
- [minis/49-state-machine-pattern](./minis/49-state-machine-pattern/)
- [minis/50-mini-service-all-features](./minis/50-mini-service-all-features/)

#### geth/

- [geth/01-stack](./geth/01-stack/)
- [geth/02-rpc-basics](./geth/02-rpc-basics/)
- [geth/03-keys-addresses](./geth/03-keys-addresses/)
- [geth/04-accounts-balances](./geth/04-accounts-balances/)
- [geth/05-tx-nonces](./geth/05-tx-nonces/)
- [geth/06-smart-contracts](./geth/06-smart-contracts/) *(Console Tutorial - Smart Contract Fundamentals)*
- [geth/06-eip1559](./geth/06-eip1559/)
- [geth/07-eth-call](./geth/07-eth-call/)
- [geth/08-abigen](./geth/08-abigen/)
- [geth/09-events](./geth/09-events/)
- [geth/10-filters](./geth/10-filters/)
- [geth/11-storage](./geth/11-storage/)
- [geth/12-proofs](./geth/12-proofs/)
- [geth/13-trace](./geth/13-trace/)
- [geth/14-explorer](./geth/14-explorer/)
- [geth/15-receipts](./geth/15-receipts/)
- [geth/16-concurrency](./geth/16-concurrency/)
- [geth/17-indexer](./geth/17-indexer/)
- [geth/18-reorgs](./geth/18-reorgs/)
- [geth/19-devnets](./geth/19-devnets/)
- [geth/20-node](./geth/20-node/)
- [geth/21-sync](./geth/21-sync/)
- [geth/22-peers](./geth/22-peers/)
- [geth/23-mempool](./geth/23-mempool/)
- [geth/24-monitor](./geth/24-monitor/)
- [geth/25-toolbox](./geth/25-toolbox/)

---

### Learning Paths

#### Go Fundamentals (minis/)

Start with `minis/01-hello-strings` and progress sequentially. Each project builds on previous concepts.

**Beginner**: 01-10 (basics, I/O, HTTP)
**Intermediate**: 11-30 (concurrency, performance, internals)
**Advanced**: 31-50 (networking, crypto, production patterns)

#### Ethereum Development (geth/)

Start with `geth/01-stack` and progress sequentially. Prerequisites are listed in each project's README.

**Foundational**: 01-06 (connectivity, accounts, transactions, console basics)
**Contracts**: 07-10 (manual calls, abigen, events, filters)
**Advanced**: 11-25 (storage, proofs, tracing, indexing, networking)

**Note**: For smart contract interaction, complete `geth/06-smart-contracts` (console tutorial) before `geth/07-eth-call` (Go implementation). The console experience provides essential conceptual foundation.
