# Geth Educational Modules - Complete Documentation

This document contains comprehensive documentation for all 25 geth educational modules (01-stack through 25-toolbox), consolidated from individual project READMEs.

---

## Table of Contents

1. [01-stack: Understanding the Ethereum Execution Stack](#01-stack-understanding-the-ethereum-execution-stack)
2. [02-rpc-basics: Deep Dive into JSON-RPC and Block Structures](#02-rpc-basics-deep-dive-into-json-rpc-and-block-structures)
3. [03-keys-addresses: Cryptographic Identity on Ethereum](#03-keys-addresses-cryptographic-identity-on-ethereum)
4. [04-accounts-balances: Understanding Account Types and State](#04-accounts-balances-understanding-account-types-and-state)
5. [05-tx-nonces: Building and Sending Legacy Transactions](#05-tx-nonces-building-and-sending-legacy-transactions)
6. [06-eip1559: Dynamic Fee Transactions (EIP-1559)](#06-eip1559-dynamic-fee-transactions-eip-1559)
7. [07-eth-call: Read-Only Contract Calls](#07-eth-call-read-only-contract-calls)
8. [08-abigen: Typed Contract Bindings](#08-abigen-typed-contract-bindings)
9. [09-events: Decoding ERC20 Transfer Events](#09-events-decoding-erc20-transfer-events)
10. [10-filters: Real-Time Block Monitoring](#10-filters-real-time-block-monitoring)
11. [11-storage: Reading Raw Storage Slots](#11-storage-reading-raw-storage-slots)
12. [12-proofs: Merkle-Patricia Trie Proofs](#12-proofs-merkle-patricia-trie-proofs)
13. [13-trace: Transaction Execution Tracing](#13-trace-transaction-execution-tracing)
14. [14-explorer: Block/Transaction Explorer](#14-explorer-blocktransaction-explorer)
15. [15-receipts: Transaction Receipts and Outcomes](#15-receipts-transaction-receipts-and-outcomes)
16. [16-concurrency: Concurrent RPC Operations with Worker Pools](#16-concurrency-concurrent-rpc-operations-with-worker-pools)
17. [17-indexer: ERC20 Transfer Indexer](#17-indexer-erc20-transfer-indexer)
18. [18-reorgs: Reorg Detection and Handling](#18-reorgs-reorg-detection-and-handling)
19. [19-devnets: Local Devnet Interaction](#19-devnets-local-devnet-interaction)
20. [20-node: Node Information and Health](#20-node-node-information-and-health)
21. [21-sync: Sync Progress Inspection](#21-sync-sync-progress-inspection)
22. [22-peers: Peer Count and P2P Health](#22-peers-peer-count-and-p2p-health)
23. [23-mempool: Mempool Inspection](#23-mempool-mempool-inspection)
24. [24-monitor: Node Health Monitoring](#24-monitor-node-health-monitoring)
25. [25-toolbox: Swiss Army CLI](#25-toolbox-swiss-army-cli)

---

## 01-stack: Understanding the Ethereum Execution Stack

**Goal:** Understand what Geth is, how it fits with consensus clients, and prove connectivity by reading chain ID + latest block.

### Big Picture: The Ethereum Stack from First Principles

Before diving into code, let's build a mental model from the ground up. Ethereum is fundamentally a **distributed state machine**—think of it like a globally synchronized database where everyone agrees on the same sequence of state transitions. But unlike a traditional database, there's no central authority. Instead, we have a **two-client architecture** that emerged from The Merge (Ethereum's transition to Proof-of-Stake):

#### The Two-Client Architecture

**Execution Client (Geth)** = The CPU + Memory + Disk
- **What it does:** Executes EVM bytecode, maintains the state trie (think: Merkle-Patricia tree indexing all account balances and contract storage), and exposes JSON-RPC endpoints
- **Computer Science analogy:** Like a CPU executing instructions, Geth executes transactions. The state trie is like a hash table, but cryptographically verifiable—you can prove "account X has balance Y" without downloading the entire blockchain
- **Fun fact:** Geth stands for "Go Ethereum"—it's written in Go, but there are other execution clients: Erigon (also Go), Nethermind (C#), Besu (Java). They all implement the same EVM spec, so they're interchangeable!

**Consensus Client** (Prysm, Lighthouse, Nimbus, etc.) = The Scheduler + Validator
- **What it does:** Runs the Beacon Chain, manages validators, drives fork choice (decides which chain is canonical), and tells the execution client "execute this block"
- **Computer Science analogy:** Like an operating system scheduler deciding which process runs next, the consensus client decides which block gets appended to the chain
- **Nerdy detail:** The Beacon Chain uses a BFT-style consensus (Casper FFG + LMD GHOST). Validators stake ETH and vote on blocks. If you vote incorrectly, you get slashed (lose ETH). This is why it's called "Proof-of-Stake"

**JSON-RPC** = The API Layer
- **What it does:** Exposes a standardized interface (JSON-RPC 2.0) for querying data and submitting transactions
- **Computer Science analogy:** Like REST APIs for web services, JSON-RPC is the protocol for interacting with Ethereum nodes
- **Protocol detail:** JSON-RPC is stateless and request-response based. Methods like `eth_blockNumber` return data, while `eth_sendTransaction` submits work

### The Complete Picture

```
┌─────────────────────────────────────────────────────────┐
│                    Your Application                      │
│              (Go code using ethclient)                   │
└────────────────────┬────────────────────────────────────┘
                     │ JSON-RPC (HTTP/WebSocket)
                     ▼
┌─────────────────────────────────────────────────────────┐
│              Execution Client (Geth)                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │   EVM Exec   │  │  State Trie  │  │  JSON-RPC    │  │
│  │   Engine     │  │  (Merkle)    │  │  Server      │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└────────────────────┬────────────────────────────────────┘
                     │ Engine API (local IPC)
                     ▼
┌─────────────────────────────────────────────────────────┐
│           Consensus Client (Prysm/Lighthouse)           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ Beacon Chain │  │   Fork       │  │  Validator   │  │
│  │   Logic      │  │   Choice     │  │  Management  │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└────────────────────┬────────────────────────────────────┘
                     │ P2P Gossip Protocol
                     ▼
              ┌──────────────┐
              │ Other Nodes  │
              │  (Peers)     │
              └──────────────┘
```

### Learning Objectives

By the end of this module, you should be able to:

1. **Draw the high-level Ethereum stack:** execution vs consensus vs networking vs JSON-RPC
2. **Use Go + `ethclient` to dial an RPC endpoint** with proper timeout handling (critical for production!)
3. **Query `chainId`, `net_version`, and the latest block header**—these are your "hello world" operations
4. **Interpret the difference between chain ID and network ID:**
   - **Chain ID** (EIP-155): Used for replay protection in transaction signing. Mainnet = 1, Sepolia = 11155111, etc.
   - **Network ID** (legacy): Older identifier, often matches chain ID but not guaranteed. Some networks use different values.
5. **Understand why "public RPC" ≠ "running a node":**
   - Public RPCs (Infura, Alchemy) are convenient but rate-limited
   - They often disable admin/debug endpoints (`debug_traceTransaction`, `admin_*`)
   - Running your own node gives you full power but requires ~1TB disk space and sync time

### Files

- **Exercise:** `exercise/exercise.go` - Your starting point with TODO comments guiding implementation
- **Solution:** `exercise/solution.go` - Complete implementation with detailed educational comments explaining every concept
- **Types:** `exercise/types.go` - Interface and struct definitions
- **Tests:** `exercise/exercise_test.go` - Test suite demonstrating patterns and verifying correctness

---

## 02-rpc-basics: Deep Dive into JSON-RPC and Block Structures

**Goal:** Build a robust JSON-RPC client, call common endpoints, and practice timeouts/retries with full block fetching.

### Big Picture: JSON-RPC as the Universal Interface

JSON-RPC is your "librarian desk" to Ethereum: you ask for data (`eth_blockNumber`, `eth_getBlockByNumber`) or submit work (transactions). Geth exposes this API; hosted RPCs proxy it with rate limits.

#### What is JSON-RPC?

**JSON-RPC** is a stateless, light-weight remote procedure call (RPC) protocol. Think of it like REST APIs, but simpler:
- **Request:** `{"jsonrpc": "2.0", "method": "eth_blockNumber", "params": [], "id": 1}`
- **Response:** `{"jsonrpc": "2.0", "result": "0x1234", "id": 1}`

**Computer Science principle:** RPC (Remote Procedure Call) allows you to call functions on a remote server as if they were local. JSON-RPC uses JSON for serialization, making it language-agnostic. Go's `ethclient` package wraps these JSON-RPC calls into convenient Go methods.

### Building on Module 01

In **module 01**, you learned to:
- Dial an RPC endpoint
- Query chain ID and network ID
- Fetch block **headers** (lightweight metadata)

In **this module**, you'll:
- Fetch **full blocks** (includes transaction data)
- Understand transaction structures
- Add **retry logic** for resilience
- Compare block numbers from different sources

**Key difference:** `HeaderByNumber()` returns just the header (~500 bytes). `BlockByNumber()` returns the full block including all transactions (can be megabytes). Use headers when you only need metadata!

### Learning Objectives

By the end of this module, you should be able to:

1. **Initialize `ethclient` with context timeouts** (building on module 01)
2. **Call common JSON-RPC methods:**
   - `eth_blockNumber` → `client.BlockNumber()`
   - `eth_getBlockByNumber` → `client.BlockByNumber()`
   - `net_version` → `client.NetworkID()`
3. **Inspect block fields:** hash, parentHash, gasUsed, transaction count
4. **Understand the difference between headers and full blocks:**
   - Headers (module 01): ~500 bytes, just metadata
   - Full blocks: 100KB-2MB, includes all transaction data
5. **Implement retry logic** for production resilience
6. **Understand JSON-RPC limitations:** rate limits, missing debug/admin endpoints

### Files

- **Exercise:** `exercise/exercise.go` - Your starting point with TODO comments guiding implementation
- **Solution:** `exercise/solution.go` - Complete implementation with detailed educational comments explaining every concept
- **Types:** `exercise/types.go` - Interface and struct definitions
- **Tests:** `exercise/exercise_test.go` - Test suite demonstrating patterns and verifying correctness

---

## 03-keys-addresses: Cryptographic Identity on Ethereum

**Goal:** Generate Ethereum keys in Go, derive addresses, and understand keystore JSON vs raw private keys.

### Big Picture: From "What" to "Who"

In **modules 01-02**, you learned to query Ethereum nodes—you were asking "what is the latest block?" or "what transactions are in this block?". Now we're moving to **"who is talking?"**—understanding cryptographic identity on Ethereum.

**EOA (Externally Owned Account) identity** = (private key, public key, address). This is the foundation of all Ethereum interactions:
- **Private key:** Your secret (never share this!)
- **Public key:** Derived from private key (can be shared)
- **Address:** Derived from public key (this is what you share publicly)

**Computer Science principle:** This is **public-key cryptography**. The private key can sign messages, and anyone with the public key can verify signatures. But you can't derive the private key from the public key—that's the mathematical foundation (discrete logarithm problem on elliptic curves).

### Learning Objectives

By the end of this module, you should be able to:

1. **Generate a secp256k1 private key** and understand the cryptographic primitives
2. **Derive the public key** from the private key (deterministic process)
3. **Derive the Ethereum address** from the public key (keccak256 hash)
4. **Save/load a key in Geth keystore JSON format** with passphrase encryption
5. **Compare raw hex vs keystore** security trade-offs
6. **Connect addresses to Solidity's `msg.sender`** and access control patterns

### Building on Previous Modules

#### From Module 01 (01-stack)
- You learned to connect to Ethereum nodes via JSON-RPC
- Now you're learning **who** is making those connections

#### From Module 02 (02-rpc-basics)
- You learned about blocks and transactions
- Transactions are **signed** with private keys—this module shows you how those keys work

### What You'll Build

In this module, you'll create a CLI that:
1. Generates a new secp256k1 private key (cryptographically secure)
2. Derives the public key from the private key
3. Derives the Ethereum address from the public key (keccak256 hash)
4. Encrypts the private key into a keystore JSON file (with passphrase)
5. Unlocks the keystore to verify you can recover the same address

**Key learning:** You'll understand the complete flow from private key → public key → address. This is fundamental to all Ethereum interactions!

### Files

- **Exercise:** `exercise/exercise.go` - Your starting point with TODO comments guiding implementation
- **Solution:** `exercise/solution.go` - Complete implementation with detailed educational comments explaining every concept
- **Types:** `exercise/types.go` - Interface and struct definitions
- **Tests:** `exercise/exercise_test.go` - Test suite demonstrating patterns and verifying correctness

---

## 04-accounts-balances: Understanding Account Types and State

**Goal:** Classify EOAs vs contracts and query balances, understanding the fundamental difference between account types on Ethereum.

### Big Picture: Two Types of Accounts

Ethereum accounts come in two flavors: **EOAs (Externally Owned Accounts)** and **Contracts**. Understanding this distinction is fundamental to Ethereum development.

**Computer Science principle:** This is a **type system** at the blockchain level. Just like programming languages have types (int, string, struct), Ethereum has account types. The type determines what operations are possible.

#### EOA (Externally Owned Account)
- **Has:** Private key, address, balance, nonce
- **Does NOT have:** Code (bytecode)
- **Can:** Send transactions, sign messages
- **Cannot:** Execute arbitrary code
- **Analogy:** An empty plot of land with a mailbox (address) and a safe (balance)

#### Contract Account
- **Has:** Address, balance, nonce, **code** (bytecode)
- **Does NOT have:** Private key (cannot initiate transactions directly)
- **Can:** Execute code when called, store state, emit events
- **Cannot:** Sign transactions (needs an EOA to call it)
- **Analogy:** A building with machinery (code) on the same street (address)

### Learning Objectives

By the end of this module, you should be able to:

1. **Fetch balances** at a specific block (or latest) in wei
2. **Detect account type** by checking for code presence
3. **Understand special cases:**
   - Precompiles (addresses 0x01-0x09) have code but are special-purpose
   - Selfdestructed contracts have nonce > 0 but code size 0
4. **Distinguish between** `eth_getBalance` and `eth_getCode` use cases

### What You'll Build

In this module, you'll create a CLI that:
1. Takes one or more addresses as command-line arguments
2. Queries the balance of each address (in wei)
3. Queries the code of each address
4. Classifies each address as EOA or Contract
5. Displays address, type, and balance

**Key learning:** You'll understand the fundamental distinction between EOAs and contracts, and how to query account state on the blockchain!

### Files

- **Exercise:** `exercise/exercise.go` - Your starting point with TODO comments guiding implementation
- **Solution:** `exercise/solution.go` - Complete implementation with detailed educational comments explaining every concept
- **Types:** `exercise/types.go` - Interface and struct definitions
- **Tests:** `exercise/exercise_test.go` - Test suite demonstrating patterns and verifying correctness

---

## 05-tx-nonces: Building and Sending Legacy Transactions

**Goal:** Build/send a legacy transaction, manage nonces, and understand replay protection.

### Big Picture: Transaction Lifecycle

Transactions are the fundamental unit of state change on Ethereum. Understanding how to build, sign, and send transactions is essential for any Ethereum developer.

**Computer Science principle:** Transactions are **immutable, ordered messages** that change blockchain state. They're like database transactions, but cryptographically signed and globally ordered.

### Learning Objectives

By the end of this module, you should be able to:

1. **Fetch pending nonce** for an address (includes pending transactions)
2. **Build a legacy transaction** with gasPrice
3. **Sign a transaction** with EIP-155 replay protection (chainID)
4. **Send a transaction** to the network
5. **Understand nonce ordering** and why gaps stall subsequent transactions

### Building on Previous Modules

#### From Module 03 (03-keys-addresses)
- You learned to generate private keys and derive addresses
- Now you're using those keys to **sign transactions**
- The address you derive is the `from` address in transactions

#### From Module 04 (04-accounts-balances)
- You learned to query account balances
- Now you're **changing** those balances by sending transactions
- The nonce you fetch is stored in the account state

#### From Module 01 (01-stack)
- You learned about chainID (EIP-155 replay protection)
- Now you're using chainID in transaction signatures
- This prevents transactions from being replayed on other chains

### What You'll Build

In this module, you'll create a CLI that:
1. Takes recipient address, ETH amount, and private key as input
2. Fetches the pending nonce for the sender
3. Fetches the suggested gas price
4. Builds a legacy transaction
5. Signs the transaction with EIP-155 (chainID) protection
6. Sends the transaction to the network
7. Displays transaction hash and details

**Key learning:** You'll understand the complete transaction lifecycle from building to broadcasting!

### Files

- **Exercise:** `exercise/exercise.go` - Your starting point with TODO comments guiding implementation
- **Solution:** `exercise/solution.go` - Complete implementation with detailed educational comments explaining every concept
- **Types:** `exercise/types.go` - Interface and struct definitions
- **Tests:** `exercise/exercise_test.go` - Test suite demonstrating patterns and verifying correctness

---

## 06-eip1559: Dynamic Fee Transactions (EIP-1559)

**Goal:** Build/send an EIP-1559 (dynamic fee) transaction and understand base fee + tip math.

### Big Picture: The London Upgrade and Dynamic Fees

Post-London (August 2021), Ethereum uses **dynamic fees** instead of fixed gas prices. This makes fee estimation more predictable and reduces fee volatility.

**Computer Science principle:** EIP-1559 introduces a **two-part fee structure**:
- **Base Fee:** Algorithmically determined, **burned** (removed from supply)
- **Priority Fee (Tip):** Paid to validators/miners, incentivizes inclusion

This is more efficient than the legacy auction model where users bid against each other.

### Learning Objectives

By the end of this module, you should be able to:

1. **Construct a DynamicFeeTx** with maxFeePerGas and maxPriorityFeePerGas
2. **Convert user inputs** (gwei) to wei safely
3. **Sign with London signer** (chainID-aware) and broadcast
4. **Understand fee math:** effectiveGasPrice = min(maxFeeCap, baseFee + tip)
5. **Understand refunds:** Excess fees are refunded to the sender

### Building on Previous Modules

#### From Module 05 (05-tx-nonces)
- You learned to build legacy transactions with fixed `gasPrice`
- Now you're building **EIP-1559 transactions** with dynamic fees
- Same nonce management, same signing process, different fee structure

#### From Module 01 (01-stack)
- You learned about chainID
- EIP-1559 transactions also use chainID for replay protection
- London signer includes chainID in the signature

### What You'll Build

In this module, you'll create a CLI that:
1. Takes recipient address, ETH amount, private key, and fee caps as input
2. Fetches the pending nonce (same as module 05)
3. Converts gwei to wei for fee caps
4. Builds an EIP-1559 DynamicFeeTx
5. Signs with London signer (includes chainID)
6. Sends the transaction to the network
7. Displays transaction hash and fee details

**Key learning:** You'll understand the modern transaction format and dynamic fee mechanics!

### Files

- **Exercise:** `exercise/exercise.go` - Your starting point with TODO comments guiding implementation
- **Solution:** `exercise/solution.go` - Complete implementation with detailed educational comments explaining every concept
- **Types:** `exercise/types.go` - Interface and struct definitions
- **Tests:** `exercise/exercise_test.go` - Test suite demonstrating patterns and verifying correctness

---

## 07-eth-call: Read-Only Contract Calls

**Goal:** Perform read-only contract calls with manual ABI encoding/decoding.

### Big Picture: Simulating Transactions Without State Changes

`eth_call` simulates a transaction **without persisting state**. You encode function selectors and arguments per ABI, send to the node, and decode the return data. No gas is spent on-chain, but the node executes the EVM locally.

**Computer Science principle:** This is like a **dry run** or **read-only query**. The EVM executes the code, but no state changes are committed. This is perfect for querying view/pure functions.

### The Difference: `eth_call` vs `eth_sendTransaction`

| Aspect | `eth_call` | `eth_sendTransaction` |
|--------|------------|----------------------|
| State changes | ❌ No (simulated) | ✅ Yes (persisted) |
| Gas cost | ❌ No (free) | ✅ Yes (paid) |
| Transaction hash | ❌ No | ✅ Yes |
| Use case | Querying data | Changing state |
| Speed | Fast (local execution) | Slower (needs mining) |

**Key insight:** `eth_call` is for **reading**, `eth_sendTransaction` is for **writing**.

### Learning Objectives

By the end of this module, you should be able to:

1. **Pack ABI** for simple view functions (ERC20 name/symbol/decimals/totalSupply)
2. **Call contracts** with `CallContract` and decode results
3. **Handle reverts** and raw return data
4. **Understand ABI encoding** (function selector + arguments)
5. **Decode return values** based on function return types

### What You'll Build

In this module, you'll create a CLI that:
1. Takes a contract address and function name as input
2. Encodes the function call using ABI (function selector)
3. Executes `eth_call` to simulate the function execution
4. Decodes the return value based on function type
5. Displays the result

**Key learning:** You'll understand how to manually encode/decode ABI data, giving you full control over contract interactions!

### Files

- **Starter:** `exercise/exercise.go` - Student entry point with TODO guidance
- **Solution:** `exercise/solution.go` - Reference implementation (run with `go test -tags solution ./07-eth-call/...`)
- **Tests:** `exercise/exercise_test.go` - Covers ABI encoding/decoding edge cases

---

## 08-abigen: Typed Contract Bindings

**Goal:** Use typed contract bindings (abigen-style) for safer calls and transactions.

### Big Picture: From Manual to Typed

abigen turns ABI into typed Go methods. Instead of manual ABI pack/unpack (module 07), you get compile-time checked methods with `CallOpts` and `TransactOpts` for read/write. This reduces boilerplate and errors.

**Computer Science principle:** This is **code generation** - converting a schema (ABI) into type-safe code. It's like generating API clients from OpenAPI specs, or database models from schemas.

### Learning Objectives

By the end of this module, you should be able to:

1. **Understand how abigen bindings** wrap `BoundContract` with typed methods
2. **Make typed view calls** (name/symbol/decimals/balanceOf)
3. **See how `CallOpts` and `TransactOpts`** carry context, block number, signer info
4. **Compare manual vs typed** approaches and their trade-offs

### Building on Previous Modules

#### From Module 07 (07-eth-call)
- You learned manual ABI encoding/decoding
- Now you're using **typed bindings** that handle encoding/decoding automatically
- Same underlying JSON-RPC calls, but with better ergonomics

### What You'll Build

In this module, you'll create a CLI that:
1. Takes a token address and optional holder address as input
2. Creates a BoundContract from ABI and address
3. Calls typed methods (Name, Symbol, Decimals, BalanceOf)
4. Displays token information and balance

**Key learning:** You'll see how typed bindings simplify contract interactions while maintaining type safety!

### Files

- **Starter:** `exercise/exercise.go`
- **Solution:** `exercise/solution.go` (build with `-tags solution`)
- **Tests:** `exercise/exercise_test.go`

---

## 09-events: Decoding ERC20 Transfer Events

**Goal:** Decode ERC20 Transfer logs and understand topics vs data.

### Big Picture: Events as Append-Only History

Events/logs are append-only "newspaper clippings" emitted during transaction execution. Indexed params go into topics (bloom-filtered for search); non-indexed go into data. This is the off-chain friendly history of state changes.

**Computer Science principle:** Events are like **write-ahead logs** or **audit trails**. They provide a searchable history of state changes without storing full state snapshots.

### Learning Objectives

By the end of this module, you should be able to:

1. **Build a filter query** for a token's Transfer events over a block range
2. **Decode indexed vs non-indexed** params with ABI
3. **Understand log roots/bloom filters** in block headers
4. **Filter events** by address and topic

### Building on Previous Modules

#### From Module 08 (08-abigen)
- You learned to call contract functions
- Now you're **listening to events** emitted by those functions
- Events complement function calls - they show what happened

#### From Module 02 (02-rpc-basics)
- You learned about blocks and transactions
- Events are stored in **transaction receipts** (module 15)
- Block headers include `logsBloom` for efficient event queries

### What You'll Build

In this module, you'll create a CLI that:
1. Takes a token address and block range as input
2. Builds a FilterQuery for Transfer events
3. Fetches logs matching the filter
4. Decodes indexed parameters (from, to) from topics
5. Decodes non-indexed parameters (value) from data
6. Displays Transfer events with block number, tx hash, from, to, value

**Key learning:** You'll understand how events work, how to filter them, and how to decode them!

### Files

- **Starter:** `exercise/exercise.go`
- **Solution:** `exercise/solution.go` (build with `-tags solution`)
- **Tests:** `exercise/exercise_test.go`

---

## 10-filters: Real-Time Block Monitoring

**Goal:** Practice filters and subscriptions (newHeads), and understand polling vs websockets.

### Big Picture: Push vs Pull

Filters let you query past logs; subscriptions push new data (heads/logs) over websockets. When WS isn't available, you poll. Detecting reorgs means comparing parent hashes to what you stored.

**Computer Science principle:** This is the **push vs pull** pattern:
- **Pull (Polling):** Client asks "any updates?" periodically
- **Push (Subscriptions):** Server tells client "here's an update!" immediately

### Learning Objectives

By the end of this module, you should be able to:

1. **Subscribe to `newHeads`** over WebSocket
2. **Poll latest headers** over HTTP fallback
3. **Understand reorg detection** via parent hash mismatch
4. **Compare push vs pull** approaches

### Building on Previous Modules

#### From Module 09 (09-events)
- You learned to filter logs (historical queries)
- Now you're **subscribing to new blocks** (real-time monitoring)
- Same filtering concepts, different protocol (WebSocket vs HTTP)

### What You'll Build

In this module, you'll create a CLI that:
1. Supports WebSocket mode (subscriptions) or HTTP mode (polling)
2. If WebSocket: Subscribe to new block headers, print as they arrive
3. If HTTP: Poll latest N blocks and print headers
4. Display block number, hash, and parent hash
5. Show how to detect reorgs (parent hash mismatch)

**Key learning:** You'll understand real-time monitoring vs polling, and how to detect chain reorganizations!

### Files

- **Starter:** `exercise/exercise.go`
- **Solution:** `exercise/solution.go` (run with `go test -tags solution ./10-filters/...`)
- **Tests:** `exercise/exercise_test.go`

---

## 11-storage: Reading Raw Storage Slots

**Goal:** Read raw storage slots directly from contracts, including mapping slots, and connect to Solidity storage layout.

### Big Picture: Storage as a Cryptographic Database

Ethereum's storage is like a **cryptographic hash table** where every contract has 2^256 possible 32-byte slots. Unlike traditional databases, storage slots are:
- **Immutable** (once written, can't be changed without a transaction)
- **Cryptographically verifiable** (committed to the state root in block headers)
- **Deterministic** (same contract code + same inputs = same storage layout)

**Computer Science principle:** Storage is organized as a **Merkle-Patricia trie**. The `stateRoot` in block headers is the root hash of this trie. Every storage slot is part of this tree, making it possible to prove "contract X has value Y in slot Z" without downloading the entire state.

### Learning Objectives

By the end of this module, you should be able to:

1. **Understand storage slot layout:**
   - Simple variables: Direct slot access (slot 0, 1, 2, ...)
   - Mappings: `slot = keccak256(abi.encode(key, baseSlot))`
   - Dynamic arrays: `base = keccak256(slot)`, then `base + index`
   - Packed variables: Multiple small types in one slot

2. **Call `eth_getStorageAt` via Go's `StorageAt` method:**
   - Read raw 32-byte values from any slot
   - Understand that values are returned as raw bytes (you must decode)

3. **Compute mapping slot hashes:**
   - Hash the key and base slot together
   - Use proper padding (32 bytes for each component)

4. **Connect to Solidity storage layout:**
   - Relate Go storage reads to Solidity variable declarations
   - Understand packed vs unpacked storage
   - Decode common types (uint256, address, bool)

### What You'll Build

In this module, you'll create a CLI that:
1. Takes a contract address and storage slot number
2. Optionally takes a mapping key
3. Computes the correct storage slot (with mapping hash if needed)
4. Reads raw 32-byte values via `eth_getStorageAt`
5. Displays the raw hex-encoded value

**Key learning:** You'll understand how Solidity's storage layout translates to actual on-chain storage slots. This is essential for:
- Debugging contract state
- Building indexers
- Verifying storage proofs
- Understanding gas costs

### Files

- **Exercise:** `exercise/exercise.go` - Your starting point with TODO comments guiding implementation
- **Solution:** `exercise/solution.go` - Complete implementation with detailed educational comments explaining every concept
- **Types:** `exercise/types.go` - Interface and struct definitions
- **Tests:** `exercise/exercise_test.go` - Test suite demonstrating patterns and verifying correctness

---

## 12-proofs: Merkle-Patricia Trie Proofs

**Goal:** Fetch and interpret Merkle-Patricia trie proofs for accounts and storage slots via `eth_getProof`.

### Big Picture: Cryptographic Proofs for Trust-Minimized Verification

**Merkle-Patricia trie proofs** are cryptographic receipts that prove "account X has balance Y and storage slot Z has value W at block N" without downloading the entire blockchain state. This enables:
- **Light clients:** Verify state without syncing full blockchain
- **Bridges:** Prove state on one chain to another chain
- **Indexers:** Verify indexed data is correct
- **Wallets:** Check balances without trusting a single RPC endpoint

**Computer Science principle:** Merkle trees allow you to prove membership in a set using only a logarithmic number of hashes. Instead of downloading 100GB of state, you download a few KB of proof nodes.

### Learning Objectives

By the end of this module, you should be able to:

1. **Call `eth_getProof` via Go's `GetProof` method:**
   - Request proofs for accounts and storage slots
   - Specify block number (or use latest)
   - Understand proof structure

2. **Interpret proof results:**
   - Account proof: balance, nonce, codeHash, storageHash
   - Storage proof: slot value and proof nodes
   - Proof nodes: Merkle tree path from root to leaf

3. **Understand trust-minimized verification:**
   - How proofs enable verification without full state
   - Why light clients need proofs
   - How bridges use proofs for cross-chain verification

4. **Connect proofs to storage slots:**
   - Proof paths use same slot calculations as module 11
   - Storage proofs prove specific slot values
   - Account proofs prove account state

### What You'll Build

In this module, you'll create a CLI that:
1. Takes an account address (and optional storage slot)
2. Calls `eth_getProof` to fetch Merkle-Patricia trie proofs
3. Displays account state (balance, nonce, codeHash, storageHash)
4. Displays storage proof (if slot provided)
5. Shows proof node counts

**Key learning:** You'll understand how cryptographic proofs enable trust-minimized verification. This is essential for:
- Building light clients
- Implementing cross-chain bridges
- Verifying indexed data
- Building trustless applications

### Files

- **Exercise:** `exercise/exercise.go` - Your starting point with TODO comments guiding implementation
- **Solution:** `exercise/solution.go` - Complete implementation with detailed educational comments explaining every concept
- **Types:** `exercise/types.go` - Interface and struct definitions
- **Tests:** `exercise/exercise_test.go` - Test suite demonstrating patterns and verifying correctness

---

## 13-trace: Transaction Execution Tracing

**Goal:** Trace transaction execution with `debug_traceTransaction` to see call tree, gas usage, and opcode-level details.

### Big Picture: Execution Instrumentation

Transaction tracing is **execution instrumentation**—observing program execution without modifying it. Ethereum's EVM is deterministic, so replaying a transaction always produces the same trace. This is like having a debugger attached to every smart contract!

**Computer Science principle:** Tracing is a form of **program analysis**. By instrumenting execution, we can:
- Measure performance (gas usage per operation)
- Debug behavior (see exact execution flow)
- Analyze security (identify suspicious call patterns)
- Optimize code (find expensive operations)

### Learning Objectives

By the end of this module, you should be able to:

1. **Call `debug_traceTransaction`** via Go's `TraceTransaction` method
2. **Understand trace structure** (call tree, gas usage, storage changes)
3. **Debug transaction failures** using traces
4. **Analyze gas consumption** at the opcode level
5. **Understand tracer types** (default, callTracer, prestateTracer, etc.)

### Building on Previous Modules

#### From Module 05-06 (Transactions)
- You learned to build and send transactions
- Now you're seeing what happened during execution
- Traces show the "inside view" of transaction processing

### What You'll Build

In this module, you'll create a function that:
1. Takes an RPC client as input
2. Calls the TraceTransaction method
3. Interprets nil (synced) vs non-nil (syncing) responses
4. Returns structured status information

**Key learning:** You'll understand how to check if an Ethereum node is ready for production use!

### Files

- **Exercise:** `exercise/exercise.go` - Your starting point with TODO comments guiding implementation
- **Solution:** `exercise/solution.go` - Complete implementation with detailed educational comments
- **Types:** `exercise/types.go` - Interface and struct definitions
- **Tests:** `exercise/exercise_test.go` - Test suite demonstrating patterns

---

## 14-explorer: Block/Transaction Explorer

**Goal:** Build a tiny block/tx explorer that fetches a block, summarizes its header, and (optionally) lists transactions.

### Why this matters
- Block explorers are just structured RPC clients: read a block, decode fields, and render a human-friendly view.
- Reinforces earlier modules:
  - 01-stack/02-rpc-basics: dialing RPC with context and handling nil responses.
  - 05/06: inspecting tx fields (gas, to, hash) without sending anything.
  - 13-trace/15-receipts: richer views layer on top of the same block anchor.

### First-principles CS
- **Data locality:** headers are small (~hundreds of bytes) vs. full blocks with tx objects (KBs). Fetch only what you need.
- **Interfaces over structs:** a tiny `RPCClient` interface keeps the explorer decoupled from concrete ethclient implementations.
- **Defensive copying:** types.Block shares internal slices; copying header fields prevents accidental mutation by callers.

### Files
- `exercise/exercise.go`: TODOs guide you to validate inputs, fetch a block, and build a small Result with optional Tx summaries.
- `exercise/solution.go`: Reference implementation with commentary on choices (defensive copies, minimal TxSummary fields).
- `exercise/exercise_test.go`: currently a scaffold—extend it to test your own explorer output.

---

## 15-receipts: Transaction Receipts and Outcomes

**Goal:** Fetch transaction receipts and classify success/failure, logs, and gas usage.

### Big Picture

Receipts record the outcome of a transaction: status, cumulative gas, and emitted logs. They live alongside blocks but outside the state trie. Decoding receipts is key for dApps, indexers, and monitoring.

### Learning Objectives
- Fetch receipts for tx hashes with context-aware RPC calls.
- Interpret `status`, `gasUsed`, `logs`, `blockNumber`, `contractAddress`.
- Tie receipts to log decoding (module 09), traces (module 13), and block exploration (module 14).

### Prerequisites
- Modules 05–09 (tx creation/sending, events), 13 (trace), 14 (explorer).

### Files
- `exercise/exercise.go`: TODOs for building a receipt fetcher.
- `exercise/solution.go`: reference implementation with defensive copying.
- `exercise/exercise_test.go`: scaffold for your own assertions.

---

## 16-concurrency: Concurrent RPC Operations with Worker Pools

**Goal:** fetch multiple resources concurrently with a worker pool while respecting RPC limits.

### Big Picture

RPCs can be slow; fan-out speeds things up, but you need to avoid rate limits and handle cancellation. Worker pools with contexts let you balance throughput and safety. You'll reuse this pattern in indexers (module 17) and monitors (module 24).

### Learning Objectives
- Build a simple worker pool (jobs/results channels + WaitGroup).
- Use contexts to cancel/timeout concurrent work.
- Understand rate limiting/backoff considerations and how to bound goroutines.

### Prerequisites
- Modules 01–10; Go concurrency basics (goroutines, channels, context).

### Key Patterns You'll Learn

#### Pattern 1: Worker Pool with Channels
Bounded concurrency prevents resource exhaustion. Unlike unbounded goroutines (one per endpoint), worker pools limit concurrency to a fixed number of workers.

#### Pattern 2: Context Propagation in Concurrent Operations
Prevents one slow operation from consuming the entire timeout budget. Each probe gets a bounded time slice.

#### Pattern 3: Mutex-Protected Map Access
Maps in Go are not safe for concurrent access. Multiple goroutines writing simultaneously causes data races and panics.

#### Pattern 4: WaitGroup for Synchronization
Main goroutine must wait for workers to complete before returning results. Without WaitGroup, workers would be killed mid-processing.

### Files

- **Exercise:** `exercise/exercise.go` - Your starting point with TODO comments guiding implementation
- **Solution:** `exercise/solution.go` - Complete implementation with detailed educational comments explaining every concept
- **Tests:** `exercise/exercise_test.go` - Test suite demonstrating patterns and verifying correctness

---

## 17-indexer: ERC20 Transfer Indexer

**Goal:** build a basic ERC20 Transfer indexer into sqlite with a simple query surface.

### Big Picture

Indexers watch events/logs and persist them into databases for fast queries. This module scans a block range for Transfer logs, decodes them, and stores into sqlite. Production indexers add pagination, reorg handling, and richer schemas.

### Learning Objectives
- Construct a filter for Transfer logs.
- Decode indexed/non-indexed params and persist to sqlite.
- Consider reorg handling (hash/number pairs) and pagination.

### Prerequisites
- Modules 09 (events), 15 (receipts), 16 (concurrency) helpful.
- Basic SQL familiarity.

### Files
- Starter: `cmd/geth-17-indexer/main.go`
- Solution: `cmd/geth-17-indexer_solution/main.go`

---

## 18-reorgs: Reorg Detection and Handling

**Goal:** detect reorgs by comparing stored block hashes to parent hashes; learn how to rewind/rescan.

### Big Picture

Reorgs happen when a different chain of blocks becomes canonical. Shallow reorgs are normal; indexers and monitors must detect them and replay affected ranges. Comparing parentHash to previously stored hash reveals mismatches.

### Learning Objectives
- Fetch sequential blocks and track hash/parentHash.
- Detect parent mismatches (reorg hint).
- Understand how to rollback and rescan safely.

### Prerequisites
- Modules 09–17 (events, indexing).

### Files
- Starter: `cmd/geth-18-reorgs/main.go`
- Solution: `cmd/geth-18-reorgs_solution/main.go`

---

## 19-devnets: Local Devnet Interaction

**Goal:** interact with a local devnet (e.g., anvil mainnet fork) and inspect balances/heads.

### Big Picture

Devnets let you fork mainnet state and safely test flows: impersonate accounts, fund addresses, and send txs without real risk. Anvil/Hardhat provide JSON-RPC compatible endpoints.

### Learning Objectives
- Dial a local devnet and query balances/head.
- Understand forked state vs live mainnet.
- Impersonation/funding basics (via anvil flags).

### Prerequisites
- Modules 01–05 (RPC, balances, tx basics).

### Files
- Starter: `cmd/geth-19-devnets/main.go`
- Solution: `cmd/geth-19-devnets_solution/main.go`

---

## 20-node: Node Information and Health

**Goal:** query basic node info (client version, peer count, sync status) from your own Geth node.

### Big Picture

Running your own node grants full control and fresh data. Basic health signals: client version, peer count, and sync progress. Some APIs (admin_*, txpool_*) require enabling extra modules or IPC.

### Learning Objectives
- Call `web3_clientVersion` and `net_peerCount`.
- Check sync progress via `SyncProgress`.
- Understand limits of public RPC vs your own node.

### Prerequisites
- Modules 01–05.

### Files
- Starter: `cmd/geth-20-node/main.go`
- Solution: `cmd/geth-20-node_solution/main.go`

---

## 21-sync: Sync Progress Inspection

**Goal:** inspect sync progress and understand full/snap/light modes.

### Big Picture

Sync modes: full replays all blocks, snap downloads snapshots then heals, light fetches proofs on demand. `SyncProgress` reports current vs highest block and state sync counters. Nil means the node believes it is synced.

**Computer Science principle:** Nil as a sentinel value—the absence of sync progress indicates completion. This is a classic pattern where "nothing to report" is meaningful information.

### Learning Objectives

By the end of this module, you should be able to:

1. **Call SyncProgress** and interpret nil vs non-nil responses
2. **Differentiate sync modes** conceptually (full/snap/light)
3. **Spot stale nodes** by checking progress or head lag
4. **Understand sentinel values** as a design pattern for status reporting

### Files

- **Exercise:** `exercise/exercise.go` - Your starting point with TODO comments
- **Solution:** `exercise/solution.go` - Complete implementation with detailed educational comments
- **Types:** `exercise/types.go` - Interface and struct definitions
- **Tests:** `exercise/exercise_test.go` - Test suite verifying correctness

---

## 22-peers: Peer Count and P2P Health

**Goal:** query peer count and understand p2p gossip health.

### Big Picture

Peers gossip txs/blocks across the network. Peer count is a coarse signal of connectivity; richer info lives in admin APIs. Public RPCs often hide peer details.

**Computer Science principle:** Peer-to-peer mesh networks rely on redundant connections for fault tolerance and data propagation. Peer count is a basic metric for assessing network health.

### Learning Objectives

By the end of this module, you should be able to:

1. **Call net_peerCount** and understand the returned value
2. **Interpret peer count** as a connectivity health metric
3. **Recognize limitations** of public RPC vs your own node
4. **Understand P2P networking** fundamentals (discovery, gossip, protocols)

### Files

- **Exercise:** `exercise/exercise.go` - Your starting point with TODO comments
- **Solution:** `exercise/solution.go` - Complete implementation with detailed educational comments
- **Types:** `exercise/types.go` - Interface and struct definitions
- **Tests:** `exercise/exercise_test.go` - Test suite verifying correctness

---

## 23-mempool: Mempool Inspection

**Goal:** inspect pending transactions (where supported) and understand mempool visibility limits.

### Big Picture

Pending txs live in the txpool before inclusion. Many public RPCs do not expose full mempool; Geth's `eth_pendingTransactions` or `txpool_*` may be restricted. Visibility varies by provider and node config.

**Computer Science principle:** Queue management with priority—mempools are priority queues where higher-fee transactions get processed first. Privacy/security trade-offs limit visibility to prevent MEV exploitation.

### Learning Objectives

By the end of this module, you should be able to:

1. **Query mempool size** and understand congestion levels
2. **Understand visibility limits** of public vs private mempools
3. **Recognize MEV implications** of mempool transparency
4. **Learn transaction replacement** rules (Replace-By-Fee)

### Files

- **Exercise:** `exercise/exercise.go` - TODOs guide implementation
- **Solution:** `exercise/solution.go` - Full implementation with educational comments
- **Types:** `exercise/types.go` - Interface and struct definitions
- **Tests:** `exercise/exercise_test.go` - Test suite

---

## 24-monitor: Node Health Monitoring

**Goal:** implement basic node health checks (block freshness/lag) and discuss alerting patterns.

### Big Picture

Monitoring a node involves tracking head freshness, RPC latency, and error rates. Simple checks catch stale nodes early; production systems export metrics to Prometheus/Grafana and alert on thresholds.

**Computer Science principle:** Threshold-based health checks convert continuous metrics (lag time) into discrete states (OK/STALE) for actionable alerting.

### Learning Objectives

1. **Fetch latest header** and extract timestamp
2. **Calculate lag** between block time and current time
3. **Classify status** using configurable thresholds
4. **Understand monitoring patterns** for production systems

### Files

- **Exercise:** `exercise/exercise.go` - TODOs guide implementation
- **Solution:** `exercise/solution.go` - Full implementation with educational comments
- **Types:** `exercise/types.go` - Interface and struct definitions
- **Tests:** `exercise/exercise_test.go` - Test suite

---

## 25-toolbox: Swiss Army CLI

**Goal:** build a Swiss Army CLI that combines status, block/tx lookup, and event decoding.

### Big Picture

Capstone module that stitches together previous lessons into one tool with subcommands. Reuses RPC basics, block/tx retrieval, receipts/logs decoding, and event filtering.

**Computer Science principle:** Composition—building complex systems from simple, well-tested components. Command pattern for operation dispatch.

### Learning Objectives

1. **Implement command routing** (dispatch to different handlers)
2. **Compose multiple RPC operations** into unified commands
3. **Reuse patterns from previous modules** (01-stack, 21-24)
4. **Build production-ready CLI tools** with subcommands

### Supported Commands

#### status
Comprehensive node status combining modules 01, 21, 22:
- Chain ID & Network ID
- Latest block number & hash
- Sync status (syncing or synced)
- Peer count

#### block <number>
Block details:
- Number, hash, parent hash
- Timestamp
- Transaction count
- Gas used/limit

#### tx <hash>
Transaction details:
- Hash, nonce, value
- Gas, gas price
- To address (if not contract creation)
- Pending status

### Files

- **Exercise:** `exercise/exercise.go` - TODOs guide implementation
- **Solution:** `exercise/solution.go` - Full implementation with command handlers
- **Types:** `exercise/types.go` - Unified interface combining all modules
- **Tests:** `exercise/exercise_test.go` - Test suite

### Next Steps

Congratulations! You've completed the geth modules series. You now understand:
- RPC client patterns and defensive programming
- Blockchain data structures (headers, blocks, transactions)
- Operational metrics (sync status, peers, mempool, health)
- Building production-ready CLI tools

**What's next?**
- Build your own tools using these patterns
- Explore advanced topics: contract calls, event filtering, indexing
- Contribute to Ethereum tooling ecosystem!

---

## Conclusion

This comprehensive documentation covers all 25 geth educational modules, organized from foundational concepts (01-stack) through advanced operations (25-toolbox). Each module builds upon previous ones, teaching:

- **Modules 01-04:** RPC basics, key management, account querying
- **Modules 05-09:** Transaction building, contract interactions, event decoding
- **Modules 10-13:** Real-time monitoring, storage, proofs, execution tracing
- **Modules 14-16:** Block explorers, receipts, concurrent operations
- **Modules 17-20:** Indexing, reorg handling, devnets, node operations
- **Modules 21-25:** Operational monitoring, health checks, and comprehensive tooling

Each module emphasizes practical patterns, educational depth, and production-ready implementations.
