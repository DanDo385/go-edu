package exercise

import (
	"context"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

// ContractCaller matches the subset of methods bound contracts rely on for read calls.
type ContractCaller interface {
	CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)
	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
}

// Config drives the typed binding demo.
type Config struct {
	ABI         string
	Contract    common.Address
	Holder      *common.Address
	BlockNumber *big.Int
}

// Result returns token metadata + optional balance.
type Result struct {
	Name        string
	Symbol      string
	Decimals    uint8
	TotalSupply *big.Int
	Balance     *big.Int
}
