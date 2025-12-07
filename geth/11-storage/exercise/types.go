package exercise

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// StorageClient captures the method needed to read raw storage slots.
type StorageClient interface {
	StorageAt(ctx context.Context, contract common.Address, slot common.Hash, blockNumber *big.Int) ([]byte, error)
}

// Config controls which contract slot to read.
type Config struct {
	Contract    common.Address
	Slot        *big.Int
	MappingKey  []byte
	BlockNumber *big.Int
}

// Result surfaces the resolved slot hash and raw value bytes.
type Result struct {
	ResolvedSlot common.Hash
	Value        []byte
}
