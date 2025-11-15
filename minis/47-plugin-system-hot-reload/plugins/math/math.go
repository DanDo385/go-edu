// Package main implements a math operations plugin.
//
// This plugin demonstrates:
// - Complex input types (maps)
// - Multiple operations in one plugin
// - Error handling for invalid operations
// - Numeric type conversions
package main

import (
	"fmt"
	"math"

	"github.com/example/go-10x-minis/minis/47-plugin-system-hot-reload/shared"
)

// Plugin is the exported symbol that the host will look up.
var Plugin MathPlugin

// Ensure MathPlugin implements shared.Plugin interface at compile time
var _ shared.Plugin = (*MathPlugin)(nil)

// MathPlugin implements mathematical operations.
type MathPlugin struct {
	supportedOps map[string]bool
}

// Name returns the plugin identifier.
func (p *MathPlugin) Name() string {
	return "math"
}

// Version returns the semantic version.
func (p *MathPlugin) Version() string {
	return "1.0.0"
}

// Init initializes supported operations.
func (p *MathPlugin) Init() error {
	p.supportedOps = map[string]bool{
		"add":      true,
		"subtract": true,
		"multiply": true,
		"divide":   true,
		"power":    true,
		"sqrt":     true,
	}
	return nil
}

// Process performs mathematical operations.
//
// Input: map[string]interface{}{
//   "op": string,        // Operation: "add", "subtract", "multiply", "divide", "power", "sqrt"
//   "a":  float64,       // First operand
//   "b":  float64,       // Second operand (not needed for "sqrt")
// }
//
// Output: float64 (result of operation)
//
// Examples:
//   {"op": "add", "a": 10, "b": 5} → 15.0
//   {"op": "multiply", "a": 7, "b": 6} → 42.0
//   {"op": "sqrt", "a": 16} → 4.0
func (p *MathPlugin) Process(input interface{}) (interface{}, error) {
	// Type assert to map
	data, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected map[string]interface{}, got %T", input)
	}

	// Extract operation
	opRaw, ok := data["op"]
	if !ok {
		return nil, fmt.Errorf("missing 'op' field")
	}
	op, ok := opRaw.(string)
	if !ok {
		return nil, fmt.Errorf("'op' must be a string, got %T", opRaw)
	}

	// Check if operation is supported
	if !p.supportedOps[op] {
		return nil, fmt.Errorf("unsupported operation: %s", op)
	}

	// Extract operand a
	aRaw, ok := data["a"]
	if !ok {
		return nil, fmt.Errorf("missing 'a' operand")
	}
	a, err := toFloat64(aRaw)
	if err != nil {
		return nil, fmt.Errorf("invalid 'a' operand: %w", err)
	}

	// Perform operation
	switch op {
	case "add":
		b, err := extractOperandB(data)
		if err != nil {
			return nil, err
		}
		return a + b, nil

	case "subtract":
		b, err := extractOperandB(data)
		if err != nil {
			return nil, err
		}
		return a - b, nil

	case "multiply":
		b, err := extractOperandB(data)
		if err != nil {
			return nil, err
		}
		return a * b, nil

	case "divide":
		b, err := extractOperandB(data)
		if err != nil {
			return nil, err
		}
		if b == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return a / b, nil

	case "power":
		b, err := extractOperandB(data)
		if err != nil {
			return nil, err
		}
		return math.Pow(a, b), nil

	case "sqrt":
		if a < 0 {
			return nil, fmt.Errorf("square root of negative number")
		}
		return math.Sqrt(a), nil

	default:
		return nil, fmt.Errorf("operation not implemented: %s", op)
	}
}

// Cleanup releases resources.
func (p *MathPlugin) Cleanup() error {
	return nil
}

// Helper function to extract operand b from the input map.
func extractOperandB(data map[string]interface{}) (float64, error) {
	bRaw, ok := data["b"]
	if !ok {
		return 0, fmt.Errorf("missing 'b' operand")
	}
	return toFloat64(bRaw)
}

// Helper function to convert interface{} to float64.
// Handles both float64 and int inputs.
func toFloat64(v interface{}) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case float32:
		return float64(val), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", v)
	}
}

// Hot Reload Demo:
//
// To test hot reload, try adding a new operation:
//
// 1. Add "modulo" to supportedOps in Init():
//    p.supportedOps["modulo"] = true
//
// 2. Add the case in Process():
//    case "modulo":
//        b, err := extractOperandB(data)
//        if err != nil {
//            return nil, err
//        }
//        if b == 0 {
//            return nil, fmt.Errorf("modulo by zero")
//        }
//        return math.Mod(a, b), nil
//
// 3. Update version to "1.1.0"
//
// 4. Rebuild and watch it hot reload!
