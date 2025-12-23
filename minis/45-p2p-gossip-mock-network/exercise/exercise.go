// ============================================================================
// Exercise 5: Implement a gossip simulator
// ============================================================================

type Simulator struct {
	nodes    []GossipNode
	network  Network
	detector *ConvergenceDetector
	mu       sync.RWMutex
}

func NewSimulator(nodeCount, fanout int, latency time.Duration, dropRate float64) *Simulator {
	// TODO: Implement this function to set up the entire simulation.

	// Step 1: Create a `mockNetwork`.
	// Step 2: Create `nodeCount` `gossipNode`s, registering each one with the network.
	// Step 3: Build the network topology. A good approach is to:
	//   a. Create a ring where every node is connected to its neighbor to guarantee a connected graph.
	//   b. Add some random connections between nodes to improve message propagation speed.
	// Step 4: Create a `ConvergenceDetector` with the list of nodes.
	// Step 5: Return the fully initialized `Simulator`.
	return nil
}

// SimulationStats contains statistics about the simulation
type SimulationStats struct {
	NodeCount    int
	MessageCount int
	DroppedCount int
}

// GetStats returns statistics about the simulation.
func (s *Simulator) GetStats() SimulationStats {
	// TODO: Implement this method.
	// - Get the message count from the network.
	// - Return a `SimulationStats` struct.
	return SimulationStats{}
}

// Shutdown stops all nodes in the simulation.
func (s *Simulator) Shutdown() {
	// TODO: Implement this method.
	// - Loop through all nodes and call their `Shutdown` method.
}