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

func Run(ctx context.Context, client LogClient, cfg Config) (*Result, error) {
	// TODO: Input Validation
	// TODO: Build FilterQuery
	// TODO: Optionally Filter by Sender Address
	// TODO: Optionally Filter by Recipient Address
	// TODO: Execute Log Query
	// TODO: Initialize Result with Preallocated Slice
	// TODO: Decode Each Log
	// TODO: Return Results
	panic("not implemented")
}

func addressTopic(addr common.Address) common.Hash {
	// TODO: Implement this function
	panic("not implemented")
}

func decodeTransferLog(lg types.Log) (TransferEvent, error) {
	// TODO: Implement this function
	panic("not implemented")
}
