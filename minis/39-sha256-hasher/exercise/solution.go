//go:build solution
// +build solution

/*
Problem: Implementing SHA-256 cryptographic hash operations

We need to:
1. Hash strings and return hex-encoded digests
2. Hash files efficiently (streaming, not loading into memory)
3. Verify file integrity using checksums
4. Demonstrate incremental hashing (streaming multiple inputs)
5. Compare hashes securely (constant-time to prevent timing attacks)

Key Concepts:
- Cryptographic hash: One-way function mapping arbitrary data to fixed-size digest
- SHA-256: 256-bit (32-byte) secure hash algorithm
- Hex encoding: Binary → human-readable (64 hex characters for SHA-256)
- Incremental hashing: Process data in chunks (essential for large files)
- Constant-time comparison: Prevent timing attacks in security-critical code

Time/Space Complexity:
- HashString: O(n) time, O(1) space (where n = string length)
- HashFile: O(n) time, O(1) space (where n = file size, streams in chunks)
- VerifyFile: O(n) time, O(1) space (same as HashFile)
- HashIncremental: O(n) time, O(1) space (where n = total length of all parts)
- CompareHashes: O(1) time (hashes are fixed size), O(1) space

Why Go is well-suited:
- crypto/sha256: Battle-tested, optimized implementation
- io.Copy: Efficient streaming without manual buffer management
- hash.Hash interface: Clean abstraction for incremental hashing
- crypto/subtle: Security primitives (constant-time operations)

Real-world applications:
- File integrity: Verify downloads haven't been corrupted or tampered
- Deduplication: Identify identical files by content (not name)
- Content-addressable storage: Git, IPFS use hashes as identifiers
- Password verification: Store Hash(password), verify by comparing hashes
- Digital signatures: Hash documents before signing (efficiency)
*/

package exercise

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"io"
	"os"
	"strings"
)

// HashString computes the SHA-256 hash of a string and returns it as a hex-encoded string.
//
// Go Concepts Demonstrated:
// - crypto/sha256.Sum256: One-shot hash computation
// - Type conversion: string → []byte
// - Array slicing: [32]byte → []byte (for hex encoding)
// - encoding/hex: Binary to hexadecimal conversion
//
// Three-Input Iteration Table:
//
// Input 1: "hello"
//   Convert to bytes: []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f}
//   SHA-256: [32]byte{0x2c, 0xf2, 0x4d, 0xba, ...}
//   Hex encode: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
//   Result: "2cf24dba..." (64 characters)
//
// Input 2: "" (empty string)
//   Convert to bytes: []byte{} (empty)
//   SHA-256: [32]byte{0xe3, 0xb0, 0xc4, 0x42, ...}
//   Hex encode: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
//   Result: "e3b0c442..." (hash of empty input, deterministic)
//
// Input 3: "The quick brown fox jumps over the lazy dog"
//   Convert to bytes: []byte{0x54, 0x68, 0x65, ...} (43 bytes)
//   SHA-256: [32]byte{0xd7, 0xa8, 0xfb, 0xb3, ...}
//   Hex encode: "d7a8fbb307d7809469ca9abcb0082e4f8d5651e46d3cdb762d02d0bf37c9e592"
//   Result: "d7a8fbb3..." (always 64 hex chars regardless of input length)
func HashString(s string) string {
	// Sum256 computes SHA-256 in one operation
	// Returns [32]byte (fixed-size array)
	hash := sha256.Sum256([]byte(s))

	// Convert array to slice for hex encoding
	// hash[:] means "all elements from hash array as a slice"
	return hex.EncodeToString(hash[:])
}

