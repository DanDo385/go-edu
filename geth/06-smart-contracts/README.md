# 06: Smart Contract Interaction Fundamentals

## What Is This Project About?

This module teaches the foundational concepts of smart contract interaction using the **Geth JavaScript console**. You’ll learn the critical distinction between **Calls** (read-only) and **Transactions** (state-changing), understand how to work with **contract addresses** and **ABIs**, and practice interacting with contracts before moving to Go implementations in later modules.

Unlike the surrounding modules which focus on Go APIs, this one is intentionally **hands-on in the console** first: you can see immediate feedback, inspect receipts/logs directly, and build strong intuition for what later Go code is actually doing under the hood.

## Why Is This Important?

Understanding contract interaction at the console level gives you:

- **Deep understanding**: See exactly what happens when you interact with contracts
- **Debugging skills**: The console is invaluable for debugging contract issues
- **Foundation for Go**: Later modules (07–09) build directly on these concepts
- **Real-world tool**: Many developers use Geth console for quick contract queries

## Real-World Problems This Solves

- **Debugging contract interactions**: Quickly test contract calls and inputs
- **Understanding RPC calls**: See what `eth_call` and `eth_sendTransaction` actually do
- **Learning contract ABIs**: Essential for any Ethereum development
- **Testing before coding**: Verify contract behavior before writing Go code

## Key Concepts You’ll Learn

- Call vs Transaction distinction (read-only vs state-changing)
- Contract addresses and ABIs
- Creating contract objects in the Geth console
- Making read-only calls
- Sending state-changing transactions
- Decoding events and logs
- Common pitfalls and debugging techniques

## Prerequisites

- Completion of `geth/01-stack` through `geth/05-tx-nonces`
- Understanding of basic Ethereum concepts (addresses, transactions)
- Geth installed on your system

## Step-by-Step Tutorial

### Step 1: Run Geth

#### Local Dev Chain (Fastest Feedback Loop)

```bash
geth --dev --http --http.api eth,net,web3,personal
```

What this does:

- `--dev`: Creates a local development chain with instant block times
- `--http`: Enables HTTP RPC endpoint (default port 8545)
- `--http.api`: Exposes APIs needed for contract interaction

Why dev chain? Instant feedback, free gas, no real money at risk.

#### Testnet (Sepolia Example)

```bash
geth \
  --sepolia \
  --http \
  --http.api eth,net,web3,personal \
  --syncmode snap
```

What this does:

- `--sepolia`: Connects to Sepolia testnet
- `--syncmode snap`: Faster sync mode
- Same APIs as dev chain

### Step 2: Open the Console

```bash
geth attach
```

What happens: You’re now in a JavaScript REPL wired directly into Ethereum.

Key insight: This console is your direct interface to the Ethereum node. Everything you do here maps to RPC calls that you’ll make from Go later.

### Step 3: Know Your Ingredients

To interact with a contract, you need exactly two things:

1) **Contract Address**: Where the contract lives on-chain

```text
0x1234567890123456789012345678901234567890
```

2) **ABI (Application Binary Interface)**: The “legend” explaining what functions exist

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

Why ABI? Without it, bytecode is opaque noise. The ABI tells you:

- What functions exist
- What parameters they take
- What they return
- Whether they’re `view`/`pure` (read-only) or state-changing

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

Key observation: Notice `"stateMutability": "view"` vs `"nonpayable"`. This tells you:

- `view`/`pure`: Read-only (**Call**)
- `nonpayable`/`payable`: State-changing (**Transaction**)

### Step 5: Create the Contract Object

```javascript
var contractAddress = "0x1234...abcd"
var myContract = eth.contract(abi).at(contractAddress)
```

What happens: Nothing touches the blockchain yet. You’ve created a local handle.

Behind the scenes: This creates a JavaScript object with methods matching your ABI functions.

### Step 6: Call (Read-Only)

If the function is `view` or `pure` in Solidity:

