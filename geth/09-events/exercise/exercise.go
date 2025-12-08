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
	// TODO: Validate input parameters
	// - Check if ctx is nil and provide a default context if needed
	// - Check if client is nil and return an appropriate error
	// - Check if cfg.Token is the zero address and return an error
	// - Token address is required because we need to know which contract to query

	// TODO: Build FilterQuery for Transfer events
	// - Create ethereum.FilterQuery struct with:
	//   * Addresses: []common.Address{cfg.Token} (which contract to query)
	//   * FromBlock: cfg.FromBlock (start of range, nil for genesis)
	//   * ToBlock: cfg.ToBlock (end of range, nil for latest)
	//   * Topics: [][]common.Hash with Transfer event signature in first position
	// - Why Topics? Topics are indexed log parameters. Topic[0] is always the event signature hash.
	// - transferSigHash is already defined: keccak256("Transfer(address,address,uint256)")
	// - This filters for Transfer events from the specified token contract

	// TODO: Optionally filter by sender address (FromHolder)
	// - If cfg.FromHolder is provided (not nil):
	//   * Convert address to topic using addressTopic helper (already implemented)
	//   * Append to query.Topics as second element: []common.Hash{addressTopic(*cfg.FromHolder)}
	// - Why optional? Sometimes you want all transfers, sometimes only from specific address
	// - Topic[1] in Transfer event is the `from` address (first indexed parameter)

	// TODO: Optionally filter by recipient address (ToHolder)
	// - If cfg.ToHolder is provided (not nil):
	//   * If query.Topics has only 1 element (event signature), append nil placeholder for Topic[1]
	//   * Append addressTopic(*cfg.ToHolder) as Topic[2]
	// - Why nil placeholder? Topic positions matter. If filtering by `to` but not `from`, use nil for Topic[1]
	// - Topic[2] in Transfer event is the `to` address (second indexed parameter)
	// - Filter logic: Topics is [][]common.Hash where outer slice is topic position, inner slice is OR options

	// TODO: Execute log query
	// - Call client.FilterLogs(ctx, query) to fetch matching logs
	// - Handle errors (network failures, invalid block range)
	// - FilterLogs uses eth_getLogs RPC method
	// - Returns []types.Log (raw log entries with topics and data)

	// TODO: Initialize result with preallocated slice
	// - Create Result struct with Events slice
	// - Preallocate capacity: make([]TransferEvent, 0, len(logs))
	// - Why preallocate? Avoids slice reallocations during append loop
	// - This is a performance optimization for large result sets

	// TODO: Decode each log entry
	// - Loop through logs returned from FilterLogs
	// - For each log, call decodeTransferLog helper (already implemented below)
	// - decodeTransferLog extracts from/to/value from topics and data
	// - Handle decoding errors (malformed logs, unexpected structure)
	// - Append decoded TransferEvent to result.Events

	// TODO: Return result
	// - Return Result struct with all decoded events
	// - Return nil error on success
	// - Events are sorted by block number, then log index (order from blockchain)

	return nil, errors.New("not implemented")
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
