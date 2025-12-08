//go:build !solution
// +build !solution

package exercise

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// RPCClient is the tiny subset of ethclient we need for this module.
type RPCClient interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
}

// Config controls which block to fetch.
// If Number is nil, fetch the latest block.
// If IncludeTxs is true, summarize transactions; otherwise only header data is returned.
type Config struct {
	Number     *big.Int
	IncludeTxs bool
}

// TxSummary captures the minimal transaction details for explorer output.
type TxSummary struct {
	Hash common.Hash
	To   *common.Address
	Gas  uint64
}

// Result is what our explorer presents to callers.
type Result struct {
	Number   uint64
	Hash     common.Hash
	Parent   common.Hash
	TxCount  int
	Txs      []TxSummary
	GasUsed  uint64
	GasLimit uint64
}

// Run is the student entry point for module 14-explorer.
//
// TODOs (implement in exercise.go, see tests/solution for guidance):
// 1) Validate inputs: ctx non-nil, client non-nil.
// 2) Fetch the block via BlockByNumber (use cfg.Number, nil means latest).
// 3) Populate Result fields from the block header (number, hash, parent, gas).
// 4) If cfg.IncludeTxs, loop block.Transactions() and build TxSummary slice.
// 5) Return the Result (and bubble up any errors with context).
func Run(ctx context.Context, client RPCClient, cfg Config) (*Result, error) {
	panic("TODO: implement Run for 14-explorer")
}
