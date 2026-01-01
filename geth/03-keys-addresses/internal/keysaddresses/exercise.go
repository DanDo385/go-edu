//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package keysaddresses

const (
	defaultOutDir   = "./keystore-demo"
	defaultPassword = "changeit"
)
// TODO: implement Run.
func Run(cfg Config) (*Result, error) { panic("TODO: implement") }
