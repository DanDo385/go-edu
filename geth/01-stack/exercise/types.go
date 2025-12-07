package exercise

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

// RPCClient captures the ethclient calls needed for module 01.
type RPCClient interface {
	ChainID(ctx context.Context) (*big.Int, error)
	NetworkID(ctx context.Context) (*big.Int, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

// Config allows overriding which block header is fetched (nil => latest).
type Config struct {
	BlockNumber *big.Int
}

// Result summarizes the Ethereum stack data retrieved from the node.
type Result struct {
	ChainID   *big.Int
	NetworkID *big.Int
	Header    *types.Header
}
