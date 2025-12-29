package main

import (
	"flag"   // Command-line argument parsing
	"fmt"    // Formatted I/O
	"log"    // Logging
	"os"     // Operating system interface

	// Import the exercise package from parent directory
	"github.com/example/go-10x-minis/minis/04-jsonl-log-filter/exercise"
)

/*
===================================================================
JSONL Log Filter: Structured Log Analysis
===================================================================

This program demonstrates JSON Lines (JSONL) parsing and filtering
by analyzing structured log files and filtering by severity level.

USAGE:
    go run ./cmd/jsonl-log-filter/main.go
    go run ./cmd/jsonl-log-filter/main.go -file=testdata/logs.jsonl
    go run ./cmd/jsonl-log-filter/main.go -level=error
    go run ./cmd/jsonl-log-filter/main.go -level=info -format=json
    go run ./cmd/jsonl-log-filter/main.go -file=custom.jsonl -level=warn

DEBUGGING:
- Set breakpoints at "// BREAKPOINT:" comments
- Use F5 to start debugging (select "Debug: Run main.go")
- Step Into (F11) to enter exercise.FilterLogs function
- Watch Variables panel to see log entries being parsed
- See /RUN_DEBUG.md for comprehensive debugging guide

PACKAGE STRUCTURE:
- This file (cmd/jsonl-log-filter/main.go) is package main
- It imports the exercise package from the parent directory
- exercise.go, solution.go are in package exercise
- Functions must be exported (capitalized) to be called from main
*/

