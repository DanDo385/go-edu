package smartcontracts

import (
	"context"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

// ContractCaller exposes the method needed for making contract calls.
// This interface is satisfied by *ethclient.Client.
type ContractCaller interface {
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
}

// Config specifies which contract to query.
type Config struct {
	// Contract is the address of the ERC20 token contract.
	Contract common.Address

	// BlockNumber specifies which block to query. nil means "latest".
	BlockNumber *big.Int
}

// Result captures decoded ERC-20 metadata.
type Result struct {
	// Name is the token name (e.g., "USD Coin")
	Name string

	// Symbol is the token symbol (e.g., "USDC")
	Symbol string

	// Decimals is the number of decimal places (e.g., 6 for USDC, 18 for DAI)
	Decimals uint8

	// TotalSupply is the total number of tokens in existence
	TotalSupply *big.Int
}
