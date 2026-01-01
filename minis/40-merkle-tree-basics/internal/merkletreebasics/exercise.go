//go:build !solution && !reference


// Package exercise contains hands-on exercises for understanding Merkle trees.
//
// LEARNING OBJECTIVES:
// - Build a Merkle tree from data blocks
// - Generate cryptographic proofs of inclusion
// - Verify proofs without accessing the full dataset
// - Understand blockchain-style data verification

package merkletreebasics

// TODO: Import required packages
// You'll need:
// - "crypto/sha256" for hashing nodes
//
// import (
//     "crypto/sha256"
// )

// ============================================================================
// MERKLE TREES: Understanding Cryptographic Data Structures
// ============================================================================
//
// What is a Merkle Tree?
// - Binary tree of hashes (each parent is hash of its children)
// - Leaf nodes: Hashes of actual data
// - Internal nodes: Hashes of concatenated child hashes
// - Root hash: Single hash representing entire dataset
//
// Invented by Ralph Merkle (1979), used in:
// - Bitcoin: Verify transactions without downloading entire block
// - Git: Efficiently detect changes in file trees
// - IPFS: Content-addressable file system
// - Certificate Transparency: Verify log consistency
//
// Key Properties:
// 1. Tamper evidence: Changing any data changes the root hash
// 2. Efficient proofs: Prove inclusion with log(n) hashes
// 3. Space efficient: Store root, verify data with small proof
//
// Example Tree (4 data blocks):
//
//           Root (R)
//          /        \
//      H(AB)        H(CD)
//      /   \        /   \
//    H(A) H(B)   H(C) H(D)
//     |    |      |    |
//     A    B      C    D     <- Original data
//
// Where:
// - H(A) = SHA256(A)
// - H(AB) = SHA256(H(A) + H(B))  (+ means concatenate)
// - Root = SHA256(H(AB) + H(CD))
//
// Merkle Proof (proving B is in tree):
// - Need: H(A), H(CD)
// - Verify: Compute H(B), then H(AB) = H(H(A) + H(B)), then R = H(H(AB) + H(CD))
// - Compare computed R with known root
// - Only 2 hashes needed instead of all 4 data blocks!
//
// Memory Considerations:
// - Each hash: 32 bytes (SHA-256)
// - Tree with n leaves: n-1 internal nodes (total 2n-1 nodes)
// - Proof size: log2(n) * 32 bytes
// - Example: 1 million items → ~20 hashes (640 bytes) for proof
//
// ============================================================================

// MerkleTree represents a complete Merkle tree structure
type MerkleTree struct {
	Root   []byte     // The root hash representing the entire tree
	Leaves [][]byte   // All leaf hashes (hashes of original data)
	Levels [][][]byte // All levels of the tree (for proof generation)
}

// ProofNode represents one element in a Merkle proof
type ProofNode struct {
	Hash   []byte // Hash value of the sibling node
	IsLeft bool   // True if this hash should be on the left when combining
}

// MerkleProof represents a proof of inclusion for a specific data block
type MerkleProof struct {
	LeafIndex int         // Index of the data in the original dataset
	Siblings  []ProofNode // Sibling hashes needed to reconstruct the root
}

