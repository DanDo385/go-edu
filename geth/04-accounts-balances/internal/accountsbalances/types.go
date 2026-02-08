package accountsbalances

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// AccountClient captures the RPC methods required for this module.
type AccountClient interface {
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
}

// AccountType represents the coarse account classification.
type AccountType string

const (
	AccountTypeEOA      AccountType = "EOA"
	AccountTypeContract AccountType = "Contract"
)

// Config controls how Run performs lookups.
type Config struct {
	// Address is the single address to query (convenience field).
	Address     common.Address
	Addresses   []common.Address
	BlockNumber *big.Int
}

// AccountState captures details for a single address.
type AccountState struct {
	Address common.Address
	Balance *big.Int
	Code    []byte
	Type    AccountType
}

// Result aggregates the queried account states.
type Result struct {
	Accounts []AccountState
	// Convenience fields for single-address queries.
	Balance *big.Int
	Nonce   uint64
}
