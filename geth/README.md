# Geth — Ethereum Development (25 Projects)

Learn Ethereum development and go-ethereum (Geth) internals through production-grade projects.

## 📚 Project Overview

Each project in this directory is **self-contained** and follows the same structure:

```
project-name/
├── README.md                    # (removed for cleaner structure)
├── cmd/
│   ├── app/main.go             # Realistic RPC client application
│   └── dev/main.go             # Debug harness (fixed RPC endpoint)
└── internal/
    └── <pkg>/
        ├── exercise.go         # YOUR implementation (stubs with TODO comments)
        ├── exercise_test.go    # Tests to verify your work
        ├── types.go            # Geth-style type definitions
        ├── solution.reference.go        # Reference implementation (excluded from builds)
        └── solution_no_err.reference.go # Alternative reference
```

## 🎯 How to Use

### 1. Set Up RPC Endpoint

```bash
# Option 1: Environment variable
export INFURA_RPC_URL=https://mainnet.infura.io/v3/YOUR_PROJECT_ID

# Option 2: Use a public endpoint
export INFURA_RPC_URL=https://eth.llamarpc.com
```

### 2. Navigate to a Project

```bash
cd geth/01-stack
```

### 3. Implement in `exercise.go`

Open `internal/<pkg>/exercise.go` and implement the functions marked with `TODO` comments.

### 4. Run Tests

```bash
go test ./...
```

Some tests may require a live RPC endpoint. Use environment variables or modify test fixtures.

### 5. Debug with `cmd/dev`

```bash
go run ./cmd/dev
```

Fixed RPC endpoint and test data — perfect for stepping through with a debugger.

### 6. Run the Application

```bash
go run ./cmd/app https://eth.llamarpc.com
```

Realistic RPC client demonstrating how your code would be consumed.

### 7. Compare with Reference

```bash
cat internal/<pkg>/solution.reference.go
```

Reference implementations are **excluded** from normal builds (build tag: `reference`).

---

## 📋 Complete Project List

### Foundation (01-06)
- **01-stack** — Ethereum stack overview, RPC connectivity
- **02-rpc-basics** — Chain ID, network ID, headers
- **03-keys-addresses** — secp256k1, private keys, addresses
- **04-accounts-balances** — Account queries, balance checks
- **05-tx-nonces** — Transaction nonces, replay protection
- **06-eip1559** — EIP-1559 dynamic fees, priority fees

### Smart Contracts (07-09)
- **07-eth-call** — eth_call, read-only contract calls
- **08-abigen** — ABI encoding/decoding, typed bindings
- **09-events** — Event logs, decoding, filtering

### State & Storage (10-12)
- **10-filters** — Log filters, WebSocket subscriptions
- **11-storage** — Storage slots, state access
- **12-proofs** — Merkle-Patricia tries, cryptographic proofs

### Advanced Queries (13-17)
- **13-trace** — Transaction tracing, internal calls
- **14-explorer** — Block explorer queries
- **15-receipts** — Transaction receipts, gas used, logs
- **16-concurrency** — Concurrent RPC calls, rate limiting
- **17-indexer** — Event indexing, historical data

### Operations (18-25)
- **18-reorgs** — Chain reorganizations, handling reorgs
- **19-devnets** — Local devnets, Ganache, Hardhat
- **20-node** — Node management, peer discovery
- **21-sync** — Sync status, sync modes (full, snap, light)
- **22-peers** — Peer management, network topology
- **23-mempool** — Mempool monitoring, pending transactions
- **24-monitor** — Node monitoring, health checks
- **25-toolbox** — Utility functions, helper libraries

---

## 🚀 Quick Start

```bash
# Set RPC endpoint
export INFURA_RPC_URL=https://eth.llamarpc.com

# Start with project 01
cd geth/01-stack

# Implement functions in internal/stack/exercise.go
# Run tests
go test ./...

# Debug with fixed inputs
go run ./cmd/dev

# Run application with RPC URL
go run ./cmd/app https://eth.llamarpc.com
```

## 🔑 Key Concepts

- **go-ethereum patterns** — Learn Geth architecture and idioms
- **RPC client patterns** — Production-grade error handling
- **Cryptographic primitives** — secp256k1, Merkle tries, proofs
- **Real-world operations** — Indexing, monitoring, reorg handling

## 🎓 Learning Path

**Ethereum Basics:** 01-06 (RPC, keys, transactions)
**Smart Contracts:** 07-09 (eth_call, events, abigen)
**Advanced:** 10-17 (storage, proofs, tracing, indexing)
**Operations:** 18-25 (reorgs, monitoring, tooling)

## 📖 Prerequisites

- Basic understanding of Ethereum concepts (accounts, transactions, blocks)
- Familiarity with RPC APIs
- Access to an Ethereum RPC endpoint (Infura, Alchemy, or local node)

---

See the [root README](../README.md) for complete documentation.
