package exercise

import (
	"bytes"
	"strings"
	"testing"
)

func TestCountLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "empty file",
			input:    "",
			expected: 0,
		},
		{
			name:     "single line without newline",
			input:    "hello",
			expected: 1,
		},
		{
			name:     "single line with newline",
			input:    "hello\n",
			expected: 1,
		},
		{
			name:     "multiple lines",
			input:    "line1\nline2\nline3\n",
			expected: 3,
		},
		{
			name:     "multiple lines without trailing newline",
			input:    "line1\nline2\nline3",
			expected: 3,
		},
		{
			name:     "empty lines",
			input:    "\n\n\n",
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			got, err := CountLines(reader)
			if err != nil {
				t.Fatalf("CountLines() error = %v", err)
			}
			if got != tt.expected {
				t.Errorf("CountLines() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestFilterLines(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		predicate      func(string) bool
		expectedOutput string
		expectedCount  int
	}{
		{
			name:  "filter ERROR lines",
			input: "INFO: starting\nERROR: failed\nINFO: continuing\nERROR: crashed\n",
			predicate: func(line string) bool {
				return strings.Contains(line, "ERROR")
			},
			expectedOutput: "ERROR: failed\nERROR: crashed\n",
			expectedCount:  2,
		},
		{
			name:  "filter all lines (predicate always true)",
			input: "line1\nline2\nline3\n",
			predicate: func(line string) bool {
				return true
			},
			expectedOutput: "line1\nline2\nline3\n",
			expectedCount:  3,
		},
		{
			name:  "filter no lines (predicate always false)",
			input: "line1\nline2\nline3\n",
			predicate: func(line string) bool {
				return false
			},
			expectedOutput: "",
			expectedCount:  0,
		},
		{
			name:  "empty input",
			input: "",
			predicate: func(line string) bool {
				return true
			},
			expectedOutput: "",
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := strings.NewReader(tt.input)
			output := &bytes.Buffer{}

			count, err := FilterLines(input, output, tt.predicate)
			if err != nil {
				t.Fatalf("FilterLines() error = %v", err)
			}

			if count != tt.expectedCount {
				t.Errorf("FilterLines() count = %d, want %d", count, tt.expectedCount)
			}

			if output.String() != tt.expectedOutput {
				t.Errorf("FilterLines() output = %q, want %q", output.String(), tt.expectedOutput)
			}
		})
	}
}

func TestWordFrequency(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]int
	}{
		{
			name:     "empty input",
			input:    "",
			expected: map[string]int{},
		},
		{
			name:  "single word",
			input: "hello",
			expected: map[string]int{
				"hello": 1,
			},
		},
		{
			name:  "multiple words",
			input: "go is great go is fast",
			expected: map[string]int{
				"go":    2,
				"is":    2,
				"great": 1,
				"fast":  1,
			},
		},
		{
			name:  "case insensitive",
			input: "Go GO go",
			expected: map[string]int{
				"go": 3,
			},
		},
		{
			name:  "multiple lines",
			input: "line one\nline two\nline three",
			expected: map[string]int{
				"line":  3,
				"one":   1,
				"two":   1,
				"three": 1,
			},
		},
		{
			name:  "extra whitespace",
			input: "  hello    world  \n\n  foo   ",
			expected: map[string]int{
				"hello": 1,
				"world": 1,
				"foo":   1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			got, err := WordFrequency(reader)
			if err != nil {
				t.Fatalf("WordFrequency() error = %v", err)
			}

			if len(got) != len(tt.expected) {
				t.Errorf("WordFrequency() got %d words, want %d words", len(got), len(tt.expected))
			}

			for word, expectedCount := range tt.expected {
				if gotCount, exists := got[word]; !exists {
					t.Errorf("WordFrequency() missing word %q", word)
				} else if gotCount != expectedCount {
					t.Errorf("WordFrequency() word %q = %d, want %d", word, gotCount, expectedCount)
				}
			}
		})
	}
}

