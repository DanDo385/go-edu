//go:build reference

package indexer

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
	type block struct {
		Height uint64
		Hash   string
		TxFrom string
		TxTo   string
	}

	chain := []block{
		{Height: 100, Hash: "0xaaa", TxFrom: "0x1", TxTo: "0x2"},
		{Height: 101, Hash: "0xbbb", TxFrom: "0x2", TxTo: "0x3"},
	}

	indexByHash := make(map[string]uint64, len(chain))
	activity := make(map[string][]uint64)

	for _, b := range chain {
		indexByHash[b.Hash] = b.Height
		activity[b.TxFrom] = append(activity[b.TxFrom], b.Height)
		activity[b.TxTo] = append(activity[b.TxTo], b.Height)
	}

	_ = indexByHash
	_ = activity
}
