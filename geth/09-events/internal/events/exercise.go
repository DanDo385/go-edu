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

var transferSigHash = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

/*
Problem: Query and decode ERC20 Transfer events from blockchain logs.

This module teaches you how to work with Ethereum events/logs. Events are append-only
records emitted during contract execution. They're searchable via bloom filters and
provide an efficient way to track state changes without querying contract storage.

Computer science principles highlighted:
  - Event-driven architecture: Logs as append-only audit trail
  - Bloom filters: Probabilistic data structures for efficient searching
  - Indexed vs non-indexed parameters: Trade-off between searchability and cost
  - Log structure: Topics (indexed) vs Data (non-indexed)
*/
func Run(ctx context.Context, client LogClient, cfg Config) (*Result, error) {
	// TODO: Implement Run
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


func addressTopic(addr common.Address) common.Hash {
	return common.BytesToHash(common.LeftPadBytes(addr.Bytes(), 32))
}

func decodeTransferLog(lg types.Log) (TransferEvent, error) {
	if len(lg.Topics) < 3 {
		return TransferEvent{}, fmt.Errorf("log %s missing topics", lg.TxHash.Hex())
	}
	if lg.Topics[0] != transferSigHash {
		return TransferEvent{}, fmt.Errorf("unexpected topic %s", lg.Topics[0].Hex())
	}

	from := common.BytesToAddress(lg.Topics[1].Bytes()[12:])
	to := common.BytesToAddress(lg.Topics[2].Bytes()[12:])
	if len(lg.Data) < 32 {
		return TransferEvent{}, fmt.Errorf("log %s data too short", lg.TxHash.Hex())
	}
	value := new(big.Int).SetBytes(lg.Data[len(lg.Data)-32:])

	return TransferEvent{
		BlockNumber: lg.BlockNumber,
		TxHash:      lg.TxHash,
		LogIndex:    lg.Index,
		From:        from,
		To:          to,
		Value:       value,
	}, nil
}
