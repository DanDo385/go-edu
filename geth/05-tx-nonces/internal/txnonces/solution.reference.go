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
Reference Solution - Legacy Transaction with Nonce
==================================================

This file demonstrates building and sending a legacy (pre-EIP-1559) Ethereum
transaction. Every Ethereum transaction requires a nonce — a per-account
sequence number that prevents replay and ensures ordering.

This connects to the Ethereum ecosystem by showing:
- go-ethereum types.NewTransaction for legacy tx format
- PendingNonceAt: fetch next usable nonce from the node's mempool view
- SuggestGasPrice: node's recommended gas price for inclusion
- types.SignTx with EIP-155 signer (Chain ID for replay protection)
- SendTransaction: broadcast to the network

The exercise builds understanding of:
- Nonce semantics: must be sequential per account; gaps cause stuck txs
- Gas pricing: legacy uses single gas price; EIP-1559 uses base fee + tip
- Private key handling: derive address from pubkey; sign with secp256k1
- Defensive copying: big.Int and byte slices must be copied to avoid mutation

Teaching notes (per .cursorrules):
- Pointer semantics: cfg.Nonce is *uint64 — optional override. *cfg.Nonce
  dereferences to read the value. nil means "fetch from chain." cfg.PrivateKey,
  cfg.GasPrice, cfg.AmountWei: same pattern. big.Int is a struct; we use
  pointers to avoid copying large multi-word integers.
- Memory/ownership: new(big.Int).Set(x) creates an independent copy. append(
  []byte(nil), cfg.Data...) copies the data slice. We never share mutable
  references with the caller or library internals.
- Error surfaces: RPC calls can fail (network, node overload). We wrap errors
  with %w for debugging. Validate nil responses — some RPCs return nil without error.
*/

/*
Run - Build, Sign, and Optionally Send Legacy Transaction

Parameters:
  - ctx: cancellation and timeout context for RPC calls
  - client: TXClient interface (BlockNumber, PendingNonceAt, etc.)
  - cfg: transaction parameters; nil pointers mean "use default or fetch"

Returns *Result with FromAddress, Nonce, signed Tx; error on any RPC or validation failure.

Algorithm:
  1. Validate client, private key, destination
  2. Derive from address from public key
  3. Resolve nonce (cfg.Nonce override or PendingNonceAt)
  4. Fetch chain ID and gas price if not provided
  5. Build unsigned tx, sign with EIP-155 signer, optionally send
*/
func Run(ctx context.Context, client TXClient, cfg Config) (*Result, error) {
	// Step 1: Input validation
	if client == nil {
		return nil, errNilClient
	}
	if cfg.PrivateKey == nil {
		return nil, errors.New("private key is required")
	}
	if cfg.To == (common.Address{}) {
		return nil, errors.New("destination address is required")
	}

	// Step 2: Derive sender address from public key
	// Ethereum address = last 20 bytes of Keccak-256(pubkey)
	from := crypto.PubkeyToAddress(cfg.PrivateKey.PublicKey)

	// Step 3: Resolve nonce
	// *cfg.Nonce: if caller provided nonce override, dereference to get uint64.
	// Otherwise fetch from node — PendingNonceAt returns next nonce for account.
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

	// Step 4: Fetch chain ID (required for EIP-155 signing)
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("chain id: %w", err)
	}
	if chainID == nil {
		return nil, errors.New("chain id: nil response")
	}

	// Step 5: Resolve gas price (use cfg or fetch from node)
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

	// Step 6: Resolve amount (defensive copy of big.Int)
	amount := big.NewInt(0)
	if cfg.AmountWei != nil {
		amount = new(big.Int).Set(cfg.AmountWei)
	}

	// Step 7: Build unsigned transaction
	// new(big.Int).Set(gasPrice): copy — don't pass library's mutable big.Int
	// append([]byte(nil), cfg.Data...): copy data slice so caller can't mutate
	tx := types.NewTransaction(
		nonce,
		cfg.To,
		amount,
		gasLimit,
		new(big.Int).Set(gasPrice),
		append([]byte(nil), cfg.Data...),
	)

	// Step 8: Sign with EIP-155 signer (Chain ID in signature for replay protection)
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("sign transaction: %w", err)
	}

	// Step 9: Optionally broadcast to network
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
