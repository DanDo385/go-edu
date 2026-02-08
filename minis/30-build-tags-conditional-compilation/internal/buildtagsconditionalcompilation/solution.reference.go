//go:build reference

package buildtagsconditionalcompilation

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

/*
Reference entrypoint for lesson 30.

This lesson teaches build constraints through multiple reference files:
- solution_unix.go / solution_windows.go
- solution_cloud.go / solution_local.go
- solution_debug.go / solution_prod.go
- solution_* arch files

The actual behavior lives in those files so each build target compiles the
correct implementation. This file intentionally provides package-level context.
*/

// ReferenceVariant reports the active reference build variant.
// It is used only as a stable, non-placeholder symbol for lesson validation.
func ReferenceVariant() string {
	return "build-tags-reference"
}
