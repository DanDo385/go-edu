# 06: Smart Contract Interaction Fundamentals

## What Is This Project About?

This module teaches the foundational concepts of smart contract interaction using the Geth console. You'll learn the critical distinction between **Calls** (read-only) and **Transactions** (state-changing), understand how to work with contract addresses and ABIs, and practice interacting with contracts before moving to Go implementations in later modules.

Unlike other modules in this series which focus on Go implementations, this module is primarily **console-based**. The Geth JavaScript console provides immediate feedback and visual understanding of contract interactions, making it the ideal learning environment before diving into Go code.

## Why Is This Important?

Understanding contract interaction at the console level gives you:

- **Deep understanding**: See exactly what happens when you interact with contracts
- **Debugging skills**: Console is invaluable for debugging contract issues
- **Foundation for Go**: Later modules (07-09) build on these concepts
- **Real-world tool**: Many developers use Geth console for quick contract queries

If ethers.js is flying a drone, Geth is learning aerodynamics. You'll understand the underlying mechanics that all Ethereum libraries abstract away.

## Real-World Problems This Solves

- **Debugging contract interactions**: Console lets you quickly test contract calls before writing production code
- **Understanding RPC calls**: See what `eth_call` and `eth_sendTransaction` actually do under the hood
- **Learning contract ABIs**: Essential knowledge for any Ethereum development
- **Testing before coding**: Verify contract behavior before writing Go code
- **Quick queries**: Check token balances, contract states, and transaction statuses interactively

## Key Concepts You'll Learn

- **Call vs Transaction distinction**: Understanding read-only operations vs state-changing operations
- **Contract addresses and ABIs**: The two ingredients needed to interact with any contract
- **Creating contract objects**: How to build a contract handle in the Geth console
- **Making read-only calls**: Querying contract state without spending gas
- **Sending state-changing transactions**: Modifying contract state with proper gas and signatures
- **Decoding events and logs**: Understanding how contracts communicate back to you
- **Common pitfalls and debugging techniques**: Avoiding and fixing common mistakes

## Prerequisites

- Completion of `geth/01-stack` through `geth/05-tx-nonces`
- Understanding of basic Ethereum concepts (addresses, transactions, gas)
- Geth installed on your system (`geth version` should work)

## Project Structure

```
geth/06-smart-contracts/
├── cmd/
│   ├── app/
│   │   └── main.go          # Optional: Go wrapper demonstrating concepts
│   └── dev/
│       └── main.go          # Optional: Debug harness
├── internal/
│   └── smartcontracts/
│       ├── exercise.go      # Minimal - this is primarily a tutorial module
│       ├── exercise_test.go # Basic tests
│       ├── solution.reference.go
│       └── solution_no_err.reference.go
├── .vscode/
│   └── launch.json          # Debug configurations
└── README.md                 # This comprehensive tutorial guide
```

---

## Step-by-Step Tutorial

### Step 1: Run Geth

#### Local Dev Chain (Fastest Feedback Loop)

```bash
geth --dev --http --http.api eth,net,web3,personal
```

**What this does:**
- `--dev`: Creates a local development chain with instant block times
- `--http`: Enables HTTP RPC endpoint (default port 8545)
- `--http.api`: Exposes APIs needed for contract interaction

**Why dev chain?** Instant feedback, free gas, no real money at risk. Perfect for learning.

#### Testnet (Sepolia Example)

```bash
geth \
  --sepolia \
  --http \
  --http.api eth,net,web3,personal \
  --syncmode snap
```

**What this does:**
- `--sepolia`: Connects to Sepolia testnet
- `--syncmode snap`: Faster sync mode
- Same APIs as dev chain

---

### Step 2: Open the Console

```bash
geth attach
```

**What happens:** You're now in a JavaScript REPL wired directly into Ethereum.

**Key insight:** This console is your direct interface to the Ethereum node. Everything you do here maps to RPC calls that you'll make from Go in later modules.

---

### Step 3: Know Your Ingredients

To interact with a contract, you need exactly **two things**:

1. **Contract Address**: Where the contract lives on-chain
   ```
   0x1234567890123456789012345678901234567890
   ```

