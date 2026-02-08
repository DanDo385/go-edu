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
Reference Solution - EIP-1559 Dynamic Fee Transaction
=====================================================

This file demonstrates building an EIP-1559 (London) transaction. EIP-1559
replaced the single gas price with a two-part fee: base fee (burned, set by
protocol) and priority tip (to miner). Users specify maxFeePerGas and
maxPriorityFeePerGas.

This connects to the Ethereum ecosystem by showing:
- types.DynamicFeeTx: ChainID, GasTipCap, GasFeeCap instead of GasPrice
- HeaderByNumber: fetch block header for current base fee
- SuggestGasTipCap: node's suggested priority fee for inclusion
- Fee cap formula: maxFeePerGas >= baseFee + maxPriorityFeePerGas

The exercise builds understanding of:
- Base fee: protocol-determined, varies by block; burned (not to miner)
- Priority tip: goes to miner; user-controlled for inclusion speed
- GasFeeCap: maximum total per gas user will pay
- Defensive copying: header.BaseFee, chainID, tip are *big.Int — copy before use

Teaching notes (per .cursorrules):
- Pointer semantics: cfg.Nonce, cfg.MaxPriorityFee, cfg.MaxFee, cfg.BlockNumber
  are optional; nil means "fetch or compute." *cfg.Nonce dereference when non-nil.
- Memory/ownership: new(big.Int).Set(x) copies; we never share library internals.
  &cfg.To passes address of the config's To — types expects *common.Address.
*/

/*
Run - Build, Sign, and Optionally Send EIP-1559 Transaction

Parameters:
  - ctx: cancellation and timeout for RPC calls
  - client: FeeClient (HeaderByNumber, SuggestGasTipCap, etc.)
  - cfg: tx parameters; nil pointers mean "fetch or compute"

Returns *Result with FromAddress, Nonce, Tx, BaseFee; error on RPC/validation failure.

Algorithm:
  1. Validate inputs, derive from address
  2. Resolve nonce (override or PendingNonceAt)
  3. Fetch chain ID and block header for base fee
  4. Resolve tip (MaxPriorityFee or SuggestGasTipCap)
  5. Compute fee cap if not provided: 2*baseFee + tip
  6. Build DynamicFeeTx, sign, optionally send
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

	// Fetch header for base fee — EIP-1559 base fee is in the block header
	header, err := client.HeaderByNumber(ctx, cfg.BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("header by number: %w", err)
	}
	if header == nil || header.BaseFee == nil {
		return nil, errors.New("header missing base fee")
	}
	baseFee := new(big.Int).Set(header.BaseFee) // Defensive copy

	// Priority tip: goes to miner; user or node suggestion
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

	// Fee cap: max total per gas. Default: 2*baseFee + tip (room for next block)
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
		To:        &cfg.To, // types expects *common.Address; nil To = contract creation
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
