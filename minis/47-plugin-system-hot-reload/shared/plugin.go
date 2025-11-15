// Package shared defines the plugin interface and shared types.
// This package must be imported by both the host application and all plugins
// to ensure type compatibility.
package shared

// Plugin is the interface that all plugins must implement.
//
// The plugin system loads .so files dynamically and looks up symbols
// that implement this interface.
//
// Design principles:
// - Stable: This interface should rarely change to avoid breaking plugins
// - Simple: Complex interfaces are hard to implement correctly
// - Versioned: Include version info for compatibility checking
type Plugin interface {
	// Name returns a unique identifier for this plugin.
	// Used for registration and lookup.
	Name() string

	// Version returns the semantic version of this plugin.
	// Format: "major.minor.patch" (e.g., "1.2.3")
	Version() string

	// Init is called once when the plugin is loaded.
	// Use this to:
	// - Allocate resources
	// - Validate configuration
	// - Set up initial state
	//
	// If Init returns an error, the plugin is not registered.
	Init() error

	// Process is the main function of the plugin.
	// It takes arbitrary input and returns arbitrary output.
	//
	// Type safety is the responsibility of the plugin implementation.
	// Most plugins will type-assert the input to expected types.
	//
	// Example:
	//   func (p *MyPlugin) Process(input interface{}) (interface{}, error) {
	//       str, ok := input.(string)
	//       if !ok {
	//           return nil, fmt.Errorf("expected string, got %T", input)
	//       }
	//       return strings.ToUpper(str), nil
	//   }
	Process(input interface{}) (interface{}, error)

	// Cleanup is called when the plugin is unloaded or the application exits.
	// Use this to:
	// - Release resources (files, connections, memory)
	// - Flush buffers
	// - Save state
	//
	// Cleanup should be idempotent (safe to call multiple times).
	Cleanup() error
}

// PluginInfo contains metadata about a loaded plugin.
type PluginInfo struct {
	Name        string
	Version     string
	Path        string
	LoadedAt    string
	LastUpdated string
}

// PluginError wraps errors from plugin operations with context.
type PluginError struct {
	PluginName string
	Operation  string // "load", "init", "process", "cleanup"
	Err        error
}

func (e *PluginError) Error() string {
	return "plugin " + e.PluginName + " " + e.Operation + " failed: " + e.Err.Error()
}

func (e *PluginError) Unwrap() error {
	return e.Err
}