2. **ABI (Application Binary Interface)**: The "legend" explaining what functions exist
   ```json
   [
     {
       "constant": true,
       "inputs": [],
       "name": "totalSupply",
       "outputs": [{"name":"","type":"uint256"}],
       "type": "function"
     }
   ]
   ```

**Why ABI?** Without it, bytecode is opaque noise. The ABI tells you:
- What functions exist
- What parameters they take
- What they return
- Whether they're `view`/`pure` (read-only) or state-changing

---

### Step 4: Load the ABI

Paste your ABI JSON into the console:

```javascript
var abi = [
  {
    "constant": true,
    "inputs": [],
    "name": "totalSupply",
    "outputs": [{"name":"","type":"uint256"}],
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [{"name":"account","type":"address"}],
    "name": "balanceOf",
    "outputs": [{"name":"","type":"uint256"}],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": false,
    "inputs": [
      {"name":"to","type":"address"},
      {"name":"amount","type":"uint256"}
    ],
    "name": "transfer",
    "outputs": [{"name":"","type":"bool"}],
    "stateMutability": "nonpayable",
    "type": "function"
  }
]
```

**Key observation:** Notice `"stateMutability": "view"` vs `"nonpayable"`. This tells you:
- `view`/`pure`: Read-only (Call)
- `nonpayable`/`payable`: State-changing (Transaction)

---

### Step 5: Create the Contract Object

```javascript
var contractAddress = "0x1234...abcd"
var myContract = eth.contract(abi).at(contractAddress)
```

**What happens:** Nothing touches the blockchain yet. You've created a **local handle**.

**Behind the scenes:** This creates a JavaScript object with methods matching your ABI functions.

---

### Step 6: Call (Read-Only)

If the function is `view` or `pure` in Solidity:

```javascript
// Read-only call - no gas, no signature needed
myContract.totalSupply()
// Returns: BigNumber { s: 1, e: 18, c: [ 1000000 ] }

myContract.balanceOf("0xabcd...5678")
// Returns: BigNumber { s: 1, e: 18, c: [ 500 ] }
```

**What happens behind the scenes:**
- No gas consumed
- No signature required
- Geth simulates execution against latest state
- Equivalent to `eth_call` RPC method

**Key insight:** This is **free and instant**. Perfect for reading contract state.

---

### Step 7: Send a Transaction (State Change)

Now we cross the Rubicon—**state modification**:

```javascript
// Step 7a: Unlock an account (if needed)
personal.unlockAccount(eth.accounts[0], "your-password", 0)
// Returns: true

// Step 7b: Send the transaction
var txHash = myContract.transfer(
  "0xabcd...5678",  // to address
  100,              // amount
  {
    from: eth.accounts[0],
    gas: 100000
  }
)
// Returns: "0x1234...txhash"
```

**What happens:**
- Signs a transaction
- Broadcasts it to network
- Modifies state once mined
- **Costs gas**

**Track the transaction:**

```javascript
// Check transaction status
eth.getTransaction(txHash)

// Get receipt (only available after mining)
eth.getTransactionReceipt(txHash)
```

**Key insight:** Receipts are where reality lives. They contain:
- Gas used
- Status (success/failure)
- Logs (events)
- Block number

---

### Step 8: Decode Events (Logs)

Contracts speak back via **events**, not return values:

```javascript
var receipt = eth.getTransactionReceipt(txHash)
receipt.logs
// Returns: Array of log objects
```

To decode properly, you need:
- Event ABI
- Topic hashes (keccak256 of event signature)

```javascript
// Event signature: Transfer(address,address,uint256)
var transferTopic = web3.sha3("Transfer(address,address,uint256)")
// Returns: "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

// Find logs matching this topic
receipt.logs.forEach(function(log) {
  if (log.topics[0] === transferTopic) {
    console.log("Transfer event found!")
    // Decode: topics[1] = from, topics[2] = to, data = value
  }
})
```

**Why this pain?** This is why libraries exist—but understanding builds expertise.

---

## Common "Why Is Nothing Working" Traps

