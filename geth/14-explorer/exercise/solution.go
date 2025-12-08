//go:build solution
// +build solution

package exercise

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

// Run contains the reference implementation for the tiny block explorer.
//
// Concepts reinforced:
// - Reusing interfaces from earlier modules (RPC calls, context for cancellation)
// - Header vs full block fetching (here we pull full blocks with tx objects)
// - Defensive copies of pointer-heavy types (types.Block shares internals)
func Run(ctx context.Context, client RPCClient, cfg Config) (*Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if client == nil {
		return nil, errors.New("client is nil")
	}

	block, err := client.BlockByNumber(ctx, cfg.Number)
	if err != nil {
		target := "latest"
		if cfg.Number != nil {
			target = cfg.Number.String()
		}
		return nil, fmt.Errorf("fetch block %s: %w", target, err)
	}
	if block == nil {
		return nil, errors.New("nil block response")
	}

	header := block.Header()
	res := &Result{
		Number:   header.Number.Uint64(),
		Hash:     block.Hash(),
		Parent:   header.ParentHash,
		TxCount:  len(block.Transactions()),
		GasUsed:  header.GasUsed,
		GasLimit: header.GasLimit,
	}

	if cfg.IncludeTxs {
		res.Txs = make([]TxSummary, 0, len(block.Transactions()))
		for _, tx := range block.Transactions() {
			to := tx.To() // pointer (may be nil for contract creation)
			res.Txs = append(res.Txs, TxSummary{
				Hash: tx.Hash(),
				To:   to,
				Gas:  tx.Gas(),
			})
		}
	}

	// types.Block shares slices internally; copy header fields above and only
	// expose minimal TxSummary to avoid leaking pointers to caller mutations.
	return res, nil
}
