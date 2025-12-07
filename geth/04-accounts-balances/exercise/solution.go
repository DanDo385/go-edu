//go:build solution
// +build solution

package exercise

import (
	"context"
	"errors"
	"fmt"
	"math/big"
)

// Run contains the reference solution for module 04-accounts-balances.
func Run(ctx context.Context, client AccountClient, cfg Config) (*Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if len(cfg.Addresses) == 0 {
		return nil, errors.New("no addresses provided")
	}

	accounts := make([]AccountState, 0, len(cfg.Addresses))

	for _, addr := range cfg.Addresses {
		bal, err := client.BalanceAt(ctx, addr, cfg.BlockNumber)
		if err != nil {
			return nil, fmt.Errorf("balance %s: %w", addr.Hex(), err)
		}
		// Copy big.Int to avoid caller mutating internal pointer
		if bal != nil {
			bal = new(big.Int).Set(bal)
		}

		code, err := client.CodeAt(ctx, addr, cfg.BlockNumber)
		if err != nil {
			return nil, fmt.Errorf("code %s: %w", addr.Hex(), err)
		}
		codeCopy := append([]byte(nil), code...)

		accType := AccountTypeEOA
		if len(codeCopy) > 0 {
			accType = AccountTypeContract
		}

		accounts = append(accounts, AccountState{
			Address: addr,
			Balance: bal,
			Code:    codeCopy,
			Type:    accType,
		})
	}

	return &Result{Accounts: accounts}, nil
}
