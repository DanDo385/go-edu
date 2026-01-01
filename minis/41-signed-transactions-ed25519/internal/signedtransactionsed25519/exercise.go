//go:build !solution && !reference

package signedtransactionsed25519

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type Transaction struct {
	From      string  `json:"from"`      // Sender's public key (hex)
	To        string  `json:"to"`        // Recipient's address
	Amount    float64 `json:"amount"`    // Transfer amount
	Nonce     int64   `json:"nonce"`     // Unique transaction number
	Timestamp int64   `json:"timestamp"` // Unix timestamp
}

type SignedTransaction struct {
	Transaction Transaction `json:"transaction"`
	Signature   string      `json:"signature"`  // Hex-encoded Ed25519 signature
	PublicKey   string      `json:"public_key"` // Hex-encoded public key
}

type Wallet struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
	Address    string // Derived from public key
}

type MultiSigTransaction struct {
	Transaction Transaction       `json:"transaction"`
	Signatures  map[string]string `json:"signatures"` // publicKey -> signature
	Required    int               `json:"required"`   // M of N signatures needed
}

// GenerateWallet implements the exercise.
//
// TODO: Implement this function
func GenerateWallet() (*Wallet, error) {
	// TODO: Implement
	return nil, nil
}

// DeriveAddress implements the exercise.
//
// TODO: Implement this function
func DeriveAddress(publicKey ed25519.PublicKey) string {
	// TODO: Implement
	return ""
}

// NewTransaction implements the exercise.
//
// TODO: Implement this function
func NewTransaction(from *Wallet, to string, amount float64, nonce int64) *Transaction {
	// TODO: Implement
	return nil
}

// Serialize implements the exercise.
//
// TODO: Implement this function
func (tx *Transaction) Serialize() ([]byte, error) {
	// TODO: Implement
	return nil, nil
}

// Sign implements the exercise.
//
// TODO: Implement this function
func (w *Wallet) Sign(tx *Transaction) (*SignedTransaction, error) {
	// TODO: Implement
	return nil, nil
}

// Verify implements the exercise.
//
// TODO: Implement this function
func (st *SignedTransaction) Verify() (bool, error) {
	// TODO: Implement
	return false, nil
}

// VerifyOwnership implements the exercise.
//
// TODO: Implement this function
func (st *SignedTransaction) VerifyOwnership() bool {
	// TODO: Implement
	return false
}

// GetTransactionID implements the exercise.
//
// TODO: Implement this function
func (st *SignedTransaction) GetTransactionID() string {
	// TODO: Implement
	return ""
}

// NewMultiSigTransaction implements the exercise.
//
// TODO: Implement this function
func NewMultiSigTransaction(tx *Transaction, required int) *MultiSigTransaction {
	// TODO: Implement
	return nil
}

// AddSignature implements the exercise.
//
// TODO: Implement this function
func (mst *MultiSigTransaction) AddSignature(wallet *Wallet) error {
	// TODO: Implement
	return nil
}

// Verify implements the exercise.
//
// TODO: Implement this function
func (mst *MultiSigTransaction) Verify() (bool, error) {
	// TODO: Implement
	return false, nil
}
