#!/bin/bash

# Batch restructure all projects in minis/ and geth/ directories
# Usage: ./batch_restructure_all.sh [--dry-run]

set -e

DRY_RUN=false
if [ "$1" == "--dry-run" ]; then
    DRY_RUN=true
    echo "DRY RUN MODE - No changes will be made"
    echo ""
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RESTRUCTURE_SCRIPT="$SCRIPT_DIR/restructure_project_v2.sh"

if [ ! -f "$RESTRUCTURE_SCRIPT" ]; then
    echo "Error: restructure_project_v2.sh not found"
    exit 1
fi

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "=========================================="
echo "Batch Restructuring All Projects"
echo "=========================================="
echo ""

# Function to restructure a single project
restructure_project() {
    local project_path=$1
    local project_name=$(basename "$project_path")

    echo -e "${CYAN}Processing: $project_path${NC}"

    if [ "$DRY_RUN" == "true" ]; then
        echo "  [DRY RUN] Would restructure: $project_path"
        return 0
    fi

    # Run the restructure script
    if "$RESTRUCTURE_SCRIPT" "$project_path" > /tmp/restructure_${project_name}.log 2>&1; then
        echo -e "${GREEN}  ✓ Success${NC}"
        return 0
    else
        echo -e "${RED}  ✗ Failed - see /tmp/restructure_${project_name}.log${NC}"
        return 1
    fi
}

# Count projects
MINIS_COUNT=$(find minis -maxdepth 1 -mindepth 1 -type d 2>/dev/null | wc -l)
GETH_COUNT=$(find geth -maxdepth 1 -mindepth 1 -type d 2>/dev/null | wc -l)
TOTAL_COUNT=$((MINIS_COUNT + GETH_COUNT))

echo "Found $MINIS_COUNT minis/ projects"
echo "Found $GETH_COUNT geth/ projects"
echo "Total: $TOTAL_COUNT projects"
echo ""

if [ "$DRY_RUN" == "false" ]; then
    read -p "Continue with restructuring? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Aborted"
        exit 1
    fi
    echo ""
fi

# Process minis/ projects
echo "=========================================="
echo "Processing minis/ projects ($MINIS_COUNT)"
echo "=========================================="
echo ""

SUCCESS_COUNT=0
FAIL_COUNT=0

for project_dir in minis/*/; do
    if [ -d "$project_dir" ]; then
        project_path="${project_dir%/}"  # Remove trailing slash

        if restructure_project "$project_path"; then
            ((SUCCESS_COUNT++))
        else
            ((FAIL_COUNT++))
        fi
        echo ""
    fi
done

# Process geth/ projects
echo "=========================================="
echo "Processing geth/ projects ($GETH_COUNT)"
echo "=========================================="
echo ""

for project_dir in geth/*/; do
    if [ -d "$project_dir" ]; then
        project_path="${project_dir%/}"  # Remove trailing slash

        if restructure_project "$project_path"; then
            ((SUCCESS_COUNT++))
        else
            ((FAIL_COUNT++))
        fi
        echo ""
    fi
done

# Summary
echo "=========================================="
echo "Batch Restructuring Complete"
echo "=========================================="
echo ""
echo -e "${GREEN}Successful: $SUCCESS_COUNT${NC}"
echo -e "${RED}Failed: $FAIL_COUNT${NC}"
echo "Total: $TOTAL_COUNT"
echo ""

if [ $FAIL_COUNT -gt 0 ]; then
    echo -e "${YELLOW}Check log files in /tmp/restructure_*.log for details${NC}"
    echo ""
fi

echo "NEXT STEPS:"
echo ""
echo "1. For each project, manually update:"
echo "   - cmd/project-name/main.go (add CLI flags, debugging comments)"
echo "   - exercise/exercise_no_err.go (implement core logic)"
echo "   - README.md (update structure, add debugging section)"
echo "   - WORDS.md (add debugging guidance, if exists)"
echo ""
echo "2. Verify all projects compile:"
echo "   - make setup"
echo ""
echo "3. Test a few projects:"
echo "   - go run ./minis/01-hello-strings/cmd/hello-strings/main.go"
echo "   - cd minis/01-hello-strings/exercise && go test"
echo "   - cd minis/01-hello-strings/exercise && go test -tags=solution"
echo ""
echo "4. Commit changes:"
echo "   - git add -A"
echo "   - git commit -m \"Restructure all projects with debugging infrastructure\""
echo "   - git push"
echo ""
