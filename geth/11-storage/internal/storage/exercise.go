//go:build !solution
// +build !solution

package storage

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
	// TODO: Implement Run
	// See solution.reference.go for reference implementation
	panic("not implemented")
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
