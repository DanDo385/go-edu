# Ethereum Go Development Track - Quick Start Guide

Learn Ethereum development with Go through 25 progressive hands-on projects, from basic RPC connections to production-grade tooling.

## Overview

This track contains **25 modules** that teach you how to build Ethereum applications using Go and the go-ethereum (geth) library. Each module is a standalone project with exercises, tests, and complete reference solutions.

**What you'll learn:**
- JSON-RPC client programming
- Cryptographic key management
- Transaction building and signing
- Smart contract interaction (ABI encoding, typed bindings)
- Event filtering and real-time monitoring
- Storage proofs and Merkle tries
- Production patterns (indexing, reorg handling, health monitoring)

## Prerequisites

### Required
- **Go 1.22+** ([install](https://go.dev/doc/install))
- **An Ethereum RPC endpoint** (choose one):
  - Public RPC: [Infura](https://infura.io), [Alchemy](https://alchemy.com), [QuickNode](https://quicknode.com)
  - Local node: Run your own Geth/Erigon node
  - Local devnet: Use [Anvil](https://book.getfoundry.sh/anvil/) or Hardhat

### Optional (for specific modules)
- **Anvil** (module 19 - devnets) - Install via Foundry
- **SQLite** (module 17 - indexer) - Usually pre-installed on macOS/Linux

## Quick Setup

### 1. Get an RPC Endpoint

**Option A: Use a Public RPC (Easiest)**

Sign up for a free account at [Infura](https://infura.io) or [Alchemy](https://alchemy.com) and get your API key.

```bash
# Set environment variable (Linux/macOS)
export INFURA_RPC_URL="https://mainnet.infura.io/v3/YOUR_PROJECT_ID"

# Or create a .env file in the geth/ directory
echo 'INFURA_RPC_URL=https://mainnet.infura.io/v3/YOUR_PROJECT_ID' > .env
```

**Option B: Run a Local Node (Advanced)**

```bash
# Requires ~1TB disk space and several hours to sync
geth --http --http.api eth,net,web3 --syncmode snap
```

### 2. Install Dependencies

```bash
# From the repository root
cd /path/to/go-edu
go mod tidy
```

### 3. Run Your First Module

```bash
# From repository root
make test P=geth/01-stack

# Or directly
cd geth/01-stack/exercise
INFURA_RPC_URL=... go test -v
```

## Project Structure

Each module follows the same structure:

```
geth/
├── 01-stack/
│   ├── README.md           # Module documentation (see READMESUMMARY.md)
│   └── exercise/
│       ├── exercise.go     # Your implementation (TODO stubs)
│       ├── solution.go     # Reference solution (build tags)
│       ├── types.go        # Interfaces and types
│       └── exercise_test.go # Test suite
├── 02-rpc-basics/
│   └── ... (same structure)
└── ... (modules 03-25)
```

## How to Use

### Option 1: Implement Exercises (Recommended)

Each module has TODO comments guiding you through implementation:

```bash
# 1. Navigate to a module
cd geth/01-stack/exercise

# 2. Read the module README
cat README.md

# 3. Implement the TODOs in exercise.go
# (exercise.go and solution.go use build tags to avoid conflicts)

# 4. Run tests against YOUR implementation
INFURA_RPC_URL=... go test -v

# 5. Compare with reference solution
INFURA_RPC_URL=... go test -tags=solution -v
```

### Option 2: Study Solutions First

If you prefer to learn by reading:

```bash
cd geth/01-stack/exercise

# Run tests with solution
INFURA_RPC_URL=... go test -tags=solution -v

# Read the solution code
cat solution.go
```

### Option 3: Use the Makefile (Convenient)

From the repository root:

```bash
# Test a specific module (auto-detects geth/)
make test P=geth/01-stack
make test P=01-stack  # Also works

# Test with solution
INFURA_RPC_URL=... go test -tags=solution ./geth/01-stack/...

# List all modules
make list-geth
```

## Module Roadmap

### Foundation (Modules 01-04)
Start here! Learn the basics of connecting to Ethereum and querying data.

| Module | Name | What You'll Learn |
|--------|------|-------------------|
| 01 | stack | Ethereum architecture, RPC connections, chain ID |
| 02 | rpc-basics | JSON-RPC methods, block structures, retry logic |
| 03 | keys-addresses | Key generation, address derivation, keystore files |
| 04 | accounts-balances | EOA vs contracts, balance queries, code detection |

### Transactions (Modules 05-06)
Build and send transactions to the network.

| Module | Name | What You'll Learn |
|--------|------|-------------------|
| 05 | tx-nonces | Legacy transactions, nonces, signing, broadcasting |
| 06 | eip1559 | EIP-1559 dynamic fees, base fee, priority fee |

### Smart Contracts (Modules 07-09)
Interact with deployed smart contracts.

| Module | Name | What You'll Learn |
|--------|------|-------------------|
| 07 | eth-call | Manual ABI encoding, read-only calls, function selectors |
| 08 | abigen | Typed contract bindings, code generation, CallOpts |
| 09 | events | Event filtering, log decoding, topics vs data |

### Real-Time Monitoring (Module 10)
Watch the blockchain in real-time.

| Module | Name | What You'll Learn |
|--------|------|-------------------|
| 10 | filters | WebSocket subscriptions, newHeads, polling fallback |

### Advanced Queries (Modules 11-13)
Deep dive into storage, proofs, and execution tracing.

| Module | Name | What You'll Learn |
|--------|------|-------------------|
| 11 | storage | Raw storage slots, mapping hashes, Solidity layout |
| 12 | proofs | Merkle-Patricia tries, eth_getProof, light clients |
| 13 | trace | debug_traceTransaction, call trees, gas analysis |

### Building Tools (Modules 14-16)
Combine concepts to build useful tools.

| Module | Name | What You'll Learn |
|--------|------|-------------------|
| 14 | explorer | Block/transaction explorer, data summarization |
| 15 | receipts | Transaction receipts, status codes, logs |
| 16 | concurrency | Worker pools, concurrent RPC calls, context timeouts |

### Production Patterns (Modules 17-19)
Learn patterns for production applications.

| Module | Name | What You'll Learn |
|--------|------|-------------------|
| 17 | indexer | ERC20 indexer, SQLite storage, event processing |
| 18 | reorgs | Chain reorganization detection, rescanning logic |
| 19 | devnets | Local development networks, Anvil, account funding |

### Node Operations (Modules 20-24)
Monitor and operate Ethereum nodes.

| Module | Name | What You'll Learn |
|--------|------|-------------------|
| 20 | node | Node info, client version, peer count |
| 21 | sync | Sync progress, sync modes (full/snap/light) |
| 22 | peers | P2P networking, gossip protocols, connectivity |
| 23 | mempool | Pending transactions, txpool visibility, MEV |
| 24 | monitor | Health checks, block lag detection, alerting |

### Capstone (Module 25)
Put it all together.

| Module | Name | What You'll Learn |
|--------|------|-------------------|
| 25 | toolbox | Swiss Army CLI combining all previous concepts |

## Running Tests

### Basic Test Commands

```bash
# Run tests for a module (YOUR implementation)
cd geth/01-stack/exercise
INFURA_RPC_URL=... go test

# Verbose output (see each test)
INFURA_RPC_URL=... go test -v

# Run specific test
INFURA_RPC_URL=... go test -v -run TestGetStatus

# Run solution tests
INFURA_RPC_URL=... go test -tags=solution -v
```

### Using Make Commands

```bash
# From repository root
make test P=geth/01-stack        # Test specific module
make test P=geth/01-stack -v     # Verbose
```

### Test Requirements

- **RPC endpoint required:** Most tests need `INFURA_RPC_URL` environment variable
- **Network access:** Tests make real RPC calls (not mocked)
- **Rate limits:** Public RPCs have rate limits; tests may fail if hit
- **Some tests need specific modules:** Module 19 needs Anvil installed

## Common Issues & Solutions

### "dial tcp: i/o timeout"
**Problem:** Can't connect to RPC endpoint
**Solution:** Check your `INFURA_RPC_URL` is set correctly and you have internet access

### "429 Too Many Requests"
**Problem:** Hit RPC rate limit
**Solution:** Wait a moment and retry, or upgrade your RPC plan

### "missing INFURA_RPC_URL"
**Problem:** Environment variable not set
**Solution:** `export INFURA_RPC_URL=https://mainnet.infura.io/v3/YOUR_KEY`

### "debug_traceTransaction not supported"
**Problem:** Public RPCs often disable debug methods
**Solution:** Use a different RPC provider or run your own node

### Tests pass with solution but fail with exercise
**Problem:** Your implementation has bugs
**Solution:** Compare your `exercise.go` with `solution.go` and read the comments

## Build Tags Explained

Go build tags let us have both `exercise.go` and `solution.go` in the same directory without conflicts:

- **Default** (no tags): Compiles `exercise.go` (your implementation)
  ```bash
  go test
  ```

- **With -tags=solution**: Compiles `solution.go` (reference solution)
  ```bash
  go test -tags=solution
  ```

**How it works:**
- `exercise.go` has: `//go:build !solution`
- `solution.go` has: `//go:build solution`

## Learning Paths

### Path 1: Beginner (Never used Go-Ethereum)
Follow modules sequentially 01→25. Each builds on previous concepts.

**Estimated time:** 20-30 hours total

### Path 2: Experienced (Know Go, new to Ethereum)
- Skim modules 01-04 (basics)
- Focus on 05-09 (transactions and contracts)
- Deep dive 11-13 (storage, proofs, tracing)
- Practice 16-25 (production patterns)

**Estimated time:** 10-15 hours

### Path 3: Ethereum Expert (New to go-ethereum)
- Quick read 01-02 (ethclient API)
- Jump to areas of interest (e.g., 17 for indexing, 13 for tracing)
- Use as reference when building your own tools

**Estimated time:** 5-10 hours

## Additional Resources

### Documentation
- **Module Details:** See [READMESUMMARY.md](./READMESUMMARY.md) for comprehensive documentation of all 25 modules
- **go-ethereum Docs:** https://geth.ethereum.org/docs/
- **JSON-RPC Spec:** https://ethereum.org/en/developers/docs/apis/json-rpc/
- **Go by Example:** https://gobyexample.com/

### Tools
- **Infura:** https://infura.io (free tier available)
- **Alchemy:** https://alchemy.com (free tier available)
- **Foundry (Anvil):** https://book.getfoundry.sh/anvil/
- **Etherscan:** https://etherscan.io (for verifying data)

### Related Tracks
- **minis/**: General Go fundamentals (strings, concurrency, HTTP, etc.)
- See main [README.md](../README.md) for the complete learning path

## Contributing

Found a bug? Have a question? Want to add a module?

1. Check existing issues at the repository
2. Open a new issue with details
3. Submit a PR following existing code style

## What's Next?

1. **Set up your RPC endpoint** (see Quick Setup above)
2. **Start with module 01-stack:** `cd geth/01-stack && cat README.md`
3. **Implement the TODOs** in `exercise.go`
4. **Run tests:** `INFURA_RPC_URL=... go test -v`
5. **Move to module 02** and repeat!

**Ready to start?** Run this command:

```bash
# From repository root
make test P=geth/01-stack
```

Good luck! 🚀
