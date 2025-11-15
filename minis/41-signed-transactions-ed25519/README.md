# Project 41: Signed Transactions with Ed25519

## 1. What Is This About?

### Real-World Scenario

You're building a cryptocurrency wallet or blockchain system:
- Users need to prove they own their digital assets
- Anyone should be able to verify that a transaction is legitimate
- Private keys must never be shared or transmitted
- Signatures must be impossible to forge

**❌ Naive approach:** Share your password with everyone to prove identity
- Password gets compromised
- Anyone can impersonate you
- No way to verify authenticity without trusting a central authority

**✅ Cryptographic signatures:** Use public-key cryptography
- Private key stays secret on your device
- Public key can be shared with everyone
- Sign transactions with private key
- Anyone can verify signature with public key
- Mathematically impossible to forge signatures without private key

This project teaches you **digital signatures** - the cryptographic foundation of:
- Cryptocurrencies (Bitcoin, Ethereum, Solana)
- Secure messaging (Signal, WhatsApp)
- Code signing (software distribution)
- SSL/TLS certificates (HTTPS)
- Git commit signing

### What You'll Learn

1. **Public-key cryptography**: Asymmetric encryption basics
2. **Ed25519**: Modern, fast, secure signature algorithm
3. **Keypair generation**: Creating private/public key pairs
4. **Digital signing**: Proving ownership without revealing secrets
5. **Signature verification**: Checking authenticity of signed data
6. **Transaction structure**: Building signed transaction objects

### The Challenge

Build a transaction signing system that:
- Generates cryptographically secure Ed25519 keypairs
- Signs arbitrary transaction data with private keys
- Verifies signatures using public keys
- Prevents signature forgery and replay attacks
- Encodes signatures in standard formats (hex, base64)
- Validates transaction integrity

---

## 2. First Principles: Understanding Digital Signatures

### What is Public-Key Cryptography?

**Public-key cryptography** uses two related but different keys:
- **Private key**: Secret, never shared (like a house key only you have)
- **Public key**: Shared openly (like your house address everyone can know)

**Analogy**: Think of it like a lockbox on your door:
- You have the **private key** that opens the lockbox
- Anyone can drop a message in (encrypt with your public key)
- Only you can open it and read it (decrypt with your private key)
- You can also leave a signed note (sign with private key)
- Anyone can verify it's really from you (verify with public key)

**Mathematical foundation**:
```
Private Key (secret)  →  Public Key (derived, can be shared)
                      ↓
                   One-way function
                  (easy to compute →, impossible to reverse ←)
```

### What is a Digital Signature?

A **digital signature** is like a handwritten signature, but:
- **Unforgeable**: Impossible to create without the private key
- **Verifiable**: Anyone with the public key can verify it
- **Non-repudiable**: Signer can't deny signing it
- **Tamper-evident**: Any modification invalidates the signature

**How it works**:
```
Original Message: "Send 10 BTC to Alice"
         ↓
Hash the message (SHA-256)
         ↓
Hash: 0x3a4b5c... (32 bytes)
         ↓
Sign hash with private key (Ed25519)
         ↓
Signature: 0x7f8e9d... (64 bytes)
```

**Verification**:
```
Received: Message + Signature + Public Key
         ↓
Hash the message
         ↓
Verify signature matches hash using public key
         ↓
Valid ✓ or Invalid ✗
```

### What is Ed25519?

**Ed25519** is a modern digital signature algorithm based on elliptic curve cryptography.

**Key properties**:
- **Fast**: 50,000+ signatures/sec on modern CPUs
- **Small keys**: 32-byte public keys, 64-byte signatures
- **Secure**: 128-bit security level (equivalent to AES-128)
- **Deterministic**: Same message always produces same signature
- **Side-channel resistant**: Safe against timing attacks

**Comparison with other algorithms**:

| Algorithm | Public Key Size | Signature Size | Speed | Security Level |
|-----------|----------------|----------------|-------|----------------|
| RSA-2048 | 256 bytes | 256 bytes | Slow | 112-bit |
| ECDSA P-256 | 64 bytes | 64 bytes | Medium | 128-bit |
| **Ed25519** | **32 bytes** | **64 bytes** | **Fast** | **128-bit** |

