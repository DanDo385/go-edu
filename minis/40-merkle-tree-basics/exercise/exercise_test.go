package exercise

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

// ========================================
// BUILDMERKLETREE TESTS
// ========================================

func TestBuildMerkleTree_EmptyData(t *testing.T) {
	data := [][]byte{}
	tree := BuildMerkleTree(data)

	if tree == nil {
		t.Fatal("BuildMerkleTree returned nil for empty data")
	}

	if tree.Root == nil {
		t.Fatal("Root should not be nil for empty tree")
	}

	// Root should be hash of empty byte slice
	expectedRoot := sha256.Sum256([]byte{})
	if !hashesEqual(tree.Root, expectedRoot[:]) {
		t.Errorf("Empty tree root incorrect.\nGot:  %x\nWant: %x", tree.Root, expectedRoot)
	}
}

func TestBuildMerkleTree_SingleBlock(t *testing.T) {
	data := [][]byte{[]byte("Single block")}
	tree := BuildMerkleTree(data)

	if tree == nil {
		t.Fatal("BuildMerkleTree returned nil")
	}

	// Root should be hash of the single block
	expectedRoot := sha256.Sum256(data[0])
	if !hashesEqual(tree.Root, expectedRoot[:]) {
		t.Errorf("Single block tree root incorrect.\nGot:  %x\nWant: %x", tree.Root, expectedRoot)
	}

	// Should have one leaf
	if len(tree.Leaves) != 1 {
		t.Errorf("Expected 1 leaf, got %d", len(tree.Leaves))
	}
}

func TestBuildMerkleTree_TwoBlocks(t *testing.T) {
	data := [][]byte{
		[]byte("Block A"),
		[]byte("Block B"),
	}
	tree := BuildMerkleTree(data)

	if tree == nil {
		t.Fatal("BuildMerkleTree returned nil")
	}

	// Verify structure
	if len(tree.Leaves) != 2 {
		t.Errorf("Expected 2 leaves, got %d", len(tree.Leaves))
	}

	// Manually compute expected root
	h1 := sha256.Sum256(data[0])
	h2 := sha256.Sum256(data[1])
	combined := append(h1[:], h2[:]...)
	expectedRoot := sha256.Sum256(combined)

	if !hashesEqual(tree.Root, expectedRoot[:]) {
		t.Errorf("Two block tree root incorrect.\nGot:  %x\nWant: %x", tree.Root, expectedRoot)
	}

	// Should have 2 levels (leaves + root)
	if len(tree.Levels) != 2 {
		t.Errorf("Expected 2 levels, got %d", len(tree.Levels))
	}
}

func TestBuildMerkleTree_FourBlocks(t *testing.T) {
	data := [][]byte{
		[]byte("Block 1"),
		[]byte("Block 2"),
		[]byte("Block 3"),
		[]byte("Block 4"),
	}
	tree := BuildMerkleTree(data)

	if tree == nil {
		t.Fatal("BuildMerkleTree returned nil")
	}

	// Verify structure
	if len(tree.Leaves) != 4 {
		t.Errorf("Expected 4 leaves, got %d", len(tree.Leaves))
	}

	// Should have 3 levels: leaves (4), middle (2), root (1)
	if len(tree.Levels) != 3 {
		t.Errorf("Expected 3 levels, got %d", len(tree.Levels))
	}

	// Manually compute expected root
	h1 := sha256.Sum256(data[0])
	h2 := sha256.Sum256(data[1])
	h3 := sha256.Sum256(data[2])
	h4 := sha256.Sum256(data[3])

	h12 := sha256.Sum256(append(h1[:], h2[:]...))
	h34 := sha256.Sum256(append(h3[:], h4[:]...))
	expectedRoot := sha256.Sum256(append(h12[:], h34[:]...))

	if !hashesEqual(tree.Root, expectedRoot[:]) {
		t.Errorf("Four block tree root incorrect.\nGot:  %x\nWant: %x", tree.Root, expectedRoot)
	}
}

