.PHONY: setup list run test lint clean help

# Colors for output
CYAN := \033[0;36m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

# Detect directories
ROOT_DIR := $(shell git rev-parse --show-toplevel 2>/dev/null || pwd)
CURRENT_DIR := $(shell pwd)
REL_DIR := $(shell realpath --relative-to=$(ROOT_DIR) $(CURRENT_DIR) 2>/dev/null || echo ".")

# Check if we are inside a project directory (e.g., minis/01-...)
IS_PROJECT := $(shell if [ -f "$(CURRENT_DIR)/go.mod" ] || [ -d "$(CURRENT_DIR)/cmd" ]; then echo "yes"; else echo "no"; fi)
PROJECT_NAME := $(shell basename "$(CURRENT_DIR)")

# ============================================================================
# MAIN COMMANDS
# ============================================================================

setup:
	@echo "$(CYAN)Initializing dependencies...$(NC)"
	go mod tidy
	@echo "\n$(CYAN)Verifying projects compile...$(NC)"
	@go build ./minis/... ./geth/... || echo "$(YELLOW)Some projects failed to build. Check 'make lint'.$(NC)"
	@echo "$(GREEN)✓ Setup complete$(NC)"

list:
	@echo "$(CYAN)Available Projects:$(NC)"
	@echo "$(CYAN)--- Minis ---$(NC)"
	@ls -1d $(ROOT_DIR)/minis/*/ 2>/dev/null | xargs -n1 basename
	@echo ""
	@echo "$(CYAN)--- Geth ---$(NC)"
	@ls -1d $(ROOT_DIR)/geth/*/ 2>/dev/null | xargs -n1 basename

run:
	@if [ "$(IS_PROJECT)" = "yes" ] && [ "$(REL_DIR)" != "." ]; then \
		echo "$(CYAN)Running current project: $(PROJECT_NAME)$(NC)"; \
		if [ -d "$(CURRENT_DIR)/cmd/app" ]; then \
			go run "$(CURRENT_DIR)/cmd/app/main.go" $(ARGS); \
		elif [ -d "$(CURRENT_DIR)/cmd" ]; then \
			go run "$(CURRENT_DIR)/cmd/..." $(ARGS); \
		else \
			echo "$(YELLOW)No cmd/ directory found in current project.$(NC)"; \
		fi; \
	else \
		if [ -z "$(P)" ]; then \
			echo "$(YELLOW)Usage: make run P=<project> [ARGS=...]$(NC)"; \
			exit 1; \
		fi; \
		PROJECT_PATH="$(ROOT_DIR)/$(P)"; \
		if [ ! -d "$$PROJECT_PATH" ] && [ -d "$(ROOT_DIR)/minis/$(P)" ]; then \
			PROJECT_PATH="$(ROOT_DIR)/minis/$(P)"; \
		fi; \
		if [ ! -d "$$PROJECT_PATH" ] && [ -d "$(ROOT_DIR)/geth/$(P)" ]; then \
			PROJECT_PATH="$(ROOT_DIR)/geth/$(P)"; \
		fi; \
		if [ ! -d "$$PROJECT_PATH" ]; then \
			echo "$(YELLOW)Error: Project '$(P)' not found$(NC)"; \
			exit 1; \
		fi; \
		echo "$(CYAN)Running $$PROJECT_PATH...$(NC)"; \
		if [ -d "$$PROJECT_PATH/cmd/app" ]; then \
			go run "$$PROJECT_PATH/cmd/app/main.go" $(ARGS); \
		else \
			go run "$$PROJECT_PATH/cmd/..." $(ARGS); \
		fi; \
	fi

test:
	@if [ "$(IS_PROJECT)" = "yes" ] && [ "$(REL_DIR)" != "." ]; then \
		echo "$(CYAN)Testing current project: $(PROJECT_NAME)$(NC)"; \
		go test -v ./...; \
	else \
		if [ -z "$(P)" ]; then \
			echo "$(CYAN)Running all tests...$(NC)"; \
			go test -v ./...; \
		else \
			PROJECT_PATH="$(ROOT_DIR)/$(P)"; \
			if [ ! -d "$$PROJECT_PATH" ] && [ -d "$(ROOT_DIR)/minis/$(P)" ]; then \
				PROJECT_PATH="$(ROOT_DIR)/minis/$(P)"; \
			fi; \
			if [ ! -d "$$PROJECT_PATH" ] && [ -d "$(ROOT_DIR)/geth/$(P)" ]; then \
				PROJECT_PATH="$(ROOT_DIR)/geth/$(P)"; \
			fi; \
			if [ ! -d "$$PROJECT_PATH" ]; then \
				echo "$(YELLOW)Error: Project '$(P)' not found$(NC)"; \
				exit 1; \
			fi; \
			echo "$(CYAN)Testing $$PROJECT_PATH...$(NC)"; \
			go test -v $$PROJECT_PATH/...; \
		fi; \
	fi

lint:
	@echo "$(CYAN)Running linters...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "$(YELLOW)golangci-lint not installed (skipping)$(NC)"; \
	fi
	@echo "$(CYAN)Running go vet...$(NC)"
	@go vet ./...

clean:
	@echo "$(CYAN)Cleaning...$(NC)"
	@go clean -testcache
	@rm -f coverage.out

help:
	@echo "$(CYAN)Go Educational Repo$(NC)"
	@echo ""
	@if [ "$(IS_PROJECT)" = "yes" ] && [ "$(REL_DIR)" != "." ]; then \
		echo "$(GREEN)Current Context: $(PROJECT_NAME)$(NC)"; \
		echo "  make run         Run this project (cmd/app)"; \
		echo "  make test        Test this project"; \
		echo "  make lint        Lint this project"; \
		echo ""; \
		echo "To return to root context, cd to root."; \
	else \
		echo "$(GREEN)Root Context$(NC)"; \
		echo "  make setup       Initialize dependencies"; \
		echo "  make list        List all projects"; \
		echo "  make run P=...   Run a specific project"; \
		echo "  make test P=...  Test a specific project (or all)"; \
		echo "  make lint        Run linters"; \
		echo "  make clean       Clean artifacts"; \
	fi
	@echo ""
