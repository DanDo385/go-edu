package reorgs

import (
	"math/big"
)

// Config controls the reorg detection run.
type Config struct {
	FromBlock *big.Int
	ToBlock   *big.Int
}

// Result contains reorg detection data.
type Result struct {
	ReorgsDetected int
}
