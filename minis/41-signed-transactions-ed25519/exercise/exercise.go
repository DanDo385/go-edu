//go:build !solution
// +build !solution

package exercise

// TODO: Import required packages
// You'll need:
// - "crypto/ed25519" for elliptic curve digital signatures
// - "encoding/hex" for hex encoding/decoding (convert bytes to readable strings)
// - "encoding/json" for transaction serialization
// - "time" for timestamping transactions
//
// import (
//     "crypto/ed25519"
//     "encoding/hex"
//     "encoding/json"
//     "time"
// )

// Transaction represents a blockchain transaction
// This is a STRUCT (value type), not a reference type
type Transaction struct {
	From      string  `json:"from"`      // Sender's public key (hex-encoded)
	To        string  `json:"to"`        // Recipient's address
	Amount    float64 `json:"amount"`    // Transfer amount
	Nonce     int64   `json:"nonce"`     // Unique transaction number (prevents replay attacks)
	Timestamp int64   `json:"timestamp"` // Unix timestamp (seconds since epoch)
}

// SignedTransaction contains a transaction and its cryptographic signature
// The signature proves the transaction was authorized by the private key holder
type SignedTransaction struct {
	Transaction Transaction `json:"transaction"` // Embedded struct (value, not pointer)
	Signature   string      `json:"signature"`   // Hex-encoded Ed25519 signature
	PublicKey   string      `json:"public_key"`  // Hex-encoded public key
}

// Wallet represents a cryptographic wallet with Ed25519 keypair
//
// Ed25519 Overview:
// - Elliptic Curve Digital Signature Algorithm (EdDSA)
// - Public key: 32 bytes (can be shared publicly)
// - Private key: 64 bytes (must be kept secret!)
// - Signature: 64 bytes
// - Fast, secure, and widely used (SSH, TLS, blockchain)
type Wallet struct {
	PublicKey  ed25519.PublicKey  // 32-byte public key (value type: []byte)
	PrivateKey ed25519.PrivateKey // 64-byte private key (value type: []byte)
	Address    string             // Human-readable address (derived from public key)
}

// GenerateWallet creates a new wallet with a fresh Ed25519 keypair.
//
// Ed25519 Key Generation:
// - Uses cryptographically secure random number generator (CSPRNG)
// - Private key: 64 random bytes
// - Public key: derived from private key via curve multiplication
// - One-way operation: public key cannot reveal private key
//
// Returns:
//   - *Wallet: Pointer to new wallet (allows modification, avoids copying large keys)
//   - error: Any error during key generation (extremely rare)
//
// TODO: Implement GenerateWallet
// Function signature: func GenerateWallet() (*Wallet, error)
//
// Steps to implement:
// 1. Generate Ed25519 keypair
//    - Use: pub, priv, err := ed25519.GenerateKey(nil)
//    - First return: pub (ed25519.PublicKey, which is []byte)
//    - Second return: priv (ed25519.PrivateKey, which is []byte)
//    - Third return: err (error, will be nil unless crypto/rand fails)
//    - nil parameter means "use crypto/rand.Reader" (secure randomness)
//
// 2. Check for errors
//    - if err != nil { return nil, fmt.Errorf("failed to generate keypair: %w", err) }
//    - Extremely rare: only fails if OS random number generator is broken
//
// 3. Derive address from public key
//    - Use: address := DeriveAddress(pub)
//    - This creates a shorter, human-readable identifier
//
// 4. Create and return Wallet struct
//    - Return &Wallet{ PublicKey: pub, PrivateKey: priv, Address: address }, nil
//    - Why &Wallet? Returns a POINTER to avoid copying 96 bytes of keys
//    - The struct is allocated on the heap (Go escape analysis)
//
// Key Go concepts:
// - ed25519.PublicKey and PrivateKey are type aliases for []byte (slices)
// - Slices are reference types (contain pointer to underlying array)
// - Passing slices is cheap (copy pointer + len + cap, not all data)
// - Returning &Wallet allocates on heap (survives function return)

// TODO: Implement the GenerateWallet function below
// func GenerateWallet() (*Wallet, error) {
//     return nil, nil
// }

