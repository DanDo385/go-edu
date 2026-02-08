//go:build reference

package devnets

/*
Reference Solution
==================

This file is the canonical reference for this exercise. It keeps failure paths
explicit when an operation can fail, so callers can decide how to handle
errors at API boundaries.

Read this alongside exercise.go and the tests to understand the intended data
flow, ownership boundaries, and invariants that keep behavior deterministic.

Teaching notes:
- Memory/ownership: make copies when returning mutable data that should not
  alias internal state; share references only when aliasing is intentional.
- Invariants: establish assumptions close to construction, and rely on them in
  smaller helper functions to keep logic easy to audit.
- Error surfaces: prefer explicit returns over hidden panics so learners can
  reason about control flow in production-style code.
*/

func Run() {
	type profile struct {
		Name      string
		ChainID   uint64
		RPC       string
		WS        string
		BlockTime int
	}

	profiles := map[string]profile{
		"local": {Name: "local", ChainID: 1337, RPC: "http://127.0.0.1:8545", WS: "ws://127.0.0.1:8546", BlockTime: 1},
		"anvil": {Name: "anvil", ChainID: 31337, RPC: "http://127.0.0.1:8545", WS: "ws://127.0.0.1:8545", BlockTime: 1},
	}

	selected := profiles["local"]
	_ = selected
}
