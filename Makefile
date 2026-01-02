.PHONY: help list list-minis list-geth run dev test bench reset clean lint check

# Colors for output
CYAN := \033[0;36m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m

# Paths
ROOT_DIR := $(shell git rev-parse --show-toplevel 2>/dev/null || pwd)
CURRENT_DIR := $(shell pwd)
REL_DIR := $(shell realpath --relative-to=$(ROOT_DIR) $(CURRENT_DIR) 2>/dev/null || echo ".")

# ----------------------------------------------------------------------------
# Discovery
# ----------------------------------------------------------------------------

list:
	@echo "$(CYAN)minis/$(NC)"; \
	ls -1d $(ROOT_DIR)/minis/*/ 2>/dev/null | xargs -n1 basename | nl -w2 -s'. ' || true; \
	echo ""; \
	echo "$(CYAN)geth/$(NC)"; \
	ls -1d $(ROOT_DIR)/geth/*/ 2>/dev/null | xargs -n1 basename | nl -w2 -s'. ' || true

list-minis:
	@ls -1d $(ROOT_DIR)/minis/*/ 2>/dev/null | xargs -n1 basename | nl -w2 -s'. ' || true

list-geth:
	@ls -1d $(ROOT_DIR)/geth/*/ 2>/dev/null | xargs -n1 basename | nl -w2 -s'. ' || true

# ----------------------------------------------------------------------------
# Run / test helpers (context-aware)
# ----------------------------------------------------------------------------

define RESOLVE_PROJECT
PROJECT=""; \
if [ -n "$(P)" ]; then \
	if [ -d "$(ROOT_DIR)/$(P)" ]; then PROJECT="$(ROOT_DIR)/$(P)"; \
	elif [ -d "$(ROOT_DIR)/minis/$(P)" ]; then PROJECT="$(ROOT_DIR)/minis/$(P)"; \
	elif [ -d "$(ROOT_DIR)/geth/$(P)" ]; then PROJECT="$(ROOT_DIR)/geth/$(P)"; \
	fi; \
else \
	case "$(REL_DIR)" in \
		minis/*|geth/*) PROJECT_REL="$$(printf "%s" "$(REL_DIR)" | cut -d/ -f1-2)"; PROJECT="$(ROOT_DIR)/$$PROJECT_REL" ;; \
	esac; \
fi; \
if [ -z "$$PROJECT" ]; then \
	echo "$(YELLOW)No project selected.$(NC)"; \
	echo "Run from inside a project dir (minis/<name> or geth/<name>)"; \
	echo "or pass P=<project> (e.g. P=minis/01-hello-strings)."; \
	exit 2; \
fi;
endef

run:
	@$(RESOLVE_PROJECT) \
	if [ -d "$$PROJECT/cmd/app" ]; then \
		echo "$(CYAN)Running: $$PROJECT/cmd/app$(NC)"; \
		cd "$$PROJECT" && go run ./cmd/app; \
	else \
		echo "$(YELLOW)No cmd/app found in $$PROJECT$(NC)"; \
		exit 2; \
	fi

dev:
	@$(RESOLVE_PROJECT) \
	if [ -d "$$PROJECT/cmd/dev" ]; then \
		echo "$(CYAN)Running: $$PROJECT/cmd/dev$(NC)"; \
		cd "$$PROJECT" && go run ./cmd/dev; \
	else \
		echo "$(YELLOW)No cmd/dev found in $$PROJECT$(NC)"; \
		exit 2; \
	fi

test:
	@$(RESOLVE_PROJECT) \
	echo "$(CYAN)Testing: $$PROJECT$(NC)"; \
	cd "$$PROJECT" && go test ./...

bench:
	@$(RESOLVE_PROJECT) \
	echo "$(CYAN)Benchmarking: $$PROJECT$(NC)"; \
	cd "$$PROJECT" && go test -bench=. -benchmem ./...

# ----------------------------------------------------------------------------
# Exercise management
# ----------------------------------------------------------------------------

reset:
	@if [ -z "$(T)" ]; then \
		echo "$(YELLOW)Usage: make reset T=<target>$(NC)"; \
		echo "Targets:"; \
		echo "  make reset T=all"; \
		echo "  make reset T=minis"; \
		echo "  make reset T=geth"; \
		echo "  make reset T=minis/01-hello-strings"; \
		exit 2; \
	fi
	@if [ "$(T)" = "all" ]; then \
		cd $(ROOT_DIR) && git checkout HEAD -- 'minis/*/internal/*/exercise.go' 'geth/*/internal/*/exercise.go' 2>/dev/null || true; \
	elif [ "$(T)" = "minis" ]; then \
		cd $(ROOT_DIR) && git checkout HEAD -- 'minis/*/internal/*/exercise.go' 2>/dev/null || true; \
	elif [ "$(T)" = "geth" ]; then \
		cd $(ROOT_DIR) && git checkout HEAD -- 'geth/*/internal/*/exercise.go' 2>/dev/null || true; \
	elif [ -d "$(ROOT_DIR)/$(T)" ]; then \
		cd $(ROOT_DIR) && git checkout HEAD -- '$(T)/internal/*/exercise.go' 2>/dev/null || true; \
	else \
		echo "$(RED)Unknown reset target: $(T)$(NC)"; \
		exit 2; \
	fi
	@echo "$(GREEN)✓ reset complete$(NC)"

clean:
	@echo "$(CYAN)Cleaning Go build/test caches...$(NC)"
	@go clean -testcache
	@rm -f coverage.out
	@echo "$(GREEN)✓ clean complete$(NC)"

# ----------------------------------------------------------------------------
# Optional quality targets (additive)
# ----------------------------------------------------------------------------

lint:
	@echo "$(CYAN)Running golangci-lint...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "$(YELLOW)golangci-lint not installed$(NC)"; \
		echo "Install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

check:
	@echo "$(CYAN)Running go vet...$(NC)"
	@go vet ./...

# ----------------------------------------------------------------------------
# Help (context-aware)
# ----------------------------------------------------------------------------

help:
	@echo "$(CYAN)go-edu Makefile$(NC)"; \
	echo ""; \
	case "$(REL_DIR)" in \
		minis/*|geth/*) \
			echo "$(GREEN)In project: $(REL_DIR)$(NC)"; \
			echo "  make run            Run cmd/app for this project"; \
			echo "  make dev            Run cmd/dev (auto examples) for this project"; \
			echo "  make test           Run tests for this project"; \
			echo "  make bench          Run benchmarks for this project"; \
			echo ""; \
			echo "  Tip: see project README.md for CLI flags/examples."; \
			;; \
		*) \
			echo "$(GREEN)At repo root$(NC)"; \
			echo "  make list                   List minis/ and geth/ projects"; \
			echo "  make run P=<project>         Run cmd/app (e.g. P=minis/01-hello-strings)"; \
			echo "  make dev P=<project>         Run cmd/dev (auto examples)"; \
			echo "  make test P=<project>        Test a project"; \
			echo "  make bench P=<project>       Bench a project"; \
			echo ""; \
			echo "  make reset T=<target>        Reset exercise.go files via git checkout"; \
			echo "  make clean                   Clean caches"; \
			echo ""; \
			echo "Optional:"; \
			echo "  make lint                    golangci-lint (if installed)"; \
			echo "  make check                   go vet"; \
			;; \
	esac
