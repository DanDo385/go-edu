//go:build reference

package node

/*
Reference Solution - Node State Machine
=======================================

This file demonstrates a minimal node lifecycle state machine: stopped →
booting → ready. Foundation for node startup and lifecycle management.

This connects to the Ethereum ecosystem by showing:
- Node states: stopped (idle), booting (startup), ready (synced/serving)
- type state string: simple enum pattern; const for each state
- current = stateX: state transitions; in real impl, guard by current state

The exercise builds understanding of:
- String-backed enum: type state string prevents mixing with plain strings
- State transitions: sequential here; production would validate valid transitions
*/
func Run() {
	type state string
	const (
		stateStopped state = "stopped"
		stateBooting state = "booting"
		stateReady   state = "ready"
	)

	current := stateStopped
	current = stateBooting
	current = stateReady
	_ = current
}
