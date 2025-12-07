package exercise

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// TXClient captures just enough of ethclient.Client for this exercise.
type TXClient interface {
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	ChainID(ctx context.Context) (*big.Int, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
}

// Config defines how to build the transaction.
type Config struct {
	PrivateKey *ecdsa.PrivateKey
	To         common.Address
	AmountWei  *big.Int

	// Optional overrides:
	Nonce    *uint64
	GasPrice *big.Int
	GasLimit uint64
	Data     []byte
	NoSend   bool // useful in tests/tutorials where you only want the signed tx
}

// Result surfaces the signed transaction plus metadata for display.
type Result struct {
	FromAddress common.Address
	Nonce       uint64
	Tx          *types.Transaction
}
