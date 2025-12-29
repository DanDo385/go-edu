#!/bin/bash

# Restructure a single Go educational project (Version 2 - Refined Spec)
# Usage: ./restructure_project_v2.sh <project_path>
# Example: ./restructure_project_v2.sh minis/01-hello-strings

set -e

if [ -z "$1" ]; then
    echo "Usage: $0 <project_path>"
    echo "Example: $0 minis/01-hello-strings"
    exit 1
fi

PROJECT_PATH="$1"
PROJECT_NAME=$(basename "$PROJECT_PATH")
TRACK=$(dirname "$PROJECT_PATH")  # minis or geth

echo "=========================================="
echo "Restructuring: $PROJECT_PATH"
echo "=========================================="

# Check if project exists
if [ ! -d "$PROJECT_PATH" ]; then
    echo "Error: Project directory not found: $PROJECT_PATH"
    exit 1
fi

cd "$PROJECT_PATH"

# Step 1: Verify cmd/ directory exists
echo "[1/7] Verifying cmd/ structure..."
if [ ! -d "cmd" ]; then
    echo "  ⚠ No cmd/ directory found - this project may need manual setup"
else
    echo "  ✓ cmd/ directory exists"
fi

# Step 2: Verify exercise/ directory exists
echo "[2/7] Verifying exercise/ structure..."
if [ ! -d "exercise" ]; then
    echo "  ⚠ No exercise/ directory found - this project may need manual setup"
else
    echo "  ✓ exercise/ directory exists"
fi

# Step 3: Create exercise_no_err.go in exercise/ directory
echo "[3/7] Creating exercise_no_err.go..."
if [ ! -f "exercise/exercise_no_err.go" ] && [ -d "exercise" ]; then
    # Extract package name from exercise.go if it exists
    PACKAGE_NAME="exercise"
    if [ -f "exercise/exercise.go" ]; then
        PACKAGE_NAME=$(grep -m 1 '^package ' exercise/exercise.go | awk '{print $2}')
    fi

    cat > exercise/exercise_no_err.go << EOF
//go:build !solution
// +build !solution

package ${PACKAGE_NAME}

/*
Core Logic Without Error Handling
==================================

This file demonstrates the core algorithm without error handling complexity.
It's designed to help you understand the fundamental logic flow.

KEY CONCEPTS:
- Focus on the algorithm, not error handling
- Simplified implementations for learning
- Use this to understand core logic before adding error handling

DEBUGGING:
- Set breakpoints at function entry points
- Step through the pure logic without error handling distractions
- Focus on understanding the algorithm
- Watch Variables panel to see data transformations

See /RUN_DEBUG.md for comprehensive debugging guide.
*/

// TODO: Implement core logic functions here
// These should mirror the logic in exercise.go but without error handling
// Focus on the algorithm, assume all operations succeed

// Example:
// func CoreFunctionName(input string) string {
//     // BREAKPOINT: Set breakpoint here to step through algorithm
//     // DEBUG: Watch Variables panel to see transformations
//
//     // Core logic without error handling
//     result := processInput(input)
//
//     return result
// }
EOF
    echo "  ✓ Created exercise/exercise_no_err.go"
else
    echo "  ✓ exercise/exercise_no_err.go already exists or no exercise/ directory"
fi

# Step 4: Rename words.md → WORDS.md
echo "[4/7] Renaming words.md → WORDS.md..."
if [ -f "words.md" ]; then
    mv words.md WORDS.md
    echo "  ✓ Renamed words.md → WORDS.md"
elif [ -f "WORDS.md" ]; then
    echo "  ✓ WORDS.md already exists"
else
    echo "  ⚠ No words.md file found"
fi

# Step 5: Create .vscode/launch.json
echo "[5/7] Creating .vscode/launch.json..."
mkdir -p .vscode

# Find the cmd subdirectory name
CMD_SUBDIR=$(find cmd -maxdepth 1 -mindepth 1 -type d 2>/dev/null | head -1 | xargs basename 2>/dev/null || echo "${PROJECT_NAME}")

