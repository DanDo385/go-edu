package main

import (
	"fmt"
	"log"
	"os"

	"github.com/example/go-10x-minis/minis/04-jsonl-log-filter/exercise"
)

func main() {
	// Open the testdata JSONL file
	file, err := os.Open("minis/04-jsonl-log-filter/testdata/logs.jsonl")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Filter logs at "warn" level and above
	minLevel := exercise.Warn
	entries, err := exercise.FilterLogs(file, minLevel)
	if err != nil {
		log.Printf("Warning: %v\n", err)
	}

	// Display results
	fmt.Printf("=== Logs (level >= %v) ===\n\n", minLevel)
	for _, entry := range entries {
		fmt.Printf("[%s] %s: %s\n",
			entry.TS.Format("2006-01-02 15:04:05"),
			levelString(entry.Level),
			entry.Msg,
		)
	}

	fmt.Printf("\nTotal entries: %d\n", len(entries))
}

func levelString(l exercise.Level) string {
	switch l {
	case exercise.Debug:
		return "DEBUG"
	case exercise.Info:
		return "INFO"
	case exercise.Warn:
		return "WARN"
	case exercise.Error:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}