func TestBuildMerkleTree_OddBlocks(t *testing.T) {
	// Test with 3 blocks (odd number)
	data := [][]byte{
		[]byte("Block 1"),
		[]byte("Block 2"),
		[]byte("Block 3"),
	}
	tree := BuildMerkleTree(data)

	if tree == nil {
		t.Fatal("BuildMerkleTree returned nil")
	}

	if len(tree.Leaves) != 3 {
		t.Errorf("Expected 3 leaves, got %d", len(tree.Leaves))
	}

	// The algorithm should duplicate the last node
	// Level 0: [H1, H2, H3, H3] (duplicated H3)
	// Level 1: [H(H1+H2), H(H3+H3)]
	// Level 2: [root]

	h1 := sha256.Sum256(data[0])
	h2 := sha256.Sum256(data[1])
	h3 := sha256.Sum256(data[2])

	h12 := sha256.Sum256(append(h1[:], h2[:]...))
	h33 := sha256.Sum256(append(h3[:], h3[:]...)) // H3 duplicated
	expectedRoot := sha256.Sum256(append(h12[:], h33[:]...))

	if !hashesEqual(tree.Root, expectedRoot[:]) {
		t.Errorf("Odd block tree root incorrect.\nGot:  %x\nWant: %x", tree.Root, expectedRoot)
	}
}

func TestBuildMerkleTree_LargeTree(t *testing.T) {
	// Test with 8 blocks
	data := make([][]byte, 8)
	for i := 0; i < 8; i++ {
		data[i] = []byte{byte(i)}
	}

	tree := BuildMerkleTree(data)

	if tree == nil {
		t.Fatal("BuildMerkleTree returned nil")
	}

	if len(tree.Leaves) != 8 {
		t.Errorf("Expected 8 leaves, got %d", len(tree.Leaves))
	}

	// 8 blocks → 4 levels: 8 → 4 → 2 → 1
	if len(tree.Levels) != 4 {
		t.Errorf("Expected 4 levels, got %d", len(tree.Levels))
	}

	// Verify level sizes
	expectedSizes := []int{8, 4, 2, 1}
	for i, expected := range expectedSizes {
		if len(tree.Levels[i]) != expected {
			t.Errorf("Level %d: expected %d nodes, got %d", i, expected, len(tree.Levels[i]))
		}
	}
}

// ========================================
// GENERATEPROOF TESTS
// ========================================

func TestGenerateProof_InvalidIndex(t *testing.T) {
	data := [][]byte{[]byte("A"), []byte("B")}
	tree := BuildMerkleTree(data)

	// Test negative index
	proof := GenerateProof(tree, -1)
	if proof != nil {
		t.Error("Expected nil proof for negative index")
	}

	// Test index out of bounds
	proof = GenerateProof(tree, 2)
	if proof != nil {
		t.Error("Expected nil proof for out of bounds index")
	}
}

func TestGenerateProof_SingleBlock(t *testing.T) {
	data := [][]byte{[]byte("Single")}
	tree := BuildMerkleTree(data)

	proof := GenerateProof(tree, 0)
	if proof == nil {
		t.Fatal("GenerateProof returned nil")
	}

	// Single block tree: proof should be empty (data hash == root)
	if len(proof.Siblings) != 0 {
		t.Errorf("Expected 0 siblings for single block, got %d", len(proof.Siblings))
	}
}

func TestGenerateProof_TwoBlocks(t *testing.T) {
	data := [][]byte{[]byte("A"), []byte("B")}
	tree := BuildMerkleTree(data)

	// Proof for index 0 (A)
	proof := GenerateProof(tree, 0)
	if proof == nil {
		t.Fatal("GenerateProof returned nil")
	}

	// Should have 1 sibling (H(B))
	if len(proof.Siblings) != 1 {
		t.Errorf("Expected 1 sibling, got %d", len(proof.Siblings))
	}

	// Sibling should be H(B), on the right
	expectedSibling := sha256.Sum256(data[1])
	if !hashesEqual(proof.Siblings[0].Hash, expectedSibling[:]) {
		t.Errorf("Sibling hash incorrect.\nGot:  %x\nWant: %x", proof.Siblings[0].Hash, expectedSibling)
	}

	if proof.Siblings[0].IsLeft {
		t.Error("Sibling should be on the right (IsLeft=false)")
	}
}