func main() {
	// ============================================================================
	// STEP 1: Parse Command-Line Arguments
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to inspect default values before parsing
	// DEBUG: Before flag.Parse(), all flags have their default values
	// DEBUG: In Variables panel, you'll see: file="minis/04-jsonl-log-filter/testdata/logs.jsonl"

	// File path for JSONL log input
	filePath := flag.String("file", "minis/04-jsonl-log-filter/testdata/logs.jsonl", "Path to the JSONL log file")

	// Minimum log level to display: debug, info, warn, error
	levelStr := flag.String("level", "warn", "Minimum log level to display: debug, info, warn, error")

	// Output format: text or json
	format := flag.String("format", "text", "Output format: text, json")

	// Show timestamps
	showTime := flag.Bool("time", true, "Show timestamps in output")

	// Use solution.go instead of exercise.go
	useSolution := flag.Bool("solution", false, "Use solution.go instead of exercise.go (requires -tags=solution)")

	flag.Parse()

	// BREAKPOINT: Set breakpoint here to inspect parsed values in Variables panel
	// DEBUG: After flag.Parse(), variables now hold actual command-line values
	// DEBUG: Watch 'filePath', 'levelStr', 'format', 'showTime', 'useSolution'

	// ============================================================================
	// STEP 2: Convert Level String to exercise.Level Type
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before level conversion
	// DEBUG: Watch 'levelStr' - this is the user's input like "warn" or "error"
	// DEBUG: We need to convert string to exercise.Level constant

	var minLevel exercise.Level
	switch *levelStr {
	case "debug":
		minLevel = exercise.Debug
	case "info":
		minLevel = exercise.Info
	case "warn":
		minLevel = exercise.Warn
	case "error":
		minLevel = exercise.Error
	default:
		log.Printf("Invalid level %q, using 'warn'\n", *levelStr)
		minLevel = exercise.Warn
	}

	// BREAKPOINT: Set breakpoint here after level conversion
	// DEBUG: Inspect 'minLevel' - should be a numeric constant (0=Debug, 1=Info, 2=Warn, 3=Error)
	// DEBUG: This will be used to filter logs

	// ============================================================================
	// STEP 3: Open JSONL Log File
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before opening file
	// DEBUG: Watch 'filePath' variable - this is the JSONL file we'll parse
	// DEBUG: Step Over (F10) to execute file opening

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("Failed to open file %q: %v", *filePath, err)
	}
	defer file.Close()

	// BREAKPOINT: Set breakpoint here after opening file
	// DEBUG: File is now open and ready to read
	// DEBUG: JSONL format: one JSON object per line

	fmt.Printf("Processing JSONL log file: %s\n", *filePath)
	fmt.Printf("Minimum level: %s\n", levelString(minLevel))
	fmt.Println()

	// ============================================================================
	// STEP 4: Filter Logs by Level
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before calling FilterLogs
	// DEBUG: Step Into (F11) to see how JSONL is parsed line by line
	// DEBUG: Watch how each line is decoded from JSON to LogEntry struct

	entries, err := exercise.FilterLogs(file, minLevel)

	// BREAKPOINT: Set breakpoint here after filtering
	// DEBUG: Inspect 'entries' slice - expand to see all filtered log entries
	// DEBUG: Each entry has fields: TS (timestamp), Level, Msg (message)
	// DEBUG: Notice only entries with Level >= minLevel are included

	if err != nil {
		log.Printf("Warning: %v\n", err)
	}

	// ============================================================================
	// STEP 5: Display Results
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before displaying results
	// DEBUG: Watch 'format' variable to see which output format is active
	// DEBUG: Step Over (F10) to see results printed

	if *format == "json" {
		// JSON output format
		fmt.Println("[")

		// BREAKPOINT: Set breakpoint here before JSON loop
		// DEBUG: Expand 'entries' in Variables panel to see all log entries

		for i, entry := range entries {
			// BREAKPOINT: Set breakpoint here inside JSON loop
			// DEBUG: Watch 'entry' change with each iteration
			// DEBUG: See entry.TS (time.Time), entry.Level (int), entry.Msg (string)

			comma := ","
			if i == len(entries)-1 {
				comma = ""
			}

			fmt.Printf("  {\"time\":\"%s\", \"level\":\"%s\", \"msg\":\"%s\"}%s\n",
				entry.TS.Format("2006-01-02T15:04:05Z"),
				levelString(entry.Level),
				entry.Msg,
				comma,
			)
		}

		fmt.Println("]")

	} else {
		// Text output format (default)
		fmt.Printf("=== Logs (level >= %s) ===\n\n", levelString(minLevel))

		// BREAKPOINT: Set breakpoint here before text loop
		// DEBUG: Expand 'entries' in Variables panel to see all log entries
		// DEBUG: Notice slice structure: []LogEntry

		for _, entry := range entries {
			// BREAKPOINT: Set breakpoint here inside text loop
			// DEBUG: Watch 'entry.TS', 'entry.Level', 'entry.Msg' change
			// DEBUG: See how time.Time is formatted with .Format()

			if *showTime {
				fmt.Printf("[%s] %s: %s\n",
					entry.TS.Format("2006-01-02 15:04:05"),
					levelString(entry.Level),
					entry.Msg,
				)
			} else {
				fmt.Printf("%s: %s\n",
					levelString(entry.Level),
					entry.Msg,
				)
			}
		}

		fmt.Printf("\nTotal entries: %d\n", len(entries))
	}

	// ============================================================================
	// STEP 6: Statistics Summary
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before computing statistics
	// DEBUG: We'll count how many entries at each level
	// DEBUG: Watch 'levelCounts' map build up

	levelCounts := make(map[exercise.Level]int)
	for _, entry := range entries {
		// BREAKPOINT: Set breakpoint here inside count loop
		// DEBUG: Watch 'entry.Level' and see how map is incremented
		// DEBUG: Map access pattern: levelCounts[key]++

		levelCounts[entry.Level]++
	}

	// BREAKPOINT: Set breakpoint here after computing counts
	// DEBUG: Inspect 'levelCounts' map - expand to see counts per level
	// DEBUG: Keys are exercise.Level constants, values are counts

	if len(entries) > 0 && *format == "text" {
		fmt.Println()
		fmt.Println("=== Statistics ===")
		if levelCounts[exercise.Debug] > 0 {
			fmt.Printf("DEBUG entries: %d\n", levelCounts[exercise.Debug])
		}
		if levelCounts[exercise.Info] > 0 {
			fmt.Printf("INFO entries:  %d\n", levelCounts[exercise.Info])
		}
		if levelCounts[exercise.Warn] > 0 {
			fmt.Printf("WARN entries:  %d\n", levelCounts[exercise.Warn])
		}
		if levelCounts[exercise.Error] > 0 {
			fmt.Printf("ERROR entries: %d\n", levelCounts[exercise.Error])
		}
	}

	// ============================================================================
	// STEP 7: Usage Note about Solution
	// ============================================================================
	if *useSolution {
		fmt.Println()
		fmt.Println("NOTE: To use solution.go, build/run with -tags=solution:")
		fmt.Println("  go run -tags=solution ./cmd/jsonl-log-filter/main.go -solution=true")
		fmt.Println("  go test -tags=solution")
		os.Exit(0)
	}

	// ============================================================================
	// DEBUGGING TIPS
	// ============================================================================
	//
	// KEY DEBUGGING TECHNIQUES FOR JSONL PARSING:
	//
	// 1. Set breakpoints at "// BREAKPOINT:" comments
	// 2. Press F5 to start debugging
	// 3. Use F10 (Step Over) to execute line by line
	// 4. Use F11 (Step Into) to enter exercise.FilterLogs function
	// 5. Watch Variables panel to see []LogEntry slice build up
	// 6. Add watch expressions: len(entries), entries[0].Msg, levelCounts
	// 7. Use Debug Console to inspect: entries[0], entries[0].TS
	//
	// UNDERSTANDING JSONL FORMAT:
	// - JSONL = JSON Lines: one JSON object per line
	// - Step Into (F11) exercise.FilterLogs to see line-by-line parsing
	// - Watch how json.Decoder reads and unmarshals each line
	// - Each line becomes a LogEntry struct
	//
	// UNDERSTANDING TIME.TIME:
	// - In Variables panel, expand 'entry.TS' to see time components
	// - Notice wall (wall clock time) and ext (monotonic time)
	// - .Format() converts time.Time to string with layout pattern
	// - Layout "2006-01-02 15:04:05" is Go's reference time
	//
	// UNDERSTANDING SLICES OF STRUCTS:
	// - entries is []LogEntry (slice of structs)
	// - Expand 'entries' to see all elements
	// - Each element is a struct with TS, Level, Msg fields
	// - Watch slice grow as logs are filtered and appended
	//
	// See /RUN_DEBUG.md for comprehensive debugging guide
}

// levelString converts exercise.Level to string for display
func levelString(l exercise.Level) string {
	// BREAKPOINT: Set breakpoint here to see level conversion
	// DEBUG: Watch 'l' parameter - should be 0, 1, 2, or 3
	// DEBUG: This maps numeric levels to human-readable strings

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
