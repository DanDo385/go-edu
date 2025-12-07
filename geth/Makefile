.PHONY: help tidy test clean

help:
	@echo "Ethereum Go Development Track"
	@echo ""
	@echo "Exercises now live entirely inside each module's exercise/ package."
	@echo "Use go test to execute your work or the reference implementation:"
	@echo "  go test ./01-stack/...            # Run your code"
	@echo "  go test -tags solution ./01-stack/...  # Run the reference solution"
	@echo ""
	@echo "Targets:"
	@echo "  tidy  - Run go mod tidy to update dependencies"
	@echo "  test  - Execute all module tests with the solution build tag"
	@echo "  clean - Remove generated artifacts (databases, *.out, etc.)"

tidy:
	@echo "Running go mod tidy..."
	@go mod tidy
	@echo "✓ Dependencies updated"

test:
	@echo "Running solution test suite (requires INFURA_RPC_URL for RPC-heavy modules)..."
	@go test -tags solution ./...
	@echo "✓ Tests completed"

clean:
	@echo "Cleaning build artifacts..."
	@find . -name "*.test" -delete
	@find . -name "*.out" -delete
	@find . -name "*.db" -delete
	@find . -name "*.db-shm" -delete
	@find . -name "*.db-wal" -delete
	@echo "✓ Clean complete"
