# geth-14-explorer

**Goal:** Build a tiny block/tx explorer that fetches a block, summarizes its header, and (optionally) lists transactions.

## Why this matters
- Block explorers are just structured RPC clients: read a block, decode fields, and render a human-friendly view.
- Reinforces earlier modules:
  - 01-stack/02-rpc-basics: dialing RPC with context and handling nil responses.
  - 05/06: inspecting tx fields (gas, to, hash) without sending anything.
  - 13-trace/15-receipts: richer views layer on top of the same block anchor.

## First-principles CS
- **Data locality:** headers are small (~hundreds of bytes) vs. full blocks with tx objects (KBs). Fetch only what you need.
- **Interfaces over structs:** a tiny `RPCClient` interface keeps the explorer decoupled from concrete ethclient implementations.
- **Defensive copying:** types.Block shares internal slices; copying header fields prevents accidental mutation by callers.

## Nerdy/fun facts
- Etherscan-like explorers cache block summaries in databases; we’re just doing an on-demand read.
- Tx objects can omit `From`—you need a signer and chain ID to recover it (see module 05/06 for signing context).
- A “block explorer” is really a “content-addressable ledger browser”: the block hash is your content address.

## Exercise shape
- `exercise.go`: TODOs guide you to validate inputs, fetch a block, and build a small Result with optional Tx summaries.
- `solution.go`: Reference implementation with commentary on choices (defensive copies, minimal TxSummary fields).
- `exercise_test.go`: currently a scaffold—extend it to test your own explorer output.

## Suggested extensions
- Pull receipts for each tx (module 15) to show status/gasUsed.
- Add trace links (module 13) for deep debugging.
- Render logs (module 09) for contract activity insights.

## Files
- `exercise/exercise.go`: student TODOs.
- `exercise/solution.go`: reference implementation.
- `exercise/exercise_test.go`: add your own assertions here.
