//go:build !solution
// +build !solution

package exercise

import (
	"github.com/example/go-10x-minis/minis/47-plugin-system-hot-reload/shared"
)

// LoadAndExecute loads a plugin from the given path, initializes it,
// executes it with the provided input, and returns the result.
func LoadAndExecute(pluginPath string, input interface{}) (interface{}, error) {
	// TODO: Implement the full plugin lifecycle.

	// Step 1: Open the plugin file. This dynamically loads the shared object (`.so`) file.
	// - `p, err := plugin.Open(pluginPath)`
	// - Handle the error (e.g., file not found).

	// Step 2: Look up the exported symbol named "Plugin".
	// - `sym, err := p.Lookup("Plugin")`
	// - This symbol is expected to be a variable in the plugin's code that holds the plugin implementation.
	// - Handle the error (e.g., symbol not found).

	// Step 3: Type-assert the symbol to the `shared.Plugin` interface.
	// - `plug, ok := sym.(shared.Plugin)`
	// - This is a crucial step to verify that the plugin correctly implements the required interface.
	// - If `!ok`, return an error because the plugin has an invalid type.

	// Step 4: Initialize the plugin.
	// - `if err := plug.Init(); err != nil { ... }`
	// - This allows the plugin to perform any setup it needs before execution.

	// Step 5: Execute the plugin's main logic.
	// - `return plug.Process(input)`
	return nil, nil
}

// DiscoverPlugins scans the given directory for .so files
// and returns a list of their absolute paths.
func DiscoverPlugins(dir string) ([]string, error) {
	// TODO: Implement plugin discovery.

	// Step 1: Create a pattern to match all `.so` files in the directory.
	// - Use `filepath.Join(dir, "*.so")` to create a platform-independent path pattern.

	// Step 2: Use `filepath.Glob` to find all files matching the pattern.
	// - `matches, err := filepath.Glob(pattern)`
	// - Handle the error.

	// Step 3: Return the list of matches.
	return nil, nil
}

// ReloadPlugin simulates reloading a plugin by loading it fresh.
func ReloadPlugin(pluginPath string) (shared.Plugin, error) {
	// TODO: Implement this function.

	// Note: Go plugins cannot be truly "unloaded" from memory. Calling `plugin.Open` on the same path
	// will typically return the same loaded instance. For a real hot-reload system, you would
	// need to build the new version of the plugin to a different file name (e.g., with a version number or timestamp)
	// and then load the new file.

	// For this exercise, the implementation is very similar to the first part of `LoadAndExecute`.

	// Step 1: Open the plugin.
	// Step 2: Look up the "Plugin" symbol.
	// Step 3: Type-assert the symbol to `shared.Plugin`.
	// Step 4: Initialize the plugin by calling `Init()`.
	// Step 5: Return the initialized plugin instance.
	return nil, nil
}
