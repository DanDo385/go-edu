//go:build !solution && !reference

package filestreamingbufio

// TODO: Import required packages
// You'll need:
// - "bufio" for buffered I/O (Scanner, Reader, Writer)
// - "io" for io.Reader, io.Writer, io.EOF interfaces
// - "strings" for string manipulation (ToLower, etc.)
//
// import (
//     "bufio"
//     "io"
//     "strings"
// )

// ============================================================================
// FILE STREAMING WITH BUFIO: Understanding Efficient I/O
// ============================================================================
//
// Why buffered I/O?
// - Reading byte-by-byte from disk/network is SLOW (many system calls)
// - Buffering reads larger chunks at once, reduces system calls
// - Example: Reading 1MB file byte-by-byte = 1,000,000 system calls
//            Reading with 4KB buffer = 250 system calls
//
// bufio.Scanner:
// - High-level API for reading line-by-line or word-by-word
// - Default buffer: 64KB (configurable with scanner.Buffer())
// - Automatically handles newlines, EOF
// - Pattern: scanner.Scan() returns false at EOF or error
// - Use scanner.Text() to get current line/token
// - Use scanner.Err() to check for errors after loop
//
// bufio.Reader:
// - Lower-level buffered reading
// - Methods: ReadByte(), ReadLine(), ReadString(), Read()
// - Use when you need more control than Scanner provides
// - Default buffer: 4096 bytes
//
// bufio.Writer:
// - Buffered writing (batches writes to reduce system calls)
// - CRITICAL: Must call Flush() to ensure all data is written
// - Not flushing = data loss! (buffer not written to disk)
// - Use defer writer.Flush() for safety
// - Default buffer: 4096 bytes
//
// Memory management:
// - Scanner/Reader/Writer allocate internal buffers on heap
// - Scanner.Text() returns string (may allocate if not cached)
// - For huge files, consider io.Copy or custom buffer reuse
//
// Streaming vs loading all:
// - Streaming: Process data as it arrives (constant memory)
// - Loading all: Read entire file into memory (O(n) memory)
// - Always prefer streaming for large files
//
// ============================================================================

// ============================================================================
// Exercise 1: CountLines
// ============================================================================

// CountLines counts the number of lines in a file.
// TODO: Implement efficient line counting using bufio.Scanner.
//
// Parameters:
//   - reader: io.Reader to read from (could be file, network, stdin)
//
// Returns:
//   - int: Number of lines (separated by '\n')
//   - error: Non-nil if reading fails
//
// Behavior:
//   - Read file line-by-line without loading entire file into memory
//   - Each line is defined as text ending with '\n'
//   - Empty file returns 0 lines
//   - File with no trailing newline: last line still counts
//
// Memory efficiency:
//   - Scanner uses ~64KB buffer (reused for all lines)
//   - Does NOT allocate memory for each line (we don't use Text())
//   - Total memory: O(1) regardless of file size
//
// Implementation:
//   - Create scanner: scanner := bufio.NewScanner(reader)
//   - Loop: for scanner.Scan() { count++ }
//   - Check error: if err := scanner.Err(); err != nil { return count, err }
//
// Why not use io.ReadAll?
//   - io.ReadAll loads entire file into memory: O(n) memory
//   - For 1GB file: uses 1GB RAM
//   - Scanner uses constant memory: O(1)
//
// func CountLines(reader io.Reader) (int, error) {
//     scanner := bufio.NewScanner(reader)
//     count := 0
//
//     // Scan advances to the next line
//     for scanner.Scan() {
//         count++
//         // Note: We don't need scanner.Text() since we're just counting
//     }
//
//     // Check if scanner stopped due to error (vs normal EOF)
//     if err := scanner.Err(); err != nil {
//         return count, err
//     }
//
//     return count, nil
// }

