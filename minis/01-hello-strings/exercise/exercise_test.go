package exercise

import (
	"testing"
)

// TestTitleCase verifies title-case conversion with various inputs
func TestTitleCase(t *testing.T) {
	// Table-driven test: the idiomatic Go testing pattern
	// Each test case is a struct with input and expected output
	tests := []struct {
		name string // Test case name (appears in output)
		in   string // Input string
		want string // Expected result
	}{
		{
			name: "simple lowercase words",
			in:   "hello world",
			want: "Hello World",
		},
		{
			name: "multiple spaces",
			in:   "hello    world",
			want: "Hello World", // strings.Fields() collapses spaces
		},
		{
			name: "already capitalized",
			in:   "Hello World",
			want: "Hello World",
		},
		{
			name: "empty string",
			in:   "",
			want: "",
		},
		{
			name: "single word",
			in:   "go",
			want: "Go",
		},
		{
			name: "with emoji",
			in:   "hello 👋 world",
			want: "Hello 👋 World",
		},
		{
			name: "accented characters",
			in:   "café résumé",
			want: "Café Résumé",
		},
		{
			name: "mixed case preserved",
			in:   "iPhone macOS",
			want: "IPhone MacOS", // Capitalizes first letter only
		},
	}

	// Run each test case as a subtest
	// Subtests provide better failure messages and can be run individually:
	// go test -run TestTitleCase/with_emoji
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel() // Uncomment to run subtests concurrently
			got := TitleCase(tt.in)
			if got != tt.want {
				// Use %q to show quotes (makes whitespace visible)
				t.Errorf("TitleCase(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// TestReverse verifies UTF-8-aware string reversal
func TestReverse(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "simple ascii",
			in:   "hello",
			want: "olleh",
		},
		{
			name: "empty string",
			in:   "",
			want: "",
		},
		{
			name: "single character",
			in:   "a",
			want: "a",
		},
		{
			name: "palindrome",
			in:   "racecar",
			want: "racecar",
		},
		{
			name: "with emoji",
			in:   "Hello 👋 World",
			want: "dlroW 👋 olleH",
		},
		{
			name: "emoji only",
			in:   "👋😀🎉",
			want: "🎉😀👋",
		},
		{
			name: "accented characters",
			in:   "café",
			want: "éfac",
		},
		{
			name: "japanese characters",
			in:   "日本語",
			want: "語本日",
		},
		{
			name: "mixed unicode",
			in:   "Hello世界",
			want: "界世olleH",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Reverse(tt.in)
			if got != tt.want {
				t.Errorf("Reverse(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// TestRuneLen verifies correct character counting (not byte counting)
func TestRuneLen(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{
			name: "ascii string",
			in:   "hello",
			want: 5,
		},
		{
			name: "empty string",
			in:   "",
			want: 0,
		},
		{
			name: "single emoji",
			in:   "👋",
			want: 1, // 1 rune, but 4 bytes!
		},
		{
			name: "multiple emoji",
			in:   "👋😀🎉",
			want: 3, // 3 runes, but 12 bytes!
		},
		{
			name: "accented characters",
			in:   "café",
			want: 4, // 4 runes, but 5 bytes (é = 2 bytes)
		},
		{
			name: "japanese",
			in:   "日本語",
			want: 3, // 3 runes, but 9 bytes (each kanji = 3 bytes)
		},
		{
			name: "mixed content",
			in:   "Hello 👋 World",
			want: 13, // 13 runes (spaces count!)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RuneLen(tt.in)
			if got != tt.want {
				// Show byte length for comparison
				byteLen := len(tt.in)
				t.Errorf("RuneLen(%q) = %d, want %d (byte len=%d)",
					tt.in, got, tt.want, byteLen)
			}
		})
	}
}

// TestReverseTwice verifies that reversing twice returns the original string
// This is a property-based test idea (though not using a framework)
func TestReverseTwice(t *testing.T) {
	inputs := []string{
		"hello",
		"",
		"a",
		"Hello 👋 World",
		"日本語",
		"café résumé",
	}

	for _, in := range inputs {
		t.Run(in, func(t *testing.T) {
			got := Reverse(Reverse(in))
			if got != in {
				t.Errorf("Reverse(Reverse(%q)) = %q, want %q", in, got, in)
			}
		})
	}
}
