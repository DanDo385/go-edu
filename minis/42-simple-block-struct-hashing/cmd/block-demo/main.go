package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/example/go-10x-minis/minis/42-simple-block-struct-hashing/exercise"
)

func main() {
	fmt.Println("=== Blockchain Block Demo ===\n")

	// Create genesis block
	fmt.Println("Creating genesis block...")
	genesis := exercise.NewGenesisBlock()
	printBlock(genesis, 0)

	// Build a chain
	chain := []*exercise.Block{genesis}

	// Add block 1
	fmt.Println("\n--- Adding Block 1 ---")
	block1 := exercise.NewBlock(chain[len(chain)-1], []string{
		"Alice sends 10 BTC to Bob",
		"Bob sends 5 BTC to Carol",
	})
	chain = append(chain, block1)
	printBlock(block1, 1)

	// Add block 2
	fmt.Println("\n--- Adding Block 2 ---")
	block2 := exercise.NewBlock(chain[len(chain)-1], []string{
		"Carol sends 3 BTC to Dave",
		"Dave sends 1 BTC to Eve",
		"Eve sends 2 BTC to Frank",
	})
	chain = append(chain, block2)
	printBlock(block2, 2)

	// Add block 3
	fmt.Println("\n--- Adding Block 3 ---")
	block3 := exercise.NewBlock(chain[len(chain)-1], []string{
		"Frank sends 1 BTC to Alice",
	})
	chain = append(chain, block3)
	printBlock(block3, 3)

	// Validate the chain
	fmt.Println("\n=== Validating Blockchain ===")
	if err := exercise.ValidateChain(chain); err != nil {
		fmt.Printf("❌ Chain validation FAILED: %v\n", err)
	} else {
		fmt.Println("✅ Chain validation PASSED - all blocks are valid")
	}

	// Show hash linking
	fmt.Println("\n=== Hash Linking Demonstration ===")
	for i := 1; i < len(chain); i++ {
		fmt.Printf("Block %d PrevHash: %s\n", i, chain[i].Header.PrevHash[:16]+"...")
		fmt.Printf("Block %d Hash:     %s\n", i-1, chain[i-1].Hash[:16]+"...")
		if chain[i].Header.PrevHash == chain[i-1].Hash {
			fmt.Printf("✅ Block %d correctly links to Block %d\n\n", i, i-1)
		} else {
			fmt.Printf("❌ Block %d does NOT link to Block %d\n\n", i, i-1)
		}
	}

	// Demonstrate tampering detection
	fmt.Println("=== Tampering Detection Demo ===")
	fmt.Println("Attempting to tamper with Block 1...")

	// Create a tampered copy of the chain
	tamperedChain := make([]*exercise.Block, len(chain))
	for i, block := range chain {
		// Deep copy (simplified)
		copied := &exercise.Block{
			Header:       block.Header,
			Transactions: make([]string, len(block.Transactions)),
			Hash:         block.Hash,
		}
		copy(copied.Transactions, block.Transactions)
		tamperedChain[i] = copied
	}

	// Tamper with block 1's transaction
	fmt.Println("Original transaction:", tamperedChain[1].Transactions[0])
	tamperedChain[1].Transactions[0] = "Alice sends 1000 BTC to Bob"
	fmt.Println("Tampered transaction:", tamperedChain[1].Transactions[0])

	// Try to validate tampered chain
	fmt.Println("\nValidating tampered chain...")
	if err := exercise.ValidateChain(tamperedChain); err != nil {
		fmt.Printf("✅ Tampering detected: %v\n", err)
		fmt.Println("   The blockchain successfully detected the fraudulent modification!")
	} else {
		fmt.Println("❌ WARNING: Tampering was NOT detected (this shouldn't happen)")
	}

	// Show block statistics
	fmt.Println("\n=== Blockchain Statistics ===")
	fmt.Printf("Total blocks: %d\n", len(chain))
	fmt.Printf("Total transactions: %d\n", countTransactions(chain))
	fmt.Printf("Chain created over: %d seconds\n", chain[len(chain)-1].Header.Timestamp-chain[0].Header.Timestamp)
	fmt.Printf("Average block size: %d bytes\n", averageBlockSize(chain))

	// Display full chain summary
	fmt.Println("\n=== Complete Chain Summary ===")
	for i, block := range chain {
		fmt.Printf("\nBlock %d:\n", i)
		fmt.Printf("  Hash:         %s\n", block.Hash[:32]+"...")
		fmt.Printf("  Previous:     %s\n", block.Header.PrevHash[:32]+"...")
		fmt.Printf("  Merkle Root:  %s\n", truncateHash(block.Header.MerkleRoot))
		fmt.Printf("  Timestamp:    %s\n", formatTimestamp(block.Header.Timestamp))
		fmt.Printf("  Transactions: %d\n", len(block.Transactions))
	}

	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("Key takeaways:")
	fmt.Println("  1. Each block contains a hash of the previous block (hash linking)")
	fmt.Println("  2. Changing any block's data changes its hash")
	fmt.Println("  3. This breaks the chain and is immediately detectable")
	fmt.Println("  4. The blockchain is tamper-evident (changes are visible)")
	fmt.Println("  5. Merkle roots ensure transaction integrity within blocks")
}

