//go:build !solution && !reference

package events

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

/*
Problem: Query and decode ERC20 Transfer events from blockchain logs.
*/

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
	return nil
}

// decodeTransferLog - TODO: implement this function
func decodeTransferLog(lg types.Log) (TransferEvent, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil, nil
}

