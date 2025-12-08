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
	// ============================================================================
	// STEP 1: Input Validation and Defaults
	// ============================================================================
	// We continue the pattern of robust validation and sensible defaults.
	if ctx == nil {
		ctx = context.Background()
	}
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if cfg.PrivateKey == nil {
		return nil, errors.New("private key is required")
	}
	// Default to sending 0 ETH if no amount is specified.
	if cfg.AmountWei == nil {
		cfg.AmountWei = big.NewInt(0)
	}
	// Default to the gas limit for a standard ETH transfer.
	if cfg.GasLimit == 0 {
		cfg.GasLimit = defaultLegacyGasLimit
	}

	// ============================================================================
	// STEP 2: Determine Sender Address and Nonce
	// ============================================================================
	// We derive the sender's address from the private key, connecting back to
	// the concepts from module 03.
	from := crypto.PubkeyToAddress(cfg.PrivateKey.PublicKey)

	// Nonce management is critical. We allow the caller to override the nonce,
	// but the standard behavior is to fetch the *pending* nonce.
	// `PendingNonceAt` includes transactions that are in the mempool but not yet
	// mined. Using this prevents "nonce too low" errors if we have other
	// pending transactions.
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

	// ============================================================================
	// STEP 3: Get Network and Gas Parameters
	// ============================================================================
	// To sign a transaction correctly (with EIP-155 replay protection), we need
	// the network's chain ID.
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("chain id: %w", err)
	}
	if chainID == nil {
		return nil, errors.New("chain id was nil")
	}

	// For a legacy transaction, we need to specify a gas price. We allow the
	// caller to override it, but otherwise we ask the node for a suggested price.
	gasPrice := cfg.GasPrice
	if gasPrice == nil {
		gasPrice, err = client.SuggestGasPrice(ctx)
		if err != nil {
			return nil, fmt.Errorf("suggest gas price: %w", err)
		}
	}

	// ============================================================================
	// STEP 4: Create and Sign the Transaction
	// ============================================================================
	// We assemble all the components into a transaction object.
	// Note the defensive copies for AmountWei, GasPrice, and Data to prevent
	// the caller from mutating the transaction after it's created.
	tx := types.NewTransaction(nonce, cfg.To, new(big.Int).Set(cfg.AmountWei), cfg.GasLimit, new(big.Int).Set(gasPrice), append([]byte(nil), cfg.Data...))

	// `types.LatestSignerForChainID` creates a signer that implements EIP-155.
	// This embeds the chain ID into the signature's `v` value, making it
	// impossible to replay the transaction on a different chain.
	signer := types.LatestSignerForChainID(chainID)
	signedTx, err := types.SignTx(tx, signer, cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("sign tx: %w", err)
	}

	// ============================================================================
	// STEP 5: Broadcast the Transaction
	// ============================================================================
	// The `NoSend` flag is a useful feature for testing and debugging, allowing
	// us to inspect a signed transaction without actually sending it.
	// If not set, we broadcast the transaction to the network.
	if !cfg.NoSend {
		if err := client.SendTransaction(ctx, signedTx); err != nil {
			// Common errors here include "nonce too low", "insufficient funds",
			// or "transaction underpriced".
			return nil, fmt.Errorf("send tx: %w", err)
		}
	}

	// ============================================================================
	// STEP 6: Return the Result
	// ============================================================================
	// We return the sender's address, the nonce used, and the final signed
	// transaction object. The caller can then inspect the transaction hash, etc.
	return &Result{
		FromAddress: from,
		Nonce:       nonce,
		Tx:          signedTx,
	}, nil
}
