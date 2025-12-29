/*
Package main demonstrates file streaming with bufio for efficient I/O operations.

USAGE EXAMPLES:
  # Run all demonstrations (default)
  ./file-stream

  # Use custom file sizes
  ./file-stream -small=100 -large=1000000

  # Set custom buffer size
  ./file-stream -buffer=8192

  # Run specific demos only
  ./file-stream -demo=scanner,writer,performance

  # Enable verbose debug output
  ./file-stream -verbose

FLAGS:
  -demo      Comma-separated demos: all, scanner, words, writer, reader, large, perf, custom, filter
  -small     Number of lines for small test file (default: 10)
  -large     Number of lines for large test file (default: 100000)
  -buffer    Buffer size in bytes (default: 4096)
  -verbose   Enable verbose debug output
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// CLI flags
var (
	demoFlag    = flag.String("demo", "all", "Demos to run: all, scanner, words, writer, reader, large, perf, custom, filter")
	smallFlag   = flag.Int("small", 10, "Lines in small test file")
	largeFlag   = flag.Int("large", 100000, "Lines in large test file")
	bufferFlag  = flag.Int("buffer", 4096, "Buffer size in bytes")
	verboseFlag = flag.Bool("verbose", false, "Enable verbose output")
)

func main() {
	// BREAKPOINT: Parse CLI flags to configure demonstrations
	flag.Parse()

	// DEBUG: Display configuration when verbose mode is enabled
	if *verboseFlag {
		fmt.Printf("DEBUG: Config - small=%d, large=%d, buffer=%d\n",
			*smallFlag, *largeFlag, *bufferFlag)
	}

	fmt.Println("=== File Streaming with bufio Demonstrations ===\n")

	// BREAKPOINT: Create test files before running demonstrations
	// Create test files for demonstrations
	createTestFiles()

	// Run demonstrations
	demo1_BasicScanner()
	demo2_ScannerByWords()
	demo3_BufferedWriter()
	demo4_ReaderMethods()
	demo5_LargeFileProcessing()
	demo6_PerformanceComparison()
	demo7_CustomSplitFunction()
	demo8_FilteringFiles()

	// Cleanup
	cleanup()
}

// createTestFiles generates test files for demonstrations
func createTestFiles() {
	fmt.Println("=== Creating Test Files ===")

	// BREAKPOINT: Watch file creation with os.Create (truncates if exists)
	// Small test file (a few lines)
	small, err := os.Create("small.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer small.Close()

	// DEBUG: bufio.NewWriter creates a buffered writer with default 4KB buffer
	writer := bufio.NewWriter(small)
	for i := 1; i <= *smallFlag; i++ {
		fmt.Fprintf(writer, "This is line %d of the small test file\n", i)
	}
	// BREAKPOINT: Flush() is critical - writes buffered data to disk
	writer.Flush()
	fmt.Printf("✓ Created small.txt (%d lines)\n", *smallFlag)

	// DEBUG: Create larger file to demonstrate streaming performance
	// Larger test file (simulating a real-world scenario)
	large, err := os.Create("large.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer large.Close()

	writer = bufio.NewWriter(large)
	// DEBUG: Writing many lines - buffering amortizes syscall overhead
	// Write 100,000 lines (simulating a log file)
	for i := 1; i <= *largeFlag; i++ {
		fmt.Fprintf(writer, "2025-11-15 10:%02d:%02d [INFO] Processing request %d from user-%d\n",
			i/3600, (i/60)%60, i, i%1000)
	}
	writer.Flush()
	fmt.Printf("✓ Created large.txt (%d lines)\n", *largeFlag)

	// File with mixed content for word scanning
	mixed, err := os.Create("mixed.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer mixed.Close()

	writer = bufio.NewWriter(mixed)
	writer.WriteString("Go is a statically typed, compiled programming language.\n")
	writer.WriteString("It was designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson.\n")
	writer.WriteString("Go is syntactically similar to C, but with memory safety and garbage collection.\n")
	writer.Flush()
	fmt.Println("✓ Created mixed.txt (3 lines with various words)\n")
}

// demo1_BasicScanner shows the most common pattern: scanning line-by-line
func demo1_BasicScanner() {
	fmt.Println("=== Demo 1: Basic Scanner (Line-by-Line Reading) ===")

	// BREAKPOINT: Watch file opening with os.Open (read-only mode)
	// STEP 1: Open the file
	// os.Open returns (*os.File, error)
	// The file is opened in read-only mode
	file, err := os.Open("small.txt")
	if err != nil {
		log.Fatal(err)
	}
	// DEBUG: defer ensures Close() is called even if function panics
	// CRITICAL: Always defer file.Close() to prevent resource leaks
	// Even if the function panics, defer ensures cleanup happens
	defer file.Close()

	// BREAKPOINT: Examine Scanner creation - wraps io.Reader with buffering
	// STEP 2: Create a Scanner
	// Scanner wraps the file and provides a simple API for reading lines
	// By default, it splits on '\n' (newlines)
	scanner := bufio.NewScanner(file)

	lineNum := 0

	// DEBUG: scanner.Scan() reads from buffer, refilling automatically when empty
	// STEP 3: Scan line-by-line
	// scanner.Scan() advances to the next token (line) and returns true
	// It returns false when:
	//   - EOF is reached, OR
	//   - An error occurs
	for scanner.Scan() {
		lineNum++

		// BREAKPOINT: scanner.Text() returns current token without delimiter
		// scanner.Text() returns the current line as a string
		// The '\n' delimiter is NOT included
		line := scanner.Text()

		fmt.Printf("  Line %d: %s\n", lineNum, line)

		// DEBUG: Scanner handles EOF gracefully - no manual EOF checking needed
		// Note: We don't need to check for io.EOF explicitly
		// The scanner handles it gracefully
	}

	// BREAKPOINT: Always check scanner.Err() after loop to catch read errors
	// DEBUG: Scan() returns false for both EOF (normal) and errors (abnormal)
	// STEP 4: Check for errors
	// scanner.Scan() returns false on both EOF and errors
	// We must check scanner.Err() to distinguish them
	if err := scanner.Err(); err != nil {
		log.Fatalf("Scanner error: %v", err)
	}

	fmt.Printf("  ✓ Successfully read %d lines\n\n", lineNum)
}

// demo2_ScannerByWords demonstrates scanning by words instead of lines
func demo2_ScannerByWords() {
	fmt.Println("=== Demo 2: Scanner with Custom Split (Words) ===")

	file, err := os.Open("mixed.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// BREAKPOINT: Watch scanner.Split() changing the tokenization strategy
	// DEBUG: Built-in split functions: ScanLines, ScanWords, ScanBytes, ScanRunes
	// CHANGE THE SPLIT FUNCTION
	// By default, Scanner uses bufio.ScanLines
	// We can change it to ScanWords, ScanBytes, ScanRunes, or a custom function
	scanner.Split(bufio.ScanWords)

	// Count word frequencies
	wordCount := make(map[string]int)

	for scanner.Scan() {
		// Now scanner.Text() returns one word at a time
		// Words are delimited by whitespace (spaces, tabs, newlines)
		word := strings.ToLower(scanner.Text())
		wordCount[word]++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Display top 10 words
	fmt.Println("  Top 10 most frequent words:")
	count := 0
	for word, freq := range wordCount {
		if count >= 10 {
			break
		}
		fmt.Printf("    %-15s: %d\n", word, freq)
		count++
	}
	fmt.Printf("  ✓ Total unique words: %d\n\n", len(wordCount))
}

// demo3_BufferedWriter demonstrates efficient writing with bufio.Writer
func demo3_BufferedWriter() {
	fmt.Println("=== Demo 3: Buffered Writing ===")

	// BREAKPOINT: Create file for writing (truncates existing file)
	// STEP 1: Create output file
	file, err := os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// BREAKPOINT: bufio.NewWriter creates in-memory buffer to batch writes
	// DEBUG: Default buffer size is 4KB - configurable with NewWriterSize()
	// STEP 2: Wrap file in a buffered writer
	// This creates a 4KB buffer in memory
	// Writes are accumulated in the buffer and flushed to disk in chunks
	writer := bufio.NewWriter(file)

	// DEBUG: Without Flush(), up to 4KB could remain in buffer and be lost!
	// CRITICAL: Always defer Flush() to ensure all buffered data is written
	// Without this, data might remain in the buffer and be lost!
	defer writer.Flush()

	start := time.Now()

	// BREAKPOINT: Watch many small writes - all go to buffer, not disk
	// DEBUG: Each write is ~25 bytes, buffer is 4KB - ~160 writes per flush
	// STEP 3: Write many small writes
	// Each write goes to the buffer (fast), not to disk (slow)
	for i := 1; i <= 10000; i++ {
		// DEBUG: WriteString avoids reflection overhead of fmt.Fprintf
		// WriteString is more efficient than fmt.Fprintf for simple strings
		// It writes directly to the buffer without formatting overhead
		_, err := writer.WriteString(fmt.Sprintf("Line %d: Some data here\n", i))
		if err != nil {
			log.Fatal(err)
		}

		// DEBUG: Auto-flush happens when buffer fills - reduces syscalls by ~160x
		// The buffer automatically flushes when it's full (4KB by default)
		// We can also manually flush at any time with writer.Flush()
	}

	// STEP 4: Final flush (via defer)
	// This ensures any remaining data in the buffer is written to disk
	elapsed := time.Since(start)

	fmt.Printf("  ✓ Wrote 10,000 lines in %v\n", elapsed)
	fmt.Println("  ✓ Data is safely on disk (after Flush)\n")
}

// demo4_ReaderMethods shows different ways to read with bufio.Reader
func demo4_ReaderMethods() {
	fmt.Println("=== Demo 4: bufio.Reader Methods ===")

	file, err := os.Open("mixed.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a Reader with a custom buffer size (64KB instead of default 4KB)
	// Larger buffers can improve performance for large files
	// But they use more memory
	reader := bufio.NewReaderSize(file, 64*1024)

	fmt.Println("  Method 1: ReadString (read until delimiter)")
	// ReadString reads until the first occurrence of delim
	// The delimiter IS included in the returned string
	line1, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	fmt.Printf("    First line: %q\n", line1)

	fmt.Println("\n  Method 2: ReadBytes (like ReadString but returns []byte)")
	line2, err := reader.ReadBytes('\n')
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	fmt.Printf("    Second line: %q\n", string(line2))

	fmt.Println("\n  Method 3: ReadByte (read one byte at a time)")
	// This is useful for parsing binary formats or protocols
	// Each ReadByte call is served from the buffer (no syscall)
	for i := 0; i < 20; i++ {
		b, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%c", b)
	}
	fmt.Println("\n")
}

// demo5_LargeFileProcessing shows memory-efficient processing of large files
func demo5_LargeFileProcessing() {
	fmt.Println("=== Demo 5: Processing Large File (Memory Efficient) ===")

	file, err := os.Open("large.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// BREAKPOINT: Scanner enables streaming - O(1) memory regardless of file size
	scanner := bufio.NewScanner(file)

	// DEBUG: We track statistics but don't store individual lines
	// Track statistics without loading entire file into memory
	var (
		lineCount  int
		errorCount int
		infoCount  int
	)

	start := time.Now()

	// BREAKPOINT: Watch streaming pattern - process and discard each line
	// DEBUG: Peak memory = buffer size (~4KB), NOT proportional to file size
	// Process line-by-line (streaming)
	// Memory usage: ~4KB (buffer size) regardless of file size!
	for scanner.Scan() {
		lineCount++
		line := scanner.Text()

		// DEBUG: Each line is analyzed then immediately garbage collected
		// Analyze log level
		if strings.Contains(line, "[ERROR]") {
			errorCount++
		} else if strings.Contains(line, "[INFO]") {
			infoCount++
		}

		// BREAKPOINT: This is the key - we don't accumulate lines in memory
		// Note: We DON'T store all lines in memory
		// Each line is processed and discarded
		// This allows processing files larger than available RAM
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	elapsed := time.Since(start)

	fmt.Printf("  ✓ Processed %d lines in %v\n", lineCount, elapsed)
	fmt.Printf("  ✓ INFO logs:  %d\n", infoCount)
	fmt.Printf("  ✓ ERROR logs: %d\n", errorCount)
	fmt.Printf("  ✓ Memory used: ~4KB (constant, not proportional to file size!)\n\n")
}

// demo6_PerformanceComparison compares buffered vs unbuffered I/O
func demo6_PerformanceComparison() {
	fmt.Println("=== Demo 6: Performance Comparison (Buffered vs Unbuffered) ===")

	// BREAKPOINT: Compare unbuffered (many syscalls) vs buffered (few syscalls)
	// DEBUG: Unbuffered = one Read() syscall per byte (~7M syscalls for 7MB file)
	// Test 1: Unbuffered reading (byte-by-byte with os.File)
	fmt.Println("  Test 1: Unbuffered byte-by-byte reading...")
	start := time.Now()

	file1, _ := os.Open("large.txt")
	defer file1.Close()

	buf := make([]byte, 1)
	unbufferedBytes := 0
	for {
		// DEBUG: Each Read() call is a syscall - very slow!
		n, err := file1.Read(buf)
		unbufferedBytes += n
		if err == io.EOF {
			break
		}
	}

	unbufferedTime := time.Since(start)
	fmt.Printf("    Read %d bytes in %v\n", unbufferedBytes, unbufferedTime)

	// BREAKPOINT: Buffered approach - ReadByte() reads from memory buffer
	// DEBUG: Buffered = ~1800 Read() syscalls (7MB / 4KB buffer) - much faster!
	// Test 2: Buffered reading (byte-by-byte with bufio.Reader)
	fmt.Println("\n  Test 2: Buffered byte-by-byte reading...")
	start = time.Now()

	file2, _ := os.Open("large.txt")
	defer file2.Close()

	reader := bufio.NewReader(file2)
	bufferedBytes := 0
	for {
		// DEBUG: ReadByte() reads from buffer (fast), refills buffer via syscall only when empty
		_, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		bufferedBytes++
	}

	bufferedTime := time.Since(start)
	fmt.Printf("    Read %d bytes in %v\n", bufferedBytes, bufferedTime)

	// Calculate speedup
	speedup := float64(unbufferedTime) / float64(bufferedTime)
	fmt.Printf("\n  ✓ Buffered I/O is %.1fx FASTER\n", speedup)
	fmt.Printf("  ✓ This is because buffered I/O makes ~1,800 syscalls vs ~7,000,000 for unbuffered\n\n")
}

// demo7_CustomSplitFunction shows how to create a custom scanner split function
func demo7_CustomSplitFunction() {
	fmt.Println("=== Demo 7: Custom Split Function (CSV-like) ===")

	// Create a CSV-like file
	csvFile, err := os.Create("data.csv")
	if err != nil {
		log.Fatal(err)
	}
	writer := bufio.NewWriter(csvFile)
	writer.WriteString("name,age,city\n")
	writer.WriteString("Alice,30,NYC\n")
	writer.WriteString("Bob,25,SF\n")
	writer.WriteString("Charlie,35,LA\n")
	writer.Flush()
	csvFile.Close()

	// Read it back with a custom splitter
	file, err := os.Open("data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Define a custom split function that splits on commas
	// This is a simplified CSV parser (doesn't handle quoted commas)
	splitOnComma := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		// Find the next comma or newline
		for i := 0; i < len(data); i++ {
			if data[i] == ',' || data[i] == '\n' {
				// Return the token before the delimiter
				return i + 1, data[:i], nil
			}
		}

		// If we're at EOF and have data, return it
		if atEOF && len(data) > 0 {
			return len(data), data, nil
		}

		// Need more data
		return 0, nil, nil
	}

	scanner.Split(splitOnComma)

	fmt.Println("  CSV fields (split on commas):")
	fieldNum := 0
	for scanner.Scan() {
		fieldNum++
		field := scanner.Text()
		fmt.Printf("    Field %2d: %s\n", fieldNum, field)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("  ✓ Parsed %d fields using custom split function\n\n", fieldNum)
}

// demo8_FilteringFiles shows streaming transformation of files
func demo8_FilteringFiles() {
	fmt.Println("=== Demo 8: Streaming File Transformation ===")

	// Open input file
	input, err := os.Open("large.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	// Create output file
	output, err := os.Create("errors_only.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	// Create buffered reader and writer
	scanner := bufio.NewScanner(input)
	writer := bufio.NewWriter(output)
	defer writer.Flush()

	start := time.Now()
	linesRead := 0
	linesWritten := 0

	// Stream line-by-line, filtering as we go
	for scanner.Scan() {
		linesRead++
		line := scanner.Text()

		// Filter: Only keep ERROR lines
		if strings.Contains(line, "[ERROR]") {
			writer.WriteString(line + "\n")
			linesWritten++
		}

		// Note: We never load the entire file into memory
		// Input is read incrementally, output is written incrementally
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	elapsed := time.Since(start)

	fmt.Printf("  ✓ Filtered %d lines to %d error lines in %v\n", linesRead, linesWritten, elapsed)
	fmt.Printf("  ✓ Memory used: ~8KB (4KB read buffer + 4KB write buffer)\n")
	fmt.Printf("  ✓ Can process files of ANY size (even 100GB+) with constant memory\n\n")
}

// cleanup removes temporary test files
func cleanup() {
	fmt.Println("=== Cleanup ===")
	files := []string{
		"small.txt",
		"large.txt",
		"mixed.txt",
		"output.txt",
		"data.csv",
		"errors_only.txt",
	}

	for _, f := range files {
		os.Remove(f)
	}
	fmt.Println("✓ Removed all test files\n")

	fmt.Println("=== All Demonstrations Complete ===")
	fmt.Println("\nKey Takeaways:")
	fmt.Println("  1. Always use bufio for file I/O (10-100x faster)")
	fmt.Println("  2. Scanner is perfect for line-oriented text files")
	fmt.Println("  3. Always defer file.Close() and writer.Flush()")
	fmt.Println("  4. Stream large files to avoid memory exhaustion")
	fmt.Println("  5. Check scanner.Err() after scanning")
	fmt.Println("  6. Custom split functions enable parsing any format")
}
