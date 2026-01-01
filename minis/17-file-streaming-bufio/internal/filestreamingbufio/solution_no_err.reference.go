//go:build reference

package filestreamingbufio

import (
	"bufio"
	"io"
	"strings"
)

/*
Core Logic Without Error Handling
==================================

This file demonstrates the core algorithm without error handling complexity.
It's designed to help you understand the fundamental logic flow.

KEY CONCEPTS:
- Focus on the algorithm, not error handling
- Simplified implementations for learning
- Use this to understand core logic before adding error handling

DEBUGGING:
- Set breakpoints at function entry points
- Step through the pure logic without error handling distractions
- Focus on understanding the algorithm
- Watch Variables panel to see data transformations

See /RUN_DEBUG.md for comprehensive debugging guide.
*/

// ============================================================================
// Core Exercise 1: Count Lines (No Error Handling)
// ============================================================================

// CoreCountLines counts lines without error handling
// TODO: Implement basic line counting
//
// Algorithm:
// 1. Create scanner
// 2. Loop while scanner.Scan() returns true
// 3. Increment counter each iteration
// 4. Return count
//
// func CoreCountLines(reader io.Reader) int {
//     // BREAKPOINT: Watch scanner creation
//     scanner := bufio.NewScanner(reader)
//     count := 0
//
//     // BREAKPOINT: Watch line-by-line processing
//     for scanner.Scan() {
//         count++
//     }
//
//     return count
// }

func CoreCountLines(reader io.Reader) int {
	// TODO: implement line counting
	return 0
}

// ============================================================================
// Core Exercise 2: Filter Lines (No Error Handling)
// ============================================================================

// CoreFilterLines filters lines without error handling
// TODO: Implement line filtering with predicate
//
// Algorithm:
// 1. Create scanner for reading
// 2. Create writer for writing
// 3. Loop through lines
// 4. If predicate matches, write line
// 5. Flush writer
// 6. Return count
//
// func CoreFilterLines(input io.Reader, output io.Writer, predicate func(string) bool) int {
//     // BREAKPOINT: Watch scanner and writer setup
//     scanner := bufio.NewScanner(input)
//     writer := bufio.NewWriter(output)
//     defer writer.Flush()
//
//     count := 0
//
//     // BREAKPOINT: Watch filtering logic
//     for scanner.Scan() {
//         line := scanner.Text()
//         if predicate(line) {
//             writer.WriteString(line + "\n")
//             count++
//         }
//     }
//
//     return count
// }

func CoreFilterLines(input io.Reader, output io.Writer, predicate func(string) bool) int {
	// TODO: implement line filtering
	return 0
}

// ============================================================================
// Core Exercise 3: Word Frequency (No Error Handling)
// ============================================================================

// CoreWordFrequency counts word frequencies without error handling
// TODO: Implement word frequency counter
//
// Algorithm:
// 1. Create scanner
// 2. Set split function to ScanWords
// 3. Create frequency map
// 4. Loop through words
// 5. Convert to lowercase
// 6. Increment count in map
// 7. Return map
//
// func CoreWordFrequency(reader io.Reader) map[string]int {
//     // BREAKPOINT: Watch scanner with ScanWords
//     scanner := bufio.NewScanner(reader)
//     scanner.Split(bufio.ScanWords)
//
//     freq := make(map[string]int)
//
//     // BREAKPOINT: Watch word counting
//     for scanner.Scan() {
//         word := strings.ToLower(scanner.Text())
//         freq[word]++
//     }
//
//     return freq
// }

func CoreWordFrequency(reader io.Reader) map[string]int {
	// TODO: implement word frequency counting
	return nil
}

// ============================================================================
// Core Exercise 4: Transform File (No Error Handling)
// ============================================================================

// CoreTransformFile transforms lines without error handling
// TODO: Implement line transformation
//
// Algorithm:
// 1. Create scanner for reading
// 2. Create writer for writing
// 3. Loop through lines
// 4. Apply transform function
// 5. Write transformed line
// 6. Flush writer
//
// func CoreTransformFile(input io.Reader, output io.Writer, transform func(string) string) {
//     // BREAKPOINT: Watch scanner and writer setup
//     scanner := bufio.NewScanner(input)
//     writer := bufio.NewWriter(output)
//     defer writer.Flush()
//
//     // BREAKPOINT: Watch line transformation
//     for scanner.Scan() {
//         line := scanner.Text()
//         transformed := transform(line)
//         writer.WriteString(transformed + "\n")
//     }
// }

func CoreTransformFile(input io.Reader, output io.Writer, transform func(string) string) {
	// TODO: implement line transformation
}

// ============================================================================
// Core Exercise 5: Read Chunks (No Error Handling)
// ============================================================================

// CoreReadChunks reads in chunks without error handling
// TODO: Implement chunked reading
//
// Algorithm:
// 1. Create buffer of chunkSize
// 2. Loop reading chunks
// 3. Call callback with each chunk
// 4. Track total bytes
// 5. Return total when EOF
//
// func CoreReadChunks(reader io.Reader, chunkSize int, callback func([]byte)) int {
//     // BREAKPOINT: Watch buffer creation
//     buffer := make([]byte, chunkSize)
//     totalBytes := 0
//
//     // BREAKPOINT: Watch chunk reading
//     for {
//         n, err := reader.Read(buffer)
//
//         if n > 0 {
//             // BREAKPOINT: Watch callback with chunk
//             callback(buffer[:n])
//             totalBytes += n
//         }
//
//         if err != nil {  // Simplified: treat any err as EOF
//             break
//         }
//     }
//
//     return totalBytes
// }

func CoreReadChunks(reader io.Reader, chunkSize int, callback func([]byte)) int {
	// TODO: implement chunked reading
	return 0
}

// After implementing core logic:
// - Compare with full implementation in exercise.go
// - Notice how error handling adds complexity
// - Use this to understand streaming I/O patterns
// - Set breakpoints to watch buffer operations
