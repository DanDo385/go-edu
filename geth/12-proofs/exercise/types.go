package exercise

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
)

// ProofClient captures the single method required to fetch MPT proofs.
type ProofClient interface {
	GetProof(ctx context.Context, account common.Address, slots []string, blockNumber *big.Int) (*gethclient.AccountResult, error)
}

// Config controls which account/storage slots to prove.
type Config struct {
	Account     common.Address
	Slots       []common.Hash
	BlockNumber *big.Int
}

// StorageProof summarizes the storage slot proof.
type StorageProof struct {
	Key        common.Hash
	Value      *big.Int
	ProofNodes []string
}

// AccountProof captures account-level data plus storage proofs.
type AccountProof struct {
	Balance     *big.Int
	Nonce       uint64
	CodeHash    common.Hash
	StorageHash common.Hash
	ProofNodes  []string
	Storage     []StorageProof
}

// Result is returned to callers.
type Result struct {
	Account AccountProof
}
