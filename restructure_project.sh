#!/bin/bash

# Restructure a single Go educational project
# Usage: ./restructure_project.sh <project_path>
# Example: ./restructure_project.sh minis/01-hello-strings

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

# Step 1: Move cmd/*/main.go → main.go
echo "[1/8] Flattening cmd/ structure..."
if [ -d "cmd" ]; then
    # Find main.go in cmd subdirectories
    MAIN_GO=$(find cmd -name "main.go" -type f | head -1)
    if [ -n "$MAIN_GO" ]; then
        mv "$MAIN_GO" main.go
        echo "  ✓ Moved $MAIN_GO → main.go"
    else
        echo "  ⚠ No main.go found in cmd/"
    fi
else
    echo "  ⚠ No cmd/ directory found"
fi

# Step 2: Move exercise/* → .
echo "[2/8] Flattening exercise/ structure..."
if [ -d "exercise" ]; then
    if [ -f "exercise/exercise.go" ]; then
        mv exercise/exercise.go .
        echo "  ✓ Moved exercise/exercise.go → exercise.go"
    fi

    if [ -f "exercise/solution.go" ]; then
        mv exercise/solution.go .
        echo "  ✓ Moved exercise/solution.go → solution.go"
    fi

    if [ -f "exercise/exercise_test.go" ]; then
        mv exercise/exercise_test.go .
        echo "  ✓ Moved exercise/exercise_test.go → exercise_test.go"
    fi
else
    echo "  ⚠ No exercise/ directory found"
fi

# Step 3: Rename words.md → WORDS.md
echo "[3/8] Renaming words.md → WORDS.md..."
if [ -f "words.md" ]; then
    mv words.md WORDS.md
    echo "  ✓ Renamed words.md → WORDS.md"
elif [ -f "WORDS.md" ]; then
    echo "  ✓ WORDS.md already exists"
else
    echo "  ⚠ No words.md file found"
fi

# Step 4: Create exercise_no_err.go placeholder
echo "[4/8] Creating exercise_no_err.go..."
if [ ! -f "exercise_no_err.go" ]; then
    # Extract package name from exercise.go if it exists
    PACKAGE_NAME="exercise"
    if [ -f "exercise.go" ]; then
        PACKAGE_NAME=$(grep -m 1 '^package ' exercise.go | awk '{print $2}')
    fi

    cat > exercise_no_err.go << EOF
//go:build !solution
// +build !solution

package ${PACKAGE_NAME}

/*
Core Logic Without Error Handling
==================================

This file demonstrates the core algorithm without error handling complexity.
It's designed to help you understand the fundamental logic flow.

DEBUGGING:
- Set breakpoints at function entry points
- Step through the pure logic without error handling distractions
- Focus on understanding the algorithm

See /RUN_DEBUG.md for comprehensive debugging guide.
*/

// TODO: Implement core logic functions here
// These should mirror the logic in exercise.go but without error handling
// Focus on the algorithm, assume all operations succeed
EOF
    echo "  ✓ Created exercise_no_err.go"
else
    echo "  ✓ exercise_no_err.go already exists"
fi

# Step 5: Create .vscode/launch.json
echo "[5/8] Creating .vscode/launch.json..."
mkdir -p .vscode

cat > .vscode/launch.json << EOF
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug: Run main.go (Default Args)",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "\${workspaceFolder}/${TRACK}/${PROJECT_NAME}",
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
            "program": "\${workspaceFolder}/${TRACK}/${PROJECT_NAME}",
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
            "program": "\${workspaceFolder}/${TRACK}/${PROJECT_NAME}",
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
            "program": "\${workspaceFolder}/${TRACK}/${PROJECT_NAME}",
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
                "^$\{selectedText}$"
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
            "program": "\${workspaceFolder}/${TRACK}/${PROJECT_NAME}",
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

# Step 6: Remove empty directories
echo "[6/8] Removing empty directories..."
if [ -d "cmd" ]; then
    find cmd -type d -empty -delete 2>/dev/null || true
    if [ -d "cmd" ] && [ -z "$(ls -A cmd)" ]; then
        rmdir cmd
        echo "  ✓ Removed empty cmd/ directory"
    fi
fi

if [ -d "exercise" ]; then
    find exercise -type d -empty -delete 2>/dev/null || true
    if [ -d "exercise" ] && [ -z "$(ls -A exercise)" ]; then
        rmdir exercise
        echo "  ✓ Removed empty exercise/ directory"
    fi
fi

# Step 7: Update import paths in files
echo "[7/8] Updating import paths..."
# This is a placeholder - actual import path updates will be done per-project
echo "  ⚠ Import path updates need manual verification"

# Step 8: Verify structure
echo "[8/8] Verifying new structure..."
echo ""
echo "Final structure:"
ls -1 | grep -E '\.(go|md)$|^\.vscode$' || true
echo ""

echo "=========================================="
echo "✓ Restructuring complete: $PROJECT_PATH"
echo "=========================================="
echo ""
echo "Next steps:"
echo "1. Verify import paths in moved files"
echo "2. Update main.go with CLI parameters and debugging comments"
echo "3. Update solution.go with detailed explanations"
echo "4. Update exercise.go with TODO format"
echo "5. Update README.md with debugging section"
echo "6. Update WORDS.md with debugging guidance (if exists)"
echo "7. Test: go run main.go"
echo "8. Test: go test"
echo "9. Test: go test -tags=solution"
echo ""
