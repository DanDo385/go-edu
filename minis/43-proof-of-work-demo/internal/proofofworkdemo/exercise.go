//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package proofofworkdemo
// TODO: implement CalculateBlockHashSolution.
func CalculateBlockHashSolution(b Block) string { panic("TODO: implement") }
// TODO: implement IsValidProofSolution.
func IsValidProofSolution(hash string, difficulty int) bool { panic("TODO: implement") }
// TODO: implement MineBlockSolution.
func MineBlockSolution(block *Block, difficulty int) int { panic("TODO: implement") }
// TODO: implement ValidateChainSolution.
func ValidateChainSolution(chain []Block, difficulty int) bool { panic("TODO: implement") }
// TODO: implement AdjustDifficultySolution.
func AdjustDifficultySolution(chain []Block, targetBlockTime int64, currentDifficulty int) int {
	panic("TODO: implement")
}
// TODO: implement MiningProbabilitySolution.
func MiningProbabilitySolution(hashRate float64, difficulty int, timeSeconds float64) float64 {
	panic("TODO: implement")
}
// TODO: implement BuildMerkleRootSolution.
func BuildMerkleRootSolution(transactions []string) string { panic("TODO: implement") }
// TODO: implement CalculatePoolRewardsSolution.
func CalculatePoolRewardsSolution(miners []Miner, blockReward float64) map[string]float64 {
	panic("TODO: implement")
}
// TODO: implement EstimateAttackCostSolution.
func EstimateAttackCostSolution(
	numBlocks int,
	difficulty int,
	honestHashRate float64,
	attackerHashRate float64,
	electricityCostPerKWh float64,
	minerWattage float64,
) float64 {
	panic("TODO: implement")
}
// TODO: implement SimulateSelfishMiningSolution.
func SimulateSelfishMiningSolution(
	honestHashRate float64,
	selfishHashRate float64,
	numBlocks int,
	difficulty int,
) (honestBlocks int, selfishBlocks int) {
	panic("TODO: implement")
}

var randState = uint64(12345)
// TODO: implement randFloat.
func randFloat() float64 { panic("TODO: implement") }
// TODO: implement MineBlockWithTimestampUpdate.
func MineBlockWithTimestampUpdate(block *Block, difficulty int) int { panic("TODO: implement") }
