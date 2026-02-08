package devnets

// Config controls the devnet setup.
type Config struct {
	RPCURL string
}

// Result contains devnet status.
type Result struct {
	ChainID int64
	Status  string
}
