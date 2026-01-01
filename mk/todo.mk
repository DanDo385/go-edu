.PHONY: todo

# This file defines the repository-wide "make todo" behavior.
# It is intended to be included from:
# - repo root:        include mk/todo.mk
# - geth/ or minis/:  include ../mk/todo.mk
# - a project folder: include ../../mk/todo.mk

# Resolve repo root. Prefer git, fall back to path of this file.
REPO_ROOT := $(shell git rev-parse --show-toplevel 2>/dev/null)
ifeq ($(strip $(REPO_ROOT)),)
REPO_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/..)
endif

# `make todo all` or `make todo <path>`
todo:
	@target="$(filter-out $@,$(MAKECMDGOALS))"; \
	if [ -z "$$target" ]; then target="all"; fi; \
	here="$$(pwd)"; \
	root="$(REPO_ROOT)"; \
	if [ "$$target" = "all" ]; then \
		if [ "$$here" = "$$root" ]; then \
			paths="$$root/geth $$root/minis"; \
		elif [ "$$here" = "$$root/geth" ]; then \
			paths="$$root/geth"; \
		elif [ "$$here" = "$$root/minis" ]; then \
			paths="$$root/minis"; \
		else \
			paths="$$here"; \
		fi; \
	else \
		if [ -e "$$target" ]; then \
			paths="$$target"; \
		elif [ -e "$$root/$$target" ]; then \
			paths="$$root/$$target"; \
		else \
			echo "Error: path not found: $$target"; \
			echo "Tried: $$target"; \
			echo "   and: $$root/$$target"; \
			exit 1; \
		fi; \
	fi; \
	go run "$$root/cmd/todo" $$paths

# Swallow extra make goals used as "arguments", e.g. `make todo all` or `make todo geth/01-stack`.
# IMPORTANT: only do this when `todo` is one of the requested goals, so it doesn't clobber
# unrelated targets (like `make test` or `make check`) in parent Makefiles.
ifeq ($(filter todo,$(MAKECMDGOALS)),todo)
.PHONY: $(filter-out todo,$(MAKECMDGOALS))
$(filter-out todo,$(MAKECMDGOALS)):
	@:
endif

