# 07-eth-call: Read-Only Contract Calls

**Goal:** Perform read-only contract calls with manual ABI encoding/decoding.

## Big Picture: Simulating Transactions Without State Changes

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

## Learning Objectives

By the end of this module, you should be able to:

1. **Pack ABI** for simple view functions (ERC20 name/symbol/decimals/totalSupply)
2. **Call contracts** with `CallContract` and decode results
3. **Handle reverts** and raw return data
4. **Understand ABI encoding** (function selector + arguments)
5. **Decode return values** based on function return types

## Prerequisites

- **Modules 01-06:** RPC basics, keys/tx basics, transaction building
- **Basic ABI understanding:** Function signatures, parameter types
- **Go basics:** Error handling, type assertions

## Building on Previous Modules

### From Module 05-06 (05-tx-nonces, 06-eip1559)
- You learned to build and send transactions (state-changing)
- Now you're learning to **query** contracts without sending transactions
- Same ABI encoding concepts, but no signing or broadcasting

### From Module 04 (04-accounts-balances)
- You learned to query account balances
- Now you're querying **contract state** (functions, not just balances)
- Contracts can have complex state that requires function calls to read

### Connection to Solidity-edu
- **02 Functions & Payable:** Read/write distinction mirrors view/pure vs state-changing functions
- **03 Events & Logging:** Pair calls with log decoding for richer off-chain views
- **08 ERC20 from Scratch:** ERC20 tokens have view functions (name, symbol, decimals, totalSupply)

## Understanding ABI Encoding

### Function Selector

**Function selector** = First 4 bytes of `keccak256(functionSignature)`

**Example:**
- Function: `name() returns (string)`
- Signature: `"name()"`
- Hash: `keccak256("name()")` = `0x06fdde03...`
- Selector: `0x06fdde03` (first 4 bytes)

**Computer Science principle:** Selectors are like **hash table keys**. They allow the EVM to quickly identify which function to call without parsing the entire function name.

### ABI Encoding Process

1. **Function selector:** First 4 bytes (identifies function)
2. **Arguments:** ABI-encoded parameters (if any)
3. **Result:** Concatenated bytes sent as `data` field

**For functions with no arguments:**
- Data = function selector only (4 bytes)
- Example: `name()` → `0x06fdde03`

**For functions with arguments:**
- Data = function selector + encoded arguments
- Example: `balanceOf(address)` → `0x70a08231` + encoded address

## Real-World Analogies

### The Database Query Analogy
- **`eth_call`:** SELECT query (read-only, no changes)
- **`eth_sendTransaction`:** INSERT/UPDATE query (changes data)
- **ABI encoding:** SQL query syntax
- **Decoding:** Parsing query results

### The Library Analogy
- **`eth_call`:** Asking a librarian to look up information (no changes to books)
- **`eth_sendTransaction`:** Checking out a book (changes library state)
- **Function selector:** Book title/call number
- **Arguments:** Specific page or chapter to read

### The CPU Analogy
- **`eth_call`:** Read-only syscall inspecting memory/state
- **`eth_sendTransaction`:** Syscall that modifies memory/state
- **ABI encoding:** Function call convention (how to pass parameters)

## Fun Facts & Nerdy Details

### Function Selector Collisions

**Problem:** Different function signatures can have the same selector (first 4 bytes)

**Example:**
- `transfer(address,uint256)` → selector `0xa9059cbb`
- `transfer(uint256,address)` → different selector (different order)

**Solution:** Solidity compiler prevents this, but it's theoretically possible with different languages.

**Fun fact:** Selector collisions are extremely rare (4 bytes = 1 in 4 billion chance), but they can happen!

### Revert Handling

**When a contract reverts:**
- `eth_call` returns an error
- Error message may contain revert reason (if contract uses `require()` or `revert()` with message)
- Raw return data contains encoded error information

**Decoding reverts:**
- Error data starts with selector `0x08c379a0` (Error(string) selector)
- Followed by ABI-encoded error message

**Nerdy detail:** Reverts are actually **successful executions** that return error data. The EVM doesn't distinguish between reverts and errors—both return data.

### Gas Estimation with `eth_call`