// HashFile computes the SHA-256 hash of a file's contents.
//
// Go Concepts Demonstrated:
// - os.Open: File I/O
// - defer: Ensure file is closed even if error occurs
// - sha256.New(): Create hash.Hash interface (stateful, incremental hashing)
// - io.Copy: Efficient streaming (automatic buffering)
// - hash.Hash.Sum(nil): Finalize hash and get digest
//
// Why streaming is important:
// - Loading 10GB file into memory: 10GB RAM usage (may fail)
// - Streaming 10GB file: ~64KB RAM usage (internal buffer)
// - Same hash result, drastically different resource usage
//
// Three-Input Iteration Table:
//
// Input 1: File with "Hello, World!" (13 bytes)
//   Open file: Success
//   Create hash: h = sha256.New()
//   io.Copy: Read file in chunks, call h.Write() for each chunk
//     - Chunk 1: "Hello, World!" (entire file fits in one read)
//   Finalize: h.Sum(nil) → []byte{0x31, 0x5f, 0x5b, ...} (32 bytes)
//   Result: 32-byte hash
//
// Input 2: File with 1MB of 'A' characters
//   Open file: Success
//   Create hash: h = sha256.New()
//   io.Copy: Read file in chunks (~32KB at a time)
//     - Chunk 1: 32KB of 'A's → h.Write()
//     - Chunk 2: 32KB of 'A's → h.Write()
//     - ... (repeat ~32 times)
//     - Final chunk: Remaining bytes → h.Write()
//   Finalize: h.Sum(nil) → []byte{...} (32 bytes)
//   Result: 32-byte hash (memory usage stayed constant)
//
// Input 3: Non-existent file
//   Open file: os.PathError (file not found)
//   Return: nil, error
//   Result: Error propagated to caller
func HashFile(filename string) ([]byte, error) {
	// Open file for reading
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	// Ensure file is closed when function returns
	// Even if io.Copy fails, file will be closed
	defer file.Close()

	// Create a new SHA-256 hash state
	// This returns a hash.Hash interface with internal state for incremental hashing
	h := sha256.New()

	// Copy file contents to hash
	// io.Copy reads from file and writes to hash in chunks
	// - Default buffer size: 32KB
	// - No need to manage buffers manually
	// - Works with files of any size
	if _, err := io.Copy(h, file); err != nil {
		return nil, err
	}

	// Finalize the hash and return the digest
	// Sum(nil) means "append hash to nil slice" → returns new slice with hash
	// Could also do Sum([]byte{}) or Sum(existingSlice) to append
	return h.Sum(nil), nil
}

// VerifyFile checks if a file's SHA-256 hash matches an expected hash.
//
// Go Concepts Demonstrated:
// - Error propagation: Return error from HashFile
// - String comparison: Case-insensitive (using strings.EqualFold)
// - Hex encoding consistency: Always lowercase from hex.EncodeToString
//
// Use cases:
// - Software downloads: Verify ISO/binary hasn't been tampered with
// - Backup verification: Ensure backup matches original
// - File synchronization: Check if file has changed
//
// Three-Input Iteration Table:
//
// Input 1: File "test.txt" with "hello", expected "2cf24dba..."
//   HashFile: Returns []byte{0x2c, 0xf2, 0x4d, ...}
//   Encode to hex: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
//   Compare: "2cf24dba..." == "2cf24dba..." (match)
//   Result: true, nil
//
// Input 2: File "test.txt" with "hello", expected "AAAA..." (wrong hash)
//   HashFile: Returns []byte{0x2c, 0xf2, 0x4d, ...}
//   Encode to hex: "2cf24dba..."
//   Compare: "2cf24dba..." == "AAAA..." (no match)
//   Result: false, nil
//
// Input 3: Non-existent file, any expected hash
//   HashFile: Returns nil, error (file not found)
//   Result: false, error (error propagated)
func VerifyFile(filename, expectedHashHex string) (bool, error) {
	// Compute actual hash of file
	actualHash, err := HashFile(filename)
	if err != nil {
		// File couldn't be read, return error
		return false, err
	}

	// Convert hash bytes to hex string
	actualHashHex := hex.EncodeToString(actualHash)

	// Compare hashes (case-insensitive for robustness)
	// Some systems might provide uppercase hex, others lowercase
	// EqualFold treats "abc" == "ABC"
	match := strings.EqualFold(actualHashHex, expectedHashHex)

	return match, nil
}

