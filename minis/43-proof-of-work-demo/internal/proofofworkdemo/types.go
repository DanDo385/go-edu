package proofofworkdemo

// Block is a minimal block structure used by the PoW demo exercises/tests.
type Block struct {
	Index     int
	Timestamp int64
	Data      string
	PrevHash  string
	Nonce     int
	Hash      string
}

// Miner represents a participant in a mining pool.
type Miner struct {
	ID     string
	Shares int
}
