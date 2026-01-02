.PHONY: setup list test bench run clean help reset lint check

# Colors for output
CYAN := \033[0;36m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

# Detect current directory relative to repo root
ROOT_DIR := $(shell git rev-parse --show-toplevel 2>/dev/null || pwd)
CURRENT_DIR := $(shell pwd)
REL_DIR := $(shell realpath --relative-to=$(ROOT_DIR) $(CURRENT_DIR) 2>/dev/null || echo ".")

# Detect if we're in a project directory
IS_PROJECT := $(shell [ -d "$(CURRENT_DIR)/internal" ] && echo "yes" || echo "no")
PROJECT_NAME := $(shell basename $(CURRENT_DIR))

# ============================================================================
# SETUP & DISCOVERY
# ============================================================================

setup:
	@echo "$(CYAN)Initializing dependencies...$(NC)"
	@cd $(ROOT_DIR) && go mod tidy
	@echo "\n$(GREEN)✓ Setup complete$(NC)"

list:
	@echo "$(CYAN)═══════════════════════════════════════$(NC)"
	@echo "$(CYAN)  Go Fundamentals (minis/)$(NC)"
	@echo "$(CYAN)═══════════════════════════════════════$(NC)"
	@ls -1d $(ROOT_DIR)/minis/*/ 2>/dev/null | xargs -n1 basename | nl -w2 -s'. ' || echo "No minis projects found"
	@echo ""
	@echo "$(CYAN)═══════════════════════════════════════$(NC)"
	@echo "$(CYAN)  Ethereum Development (geth/)$(NC)"
	@echo "$(CYAN)═══════════════════════════════════════$(NC)"
	@ls -1d $(ROOT_DIR)/geth/*/ 2>/dev/null | xargs -n1 basename | nl -w2 -s'. ' || echo "No geth projects found"

# ============================================================================
# RUNNING PROJECTS
# ============================================================================

run:
	@if [ -z "$(P)" ]; then \
		if [ "$(IS_PROJECT)" = "yes" ]; then \
			if [ -d "$(CURRENT_DIR)/cmd/app" ]; then \
				echo "$(CYAN)Running $(PROJECT_NAME) CLI...$(NC)"; \
				go run ./cmd/app/main.go; \
			elif [ -d "$(CURRENT_DIR)/cmd/dev" ]; then \
				echo "$(CYAN)Running $(PROJECT_NAME) dev demo...$(NC)"; \
				go run ./cmd/dev/main.go; \
			else \
				echo "$(YELLOW)No cmd/ directory found. Try: make test$(NC)"; \
			fi; \
		else \
			echo "$(YELLOW)Usage: make run P=<project>$(NC)"; \
			echo "Examples:"; \
			echo "  make run P=minis/01-hello-strings"; \
			echo "  make run P=geth/01-stack"; \
			echo "  make run P=01-hello-strings  (assumes minis/)"; \
			echo ""; \
			echo "Or navigate to a project directory and run: make run"; \
		fi; \
	else \
		PROJECT_PATH="$(ROOT_DIR)/$(P)"; \
		if [ ! -d "$$PROJECT_PATH" ] && [ -d "$(ROOT_DIR)/minis/$(P)" ]; then \
			PROJECT_PATH="$(ROOT_DIR)/minis/$(P)"; \
		fi; \
		if [ ! -d "$$PROJECT_PATH" ]; then \
			echo "$(YELLOW)Error: Project '$(P)' not found$(NC)"; \
			echo "Run 'make list' to see available projects"; \
			exit 1; \
		fi; \
		if [ -d "$$PROJECT_PATH/cmd/app" ]; then \
			echo "$(CYAN)Running $$PROJECT_PATH/cmd/app...$(NC)"; \
			cd $$PROJECT_PATH && go run ./cmd/app/main.go; \
		elif [ -d "$$PROJECT_PATH/cmd/dev" ]; then \
			echo "$(CYAN)Running $$PROJECT_PATH/cmd/dev...$(NC)"; \
			cd $$PROJECT_PATH && go run ./cmd/dev/main.go; \
		else \
			echo "$(YELLOW)No cmd/ directory found$(NC)"; \
			echo "Try: make test P=$(P)"; \
		fi; \
	fi

# ============================================================================
# TESTING & BENCHMARKS
# ============================================================================

test:
	@if [ -z "$(P)" ]; then \
		if [ "$(IS_PROJECT)" = "yes" ]; then \
			echo "$(CYAN)Testing $(PROJECT_NAME)...$(NC)"; \
			go test -v ./...; \
		else \
			echo "$(CYAN)Running all tests...$(NC)"; \
			cd $(ROOT_DIR) && go test -v ./...; \
		fi; \
	else \
		PROJECT_PATH="$(ROOT_DIR)/$(P)"; \
		if [ ! -d "$$PROJECT_PATH" ] && [ -d "$(ROOT_DIR)/minis/$(P)" ]; then \
			PROJECT_PATH="$(ROOT_DIR)/minis/$(P)"; \
		fi; \
		if [ ! -d "$$PROJECT_PATH" ]; then \
			echo "$(YELLOW)Error: Project '$(P)' not found$(NC)"; \
			exit 1; \
		fi; \
		echo "$(CYAN)Testing $$PROJECT_PATH...$(NC)"; \
		cd $$PROJECT_PATH && go test -v ./...; \
	fi

bench:
	@if [ -z "$(P)" ]; then \
		if [ "$(IS_PROJECT)" = "yes" ]; then \
			echo "$(CYAN)Benchmarking $(PROJECT_NAME)...$(NC)"; \
			go test -bench=. -benchmem ./...; \
		else \
			echo "$(CYAN)Running all benchmarks...$(NC)"; \
			cd $(ROOT_DIR) && go test -bench=. -benchmem ./...; \
		fi; \
	else \
		PROJECT_PATH="$(ROOT_DIR)/$(P)"; \
		if [ ! -d "$$PROJECT_PATH" ] && [ -d "$(ROOT_DIR)/minis/$(P)" ]; then \
			PROJECT_PATH="$(ROOT_DIR)/minis/$(P)"; \
		fi; \
		if [ ! -d "$$PROJECT_PATH" ]; then \
			echo "$(YELLOW)Error: Project '$(P)' not found$(NC)"; \
			exit 1; \
		fi; \
		echo "$(CYAN)Benchmarking $$PROJECT_PATH...$(NC)"; \
		cd $$PROJECT_PATH && go test -bench=. -benchmem ./...; \
	fi

# ============================================================================
# CODE QUALITY
# ============================================================================

lint:
	@echo "$(CYAN)Running golangci-lint...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		cd $(ROOT_DIR) && golangci-lint run ./...; \
		echo "$(GREEN)✓ Lint complete$(NC)"; \
	else \
		echo "$(YELLOW)golangci-lint not installed. Install with:$(NC)"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

check:
	@echo "$(CYAN)Running go vet...$(NC)"
	@cd $(ROOT_DIR) && go vet ./...
	@echo "$(GREEN)✓ go vet complete$(NC)"
	@echo "$(CYAN)Running staticcheck...$(NC)"
	@if command -v staticcheck >/dev/null 2>&1; then \
		cd $(ROOT_DIR) && staticcheck ./...; \
		echo "$(GREEN)✓ staticcheck complete$(NC)"; \
	else \
		echo "$(YELLOW)staticcheck not installed. Install with:$(NC)"; \
		echo "  go install honnef.co/go/tools/cmd/staticcheck@latest"; \
	fi

# ============================================================================
# EXERCISE MANAGEMENT - Reset to TODO state
# ============================================================================

reset:
	@if [ -z "$(T)" ]; then \
		echo "$(YELLOW)Usage: make reset T=<target>$(NC)"; \
		echo ""; \
		echo "Reset exercise.go files to their original TODO state."; \
		echo ""; \
		echo "Targets:"; \
		echo "  make reset T=all              Reset all exercises"; \
		echo "  make reset T=geth             Reset all geth exercises"; \
		echo "  make reset T=minis            Reset all minis exercises"; \
		echo "  make reset T=geth/01-stack    Reset specific project"; \
		echo ""; \
		echo "$(YELLOW)Note:$(NC) This uses 'git checkout' to restore files."; \
		echo "Any uncommitted changes to exercise.go files will be lost."; \
		exit 0; \
	elif [ "$(T)" = "all" ]; then \
		echo "$(CYAN)Resetting all exercises to TODO state...$(NC)"; \
		cd $(ROOT_DIR) && git checkout HEAD -- 'minis/*/internal/*/exercise.go' 'geth/*/internal/*/exercise.go' 2>/dev/null || true; \
		echo "$(GREEN)✓ All exercises reset$(NC)"; \
	elif [ "$(T)" = "geth" ]; then \
		echo "$(CYAN)Resetting all geth exercises...$(NC)"; \
		cd $(ROOT_DIR) && git checkout HEAD -- 'geth/*/internal/*/exercise.go' 2>/dev/null || true; \
		echo "$(GREEN)✓ Geth exercises reset$(NC)"; \
	elif [ "$(T)" = "minis" ]; then \
		echo "$(CYAN)Resetting all minis exercises...$(NC)"; \
		cd $(ROOT_DIR) && git checkout HEAD -- 'minis/*/internal/*/exercise.go' 2>/dev/null || true; \
		echo "$(GREEN)✓ Minis exercises reset$(NC)"; \
	elif [ -d "$(ROOT_DIR)/$(T)" ]; then \
		echo "$(CYAN)Resetting exercises in $(T)...$(NC)"; \
		cd $(ROOT_DIR) && git checkout HEAD -- '$(T)/internal/*/exercise.go' 2>/dev/null || true; \
		echo "$(GREEN)✓ $(T) exercises reset$(NC)"; \
	else \
		echo "$(RED)Error: Unknown target '$(T)'$(NC)"; \
		echo "Run 'make reset' for usage."; \
		exit 1; \
	fi

# ============================================================================
# CLEANUP
# ============================================================================

clean:
	@echo "$(CYAN)Cleaning build cache...$(NC)"
	@cd $(ROOT_DIR) && go clean -testcache
	@rm -f $(ROOT_DIR)/coverage.out
	@find $(ROOT_DIR) -name "*.test" -delete 2>/dev/null || true
	@find $(ROOT_DIR) -name "*.out" -delete 2>/dev/null || true
	@echo "$(GREEN)✓ Build cache cleaned$(NC)"

# ============================================================================
# HELP - Context-aware help based on current directory
# ============================================================================

help:
	@if [ "$(IS_PROJECT)" = "yes" ]; then \
		echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"; \
		echo "$(CYAN)  Project: $(PROJECT_NAME)$(NC)"; \
		echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"; \
		echo ""; \
		echo "$(GREEN)Project Commands:$(NC)"; \
		echo "  make test            Run tests for this project"; \
		echo "  make bench           Run benchmarks for this project"; \
		echo "  make run             Run cmd/app/main.go (or cmd/dev/main.go)"; \
		echo ""; \
		echo "$(GREEN)Development:$(NC)"; \
		echo "  go test -v ./...                    # Run tests"; \
		echo "  go run ./cmd/app/main.go [args]     # Run CLI app"; \
		echo "  go run ./cmd/dev/main.go            # Run dev demo"; \
		echo ""; \
		echo "$(GREEN)Files:$(NC)"; \
		echo "  internal/*/exercise.go               # YOUR CODE GOES HERE"; \
		echo "  internal/*/exercise_test.go        # Test cases"; \
		echo "  internal/*/solution.reference.go   # Reference solution"; \
		echo "  README.md                           # Project documentation"; \
		echo ""; \
		echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"; \
	else \
		echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"; \
		echo "$(CYAN)  go-edu - Makefile Commands$(NC)"; \
		echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"; \
		echo ""; \
		echo "$(GREEN)Setup & Discovery:$(NC)"; \
		echo "  make setup           Initialize dependencies"; \
		echo "  make list            Show all available projects"; \
		echo ""; \
		echo "$(GREEN)Running Projects:$(NC)"; \
		echo "  make run P=<path>    Run specific project"; \
		echo "                       Examples:"; \
		echo "                         make run P=minis/01-hello-strings"; \
		echo "                         make run P=geth/01-stack"; \
		echo "                         make run P=01-hello-strings  (assumes minis/)"; \
		echo ""; \
		echo "$(GREEN)Testing:$(NC)"; \
		echo "  make test            Run all tests"; \
		echo "  make test P=<path>  Test specific project"; \
		echo "  make bench          Run all benchmarks"; \
		echo "  make bench P=<path> Benchmark specific project"; \
		echo ""; \
		echo "$(GREEN)Code Quality:$(NC)"; \
		echo "  make lint            Run golangci-lint"; \
		echo "  make check           Run go vet and staticcheck"; \
		echo ""; \
		echo "$(GREEN)Exercise Management:$(NC)"; \
		echo "  make reset T=all              Reset all exercises to TODO state"; \
		echo "  make reset T=geth             Reset all geth exercises"; \
		echo "  make reset T=minis            Reset all minis exercises"; \
		echo "  make reset T=geth/01-stack    Reset specific project"; \
		echo ""; \
		echo "$(GREEN)Cleanup:$(NC)"; \
		echo "  make clean           Clean build cache and artifacts"; \
		echo ""; \
		echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"; \
		echo "$(YELLOW)Quick Start:$(NC)"; \
		echo "  1. make setup                         # Initialize"; \
		echo "  2. make list                          # See all projects"; \
		echo "  3. cd minis/01-hello-strings          # Pick a project"; \
		echo "  4. code internal/hellostrings/exercise.go  # Implement TODOs"; \
		echo "  5. make test                           # Run tests"; \
		echo ""; \
		echo "$(YELLOW)For geth/ projects:$(NC)"; \
		echo "  export INFURA_RPC_URL=https://mainnet.infura.io/v3/YOUR_KEY"; \
		echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"; \
	fi
