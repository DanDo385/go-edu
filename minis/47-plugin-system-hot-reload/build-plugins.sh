#!/bin/bash

# Build script for Go plugins
# This script builds all plugins in the plugins/ directory

set -e  # Exit on error

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PLUGINS_SRC="$PROJECT_ROOT/plugins"
PLUGINS_OUT="$PROJECT_ROOT/plugins"

echo "=== Building Go Plugins ==="
echo "Project root: $PROJECT_ROOT"
echo "Plugin source: $PLUGINS_SRC"
echo "Plugin output: $PLUGINS_OUT"
echo

# Check if we're on a supported platform
OS=$(go env GOOS)
if [ "$OS" != "linux" ] && [ "$OS" != "darwin" ]; then
    echo "ERROR: Go plugins are only supported on Linux and macOS"
    echo "Current OS: $OS"
    exit 1
fi

# Create output directory if it doesn't exist
mkdir -p "$PLUGINS_OUT"

# Find and build all plugin source files
for plugin_dir in "$PLUGINS_SRC"/*/; do
    if [ ! -d "$plugin_dir" ]; then
        continue
    fi

    plugin_name=$(basename "$plugin_dir")
    source_file="$plugin_dir$plugin_name.go"
    output_file="$PLUGINS_OUT/$plugin_name.so"

    # Skip if no source file exists
    if [ ! -f "$source_file" ]; then
        echo "Skipping $plugin_name (no source file found)"
        continue
    fi

    echo "Building plugin: $plugin_name"
    echo "  Source: $source_file"
    echo "  Output: $output_file"

    # Build the plugin
    go build -buildmode=plugin -o "$output_file" "$source_file"

    if [ $? -eq 0 ]; then
        echo "  ✓ Built successfully"
        ls -lh "$output_file" | awk '{print "  Size: " $5}'
    else
        echo "  ✗ Build failed"
        exit 1
    fi
    echo
done

echo "=== Build Complete ==="
echo
echo "Built plugins:"
ls -lh "$PLUGINS_OUT"/*.so 2>/dev/null || echo "No plugins found"
echo
echo "To run the demo:"
echo "  go run cmd/plugin-demo/main.go"
echo
echo "To test hot reload:"
echo "  1. Run the demo in one terminal"
echo "  2. Edit a plugin source file (e.g., plugins/greeter/greeter.go)"
echo "  3. Re-run this script to rebuild"
echo "  4. Watch the demo automatically detect and reload the plugin!"
