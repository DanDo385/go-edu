## go-edu

Go learning repo built as a collection of **self-contained projects** under `minis/` (Go fundamentals) and `geth/` (Ethereum/go-ethereum-flavored exercises).

### Project layout (every project)

Each project follows the same shape:

```text
<project>/
  .vscode/
    launch.json
    settings.json
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

- **`internal/<pkg>`**: the exercise package and tests.
- **`cmd/app`**: “real” entry point (CLI/server/demo).
- **`cmd/dev`**: deterministic “debug harness” entry point.
- **`.vscode/`**: per-project debug/run/test configurations.

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
  - Open the project folder (e.g. `minis/01-hello-strings/`).
  - Use:
    - **“Debug: cmd/dev”** (debug harness)
    - **“Run: cmd/app”** (no-debug run)
    - **“Test: internal package”** (debug tests)

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

- [geth/01-stack](./geth/01-stack/) - RPC connectivity and chain information
- [geth/02-rpc-basics](./geth/02-rpc-basics/) - RPC methods and retry logic
- [geth/03-keys-addresses](./geth/03-keys-addresses/) - Key generation and addresses
- [geth/04-accounts-balances](./geth/04-accounts-balances/) - Account state queries
- [geth/05-tx-nonces](./geth/05-tx-nonces/) - Transaction nonces and sequencing
- [geth/06-smart-contracts](./geth/06-smart-contracts/) - **Smart contract interaction fundamentals (Geth console tutorial)**
- [geth/07-eth-call](./geth/07-eth-call/) - Contract calls in Go (eth_call)
- [geth/08-abigen](./geth/08-abigen/) - Type-safe contract bindings
- [geth/09-events](./geth/09-events/) - Contract events and logs
- [geth/10-filters](./geth/10-filters/) - Log filtering techniques
- [geth/11-storage](./geth/11-storage/) - Direct storage slot reading
- [geth/12-proofs](./geth/12-proofs/) - Merkle proofs and verification
- [geth/13-trace](./geth/13-trace/) - Transaction tracing and debugging
- [geth/14-explorer](./geth/14-explorer/) - Block explorer implementation
- [geth/15-receipts](./geth/15-receipts/) - Transaction receipt handling
- [geth/16-concurrency](./geth/16-concurrency/) - Concurrent RPC patterns
- [geth/17-indexer](./geth/17-indexer/) - Event indexer implementation
- [geth/18-reorgs](./geth/18-reorgs/) - Chain reorganization handling
- [geth/19-devnets](./geth/19-devnets/) - Local development networks
- [geth/20-node](./geth/20-node/) - Geth node management
- [geth/21-sync](./geth/21-sync/) - Chain synchronization modes
- [geth/22-peers](./geth/22-peers/) - P2P peer management
- [geth/23-mempool](./geth/23-mempool/) - Mempool monitoring
- [geth/24-monitor](./geth/24-monitor/) - Network and node monitoring
- [geth/25-toolbox](./geth/25-toolbox/) - Developer utilities
- [geth/26-eip1559](./geth/26-eip1559/) - EIP-1559 transaction fee mechanics

---

## Testing and Debugging Guide

### Overview

Each project provides two ways to run and debug your implementation:

1. **cmd/app/main.go** - Application entry point with CLI arguments
2. **cmd/dev/main.go** - Debug harness with fixed inputs (recommended for learning)

### Using cmd/dev/main.go (Recommended for Learning)

The `cmd/dev/main.go` file is a debug harness designed for stepping through code with breakpoints.

#### Why Use cmd/dev/main.go?

- **Fixed inputs**: No need to remember command-line arguments
- **Deterministic**: Same inputs every time, making debugging predictable
- **Focused**: Contains only the essential code to test your implementation
- **Breakpoint-friendly**: Includes "// BREAKPOINT:" comments at key locations

#### How to Use

1. **Open** `cmd/dev/main.go` in VS Code
2. **Set breakpoints** at "// BREAKPOINT:" comments or anywhere in your code
3. **Press F5** and select "Debug: cmd/dev (Debug Harness)"
4. **Step through** using:
   - F10 (Step Over) - Execute current line
   - F11 (Step Into) - Enter function calls
   - Shift+F11 (Step Out) - Return to caller
5. **Watch variables** in the Variables panel

#### Example Workflow

```bash
# 1. Navigate to project
cd minis/01-hello-strings

# 2. Open cmd/dev/main.go in VS Code

