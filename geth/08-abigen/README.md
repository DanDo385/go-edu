# 08-abigen: Typed Contract Bindings

**Goal:** Use typed contract bindings (abigen-style) for safer calls and transactions.

## Big Picture: From Manual to Typed

abigen turns ABI into typed Go methods. Instead of manual ABI pack/unpack (module 07), you get compile-time checked methods with `CallOpts` and `TransactOpts` for read/write. This reduces boilerplate and errors.

**Computer Science principle:** This is **code generation** - converting a schema (ABI) into type-safe code. It's like generating API clients from OpenAPI specs, or database models from schemas.

### The Evolution: Manual → Typed

**Module 07 (Manual):**
```go
data, _ := abi.Pack("name")
raw, _ := client.CallContract(ctx, callMsg, nil)
var name string
abi.UnpackIntoInterface(&name, "name", raw)
```

**Module 08 (Typed):**
```go
name, _ := contract.Name(ctx)
```

**Key benefit:** Compile-time type safety + less boilerplate!

## Learning Objectives

By the end of this module, you should be able to:

1. **Understand how abigen bindings** wrap `BoundContract` with typed methods
2. **Make typed view calls** (name/symbol/decimals/balanceOf)
3. **See how `CallOpts` and `TransactOpts`** carry context, block number, signer info
4. **Compare manual vs typed** approaches and their trade-offs

## Prerequisites

- **Module 07 (07-eth-call):** Manual ABI calls and Go interfaces
- **Go basics:** Interfaces, type assertions, error handling

## Building on Previous Modules

### From Module 07 (07-eth-call)
- You learned manual ABI encoding/decoding
- Now you're using **typed bindings** that handle encoding/decoding automatically
- Same underlying JSON-RPC calls, but with better ergonomics

### From Module 04 (04-accounts-balances)
- You learned to query account balances
- Now you're querying **contract functions** with type safety
- Balance queries can use typed bindings too!

### Connection to Solidity-edu
- **Events & Logging:** Bindings can also decode events (module 09)
- **Contract interactions:** Matches front-end libraries but in Go

## Understanding BoundContract

### What is BoundContract?

**BoundContract** = Contract address + ABI + Backend (RPC client)

It provides:
- **Call():** For view functions (read-only)
- **Transact():** For state-changing functions (write)
- **FilterLogs():** For event queries

**Computer Science principle:** This is the **adapter pattern** - wrapping low-level RPC calls with a high-level interface.

### CallOpts vs TransactOpts

**CallOpts** (for read operations):
- `Context`: For cancellation/timeouts
- `BlockNumber`: Which block to query (nil = latest)
- `From`: Optional sender address (for view functions that check `msg.sender`)

**TransactOpts** (for write operations):
- `Context`: For cancellation/timeouts
- `From`: Sender address
- `Signer`: Transaction signer function
- `Value`: ETH to send
- `GasPrice` / `GasFeeCap` / `GasTipCap`: Gas pricing
- `GasLimit`: Maximum gas to consume
- `Nonce`: Transaction nonce

**Key insight:** `CallOpts` is lightweight (just context), `TransactOpts` includes signing info.

## Real-World Analogies

### The Typed Remote Control Analogy
- **Manual ABI:** Raw hex payloads (like sending raw IR codes)
- **Typed Bindings:** Labeled buttons (like a TV remote with labeled buttons)
- **BoundContract:** The remote control itself
- **CallOpts/TransactOpts:** Settings (volume, channel, etc.)

### The CPU Analogy
- **Manual ABI:** Hand-rolled assembly (full control, error-prone)
- **Typed Bindings:** Syscall stubs (type-safe, less boilerplate)
- **BoundContract:** System call interface

### The Database ORM Analogy
- **Manual ABI:** Raw SQL queries
- **Typed Bindings:** ORM methods (type-safe, less boilerplate)
- **BoundContract:** Database connection + schema

## Fun Facts & Nerdy Details

### abigen Code Generation

