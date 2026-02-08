package indexer

import (
	"math/big"
)

// Config controls the indexer run.
type Config struct {
	FromBlock *big.Int
	ToBlock   *big.Int
}

// Result contains indexed data.
type Result struct {
	BlocksProcessed int
}
