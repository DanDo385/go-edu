# geth-16-concurrency

**Goal:** fetch multiple resources concurrently with a worker pool while respecting RPC limits.

## Big Picture

RPCs can be slow; fan-out speeds things up, but you need to avoid rate limits and handle cancellation. Worker pools with contexts let you balance throughput and safety. You’ll reuse this pattern in indexers (module 17) and monitors (module 24).

## Learning Objectives
- Build a simple worker pool (jobs/results channels + WaitGroup).
- Use contexts to cancel/timeout concurrent work.
- Understand rate limiting/backoff considerations and how to bound goroutines.

## Prerequisites
- Modules 01–10; Go concurrency basics (goroutines, channels, context).

## Real-World Analogy
- Multiple clerks fetching ledger pages in parallel; close the office when time is up (context cancel).

## Steps
1. Parse endpoints + worker count.
2. Spin up worker goroutines that probe endpoints via a `Prober`.
3. Feed jobs, close channel, wait for completion.
4. Aggregate successes/failures with a mutex.

## Fun Facts & Comparisons
- Too much fan-out can trigger provider rate limits; add backoff/token bucket in production.
- ethers.js/JS often uses Promise.all; Go idiom is worker pools with channels.
- Go’s goroutines are ~2KB stacks; still, unbounded creation can exhaust resources.

## Related Solidity-edu Modules
- 17-indexer — fan out block fetches.
- 24-monitor — continuous health checks.
- 06-eip1559 — reuse contexts to time-bound RPCs.

## Files
- `exercise/exercise.go`: TODOs to build the worker pool.
- `exercise/solution.go`: reference implementation with comments.
- `exercise/exercise_test.go`: add your own stress tests and race checks.