func CountLines(reader io.Reader) (int, error) {
	// TODO: Implement CountLines
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// ============================================================================
// Exercise 2: FilterLines
// ============================================================================

// FilterLines reads from input, writes lines matching predicate to output.
// TODO: Implement streaming filter using bufio.Scanner and bufio.Writer.
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
// Behavior:
//   - Read input line-by-line
//   - For each line, test predicate(line)
//   - If true, write line to output (with newline)
//   - If false, skip line
//   - CRITICAL: Must flush writer at end
//
// Memory efficiency:
//   - Scanner buffer: ~64KB
//   - Writer buffer: ~4KB
//   - Total memory: O(1) regardless of file size
//   - Processes files of any size efficiently
//
// Example:
//   Input file: "ERROR: failed\nINFO: success\nERROR: timeout\n"
//   Predicate: func(s string) bool { return strings.Contains(s, "ERROR") }
//   Output file: "ERROR: failed\nERROR: timeout\n"
//   Returns: 2, nil
//
// Common mistake:
//   - Forgetting writer.Flush() → data stays in buffer, not written to disk
//   - Use defer writer.Flush() for safety
//
// func FilterLines(input io.Reader, output io.Writer, predicate func(string) bool) (int, error) {
//     scanner := bufio.NewScanner(input)
//     writer := bufio.NewWriter(output)
//     defer writer.Flush() // CRITICAL: Flush buffered data
//
//     count := 0
//
//     for scanner.Scan() {
//         line := scanner.Text()
//
//         if predicate(line) {
//             // Write line with newline
//             _, err := writer.WriteString(line + "\n")
//             if err != nil {
//                 return count, err
//             }
//             count++
//         }
//     }
//
//     if err := scanner.Err(); err != nil {
//         return count, err
//     }
//
//     return count, nil
// }

func FilterLines(input io.Reader, output io.Writer, predicate func(string) bool) (int, error) {
	// TODO: Implement FilterLines
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// ============================================================================
// Exercise 3: WordFrequency
// ============================================================================

// WordFrequency counts word frequencies in the input.
// TODO: Implement word frequency counter using Scanner with custom split function.
//
// Parameters:
//   - reader: io.Reader to read from
//
// Returns:
//   - map[string]int: Word frequencies (lowercase words)
//   - error: Non-nil if reading fails
//
// Behavior:
//   - Split input by whitespace (spaces, tabs, newlines)
//   - Convert words to lowercase for case-insensitive counting
//   - Count occurrences of each word
//   - Return map of word → count
//
// Scanner split functions:
//   - Default: ScanLines (split by '\n')
//   - ScanWords: Split by whitespace (what we need)
//   - ScanRunes: Split by individual runes
//   - ScanBytes: Split by individual bytes
//   - Custom: Write your own bufio.SplitFunc
//
// Example:
//   Input: "Go is great. Go is fast."
//   Output: map[string]int{"go": 2, "is": 2, "great.": 1, "fast.": 1}
//
// Implementation:
//   scanner := bufio.NewScanner(reader)
//   scanner.Split(bufio.ScanWords) // Change from lines to words
//   freq := make(map[string]int)
//
//   for scanner.Scan() {
//       word := strings.ToLower(scanner.Text())
//       freq[word]++
//   }
//
// Memory:
//   - Map grows as new words found: O(unique words)
//   - For most documents, unique words << total words
//   - Scanner buffer: O(1)
//
// func WordFrequency(reader io.Reader) (map[string]int, error) {
//     scanner := bufio.NewScanner(reader)
//     scanner.Split(bufio.ScanWords) // Split on whitespace
//
//     freq := make(map[string]int)
//
//     for scanner.Scan() {
//         word := strings.ToLower(scanner.Text())
//         freq[word]++
//     }
//
//     if err := scanner.Err(); err != nil {
//         return nil, err
//     }
//
//     return freq, nil
// }

func WordFrequency(reader io.Reader) (map[string]int, error) {
	// TODO: Implement WordFrequency
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// ============================================================================
// Exercise 4: TransformFile
// ============================================================================

// TransformFile reads from input, applies transform to each line, writes to output.
// TODO: Implement streaming transformation pipeline.
//
// Parameters:
//   - input: io.Reader to read from
//   - output: io.Writer to write to
//   - transform: Function to transform each line
//
// Returns:
//   - error: Non-nil if reading/writing fails
//
// Behavior:
//   - Read input line-by-line
//   - Apply transform function to each line
//   - Write transformed line to output (with newline)
//   - Flush buffer at end
//
// Example:
//   Input: "hello\nworld\n"
//   Transform: strings.ToUpper
//   Output: "HELLO\nWORLD\n"
//
// This is a simple streaming pipeline:
//   input → scanner → transform → writer → output
//
// func TransformFile(input io.Reader, output io.Writer, transform func(string) string) error {
//     scanner := bufio.NewScanner(input)
//     writer := bufio.NewWriter(output)
//     defer writer.Flush()
//
//     for scanner.Scan() {
//         line := scanner.Text()
//         transformed := transform(line)
//
//         _, err := writer.WriteString(transformed + "\n")
//         if err != nil {
//             return err
//         }
//     }
//
//     return scanner.Err()
// }

func TransformFile(input io.Reader, output io.Writer, transform func(string) string) error {
	// TODO: Implement TransformFile
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// ============================================================================
// Exercise 5: ReadChunks
// ============================================================================

// ReadChunks reads data in fixed-size chunks and calls callback for each chunk.
// TODO: Implement chunked reading for binary data or custom parsing.
//
// Parameters:
//   - reader: io.Reader to read from
//   - chunkSize: Size of each chunk in bytes (e.g., 1024 for 1KB chunks)
//   - callback: Function called with each chunk (may be smaller at EOF)
//
// Returns:
//   - int: Total bytes read
//   - error: Non-nil if reading fails
//
// Behavior:
//   - Create buffer of size chunkSize
//   - Repeatedly call reader.Read(buffer)
//   - Call callback with buffer[:n] (n = bytes actually read)
//   - Continue until EOF (io.EOF is expected, not an error)
//   - Return total bytes read
//
// Why this is useful:
//   - Processing binary files (images, videos, archives)
//   - Custom parsing logic that doesn't fit line/word model
//   - Progress reporting (callback can update progress bar)
//   - Memory-constrained environments (limit chunk size)
//
// Memory:
//   - Allocates one buffer of size chunkSize: O(chunkSize)
//   - Reuses buffer for all reads (no allocation per chunk)
//   - Total memory: O(1) regardless of file size
//
// Example:
//   total, err := ReadChunks(file, 1024, func(chunk []byte) {
//       fmt.Printf("Read %d bytes: %x\n", len(chunk), chunk)
//   })
//
// func ReadChunks(reader io.Reader, chunkSize int, callback func([]byte)) (int, error) {
//     buffer := make([]byte, chunkSize)
//     totalBytes := 0
//
//     for {
//         // Read up to chunkSize bytes
//         n, err := reader.Read(buffer)
//
//         if n > 0 {
//             // Call callback with actual bytes read
//             // IMPORTANT: Pass buffer[:n], not buffer
//             callback(buffer[:n])
//             totalBytes += n
//         }
//
//         if err == io.EOF {
//             // End of file is expected, not an error
//             break
//         }
//         if err != nil {
//             // Actual error occurred
//             return totalBytes, err
//         }
//     }
//
//     return totalBytes, nil
// }
//
// Alternative using bufio.Reader (better performance):
// func ReadChunks(reader io.Reader, chunkSize int, callback func([]byte)) (int, error) {
//     bufferedReader := bufio.NewReaderSize(reader, chunkSize*2)
//     buffer := make([]byte, chunkSize)
//     totalBytes := 0
//
//     for {
//         n, err := bufferedReader.Read(buffer)
//         if n > 0 {
//             callback(buffer[:n])
//             totalBytes += n
//         }
//         if err == io.EOF {
//             break
//         }
//         if err != nil {
//             return totalBytes, err
//         }
//     }
//
//     return totalBytes, nil
// }

func ReadChunks(reader io.Reader, chunkSize int, callback func([]byte)) (int, error) {
	// TODO: Implement ReadChunks
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// After implementing all functions:
// - Run: go test -v ./...
// - Test with large files to verify memory efficiency
// - Compare with solution.go to see detailed implementations
// - Try different buffer sizes to see performance impact
//
// Key takeaways:
// 1. Always use buffered I/O for file operations
// 2. Scanner is perfect for line/word-based processing
// 3. Reader/Writer give more control when needed
// 4. Always Flush() buffered writers
// 5. Handle io.EOF correctly (it's expected, not an error)
// 6. Streaming keeps memory usage constant
