package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// DEMONSTRATION: Merkle Tree Construction and Verification
//
// This program demonstrates:
// 1. Building a Merkle tree from data blocks
// 2. Generating proofs of inclusion
// 3. Verifying proofs without the full dataset
// 4. Detecting tampered data

// MerkleTree represents a complete Merkle tree structure
type MerkleTree struct {
	Root   []byte     // The root hash
	Leaves [][]byte   // All leaf hashes
	Levels [][][]byte // All levels of the tree (for proof generation)
}

// ProofNode represents one element in a Merkle proof
type ProofNode struct {
	Hash   []byte
	IsLeft bool // True if this hash should be on the left when combining
}

// MerkleProof represents a proof of inclusion for a specific data block
type MerkleProof struct {
	LeafIndex int         // Index of the data in the original dataset
	Siblings  []ProofNode // Hashes needed to reconstruct the root
}

func main() {
	fmt.Println("=== MERKLE TREE DEMONSTRATION ===\n")

	// Demo 1: Building a Merkle tree
	demo1_BuildingTree()

	fmt.Println()

	// Demo 2: Generating and verifying proofs
	demo2_ProofGeneration()

	fmt.Println()

	// Demo 3: Detecting tampered data
	demo3_TamperDetection()

	fmt.Println()

	// Demo 4: Blockchain-style transaction batching
	demo4_BlockchainTransactions()
}

// demo1_BuildingTree shows how to construct a Merkle tree from data
func demo1_BuildingTree() {
	fmt.Println("--- Demo 1: Building a Merkle Tree ---")

	// Sample data: Four files
	data := [][]byte{
		[]byte("File A: Important document"),
		[]byte("File B: Contract details"),
		[]byte("File C: Receipt #12345"),
		[]byte("File D: User agreement"),
	}

	fmt.Println("Original Data:")
	for i, d := range data {
		fmt.Printf("  [%d] %s\n", i, string(d))
	}

	// Build the tree
	tree := BuildMerkleTree(data)

	fmt.Printf("\nMerkle Root: %s\n", hashToString(tree.Root))
	fmt.Println("\nTree Structure:")
	printTree(tree)

	fmt.Println("\nKEY INSIGHT:")
	fmt.Println("The root hash represents ALL data below it.")
	fmt.Println("If even one byte changes, the root hash changes completely!")
}

// demo2_ProofGeneration shows proof creation and verification
func demo2_ProofGeneration() {
	fmt.Println("--- Demo 2: Proof of Inclusion ---")

	data := [][]byte{
		[]byte("Transaction 1: Alice pays Bob $10"),
		[]byte("Transaction 2: Bob pays Charlie $5"),
		[]byte("Transaction 3: Charlie pays Dave $3"),
		[]byte("Transaction 4: Dave pays Eve $1"),
	}

	tree := BuildMerkleTree(data)
	fmt.Printf("Merkle Root: %s\n", hashToString(tree.Root))

	// Generate proof for Transaction 2
	targetIndex := 1
	targetData := data[targetIndex]

	fmt.Printf("\nProving that this transaction is in the tree:\n")
	fmt.Printf("  \"%s\"\n", string(targetData))

	proof := GenerateProof(tree, targetIndex)

	fmt.Printf("\nProof contains %d sibling hashes:\n", len(proof.Siblings))
	for i, node := range proof.Siblings {
		position := "right"
		if node.IsLeft {
			position = "left"
		}
		fmt.Printf("  [%d] %s (position: %s)\n", i, hashToString(node.Hash), position)
	}

	// Verify the proof
	isValid := VerifyProof(targetData, proof, tree.Root)
	fmt.Printf("\nProof verification: %v\n", isValid)

	// Calculate efficiency
	totalDataSize := 0
	for _, d := range data {
		totalDataSize += len(d)
	}
	proofSize := len(targetData) + len(proof.Siblings)*32

	fmt.Printf("\nEFFICIENCY ANALYSIS:\n")
	fmt.Printf("  Total data size: %d bytes\n", totalDataSize)
	fmt.Printf("  Proof size: %d bytes\n", proofSize)
	fmt.Printf("  Savings: %.1f%%\n", (1.0-float64(proofSize)/float64(totalDataSize))*100)
}

