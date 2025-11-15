// Package main implements a greeter plugin.
//
// This plugin demonstrates:
// - Basic plugin structure
// - Type assertion for input validation
// - Stateful plugin (greeting style can be configured)
// - String manipulation
package main

import (
	"fmt"
	"strings"

	"github.com/example/go-10x-minis/minis/47-plugin-system-hot-reload/shared"
)

// Plugin is the exported symbol that the host will look up.
// It must be named "Plugin" and implement the shared.Plugin interface.
var Plugin GreeterPlugin

// Ensure GreeterPlugin implements shared.Plugin interface at compile time
var _ shared.Plugin = (*GreeterPlugin)(nil)

// GreeterPlugin implements the shared.Plugin interface.
type GreeterPlugin struct {
	greetingStyle string
	greetCount    int
}

// Name returns the plugin identifier.
func (p *GreeterPlugin) Name() string {
	return "greeter"
}

// Version returns the semantic version.
// Change this when you update the plugin to test hot reload!
func (p *GreeterPlugin) Version() string {
	return "1.0.0"
}

// Init initializes the plugin.
// Called once when the plugin is loaded.
func (p *GreeterPlugin) Init() error {
	p.greetingStyle = "friendly"
	p.greetCount = 0
	return nil
}

// Process generates a greeting message.
//
// Input: string (name to greet)
// Output: string (greeting message)
//
// Example:
//   Input: "Alice"
//   Output: "Hello, Alice! Great to see you!"
func (p *GreeterPlugin) Process(input interface{}) (interface{}, error) {
	// Type assert input to string
	name, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("expected string input, got %T", input)
	}

	// Validate input
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	// Increment counter
	p.greetCount++

	// Generate greeting based on style
	var greeting string
	switch p.greetingStyle {
	case "friendly":
		greeting = fmt.Sprintf("Hello, %s! Great to see you!", name)
	case "formal":
		greeting = fmt.Sprintf("Good day, %s.", name)
	case "casual":
		greeting = fmt.Sprintf("Hey %s, what's up?", name)
	case "enthusiastic":
		greeting = fmt.Sprintf("HI %s!!! SO HAPPY TO MEET YOU!!!", strings.ToUpper(name))
	default:
		greeting = fmt.Sprintf("Hi, %s", name)
	}

	return greeting, nil
}

// Cleanup releases resources.
// Called when the plugin is unloaded.
func (p *GreeterPlugin) Cleanup() error {
	// Log statistics before cleanup
	if p.greetCount > 0 {
		fmt.Printf("[greeter] Processed %d greetings before cleanup\n", p.greetCount)
	}
	return nil
}

// Hot Reload Demo:
//
// To test hot reload, try changing the Version() or the greeting style:
//
// 1. Original version (1.0.0):
//    func (p *GreeterPlugin) Version() string {
//        return "1.0.0"
//    }
//
// 2. Change to enthusiastic style (1.1.0):
//    func (p *GreeterPlugin) Init() error {
//        p.greetingStyle = "enthusiastic"
//        p.greetCount = 0
//        return nil
//    }
//    func (p *GreeterPlugin) Version() string {
//        return "1.1.0"
//    }
//
// 3. Rebuild:
//    go build -buildmode=plugin -o plugins/greeter.so plugins/greeter/greeter.go
//
// 4. Watch the host application detect the change and hot reload!
