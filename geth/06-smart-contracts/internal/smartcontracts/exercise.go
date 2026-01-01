//go:build !solution && !reference

package smartcontracts

/*
geth/06-smart-contracts: Smart Contract Interaction Fundamentals

This module is primarily a tutorial-based learning experience using the Geth console.
The exercises are designed to be completed in the Geth JavaScript console, not in Go.

After completing this module, you'll understand:
  - Call vs Transaction distinction (read-only vs state-changing)
  - Contract addresses and ABIs
  - Console-based contract interaction
  - Event decoding
  - Common pitfalls and debugging techniques

The Go code in this module serves as a reference showing how console concepts
map to Go implementations. The actual learning happens in the Geth console
following the tutorial in README.md.

To complete this module:
1. Read README.md thoroughly
2. Start Geth in dev mode: geth --dev --http --http.api eth,net,web3,personal
3. Attach to console: geth attach
4. Follow the step-by-step tutorial in README.md
5. Complete the exercises in the console
6. After mastering console concepts, proceed to geth/07-eth-call for Go implementation

This approach ensures you understand what happens under the hood before
writing Go code that abstracts these details away.
*/

// Placeholder type for module structure
// The actual exercises are console-based (see README.md)
type Tutorial struct {
	// This module teaches concepts via Geth console
	// See README.md for step-by-step tutorial
}

/*
Console Exercises (from README.md):

Exercise 1: Connect and Query
  - Start Geth in dev mode
  - Attach to console
  - Check chain ID and balance
  - Query latest block number

Exercise 2: Create Contract Object
  - Load ERC20 ABI
  - Create contract object for USDC or test contract
  - Call totalSupply() and balanceOf()

Exercise 3: Send Transaction
  - Deploy or use existing ERC20 on dev chain
  - Unlock account
  - Send transfer() transaction
  - Get transaction receipt
  - Verify success

Exercise 4: Decode Events
  - Extract logs from receipt
  - Calculate Transfer event topic hash
  - Find matching log
  - Decode from, to, and value

Exercise 5: Common Errors
  - Trigger and fix common mistakes
  - Wrong network, missing 'from' field, etc.

For Go implementations of these concepts, see:
  - geth/07-eth-call: Making contract calls from Go
  - geth/08-abigen: Type-safe Go bindings
  - geth/09-events: Event subscriptions and decoding
*/
