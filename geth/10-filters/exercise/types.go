package exercise

import (
	"context"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// HeadClient captures the subset of ethclient used for monitoring heads.
type HeadClient interface {
	SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

// Config tunes the monitoring strategy.
type Config struct {
	MaxHeads     int
	PollInterval time.Duration
	PollMode     bool // if true, use HTTP polling; otherwise prefer ws subscription
}

// HeadInfo contains the metadata surfaced to the CLI/tests.
type HeadInfo struct {
	Number     uint64
	Hash       common.Hash
	ParentHash common.Hash
	Reorg      bool
}

// Result aggregates collected head info.
type Result struct {
	Heads []HeadInfo
	Mode  string // "subscription" or "polling"
}
