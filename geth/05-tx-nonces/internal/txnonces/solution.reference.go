//go:build reference

package txnonces

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const defaultLegacyGasLimit = 21000

var errNilClient = errors.New("nil tx client")

/*
Reference Solution

Structure:
- Resolve nonce, chain id, gas settings.
- Build and sign a legacy transaction.
- Optionally broadcast.

Invariants:
- Override values in config take precedence over RPC suggestions.
- Signed tx uses the resolved chain id.
*/
func Run(ctx context.Context, client TXClient, cfg Config) (*Result, error) {
	if client == nil {
		return nil, errNilClient
	}
	if cfg.PrivateKey == nil {
		return nil, errors.New("private key is required")
	}
	if cfg.To == (common.Address{}) {
		return nil, errors.New("destination address is required")
	}

	from := crypto.PubkeyToAddress(cfg.PrivateKey.PublicKey)

	var nonce uint64
	if cfg.Nonce != nil {
		nonce = *cfg.Nonce
	} else {
		n, err := client.PendingNonceAt(ctx, from)
		if err != nil {
			return nil, fmt.Errorf("pending nonce: %w", err)
		}
		nonce = n
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("chain id: %w", err)
	}
	if chainID == nil {
		return nil, errors.New("chain id: nil response")
	}

	gasPrice := cfg.GasPrice
	if gasPrice == nil {
		g, err := client.SuggestGasPrice(ctx)
		if err != nil {
			return nil, fmt.Errorf("suggest gas price: %w", err)
		}
		if g == nil {
			return nil, errors.New("suggest gas price: nil response")
		}
		gasPrice = g
	}

	gasLimit := cfg.GasLimit
	if gasLimit == 0 {
		gasLimit = defaultLegacyGasLimit
	}

	amount := big.NewInt(0)
	if cfg.AmountWei != nil {
		amount = new(big.Int).Set(cfg.AmountWei)
	}

	tx := types.NewTransaction(
		nonce,
		cfg.To,
		amount,
		gasLimit,
		new(big.Int).Set(gasPrice),
		append([]byte(nil), cfg.Data...),
	)

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("sign transaction: %w", err)
	}

	if !cfg.NoSend {
		if err := client.SendTransaction(ctx, signedTx); err != nil {
			return nil, fmt.Errorf("send transaction: %w", err)
		}
	}

	return &Result{
		FromAddress: from,
		Nonce:       nonce,
		Tx:          signedTx,
	}, nil
}
