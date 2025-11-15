package exercise

import (
	"fmt"
	"math"
)

// Example demonstrating hash calculation
func ExampleCalculateBlockHash() {
	block := Block{
		Index:     1,
		Timestamp: 1609459200,
		Data:      "Alice -> Bob: 50 BTC",
		PrevHash:  "0000000000000000000000000000000000000000000000000000000000000000",
		Nonce:     12345,
	}

	hash := CalculateBlockHashSolution(block)
	fmt.Printf("Hash length: %d characters\n", len(hash))
	fmt.Printf("Hash starts with: %s...\n", hash[:8])
	// Output:
	// Hash length: 64 characters
	// Hash starts with: 8e7ea5db...
}

// Example demonstrating proof of work validation
func ExampleIsValidProof() {
	// Valid proof with 4 leading zeros
	validHash := "0000abc123def456789"
	fmt.Printf("Hash '%s' with difficulty 4: %v\n",
		validHash, IsValidProofSolution(validHash, 4))

	// Invalid proof (only 3 zeros)
	invalidHash := "000abc123def456789"
	fmt.Printf("Hash '%s' with difficulty 4: %v\n",
		invalidHash, IsValidProofSolution(invalidHash, 4))

	// Output:
	// Hash '0000abc123def456789' with difficulty 4: true
	// Hash '000abc123def456789' with difficulty 4: false
}

// Example demonstrating block mining
func ExampleMineBlock() {
	block := Block{
		Index:     1,
		Timestamp: 1609459200,
		Data:      "Genesis transaction",
		PrevHash:  "0",
		Nonce:     0,
	}

	difficulty := 3
	attempts := MineBlockSolution(&block, difficulty)

	fmt.Printf("Mined block in %d attempts\n", attempts)
	fmt.Printf("Final nonce: %d\n", block.Nonce)
	fmt.Printf("Hash starts with: %s...\n", block.Hash[:8])
	fmt.Printf("Valid proof: %v\n", IsValidProofSolution(block.Hash, difficulty))
	// Output varies due to randomness, but structure is:
	// Mined block in X attempts
	// Final nonce: Y
	// Hash starts with: 000...
	// Valid proof: true
}

// Example demonstrating blockchain validation
func ExampleValidateChain() {
	// Create a small valid chain
	chain := []Block{
		{
			Index:     0,
			Timestamp: 1000,
			Data:      "Genesis",
			PrevHash:  "0",
			Nonce:     0,
		},
	}
	MineBlockSolution(&chain[0], 3)

	// Add second block
	chain = append(chain, Block{
		Index:     1,
		Timestamp: 1010,
		Data:      "Block 1",
		PrevHash:  chain[0].Hash,
		Nonce:     0,
	})
	MineBlockSolution(&chain[1], 3)

	// Validate the chain
	isValid := ValidateChainSolution(chain, 3)
	fmt.Printf("Chain with %d blocks is valid: %v\n", len(chain), isValid)

	// Tamper with a block
	chain[0].Data = "TAMPERED"
	isValid = ValidateChainSolution(chain, 3)
	fmt.Printf("Chain after tampering is valid: %v\n", isValid)

	// Output:
	// Chain with 2 blocks is valid: true
	// Chain after tampering is valid: false
}

// Example demonstrating difficulty adjustment
func ExampleAdjustDifficulty() {
	// Create blocks that were mined too fast
	fastChain := []Block{
		{Timestamp: 1000},
		{Timestamp: 1002}, // 2 seconds each (too fast)
		{Timestamp: 1004},
		{Timestamp: 1006},
		{Timestamp: 1008},
		{Timestamp: 1010},
		{Timestamp: 1012},
		{Timestamp: 1014},
		{Timestamp: 1016},
		{Timestamp: 1018},
		{Timestamp: 1020}, // Total: 20 seconds for 10 blocks
	}

	targetBlockTime := int64(10) // Want 10 seconds per block
	currentDifficulty := 3

	newDifficulty := AdjustDifficultySolution(fastChain, targetBlockTime, currentDifficulty)
	fmt.Printf("Fast mining: difficulty %d -> %d\n", currentDifficulty, newDifficulty)

	// Create blocks that were mined at the right speed
	goodChain := []Block{
		{Timestamp: 1000},
		{Timestamp: 1010}, // 10 seconds each (just right)
		{Timestamp: 1020},
		{Timestamp: 1030},
		{Timestamp: 1040},
		{Timestamp: 1050},
		{Timestamp: 1060},
		{Timestamp: 1070},
		{Timestamp: 1080},
		{Timestamp: 1090},
		{Timestamp: 1100}, // Total: 100 seconds for 10 blocks
	}

	newDifficulty = AdjustDifficultySolution(goodChain, targetBlockTime, currentDifficulty)
	fmt.Printf("Good mining: difficulty %d -> %d\n", currentDifficulty, newDifficulty)

	// Output:
	// Fast mining: difficulty 3 -> 4
	// Good mining: difficulty 3 -> 3
}

