//go:build !solution && !reference

package keysaddresses

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	defaultOutDir   = "./keystore-demo"
	defaultPassword = "changeit"
)

// Run contains the reference solution for module 03-keys-addresses.
func Run(cfg Config) (*Result, error) {
	// ============================================================================
	// STEP 1: Set Default Configuration
	// TODO: Implement

	// ============================================================================
	// We provide default values for the output directory and passphrase. This makes
	// TODO: Implement

	// ============================================================================
	// STEP 2: Generate a New Private Key
	// TODO: Implement

	// ============================================================================
	// This is the cryptographic heart of an Ethereum account. `crypto.GenerateKey`
	// TODO: Implement

	// ============================================================================
	// STEP 3: Derive the Ethereum Address
	// TODO: Implement

	// ============================================================================
	// The Ethereum address is derived from the public key, which in turn is
	// TODO: Implement

	// ============================================================================
	// STEP 4: Create Keystore Directory
	// TODO: Implement

	// ============================================================================
	// Before writing the keystore file, we need to make sure the directory exists.
	// TODO: Implement

	// ============================================================================
	// STEP 5: Create Keystore and Import Key
	// TODO: Implement

	// ============================================================================
	// This is where we securely store the private key. `keystore.NewKeyStore`
	// TODO: Implement

	// ============================================================================
	// STEP 6: Unlock Account and Defer Lock
	// TODO: Implement

	// ============================================================================
	// To prove that the import was successful, we immediately try to unlock the
	// TODO: Implement

	// ============================================================================
	// STEP 7: Get Keystore Path
	// TODO: Implement

	// ============================================================================
	// We need to return the path to the newly created keystore file. The `account`
	// TODO: Implement

	// ============================================================================
	// STEP 8: Construct and Return Result
	// TODO: Implement

	// ============================================================================
	// We return the derived address, the raw private key as a hex string (for
	// TODO: Implement

	panic("unimplemented")
}
