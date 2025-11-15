package exercise

import (
	"math"
)

// Block represents a single block in the blockchain
type Block struct {
	Index     int
	Timestamp int64
	Data      string
	PrevHash  string
	Nonce     int
	Hash      string
}

// Exercise 1: Calculate Block Hash
// Implement a function that calculates the SHA-256 hash of a block.
//
// Requirements:
// - Concatenate fields in order: Index, Timestamp, Data, PrevHash, Nonce
// - Compute SHA-256 hash
// - Return as hexadecimal string
//
// Hints:
// - Use crypto/sha256 for hashing
// - Use fmt.Sprintf() to concatenate fields
// - Use encoding/hex to convert bytes to hex string
func CalculateBlockHash(b Block) string {
	// TODO: Implement this function
	// 1. Concatenate all block fields into a single string
	// 2. Compute SHA-256 hash of that string
	// 3. Convert the hash to a hexadecimal string
	// 4. Return the hex string

	return "" // Replace with your implementation
}

// Exercise 2: Validate Proof of Work
// Implement a function that checks if a hash meets the difficulty requirement.
//
// Requirements:
// - Check if hash starts with 'difficulty' number of zeros
// - Return true if valid, false otherwise
//
// Hints:
// - Use strings.Repeat() to create target string of zeros
// - Use strings.HasPrefix() to check the prefix
func IsValidProof(hash string, difficulty int) bool {
	// TODO: Implement this function
	// 1. Create a target string with 'difficulty' number of zeros
	// 2. Check if hash starts with this target
	// 3. Return true if it does, false otherwise

	return false // Replace with your implementation
}

// Exercise 3: Mine a Block
// Implement the mining algorithm that finds a valid nonce.
//
// Requirements:
// - Start with nonce = 0
// - Calculate hash, check if valid
// - If not valid, increment nonce and repeat
// - Update block's Nonce and Hash fields when found
// - Return the number of attempts
//
// Hints:
// - Use a loop that continues until proof is valid
// - Use CalculateBlockHash() and IsValidProof()
// - Don't forget to update both block.Nonce and block.Hash
func MineBlock(block *Block, difficulty int) int {
	// TODO: Implement this function
	// 1. Initialize attempts counter
	// 2. Loop until valid proof found:
	//    a. Set block.Nonce to current attempt
	//    b. Calculate hash
	//    c. Check if hash is valid
	//    d. If valid, update block.Hash and return attempts
	//    e. If not valid, increment attempts and continue
	// 3. Return number of attempts

	return 0 // Replace with your implementation
}

// Exercise 4: Validate Blockchain
// Implement a function that validates an entire blockchain.
//
// Requirements:
// - Check each block's hash is correct
// - Check each block's proof of work is valid
// - Check each block links to previous block
// - Return true if all checks pass, false otherwise
//
// Hints:
// - Loop through blocks starting at index 1 (skip genesis)
// - Use CalculateBlockHash() and IsValidProof()
// - Check block.PrevHash == previousBlock.Hash
func ValidateChain(chain []Block, difficulty int) bool {
	// TODO: Implement this function
	// 1. Handle edge case: empty or single block chain
	// 2. Loop through each block (starting from index 1):
	//    a. Get current and previous block
	//    b. Verify stored hash matches calculated hash
	//    c. Verify proof of work is valid
	//    d. Verify PrevHash links to previous block
	//    e. Return false if any check fails
	// 3. Return true if all blocks pass all checks

	return false // Replace with your implementation
}

// Exercise 5: Adjust Difficulty
// Implement difficulty adjustment to maintain target block time.
//
// Requirements:
// - Calculate actual time for last N blocks
// - Compare to expected time (targetBlockTime * N)
// - Increase difficulty if too fast (actual < expected/2)
// - Decrease difficulty if too slow (actual > expected*2)
// - Keep difficulty unchanged if in acceptable range
// - Minimum difficulty is 1
//
// Hints:
// - Use timestamps of first and last block in window
// - Number of time intervals = number of blocks - 1
// - Return currentDifficulty if chain is too short
func AdjustDifficulty(chain []Block, targetBlockTime int64, currentDifficulty int) int {
	// TODO: Implement this function
	// 1. Check if chain has enough blocks (need at least 2)
	// 2. Determine window size (e.g., last 10 blocks)
	// 3. Calculate actual time: lastBlock.Timestamp - firstBlock.Timestamp
	// 4. Calculate expected time: targetBlockTime * (number of intervals)
	// 5. Compare actual vs expected:
	//    - If actual < expected/2: increase difficulty
	//    - If actual > expected*2: decrease difficulty (min 1)
	//    - Otherwise: keep current difficulty
	// 6. Return adjusted difficulty

	return currentDifficulty // Replace with your implementation
}

