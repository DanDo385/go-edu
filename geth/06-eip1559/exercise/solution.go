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

const defaultDynamicGasLimit = 21000

// Run contains the reference solution for module 06-eip1559.
func Run(ctx context.Context, client FeeClient, cfg Config) (*Result, error) {
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
		cfg.GasLimit = defaultDynamicGasLimit
	}

	from := crypto.PubkeyToAddress(cfg.PrivateKey.PublicKey)

	var nonce uint64
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

	header, err := client.HeaderByNumber(ctx, cfg.BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("header by number: %w", err)
	}
	if header == nil || header.BaseFee == nil {
		return nil, errors.New("base fee unavailable (upgrade to London block)")
	}
	baseFee := new(big.Int).Set(header.BaseFee)

	tipCap := cfg.MaxPriorityFee
	if tipCap == nil {
		tipCap, err = client.SuggestGasTipCap(ctx)
		if err != nil {
			return nil, fmt.Errorf("suggest gas tip cap: %w", err)
		}
	}
	tipCap = new(big.Int).Set(tipCap)

	maxFee := cfg.MaxFee
	if maxFee == nil {
		// Rule of thumb: pay at most 2x current base fee + tip to cover variance.
		twoBase := new(big.Int).Mul(baseFee, big.NewInt(2))
		maxFee = new(big.Int).Add(twoBase, tipCap)
	} else {
		maxFee = new(big.Int).Set(maxFee)
	}

	gasLimit := cfg.GasLimit
	dataCopy := append([]byte(nil), cfg.Data...)

	txData := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: maxFee,
		Gas:       gasLimit,
		To:        &cfg.To,
		Value:     new(big.Int).Set(cfg.AmountWei),
		Data:      dataCopy,
	}

	tx := types.NewTx(txData)
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
		BaseFee:     baseFee,
	}, nil
}