**abigen CLI tool:**
```bash
abigen --abi token.abi --bin token.bin --pkg token --out token.go
```

**What it generates:**
- Typed contract struct
- Methods for each function (with proper types)
- Event structs and filters
- Helper functions

**Fun fact:** abigen reads Solidity compiler output (ABI JSON + bytecode) and generates Go code. It's like protobuf code generation!

### Type Safety Benefits

**Compile-time checks:**
- Wrong function name → Compile error
- Wrong argument types → Compile error
- Wrong return type → Compile error

**Runtime benefits:**
- Less error handling (types are guaranteed)
- Better IDE autocomplete
- Easier refactoring

**Nerdy detail:** Go's type system catches errors at compile time, preventing runtime failures. This is especially valuable for contract interactions where errors can be expensive!

### BoundContract Internals

**Under the hood:**
1. BoundContract stores ABI + address
2. `Call()` encodes function call using ABI
3. Executes `eth_call` via backend
4. Decodes return value using ABI
5. Returns typed Go value

**Same JSON-RPC:** Still uses `eth_call` under the hood, just with better ergonomics!

## Comparisons

### Manual ABI vs Typed Bindings
| Aspect | Manual (module 07) | Typed (this module) |
|--------|-------------------|-------------------|
| Type safety | ❌ Runtime errors | ✅ Compile-time checks |
| Boilerplate | ❌ High | ✅ Low |
| Flexibility | ✅ High | ❌ Lower |
| Use case | One-off calls | Production code |
| Learning | ✅ Understand internals | ✅ Faster development |

### CallOpts vs TransactOpts
| Aspect | CallOpts | TransactOpts |
|--------|----------|--------------|
| Use case | Read operations | Write operations |
| Signing | ❌ No | ✅ Yes |
| Gas pricing | ❌ No | ✅ Yes |
| Complexity | Low | High |

### Go abigen vs JavaScript ethers.js
- **Go:** Compile-time type safety, code generation
- **JavaScript:** Runtime type safety, dynamic typing
- **Same concept:** Both provide typed contract interfaces

## Related Solidity-edu Modules

- **Events & Logging:** Bindings can also decode events (module 09)
- **Contract interactions:** Matches front-end libraries but in Go
- **08 ERC20 from Scratch:** ERC20 tokens are perfect for learning bindings

## What You'll Build

In this module, you'll create a CLI that:
1. Takes a token address and optional holder address as input
2. Creates a BoundContract from ABI and address
3. Calls typed methods (Name, Symbol, Decimals, BalanceOf)
4. Displays token information and balance

**Key learning:** You'll see how typed bindings simplify contract interactions while maintaining type safety!

## Files

- **Starter:** `exercise/exercise.go`
- **Solution:** `exercise/solution.go` (build with `-tags solution`)
- **Tests:** `exercise/exercise_test.go`

## How to Run Tests

To run the tests for this module:

```bash
# From the project root (go-edu/)
cd geth/08-abigen
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
2. **ABI Parsing** - Understand contract interface definitions
3. **BoundContract Creation** - Learn the adapter pattern
4. **CallOpts Configuration** - Understand call contexts
5. **Typed Contract Calls** - Learn automatic encoding/decoding

### The Solution File (`exercise/solution.go`)

The solution file contains detailed educational comments explaining:
- **Why** each step is necessary (the reasoning behind the code)
- **How** concepts repeat and build on each other (pattern recognition)
- **What** fundamental principles are being demonstrated (computer science concepts)

### Key Patterns You'll Learn

#### Pattern 1: ABI Parsing and Validation
```go
// Parse ABI from JSON string
parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
if err != nil {
    return nil, fmt.Errorf("parse ABI: %w", err)
}

