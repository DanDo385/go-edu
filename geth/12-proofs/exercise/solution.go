//go:build solution
// +build solution

package exercise

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

/*
Problem: Fetch Merkle-Patricia trie proofs for accounts and storage slots.

Merkle-Patricia trie proofs are cryptographic receipts that prove "account X has
balance Y and storage slot Z has value W" without downloading the entire blockchain
state. This enables:
  - Light clients that verify state without full sync
  - Cross-chain bridges that prove state on one chain to another
  - Indexers that verify indexed data is correct
  - Trust-minimized verification of contract state

Computer science principles highlighted:
  - Merkle trees provide logarithmic proof size (log N instead of N)
  - Cryptographic commitment (root hash commits to all data)
  - Path-based verification (prove membership by providing path to root)
*/
func Run(ctx context.Context, client ProofClient, cfg Config) (*Result, error) {
	// ============================================================================
	// STEP 1: Input Validation - Defensive Programming Pattern
	// ============================================================================
	// Why validate inputs? Proof RPC calls are computationally expensive on the
	// node side (they must traverse the Merkle-Patricia trie and collect proof
	// nodes). Validating inputs early ensures we don't waste resources.
	//
	// Context handling: Same pattern as modules 01 and 11. If ctx is nil, we
	// provide context.Background() as a safe default. This ensures our RPC calls
	// always have a valid context for cancellation and timeout handling.
	//
	// This pattern repeats: validate → default → proceed.
	if ctx == nil {
		ctx = context.Background()
	}

	// Client validation: The ProofClient interface is our dependency for making
	// RPC calls. If it's nil, we can't proceed. Returning early with a descriptive
	// error is Go's idiomatic error handling pattern.
	//
	// Building on modules 01 and 11: Same validation pattern, different interface.
	// Each module uses specialized interfaces (RPCClient, StorageClient, ProofClient)
	// but the validation logic is identical.
	if client == nil {
		return nil, errors.New("client is nil")
	}

	// Account validation: Proofs are for specific accounts. The zero address
	// (0x0000...0000) is not a valid account to prove. While you could technically
	// get a proof showing "account 0x0 doesn't exist," it's more likely a caller
	// error, so we reject it early.
	//
	// Connection to module 11: Same address validation pattern. Consistent checks
	// across modules make code predictable and easier to understand.
	if cfg.Account == (common.Address{}) {
		return nil, errors.New("account address required")
	}

	// ============================================================================
	// STEP 2: Convert Slots to Hex Strings - Interface Adaptation Pattern
	// ============================================================================
	// The eth_getProof RPC method expects storage slot identifiers as hex strings,
	// not as common.Hash types. This is a quirk of the JSON-RPC interface—some
	// methods use typed parameters, others use strings.
	//
	// Why this matters: When building Ethereum tooling, you'll frequently need to
	// convert between Go types (common.Hash, big.Int) and RPC-friendly formats
	// (hex strings). This conversion layer is essential for clean interfaces.
	//
	// Implementation details:
	//   - Only allocate slotStrings if we have slots to convert (avoid empty alloc)
	//   - Pre-allocate with exact capacity to avoid slice growth
	//   - Use .Hex() which adds "0x" prefix automatically
	//
	// Building on previous concepts: This is similar to how module 11 converted
	// slot numbers to hashes. Here we're converting hashes to strings. Each layer
	// of the stack requires its own representation format.
	var slotStrings []string
	if len(cfg.Slots) > 0 {
		slotStrings = make([]string, len(cfg.Slots))
		for i, slot := range cfg.Slots {
			slotStrings[i] = slot.Hex()
		}
	}

	// ============================================================================
	// STEP 3: Fetch Proof - Understanding eth_getProof
	// ============================================================================
	// This calls eth_getProof under the hood, which is one of the most powerful
	// RPC methods in Ethereum. It returns:
	//
	// 1. Account Proof:
	//    - Account state (balance, nonce, codeHash, storageHash)
	//    - Proof nodes (path from stateRoot to account in the account trie)
	//
	// 2. Storage Proofs (one per requested slot):
	//    - Slot value
	//    - Proof nodes (path from storageRoot to slot in the storage trie)
	//
	// Computer Science insight: This is a Merkle proof! The proof nodes are
	// the siblings along the path from leaf to root. With these nodes, you can:
	//   1. Hash your way up the tree
	//   2. Verify the final hash matches the known root
	//   3. Prove the data is part of the committed state
	//
	// Why this is powerful:
	//   - Light clients: Verify state without downloading 100GB+ of blockchain data
	//   - Bridges: Prove state on one chain to another chain's smart contract
	//   - Indexers: Verify indexed data matches on-chain state
	//
	// Parameters:
	//   - ctx: Context for cancellation/timeout
	//   - cfg.Account: The account to prove
	//   - slotStrings: Storage slots to prove (can be empty for account-only proof)
	//   - cfg.BlockNumber: Which block to prove from (nil = latest)
	//
	// Building on module 11: We used the same storage slot calculations. Here
	// we're getting cryptographic proofs for those slots!
	proof, err := client.GetProof(ctx, cfg.Account, slotStrings, cfg.BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("get proof: %w", err)
	}

	// Proof response validation: Even if the RPC call succeeds, the response
	// might be nil (though this shouldn't happen with well-behaved clients).
	// We validate to prevent nil pointer dereferences later.
	//
	// Building on previous patterns: Same validate-after-RPC pattern from modules
	// 01 and 11. Always check both error and nil response.
	if proof == nil {
		return nil, errors.New("nil proof response")
	}

	// ============================================================================
	// STEP 4: Build AccountProof with Defensive Copying
	// ============================================================================
	// We construct our result with defensive copies of all mutable data. This is
	// a critical security and correctness pattern!
	//
	// Why defensive copying?
	//   1. The proof response might contain pointers to client's internal data
	//   2. If we return those pointers directly, callers could mutate them
	//   3. This could affect other callers or cause data races in concurrent code
	//   4. By copying, we ensure each caller gets independent, isolated data
	//
	// What needs copying?
	//   - Balance (big.Int): Mutable type, must copy with new(big.Int).Set()
	//   - Nonce (uint64): Value type, automatically copied
	//   - CodeHash (common.Hash): Value type (array), automatically copied
	//   - StorageHash (common.Hash): Value type (array), automatically copied
	//   - ProofNodes ([]string): Slice (reference type), must copy
	//
	// Building on module 01: We learned defensive copying for big.Int there.
	// Here we extend it to slices. This pattern repeats throughout the course.
	res := &Result{
		Account: AccountProof{
			Balance:     new(big.Int).Set(proof.Balance),           // Defensive copy: big.Int is mutable
			Nonce:       proof.Nonce,                               // Direct copy: uint64 is value type
			CodeHash:    proof.CodeHash,                            // Direct copy: Hash is value type (array)
			StorageHash: proof.StorageHash,                         // Direct copy: Hash is value type (array)
			ProofNodes:  append([]string(nil), proof.AccountProof...), // Defensive copy: slice is reference type
		},
	}

	// ============================================================================
	// STEP 5: Process Storage Proofs - Understanding Nested Tries
	// ============================================================================
	// Storage proofs are nested inside account proofs! Here's how it works:
	//
	// 1. The account trie (global) contains all accounts
	//    - Root: stateRoot (in block header)
	//    - Leaves: Account objects (including storageHash)
	//
	// 2. Each contract has its own storage trie
	//    - Root: storageHash (in account object from step 1)
	//    - Leaves: Storage slot values
	//
	// So to prove a storage slot:
	//   1. Prove the account exists (account proof)
	//   2. Extract storageHash from account
	//   3. Prove the slot exists in that storage trie (storage proof)
	//
	// This two-level structure is why we get both account and storage proofs!
	//
	// Implementation: We iterate through storage proofs and build our own
	// StorageProof structs with defensive copies:
	//   - Value: big.Int (mutable) → must copy
	//   - Key: Convert from string to Hash (immutable after conversion)
	//   - ProofNodes: Slice (reference) → must copy
	//
	// Why check len(proof.StorageProof) > 0? If no storage slots were requested,
	// this will be empty. We only allocate and process if there's data.
	//
	// Connection to module 11: The slot keys here are the same slot calculations
	// we learned in module 11 (including mapping slot hashes). The proof proves
	// those calculations were correct!
	if len(proof.StorageProof) > 0 {
		res.Account.Storage = make([]StorageProof, 0, len(proof.StorageProof))
		for _, sp := range proof.StorageProof {
			// Copy value: big.Int is mutable
			value := new(big.Int).Set(sp.Value)

			// Convert key: string → Hash (creates new Hash value)
			key := common.HexToHash(sp.Key)

			// Copy proof nodes: slice is reference type
			proofNodes := append([]string(nil), sp.Proof...)

			// Append to result
			res.Account.Storage = append(res.Account.Storage, StorageProof{
				Key:        key,
				Value:      value,
				ProofNodes: proofNodes,
			})
		}
	}

	// ============================================================================
	// STEP 6: Return Complete Result
	// ============================================================================
	// We return a Result containing:
	//   - Account proof (balance, nonce, hashes, proof nodes)
	//   - Storage proofs (one per requested slot, with values and proof nodes)
	//
	// What can callers do with this?
	//   1. Verify the proofs against known state roots (trust-minimized verification)
	//   2. Store proofs for later verification (archival)
	//   3. Submit proofs to smart contracts (cross-chain bridges)
	//   4. Build light clients that verify state without full sync
	//
	// Building on previous concepts:
	//   - We validated inputs (Step 1) → now we return validated data
	//   - We handled errors consistently (Steps 1-3) → now we return success
	//   - We used defensive copying (Steps 4-5) → callers get isolated data
	//   - We used context for cancellation (Step 3) → operation completed successfully
	//
	// The progression across modules:
	//   - Module 01: Read block headers (lightweight state commitment)
	//   - Module 11: Read storage slots (raw state data)
	//   - Module 12: Get proofs (cryptographic verification of state)
	//   - Future: Verify proofs, build light clients, implement bridges
	return res, nil
}
