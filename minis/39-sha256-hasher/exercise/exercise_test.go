package exercise

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strings"
	"testing"
)

func TestHashString_Basic(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "hello",
			expected: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		{
			input:    "world",
			expected: "486ea46224d1bb4fb680f34f7c9ad96a8f24ec88be73ea8e5a6c65260e9cb8a7",
		},
		{
			input:    "",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			input:    "The quick brown fox jumps over the lazy dog",
			expected: "d7a8fbb307d7809469ca9abcb0082e4f8d5651e46d3cdb762d02d0bf37c9e592",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := HashString(tt.input)
			if result != tt.expected {
				t.Errorf("HashString(%q) = %s, want %s", tt.input, result, tt.expected)
			}
			// Verify it's 64 characters (256 bits = 32 bytes = 64 hex chars)
			if len(result) != 64 {
				t.Errorf("HashString(%q) returned %d characters, want 64", tt.input, len(result))
			}
		})
	}
}

func TestHashString_CaseSensitive(t *testing.T) {
	hash1 := HashString("hello")
	hash2 := HashString("Hello")

	if hash1 == hash2 {
		t.Error("HashString should be case-sensitive: hash('hello') should differ from hash('Hello')")
	}
}

func TestHashString_Deterministic(t *testing.T) {
	input := "deterministic test"
	hash1 := HashString(input)
	hash2 := HashString(input)
	hash3 := HashString(input)

	if hash1 != hash2 || hash2 != hash3 {
		t.Error("HashString should be deterministic: same input should always produce same hash")
	}
}

func TestHashFile_Basic(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "hashtest-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write test data
	content := "Hello, World!"
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Hash the file
	hash, err := HashFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("HashFile failed: %v", err)
	}

	// Verify hash
	expected := sha256.Sum256([]byte(content))
	if !bytesEqual(hash, expected[:]) {
		t.Errorf("HashFile returned %x, want %x", hash, expected)
	}

	// Verify it's 32 bytes (256 bits)
	if len(hash) != 32 {
		t.Errorf("HashFile returned %d bytes, want 32", len(hash))
	}
}

func TestHashFile_EmptyFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "hashtest-empty-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	hash, err := HashFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("HashFile failed: %v", err)
	}

	// Hash of empty file
	expected := sha256.Sum256([]byte{})
	if !bytesEqual(hash, expected[:]) {
		t.Errorf("HashFile of empty file returned %x, want %x", hash, expected)
	}
}

func TestHashFile_LargeFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "hashtest-large-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write 1MB of 'A's
	largeContent := strings.Repeat("A", 1024*1024)
	if _, err := tmpFile.WriteString(largeContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	hash, err := HashFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("HashFile failed: %v", err)
	}

	// Verify against expected hash
	expected := sha256.Sum256([]byte(largeContent))
	if !bytesEqual(hash, expected[:]) {
		t.Errorf("HashFile of large file returned %x, want %x", hash, expected)
	}
}

func TestHashFile_NonExistent(t *testing.T) {
	_, err := HashFile("/nonexistent/file/path.txt")
	if err == nil {
		t.Error("HashFile should return error for non-existent file")
	}
}

func TestVerifyFile_Match(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "verifytest-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := "verify me"
	tmpFile.WriteString(content)
	tmpFile.Close()

	expectedHashBytes := sha256.Sum256([]byte(content))
	expectedHash := hex.EncodeToString(expectedHashBytes[:])

	match, err := VerifyFile(tmpFile.Name(), expectedHash)
	if err != nil {
		t.Fatalf("VerifyFile failed: %v", err)
	}
	if !match {
		t.Error("VerifyFile should return true for matching hash")
	}
}

func TestVerifyFile_NoMatch(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "verifytest-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString("actual content")
	tmpFile.Close()

	wrongHash := "0000000000000000000000000000000000000000000000000000000000000000"

	match, err := VerifyFile(tmpFile.Name(), wrongHash)
	if err != nil {
		t.Fatalf("VerifyFile failed: %v", err)
	}
	if match {
		t.Error("VerifyFile should return false for non-matching hash")
	}
}

