//go:build !solution
// +build !solution

package exercise

import (
	"context"
	"errors"
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
func Run(ctx context.Context, client RPCClient, cfg Config) (*Result, error) {
	// TODO: Validate input parameters
	// - Check if ctx is nil and provide context.Background() as default
	// - Check if client is nil and return an appropriate error
	// Why validate? Block fetching is an RPC call; validate before network operations

	// TODO: Fetch the block via BlockByNumber
	// - Call client.BlockByNumber(ctx, cfg.Number)
	// - cfg.Number can be nil (means latest block) or a specific block number
	// - Handle potential errors from the RPC call
	// - Validate that the block response is not nil

	// TODO: Extract header and populate Result with basic block info
	// - Get the block header with block.Header()
	// - Create a Result struct with:
	//   - Number: header.Number.Uint64()
	//   - Hash: block.Hash()
	//   - Parent: header.ParentHash
	//   - TxCount: len(block.Transactions())
	//   - GasUsed: header.GasUsed
	//   - GasLimit: header.GasLimit

	// TODO: Optionally include transaction summaries
	// - Check if cfg.IncludeTxs is true
	// - If yes, create a Txs slice with capacity len(block.Transactions())
	// - Loop through block.Transactions() and for each tx:
	//   - Create a TxSummary with Hash (tx.Hash()), To (tx.To()), Gas (tx.Gas())
	//   - Note: tx.To() returns *common.Address (pointer, can be nil for contract creation)
	//   - Append to result.Txs
	// Why optional? Transaction details can be large; only include if requested

	// TODO: Return the completed Result
	// - Return the Result struct and nil error on success

	return nil, errors.New("not implemented")
}
