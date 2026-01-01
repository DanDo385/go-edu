package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
	"time"
)

// Block represents a single block in the blockchain
type Block struct {
	Index     int    // Position in the chain
	Timestamp int64  // When the block was created (Unix timestamp)
	Data      string // The actual content (transactions, messages, etc.)
	PrevHash  string // Hash of the previous block (creates the chain)
	Nonce     int    // The proof of work nonce
	Hash      string // Hash of this block
}

// Blockchain represents the entire chain of blocks
type Blockchain struct {
	Blocks     []*Block
	Difficulty int
}

// CalculateHash computes the SHA-256 hash of the block
func (b *Block) CalculateHash() string {
	record := fmt.Sprintf("%d%d%s%s%d",
		b.Index,
		b.Timestamp,
		b.Data,
		b.PrevHash,
		b.Nonce,
	)
	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}

// Mine performs proof of work to find a valid nonce
// Returns the number of attempts it took
func (b *Block) Mine(difficulty int) int {
	target := strings.Repeat("0", difficulty)
	attempts := 0

	fmt.Printf("\n⛏️  Mining block %d (difficulty %d)...\n", b.Index, difficulty)
	startTime := time.Now()

	for {
		b.Hash = b.CalculateHash()
		attempts++

		// Check if we found a valid hash
		if strings.HasPrefix(b.Hash, target) {
			duration := time.Since(startTime)
			hashRate := float64(attempts) / duration.Seconds()

			fmt.Printf("✅ Block mined!\n")
			fmt.Printf("   Nonce: %d\n", b.Nonce)
			fmt.Printf("   Hash: %s\n", b.Hash)
			fmt.Printf("   Attempts: %s\n", formatNumber(attempts))
			fmt.Printf("   Time: %s\n", duration.Round(time.Millisecond))
			fmt.Printf("   Hash Rate: %s H/s\n\n", formatHashRate(hashRate))

			return attempts
		}

		b.Nonce++

		// If we've exhausted the nonce space, update timestamp
		if b.Nonce == math.MaxInt32 {
			b.Timestamp = time.Now().Unix()
			b.Nonce = 0
		}

		// Show progress every 100,000 attempts
		if attempts%100000 == 0 {
			elapsed := time.Since(startTime)
			currentHashRate := float64(attempts) / elapsed.Seconds()
			fmt.Printf("   ... %s attempts (%.0f H/s)\n", formatNumber(attempts), currentHashRate)
		}
	}
}

// NewBlockchain creates a new blockchain with a genesis block
func NewBlockchain(difficulty int) *Blockchain {
	genesis := &Block{
		Index:     0,
		Timestamp: time.Now().Unix(),
		Data:      "Genesis Block",
		PrevHash:  "0",
		Nonce:     0,
	}

	// Mine the genesis block
	genesis.Mine(difficulty)

	return &Blockchain{
		Blocks:     []*Block{genesis},
		Difficulty: difficulty,
	}
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(data string) int {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]

	newBlock := &Block{
		Index:     prevBlock.Index + 1,
		Timestamp: time.Now().Unix(),
		Data:      data,
		PrevHash:  prevBlock.Hash,
		Nonce:     0,
	}

	attempts := newBlock.Mine(bc.Difficulty)
	bc.Blocks = append(bc.Blocks, newBlock)

	return attempts
}

// IsValid checks if the blockchain is valid
func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		current := bc.Blocks[i]
		previous := bc.Blocks[i-1]

		// Check if the stored hash is correct
		if current.Hash != current.CalculateHash() {
			fmt.Printf("❌ Block %d: Hash mismatch\n", i)
			return false
		}

		// Check if the block links to the previous block
		if current.PrevHash != previous.Hash {
			fmt.Printf("❌ Block %d: Previous hash mismatch\n", i)
			return false
		}

		// Check if the proof of work is valid
		target := strings.Repeat("0", bc.Difficulty)
		if !strings.HasPrefix(current.Hash, target) {
			fmt.Printf("❌ Block %d: Invalid proof of work\n", i)
			return false
		}
	}

	return true
}

