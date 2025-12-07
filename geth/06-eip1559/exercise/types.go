package exercise

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// FeeClient abstracts the ethclient methods needed for dynamic fee tx building.
type FeeClient interface {
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
	ChainID(ctx context.Context) (*big.Int, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
}

// Config controls how Run constructs the transaction.
type Config struct {
	PrivateKey *ecdsa.PrivateKey
	To         common.Address
	AmountWei  *big.Int

	Nonce          *uint64
	GasLimit       uint64
	MaxPriorityFee *big.Int
	MaxFee         *big.Int
	Data           []byte
	NoSend         bool
	BlockNumber    *big.Int // optional header to sample fees from (nil => latest)
}

// Result surfaces useful transaction metadata.
type Result struct {
	FromAddress common.Address
	Nonce       uint64
	Tx          *types.Transaction
	BaseFee     *big.Int
}
