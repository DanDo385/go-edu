#!/usr/bin/env bash

set -euo pipefail

usage() {
	echo "Usage: $0 [--force-template] <readme-path> [<readme-path> ...]"
	echo
	echo "Default mode: report missing teaching sections without editing files."
	echo "--force-template: append placeholder teaching sections (manual rewrite still required)."
}

if [[ "$#" -eq 0 ]]; then
	usage
	exit 1
fi

force_template=0
if [[ "${1:-}" == "--force-template" ]]; then
	force_template=1
	shift
fi

if [[ "$#" -eq 0 ]]; then
	usage
	exit 1
fi

has_heading() {
	local file="$1"
	shift
	local pattern
	for pattern in "$@"; do
		if rg -q "$pattern" "$file"; then
			return 0
		fi
	done
	return 1
}

append_placeholder_sections() {
	cat <<'EOT'
## Core Concepts

- TODO: Replace with lesson-specific concepts (no generic boilerplate).

## CS Connection

- TODO: Explain concrete memory/runtime links for this lesson.

## End-State Understanding

- TODO: State what learners can explain and implement after this lesson.

## Step-by-Step Build Path

### Step 1: Problem This Step Solves
- TODO

### Step 2: Why This Approach
- TODO

### Step 3: Memory / Pointer Impact
- TODO (required when `*` or `&` appears in lesson code)

### Step 4: What Changed
- TODO

## Pointer and Indirection Checklist (`*` and `&`)

- TODO: Explain exact meaning of `*` / `&` in this lesson.
- TODO: Describe memory-before and memory-after state transitions.
- TODO: Link to `docs/MEMORY_POINTERS_PRIMER.md`.

EOT
}

backfill_file() {
	local file="$1"
	local missing=0

	if [[ ! -f "$file" ]]; then
		echo "[SKIP] $file (not found)"
		return 0
	fi

	if ! has_heading "$file" '^##\s+Core Concepts(\b|:)' '^##\s+Core Concepts Involved(\b|:)'; then
		missing=1
		echo "[MISS] $file -> Core Concepts"
	fi
	if ! has_heading "$file" '^##\s+CS Connection(\b|:)' '^##\s+Connection to (Fundamental )?CS Ideas(\b|:)'; then
		missing=1
		echo "[MISS] $file -> CS Connection"
	fi
	if ! has_heading "$file" '^##\s+End-State(\b|:)' '^##\s+End-State Understanding(\b|:)'; then
		missing=1
		echo "[MISS] $file -> End-State Understanding"
	fi

	if [[ "$missing" -eq 0 ]]; then
		echo "[OK]   $file (teaching headings present)"
		return 0
	fi

	if [[ "$force_template" -eq 0 ]]; then
		echo "[INFO] $file not modified. Re-run with --force-template to append placeholders."
		return 0
	fi

	{
		echo
		echo "<!-- PLACEHOLDER: replace with lesson-specific teaching content -->"
		append_placeholder_sections
	} >> "$file"

	echo "[EDIT] $file (placeholder template appended)"
}

for readme in "$@"; do
	backfill_file "$readme"
done
