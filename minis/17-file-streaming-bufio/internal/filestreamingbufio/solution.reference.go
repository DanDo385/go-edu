//go:build reference

package filestreamingbufio

/*
Reference Solution
==================

This file is the canonical reference for this exercise. It keeps failure paths
explicit when an operation can fail, so callers can decide how to handle
errors at API boundaries.

Read this alongside exercise.go and the tests to understand the intended data
flow, ownership boundaries, and invariants that keep behavior deterministic.

Teaching notes:
- Memory/ownership: make copies when returning mutable data that should not
  alias internal state; share references only when aliasing is intentional.
- Invariants: establish assumptions close to construction, and rely on them in
  smaller helper functions to keep logic easy to audit.
- Error surfaces: prefer explicit returns over hidden panics so learners can
  reason about control flow in production-style code.
*/

/*
Project 17: File Streaming with bufio - Solutions

This file contains complete solutions with extensive debugging support.

Key Go Concepts Demonstrated:
1. Buffered I/O for efficient file operations
2. Scanner for line-by-line reading
3. Writer with buffering and flushing
4. Custom split functions
5. Memory-efficient streaming

DEBUGGING GUIDE:
- BREAKPOINT comments mark ideal breakpoint locations
- DEBUG comments explain what to observe in debugger
- Use Step Over (F10) to execute line by line
- Use Step Into (F11) to enter function calls
- Watch panel shows variable values in real-time
*/

import (
	"bufio"
	"io"
	"strings"
)

// CountLines counts the number of lines in a file.
func CountLines(reader io.Reader) (int, error) {
	// BREAKPOINT: Set breakpoint to observe scanner creation
	// DEBUG: Watch scanner being initialized with default buffer (64KB)
	// DEBUG: Scanner uses ScanLines split function by default
	scanner := bufio.NewScanner(reader)

	// BREAKPOINT: Set breakpoint before loop to watch line counting
	// DEBUG: Watch count variable incrementing
	count := 0

	// BREAKPOINT: Set breakpoint in loop to observe line-by-line reading
	// DEBUG: scanner.Scan() advances to next line
	// DEBUG: Returns false when EOF reached or error occurs
	for scanner.Scan() {
		// DEBUG: Each Scan() call reads one line into internal buffer
		// DEBUG: We don't need scanner.Text() since just counting
		count++
	}

	// BREAKPOINT: Set breakpoint to check for scanner errors
	// DEBUG: scanner.Err() returns nil for normal EOF
	// DEBUG: Non-nil means actual read error occurred
	if err := scanner.Err(); err != nil {
		return count, err
	}

	return count, nil
}

// FilterLines reads from input, writes lines matching the predicate to output.
func FilterLines(input io.Reader, output io.Writer, predicate func(string) bool) (int, error) {
	// BREAKPOINT: Set breakpoint to observe scanner creation
	// DEBUG: Scanner for efficient line-by-line reading
	scanner := bufio.NewScanner(input)

	// BREAKPOINT: Set breakpoint to observe buffered writer creation
	// DEBUG: Writer batches writes to reduce system calls
	// DEBUG: Default buffer size is 4KB
	writer := bufio.NewWriter(output)
	// DEBUG: defer ensures Flush() called even if function panics
	// DEBUG: Without Flush(), buffered data stays in memory, not written
	defer writer.Flush()

	// BREAKPOINT: Set breakpoint to watch filtered line count
	// DEBUG: Only lines matching predicate increment count
	count := 0

	// BREAKPOINT: Set breakpoint in loop to observe line filtering
	// DEBUG: Watch scanner.Text() returning current line
	for scanner.Scan() {
		line := scanner.Text()

		// BREAKPOINT: Set breakpoint to observe predicate evaluation
		// DEBUG: Watch predicate(line) return true or false
		// DEBUG: Only matching lines are written
		if predicate(line) {
			// BREAKPOINT: Set breakpoint before write
			// DEBUG: WriteString writes to buffer, not directly to output
			// DEBUG: Watch count incrementing only for matching lines
			_, err := writer.WriteString(line + "\n")
			if err != nil {
				return count, err
			}
			count++
		}
	}

	// BREAKPOINT: Set breakpoint to check scanner error
	// DEBUG: Check if scanner stopped due to error vs normal EOF
	if err := scanner.Err(); err != nil {
		return count, err
	}

	// DEBUG: Flush happens automatically via defer
	// DEBUG: All buffered data written to output before return
	return count, nil
}

// WordFrequency counts word frequencies in the input.
func WordFrequency(reader io.Reader) (map[string]int, error) {
	// BREAKPOINT: Set breakpoint to observe scanner creation
	// DEBUG: Scanner initialized with input reader
	scanner := bufio.NewScanner(reader)

	// BREAKPOINT: Set breakpoint to observe custom split function
	// DEBUG: scanner.Split changes from ScanLines to ScanWords
	// DEBUG: ScanWords splits on any whitespace (spaces, tabs, newlines)
	scanner.Split(bufio.ScanWords)

	// BREAKPOINT: Set breakpoint to watch frequency map creation
	// DEBUG: Map starts empty, grows as new words found
	freq := make(map[string]int)

	// BREAKPOINT: Set breakpoint in loop to watch word processing
	// DEBUG: scanner.Text() returns one word per iteration
	// DEBUG: Watch words being converted to lowercase
	for scanner.Scan() {
		// DEBUG: ToLower ensures case-insensitive counting
		// DEBUG: "Hello" and "hello" count as same word
		word := strings.ToLower(scanner.Text())
		// BREAKPOINT: Watch map entry increment
		// DEBUG: If word not in map, starts at 0, then increments to 1
		// DEBUG: If word exists, increments existing count
		freq[word]++
	}

	// BREAKPOINT: Set breakpoint to check scanner error
	// DEBUG: Verify no errors occurred during scanning
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// BREAKPOINT: Set breakpoint to observe final frequency map
	// DEBUG: Watch freq map with all word counts
	return freq, nil
}