// Exercise 6: Calculate Mining Probability
// Calculate the probability of finding a block within given time.
//
// Requirements:
// - Calculate expected attempts based on difficulty (16^difficulty)
// - Calculate lambda: hashRate / expectedAttempts
// - Use Poisson formula: P = 1 - e^(-lambda * time)
// - Return probability as a float between 0 and 1
//
// Hints:
// - Use math.Pow(16, difficulty) for expected attempts
// - Use math.Exp(x) for e^x
// - Lambda represents the rate parameter of the Poisson process
func MiningProbability(hashRate float64, difficulty int, timeSeconds float64) float64 {
	// TODO: Implement this function
	// 1. Calculate expected attempts: 16^difficulty
	// 2. Calculate lambda: hashRate / expectedAttempts
	// 3. Calculate probability: 1 - e^(-lambda * time)
	// 4. Return the probability

	return 0.0 // Replace with your implementation
}

// Helper function for testing - creates a simple blockchain
func CreateTestChain(numBlocks int, difficulty int) []Block {
	chain := make([]Block, numBlocks)

	// Genesis block
	chain[0] = Block{
		Index:     0,
		Timestamp: 1000,
		Data:      "Genesis Block",
		PrevHash:  "0",
		Nonce:     0,
	}
	MineBlock(&chain[0], difficulty)

	// Subsequent blocks
	for i := 1; i < numBlocks; i++ {
		chain[i] = Block{
			Index:     i,
			Timestamp: chain[i-1].Timestamp + 10, // 10 second intervals
			Data:      "Block " + string(rune('0'+i)),
			PrevHash:  chain[i-1].Hash,
			Nonce:     0,
		}
		MineBlock(&chain[i], difficulty)
	}

	return chain
}

// STRETCH GOAL: Merkle Root
// Build a Merkle tree from a list of transactions and return the root hash.
//
// A Merkle tree is a binary tree of hashes where:
// - Leaf nodes are hashes of individual transactions
// - Internal nodes are hashes of their children concatenated
// - The root is a single hash representing all transactions
//
// This allows efficient verification that a transaction is in a block.
func BuildMerkleRoot(transactions []string) string {
	// TODO: (Stretch goal - optional)
	// 1. If no transactions, return empty or zero hash
	// 2. Create leaf nodes (hash of each transaction)
	// 3. While more than one node remains:
	//    a. Pair up nodes
	//    b. Hash each pair together
	//    c. If odd number, duplicate last node
	// 4. Return the final root hash

	return ""
}

// STRETCH GOAL: Mining Pool Share Calculation
// Calculate how rewards should be distributed among pool miners
// based on their contributed hash power.
//
// A mining pool combines hash power from multiple miners.
// When the pool finds a block, the reward is split proportionally.
type Miner struct {
	ID       string
	HashRate float64
	Shares   int // Number of "near-valid" hashes submitted
}

func CalculatePoolRewards(miners []Miner, blockReward float64) map[string]float64 {
	// TODO: (Stretch goal - optional)
	// 1. Calculate total shares from all miners
	// 2. For each miner, calculate their proportion: miner.Shares / totalShares
	// 3. Multiply proportion by blockReward
	// 4. Return map of miner ID to reward amount

	return nil
}

// STRETCH GOAL: Estimate Attack Cost
// Calculate the cost for an attacker to rewrite N blocks.
//
// To rewrite history, an attacker must:
// 1. Re-mine the target block with modified data
// 2. Re-mine all subsequent blocks
// 3. Catch up to (and surpass) the honest chain
//
// This requires enormous computational resources.
func EstimateAttackCost(
	numBlocks int,
	difficulty int,
	honestHashRate float64,
	attackerHashRate float64,
	electricityCostPerKWh float64,
	minerWattage float64,
) float64 {
	// TODO: (Stretch goal - optional)
	// 1. Calculate expected time for attacker to mine N blocks
	// 2. During that time, honest miners add more blocks
	// 3. Calculate total blocks attacker must mine
	// 4. Calculate energy consumption
	// 5. Return total electricity cost

	return 0.0
}

// STRETCH GOAL: Selfish Mining Simulation
// Simulate a selfish mining attack where a miner withholds blocks
// to gain advantage over honest miners.
//
// Selfish mining: Instead of broadcasting found blocks immediately,
// attacker keeps them secret and continues mining on their private chain.
// If they get ahead, they broadcast their longer chain, making honest
// miners' work worthless.
func SimulateSelfishMining(
	honestHashRate float64,
	selfishHashRate float64,
	numBlocks int,
	difficulty int,
) (honestBlocks int, selfishBlocks int) {
	// TODO: (Stretch goal - optional)
	// 1. Both honest and selfish miners mine blocks
	// 2. Honest miners broadcast immediately
	// 3. Selfish miner withholds blocks
	// 4. Selfish miner broadcasts when they're ahead
	// 5. Return count of blocks each side successfully added

	return 0, 0
}

// Additional utility functions you might need:

// CountLeadingZeros counts the number of leading zero characters in a hex string
func CountLeadingZeros(hash string) int {
	count := 0
	for _, char := range hash {
		if char == '0' {
			count++
		} else {
			break
		}
	}
	return count
}

// ExpectedMiningTime calculates the expected time to mine a block
func ExpectedMiningTime(hashRate float64, difficulty int) float64 {
	expectedAttempts := math.Pow(16, float64(difficulty))
	return expectedAttempts / hashRate
}
