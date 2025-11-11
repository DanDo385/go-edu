package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/example/go-10x-minis/minis/02-arrays-maps-basics/exercise"
)

func main() {
	// Demo 1: Read from testdata file
	fmt.Println("=== Reading from testdata/input.txt ===")
	file, err := os.Open("minis/02-arrays-maps-basics/testdata/input.txt")
	if err != nil {
		log.Printf("Could not open testdata file: %v (this is OK for demo)\n", err)
	} else {
		defer file.Close()
		freq, mostCommon, err := exercise.FreqFromReader(file)
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}
		fmt.Printf("Frequency map: %v\n", freq)
		fmt.Printf("Most common word: %q\n\n", mostCommon)
	}

	// Demo 2: Read from a string (shows io.Reader flexibility)
	fmt.Println("=== Reading from string ===")
	input := `go
is
awesome
go
rocks
go
`
	freq, mostCommon, err := exercise.FreqFromReader(strings.NewReader(input))
	if err != nil {
		log.Fatalf("Error reading string: %v", err)
	}
	fmt.Printf("Frequency map: %v\n", freq)
	fmt.Printf("Most common word: %q (appears %d times)\n", mostCommon, freq[mostCommon])
}
