package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/example/go-10x-minis/minis/04-jsonl-log-filter/internal/jsonllogfilter"
)

/*
JSONL Log Filter CLI

Usage:
  go run ./cmd/app/main.go <file> <level>

Arguments:
  file    Path to JSONL file
  level   Minimum log level: debug, info, warn, error

Examples:
  go run ./cmd/app/main.go testdata/logs.jsonl warn
  go run ./cmd/app/main.go testdata/logs.jsonl error
*/
func main() {
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	levelStr := strings.ToLower(os.Args[2])
	var level jsonllogfilter.Level
	switch levelStr {
	case "debug":
		level = jsonllogfilter.LevelDebug
	case "info":
		level = jsonllogfilter.LevelInfo
	case "warn":
		level = jsonllogfilter.LevelWarn
	case "error":
		level = jsonllogfilter.LevelError
	default:
		fmt.Printf("Unknown level: %s\n", levelStr)
		printUsage()
		os.Exit(1)
	}

	fmt.Printf("Filtering logs from: %s (min level: %s)\n\n", os.Args[1], levelStr)

	entries, err := jsonllogfilter.FilterLogs(file, level)
	if err != nil {
		fmt.Printf("Error filtering logs: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d entries at level %s or above:\n\n", len(entries), levelStr)
	for _, entry := range entries {
		fmt.Printf("[%s] %s: %s\n", entry.TS.Format("2006-01-02 15:04:05"), levelToString(entry.Level), entry.Msg)
	}
}

func levelToString(l jsonllogfilter.Level) string {
	switch l {
	case jsonllogfilter.LevelDebug:
		return "DEBUG"
	case jsonllogfilter.LevelInfo:
		return "INFO"
	case jsonllogfilter.LevelWarn:
		return "WARN"
	case jsonllogfilter.LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func printUsage() {
	fmt.Println("JSONL Log Filter CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run ./cmd/app/main.go <file> <level>")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  file    Path to JSONL file")
	fmt.Println("  level   Minimum log level: debug, info, warn, error")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run ./cmd/app/main.go testdata/logs.jsonl warn")
	fmt.Println("  go run ./cmd/app/main.go testdata/logs.jsonl error")
}
