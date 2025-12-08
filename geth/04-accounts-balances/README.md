# 04-accounts-balances: Understanding Account Types and State

**Goal:** Classify EOAs vs contracts and query balances, understanding the fundamental difference between account types on Ethereum.

## Big Picture: Two Types of Accounts

Ethereum accounts come in two flavors: **EOAs (Externally Owned Accounts)** and **Contracts**. Understanding this distinction is fundamental to Ethereum development.

**Computer Science principle:** This is a **type system** at the blockchain level. Just like programming languages have types (int, string, struct), Ethereum has account types. The type determines what operations are possible.

### EOA (Externally Owned Account)
- **Has:** Private key, address, balance, nonce
- **Does NOT have:** Code (bytecode)
- **Can:** Send transactions, sign messages
- **Cannot:** Execute arbitrary code
- **Analogy:** An empty plot of land with a mailbox (address) and a safe (balance)

### Contract Account
- **Has:** Address, balance, nonce, **code** (bytecode)
- **Does NOT have:** Private key (cannot initiate transactions directly)
- **Can:** Execute code when called, store state, emit events
- **Cannot:** Sign transactions (needs an EOA to call it)
- **Analogy:** A building with machinery (code) on the same street (address)

**Key insight:** Contracts are **stateless** in the sense that they don't have private keys, but they have **state** (storage) that can be modified through function calls.

## Learning Objectives

By the end of this module, you should be able to:

1. **Fetch balances** at a specific block (or latest) in wei
2. **Detect account type** by checking for code presence
3. **Understand special cases:**
   - Precompiles (addresses 0x01-0x09) have code but are special-purpose
   - Selfdestructed contracts have nonce > 0 but code size 0
4. **Distinguish between** `eth_getBalance` and `eth_getCode` use cases

## Prerequisites

- **Module 03 (03-keys-addresses):** You should understand how addresses are derived from keys
- **Module 01-02:** Basic RPC familiarity, understanding of blocks and state
- **Go basics:** Loops, error handling, big integers

## Building on Previous Modules

### From Module 03 (03-keys-addresses)
- You learned to generate addresses from private keys
- Now you're querying those addresses on the blockchain to see their state
- The addresses you generated are **EOAs** (they have no code until you deploy a contract)

### From Module 01-02 (01-stack, 02-rpc-basics)
- You learned to connect to Ethereum nodes and query blocks
- Now you're querying **account state** within those blocks
- Balances are part of the state trie (referenced by `stateRoot` in block headers)

### Connection to Solidity-edu
- **Functions & Payable:** `msg.sender` can be either an EOA or a contract (contracts can call other contracts)
- **Access Control / Ownable:** Ownership checks typically assume EOAs or smart wallets (contracts that act like EOAs)
- **Modifiers & Restrictions:** Access control relies on understanding who `msg.sender` is (EOA vs contract)

## Understanding Account State

### The Account Trie

Accounts are stored in a **Merkle-Patricia Trie** (MPT), indexed by address. Each account has:

```go
type Account struct {
    Nonce    uint64      // Transaction count (prevents replay)
    Balance  *big.Int    // Balance in wei
    Root     common.Hash // Storage root (for contracts)
    CodeHash common.Hash // Hash of contract bytecode
}
```

**Computer Science principle:** The account trie is a **key-value store** where:
- **Key:** Account address (20 bytes)
- **Value:** Account struct (nonce, balance, storage root, code hash)

**Fun fact:** The `stateRoot` in block headers (from module 01) is the Merkle root of this entire account trie. Change any account, and the root changes!

### Balance Storage

Balances are stored in **wei** (the smallest unit of ETH):
- 1 ETH = 10^18 wei
- 1 gwei = 10^9 wei

**Why wei?** Using the smallest unit avoids floating-point precision issues. All calculations use integers (big.Int in Go).

**Nerdy detail:** Wei is named after Wei Dai, creator of b-money (an early cryptocurrency proposal). Ethereum uses wei to honor this contribution to the field.

### Code Storage

Contract bytecode is stored separately from the account trie:
- **Account trie:** Stores `codeHash` (hash of bytecode)
- **Code storage:** Actual bytecode stored in a separate database

**Why separate?** Code is large (can be 24KB max). Storing hashes in the trie keeps it small, while code is stored separately and retrieved when needed.

## Real-World Analogies