// AdjustDifficulty adjusts mining difficulty based on block times
func (bc *Blockchain) AdjustDifficulty(targetBlockTime int64) {
	if len(bc.Blocks) < 11 {
		return // Need at least 11 blocks for adjustment
	}

	// Calculate actual time for last 10 blocks
	actualTime := bc.Blocks[len(bc.Blocks)-1].Timestamp -
		bc.Blocks[len(bc.Blocks)-10].Timestamp
	expectedTime := targetBlockTime * 9 // 9 intervals between 10 blocks

	fmt.Printf("\n📊 Difficulty Adjustment Check:\n")
	fmt.Printf("   Expected time for 10 blocks: %d seconds\n", expectedTime)
	fmt.Printf("   Actual time: %d seconds\n", actualTime)
	fmt.Printf("   Current difficulty: %d\n", bc.Difficulty)

	// Adjust difficulty if significantly off target
	if actualTime < expectedTime/2 {
		bc.Difficulty++
		fmt.Printf("   ⬆️  Increasing difficulty to %d (mining too fast)\n\n", bc.Difficulty)
	} else if actualTime > expectedTime*2 {
		if bc.Difficulty > 1 {
			bc.Difficulty--
			fmt.Printf("   ⬇️  Decreasing difficulty to %d (mining too slow)\n\n", bc.Difficulty)
		}
	} else {
		fmt.Printf("   ✅ Difficulty unchanged (mining rate is good)\n\n")
	}
}

// PrintChain displays the entire blockchain
func (bc *Blockchain) PrintChain() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("BLOCKCHAIN SUMMARY")
	fmt.Println(strings.Repeat("=", 80))

	for _, block := range bc.Blocks {
		fmt.Printf("\nBlock %d:\n", block.Index)
		fmt.Printf("  Timestamp: %s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05"))
		fmt.Printf("  Data: %s\n", block.Data)
		fmt.Printf("  PrevHash: %s\n", block.PrevHash)
		fmt.Printf("  Nonce: %d\n", block.Nonce)
		fmt.Printf("  Hash: %s\n", block.Hash)

		// Highlight the leading zeros
		target := strings.Repeat("0", bc.Difficulty)
		if strings.HasPrefix(block.Hash, target) {
			fmt.Printf("  ✅ Valid PoW (%d leading zeros)\n", bc.Difficulty)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
}

// DemonstrateAttack shows what happens when someone tries to tamper with the blockchain
func (bc *Blockchain) DemonstrateAttack() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("TAMPERING DEMONSTRATION")
	fmt.Println(strings.Repeat("=", 80))

	if len(bc.Blocks) < 3 {
		fmt.Println("Need at least 3 blocks to demonstrate tampering")
		return
	}

	// Try to modify an old block
	targetBlock := bc.Blocks[1]
	fmt.Printf("\n🔧 Attempting to modify Block %d...\n", targetBlock.Index)
	fmt.Printf("   Original data: '%s'\n", targetBlock.Data)

	originalData := targetBlock.Data
	targetBlock.Data = "HACKED DATA"

	fmt.Printf("   Modified data: '%s'\n", targetBlock.Data)
	fmt.Printf("\n❓ Is blockchain still valid?\n")

	if bc.IsValid() {
		fmt.Println("✅ Blockchain is still valid (this shouldn't happen!)")
	} else {
		fmt.Println("❌ Blockchain is now INVALID!")
		fmt.Println("\n💡 Why? Because:")
		fmt.Println("   1. Block 1's hash changed (due to modified data)")
		fmt.Println("   2. Block 2's PrevHash still points to the OLD hash of Block 1")
		fmt.Println("   3. The chain is broken!")
		fmt.Println("\n🛡️  To successfully tamper, attacker would need to:")
		fmt.Println("   1. Re-mine Block 1 (find new valid nonce)")
		fmt.Println("   2. Re-mine Block 2 (with updated PrevHash)")
		fmt.Println("   3. Re-mine ALL subsequent blocks")
		fmt.Println("   4. Do this faster than honest miners add new blocks")
		fmt.Println("   5. This requires >50% of network hash rate!")
	}

	// Restore original data
	targetBlock.Data = originalData
	fmt.Println("\n🔄 Restoring original data...\n")
}