func TestGenerateProof_FourBlocks(t *testing.T) {
	data := [][]byte{
		[]byte("Block 1"),
		[]byte("Block 2"),
		[]byte("Block 3"),
		[]byte("Block 4"),
	}
	tree := BuildMerkleTree(data)

	// Proof for index 1 (Block 2)
	proof := GenerateProof(tree, 1)
	if proof == nil {
		t.Fatal("GenerateProof returned nil")
	}

	// Should have 2 siblings (one from each level)
	if len(proof.Siblings) != 2 {
		t.Errorf("Expected 2 siblings, got %d", len(proof.Siblings))
	}

	// First sibling: H(Block 1), should be left
	h1 := sha256.Sum256(data[0])
	if !hashesEqual(proof.Siblings[0].Hash, h1[:]) {
		t.Errorf("First sibling incorrect.\nGot:  %x\nWant: %x", proof.Siblings[0].Hash, h1)
	}
	if !proof.Siblings[0].IsLeft {
		t.Error("First sibling should be on the left")
	}

	// Second sibling: H(H(Block3)+H(Block4)), should be right
	h3 := sha256.Sum256(data[2])
	h4 := sha256.Sum256(data[3])
	h34 := sha256.Sum256(append(h3[:], h4[:]...))
	if !hashesEqual(proof.Siblings[1].Hash, h34[:]) {
		t.Errorf("Second sibling incorrect.\nGot:  %x\nWant: %x", proof.Siblings[1].Hash, h34)
	}
	if proof.Siblings[1].IsLeft {
		t.Error("Second sibling should be on the right")
	}
}

func TestGenerateProof_AllIndices(t *testing.T) {
	data := [][]byte{
		[]byte("A"),
		[]byte("B"),
		[]byte("C"),
		[]byte("D"),
	}
	tree := BuildMerkleTree(data)

	// Generate proof for each index and verify
	for i := range data {
		proof := GenerateProof(tree, i)
		if proof == nil {
			t.Fatalf("GenerateProof returned nil for index %d", i)
		}

		// Verify the proof
		isValid := VerifyProof(data[i], proof, tree.Root)
		if !isValid {
			t.Errorf("Proof verification failed for index %d", i)
		}
	}
}

// ========================================
// VERIFYPROOF TESTS
// ========================================

func TestVerifyProof_NilProof(t *testing.T) {
	data := []byte("Test")
	root := sha256.Sum256(data)

	isValid := VerifyProof(data, nil, root[:])
	if isValid {
		t.Error("Expected false for nil proof")
	}
}

func TestVerifyProof_ValidProof(t *testing.T) {
	data := [][]byte{
		[]byte("Transaction 1"),
		[]byte("Transaction 2"),
		[]byte("Transaction 3"),
		[]byte("Transaction 4"),
	}
	tree := BuildMerkleTree(data)

	// Generate and verify proof for each block
	for i, block := range data {
		proof := GenerateProof(tree, i)
		isValid := VerifyProof(block, proof, tree.Root)

		if !isValid {
			t.Errorf("Valid proof failed verification for block %d", i)
		}
	}
}

func TestVerifyProof_TamperedData(t *testing.T) {
	data := [][]byte{
		[]byte("Block 1"),
		[]byte("Block 2"),
		[]byte("Block 3"),
		[]byte("Block 4"),
	}
	tree := BuildMerkleTree(data)

	// Generate proof for Block 2
	proof := GenerateProof(tree, 1)

	// Try to verify with tampered data
	tamperedData := []byte("Block 2 MODIFIED")
	isValid := VerifyProof(tamperedData, proof, tree.Root)

	if isValid {
		t.Error("Tampered data should fail verification")
	}
}

func TestVerifyProof_WrongRoot(t *testing.T) {
	data := [][]byte{
		[]byte("Block 1"),
		[]byte("Block 2"),
	}
	tree := BuildMerkleTree(data)
	proof := GenerateProof(tree, 0)

	// Use a wrong root
	wrongRoot := make([]byte, 32)
	for i := range wrongRoot {
		wrongRoot[i] = 0xFF
	}

	isValid := VerifyProof(data[0], proof, wrongRoot)
	if isValid {
		t.Error("Proof should fail with wrong root")
	}
}

