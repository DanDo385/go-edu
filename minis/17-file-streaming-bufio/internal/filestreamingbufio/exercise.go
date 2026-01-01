//go:build !solution && !reference

package filestreamingbufio

import (
	"bufio"
	"io"
	"strings"
)

// CountLines implements the exercise.
//
// TODO: Implement this function
func CountLines(reader io.Reader) (int, error) {
	// TODO: Implement
	return 0, nil
}

// FilterLines implements the exercise.
//
// TODO: Implement this function
func FilterLines(input io.Reader, output io.Writer, predicate func(string) bool) (int, error) {
	// TODO: Implement
	return 0, nil
}

// WordFrequency implements the exercise.
//
// TODO: Implement this function
func WordFrequency(reader io.Reader) (map[string]int, error) {
	// TODO: Implement
	return nil, nil
}

// TransformFile implements the exercise.
//
// TODO: Implement this function
func TransformFile(input io.Reader, output io.Writer, transform func(string) string) error {
	// TODO: Implement
	return nil
}

// ReadChunks implements the exercise.
//
// TODO: Implement this function
func ReadChunks(reader io.Reader, chunkSize int, callback func([]byte)) (int, error) {
	// TODO: Implement
	return 0, nil
}

// ReadChunksBuffered implements the exercise.
//
// TODO: Implement this function
func ReadChunksBuffered(reader io.Reader, chunkSize int, callback func([]byte)) (int, error) {
	// TODO: Implement
	return 0, nil
}