### The Real Estate Analogy
- **EOA:** Empty plot of land with a mailbox (address) and a safe (balance)
- **Contract:** Building with machinery (code) on the same street (address)
- **Balance:** Money in the safe (can be in either type)
- **Code:** The machinery/blueprint (only contracts have this)

### The Bank Account Analogy
- **EOA:** Personal bank account (you control it with your private key)
- **Contract:** Trust account or automated system (runs code, but needs you to initiate)
- **Balance:** Money in the account (same for both types)
- **Code:** The rules/automation (only contracts have this)

### The Computer Process Analogy
- **EOA:** A user account (can run commands, but no persistent code)
- **Contract:** A daemon/service (has code that runs, but needs to be invoked)
- **Balance:** Resources allocated to the process
- **Code:** The executable program (only contracts have this)

## Fun Facts & Nerdy Details

### Precompiled Contracts

Addresses `0x01` through `0x09` are **precompiled contracts**:
- They have code (special bytecode)
- They're implemented natively in the EVM (not stored as bytecode)
- Examples:
  - `0x01`: ECRecover (signature verification)
  - `0x02`: SHA256 hash
  - `0x03`: RIPEMD160 hash
  - `0x04`: Identity (data copy)
  - `0x05`: ModExp (modular exponentiation)
  - `0x06-0x08`: Elliptic curve operations
  - `0x09`: Blake2 compression

**Why precompiles?** These operations are expensive to compute in EVM bytecode, so they're implemented natively for efficiency.

### Selfdestructed Contracts

Contracts that call `selfdestruct()` have a special state:
- **Nonce:** Still > 0 (shows it was deployed)
- **Code:** Empty (0 bytes) - the code was deleted
- **Balance:** Can be > 0 (sent to beneficiary on selfdestruct)

**Heuristic warning:** Don't assume `code length == 0` always means EOA. Check nonce too!

### Account Creation

**EOA creation:**
- Happens automatically when you generate a keypair (module 03)
- Account appears on-chain when it receives its first transaction
- No deployment cost

**Contract creation:**
- Requires a transaction with bytecode in the `data` field
- Costs gas (deployment is expensive!)
- Account appears when the deployment transaction is mined

## Comparisons

### `eth_getBalance` vs `eth_getCode`
| Method | Purpose | Returns | Cost |
|--------|---------|---------|------|
| `eth_getBalance` | Get account balance | `big.Int` (wei) | Cheap |
| `eth_getCode` | Get contract bytecode | `[]byte` | Moderate (can be large) |

### EOA vs Contract
| Aspect | EOA | Contract |
|--------|-----|----------|
| Private key | ✅ Yes | ❌ No |
| Code | ❌ No | ✅ Yes |
| Can initiate tx | ✅ Yes | ❌ No (needs EOA caller) |
| Can be called | ❌ No | ✅ Yes |
| Storage | ❌ No | ✅ Yes (if contract has storage) |

### Go `ethclient` vs JavaScript `ethers.js`
- **Go:** `client.BalanceAt(ctx, addr, nil)` → Returns `*big.Int`
- **JavaScript:** `provider.getBalance(addr)` → Returns `BigNumber`
- **Same JSON-RPC:** Both call `eth_getBalance` under the hood

## Related Solidity-edu Modules

- **02 Functions & Payable:** Balances change when contracts receive ETH via `payable` functions. This module shows you how to query those balances!
- **04 Modifiers & Access Control:** Ownership checks use `msg.sender == owner`. Understanding EOA vs contract helps you understand who `msg.sender` can be.
- **08 ERC20 from Scratch:** ERC20 tokens maintain balances in contract storage. This module shows you how to query account balances (both ETH and token balances).

## What You'll Build

In this module, you'll create a CLI that:
1. Takes one or more addresses as command-line arguments
2. Queries the balance of each address (in wei)
3. Queries the code of each address
4. Classifies each address as EOA or Contract
5. Displays address, type, and balance

**Key learning:** You'll understand the fundamental distinction between EOAs and contracts, and how to query account state on the blockchain!

## Code Structure & Patterns

### The Exercise File (`exercise/exercise.go`)

The exercise file contains TODO comments guiding you through the implementation. Each TODO represents a fundamental concept:

1.  **Input Validation** - Reinforce the pattern of validating inputs.
2.  **Iterating over Inputs** - Learn how to process a list of inputs.
3.  **Querying Account State** - Use `BalanceAt` and `CodeAt` to inspect accounts.
4.  **Defensive Copying** - Practice defensive copying for `big.Int` and `[]byte`.
5.  **Classifying Accounts** - Use the retrieved state to determine the account type.