func TestTransformFile(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		transform      func(string) string
		expectedOutput string
	}{
		{
			name:  "uppercase transform",
			input: "hello\nworld\n",
			transform: func(line string) string {
				return strings.ToUpper(line)
			},
			expectedOutput: "HELLO\nWORLD\n",
		},
		{
			name:  "add prefix",
			input: "line1\nline2\n",
			transform: func(line string) string {
				return ">> " + line
			},
			expectedOutput: ">> line1\n>> line2\n",
		},
		{
			name:  "reverse lines",
			input: "abc\ndef\n",
			transform: func(line string) string {
				runes := []rune(line)
				for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
					runes[i], runes[j] = runes[j], runes[i]
				}
				return string(runes)
			},
			expectedOutput: "cba\nfed\n",
		},
		{
			name:  "empty input",
			input: "",
			transform: func(line string) string {
				return strings.ToUpper(line)
			},
			expectedOutput: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := strings.NewReader(tt.input)
			output := &bytes.Buffer{}

			err := TransformFile(input, output, tt.transform)
			if err != nil {
				t.Fatalf("TransformFile() error = %v", err)
			}

			if output.String() != tt.expectedOutput {
				t.Errorf("TransformFile() output = %q, want %q", output.String(), tt.expectedOutput)
			}
		})
	}
}

func TestReadChunks(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		chunkSize     int
		expectedTotal int
		expectedCalls int
	}{
		{
			name:          "exact chunks",
			input:         "123456789012",
			chunkSize:     4,
			expectedTotal: 12,
			expectedCalls: 3, // 4+4+4
		},
		{
			name:          "uneven chunks",
			input:         "1234567890",
			chunkSize:     3,
			expectedTotal: 10,
			expectedCalls: 4, // 3+3+3+1
		},
		{
			name:          "single chunk",
			input:         "hello",
			chunkSize:     10,
			expectedTotal: 5,
			expectedCalls: 1,
		},
		{
			name:          "empty input",
			input:         "",
			chunkSize:     10,
			expectedTotal: 0,
			expectedCalls: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			calls := 0
			totalBytes := 0

			gotTotal, err := ReadChunks(reader, tt.chunkSize, func(chunk []byte) {
				calls++
				totalBytes += len(chunk)
			})

			if err != nil {
				t.Fatalf("ReadChunks() error = %v", err)
			}

			if gotTotal != tt.expectedTotal {
				t.Errorf("ReadChunks() total = %d, want %d", gotTotal, tt.expectedTotal)
			}

			if calls != tt.expectedCalls {
				t.Errorf("ReadChunks() calls = %d, want %d", calls, tt.expectedCalls)
			}

			if totalBytes != tt.expectedTotal {
				t.Errorf("ReadChunks() total bytes in chunks = %d, want %d", totalBytes, tt.expectedTotal)
			}
		})
	}
}

// Benchmark tests to demonstrate performance benefits

func BenchmarkCountLines(b *testing.B) {
	// Generate a large input
	var sb strings.Builder
	for i := 0; i < 10000; i++ {
		sb.WriteString("This is a test line with some content\n")
	}
	input := sb.String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(input)
		CountLines(reader)
	}
}

func BenchmarkWordFrequency(b *testing.B) {
	// Generate a large input
	var sb strings.Builder
	for i := 0; i < 1000; i++ {
		sb.WriteString("the quick brown fox jumps over the lazy dog\n")
	}
	input := sb.String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(input)
		WordFrequency(reader)
	}
}

func BenchmarkFilterLines(b *testing.B) {
	// Generate a large input
	var sb strings.Builder
	for i := 0; i < 5000; i++ {
		sb.WriteString("INFO: normal log line\n")
		sb.WriteString("ERROR: error log line\n")
	}
	input := sb.String()

	predicate := func(line string) bool {
		return strings.Contains(line, "ERROR")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(input)
		output := &bytes.Buffer{}
		FilterLines(reader, output, predicate)
	}
}
