//go:build solution
// +build solution

package exercise

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// Transaction represents a blockchain transaction
type Transaction struct {
	From      string  `json:"from"`      // Sender's public key (hex)
	To        string  `json:"to"`        // Recipient's address
	Amount    float64 `json:"amount"`    // Transfer amount
	Nonce     int64   `json:"nonce"`     // Unique transaction number
	Timestamp int64   `json:"timestamp"` // Unix timestamp
}

// SignedTransaction contains a transaction and its cryptographic signature
type SignedTransaction struct {
	Transaction Transaction `json:"transaction"`
	Signature   string      `json:"signature"`  // Hex-encoded Ed25519 signature
	PublicKey   string      `json:"public_key"` // Hex-encoded public key
}

// Wallet represents a cryptographic wallet with Ed25519 keypair
type Wallet struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
	Address    string // Derived from public key
}

// GenerateWallet creates a new wallet with a fresh Ed25519 keypair.
func GenerateWallet() (*Wallet, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate keypair: %w", err)
	}

	return &Wallet{
		PublicKey:  pub,
		PrivateKey: priv,
		Address:    DeriveAddress(pub),
	}, nil
}

// DeriveAddress derives a human-readable address from a public key.
// The address is the hex-encoded first 20 bytes of the public key.
func DeriveAddress(publicKey ed25519.PublicKey) string {
	// Take first 20 bytes of public key (similar to Ethereum)
	addressBytes := publicKey[:20]
	return hex.EncodeToString(addressBytes)
}

// NewTransaction creates a new unsigned transaction.
func NewTransaction(from *Wallet, to string, amount float64, nonce int64) *Transaction {
	return &Transaction{
		From:      hex.EncodeToString(from.PublicKey),
		To:        to,
		Amount:    amount,
		Nonce:     nonce,
		Timestamp: time.Now().Unix(),
	}
}

// Serialize converts a transaction to bytes for signing/verification.
// Uses JSON encoding to ensure deterministic serialization.
func (tx *Transaction) Serialize() ([]byte, error) {
	data, err := json.Marshal(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize transaction: %w", err)
	}
	return data, nil
}

// Sign signs a transaction with the wallet's private key.
func (w *Wallet) Sign(tx *Transaction) (*SignedTransaction, error) {
	// Serialize transaction
	txBytes, err := tx.Serialize()
	if err != nil {
		return nil, err
	}

	// Sign the serialized transaction
	signature := ed25519.Sign(w.PrivateKey, txBytes)

	return &SignedTransaction{
		Transaction: *tx,
		Signature:   hex.EncodeToString(signature),
		PublicKey:   hex.EncodeToString(w.PublicKey),
	}, nil
}

// Verify verifies the signature on a signed transaction.
func (st *SignedTransaction) Verify() (bool, error) {
	// Decode signature from hex
	signature, err := hex.DecodeString(st.Signature)
	if err != nil {
		return false, fmt.Errorf("invalid signature encoding: %w", err)
	}

	// Decode public key from hex
	publicKey, err := hex.DecodeString(st.PublicKey)
	if err != nil {
		return false, fmt.Errorf("invalid public key encoding: %w", err)
	}

	// Serialize transaction (must use same method as Sign!)
	txBytes, err := st.Transaction.Serialize()
	if err != nil {
		return false, err
	}

	// Verify signature
	valid := ed25519.Verify(publicKey, txBytes, signature)
	return valid, nil
}

// VerifyOwnership checks if the transaction's 'from' field matches
// the public key that signed it.
func (st *SignedTransaction) VerifyOwnership() bool {
	return st.Transaction.From == st.PublicKey
}

// GetTransactionID computes a unique ID for this transaction.
// The ID is the hex-encoded first 16 bytes of the signature.
func (st *SignedTransaction) GetTransactionID() string {
	// Decode signature
	signature, err := hex.DecodeString(st.Signature)
	if err != nil || len(signature) < 16 {
		return ""
	}

	// Use first 16 bytes as transaction ID
	return hex.EncodeToString(signature[:16])
}

// MultiSigTransaction represents a transaction requiring multiple signatures
type MultiSigTransaction struct {
	Transaction Transaction       `json:"transaction"`
	Signatures  map[string]string `json:"signatures"` // publicKey -> signature
	Required    int               `json:"required"`   // M of N signatures needed
}

// NewMultiSigTransaction creates a new multi-signature transaction.
func NewMultiSigTransaction(tx *Transaction, required int) *MultiSigTransaction {
	return &MultiSigTransaction{
		Transaction: *tx,
		Signatures:  make(map[string]string),
		Required:    required,
	}
}

// AddSignature adds a signature to the multi-sig transaction.
func (mst *MultiSigTransaction) AddSignature(wallet *Wallet) error {
	// Serialize transaction
	txBytes, err := mst.Transaction.Serialize()
	if err != nil {
		return err
	}

	// Sign transaction
	signature := ed25519.Sign(wallet.PrivateKey, txBytes)

	// Add to signatures map
	pubKeyHex := hex.EncodeToString(wallet.PublicKey)
	mst.Signatures[pubKeyHex] = hex.EncodeToString(signature)

	return nil
}

// Verify verifies that the multi-sig transaction has enough valid signatures.
func (mst *MultiSigTransaction) Verify() (bool, error) {
	// Serialize transaction
	txBytes, err := mst.Transaction.Serialize()
	if err != nil {
		return false, err
	}

	validCount := 0

	// Verify each signature
	for pubKeyHex, sigHex := range mst.Signatures {
		// Decode public key
		pubKey, err := hex.DecodeString(pubKeyHex)
		if err != nil {
			continue // Skip invalid encoding
		}

		// Decode signature
		sig, err := hex.DecodeString(sigHex)
		if err != nil {
			continue // Skip invalid encoding
		}

		// Verify signature
		if ed25519.Verify(pubKey, txBytes, sig) {
			validCount++
		}
	}

	// Check if we have enough valid signatures
	return validCount >= mst.Required, nil
}
