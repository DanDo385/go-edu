#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
EXIT_CODE=0
READMES_CHECKED=0

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

uses_pointer_tokens() {
	local lesson_dir="$1"
	# Detect explicit pointer/address operators in lesson Go files.
	rg -q '[&*]' "$lesson_dir/internal"/*.go "$lesson_dir/internal"/*/*.go 2>/dev/null
}

check_readme() {
	local lesson_dir="$1"
	local file="$2"
	local rel="${file#"$ROOT_DIR"/}"
	local -a errors=()
	local ok=1

	READMES_CHECKED=$((READMES_CHECKED + 1))

	if ! has_heading "$file" '^##\s+Core Concepts(\b|:)' '^##\s+Core Concepts Involved(\b|:)'; then
		errors+=("missing 'Core Concepts' section")
		ok=0
	fi

	if ! has_heading "$file" '^##\s+CS Connection(\b|:)' '^##\s+Connection to (Fundamental )?CS Ideas(\b|:)'; then
		errors+=("missing 'CS Connection' section")
		ok=0
	fi

	if ! has_heading "$file" '^##\s+End-State(\b|:)' '^##\s+End-State Understanding(\b|:)'; then
		errors+=("missing 'End-State' section")
		ok=0
	fi

	if ! has_heading "$file" '^##\s+Step-by-Step' '^##\s+Implementation Steps'; then
		errors+=("missing step-by-step teaching section")
		ok=0
	fi

	if ! rg -q '^###\s+Step\s+1:' "$file" || ! rg -q '^###\s+Step\s+4:' "$file"; then
		errors+=("step section must include at least Step 1 and Step 4 headings")
		ok=0
	fi

	# Generic boilerplate injected by automation should not persist.
	if rg -q 'Data representation and invariant design for|Explain how this lesson''s data flow maps to concrete runtime state changes\.|Compare the student path \(`exercise\.go`\) against `solution\.reference\.go`' "$file"; then
		errors+=("contains generic boilerplate teaching text")
		ok=0
	fi

	if uses_pointer_tokens "$lesson_dir"; then
		if ! has_heading "$file" '^##\s+Pointer' '^##\s+Deep Dive: `\*` and `&`'; then
			errors+=("lesson code uses '*' or '&' but README lacks pointer/indirection section")
			ok=0
		fi
			if ! rg -q 'MEMORY_POINTERS_PRIMER\.md|\* and &|`\*`|`&`' "$file"; then
				errors+=("pointer lesson is missing explicit star/ampersand explanation or primer link")
				ok=0
			fi
		fi

	if [[ "$ok" -eq 1 ]]; then
		printf '[OK] %s\n' "$rel"
		return 0
	fi

	printf '[FAIL] %s\n' "$rel"
	local err
	for err in "${errors[@]}"; do
		printf '  - %s\n' "$err"
	done
	EXIT_CODE=1
	return 1
}

for group in minis geth; do
	group_dir="$ROOT_DIR/$group"
	[[ -d "$group_dir" ]] || continue
	for lesson_dir in "$group_dir"/*; do
		[[ -d "$lesson_dir" ]] || continue
		readme="$lesson_dir/README.md"
		if [[ ! -f "$readme" ]]; then
			printf '[FAIL] %s\n' "${lesson_dir#"$ROOT_DIR"/}"
			printf '  - missing README.md\n'
			EXIT_CODE=1
			continue
		fi
		check_readme "$lesson_dir" "$readme" || true
	done
done

if [[ "$READMES_CHECKED" -eq 0 ]]; then
	echo "No lesson README files found under minis/ or geth/."
	exit 1
fi

if [[ "$EXIT_CODE" -eq 0 ]]; then
	echo "All README teaching sections passed semantic checks ($READMES_CHECKED files checked)."
else
	echo "Teaching quality gaps found. Fix failing README files before merge."
fi

exit "$EXIT_CODE"