// Use parsed ABI to create BoundContract
contract := bind.NewBoundContract(address, parsedABI, backend, nil, nil)
```

**Why:** ABI defines the contract interface. Without it, we can't encode calls or decode returns.

**Building on:** Module 07 manually encoded/decoded. Here we let ABI do it automatically.

**Repeats in:** All typed contract interaction code (production Go + Ethereum projects).

#### Pattern 2: BoundContract Adapter Pattern
```go
// Create BoundContract with address + ABI + backend
contract := bind.NewBoundContract(cfg.Contract, parsedABI, backend, nil, nil)

// Use it to make calls (encoding/decoding automatic)
contract.Call(opts, &out, "name")
```

**Why:** Adapter pattern wraps low-level RPC with high-level interface. Automatic encoding/decoding based on ABI.

**Building on:** Module 07 showed manual encoding. Module 08 shows how libraries abstract this.

**Repeats in:** All abigen-generated contract bindings use BoundContract internally.

#### Pattern 3: CallOpts Configuration
```go
callOpts := &bind.CallOpts{
    Context:     ctx,
    BlockNumber: cfg.BlockNumber,
}
if cfg.Holder != nil {
    callOpts.From = *cfg.Holder
}
```

**Why:** Separates call configuration from contract definition. Different calls can use different options.

**Building on:** Separation of concerns from all previous modules. Config is separate from logic.

**Repeats in:** Every contract call in Go + Ethereum (both read and write operations).

#### Pattern 4: Type-Safe Helper Functions
```go
func callString(contract *bind.BoundContract, opts *bind.CallOpts, method string, params ...interface{}) (string, error) {
    var out []interface{}
    contract.Call(opts, &out, method, params...)
    return *abi.ConvertType(out[0], new(string)).(*string), nil
}
```

**Why:** Wraps contract.Call with type conversion. One helper per return type eliminates boilerplate.

**Building on:** DRY principle from module 07. Here we apply it to BoundContract calls.

**Repeats in:** Custom contract wrappers (when abigen-generated code isn't suitable).

## Error Handling: Building Robust Systems

### Common BoundContract Errors

**1. "invalid ABI JSON"**
```
Cause: ABI string is malformed JSON or invalid structure
Solution: Verify ABI format, use compiler-generated ABI
Prevention: Always validate ABI before parsing
```

**2. "execution reverted"**
```
Cause: Contract function reverted (require/revert statement)
Solution: Check function requirements, verify contract state
Prevention: Understand contract logic before calling
```

**3. "method 'xyz' not found"**
```
Cause: Method name doesn't exist in ABI
Solution: Verify method name matches ABI exactly (case-sensitive)
Prevention: Use abigen-generated bindings for compile-time checks
```

**4. "cannot use type X as type Y"**
```
Cause: Type conversion failed (wrong type for return value)
Solution: Verify ABI return type matches helper function type
Prevention: Use correct helper (callString vs callUint256, etc.)
```

### Error Wrapping Strategy

```go
// Layer 1: BoundContract error
err := contract.Call(opts, &out, "name")
// Error: "execution reverted"

// Layer 2: Add method context
return fmt.Errorf("call name: %w", err)
// Error: "call name: execution reverted"

// Layer 3: Caller adds domain context
return fmt.Errorf("query token metadata: %w", err)
// Error: "query token metadata: call name: execution reverted"
```

This creates a traceable error chain showing exactly what failed and why.

## Testing Strategy

The test file (`exercise_test.go`) demonstrates several important patterns:

1. **Mock implementations:** `mockContractCaller` implements `ContractCaller` interface
2. **ABI validation:** Tests verify correct ABI parsing
3. **Call encoding:** Tests verify parameters are encoded correctly
4. **Result decoding:** Tests verify return values are decoded correctly
5. **Error handling:** Tests verify errors propagate correctly

**Key insight:** BoundContract makes testing easier because you can mock the backend without implementing full ABI encoding/decoding logic.

## Common Pitfalls & How to Avoid Them

### Pitfall 1: Incorrect ABI JSON
```go
// BAD: Malformed JSON (missing quotes)
const abi = `[{name:transfer,type:function}]`