**`eth_call` can estimate gas:**
- Set `Gas` field in `CallMsg` to estimate required gas
- Node simulates execution and reports gas used
- Useful for estimating transaction costs before sending

**Production tip:** Always estimate gas before sending transactions to avoid failures!

## Comparisons

### Manual ABI vs Typed Bindings
| Aspect | Manual (this module) | Typed Bindings (module 08) |
|--------|---------------------|---------------------------|
| Type safety | ❌ Runtime errors | ✅ Compile-time checks |
| Boilerplate | ❌ High | ✅ Low |
| Flexibility | ✅ High | ❌ Lower |
| Use case | One-off calls | Production code |

### `eth_call` vs `eth_sendTransaction`
| Aspect | `eth_call` | `eth_sendTransaction` |
|--------|------------|----------------------|
| Execution | Local (simulated) | On-chain (persisted) |
| Gas cost | Free | Paid |
| Speed | Fast | Slower (needs mining) |
| State changes | No | Yes |

### Go `ethclient` vs JavaScript `ethers.js`
- **Go:** `client.CallContract(ctx, callMsg, nil)` → Returns `[]byte`
- **JavaScript:** `contract.name()` → Auto-encodes/decodes, returns typed value
- **Same JSON-RPC:** Both use `eth_call` under the hood

## Related Solidity-edu Modules

- **02 Functions & Payable:** Read/write distinction mirrors view/pure vs state-changing functions
- **03 Events & Logging:** Pair calls with log decoding for richer off-chain views
- **08 ERC20 from Scratch:** ERC20 tokens have view functions you'll call in this module

## What You'll Build

In this module, you'll create a CLI that:
1. Takes a contract address and function name as input
2. Encodes the function call using ABI (function selector)
3. Executes `eth_call` to simulate the function execution
4. Decodes the return value based on function type
5. Displays the result

**Key learning:** You'll understand how to manually encode/decode ABI data, giving you full control over contract interactions!

## Files

- **Starter:** `exercise/exercise.go` - Student entry point with TODO guidance
- **Solution:** `exercise/solution.go` - Reference implementation (run with `go test -tags solution ./07-eth-call/...`)
- **Tests:** `exercise/exercise_test.go` - Covers ABI encoding/decoding edge cases

## How to Run Tests

To run the tests for this module:

```bash
# From the project root (go-edu/)
cd geth/07-eth-call
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
2. **Helper Function Creation** - Understand DRY principle and closures
3. **Contract Calling** - Learn eth_call basics and CallMsg structure
4. **ABI Decoding** - Understand static vs dynamic type encoding
5. **Error Handling** - Learn error wrapping and context preservation

### The Solution File (`exercise/solution.go`)

The solution file contains detailed educational comments explaining:
- **Why** each step is necessary (the reasoning behind the code)
- **How** concepts repeat and build on each other (pattern recognition)
- **What** fundamental principles are being demonstrated (computer science concepts)

### Key Patterns You'll Learn

#### Pattern 1: Helper Functions with Closures
```go
// BAD: Repeat CallMsg construction for each call
msg1 := ethereum.CallMsg{To: &addr, Data: selector1}
client.CallContract(ctx, msg1, nil)
msg2 := ethereum.CallMsg{To: &addr, Data: selector2}
client.CallContract(ctx, msg2, nil)

// GOOD: Extract into helper closure
call := func(selector []byte) ([]byte, error) {
    msg := ethereum.CallMsg{To: &cfg.Contract, Data: selector}
    return client.CallContract(ctx, msg, cfg.BlockNumber)
}
nameBytes, _ := call(selectorName)
symbolBytes, _ := call(selectorSymbol)
```

**Why:** Reduces duplication (DRY principle). If call logic changes, update one place.

**Building on:** Module 01-stack taught error handling patterns. Here we add code reuse.

**Repeats in:** Every module that makes multiple similar calls (08-abigen, 09-events).

#### Pattern 2: Function Selector Computation
```go
func selector(sig string) []byte {
    hash := crypto.Keccak256([]byte(sig))
    return hash[:4]
}

