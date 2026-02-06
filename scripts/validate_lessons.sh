#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
EXIT_CODE=0
LESSON_COUNT=0

is_todo_reference() {
	local file="$1"
	rg -q 'TODO: Implement|This file intentionally left minimal|TODO stub|core logic functions here' "$file"
}

has_meaningful_code() {
	local file="$1"
	# Count non-empty, non-comment, non-build-tag lines.
	local count
	count=$(awk '
		/^[[:space:]]*$/ { next }
		/^[[:space:]]*\/\// { next }
		{ print }
	' "$file" | wc -l | tr -d ' ')
	[[ "$count" -ge 8 ]]
}

check_lesson() {
	local lesson_dir="$1"
	local lesson_rel="${lesson_dir#"$ROOT_DIR"/}"
	local lesson_ok=1
	local -a errors=()
	local -a exercise_files=()

	LESSON_COUNT=$((LESSON_COUNT + 1))

	if [[ ! -f "$lesson_dir/README.md" ]]; then
		errors+=("missing README.md")
		lesson_ok=0
	fi

	if [[ -d "$lesson_dir/internal" ]]; then
		while IFS= read -r path; do
			exercise_files+=("$path")
		done < <(find "$lesson_dir/internal" -type f -name "exercise.go" | sort)
	fi

	if [[ "${#exercise_files[@]}" -eq 0 ]]; then
		errors+=("missing internal/**/exercise.go")
		lesson_ok=0
	fi

	for exercise_file in "${exercise_files[@]}"; do
		local pkg_dir
		local pkg_rel
		local ref_file
		pkg_dir="$(dirname "$exercise_file")"
		pkg_rel="${pkg_dir#"$ROOT_DIR"/}"
		ref_file="$pkg_dir/solution.reference.go"

		if [[ ! -f "$ref_file" ]]; then
			errors+=("missing ${pkg_rel}/solution.reference.go")
			lesson_ok=0
		else
			if is_todo_reference "$ref_file"; then
				errors+=("${pkg_rel}/solution.reference.go appears TODO-only or placeholder")
				lesson_ok=0
			fi
			if ! rg -q '^func[[:space:]]+' "$ref_file"; then
				errors+=("${pkg_rel}/solution.reference.go has no function implementations")
				lesson_ok=0
			fi
			if ! has_meaningful_code "$ref_file"; then
				errors+=("${pkg_rel}/solution.reference.go has insufficient executable content")
				lesson_ok=0
			fi
		fi

		if [[ ! -f "$pkg_dir/exercise_test.go" ]]; then
			errors+=("missing ${pkg_rel}/exercise_test.go")
			lesson_ok=0
		fi
	done

	if [[ "$lesson_ok" -eq 1 ]]; then
		printf '[OK] %s\n' "$lesson_rel"
		return 0
	fi

	printf '[FAIL] %s\n' "$lesson_rel"
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
		check_lesson "$lesson_dir" || true
	done
done

if [[ "$LESSON_COUNT" -eq 0 ]]; then
	echo "No lessons found under minis/ or geth/."
	exit 1
fi

if [[ "$EXIT_CODE" -eq 0 ]]; then
	echo "All lesson contracts are valid ($LESSON_COUNT lessons checked)."
fi

exit "$EXIT_CODE"
