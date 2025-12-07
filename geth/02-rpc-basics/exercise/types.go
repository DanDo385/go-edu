package exercise

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

// RPCClient captures the go-ethereum methods needed for this exercise.
type RPCClient interface {
	NetworkID(ctx context.Context) (*big.Int, error)
	BlockNumber(ctx context.Context) (uint64, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
}

// Config configures retry behavior.
type Config struct {
	Retries int
}

// Result summarizes the information fetched from the RPC endpoint.
type Result struct {
	NetworkID   *big.Int
	BlockNumber uint64
	Block       *types.Block
}
