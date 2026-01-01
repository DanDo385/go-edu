//go:build !solution && !reference

package filestreamingbufio

import (
	"bufio"
	"io"
	"strings"
)

func CountLines(reader io.Reader) (int, error) {
	// TODO: Implement this function
	// TODO: Handle errors appropriately
	// TODO: Add necessary validations
	panic("not implemented")
}

func FilterLines(input io.Reader, output io.Writer, predicate func(string) bool) (int, error) {
	// TODO: Implement this function
	panic("not implemented")
}

func WordFrequency(reader io.Reader) (map[string]int, error) {
	// TODO: Implement this function
	panic("not implemented")
}

func TransformFile(input io.Reader, output io.Writer, transform func(string) string) error {
	// TODO: Implement this function
	panic("not implemented")
}

func ReadChunks(reader io.Reader, chunkSize int, callback func([]byte)) (int, error) {
	// TODO: Implement this function
	panic("not implemented")
}

func ReadChunksBuffered(reader io.Reader, chunkSize int, callback func([]byte)) (int, error) {
	// TODO: Implement this function
	panic("not implemented")
}
