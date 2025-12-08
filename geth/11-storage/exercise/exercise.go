//go:build !solution
// +build !solution

package exercise

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var zeroHash = common.Hash{}

/*
Problem: Read raw storage slots from Ethereum contracts, including mapping slots.

Storage is the cryptographic database where contracts store their persistent state.
Every contract has 2^256 possible 32-byte slots organized as a Merkle-Patricia trie.
Understanding storage layout is essential for:
  - Debugging contract state
  - Building indexers that track specific contract data
  - Verifying storage proofs (module 12)
  - Optimizing gas costs (packed storage)

Computer science principles highlighted:
  - Cryptographic commitment via Merkle trees (storage root commits to all slots)
  - Deterministic slot calculation (mapping slots via keccak256 hash)
  - Key-value store abstraction (2^256 address space maps to 32-byte values)
*/
func Run(ctx context.Context, client StorageClient, cfg Config) (*Result, error) {
	// TODO: Validate input parameters
	// - Check if ctx is nil and provide a default context if needed
	// - Check if client is nil and return an appropriate error
	// - Validate that cfg.Contract is not the zero address
	// - Validate that cfg.Slot is not nil (slot number is required)
	// Why validate? Storage reads are expensive RPC calls; fail fast on bad inputs

	// TODO: Convert slot number to storage hash
	// - Call slotToHash(cfg.Slot) to convert big.Int slot to common.Hash
	// - Storage slots are identified by 32-byte hashes, not integers directly
	// - This conversion pads the slot number to 32 bytes (big-endian)

	// TODO: Handle mapping key if provided
	// - Check if cfg.MappingKey has length > 0
	// - If yes, call mappingSlotHash(cfg.MappingKey, slotHash) to compute the mapped slot
	// - Mapping slots use: keccak256(abi.encode(key, baseSlot))
	// - This is how Solidity mappings work under the hood!
	// Why? Mappings don't store data in the base slot; they hash the key to find the actual slot

	// TODO: Read the storage value from the contract
	// - Call client.StorageAt(ctx, cfg.Contract, slotHash, cfg.BlockNumber)
	// - Use the resolved slotHash (either direct or mapping-derived)
	// - cfg.BlockNumber can be nil for latest block
	// - Handle potential errors from the RPC call
	// - StorageAt returns []byte (raw 32-byte value from storage)

	// TODO: Construct and return the Result
	// - Create a Result struct with:
	//   - ResolvedSlot: the final slot hash that was queried
	//   - Value: the raw bytes returned from storage
	// - Why track ResolvedSlot? For mapping queries, this shows the computed slot location
	// - Return the result and nil error on success

	return nil, errors.New("not implemented")
}

func slotToHash(slot *big.Int) common.Hash {
	if slot == nil {
		return zeroHash
	}
	return common.BigToHash(slot)
}

func mappingSlotHash(key []byte, slot common.Hash) common.Hash {
	keyPadded := common.LeftPadBytes(key, 32)
	data := append(keyPadded, slot.Bytes()...)
	return common.BytesToHash(crypto.Keccak256(data))
}
