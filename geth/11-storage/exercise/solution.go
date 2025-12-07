//go:build solution
// +build solution

package exercise

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var zeroHash = common.Hash{}

// Run contains the reference solution for module 11-storage.
func Run(ctx context.Context, client StorageClient, cfg Config) (*Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if cfg.Contract == (common.Address{}) {
		return nil, errors.New("contract address required")
	}
	if cfg.Slot == nil {
		return nil, errors.New("slot is required")
	}

	slotHash := slotToHash(cfg.Slot)
	if len(cfg.MappingKey) > 0 {
		slotHash = mappingSlotHash(cfg.MappingKey, slotHash)
	}

	value, err := client.StorageAt(ctx, cfg.Contract, slotHash, cfg.BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("storage at slot %s: %w", slotHash.Hex(), err)
	}

	return &Result{
		ResolvedSlot: slotHash,
		Value:        value,
	}, nil
}

func slotToHash(slot *big.Int) common.Hash {
	if slot == nil {
		return zeroHash
	}
	return common.BigToHash(slot)
}

func mappingSlotHash(key []byte, slot common.Hash) common.Hash {
	keyPadded := common.LeftPadBytes(key, 32)
	data := append(keyPadded, slot.Bytes()...)
	return common.BytesToHash(crypto.Keccak256(data))
}