```javascript
// Read-only call - no gas, no signature needed
myContract.totalSupply()

myContract.balanceOf("0xabcd...5678")
```

What happens behind the scenes:

- No gas consumed
- No signature required
- Geth simulates execution against latest state
- Equivalent to `eth_call` RPC method

Key insight: This is free and instant. Perfect for reading contract state.

### Step 7: Send a Transaction (State Change)

Now we cross the Rubicon—state modification:

```javascript
// Step 7a: Unlock an account (if needed)
personal.unlockAccount(eth.accounts[0], "your-password", 0)

// Step 7b: Send the transaction
var txHash = myContract.transfer(
  "0xabcd...5678",  // to address
  100,              // amount
  {
    from: eth.accounts[0],
    gas: 100000
  }
)
```

What happens:

- Signs a transaction
- Broadcasts it to network
- Modifies state once mined
- Costs gas

Track the transaction:

```javascript
eth.getTransaction(txHash)
eth.getTransactionReceipt(txHash)
```

Key insight: Receipts are where reality lives. They contain:

- Gas used
- Status (success/failure)
- Logs (events)
- Block number

### Step 8: Decode Events (Logs)

Contracts speak back via events, not return values:

```javascript
var receipt = eth.getTransactionReceipt(txHash)
receipt.logs
```

To decode properly, you need:

- Event ABI
- Topic hashes (keccak256 of event signature)

```javascript
// Event signature: Transfer(address,address,uint256)
var transferTopic = web3.sha3("Transfer(address,address,uint256)")

receipt.logs.forEach(function(log) {
  if (log.topics[0] === transferTopic) {
    console.log("Transfer event found!")
    // Decode: topics[1] = from, topics[2] = to, data = value
  }
})
```

Why this pain? This is why libraries exist—but understanding builds expertise.

## Common “Why Is Nothing Working” Traps

- ❌ ABI doesn’t match deployed bytecode
  - Solution: Verify ABI matches the deployed contract
- ❌ Wrong network (mainnet vs Sepolia)
  - Solution: Check `eth.net.version` or `eth.chainId`
- ❌ Forgot `from:` field
  - Solution: Always specify sender for transactions
- ❌ Gas too low
  - Solution: Estimate gas first: `myContract.transfer.estimateGas(...)`
- ❌ Account still locked
  - Solution: Unlock before sending: `personal.unlockAccount(...)`
- ❌ Calling a non-view function without sending tx
  - Solution: Check `stateMutability`—if not `view`/`pure`, you must send a transaction

Geth is unforgiving. That’s a feature. It teaches you exactly what’s happening.

## Why This Matters (Especially for Go Developers)

Given your Go + Solidity + Geth trajectory, this skill plugs directly into:

- Custom Go services calling contracts via `ethclient`
- Mempool tooling that needs to understand contract interactions
- Geth forks and custom node implementations
- On-chain infra roles (clients, nodes, validators)
- DeFi backends that don’t rely on JS glue

If ethers.js is flying a drone, Geth is learning aerodynamics.

## Where This Goes Next

Once this feels natural, the next mental upgrade is:

- **Module 07 (eth-call)**: Doing the same interactions from Go using `go-ethereum/ethclient`
- **Module 08 (abigen)**: Typed bindings and how they map to `eth_call`
- **Module 09 (events)**: Watching logs/events and decoding them in Go

## Exercises

- Exercise 1: Connect to dev chain and query an ERC20 token’s `totalSupply`
- Exercise 2: Create a contract object and call multiple `view` functions
- Exercise 3: Send a transaction to transfer tokens (on dev chain)
- Exercise 4: Decode `Transfer` events from a transaction receipt
- Exercise 5: Identify common errors and their solutions

## Additional Resources

- Geth Console Documentation
- ABI Specification
- ERC20 Token Standard

## Exercise Files

The `internal/smartcontracts/exercise.go` file is intentionally minimal: this module is primarily a **tutorial-based learning experience** using the Geth console.
