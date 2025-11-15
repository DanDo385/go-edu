package exercise

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
)

// Solution for Exercise 1: Calculate Block Hash
func CalculateBlockHashSolution(b Block) string {
	// Concatenate all block fields in a consistent order
	record := fmt.Sprintf("%d%d%s%s%d",
		b.Index,
		b.Timestamp,
		b.Data,
		b.PrevHash,
		b.Nonce,
	)

	// Compute SHA-256 hash
	hash := sha256.Sum256([]byte(record))

	// Convert to hexadecimal string
	return hex.EncodeToString(hash[:])
}

// Solution for Exercise 2: Validate Proof of Work
func IsValidProofSolution(hash string, difficulty int) bool {
	// Create target string with required number of zeros
	target := strings.Repeat("0", difficulty)

	// Check if hash starts with the target
	return strings.HasPrefix(hash, target)
}

// Solution for Exercise 3: Mine a Block
func MineBlockSolution(block *Block, difficulty int) int {
	attempts := 0

	// Keep trying until we find a valid proof
	for {
		// Set the nonce
		block.Nonce = attempts

		// Calculate hash with current nonce
		block.Hash = CalculateBlockHashSolution(*block)

		// Increment attempts counter
		attempts++

		// Check if we found a valid proof
		if IsValidProofSolution(block.Hash, difficulty) {
			return attempts
		}

		// Safety check: if nonce space exhausted, we could update timestamp
		// For this implementation, we assume difficulty isn't impossibly high
	}
}

// Solution for Exercise 4: Validate Blockchain
func ValidateChainSolution(chain []Block, difficulty int) bool {
	// Handle edge case: empty chain or single block
	if len(chain) == 0 {
		return true
	}
	if len(chain) == 1 {
		// Just validate the genesis block
		calculatedHash := CalculateBlockHashSolution(chain[0])
		return chain[0].Hash == calculatedHash &&
			IsValidProofSolution(chain[0].Hash, difficulty)
	}

	// Validate each block starting from index 1
	for i := 1; i < len(chain); i++ {
		current := chain[i]
		previous := chain[i-1]

		// Check 1: Stored hash matches calculated hash
		calculatedHash := CalculateBlockHashSolution(current)
		if current.Hash != calculatedHash {
			return false // Hash was tampered with
		}

		// Check 2: Proof of work is valid
		if !IsValidProofSolution(current.Hash, difficulty) {
			return false // Didn't do the work
		}

		// Check 3: Block links to previous block
		if current.PrevHash != previous.Hash {
			return false // Chain is broken
		}
	}

	// Also validate the genesis block
	calculatedHash := CalculateBlockHashSolution(chain[0])
	if chain[0].Hash != calculatedHash {
		return false
	}
	if !IsValidProofSolution(chain[0].Hash, difficulty) {
		return false
	}

	return true
}

// Solution for Exercise 5: Adjust Difficulty
func AdjustDifficultySolution(chain []Block, targetBlockTime int64, currentDifficulty int) int {
	// Need at least 2 blocks to calculate time
	if len(chain) < 2 {
		return currentDifficulty
	}

	// Use last 10 blocks for adjustment (or fewer if chain is shorter)
	windowSize := 10
	if len(chain) < windowSize {
		windowSize = len(chain)
	}

	// Calculate actual time taken
	startIndex := len(chain) - windowSize
	actualTime := chain[len(chain)-1].Timestamp - chain[startIndex].Timestamp

	// Calculate expected time
	// Number of intervals = number of blocks - 1
	numIntervals := int64(windowSize - 1)
	expectedTime := targetBlockTime * numIntervals

	// Adjust difficulty based on actual vs expected time
	// Use a factor of 2x to avoid overreacting to variance

	if actualTime < expectedTime/2 {
		// Mining too fast, increase difficulty
		return currentDifficulty + 1
	} else if actualTime > expectedTime*2 {
		// Mining too slow, decrease difficulty (but minimum is 1)
		if currentDifficulty > 1 {
			return currentDifficulty - 1
		}
		return 1
	}

	// Mining rate is acceptable, keep current difficulty
	return currentDifficulty
}

// Solution for Exercise 6: Calculate Mining Probability
func MiningProbabilitySolution(hashRate float64, difficulty int, timeSeconds float64) float64 {
	// Calculate expected number of attempts for this difficulty
	// For hexadecimal leading zeros, each zero multiplies attempts by 16
	expectedAttempts := math.Pow(16, float64(difficulty))

	// Calculate lambda (rate parameter)
	// Lambda represents the expected number of blocks found per second
	lambda := hashRate / expectedAttempts

	// Use Poisson distribution formula: P(at least 1 success) = 1 - e^(-λt)
	// This models the probability of finding at least one valid block
	// within the given time period
	probability := 1 - math.Exp(-lambda*timeSeconds)

	return probability
}

