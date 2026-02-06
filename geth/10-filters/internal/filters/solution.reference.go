//go:build reference

package filters

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

var errNilClient = errors.New("nil head client")

/*
Reference Solution

Structure:
- Resolve defaults and choose mode (subscription vs polling).
- Normalize each header into HeadInfo.
- Detect reorgs by parent/hash continuity.
*/
func Run(ctx context.Context, client HeadClient, cfg Config) (*Result, error) {
	if client == nil {
		return nil, errNilClient
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
	ch := make(chan *types.Header, cfg.MaxHeads)
	sub, err := client.SubscribeNewHead(ctx, ch)
	if err != nil {
		return nil, fmt.Errorf("subscribe new head: %w", err)
	}
	defer sub.Unsubscribe()

	heads := make([]HeadInfo, 0, cfg.MaxHeads)
	for len(heads) < cfg.MaxHeads {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case err, ok := <-sub.Err():
			if !ok {
				return &Result{Heads: heads, Mode: "subscription"}, nil
			}
			if err != nil {
				return nil, fmt.Errorf("subscription error: %w", err)
			}
		case h := <-ch:
			if h == nil || h.Number == nil {
				continue
			}
			heads = appendHeadInfo(heads, h)
		}
	}

	return &Result{Heads: heads, Mode: "subscription"}, nil
}

func pollHeads(ctx context.Context, client HeadClient, cfg Config) (*Result, error) {
	ticker := time.NewTicker(cfg.PollInterval)
	defer ticker.Stop()

	heads := make([]HeadInfo, 0, cfg.MaxHeads)
	seen := make(map[common.Hash]struct{}, cfg.MaxHeads)

	for len(heads) < cfg.MaxHeads {
		h, err := client.HeaderByNumber(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("poll latest header: %w", err)
		}
		if h != nil {
			hash := h.Hash()
			if _, exists := seen[hash]; !exists {
				heads = appendHeadInfo(heads, h)
				seen[hash] = struct{}{}
				if len(heads) >= cfg.MaxHeads {
					break
				}
			}
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}
	}

	return &Result{Heads: heads, Mode: "polling"}, nil
}

func appendHeadInfo(existing []HeadInfo, h *types.Header) []HeadInfo {
	info := HeadInfo{
		Number:     h.Number.Uint64(),
		Hash:       h.Hash(),
		ParentHash: h.ParentHash,
	}
	if len(existing) > 0 {
		prev := existing[len(existing)-1]
		info.Reorg = prev.Hash != h.ParentHash
	}
	return append(existing, info)
}
