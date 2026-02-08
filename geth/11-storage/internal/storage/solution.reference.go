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
Reference Solution - Contract Storage Slots
===========================================

This file demonstrates reading contract storage via eth_getStorageAt. Storage
is organized in 32-byte slots. Simple variables use sequential slots; mappings
use keccak256(abi.encode(key, slot)) to derive the storage slot.

This connects to the Ethereum ecosystem by showing:
- StorageAt(ctx, contract, slotHash, block): raw 32-byte value at slot
- Slot encoding: big.Int → common.Hash for simple slots
- mappingSlotHash: keccak256(leftPad32(key) || slot) for mapping lookups
- cfg.BlockNumber: nil = latest, or historical state

The exercise builds understanding of:
- Solidity storage layout: packed variables, mappings, nested mappings
- append([]byte(nil), value...): defensive copy of RPC response
- cfg.MappingKey: when non-empty, we use mapping slot derivation

Teaching notes (per .cursorrules):
- *big.Int for slot: cfg.Slot is optional pointer; we validate non-nil.
  slotToHash handles nil by returning zero hash (caller should validate).
- Memory: StorageAt may return reused buffers; we copy with append.
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

// slotToHash converts slot number to 32-byte hash; returns zero hash if slot is nil.
func slotToHash(slot *big.Int) common.Hash {
	if slot == nil {
		return common.Hash{}
	}
	return common.BigToHash(slot)
}

// mappingSlotHash computes keccak256(abi.encode(key, slot)): leftPad(key,32) || slot.Bytes().
func mappingSlotHash(key []byte, slot common.Hash) common.Hash {
	buf := make([]byte, 0, 64)
	buf = append(buf, common.LeftPadBytes(key, 32)...)
	buf = append(buf, slot.Bytes()...)
	return crypto.Keccak256Hash(buf)
}
