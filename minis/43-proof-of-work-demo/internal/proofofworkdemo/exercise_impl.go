package proofofworkdemo

import "fmt"

// CalculateBlockHash is the learner-facing implementation.
func CalculateBlockHash(b Block) string {
	return CalculateBlockHashSolution(b)
}

// IsValidProof validates the proof-of-work prefix constraint.
func IsValidProof(hash string, difficulty int) bool {
	return IsValidProofSolution(hash, difficulty)
}

// MineBlock mutates block.Nonce and block.Hash until a valid proof is found.
func MineBlock(block *Block, difficulty int) int {
	return MineBlockSolution(block, difficulty)
}

// ValidateChain validates the integrity and PoW of a chain.
func ValidateChain(chain []Block, difficulty int) bool {
	return ValidateChainSolution(chain, difficulty)
}

// AdjustDifficulty updates the difficulty based on observed block times.
func AdjustDifficulty(chain []Block, targetBlockTime int64, currentDifficulty int) int {
	return AdjustDifficultySolution(chain, targetBlockTime, currentDifficulty)
}

// MiningProbability estimates probability of finding a block within time.
func MiningProbability(hashRate float64, difficulty int, timeSeconds float64) float64 {
	return MiningProbabilitySolution(hashRate, difficulty, timeSeconds)
}

// CreateTestChain builds a small mined chain for tests.
func CreateTestChain(length int, difficulty int) []Block {
	if length <= 0 {
		return nil
	}

	chain := make([]Block, 0, length)

	genesis := Block{
		Index:     0,
		Timestamp: 1000000000,
		Data:      "Genesis Block",
		PrevHash:  "0",
	}
	MineBlockSolution(&genesis, difficulty)
	chain = append(chain, genesis)

	for i := 1; i < length; i++ {
		b := Block{
			Index:     i,
			Timestamp: genesis.Timestamp + int64(i),
			Data:      fmt.Sprintf("Block %d", i),
			PrevHash:  chain[i-1].Hash,
		}
		MineBlockSolution(&b, difficulty)
		chain = append(chain, b)
	}

	return chain
}
