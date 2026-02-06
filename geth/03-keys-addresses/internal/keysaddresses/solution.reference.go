//go:build reference

package keysaddresses

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	defaultOutputDir  = "./keystore-demo"
	defaultPassphrase = "changeit"
)

/*
Reference Solution

Structure:
- Normalize config defaults.
- Create a keystore and new account.
- Decrypt generated keyfile to surface private key hex for learning purposes.

Invariant:
- Result address, keyfile, and private key must represent the same account.
*/
func Run(cfg Config) (*Result, error) {
	outDir := cfg.OutputDir
	if outDir == "" {
		outDir = defaultOutputDir
	}
	passphrase := cfg.Passphrase
	if passphrase == "" {
		passphrase = defaultPassphrase
	}

	if err := os.MkdirAll(outDir, 0o700); err != nil {
		return nil, fmt.Errorf("create output dir: %w", err)
	}

	ks := keystore.NewKeyStore(outDir, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(passphrase)
	if err != nil {
		return nil, fmt.Errorf("create account: %w", err)
	}

	keyJSON, err := os.ReadFile(account.URL.Path)
	if err != nil {
		return nil, fmt.Errorf("read keystore file: %w", err)
	}

	key, err := keystore.DecryptKey(keyJSON, passphrase)
	if err != nil {
		return nil, fmt.Errorf("decrypt keystore file: %w", err)
	}

	return &Result{
		Address:       account.Address,
		PrivateKeyHex: hex.EncodeToString(crypto.FromECDSA(key.PrivateKey)),
		KeystorePath:  account.URL.Path,
	}, nil
}
