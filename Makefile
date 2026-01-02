.PHONY: setup list list-minis list-geth test bench run clean help reset lint check dev
.DEFAULT_GOAL := help

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

# Detect current context
IN_PROJECT := $(shell test -f "$(CURRENT_DIR)/go.mod" || test -d "$(CURRENT_DIR)/cmd" && echo "yes" || echo "no")
PROJECT_NAME := $(shell basename $(CURRENT_DIR))
IS_GETH := $(shell echo $(REL_DIR) | grep -q "^geth/" && echo "yes" || echo "no")
IS_MINIS := $(shell echo $(REL_DIR) | grep -q "^minis/" && echo "yes" || echo "no")

# ============================================================================
# SETUP & DISCOVERY
# ============================================================================

setup:
	@echo "$(CYAN)Initializing dependencies...$(NC)"
	go mod tidy
	@echo "\n$(CYAN)Verifying all minis/ projects compile...$(NC)"
	@for d in $(ROOT_DIR)/minis/*/; do \
		echo "Building $$(basename $$d)..."; \
		go build ./$$d/... 2>/dev/null || echo "  (no main package)"; \
	done
	@echo "\n$(CYAN)Verifying all geth/ projects compile...$(NC)"
	@for d in $(ROOT_DIR)/geth/*/; do \
		echo "Building $$(basename $$d)..."; \
		go build ./$$d/... 2>/dev/null || echo "  (no main package)"; \
	done
	@echo "\n$(GREEN)✓ All projects verified successfully$(NC)"

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

list-minis:
	@echo "$(CYAN)Go Fundamentals (minis/):$(NC)"
	@ls -1d $(ROOT_DIR)/minis/*/ 2>/dev/null | xargs -n1 basename | nl -w2 -s'. '

list-geth:
	@echo "$(CYAN)Ethereum Development (geth/):$(NC)"
	@ls -1d $(ROOT_DIR)/geth/*/ 2>/dev/null | xargs -n1 basename | nl -w2 -s'. '

# ============================================================================
# RUNNING PROJECTS
# ============================================================================

run:
	@if [ "$(IN_PROJECT)" = "yes" ] && [ -z "$(P)" ]; then \
		if [ -d "$(CURRENT_DIR)/cmd/app" ]; then \
			echo "$(CYAN)Running current project: $(PROJECT_NAME)$(NC)"; \
			go run $(CURRENT_DIR)/cmd/app/main.go; \
		else \
			echo "$(YELLOW)No cmd/app directory in current project$(NC)"; \
			echo "Try: make test"; \
		fi; \
	elif [ -z "$(P)" ]; then \
		echo "$(YELLOW)Usage: make run P=<project>$(NC)"; \
		echo ""; \
		echo "From root: Specify project path"; \
		echo "  make run P=minis/01-hello-strings"; \
		echo "  make run P=geth/01-stack"; \
		echo ""; \
		echo "From project directory: Just run 'make run'"; \
		exit 1; \
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
			echo "$(CYAN)Running $$PROJECT_PATH$(NC)"; \
			go run $$PROJECT_PATH/cmd/app/main.go; \
		else \
			echo "$(YELLOW)No cmd/app directory found$(NC)"; \
			echo "Try: make test P=$(P)"; \
		fi; \
	fi

dev:
	@if [ "$(IN_PROJECT)" = "yes" ]; then \
		if [ -d "$(CURRENT_DIR)/cmd/dev" ]; then \
			echo "$(CYAN)Running dev harness: $(PROJECT_NAME)$(NC)"; \
			go run $(CURRENT_DIR)/cmd/dev/main.go; \
		else \
			echo "$(YELLOW)No cmd/dev directory in current project$(NC)"; \
		fi; \
	elif [ ! -z "$(P)" ]; then \
		PROJECT_PATH="$(ROOT_DIR)/$(P)"; \
		if [ ! -d "$$PROJECT_PATH" ] && [ -d "$(ROOT_DIR)/minis/$(P)" ]; then \
			PROJECT_PATH="$(ROOT_DIR)/minis/$(P)"; \
		fi; \
		if [ -d "$$PROJECT_PATH/cmd/dev" ]; then \
			echo "$(CYAN)Running dev harness: $(P)$(NC)"; \
			go run $$PROJECT_PATH/cmd/dev/main.go; \
		else \
			echo "$(YELLOW)No cmd/dev directory found$(NC)"; \
		fi; \
	else \
		echo "$(YELLOW)Usage:$(NC)"; \
		echo "  From project directory: make dev"; \
		echo "  From root: make dev P=minis/01-hello-strings"; \
	fi

# ============================================================================
# TESTING & BENCHMARKS
# ============================================================================

test:
	@if [ "$(IN_PROJECT)" = "yes" ] && [ -z "$(P)" ]; then \
		echo "$(CYAN)Testing current project: $(PROJECT_NAME)$(NC)"; \
		go test -v ./...; \
	elif [ -z "$(P)" ]; then \
		echo "$(CYAN)Running all tests...$(NC)"; \
		cd $(ROOT_DIR) && go test -v ./...; \
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
		go test -v $$PROJECT_PATH/...; \
	fi

