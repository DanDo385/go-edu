//go:build reference

package buildtagsconditionalcompilation

/*
Reference Solution - Build Tags and Conditional Compilation
==========================================================

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
