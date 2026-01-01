//go:build !solution && !reference


package sha256hasher

// TODO: Import required packages
// You'll need:
// - "crypto/sha256" for SHA-256 hashing
// - "crypto/subtle" for constant-time comparison (security)
// - "encoding/hex" for hex encoding/decoding
// - "io" for streaming file I/O
// - "os" for file operations
//
// import (
//     "crypto/sha256"
//     "crypto/subtle"
//     "encoding/hex"
//     "io"
//     "os"
// )

// ============================================================================
// SHA-256 HASHING: Understanding Cryptographic Hash Functions
// ============================================================================
//
// What is SHA-256?
// - Secure Hash Algorithm 256-bit
// - Cryptographic hash function (one-way: can't reverse)
// - Fixed output: Always 32 bytes (256 bits) regardless of input size
// - Deterministic: Same input always produces same output
// - Avalanche effect: Tiny input change → completely different hash
//
// Key Properties:
// 1. Pre-image resistance: Given hash, can't find original input
// 2. Second pre-image resistance: Can't find different input with same hash
// 3. Collision resistance: Extremely hard to find two inputs with same hash
//
// Common Uses:
// - File integrity verification (checksums)
// - Password hashing (with salt)
// - Blockchain (Bitcoin uses SHA-256 for mining)
// - Digital signatures (sign the hash, not the full document)
// - Content-addressable storage (Git, IPFS)
//
// Memory Considerations:
// - Hash state: ~200 bytes (internal SHA-256 state)
// - Output: Always 32 bytes
// - Hex encoding: 64 characters (2 hex digits per byte)
// - Streaming: Process data in chunks (don't load entire file into memory)
//
// ============================================================================

// HashString computes the SHA-256 hash of a string and returns it as a hex-encoded string.
//
// TODO: Implement this function
//
// This is the simplest hash operation: single input → single hash.
//
// Algorithm:
// 1. Convert string to []byte (strings are immutable, must convert to bytes)
// 2. Compute SHA-256 hash using sha256.Sum256()
// 3. Convert hash to hex string using hex.EncodeToString()
//
// Memory flow:
// - Input string: Variable size (stored on heap if large)
// - []byte conversion: Creates a copy (strings are immutable)
// - sha256.Sum256(): Returns [32]byte (array on stack)
// - hex.EncodeToString(): Allocates 64-byte string (on heap)
//
// Example:
//   HashString("hello") → "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
//
// Why hex encoding?
// - Binary data isn't printable or JSON-safe
// - Hex is human-readable: "2cf24d..." instead of raw bytes
// - Each byte becomes 2 hex characters: 0xFF → "ff"
//
// func HashString(s string) string {
//     // Step 1: Convert string to bytes
//     // Why []byte(s)? Strings are immutable, hashing needs byte slice
//     data := []byte(s)
//
//     // Step 2: Compute SHA-256 hash
//     // Sum256 returns [32]byte (fixed-size array)
//     // This is a one-shot operation (entire input at once)
//     hash := sha256.Sum256(data)
//
//     // Step 3: Convert hash to hex string
//     // hash[:] converts [32]byte array to []byte slice
//     // Why? EncodeToString expects a slice, not array
//     return hex.EncodeToString(hash[:])
// }

// HashFile computes the SHA-256 hash of a file's contents.
//
// TODO: Implement this function
//
// This demonstrates STREAMING hash computation (essential for large files).
//
// Algorithm:
// 1. Open file for reading
// 2. Create SHA-256 hash state (stateful, incremental)
// 3. Stream file contents to hash (io.Copy handles buffering)
// 4. Finalize hash and return bytes
//
// Why streaming?
// - Loading 10GB file into memory: 10GB RAM usage (may fail!)
// - Streaming 10GB file: ~64KB RAM usage (internal buffer)
// - Same hash result, drastically different memory footprint
//
// Memory flow:
// - os.Open(): File descriptor (~200 bytes)
// - sha256.New(): Hash state (~200 bytes)
// - io.Copy(): Internal buffer (~32KB, reused for all chunks)
// - hash.Sum(nil): Final hash (32 bytes allocated)
//
// Example:
//   HashFile("document.pdf") → []byte{0x2c, 0xf2, 0x4d, ...} (32 bytes)
//
// Error cases:
// - File doesn't exist: os.Open() returns os.PathError
// - Permission denied: os.Open() returns os.PathError
// - I/O error during read: io.Copy() returns error
//
// func HashFile(filename string) ([]byte, error) {
//     // Step 1: Open file for reading
//     // os.Open() returns (*os.File, error)
//     // File implements io.Reader interface
//     file, err := os.Open(filename)
//     if err != nil {
//         return nil, err  // Return error if file can't be opened
//     }
//     // Ensure file is closed when function exits
//     // defer runs even if function returns early or panics
//     defer file.Close()
//
//     // Step 2: Create SHA-256 hash state
//     // sha256.New() returns hash.Hash interface
//     // Unlike Sum256, this is STATEFUL (supports incremental updates)
//     h := sha256.New()
//
//     // Step 3: Stream file contents to hash
//     // io.Copy(dst, src) reads from src, writes to dst
//     // Reads file in chunks (~32KB), calls h.Write() for each chunk
//     // Returns (bytesWritten, error)
//     if _, err := io.Copy(h, file); err != nil {
//         return nil, err  // Return error if read fails
//     }
//
//     // Step 4: Finalize hash and get digest
//     // h.Sum(nil) appends hash to nil slice → returns new slice with hash
//     // Could also do h.Sum(existingSlice) to append to existing slice
//     return h.Sum(nil), nil
// }

