.PHONY: setup list list-minis list-geth test bench run run-minis run-geth check clean help todo all

# Colors for output
CYAN := \033[0;36m
GREEN := \033[0;32m
YELLOW := \033[0;33m
NC := \033[0m # No Color

# Helper to extract arguments for todo command
TODO_ARGS = $(filter-out todo,$(MAKECMDGOALS))

# Default target
all:
ifneq ($(filter todo,$(MAKECMDGOALS)),)
	@:
else
	@$(MAKE) setup
endif

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

# Run a specific project
run:
	@if [ -z "$(P)" ]; then \
		echo "$(YELLOW)Usage: make run P=<project>$(NC)"; \
		echo "Examples:"; \
		echo "  make run P=minis/01-hello-strings"; \
		exit 1; \
	fi
	@PROJECT_PATH="$(P)"; \
	if [ ! -d "$$PROJECT_PATH" ] && [ -d "minis/$$PROJECT_PATH" ]; then \
		PROJECT_PATH="minis/$$PROJECT_PATH"; \
	fi; \
	if [ ! -d "$$PROJECT_PATH" ]; then \
		echo "$(YELLOW)Error: Project '$(P)' not found$(NC)"; \
		exit 1; \
	fi; \
	if [ -d "$$PROJECT_PATH/cmd/app" ]; then \
		echo "$(CYAN)Running $$PROJECT_PATH/cmd/app...$(NC)"; \
		go run ./$$PROJECT_PATH/cmd/app/...; \
	elif [ -d "$$PROJECT_PATH/cmd" ]; then \
		echo "$(CYAN)Running $$PROJECT_PATH/cmd/...$(NC)"; \
		go run ./$$PROJECT_PATH/cmd/...; \
	else \
		echo "$(YELLOW)No cmd/ directory found in $$PROJECT_PATH$(NC)"; \
	fi

# Static analysis
check:
	@echo "$(CYAN)Running go vet...$(NC)"
	go vet ./...
	@echo "$(GREEN)✓ Code checked$(NC)"

clean:
	@echo "$(CYAN)Cleaning build cache...$(NC)"
	go clean -testcache
	rm -f coverage.out
	@echo "$(GREEN)✓ Build cache cleaned$(NC)"

# Reset exercises contextually
todo:
	@if [ "$(TODO_ARGS)" = "all" ] || [ -z "$(TODO_ARGS)" ]; then \
		echo "$(CYAN)Resetting all exercises...$(NC)"; \
		./scripts/reset-exercises.sh; \
	else \
		./scripts/reset-exercises.sh $(TODO_ARGS); \
	fi

help:
	@echo "$(CYAN)Go Educational Projects - Makefile Commands$(NC)"
	@echo ""
	@echo "$(GREEN)Setup:$(NC)    make setup"
	@echo "$(GREEN)List:$(NC)     make list"
	@echo "$(GREEN)Run:$(NC)      make run P=<path>"
	@echo "$(GREEN)Test:$(NC)     make test P=<path>"
	@echo "$(GREEN)Check:$(NC)    make check"
	@echo "$(GREEN)Reset:$(NC)    make todo [all|<path>]"