// DeriveAddress derives a human-readable address from a public key.
//
// Real-World Address Derivation:
// - Bitcoin: SHA256(public key) → RIPEMD160 → Base58Check
// - Ethereum: Keccak256(public key) → take last 20 bytes → hex encode
// - Solana: Base58 encode public key directly
//
// Our simplified version: Take first 20 bytes of public key → hex encode
//
// Parameters:
//   - publicKey: Ed25519 public key (32 bytes, type is []byte)
//
// Returns:
//   - string: Hex-encoded address (40 characters: 20 bytes * 2 hex digits)
//
// TODO: Implement DeriveAddress
// Function signature: func DeriveAddress(publicKey ed25519.PublicKey) string
//
// Steps to implement:
// 1. Extract first 20 bytes of public key
//    - Use slice syntax: addressBytes := publicKey[:20]
//    - This creates a NEW slice header pointing to same underlying array
//    - Range: index 0 (inclusive) to 20 (exclusive)
//    - No allocation: shares memory with original publicKey
//
// 2. Convert to hex string
//    - Use: hex.EncodeToString(addressBytes)
//    - Allocates a new string: 20 bytes → 40 hex characters
//    - Example: [0xAB, 0xCD] → "abcd"
//
// 3. Return the hex string
//
// Key Go concepts:
// - Slice syntax: s[:n] takes first n elements
// - hex.EncodeToString allocates (creates new string)
// - ed25519.PublicKey is just []byte (type alias for clarity)

// TODO: Implement the DeriveAddress function below
// func DeriveAddress(publicKey ed25519.PublicKey) string {
//     return ""
// }

// NewTransaction creates a new unsigned transaction.
//
// Transaction Fields Explained:
// - From: Sender's public key (proves who initiated the transaction)
// - To: Recipient's address (where funds go)
// - Amount: Transfer amount (in whatever unit the blockchain uses)
// - Nonce: Unique number (prevents replay attacks - same tx can't be submitted twice)
// - Timestamp: When transaction was created (helps with ordering and expiration)
//
// Parameters:
//   - from: Sender's wallet (passed by POINTER - we only need to read PublicKey)
//   - to: Recipient's address (passed by VALUE - string is cheap to copy)
//   - amount: Transfer amount (passed by VALUE - float64 is 8 bytes)
//   - nonce: Unique transaction number (passed by VALUE - int64 is 8 bytes)
//
// Returns:
//   - *Transaction: Pointer to new transaction (avoids copying 48+ bytes)
//
// TODO: Implement NewTransaction
// Function signature: func NewTransaction(from *Wallet, to string, amount float64, nonce int64) *Transaction
//
// Steps to implement:
// 1. Create Transaction struct with all fields
//    - From: hex.EncodeToString(from.PublicKey)
//      * Converts []byte public key to hex string
//      * Example: [0x12, 0x34] → "1234"
//    - To: to (recipient address, already a string)
//    - Amount: amount (transfer amount)
//    - Nonce: nonce (unique transaction number)
//    - Timestamp: time.Now().Unix()
//      * time.Now() returns current time (time.Time struct)
//      * .Unix() converts to seconds since January 1, 1970 (int64)
//
// 2. Return pointer to transaction
//    - Use: return &Transaction{ ... }
//    - The & (address-of operator) returns a pointer to the struct
//    - Struct is allocated on the heap (escape analysis)
//
// Key Go concepts:
// - time.Now() returns a value (time.Time struct)
// - .Unix() is a method on time.Time (returns int64)
// - &T{ fields } creates T on heap and returns *T
// - Pointer return avoids copying the struct

// TODO: Implement the NewTransaction function below
// func NewTransaction(from *Wallet, to string, amount float64, nonce int64) *Transaction {
//     return nil
// }

// Serialize converts a transaction to bytes for signing/verification.
//
// Why Serialization Matters:
// - Digital signatures sign BYTES, not Go structs
// - Must convert struct → deterministic byte sequence
// - Same serialization method MUST be used for sign and verify
// - Different serialization = different bytes = invalid signature!
//
// Serialization Options:
// - JSON: Human-readable, widely supported, but not the most compact
// - Protobuf: Compact, fast, but requires schema definition
// - Custom binary: Maximum control, but error-prone
//
// We use JSON for simplicity and debuggability.
//
// Returns:
//   - []byte: JSON-encoded transaction (allocated on heap)
//   - error: Any error during serialization (rare: only fails if custom JSON marshalers panic)
//
// TODO: Implement Serialize method
// Method signature: func (tx *Transaction) Serialize() ([]byte, error)
//
// Steps to implement:
// 1. Marshal transaction to JSON
//    - Use: data, err := json.Marshal(tx)
//    - json.Marshal takes interface{} (any type) and returns []byte
//    - Uses struct tags like `json:"from"` to control field names
//    - Allocates a new byte slice on the heap
//
// 2. Check for errors
//    - if err != nil { return nil, fmt.Errorf("failed to serialize: %w", err) }
//    - Rare: only fails if struct has circular references or custom marshalers panic
//
// 3. Return JSON bytes
//    - return data, nil
//
// Key Go concepts:
// - (tx *Transaction) is a method receiver (pointer type)
// - json.Marshal uses reflection to inspect struct fields
// - Struct tags control JSON field names and omitempty behavior
// - %w wraps error (allows errors.Is and errors.As to work)

