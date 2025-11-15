// Package main implements a text transformation plugin.
//
// This plugin demonstrates:
// - String manipulation
// - Multiple transformation modes
// - Chainable transformations
// - Unicode handling
package main

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/example/go-10x-minis/minis/47-plugin-system-hot-reload/shared"
)

// Plugin is the exported symbol.
var Plugin TransformerPlugin

// Ensure TransformerPlugin implements shared.Plugin interface at compile time
var _ shared.Plugin = (*TransformerPlugin)(nil)

// TransformerPlugin implements text transformations.
type TransformerPlugin struct {
	mode string // "upper", "lower", "title", "reverse", "leetspeak"
}

// Name returns the plugin identifier.
func (p *TransformerPlugin) Name() string {
	return "transformer"
}

// Version returns the semantic version.
func (p *TransformerPlugin) Version() string {
	return "1.0.0"
}

// Init initializes the plugin with default mode.
func (p *TransformerPlugin) Init() error {
	p.mode = "upper" // Default transformation mode
	return nil
}

// Process transforms text based on the current mode.
//
// Input: string (text to transform)
// Output: string (transformed text)
//
// Transformation modes:
// - "upper": Convert to uppercase
// - "lower": Convert to lowercase
// - "title": Convert to title case
// - "reverse": Reverse the string
// - "leetspeak": Convert to l33t sp34k
//
// Example:
//   Input: "hello world"
//   Output (upper mode): "HELLO WORLD"
//   Output (reverse mode): "dlrow olleh"
func (p *TransformerPlugin) Process(input interface{}) (interface{}, error) {
	// Type assert to string
	text, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("expected string input, got %T", input)
	}

	// Apply transformation based on mode
	switch p.mode {
	case "upper":
		return strings.ToUpper(text), nil

	case "lower":
		return strings.ToLower(text), nil

	case "title":
		return strings.Title(text), nil

	case "reverse":
		return reverseString(text), nil

	case "leetspeak":
		return toLeetSpeak(text), nil

	case "alternating":
		return toAlternatingCase(text), nil

	case "snake":
		return toSnakeCase(text), nil

	default:
		return nil, fmt.Errorf("unknown transformation mode: %s", p.mode)
	}
}

// Cleanup releases resources.
func (p *TransformerPlugin) Cleanup() error {
	return nil
}

// reverseString reverses a string (handles Unicode correctly).
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// toLeetSpeak converts text to l33t sp34k.
func toLeetSpeak(s string) string {
	replacements := map[rune]rune{
		'a': '4', 'A': '4',
		'e': '3', 'E': '3',
		'i': '1', 'I': '1',
		'o': '0', 'O': '0',
		's': '5', 'S': '5',
		't': '7', 'T': '7',
		'l': '1', 'L': '1',
	}

	runes := []rune(s)
	for i, r := range runes {
		if replacement, ok := replacements[r]; ok {
			runes[i] = replacement
		}
	}
	return string(runes)
}

// toAlternatingCase converts text to AlTeRnAtInG cAsE.
func toAlternatingCase(s string) string {
	runes := []rune(s)
	upper := true
	for i, r := range runes {
		if unicode.IsLetter(r) {
			if upper {
				runes[i] = unicode.ToUpper(r)
			} else {
				runes[i] = unicode.ToLower(r)
			}
			upper = !upper
		}
	}
	return string(runes)
}

// toSnakeCase converts text to snake_case.
func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsSpace(r) {
			result = append(result, '_')
		} else if unicode.IsUpper(r) {
			if i > 0 && !unicode.IsSpace(rune(s[i-1])) {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// Hot Reload Demo:
//
// To test hot reload, try changing the default mode:
//
// Version 1.0.0 (original):
//   func (p *TransformerPlugin) Init() error {
//       p.mode = "upper"
//       return nil
//   }
//
// Version 1.1.0 (leetspeak mode):
//   func (p *TransformerPlugin) Init() error {
//       p.mode = "leetspeak"
//       return nil
//   }
//   func (p *TransformerPlugin) Version() string {
//       return "1.1.0"
//   }
//
// Version 1.2.0 (reverse mode):
//   func (p *TransformerPlugin) Init() error {
//       p.mode = "reverse"
//       return nil
//   }
//   func (p *TransformerPlugin) Version() string {
//       return "1.2.0"
//   }
//
// Rebuild and see the transformation change in real-time!
