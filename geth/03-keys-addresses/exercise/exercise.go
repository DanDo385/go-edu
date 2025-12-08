//go:build !solution
// +build !solution

package exercise

import (
	"errors"
)

/*
Problem: Generate an Ethereum key, derive its address, and store it in an
encrypted keystore file.

This module is your first step into cryptographic identity. You'll generate a
secp256k1 private key, derive the corresponding public key and Ethereum address,
and then learn how to securely store the private key using Geth's keystore
format. This is the foundation for creating and managing user accounts.

Computer science principles highlighted:
  - Public-key cryptography (secp256k1)
  - Hashing functions for address derivation (Keccak-256)
  - Key-stretching functions for secure password storage (scrypt)
  - Symmetric encryption for data confidentiality (AES)
*/
func Run(cfg Config) (*Result, error) {
	// TODO: Set default values for OutputDir and Passphrase in the Config
	// - If cfg.OutputDir is empty, set it to "./keystore-demo"
	// - If cfg.Passphrase is empty, set it to "changeit"

	// TODO: Generate a new secp256k1 private key
	// - Use `crypto.GenerateKey()`
	// - Handle and wrap any errors

	// TODO: Derive the Ethereum address from the private key's public key
	// - The public key can be accessed via `privKey.PublicKey`
	// - Use `crypto.PubkeyToAddress()` to perform the derivation

	// TODO: Create the output directory for the keystore
	// - Use `os.MkdirAll()` to create the directory specified in `cfg.OutputDir`
	// - The second argument should be the directory permissions, e.g., 0o700
	// - Handle and wrap any errors

	// TODO: Create a new keystore and import the private key
	// - Use `keystore.NewKeyStore()` to create a keystore instance. Pass it
	//   the output directory and the standard scrypt parameters:
	//   `keystore.StandardScryptN` and `keystore.StandardScryptP`.
	// - Use `ks.ImportECDSA()` to import the private key. This will encrypt
	//   the key with the passphrase and save it to a file.
	// - Handle and wrap any errors

	// TODO: Unlock the account to verify the passphrase
	// - Use `ks.Unlock()` to decrypt the key in memory.
	// - Use a `defer ks.Lock(...)` to ensure the key is cleared from memory
	//   when the function exits.
	// - Handle and wrap any errors

	// TODO: Get the path to the generated keystore file
	// - The path is available in `account.URL.Path`
	// - You may need to clean up the path using `filepath.Clean`

	// TODO: Construct and return the Result struct
	// - The result should contain the derived address, the private key as a
	//   hex string, and the path to the keystore file.
	// - To convert the private key to a hex string, use:
	//   `hex.EncodeToString(crypto.FromECDSA(privKey))`

	return nil, errors.New("not implemented")
}