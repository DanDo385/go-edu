//go:build !solution && !reference

package txnonces

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
	// TODO: Implement

	// ============================================================================
	// We continue the pattern of robust validation and sensible defaults.
	// TODO: Implement

	// ============================================================================
	// STEP 2: Determine Sender Address and Nonce
	// TODO: Implement

	// ============================================================================
	// We derive the sender's address from the private key, connecting back to
	// TODO: Implement

	// ============================================================================
	// STEP 3: Get Network and Gas Parameters
	// TODO: Implement

	// ============================================================================
	// To sign a transaction correctly (with EIP-155 replay protection), we need
	// TODO: Implement

	// ============================================================================
	// STEP 4: Create and Sign the Transaction
	// TODO: Implement

	// ============================================================================
	// We assemble all the components into a transaction object.
	// TODO: Implement

	// ============================================================================
	// STEP 5: Broadcast the Transaction
	// TODO: Implement

	// ============================================================================
	// The `NoSend` flag is a useful feature for testing and debugging, allowing
	// TODO: Implement

	// ============================================================================
	// STEP 6: Return the Result
	// TODO: Implement

	// ============================================================================
	// We return the sender's address, the nonce used, and the final signed
	// TODO: Implement

	panic("unimplemented")
}
