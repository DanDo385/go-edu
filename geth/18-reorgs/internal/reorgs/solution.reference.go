//go:build reference

package reorgs

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
	type head struct {
		Hash   string
		Parent string
	}

	canonical := []head{{Hash: "A", Parent: ""}, {Hash: "B", Parent: "A"}, {Hash: "C", Parent: "B"}}
	incoming := []head{{Hash: "X", Parent: "B"}, {Hash: "Y", Parent: "X"}}

	tip := canonical[len(canonical)-1]
	if incoming[0].Parent != tip.Hash {
		for len(canonical) > 0 && canonical[len(canonical)-1].Hash != incoming[0].Parent {
			canonical = canonical[:len(canonical)-1]
		}
	}
	canonical = append(canonical, incoming...)
	_ = canonical
}
