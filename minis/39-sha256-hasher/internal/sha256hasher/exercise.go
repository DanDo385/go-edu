//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package sha256hasher
// TODO: implement HashString.
func HashString(s string) string { panic("TODO: implement") }
// TODO: implement HashFile.
func HashFile(filename string) ([]byte, error) { panic("TODO: implement") }
// TODO: implement VerifyFile.
func VerifyFile(filename, expectedHashHex string) (bool, error) { panic("TODO: implement") }
// TODO: implement HashIncremental.
func HashIncremental(parts ...string) string { panic("TODO: implement") }
// TODO: implement CompareHashes.
func CompareHashes(hash1, hash2 string) bool { panic("TODO: implement") }
