//go:build !solution && !reference

package events

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

/*
Problem: Query and decode ERC20 Transfer events from blockchain logs.
*/

// transferSigHash is the Keccak256 hash of the Transfer(address,address,uint256) event signature.
// TODO: compute via crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
var transferSigHash = common.Hash{}

// Run - TODO: implement this function
func Run(ctx context.Context, client LogClient, cfg Config) (*Result, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil, nil
}

// addressTopic - TODO: implement this function
func addressTopic(addr common.Address) common.Hash {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return common.Hash{}
}

// decodeTransferLog - TODO: implement this function
func decodeTransferLog(lg types.Log) (TransferEvent, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return TransferEvent{}, nil
}
