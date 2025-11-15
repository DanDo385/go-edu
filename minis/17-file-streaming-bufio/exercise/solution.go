//go:build solution
// +build solution

package exercise

import (
	"bufio"
	"io"
	"strings"
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
	// Create a scanner to read line-by-line
	// Scanner uses buffered I/O internally (efficient)
	scanner := bufio.NewScanner(reader)

	count := 0

	// Scan advances to the next line
	// Returns false when EOF is reached or an error occurs
	for scanner.Scan() {
		count++
		// Note: We don't need to process scanner.Text()
		// Just counting lines, so we only increment
	}

	// Check if scanner stopped due to an error (vs normal EOF)
	if err := scanner.Err(); err != nil {
		return count, err
	}

	return count, nil
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
	// Create buffered reader for efficient reading
	scanner := bufio.NewScanner(input)

	// Create buffered writer for efficient writing
	// This batches writes to reduce system calls
	writer := bufio.NewWriter(output)
	// CRITICAL: Defer flush to ensure all buffered data is written
	defer writer.Flush()

	count := 0

	for scanner.Scan() {
		line := scanner.Text()

		// Test predicate
		if predicate(line) {
			// Write the line (to buffer, not directly to output)
			_, err := writer.WriteString(line + "\n")
			if err != nil {
				return count, err
			}
			count++
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return count, err
	}

	return count, nil
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
	scanner := bufio.NewScanner(reader)

	// Change split function to split on words instead of lines
	// ScanWords splits on whitespace (spaces, tabs, newlines)
	scanner.Split(bufio.ScanWords)

	freq := make(map[string]int)

	for scanner.Scan() {
		// Get word and convert to lowercase for case-insensitive counting
		word := strings.ToLower(scanner.Text())
		freq[word]++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return freq, nil
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
	scanner := bufio.NewScanner(input)
	writer := bufio.NewWriter(output)
	defer writer.Flush()

	for scanner.Scan() {
		line := scanner.Text()

		// Apply transformation
		transformed := transform(line)

		// Write transformed line
		_, err := writer.WriteString(transformed + "\n")
		if err != nil {
			return err
		}
	}

	return scanner.Err()
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
	// Create a buffer of the requested chunk size
	buffer := make([]byte, chunkSize)
	totalBytes := 0

	for {
		// Read up to chunkSize bytes
		// n is the actual number of bytes read (may be less than chunkSize)
		n, err := reader.Read(buffer)

		if n > 0 {
			// Call callback with the portion of buffer that was filled
			// Important: Pass buffer[:n], not buffer, since we may have read less than chunkSize
			callback(buffer[:n])
			totalBytes += n
		}

		// Check for errors
		if err == io.EOF {
			// EOF is expected, not an error
			break
		}
		if err != nil {
			// Actual error
			return totalBytes, err
		}
	}

	return totalBytes, nil
}

// Advanced implementation using bufio.Reader for better performance:
func ReadChunksBuffered(reader io.Reader, chunkSize int, callback func([]byte)) (int, error) {
	// Wrap in bufio.Reader for buffered reading
	// This can improve performance by reducing system calls
	bufferedReader := bufio.NewReaderSize(reader, chunkSize*2)

	buffer := make([]byte, chunkSize)
	totalBytes := 0

	for {
		n, err := bufferedReader.Read(buffer)

		if n > 0 {
			callback(buffer[:n])
			totalBytes += n
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return totalBytes, err
		}
	}

	return totalBytes, nil
}
