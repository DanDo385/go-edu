.PHONY: setup list test bench run clean help

setup:
	go mod tidy
	@echo "Checking all projects compile..."
	@for d in minis/*/; do \
		echo "Building $$(basename $$d)..."; \
		go build ./$$d/... || exit 1; \
	done
	@echo "✓ All projects compile successfully"

list:
	@echo "Available projects:"
	@ls -1 minis/ | nl

test:
	go test -v ./...

bench:
	go test -bench=. -benchmem ./...

run:
	@if [ -z "$(P)" ]; then \
		echo "Usage: make run P=01-hello-strings"; \
		exit 1; \
	fi
	@if [ ! -d "minis/$(P)" ]; then \
		echo "Error: Project '$(P)' not found"; \
		echo "Run 'make list' to see available projects"; \
		exit 1; \
	fi
	go run ./minis/$(P)/cmd/...

clean:
	go clean -testcache
	rm -f coverage.out

help:
	@echo "Available targets:"
	@echo "  make setup    - Initialize dependencies and verify builds"
	@echo "  make list     - Show all available projects"
	@echo "  make test     - Run all tests"
	@echo "  make bench    - Run benchmarks"
	@echo "  make run P=XX - Run specific project (e.g., make run P=01-hello-strings)"
	@echo "  make clean    - Clean build cache"
