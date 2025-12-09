//go:build !solution
// +build !solution

package exercise

// TODO: Import required packages
// You'll need:
// - "bufio" for buffered I/O and line-by-line reading
// - "io" for the Reader interface
// - "strings" for string manipulation (ToLower, TrimSpace)
//
// import (
//     "bufio"
//     "io"
//     "strings"
// )

// FreqFromReader reads words from r (one per line) and returns:
// - A frequency map (word -> count)
// - The most common word (arbitrary choice if there's a tie)
// - An error if reading fails
//
// TODO: Implement FreqFromReader function
// Function signature: func FreqFromReader(r io.Reader) (map[string]int, string, error)
//
// Steps to implement:
//
// 1. Create an empty frequency map
//    - Use: freq := make(map[string]int)
//    - Why make()? Maps are REFERENCE types (like slices)
//    - The map variable stores a pointer to the hash table
//    - nil maps can't be written to, so we must initialize with make()
//    - Alternative: freq := map[string]int{} also works
//
// 2. Create a scanner to read line-by-line
//    - Use: scanner := bufio.NewScanner(r)
//    - bufio.Scanner is a STRUCT VALUE (not a pointer!)
//    - But it contains a pointer to internal buffers
//    - Passing scanner to functions copies the struct, but the pointer inside
//      still points to the same buffer (shared state)
//
// 3. Loop through each line
//    - Use: for scanner.Scan() { ... }
//    - Scan() MODIFIES the scanner's internal state (advances position)
//    - Returns true if a line was read, false if EOF or error
//    - scanner.Text() returns the current line as a string
//      * This returns a NEW string (copy of data from buffer)
//      * Safe to store - won't be overwritten on next Scan()
//
// 4. Normalize each word
//    - Use: word := strings.ToLower(strings.TrimSpace(line))
//    - strings.TrimSpace(line) creates a NEW string (substring, may share data)
//    - strings.ToLower(...) creates another NEW string
//    - Strings are IMMUTABLE - these functions never modify the input
//    - Each function call may allocate memory for the result
//
// 5. Skip blank lines
//    - Use: if word == "" { continue }
//    - String comparison (==) compares CONTENTS, not pointers
//    - Empty string "" is a zero-length string (not nil!)
//
// 6. Increment count in map
//    - Use: freq[word]++
//    - CRITICAL Go feature: Maps have "zero value" behavior
//    - If freq[word] doesn't exist, reading it returns 0 (int's zero value)
//    - So freq[word]++ works even for first occurrence!
//    - Behind the scenes:
//      * Go checks if key exists in hash table
//      * If not, creates new entry with value = 0
//      * Increments the value
//      * Stores the updated value back in the hash table
//    - This modifies the map's internal data (which lives on the heap)
//
// 7. Check for scanner errors
//    - Use: if err := scanner.Err(); err != nil { return freq, "", err }
//    - scanner.Err() returns any error encountered during scanning
//    - Must check AFTER the loop (not inside it!)
//    - If error occurred, return what we have so far + the error
//
// 8. Find the most common word
//    - Declare variables: var maxWord string; var maxCount int
//    - Loop through map: for word, count := range freq { ... }
//      * range on a map returns (key, value) pairs
//      * Order is RANDOM (maps don't maintain insertion order)
//      * word and count are COPIES of map entries
//      * Modifying word or count won't change the map
//    - Compare: if count > maxCount { maxWord = word; maxCount = count }
//    - After loop, maxWord contains the most frequent word
//
// 9. Return results
//    - Use: return freq, maxWord, nil
//    - freq is a map (reference type) - passes pointer to caller
//    - maxWord is a string (value type, but shares data efficiently)
//    - nil means no error occurred
//
// Key Go concepts:
// - Maps are reference types (store pointer to hash table)
// - Map access freq[key] returns zero value if key doesn't exist
// - Strings are immutable (functions create new strings, never modify)
// - Strings are passed efficiently (copy pointer + length, not data)
// - range on maps is non-deterministic (random iteration order)
// - Multiple return values (result, result, error) is idiomatic Go
//
// Memory notes:
// - freq map grows as needed (starts small, resizes when full)
// - Each unique word allocates: map entry + string data
// - scanner uses a fixed-size buffer (default 64KB)
// - No memory leaks: GC cleans up when function returns

// TODO: Implement the FreqFromReader function below
// func FreqFromReader(r io.Reader) (map[string]int, string, error) {
//     return nil, "", nil
// }

// After implementing:
// - Run: go test ./...
// - Check: go test -v for verbose output
// - Compare with solution.go to see detailed explanations
