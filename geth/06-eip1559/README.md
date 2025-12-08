# 06-eip1559: Dynamic Fee Transactions (EIP-1559)

**Goal:** Build/send an EIP-1559 (dynamic fee) transaction and understand base fee + tip math.

## Big Picture: The London Upgrade and Dynamic Fees

Post-London (August 2021), Ethereum uses **dynamic fees** instead of fixed gas prices. This makes fee estimation more predictable and reduces fee volatility.

**Computer Science principle:** EIP-1559 introduces a **two-part fee structure**:
- **Base Fee:** Algorithmically determined, **burned** (removed from supply)
- **Priority Fee (Tip):** Paid to validators/miners, incentivizes inclusion

This is more efficient than the legacy auction model where users bid against each other.

## Learning Objectives

By the end of this module, you should be able to:

1. **Construct a DynamicFeeTx** with maxFeePerGas and maxPriorityFeePerGas
2. **Convert user inputs** (gwei) to wei safely
3. **Sign with London signer** (chainID-aware) and broadcast
4. **Understand fee math:** effectiveGasPrice = min(maxFeeCap, baseFee + tip)
5. **Understand refunds:** Excess fees are refunded to the sender

## Prerequisites

- **Module 05 (05-tx-nonces):** Legacy transaction flow, nonces, signing
- **Comfort with:** ETH units, RPC basics, transaction concepts
- **Go basics:** Big integers, error handling

## Building on Previous Modules

### From Module 05 (05-tx-nonces)
- You learned to build legacy transactions with fixed `gasPrice`
- Now you're building **EIP-1559 transactions** with dynamic fees
- Same nonce management, same signing process, different fee structure

### From Module 01 (01-stack)
- You learned about chainID
- EIP-1559 transactions also use chainID for replay protection
- London signer includes chainID in the signature

### Connection to Solidity-edu
- **Functions & Payable:** Gas budgeting affects send/withdraw patterns
- **Gas & Storage lessons:** Understanding gas pricing is critical for optimizing costs

## Understanding EIP-1559 Fee Structure

### The Two-Part Fee System

**Base Fee:**
- **Set by protocol:** Algorithmically determined based on block fullness
- **Burned:** Removed from ETH supply (deflationary mechanism)
- **Variable:** Adjusts up/down by 12.5% per block based on target (50% full blocks)
- **Predictable:** Users can estimate base fee for next block

**Priority Fee (Tip):**
- **Set by user:** How much you're willing to pay for faster inclusion
- **Paid to validators:** Incentivizes block inclusion
- **Optional:** Can be 0 (transaction will still be included, just slower)

### Fee Caps and Refunds

**maxFeePerGas:** Maximum you're willing to pay total (baseFee + tip)
**maxPriorityFeePerGas:** Maximum tip you're willing to pay

**Effective gas price calculation:**
```
effectiveGasPrice = min(maxFeePerGas, baseFee + maxPriorityFeePerGas)
```

**Refund calculation:**
```
refund = (maxFeePerGas - effectiveGasPrice) * gasUsed
```

**Key insight:** You set caps to protect yourself from fee spikes. If baseFee rises unexpectedly, you're protected by maxFeePerGas.

## Real-World Analogies