// Example demonstrating mining probability calculation
func ExampleMiningProbability() {
	hashRate := 1000.0    // 1,000 hashes per second
	difficulty := 3       // Requires ~4,096 attempts on average
	timeSeconds := 10.0   // 10 seconds

	probability := MiningProbabilitySolution(hashRate, difficulty, timeSeconds)

	fmt.Printf("Hash rate: %.0f H/s\n", hashRate)
	fmt.Printf("Difficulty: %d (expected attempts: %.0f)\n",
		difficulty, math.Pow(16, float64(difficulty)))
	fmt.Printf("Time: %.0f seconds\n", timeSeconds)
	fmt.Printf("Probability of finding block: %.2f%%\n", probability*100)

	// Calculate expected time
	expectedAttempts := math.Pow(16, float64(difficulty))
	expectedTime := expectedAttempts / hashRate
	fmt.Printf("Expected time to find block: %.1f seconds\n", expectedTime)

	// Output:
	// Hash rate: 1000 H/s
	// Difficulty: 3 (expected attempts: 4096)
	// Time: 10 seconds
	// Probability of finding block: 91.79%
	// Expected time to find block: 4.1 seconds
}

// Example showing how difficulty scales exponentially
func Example_difficultyScaling() {
	hashRate := 10000.0 // 10,000 hashes per second

	fmt.Println("Difficulty scaling (each +1 multiplies time by ~16x):")
	fmt.Println("Difficulty | Expected Attempts | Expected Time")
	fmt.Println("-----------|-------------------|---------------")

	for difficulty := 1; difficulty <= 6; difficulty++ {
		expectedAttempts := math.Pow(16, float64(difficulty))
		expectedTime := expectedAttempts / hashRate

		var timeStr string
		if expectedTime < 1 {
			timeStr = fmt.Sprintf("%.0f ms", expectedTime*1000)
		} else if expectedTime < 60 {
			timeStr = fmt.Sprintf("%.1f sec", expectedTime)
		} else if expectedTime < 3600 {
			timeStr = fmt.Sprintf("%.1f min", expectedTime/60)
		} else {
			timeStr = fmt.Sprintf("%.1f hours", expectedTime/3600)
		}

		fmt.Printf("    %d      |   %10.0f      | %s\n",
			difficulty, expectedAttempts, timeStr)
	}

	// Output:
	// Difficulty scaling (each +1 multiplies time by ~16x):
	// Difficulty | Expected Attempts | Expected Time
	// -----------|-------------------|---------------
	//     1      |           16      | 2 ms
	//     2      |          256      | 26 ms
	//     3      |         4096      | 0.4 sec
	//     4      |        65536      | 6.6 sec
	//     5      |      1048576      | 1.7 min
	//     6      |     16777216      | 27.9 min
}

// Example showing the avalanche effect
func Example_avalancheEffect() {
	block1 := Block{
		Index:     1,
		Timestamp: 1609459200,
		Data:      "Hello, World!",
		PrevHash:  "0",
		Nonce:     0,
	}

	block2 := Block{
		Index:     1,
		Timestamp: 1609459200,
		Data:      "Hello, World?", // Changed ! to ?
		PrevHash:  "0",
		Nonce:     0,
	}

	hash1 := CalculateBlockHashSolution(block1)
	hash2 := CalculateBlockHashSolution(block2)

	fmt.Printf("Original:  %s\n", hash1[:32])
	fmt.Printf("Modified:  %s\n", hash2[:32])
	fmt.Printf("First 32 chars match: %v\n", hash1[:32] == hash2[:32])

	// Output varies, but hashes will be completely different:
	// Original:  <some hash>
	// Modified:  <completely different hash>
	// First 32 chars match: false
}

// Example demonstrating why blockchain is immutable
func Example_blockchainImmutability() {
	// Build a small chain
	chain := []Block{{
		Index: 0, Timestamp: 1000, Data: "Genesis", PrevHash: "0",
	}}
	MineBlockSolution(&chain[0], 2)

	for i := 1; i <= 3; i++ {
		block := Block{
			Index:     i,
			Timestamp: 1000 + int64(i*10),
			Data:      fmt.Sprintf("Block %d", i),
			PrevHash:  chain[i-1].Hash,
		}
		MineBlockSolution(&block, 2)
		chain = append(chain, block)
	}

	fmt.Printf("Original chain valid: %v\n", ValidateChainSolution(chain, 2))

	// Try to modify block 1
	originalData := chain[1].Data
	chain[1].Data = "HACKED DATA"

	fmt.Printf("After modifying block 1: %v\n", ValidateChainSolution(chain, 2))
	fmt.Println("Why? Block 1's hash changed, so block 2's PrevHash is now wrong")

	// Restore
	chain[1].Data = originalData
	fmt.Printf("After restoring: %v\n", ValidateChainSolution(chain, 2))

	// Output:
	// Original chain valid: true
	// After modifying block 1: false
	// Why? Block 1's hash changed, so block 2's PrevHash is now wrong
	// After restoring: true
}

