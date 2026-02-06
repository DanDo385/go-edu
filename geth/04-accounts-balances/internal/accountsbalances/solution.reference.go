//go:build reference

package accountsbalances

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var errNilClient = errors.New("nil account client")

/*
Reference Solution

Structure:
- Validate input addresses.
- For each address, fetch balance and bytecode.
- Classify account type (EOA/Contract).

Pointer notes:
- `*big.Int` and `[]byte` are copied to avoid aliasing mutable client-owned memory.
*/
func Run(ctx context.Context, client AccountClient, cfg Config) (*Result, error) {
	if client == nil {
		return nil, errNilClient
	}
	if len(cfg.Addresses) == 0 {
		return nil, errors.New("at least one address is required")
	}

	accounts := make([]AccountState, 0, len(cfg.Addresses))
	for i, addr := range cfg.Addresses {
		if addr == (common.Address{}) {
			return nil, fmt.Errorf("address %d is zero address", i)
		}

		bal, err := client.BalanceAt(ctx, addr, cfg.BlockNumber)
		if err != nil {
			return nil, fmt.Errorf("balance for %s: %w", addr.Hex(), err)
		}
		if bal == nil {
			bal = big.NewInt(0)
		}

		code, err := client.CodeAt(ctx, addr, cfg.BlockNumber)
		if err != nil {
			return nil, fmt.Errorf("code for %s: %w", addr.Hex(), err)
		}

		kind := AccountTypeEOA
		if len(code) > 0 {
			kind = AccountTypeContract
		}

		accounts = append(accounts, AccountState{
			Address: addr,
			Balance: new(big.Int).Set(bal),
			Code:    append([]byte(nil), code...),
			Type:    kind,
		})
	}

	return &Result{Accounts: accounts}, nil
}
