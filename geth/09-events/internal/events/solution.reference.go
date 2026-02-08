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
==================

This file is the canonical reference for this exercise. It keeps failure paths
explicit when an operation can fail, so callers can decide how to handle
errors at API boundaries.

Read this alongside exercise.go and the tests to understand the intended data
flow, ownership boundaries, and invariants that keep behavior deterministic.

Teaching notes:
- Memory/ownership: make copies when returning mutable data that should not
  alias internal state; share references only when aliasing is intentional.
- Invariants: establish assumptions close to construction, and rely on them in
  smaller helper functions to keep logic easy to audit.
- Error surfaces: prefer explicit returns over hidden panics so learners can
  reason about control flow in production-style code.
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

// addressTopic implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func addressTopic(addr common.Address) common.Hash {
	return common.BytesToHash(common.LeftPadBytes(addr.Bytes(), 32))
}

// decodeTransferLog implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
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
