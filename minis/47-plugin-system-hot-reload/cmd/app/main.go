package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"plugin"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/example/go-10x-minis/minis/47-plugin-system-hot-reload/shared"
)

// PluginManager manages the lifecycle of dynamically loaded plugins.
type PluginManager struct {
	plugins   map[string]*LoadedPlugin
	pluginDir string
	mu        sync.RWMutex
	stopWatch chan struct{}
}

// LoadedPlugin wraps a plugin with metadata.
type LoadedPlugin struct {
	plugin   shared.Plugin
	path     string
	loadTime time.Time
	modTime  time.Time
}

// NewPluginManager creates a new plugin manager.
func NewPluginManager(pluginDir string) (*PluginManager, error) {
	// Create plugin directory if it doesn't exist
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create plugin directory: %w", err)
	}

	pm := &PluginManager{
		plugins:   make(map[string]*LoadedPlugin),
		pluginDir: pluginDir,
		stopWatch: make(chan struct{}),
	}

	// Load existing plugins
	if err := pm.LoadAll(); err != nil {
		log.Printf("Warning: failed to load some plugins: %v", err)
	}

	// Start watching for changes (using polling)
	pm.StartWatching()

	return pm, nil
}

// LoadAll discovers and loads all plugins in the plugin directory.
func (pm *PluginManager) LoadAll() error {
	// Find all .so files
	matches, err := filepath.Glob(filepath.Join(pm.pluginDir, "*.so"))
	if err != nil {
		return err
	}

	var errors []error
	for _, path := range matches {
		if err := pm.LoadPlugin(path); err != nil {
			errors = append(errors, err)
			log.Printf("Failed to load plugin %s: %v", path, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to load %d plugins", len(errors))
	}

	return nil
}

// LoadPlugin loads a single plugin from a .so file.
func (pm *PluginManager) LoadPlugin(path string) error {
	// Get file info for modification time
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat failed: %w", err)
	}

	// Open the plugin
	p, err := plugin.Open(path)
	if err != nil {
		return &shared.PluginError{
			PluginName: filepath.Base(path),
			Operation:  "load",
			Err:        err,
		}
	}

	// Look up the exported Plugin symbol
	sym, err := p.Lookup("Plugin")
	if err != nil {
		return &shared.PluginError{
			PluginName: filepath.Base(path),
			Operation:  "load",
			Err:        fmt.Errorf("missing 'Plugin' symbol: %w", err),
		}
	}

	// Type assert to the Plugin interface
	plug, ok := sym.(shared.Plugin)
	if !ok {
		return &shared.PluginError{
			PluginName: filepath.Base(path),
			Operation:  "load",
			Err:        fmt.Errorf("symbol 'Plugin' is not of type shared.Plugin (got %T)", sym),
		}
	}

	// Initialize the plugin
	if err := plug.Init(); err != nil {
		return &shared.PluginError{
			PluginName: plug.Name(),
			Operation:  "init",
			Err:        err,
		}
	}

	// Store the loaded plugin
	pm.mu.Lock()
	defer pm.mu.Unlock()

	name := plug.Name()
	if old, exists := pm.plugins[name]; exists {
		// Clean up old version
		if err := old.plugin.Cleanup(); err != nil {
			log.Printf("Warning: cleanup of old %s failed: %v", name, err)
		}
		log.Printf("Replaced plugin %s: %s -> %s", name, old.plugin.Version(), plug.Version())
	} else {
		log.Printf("Loaded plugin %s v%s from %s", name, plug.Version(), filepath.Base(path))
	}

	pm.plugins[name] = &LoadedPlugin{
		plugin:   plug,
		path:     path,
		loadTime: time.Now(),
		modTime:  info.ModTime(),
	}

	return nil
}

// StartWatching begins monitoring the plugin directory for changes.
// Uses polling instead of inotify for better cross-platform compatibility.
func (pm *PluginManager) StartWatching() {
	go pm.watchLoop()
}

// watchLoop polls the plugin directory for file changes.
func (pm *PluginManager) watchLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Check all .so files for changes
			matches, err := filepath.Glob(filepath.Join(pm.pluginDir, "*.so"))
			if err != nil {
				log.Printf("Glob error: %v", err)
				continue
			}

			for _, path := range matches {
				info, err := os.Stat(path)
				if err != nil {
					continue
				}

				// Check if file was modified
				pm.mu.RLock()
				shouldReload := false
				for _, loaded := range pm.plugins {
					if loaded.path == path && info.ModTime().After(loaded.modTime) {
						shouldReload = true
						break
					}
				}
				pm.mu.RUnlock()

				if shouldReload {
					pm.handleFileChange(path)
				}
			}

		case <-pm.stopWatch:
			return
		}
	}
}

// handleFileChange processes a plugin file change event.
func (pm *PluginManager) handleFileChange(path string) {
	// Wait a bit to ensure the file write is complete
	time.Sleep(100 * time.Millisecond)

	// Check if the file actually changed
	info, err := os.Stat(path)
	if err != nil {
		log.Printf("Error stating file %s: %v", path, err)
		return
	}

	pm.mu.RLock()
	for _, loaded := range pm.plugins {
		if loaded.path == path {
			// File hasn't actually changed
			if !info.ModTime().After(loaded.modTime) {
				pm.mu.RUnlock()
				return
			}
			break
		}
	}
	pm.mu.RUnlock()

	log.Printf("Detected change in %s, reloading...", filepath.Base(path))

	// Reload the plugin
	if err := pm.LoadPlugin(path); err != nil {
		log.Printf("Failed to reload plugin %s: %v", filepath.Base(path), err)
		return
	}

	log.Printf("Successfully hot-reloaded %s", filepath.Base(path))
}

