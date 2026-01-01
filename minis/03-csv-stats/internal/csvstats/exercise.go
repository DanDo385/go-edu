//go:build !solution && !reference

package csvstats

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

type Stat struct {
	Count int
	Sum   float64
	Avg   float64
}

// SummarizeCSV implements the exercise.
//
// TODO: Implement this function
func SummarizeCSV(r io.Reader) (map[string]Stat, error) {
	// TODO: Implement
	return nil, nil
}
