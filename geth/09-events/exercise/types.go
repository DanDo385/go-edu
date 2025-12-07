package exercise

import (
	"context"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// LogClient captures the subset of ethclient capabilities needed to filter logs.
type LogClient interface {
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
}

// Config configures the ERC20 Transfer log query.
type Config struct {
	Token      common.Address
	FromBlock  *big.Int
	ToBlock    *big.Int
	FromHolder *common.Address
	ToHolder   *common.Address
}

// TransferEvent represents a decoded Transfer log.
type TransferEvent struct {
	BlockNumber uint64
	TxHash      common.Hash
	LogIndex    uint
	From        common.Address
	To          common.Address
	Value       *big.Int
}

// Result aggregates decoded events.
type Result struct {
	Events []TransferEvent
}
