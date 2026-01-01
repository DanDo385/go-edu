package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/example/go-10x-minis/minis/39-sha256-hasher/exercise"
)

func main() {
	fmt.Println("=== SHA-256 Cryptographic Hash Demonstrations ===\n")

	demo1_BasicHashing()
	demo2_AvalancheEffect()
	demo3_FileHashing()
	demo4_IncrementalHashing()
	demo5_DigestReuse()
	demo6_CollisionResistance()
	demo7_HashComparison()
}

// Demo 1: Basic string hashing
func demo1_BasicHashing() {
	fmt.Println("--- Demo 1: Basic String Hashing ---")

	inputs := []string{
		"hello",
		"Hello",
		"hello world",
		"",
		"The quick brown fox jumps over the lazy dog",
	}

	for _, input := range inputs {
		hash := sha256.Sum256([]byte(input))
		hexHash := hex.EncodeToString(hash[:])

		fmt.Printf("Input:  %q\n", input)
		fmt.Printf("SHA-256: %s\n", hexHash)
		fmt.Printf("Length:  %d characters (always 64 hex chars = 256 bits)\n\n", len(hexHash))
	}
}

// Demo 2: Avalanche effect - tiny input change causes massive hash change
func demo2_AvalancheEffect() {
	fmt.Println("--- Demo 2: Avalanche Effect ---")
	fmt.Println("Tiny input changes cause completely different hashes\n")

	pairs := []struct {
		input1, input2 string
		description    string
	}{
		{"hello", "Hello", "Changed case of first letter (h → H)"},
		{"hello", "hallo", "Changed one letter (e → a)"},
		{"hello", "hello!", "Added exclamation mark"},
		{"hello", "hell", "Removed one letter"},
		{"bitcoin", "bitcain", "Changed one letter (o → a)"},
	}

	for _, pair := range pairs {
		hash1 := sha256.Sum256([]byte(pair.input1))
		hash2 := sha256.Sum256([]byte(pair.input2))

		hex1 := hex.EncodeToString(hash1[:])
		hex2 := hex.EncodeToString(hash2[:])

		// Count different bits
		diffBits := countDifferentBits(hash1[:], hash2[:])
		totalBits := len(hash1) * 8

		fmt.Printf("Change: %s\n", pair.description)
		fmt.Printf("  Input 1: %q\n", pair.input1)
		fmt.Printf("  Hash 1:  %s\n", hex1)
		fmt.Printf("  Input 2: %q\n", pair.input2)
		fmt.Printf("  Hash 2:  %s\n", hex2)
		fmt.Printf("  Bits changed: %d/%d (%.2f%%)\n", diffBits, totalBits, float64(diffBits)/float64(totalBits)*100)
		fmt.Println()
	}
}

// Demo 3: File hashing
func demo3_FileHashing() {
	fmt.Println("--- Demo 3: File Hashing ---")
	fmt.Println("Creating temporary files and hashing their contents\n")

	// Create temporary test files
	testFiles := []struct {
		name    string
		content string
	}{
		{"test1.txt", "Hello, World!"},
		{"test2.txt", "Hello, World!"},
		{"test3.txt", "Hello, world!"},
		{"large.txt", strings.Repeat("A", 1000000)},
	}

	tmpDir, err := os.MkdirTemp("", "sha256-demo-")
	if err != nil {
		fmt.Printf("Error creating temp dir: %v\n", err)
		return
	}
	defer os.RemoveAll(tmpDir)

	for _, tf := range testFiles {
		path := tmpDir + "/" + tf.name
		err := os.WriteFile(path, []byte(tf.content), 0644)
		if err != nil {
			fmt.Printf("Error writing file %s: %v\n", tf.name, err)
			continue
		}

		hash, err := exercise.HashFile(path)
		if err != nil {
			fmt.Printf("Error hashing file %s: %v\n", tf.name, err)
			continue
		}

		hexHash := hex.EncodeToString(hash)
		fileInfo, _ := os.Stat(path)

		fmt.Printf("File: %s (%d bytes)\n", tf.name, fileInfo.Size())
		fmt.Printf("SHA-256: %s\n", hexHash)

		// Note duplicates
		if tf.name == "test2.txt" {
			fmt.Printf("  → Same hash as test1.txt (identical content)\n")
		} else if tf.name == "test3.txt" {
			fmt.Printf("  → Different hash from test1.txt (different case)\n")
		}
		fmt.Println()
	}

	fmt.Println("Key observation: Files with identical content have identical hashes")
	fmt.Println("                 (Used for deduplication and integrity verification)\n")
}

