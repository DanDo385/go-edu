//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package csvstats

import "io"

type Stat struct {
	Count int
	Sum   float64
	Avg   float64
}
// TODO: implement SummarizeCSV.
func SummarizeCSV(r io.Reader) (map[string]Stat, error) { panic("TODO: implement") }
