//go:build reference

package node

/*
Reference Solution

This lesson currently has no callable API beyond `Run()`. The reference models
basic node lifecycle state transitions so the teaching intent is concrete.
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
