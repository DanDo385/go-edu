//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package signedtransactionsed25519

import "crypto/ed25519"

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
// TODO: implement GenerateWallet.
func GenerateWallet() (*Wallet, error) { panic("TODO: implement") }
// TODO: implement DeriveAddress.
func DeriveAddress(publicKey ed25519.PublicKey) string { panic("TODO: implement") }
// TODO: implement NewTransaction.
func NewTransaction(from *Wallet, to string, amount float64, nonce int64) *Transaction {
	panic("TODO: implement")
}
// TODO: implement Serialize.
func (tx *Transaction) Serialize() ([]byte, error) { panic("TODO: implement") }
// TODO: implement Sign.
func (w *Wallet) Sign(tx *Transaction) (*SignedTransaction, error) { panic("TODO: implement") }
// TODO: implement Verify.
func (st *SignedTransaction) Verify() (bool, error) { panic("TODO: implement") }
// TODO: implement VerifyOwnership.
func (st *SignedTransaction) VerifyOwnership() bool { panic("TODO: implement") }
// TODO: implement GetTransactionID.
func (st *SignedTransaction) GetTransactionID() string { panic("TODO: implement") }

type MultiSigTransaction struct {
	Transaction Transaction       `json:"transaction"`
	Signatures  map[string]string `json:"signatures"` // publicKey -> signature
	Required    int               `json:"required"`   // M of N signatures needed
}
// TODO: implement NewMultiSigTransaction.
func NewMultiSigTransaction(tx *Transaction, required int) *MultiSigTransaction {
	panic("TODO: implement")
}
// TODO: implement AddSignature.
func (mst *MultiSigTransaction) AddSignature(wallet *Wallet) error { panic("TODO: implement") }
// TODO: implement Verify.
func (mst *MultiSigTransaction) Verify() (bool, error) { panic("TODO: implement") }
