package keysaddresses

import "github.com/ethereum/go-ethereum/common"

// Config controls the keystore output directory and passphrase.
type Config struct {
	// PrivateKeyHex is an optional hex-encoded private key.
	// When empty, a new key pair is generated.
	PrivateKeyHex string

	// OutputDir is where the encrypted keyfile will be stored.
	// Defaults to "./keystore-demo" when empty.
	OutputDir string

	// Passphrase encrypts/decrypts the keystore. The default "changeit"
	// is intentionally weak so demos can be recovered quickly.
	Passphrase string
}

// Result summarizes the outputs of a Run invocation.
type Result struct {
	Address       common.Address
	PublicKey     []byte
	PrivateKeyHex string
	KeystorePath  string
}
