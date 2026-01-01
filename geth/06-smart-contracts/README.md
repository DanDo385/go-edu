# 06: Smart Contract Interaction Fundamentals

## What Is This Project About?

This module teaches the foundational concepts of smart contract interaction using the Geth console. You'll learn the critical distinction between **Calls** (read-only) and **Transactions** (state-changing), understand how to work with contract addresses and ABIs, and practice interacting with contracts before moving to Go implementations in later modules.

This module is intentionally console-based rather than Go-based. By learning contract interaction concepts first through the Geth console, you'll gain a deep understanding of what happens under the hood before we abstract it away with Go libraries in modules 07-09.

## Why Is This Important?

Understanding contract interaction at the console level gives you:

- **Deep understanding**: See exactly what happens when you interact with contracts
- **Debugging skills**: Console is invaluable for debugging contract issues
- **Foundation for Go**: Later modules (07-09) build on these concepts
- **Real-world tool**: Many developers use Geth console for quick contract queries

## Real-World Problems This Solves

- **Debugging contract interactions**: Console lets you quickly test contract calls
- **Understanding RPC calls**: See what `eth_call` and `eth_sendTransaction` actually do
- **Learning contract ABIs**: Essential for any Ethereum development
- **Testing before coding**: Verify contract behavior before writing Go code

## Key Concepts You'll Learn

- **Call vs Transaction distinction** (read-only vs state-changing)
- **Contract addresses and ABIs**: The two ingredients needed for contract interaction
- **Creating contract objects** in Geth console
- **Making read-only calls**: Query contract state without spending gas
- **Sending state-changing transactions**: Modify contract state
- **Decoding events and logs**: Understanding contract communication
- **Common pitfalls and debugging techniques**: Real-world troubleshooting

## Prerequisites

- Completion of `geth/01-stack` through `geth/05-tx-nonces`
- Understanding of basic Ethereum concepts (addresses, transactions)
- Geth installed on your system (`geth version` should work)
- Basic familiarity with JavaScript (Geth console uses JavaScript)

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

**Why dev chain?** Instant feedback, free gas, no real money at risk.

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

### Step 2: Open the Console

```bash
geth attach
```

**What happens:** You're now in a JavaScript REPL wired directly into Ethereum.

**Key insight:** This console is your direct interface to the Ethereum node. Everything you do here maps to RPC calls that you'll make from Go later.

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
- Whether they're view/pure (read-only) or state-changing

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

### Step 5: Create the Contract Object

```javascript
var contractAddress = "0x1234...abcd"
var myContract = eth.contract(abi).at(contractAddress)
```

**What happens:** Nothing touches the blockchain yet. You've created a local handle.

**Behind the scenes:** This creates a JavaScript object with methods matching your ABI functions.

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

**Key insight:** This is free and instant. Perfect for reading contract state.

### Step 7: Send a Transaction (State Change)

Now we cross the Rubicon - state modification:

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
- Costs gas

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

### Step 8: Decode Events (Logs)

Contracts speak back via events, not return values:

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

## Common "Why Is Nothing Working" Traps

❌ **ABI doesn't match deployed bytecode**
- **Solution:** Verify ABI matches the deployed contract

❌ **Wrong network** (mainnet vs Sepolia)
- **Solution:** Check `eth.net.version` or `eth.chainId`

❌ **Forgot `from:` field**
- **Solution:** Always specify sender for transactions

❌ **Gas too low**
- **Solution:** Estimate gas first: `myContract.transfer.estimateGas(...)`

❌ **Account still locked**
- **Solution:** Unlock before sending: `personal.unlockAccount(...)`

❌ **Calling a non-view function without sending tx**
- **Solution:** Check function `stateMutability` - if not `view`/`pure`, must send transaction

**Geth is unforgiving. That's a feature. It teaches you exactly what's happening.**

## Why This Matters (Especially for Go Developers)

Given your Go + Solidity + Geth trajectory, this skill plugs directly into:

- **Custom Go services** calling contracts via `ethclient`
- **MEV / mempool tooling** that needs to understand contract interactions
- **Geth forks and custom node implementations**
- **On-chain infra roles** (clients, nodes, validators)
- **DeFi backend systems** that don't rely on JS glue

If `ethers.js` is flying a drone, Geth is learning aerodynamics.

## Where This Goes Next

Once this feels natural, the next mental upgrade is:

- **Module 07 (eth-call)**: Doing the same interactions from Go using `go-ethereum/ethclient`
- **Module 08 (abigen)**: Understanding `eth_call` vs `SendTransaction` at the RPC level
- **Module 09 (events)**: Watching pending txs in the mempool and simulating execution

## Exercises

### Exercise 1: Connect to dev chain and query an ERC20 token's totalSupply

1. Start a dev chain: `geth --dev --http --http.api eth,net,web3,personal`
2. Deploy a simple ERC20 contract (or use an existing one)
3. Load the contract ABI
4. Create a contract object
5. Call `totalSupply()` and display the result

### Exercise 2: Create a contract object and call multiple view functions

1. Use the same contract from Exercise 1
2. Call `name()`, `symbol()`, `decimals()`, and `balanceOf(address)`
3. Display all results formatted nicely

### Exercise 3: Send a transaction to transfer tokens (on dev chain)

1. Unlock an account with tokens
2. Send a `transfer()` transaction
3. Wait for the transaction to be mined
4. Verify the balance changed by calling `balanceOf()` again

### Exercise 4: Decode Transfer events from a transaction receipt

1. After Exercise 3, get the transaction receipt
2. Extract logs from the receipt
3. Find Transfer events by topic hash
4. Decode the event data (from, to, value)

### Exercise 5: Identify common errors and their solutions

1. Try calling a non-view function without sending a transaction (observe error)
2. Try sending a transaction without unlocking the account (observe error)
3. Try using wrong ABI (observe error)
4. Document what each error means and how to fix it

## Additional Resources

- [Geth Console Documentation](https://geth.ethereum.org/docs/interface/javascript-console)
- [ABI Specification](https://docs.soliditylang.org/en/latest/abi-spec.html)
- [ERC20 Token Standard](https://eips.ethereum.org/EIPS/eip-20)
- [Ethereum JSON-RPC API](https://ethereum.org/en/developers/docs/apis/json-rpc/)

## Project Structure

```
06-smart-contracts/
├── cmd/
│   ├── app/          # Optional: Go wrapper demonstrating concepts
│   └── dev/          # Optional: Debug harness
├── internal/
│   └── smartcontracts/
│       ├── exercise.go
│       ├── exercise_test.go
│       ├── solution.reference.go
│       └── solution_no_err.reference.go
└── .vscode/
    └── launch.json   # Debug configurations
```

**Note:** This module is primarily console-based. The Go files are minimal and serve as reference implementations showing how the console concepts map to Go code.
