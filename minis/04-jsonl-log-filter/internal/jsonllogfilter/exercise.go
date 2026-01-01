//go:build !solution && !reference

package jsonllogfilter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

func (l *Level) UnmarshalJSON(data []byte) error {
	// TODO: Implement this function
	// TODO: Handle errors appropriately
	// TODO: Add necessary validations
	panic("not implemented")
}

func FilterLogs(r io.Reader, minLevel Level) ([]Entry, error) {
	// TODO: Implement this function
	panic("not implemented")
}