# 3. Set breakpoint in internal/hellostrings/exercise.go

# 4. Press F5, select "Debug: cmd/dev"

# 5. Debugger stops at breakpoint - step through implementation
```

### Using cmd/app/main.go (CLI Arguments)

The `cmd/app/main.go` file is the application entry point that accepts command-line arguments.

#### Project-Specific CLI Arguments

Each project has different CLI arguments based on its purpose:

**minis/ Projects:**

- **01-hello-strings**: `[input_string] [function]`
  ```bash
  go run ./cmd/app/main.go "hello world" titlecase
  ```

- **06-worker-pool-wordcount**: `[url1] [url2] ... [urlN]`
  ```bash
  go run ./cmd/app/main.go https://example.com https://example.org
  ```

- **08-http-client-retries**: `[url] [max-retries]`
  ```bash
  go run ./cmd/app/main.go https://api.example.com 3
  ```

**geth/ Projects:**

- **01-stack**: `<RPC_URL> [block_number]`
  ```bash
  go run ./cmd/app/main.go https://eth.llamarpc.com
  go run ./cmd/app/main.go https://eth.llamarpc.com 12345
  ```

- **05-tx-nonces**: `<RPC_URL> <address>`
  ```bash
  go run ./cmd/app/main.go https://eth.llamarpc.com 0x742d35Cc6634C0532925a3b844Bc454e4438f44e
  ```

- **06-smart-contracts**: Console-based tutorial (see README.md)
  ```bash
  # This module uses Geth console, not cmd/app
  # See geth/06-smart-contracts/README.md for tutorial
  ```

- **07-eth-call**: `<RPC_URL> <contract_address>`
  ```bash
  go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
  ```

- **09-events**: `<RPC_URL> <contract_address> <event_signature>`
  ```bash
  go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991... Transfer(address,address,uint256)
  ```

#### Debugging cmd/app/main.go

1. **Open** `.vscode/launch.json`
2. **Find** "Debug: cmd/app" configuration
3. **Edit** the `args` array to include your CLI arguments:
   ```json
   "args": ["https://eth.llamarpc.com", "12345"]
   ```
4. **Press F5** and select "Debug: cmd/app"
5. Set breakpoints and step through

### Testing with go test

#### Running Tests

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

#### Debugging Tests

1. Set breakpoint in test function or exercise code
2. Press F5 and select "Test: Run All Tests" or "Test: Current Test Function"
3. Step through test execution

### VS Code Debug Configurations

Each project includes `.vscode/launch.json` with these configurations:

- **Debug: cmd/app** - Debug application with CLI arguments
- **Debug: cmd/dev (Debug Harness)** - Debug harness with fixed inputs (recommended)
- **Test: Run All Tests** - Run all tests with debugger
- **Test: Current Test Function** - Debug specific test (edit test name in launch.json)
- **Test: View Reference Implementation** - Run tests with `solution.reference.go`
- **Debug: Current File** - Debug currently open file

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

```bash
# Run go mod tidy in project directory
go mod tidy
# Ensure you're in the correct directory
```

#### "Build constraints exclude all Go files"

- Check build tags in `exercise.go` (should be `//go:build !solution && !reference`)
- Ensure you're not building with conflicting tags

#### "RPC connection failed" (geth projects)

- Check RPC URL is correct and accessible
- Try a different public RPC endpoint: https://chainlist.org/
- Ensure network connectivity
- Verify firewall settings

#### "Geth console not working" (geth/06-smart-contracts)

- Ensure Geth is installed: `geth version`
- Check Geth is running: `geth attach` should connect
- Verify RPC endpoint is accessible
- Check firewall settings

### Special Note: geth/06-smart-contracts

**Module 06 is console-tutorial based.** The main learning happens in the Geth JavaScript console, not in Go code. This teaches you:

- How contracts work at the RPC level
- The difference between `eth_call` (read-only) and `eth_sendTransaction` (state-changing)
- How to decode events and logs manually
- What happens under the hood before you use Go abstractions

After completing the console tutorial, modules 07-09 teach you how to do the same things in Go:

- **07-eth-call**: Making contract calls from Go
- **08-abigen**: Type-safe Go bindings
- **09-events**: Event subscriptions and decoding

To use module 06:
1. Read `geth/06-smart-contracts/README.md`
2. Start Geth: `geth --dev --http --http.api eth,net,web3,personal`
3. Attach: `geth attach`
4. Follow the step-by-step console tutorial

---
