package explorer

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// RPCClient captures the ethclient calls needed for the explorer module.
type RPCClient interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
}

// Config controls which block is fetched and whether transactions are included.
type Config struct {
	// Number is nil for latest, or a specific block number.
	Number *big.Int
	// IncludeTxs controls whether transaction summaries are populated.
	IncludeTxs bool
}

// TxSummary is a lightweight view of a transaction.
type TxSummary struct {
	Hash common.Hash
	To   *common.Address
	Gas  uint64
}

// Result is a minimal block-explorer view of a block.
type Result struct {
	Number   uint64
	Hash     common.Hash
	Parent   common.Hash
	TxCount  int
	GasUsed  uint64
	GasLimit uint64
	Txs      []TxSummary
}
