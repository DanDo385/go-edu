# 05-tx-nonces: Building and Sending Legacy Transactions

**Goal:** Build/send a legacy transaction, manage nonces, and understand replay protection.

## Big Picture: Transaction Lifecycle

Transactions are the fundamental unit of state change on Ethereum. Understanding how to build, sign, and send transactions is essential for any Ethereum developer.

**Computer Science principle:** Transactions are **immutable, ordered messages** that change blockchain state. They're like database transactions, but cryptographically signed and globally ordered.

### Transaction Components

A legacy transaction (pre-EIP-1559) contains:
- **Nonce:** Sequence number (prevents replay and ensures ordering)
- **Gas Price:** Price per unit of gas (in wei)
- **Gas Limit:** Maximum gas to consume
- **To:** Recipient address (nil for contract creation)
- **Value:** Amount of ETH to send (in wei)
- **Data:** Calldata (function calls, contract bytecode, etc.)
- **v, r, s:** Signature components (from ECDSA signing)

**Key insight:** Transactions are **signed messages**. The signature proves you own the private key, and the nonce ensures ordering.

## Learning Objectives

By the end of this module, you should be able to:

1. **Fetch pending nonce** for an address (includes pending transactions)
2. **Build a legacy transaction** with gasPrice
3. **Sign a transaction** with EIP-155 replay protection (chainID)
4. **Send a transaction** to the network
5. **Understand nonce ordering** and why gaps stall subsequent transactions

## Prerequisites

- **Modules 01-04:** RPC basics, keys/addresses, account balances
- **Basic familiarity:** ETH units (wei vs ETH), transaction concepts
- **Go basics:** Error handling, big integers, crypto operations

## Building on Previous Modules

### From Module 03 (03-keys-addresses)
- You learned to generate private keys and derive addresses
- Now you're using those keys to **sign transactions**
- The address you derive is the `from` address in transactions

### From Module 04 (04-accounts-balances)
- You learned to query account balances
- Now you're **changing** those balances by sending transactions
- The nonce you fetch is stored in the account state

### From Module 01 (01-stack)
- You learned about chainID (EIP-155 replay protection)
- Now you're using chainID in transaction signatures
- This prevents transactions from being replayed on other chains

### Connection to Solidity-edu
- **Functions & Payable:** Sending ETH and tracking balances
- **Errors & Reverts:** Transaction status and failure modes connect to receipts (module 15)
- **Gas & Storage:** Understanding gas limits and pricing

## Understanding Nonces: The Ordering Mechanism

### What is a Nonce?

**Nonce** = "number used once" - a sequence number for each address.

**Computer Science principle:** Nonces ensure **ordering** and **uniqueness**. They prevent:
1. **Replay attacks:** Can't reuse the same transaction
2. **Out-of-order execution:** Transactions must be processed sequentially
3. **Double-spending:** Can't send the same transaction twice

### Nonce Rules

- **Start at 0:** First transaction from an address uses nonce 0
- **Increment by 1:** Each subsequent transaction increments the nonce
- **No gaps:** If nonce 5 is missing, nonce 6+ will be stuck pending
- **No duplicates:** Can't reuse a nonce (transaction will fail)

**Fun fact:** Nonces are per-address, not global. Address A can use nonce 0 at the same time address B uses nonce 0.

### Pending vs Confirmed Nonces

- **PendingNonceAt():** Returns the next nonce including pending transactions
- **NonceAt(blockNumber):** Returns nonce at a specific block (confirmed only)

**Production tip:** Use `PendingNonceAt()` when sending transactions to avoid nonce conflicts with pending transactions.

## Real-World Analogies

### The Post Office Queue Analogy
- **Nonce:** Your ticket number
- **Queue:** Transactions waiting to be processed
- **Counter:** Block proposer/miner
- **Problem:** If ticket 5 is missing, tickets 6+ wait forever

