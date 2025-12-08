//go:build solution
// +build solution

package exercise

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
	// ============================================================================
	// We provide default values for the output directory and passphrase. This makes
	// the function easier to use for simple cases, while still allowing callers
	// to customize the behavior.
	if cfg.OutputDir == "" {
		cfg.OutputDir = defaultOutDir
	}
	if cfg.Passphrase == "" {
		cfg.Passphrase = defaultPassword
	}

	// ============================================================================
	// STEP 2: Generate a New Private Key
	// ============================================================================
	// This is the cryptographic heart of an Ethereum account. `crypto.GenerateKey`
	// creates a new private key using the secp256k1 elliptic curve.
	//
	// Computer Science Principle: Public-Key Cryptography. The private key can be
	// used to generate a public key, but the private key cannot be derived from
	// the public key. This is the foundation of digital signatures.
	//
	// `crypto.GenerateKey` uses `crypto/rand` to ensure the key is generated
	// from a cryptographically secure random number source.
	privKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("generate secp256k1 key: %w", err)
	}

	// ============================================================================
	// STEP 3: Derive the Ethereum Address
	// ============================================================================
	// The Ethereum address is derived from the public key, which in turn is
	// derived from the private key. `crypto.PubkeyToAddress` handles this logic:
	//   1. It takes the public key (part of `privKey.PublicKey`).
	//   2. It computes the Keccak-256 hash of the public key.
	//   3. It takes the last 20 bytes of the hash as the address.
	//
	// This process is deterministic: the same public key will always produce the
	// same address.
	addr := crypto.PubkeyToAddress(privKey.PublicKey)

	// ============================================================================
	// STEP 4: Create Keystore Directory
	// ============================================================================
	// Before writing the keystore file, we need to make sure the directory exists.
	// `os.MkdirAll` is idempotent: it will create the directory if it doesn't
	// exist, and do nothing if it already exists. The `0o700` permissions mean
	// that only the owner can read, write, and execute the directory.
	if err := os.MkdirAll(cfg.OutputDir, 0o700); err != nil {
		return nil, fmt.Errorf("create keystore directory: %w", err)
	}

	// ============================================================================
	// STEP 5: Create Keystore and Import Key
	// ============================================================================
	// This is where we securely store the private key. `keystore.NewKeyStore`
	// creates an object that manages a directory of encrypted key files.
	//
	// `keystore.StandardScryptN` and `keystore.StandardScryptP` are the standard
	// parameters for the scrypt key derivation function. Scrypt is designed to be
	// slow, making brute-force attacks on the passphrase more difficult.
	//
	// `ks.ImportECDSA` takes the private key and passphrase, encrypts the key,
	// and saves it to a new file in the keystore directory.
	ks := keystore.NewKeyStore(cfg.OutputDir, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.ImportECDSA(privKey, cfg.Passphrase)
	if err != nil {
		return nil, fmt.Errorf("import key into keystore: %w", err)
	}

	// ============================================================================
	// STEP 6: Unlock Account and Defer Lock
	// ============================================================================
	// To prove that the import was successful, we immediately try to unlock the
	// account with the same passphrase. `ks.Unlock` decrypts the key and stores
	// it in memory.
	//
	// CRITICAL: For security, we must ensure the decrypted key is cleared from
	// memory when we're done with it. `defer ks.Lock(account.Address)` schedules
	// the `Lock` function to be called just before `Run` returns. This is a
	// fundamental pattern for resource cleanup in Go.
	if err := ks.Unlock(account, cfg.Passphrase); err != nil {
		return nil, fmt.Errorf("unlock keystore: %w", err)
	}
	defer ks.Lock(account.Address)

	// ============================================================================
	// STEP 7: Get Keystore Path
	// ============================================================================
	// We need to return the path to the newly created keystore file. The `account`
	// object contains a URL field that points to the file. We do some basic
	// path cleaning to make it a clean, absolute path.
	keystorePath := account.URL.Path
	if keystorePath == "" {
		keystorePath = filepath.Join(cfg.OutputDir, filepath.Base(account.URL.String()))
	} else {
		keystorePath = filepath.Clean(keystorePath)
	}

	// ============================================================================
	// STEP 8: Construct and Return Result
	// ============================================================================
	// We return the derived address, the raw private key as a hex string (for
	// demonstration purposes - NEVER log or store raw private keys in a real
	// application!), and the path to the encrypted keystore file.
	return &Result{
		Address:       addr,
		PrivateKeyHex: hex.EncodeToString(crypto.FromECDSA(privKey)),
		KeystorePath:  keystorePath,
	}, nil
}