### The Solution File (`exercise/solution.go`)

The solution file contains detailed educational comments explaining:
- **Why** each step is necessary (the reasoning behind the code).
- **How** concepts repeat and build on each other (pattern recognition).
- **What** fundamental principles are being demonstrated (account types, state).

### Key Patterns You'll Learn

#### Pattern 1: Looping and Querying
```go
for _, addr := range cfg.Addresses {
    bal, err := client.BalanceAt(ctx, addr, cfg.BlockNumber)
    // ...
    code, err := client.CodeAt(ctx, addr, cfg.BlockNumber)
    // ...
}
```
**Why:** This is a common pattern for batch-processing a list of items. For each address, we perform a series of RPC calls to gather the necessary information.

**Building on:** Go's `for...range` loop and the RPC call pattern from previous modules.

**Repeats in:** Any application that needs to process multiple accounts, transactions, or blocks.

#### Pattern 2: Defensive Copying of Slices and `big.Int`
```go
if bal != nil {
    bal = new(big.Int).Set(bal)
}
// ...
codeCopy := append([]byte(nil), code...)
```
**Why:** `BalanceAt` returns a `*big.Int` and `CodeAt` returns a `[]byte`. Both are reference types. If we don't copy them, the caller of our function could accidentally mutate the RPC client's internal data.

**Building on:** The defensive copying pattern from Modules 01 and 02.

**Repeats in:** Every function that returns data from external libraries or shared state.

#### Pattern 3: Account Type Classification
```go
accType := AccountTypeEOA
if len(codeCopy) > 0 {
    accType = AccountTypeContract
}
```
**Why:** The defining characteristic of a contract account is that it has code. If the code length is greater than zero, it's a contract.

**Building on:** The understanding of EOA vs. contract accounts.

**Repeats in:** Any application that needs to differentiate between account types (e.g., block explorers, wallets).

## Deep Dive: Defensive Copying of Slices

In Go, a slice is a descriptor for a contiguous segment of an underlying array. It consists of a pointer to the array, the length of the segment, and its capacity.

When you pass a slice to a function, you are passing a copy of the slice header, but not the underlying array. This means that if the function modifies the elements of the slice, the changes will be visible to the caller.

### The Bug (Without Defensive Copying)
```go
// BAD: Returning a slice that points to the client's internal buffer
return &Result{Code: code}
```

If the caller modifies `result.Code`, they will be modifying the RPC client's internal buffer, which can lead to unpredictable behavior and data races.

### The Fix (With Defensive Copying)
```go
// GOOD: Creating a new slice and copying the data
codeCopy := append([]byte(nil), code...)
return &Result{Code: codeCopy}
```
`append([]byte(nil), code...)` is an idiomatic way to create a copy of a slice. It creates a new slice and copies the elements from the old slice to the new one.

## Error Handling: Building Robust Systems

Error handling in this module is straightforward. We wrap errors from the RPC client with additional context to make them easier to debug.

```go
bal, err := client.BalanceAt(ctx, addr, cfg.BlockNumber)
if err != nil {
    return nil, fmt.Errorf("balance %s: %w", addr.Hex(), err)
}
```

If `BalanceAt` fails, the error message will include the address that caused the failure, making it easy to identify the source of the problem.

## Testing Strategy

The test file (`exercise/exercise_test.go`) demonstrates several important patterns:

1.  **Mock implementations:** `mockAccountClient` implements the `AccountClient` interface, allowing us to test our logic without a real Ethereum node.
2.  **Table-driven tests:** We can test multiple scenarios with different account states.
3.  **Defensive copy verification:** The tests ensure that the returned `AccountState` contains copies of the balance and code, not pointers to the mock's internal data.

## Files

- **Exercise:** `exercise/exercise.go` - Your starting point with TODO comments guiding implementation
- **Solution:** `exercise/solution.go` - Complete implementation with detailed educational comments explaining every concept
- **Types:** `exercise/types.go` - Interface and struct definitions
- **Tests:** `exercise/exercise_test.go` - Test suite demonstrating patterns and verifying correctness

## Next Steps

After completing this module, you'll move to **05-tx-nonces** where you'll:
- Build transactions from scratch
- Understand transaction nonces (sequence numbers)
- Sign transactions with private keys
- Broadcast transactions to the network
- Connect the addresses from module 03 to actual on-chain transactions
