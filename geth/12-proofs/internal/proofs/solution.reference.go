//go:build reference

package proofs

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var errNilClient = errors.New("nil proof client")

/*
Reference Solution - eth_getProof (Account and Storage)
======================================================

This file demonstrates eth_getProof: Merkle-Patricia proof for an account and
optional storage slots at a block. Used by light clients and bridges to verify
state without full node data.

This connects to the Ethereum ecosystem by showing:
- GetProof(account, slots, block): account proof + per-slot storage proofs
- AccountProof: balance, nonce, codeHash, storageHash, accountProof nodes
- StorageProof: key, value, proof nodes for each requested slot
- Defensive copies: Balance, StorageProof[].Value, Proof nodes — RPC may reuse

The exercise builds understanding of:
- Merkle proofs: verify account/storage inclusion using block state root
- append([]string(nil), sp.Proof...): copy string slices from response
- Slots as hex strings: geth expects []string of slot hashes

Teaching notes (per .cursorrules):
- resp.Balance, sp.Value: *big.Int from RPC; new(big.Int).Set() for copy.
- resp may be reused by client; we copy all nested data into Result.
*/
func Run(ctx context.Context, client ProofClient, cfg Config) (*Result, error) {
	if client == nil {
		return nil, errNilClient
	}
	if cfg.Account == (common.Address{}) {
		return nil, errors.New("account is required")
	}

	slots := make([]string, 0, len(cfg.Slots))
	for _, slot := range cfg.Slots {
		slots = append(slots, slot.Hex())
	}

	resp, err := client.GetProof(ctx, cfg.Account, slots, cfg.BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("get proof: %w", err)
	}
	if resp == nil {
		return nil, errors.New("get proof: nil response")
	}

	balance := big.NewInt(0)
	if resp.Balance != nil {
		balance = new(big.Int).Set(resp.Balance)
	}

	storage := make([]StorageProof, 0, len(resp.StorageProof))
	for _, sp := range resp.StorageProof {
		v := big.NewInt(0)
		if sp.Value != nil {
			v = new(big.Int).Set(sp.Value)
		}
		storage = append(storage, StorageProof{
			Key:        common.HexToHash(sp.Key),
			Value:      v,
			ProofNodes: append([]string(nil), sp.Proof...),
		})
	}

	return &Result{
		Account: AccountProof{
			Balance:     balance,
			Nonce:       resp.Nonce,
			CodeHash:    resp.CodeHash,
			StorageHash: resp.StorageHash,
			ProofNodes:  append([]string(nil), resp.AccountProof...),
			Storage:     storage,
		},
	}, nil
}