// HashIncremental demonstrates incremental hashing by combining multiple strings.
//
// Go Concepts Demonstrated:
// - Variadic functions: ...string allows variable number of arguments
// - hash.Hash.Write(): Incremental state updates
// - Equivalence: Hash(a+b+c) == Hash(a) then Hash(b) then Hash(c) finalized
//
// Why incremental hashing matters:
// - Large data: Don't need to concatenate first (saves memory)
// - Streaming: Process data as it arrives (network, file, etc.)
// - Efficiency: Single allocation for hash state vs multiple concatenations
//
// Example use case:
//   Hashing HTTP request: Hash(method + url + headers + body)
//   Without incremental: Concatenate all → hash (2 allocations)
//   With incremental: Hash each part (1 allocation for hash state)
//
// Three-Input Iteration Table:
//
// Input 1: HashIncremental("hello", "world")
//   Create hash: h = sha256.New()
//   Write "hello": h.Write([]byte("hello"))
//   Write "world": h.Write([]byte("world"))
//   Finalize: h.Sum(nil) → same as Hash("helloworld")
//   Encode to hex: "936a185c..."
//   Result: "936a185c..." (64 characters)
//
// Input 2: HashIncremental("a", "b", "c", "d", "e")
//   Create hash: h = sha256.New()
//   Write "a": h.Write([]byte("a"))
//   Write "b": h.Write([]byte("b"))
//   Write "c": h.Write([]byte("c"))
//   Write "d": h.Write([]byte("d"))
//   Write "e": h.Write([]byte("e"))
//   Finalize: h.Sum(nil) → same as Hash("abcde")
//   Result: Same as HashString("abcde")
//
// Input 3: HashIncremental() (no arguments)
//   Create hash: h = sha256.New()
//   Loop: No iterations (parts is empty)
//   Finalize: h.Sum(nil) → hash of empty input
//   Result: "e3b0c44298..." (hash of empty string)
func HashIncremental(parts ...string) string {
	// Create new hash state
	h := sha256.New()

	// Write each part to the hash
	// This is equivalent to hashing the concatenation of all parts
	for _, part := range parts {
		// hash.Hash.Write() never returns an error
		// (it's part of the hash.Hash interface contract)
		h.Write([]byte(part))
	}

	// Finalize and get the hash
	hash := h.Sum(nil)

	// Convert to hex string
	return hex.EncodeToString(hash)
}

// CompareHashes compares two hashes in constant time to prevent timing attacks.
//
// Go Concepts Demonstrated:
// - crypto/subtle: Constant-time cryptographic operations
// - hex.DecodeString: Hex → binary conversion
// - Error handling: Graceful handling of invalid hex strings
//
// Security consideration: Timing attacks
//
// Vulnerable comparison (DO NOT USE for security):
//   if hash1 == hash2 {  // Stops at first different byte!
//
// Attack: Attacker measures time to compare:
//   - hash1 = "aaaa...", hash2 = "aaaa..." → Slow (compares all)
//   - hash1 = "aaaa...", hash2 = "baaa..." → Fast (stops at first byte)
//   - Attacker can guess hash byte-by-byte by measuring timing!
//
// Constant-time comparison (USE for security):
//   subtle.ConstantTimeCompare(hash1, hash2)  // Always compares all bytes
//   - Takes same time regardless of where difference occurs
//   - Prevents timing-based guessing
//
// When to use constant-time:
// ✓ Password hash verification
// ✓ HMAC verification
// ✓ Cryptographic signatures
// ✗ File integrity checks (no adversarial timing attack)
// ✗ Deduplication (performance matters more than timing leakage)
//
// Three-Input Iteration Table:
//
// Input 1: hash1 = "2cf24dba...", hash2 = "2cf24dba..." (identical)
//   Decode hash1: []byte{0x2c, 0xf2, 0x4d, ...}
//   Decode hash2: []byte{0x2c, 0xf2, 0x4d, ...}
//   ConstantTimeCompare: 1 (equal)
//   Result: true
//
// Input 2: hash1 = "2cf24dba...", hash2 = "aaaa..." (different)
//   Decode hash1: []byte{0x2c, 0xf2, ...}
//   Decode hash2: []byte{0xaa, 0xaa, ...}
//   ConstantTimeCompare: 0 (not equal)
//   Result: false
//
// Input 3: hash1 = "ZZZZ" (invalid hex), hash2 = "2cf24dba..."
//   Decode hash1: error (invalid hex character 'Z')
//   Result: false (treat decode error as non-match)
func CompareHashes(hash1, hash2 string) bool {
	// Decode hex strings to bytes
	bytes1, err1 := hex.DecodeString(hash1)
	bytes2, err2 := hex.DecodeString(hash2)

	// If either decode fails, hashes don't match
	if err1 != nil || err2 != nil {
		return false
	}

	// Constant-time comparison
	// Returns 1 if equal, 0 if not equal
	// Always compares all bytes (doesn't short-circuit)
	return subtle.ConstantTimeCompare(bytes1, bytes2) == 1
}

