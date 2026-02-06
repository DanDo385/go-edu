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
Reference Solution

Structure:
- Validate account input.
- Convert slot hashes to hex strings required by debug API.
- Translate gethclient response into lesson-local types.

Pointer notes:
- `*big.Int` balances/values are copied to keep result ownership local.
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