### The Database Transaction Analogy
- **Nonce:** Transaction sequence number
- **Ordering:** Transactions must execute in order
- **Atomicity:** Each transaction is all-or-nothing

### The Checkbook Analogy
- **Nonce:** Check number
- **Ordering:** Checks must be cashed in order
- **Gaps:** Missing check numbers cause problems

## Fun Facts & Nerdy Details

### EIP-155: Replay Protection

**Before EIP-155:** Transactions could be replayed across chains (e.g., mainnet → testnet)

**After EIP-155:** ChainID is included in the signature:
- Signature includes `(r, s, v)` where `v` encodes chainID
- Transaction signed for chainID 1 (mainnet) can't be replayed on chainID 11155111 (Sepolia)

**Nerdy detail:** The `v` value is `recoveryID + chainID * 2 + 35`. This encodes both the recovery ID (for ECDSA) and the chain ID.

### Gas Price Mechanics

- **Gas Price:** Price per unit of gas (in wei)
- **Gas Limit:** Maximum gas to consume
- **Total Cost:** `gasPrice * gasUsed` (you pay for gas used, not limit)

**Fun fact:** Legacy transactions use a fixed gas price. EIP-1559 (module 06) introduces dynamic fees with base fee + tip.

### Transaction Signing Process

1. **Serialize transaction:** RLP-encode all fields except signature
2. **Hash:** Keccak256 hash of serialized data
3. **Sign:** ECDSA sign the hash with private key
4. **Encode:** Add signature to transaction (v, r, s)

**Computer Science principle:** Signing the hash (not the raw data) is more efficient and secure. The hash is fixed-size (32 bytes) regardless of transaction size.

## Comparisons

### Legacy vs EIP-1559 Transactions
| Aspect | Legacy (this module) | EIP-1559 (module 06) |
|--------|---------------------|---------------------|
| Gas pricing | Fixed `gasPrice` | Dynamic `baseFee + tip` |
| Fee structure | Single price | Max fee cap + priority tip |
| Efficiency | Less efficient | More efficient |
| Status | Still works | Recommended for production |

### PendingNonceAt vs NonceAt
| Method | Includes Pending | Use Case |
|--------|------------------|----------|
| `PendingNonceAt()` | ✅ Yes | Sending new transactions |
| `NonceAt(block)` | ❌ No | Historical queries |

### Go `ethclient` vs JavaScript `ethers.js`
- **Go:** `client.PendingNonceAt(ctx, addr)` → Returns `uint64`
- **JavaScript:** `provider.getTransactionCount(addr, "pending")` → Returns `BigNumber`
- **Same JSON-RPC:** Both call `eth_getTransactionCount` with `"pending"` block

## Related Solidity-edu Modules

- **02 Functions & Payable:** Sending ETH and tracking balances
- **05 Errors & Reverts:** Transaction status and failure modes connect to receipts (module 15)
- **06 Mappings, Arrays & Gas:** Understanding gas limits and pricing

## What You'll Build

In this module, you'll create a CLI that:
1. Takes recipient address, ETH amount, and private key as input
2. Fetches the pending nonce for the sender
3. Fetches the suggested gas price
4. Builds a legacy transaction
5. Signs the transaction with EIP-155 (chainID) protection
6. Sends the transaction to the network
7. Displays transaction hash and details

**Key learning:** You'll understand the complete transaction lifecycle from building to broadcasting!

## Code Structure & Patterns

### The Exercise File (`exercise/exercise.go`)

The exercise file contains TODO comments guiding you through the implementation. Each TODO represents a fundamental concept:

1.  **Input Validation** - Validate private key, amount, and gas limit.
2.  **Nonce Management** - Fetch the pending nonce for the sender's address.
3.  **Gas Price** - Fetch the suggested gas price from the network.
4.  **Transaction Creation** - Assemble a new legacy transaction.
5.  **Transaction Signing** - Sign the transaction with the private key and chain ID.
6.  **Transaction Sending** - Broadcast the signed transaction to the network.