func TestVerifyProof_SingleBlock(t *testing.T) {
	data := [][]byte{[]byte("Single block")}
	tree := BuildMerkleTree(data)
	proof := GenerateProof(tree, 0)

	isValid := VerifyProof(data[0], proof, tree.Root)
	if !isValid {
		t.Error("Single block proof should be valid")
	}
}

// ========================================
// GETMERKLEROOT TESTS
// ========================================

func TestGetMerkleRoot_SimpleCase(t *testing.T) {
	data := [][]byte{
		[]byte("A"),
		[]byte("B"),
	}

	root := GetMerkleRoot(data)
	if root == nil {
		t.Fatal("GetMerkleRoot returned nil")
	}

	// Verify it matches BuildMerkleTree result
	tree := BuildMerkleTree(data)
	if !hashesEqual(root, tree.Root) {
		t.Error("GetMerkleRoot doesn't match BuildMerkleTree root")
	}
}

func TestGetMerkleRoot_EmptyData(t *testing.T) {
	data := [][]byte{}
	root := GetMerkleRoot(data)

	if root == nil {
		t.Fatal("GetMerkleRoot returned nil for empty data")
	}

	expectedRoot := sha256.Sum256([]byte{})
	if !hashesEqual(root, expectedRoot[:]) {
		t.Errorf("Empty data root incorrect.\nGot:  %x\nWant: %x", root, expectedRoot)
	}
}

// ========================================
// INTEGRATION TESTS
// ========================================

func TestMerkleTree_BitcoinStyleTransactions(t *testing.T) {
	// Simulate a block with 8 transactions (like Bitcoin)
	transactions := [][]byte{
		[]byte("tx1: Alice → Bob: 1.5 BTC"),
		[]byte("tx2: Bob → Charlie: 0.5 BTC"),
		[]byte("tx3: Charlie → Dave: 0.3 BTC"),
		[]byte("tx4: Dave → Eve: 0.2 BTC"),
		[]byte("tx5: Eve → Frank: 0.1 BTC"),
		[]byte("tx6: Frank → Grace: 0.05 BTC"),
		[]byte("tx7: Grace → Henry: 0.03 BTC"),
		[]byte("tx8: Henry → Alice: 0.01 BTC"),
	}

	// Build Merkle tree (what full node does)
	tree := BuildMerkleTree(transactions)

	// Light client wants to verify tx3 without downloading all transactions
	txIndex := 2
	proof := GenerateProof(tree, txIndex)

	// Light client only needs:
	// 1. The transaction itself
	// 2. The Merkle proof
	// 3. The root hash (from block header)

	isValid := VerifyProof(transactions[txIndex], proof, tree.Root)
	if !isValid {
		t.Error("Bitcoin-style SPV verification failed")
	}

	// Verify proof size is logarithmic
	expectedProofLen := 3 // log2(8) = 3
	if len(proof.Siblings) != expectedProofLen {
		t.Errorf("Expected %d siblings for 8 transactions, got %d", expectedProofLen, len(proof.Siblings))
	}
}

func TestMerkleTree_TamperDetection(t *testing.T) {
	originalData := [][]byte{
		[]byte("Document 1"),
		[]byte("Document 2"),
		[]byte("Document 3"),
		[]byte("Document 4"),
	}

	originalTree := BuildMerkleTree(originalData)

	// Tamper with one document
	tamperedData := [][]byte{
		[]byte("Document 1"),
		[]byte("Document 2 MODIFIED"), // Changed
		[]byte("Document 3"),
		[]byte("Document 4"),
	}

	tamperedTree := BuildMerkleTree(tamperedData)

	// Roots should be different
	if hashesEqual(originalTree.Root, tamperedTree.Root) {
		t.Error("Tampering should change the root hash")
	}

	// Original proofs should fail with tampered tree
	for i, doc := range originalData {
		proof := GenerateProof(originalTree, i)
		isValid := VerifyProof(doc, proof, tamperedTree.Root)

		if isValid {
			t.Errorf("Original document %d should not verify against tampered tree", i)
		}
	}
}

