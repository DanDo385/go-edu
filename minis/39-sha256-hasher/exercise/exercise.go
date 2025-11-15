//go:build !solution
// +build !solution

package exercise

// HashString computes the SHA-256 hash of a string and returns it as a hex-encoded string.
//
// Example:
//   HashString("hello") → "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
//
// TODO: Implement this function
// Hints:
//   - Use crypto/sha256.Sum256() for the hash
//   - Use encoding/hex.EncodeToString() to convert to hex
//   - Remember to convert string to []byte
func HashString(s string) string {
	// TODO: implement
	return ""
}

// HashFile computes the SHA-256 hash of a file's contents.
//
// This function should:
//   - Open the file
//   - Stream the contents through the hash (don't load entire file into memory)
//   - Return the hash as a byte slice
//
// Returns an error if the file cannot be read.
//
// TODO: Implement this function
// Hints:
//   - Use os.Open() to open the file
//   - Use sha256.New() to create a hash.Hash
//   - Use io.Copy() to stream file contents to the hash
//   - Use hash.Sum(nil) to finalize and get the digest
func HashFile(filename string) ([]byte, error) {
	// TODO: implement
	return nil, nil
}

// VerifyFile checks if a file's SHA-256 hash matches an expected hash.
//
// Parameters:
//   - filename: Path to the file to verify
//   - expectedHashHex: Expected hash as a hex-encoded string
//
// Returns true if hashes match, false otherwise.
// Returns an error if file cannot be read.
//
// TODO: Implement this function
// Hints:
//   - Use HashFile() to get the actual hash
//   - Use hex.EncodeToString() to convert to hex
//   - Compare strings (case-insensitive for robustness)
func VerifyFile(filename, expectedHashHex string) (bool, error) {
	// TODO: implement
	return false, nil
}

// HashIncremental demonstrates incremental hashing by combining multiple strings.
//
// This function should hash all input strings as if they were concatenated,
// but use incremental hashing (multiple Write() calls) instead of concatenating first.
//
// Example:
//   HashIncremental("hello", "world") should produce same hash as HashString("helloworld")
//
// TODO: Implement this function
// Hints:
//   - Use sha256.New() to create a hash.Hash
//   - Call Write() for each input string
//   - Use Sum(nil) to get final hash
//   - Convert to hex string
func HashIncremental(parts ...string) string {
	// TODO: implement
	return ""
}

// CompareHashes compares two hashes in constant time to prevent timing attacks.
//
// This is important for security-critical applications where an attacker
// might try to learn information by measuring comparison time.
//
// Parameters:
//   - hash1, hash2: Hashes as hex-encoded strings
//
// Returns true if hashes are equal.
//
// TODO: Implement this function
// Hints:
//   - Use hex.DecodeString() to convert hex to bytes
//   - Use crypto/subtle.ConstantTimeCompare() for secure comparison
//   - Handle hex decode errors gracefully
func CompareHashes(hash1, hash2 string) bool {
	// TODO: implement
	return false
}
