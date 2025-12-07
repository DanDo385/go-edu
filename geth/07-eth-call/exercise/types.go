package exercise

import (
	"context"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

// CallClient exposes the single method we need from ethclient.
type CallClient interface {
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
}

// Config specifies which contract to query.
type Config struct {
	Contract    common.Address
	BlockNumber *big.Int
}

// Result captures decoded ERC-20 metadata.
type Result struct {
	Name        string
	Symbol      string
	Decimals    uint8
	TotalSupply *big.Int
}
