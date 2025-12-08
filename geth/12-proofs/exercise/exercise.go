//go:build !solution
// +build !solution

package exercise

import (
	"context"
	"errors"
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
	// TODO: Validate input parameters
	// - Check if ctx is nil and provide a default context if needed
	// - Check if client is nil and return an appropriate error
	// - Validate that cfg.Account is not the zero address
	// Why validate? Proof RPC calls are expensive; fail fast on bad inputs

	// TODO: Convert storage slots to hex strings
	// - Check if cfg.Slots has length > 0
	// - Create a slotStrings slice with capacity len(cfg.Slots)
	// - Loop through cfg.Slots and convert each common.Hash to hex string using .Hex()
	// - The GetProof RPC method expects slot identifiers as hex strings
	// Why convert? The RPC interface uses strings, not Hash types directly

	// TODO: Call GetProof to fetch account and storage proofs
	// - Call client.GetProof(ctx, cfg.Account, slotStrings, cfg.BlockNumber)
	// - This returns an AccountResult pointer containing:
	//   - Account state (balance, nonce, codeHash, storageHash)
	//   - Account proof nodes (Merkle path from state root to account)
	//   - Storage proof nodes (Merkle path from storage root to each slot)
	// - Handle potential errors from the RPC call
	// - Validate that the proof response is not nil

	// TODO: Build the Result with defensive copying
	// - Create an AccountProof struct with:
	//   - Balance: Copy using new(big.Int).Set(proof.Balance)
	//   - Nonce: Direct copy (uint64 is a value type)
	//   - CodeHash: Direct copy (common.Hash is a value type)
	//   - StorageHash: Direct copy (common.Hash is a value type)
	//   - ProofNodes: Copy slice using append([]string(nil), proof.AccountProof...)
	// Why defensive copying? Prevent caller mutations from affecting our data

	// TODO: Process storage proofs if present
	// - Check if len(proof.StorageProof) > 0
	// - Create a slice for Storage proofs with capacity len(proof.StorageProof)
	// - Loop through proof.StorageProof and for each StorageResult:
	//   - Copy Value using new(big.Int).Set(sp.Value)
	//   - Convert Key string to common.Hash using common.HexToHash(sp.Key)
	//   - Copy ProofNodes using append([]string(nil), sp.Proof...)
	//   - Append to the Storage slice
	// Why process these? Storage proofs prove specific slot values in the contract

	// TODO: Return the complete Result
	// - Create a Result struct with the Account field set to the AccountProof
	// - Return the result and nil error on success

	return nil, errors.New("not implemented")
}