// Demo 4: Incremental hashing (streaming)
func demo4_IncrementalHashing() {
	fmt.Println("--- Demo 4: Incremental Hashing ---")
	fmt.Println("Demonstrates hashing data in chunks (like streaming a large file)\n")

	data := "Hello, World!"

	// Method 1: Hash all at once
	hashAllAtOnce := sha256.Sum256([]byte(data))

	// Method 2: Hash incrementally
	h := sha256.New()
	h.Write([]byte("Hello, "))
	h.Write([]byte("World!"))
	hashIncremental := h.Sum(nil)

	fmt.Printf("Original data: %q\n\n", data)

	fmt.Printf("Method 1 - All at once:\n")
	fmt.Printf("  sha256.Sum256(data)\n")
	fmt.Printf("  Hash: %s\n\n", hex.EncodeToString(hashAllAtOnce[:]))

	fmt.Printf("Method 2 - Incremental:\n")
	fmt.Printf("  h.Write(\"Hello, \")\n")
	fmt.Printf("  h.Write(\"World!\")\n")
	fmt.Printf("  Hash: %s\n\n", hex.EncodeToString(hashIncremental))

	if hex.EncodeToString(hashAllAtOnce[:]) == hex.EncodeToString(hashIncremental) {
		fmt.Println("✓ Both methods produce identical hashes!")
		fmt.Println("  This allows efficient hashing of large files without loading them entirely into memory\n")
	}

	// Demonstrate with reader
	fmt.Println("Method 3 - Using io.Copy (common pattern for files):")
	h2 := sha256.New()
	reader := strings.NewReader(data)
	io.Copy(h2, reader)
	hashFromReader := h2.Sum(nil)
	fmt.Printf("  io.Copy(hash, reader)\n")
	fmt.Printf("  Hash: %s\n\n", hex.EncodeToString(hashFromReader))
}

// Demo 5: Digest reuse and reset
func demo5_DigestReuse() {
	fmt.Println("--- Demo 5: Hash Digest Reuse and Reset ---")
	fmt.Println("Common mistake: Forgetting to reset hash state between uses\n")

	h := sha256.New()

	// First hash
	fmt.Println("Step 1: Hash \"hello\"")
	h.Write([]byte("hello"))
	hash1 := h.Sum(nil)
	fmt.Printf("  Hash: %s\n\n", hex.EncodeToString(hash1))

	// MISTAKE: Continuing without reset
	fmt.Println("Step 2: Write \"world\" WITHOUT reset (WRONG)")
	h.Write([]byte("world"))
	hash2 := h.Sum(nil)
	fmt.Printf("  Hash: %s\n", hex.EncodeToString(hash2))
	expectedConcatenated := sha256.Sum256([]byte("helloworld"))
	fmt.Printf("  Expected (Hash of \"helloworld\"): %s\n", hex.EncodeToString(expectedConcatenated[:]))
	fmt.Printf("  → This is Hash(\"helloworld\"), not Hash(\"world\")!\n\n")

	// CORRECT: Reset before new hash
	fmt.Println("Step 3: Reset and hash \"world\" (CORRECT)")
	h.Reset()
	h.Write([]byte("world"))
	hash3 := h.Sum(nil)
	fmt.Printf("  Hash: %s\n", hex.EncodeToString(hash3))
	expectedWorld := sha256.Sum256([]byte("world"))
	fmt.Printf("  Expected (Hash of \"world\"): %s\n", hex.EncodeToString(expectedWorld[:]))
	fmt.Printf("  ✓ Correct! Always call Reset() before reusing hash.Hash\n\n")
}

