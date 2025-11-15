//go:build !solution
// +build !solution

package exercise

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
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
	Signature   string      `json:"signature"` // Hex-encoded Ed25519 signature
	PublicKey   string      `json:"public_key"` // Hex-encoded public key
}

// Wallet represents a cryptographic wallet with Ed25519 keypair
type Wallet struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
	Address    string // Derived from public key
}

// GenerateWallet creates a new wallet with a fresh Ed25519 keypair.
// It generates a cryptographically secure random keypair and derives
// an address from the public key.
//
// Returns:
//   - *Wallet: A new wallet with public key, private key, and address
//   - error: Any error during key generation
//
// Example usage:
//
//	wallet, err := GenerateWallet()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Address: %s\n", wallet.Address)
func GenerateWallet() (*Wallet, error) {
	// TODO: Implement wallet generation
	// Steps:
	// 1. Use ed25519.GenerateKey(nil) to generate a keypair
	// 2. Handle any errors (though GenerateKey rarely fails)
	// 3. Derive an address using DeriveAddress function
	// 4. Return a Wallet struct with all fields populated
	return nil, nil
}

// DeriveAddress derives a human-readable address from a public key.
// The address is the hex-encoded first 20 bytes of the public key.
//
// Parameters:
//   - publicKey: Ed25519 public key (32 bytes)
//
// Returns:
//   - string: Hex-encoded address (40 characters)
//
// Example:
//
//	pub := wallet.PublicKey
//	addr := DeriveAddress(pub)
//	// addr = "a3b4c5d6e7f8..."
func DeriveAddress(publicKey ed25519.PublicKey) string {
	// TODO: Implement address derivation
	// Steps:
	// 1. Take the first 20 bytes of the public key
	// 2. Convert to hex string using hex.EncodeToString
	// 3. Return the hex string
	//
	// Note: In real blockchains, address derivation is more complex:
	//   - Bitcoin: SHA256 + RIPEMD160 + Base58Check
	//   - Ethereum: Keccak256 + take last 20 bytes
	// We use a simplified version for educational purposes.
	return ""
}

// NewTransaction creates a new unsigned transaction.
//
// Parameters:
//   - from: Sender's wallet (provides public key for 'from' field)
//   - to: Recipient's address
//   - amount: Transfer amount (must be positive)
//   - nonce: Unique transaction number (prevents replay attacks)
//
// Returns:
//   - *Transaction: The created transaction
//
// Example usage:
//
//	tx := NewTransaction(wallet, "recipient_addr", 10.5, 1)
func NewTransaction(from *Wallet, to string, amount float64, nonce int64) *Transaction {
	// TODO: Implement transaction creation
	// Steps:
	// 1. Create a Transaction struct
	// 2. Set From to hex-encoded public key: hex.EncodeToString(from.PublicKey)
	// 3. Set To to the recipient address
	// 4. Set Amount to the specified amount
	// 5. Set Nonce to the specified nonce
	// 6. Set Timestamp to current Unix time: time.Now().Unix()
	// 7. Return pointer to the transaction
	return nil
}

// Serialize converts a transaction to bytes for signing/verification.
// Uses JSON encoding to ensure deterministic serialization.
//
// Returns:
//   - []byte: JSON-encoded transaction
//   - error: Any error during serialization
func (tx *Transaction) Serialize() ([]byte, error) {
	// TODO: Implement serialization
	// Steps:
	// 1. Use json.Marshal(tx) to convert to JSON bytes
	// 2. Return the bytes and any error
	//
	// Important: Always use the same serialization method for both
	// signing and verification, or signatures will be invalid!
	return nil, nil
}

// Sign signs a transaction with the wallet's private key.
// The signature proves that the transaction was created by the
// owner of the private key, without revealing the private key itself.
//
// Parameters:
//   - tx: Transaction to sign
//
// Returns:
//   - *SignedTransaction: Transaction with signature and public key
//   - error: Any error during signing
//
// Example usage:
//
//	tx := NewTransaction(wallet, "recipient", 10.5, 1)
//	signedTx, err := wallet.Sign(tx)
//	if err != nil {
//	    log.Fatal(err)
//	}
func (w *Wallet) Sign(tx *Transaction) (*SignedTransaction, error) {
	// TODO: Implement transaction signing
	// Steps:
	// 1. Serialize the transaction using tx.Serialize()
	// 2. Handle any serialization errors
	// 3. Sign the serialized bytes: ed25519.Sign(w.PrivateKey, txBytes)
	// 4. Encode signature to hex: hex.EncodeToString(signature)
	// 5. Encode public key to hex: hex.EncodeToString(w.PublicKey)
	// 6. Create and return SignedTransaction with:
	//    - Transaction: *tx
	//    - Signature: hex-encoded signature
	//    - PublicKey: hex-encoded public key
	return nil, nil
}

