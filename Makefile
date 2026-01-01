# Makefile for Go Exercises
#
# Usage:
#   make reset-all     # Reset all exercise.go files to TODO state
#   make reset-minis   # Reset only minis exercises
#   make reset-geth    # Reset only geth exercises
#   make dry-run       # Preview what would be reset (without changing files)

.PHONY: reset-all reset-minis reset-geth dry-run help

# Default target
help:
	@echo "Exercise Management Commands:"
	@echo "  make reset-all     Reset all exercise.go files to TODO state"
	@echo "  make reset-minis   Reset only minis exercises"
	@echo "  make reset-geth    Reset only geth exercises"
	@echo "  make dry-run       Preview what would be reset"
	@echo ""
	@echo "The reset commands regenerate exercise.go files from solution.reference.go,"
	@echo "leaving only function signatures with TODO placeholders."

# Reset all exercises
reset-all:
	@echo "Resetting all exercises..."
	@go run scripts/reset-exercises.go -target=all

# Reset minis exercises only
reset-minis:
	@echo "Resetting minis exercises..."
	@go run scripts/reset-exercises.go -target=minis

# Reset geth exercises only
reset-geth:
	@echo "Resetting geth exercises..."
	@go run scripts/reset-exercises.go -target=geth

# Dry run - show what would be reset without making changes
dry-run:
	@echo "Dry run - showing what would be reset..."
	@go run scripts/reset-exercises.go -target=all -dry-run