// Demo 6: Collision resistance demonstration
func demo6_CollisionResistance() {
	fmt.Println("--- Demo 6: Collision Resistance ---")
	fmt.Println("Finding collisions is computationally infeasible\n")

	// Generate hashes for many inputs
	inputs := []string{
		"hello",
		"world",
		"bitcoin",
		"blockchain",
		"ethereum",
		"satoshi",
		"nakamoto",
	}

	// Try to find similar hashes (won't find exact collisions)
	fmt.Println("Attempting to find hash collisions (spoiler: won't find any):")
	hashes := make(map[string]string)

	for _, input := range inputs {
		hash := sha256.Sum256([]byte(input))
		hexHash := hex.EncodeToString(hash[:])
		hashes[hexHash] = input

		// Show first 16 characters
		fmt.Printf("  %s → %s...\n", input, hexHash[:16])
	}

	fmt.Printf("\nGenerated %d hashes, found %d unique hashes (no collisions)\n", len(inputs), len(hashes))
	fmt.Println("\nTo find a collision in SHA-256:")
	fmt.Println("  - Need approximately 2^128 hash computations (birthday attack)")
	fmt.Println("  - With 1 trillion hashes/second: ~10 billion years")
	fmt.Println("  - Current technology: Infeasible\n")
}

// Demo 7: Hash comparison for integrity verification
func demo7_HashComparison() {
	fmt.Println("--- Demo 7: Hash Comparison for Integrity Verification ---")
	fmt.Println("Simulating file download verification\n")

	// Simulate original file
	originalContent := "This is the legitimate Ubuntu ISO file contents..."
	originalHash := sha256.Sum256([]byte(originalContent))
	publishedHash := hex.EncodeToString(originalHash[:])

	fmt.Println("Scenario: Verifying downloaded file integrity")
	fmt.Printf("Published SHA-256 (from ubuntu.com): %s\n\n", publishedHash)

	// Test case 1: Identical download
	fmt.Println("Test 1: Perfect download (file unchanged)")
	downloadedContent1 := originalContent
	downloadedHash1 := sha256.Sum256([]byte(downloadedContent1))
	match1 := hex.EncodeToString(downloadedHash1[:]) == publishedHash
	fmt.Printf("  Downloaded hash: %s\n", hex.EncodeToString(downloadedHash1[:]))
	if match1 {
		fmt.Println("  ✓ VERIFIED: File is authentic and uncorrupted\n")
	}

	// Test case 2: Corrupted download
	fmt.Println("Test 2: Corrupted download (one bit flipped)")
	downloadedContent2 := "This is the legitimate Ubuntu ISO file content..."
	downloadedHash2 := sha256.Sum256([]byte(downloadedContent2))
	match2 := hex.EncodeToString(downloadedHash2[:]) == publishedHash
	fmt.Printf("  Downloaded hash: %s\n", hex.EncodeToString(downloadedHash2[:]))
	if !match2 {
		fmt.Println("  ✗ VERIFICATION FAILED: File is corrupted or tampered!")
		fmt.Println("    → Do not trust this file, re-download from official source\n")
	}

	// Test case 3: Malicious replacement
	fmt.Println("Test 3: Malicious file (replaced with malware)")
	downloadedContent3 := "This is malware pretending to be Ubuntu..."
	downloadedHash3 := sha256.Sum256([]byte(downloadedContent3))
	match3 := hex.EncodeToString(downloadedHash3[:]) == publishedHash
	fmt.Printf("  Downloaded hash: %s\n", hex.EncodeToString(downloadedHash3[:]))
	if !match3 {
		fmt.Println("  ✗ VERIFICATION FAILED: File does not match published hash!")
		fmt.Println("    → This could be a malicious replacement, DELETE immediately\n")
	}

	fmt.Println("Key Takeaway: Always verify SHA-256 checksums when downloading software")
	fmt.Println("              from the internet, especially operating systems and security tools.\n")
}

// Helper: Count different bits between two byte slices
func countDifferentBits(a, b []byte) int {
	if len(a) != len(b) {
		return -1
	}

	count := 0
	for i := 0; i < len(a); i++ {
		xor := a[i] ^ b[i]
		// Count set bits in XOR result
		for xor != 0 {
			count += int(xor & 1)
			xor >>= 1
		}
	}
	return count
}
