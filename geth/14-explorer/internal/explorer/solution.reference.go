//go:build reference

package explorer

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

var errNilClient = errors.New("nil explorer rpc client")

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

func Run(ctx context.Context, client RPCClient, cfg Config) (*Result, error) {
	if client == nil {
		return nil, errNilClient
	}

	block, err := client.BlockByNumber(ctx, cfg.Number)
	if err != nil {
		return nil, fmt.Errorf("block by number: %w", err)
	}
	if block == nil {
		return nil, errors.New("block by number: nil response")
	}

	res := &Result{
		Number:   block.NumberU64(),
		Hash:     block.Hash(),
		Parent:   block.ParentHash(),
		TxCount:  len(block.Transactions()),
		GasUsed:  block.GasUsed(),
		GasLimit: block.GasLimit(),
	}

	if cfg.IncludeTxs {
		res.Txs = make([]TxSummary, 0, len(block.Transactions()))
		for _, tx := range block.Transactions() {
			var toCopy *common.Address
			if to := tx.To(); to != nil {
				v := *to
				toCopy = &v
			}
			res.Txs = append(res.Txs, TxSummary{
				Hash: tx.Hash(),
				To:   toCopy,
				Gas:  tx.Gas(),
			})
		}
	}

	return res, nil
}
