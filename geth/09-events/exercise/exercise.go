//go:build !solution
// +build !solution

package exercise

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

// Run contains the reference solution for module 09-events.
func Run(ctx context.Context, client LogClient, cfg Config) (*Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if cfg.Token == (common.Address{}) {
		return nil, errors.New("token address required")
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{cfg.Token},
		FromBlock: cfg.FromBlock,
		ToBlock:   cfg.ToBlock,
		Topics:    [][]common.Hash{{transferSigHash}},
	}

	if cfg.FromHolder != nil {
		query.Topics = append(query.Topics, []common.Hash{addressTopic(*cfg.FromHolder)})
	}
	if cfg.ToHolder != nil {
		if len(query.Topics) == 1 {
			query.Topics = append(query.Topics, nil)
		}
		query.Topics = append(query.Topics, []common.Hash{addressTopic(*cfg.ToHolder)})
	}

	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("filter logs: %w", err)
	}

	result := &Result{
		Events: make([]TransferEvent, 0, len(logs)),
	}

	for _, lg := range logs {
		event, err := decodeTransferLog(lg)
		if err != nil {
			return nil, err
		}
		result.Events = append(result.Events, event)
	}
	return result, nil
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