cat > .vscode/launch.json << EOF
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug: Run main.go (Default Args)",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "\${workspaceFolder}/${TRACK}/${PROJECT_NAME}/cmd/${CMD_SUBDIR}",
            "args": [],
            "showLog": true,
            "trace": "log",
            "buildFlags": "-gcflags='-N -l'",
            "comment": "Runs main.go with default CLI arguments. Set breakpoints and step through code."
        },
        {
            "name": "Debug: Run main.go (Custom Args)",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "\${workspaceFolder}/${TRACK}/${PROJECT_NAME}/cmd/${CMD_SUBDIR}",
            "args": [
                "-help"
            ],
            "showLog": true,
            "trace": "log",
            "buildFlags": "-gcflags='-N -l'",
            "comment": "Runs main.go with custom CLI arguments. Edit 'args' array to change arguments."
        },
        {
            "name": "Debug: Test exercise.go",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "\${workspaceFolder}/${TRACK}/${PROJECT_NAME}/exercise",
            "args": [
                "-test.v"
            ],
            "showLog": true,
            "trace": "log",
            "buildFlags": "-gcflags='-N -l'",
            "comment": "Runs tests for exercise.go. Set breakpoints in test functions or exercise.go."
        },
        {
            "name": "Debug: Test solution.go",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "\${workspaceFolder}/${TRACK}/${PROJECT_NAME}/exercise",
            "args": [
                "-tags=solution",
                "-test.v"
            ],
            "showLog": true,
            "trace": "log",
            "buildFlags": "-gcflags='-N -l'",
            "comment": "Runs tests for solution.go. Set breakpoints to see reference implementation."
        },
        {
            "name": "Debug: Current File",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "\${fileDirname}",
            "args": [],
            "showLog": true,
            "buildFlags": "-gcflags='-N -l'",
            "comment": "Debugs the currently open file. Useful for quick debugging."
        },
        {
            "name": "Debug: Current Test",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "\${fileDirname}",
            "args": [
                "-test.v",
                "-test.run",
                "^\${selectedText}\$"
            ],
            "showLog": true,
            "buildFlags": "-gcflags='-N -l'",
            "comment": "Debugs the currently selected test function. Select test name first."
        },
        {
            "name": "Debug: Concurrent Code (with Race Detector)",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "\${workspaceFolder}/${TRACK}/${PROJECT_NAME}/exercise",
            "args": [
                "-race",
                "-test.v"
            ],
            "showLog": true,
            "trace": "log",
            "buildFlags": "-gcflags='-N -l'",
            "comment": "Runs tests with race detector. Detects data races in concurrent code."
        }
    ]
}
EOF
echo "  ✓ Created .vscode/launch.json"

# Step 6: Verify structure
echo "[6/7] Verifying structure..."
echo ""
echo "Expected structure:"
echo "  ${PROJECT_PATH}/"
echo "  ├── cmd/"
echo "  │   └── ${CMD_SUBDIR}/"
echo "  │       └── main.go"
echo "  ├── exercise/"
echo "  │   ├── exercise.go"
echo "  │   ├── exercise_no_err.go  (NEW)"
echo "  │   ├── solution.go"
echo "  │   └── exercise_test.go"
echo "  ├── .vscode/"
echo "  │   └── launch.json  (NEW)"
echo "  ├── README.md"
echo "  └── WORDS.md  (if exists)"
echo ""

# Step 7: Provide next steps
echo "[7/7] Next steps:"
echo ""
echo "=========================================="
echo "✓ Basic restructuring complete: $PROJECT_PATH"
echo "=========================================="
echo ""
echo "MANUAL TASKS REQUIRED:"
echo "1. Update cmd/${CMD_SUBDIR}/main.go:"
echo "   - Add CLI parameters using flag package"
echo "   - Import exercise package: import \"github.com/example/go-10x-minis/${TRACK}/${PROJECT_NAME}/exercise\""
echo "   - Call exercise functions: exercise.FunctionName()"
echo "   - Add extensive BREAKPOINT and DEBUG comments"
echo ""
echo "2. Update exercise/solution.go:"
echo "   - Add BREAKPOINT comments at key algorithm steps"
echo "   - Add DEBUG comments explaining what to watch"
echo "   - Ensure functions are exported (capitalized)"
echo ""
echo "3. Update exercise/exercise.go:"
echo "   - Ensure functions are exported (capitalized)"
echo "   - Add TODO format with debugging hints"
echo ""
echo "4. Implement exercise/exercise_no_err.go:"
echo "   - Add core logic implementations"
echo "   - No error handling"
echo "   - Debugging comments"
echo ""
echo "5. Update README.md:"
echo "   - Update structure section"
echo "   - Update running instructions (go run ./cmd/${CMD_SUBDIR}/main.go)"
echo "   - Add debugging section"
echo ""
echo "6. Update WORDS.md (if exists):"
echo "   - Add debugging guidance"
echo ""
echo "7. Verify:"
echo "   - go run ./cmd/${CMD_SUBDIR}/main.go"
echo "   - cd exercise && go test"
echo "   - cd exercise && go test -tags=solution"
echo ""