// BuildMerkleTree constructs a Merkle tree from data blocks.
//
// TODO: Implement this function
//
// This is the core algorithm for building the tree bottom-up.
//
// Algorithm:
// 1. Hash each data block to create leaf nodes (Level 0)
// 2. Build parent level by pairing and hashing children
// 3. If odd number of nodes, duplicate the last one
// 4. Repeat until only one node remains (the root)
// 5. Store all levels for later proof generation
//
// Step-by-step example with data = ["A", "B", "C"]:
//
// Step 1: Create leaves by hashing each data block
//   Level 0: [H(A), H(B), H(C)]
//
// Step 2: Pair up and hash (odd count, so duplicate H(C))
//   Level 0: [H(A), H(B), H(C), H(C)]  <- duplicated
//   Pairs: (H(A), H(B)) and (H(C), H(C))
//   Level 1: [H(H(A)+H(B)), H(H(C)+H(C))]
//
// Step 3: Pair up and hash again
//   Pairs: (H(H(A)+H(B)), H(H(C)+H(C)))
//   Level 2: [H(H(H(A)+H(B))+H(H(C)+H(C)))]  <- Root!
//
// Memory allocation:
// - Leaves slice: n * 32 bytes (n = number of data blocks)
// - Levels slice: log2(n) levels, each with varying number of hashes
// - Total: O(n) space (approximately 2n hashes)
//
// Edge cases:
// - Empty data: Return tree with hash of empty bytes as root
// - Single block: Hash of that block IS the root (no pairing needed)
// - Power of 2 blocks: No duplication needed at any level
// - Non-power of 2: Some duplication at each level
//
// Why duplicate odd nodes?
// - Bitcoin-style approach (used in Bitcoin, Ethereum)
// - Ensures tree is always balanced
// - Alternative: Promote odd node to next level (used in some systems)
//
// func BuildMerkleTree(data [][]byte) *MerkleTree {
//     // Handle empty data case
//     if len(data) == 0 {
//         emptyHash := hash([]byte{})
//         return &MerkleTree{
//             Root:   emptyHash,
//             Leaves: [][]byte{},
//             Levels: [][][]byte{{emptyHash}},
//         }
//     }
//
//     // Step 1: Create leaf nodes by hashing each data block
//     leaves := make([][]byte, len(data))
//     for i, d := range data {
//         leaves[i] = hash(d)  // hash() is a helper function (provided below)
//     }
//
//     // Step 2: Initialize levels with the leaf level
//     levels := [][][]byte{leaves}
//     currentLevel := leaves
//
//     // Step 3: Build tree level by level until we reach the root
//     for len(currentLevel) > 1 {
//         nextLevel := [][]byte{}
//
//         // If odd number of nodes, duplicate the last one
//         // This is the Bitcoin approach for balanced trees
//         if len(currentLevel)%2 == 1 {
//             currentLevel = append(currentLevel, currentLevel[len(currentLevel)-1])
//         }
//
//         // Pair nodes and hash them to create parent nodes
//         for i := 0; i < len(currentLevel); i += 2 {
//             left := currentLevel[i]
//             right := currentLevel[i+1]
//
//             // Parent hash = hash(left || right)
//             // Concatenate the two hashes and hash the result
//             combined := append(left, right...)
//             parentHash := hash(combined)
//             nextLevel = append(nextLevel, parentHash)
//         }
//
//         // Store this level and move up the tree
//         levels = append(levels, nextLevel)
//         currentLevel = nextLevel
//     }
//
//     return &MerkleTree{
//         Root:   currentLevel[0],  // Last remaining node is the root
//         Leaves: leaves,
//         Levels: levels,
//     }
// }

// GenerateProof creates a proof of inclusion for the data at the given index.
//
// TODO: Implement this function
//
// This generates the minimum set of hashes needed to verify data is in tree.
//
// Algorithm:
// 1. Start at the leaf index
// 2. For each level from bottom to top:
//    a. Find sibling index (if even, sibling is +1; if odd, sibling is -1)
//    b. Add sibling hash and its position (left or right) to proof
//    c. Move to parent index (divide by 2)
// 3. Return the proof
//
// Example: Tree with 4 leaves [A, B, C, D], proving index 1 (B):
//
//           Root
//          /    \
//      H(AB)    H(CD)
//      /  \      /  \
//   H(A) H(B) H(C) H(D)
//    0    1    2    3     <- Leaf indices
//
// Proving B (index 1):
//   Level 0: Index 1, sibling is index 0 (H(A)), IsLeft=true
//   Level 1: Index 0 (parent of 0,1), sibling is index 1 (H(CD)), IsLeft=false
//   Proof: [{H(A), true}, {H(CD), false}]
//
// How to use this proof:
//   Start with H(B)
//   Combine with H(A) on left: H(H(A) + H(B))
//   Combine with H(CD) on right: H(H(AB) + H(CD))
//   Result should equal Root
//
// Memory:
// - Proof size: log2(n) * (32 bytes + 1 bool) ≈ log2(n) * 33 bytes
// - Example: 1 million items → ~20 proof nodes ≈ 660 bytes
//
// Edge cases:
// - Invalid index: Return nil
// - Single leaf: Empty proof (data hash == root, no siblings needed)
// - Index 0: All siblings are on the right
// - Last index: All siblings are on the left
//
// func GenerateProof(tree *MerkleTree, index int) *MerkleProof {
//     // Validate index bounds
//     if index < 0 || index >= len(tree.Leaves) {
//         return nil
//     }
//
//     // Initialize proof structure
//     proof := &MerkleProof{
//         LeafIndex: index,
//         Siblings:  []ProofNode{},
//     }
//
//     currentIndex := index
//
//     // Traverse from leaf to root, collecting siblings
//     for level := 0; level < len(tree.Levels)-1; level++ {
//         currentLevelNodes := tree.Levels[level]
//
//         // Duplicate last node if odd count (same as BuildMerkleTree)
//         if len(currentLevelNodes)%2 == 1 {
//             currentLevelNodes = append(currentLevelNodes, currentLevelNodes[len(currentLevelNodes)-1])
//         }
//
//         // Find sibling index
//         // If current index is even (left child), sibling is index+1 (right)
//         // If current index is odd (right child), sibling is index-1 (left)
//         var siblingIndex int
//         if currentIndex%2 == 0 {
//             siblingIndex = currentIndex + 1
//         } else {
//             siblingIndex = currentIndex - 1
//         }
//
//         // Add sibling to proof
//         siblingHash := currentLevelNodes[siblingIndex]
//         isLeft := siblingIndex < currentIndex  // Sibling is left if its index is smaller
//
//         proof.Siblings = append(proof.Siblings, ProofNode{
//             Hash:   siblingHash,
//             IsLeft: isLeft,
//         })
//
//         // Move to parent index for next level
//         // Parent index is always currentIndex / 2 (integer division)
//         currentIndex = currentIndex / 2
//     }
//
//     return proof
// }

