//go:build !solution && !reference

package storage

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

/*
Problem: Read raw storage slots from Ethereum contracts, including mapping slots.
*/

// Run - TODO: implement this function
func Run(ctx context.Context, client StorageClient, cfg Config) (*Result, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil, nil
}

// slotToHash - TODO: implement this function
func slotToHash(slot *big.Int) common.Hash {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil
}

// mappingSlotHash - TODO: implement this function
func mappingSlotHash(key []byte, slot common.Hash) common.Hash {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil
}

