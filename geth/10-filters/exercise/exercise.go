//go:build !solution
// +build !solution

package exercise

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const defaultMaxHeads = 5
const defaultPollInterval = time.Second

// Run contains the reference solution for module 10-filters.
func Run(ctx context.Context, client HeadClient, cfg Config) (*Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if cfg.MaxHeads <= 0 {
		cfg.MaxHeads = defaultMaxHeads
	}
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = defaultPollInterval
	}

	if cfg.PollMode {
		return pollHeads(ctx, client, cfg)
	}
	return subscribeHeads(ctx, client, cfg)
}

func subscribeHeads(ctx context.Context, client HeadClient, cfg Config) (*Result, error) {
	headCh := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(ctx, headCh)
	if err != nil {
		return nil, fmt.Errorf("subscribe new head: %w", err)
	}
	defer sub.Unsubscribe()

	result := &Result{
		Heads: make([]HeadInfo, 0, cfg.MaxHeads),
		Mode:  "subscription",
	}

	var prevHash common.Hash
	for len(result.Heads) < cfg.MaxHeads {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context canceled: %w", ctx.Err())
		case err := <-sub.Err():
			if err != nil {
				return nil, fmt.Errorf("subscription error: %w", err)
			}
		case head := <-headCh:
			if head == nil {
				continue
			}
			hash := head.Hash()
			reorg := (prevHash != (common.Hash{})) && (head.ParentHash != prevHash)
			result.Heads = append(result.Heads, HeadInfo{
				Number:     head.Number.Uint64(),
				Hash:       hash,
				ParentHash: head.ParentHash,
				Reorg:      reorg,
			})
			prevHash = hash
		}
	}
	return result, nil
}

func pollHeads(ctx context.Context, client HeadClient, cfg Config) (*Result, error) {
	result := &Result{
		Heads: make([]HeadInfo, 0, cfg.MaxHeads),
		Mode:  "polling",
	}
	var prevHash common.Hash
	var prevNumber uint64

	for len(result.Heads) < cfg.MaxHeads {
		head, err := client.HeaderByNumber(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("header by number: %w", err)
		}
		if head == nil {
			return nil, errors.New("received nil header")
		}
		number := head.Number.Uint64()
		hash := head.Hash()

		if number == prevNumber {
			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("context canceled: %w", ctx.Err())
			case <-time.After(cfg.PollInterval):
			}
			continue
		}
		reorg := (prevHash != (common.Hash{})) && (head.ParentHash != prevHash)
		result.Heads = append(result.Heads, HeadInfo{
			Number:     number,
			Hash:       hash,
			ParentHash: head.ParentHash,
			Reorg:      reorg,
		})
		prevHash = hash
		prevNumber = number
	}
	return result, nil
}
