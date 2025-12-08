//go:build !solution
// +build !solution

package exercise

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Run contains the reference solution for module 12-proofs.
func Run(ctx context.Context, client ProofClient, cfg Config) (*Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if cfg.Account == (common.Address{}) {
		return nil, errors.New("account address required")
	}

	var slotStrings []string
	if len(cfg.Slots) > 0 {
		slotStrings = make([]string, len(cfg.Slots))
		for i, slot := range cfg.Slots {
			slotStrings[i] = slot.Hex()
		}
	}

	proof, err := client.GetProof(ctx, cfg.Account, slotStrings, cfg.BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("get proof: %w", err)
	}
	if proof == nil {
		return nil, errors.New("nil proof response")
	}

	res := &Result{
		Account: AccountProof{
			Balance:     new(big.Int).Set(proof.Balance),
			Nonce:       proof.Nonce,
			CodeHash:    proof.CodeHash,
			StorageHash: proof.StorageHash,
			ProofNodes:  append([]string(nil), proof.AccountProof...),
		},
	}

	if len(proof.StorageProof) > 0 {
		res.Account.Storage = make([]StorageProof, 0, len(proof.StorageProof))
		for _, sp := range proof.StorageProof {
			value := new(big.Int).Set(sp.Value)
			key := common.HexToHash(sp.Key)
			res.Account.Storage = append(res.Account.Storage, StorageProof{
				Key:        key,
				Value:      value,
				ProofNodes: append([]string(nil), sp.Proof...),
			})
		}
	}

	return res, nil
}
