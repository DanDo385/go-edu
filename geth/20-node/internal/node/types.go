package node

// Config controls the node configuration.
type Config struct {
	DataDir  string
	SyncMode string
}

// Result contains node status.
type Result struct {
	Status string
}