func TestMerkleTree_ProofSizeScaling(t *testing.T) {
	// Test that proof size grows logarithmically with number of leaves
	testCases := []struct {
		numLeaves       int
		expectedProofLen int // ceil(log2(numLeaves))
	}{
		{1, 0},
		{2, 1},
		{4, 2},
		{8, 3},
		{16, 4},
		{32, 5},
		{64, 6},
		{128, 7},
	}

	for _, tc := range testCases {
		// Create data
		data := make([][]byte, tc.numLeaves)
		for i := 0; i < tc.numLeaves; i++ {
			data[i] = []byte{byte(i)}
		}

		tree := BuildMerkleTree(data)
		proof := GenerateProof(tree, 0)

		if len(proof.Siblings) != tc.expectedProofLen {
			t.Errorf("Leaves=%d: expected proof length %d, got %d",
				tc.numLeaves, tc.expectedProofLen, len(proof.Siblings))
		}
	}
}

// ========================================
// EDGE CASE TESTS
// ========================================

func TestMerkleTree_IdenticalBlocks(t *testing.T) {
	// All blocks are identical
	data := [][]byte{
		[]byte("Same"),
		[]byte("Same"),
		[]byte("Same"),
		[]byte("Same"),
	}

	tree := BuildMerkleTree(data)
	if tree == nil {
		t.Fatal("BuildMerkleTree returned nil")
	}

	// All leaves should have the same hash
	firstLeaf := tree.Leaves[0]
	for i, leaf := range tree.Leaves {
		if !hashesEqual(leaf, firstLeaf) {
			t.Errorf("Leaf %d differs from first leaf (all data is identical)", i)
		}
	}

	// Proofs should still work
	for i := range data {
		proof := GenerateProof(tree, i)
		isValid := VerifyProof(data[i], proof, tree.Root)
		if !isValid {
			t.Errorf("Proof failed for identical block %d", i)
		}
	}
}

func TestMerkleTree_LargeData(t *testing.T) {
	// Test with 1000 blocks
	numBlocks := 1000
	data := make([][]byte, numBlocks)
	for i := 0; i < numBlocks; i++ {
		data[i] = []byte{byte(i % 256), byte(i / 256)}
	}

	tree := BuildMerkleTree(data)
	if tree == nil {
		t.Fatal("BuildMerkleTree returned nil")
	}

	// Verify a few random indices
	testIndices := []int{0, 100, 500, 999}
	for _, idx := range testIndices {
		proof := GenerateProof(tree, idx)
		if proof == nil {
			t.Fatalf("GenerateProof returned nil for index %d", idx)
		}

		isValid := VerifyProof(data[idx], proof, tree.Root)
		if !isValid {
			t.Errorf("Proof verification failed for index %d in large tree", idx)
		}

		// Proof should be ~10 hashes for 1000 leaves (log2(1000) ≈ 10)
		if len(proof.Siblings) > 12 {
			t.Errorf("Proof too large: %d siblings for 1000 leaves", len(proof.Siblings))
		}
	}
}

// ========================================
// BENCHMARK TESTS
// ========================================

func BenchmarkBuildMerkleTree_100(b *testing.B) {
	data := make([][]byte, 100)
	for i := 0; i < 100; i++ {
		data[i] = []byte{byte(i)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BuildMerkleTree(data)
	}
}

func BenchmarkBuildMerkleTree_1000(b *testing.B) {
	data := make([][]byte, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = []byte{byte(i % 256), byte(i / 256)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BuildMerkleTree(data)
	}
}

func BenchmarkGenerateProof_1000(b *testing.B) {
	data := make([][]byte, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = []byte{byte(i % 256), byte(i / 256)}
	}
	tree := BuildMerkleTree(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateProof(tree, 500)
	}
}

func BenchmarkVerifyProof_1000(b *testing.B) {
	data := make([][]byte, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = []byte{byte(i % 256), byte(i / 256)}
	}
	tree := BuildMerkleTree(data)
	proof := GenerateProof(tree, 500)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		VerifyProof(data[500], proof, tree.Root)
	}
}

// ========================================
// HELPER FUNCTIONS FOR TESTS
// ========================================

func hashToHex(h []byte) string {
	return hex.EncodeToString(h)
}
