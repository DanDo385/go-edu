//go:build !solution && !reference

package events

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

/*
Problem: Query and decode ERC20 Transfer events from blockchain logs.
*/

var transferSigHash = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

// Run - TODO: implement this function
func Run(ctx context.Context, client LogClient, cfg Config) (*Result, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	var zero0 *Result
	var zero1 error
	return zero0, zero1
}

// addressTopic - TODO: implement this function
func addressTopic(addr common.Address) common.Hash {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	var zero0 common.Hash
	return zero0
}

// decodeTransferLog - TODO: implement this function
func decodeTransferLog(lg types.Log) (TransferEvent, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	var zero0 TransferEvent
	var zero1 error
	return zero0, zero1
}