// demo3_TamperDetection shows how Merkle trees detect tampering
func demo3_TamperDetection() {
	fmt.Println("--- Demo 3: Tamper Detection ---")

	originalData := [][]byte{
		[]byte("Block 1: Alice has 100 coins"),
		[]byte("Block 2: Bob has 50 coins"),
		[]byte("Block 3: Charlie has 75 coins"),
		[]byte("Block 4: Dave has 25 coins"),
	}

	tree := BuildMerkleTree(originalData)
	fmt.Printf("Original Merkle Root: %s\n", hashToString(tree.Root))

	// An attacker tries to tamper with the data
	tamperedData := [][]byte{
		[]byte("Block 1: Alice has 100 coins"),
		[]byte("Block 2: Bob has 50 coins"),
		[]byte("Block 3: Charlie has 750 coins"), // CHANGED: 75 → 750
		[]byte("Block 4: Dave has 25 coins"),
	}

	tamperedTree := BuildMerkleTree(tamperedData)
	fmt.Printf("Tampered Merkle Root: %s\n", hashToString(tamperedTree.Root))

	fmt.Printf("\nRoots match: %v\n", hashesEqual(tree.Root, tamperedTree.Root))

	fmt.Println("\nKEY INSIGHT:")
	fmt.Println("Even a tiny change (75 → 750) completely changes the root hash.")
	fmt.Println("This makes Merkle trees perfect for detecting any data tampering!")

	// Show which block was tampered
	fmt.Println("\nIdentifying the tampered block:")
	for i := range originalData {
		origHash := hash(originalData[i])
		tampHash := hash(tamperedData[i])
		if !hashesEqual(origHash, tampHash) {
			fmt.Printf("  Block %d was modified!\n", i)
			fmt.Printf("    Original: %s\n", string(originalData[i]))
			fmt.Printf("    Tampered: %s\n", string(tamperedData[i]))
		}
	}
}

// demo4_BlockchainTransactions simulates how Bitcoin uses Merkle trees
func demo4_BlockchainTransactions() {
	fmt.Println("--- Demo 4: Blockchain-Style Transaction Verification ---")

	// Simulate a block with 8 transactions
	transactions := [][]byte{
		[]byte("tx1: Alice → Bob: 5 BTC"),
		[]byte("tx2: Bob → Charlie: 2 BTC"),
		[]byte("tx3: Charlie → Dave: 1 BTC"),
		[]byte("tx4: Dave → Eve: 0.5 BTC"),
		[]byte("tx5: Eve → Frank: 0.3 BTC"),
		[]byte("tx6: Frank → Grace: 0.2 BTC"),
		[]byte("tx7: Grace → Henry: 0.1 BTC"),
		[]byte("tx8: Henry → Alice: 0.05 BTC"),
	}

	fmt.Printf("Block contains %d transactions\n", len(transactions))

	// Build Merkle tree (what full nodes do)
	tree := BuildMerkleTree(transactions)
	merkleRoot := tree.Root

	fmt.Printf("Merkle Root (stored in block header): %s\n", hashToString(merkleRoot))

	fmt.Println("\n--- Light Client (SPV) Verification ---")
	fmt.Println("A light client (e.g., mobile wallet) wants to verify tx3 without downloading all transactions.")

	// Full node generates proof
	txIndex := 2 // tx3
	proof := GenerateProof(tree, txIndex)

	fmt.Printf("\nFull node sends to light client:\n")
	fmt.Printf("  1. The transaction: %s\n", string(transactions[txIndex]))
	fmt.Printf("  2. Merkle proof: %d hashes (%d bytes)\n", len(proof.Siblings), len(proof.Siblings)*32)

	// Light client verifies (only needs the root from block header)
	isValid := VerifyProof(transactions[txIndex], proof, merkleRoot)

	fmt.Printf("\nLight client verification: %v\n", isValid)

	fmt.Println("\nSCALABILITY ANALYSIS:")
	fmt.Printf("  Transactions in block: %d\n", len(transactions))
	fmt.Printf("  Proof size: %d hashes\n", len(proof.Siblings))
	fmt.Printf("  Light client downloads: ~%d bytes (vs ~%d bytes for all transactions)\n",
		len(transactions[txIndex])+len(proof.Siblings)*32,
		totalSize(transactions))

	// Show scalability
	fmt.Println("\n  Proof size growth (logarithmic):")
	for n := 100; n <= 100000; n *= 10 {
		proofLen := ceilLog2(n)
		fmt.Printf("    %7d transactions → %2d hashes in proof\n", n, proofLen)
	}
}

// ========================================
// CORE MERKLE TREE IMPLEMENTATION
// ========================================

