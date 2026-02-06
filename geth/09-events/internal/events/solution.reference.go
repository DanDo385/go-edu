//go:build reference

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

var (
	errNilClient    = errors.New("nil log client")
	transferSigHash = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
)

/*
Reference Solution

Structure:
- Build a FilterQuery for Transfer logs.
- Query logs once.
- Decode each log into a typed TransferEvent.

Invariants:
- Topic[0] must be Transfer signature.
- Topic[1]/Topic[2] are indexed from/to addresses.
*/
func Run(ctx context.Context, client LogClient, cfg Config) (*Result, error) {
	if client == nil {
		return nil, errNilClient
	}
	if cfg.Token == (common.Address{}) {
		return nil, errors.New("token address is required")
	}

	topics := [][]common.Hash{{transferSigHash}}
	if cfg.FromHolder != nil || cfg.ToHolder != nil {
		if cfg.FromHolder != nil {
			topics = append(topics, []common.Hash{addressTopic(*cfg.FromHolder)})
		} else {
			topics = append(topics, nil)
		}
		if cfg.ToHolder != nil {
			topics = append(topics, []common.Hash{addressTopic(*cfg.ToHolder)})
		}
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{cfg.Token},
		FromBlock: cfg.FromBlock,
		ToBlock:   cfg.ToBlock,
		Topics:    topics,
	}

	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("filter logs: %w", err)
	}

	events := make([]TransferEvent, 0, len(logs))
	for i, lg := range logs {
		ev, err := decodeTransferLog(lg)
		if err != nil {
			return nil, fmt.Errorf("decode log %d: %w", i, err)
		}
		events = append(events, ev)
	}

	return &Result{Events: events}, nil
}

func addressTopic(addr common.Address) common.Hash {
	return common.BytesToHash(common.LeftPadBytes(addr.Bytes(), 32))
}

func decodeTransferLog(lg types.Log) (TransferEvent, error) {
	if len(lg.Topics) < 3 {
		return TransferEvent{}, errors.New("transfer log requires at least 3 topics")
	}
	if lg.Topics[0] != transferSigHash {
		return TransferEvent{}, errors.New("unexpected event signature")
	}
	if len(lg.Data) < 32 {
		return TransferEvent{}, errors.New("transfer value payload too short")
	}

	from := common.BytesToAddress(lg.Topics[1].Bytes()[12:])
	to := common.BytesToAddress(lg.Topics[2].Bytes()[12:])
	value := new(big.Int).SetBytes(lg.Data[:32])

	return TransferEvent{
		BlockNumber: lg.BlockNumber,
		TxHash:      lg.TxHash,
		LogIndex:    lg.Index,
		From:        from,
		To:          to,
		Value:       value,
	}, nil
}