// TransformFile reads from input, applies transform to each line, writes to output.
func TransformFile(input io.Reader, output io.Writer, transform func(string) string) error {
	// BREAKPOINT: Set breakpoint to observe scanner and writer creation
	// DEBUG: Scanner for reading, writer for writing
	scanner := bufio.NewScanner(input)
	writer := bufio.NewWriter(output)
	// DEBUG: defer Flush() ensures all buffered data written
	defer writer.Flush()

	// BREAKPOINT: Set breakpoint in loop to watch transformation
	// DEBUG: Each line read, transformed, and written
	for scanner.Scan() {
		// DEBUG: scanner.Text() gets current line
		line := scanner.Text()

		// BREAKPOINT: Set breakpoint to observe transformation
		// DEBUG: Watch transform function being applied
		// DEBUG: Step into transform() to see operation (e.g., ToUpper)
		transformed := transform(line)

		// BREAKPOINT: Set breakpoint before write
		// DEBUG: Watch transformed line being written with newline
		_, err := writer.WriteString(transformed + "\n")
		if err != nil {
			return err
		}
	}

	// BREAKPOINT: Set breakpoint to check scanner completion
	// DEBUG: Returns scanner error or nil
	return scanner.Err()
}

// ReadChunks reads data in fixed-size chunks and calls the callback for each chunk.
func ReadChunks(reader io.Reader, chunkSize int, callback func([]byte)) (int, error) {
	// BREAKPOINT: Set breakpoint to observe buffer allocation
	// DEBUG: Watch buffer being created with chunkSize capacity
	// DEBUG: This buffer is reused for all reads (memory efficient)
	buffer := make([]byte, chunkSize)
	totalBytes := 0

	// BREAKPOINT: Set breakpoint in loop to watch chunk reading
	// DEBUG: Loop continues until EOF or error
	for {
		// BREAKPOINT: Set breakpoint on Read() call
		// DEBUG: Read fills buffer up to chunkSize bytes
		// DEBUG: n is actual bytes read (may be less than chunkSize)
		// DEBUG: Watch n varying as we approach end of file
		n, err := reader.Read(buffer)

		// BREAKPOINT: Set breakpoint to process read data
		// DEBUG: Only process if n > 0 (some data was read)
		if n > 0 {
			// DEBUG: buffer[:n] is crucial - only pass filled portion
			// DEBUG: buffer may not be completely filled on last read
			callback(buffer[:n])
			// DEBUG: Watch totalBytes accumulating
			totalBytes += n
		}

		// BREAKPOINT: Set breakpoint to check error handling
		// DEBUG: io.EOF is expected, not an error condition
		// DEBUG: Watch err value - io.EOF vs actual errors
		if err == io.EOF {
			// DEBUG: Reached end of file normally
			break
		}
		if err != nil {
			// DEBUG: Actual error occurred during reading
			return totalBytes, err
		}
	}

	// BREAKPOINT: Set breakpoint to observe completion
	// DEBUG: Watch totalBytes - should match file size
	return totalBytes, nil
}

// ReadChunksBuffered is an advanced implementation using bufio.Reader for better performance.
func ReadChunksBuffered(reader io.Reader, chunkSize int, callback func([]byte)) (int, error) {
	// BREAKPOINT: Set breakpoint to observe buffered reader creation
	// DEBUG: bufio.Reader adds internal buffer (chunkSize*2)
	// DEBUG: This reduces system calls by batching reads
	bufferedReader := bufio.NewReaderSize(reader, chunkSize*2)

	// BREAKPOINT: Set breakpoint to watch buffer creation
	// DEBUG: Same as ReadChunks, but reading from bufferedReader
	buffer := make([]byte, chunkSize)
	totalBytes := 0

	// BREAKPOINT: Set breakpoint in loop
	// DEBUG: bufferedReader.Read may return data from internal buffer
	// DEBUG: Fewer actual I/O operations than unbuffered version
	for {
		// DEBUG: Read from buffered reader (may not hit disk every time)
		n, err := bufferedReader.Read(buffer)

		if n > 0 {
			// DEBUG: Process chunk same as unbuffered version
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

/*
Common Implementation Patterns:

1. Scanner Pattern (line/word processing):
   - Use bufio.Scanner for text processing
   - Default 64KB buffer, adjustable with scanner.Buffer()
   - Built-in split functions: ScanLines, ScanWords, ScanRunes, ScanBytes
   - Custom split functions possible for special parsing

2. Buffered Writer Pattern:
   - Use bufio.Writer to batch writes
   - Reduces system calls dramatically
   - MUST call Flush() to ensure all data written
   - Use defer writer.Flush() for safety

3. Chunked Reading Pattern:
   - Allocate fixed buffer, reuse for all reads
   - Read() returns actual bytes read (may be less than buffer size)
   - Use buffer[:n] to get only filled portion
   - io.EOF is normal completion, not error

Critical Implementation Details:

- Scanner.Text() allocates new string each call
- For huge files, consider custom buffer reuse
- Buffered I/O reduces syscalls: 1M unbuffered reads vs 250 buffered
- Always check scanner.Err() after loop (Scan() returns false on error or EOF)
- Flush() writes buffered data; without it, data may be lost

Debugging Tips:

- Watch buffer sizes and capacities
- Monitor scanner internal buffer state
- Track bytes read vs bytes processed
- Observe Flush() timing (should be at end)
- Compare buffered vs unbuffered performance
*/