// Example showing mining pool reward distribution
func ExampleCalculatePoolRewards() {
	miners := []Miner{
		{ID: "Alice", Shares: 100},
		{ID: "Bob", Shares: 150},
		{ID: "Charlie", Shares: 50},
	}

	blockReward := 12.5 // BTC

	rewards := CalculatePoolRewardsSolution(miners, blockReward)

	fmt.Println("Mining Pool Reward Distribution:")
	totalShares := 100 + 150 + 50
	for _, miner := range miners {
		percentage := float64(miner.Shares) / float64(totalShares) * 100
		fmt.Printf("%s: %.2f%% (%d shares) = %.2f BTC\n",
			miner.ID, percentage, miner.Shares, rewards[miner.ID])
	}

	// Output:
	// Mining Pool Reward Distribution:
	// Alice: 33.33% (100 shares) = 4.17 BTC
	// Bob: 50.00% (150 shares) = 6.25 BTC
	// Charlie: 16.67% (50 shares) = 2.08 BTC
}

// Example showing Merkle tree construction
func ExampleBuildMerkleRoot() {
	transactions := []string{
		"Alice -> Bob: 10 BTC",
		"Bob -> Charlie: 5 BTC",
		"Charlie -> Dave: 3 BTC",
		"Dave -> Alice: 2 BTC",
	}

	merkleRoot := BuildMerkleRootSolution(transactions)

	fmt.Printf("Transactions: %d\n", len(transactions))
	fmt.Printf("Merkle root: %s...\n", merkleRoot[:16])
	fmt.Printf("Root hash length: %d characters\n", len(merkleRoot))

	// Change one transaction
	transactions[0] = "Alice -> Bob: 11 BTC" // Changed amount
	newRoot := BuildMerkleRootSolution(transactions)

	fmt.Printf("Changed one transaction, new root: %s...\n", newRoot[:16])
	fmt.Printf("Roots match: %v\n", merkleRoot == newRoot)

	// Output will vary but show:
	// Transactions: 4
	// Merkle root: <some hash>...
	// Root hash length: 64 characters
	// Changed one transaction, new root: <different hash>...
	// Roots match: false
}

// Example showing expected mining time vs actual variance
func Example_miningVariance() {
	hashRate := 5000.0
	difficulty := 3
	expectedAttempts := math.Pow(16, float64(difficulty))
	expectedTime := expectedAttempts / hashRate

	fmt.Printf("Expected attempts: %.0f\n", expectedAttempts)
	fmt.Printf("Expected time: %.2f seconds\n", expectedTime)
	fmt.Println("\nProbability of finding block within different times:")

	times := []float64{0.5, 1.0, 2.0, 5.0, 10.0}
	for _, t := range times {
		prob := MiningProbabilitySolution(hashRate, difficulty, t*expectedTime)
		fmt.Printf("  %.1fx expected time (%.2fs): %.1f%% chance\n",
			t, t*expectedTime, prob*100)
	}

	// Output:
	// Expected attempts: 4096
	// Expected time: 0.82 seconds
	//
	// Probability of finding block within different times:
	//   0.5x expected time (0.41s): 39.3% chance
	//   1.0x expected time (0.82s): 63.2% chance
	//   2.0x expected time (1.64s): 86.5% chance
	//   5.0x expected time (4.10s): 99.3% chance
	//   10.0x expected time (8.19s): 100.0% chance
}

// Example showing network hash rate and difficulty relationship
func Example_networkHashRate() {
	targetBlockTime := 600.0 // 10 minutes in seconds
	difficulties := []int{1, 2, 3, 4, 5}

	fmt.Println("Network hash rate needed for 10-minute blocks:")
	fmt.Println("Difficulty | Expected Attempts | Required Hash Rate")
	fmt.Println("-----------|-------------------|-------------------")

	for _, diff := range difficulties {
		expectedAttempts := math.Pow(16, float64(diff))
		requiredHashRate := expectedAttempts / targetBlockTime

		var hashRateStr string
		if requiredHashRate < 1000 {
			hashRateStr = fmt.Sprintf("%.0f H/s", requiredHashRate)
		} else if requiredHashRate < 1000000 {
			hashRateStr = fmt.Sprintf("%.1f KH/s", requiredHashRate/1000)
		} else if requiredHashRate < 1000000000 {
			hashRateStr = fmt.Sprintf("%.1f MH/s", requiredHashRate/1000000)
		} else {
			hashRateStr = fmt.Sprintf("%.1f GH/s", requiredHashRate/1000000000)
		}

		fmt.Printf("    %d      |   %10.0f      | %s\n",
			diff, expectedAttempts, hashRateStr)
	}

	// Output:
	// Network hash rate needed for 10-minute blocks:
	// Difficulty | Expected Attempts | Required Hash Rate
	// -----------|-------------------|-------------------
	//     1      |           16      | 0 H/s
	//     2      |          256      | 0 H/s
	//     3      |         4096      | 6.8 H/s
	//     4      |        65536      | 109.2 H/s
	//     5      |      1048576      | 1.7 KH/s
}