selectorName := selector("name()") // 0x06fdde03
```

**Why:** Function selectors are the first 4 bytes of the keccak256 hash of the function signature. This is how the EVM identifies which function to call.

**Building on:** New concept, but builds on understanding of hashing from blockchain fundamentals.

**Repeats in:** Every module that does manual ABI encoding (module 07, manual contract interactions).

#### Pattern 3: Two-Phase Call Pattern (Network + Decode)
```go
// Phase 1: Network operation
rawBytes, err := call(selector)
if err != nil {
    return nil, fmt.Errorf("call failed: %w", err)
}

// Phase 2: Data processing
value, err := decode(rawBytes)
if err != nil {
    return nil, fmt.Errorf("decode failed: %w", err)
}
```

**Why:** Separates concerns. Network failures vs parsing failures have different causes and solutions.

**Building on:** Module 01-stack separated RPC calls from data processing. Same pattern here.

**Repeats in:** All contract interaction modules (07, 08, 09).

#### Pattern 4: Understanding Static vs Dynamic ABI Types
```go
// Static types (uint8, uint256): Simple 32-byte decoding
func decodeUint8(data []byte) (uint8, error) {
    return uint8(data[len(data)-1]), nil
}

// Dynamic types (string): Offset + length + data
func decodeString(data []byte) (string, error) {
    offset := readUint256(data[0:32])    // Where is the data?
    length := readUint256(data[offset:]) // How long is it?
    return string(data[offset+32:offset+32+length]), nil
}
```

**Why:** Static types fit in one 32-byte word. Dynamic types need structure (offset → length → data) because size varies.

**Building on:** New ABI-specific concept. Fundamental for manual contract interactions.

**Repeats in:** Any manual ABI encoding/decoding (module 07, advanced use cases).

## Error Handling: Building Robust Systems

### Common Contract Call Errors

**1. "execution reverted"**
```
Cause: Contract function reverted (require/revert statement failed)
Solution: Check contract conditions, verify contract is deployed
Prevention: Test contract address, understand function requirements
```

**2. "invalid contract address"**
```
Cause: Contract address is invalid or contract not deployed
Solution: Verify address format, check if contract exists at address
Prevention: Always validate address format before calling
```

**3. "abi: attempting to unmarshall an empty string"**
```
Cause: Contract call returned no data (contract might not exist)
Solution: Verify contract exists at address, check if function is view/pure
Prevention: Use HeaderByNumber to verify block exists before historical queries
```

**4. "data too short for string"**
```
Cause: ABI decoding failed (data truncated or malformed)
Solution: Verify contract implements ERC20 standard correctly
Prevention: Use try-catch in Solidity, validate return data length
```

### Error Wrapping Strategy

```go
// Layer 1: RPC error
err := client.CallContract(ctx, msg, nil)
// Error: "execution reverted"

// Layer 2: Add context
return fmt.Errorf("call name(): %w", err)
// Error: "call name(): execution reverted"

// Layer 3: Caller adds more context
return fmt.Errorf("query token metadata: %w", err)
// Error: "query token metadata: call name(): execution reverted"
```

This creates a traceable error chain that shows exactly where and why the failure occurred.

## Testing Strategy

The test file (`exercise_test.go`) demonstrates several important patterns:

1. **Mock implementations:** `mockCallClient` implements `CallClient` interface
2. **Table-driven tests:** Multiple test cases with different scenarios
3. **ABI encoding verification:** Tests verify correct function selectors
4. **Decoding tests:** Tests verify correct parsing of return values
5. **Error case testing:** Tests verify error handling works correctly

**Key insight:** Because we use interfaces, we can test our logic without needing a real Ethereum node. This makes tests fast, reliable, and deterministic.

**Example test case:**
```go
{
    name: "successful token query",
    setup: func(m *mockCallClient) {
        m.responses["name"] = encodeString("My Token")
        m.responses["symbol"] = encodeString("MTK")
        m.responses["decimals"] = encodeUint8(18)
        m.responses["totalSupply"] = encodeUint256(big.NewInt(1000000))
    },
    validate: func(t *testing.T, res *Result) {
        if res.Name != "My Token" {
            t.Errorf("name = %s, want My Token", res.Name)
        }
        if res.Symbol != "MTK" {
            t.Errorf("symbol = %s, want MTK", res.Symbol)
        }
        if res.Decimals != 18 {
            t.Errorf("decimals = %d, want 18", res.Decimals)
        }
    },
}
```

## Common Pitfalls & How to Avoid Them

### Pitfall 1: Incorrect Function Signature
```go
// BAD: Wrong signature (missing parameter)
selector("balanceOf()") // Wrong! Missing address parameter