// formatNumber formats a number with commas
func formatNumber(n int) string {
	str := fmt.Sprintf("%d", n)
	if len(str) <= 3 {
		return str
	}

	var result []byte
	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(digit))
	}
	return string(result)
}

// formatHashRate formats hash rate with appropriate unit
func formatHashRate(hashRate float64) string {
	if hashRate >= 1_000_000_000 {
		return fmt.Sprintf("%.2f GH", hashRate/1_000_000_000)
	} else if hashRate >= 1_000_000 {
		return fmt.Sprintf("%.2f MH", hashRate/1_000_000)
	} else if hashRate >= 1_000 {
		return fmt.Sprintf("%.2f KH", hashRate/1_000)
	}
	return fmt.Sprintf("%.2f", hashRate)
}

// DifficultyDemo demonstrates how difficulty affects mining time
func DifficultyDemo() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("DIFFICULTY DEMONSTRATION")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("\nThis demo shows how difficulty exponentially affects mining time.")
	fmt.Println("Each additional zero roughly increases mining time by 16x.\n")

	for difficulty := 1; difficulty <= 5; difficulty++ {
		fmt.Printf("Testing difficulty %d (%d leading zeros):\n", difficulty, difficulty)

		block := &Block{
			Index:     1,
			Timestamp: time.Now().Unix(),
			Data:      fmt.Sprintf("Test block for difficulty %d", difficulty),
			PrevHash:  "0000000000000000000000000000000000000000000000000000000000000000",
			Nonce:     0,
		}

		startTime := time.Now()
		attempts := block.Mine(difficulty)
		duration := time.Since(startTime)

		expectedAttempts := int(math.Pow(16, float64(difficulty)))
		fmt.Printf("   Expected attempts: ~%s\n", formatNumber(expectedAttempts))
		fmt.Printf("   Actual attempts: %s\n", formatNumber(attempts))
		fmt.Printf("   Ratio: %.2fx\n\n", float64(attempts)/float64(expectedAttempts))

		// Stop if mining is taking too long
		if duration > 10*time.Second {
			fmt.Println("⚠️  Stopping demo - difficulty is getting too high for quick demonstration")
			break
		}
	}
}

// ChainDemo demonstrates building a blockchain with difficulty adjustment
func ChainDemo() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("BLOCKCHAIN BUILDING DEMONSTRATION")
	fmt.Println(strings.Repeat("=", 80))

	difficulty := 4
	targetBlockTime := int64(2) // 2 seconds per block
	bc := NewBlockchain(difficulty)

	fmt.Println("\n📦 Building blockchain with automatic difficulty adjustment...")
	fmt.Printf("   Target block time: %d seconds\n", targetBlockTime)
	fmt.Printf("   Initial difficulty: %d\n\n", difficulty)

	// Mine several blocks
	transactions := []string{
		"Alice -> Bob: 50 BTC",
		"Bob -> Charlie: 25 BTC",
		"Charlie -> Alice: 10 BTC",
		"Alice -> Dave: 5 BTC",
		"Dave -> Bob: 2 BTC",
		"Bob -> Alice: 15 BTC",
		"Charlie -> Dave: 8 BTC",
		"Dave -> Charlie: 3 BTC",
		"Alice -> Bob: 12 BTC",
		"Bob -> Charlie: 7 BTC",
	}

	var totalAttempts int
	startTime := time.Now()

	for i, tx := range transactions {
		attempts := bc.AddBlock(tx)
		totalAttempts += attempts

		// Adjust difficulty every 10 blocks
		if (i+2)%10 == 0 { // +2 because index 0 is genesis
			bc.AdjustDifficulty(targetBlockTime)
		}
	}

	totalTime := time.Since(startTime)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("MINING STATISTICS")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("\nTotal blocks: %d\n", len(bc.Blocks))
	fmt.Printf("Total time: %s\n", totalTime.Round(time.Millisecond))
	fmt.Printf("Total attempts: %s\n", formatNumber(totalAttempts))
	fmt.Printf("Average time per block: %s\n", (totalTime/time.Duration(len(bc.Blocks))).Round(time.Millisecond))
	fmt.Printf("Average hash rate: %s H/s\n", formatHashRate(float64(totalAttempts)/totalTime.Seconds()))
	fmt.Printf("Final difficulty: %d\n", bc.Difficulty)

	// Validate the blockchain
	fmt.Printf("\n🔍 Validating blockchain...\n")
	if bc.IsValid() {
		fmt.Println("✅ Blockchain is valid!")
	} else {
		fmt.Println("❌ Blockchain is invalid!")
	}

	// Print the chain
	bc.PrintChain()

	// Demonstrate tampering
	bc.DemonstrateAttack()
}

