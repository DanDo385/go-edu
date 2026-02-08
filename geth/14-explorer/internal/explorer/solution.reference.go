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
Reference Solution - Block Explorer (eth_getBlockByNumber)
==========================================================

This file demonstrates fetching block metadata for explorer-style UIs: number,
hash, parent, transaction count, gas used/limit. Optionally includes tx summaries.

This connects to the Ethereum ecosystem by showing:
- BlockByNumber(ctx, number): nil = latest, or *big.Int for specific block
- block.NumberU64(), Hash(), ParentHash(), Transactions(), GasUsed, GasLimit
- tx.To(): returns *common.Address; nil for contract creation txs
- Defensive copy: v := *to; toCopy = &v for To field so we own the address

The exercise builds understanding of:
- Block structure: header fields + transaction list
- IncludeTxs: when true, build TxSummary slice with Hash, To, Gas
- Pointer copy: tx.To() may point to internal; we copy before storing

Teaching notes (per .cursorrules):
- tx.To() returns *common.Address; nil means contract creation. We copy via
  v := *to; toCopy = &v so Result doesn't share pointers with block internals.
- block is from RPC; we extract value types and copy references we store.
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