/*
Alternatives & Trade-offs:

1. HashString using sha256.New() instead of Sum256():
   func HashString(s string) string {
       h := sha256.New()
       h.Write([]byte(s))
       return hex.EncodeToString(h.Sum(nil))
   }
   Pros: Consistent with streaming pattern
   Cons: More allocations for small inputs; Sum256 is optimized for one-shot
   Trade-off: Sum256 is cleaner and faster for complete inputs

2. HashFile loading entire file into memory:
   func HashFile(filename string) ([]byte, error) {
       data, err := os.ReadFile(filename)
       if err != nil { return nil, err }
       hash := sha256.Sum256(data)
       return hash[:], nil
   }
   Pros: Simpler code (3 lines)
   Cons: Fails for large files (out of memory); not scalable
   Trade-off: Never do this in production! Streaming is essential

3. VerifyFile using bytes.Equal instead of strings.EqualFold:
   Pros: Exact comparison (stricter)
   Cons: Case-sensitive (might fail if expected hash is uppercase)
   Trade-off: EqualFold is more robust for real-world usage

4. CompareHashes using == operator:
   Pros: Simpler, faster
   Cons: Vulnerable to timing attacks (security risk)
   Trade-off: Only use subtle.ConstantTimeCompare for security-critical code
              For non-security use cases (file integrity), == is fine

5. HashIncremental pre-concatenating strings:
   func HashIncremental(parts ...string) string {
       concatenated := strings.Join(parts, "")
       return HashString(concatenated)
   }
   Pros: Simpler (one line)
   Cons: Extra allocation for concatenated string; defeats purpose of "incremental"
   Trade-off: Use incremental for large data or streaming scenarios

Go vs X:

Go vs Python (hashlib):
  import hashlib
  hash = hashlib.sha256(b"hello").hexdigest()
  Pros: Very concise
  Cons: Less explicit about encoding; no compiler help
  Go: More verbose but type-safe; clear distinction between bytes and strings

Go vs Rust (sha2 crate):
  use sha2::{Sha256, Digest};
  let hash = Sha256::digest(b"hello");
  Pros: Zero-cost abstractions; very fast
  Cons: More complex trait system; steeper learning curve
  Go: Simpler interface; hash.Hash is straightforward

Go vs Java (MessageDigest):
  MessageDigest md = MessageDigest.getInstance("SHA-256");
  byte[] hash = md.digest("hello".getBytes());
  Pros: Similar structure to Go's hash.Hash
  Cons: More verbose; exceptions instead of errors
  Go: Cleaner error handling; fewer lines of code

Go vs JavaScript (crypto):
  const crypto = require('crypto');
  const hash = crypto.createHash('sha256').update('hello').digest('hex');
  Pros: Very concise; fluent API
  Cons: Dynamic typing; easy to mess up encoding
  Go: Type safety prevents encoding bugs; compiler catches errors

Performance Benchmarks (approximate, on modern CPU):

HashString("hello"):
- Go:    ~500 ns/op
- Python: ~2000 ns/op (4x slower)
- Rust:   ~400 ns/op (faster, but similar)

HashFile(1GB file):
- Go streaming:  ~1.5s (constant ~64KB memory)
- Go all-in-mem: ~1.2s (1GB memory, may OOM)
- Python:        ~2.5s (similar memory if using streaming)

Key insight: Go's performance is excellent, and streaming API
             makes it practical for any file size.
*/