// VerifyProof verifies that data is in the tree by checking the proof against the root.
//
// TODO: Implement this function
//
// This is how light clients verify data without downloading the entire tree.
//
// Algorithm:
// 1. Compute hash of the data (this is the leaf hash)
// 2. For each sibling in proof:
//    a. If sibling is left: hash = H(sibling + current)
//    b. If sibling is right: hash = H(current + sibling)
//    c. Update current hash
// 3. Compare final computed hash with expected root
// 4. Return true if they match, false otherwise
//
// Example verification (proving B is in tree):
//   Data: "B"
//   Proof: [{H(A), true}, {H(CD), false}]
//   Root: H(H(AB) + H(CD))
//
// Step by step:
//   1. current = H(B)
//   2. Sibling 1: H(A), IsLeft=true
//      current = H(H(A) + H(B))  [sibling on left]
//   3. Sibling 2: H(CD), IsLeft=false
//      current = H(H(AB) + H(CD))  [sibling on right]
//   4. current == root? → true (verified!)
//
// Security critical:
// - Order matters! H(A+B) ≠ H(B+A)
// - Must respect IsLeft flag to prevent second-preimage attacks
// - Attacker cannot forge proof without knowing hash preimages
//
// Why this works:
// - We're reconstructing the path from leaf to root
// - Each sibling provides the "missing piece" to compute parent
// - If data was modified, the hash chain breaks
// - If proof is for different data, final hash won't match root
//
// func VerifyProof(data []byte, proof *MerkleProof, root []byte) bool {
//     // Handle nil proof case
//     if proof == nil {
//         return false
//     }
//
//     // Step 1: Start with hash of the data
//     current := hash(data)
//
//     // Step 2: Iteratively combine with siblings to move up the tree
//     for _, sibling := range proof.Siblings {
//         if sibling.IsLeft {
//             // Sibling goes on the left: hash(sibling + current)
//             combined := append(sibling.Hash, current...)
//             current = hash(combined)
//         } else {
//             // Sibling goes on the right: hash(current + sibling)
//             combined := append(current, sibling.Hash...)
//             current = hash(combined)
//         }
//     }
//
//     // Step 3: Compare computed root with expected root
//     return hashesEqual(current, root)
// }

// GetMerkleRoot is a convenience function that builds a tree and returns the root.
//
// TODO: Implement this function
//
// This is useful when you only care about the root hash (e.g., block header).
//
// Use case:
// - Block headers store only the Merkle root (not full tree)
// - Miners compute root to include in block header
// - Light clients can verify transactions using root + proof
//
// func GetMerkleRoot(data [][]byte) []byte {
//     tree := BuildMerkleTree(data)
//     return tree.Root
// }

// ========================================
// UTILITY FUNCTIONS (PROVIDED)
// ========================================

// hash computes SHA-256 hash of data
func hash(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

// hashesEqual compares two hashes for equality
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

// After implementing all functions:
// - Run: go test ./minis/40-merkle-tree-basics/exercise/...
// - Test solution: go test -tags solution ./minis/40-merkle-tree-basics/exercise/...
// - Compare with solution.go to see detailed implementations
//
// Key Concepts Learned:
// - Merkle trees: Binary hash trees for efficient verification
// - Tree construction: Bottom-up hashing from leaves to root
// - Proof generation: Collect sibling hashes along path to root
// - Proof verification: Reconstruct root from data + proof
// - Space efficiency: log(n) proof size for n items
//
// Real-World Applications:
// - Bitcoin: SPV (Simplified Payment Verification) wallets
// - Ethereum: State trees, transaction trees, receipt trees
// - Git: Efficiently detect changed files in repositories
// - IPFS: Content-addressable distributed file system
// - Certificate Transparency: Verify certificate logs
// - ZFS: File system integrity verification
//
// Comparison to alternatives:
// - Hash list: O(n) proof size (need all hashes)
// - Merkle tree: O(log n) proof size (exponentially better)
// - Example: 1M items → Hash list: 32MB, Merkle tree: 660 bytes
