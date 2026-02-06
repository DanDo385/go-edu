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

Structure:
- Fetch one block.
- Build a stable, lightweight summary object.
- Optionally include per-transaction summaries.

Pointer notes:
  - `tx.To()` returns `*common.Address`; we copy that address value before storing
    it in `TxSummary` so result data does not alias transaction internals.
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
