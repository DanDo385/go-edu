//go:build reference

package receipts

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var errNilClient = errors.New("nil receipt client")

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

func Run(ctx context.Context, client ReceiptClient, cfg Config) (*Result, error) {
	if client == nil {
		return nil, errNilClient
	}
	if cfg.TxHash == (common.Hash{}) {
		return nil, errors.New("transaction hash is required")
	}

	rcpt, err := client.TransactionReceipt(ctx, cfg.TxHash)
	if err != nil {
		return nil, fmt.Errorf("transaction receipt: %w", err)
	}
	if rcpt == nil {
		return nil, errors.New("transaction receipt: nil response")
	}

	logs := make([]LogSummary, 0, len(rcpt.Logs))
	for _, lg := range rcpt.Logs {
		if lg == nil {
			continue
		}
		logs = append(logs, LogSummary{
			Address: lg.Address,
			Topics:  append([]common.Hash(nil), lg.Topics...),
			Data:    append([]byte(nil), lg.Data...),
			Index:   lg.Index,
		})
	}

	var blockNumber *big.Int
	if rcpt.BlockNumber != nil {
		blockNumber = new(big.Int).Set(rcpt.BlockNumber)
	}

	return &Result{
		TxHash:        cfg.TxHash,
		BlockNumber:   blockNumber,
		StatusOK:      rcpt.Status == types.ReceiptStatusSuccessful,
		GasUsed:       rcpt.GasUsed,
		CumulativeGas: rcpt.CumulativeGasUsed,
		Contract:      rcpt.ContractAddress,
		Logs:          logs,
		PostStateRoot: append([]byte(nil), rcpt.PostState...),
	}, nil
}
