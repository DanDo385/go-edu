# 06: Smart Contract Interaction Fundamentals

## What Is This Project About?

This module teaches the foundational concepts of smart contract interaction using the Geth console. You'll learn the critical distinction between **Calls** (read-only) and **Transactions** (state-changing), understand how to work with contract addresses and ABIs, and practice interacting with contracts before moving to Go implementations in later modules.

This is primarily a tutorial-based learning experience. Unlike other modules where you write Go code from the start, here you'll work directly in the Geth JavaScript console to understand what happens under the hood when you interact with smart contracts. This hands-on experience with the console provides the conceptual foundation needed for modules 07-09.

## Why Is This Important?

Understanding contract interaction at the console level gives you:

- **Deep understanding**: See exactly what happens when you interact with contracts at the RPC level
- **Debugging skills**: The Geth console is invaluable for debugging contract issues in production
- **Foundation for Go**: Later modules (07-09) build directly on these concepts
- **Real-world tool**: Many professional Ethereum developers use the Geth console for quick contract queries and troubleshooting

## Real-World Problems This Solves

- **Debugging contract interactions**: The console lets you quickly test contract calls without writing any code
- **Understanding RPC calls**: See what `eth_call` and `eth_sendTransaction` actually do behind the scenes
- **Learning contract ABIs**: Essential knowledge for any Ethereum development, regardless of language
- **Testing before coding**: Verify contract behavior and responses before writing Go code
- **Production troubleshooting**: Query live contracts to investigate issues or verify state

## Key Concepts You'll Learn

- **Call vs Transaction distinction**: Understand the fundamental difference between read-only calls and state-changing transactions
- **Contract addresses and ABIs**: Learn how to locate and describe smart contracts
- **Creating contract objects**: Use the Geth console to create JavaScript handles for contracts
- **Making read-only calls**: Query contract state without consuming gas
- **Sending state-changing transactions**: Modify contract state and track transaction receipts
- **Decoding events and logs**: Understand how contracts communicate through events
- **Common pitfalls and debugging**: Learn to identify and fix typical mistakes

## Prerequisites

- Completion of `geth/01-stack` through `geth/05-tx-nonces`
- Understanding of basic Ethereum concepts (addresses, transactions, gas)
- Geth installed on your system: `geth version`

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
│       ├── exercise.go      # Minimal exercise file (primarily tutorial)
│       └── exercise_test.go # Tests if applicable
└── README.md                # This comprehensive tutorial guide
```

**Note**: This module is primarily console-based. The Go files are optional and serve to demonstrate how the console concepts map to Go code.

---

## Step-by-Step Tutorial

### Step 1: Run Geth

#### Option A: Local Dev Chain (Fastest Feedback Loop - Recommended)

```bash
geth --dev --http --http.api eth,net,web3,personal
```

**What this does**:
- `--dev`: Creates a local development chain with instant block times (no mining delays)
- `--http`: Enables HTTP RPC endpoint (default port 8545)
- `--http.api`: Exposes APIs needed for contract interaction

**Why dev chain?** Instant feedback, free gas, no real money at risk. Perfect for learning.

#### Option B: Testnet (Sepolia Example)

```bash
geth \
  --sepolia \
  --http \
  --http.api eth,net,web3,personal \
  --syncmode snap
