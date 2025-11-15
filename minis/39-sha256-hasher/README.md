# Project 39: SHA-256 Hasher

## 1. What Is This About?

### Real-World Scenario

Imagine you're building a system that needs to:
- Verify file integrity after downloads (ensure no corruption)
- Store passwords securely (never store plaintext!)
- Create content-addressable storage (like Git)
- Detect duplicate files across millions of documents
- Generate unique identifiers for blockchain transactions

All these scenarios require **cryptographic hash functions**, and SHA-256 is one of the most widely used.

### What You'll Learn

1. **Cryptographic hashing**: One-way functions that map arbitrary data to fixed-size digests
2. **SHA-256 algorithm**: Industry-standard secure hash algorithm producing 256-bit hashes
3. **Hex encoding**: Converting binary hashes to readable hexadecimal strings
4. **Hash properties**: Determinism, collision resistance, avalanche effect
5. **Practical applications**: File verification, password hashing patterns, deduplication
6. **Digest reuse**: Efficient incremental hashing for streaming data

### The Challenge

Build a command-line tool that can:
- Hash strings and produce SHA-256 digests
- Hash files efficiently (even large files)
- Demonstrate incremental hashing (streaming)
- Verify file integrity using checksums
- Show the avalanche effect (tiny input changes → massive hash changes)

---

## 2. First Principles: Understanding Cryptographic Hashing

### What Is a Hash Function?

A **hash function** is a mathematical function that takes an input (of any size) and produces a fixed-size output called a **digest** or **hash**.

**Simple Analogy**: Think of a hash function like a fingerprint machine:
- You put your hand in (input: your unique hand)
- It produces a fingerprint (output: fixed-size identifier)
- Same hand always produces same fingerprint (deterministic)
- Fingerprint is unique enough that two different hands won't produce the same print (collision resistant)
- You can't reconstruct the hand from just the fingerprint (one-way)

**Mathematical definition**:
```
H: {0,1}* → {0,1}^n

Where:
- {0,1}* represents inputs of any length (arbitrary binary data)
- {0,1}^n represents outputs of fixed length n bits (for SHA-256, n=256)
- H is the hash function
```

### Properties of Cryptographic Hash Functions

A **cryptographic hash function** must have these properties:

#### 1. Deterministic
The same input always produces the same output.
```
H("hello") = 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
H("hello") = 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824  // Always the same
```

**Why it matters**: This allows us to verify data integrity. If you hash a file today and hash it again tomorrow, identical hashes mean the file hasn't changed.

#### 2. Fast to Compute
Computing the hash should be efficient, even for large inputs.

**Performance**: SHA-256 can hash data at hundreds of megabytes per second on modern CPUs.

**Why it matters**: You can hash entire files, disk images, or streaming data without unacceptable delays.

#### 3. Pre-image Resistance (One-way)
Given a hash `h`, it should be computationally infeasible to find any input `m` such that `H(m) = h`.

**Example**: If I give you hash `2cf24dba...9824`, you can't reverse it to find "hello" (without trying every possible input).

**Why it matters**: Password storage. We store `H(password)` in databases. Even if the database leaks, attackers can't recover original passwords.

#### 4. Second Pre-image Resistance
Given an input `m1`, it should be computationally infeasible to find a different input `m2` such that `H(m1) = H(m2)`.

**Example**: If you know `H("hello") = 2cf24...`, you can't find another string that produces the same hash.

**Why it matters**: Prevents attackers from creating malicious files with the same hash as legitimate files (malware masquerading as trusted software).

#### 5. Collision Resistance
It should be computationally infeasible to find any two different inputs `m1` and `m2` such that `H(m1) = H(m2)`.

**Mathematical insight**: By the pigeonhole principle, collisions MUST exist (infinite possible inputs, finite possible outputs). But finding them should be practically impossible.

**Attack complexity**: For SHA-256, finding a collision requires approximately 2^128 operations (astronomically large).

**Why it matters**: Digital signatures rely on this. If you could find two documents with the same hash, you could get one signed and swap it for the other.

#### 6. Avalanche Effect
A small change in input should produce a completely different hash (approximately 50% of bits flip).

**Example**:
```
H("hello")  = 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
H("Hello")  = 185f8db32271fe25f561a6fc938b2e264306ec304eda518007d1764826381969
              ↑ Only changed one letter (h→H), but hash is completely different!
```

