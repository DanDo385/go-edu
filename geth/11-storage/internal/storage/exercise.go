//go:build !solution && !reference

package storage

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func Run(ctx context.Context, client StorageClient, cfg Config) (*Result, error) {
	// TODO: Input Validation
	// TODO: Slot Hash Conversion
	// TODO: Mapping Slot Calculation
	// TODO: Storage Read
	// TODO: Result Construction
	panic("not implemented")
}

func slotToHash(slot *big.Int) common.Hash {
	// TODO: Implement this function
	panic("not implemented")
}

func mappingSlotHash(key []byte, slot common.Hash) common.Hash {
	// TODO: Implement this function
	panic("not implemented")
}