func printBlock(block *exercise.Block, blockNum int) {
	fmt.Printf("Block %d created:\n", blockNum)
	fmt.Printf("  Index:        %d\n", block.Header.Index)
	fmt.Printf("  Timestamp:    %s\n", formatTimestamp(block.Header.Timestamp))
	fmt.Printf("  PrevHash:     %s\n", truncateHash(block.Header.PrevHash))
	fmt.Printf("  MerkleRoot:   %s\n", truncateHash(block.Header.MerkleRoot))
	fmt.Printf("  Hash:         %s\n", truncateHash(block.Hash))
	fmt.Printf("  Transactions: %d\n", len(block.Transactions))
	for i, tx := range block.Transactions {
		fmt.Printf("    [%d] %s\n", i, tx)
	}
}

func truncateHash(hash string) string {
	if len(hash) > 16 {
		return hash[:16] + "..."
	}
	return hash
}

func formatTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}

func countTransactions(chain []*exercise.Block) int {
	total := 0
	for _, block := range chain {
		total += len(block.Transactions)
	}
	return total
}

func averageBlockSize(chain []*exercise.Block) int {
	if len(chain) == 0 {
		return 0
	}

	totalSize := 0
	for _, block := range chain {
		totalSize += len(block.Serialize())
	}

	return totalSize / len(chain)
}

// Demonstration of how hash changes propagate
func demonstrateHashAvalanche() {
	fmt.Println("\n=== Hash Avalanche Effect ===")
	fmt.Println("Showing how small changes create completely different hashes:\n")

	// Create two nearly identical blocks
	genesis := exercise.NewGenesisBlock()

	block1 := exercise.NewBlock(genesis, []string{"Alice sends 10 BTC to Bob"})
	block2 := exercise.NewBlock(genesis, []string{"Alice sends 11 BTC to Bob"}) // Only 1 digit different

	fmt.Println("Transaction 1:", block1.Transactions[0])
	fmt.Println("Hash 1:       ", block1.Hash[:32]+"...")
	fmt.Println()
	fmt.Println("Transaction 2:", block2.Transactions[0])
	fmt.Println("Hash 2:       ", block2.Hash[:32]+"...")
	fmt.Println()

	// Count how many characters are different
	differentChars := 0
	for i := 0; i < len(block1.Hash) && i < len(block2.Hash); i++ {
		if block1.Hash[i] != block2.Hash[i] {
			differentChars++
		}
	}

	fmt.Printf("Characters different: %d out of %d (%.1f%%)\n",
		differentChars, len(block1.Hash),
		float64(differentChars)/float64(len(block1.Hash))*100)
	fmt.Println("This is the avalanche effect: tiny input change → completely different hash")
}

// Demonstrate the importance of proper serialization
func demonstrateSerializationImportance() {
	fmt.Println("\n=== Serialization Importance ===")

	genesis := exercise.NewGenesisBlock()
	block := exercise.NewBlock(genesis, []string{"Test transaction"})

	// Show serialized bytes (first 64 bytes)
	serialized := block.Serialize()
	fmt.Printf("Block serialized to %d bytes\n", len(serialized))
	fmt.Printf("First 64 bytes (hex): %x\n", serialized[:min(64, len(serialized))])

	// Compute hash multiple times to show determinism
	hash1 := block.ComputeHash()
	time.Sleep(10 * time.Millisecond) // Wait a bit
	hash2 := block.ComputeHash()

	if hash1 == hash2 {
		fmt.Println("✅ Hash is deterministic (same block → same hash)")
	} else {
		fmt.Println("❌ Hash is NOT deterministic (same block → different hash)")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Visualize the blockchain
func visualizeChain(chain []*exercise.Block) {
	fmt.Println("\n=== Blockchain Visualization ===\n")

	for i, block := range chain {
		// Draw block
		fmt.Println("┌" + strings.Repeat("─", 60) + "┐")
		fmt.Printf("│ Block %d%s│\n", i, strings.Repeat(" ", 53-len(fmt.Sprintf("Block %d", i))))
		fmt.Println("├" + strings.Repeat("─", 60) + "┤")
		fmt.Printf("│ Hash: %-51s │\n", block.Hash[:51])
		fmt.Printf("│ PrevHash: %-47s │\n", block.Header.PrevHash[:47])
		fmt.Printf("│ Transactions: %-43d │\n", len(block.Transactions))
		fmt.Println("└" + strings.Repeat("─", 60) + "┘")

		// Draw link to next block
		if i < len(chain)-1 {
			fmt.Println("    │")
			fmt.Println("    │ (hash link)")
			fmt.Println("    ↓")
		}
	}
}
