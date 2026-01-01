.PHONY: setup list list-minis list-geth test bench run run-minis run-geth clean help todo lint check

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

# Lint code using go vet (built-in) and optionally golangci-lint
lint:
	@if [ -z "$(P)" ]; then \
		echo "$(CYAN)Running go vet on all packages...$(NC)"; \
		go vet ./...; \
		if command -v golangci-lint >/dev/null 2>&1; then \
			echo "\n$(CYAN)Running golangci-lint...$(NC)"; \
			golangci-lint run ./...; \
		else \
			echo "\n$(YELLOW)Note: golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
		fi; \
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
		echo "$(CYAN)Linting $$PROJECT_PATH...$(NC)"; \
		go vet ./$$PROJECT_PATH/...; \
		if command -v golangci-lint >/dev/null 2>&1; then \
			golangci-lint run ./$$PROJECT_PATH/...; \
		fi; \
	fi

# Check code (vet + build) with optional project path
check:
	@if [ -z "$(P)" ]; then \
		echo "$(CYAN)Running go vet...$(NC)"; \
		go vet ./... || exit 1; \
		echo "\n$(CYAN)Building all packages...$(NC)"; \
		go build ./... || exit 1; \
		echo "\n$(GREEN)✓ All checks passed$(NC)"; \
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
		echo "$(CYAN)Checking $$PROJECT_PATH...$(NC)"; \
		go vet ./$$PROJECT_PATH/... || exit 1; \
		go build ./$$PROJECT_PATH/... || exit 1; \
		echo "\n$(GREEN)✓ Checks passed$(NC)"; \
	fi

# Contextual todo command - resets exercise.go files to TODO state
# Usage:
#   make todo all                    # Reset all exercises (from root)
#   make todo all                    # Reset all in current directory (from geth/ or minis/)
#   make todo all                    # Reset current project (from project directory)
#   make todo <path>                 # Reset specific path (from any location)
todo:
	@if [ -z "$(filter-out todo,$(MAKECMDGOALS))" ]; then \
		echo "$(YELLOW)Usage: make todo [all|<path>]$(NC)"; \
		echo "Examples:"; \
		echo "  make todo all                    # Reset all exercises (contextual)"; \
		echo "  make todo geth/01-stack          # Reset specific path"; \
		echo "  make todo minis/02-arrays-maps-basics"; \
		exit 1; \
	fi
	@TARGET="$(filter-out todo,$(MAKECMDGOALS))"; \
	CURRENT_DIR=$$(pwd); \
	REPO_ROOT=$$(git rev-parse --show-toplevel 2>/dev/null || echo "$$(pwd)"); \
	cd "$$REPO_ROOT"; \
	if [ "$$TARGET" = "all" ]; then \
		RELATIVE_PATH=$$(realpath --relative-to="$$REPO_ROOT" "$$CURRENT_DIR" 2>/dev/null || echo "$$CURRENT_DIR"); \
		if [ "$$RELATIVE_PATH" = "." ] || [ "$$RELATIVE_PATH" = "$$REPO_ROOT" ]; then \
			echo "$(CYAN)Resetting all exercises in both geth/ and minis/...$(NC)"; \
			./scripts/reset-exercises.sh; \
		elif [ "$$RELATIVE_PATH" = "geth" ]; then \
			echo "$(CYAN)Resetting all exercises in geth/...$(NC)"; \
			./scripts/reset-exercises.sh geth; \
		elif [ "$$RELATIVE_PATH" = "minis" ]; then \
			echo "$(CYAN)Resetting all exercises in minis/...$(NC)"; \
			./scripts/reset-exercises.sh minis; \
		elif echo "$$RELATIVE_PATH" | grep -q "^geth/"; then \
			PROJECT_PATH=$$(echo "$$RELATIVE_PATH" | cut -d'/' -f1-2); \
			if [ -d "$$PROJECT_PATH" ]; then \
				echo "$(CYAN)Resetting exercise in $$PROJECT_PATH...$(NC)"; \
				./scripts/reset-exercises.sh "$$PROJECT_PATH"; \
			else \
				echo "$(YELLOW)Error: Project directory not found: $$PROJECT_PATH$(NC)"; \
				exit 1; \
			fi; \
		elif echo "$$RELATIVE_PATH" | grep -q "^minis/"; then \
			PROJECT_PATH=$$(echo "$$RELATIVE_PATH" | cut -d'/' -f1-2); \
			if [ -d "$$PROJECT_PATH" ]; then \
				echo "$(CYAN)Resetting exercise in $$PROJECT_PATH...$(NC)"; \
				./scripts/reset-exercises.sh "$$PROJECT_PATH"; \
			else \
				echo "$(YELLOW)Error: Project directory not found: $$PROJECT_PATH$(NC)"; \
				exit 1; \
			fi; \
		else \
			echo "$(YELLOW)Error: Could not determine context$(NC)"; \
			echo "Current directory: $$CURRENT_DIR"; \
			echo "Relative path: $$RELATIVE_PATH"; \
			echo "Run from root, geth/, minis/, or a project directory"; \
			exit 1; \
		fi; \
		echo "$(GREEN)✓ Exercises reset$(NC)"; \
	else \
		TARGET_PATH="$$TARGET"; \
		if [ ! -d "$$TARGET_PATH" ] && [ -d "minis/$$TARGET_PATH" ]; then \
			TARGET_PATH="minis/$$TARGET_PATH"; \
		fi; \
		if [ ! -d "$$TARGET_PATH" ]; then \
			echo "$(YELLOW)Error: Path '$$TARGET' not found$(NC)"; \
			echo "Run 'make list' to see available projects"; \
			exit 1; \
		fi; \
		echo "$(CYAN)Resetting exercises in $$TARGET_PATH...$(NC)"; \
		./scripts/reset-exercises.sh "$$TARGET_PATH"; \
		echo "$(GREEN)✓ Exercises reset$(NC)"; \
	fi

# Prevent make from treating the todo argument as a target
%:
	@:

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
	@echo "  make lint            Run linters (go vet + golangci-lint if available)"
	@echo "  make lint P=<path>   Lint specific project"
	@echo "  make check           Run vet + build checks"
	@echo "  make check P=<path>  Check specific project"
	@echo ""
	@echo "$(GREEN)Cleanup:$(NC)"
	@echo "  make clean           Clean build cache"
	@echo ""
	@echo "$(GREEN)Exercise Management:$(NC)"
	@echo "  make todo all              Reset exercises (contextual - works from any directory)"
	@echo "                             Examples:"
	@echo "                               From root: make todo all  (resets all)"
	@echo "                               From geth/: make todo all  (resets all geth)"
	@echo "                               From project/: make todo all  (resets that project)"
	@echo "  make todo <path>           Reset exercises in specific path"
	@echo "                             Examples:"
	@echo "                               make todo geth/01-stack"
	@echo "                               make todo minis/02-arrays-maps-basics"
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