bench:
	@if [ "$(IN_PROJECT)" = "yes" ] && [ -z "$(P)" ]; then \
		echo "$(CYAN)Benchmarking current project: $(PROJECT_NAME)$(NC)"; \
		go test -bench=. -benchmem ./...; \
	elif [ -z "$(P)" ]; then \
		echo "$(CYAN)Running all benchmarks...$(NC)"; \
		cd $(ROOT_DIR) && go test -bench=. -benchmem ./...; \
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
		go test -bench=. -benchmem $$PROJECT_PATH/...; \
	fi

# ============================================================================
# CODE QUALITY
# ============================================================================

lint:
	@echo "$(CYAN)Running golangci-lint...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
		echo "$(GREEN)✓ Lint complete$(NC)"; \
	else \
		echo "$(YELLOW)golangci-lint not installed. Install with:$(NC)"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

check:
	@echo "$(CYAN)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✓ go vet complete$(NC)"
	@echo "$(CYAN)Running staticcheck...$(NC)"
	@if command -v staticcheck >/dev/null 2>&1; then \
		staticcheck ./...; \
		echo "$(GREEN)✓ staticcheck complete$(NC)"; \
	else \
		echo "$(YELLOW)staticcheck not installed. Install with:$(NC)"; \
		echo "  go install honnef.co/go/tools/cmd/staticcheck@latest"; \
	fi

# ============================================================================
# EXERCISE MANAGEMENT - Reset to TODO state
# ============================================================================

# Reset exercise files to their original TODO state using git
# Usage:
#   make reset T=all              - Reset all exercises
#   make reset T=geth             - Reset all geth exercises
#   make reset T=minis            - Reset all minis exercises
#   make reset T=geth/01-stack    - Reset specific project
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
	go clean -testcache
	rm -f coverage.out
	@find . -name "*.test" -delete 2>/dev/null || true
	@find . -name "*.out" -delete 2>/dev/null || true
	@echo "$(GREEN)✓ Build cache cleaned$(NC)"

# ============================================================================
# HELP
# ============================================================================

help:
	@echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"
	@if [ "$(IN_PROJECT)" = "yes" ]; then \
		echo "$(CYAN)  📍 Current Project: $(PROJECT_NAME)$(NC)"; \
	else \
		echo "$(CYAN)  go-edu - Makefile Commands$(NC)"; \
	fi
	@echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"
	@echo ""
	@if [ "$(IN_PROJECT)" = "yes" ]; then \
		echo "$(GREEN)Project Commands (you are in a project directory):$(NC)"; \
		echo "  make run             Run this project's CLI (cmd/app)"; \
		echo "  make dev             Run debug harness (cmd/dev)"; \
		echo "  make test            Run tests for this project"; \
		echo "  make bench           Run benchmarks for this project"; \
		echo ""; \
		echo "$(GREEN)Direct Go Commands:$(NC)"; \
		echo "  go test -v ./...                    Run tests"; \
		echo "  go run ./cmd/app/main.go [args]     Run CLI"; \
		echo "  go run ./cmd/dev/main.go            Run debug harness"; \
		echo ""; \
		if [ "$(IS_GETH)" = "yes" ]; then \
			echo "$(YELLOW)💡 Geth Project Tips:$(NC)"; \
			echo "  - Set RPC URL: export INFURA_RPC_URL=https://..."; \
			echo "  - See README.md for CLI argument examples"; \
			echo ""; \
		fi; \
		if [ "$(IS_MINIS)" = "yes" ]; then \
			echo "$(YELLOW)💡 Minis Project Tips:$(NC)"; \
			echo "  - See README.md for CLI argument examples"; \
			echo "  - Use VS Code debugger (F5) for stepping through code"; \
			echo ""; \
		fi; \
		echo "$(GREEN)Exercise Management:$(NC)"; \
		echo "  make reset T=$(REL_DIR)             Reset this exercise"; \
		echo ""; \
		echo "$(GREEN)Navigation:$(NC)"; \
		echo "  cd $(ROOT_DIR)                       Go to repository root"; \
		echo "  make list                            List all projects"; \
		echo ""; \
	else \
		echo "$(GREEN)Setup & Discovery:$(NC)"; \
		echo "  make setup           Initialize dependencies and verify builds"; \
		echo "  make list            Show all available projects"; \
		echo "  make list-minis      Show only minis/ projects"; \
		echo "  make list-geth       Show only geth/ projects"; \
		echo ""; \
		echo "$(GREEN)Running Projects:$(NC)"; \
		echo "  make run P=<path>    Run specific project"; \
		echo "                       Examples:"; \
		echo "                         make run P=minis/01-hello-strings"; \
		echo "                         make run P=geth/01-stack"; \
		echo "  make dev P=<path>    Run project debug harness"; \
		echo ""; \
		echo "$(GREEN)Testing:$(NC)"; \
		echo "  make test            Run all tests"; \
		echo "  make test P=<path>   Test specific project"; \
		echo "  make bench           Run all benchmarks"; \
		echo "  make bench P=<path>  Benchmark specific project"; \
		echo ""; \
		echo "$(GREEN)Code Quality:$(NC)"; \
		echo "  make lint            Run golangci-lint on all projects"; \
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
		echo "  4. make help                          # See project-specific help"; \
		echo "  5. make run                           # Run the project"; \
		echo "  6. make test                          # Run tests"; \
		echo ""; \
	fi
	@echo "$(GREEN)Additional Help:$(NC)"
	@echo "  README.md            Project documentation"
	@echo "  .vscode/README.md    VS Code setup and debugging guide"
	@echo ""
	@echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"