// BuildMerkleTree constructs a Merkle tree from data blocks
func BuildMerkleTree(data [][]byte) *MerkleTree {
	if len(data) == 0 {
		// Empty tree: root is hash of empty string
		emptyHash := hash([]byte{})
		return &MerkleTree{
			Root:   emptyHash,
			Leaves: [][]byte{},
			Levels: [][][]byte{{emptyHash}},
		}
	}

	// Level 0: Hash all data blocks (leaf nodes)
	leaves := make([][]byte, len(data))
	for i, d := range data {
		leaves[i] = hash(d)
	}

	// Build tree level by level
	levels := [][][]byte{leaves}
	currentLevel := leaves

	for len(currentLevel) > 1 {
		nextLevel := [][]byte{}

		// If odd number of nodes, duplicate the last one
		if len(currentLevel)%2 == 1 {
			currentLevel = append(currentLevel, currentLevel[len(currentLevel)-1])
		}

		// Pair and hash
		for i := 0; i < len(currentLevel); i += 2 {
			left := currentLevel[i]
			right := currentLevel[i+1]
			parentHash := hash(append(left, right...))
			nextLevel = append(nextLevel, parentHash)
		}

		levels = append(levels, nextLevel)
		currentLevel = nextLevel
	}

	return &MerkleTree{
		Root:   currentLevel[0],
		Leaves: leaves,
		Levels: levels,
	}
}

// GenerateProof creates a proof of inclusion for the data at the given index
func GenerateProof(tree *MerkleTree, index int) *MerkleProof {
	if index < 0 || index >= len(tree.Leaves) {
		return nil
	}

	proof := &MerkleProof{
		LeafIndex: index,
		Siblings:  []ProofNode{},
	}

	currentIndex := index

	// Traverse from leaf to root, collecting sibling hashes
	for level := 0; level < len(tree.Levels)-1; level++ {
		currentLevelNodes := tree.Levels[level]

		// Handle odd number of nodes (duplicated last node)
		if len(currentLevelNodes)%2 == 1 {
			currentLevelNodes = append(currentLevelNodes, currentLevelNodes[len(currentLevelNodes)-1])
		}

		// Find sibling
		var siblingIndex int
		if currentIndex%2 == 0 {
			// Current is left child, sibling is right
			siblingIndex = currentIndex + 1
		} else {
			// Current is right child, sibling is left
			siblingIndex = currentIndex - 1
		}

		// Add sibling to proof
		siblingHash := currentLevelNodes[siblingIndex]
		isLeft := siblingIndex < currentIndex

		proof.Siblings = append(proof.Siblings, ProofNode{
			Hash:   siblingHash,
			IsLeft: isLeft,
		})

		// Move to parent index
		currentIndex = currentIndex / 2
	}

	return proof
}

// VerifyProof verifies that data is in the tree by checking the proof against the root
func VerifyProof(data []byte, proof *MerkleProof, root []byte) bool {
	if proof == nil {
		return false
	}

	// Start with hash of the data
	currentHash := hash(data)

	// Combine with siblings to reconstruct the root
	for _, sibling := range proof.Siblings {
		if sibling.IsLeft {
			// Sibling goes on the left
			currentHash = hash(append(sibling.Hash, currentHash...))
		} else {
			// Sibling goes on the right
			currentHash = hash(append(currentHash, sibling.Hash...))
		}
	}

	// Compare computed root with expected root
	return hashesEqual(currentHash, root)
}

// ========================================
// UTILITY FUNCTIONS
// ========================================

// hash computes SHA-256 hash of data
func hash(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

// hashToString converts hash bytes to hex string (shortened for display)
func hashToString(h []byte) string {
	full := hex.EncodeToString(h)
	if len(full) > 16 {
		return full[:16] + "..." // Show first 16 chars
	}
	return full
}

// hashesEqual compares two hashes
func hashesEqual(a, b []byte) bool {
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

// printTree displays the tree structure visually
func printTree(tree *MerkleTree) {
	for level := len(tree.Levels) - 1; level >= 0; level-- {
		indent := strings.Repeat("  ", len(tree.Levels)-level-1)
		fmt.Printf("  Level %d:%s", level, indent)

		for i, h := range tree.Levels[level] {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%s", hashToString(h))
		}
		fmt.Println()
	}
}

// totalSize calculates total bytes in a slice of byte slices
func totalSize(data [][]byte) int {
	total := 0
	for _, d := range data {
		total += len(d)
	}
	return total
}

// ceilLog2 returns ceil(log2(n))
func ceilLog2(n int) int {
	if n <= 0 {
		return 0
	}
	count := 0
	n = n - 1
	for n > 0 {
		n >>= 1
		count++
	}
	return count
}