### The Bus Fare Analogy
- **Base Fee:** Standard bus fare (set by transit authority, everyone pays)
- **Priority Tip:** Tip for priority boarding (optional, goes to driver)
- **Max Fee:** Your total budget cap (base fare + max tip you're willing to pay)
- **Refund:** If base fare is lower than expected, you get change back

### The Restaurant Analogy
- **Base Fee:** Menu price (set by restaurant)
- **Priority Tip:** Gratuity (optional, goes to server)
- **Max Fee:** Your total budget (menu price + max tip)
- **Refund:** If bill is less than budget, you get change

### The Auction Analogy (Legacy vs EIP-1559)
- **Legacy:** Everyone bids against each other (volatile, unpredictable)
- **EIP-1559:** Base price + optional tip (predictable, efficient)

## Fun Facts & Nerdy Details

### Base Fee Algorithm

The base fee adjusts based on block fullness:
- **Target:** 50% block fullness (15M gas out of 30M limit)
- **If block > 50% full:** Base fee increases by 12.5%
- **If block < 50% full:** Base fee decreases by 12.5%

**Mathematical formula:**
```
baseFee_new = baseFee_old * (1 + (gasUsed - targetGas) / targetGas / 8)
```

**Fun fact:** The 12.5% adjustment rate was chosen to balance responsiveness with stability. Too fast = volatile, too slow = unresponsive.

### Fee Burning: Deflationary Mechanism

**Before EIP-1559:** All fees went to miners/validators (inflationary)

**After EIP-1559:** Base fee is burned (removed from supply)

**Impact:** During high network activity, more ETH is burned than issued, making ETH deflationary!

**Nerdy detail:** As of 2024, Ethereum has burned over 4 million ETH (worth billions of dollars) through base fee burning.

### Gas Price Units

- **wei:** Smallest unit (1 ETH = 10^18 wei)
- **gwei:** Common unit for gas prices (1 gwei = 10^9 wei)
- **ETH:** Rarely used for gas prices (too large)

**Why gwei?** It's a convenient middle ground. Gas prices are typically 10-100 gwei, which is easier to work with than trillions of wei.

## Comparisons

### Legacy vs EIP-1559 Transactions
| Aspect | Legacy (module 05) | EIP-1559 (this module) |
|--------|-------------------|----------------------|
| Fee structure | Single `gasPrice` | `baseFee + tip` |
| Predictability | Low (auction model) | High (algorithmic base fee) |
| Refunds | No | Yes (excess refunded) |
| Status | Still works | Recommended for production |
| Signer | EIP155Signer | LondonSigner |

### maxFeePerGas vs maxPriorityFeePerGas
| Field | Purpose | Who Sets It |
|-------|---------|-------------|
| `maxFeePerGas` | Total budget cap | User |
| `maxPriorityFeePerGas` | Tip cap | User |
| `baseFee` | Base fee | Protocol (algorithmic) |

### Go `ethclient` vs JavaScript `ethers.js`
- **Go:** `types.NewTx(&types.DynamicFeeTx{...})` → Build transaction struct
- **JavaScript:** `populateTransaction()` + `signer.sendTransaction()` → Auto-populates EIP-1559
- **Same JSON-RPC:** Both use `eth_sendRawTransaction` under the hood

## Related Solidity-edu Modules

- **02 Functions & Payable:** Gas budgeting affects send/withdraw patterns
- **06 Mappings, Arrays & Gas:** Understanding gas pricing is critical for optimizing costs
- **07 Reentrancy & Security:** Gas costs affect security patterns (gas griefing attacks)

## What You'll Build

In this module, you'll create a CLI that:
1. Takes recipient address, ETH amount, private key, and fee caps as input
2. Fetches the pending nonce (same as module 05)
3. Converts gwei to wei for fee caps
4. Builds an EIP-1559 DynamicFeeTx
5. Signs with London signer (includes chainID)
6. Sends the transaction to the network
7. Displays transaction hash and fee details

**Key learning:** You'll understand the modern transaction format and dynamic fee mechanics!

## Files

- **Exercise:** `exercise/exercise.go` - Your starting point with TODO comments guiding implementation
- **Solution:** `exercise/solution.go` - Complete implementation with detailed educational comments explaining every concept
- **Types:** `exercise/types.go` - Interface and struct definitions
- **Tests:** `exercise/exercise_test.go` - Test suite demonstrating patterns and verifying correctness

## How to Run Tests

To run the tests for this module:

```bash
# From the project root (go-edu/)
cd geth/06-eip1559
go test ./exercise/

# Run with verbose output to see test details
go test -v ./exercise/

# Run solution tests (build with solution tag)
go test -tags solution -v ./exercise/

# Run specific test
go test -v ./exercise/ -run TestRun
```

## Code Structure & Patterns

### The Exercise File (`exercise/exercise.go`)

The exercise file contains TODO comments guiding you through the implementation. Each TODO represents a fundamental concept:

1. **Input Validation** - Learn defensive programming patterns
2. **Address Derivation** - Understand cryptographic address derivation from keys
3. **Nonce Management** - Learn automatic vs manual nonce handling
4. **Fee Estimation** - Understand EIP-1559 fee structure (base fee + tip)
5. **Transaction Construction** - Build DynamicFeeTx with all required fields
6. **Transaction Signing** - Cryptographically sign the transaction
7. **Transaction Broadcasting** - Send to network (with NoSend option for testing)

### The Solution File (`exercise/solution.go`)

The solution file contains detailed educational comments explaining:
- **Why** each step is necessary (the reasoning behind the code)
- **How** concepts repeat and build on each other (pattern recognition)
- **What** fundamental principles are being demonstrated (computer science concepts)

### Key Patterns You'll Learn

#### Pattern 1: Defensive Copying for Mutable Types
```go
// BAD: Shares pointer, mutations affect original
baseFee := header.BaseFee

// GOOD: Creates independent copy
baseFee := new(big.Int).Set(header.BaseFee)
```

**Why:** `big.Int` is mutable. If we return pointers to internal data, callers could mutate it, affecting other users of the client.

**Building on:** Module 01-stack taught defensive copying basics. Here we apply it to transaction fee data.

**Repeats in:** Every module that works with `big.Int` values (all transaction-related modules).

#### Pattern 2: Config-Based Defaults
```go
if cfg.GasLimit == 0 {
    cfg.GasLimit = defaultDynamicGasLimit
}
```

**Why:** Provides sensible defaults while allowing advanced users to override. Makes the API ergonomic for common cases.

**Building on:** Module 01-stack used Config for block numbers. Here we extend it for transaction parameters.

**Repeats in:** Every module with configurable behavior.

#### Pattern 3: Error Wrapping with Context
```go
if err != nil {
    return nil, fmt.Errorf("pending nonce: %w", err)
}
```

**Why:** Adds context to errors while preserving the error chain. Makes debugging easier.

**Building on:** Module 01-stack introduced this pattern. It's consistent across all modules.

**Repeats in:** Every error return in production Go code.

#### Pattern 4: RPC Call → Error Check → Nil Check → Defensive Copy
```go
header, err := client.HeaderByNumber(ctx, cfg.BlockNumber)
if err != nil {
    return nil, fmt.Errorf("header by number: %w", err)
}
if header == nil || header.BaseFee == nil {
    return nil, errors.New("base fee unavailable")
}
baseFee := new(big.Int).Set(header.BaseFee)
```

**Why:** RPC calls can fail or return nil. Always validate before using. Always copy mutable return values.

**Building on:** Module 01-stack taught RPC call patterns. Here we add BaseFee validation.

**Repeats in:** Every RPC call throughout the entire course.

## Deep Dive: EIP-1559 Fee Mechanics

### Fee Calculation Examples

**Example 1: Normal Case**
```
baseFee = 50 gwei
tipCap = 2 gwei
maxFee = 150 gwei

effectiveGasPrice = min(maxFee, baseFee + tipCap)
                  = min(150, 50 + 2)
                  = 52 gwei

You pay: 52 gwei per gas
Burned: 50 gwei per gas
Validator gets: 2 gwei per gas
```

**Example 2: Base Fee Spike**
```
baseFee = 100 gwei (spiked!)
tipCap = 2 gwei
maxFee = 150 gwei

effectiveGasPrice = min(150, 100 + 2)
                  = 102 gwei

You pay: 102 gwei per gas
Burned: 100 gwei per gas
Validator gets: 2 gwei per gas
```

**Example 3: Max Fee Too Low**
```
baseFee = 160 gwei (very high!)
tipCap = 2 gwei
maxFee = 150 gwei

effectiveGasPrice = min(150, 160 + 2)
                  = 150 gwei

Problem: Your maxFee (150) < baseFee (160)
Result: Transaction won't be included until baseFee drops below 150
```

### The "2x Base Fee" Rule Explained

**Why multiply by 2?**

EIP-1559 allows base fee to increase by 12.5% per block. Over 6 blocks:
```
Block 0: 50 gwei
Block 1: 50 * 1.125 = 56.25 gwei
Block 2: 56.25 * 1.125 = 63.28 gwei
Block 3: 63.28 * 1.125 = 71.19 gwei
Block 4: 71.19 * 1.125 = 80.09 gwei
Block 5: 80.09 * 1.125 = 90.10 gwei
Block 6: 90.10 * 1.125 = 101.36 gwei
```

After 6 blocks, base fee ~doubled. Setting `maxFee = 2 * baseFee + tip` gives you a buffer for ~6 blocks of worst-case increases.

**If your transaction isn't included in 6 blocks:**
- Your tip is probably too low (validators prefer higher tips)
- Or network congestion is extreme (rare)

## Error Handling: Building Robust Systems

### Common Transaction Errors

**1. "nonce too low"**
```
Cause: You already sent a transaction with this nonce (or higher)
Solution: Use PendingNonceAt to get the next available nonce
Prevention: Always query for nonce, don't hard-code
```

**2. "insufficient funds for gas * price + value"**
```
Cause: Account balance < (gasLimit * maxFee) + value
Solution: Check balance before sending, or reduce value/fees
Prevention: Always validate sufficient balance
```

**3. "max fee per gas less than block base fee"**
```
Cause: Your maxFee is below the current baseFee
Solution: Increase maxFee or wait for baseFee to drop
Prevention: Use suggested fees or the "2x base fee" rule
```

**4. "replacement transaction underpriced"**
```
Cause: Trying to replace a pending tx with same nonce but lower fees
Solution: Increase fees by at least 10% to replace
Prevention: Check mempool before replacing transactions
```

### Error Wrapping Strategy

```go
// Layer 1: RPC error
err := client.SendTransaction(ctx, tx)
// Error: "insufficient funds"

// Layer 2: Add context
return fmt.Errorf("send tx: %w", err)
// Error: "send tx: insufficient funds"

// Layer 3: Caller adds more context
return fmt.Errorf("failed to transfer ETH: %w", err)
// Error: "failed to transfer ETH: send tx: insufficient funds"
```

This creates a traceable error chain that shows exactly where and why the failure occurred.

## Testing Strategy

The test file (`exercise_test.go`) demonstrates several important patterns:

1. **Mock implementations:** `mockFeeClient` implements `FeeClient` interface
2. **Table-driven tests:** Multiple test cases with different scenarios
3. **Defensive copy verification:** Tests ensure immutability
4. **Error case testing:** Tests verify error handling works correctly
5. **Fee math testing:** Tests verify "2x base fee + tip" calculation

**Key insight:** Because we use interfaces, we can test our logic without needing a real Ethereum node. This makes tests fast, reliable, and deterministic.

**Example test case:**
```go
{
    name: "automatic fee calculation",
    setup: func(m *mockFeeClient) {
        m.baseFee = big.NewInt(50_000_000_000) // 50 gwei
        m.tipCap = big.NewInt(2_000_000_000)   // 2 gwei
    },
    validate: func(t *testing.T, res *Result) {
        // Verify maxFee = 2 * baseFee + tipCap
        want := new(big.Int).Mul(big.NewInt(50_000_000_000), big.NewInt(2))
        want.Add(want, big.NewInt(2_000_000_000))

        // Extract maxFee from transaction
        dynamicTx := res.Tx.Type() == types.DynamicFeeTxType
        got := res.Tx.GasFeeCap()

        if got.Cmp(want) != 0 {
            t.Errorf("maxFee = %s, want %s", got, want)
        }
    },
}
```

## Common Pitfalls & How to Avoid Them

### Pitfall 1: Not Copying big.Int Values
```go
// BAD: Shares pointer, mutations affect original
result.BaseFee = header.BaseFee

// GOOD: Creates independent copy
result.BaseFee = new(big.Int).Set(header.BaseFee)
```

**Why it's a problem:** `big.Int` is mutable. Sharing pointers causes data races and unexpected mutations.

**Fix:** Always use `new(big.Int).Set()` to create copies.

### Pitfall 2: Setting maxFee Too Low
```go
// BAD: Might be below baseFee if it spikes
maxFee := baseFee

// GOOD: Provides buffer for base fee increases
maxFee := new(big.Int).Mul(baseFee, big.NewInt(2))
maxFee.Add(maxFee, tipCap)
```

**Why it's a problem:** Base fee can increase 12.5% per block. If it spikes, your transaction won't be included.

**Fix:** Use the "2x base fee + tip" rule of thumb.

### Pitfall 3: Using NonceAt Instead of PendingNonceAt
```go
// BAD: Doesn't account for pending transactions
nonce, _ := client.NonceAt(ctx, from, nil)

// GOOD: Includes pending transactions in mempool
nonce, _ := client.PendingNonceAt(ctx, from)
```

**Why it's a problem:** If you have pending transactions, `NonceAt` returns an already-used nonce, causing "nonce too low" errors.

**Fix:** Always use `PendingNonceAt` for transaction creation.

### Pitfall 4: Not Validating BaseFee Exists
```go
// BAD: Panics on pre-London blocks
baseFee := header.BaseFee // Could be nil!

// GOOD: Validates before using
if header == nil || header.BaseFee == nil {
    return errors.New("base fee unavailable")
}
baseFee := new(big.Int).Set(header.BaseFee)
```

**Why it's a problem:** Pre-London blocks don't have BaseFee. Accessing it causes nil pointer panics.

**Fix:** Always validate `header.BaseFee != nil` before using.

### Pitfall 5: Forgetting to Copy Transaction Data
```go
// BAD: Shares slice backing array
txData := &types.DynamicFeeTx{
    Data: cfg.Data, // Caller could mutate this!
}

// GOOD: Creates independent copy
dataCopy := append([]byte(nil), cfg.Data...)
txData := &types.DynamicFeeTx{
    Data: dataCopy,
}
```

**Why it's a problem:** Slices are references. If caller mutates `cfg.Data` after calling us, it affects our transaction.

**Fix:** Always copy byte slices with `append([]byte(nil), slice...)`.

## How Concepts Build on Each Other

This module builds on patterns from previous modules while introducing new concepts:

1. **From Module 01-stack:**
   - Context validation → Same pattern here
   - RPC call pattern → Extended for fees and nonces
   - Defensive copying → Applied to fee data
   - Error wrapping → Consistent usage

2. **New in this module:**
   - Cryptographic signing (private keys, signatures)
   - Nonce management (PendingNonceAt vs NonceAt)
   - EIP-1559 fee structure (base fee, tip, max fee)
   - Transaction construction (DynamicFeeTx fields)
   - Address derivation (key → address)

3. **Patterns that repeat throughout the course:**
   - Input validation → Every function
   - Defensive copying → All mutable types
   - Error wrapping → All error returns
   - Config-based defaults → All configurable functions
   - RPC call pattern → All network operations

**The progression:**
- Module 01: Read data (headers, chain ID)
- Module 06: Write data (transactions)
- Future modules: More complex writes (contract calls, state changes)

Each module layers new concepts on top of existing patterns, building your understanding incrementally.

## Next Steps

After completing this module, you'll move to **07-eth-call** where you'll:
- Call contract functions without sending transactions
- Encode function calls manually (ABI encoding)
- Decode return values
- Understand the difference between `eth_call` and `eth_sendTransaction`