// HashRateDemo demonstrates hash rate calculation
func HashRateDemo() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("HASH RATE DEMONSTRATION")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("\nMeasuring your computer's hash rate...\n")

	block := &Block{
		Index:     1,
		Timestamp: time.Now().Unix(),
		Data:      "Hash rate test",
		PrevHash:  "0000000000000000000000000000000000000000000000000000000000000000",
		Nonce:     0,
	}

	// Compute many hashes to get accurate measurement
	iterations := 1_000_000
	startTime := time.Now()

	for i := 0; i < iterations; i++ {
		block.Nonce = i
		_ = block.CalculateHash()
	}

	duration := time.Since(startTime)
	hashRate := float64(iterations) / duration.Seconds()

	fmt.Printf("Computed %s hashes in %s\n", formatNumber(iterations), duration.Round(time.Millisecond))
	fmt.Printf("Hash rate: %s H/s\n\n", formatHashRate(hashRate))

	// Estimate mining times for different difficulties
	fmt.Println("Estimated mining times at this hash rate:")
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("%-15s %-20s %-20s\n", "Difficulty", "Expected Attempts", "Expected Time")
	fmt.Println(strings.Repeat("-", 60))

	for difficulty := 1; difficulty <= 8; difficulty++ {
		expectedAttempts := math.Pow(16, float64(difficulty))
		expectedTime := expectedAttempts / hashRate

		fmt.Printf("%-15d %-20s %-20s\n",
			difficulty,
			formatNumber(int(expectedAttempts)),
			formatDuration(expectedTime),
		)

		// Stop if time gets too long
		if expectedTime > 3600 {
			fmt.Println("\n⚠️  Higher difficulties would take hours to days...")
			break
		}
	}
	fmt.Println()
}

// formatDuration formats seconds into human-readable duration
func formatDuration(seconds float64) string {
	if seconds < 1 {
		return fmt.Sprintf("%.0f ms", seconds*1000)
	} else if seconds < 60 {
		return fmt.Sprintf("%.1f sec", seconds)
	} else if seconds < 3600 {
		return fmt.Sprintf("%.1f min", seconds/60)
	} else if seconds < 86400 {
		return fmt.Sprintf("%.1f hours", seconds/3600)
	}
	return fmt.Sprintf("%.1f days", seconds/86400)
}

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                       PROOF OF WORK DEMONSTRATION                          ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════════════╝")

	// Run demonstrations
	HashRateDemo()
	DifficultyDemo()
	ChainDemo()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("DEMONSTRATION COMPLETE")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("\n💡 Key Takeaways:")
	fmt.Println("   • Each additional zero in difficulty increases mining time by ~16x")
	fmt.Println("   • Hash rate determines mining speed (more hash power = faster mining)")
	fmt.Println("   • Blockchain links blocks cryptographically (tampering breaks the chain)")
	fmt.Println("   • Difficulty adjusts to maintain consistent block times")
	fmt.Println("   • Proof of Work makes it expensive to attack but easy to verify")
	fmt.Println()
}
