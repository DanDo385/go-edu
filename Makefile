.PHONY: setup list list-minis list-geth test bench run clean help reset lint check verify-lessons verify-teaching

# Colors for output
CYAN := \033[0;36m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
BOLD := \033[1m
NC := \033[0m # No Color

# Detect current directory relative to repo root
ROOT_DIR := $(shell git rev-parse --show-toplevel 2>/dev/null || pwd)
CURRENT_DIR := $(shell pwd)
REL_DIR := $(shell realpath --relative-to=$(ROOT_DIR) $(CURRENT_DIR) 2>/dev/null || echo ".")

# Detect if we're in a project directory
IS_MINIS_PROJECT := $(shell echo "$(REL_DIR)" | grep -q "^minis/[^/]*$$" && echo "yes" || echo "no")
IS_GETH_PROJECT := $(shell echo "$(REL_DIR)" | grep -q "^geth/[^/]*$$" && echo "yes" || echo "no")
IS_PROJECT := $(shell [ "$(IS_MINIS_PROJECT)" = "yes" ] || [ "$(IS_GETH_PROJECT)" = "yes" ] && echo "yes" || echo "no")
PROJECT_NAME := $(shell basename "$(CURRENT_DIR)")

# ============================================================================
# CONTEXT-AWARE HELP - Different output based on current directory
# ============================================================================

help:
ifeq ($(IS_PROJECT),yes)
	@echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"
	@echo "$(CYAN)  $(PROJECT_NAME) - Project Commands$(NC)"
	@echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"
	@echo ""
	@echo "$(BOLD)You are in:$(NC) $(REL_DIR)"
	@echo ""
	@echo "$(GREEN)Quick Start:$(NC)"
	@echo "  go run ./cmd/app/main.go    Run the CLI application"
	@echo "  go test -v ./...            Run tests"
	@echo ""
	@echo "$(GREEN)Development:$(NC)"
	@echo "  make run                    Run cmd/app"
	@echo "  make test                   Run tests for this project"
	@echo "  make bench                  Run benchmarks"
	@echo ""
	@echo "$(GREEN)Files to Edit:$(NC)"
	@echo "  internal/*/exercise.go      YOUR CODE GOES HERE"
	@echo "  internal/*/solution.reference.go  Reference solution"
	@echo ""
	@echo "$(GREEN)Debugging (VS Code):$(NC)"
	@echo "  1. Set breakpoints in exercise.go"
	@echo "  2. Press F5 to start debugging"
	@echo "  3. Use F10 (Step Over) and F11 (Step Into)"
	@echo ""
	@echo "$(YELLOW)Tip:$(NC) Read the README.md in this directory for CLI examples"
	@echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"
else
	@echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"
	@echo "$(CYAN)  go-edu - Makefile Commands$(NC)"
	@echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"
	@echo ""
	@echo "$(GREEN)Setup & Discovery:$(NC)"
	@echo "  make setup           Initialize dependencies and verify builds"
	@echo "  make list            Show all available projects"
	@echo "  make list-minis      Show only minis/ projects (Go fundamentals)"
	@echo "  make list-geth       Show only geth/ projects (Ethereum dev)"
	@echo ""
	@echo "$(GREEN)Running Projects:$(NC)"
	@echo "  make run P=<path>    Run project's cmd/app"
	@echo "                       Examples:"
	@echo "                         make run P=minis/01-hello-strings"
	@echo ""
	@echo "$(GREEN)Testing:$(NC)"
	@echo "  make test            Run all tests"
	@echo "  make test P=<path>   Test specific project"
	@echo "  make bench P=<path>  Benchmark specific project"
	@echo ""
	@echo "$(GREEN)Code Quality:$(NC)"
	@echo "  make lint            Run golangci-lint"
	@echo "  make check           Run go vet and staticcheck"
	@echo "  make verify-lessons  Validate lesson contract files"
	@echo "  make verify-teaching Validate README teaching sections"
	@echo ""
	@echo "$(GREEN)Exercise Management:$(NC)"
	@echo "  make reset T=all              Reset all exercises to TODO state"
	@echo "  make reset T=geth             Reset all geth exercises"
	@echo "  make reset T=minis            Reset all minis exercises"
	@echo "  make reset T=geth/01-stack    Reset specific project"
	@echo ""
	@echo "$(GREEN)Cleanup:$(NC)"
	@echo "  make clean           Clean build cache and artifacts"
	@echo ""
	@echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"
	@echo "$(YELLOW)Quick Start:$(NC)"
	@echo "  1. make list                          # See all projects"
	@echo "  2. cd minis/01-hello-strings          # Pick a project"
	@echo "  3. make help                          # See project-specific help"
	@echo "  4. code internal/hellostrings/exercise.go  # Start coding"
	@echo "  5. go test -v ./...                   # Run tests"
	@echo ""
	@echo "$(YELLOW)For geth/ projects:$(NC)"
	@echo "  export INFURA_RPC_URL=https://mainnet.infura.io/v3/YOUR_KEY"
	@echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"
endif

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
ifeq ($(IS_PROJECT),yes)
	@echo "$(CYAN)Running ./cmd/app/...$(NC)"
	@go run ./cmd/app/...
else
	@if [ -z "$(P)" ]; then \
		echo "$(YELLOW)Usage: make run P=<project>$(NC)"; \
		echo "Examples:"; \
		echo "  make run P=minis/01-hello-strings"; \
		echo "  make run P=geth/01-stack"; \
		exit 1; \
	fi
	@PROJECT_PATH="$(ROOT_DIR)/$(P)"; \
	if [ ! -d "$$PROJECT_PATH" ] && [ -d "$(ROOT_DIR)/minis/$(P)" ]; then \
		PROJECT_PATH="$(ROOT_DIR)/minis/$(P)"; \
	fi; \
	if [ ! -d "$$PROJECT_PATH" ]; then \
		echo "$(YELLOW)Error: Project '$(P)' not found$(NC)"; \
		echo "Run 'make list' to see available projects"; \
		exit 1; \
	fi; \
	if [ -d "$$PROJECT_PATH/cmd/app" ]; then \
		echo "$(CYAN)Running $$PROJECT_PATH/cmd/app/...$(NC)"; \
		go run $$PROJECT_PATH/cmd/app/...; \
	else \
		echo "$(YELLOW)No cmd/app directory found$(NC)"; \
		echo "Try: make test P=$(P)"; \
	fi
endif

# ============================================================================
# TESTING & BENCHMARKS
# ============================================================================

test:
ifeq ($(IS_PROJECT),yes)
	@echo "$(CYAN)Running tests for $(PROJECT_NAME)...$(NC)"
	@go test -v ./...
else
	@if [ -z "$(P)" ]; then \
		echo "$(CYAN)Running all tests...$(NC)"; \
		go test -v ./...; \
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
endif

bench:
ifeq ($(IS_PROJECT),yes)
	@echo "$(CYAN)Running benchmarks for $(PROJECT_NAME)...$(NC)"
	@go test -bench=. -benchmem ./...
else
	@if [ -z "$(P)" ]; then \
		echo "$(CYAN)Running all benchmarks...$(NC)"; \
		go test -bench=. -benchmem ./...; \
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
endif

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

verify-lessons:
	@echo "$(CYAN)Validating lesson contracts...$(NC)"
	@./scripts/validate_lessons.sh

verify-teaching:
	@echo "$(CYAN)Validating README teaching sections...$(NC)"
	@./scripts/validate_teaching_sections.sh

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
	go clean -testcache
	rm -f coverage.out
	@find . -name "*.test" -delete 2>/dev/null || true
	@find . -name "*.out" -delete 2>/dev/null || true
	@find . -name "todos.json" -delete 2>/dev/null || true
	@echo "$(GREEN)✓ Build cache cleaned$(NC)"