**Why it matters**: Makes it impossible to predict hash changes based on input changes. Essential for security.

### What Makes SHA-256 Special?

**SHA-256** (Secure Hash Algorithm 256-bit) is part of the SHA-2 family, designed by the NSA and published in 2001.

**Key facts**:
- **Output size**: 256 bits (32 bytes, 64 hex characters)
- **Block size**: 512 bits (processes input in 64-byte chunks)
- **Rounds**: 64 compression rounds per block
- **Security**: No known practical attacks (as of 2024)
- **Usage**: Bitcoin, TLS/SSL certificates, file integrity, Git commits

**Comparison with other hash functions**:

| Algorithm | Output Size | Security | Status |
|-----------|-------------|----------|--------|
| MD5 | 128 bits | Broken (collisions found) | Deprecated |
| SHA-1 | 160 bits | Broken (collisions found) | Deprecated |
| **SHA-256** | 256 bits | Secure | Recommended |
| SHA-512 | 512 bits | Secure | Recommended |
| SHA-3 | Variable | Secure | Recommended |

---

## 3. How SHA-256 Works (Simplified)

### High-Level Algorithm

```
1. Pad the message to a multiple of 512 bits
2. Parse the padded message into 512-bit blocks
3. Initialize 8 hash values (H0 through H7) with constants
4. For each 512-bit block:
   a. Prepare 64 words (W0 through W63) from the block
   b. Initialize 8 working variables (a,b,c,d,e,f,g,h) with hash values
   c. Perform 64 rounds of compression using bitwise operations
   d. Add compressed values back to hash values
5. Concatenate final hash values to produce 256-bit digest
```

### Step 1: Padding

**Goal**: Make the message length a multiple of 512 bits.

**Process**:
1. Append a single '1' bit
2. Append '0' bits until length ≡ 448 (mod 512)
3. Append original message length as 64-bit big-endian integer

**Example**: Hash "abc" (24 bits)
```
Original:  01100001 01100010 01100011  (3 bytes)
Add '1':   01100001 01100010 01100011 1
Add '0's:  01100001 01100010 01100011 10000000 00000000 ... (until 448 bits total)
Add len:   ... 00000000 00000000 00000000 00011000  (24 in binary)
Result:    512 bits total
```

### Step 2: Initialize Hash Values

SHA-256 starts with 8 constants (first 32 bits of fractional parts of square roots of first 8 primes):
```
H0 = 0x6a09e667  (√2)
H1 = 0xbb67ae85  (√3)
H2 = 0x3c6ef372  (√5)
H3 = 0xa54ff53a  (√7)
H4 = 0x510e527f  (√11)
H5 = 0x9b05688c  (√13)
H6 = 0x1f83d9ab  (√17)
H7 = 0x5be0cd19  (√19)
```

**Why these numbers?**: Using constants derived from irrational numbers ensures no hidden structure (nothing-up-my-sleeve numbers).

### Step 3: Compression Function

For each 512-bit block, SHA-256 performs 64 rounds of operations involving:

**Bitwise operations**:
- **AND**: `a AND b` (both bits must be 1)
- **OR**: `a OR b` (at least one bit is 1)
- **XOR**: `a XOR b` (exclusive or, bits must differ)
- **NOT**: `NOT a` (flip all bits)
- **Rotate**: `ROTR^n(x)` (circular right shift by n positions)
- **Shift**: `SHR^n(x)` (right shift by n, fill with zeros)

**Key functions**:
```
Ch(x,y,z)  = (x AND y) XOR (NOT x AND z)   // Choose
Maj(x,y,z) = (x AND y) XOR (x AND z) XOR (y AND z)  // Majority
Σ0(x) = ROTR^2(x) XOR ROTR^13(x) XOR ROTR^22(x)
Σ1(x) = ROTR^6(x) XOR ROTR^11(x) XOR ROTR^25(x)
σ0(x) = ROTR^7(x) XOR ROTR^18(x) XOR SHR^3(x)
σ1(x) = ROTR^17(x) XOR ROTR^19(x) XOR SHR^10(x)
```

**Each round**:
```
T1 = h + Σ1(e) + Ch(e,f,g) + Kt + Wt
T2 = Σ0(a) + Maj(a,b,c)
h = g
g = f
f = e
e = d + T1
d = c
c = b
b = a
a = T1 + T2
```

