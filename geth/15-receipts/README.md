# geth-15-receipts

**Goal:** Fetch transaction receipts and classify success/failure, logs, and gas usage.

## Big Picture

Receipts record the outcome of a transaction: status, cumulative gas, and emitted logs. They live alongside blocks but outside the state trie. Decoding receipts is key for dApps, indexers, and monitoring.

## Learning Objectives
- Fetch receipts for tx hashes with context-aware RPC calls.
- Interpret `status`, `gasUsed`, `logs`, `blockNumber`, `contractAddress`.
- Tie receipts to log decoding (module 09), traces (module 13), and block exploration (module 14).

## Prerequisites
- Modules 05–09 (tx creation/sending, events), 13 (trace), 14 (explorer).

## Real-World Analogy
- Delivery receipt with a success stamp and a list of items delivered (logs).
- CPU analogy: syscall return struct with status + emitted events.

## Steps
1. Parse tx hashes.
2. Call `TransactionReceipt` for each.
3. Print status, gasUsed, log count, block number.
4. Optional: decode event topics into ABI-friendly structs (see module 09).

## Fun Facts & Comparisons
- Status 1 = success, 0 = revert (post-Byzantium). Pre-Byzantium had no status.
- CumulativeGasUsed is per-block order; useful for gas accounting.
- ethers.js: `provider.getTransactionReceipt` same RPC.
- Receipts are stored in a separate trie per block (Merkle-Patricia); the `receiptRoot` commits to them.

## Related Solidity-edu Modules
- 09-events — decode logs from receipts.
- 13-trace — traces are a deeper view of execution; receipts are the summary.
- 14-explorer — stitch receipts into your block viewer.
- 05/06 — tx construction; receipts confirm execution.

## Files
- `exercise/exercise.go`: TODOs for building a receipt fetcher.
- `exercise/solution.go`: reference implementation with defensive copying.
- `exercise/exercise_test.go`: scaffold for your own assertions.
