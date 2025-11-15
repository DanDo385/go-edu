.PHONY: all build-plugins run test clean help

# Default target
all: build-plugins

# Build all plugins
build-plugins:
	@echo "Building all plugins..."
	@./build-plugins.sh

# Run the plugin demo
run: build-plugins
	@echo "Starting plugin demo..."
	@go run cmd/plugin-demo/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with solution
test-solution:
	@echo "Running tests with solution..."
	@go test -v -tags solution ./exercise/

# Clean built plugins
clean:
	@echo "Cleaning built plugins..."
	@rm -f plugins/*.so
	@go clean -testcache

# Watch mode for development (requires entr or similar)
watch:
	@echo "Watching for changes... (requires 'entr' to be installed)"
	@echo "Press Ctrl+C to stop"
	@ls plugins/*/*.go | entr -r make build-plugins

# Help
help:
	@echo "Plugin System - Available targets:"
	@echo "  make build-plugins  - Build all plugins"
	@echo "  make run           - Build plugins and run demo"
	@echo "  make test          - Run tests"
	@echo "  make test-solution - Run tests with solution"
	@echo "  make clean         - Remove built plugins"
	@echo "  make watch         - Watch for changes and rebuild (requires entr)"
	@echo ""
	@echo "Hot Reload Workflow:"
	@echo "  Terminal 1: make run"
	@echo "  Terminal 2: Edit plugins/greeter/greeter.go"
	@echo "  Terminal 2: make build-plugins"
	@echo "  Terminal 1: Watch plugin auto-reload!"