### The Solution File (`exercise/solution.go`)

The solution file contains detailed educational comments explaining:
- **Why** each step is necessary (the reasoning behind the code).
- **How** concepts repeat and build on each other (pattern recognition).
- **What** fundamental principles are being demonstrated (transaction lifecycle, replay protection).

### Key Patterns You'll Learn

#### Pattern 1: Nonce Fetching
```go
nonce, err := client.PendingNonceAt(ctx, from)
if err != nil {
    return nil, fmt.Errorf("pending nonce: %w", err)
}
```
**Why:** Using the pending nonce is crucial to avoid conflicts with transactions that are in the mempool but not yet mined.

**Building on:** RPC calls from previous modules.

**Repeats in:** Any application that sends transactions.

#### Pattern 2: Transaction Creation
```go
tx := types.NewTransaction(nonce, cfg.To, cfg.AmountWei, cfg.GasLimit, gasPrice, cfg.Data)
```
**Why:** This function assembles all the necessary fields into a `types.Transaction` object.

**Building on:** Go's struct creation and the `go-ethereum` types package.

**Repeats in:** Any application that creates new transactions.

#### Pattern 3: Transaction Signing
```go
signer := types.LatestSignerForChainID(chainID)
signedTx, err := types.SignTx(tx, signer, cfg.PrivateKey)
if err != nil {
    return nil, fmt.Errorf("sign tx: %w", err)
}
```
**Why:** Signing a transaction proves that you own the private key and authorize the state change. The `signer` object incorporates the chain ID to prevent replay attacks (EIP-155).

**Building on:** Cryptographic principles from Module 03.

**Repeats in:** Any application that sends transactions on behalf of a user.

## Deep Dive: Nonce Management

The nonce is the most critical and often misunderstood part of sending transactions.

- **Strict Ordering:** Transactions from a single account are processed strictly in nonce order. If you submit nonce 5, it will not be processed until nonce 4 is confirmed.
- **Gaps are Problematic:** If you have a pending transaction with nonce 5, and you submit another with nonce 7, the second one will get stuck in the mempool until a transaction with nonce 6 is mined.
- **Replacing Transactions:** You can replace a pending transaction by sending a new one with the same nonce but a higher gas price.

## Error Handling: Building Robust Systems

Error handling in this module involves wrapping errors from the RPC client and the signing process.

```go
if err := client.SendTransaction(ctx, signedTx); err != nil {
    return nil, fmt.Errorf("send tx: %w", err)
}
```

If `SendTransaction` fails, the error will be wrapped with "send tx: ". This helps pinpoint the failure to the broadcasting step. Common errors include "nonce too low" or "insufficient funds".

## Testing Strategy

The test file (`exercise/exercise_test.go`) demonstrates several important patterns:

1.  **Mock implementations:** `mockTXClient` implements the `TXClient` interface, allowing us to test our logic without a real Ethereum node.
2.  **Configuration testing:** We test various configurations, such as providing a nonce manually vs. fetching it automatically.
3.  **Signature verification:** We can verify that the transaction was signed correctly.
4.  **"No Send" mode:** The `NoSend` flag allows us to test the transaction creation and signing logic without actually broadcasting the transaction.

## Files

- **Exercise:** `exercise/exercise.go` - Your starting point with TODO comments guiding implementation
- **Solution:** `exercise/solution.go` - Complete implementation with detailed educational comments explaining every concept
- **Types:** `exercise/types.go` - Interface and struct definitions
- **Tests:** `exercise/exercise_test.go` - Test suite demonstrating patterns and verifying correctness

## Next Steps

After completing this module, you'll move to **06-eip1559** where you'll:
- Build EIP-1559 dynamic fee transactions
- Understand base fee + priority tip mechanics
- Learn about max fee caps and refunds
- Use the modern transaction format (recommended for production)