**Why Ed25519 is popular**:
- Used by Bitcoin/Ethereum (via similar curves)
- SSH keys (`ssh-keygen -t ed25519`)
- Signal protocol (secure messaging)
- Tor hidden services
- Cryptocurrency wallets

### How Ed25519 Signing Works

**Mathematical details** (simplified):

Ed25519 uses the twisted Edwards curve:
```
-x² + y² = 1 + dx²y²
```

**Signing algorithm**:
```
1. Generate random 32-byte private key: sk
2. Derive public key: pk = sk × G (scalar multiplication on curve)
3. Hash message: h = SHA512(message)
4. Compute signature: (R, S) where:
   R = random point on curve
   S = secret nonce + hash × private key (mod order)
5. Signature = (R, S) = 64 bytes
```

**Verification algorithm**:
```
1. Parse signature: (R, S)
2. Hash message: h = SHA512(message)
3. Compute: P = S × G - h × pk
4. Check: P == R?
   If yes → Valid signature ✓
   If no → Invalid signature ✗
```

**Key insight**:
- Signing requires private key (only signer can do this)
- Verification requires public key (anyone can do this)
- Forging requires solving discrete logarithm problem (mathematically infeasible)

### Transaction Signing in Practice

A **signed transaction** typically contains:
```go
type SignedTransaction struct {
    // Transaction data
    From      string  // Sender's public key
    To        string  // Recipient's address
    Amount    float64 // Transfer amount
    Nonce     int64   // Unique transaction number (prevents replay)
    Timestamp int64   // Unix timestamp

    // Cryptographic proof
    Signature []byte  // Ed25519 signature (64 bytes)
}
```

**Why nonce?** Prevents **replay attacks**:
```
Without nonce:
  Alice signs: "Send 10 BTC to Bob"
  Attacker intercepts and replays this transaction 100 times
  Alice loses 1000 BTC! ✗

With nonce:
  Alice signs: "Send 10 BTC to Bob (nonce=1)"
  Attacker replays transaction
  Network rejects: "nonce=1 already used" ✓
```

---

## 3. Breaking Down the Solution

### Step 1: Generate Keypair

```go
import "crypto/ed25519"

// Generate new keypair
publicKey, privateKey, err := ed25519.GenerateKey(nil)
// publicKey:  32 bytes (can share)
// privateKey: 64 bytes (keep secret!)
```

**What happens internally**:
1. Read 32 random bytes from OS entropy source (`/dev/urandom`)
2. This is your private key (seed)
3. Derive public key using curve point multiplication
4. Return both keys

**Encoding keys for storage**:
```go
import "encoding/hex"

// Encode to hex (human-readable)
pubHex := hex.EncodeToString(publicKey)
// "a3b4c5d6e7f8..."

// Decode from hex
pubBytes, _ := hex.DecodeString(pubHex)
```

### Step 2: Create Transaction

```go
type Transaction struct {
    From      string
    To        string
    Amount    float64
    Nonce     int64
    Timestamp int64
}

tx := Transaction{
    From:      hex.EncodeToString(publicKey),
    To:        "recipient_address",
    Amount:    10.5,
    Nonce:     1,
    Timestamp: time.Now().Unix(),
}
```

### Step 3: Serialize for Signing

**Critical**: Must serialize deterministically (same data → same bytes)

```go
import "encoding/json"

// Serialize to JSON (canonical format)
txBytes, err := json.Marshal(tx)
// {"From":"a3b4c5...","To":"recipient...","Amount":10.5,...}
```