// GOOD: Correct signature with parameter types
selector("balanceOf(address)") // Correct!
```

**Why it's a problem:** Function selector includes parameter types. Wrong signature = wrong selector = function not found.

**Fix:** Always include parameter types in signature: `functionName(type1,type2)` with no spaces.

### Pitfall 2: Not Validating Contract Address
```go
// BAD: Don't check if address is valid
client.CallContract(ctx, msg, nil) // Might call zero address

// GOOD: Validate contract address
if cfg.Contract == (common.Address{}) {
    return errors.New("contract address required")
}
```

**Why it's a problem:** Calling zero address or invalid address returns empty data, causing decode errors.

**Fix:** Always validate contract address before calling.

### Pitfall 3: Ignoring Block Number Parameter
```go
// BAD: Hard-code nil (always queries latest)
client.CallContract(ctx, msg, nil)

// GOOD: Use cfg.BlockNumber (allows historical queries)
client.CallContract(ctx, msg, cfg.BlockNumber)
```

**Why it's a problem:** Tests need deterministic behavior. Hard-coding nil means you can't query historical state.

**Fix:** Always pass block number from config, use nil as default in config.

### Pitfall 4: Incorrect String Decoding
```go
// BAD: Treat string as static type (wrong!)
func decodeString(data []byte) string {
    return string(data[:32]) // Wrong! Ignores offset and length
}

// GOOD: Parse offset, length, then data
func decodeString(data []byte) (string, error) {
    offset := readOffset(data[:32])
    length := readLength(data[offset:offset+32])
    return string(data[offset+32:offset+32+length]), nil
}
```

**Why it's a problem:** Strings are dynamic types. They use offset + length + data encoding, not direct bytes.

**Fix:** Always parse offset and length for dynamic types (strings, bytes, arrays).

### Pitfall 5: Not Checking Data Length Before Decoding
```go
// BAD: Assume data is long enough
func decodeUint256(data []byte) *big.Int {
    return new(big.Int).SetBytes(data[len(data)-32:]) // Panics if len(data) < 32!
}

// GOOD: Validate length first
func decodeUint256(data []byte) (*big.Int, error) {
    if len(data) < 32 {
        return nil, errors.New("data too short for uint256")
    }
    return new(big.Int).SetBytes(data[len(data)-32:]), nil
}
```

**Why it's a problem:** Short data causes slice index out of bounds panic.

**Fix:** Always validate data length before accessing slices.

## How Concepts Build on Each Other

This module builds on patterns from previous modules while introducing new concepts:

1. **From Module 01-stack:**
   - Context validation → Same pattern here
   - RPC call pattern → Extended for contract calls
   - Error wrapping → Consistent usage
   - Defensive programming → Applied to ABI decoding

2. **New in this module:**
   - Function selectors (keccak256 hashing)
   - ABI encoding/decoding (static vs dynamic types)
   - eth_call vs eth_sendTransaction (read vs write)
   - Helper functions and closures (DRY principle)
   - Manual memory management (parsing byte arrays)

3. **Patterns that repeat throughout the course:**
   - Input validation → Every function
   - Error wrapping → All error returns
   - Two-phase operations (network + processing) → All RPC interactions
   - Config-based behavior → All configurable functions

**The progression:**
- Module 01: Read chain metadata (headers, IDs)
- Module 06: Write transactions (state changes)
- Module 07: Read contract state (view functions)
- Module 08: Typed bindings (convenience layer)
- Module 09: Read events (historical logs)

Each module layers new concepts on top of existing patterns, building your understanding incrementally.

## Next Steps

After completing this module, you'll move to **08-abigen** where you'll:
- Use typed contract bindings (generated code)
- Reduce boilerplate with compile-time type safety
- See how `abigen` generates Go code from ABIs
- Understand the trade-offs between manual and typed approaches
