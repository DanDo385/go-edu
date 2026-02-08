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
Reference Solution - Contract Event Logs (Transfer)
==================================================

This file demonstrates querying ERC20 Transfer event logs via eth_getLogs.
Events are emitted as logs: topic0 = event signature hash, topics[1..n] =
indexed args, Data = non-indexed args. Transfer(address,address,uint256) has
indexed from/to and non-indexed value.

This connects to the Ethereum ecosystem by showing:
- Keccak256Hash of full signature: topic0 identifies the event type
- FilterQuery: Addresses, FromBlock, ToBlock, Topics (OR within each, AND across)
- topics[1], topics[2]: indexed addresses left-padded to 32 bytes; addr in last 20
- log.Data: ABI-encoded non-indexed value (32 bytes for uint256)

The exercise builds understanding of:
- Event layout: indexed vs non-indexed (topics vs Data)
- addressTopic: LeftPadBytes(addr.Bytes(), 32) then BytesToHash for topic format
- Optional filters: cfg.FromHolder, cfg.ToHolder add topic1/topic2 when set
- Defensive copy: new(big.Int).SetBytes for value

Teaching notes (per .cursorrules):
- topics layout: [topic0, topic1?, topic2?] — nil in a slot means "any" for that
  indexed param. topics = append(topics, nil) for "match topic0 only" when
  filtering by one holder.
- Slice sharing: FilterQuery holds references; we copy decoded values into
  TransferEvent, not the raw log.
*/
func Run(ctx context.Context, client LogClient, cfg Config) (*Result, error) {
	if client == nil {
		return nil, errNilClient
	}
	if cfg.Token == (common.Address{}) {
		return nil, errors.New("token address is required")
	}

	// topic0 = event sig; topic1/topic2 = indexed from/to when filtering
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

// addressTopic encodes address as 32-byte topic: left-pad to 32, then BytesToHash.
func addressTopic(addr common.Address) common.Hash {
	return common.BytesToHash(common.LeftPadBytes(addr.Bytes(), 32))
}

// decodeTransferLog extracts from (topics[1]), to (topics[2]), value (Data[:32]) from a Transfer log.
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
