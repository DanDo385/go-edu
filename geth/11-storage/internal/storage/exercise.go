//go:build !solution && !reference

package storage

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

/*
Problem: Read raw storage slots from Ethereum contracts, including mapping slots.
*/

var zeroHash = common.Hash{}

// Run - TODO: implement this function
func Run(ctx context.Context, client StorageClient, cfg Config) (*Result, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	var zero0 *Result
	var zero1 error
	return zero0, zero1
}

// slotToHash - TODO: implement this function
func slotToHash(slot *big.Int) common.Hash {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	var zero0 common.Hash
	return zero0
}

// mappingSlotHash - TODO: implement this function
func mappingSlotHash(key []byte, slot common.Hash) common.Hash {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	var zero0 common.Hash
	return zero0
}