// TODO: Implement the Serialize method below
// func (tx *Transaction) Serialize() ([]byte, error) {
//     return nil, nil
// }

// Sign signs a transaction with the wallet's private key.
//
// Digital Signature Process (Ed25519):
// 1. Serialize transaction → bytes
// 2. Hash the bytes (SHA-512 internally in Ed25519)
// 3. Use private key to generate signature (curve operations)
// 4. Return signature (64 bytes)
//
// Security Properties:
// - Authentication: Proves transaction came from private key holder
// - Integrity: Any modification to transaction invalidates signature
// - Non-repudiation: Signer cannot deny creating signature
// - Does NOT provide confidentiality: Transaction is still readable
//
// Parameters:
//   - tx: Transaction to sign (passed by POINTER - we only read it)
//
// Returns:
//   - *SignedTransaction: Transaction with signature and public key
//   - error: Any error during signing (serialization failure)
//
// TODO: Implement Sign method
// Method signature: func (w *Wallet) Sign(tx *Transaction) (*SignedTransaction, error)
//
// Steps to implement:
// 1. Serialize the transaction
//    - txBytes, err := tx.Serialize()
//    - if err != nil { return nil, err }
//
// 2. Sign the serialized bytes
//    - Use: signature := ed25519.Sign(w.PrivateKey, txBytes)
//    - First param: private key ([]byte, 64 bytes)
//    - Second param: message to sign ([]byte)
//    - Returns: signature ([]byte, 64 bytes)
//    - This is a VALUE return (copies 64 bytes)
//
// 3. Encode signature and public key to hex
//    - sigHex := hex.EncodeToString(signature)
//    - pubKeyHex := hex.EncodeToString(w.PublicKey)
//    - Hex encoding makes binary data human-readable and JSON-safe
//
// 4. Create and return SignedTransaction
//    - return &SignedTransaction{
//        Transaction: *tx,           // Dereference pointer to copy struct
//        Signature:   sigHex,
//        PublicKey:   pubKeyHex,
//      }, nil
//
// Key Go concepts:
// - ed25519.Sign takes private key and message, returns signature
// - All are []byte slices (reference types, but values are copied)
// - *tx dereferences pointer to get the Transaction value
// - Embedding Transaction (not *Transaction) means copy, not reference

// TODO: Implement the Sign method below
// func (w *Wallet) Sign(tx *Transaction) (*SignedTransaction, error) {
//     return nil, nil
// }

// Verify verifies the signature on a signed transaction.
//
// Verification Process (Ed25519):
// 1. Deserialize signature and public key from hex
// 2. Serialize transaction (same method as signing!)
// 3. Use ed25519.Verify(publicKey, message, signature)
// 4. Return true if signature is valid, false otherwise
//
// What Verification Proves:
// - The transaction was signed by the holder of the private key
// - The transaction has not been modified since signing
// - The signature is mathematically valid
//
// What Verification Does NOT Prove:
// - That the signer is authorized to send funds
// - That the sender has sufficient balance
// - That the recipient address is valid
//
// Returns:
//   - bool: true if signature is valid, false otherwise
//   - error: Any error during verification (e.g., invalid hex encoding)
//
// TODO: Implement Verify method
// Method signature: func (st *SignedTransaction) Verify() (bool, error)
//
// Steps to implement:
// 1. Decode signature from hex
//    - signature, err := hex.DecodeString(st.Signature)
//    - if err != nil { return false, fmt.Errorf("invalid signature encoding: %w", err) }
//    - hex.DecodeString returns ([]byte, error)
//    - Allocates a new byte slice with decoded data
//
// 2. Decode public key from hex
//    - publicKey, err := hex.DecodeString(st.PublicKey)
//    - if err != nil { return false, fmt.Errorf("invalid public key encoding: %w", err) }
//
// 3. Serialize transaction
//    - txBytes, err := st.Transaction.Serialize()
//    - if err != nil { return false, err }
//    - CRITICAL: Must use SAME serialization method as Sign!
//
// 4. Verify signature
//    - valid := ed25519.Verify(publicKey, txBytes, signature)
//    - First param: public key ([]byte, 32 bytes)
//    - Second param: message that was signed ([]byte)
//    - Third param: signature to verify ([]byte, 64 bytes)
//    - Returns: bool (true if valid, false otherwise)
//
// 5. Return result
//    - return valid, nil
//
// Key Go concepts:
// - hex.DecodeString is inverse of hex.EncodeToString
// - ed25519.Verify is the verification counterpart to ed25519.Sign
// - Must preserve same serialization for sign and verify
// - Multiple return values: (result, error) is idiomatic Go