Where `Kt` are round constants (derived from cube roots of first 64 primes).

### Step 4: Produce Final Hash

After processing all blocks:
```
H0 = H0 + a
H1 = H1 + b
H2 = H2 + c
H3 = H3 + d
H4 = H4 + e
H5 = H5 + f
H6 = H6 + g
H7 = H7 + h

Final hash = H0 || H1 || H2 || H3 || H4 || H5 || H6 || H7  (256 bits)
```

**Note**: You don't need to implement this from scratch! Go's `crypto/sha256` package provides a battle-tested implementation.

---

## 4. Hex Encoding: Making Hashes Readable

### The Problem with Binary

SHA-256 produces 256 bits (32 bytes) of binary data:
```
[0x2c, 0xf2, 0x4d, 0xba, 0x5f, 0xb0, 0xa3, 0x0e, ...]  // 32 bytes
```

**Problems**:
- Not human-readable
- Hard to copy/paste (contains non-printable characters)
- Difficult to communicate verbally or in text

### Hexadecimal Encoding

**Hex** (base-16) encodes each byte as two characters from `[0-9a-f]`.

**Conversion**:
```
Binary byte:   10110101  (181 in decimal)
High nibble:   1011 = 11 = 'b'
Low nibble:    0101 = 5  = '5'
Hex:           "b5"
```

**Full SHA-256 hash**:
```
Binary:  [0x2c, 0xf2, 0x4d, ...]  (32 bytes)
Hex:     "2cf24dba..."            (64 characters)
```

**Why hex is standard**:
- Compact: 64 characters vs 256 bits
- Human-readable: Uses familiar characters
- Easy to compare: Visual inspection
- URL-safe: No special encoding needed

### Encoding in Go

```go
import (
	"crypto/sha256"
	"encoding/hex"
)

// Compute hash
hash := sha256.Sum256([]byte("hello"))  // Returns [32]byte

// Convert to hex string
hexHash := hex.EncodeToString(hash[:])  // "2cf24dba..."
```

**Alternative (using fmt)**:
```go
hexHash := fmt.Sprintf("%x", hash)  // Same result
```

---

## 5. Using SHA-256 in Go

### Basic Hashing (One-shot)

**Simple string hashing**:
```go
import "crypto/sha256"

func HashString(s string) [32]byte {
	return sha256.Sum256([]byte(s))
}

hash := HashString("hello")
// hash = [32]byte{0x2c, 0xf2, 0x4d, ...}
```

**Why `[32]byte`?**: SHA-256 always produces exactly 32 bytes.

### Incremental Hashing (Streaming)

For large files or streaming data, you don't want to load everything into memory.

**Pattern**:
```go
h := sha256.New()  // Create hash.Hash interface
io.Copy(h, reader) // Stream data to hash
hash := h.Sum(nil) // Finalize and get digest
```

**Full example (hashing a file)**:
```go
func HashFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
```

**How it works**:
1. `sha256.New()` creates a hash state (contains internal buffers, block processing state)
2. `io.Copy()` reads file in chunks, calling `h.Write()` for each chunk
3. Hash state is updated incrementally (processes 512-bit blocks)
4. `h.Sum(nil)` finalizes the hash (pads final block, returns digest)

**Memory usage**: Only buffers 64 bytes at a time (one SHA-256 block), regardless of file size!

### Important: Digest Reuse and Reset

**Question**: Can you reuse a hash.Hash for multiple inputs?

**Answer**: Yes, but you must `Reset()` it first!

```go
h := sha256.New()

// Hash first string
h.Write([]byte("hello"))
hash1 := h.Sum(nil)

// WRONG: Continuing to write appends to "hello"
h.Write([]byte("world"))
hash2 := h.Sum(nil)  // This is Hash("helloworld"), not Hash("world")!

// CORRECT: Reset before new hash
h.Reset()
h.Write([]byte("world"))
hash3 := h.Sum(nil)  // This is Hash("world") ✓
```

**Why reuse?**: Reduces allocations. If hashing many small items, reusing the hash.Hash is more efficient than creating new ones.

### Hash.Hash Interface