// GetPlugin retrieves a plugin by name (thread-safe).
func (pm *PluginManager) GetPlugin(name string) (shared.Plugin, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	loaded, ok := pm.plugins[name]
	if !ok {
		return nil, false
	}
	return loaded.plugin, true
}

// ListPlugins returns information about all loaded plugins.
func (pm *PluginManager) ListPlugins() []shared.PluginInfo {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	infos := make([]shared.PluginInfo, 0, len(pm.plugins))
	for _, loaded := range pm.plugins {
		infos = append(infos, shared.PluginInfo{
			Name:        loaded.plugin.Name(),
			Version:     loaded.plugin.Version(),
			Path:        loaded.path,
			LoadedAt:    loaded.loadTime.Format(time.RFC3339),
			LastUpdated: loaded.modTime.Format(time.RFC3339),
		})
	}

	return infos
}

// Process runs a plugin by name with the given input.
func (pm *PluginManager) Process(name string, input interface{}) (interface{}, error) {
	plug, ok := pm.GetPlugin(name)
	if !ok {
		return nil, fmt.Errorf("plugin not found: %s", name)
	}

	// Wrap in panic recovery for safety
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Plugin %s panicked: %v", name, r)
		}
	}()

	return plug.Process(input)
}

// Cleanup shuts down all plugins gracefully.
func (pm *PluginManager) Cleanup() {
	// Stop the watcher
	close(pm.stopWatch)

	pm.mu.Lock()
	defer pm.mu.Unlock()

	for name, loaded := range pm.plugins {
		if err := loaded.plugin.Cleanup(); err != nil {
			log.Printf("Cleanup error for plugin %s: %v", name, err)
		}
	}
}

// Interactive demo
func main() {
	fmt.Println("=== Go Plugin System with Hot Reload Demo ===\n")

	// Get plugin directory
	pluginDir := "./plugins"
	if len(os.Args) > 1 {
		pluginDir = os.Args[1]
	}

	fmt.Printf("Plugin directory: %s\n", pluginDir)
	fmt.Printf("Watching for .so file changes...\n\n")

	// Create plugin manager
	pm, err := NewPluginManager(pluginDir)
	if err != nil {
		log.Fatalf("Failed to create plugin manager: %v", err)
	}
	defer pm.Cleanup()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n\nShutting down...")
		pm.Cleanup()
		os.Exit(0)
	}()

	// Show loaded plugins
	listPlugins(pm)

	// Demo: Run some plugins
	if len(pm.plugins) > 0 {
		demoPlugins(pm)
	} else {
		fmt.Println("No plugins loaded. Add .so files to the plugins directory.")
		fmt.Println("\nExample: Build a plugin with:")
		fmt.Println("  go build -buildmode=plugin -o plugins/greeter.so plugins/greeter/greeter.go")
	}

	fmt.Println("\n=== Interactive Mode ===")
	fmt.Println("Commands:")
	fmt.Println("  list              - List all loaded plugins")
	fmt.Println("  run <plugin> <input> - Run a plugin with input")
	fmt.Println("  reload <file>     - Manually reload a plugin")
	fmt.Println("  quit              - Exit")
	fmt.Println("\nNote: Plugins will auto-reload when .so files change!")
	fmt.Println()

	// Interactive loop
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		cmd := parts[0]

		switch cmd {
		case "list":
			listPlugins(pm)

		case "run":
			if len(parts) < 2 {
				fmt.Println("Usage: run <plugin> [input]")
				continue
			}
			pluginName := parts[1]
			input := ""
			if len(parts) > 2 {
				input = strings.Join(parts[2:], " ")
			}

			result, err := pm.Process(pluginName, input)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Result: %v\n", result)
			}

		case "reload":
			if len(parts) < 2 {
				fmt.Println("Usage: reload <filename>")
				continue
			}
			path := filepath.Join(pluginDir, parts[1])
			if err := pm.LoadPlugin(path); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Reloaded successfully")
			}

		case "quit", "exit":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Printf("Unknown command: %s\n", cmd)
		}
	}
}

func listPlugins(pm *PluginManager) {
	infos := pm.ListPlugins()
	if len(infos) == 0 {
		fmt.Println("No plugins loaded.")
		return
	}

	fmt.Printf("\nLoaded Plugins (%d):\n", len(infos))
	fmt.Println(strings.Repeat("-", 70))
	for i, info := range infos {
		fmt.Printf("%d. %s (v%s)\n", i+1, info.Name, info.Version)
		fmt.Printf("   Path: %s\n", info.Path)
		fmt.Printf("   Loaded: %s\n", info.LoadedAt)
	}
	fmt.Println()
}

func demoPlugins(pm *PluginManager) {
	fmt.Println("=== Plugin Demo ===\n")

	// Try each plugin with sample inputs
	demos := []struct {
		plugin string
		input  interface{}
	}{
		{"greeter", "Alice"},
		{"greeter", "Bob"},
		{"math", map[string]interface{}{"op": "add", "a": 10.0, "b": 5.0}},
		{"math", map[string]interface{}{"op": "multiply", "a": 7.0, "b": 6.0}},
		{"transformer", "hello world"},
	}

	for _, demo := range demos {
		plug, ok := pm.GetPlugin(demo.plugin)
		if !ok {
			continue
		}

		result, err := pm.Process(demo.plugin, demo.input)
		if err != nil {
			fmt.Printf("[%s] Error: %v\n", demo.plugin, err)
		} else {
			fmt.Printf("[%s v%s] Input: %v → Output: %v\n",
				plug.Name(), plug.Version(), demo.input, result)
		}
	}
	fmt.Println()
}