// TODO: Implement the Verify method below
// func (st *SignedTransaction) Verify() (bool, error) {
//     return false, nil
// }

// VerifyOwnership checks if the transaction's 'from' field matches
// the public key that signed it.
//
// Why This Matters:
// - Someone could create a transaction with from="alice" but sign it with bob's key
// - The signature would be valid (bob signed it), but bob doesn't own alice's funds!
// - This check prevents such attacks
//
// Security Note:
// - ALWAYS verify both: 1) signature is valid, 2) signer owns the funds
// - Verifying signature alone is NOT sufficient!
//
// Returns:
//   - bool: true if 'from' matches signing public key, false otherwise
//
// TODO: Implement VerifyOwnership method
// Method signature: func (st *SignedTransaction) VerifyOwnership() bool
//
// Steps to implement:
// 1. Compare st.Transaction.From with st.PublicKey
//    - Both are hex-encoded strings
//    - Use: return st.Transaction.From == st.PublicKey
//
// Key Go concepts:
// - String comparison with == compares contents, not pointers
// - Simple one-liner, but critical for security!

// TODO: Implement the VerifyOwnership method below
// func (st *SignedTransaction) VerifyOwnership() bool {
//     return false
// }

// GetTransactionID computes a unique ID for this transaction.
//
// Transaction ID Purpose:
// - Unique identifier for the transaction
// - Used for indexing, lookup, and deduplication
// - Prevents same transaction from being included twice
//
// Real-World Transaction IDs:
// - Bitcoin: Double SHA256 of serialized transaction
// - Ethereum: Keccak256 of serialized transaction
// - Our simplified version: First 16 bytes of signature (unique enough)
//
// Returns:
//   - string: Transaction ID (32 hex characters)
//
// TODO: Implement GetTransactionID method
// Method signature: func (st *SignedTransaction) GetTransactionID() string
//
// Steps to implement:
// 1. Decode signature from hex
//    - signature, err := hex.DecodeString(st.Signature)
//    - if err != nil || len(signature) < 16 { return "" }
//    - Check length to avoid panic on signature[:16]
//
// 2. Take first 16 bytes
//    - idBytes := signature[:16]
//    - Slice syntax: start at 0, stop before 16
//
// 3. Encode to hex
//    - return hex.EncodeToString(idBytes)
//    - 16 bytes → 32 hex characters
//
// Key Go concepts:
// - Slicing: s[:n] creates new slice pointing to first n elements
// - Length check prevents panic (Go panics on out-of-bounds access)
// - hex encoding makes binary data readable and safe for URLs/JSON

// TODO: Implement the GetTransactionID method below
// func (st *SignedTransaction) GetTransactionID() string {
//     return ""
// }

// MultiSigTransaction represents a transaction requiring multiple signatures.
//
// Multi-Signature Explained:
// - M-of-N signatures: Need M valid signatures from N possible signers
// - Example: 2-of-3 means any 2 of 3 people can authorize transaction
// - Use cases: Corporate accounts, escrow, joint accounts
//
// Security Benefits:
// - No single point of failure (one compromised key doesn't lose funds)
// - Requires collusion to steal funds
// - Can implement complex authorization policies
type MultiSigTransaction struct {
	Transaction Transaction       `json:"transaction"`
	Signatures  map[string]string `json:"signatures"` // publicKey (hex) → signature (hex)
	Required    int               `json:"required"`   // M of N signatures needed
}

// NewMultiSigTransaction creates a new multi-signature transaction.
//
// Parameters:
//   - tx: Transaction to be signed by multiple parties
//   - required: Number of signatures required (M in "M of N")
//
// Returns:
//   - *MultiSigTransaction: Multi-sig transaction ready for signing
//
// TODO: Implement NewMultiSigTransaction
// Function signature: func NewMultiSigTransaction(tx *Transaction, required int) *MultiSigTransaction
//
// Steps to implement:
// 1. Create MultiSigTransaction struct
//    - Transaction: *tx (dereference pointer to copy struct)
//    - Signatures: make(map[string]string) (initialize empty map)
//    - Required: required
//
// 2. Return pointer to struct
//    - return &MultiSigTransaction{ ... }
//
// Key Go concepts:
// - Maps must be initialized with make() or literal {}
// - nil maps cannot be written to (panic!)
// - map[string]string is a reference type (stores pointer to hash table)