```go
type Hash interface {
	io.Writer                    // Write(p []byte) (n int, err error)
	Sum(b []byte) []byte         // Append hash to b and return
	Reset()                      // Reset to initial state
	Size() int                   // Size of hash in bytes (32 for SHA-256)
	BlockSize() int              // Block size in bytes (64 for SHA-256)
}
```

**Key methods**:

**`Write(p []byte)`**: Add more data to hash
- Returns (n, err) but **never returns an error** for hash functions
- Can call multiple times (data is accumulated)

**`Sum(b []byte)`**: Finalize and return hash
- Appends hash to `b` and returns the result
- Common pattern: `hash := h.Sum(nil)` (append to nil → returns new slice)
- **Does not modify internal state** (can call multiple times)

**`Reset()`**: Clear internal state
- Returns hash to initial state (like calling `sha256.New()` again)
- Allows reusing the same hash.Hash

**`Size()`**: Returns output size
- SHA-256: Always 32 bytes
- SHA-512: Always 64 bytes

**`BlockSize()`**: Returns internal block size
- SHA-256: 64 bytes (512 bits)
- Used for HMAC and other advanced operations

---

## 6. Practical Applications

### Application 1: File Integrity Verification

**Scenario**: You download a large file (e.g., Linux ISO). How do you know it wasn't corrupted during download or replaced with malware?

**Solution**: The website provides a SHA-256 checksum:
```
ubuntu-22.04.iso
SHA-256: a4acfda10b18da50e2ec50ccaf860d7f20b389df8765611142305c0e911d16fd
```

**Verification process**:
```bash
$ sha256sum ubuntu-22.04.iso
a4acfda10b18da50e2ec50ccaf860d7f20b389df8765611142305c0e911d16fd  ubuntu-22.04.iso
```

If hashes match → file is authentic and uncorrupted ✓

**In Go**:
```go
func VerifyFile(filename, expectedHash string) (bool, error) {
	hash, err := HashFile(filename)
	if err != nil {
		return false, err
	}

	actualHash := hex.EncodeToString(hash)
	return actualHash == expectedHash, nil
}
```

### Application 2: Deduplication

**Scenario**: You're building a backup system with millions of files. Many files are duplicates (same content, different names/locations).

**Solution**: Use SHA-256 as a content identifier:
```go
type FileStore struct {
	hashes map[string]string  // hash → file path
}

func (fs *FileStore) AddFile(path string) error {
	hash, err := HashFile(path)
	if err != nil {
		return err
	}

	hashStr := hex.EncodeToString(hash)
	if existingPath, exists := fs.hashes[hashStr]; exists {
		// Duplicate! Point to existing file instead of storing again
		fmt.Printf("Duplicate: %s is same as %s\n", path, existingPath)
		return nil
	}

	fs.hashes[hashStr] = path
	// Store file...
	return nil
}
```

**Real-world usage**: Dropbox, Git, and Docker all use content-addressable storage based on hashing.

### Application 3: Password Hashing (with caveats!)

**WARNING**: Never use plain SHA-256 for password hashing in production! Use bcrypt, scrypt, or Argon2 instead.

**Why SHA-256 alone is insufficient**:
1. **Too fast**: Attackers can try billions of passwords per second
2. **No salt**: Same password → same hash (rainbow table attacks)
3. **No stretching**: One hash operation per password attempt

**Proper approach (using bcrypt)**:
```go
import "golang.org/x/crypto/bcrypt"

// Hashing
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// Verification
err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
if err != nil {
	// Wrong password
}
```

**What bcrypt does right**:
- Slow by design (adjustable work factor)
- Automatic salt generation
- Resistant to GPU/ASIC attacks

**When to use SHA-256**: Non-password use cases (file integrity, content addressing, HMACs, etc.)

### Application 4: Content-Addressable Storage (Git)

**How Git uses SHA-1 (similar to SHA-256)**:

Every Git object (commit, tree, blob) is identified by its hash:
```
blob 14\0Hello, world!  →  SHA-1  →  af5626b4a114abcb82d63db7c8082c3c4756e51b
```

**Storage**:
```
.git/objects/af/5626b4a114abcb82d63db7c8082c3c4756e51b
              ^^  ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
           First 2  Rest of hash (38 characters)
```

