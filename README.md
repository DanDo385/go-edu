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

- [geth/01-stack](./geth/01-stack/)
- [geth/02-rpc-basics](./geth/02-rpc-basics/)
- [geth/03-keys-addresses](./geth/03-keys-addresses/)
- [geth/04-accounts-balances](./geth/04-accounts-balances/)
- [geth/05-tx-nonces](./geth/05-tx-nonces/)
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
