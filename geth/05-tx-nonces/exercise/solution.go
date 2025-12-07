//go:build solution
// +build solution

package exercise

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const defaultLegacyGasLimit = 21000

// Run contains the reference solution for module 05-tx-nonces.
func Run(ctx context.Context, client TXClient, cfg Config) (*Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if cfg.PrivateKey == nil {
		return nil, errors.New("private key is required")
	}
	if cfg.AmountWei == nil {
		cfg.AmountWei = big.NewInt(0)
	}
	if cfg.GasLimit == 0 {
		cfg.GasLimit = defaultLegacyGasLimit
	}

	from := crypto.PubkeyToAddress(cfg.PrivateKey.PublicKey)
	nonce := uint64(0)
	var err error
	if cfg.Nonce != nil {
		nonce = *cfg.Nonce
	} else {
		nonce, err = client.PendingNonceAt(ctx, from)
		if err != nil {
			return nil, fmt.Errorf("pending nonce: %w", err)
		}
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("chain id: %w", err)
	}
	if chainID == nil {
		return nil, errors.New("chain id was nil")
	}

	gasPrice := cfg.GasPrice
	if gasPrice == nil {
		gasPrice, err = client.SuggestGasPrice(ctx)
		if err != nil {
			return nil, fmt.Errorf("suggest gas price: %w", err)
		}
	}

	tx := types.NewTransaction(nonce, cfg.To, new(big.Int).Set(cfg.AmountWei), cfg.GasLimit, new(big.Int).Set(gasPrice), append([]byte(nil), cfg.Data...))
	signer := types.LatestSignerForChainID(chainID)
	signedTx, err := types.SignTx(tx, signer, cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("sign tx: %w", err)
	}

	if !cfg.NoSend {
		if err := client.SendTransaction(ctx, signedTx); err != nil {
			return nil, fmt.Errorf("send tx: %w", err)
		}
	}

	return &Result{
		FromAddress: from,
		Nonce:       nonce,
		Tx:          signedTx,
	}, nil
}