```

**What this does**:
- `--sepolia`: Connects to Sepolia testnet
- `--syncmode snap`: Faster sync mode
- Same APIs as dev chain

**Note**: Testnet requires syncing and testnet ETH. Use dev chain for this tutorial.

---

### Step 2: Open the Console

In a separate terminal:

```bash
geth attach
```

**What happens**: You're now in a JavaScript REPL wired directly into Ethereum.

**Key insight**: This console is your direct interface to the Ethereum node. Everything you do here maps to RPC calls that you'll make from Go in later modules.

You should see a prompt like:

```
Welcome to the Geth JavaScript console!
>
```

---

### Step 3: Know Your Ingredients

To interact with a contract, you need exactly **two things**:

#### 1. Contract Address

Where the contract lives on-chain:

```
0x1234567890123456789012345678901234567890
```

#### 2. ABI (Application Binary Interface)

The "legend" explaining what functions exist:

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

**Where to get ABIs**:
- Etherscan (Verified contracts have an ABI tab)
- Your own Solidity compilation output
- Project documentation

---

### Step 4: Load the ABI

For this tutorial, we'll use a simplified ERC20 token ABI. Paste this into your Geth console:

```javascript
var abi = [
  {
    "constant": true,
    "inputs": [],
    "name": "totalSupply",
    "outputs": [{"name":"","type":"uint256"}],
    "stateMutability": "view",
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

**Key observation**: Notice `"stateMutability": "view"` vs `"nonpayable"`. This tells you:
- `view`/`pure`: Read-only → **Call**
- `nonpayable`/`payable`: State-changing → **Transaction**

---

### Step 5: Create the Contract Object

```javascript
// Use a real contract address (example: USDC on mainnet)
var contractAddress = "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"

// Create the contract object
var myContract = eth.contract(abi).at(contractAddress)
```

**What happens**: Nothing touches the blockchain yet. You've created a local JavaScript handle.

**Behind the scenes**: This creates a JavaScript object with methods matching your ABI functions.

**Verify it worked**:

```javascript
> myContract
{
  abi: [...],
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
  transactionHash: null,
  allEvents: function(),
  balanceOf: function(),
  totalSupply: function(),
  transfer: function()
}
```

---

### Step 6: Call (Read-Only)

If the function is `view` or `pure` in Solidity, you use **Call**:

```javascript
// Read-only call - no gas, no signature needed
myContract.totalSupply()
// Returns: 42427041540499498 (wei units for USDC)

// Query an account balance
myContract.balanceOf("0x742d35Cc6634C0532925a3b844Bc454e4438f44e")
// Returns: BigNumber { ... }
```

**What happens behind the scenes**:
1. No gas consumed
2. No signature required
3. Geth simulates execution against latest state
4. Equivalent to `eth_call` RPC method
5. Returns immediately

**Key insight**: This is **free and instant**. Perfect for reading contract state. No transaction is created.

**BREAKPOINT**: Try calling these functions. Watch the Variables panel in VS Code if debugging Go code later.

---

### Step 7: Send a Transaction (State Change)

Now we cross the Rubicon—**state modification**:

```javascript
// Step 7a: Unlock an account (if needed on dev chain)
personal.unlockAccount(eth.accounts[0], "your-password", 0)
// Returns: true

// Step 7b: Send the transaction
var txHash = myContract.transfer(
  "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",  // to address
  100000000,                                       // amount (USDC has 6 decimals)
  {
    from: eth.accounts[0],
    gas: 100000
  }
)
// Returns: "0x1234...txhash"
```

**What happens**:
1. Signs a transaction with the sender's private key
2. Broadcasts it to the network
3. Modifies state once mined
4. Costs gas (paid in ETH)
5. Creates a transaction receipt

**Track the transaction**:

```javascript
// Check transaction status
eth.getTransaction(txHash)
// Shows: blockNumber, gas, input data, etc.

// Get receipt (only available after mining)
eth.getTransactionReceipt(txHash)
// Shows: gasUsed, status, logs, etc.
```

**Key insight**: Receipts are where reality lives. They contain:
- `gasUsed`: Actual gas consumed
- `status`: `0x1` (success) or `0x0` (failure)
- `logs`: Events emitted
- `blockNumber`: Which block included this transaction

**BREAKPOINT**: This is the critical moment when blockchain state changes. In debugging, you'd set a breakpoint here to inspect transaction parameters.

---

### Step 8: Decode Events (Logs)

Contracts speak back via **events**, not return values:

```javascript
var receipt = eth.getTransactionReceipt(txHash)
receipt.logs
// Returns: Array of log objects
```

**Example log structure**:

```javascript
[{
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
  topics: [
    "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
    "0x000000000000000000000000abcd...",
    "0x0000000000000000000000001234..."
  ],
  data: "0x0000000000000000000000000000000000000000000000000000000005f5e100"
}]
```

To decode properly, you need:
1. Event ABI
2. Topic hashes (keccak256 of event signature)

```javascript
// Event signature: Transfer(address,address,uint256)
var transferTopic = web3.sha3("Transfer(address,address,uint256)")
// Returns: "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

// Find logs matching this topic
receipt.logs.forEach(function(log) {
  if (log.topics[0] === transferTopic) {
    console.log("Transfer event found!")
    // topics[1] = from address (indexed)
    // topics[2] = to address (indexed)
    // data = value (non-indexed)
  }
})
```

**Why this pain?** This manual process shows you exactly what libraries like `go-ethereum` do automatically. Understanding it builds expertise.

**Key insight**: Events are how contracts communicate results. Function return values are only accessible via `eth_call` (read-only). For state-changing transactions, you **must** use events to get outputs.

---

## Common "Why Is Nothing Working" Traps

### ❌ ABI doesn't match deployed bytecode

**Problem**: You load an ABI but the contract doesn't have those functions.

**Solution**: Verify ABI matches the deployed contract. Check Etherscan's "Contract" tab.

### ❌ Wrong network (mainnet vs Sepolia vs dev)

**Problem**: Your contract exists on one network but you're connected to another.

**Solution**: Check which network you're on:

```javascript
eth.chainId  // 1 = mainnet, 11155111 = Sepolia, 1337 = dev
```

### ❌ Forgot `from:` field

**Problem**: Transaction fails with "from address required"

**Solution**: Always specify sender for transactions:

```javascript
myContract.transfer(to, amount, { from: eth.accounts[0], gas: 100000 })
```

### ❌ Gas too low

**Problem**: Transaction reverts with "out of gas"

**Solution**: Estimate gas first:

```javascript
myContract.transfer.estimateGas(to, amount, { from: eth.accounts[0] })
// Returns: estimated gas units needed
```

### ❌ Account still locked

**Problem**: "authentication needed: password or unlock"

**Solution**: Unlock before sending:

```javascript
personal.unlockAccount(eth.accounts[0], "password", 300)  // unlock for 300 seconds
```

### ❌ Calling a non-view function without sending tx

**Problem**: You do `myContract.transfer(...)` directly but nothing happens

**Solution**: Check function `stateMutability`:
- If `view`/`pure`: Call directly → `myContract.balanceOf(addr)`
- If `nonpayable`/`payable`: Must send transaction → `myContract.transfer(..., {from: ..., gas: ...})`

**Geth is unforgiving. That's a feature. It teaches you exactly what's happening.**

---

## Why This Matters (Especially for Go Developers)

Given your Go + Solidity + Geth trajectory, this skill plugs directly into:

1. **Custom Go services** calling contracts via `ethclient.Client`
2. **MEV / mempool tooling** that needs to understand contract interactions
3. **Geth forks** and custom node implementations
4. **On-chain infra roles** (clients, nodes, validators)
5. **DeFi backend systems** that don't rely on JavaScript libraries

**If `ethers.js` is flying a drone, Geth console is learning aerodynamics.**

You'll understand:
- What `eth_call` actually does under the hood
- Why transaction receipts matter more than return values
- How ABIs map to RPC calls
- What happens during contract execution

---

## Where This Goes Next

Once this feels natural, the next mental upgrade is:

- **Module 07 (eth-call)**: Doing the same interactions from Go using `go-ethereum/ethclient`
- **Module 08 (abigen)**: Using `abigen` to generate type-safe Go bindings from ABIs
- **Module 09 (events)**: Subscribing to contract events and decoding logs programmatically

**This module is the conceptual foundation. Modules 07-09 are the Go implementation.**

---

## Exercises

### Exercise 1: Connect and Query

1. Start Geth in dev mode
2. Attach to the console
3. Check your chain ID and account balance
4. Query the latest block number

```javascript
eth.chainId
eth.accounts[0]
eth.getBalance(eth.accounts[0])
eth.blockNumber
```

### Exercise 2: Create a Contract Object

1. Use the ERC20 ABI provided above
2. Use USDC contract address: `0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48` (mainnet) or deploy your own on dev chain
3. Create the contract object
4. Call `totalSupply()` and `balanceOf(someAddress)`

### Exercise 3: Send a Transaction (Dev Chain Only)

1. Deploy a simple ERC20 contract on your dev chain (or use a pre-deployed test contract)
2. Unlock your dev account
3. Send a `transfer()` transaction
4. Get the transaction receipt
5. Verify the transaction succeeded (`status: 0x1`)

### Exercise 4: Decode Transfer Events

1. From the receipt in Exercise 3, extract the logs
2. Calculate the Transfer event topic hash: `web3.sha3("Transfer(address,address,uint256)")`
3. Find the log with matching topic
4. Identify the `from`, `to`, and `value` from topics and data

### Exercise 5: Common Errors

Intentionally trigger and fix these errors:
1. Call a function that doesn't exist
2. Send a transaction without `from:` field
3. Send a transaction with insufficient gas
4. Try to transfer tokens you don't have

---

## Additional Resources

- [Geth Console Documentation](https://geth.ethereum.org/docs/interacting-with-geth/javascript-console)
- [ABI Specification](https://docs.soliditylang.org/en/latest/abi-spec.html)
- [ERC20 Token Standard](https://eips.ethereum.org/EIPS/eip-20)
- [JSON-RPC API](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [Understanding eth_call vs eth_sendTransaction](https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_call)

---

## How to Run (Optional Go Wrapper)

Since this module is primarily tutorial-based, the Go files are optional. If you want to see how these console concepts map to Go code:

### Using cmd/app/main.go (CLI Arguments)

```bash
# Example: Query a contract
go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
```

### Using cmd/dev/main.go (Debug Harness)

```bash
# Run with fixed test inputs
go run ./cmd/dev/main.go

# Or use VS Code debugger (F5) with "Debug: cmd/dev" configuration
```

### Testing

```bash
# Run tests (if applicable)
go test ./...

# Run with verbose output
go test -v ./...
```

---

## Summary

This module teaches you **how to think about smart contract interaction**:

1. **Calls** are free, instant, read-only simulations
2. **Transactions** are signed, broadcast, state-changing, and cost gas
3. **ABIs** are the Rosetta Stone for contract interaction
4. **Events** are how contracts communicate state changes
5. **Receipts** contain the truth about what actually happened

**Master the console, and you'll understand what every Ethereum library is doing under the hood.**

Next stop: Doing all of this from Go code in modules 07-09. 🚀
