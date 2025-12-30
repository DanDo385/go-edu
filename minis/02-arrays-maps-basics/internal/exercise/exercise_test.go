package exercise

import (
	"reflect"
	"strings"
	"testing"
)

func TestFreqFromReader(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		wantFreq       map[string]int
		wantMostCommon string
	}{
		{
			name:  "simple words",
			input: "hello\nworld\nhello\n",
			wantFreq: map[string]int{
				"hello": 2,
				"world": 1,
			},
			wantMostCommon: "hello",
		},
		{
			name:           "empty input",
			input:          "",
			wantFreq:       map[string]int{},
			wantMostCommon: "",
		},
		{
			name:  "single word",
			input: "go\n",
			wantFreq: map[string]int{
				"go": 1,
			},
			wantMostCommon: "go",
		},
		{
			name:  "case insensitive",
			input: "Go\ngo\nGO\nGo\n",
			wantFreq: map[string]int{
				"go": 4,
			},
			wantMostCommon: "go",
		},
		{
			name:  "with blank lines",
			input: "hello\n\nworld\n\nhello\n",
			wantFreq: map[string]int{
				"hello": 2,
				"world": 1,
			},
			wantMostCommon: "hello",
		},
		{
			name:  "with whitespace",
			input: "  hello  \n\tworld\t\n  hello\n",
			wantFreq: map[string]int{
				"hello": 2,
				"world": 1,
			},
			wantMostCommon: "hello",
		},
		{
			name:  "all unique words",
			input: "apple\nbanana\ncherry\n",
			wantFreq: map[string]int{
				"apple":  1,
				"banana": 1,
				"cherry": 1,
			},
			// For ties, any word is acceptable (map iteration is random)
			// We'll check if the result is one of the valid options
			wantMostCommon: "apple", // Could be any of the three
		},
		{
			name:  "repeated word",
			input: "go\ngo\ngo\ngo\ngo\n",
			wantFreq: map[string]int{
				"go": 5,
			},
			wantMostCommon: "go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a strings.Reader from test input
			// This demonstrates io.Reader flexibility: no file needed!
			r := strings.NewReader(tt.input)

			gotFreq, gotMostCommon, err := FreqFromReader(r)

			// Check for unexpected errors
			if err != nil {
				t.Fatalf("FreqFromReader() error = %v, want nil", err)
			}

			// Compare frequency maps
			// reflect.DeepEqual compares maps correctly (order-independent)
			if !reflect.DeepEqual(gotFreq, tt.wantFreq) {
				t.Errorf("FreqFromReader() freq = %v, want %v", gotFreq, tt.wantFreq)
			}

			// For the "all unique" case, accept any word from the map
			if tt.name == "all unique words" {
				if _, exists := tt.wantFreq[gotMostCommon]; !exists {
					t.Errorf("FreqFromReader() mostCommon = %q, not in freq map %v",
						gotMostCommon, tt.wantFreq)
				}
			} else {
				// For other cases, expect exact match
				if gotMostCommon != tt.wantMostCommon {
					t.Errorf("FreqFromReader() mostCommon = %q, want %q",
						gotMostCommon, tt.wantMostCommon)
				}
			}
		})
	}
}

// TestFreqFromReader_EmptyLines specifically tests blank line handling
func TestFreqFromReader_EmptyLines(t *testing.T) {
	input := "\n\n\n"
	r := strings.NewReader(input)

	freq, mostCommon, err := FreqFromReader(r)

	if err != nil {
		t.Fatalf("FreqFromReader() error = %v, want nil", err)
	}

	if len(freq) != 0 {
		t.Errorf("FreqFromReader() freq = %v, want empty map", freq)
	}

	if mostCommon != "" {
		t.Errorf("FreqFromReader() mostCommon = %q, want empty string", mostCommon)
	}
}

// TestFreqFromReader_UnicodeWords tests non-ASCII input
func TestFreqFromReader_UnicodeWords(t *testing.T) {
	input := "日本語\n日本語\ncafé\n"
	r := strings.NewReader(input)

	freq, mostCommon, err := FreqFromReader(r)

	if err != nil {
		t.Fatalf("FreqFromReader() error = %v, want nil", err)
	}

	wantFreq := map[string]int{
		"日本語": 2,
		"café": 1,
	}

	if !reflect.DeepEqual(freq, wantFreq) {
		t.Errorf("FreqFromReader() freq = %v, want %v", freq, wantFreq)
	}

	if mostCommon != "日本語" {
		t.Errorf("FreqFromReader() mostCommon = %q, want \"日本語\"", mostCommon)
	}
}
