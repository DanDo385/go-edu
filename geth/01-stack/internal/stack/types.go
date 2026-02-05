package stack

import (
	"context" // context for the RPC client
	"math/big" // big integer for the chain ID and network ID

	"github.com/ethereum/go-ethereum/core/types" // types of the header
) 

// RPCClient captures the ethclient calls needed for module 01.
type RPCClient interface {
	ChainID(ctx context.Context) (*big.Int, error) // get the chain ID from the context
	NetworkID(ctx context.Context) (*big.Int, error) // get the network ID from the context
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) // get the header by number from the context
}

// Config allows overriding which block header is fetched (nil => latest).
type Config struct { // get the block number from the config
	BlockNumber *big.Int // block number to get the header from
}

// Result summarizes the Ethereum stack data retrieved from the node.
type Result struct {
	ChainID   *big.Int
	NetworkID *big.Int
	Header    *types.Header
}