// GOOD: Valid JSON with proper quoting
const abi = `[{"name":"transfer","type":"function"}]`
```

**Why it's a problem:** abi.JSON() expects valid JSON. Malformed JSON causes parsing errors.

**Fix:** Always use compiler-generated ABI or validate JSON format.

### Pitfall 2: Wrong Number of Return Values
```go
// BAD: Expect single value when function returns multiple
name, _ := callString(contract, opts, "getInfo") // getInfo returns (string, uint256)

// GOOD: Handle all return values
var out []interface{}
contract.Call(opts, &out, "getInfo")
name := out[0].(string)
value := out[1].(*big.Int)
```

**Why it's a problem:** Solidity functions can return multiple values. Ignoring some causes incorrect decoding.

**Fix:** Check ABI for number of return values. Handle all of them.

### Pitfall 3: Not Setting CallOpts.From for View Functions
```go
// BAD: Don't set From when function checks msg.sender
callOpts := &bind.CallOpts{Context: ctx}
allowance, _ := callUint256(contract, callOpts, "allowance", owner, spender)

// GOOD: Set From if function needs msg.sender
callOpts := &bind.CallOpts{
    Context: ctx,
    From:    owner, // Function might check msg.sender
}
```

**Why it's a problem:** Some view functions check msg.sender for access control. Without From, they might revert or return wrong data.

**Fix:** Set From if the view function uses msg.sender internally.

### Pitfall 4: Reusing CallOpts Across Goroutines
```go
// BAD: Share CallOpts between goroutines (race condition)
opts := &bind.CallOpts{Context: ctx}
go func() { contract.Call(opts, &out1, "name") }()
go func() { contract.Call(opts, &out2, "symbol") }()

// GOOD: Create separate CallOpts for each goroutine
go func() {
    opts := &bind.CallOpts{Context: ctx}
    contract.Call(opts, &out1, "name")
}()
go func() {
    opts := &bind.CallOpts{Context: ctx}
    contract.Call(opts, &out2, "symbol")
}()
```

**Why it's a problem:** CallOpts is not thread-safe. Concurrent access causes data races.

**Fix:** Create separate CallOpts for each goroutine or protect with mutex.

### Pitfall 5: Ignoring ABI Parameter Types
```go
// BAD: Pass wrong type (int instead of address)
balance, _ := callUint256(contract, opts, "balanceOf", 12345)

// GOOD: Pass correct type matching ABI
addr := common.HexToAddress("0x...")
balance, _ := callUint256(contract, opts, "balanceOf", addr)
```

**Why it's a problem:** BoundContract encodes parameters based on ABI types. Wrong Go type causes encoding errors.

**Fix:** Always pass Go types that match ABI parameter types (address → common.Address, uint256 → *big.Int).

## How Concepts Build on Each Other

This module builds on patterns from previous modules while introducing new concepts:

1. **From Module 07-eth-call:**
   - Manual ABI encoding/decoding → Now automatic with BoundContract
   - Function selectors → Now handled by ABI package
   - CallMsg construction → Now done by BoundContract.Call
   - Same underlying eth_call, different abstraction level

2. **New in this module:**
   - ABI parsing (abi.JSON)
   - BoundContract adapter pattern
   - CallOpts configuration
   - Type conversion helpers (abi.ConvertType)
   - Automatic parameter encoding

3. **Patterns that repeat throughout the course:**
   - Input validation → Every function
   - Error wrapping → All error returns
   - Config-based behavior → All configurable functions
   - Separation of concerns → Config vs logic vs encoding

**The progression:**
- Module 01: Read chain metadata
- Module 06: Write transactions
- Module 07: Manual contract calls (low-level)
- Module 08: Typed contract calls (high-level)
- Module 09: Event decoding
- Future: abigen-generated bindings (fully automated)

Each module shows a different level of abstraction, building your understanding from low-level to high-level APIs.

## Next Steps

After completing this module, you'll move to **09-events** where you'll:
- Decode ERC20 Transfer events
- Understand topics vs data in logs
- Filter events by block range
- See how events complement function calls
