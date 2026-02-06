//go:build reference

package storage

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var errNilClient = errors.New("nil storage client")

/*
Reference Solution

Structure:
- Validate contract+slot inputs.
- Resolve canonical slot hash.
- If a mapping key is provided, derive `keccak256(pad(key) || slot)`.
- Read slot bytes and return a defensive copy.

Pointer notes:
- Slot arithmetic uses `*big.Int`; we never mutate caller-provided values.
*/
func Run(ctx context.Context, client StorageClient, cfg Config) (*Result, error) {
	if client == nil {
		return nil, errNilClient
	}
	if cfg.Contract == (common.Address{}) {
		return nil, errors.New("contract address is required")
	}
	if cfg.Slot == nil {
		return nil, errors.New("slot is required")
	}

	resolved := slotToHash(cfg.Slot)
	if len(cfg.MappingKey) > 0 {
		resolved = mappingSlotHash(cfg.MappingKey, resolved)
	}

	value, err := client.StorageAt(ctx, cfg.Contract, resolved, cfg.BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("storage at slot %s: %w", resolved.Hex(), err)
	}

	return &Result{
		ResolvedSlot: resolved,
		Value:        append([]byte(nil), value...),
	}, nil
}

func slotToHash(slot *big.Int) common.Hash {
	if slot == nil {
		return common.Hash{}
	}
	return common.BigToHash(slot)
}

func mappingSlotHash(key []byte, slot common.Hash) common.Hash {
	buf := make([]byte, 0, 64)
	buf = append(buf, common.LeftPadBytes(key, 32)...)
	buf = append(buf, slot.Bytes()...)
	return crypto.Keccak256Hash(buf)
}