**Benefits**:
- Content-based naming (same content → same hash → deduplicated)
- Integrity verification (corrupt file → hash doesn't match → detected)
- Distributed system (everyone can verify hashes independently)

**Note**: Git is transitioning from SHA-1 to SHA-256 due to SHA-1 collision vulnerabilities.

### Application 5: Blockchain (Bitcoin)

**Bitcoin mining**: Find a nonce such that:
```
SHA-256(SHA-256(block_header + nonce)) < target
```

**Example**:
```
Block header: "Bitcoin Block #12345..."
Nonce: Try 0, 1, 2, 3, ... until hash has required leading zeros

Target: 0000000000000000000abc... (18 leading zeros)
Hash:   0000000000000000000def... (Success! 18 leading zeros)
```

**Why double SHA-256?**: Extra security against length-extension attacks.

**Difficulty**: Finding a hash with N leading zeros requires ~2^(4N) attempts on average.

**Current Bitcoin difficulty**: ~73 trillion trillion hashes to find one valid block!

---

## 7. Common Patterns and Best Practices

### Pattern 1: Hash and Hex in One Function

```go
func HashToHex(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// Usage
fmt.Println(HashToHex([]byte("hello")))
// Output: 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
```

### Pattern 2: Streaming Large Files

```go
func HashLargeFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	h := sha256.New()

	// Uses buffered I/O automatically
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
```

**Memory usage**: Constant (~64KB buffer), regardless of file size (1MB or 10GB).

### Pattern 3: Hashing Multiple Inputs (Concatenation)

```go
func HashMultiple(parts ...string) string {
	h := sha256.New()
	for _, part := range parts {
		h.Write([]byte(part))
	}
	return hex.EncodeToString(h.Sum(nil))
}

// Hash("hello" + "world") is same as Hash("helloworld")
hash1 := HashMultiple("hello", "world")
hash2 := HashString("helloworld")
// hash1 == hash2 ✓
```

**Warning**: Concatenation can be ambiguous:
```go
HashMultiple("hello", "world")  // "helloworld"
HashMultiple("hell", "oworld")  // "helloworld" (same hash!)
```

**Solution**: Use delimiters or length prefixes in protocols.

### Pattern 4: Comparing Hashes Securely

```go
import "crypto/subtle"

func SecureHashCompare(hash1, hash2 []byte) bool {
	return subtle.ConstantTimeCompare(hash1, hash2) == 1
}
```

**Why `subtle.ConstantTimeCompare`?**
- Normal comparison (`==`) can leak timing information
- Attacker can measure how long comparison takes
- For cryptographic operations, use constant-time comparison

**When to use**:
- Password hash comparison
- HMAC verification
- Any security-critical comparison

**When NOT needed**:
- File integrity checks (no adversarial timing attacks)
- Content deduplication

### Pattern 5: Hash Directory Tree

```go
func HashDirectory(dir string) (string, error) {
	h := sha256.New()

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil  // Skip directories
		}

		// Include filename in hash
		rel, _ := filepath.Rel(dir, path)
		h.Write([]byte(rel))

		// Include file contents
		fileHash, err := HashFile(path)
		if err != nil {
			return err
		}
		h.Write(fileHash)

		return nil
	})

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
```

**Use case**: Verify entire directory hasn't changed (e.g., deployed code).

---

## 8. Collision Resistance: Why It Matters

### The Birthday Paradox

**Question**: How many people need to be in a room before there's a 50% chance two share a birthday?

**Answer**: Only 23 people!

**Intuition**: Not comparing one person to everyone else, but all pairs.

**Math**: With n people, there are n(n-1)/2 pairs. With 365 possible birthdays, 23 people create 253 pairs, enough for ~50% collision probability.

### Collision Probability for SHA-256

**Output space**: 2^256 possible hashes

**Birthday attack**: To find any collision (not targeting specific hash), need approximately:
```
√(2^256) = 2^128 hashes
```

**How hard is 2^128?**
- If entire world's computing power (~1 zettaFLOPS) could generate 1 trillion hashes/second
- It would take ~10 billion years to reach 2^128 hashes
- The universe is only ~13.8 billion years old!

**Conclusion**: SHA-256 collisions are theoretically possible but practically infeasible with current technology.

### Length Extension Attacks

**Vulnerability**: Given `H(message)`, attacker can compute `H(message + padding + extra)` without knowing `message`.

**Affected**: SHA-256, SHA-512 (Merkle-Damgård construction)

**Not affected**: SHA-3 (Sponge construction)

**Example attack**:
```
// Server
secret = "my_secret"
message = "send $100"
auth_tag = SHA256(secret + message)  // Send message + auth_tag to client

// Attacker (doesn't know secret)
// Can compute: SHA256(secret + message + padding + " to attacker")
// Without knowing secret!
```

**Mitigation**: Use HMAC instead of plain hashing:
```go
import "crypto/hmac"

mac := hmac.New(sha256.New, []byte(secret))
mac.Write([]byte(message))
tag := mac.Sum(nil)
```

HMAC is not vulnerable to length extension attacks.

---

## 9. Common Mistakes to Avoid

### Mistake 1: Treating Hash as Encryption

**❌ Wrong thinking**: "I'll hash my data to keep it secret"

**Why wrong**: Hashing is **one-way**. You can't decrypt it.

**Correct use cases**:
- Password storage (you verify, not decrypt)
- Data integrity (you compare hashes, not recover data)

**For encryption, use**: AES, ChaCha20, etc. (crypto/cipher package)

### Mistake 2: Using SHA-256 for Passwords Directly

**❌ Wrong**:
```go
passwordHash := sha256.Sum256([]byte(password))
// Store passwordHash in database
```

**Why wrong**:
- Too fast (GPU can try billions/sec)
- No salt (same password → same hash → rainbow tables)

**✅ Correct**:
```go
import "golang.org/x/crypto/bcrypt"
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
```

### Mistake 3: Forgetting to Reset Hash State

**❌ Wrong**:
```go
h := sha256.New()

h.Write([]byte("first"))
hash1 := h.Sum(nil)

h.Write([]byte("second"))
hash2 := h.Sum(nil)  // Hash of "firstsecond", not "second"!
```

**✅ Correct**:
```go
h := sha256.New()

h.Write([]byte("first"))
hash1 := h.Sum(nil)

h.Reset()  // Clear state
h.Write([]byte("second"))
hash2 := h.Sum(nil)  // Hash of "second" ✓
```

### Mistake 4: Loading Entire Large File into Memory

**❌ Wrong**:
```go
data, _ := os.ReadFile("large_file.iso")  // Could be GBs!
hash := sha256.Sum256(data)               // Out of memory!
```

**✅ Correct**:
```go
file, _ := os.Open("large_file.iso")
defer file.Close()

h := sha256.New()
io.Copy(h, file)  // Streams file in chunks
hash := h.Sum(nil)
```

### Mistake 5: Using Hashes as Random Numbers

**❌ Wrong**:
```go
hash := sha256.Sum256([]byte("seed"))
randomNumber := binary.BigEndian.Uint64(hash[:8])  // NOT cryptographically random!
```

**Why wrong**: Deterministic. Given same input, same "random" number.

**✅ Correct** (for crypto-random numbers):
```go
import "crypto/rand"

var randomBytes [32]byte
rand.Read(randomBytes[:])
```

### Mistake 6: Comparing Hashes with == in Security Code

**❌ Wrong** (timing attack vulnerable):
```go
if storedHash == computedHash {
	// Authenticated
}
```

**Why wrong**: Comparison stops at first differing byte. Attacker can measure timing to guess hash byte-by-byte.

**✅ Correct**:
```go
import "crypto/subtle"

if subtle.ConstantTimeCompare(storedHash, computedHash) == 1 {
	// Authenticated
}
```

---

## 10. Stretch Goals

### Goal 1: HMAC Implementation ⭐⭐

Implement HMAC (Hash-based Message Authentication Code) using SHA-256.

**What is HMAC?**: A way to verify both integrity and authenticity of a message using a secret key.

**Formula**:
```
HMAC(key, message) = H((key ⊕ opad) || H((key ⊕ ipad) || message))

Where:
- H = hash function (SHA-256)
- opad = 0x5c repeated (outer padding)
- ipad = 0x36 repeated (inner padding)
- || = concatenation
- ⊕ = XOR
```

**Starter code**:
```go
func HMAC_SHA256(key, message []byte) []byte {
	// TODO: Implement HMAC
	// Hint: Use two sha256.New() instances
}
```

### Goal 2: Merkle Tree ⭐⭐⭐

Build a Merkle tree (used in Git, Bitcoin, BitTorrent).

**Structure**:
```
         Root Hash
        /          \
    H(AB)          H(CD)
    /    \         /    \
  H(A)  H(B)    H(C)  H(D)
   |     |       |      |
   A     B       C      D
```

**Properties**:
- Leaf nodes: Hashes of data blocks
- Internal nodes: Hashes of children concatenated
- Root hash: Represents entire tree

**Use case**: Efficiently verify a single item is part of a large dataset (O(log n) proof size).

**Challenge**: Implement tree construction and proof verification.

### Goal 3: Hash-based Deduplication Tool ⭐⭐

Build a CLI tool that scans directories and finds duplicate files.

**Features**:
- Hash all files
- Group by hash
- Report duplicates with total wasted space
- Option to delete duplicates (keep one copy)

**Example output**:
```
Scanning directory: /home/user/Documents
Found 1,234 files (5.2 GB)

Duplicates:
  report.pdf (3 copies, 5 MB each, wasted: 10 MB)
    - /home/user/Documents/report.pdf
    - /home/user/Documents/backup/report.pdf
    - /home/user/Downloads/report.pdf

Total wasted space: 127 MB across 45 duplicate files
```

### Goal 4: Demonstrate Avalanche Effect ⭐

Create a visualization showing how changing one bit in input affects the hash.

**Challenge**:
1. Hash a string
2. Flip each bit one at a time
3. Hash each modified version
4. Count how many bits differ in each hash vs original
5. Verify ~50% of bits flip for each single-bit input change

**Example output**:
```
Original: "hello"
Hash: 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824

Bit flip #0 (0x68 → 0x69): 129/256 bits changed (50.39%)
Bit flip #1 (0x68 → 0x6a): 124/256 bits changed (48.44%)
Bit flip #2 (0x68 → 0x6c): 131/256 bits changed (51.17%)
...
Average: 127.4/256 bits changed (49.77%) ✓ Avalanche effect confirmed
```

### Goal 5: Build a Simple Blockchain ⭐⭐⭐

Implement a basic blockchain using SHA-256.

**Structure**:
```go
type Block struct {
	Index     int
	Timestamp time.Time
	Data      string
	PrevHash  string
	Hash      string
	Nonce     int  // For proof-of-work
}

type Blockchain struct {
	Blocks []*Block
}
```

**Features**:
1. Genesis block creation
2. Add blocks with proof-of-work (find hash with N leading zeros)
3. Validate chain (each block's PrevHash matches previous block's Hash)
4. Detect tampering (changing any block invalidates all subsequent blocks)

---

## How to Run

```bash
# Run the demonstration program
cd /home/user/go-edu
make run P=39-sha256-hasher

# Run tests
go test ./minis/39-sha256-hasher/exercise/...

# Run specific test
go test -run TestHashString ./minis/39-sha256-hasher/exercise/

# Benchmark hashing performance
go test -bench=. ./minis/39-sha256-hasher/exercise/
```

---

## Summary

**What you learned**:
- ✅ Cryptographic hash functions are one-way, deterministic, and collision-resistant
- ✅ SHA-256 produces 256-bit (32-byte) digests, encoded as 64 hex characters
- ✅ Hash functions have critical properties: preimage resistance, collision resistance, avalanche effect
- ✅ `crypto/sha256` provides both one-shot (`Sum256`) and streaming (`New()`) APIs
- ✅ Use `io.Copy()` for efficient file hashing without loading into memory
- ✅ Always `Reset()` hash.Hash before reusing for new inputs
- ✅ Hex encoding makes binary hashes human-readable
- ✅ Real-world uses: file integrity, deduplication, content-addressing, blockchain

**Why this matters**:
SHA-256 and cryptographic hashing are fundamental to modern security:
- **TLS/SSL**: Certificate fingerprints
- **Git**: Content-addressable storage
- **Bitcoin**: Proof-of-work mining
- **Password storage**: Foundation for bcrypt/scrypt (though not used directly)
- **File integrity**: Software distribution, backups

**Next steps**:
- Learn about HMAC for authenticated hashing
- Explore bcrypt/Argon2 for password hashing
- Study Merkle trees for efficient data verification
- Understand digital signatures (ECDSA, Ed25519)

Hash on! 🔒