// STRETCH GOAL Solution: Build Merkle Root
func BuildMerkleRootSolution(transactions []string) string {
	if len(transactions) == 0 {
		return ""
	}

	// Create leaf nodes (hash each transaction)
	var nodes []string
	for _, tx := range transactions {
		hash := sha256.Sum256([]byte(tx))
		nodes = append(nodes, hex.EncodeToString(hash[:]))
	}

	// Build tree level by level until we have a single root
	for len(nodes) > 1 {
		var nextLevel []string

		// Process nodes in pairs
		for i := 0; i < len(nodes); i += 2 {
			var combined string

			if i+1 < len(nodes) {
				// Pair exists, combine both
				combined = nodes[i] + nodes[i+1]
			} else {
				// Odd number of nodes, duplicate the last one
				combined = nodes[i] + nodes[i]
			}

			// Hash the combined pair
			hash := sha256.Sum256([]byte(combined))
			nextLevel = append(nextLevel, hex.EncodeToString(hash[:]))
		}

		nodes = nextLevel
	}

	return nodes[0]
}

// STRETCH GOAL Solution: Calculate Pool Rewards
func CalculatePoolRewardsSolution(miners []Miner, blockReward float64) map[string]float64 {
	rewards := make(map[string]float64)

	// Calculate total shares
	totalShares := 0
	for _, miner := range miners {
		totalShares += miner.Shares
	}

	// Avoid division by zero
	if totalShares == 0 {
		return rewards
	}

	// Distribute rewards proportionally
	for _, miner := range miners {
		proportion := float64(miner.Shares) / float64(totalShares)
		rewards[miner.ID] = proportion * blockReward
	}

	return rewards
}

// STRETCH GOAL Solution: Estimate Attack Cost
func EstimateAttackCostSolution(
	numBlocks int,
	difficulty int,
	honestHashRate float64,
	attackerHashRate float64,
	electricityCostPerKWh float64,
	minerWattage float64,
) float64 {
	// Calculate expected attempts per block
	expectedAttemptsPerBlock := math.Pow(16, float64(difficulty))

	// Calculate time for attacker to mine N blocks
	attackerTimePerBlock := expectedAttemptsPerBlock / attackerHashRate
	attackerTotalTime := attackerTimePerBlock * float64(numBlocks)

	// During this time, honest miners also mine blocks
	honestBlocksAdded := (honestHashRate / expectedAttemptsPerBlock) * attackerTotalTime

	// Attacker must mine original N blocks PLUS the blocks honest miners added
	totalBlocksToMine := float64(numBlocks) + honestBlocksAdded

	// Calculate total time needed
	totalTime := (expectedAttemptsPerBlock * totalBlocksToMine) / attackerHashRate

	// Calculate energy consumption
	// Convert seconds to hours, multiply by wattage to get watt-hours
	// Divide by 1000 to get kilowatt-hours
	energyKWh := (totalTime / 3600) * minerWattage / 1000

	// Calculate cost
	totalCost := energyKWh * electricityCostPerKWh

	return totalCost
}

// STRETCH GOAL Solution: Simulate Selfish Mining
func SimulateSelfishMiningSolution(
	honestHashRate float64,
	selfishHashRate float64,
	numBlocks int,
	difficulty int,
) (honestBlocks int, selfishBlocks int) {
	totalHashRate := honestHashRate + selfishHashRate

	// Probability that selfish miner finds next block
	selfishProbability := selfishHashRate / totalHashRate

	// Simplified simulation:
	// This is a basic model. Real selfish mining is more complex
	// and depends on network propagation delays and strategy.

	privateChainLength := 0  // Blocks selfish miner has but hasn't broadcast
	publicChainLength := 0   // Blocks on the honest chain
	selfishBlocksKept := 0   // Blocks selfish miner successfully added

	for publicChainLength+selfishBlocksKept < numBlocks {
		// Determine who finds the next block
		if randFloat() < selfishProbability {
			// Selfish miner found a block
			privateChainLength++

			// Strategy: If private chain is longer, broadcast it
			if privateChainLength > publicChainLength {
				// Selfish miner's chain wins
				selfishBlocksKept += privateChainLength
				privateChainLength = 0
				publicChainLength = 0
			}
		} else {
			// Honest miner found a block
			publicChainLength++

			// If selfish miner has private blocks, there's a race
			if privateChainLength > 0 {
				if privateChainLength >= publicChainLength {
					// Selfish chain wins
					selfishBlocksKept += privateChainLength
					privateChainLength = 0
					publicChainLength = 0
				} else {
					// Honest chain wins, selfish miner loses their private blocks
					honestBlocks += publicChainLength
					privateChainLength = 0
					publicChainLength = 0
				}
			} else {
				// No competition, honest block is accepted
				honestBlocks++
				publicChainLength = 0
			}
		}
	}

	selfishBlocks = selfishBlocksKept
	return honestBlocks, selfishBlocks
}

// Simple pseudo-random float generator for simulation
// In real code, use math/rand properly seeded
var randState = uint64(12345)

func randFloat() float64 {
	// Linear congruential generator (simple PRNG)
	randState = (1103515245*randState + 12345) % (1 << 31)
	return float64(randState) / float64(1<<31)
}

// Alternative implementation using a different approach for Exercise 3
// This version includes timestamp updates if nonce space is exhausted
func MineBlockWithTimestampUpdate(block *Block, difficulty int) int {
	attempts := 0
	target := strings.Repeat("0", difficulty)

	for {
		block.Hash = CalculateBlockHashSolution(*block)
		attempts++

		if strings.HasPrefix(block.Hash, target) {
			return attempts
		}

		block.Nonce++

		// If we've exhausted the nonce space, update timestamp
		if block.Nonce >= math.MaxInt32 {
			block.Timestamp++
			block.Nonce = 0
		}
	}
}
