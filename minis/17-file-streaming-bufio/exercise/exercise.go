//go:build !solution
// +build !solution

package exercise

import (
	"io"
)

// CountLines counts the number of lines in a file.
// It should handle files of any size efficiently (streaming, not loading entire file).
//
// Parameters:
//   - reader: io.Reader to read from (could be a file, network stream, etc.)
//
// Returns:
//   - int: Number of lines (separated by '\n')
//   - error: Non-nil if reading fails
//
// Example:
//   file, _ := os.Open("data.txt")
//   count, err := CountLines(file)
//   // If file contains "line1\nline2\nline3", returns 3
func CountLines(reader io.Reader) (int, error) {
	// TODO: implement using bufio.Scanner
	// Hint: Use bufio.NewScanner(reader) and scan line-by-line
	return 0, nil
}

// FilterLines reads from input, writes lines matching the predicate to output.
// It should stream data (not load entire file into memory).
//
// Parameters:
//   - input: io.Reader to read from
//   - output: io.Writer to write to
//   - predicate: Function that returns true for lines to keep
//
// Returns:
//   - int: Number of lines written
//   - error: Non-nil if reading/writing fails
//
// Example:
//   in, _ := os.Open("input.txt")
//   out, _ := os.Create("output.txt")
//   count, err := FilterLines(in, out, func(line string) bool {
//       return strings.Contains(line, "ERROR")
//   })
//   // Copies only lines containing "ERROR"
func FilterLines(input io.Reader, output io.Writer, predicate func(string) bool) (int, error) {
	// TODO: implement using bufio.Scanner and bufio.Writer
	// Hint: Don't forget to Flush() the writer!
	return 0, nil
}

// WordFrequency counts word frequencies in the input.
// Words are case-insensitive and separated by whitespace.
//
// Parameters:
//   - reader: io.Reader to read from
//
// Returns:
//   - map[string]int: Word frequencies (lowercase words)
//   - error: Non-nil if reading fails
//
// Example:
//   input := strings.NewReader("Go is great. Go is fast.")
//   freq, _ := WordFrequency(input)
//   // Returns: {"go": 2, "is": 2, "great.": 1, "fast.": 1}
func WordFrequency(reader io.Reader) (map[string]int, error) {
	// TODO: implement using bufio.Scanner with ScanWords
	// Hint: Use scanner.Split(bufio.ScanWords) to split by words
	return nil, nil
}

// TransformFile reads from input, applies transform to each line, writes to output.
// It should handle large files efficiently.
//
// Parameters:
//   - input: io.Reader to read from
//   - output: io.Writer to write to
//   - transform: Function to transform each line
//
// Returns:
//   - error: Non-nil if reading/writing fails
//
// Example:
//   transform := func(line string) string {
//       return strings.ToUpper(line)
//   }
//   TransformFile(in, out, transform)
//   // "hello" becomes "HELLO"
func TransformFile(input io.Reader, output io.Writer, transform func(string) string) error {
	// TODO: implement using bufio.Scanner and bufio.Writer
	return nil
}

// ReadChunks reads data in fixed-size chunks and calls the callback for each chunk.
// This is useful for processing binary data or implementing custom parsing.
//
// Parameters:
//   - reader: io.Reader to read from
//   - chunkSize: Size of each chunk in bytes
//   - callback: Function called with each chunk (may be smaller than chunkSize at EOF)
//
// Returns:
//   - int: Total bytes read
//   - error: Non-nil if reading fails
//
// Example:
//   total, err := ReadChunks(file, 1024, func(chunk []byte) {
//       fmt.Printf("Read %d bytes\n", len(chunk))
//   })
func ReadChunks(reader io.Reader, chunkSize int, callback func([]byte)) (int, error) {
	// TODO: implement using bufio.Reader or direct Read() calls
	// Hint: Create a buffer of size chunkSize and read into it
	return 0, nil
}
