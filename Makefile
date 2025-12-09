.PHONY: setup list list-minis list-geth test bench run run-minis run-geth clean help

# Colors for output
CYAN := \033[0;36m
GREEN := \033[0;32m
YELLOW := \033[0;33m
NC := \033[0m # No Color

setup:
	@echo "$(CYAN)Initializing dependencies...$(NC)"
	go mod tidy
	@echo "\n$(CYAN)Verifying all minis/ projects compile...$(NC)"
	@for d in minis/*/; do \
		echo "Building $$(basename $$d)..."; \
		go build ./$$d/... 2>/dev/null || echo "  (no main package)"; \
	done
	@echo "\n$(CYAN)Verifying all geth/ projects compile...$(NC)"
	@for d in geth/*/; do \
		echo "Building $$(basename $$d)..."; \
		go build ./$$d/... 2>/dev/null || echo "  (no main package)"; \
	done
	@echo "\n$(GREEN)✓ All projects verified successfully$(NC)"

# List all projects from both tracks
list:
	@echo "$(CYAN)═══════════════════════════════════════$(NC)"
	@echo "$(CYAN)  Go Fundamentals (minis/)$(NC)"
	@echo "$(CYAN)═══════════════════════════════════════$(NC)"
	@ls -1d minis/*/ 2>/dev/null | sed 's|minis/||' | sed 's|/||' | nl -w2 -s'. ' || echo "No minis projects found"
	@echo "\n$(CYAN)═══════════════════════════════════════$(NC)"
	@echo "$(CYAN)  Ethereum Development (geth/)$(NC)"
	@echo "$(CYAN)═══════════════════════════════════════$(NC)"
	@ls -1d geth/*/ 2>/dev/null | sed 's|geth/||' | sed 's|/||' | nl -w2 -s'. ' || echo "No geth projects found"

list-minis:
	@echo "$(CYAN)Go Fundamentals (minis/):$(NC)"
	@ls -1d minis/*/ 2>/dev/null | sed 's|minis/||' | sed 's|/||' | nl -w2 -s'. '

list-geth:
	@echo "$(CYAN)Ethereum Development (geth/):$(NC)"
	@ls -1d geth/*/ 2>/dev/null | sed 's|geth/||' | sed 's|/||' | nl -w2 -s'. '

# Test with optional project path
test:
	@if [ -z "$(P)" ]; then \
		echo "$(CYAN)Running all tests...$(NC)"; \
		go test -v ./...; \
	else \
		PROJECT_PATH="$(P)"; \
		if [ ! -d "$$PROJECT_PATH" ] && [ -d "minis/$$PROJECT_PATH" ]; then \
			PROJECT_PATH="minis/$$PROJECT_PATH"; \
		fi; \
		if [ ! -d "$$PROJECT_PATH" ]; then \
			echo "$(YELLOW)Error: Project '$(P)' not found$(NC)"; \
			echo "Run 'make list' to see available projects"; \
			exit 1; \
		fi; \
		echo "$(CYAN)Testing $$PROJECT_PATH...$(NC)"; \
		go test -v ./$$PROJECT_PATH/...; \
	fi

# Run benchmarks with optional project path
bench:
	@if [ -z "$(P)" ]; then \
		echo "$(CYAN)Running all benchmarks...$(NC)"; \
		go test -bench=. -benchmem ./...; \
	else \
		PROJECT_PATH="$(P)"; \
		if [ ! -d "$$PROJECT_PATH" ] && [ -d "minis/$$PROJECT_PATH" ]; then \
			PROJECT_PATH="minis/$$PROJECT_PATH"; \
		fi; \
		if [ ! -d "$$PROJECT_PATH" ]; then \
			echo "$(YELLOW)Error: Project '$(P)' not found$(NC)"; \
			echo "Run 'make list' to see available projects"; \
			exit 1; \
		fi; \
		echo "$(CYAN)Benchmarking $$PROJECT_PATH...$(NC)"; \
		go test -bench=. -benchmem ./$$PROJECT_PATH/...; \
	fi

# Run a specific project (auto-detects minis/ or geth/)
run:
	@if [ -z "$(P)" ]; then \
		echo "$(YELLOW)Usage: make run P=<project>$(NC)"; \
		echo "Examples:"; \
		echo "  make run P=minis/01-hello-strings"; \
		echo "  make run P=geth/01-stack"; \
		echo "  make run P=01-hello-strings  (assumes minis/)"; \
		exit 1; \
	fi
	@PROJECT_PATH="$(P)"; \
	if [ ! -d "$$PROJECT_PATH" ] && [ -d "minis/$$PROJECT_PATH" ]; then \
		PROJECT_PATH="minis/$$PROJECT_PATH"; \
	fi; \
	if [ ! -d "$$PROJECT_PATH" ]; then \
		echo "$(YELLOW)Error: Project '$(P)' not found$(NC)"; \
		echo "Run 'make list' to see available projects"; \
		exit 1; \
	fi; \
	if [ -d "$$PROJECT_PATH/cmd" ]; then \
		echo "$(CYAN)Running $$PROJECT_PATH/cmd/...$(NC)"; \
		go run ./$$PROJECT_PATH/cmd/...; \
	else \
		echo "$(YELLOW)No cmd/ directory found in $$PROJECT_PATH$(NC)"; \
		echo "Try running tests instead: make test P=$$PROJECT_PATH"; \
	fi

# Explicit minis runner
run-minis:
	@if [ -z "$(P)" ]; then \
		echo "$(YELLOW)Usage: make run-minis P=<project-name>$(NC)"; \
		echo "Example: make run-minis P=01-hello-strings"; \
		exit 1; \
	fi
	@if [ ! -d "minis/$(P)" ]; then \
		echo "$(YELLOW)Error: Project 'minis/$(P)' not found$(NC)"; \
		echo "Run 'make list-minis' to see available projects"; \
		exit 1; \
	fi
	@if [ -d "minis/$(P)/cmd" ]; then \
		echo "$(CYAN)Running minis/$(P)/cmd/...$(NC)"; \
		go run ./minis/$(P)/cmd/...; \
	else \
		echo "$(YELLOW)No cmd/ directory found in minis/$(P)$(NC)"; \
		echo "Try: make test P=minis/$(P)"; \
	fi

# Explicit geth runner
run-geth:
	@if [ -z "$(P)" ]; then \
		echo "$(YELLOW)Usage: make run-geth P=<project-name>$(NC)"; \
		echo "Example: make run-geth P=01-stack"; \
		exit 1; \
	fi
	@if [ ! -d "geth/$(P)" ]; then \
		echo "$(YELLOW)Error: Project 'geth/$(P)' not found$(NC)"; \
		echo "Run 'make list-geth' to see available projects"; \
		exit 1; \
	fi
	@if [ -d "geth/$(P)/cmd" ]; then \
		echo "$(CYAN)Running geth/$(P)/cmd/...$(NC)"; \
		go run ./geth/$(P)/cmd/...; \
	else \
		echo "$(YELLOW)No cmd/ directory found in geth/$(P)$(NC)"; \
		echo "Try: make test P=geth/$(P)"; \
	fi

clean:
	@echo "$(CYAN)Cleaning build cache...$(NC)"
	go clean -testcache
	rm -f coverage.out
	@echo "$(GREEN)✓ Build cache cleaned$(NC)"

help:
	@echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"
	@echo "$(CYAN)  Go Educational Projects - Makefile Commands$(NC)"
	@echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"
	@echo ""
	@echo "$(GREEN)Setup & Discovery:$(NC)"
	@echo "  make setup           Initialize dependencies and verify builds"
	@echo "  make list            Show all available projects (both tracks)"
	@echo "  make list-minis      Show only minis/ projects"
	@echo "  make list-geth       Show only geth/ projects"
	@echo ""
	@echo "$(GREEN)Running Projects:$(NC)"
	@echo "  make run P=<path>    Run specific project (auto-detects track)"
	@echo "                       Examples:"
	@echo "                         make run P=minis/01-hello-strings"
	@echo "                         make run P=geth/01-stack"
	@echo "                         make run P=01-hello-strings  (assumes minis/)"
	@echo ""
	@echo "  make run-minis P=XX  Run minis project explicitly"
	@echo "                       Example: make run-minis P=01-hello-strings"
	@echo ""
	@echo "  make run-geth P=XX   Run geth project explicitly"
	@echo "                       Example: make run-geth P=01-stack"
	@echo ""
	@echo "$(GREEN)Testing:$(NC)"
	@echo "  make test            Run all tests (both tracks)"
	@echo "  make test P=<path>   Test specific project"
	@echo "                       Examples:"
	@echo "                         make test P=minis/03-csv-stats"
	@echo "                         make test P=geth/02-rpc-basics"
	@echo "                         make test P=03-csv-stats  (assumes minis/)"
	@echo ""
	@echo "$(GREEN)Benchmarking:$(NC)"
	@echo "  make bench           Run all benchmarks"
	@echo "  make bench P=<path>  Benchmark specific project"
	@echo "                       Example: make bench P=minis/07-generic-lru-cache"
	@echo ""
	@echo "$(GREEN)Cleanup:$(NC)"
	@echo "  make clean           Clean build cache"
	@echo ""
	@echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"
	@echo "$(YELLOW)Quick Start:$(NC)"
	@echo "  1. make setup                    # Initialize"
	@echo "  2. make list                     # See all projects"
	@echo "  3. make run P=minis/01-hello-strings  # Run first project"
	@echo ""
	@echo "$(YELLOW)For geth/ projects:$(NC)"
	@echo "  Export RPC URL: export INFURA_RPC_URL=https://mainnet.infura.io/v3/YOUR_KEY"
	@echo "  Then run: make run P=geth/01-stack"
	@echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"
