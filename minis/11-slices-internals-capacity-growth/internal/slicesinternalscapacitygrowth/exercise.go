//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package slicesinternalscapacitygrowth
// TODO: implement GrowSlice.
func GrowSlice(s []int, elem int) (newSlice []int, oldCap, newCap int) { panic("TODO: implement") }
// TODO: implement SharesBackingArray.
func SharesBackingArray(a, b []int) bool { panic("TODO: implement") }
// TODO: implement SafeTruncate.
func SafeTruncate(s []int, n int) []int { panic("TODO: implement") }
// TODO: implement PreallocateVsDynamic.
func PreallocateVsDynamic(n int) (dynamicAllocs, preallocAllocs int) { panic("TODO: implement") }
// TODO: implement ReSliceWithCapLimit.
func ReSliceWithCapLimit(s []int, start, end int) []int { panic("TODO: implement") }