**Why JSON?**
- Human-readable
- Deterministic ordering (in Go's encoding/json)
- Standard format

**Alternative**: Use hash of transaction
```go
import "crypto/sha256"

hash := sha256.Sum256(txBytes)
// Sign the hash instead of full transaction
// Faster for large transactions
```

### Step 4: Sign Transaction

```go
signature := ed25519.Sign(privateKey, txBytes)
// signature is 64 bytes
```

**What happens internally**:
1. Derive signing key from private key
2. Generate random nonce (deterministically from message)
3. Compute R point on curve
4. Compute S scalar
5. Return (R || S) as 64-byte signature

### Step 5: Verify Signature

```go
valid := ed25519.Verify(publicKey, txBytes, signature)
if valid {
    fmt.Println("Signature is valid ✓")
    // Transaction is authentic
} else {
    fmt.Println("Signature is invalid ✗")
    // Transaction was tampered with or forged
}
```

**What happens internally**:
1. Parse signature into (R, S)
2. Hash message
3. Perform elliptic curve operations
4. Check if computed point equals R

### Step 6: Complete Signed Transaction

```go
type SignedTransaction struct {
    Transaction
    Signature []byte
}

signedTx := SignedTransaction{
    Transaction: tx,
    Signature:   signature,
}

// Serialize for transmission
txJSON, _ := json.Marshal(signedTx)
// Send to blockchain network...
```

---

## 4. Complete Solution Walkthrough

### Keypair Generation

```go
package main

import (
    "crypto/ed25519"
    "encoding/hex"
    "fmt"
)

func GenerateKeypair() (publicKey ed25519.PublicKey, privateKey ed25519.PrivateKey) {
    pub, priv, err := ed25519.GenerateKey(nil)
    if err != nil {
        panic(err) // Should never happen
    }
    return pub, priv
}

func main() {
    pub, priv := GenerateKeypair()

    fmt.Printf("Public key:  %s\n", hex.EncodeToString(pub))
    fmt.Printf("Private key: %s\n", hex.EncodeToString(priv))

    fmt.Printf("\nKey sizes:\n")
    fmt.Printf("  Public:  %d bytes\n", len(pub))  // 32
    fmt.Printf("  Private: %d bytes\n", len(priv)) // 64
}
```

**Output**:
```
Public key:  a3b4c5d6e7f8... (64 hex chars = 32 bytes)
Private key: 1a2b3c4d5e6f... (128 hex chars = 64 bytes)

Key sizes:
  Public:  32 bytes
  Private: 64 bytes
```

### Transaction Signing

```go
type Transaction struct {
    From      string  `json:"from"`
    To        string  `json:"to"`
    Amount    float64 `json:"amount"`
    Nonce     int64   `json:"nonce"`
    Timestamp int64   `json:"timestamp"`
}

func (tx *Transaction) Serialize() []byte {
    data, err := json.Marshal(tx)
    if err != nil {
        panic(err)
    }
    return data
}

func SignTransaction(tx *Transaction, privateKey ed25519.PrivateKey) []byte {
    txBytes := tx.Serialize()
    signature := ed25519.Sign(privateKey, txBytes)
    return signature
}

func VerifyTransaction(tx *Transaction, signature []byte, publicKey ed25519.PublicKey) bool {
    txBytes := tx.Serialize()
    return ed25519.Verify(publicKey, txBytes, signature)
}
```

### Complete Example

```go
func main() {
    // 1. Generate keypair
    pub, priv := GenerateKeypair()
    fmt.Println("Generated keypair")

    // 2. Create transaction
    tx := Transaction{
        From:      hex.EncodeToString(pub),
        To:        "recipient_address_xyz",
        Amount:    10.5,
        Nonce:     1,
        Timestamp: time.Now().Unix(),
    }
    fmt.Printf("Transaction: %+v\n", tx)

    // 3. Sign transaction
    signature := SignTransaction(&tx, priv)
    fmt.Printf("Signature: %s\n", hex.EncodeToString(signature))

    // 4. Verify signature
    valid := VerifyTransaction(&tx, signature, pub)
    fmt.Printf("Signature valid: %v\n", valid)

    // 5. Test tampering
    tx.Amount = 1000.0 // Modify transaction
    valid = VerifyTransaction(&tx, signature, pub)
    fmt.Printf("Signature valid after tampering: %v\n", valid)
}
```

**Output**:
```
Generated keypair
Transaction: {From:a3b4c5... To:recipient_address_xyz Amount:10.5 Nonce:1 Timestamp:1699564800}
Signature: 7f8e9d2c3b4a...
Signature valid: true
Signature valid after tampering: false
```

---

## 5. Key Concepts Explained

### Concept 1: Why 32-byte Public Key, 64-byte Private Key?

**Ed25519 key structure**:
```
Private Key (64 bytes):
  [0:32]  → Seed (actual secret)
  [32:64] → Public key (precomputed for speed)

Public Key (32 bytes):
  [0:32]  → Point on curve (compressed format)
```

**Why include public key in private key?**
- Performance: Signing needs both keys
- Convenience: One value contains everything

### Concept 2: Deterministic Signatures

Ed25519 signatures are **deterministic**:
```go
msg := []byte("hello")
sig1 := ed25519.Sign(priv, msg)
sig2 := ed25519.Sign(priv, msg)
// sig1 == sig2 (always!)
```

**Contrast with ECDSA** (used in Bitcoin):
```go
// ECDSA (pseudocode)
sig1 := ecdsa.Sign(priv, msg) // Uses random nonce
sig2 := ecdsa.Sign(priv, msg) // Different nonce
// sig1 != sig2 (different each time)
```

**Why deterministic is better**:
- No need for random number generator
- Safer (bad RNG can leak private key)
- Reproducible (useful for testing)

### Concept 3: Signature Malleability

**Problem**: Some signature schemes allow modifying signatures without invalidating them.

**Example** (ECDSA):
```
Valid signature: (r, s)
Also valid: (r, -s mod n)
// Same message, different signature bytes!
```

**Why this is bad**:
- Transaction ID depends on signature
- Attacker can change signature → changes transaction ID
- Breaks transaction tracking

**Ed25519 solution**:
- Canonical signatures only
- Verification rejects non-canonical signatures
- No malleability issues

### Concept 4: Side-Channel Attacks

**Timing attack example**:
```go
// Vulnerable comparison (pseudocode)
func BadVerify(sig1, sig2 []byte) bool {
    for i := range sig1 {
        if sig1[i] != sig2[i] {
            return false // Returns early!
        }
    }
    return true
}
```

**Problem**: Time to return reveals which byte differs!
```
Attacker measures response time:
  sig[0] wrong: 10μs (fails immediately)
  sig[0] right, sig[1] wrong: 20μs (fails at byte 1)
  sig[0] right, sig[1] right, sig[2] wrong: 30μs
  → Can brute-force signature byte by byte!
```

**Ed25519 protection**:
- Constant-time operations
- Always processes entire signature
- No timing leaks

### Concept 5: Why Sign the Hash?

**Option 1: Sign full message**
```go
signature := ed25519.Sign(priv, largeMessage)
// For 1MB message, slow!
```

**Option 2: Sign hash of message**
```go
hash := sha256.Sum256(largeMessage)
signature := ed25519.Sign(priv, hash[:])
// Much faster!
```

**Verification**:
```go
// Verifier also hashes the message
hash := sha256.Sum256(receivedMessage)
valid := ed25519.Verify(pub, hash[:], signature)
```

**Trade-off**:
- Faster signing/verification
- Hash collisions break security (but SHA-256 is collision-resistant)
- Standard practice in cryptocurrencies

---

## 6. Common Patterns You Can Reuse

### Pattern 1: Keypair Management

```go
type Wallet struct {
    PublicKey  ed25519.PublicKey
    PrivateKey ed25519.PrivateKey
    Address    string // Derived from public key
}

func NewWallet() *Wallet {
    pub, priv, _ := ed25519.GenerateKey(nil)
    return &Wallet{
        PublicKey:  pub,
        PrivateKey: priv,
        Address:    DeriveAddress(pub),
    }
}

func DeriveAddress(pubKey ed25519.PublicKey) string {
    // Common pattern: Hash public key → address
    hash := sha256.Sum256(pubKey)
    return hex.EncodeToString(hash[:20]) // First 20 bytes
}

func (w *Wallet) Sign(data []byte) []byte {
    return ed25519.Sign(w.PrivateKey, data)
}
```

### Pattern 2: Multisignature Transactions

```go
type MultiSigTransaction struct {
    Transaction
    RequiredSigs int                    // M of N signatures
    Signatures   map[string][]byte      // pubkey → signature
}

func (mst *MultiSigTransaction) AddSignature(pubKey ed25519.PublicKey, sig []byte) {
    if mst.Signatures == nil {
        mst.Signatures = make(map[string][]byte)
    }
    mst.Signatures[hex.EncodeToString(pubKey)] = sig
}

func (mst *MultiSigTransaction) Verify(authorizedKeys []ed25519.PublicKey) bool {
    if len(mst.Signatures) < mst.RequiredSigs {
        return false
    }

    txBytes := mst.Transaction.Serialize()
    validCount := 0

    for _, pubKey := range authorizedKeys {
        pubHex := hex.EncodeToString(pubKey)
        if sig, ok := mst.Signatures[pubHex]; ok {
            if ed25519.Verify(pubKey, txBytes, sig) {
                validCount++
            }
        }
    }

    return validCount >= mst.RequiredSigs
}
```

### Pattern 3: Persistent Key Storage

```go
import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/json"
    "os"
)

type EncryptedKeystore struct {
    PublicKey     string `json:"public_key"`
    EncryptedKey  string `json:"encrypted_private_key"`
    Salt          string `json:"salt"`
}

func SaveKeypair(pub ed25519.PublicKey, priv ed25519.PrivateKey, password, filename string) error {
    // Derive encryption key from password
    key := deriveKey(password, generateSalt())

    // Encrypt private key with AES
    encrypted := encryptAES(priv, key)

    ks := EncryptedKeystore{
        PublicKey:    hex.EncodeToString(pub),
        EncryptedKey: hex.EncodeToString(encrypted),
    }

    data, _ := json.Marshal(ks)
    return os.WriteFile(filename, data, 0600) // 0600 = owner read/write only
}

func LoadKeypair(password, filename string) (ed25519.PublicKey, ed25519.PrivateKey, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, nil, err
    }

    var ks EncryptedKeystore
    json.Unmarshal(data, &ks)

    encrypted, _ := hex.DecodeString(ks.EncryptedKey)
    key := deriveKey(password, salt)
    priv := decryptAES(encrypted, key)
    pub, _ := hex.DecodeString(ks.PublicKey)

    return pub, priv, nil
}
```

### Pattern 4: Message Authentication Codes (MAC)

```go
// For authenticated messages without signing
type AuthenticatedMessage struct {
    Data      []byte
    MAC       []byte // Message Authentication Code
    PublicKey ed25519.PublicKey
}

func CreateAuthMessage(data []byte, priv ed25519.PrivateKey, pub ed25519.PublicKey) *AuthenticatedMessage {
    // Sign data
    mac := ed25519.Sign(priv, data)

    return &AuthenticatedMessage{
        Data:      data,
        MAC:       mac,
        PublicKey: pub,
    }
}

func (am *AuthenticatedMessage) Verify() bool {
    return ed25519.Verify(am.PublicKey, am.Data, am.MAC)
}
```

### Pattern 5: Batch Verification

```go
// Verify multiple signatures efficiently
type SignedMessage struct {
    Message   []byte
    Signature []byte
    PublicKey ed25519.PublicKey
}

func BatchVerify(messages []SignedMessage) []bool {
    results := make([]bool, len(messages))

    // Parallel verification
    var wg sync.WaitGroup
    for i, msg := range messages {
        wg.Add(1)
        go func(idx int, m SignedMessage) {
            defer wg.Done()
            results[idx] = ed25519.Verify(m.PublicKey, m.Message, m.Signature)
        }(i, msg)
    }
    wg.Wait()

    return results
}
```

---

## 7. Real-World Applications

### Cryptocurrency Wallets

**Companies**: Coinbase, MetaMask, Ledger

```go
type CryptoWallet struct {
    wallet *Wallet
    utxos  []UTXO // Unspent transaction outputs
}

func (cw *CryptoWallet) SendTransaction(to string, amount float64) (*SignedTransaction, error) {
    tx := Transaction{
        From:      cw.wallet.Address,
        To:        to,
        Amount:    amount,
        Nonce:     cw.getNextNonce(),
        Timestamp: time.Now().Unix(),
    }

    sig := cw.wallet.Sign(tx.Serialize())

    return &SignedTransaction{
        Transaction: tx,
        Signature:   sig,
    }, nil
}
```

### Code Signing

**Companies**: Apple, Microsoft, Google (app distribution)

```go
type SignedPackage struct {
    Binary    []byte
    Version   string
    Signature []byte
    Publisher ed25519.PublicKey
}

func SignSoftware(binary []byte, version string, priv ed25519.PrivateKey, pub ed25519.PublicKey) *SignedPackage {
    // Hash the binary
    hash := sha256.Sum256(binary)

    // Sign the hash
    sig := ed25519.Sign(priv, hash[:])

    return &SignedPackage{
        Binary:    binary,
        Version:   version,
        Signature: sig,
        Publisher: pub,
    }
}

func (sp *SignedPackage) Verify() bool {
    hash := sha256.Sum256(sp.Binary)
    return ed25519.Verify(sp.Publisher, hash[:], sp.Signature)
}
```

### Git Commit Signing

**Usage**: Verify commit authenticity

```bash
# Configure Git to use Ed25519
git config --global gpg.format ssh
git config --global user.signingkey ~/.ssh/id_ed25519.pub

# Sign commits
git commit -S -m "Important change"
```

**Go implementation**:
```go
type GitCommit struct {
    Tree      string
    Parent    string
    Author    string
    Committer string
    Message   string
    Signature []byte
}

func (gc *GitCommit) Sign(priv ed25519.PrivateKey) {
    commitData := fmt.Sprintf("tree %s\nparent %s\nauthor %s\ncommitter %s\n\n%s",
        gc.Tree, gc.Parent, gc.Author, gc.Committer, gc.Message)
    gc.Signature = ed25519.Sign(priv, []byte(commitData))
}
```

### Secure Messaging

**Companies**: Signal, WhatsApp (via similar crypto)

```go
type SecureMessage struct {
    From      ed25519.PublicKey
    To        ed25519.PublicKey
    Content   []byte
    Signature []byte
    Timestamp int64
}

func SendSecureMessage(content []byte, senderPriv ed25519.PrivateKey, senderPub, recipientPub ed25519.PublicKey) *SecureMessage {
    msg := SecureMessage{
        From:      senderPub,
        To:        recipientPub,
        Content:   content,
        Timestamp: time.Now().Unix(),
    }

    msgBytes, _ := json.Marshal(msg)
    msg.Signature = ed25519.Sign(senderPriv, msgBytes)

    return &msg
}
```

### API Authentication

**Pattern**: Sign API requests

```go
type SignedAPIRequest struct {
    Method    string
    Path      string
    Body      []byte
    Timestamp int64
    APIKey    string // Public key
    Signature []byte
}

func CreateSignedRequest(method, path string, body []byte, priv ed25519.PrivateKey, pub ed25519.PublicKey) *SignedAPIRequest {
    req := SignedAPIRequest{
        Method:    method,
        Path:      path,
        Body:      body,
        Timestamp: time.Now().Unix(),
        APIKey:    hex.EncodeToString(pub),
    }

    // Sign canonical request
    canonical := fmt.Sprintf("%s\n%s\n%d\n%s", method, path, req.Timestamp, body)
    req.Signature = ed25519.Sign(priv, []byte(canonical))

    return &req
}
```

---

## 8. Common Mistakes to Avoid

### Mistake 1: Exposing Private Keys

**❌ Wrong**:
```go
func GetWalletInfo() map[string]string {
    return map[string]string{
        "public_key":  hex.EncodeToString(wallet.PublicKey),
        "private_key": hex.EncodeToString(wallet.PrivateKey), // NEVER!
    }
}
```

**✅ Correct**:
```go
func GetWalletInfo() map[string]string {
    return map[string]string{
        "public_key": hex.EncodeToString(wallet.PublicKey),
        "address":    wallet.Address,
        // Private key never leaves the wallet
    }
}
```

### Mistake 2: Signing Without Nonce

**❌ Wrong**:
```go
type Transaction struct {
    From   string
    To     string
    Amount float64
    // Missing nonce - vulnerable to replay attacks!
}
```

**✅ Correct**:
```go
type Transaction struct {
    From      string
    To        string
    Amount    float64
    Nonce     int64  // Prevents replay attacks
    Timestamp int64  // Additional uniqueness
}
```

### Mistake 3: Not Validating Before Signing

**❌ Wrong**:
```go
func SignTransaction(tx *Transaction, priv ed25519.PrivateKey) []byte {
    // Sign without validation
    return ed25519.Sign(priv, tx.Serialize())
}
```

**✅ Correct**:
```go
func SignTransaction(tx *Transaction, priv ed25519.PrivateKey) ([]byte, error) {
    // Validate first
    if err := tx.Validate(); err != nil {
        return nil, fmt.Errorf("invalid transaction: %w", err)
    }

    return ed25519.Sign(priv, tx.Serialize()), nil
}

func (tx *Transaction) Validate() error {
    if tx.Amount <= 0 {
        return errors.New("amount must be positive")
    }
    if tx.From == "" || tx.To == "" {
        return errors.New("from and to addresses required")
    }
    return nil
}
```

### Mistake 4: Modifying After Signing

**❌ Wrong**:
```go
signedTx := SignTransaction(&tx, priv)
tx.Amount = 100.0 // Modify after signing!
// Signature now invalid
```

**✅ Correct**:
```go
// Make transaction immutable after signing
type SignedTransaction struct {
    tx        Transaction // Unexported
    Signature []byte
}

func (st *SignedTransaction) GetTransaction() Transaction {
    return st.tx // Return copy, not reference
}
```

### Mistake 5: Using Wrong Encoding

**❌ Wrong**:
```go
// Different encodings produce different bytes!
txJSON := json.Marshal(tx)    // {"amount":10.5}
txGOB := gob.Encode(tx)        // Binary format (different!)
// Signatures won't match!
```

**✅ Correct**:
```go
// Always use the SAME encoding for sign and verify
func (tx *Transaction) Serialize() []byte {
    data, _ := json.Marshal(tx) // Consistent encoding
    return data
}
```

### Mistake 6: Not Checking Signature Length

**❌ Wrong**:
```go
func Verify(tx *Transaction, sig []byte, pub ed25519.PublicKey) bool {
    return ed25519.Verify(pub, tx.Serialize(), sig)
    // Panics if sig is wrong length!
}
```

**✅ Correct**:
```go
func Verify(tx *Transaction, sig []byte, pub ed25519.PublicKey) (bool, error) {
    if len(sig) != ed25519.SignatureSize {
        return false, fmt.Errorf("invalid signature length: %d", len(sig))
    }
    return ed25519.Verify(pub, tx.Serialize(), sig), nil
}
```

### Mistake 7: Storing Keys in Plain Text

**❌ Wrong**:
```go
// config.json
{
    "private_key": "abc123..." // Plain text!
}
```

**✅ Correct**:
```go
// Encrypt private keys before storing
type EncryptedWallet struct {
    PublicKey    string
    EncryptedKey string // Encrypted with password
    Salt         string
}
```

---

## 9. Stretch Goals

### Goal 1: Implement Key Derivation ⭐

Derive multiple keys from a single seed (HD wallets).

**Hint**:
```go
import "crypto/sha512"

func DeriveKey(seed []byte, index int) (ed25519.PublicKey, ed25519.PrivateKey) {
    // BIP32-style derivation
    data := append(seed, intToBytes(index)...)
    hash := sha512.Sum512(data)

    pub, priv, _ := ed25519.GenerateKey(bytes.NewReader(hash[:32]))
    return pub, priv
}
```

### Goal 2: Add Threshold Signatures ⭐⭐

Require M of N signatures to authorize transaction.

**Hint**:
```go
type ThresholdWallet struct {
    Threshold   int                     // M signatures required
    Participants []ed25519.PublicKey    // N total keys
}

func (tw *ThresholdWallet) Sign(tx *Transaction, signers []ed25519.PrivateKey) (*MultiSigTx, error) {
    if len(signers) < tw.Threshold {
        return nil, errors.New("not enough signers")
    }

    sigs := make([][]byte, len(signers))
    for i, priv := range signers {
        sigs[i] = ed25519.Sign(priv, tx.Serialize())
    }

    return &MultiSigTx{Transaction: tx, Signatures: sigs}, nil
}
```

### Goal 3: Implement Time-Locked Transactions ⭐⭐

Transactions that become valid only after a specific time.

**Hint**:
```go
type TimeLockedTransaction struct {
    Transaction
    UnlockTime int64 // Unix timestamp
}

func (tlt *TimeLockedTransaction) IsUnlocked() bool {
    return time.Now().Unix() >= tlt.UnlockTime
}

func (tlt *TimeLockedTransaction) Verify(sig []byte, pub ed25519.PublicKey) bool {
    if !tlt.IsUnlocked() {
        return false
    }
    return ed25519.Verify(pub, tlt.Serialize(), sig)
}
```

### Goal 4: Build a Simple Blockchain ⭐⭐⭐

Chain transactions together with signatures.

**Hint**:
```go
type Block struct {
    Index        int
    Timestamp    int64
    Transactions []SignedTransaction
    PreviousHash []byte
    Hash         []byte
}

func (b *Block) CalculateHash() []byte {
    data, _ := json.Marshal(b)
    hash := sha256.Sum256(data)
    return hash[:]
}

func (b *Block) VerifyTransactions() bool {
    for _, tx := range b.Transactions {
        pub, _ := hex.DecodeString(tx.From)
        if !ed25519.Verify(pub, tx.Transaction.Serialize(), tx.Signature) {
            return false
        }
    }
    return true
}
```

### Goal 5: Implement BIP39 Mnemonic Seeds ⭐⭐⭐

Generate keys from human-readable word phrases.

**Hint**:
```go
// Use github.com/tyler-smith/go-bip39

func GenerateFromMnemonic(mnemonic string) (ed25519.PublicKey, ed25519.PrivateKey, error) {
    seed := bip39.NewSeed(mnemonic, "")

    // Derive Ed25519 key from BIP39 seed
    hash := sha512.Sum512(seed)
    pub, priv, _ := ed25519.GenerateKey(bytes.NewReader(hash[:32]))

    return pub, priv, nil
}

// Example:
// mnemonic := "witch collapse practice feed shame open despair creek road again ice least"
// pub, priv, _ := GenerateFromMnemonic(mnemonic)
```

---

## How to Run

```bash
# Run the demo
go run ./minis/41-signed-transactions-ed25519/cmd/sign-demo/main.go

# Run exercises
go test ./minis/41-signed-transactions-ed25519/exercise/...

# Run with solution
go test -tags solution ./minis/41-signed-transactions-ed25519/exercise/...

# Benchmark signature performance
go test -bench=. -benchmem ./minis/41-signed-transactions-ed25519/exercise/...
```

---

## Summary

**What you learned**:
- ✅ Public-key cryptography fundamentals
- ✅ Ed25519 signature algorithm
- ✅ Keypair generation and management
- ✅ Transaction signing and verification
- ✅ Replay attack prevention with nonces
- ✅ Signature encoding and serialization

**Why this matters**:
Digital signatures are the foundation of modern security. Every cryptocurrency transaction, HTTPS connection, code deployment, and secure message relies on these principles. Ed25519 is the modern standard, used by billions of devices daily.

**Key insights**:
- Private keys must **never** be shared or transmitted
- Signatures prove ownership **without revealing secrets**
- Nonces prevent **replay attacks**
- Ed25519 is **fast, secure, and compact**
- Same serialization for **sign and verify** is critical

**Real-world impact**:
- Bitcoin/Ethereum process **millions of signed transactions daily**
- Signal/WhatsApp secure **billions of messages** with signatures
- GitHub verifies **thousands of signed commits per minute**
- Ed25519 enables **trust without central authority**

**Next steps**:
- Project 42: Build block structures with hashing
- Project 43: Implement proof-of-work consensus
- Project 44: Create transaction mempools
- Combine all concepts to build a complete blockchain!

Master cryptographic signatures, secure the decentralized future! 🔐