// Verify verifies the signature on a signed transaction.
// Returns true if the signature is valid, false otherwise.
//
// A valid signature proves:
//   - The transaction was signed by the holder of the private key
//   - The transaction has not been modified since signing
//   - The signer cannot deny creating the signature
//
// Returns:
//   - bool: true if signature is valid, false otherwise
//   - error: Any error during verification (e.g., invalid encoding)
//
// Example usage:
//
//	valid, err := signedTx.Verify()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if valid {
//	    fmt.Println("Signature is valid!")
//	}
func (st *SignedTransaction) Verify() (bool, error) {
	// TODO: Implement signature verification
	// Steps:
	// 1. Decode the hex signature: hex.DecodeString(st.Signature)
	// 2. Handle any decoding errors
	// 3. Decode the hex public key: hex.DecodeString(st.PublicKey)
	// 4. Handle any decoding errors
	// 5. Serialize the transaction: st.Transaction.Serialize()
	// 6. Handle any serialization errors
	// 7. Verify: ed25519.Verify(publicKey, txBytes, signature)
	// 8. Return the verification result (bool) and nil error
	//
	// Important: Use the same serialization method as Sign!
	return false, nil
}

// VerifyOwnership checks if the transaction's 'from' field matches
// the public key that signed it. This prevents someone from signing
// a transaction claiming to be from a different address.
//
// Returns:
//   - bool: true if 'from' matches signing public key
//
// Example usage:
//
//	if !signedTx.VerifyOwnership() {
//	    fmt.Println("Transaction sender doesn't match signature!")
//	}
func (st *SignedTransaction) VerifyOwnership() bool {
	// TODO: Implement ownership verification
	// Steps:
	// 1. Check if st.Transaction.From equals st.PublicKey
	// 2. Return true if they match, false otherwise
	//
	// Why this matters: Someone could sign a transaction with their
	// private key but put a different address in the 'from' field.
	// This check prevents that attack.
	return false
}

// GetTransactionID computes a unique ID for this transaction.
// The ID is the hex-encoded first 16 bytes of the signature.
//
// Returns:
//   - string: Transaction ID (32 hex characters)
//
// Example usage:
//
//	txID := signedTx.GetTransactionID()
//	fmt.Printf("Transaction ID: %s\n", txID)
func (st *SignedTransaction) GetTransactionID() string {
	// TODO: Implement transaction ID generation
	// Steps:
	// 1. Decode signature from hex: hex.DecodeString(st.Signature)
	// 2. Take first 16 bytes of signature
	// 3. Encode to hex: hex.EncodeToString(bytes[:16])
	// 4. Return hex string
	//
	// Note: In real blockchains, transaction IDs are typically:
	//   - Bitcoin: Double SHA256 of transaction
	//   - Ethereum: Keccak256 of transaction
	// We use signature prefix for simplicity.
	return ""
}

// MultiSigTransaction represents a transaction requiring multiple signatures
type MultiSigTransaction struct {
	Transaction Transaction       `json:"transaction"`
	Signatures  map[string]string `json:"signatures"` // publicKey -> signature
	Required    int               `json:"required"`   // M of N signatures needed
}

// NewMultiSigTransaction creates a new multi-signature transaction.
//
// Parameters:
//   - tx: The transaction to be signed by multiple parties
//   - required: Number of signatures required (M in "M of N")
//
// Returns:
//   - *MultiSigTransaction: Multi-sig transaction ready for signing
func NewMultiSigTransaction(tx *Transaction, required int) *MultiSigTransaction {
	// TODO: Implement multi-sig transaction creation
	// Steps:
	// 1. Create MultiSigTransaction struct
	// 2. Set Transaction to *tx
	// 3. Initialize Signatures as empty map: make(map[string]string)
	// 4. Set Required to the specified value
	// 5. Return pointer to struct
	return nil
}

// AddSignature adds a signature to the multi-sig transaction.
//
// Parameters:
//   - wallet: Wallet signing the transaction
//
// Returns:
//   - error: Any error during signing
func (mst *MultiSigTransaction) AddSignature(wallet *Wallet) error {
	// TODO: Implement adding signature to multi-sig
	// Steps:
	// 1. Serialize transaction: mst.Transaction.Serialize()
	// 2. Handle serialization errors
	// 3. Sign: ed25519.Sign(wallet.PrivateKey, txBytes)
	// 4. Encode signature to hex
	// 5. Encode wallet's public key to hex
	// 6. Add to signatures map: mst.Signatures[pubKeyHex] = signatureHex
	// 7. Return nil
	return nil
}

// Verify verifies that the multi-sig transaction has enough valid signatures.
//
// Returns:
//   - bool: true if transaction has required number of valid signatures
//   - error: Any error during verification
func (mst *MultiSigTransaction) Verify() (bool, error) {
	// TODO: Implement multi-sig verification
	// Steps:
	// 1. Serialize transaction
	// 2. Initialize validCount := 0
	// 3. For each publicKey, signature in mst.Signatures:
	//    a. Decode public key from hex
	//    b. Decode signature from hex
	//    c. If ed25519.Verify succeeds, increment validCount
	// 4. Return (validCount >= mst.Required), nil
	return false, nil
}
