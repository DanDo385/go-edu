.PHONY: help list test run check lint clean

include mk/todo.mk

help:
	@echo "Targets:"
	@echo "  make todo all            Reset exercises to starter TODO state (contextual)"
	@echo "  make todo <path>         Reset exercises under <path> (root-relative works anywhere)"
	@echo "  make list                List projects"
	@echo "  make test [P=<path>]     Run tests (optionally for one project)"
	@echo "  make run P=<path>        Run a project's cmd/app (if present)"
	@echo "  make check               Quick sanity checks (test + vet)"
	@echo "  make lint                Static checks (go vet)"
	@echo "  make clean               Clean Go test cache"

list:
	@cd "$(REPO_ROOT)" && \
	echo "minis/:" && ls -1d minis/*/ 2>/dev/null | sed 's|minis/||' | sed 's|/||' && \
	echo "" && \
	echo "geth/:" && ls -1d geth/*/ 2>/dev/null | sed 's|geth/||' | sed 's|/||'

test:
	@cd "$(REPO_ROOT)" && \
	if [ -z "$(P)" ]; then \
		go test ./...; \
	else \
		p="$(P)"; \
		if [ ! -d "$$p" ] && [ -d "$(REPO_ROOT)/$$p" ]; then p="$(REPO_ROOT)/$$p"; fi; \
		if [ ! -d "$$p" ]; then echo "Error: project not found: $(P)"; exit 1; fi; \
		go test "./$${p#$(REPO_ROOT)/}/..."; \
	fi

run:
	@cd "$(REPO_ROOT)" && \
	if [ -z "$(P)" ]; then echo "Usage: make run P=<path>"; exit 1; fi; \
	p="$(P)"; \
	if [ ! -d "$$p" ] && [ -d "$(REPO_ROOT)/$$p" ]; then p="$(REPO_ROOT)/$$p"; fi; \
	if [ ! -d "$$p" ]; then echo "Error: project not found: $(P)"; exit 1; fi; \
	if [ -d "$$p/cmd/app" ]; then \
		go run "./$${p#$(REPO_ROOT)/}/cmd/app"; \
	else \
		echo "Error: $$p/cmd/app not found (try go test ./... instead)"; \
		exit 1; \
	fi

check:
	@cd "$(REPO_ROOT)" && \
	echo "Running reference test suite for minis/ ..." && \
	go test -tags=reference ./minis/... && \
	if [ -n "$$INFURA_RPC_URL" ]; then \
		echo "" && echo "Running reference test suite for geth/ (INFURA_RPC_URL set) ..." && \
		go test -tags=reference ./geth/...; \
	else \
		echo "" && echo "Skipping geth/ reference tests (set INFURA_RPC_URL to enable)"; \
	fi && \
	echo "" && echo "Running go vet ..." && \
	go vet ./...

lint:
	@cd "$(REPO_ROOT)" && go vet ./...

clean:
	@cd "$(REPO_ROOT)" && go clean -testcache
