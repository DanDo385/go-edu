# Learning Arc

This course is designed as one continuous sequence. The three sections below are not separate silos; each stage depends on mental models from the previous stage.

## Opening Lecture

Before writing code in any lesson, keep these ideas active:

1. Stack vs heap allocation.
2. Values vs pointers.
3. Copying vs aliasing.
4. Escape analysis and allocation placement.
5. Goroutine scheduling and shared state.
6. Channels vs shared-memory synchronization.

At a CS level, this is about concrete memory behavior:

- A variable is a value in storage.
- A pointer is a value whose payload is an address.
- Passing arguments in Go is always pass-by-value.
- Some copied values contain addresses, which can still alias shared backing memory.

End-state understanding for the full repository:

1. You can reason about where values live and why.
2. You can explain when `*T`, `&x`, and `*p` imply indirection and when `*` is only multiplication.
3. You can implement classic DSAs without relying on hidden magic.
4. You can build concurrent services with cancellation, bounded work, and explicit ownership boundaries.
5. You can map those concepts into blockchain-client architecture.

## Stage 1: Data Structures and Algorithms

Focus: data layout, complexity, and deterministic transformations.

Suggested order:

1. `minis/01-hello-strings`
2. `minis/02-arrays-maps-basics`
3. `minis/03-csv-stats`
4. `minis/07-generic-lru-cache`
5. `minis/11-slices-internals-capacity-growth`
6. `minis/12-pointers-zero-values-nil-gotchas`
7. `minis/39-sha256-hasher`
8. `minis/40-merkle-tree-basics`
9. `minis/42-simple-block-struct-hashing`
10. `minis/43-proof-of-work-demo`
11. `minis/44-mempool-in-memory`
12. `minis/46-generics-map-reduce`
13. `minis/49-state-machine-pattern`

What this stage builds:

1. Arrays, slices, maps, and their memory tradeoffs.
2. Cache behavior and eviction policy (LRU).
3. Hash-based composition (Merkle and block hashing).
4. Algorithmic constraints such as throughput, latency, and bounded memory.

## Stage 2: Go as a Systems Language

Focus: processes, networking, synchronization, lifecycle, and performance.

Suggested order:

1. `minis/04-jsonl-log-filter`
2. `minis/05-cli-todo-files`
3. `minis/06-worker-pool-wordcount`
4. `minis/08-http-client-retries`
5. `minis/09-http-server-graceful`
6. `minis/10-grpc-telemetry-service`
7. `minis/13-interfaces-duck-typing`
8. `minis/14-methods-value-vs-pointer-receivers`
9. `minis/15-error-wrapping-sentinel-errors`
10. `minis/16-context-cancellation-timeouts`
11. `minis/17-file-streaming-bufio`
12. `minis/18-goroutines-1M-demo`
13. `minis/19-channels-basics`
14. `minis/20-select-fanin-fanout`
15. `minis/21-race-detection-demo`
16. `minis/22-worker-pool-with-backpressure`
17. `minis/23-bounded-channel-semaphore`
18. `minis/24-sync-mutex-vs-rwmutex`
19. `minis/25-atomic-counters-vs-mutex`
20. `minis/26-sync-once-singleton`
21. `minis/27-sync-pool-allocator`
22. `minis/28-pprof-cpu-mem-benchmarks`
23. `minis/29-escape-analysis-inlining`
24. `minis/30-build-tags-conditional-compilation`
25. `minis/31-static-file-server`
26. `minis/32-websocket-chatroom`
27. `minis/33-tcp-echo-server-client`
28. `minis/34-rate-limiter-token-bucket`
29. `minis/35-jwt-auth-middleware`
30. `minis/36-caching-reverse-proxy`
31. `minis/37-http-middleware-chain`
32. `minis/38-config-loader-env-yaml`
33. `minis/47-plugin-system-hot-reload`
34. `minis/48-reflection-introspection`
35. `minis/50-mini-service-all-features`

What this stage builds:

1. Concurrency patterns under load and cancellation.
2. Safe shared-state design and race avoidance.
3. Service lifecycle design for long-running systems.
4. Instrumentation and performance reasoning.

## Stage 3: Infrastructure and Blockchain Context

Focus: apply systems knowledge to blockchain-client style architecture.

Suggested order:

1. `geth/01-stack`
2. `geth/02-rpc-basics`
3. `geth/03-keys-addresses`
4. `geth/04-accounts-balances`
5. `geth/05-tx-nonces`
6. `geth/06-eip1559`
7. `geth/06-smart-contracts`
8. `geth/07-eth-call`
9. `geth/08-abigen`
10. `geth/09-events`
11. `geth/10-filters`
12. `geth/11-storage`
13. `geth/12-proofs`
14. `geth/13-trace`
15. `geth/14-explorer`
16. `geth/15-receipts`
17. `geth/16-concurrency`
18. `geth/17-indexer`
19. `geth/18-reorgs`
20. `geth/19-devnets`
21. `geth/20-node`
22. `geth/21-sync`
23. `geth/22-peers`
24. `geth/23-mempool`
25. `geth/24-monitor`
26. `geth/25-toolbox`

What this stage builds:

1. Real RPC data flow and decoding.
2. Event-driven indexing pipelines.
3. Reorg-aware state handling.
4. Peer/mempool/network operations.
5. Monitoring and operations-grade tooling.

## Required Lesson Contract

Every non-trivial lesson should include:

1. A student implementation file (`exercise.go`).
2. A clean reference implementation (`solution.reference.go`).
3. Tests that prove behavior and edge cases.
4. A README with conceptual framing and verification steps.

Authoring details live in `docs/LESSON_BLUEPRINT.md`.