// TODO: Implement the NewMultiSigTransaction function below
// func NewMultiSigTransaction(tx *Transaction, required int) *MultiSigTransaction {
//     return nil
// }

// AddSignature adds a signature to the multi-sig transaction.
//
// Process:
// 1. Wallet signs the transaction (same as regular signing)
// 2. Add signature to map keyed by public key
// 3. Multiple wallets can call this to add their signatures
//
// Parameters:
//   - wallet: Wallet signing the transaction
//
// Returns:
//   - error: Any error during signing
//
// TODO: Implement AddSignature method
// Method signature: func (mst *MultiSigTransaction) AddSignature(wallet *Wallet) error
//
// Steps to implement:
// 1. Serialize transaction
//    - txBytes, err := mst.Transaction.Serialize()
//    - if err != nil { return err }
//
// 2. Sign with wallet's private key
//    - signature := ed25519.Sign(wallet.PrivateKey, txBytes)
//
// 3. Encode to hex
//    - sigHex := hex.EncodeToString(signature)
//    - pubKeyHex := hex.EncodeToString(wallet.PublicKey)
//
// 4. Add to signatures map
//    - mst.Signatures[pubKeyHex] = sigHex
//    - Key is public key (identifies signer)
//    - Value is signature (proves authorization)
//
// 5. Return success
//    - return nil
//
// Key Go concepts:
// - Map assignment: m[key] = value
// - If key exists, overwrites; if not, creates new entry
// - mst is a pointer, so we're modifying the original struct

// TODO: Implement the AddSignature method below
// func (mst *MultiSigTransaction) AddSignature(wallet *Wallet) error {
//     return nil
// }

// Verify verifies that the multi-sig transaction has enough valid signatures.
//
// Verification Process:
// 1. Serialize transaction
// 2. For each signature in map:
//    a. Decode public key and signature from hex
//    b. Verify signature using ed25519.Verify
//    c. Count valid signatures
// 3. Check if validCount >= Required
//
// Returns:
//   - bool: true if transaction has required number of valid signatures
//   - error: Any error during verification
//
// TODO: Implement Verify method for MultiSigTransaction
// Method signature: func (mst *MultiSigTransaction) Verify() (bool, error)
//
// Steps to implement:
// 1. Serialize transaction
//    - txBytes, err := mst.Transaction.Serialize()
//    - if err != nil { return false, err }
//
// 2. Initialize valid signature counter
//    - validCount := 0
//
// 3. Loop through signatures map
//    - for pubKeyHex, sigHex := range mst.Signatures { ... }
//    - range on map returns (key, value) pairs
//    - Order is random (maps are unordered in Go)
//
// 4. For each signature:
//    a. Decode public key
//       - pubKey, err := hex.DecodeString(pubKeyHex)
//       - if err != nil { continue } // Skip invalid encoding
//    b. Decode signature
//       - sig, err := hex.DecodeString(sigHex)
//       - if err != nil { continue }
//    c. Verify signature
//       - if ed25519.Verify(pubKey, txBytes, sig) { validCount++ }
//
// 5. Check if we have enough valid signatures
//    - return validCount >= mst.Required, nil
//
// Key Go concepts:
// - range on maps: for k, v := range m
// - continue skips to next iteration (don't count invalid signatures)
// - Graceful error handling: skip invalid entries instead of failing entire verification

// TODO: Implement the Verify method for MultiSigTransaction below
// func (mst *MultiSigTransaction) Verify() (bool, error) {
//     return false, nil
// }

// After implementing all functions:
// - Run: go test ./minis/41-signed-transactions-ed25519/exercise/...
// - Test solution: go test -tags solution ./minis/41-signed-transactions-ed25519/exercise/...
// - Compare with solution.go to see detailed explanations
//
// Key Cryptographic Concepts Learned:
// - Ed25519: Modern elliptic curve digital signature algorithm
// - Public/Private key pairs: Asymmetric cryptography
// - Digital signatures: Authentication + integrity + non-repudiation
// - Multi-sig: M-of-N signature schemes for shared control
// - Transaction serialization: Deterministic byte representation
//
// Key Go Concepts Learned:
// - crypto/ed25519 package: Modern cryptography in stdlib
// - hex encoding: Binary data → human-readable strings
// - json.Marshal: Struct → JSON bytes
// - Method receivers: (w *Wallet) vs (st SignedTransaction)
// - Maps: Reference types, must initialize with make()
// - Error wrapping: fmt.Errorf with %w
