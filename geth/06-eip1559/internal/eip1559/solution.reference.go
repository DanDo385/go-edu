//go:build reference

package eip1559

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const defaultDynamicGasLimit = 21000

var errNilClient = errors.New("nil fee client")

/*
Reference Solution

Structure:
- Resolve nonce and chain id.
- Resolve base fee and tip.
- Compute/override max fee and build dynamic-fee tx.

Invariants:
- Default max fee is `2*baseFee + tip`.
- Overrides short-circuit RPC suggestions.
*/
func Run(ctx context.Context, client FeeClient, cfg Config) (*Result, error) {
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

	header, err := client.HeaderByNumber(ctx, cfg.BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("header by number: %w", err)
	}
	if header == nil || header.BaseFee == nil {
		return nil, errors.New("header missing base fee")
	}
	baseFee := new(big.Int).Set(header.BaseFee)

	tip := cfg.MaxPriorityFee
	if tip == nil {
		t, err := client.SuggestGasTipCap(ctx)
		if err != nil {
			return nil, fmt.Errorf("suggest gas tip cap: %w", err)
		}
		if t == nil {
			return nil, errors.New("suggest gas tip cap: nil response")
		}
		tip = t
	}

	feeCap := cfg.MaxFee
	if feeCap == nil {
		feeCap = new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), tip)
	}

	gasLimit := cfg.GasLimit
	if gasLimit == 0 {
		gasLimit = defaultDynamicGasLimit
	}

	amount := big.NewInt(0)
	if cfg.AmountWei != nil {
		amount = new(big.Int).Set(cfg.AmountWei)
	}

	txData := &types.DynamicFeeTx{
		ChainID:   new(big.Int).Set(chainID),
		Nonce:     nonce,
		To:        &cfg.To,
		Value:     amount,
		Gas:       gasLimit,
		GasTipCap: new(big.Int).Set(tip),
		GasFeeCap: new(big.Int).Set(feeCap),
		Data:      append([]byte(nil), cfg.Data...),
	}
	tx := types.NewTx(txData)

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
		BaseFee:     baseFee,
	}, nil
}