| Problem | Symptom | Solution |
|---------|---------|----------|
| ❌ ABI doesn't match deployed bytecode | Unexpected revert or wrong data | Verify ABI matches the deployed contract |
| ❌ Wrong network (mainnet vs Sepolia) | Contract not found or wrong state | Check `eth.chainId()` or `net.version` |
| ❌ Forgot `from:` field | Transaction rejected | Always specify sender for transactions |
| ❌ Gas too low | Out of gas error | Estimate gas first: `myContract.transfer.estimateGas(...)` |
| ❌ Account still locked | Transaction signing failed | Unlock before sending: `personal.unlockAccount(...)` |
| ❌ Calling non-view function without tx | No state change | Check function `stateMutability`—if not `view`/`pure`, must send transaction |

**Geth is unforgiving. That's a feature.** It teaches you exactly what's happening.

---

## Why This Matters (Especially for Go Developers)

Given your Go + Solidity + Geth trajectory, this skill plugs directly into:

- **Custom Go services** calling contracts via `ethclient`
- **MEV / mempool tooling** that needs to understand contract interactions
- **Geth forks** and custom node implementations
- **On-chain infra roles** (clients, nodes, validators)
- **DeFi backend systems** that don't rely on JS glue

---

## Where This Goes Next

Once this feels natural, the next mental upgrade is:

| Module | Focus |
|--------|-------|
| **07 (eth-call)** | Doing the same interactions from Go using `go-ethereum/ethclient` |
| **08 (abigen)** | Understanding `eth_call` vs `SendTransaction` at the RPC level with typed bindings |
| **09 (events)** | Watching events and logs programmatically in Go |

---

## Exercises

### Exercise 1: Connect and Query

1. Start a Geth dev chain
2. Attach to the console
3. Query an ERC20 token's `totalSupply` (use USDC on mainnet if testing against real network)

### Exercise 2: Create a Contract Object

1. Find the ABI for a known contract (e.g., USDC, DAI, or any ERC20)
2. Create a contract object in the console
3. Call multiple view functions (`name()`, `symbol()`, `decimals()`, `balanceOf()`)

### Exercise 3: Send a Transaction (Dev Chain Only)

1. On the dev chain, create an account or use `eth.accounts[0]`
2. Deploy a simple ERC20 or use the pre-funded dev account
3. Send a `transfer` transaction
4. Check the transaction receipt

### Exercise 4: Decode Transfer Events

1. After a successful transfer, get the transaction receipt
2. Find the Transfer event topic hash
3. Match and decode the event data from logs

### Exercise 5: Debug Common Errors

1. Intentionally cause common errors (wrong ABI, locked account, low gas)
2. Read the error messages and understand what they mean
3. Fix each error

---

## How to Run

### Console Tutorial (Primary Method)

```bash
# 1. Start Geth dev chain
geth --dev --http --http.api eth,net,web3,personal

# 2. In another terminal, attach to console
geth attach

# 3. Follow the tutorial steps above
```

### Using cmd/app/main.go (Optional Go Demo)

```bash
# Basic usage (demonstrates concepts from console in Go)
go run ./cmd/app/main.go [RPC_URL] [CONTRACT_ADDRESS]

# Example with USDC on mainnet
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
```

### Using cmd/dev/main.go (Debug Harness)

```bash
# Run with fixed test inputs
go run ./cmd/dev/main.go

# Or use VS Code debugger (F5) with "Debug: cmd/dev" configuration
```

## How to Debug

1. Set breakpoints at `// BREAKPOINT:` comments
2. Use VS Code debugger (F5) and select appropriate configuration:
   - "Debug: cmd/app" - Debug with CLI arguments
   - "Debug: cmd/dev" - Debug with fixed inputs (recommended for learning)
   - "Test: Run All Tests" - Debug tests
3. Step through code using F10 (Step Over) and F11 (Step Into)
4. Watch variables in the Variables panel

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with reference implementation
go test -tags=reference -v ./...
```

---

## Additional Resources

- [Geth Console Documentation](https://geth.ethereum.org/docs/interacting-with-geth/javascript-console)
- [ABI Specification](https://docs.soliditylang.org/en/latest/abi-spec.html)
- [ERC20 Token Standard](https://eips.ethereum.org/EIPS/eip-20)
- [Ethereum JSON-RPC Specification](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [go-ethereum ethclient Package](https://pkg.go.dev/github.com/ethereum/go-ethereum/ethclient)
