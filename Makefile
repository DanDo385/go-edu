.PHONY: setup list list-minis list-geth test bench run run-minis run-geth clean help todo lint check vet fmt

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

# Code quality checks
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "$(CYAN)Running golangci-lint...$(NC)"; \
		golangci-lint run ./...; \
		echo "$(GREEN)✓ Linting complete$(NC)"; \
	else \
		echo "$(YELLOW)golangci-lint not installed$(NC)"; \
		echo "Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		echo "Running basic checks instead..."; \
		$(MAKE) vet; \
	fi

vet:
	@echo "$(CYAN)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✓ Vet checks passed$(NC)"

fmt:
	@echo "$(CYAN)Checking code formatting...$(NC)"
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "$(YELLOW)The following files need formatting:$(NC)"; \
		gofmt -l .; \
		echo "$(YELLOW)Run 'make fmt-fix' to auto-format$(NC)"; \
		exit 1; \
	else \
		echo "$(GREEN)✓ All files properly formatted$(NC)"; \
	fi

fmt-fix:
	@echo "$(CYAN)Formatting code...$(NC)"
	@gofmt -w .
	@echo "$(GREEN)✓ Code formatted$(NC)"

check: vet fmt
	@echo "$(GREEN)✓ All checks passed$(NC)"

# Context-aware TODO reset command
# Usage:
#   make todo              # Context-aware: resets based on current directory
#   make todo P=<path>     # Reset specific path (e.g., P=geth/01-stack or P=minis)
todo:
	@if [ -n "$(P)" ]; then \
		./scripts/reset-exercises.sh "$(P)"; \
	else \
		CWD=$$(pwd); \
		REPO_ROOT=$$(git rev-parse --show-toplevel 2>/dev/null || pwd); \
		if [ "$$CWD" = "$$REPO_ROOT" ]; then \
			echo "$(CYAN)Resetting all exercises (from repository root)...$(NC)"; \
			./scripts/reset-exercises.sh; \
		elif echo "$$CWD" | grep -q "$$REPO_ROOT/geth/[^/]*$$"; then \
			PROJECT=$$(basename "$$CWD"); \
			echo "$(CYAN)Resetting geth/$$PROJECT (from project directory)...$(NC)"; \
			./scripts/reset-exercises.sh "geth/$$PROJECT"; \
		elif echo "$$CWD" | grep -q "$$REPO_ROOT/minis/[^/]*$$"; then \
			PROJECT=$$(basename "$$CWD"); \
			echo "$(CYAN)Resetting minis/$$PROJECT (from project directory)...$(NC)"; \
			./scripts/reset-exercises.sh "minis/$$PROJECT"; \
		elif echo "$$CWD" | grep -q "$$REPO_ROOT/geth$$"; then \
			echo "$(CYAN)Resetting all geth/ exercises (from geth directory)...$(NC)"; \
			./scripts/reset-exercises.sh geth; \
		elif echo "$$CWD" | grep -q "$$REPO_ROOT/minis$$"; then \
			echo "$(CYAN)Resetting all minis/ exercises (from minis directory)...$(NC)"; \
			./scripts/reset-exercises.sh minis; \
		else \
			echo "$(YELLOW)Cannot determine context. Use 'make todo P=<path>' to specify target.$(NC)"; \
			echo "Examples:"; \
			echo "  make todo P=minis              # Reset all minis"; \
			echo "  make todo P=geth               # Reset all geth"; \
			echo "  make todo P=geth/01-stack      # Reset specific project"; \
			exit 1; \
		fi; \
	fi

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
	@echo "$(GREEN)Code Quality:$(NC)"
	@echo "  make lint            Run linter (golangci-lint or go vet)"
	@echo "  make vet             Run go vet"
	@echo "  make fmt             Check code formatting"
	@echo "  make fmt-fix         Auto-format code"
	@echo "  make check           Run vet + fmt checks"
	@echo ""
	@echo "$(GREEN)Cleanup:$(NC)"
	@echo "  make clean           Clean build cache"
	@echo ""
	@echo "$(GREEN)Exercise Management:$(NC)"
	@echo "  make todo            Context-aware reset to TODO format"
	@echo "                       - From root: resets all exercises"
	@echo "                       - From geth/: resets all geth exercises"
	@echo "                       - From minis/: resets all minis exercises"
	@echo "                       - From project dir: resets that project only"
	@echo "  make todo P=<path>   Reset specific path"
	@echo "                       Examples:"
	@echo "                         make todo P=minis"
	@echo "                         make todo P=geth"
	@echo "                         make todo P=geth/01-stack"
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
