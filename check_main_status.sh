#!/bin/bash

# Script to check status of all main.go files
# Shows which need enhancement with CLI args and debugging comments

echo "Checking all projects for CLI args and debugging comments..."
echo ""

check_project() {
    local project=$1
    local main_file="${project}/cmd/*/main.go"

    if [ ! -f $main_file 2>/dev/null ]; then
        return
    fi

    local has_flag=$(grep -c "\"flag\"" $main_file 2>/dev/null || echo "0")
    local has_breakpoint=$(grep -c "// BREAKPOINT:" $main_file 2>/dev/null || echo "0")
    local has_debug=$(grep -c "// DEBUG:" $main_file 2>/dev/null || echo "0")
    local lines=$(wc -l < $main_file 2>/dev/null || echo "0")

    local status="❌ NEEDS WORK"
    if [ "$has_flag" -gt "0" ] && [ "$has_breakpoint" -gt "3" ] && [ "$has_debug" -gt "3" ]; then
        status="✅ COMPLETE"
    elif [ "$has_flag" -gt "0" ]; then
        status="🟡 PARTIAL (has flags)"
    fi

    printf "%-40s %12s  Lines: %4d  Flags: %2d  BP: %2d  Debug: %2d\n" \
        "$(basename $project)" "$status" "$lines" "$has_flag" "$has_breakpoint" "$has_debug"
}

echo "MINIS Projects:"
echo "----------------------------------------"
for dir in minis/*/; do
    check_project "$dir"
done

echo ""
echo "GETH Projects:"
echo "----------------------------------------"
for dir in geth/*/; do
    check_project "$dir"
done