func TestVerifyFile_CaseInsensitive(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "verifytest-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := "case test"
	tmpFile.WriteString(content)
	tmpFile.Close()

	hashBytes := sha256.Sum256([]byte(content))
	hashLower := hex.EncodeToString(hashBytes[:])
	hashUpper := strings.ToUpper(hashLower)

	// Test lowercase
	match1, err := VerifyFile(tmpFile.Name(), hashLower)
	if err != nil || !match1 {
		t.Error("VerifyFile should accept lowercase hash")
	}

	// Test uppercase
	match2, err := VerifyFile(tmpFile.Name(), hashUpper)
	if err != nil || !match2 {
		t.Error("VerifyFile should accept uppercase hash (case-insensitive)")
	}
}

func TestHashIncremental_Single(t *testing.T) {
	input := "hello"
	result := HashIncremental(input)
	expected := HashString(input)

	if result != expected {
		t.Errorf("HashIncremental(%q) = %s, want %s", input, result, expected)
	}
}

func TestHashIncremental_Multiple(t *testing.T) {
	result := HashIncremental("hello", "world")
	expected := HashString("helloworld")

	if result != expected {
		t.Errorf("HashIncremental(\"hello\", \"world\") = %s, want %s (hash of 'helloworld')", result, expected)
	}
}

func TestHashIncremental_Many(t *testing.T) {
	result := HashIncremental("a", "b", "c", "d", "e", "f")
	expected := HashString("abcdef")

	if result != expected {
		t.Errorf("HashIncremental multiple parts should equal hash of concatenation")
	}
}

func TestHashIncremental_Empty(t *testing.T) {
	result := HashIncremental()
	expected := HashString("")

	if result != expected {
		t.Errorf("HashIncremental() (no args) should equal hash of empty string")
	}
}

func TestHashIncremental_EmptyStrings(t *testing.T) {
	result := HashIncremental("", "", "")
	expected := HashString("")

	if result != expected {
		t.Errorf("HashIncremental with empty strings should equal hash of empty string")
	}
}

func TestCompareHashes_Equal(t *testing.T) {
	hash := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"

	if !CompareHashes(hash, hash) {
		t.Error("CompareHashes should return true for identical hashes")
	}
}

func TestCompareHashes_Different(t *testing.T) {
	hash1 := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	hash2 := "486ea46224d1bb4fb680f34f7c9ad96a8f24ec88be73ea8e5a6c65260e9cb8a7"

	if CompareHashes(hash1, hash2) {
		t.Error("CompareHashes should return false for different hashes")
	}
}

func TestCompareHashes_CaseInsensitive(t *testing.T) {
	hashLower := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	hashUpper := "2CF24DBA5FB0A30E26E83B2AC5B9E29E1B161E5C1FA7425E73043362938B9824"
	hashMixed := "2Cf24DbA5fB0a30E26e83B2aC5b9E29e1B161E5c1Fa7425E73043362938B9824"

	if !CompareHashes(hashLower, hashUpper) {
		t.Error("CompareHashes should be case-insensitive")
	}
	if !CompareHashes(hashLower, hashMixed) {
		t.Error("CompareHashes should be case-insensitive")
	}
}

func TestCompareHashes_InvalidHex(t *testing.T) {
	validHash := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	invalidHash := "ZZZZ"

	if CompareHashes(validHash, invalidHash) {
		t.Error("CompareHashes should return false for invalid hex")
	}
}

func TestCompareHashes_DifferentLengths(t *testing.T) {
	hash1 := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	hash2 := "2cf24dba"

	if CompareHashes(hash1, hash2) {
		t.Error("CompareHashes should return false for different length hashes")
	}
}

// Benchmark tests

func BenchmarkHashString(b *testing.B) {
	input := "The quick brown fox jumps over the lazy dog"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashString(input)
	}
}

func BenchmarkHashIncremental(b *testing.B) {
	parts := []string{"The", " quick", " brown", " fox", " jumps", " over", " the", " lazy", " dog"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashIncremental(parts...)
	}
}

func BenchmarkHashFile_Small(b *testing.B) {
	tmpFile, err := os.CreateTemp("", "benchmark-small-*.txt")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString("Small file content for benchmarking")
	tmpFile.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashFile(tmpFile.Name())
	}
}

func BenchmarkHashFile_Large(b *testing.B) {
	tmpFile, err := os.CreateTemp("", "benchmark-large-*.txt")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// 1MB file
	largeContent := strings.Repeat("A", 1024*1024)
	tmpFile.WriteString(largeContent)
	tmpFile.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashFile(tmpFile.Name())
	}
}

func BenchmarkCompareHashes(b *testing.B) {
	hash1 := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	hash2 := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CompareHashes(hash1, hash2)
	}
}

// Helper function
func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
