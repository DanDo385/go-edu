//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package filestreamingbufio

import "io"
// TODO: implement CountLines.
func CountLines(reader io.Reader) (int, error) { panic("TODO: implement") }
// TODO: implement FilterLines.
func FilterLines(input io.Reader, output io.Writer, predicate func(string) bool) (int, error) {
	panic("TODO: implement")
}
// TODO: implement WordFrequency.
func WordFrequency(reader io.Reader) (map[string]int, error) { panic("TODO: implement") }
// TODO: implement TransformFile.
func TransformFile(input io.Reader, output io.Writer, transform func(string) string) error {
	panic("TODO: implement")
}
// TODO: implement ReadChunks.
func ReadChunks(reader io.Reader, chunkSize int, callback func([]byte)) (int, error) {
	panic("TODO: implement")
}
// TODO: implement ReadChunksBuffered.
func ReadChunksBuffered(reader io.Reader, chunkSize int, callback func([]byte)) (int, error) {
	panic("TODO: implement")
}
