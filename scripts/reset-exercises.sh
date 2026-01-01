#!/bin/bash
# Script to reset exercise.go files back to their starter TODO state.
# (Wrapper around `git restore` for tracked exercise.go files.)
#
# NOTE: Prefer `make todo ...` unless you're scripting.
#
# Usage:
#   ./scripts/reset-exercises.sh                  # Reset all exercises (minis + geth)
#   ./scripts/reset-exercises.sh minis            # Reset only minis
#   ./scripts/reset-exercises.sh geth             # Reset only geth
#   ./scripts/reset-exercises.sh geth/01-stack    # Reset a specific project/path

set -euo pipefail

ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"

if [ $# -eq 0 ]; then
	( cd "$ROOT" && make todo all )
	exit 0
fi

arg="$1"
( cd "$ROOT" && make todo "$arg" )