// VerifyFile checks if a file's SHA-256 hash matches an expected hash.
//
// TODO: Implement this function
//
// This is how you verify file integrity (downloads, backups, etc.).
//
// Algorithm:
// 1. Hash the file (using HashFile)
// 2. Encode actual hash to hex
// 3. Compare with expected hash (case-insensitive)
//
// Use cases:
// - Software downloads: Verify installer hasn't been tampered with
// - Backup verification: Ensure backup matches original
// - File synchronization: Check if file has changed
//
// Security note:
// - This uses simple string comparison (not constant-time)
// - OK for file integrity (no timing attack vector)
// - For crypto verification, use CompareHashes() instead
//
// func VerifyFile(filename, expectedHashHex string) (bool, error) {
//     // Step 1: Compute actual hash of file
//     actualHash, err := HashFile(filename)
//     if err != nil {
//         return false, err  // File couldn't be read
//     }
//
//     // Step 2: Convert hash bytes to hex string
//     actualHashHex := hex.EncodeToString(actualHash)
//
//     // Step 3: Compare hashes (case-insensitive for robustness)
//     // Some sources provide uppercase hex, others lowercase
//     // strings.EqualFold treats "ABC" == "abc"
//     // Alternative: strings.ToLower(actualHashHex) == strings.ToLower(expectedHashHex)
//     match := strings.EqualFold(actualHashHex, expectedHashHex)
//
//     return match, nil
// }

// HashIncremental demonstrates incremental hashing by combining multiple strings.
//
// TODO: Implement this function
//
// This shows how to hash data in CHUNKS (without concatenating first).
//
// Algorithm:
// 1. Create SHA-256 hash state
// 2. Write each string part to hash (incremental)
// 3. Finalize and convert to hex
//
// Why incremental hashing?
// - Memory efficiency: Don't need to concatenate large strings first
// - Streaming: Process data as it arrives (network, file, etc.)
// - Flexibility: Hash multiple sources without buffering
//
// Mathematical property:
//   Hash(a + b + c) == Hash(a) then Hash(b) then Hash(c) finalized
//   The hash state accumulates all inputs before finalizing
//
// Example:
//   HashIncremental("hello", "world") == HashString("helloworld")
//
// Memory comparison:
// - Concatenating first: "hello" + "world" = "helloworld" (allocates 10 bytes)
// - Incremental: Hash state only (~200 bytes, fixed regardless of input size)
//
// func HashIncremental(parts ...string) string {
//     // Step 1: Create hash state
//     h := sha256.New()
//
//     // Step 2: Write each part to hash
//     // hash.Hash.Write() updates internal state
//     // Write() never returns an error (it's part of io.Writer interface contract)
//     for _, part := range parts {
//         h.Write([]byte(part))
//     }
//
//     // Step 3: Finalize hash
//     hash := h.Sum(nil)  // Get final digest
//
//     // Step 4: Convert to hex
//     return hex.EncodeToString(hash)
// }

// CompareHashes compares two hashes in constant time to prevent timing attacks.
//
// TODO: Implement this function
//
// This demonstrates SECURE hash comparison (important for cryptography).
//
// Algorithm:
// 1. Decode both hex strings to bytes
// 2. Use constant-time comparison (crypto/subtle)
// 3. Return true if equal, false otherwise
//
// Security: Timing Attacks Explained
//
// Vulnerable comparison (DO NOT USE for crypto):
//   if hash1 == hash2 {  // STOPS at first different byte!
//
// Attack scenario:
//   - hash1 = "aaaa...", hash2 = "aaaa..." → Slow (compares all bytes)
//   - hash1 = "aaaa...", hash2 = "baaa..." → Fast (stops at first byte)
//   - Attacker measures timing to guess hash byte-by-byte!
//
// Constant-time comparison (USE for crypto):
//   subtle.ConstantTimeCompare(hash1, hash2)  // Always compares ALL bytes
//   - Takes same time regardless of where difference occurs
//   - Prevents timing-based guessing attacks
//
// When to use constant-time:
// ✓ Password hash verification
// ✓ HMAC verification
// ✓ Cryptographic signature verification
// ✗ File integrity checks (no adversarial timing attack)
// ✗ Deduplication (performance matters more)
//
// func CompareHashes(hash1, hash2 string) bool {
//     // Step 1: Decode hex strings to bytes
//     bytes1, err1 := hex.DecodeString(hash1)
//     bytes2, err2 := hex.DecodeString(hash2)
//
//     // Step 2: Handle decode errors
//     // If either decode fails, hashes don't match
//     if err1 != nil || err2 != nil {
//         return false
//     }
//
//     // Step 3: Constant-time comparison
//     // subtle.ConstantTimeCompare returns 1 if equal, 0 if not
//     // Always compares all bytes (prevents timing attacks)
//     return subtle.ConstantTimeCompare(bytes1, bytes2) == 1
// }

// After implementing all functions:
// - Run: go test ./minis/39-sha256-hasher/exercise/...
// - Test solution: go test -tags solution ./minis/39-sha256-hasher/exercise/...
// - Compare with solution.go to see detailed implementations
//
// Key Concepts Learned:
// - Cryptographic hashing: One-way, deterministic, collision-resistant
// - SHA-256: 256-bit secure hash algorithm
// - Streaming: Process large files without loading into memory
// - Hex encoding: Binary → human-readable representation
// - Constant-time comparison: Prevent timing attacks
// - Incremental hashing: Update hash state progressively
//
// Real-World Applications:
// - Bitcoin: SHA-256 for proof-of-work mining
// - Git: SHA-1 (transitioning to SHA-256) for commit IDs
// - IPFS: Content-addressable storage using hashes
// - TLS: Certificate fingerprints
// - Package managers: Verify download integrity
