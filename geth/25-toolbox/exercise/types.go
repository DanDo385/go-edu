package exercise

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

// ToolboxClient captures all the ethclient calls needed for the toolbox (combines all previous modules).
type ToolboxClient interface {
	// From module 01: chain metadata
	ChainID(ctx context.Context) (*big.Int, error)
	NetworkID(ctx context.Context) (*big.Int, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)

	// From module 21: sync status
	SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error)

	// From module 22: peer count
	PeerCount(ctx context.Context) (uint64, error)

	// From module 24: block/tx lookups
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	TransactionByHash(ctx context.Context, hash string) (*types.Transaction, bool, error)
}

// Config allows configuration for different toolbox commands.
type Config struct {
	// Command specifies which operation to perform (status, block, tx, etc.)
	Command string
	// Args contains command-specific arguments
	Args []string
}

// Result contains the output from a toolbox command.
type Result struct {
	// Command that was executed
	Command string
	// Output contains the result data (structure varies by command)
	Output interface{}
	// Status indicates success/failure
	Status string
}
